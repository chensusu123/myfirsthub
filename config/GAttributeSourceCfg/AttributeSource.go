package GAttributeSourceCfg

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

// AttributeSourceConfigRow from attribute【属性】.xlsx attribute_source
type AttributeSourceConfigRow struct {
	Id      int32 `json:"id"`      // 属性ID
	Jump_id int32 `json:"jump_id"` // 跳转来源
}

// AttributeSourceConfig from attribute【属性】.xlsx attribute_source
type AttributeSourceConfig struct {
	ConfigRows map[int32]*AttributeSourceConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *AttributeSourceConfig {
	ret := &AttributeSourceConfig{ConfigRows: map[int32]*AttributeSourceConfigRow{}}
	return ret
}

// GetAttributeSourceConfig get one config by configId
func (c *AttributeSourceConfig) GetAttributeSourceConfig(configId int32) *AttributeSourceConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *AttributeSourceConfig) Get(configId int32) *AttributeSourceConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllAttributeSourceConfig get all config slice
func (c *AttributeSourceConfig) GetAllAttributeSourceConfig() (res []*AttributeSourceConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *AttributeSourceConfig) GetAll() (res []*AttributeSourceConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *AttributeSourceConfig

// GetAttributeSourceConfig pkg func. get one config by configId
func GetAttributeSourceConfig(configId int32) *AttributeSourceConfigRow {
	return gConfigData.GetAttributeSourceConfig(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *AttributeSourceConfigRow {
	return gConfigData.Get(configId)
}

// GetAllAttributeSourceConfig pkg func. get all config slice
func GetAllAttributeSourceConfig() []*AttributeSourceConfigRow {
	return gConfigData.GetAllAttributeSourceConfig()
}

// GetAll pkg func. get all config slice
func GetAll() []*AttributeSourceConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*AttributeSourceConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "AttributeSourceConfigRow from attribute【属性】.xlsx attribute_source"
}

// GetRawValue get raw data
func GetRawValue() *AttributeSourceConfig {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "attribute_source"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("attribute_source.json",
		"attribute【属性】.xlsx", "attribute_source",
		&gAttributeSourceParser{}, &gAttributeSourceLoader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*AttributeSourceConfig)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *AttributeSourceConfig) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*AttributeSourceConfig))(c)
		return true
	})
}

// RegisterAttributeSourceInitCallBack reg config update func (old func)
var RegisterAttributeSourceInitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*AttributeSourceConfig)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*AttributeSourceConfig)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *AttributeSourceConfig) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*AttributeSourceConfig) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*AttributeSourceConfig) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gAttributeSourceLoader struct {
}

// NewContainer new data container pointer
func (*gAttributeSourceLoader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gAttributeSourceLoader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*AttributeSourceConfig)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gAttributeSourceLoader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*AttributeSourceConfig)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gAttributeSourceLoader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*AttributeSourceConfigRow)
	if !ok {
		err = errors.New("invalid type. not *AttributeSourceConfigRow")
		logger.ErrorWF("invalid type. not *AttributeSourceConfigRow", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_source"))
		return
	}
	config, ok := container.(*AttributeSourceConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeSourceConfig")
		logger.ErrorWF("invalid type. not *AttributeSourceConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_source"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gAttributeSourceLoader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*AttributeSourceConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeSourceConfig")
		logger.ErrorWF("invalid type. not *AttributeSourceConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_source"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gAttributeSourceLoader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*AttributeSourceConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeSourceConfig")
		logger.ErrorWF("invalid type. not *AttributeSourceConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_source"))
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
type gAttributeSourceParser struct {
}

// New new config row data
func (*gAttributeSourceParser) New() interface{} {
	return &AttributeSourceConfigRow{}
}

// Fields get config fields names
func (*gAttributeSourceParser) Fields() []string {
	return gAttributeSourceFields
}

// Parse parse raw data to row data
func (*gAttributeSourceParser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*AttributeSourceConfigRow)
	if !ok {
		err = errors.New("invalid type. not *AttributeSourceConfigRow")
		logger.ErrorWF("invalid type. not *AttributeSourceConfigRow", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_source"))
		return
	}
	// compare length
	if len(data) != len(gAttributeSourceFields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*AttributeSourceConfigRow",
			zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_source"), zap.Int("need_count", len(gAttributeSourceFields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 属性ID
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 属性ID to int32 failed")
			logger.ErrorWF("parse field id 属性ID to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_source"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 jump_id : 跳转来源
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field jump_id 跳转来源 to int32 failed")
			logger.ErrorWF("parse field jump_id 跳转来源 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_source"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Jump_id = int32(tmp)
	}
	return
}

var gAttributeSourceFields = []string{
	"id",
	"jump_id",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gAttributeSourceParser{}
	loader := &gAttributeSourceLoader{}
	var data [][]string
	data, err = load("attribute【属性】.xlsx", "attribute_source", gAttributeSourceFields)
	if err != nil {
		logger.ErrorWF("load attribute【属性】.xlsx attribute_source data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load attribute【属性】.xlsx attribute_source data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gAttributeSourceFields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse attribute【属性】.xlsx attribute_source failed.", zap.Int("row", k), zap.Strings("need", gAttributeSourceFields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse attribute【属性】.xlsx attribute_source row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add attribute【属性】.xlsx attribute_source row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check attribute【属性】.xlsx attribute_source data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load attribute【属性】.xlsx attribute_source data success.")
	return
}
