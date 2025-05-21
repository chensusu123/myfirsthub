package GJumpUseCfg

import (
	"errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

// JumpUseConfigRow from jump_config【跳转】.xlsx jump_use
type JumpUseConfigRow struct {
	Jump_key_id    int32    `json:"jump_key_id"`    // 序号
	Jump_use_id    int32    `json:"jump_use_id"`    // 跳转id
	Jump_type      int32    `json:"jump_type"`      // 跳转指令
	Jump_param     []string `json:"jump_param"`     // 跳转参数
	Jump_condition int32    `json:"jump_condition"` // 跳转出现条件
	Jump_show_type int32    `json:"jump_show_type"` // 跳转显示类型
}

// JumpUseConfig from jump_config【跳转】.xlsx jump_use
type JumpUseConfig struct {
	ConfigRows map[int32]*JumpUseConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *JumpUseConfig {
	ret := &JumpUseConfig{ConfigRows: map[int32]*JumpUseConfigRow{}}
	return ret
}

// GetJumpUseConfig get one config by configId
func (c *JumpUseConfig) GetJumpUseConfig(configId int32) *JumpUseConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *JumpUseConfig) Get(configId int32) *JumpUseConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllJumpUseConfig get all config slice
func (c *JumpUseConfig) GetAllJumpUseConfig() (res []*JumpUseConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *JumpUseConfig) GetAll() (res []*JumpUseConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *JumpUseConfig

// GetJumpUseConfig pkg func. get one config by configId
func GetJumpUseConfig(configId int32) *JumpUseConfigRow {
	return gConfigData.GetJumpUseConfig(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *JumpUseConfigRow {
	return gConfigData.Get(configId)
}

// GetAllJumpUseConfig pkg func. get all config slice
func GetAllJumpUseConfig() []*JumpUseConfigRow {
	return gConfigData.GetAllJumpUseConfig()
}

// GetAll pkg func. get all config slice
func GetAll() []*JumpUseConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*JumpUseConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "JumpUseConfigRow from jump_config【跳转】.xlsx jump_use"
}

// GetRawValue get raw data
func GetRawValue() *JumpUseConfig {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "jump_use"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("jump_use.json",
		"jump_config【跳转】.xlsx", "jump_use",
		&gJumpUseParser{}, &gJumpUseLoader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*JumpUseConfig)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *JumpUseConfig) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*JumpUseConfig))(c)
		return true
	})
}

// RegisterJumpUseInitCallBack reg config update func (old func)
var RegisterJumpUseInitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*JumpUseConfig)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*JumpUseConfig)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *JumpUseConfig) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*JumpUseConfig) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*JumpUseConfig) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gJumpUseLoader struct {
}

// NewContainer new data container pointer
func (*gJumpUseLoader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gJumpUseLoader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*JumpUseConfig)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gJumpUseLoader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*JumpUseConfig)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gJumpUseLoader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*JumpUseConfigRow)
	if !ok {
		err = errors.New("invalid type. not *JumpUseConfigRow")
		logger.ErrorWF("invalid type. not *JumpUseConfigRow", zap.String("xlsx", "jump_config【跳转】.xlsx"),
			zap.String("sheet", "jump_use"))
		return
	}
	config, ok := container.(*JumpUseConfig)
	if !ok {
		err = errors.New("invalid type. not *JumpUseConfig")
		logger.ErrorWF("invalid type. not *JumpUseConfig", zap.String("xlsx", "jump_config【跳转】.xlsx"),
			zap.String("sheet", "jump_use"))
		return
	}
	config.ConfigRows[row.Jump_key_id] = row
	return
}

// GetValue get real map value for json parse
func (*gJumpUseLoader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*JumpUseConfig)
	if !ok {
		err = errors.New("invalid type. not *JumpUseConfig")
		logger.ErrorWF("invalid type. not *JumpUseConfig", zap.String("xlsx", "jump_config【跳转】.xlsx"),
			zap.String("sheet", "jump_use"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gJumpUseLoader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*JumpUseConfig)
	if !ok {
		err = errors.New("invalid type. not *JumpUseConfig")
		logger.ErrorWF("invalid type. not *JumpUseConfig", zap.String("xlsx", "jump_config【跳转】.xlsx"),
			zap.String("sheet", "jump_use"))
		return
	}
	for _, row := range config.ConfigRows {
		err = rf(row)
		if err != nil {
			return err
		}
	}
	return
}

// implete ConfigParser interface
type gJumpUseParser struct {
}

// New new config row data
func (*gJumpUseParser) New() interface{} {
	return &JumpUseConfigRow{}
}

// Fields get config fields names
func (*gJumpUseParser) Fields() []string {
	return gJumpUseFields
}

// Parse parse raw data to row data
func (*gJumpUseParser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*JumpUseConfigRow)
	if !ok {
		err = errors.New("invalid type. not *JumpUseConfigRow")
		logger.ErrorWF("invalid type. not *JumpUseConfigRow", zap.String("xlsx", "jump_config【跳转】.xlsx"),
			zap.String("sheet", "jump_use"))
		return
	}
	// compare length
	if len(data) != len(gJumpUseFields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*JumpUseConfigRow",
			zap.String("xlsx", "jump_config【跳转】.xlsx"),
			zap.String("sheet", "jump_use"), zap.Int("need_count", len(gJumpUseFields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 jump_key_id : 序号
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field jump_key_id 序号 to int32 failed")
			logger.ErrorWF("parse field jump_key_id 序号 to int32 failed.",
				zap.String("xlsx", "jump_config【跳转】.xlsx"), zap.String("sheet", "jump_use"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Jump_key_id = int32(tmp)
	}

	// parse column 1 jump_use_id : 跳转id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field jump_use_id 跳转id to int32 failed")
			logger.ErrorWF("parse field jump_use_id 跳转id to int32 failed.",
				zap.String("xlsx", "jump_config【跳转】.xlsx"), zap.String("sheet", "jump_use"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Jump_use_id = int32(tmp)
	}

	// parse column 2 jump_type : 跳转指令
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field jump_type 跳转指令 to int32 failed")
			logger.ErrorWF("parse field jump_type 跳转指令 to int32 failed.",
				zap.String("xlsx", "jump_config【跳转】.xlsx"), zap.String("sheet", "jump_use"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Jump_type = int32(tmp)
	}

	// parse column 3 jump_param : 跳转参数
	if data[3] != "" {
		config.Jump_param = strings.Split(data[3], ",")
	}

	// parse column 4 jump_condition : 跳转出现条件
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field jump_condition 跳转出现条件 to int32 failed")
			logger.ErrorWF("parse field jump_condition 跳转出现条件 to int32 failed.",
				zap.String("xlsx", "jump_config【跳转】.xlsx"), zap.String("sheet", "jump_use"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Jump_condition = int32(tmp)
	}

	// parse column 5 jump_show_type : 跳转显示类型
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field jump_show_type 跳转显示类型 to int32 failed")
			logger.ErrorWF("parse field jump_show_type 跳转显示类型 to int32 failed.",
				zap.String("xlsx", "jump_config【跳转】.xlsx"), zap.String("sheet", "jump_use"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Jump_show_type = int32(tmp)
	}
	return
}

var gJumpUseFields = []string{
	"jump_key_id",
	"jump_use_id",
	"jump_type",
	"jump_param",
	"jump_condition",
	"jump_show_type",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gJumpUseParser{}
	loader := &gJumpUseLoader{}
	var data [][]string
	data, err = load("jump_config【跳转】.xlsx", "jump_use", gJumpUseFields)
	if err != nil {
		logger.ErrorWF("load jump_config【跳转】.xlsx jump_use data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load jump_config【跳转】.xlsx jump_use data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gJumpUseFields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse jump_config【跳转】.xlsx jump_use failed.", zap.Int("row", k), zap.Strings("need", gJumpUseFields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse jump_config【跳转】.xlsx jump_use row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add jump_config【跳转】.xlsx jump_use row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check jump_config【跳转】.xlsx jump_use data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load jump_config【跳转】.xlsx jump_use data success.")
	return
}
