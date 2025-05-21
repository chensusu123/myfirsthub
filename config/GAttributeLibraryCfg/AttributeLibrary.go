package GAttributeLibraryCfg

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

// AttributeLibraryConfigRow from attribute【属性】.xlsx attribute_library
type AttributeLibraryConfigRow struct {
	Id           int32   `json:"id"`           // 自定义的id，其他业务使用时，自己去映射（@马健
	Attribute_id []int32 `json:"attribute_id"` // 业务内属性id（业务内曾经使用过的属性id，这里的属性id只能增加，不能减少@马健
}

// AttributeLibraryConfig from attribute【属性】.xlsx attribute_library
type AttributeLibraryConfig struct {
	ConfigRows map[int32]*AttributeLibraryConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *AttributeLibraryConfig {
	ret := &AttributeLibraryConfig{ConfigRows: map[int32]*AttributeLibraryConfigRow{}}
	return ret
}

// GetAttributeLibraryConfig get one config by configId
func (c *AttributeLibraryConfig) GetAttributeLibraryConfig(configId int32) *AttributeLibraryConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *AttributeLibraryConfig) Get(configId int32) *AttributeLibraryConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllAttributeLibraryConfig get all config slice
func (c *AttributeLibraryConfig) GetAllAttributeLibraryConfig() (res []*AttributeLibraryConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *AttributeLibraryConfig) GetAll() (res []*AttributeLibraryConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *AttributeLibraryConfig

// GetAttributeLibraryConfig pkg func. get one config by configId
func GetAttributeLibraryConfig(configId int32) *AttributeLibraryConfigRow {
	return gConfigData.GetAttributeLibraryConfig(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *AttributeLibraryConfigRow {
	return gConfigData.Get(configId)
}

// GetAllAttributeLibraryConfig pkg func. get all config slice
func GetAllAttributeLibraryConfig() []*AttributeLibraryConfigRow {
	return gConfigData.GetAllAttributeLibraryConfig()
}

// GetAll pkg func. get all config slice
func GetAll() []*AttributeLibraryConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*AttributeLibraryConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "AttributeLibraryConfigRow from attribute【属性】.xlsx attribute_library"
}

// GetRawValue get raw data
func GetRawValue() *AttributeLibraryConfig {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "attribute_library"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("attribute_library.json",
		"attribute【属性】.xlsx", "attribute_library",
		&gAttributeLibraryParser{}, &gAttributeLibraryLoader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*AttributeLibraryConfig)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *AttributeLibraryConfig) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*AttributeLibraryConfig))(c)
		return true
	})
}

// RegisterAttributeLibraryInitCallBack reg config update func (old func)
var RegisterAttributeLibraryInitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*AttributeLibraryConfig)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*AttributeLibraryConfig)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *AttributeLibraryConfig) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*AttributeLibraryConfig) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*AttributeLibraryConfig) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gAttributeLibraryLoader struct {
}

// NewContainer new data container pointer
func (*gAttributeLibraryLoader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gAttributeLibraryLoader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*AttributeLibraryConfig)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gAttributeLibraryLoader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*AttributeLibraryConfig)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gAttributeLibraryLoader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*AttributeLibraryConfigRow)
	if !ok {
		err = errors.New("invalid type. not *AttributeLibraryConfigRow")
		logger.ErrorWF("invalid type. not *AttributeLibraryConfigRow", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_library"))
		return
	}
	config, ok := container.(*AttributeLibraryConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeLibraryConfig")
		logger.ErrorWF("invalid type. not *AttributeLibraryConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_library"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gAttributeLibraryLoader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*AttributeLibraryConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeLibraryConfig")
		logger.ErrorWF("invalid type. not *AttributeLibraryConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_library"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gAttributeLibraryLoader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*AttributeLibraryConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeLibraryConfig")
		logger.ErrorWF("invalid type. not *AttributeLibraryConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_library"))
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
type gAttributeLibraryParser struct {
}

// New new config row data
func (*gAttributeLibraryParser) New() interface{} {
	return &AttributeLibraryConfigRow{}
}

// Fields get config fields names
func (*gAttributeLibraryParser) Fields() []string {
	return gAttributeLibraryFields
}

// Parse parse raw data to row data
func (*gAttributeLibraryParser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*AttributeLibraryConfigRow)
	if !ok {
		err = errors.New("invalid type. not *AttributeLibraryConfigRow")
		logger.ErrorWF("invalid type. not *AttributeLibraryConfigRow", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_library"))
		return
	}
	// compare length
	if len(data) != len(gAttributeLibraryFields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*AttributeLibraryConfigRow",
			zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_library"), zap.Int("need_count", len(gAttributeLibraryFields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 自定义的id，其他业务使用时，自己去映射（@马健
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 自定义的id，其他业务使用时，自己去映射（@马健 to int32 failed")
			logger.ErrorWF("parse field id 自定义的id，其他业务使用时，自己去映射（@马健 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_library"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 attribute_id : 业务内属性id（业务内曾经使用过的属性id，这里的属性id只能增加，不能减少@马健
	if data[1] != "" {

		vals := strings.Split(data[1], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field attribute_id 业务内属性id（业务内曾经使用过的属性id，这里的属性id只能增加，不能减少@马健 to []int32 failed")
				logger.ErrorWF("parse array field attribute_id 业务内属性id（业务内曾经使用过的属性id，这里的属性id只能增加，不能减少@马健 to []int32 failed.",
					zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_library"),
					// zap.String("field_data",data[1]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Attribute_id = append(config.Attribute_id, int32(tmp))
		}
	}
	return
}

var gAttributeLibraryFields = []string{
	"id",
	"attribute_id",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gAttributeLibraryParser{}
	loader := &gAttributeLibraryLoader{}
	var data [][]string
	data, err = load("attribute【属性】.xlsx", "attribute_library", gAttributeLibraryFields)
	if err != nil {
		logger.ErrorWF("load attribute【属性】.xlsx attribute_library data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load attribute【属性】.xlsx attribute_library data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gAttributeLibraryFields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse attribute【属性】.xlsx attribute_library failed.", zap.Int("row", k), zap.Strings("need", gAttributeLibraryFields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse attribute【属性】.xlsx attribute_library row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add attribute【属性】.xlsx attribute_library row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check attribute【属性】.xlsx attribute_library data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load attribute【属性】.xlsx attribute_library data success.")
	return
}
