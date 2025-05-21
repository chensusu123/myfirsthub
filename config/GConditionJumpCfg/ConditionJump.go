package GConditionJumpCfg

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

// ConditionJumpConfigRow from condition【条件】.xlsx condition_jump
type ConditionJumpConfigRow struct {
	Key_id          int32           `json:"key_id"`          // 主键id
	Condition_id    int32           `json:"condition_id"`    // 条件id
	Condition_type  int32           `json:"condition_type"`  // 条件类型
	Condition_value map[int32]int32 `json:"condition_value"` // 条件参数
	Condition_desc  string          `json:"condition_desc"`  // 条件描述
}

// ConditionJumpConfig from condition【条件】.xlsx condition_jump
type ConditionJumpConfig struct {
	ConfigRows map[int32]*ConditionJumpConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *ConditionJumpConfig {
	ret := &ConditionJumpConfig{ConfigRows: map[int32]*ConditionJumpConfigRow{}}
	return ret
}

// GetConditionJumpConfig get one config by configId
func (c *ConditionJumpConfig) GetConditionJumpConfig(configId int32) *ConditionJumpConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *ConditionJumpConfig) Get(configId int32) *ConditionJumpConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllConditionJumpConfig get all config slice
func (c *ConditionJumpConfig) GetAllConditionJumpConfig() (res []*ConditionJumpConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *ConditionJumpConfig) GetAll() (res []*ConditionJumpConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *ConditionJumpConfig

// GetConditionJumpConfig pkg func. get one config by configId
func GetConditionJumpConfig(configId int32) *ConditionJumpConfigRow {
	return gConfigData.GetConditionJumpConfig(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *ConditionJumpConfigRow {
	return gConfigData.Get(configId)
}

// GetAllConditionJumpConfig pkg func. get all config slice
func GetAllConditionJumpConfig() []*ConditionJumpConfigRow {
	return gConfigData.GetAllConditionJumpConfig()
}

// GetAll pkg func. get all config slice
func GetAll() []*ConditionJumpConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*ConditionJumpConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "ConditionJumpConfigRow from condition【条件】.xlsx condition_jump"
}

// GetRawValue get raw data
func GetRawValue() *ConditionJumpConfig {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "condition_jump"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("condition_jump.json",
		"condition【条件】.xlsx", "condition_jump",
		&gConditionJumpParser{}, &gConditionJumpLoader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*ConditionJumpConfig)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *ConditionJumpConfig) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*ConditionJumpConfig))(c)
		return true
	})
}

// RegisterConditionJumpInitCallBack reg config update func (old func)
var RegisterConditionJumpInitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*ConditionJumpConfig)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*ConditionJumpConfig)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *ConditionJumpConfig) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*ConditionJumpConfig) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*ConditionJumpConfig) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gConditionJumpLoader struct {
}

// NewContainer new data container pointer
func (*gConditionJumpLoader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gConditionJumpLoader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*ConditionJumpConfig)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gConditionJumpLoader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*ConditionJumpConfig)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gConditionJumpLoader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*ConditionJumpConfigRow)
	if !ok {
		err = errors.New("invalid type. not *ConditionJumpConfigRow")
		logger.ErrorWF("invalid type. not *ConditionJumpConfigRow", zap.String("xlsx", "condition【条件】.xlsx"),
			zap.String("sheet", "condition_jump"))
		return
	}
	config, ok := container.(*ConditionJumpConfig)
	if !ok {
		err = errors.New("invalid type. not *ConditionJumpConfig")
		logger.ErrorWF("invalid type. not *ConditionJumpConfig", zap.String("xlsx", "condition【条件】.xlsx"),
			zap.String("sheet", "condition_jump"))
		return
	}
	config.ConfigRows[row.Key_id] = row
	return
}

// GetValue get real map value for json parse
func (*gConditionJumpLoader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*ConditionJumpConfig)
	if !ok {
		err = errors.New("invalid type. not *ConditionJumpConfig")
		logger.ErrorWF("invalid type. not *ConditionJumpConfig", zap.String("xlsx", "condition【条件】.xlsx"),
			zap.String("sheet", "condition_jump"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gConditionJumpLoader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*ConditionJumpConfig)
	if !ok {
		err = errors.New("invalid type. not *ConditionJumpConfig")
		logger.ErrorWF("invalid type. not *ConditionJumpConfig", zap.String("xlsx", "condition【条件】.xlsx"),
			zap.String("sheet", "condition_jump"))
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
type gConditionJumpParser struct {
}

// New new config row data
func (*gConditionJumpParser) New() interface{} {
	return &ConditionJumpConfigRow{}
}

// Fields get config fields names
func (*gConditionJumpParser) Fields() []string {
	return gConditionJumpFields
}

// Parse parse raw data to row data
func (*gConditionJumpParser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*ConditionJumpConfigRow)
	if !ok {
		err = errors.New("invalid type. not *ConditionJumpConfigRow")
		logger.ErrorWF("invalid type. not *ConditionJumpConfigRow", zap.String("xlsx", "condition【条件】.xlsx"),
			zap.String("sheet", "condition_jump"))
		return
	}
	// compare length
	if len(data) != len(gConditionJumpFields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*ConditionJumpConfigRow",
			zap.String("xlsx", "condition【条件】.xlsx"),
			zap.String("sheet", "condition_jump"), zap.Int("need_count", len(gConditionJumpFields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 key_id : 主键id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field key_id 主键id to int32 failed")
			logger.ErrorWF("parse field key_id 主键id to int32 failed.",
				zap.String("xlsx", "condition【条件】.xlsx"), zap.String("sheet", "condition_jump"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Key_id = int32(tmp)
	}

	// parse column 1 condition_id : 条件id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field condition_id 条件id to int32 failed")
			logger.ErrorWF("parse field condition_id 条件id to int32 failed.",
				zap.String("xlsx", "condition【条件】.xlsx"), zap.String("sheet", "condition_jump"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Condition_id = int32(tmp)
	}

	// parse column 2 condition_type : 条件类型
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field condition_type 条件类型 to int32 failed")
			logger.ErrorWF("parse field condition_type 条件类型 to int32 failed.",
				zap.String("xlsx", "condition【条件】.xlsx"), zap.String("sheet", "condition_jump"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Condition_type = int32(tmp)
	}

	// parse column 3 condition_value : 条件参数
	if data[3] != "" {

		config.Condition_value = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field condition_value 条件参数 to key int32 failed")
				logger.ErrorWF("parse map field condition_value 条件参数 to key int32 failed.",
					zap.String("xlsx", "condition【条件】.xlsx"), zap.String("sheet", "condition_jump"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field condition_value 条件参数 to value int32 failed")
				logger.ErrorWF("parse map field condition_value 条件参数 to value int32 failed.",
					zap.String("xlsx", "condition【条件】.xlsx"), zap.String("sheet", "condition_jump"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Condition_value[key] = value
		}
	}

	// parse column 4 condition_desc : 条件描述
	if data[4] != "" {
		config.Condition_desc = data[4]
	}
	return
}

var gConditionJumpFields = []string{
	"key_id",
	"condition_id",
	"condition_type",
	"condition_value",
	"condition_desc",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gConditionJumpParser{}
	loader := &gConditionJumpLoader{}
	var data [][]string
	data, err = load("condition【条件】.xlsx", "condition_jump", gConditionJumpFields)
	if err != nil {
		logger.ErrorWF("load condition【条件】.xlsx condition_jump data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load condition【条件】.xlsx condition_jump data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gConditionJumpFields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse condition【条件】.xlsx condition_jump failed.", zap.Int("row", k), zap.Strings("need", gConditionJumpFields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse condition【条件】.xlsx condition_jump row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add condition【条件】.xlsx condition_jump row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check condition【条件】.xlsx condition_jump data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load condition【条件】.xlsx condition_jump data success.")
	return
}
