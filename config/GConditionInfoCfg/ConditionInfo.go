package GConditionInfoCfg

import (
	"errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"
)

// ConditionInfoConfigRow from condition【条件】.xlsx condition_info
type ConditionInfoConfigRow struct {
	Key_id           int32 `json:"key_id"`           // 主键id
	Condition_type   int32 `json:"condition_type"`   // 条件类型
	Condition_value1 int32 `json:"condition_value1"` // 条件参数1
	Condition_value2 int32 `json:"condition_value2"` // 条件参数2
}

// ConditionInfoConfig from condition【条件】.xlsx condition_info
type ConditionInfoConfig struct {
	ConfigRows map[int32]*ConditionInfoConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *ConditionInfoConfig {
	ret := &ConditionInfoConfig{ConfigRows: map[int32]*ConditionInfoConfigRow{}}
	return ret
}

// GetConditionInfoConfig get one config by configId
func (c *ConditionInfoConfig) GetConditionInfoConfig(configId int32) *ConditionInfoConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *ConditionInfoConfig) Get(configId int32) *ConditionInfoConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllConditionInfoConfig get all config slice
func (c *ConditionInfoConfig) GetAllConditionInfoConfig() (res []*ConditionInfoConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *ConditionInfoConfig) GetAll() (res []*ConditionInfoConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *ConditionInfoConfig

// GetConditionInfoConfig pkg func. get one config by configId
func GetConditionInfoConfig(configId int32) *ConditionInfoConfigRow {
	return gConfigData.GetConditionInfoConfig(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *ConditionInfoConfigRow {
	return gConfigData.Get(configId)
}

// GetAllConditionInfoConfig pkg func. get all config slice
func GetAllConditionInfoConfig() []*ConditionInfoConfigRow {
	return gConfigData.GetAllConditionInfoConfig()
}

// GetAll pkg func. get all config slice
func GetAll() []*ConditionInfoConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*ConditionInfoConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "ConditionInfoConfigRow from condition【条件】.xlsx condition_info"
}

// GetRawValue get raw data
func GetRawValue() *ConditionInfoConfig {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "condition_info"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("condition_info.json",
		"condition【条件】.xlsx", "condition_info",
		&gConditionInfoParser{}, &gConditionInfoLoader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*ConditionInfoConfig)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *ConditionInfoConfig) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*ConditionInfoConfig))(c)
		return true
	})
}

// RegisterConditionInfoInitCallBack reg config update func (old func)
var RegisterConditionInfoInitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*ConditionInfoConfig)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*ConditionInfoConfig)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *ConditionInfoConfig) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*ConditionInfoConfig) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*ConditionInfoConfig) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gConditionInfoLoader struct {
}

// NewContainer new data container pointer
func (*gConditionInfoLoader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gConditionInfoLoader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*ConditionInfoConfig)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gConditionInfoLoader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*ConditionInfoConfig)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gConditionInfoLoader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*ConditionInfoConfigRow)
	if !ok {
		err = errors.New("invalid type. not *ConditionInfoConfigRow")
		logger.ErrorWF("invalid type. not *ConditionInfoConfigRow", zap.String("xlsx", "condition【条件】.xlsx"),
			zap.String("sheet", "condition_info"))
		return
	}
	config, ok := container.(*ConditionInfoConfig)
	if !ok {
		err = errors.New("invalid type. not *ConditionInfoConfig")
		logger.ErrorWF("invalid type. not *ConditionInfoConfig", zap.String("xlsx", "condition【条件】.xlsx"),
			zap.String("sheet", "condition_info"))
		return
	}
	config.ConfigRows[row.Key_id] = row
	return
}

// GetValue get real map value for json parse
func (*gConditionInfoLoader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*ConditionInfoConfig)
	if !ok {
		err = errors.New("invalid type. not *ConditionInfoConfig")
		logger.ErrorWF("invalid type. not *ConditionInfoConfig", zap.String("xlsx", "condition【条件】.xlsx"),
			zap.String("sheet", "condition_info"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gConditionInfoLoader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*ConditionInfoConfig)
	if !ok {
		err = errors.New("invalid type. not *ConditionInfoConfig")
		logger.ErrorWF("invalid type. not *ConditionInfoConfig", zap.String("xlsx", "condition【条件】.xlsx"),
			zap.String("sheet", "condition_info"))
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
type gConditionInfoParser struct {
}

// New new config row data
func (*gConditionInfoParser) New() interface{} {
	return &ConditionInfoConfigRow{}
}

// Fields get config fields names
func (*gConditionInfoParser) Fields() []string {
	return gConditionInfoFields
}

// Parse parse raw data to row data
func (*gConditionInfoParser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*ConditionInfoConfigRow)
	if !ok {
		err = errors.New("invalid type. not *ConditionInfoConfigRow")
		logger.ErrorWF("invalid type. not *ConditionInfoConfigRow", zap.String("xlsx", "condition【条件】.xlsx"),
			zap.String("sheet", "condition_info"))
		return
	}
	// compare length
	if len(data) != len(gConditionInfoFields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*ConditionInfoConfigRow",
			zap.String("xlsx", "condition【条件】.xlsx"),
			zap.String("sheet", "condition_info"), zap.Int("need_count", len(gConditionInfoFields)),
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
				zap.String("xlsx", "condition【条件】.xlsx"), zap.String("sheet", "condition_info"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Key_id = int32(tmp)
	}

	// parse column 1 condition_type : 条件类型
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field condition_type 条件类型 to int32 failed")
			logger.ErrorWF("parse field condition_type 条件类型 to int32 failed.",
				zap.String("xlsx", "condition【条件】.xlsx"), zap.String("sheet", "condition_info"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Condition_type = int32(tmp)
	}

	// parse column 2 condition_value1 : 条件参数1
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field condition_value1 条件参数1 to int32 failed")
			logger.ErrorWF("parse field condition_value1 条件参数1 to int32 failed.",
				zap.String("xlsx", "condition【条件】.xlsx"), zap.String("sheet", "condition_info"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Condition_value1 = int32(tmp)
	}

	// parse column 3 condition_value2 : 条件参数2
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field condition_value2 条件参数2 to int32 failed")
			logger.ErrorWF("parse field condition_value2 条件参数2 to int32 failed.",
				zap.String("xlsx", "condition【条件】.xlsx"), zap.String("sheet", "condition_info"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Condition_value2 = int32(tmp)
	}
	return
}

var gConditionInfoFields = []string{
	"key_id",
	"condition_type",
	"condition_value1",
	"condition_value2",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gConditionInfoParser{}
	loader := &gConditionInfoLoader{}
	var data [][]string
	data, err = load("condition【条件】.xlsx", "condition_info", gConditionInfoFields)
	if err != nil {
		logger.ErrorWF("load condition【条件】.xlsx condition_info data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load condition【条件】.xlsx condition_info data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gConditionInfoFields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse condition【条件】.xlsx condition_info failed.", zap.Int("row", k), zap.Strings("need", gConditionInfoFields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse condition【条件】.xlsx condition_info row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add condition【条件】.xlsx condition_info row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check condition【条件】.xlsx condition_info data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load condition【条件】.xlsx condition_info data success.")
	return
}
