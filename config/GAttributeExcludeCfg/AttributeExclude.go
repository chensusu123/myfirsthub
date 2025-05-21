package GAttributeExcludeCfg

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

// AttributeExcludeConfigRow from attribute【属性】.xlsx attribute_exclude
type AttributeExcludeConfigRow struct {
	Id              int32           `json:"id"`              // 属性ID
	Exclude_source  []int32         `json:"exclude_source"`  // 来源互斥(此列不使用，使用MAP）
	Exclude_select  int32           `json:"exclude_select"`  // 互斥时，使用哪一种 1-使用数值大的
	Source_conflict map[int32]int32 `json:"source_conflict"` // 来源互斥(值相同时，value大的优先显示）
}

// AttributeExcludeConfig from attribute【属性】.xlsx attribute_exclude
type AttributeExcludeConfig struct {
	ConfigRows map[int32]*AttributeExcludeConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *AttributeExcludeConfig {
	ret := &AttributeExcludeConfig{ConfigRows: map[int32]*AttributeExcludeConfigRow{}}
	return ret
}

// GetAttributeExcludeConfig get one config by configId
func (c *AttributeExcludeConfig) GetAttributeExcludeConfig(configId int32) *AttributeExcludeConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *AttributeExcludeConfig) Get(configId int32) *AttributeExcludeConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllAttributeExcludeConfig get all config slice
func (c *AttributeExcludeConfig) GetAllAttributeExcludeConfig() (res []*AttributeExcludeConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *AttributeExcludeConfig) GetAll() (res []*AttributeExcludeConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *AttributeExcludeConfig

// GetAttributeExcludeConfig pkg func. get one config by configId
func GetAttributeExcludeConfig(configId int32) *AttributeExcludeConfigRow {
	return gConfigData.GetAttributeExcludeConfig(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *AttributeExcludeConfigRow {
	return gConfigData.Get(configId)
}

// GetAllAttributeExcludeConfig pkg func. get all config slice
func GetAllAttributeExcludeConfig() []*AttributeExcludeConfigRow {
	return gConfigData.GetAllAttributeExcludeConfig()
}

// GetAll pkg func. get all config slice
func GetAll() []*AttributeExcludeConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*AttributeExcludeConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "AttributeExcludeConfigRow from attribute【属性】.xlsx attribute_exclude"
}

// GetRawValue get raw data
func GetRawValue() *AttributeExcludeConfig {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "attribute_exclude"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("attribute_exclude.json",
		"attribute【属性】.xlsx", "attribute_exclude",
		&gAttributeExcludeParser{}, &gAttributeExcludeLoader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*AttributeExcludeConfig)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *AttributeExcludeConfig) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*AttributeExcludeConfig))(c)
		return true
	})
}

// RegisterAttributeExcludeInitCallBack reg config update func (old func)
var RegisterAttributeExcludeInitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*AttributeExcludeConfig)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*AttributeExcludeConfig)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *AttributeExcludeConfig) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*AttributeExcludeConfig) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*AttributeExcludeConfig) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gAttributeExcludeLoader struct {
}

// NewContainer new data container pointer
func (*gAttributeExcludeLoader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gAttributeExcludeLoader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*AttributeExcludeConfig)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gAttributeExcludeLoader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*AttributeExcludeConfig)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gAttributeExcludeLoader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*AttributeExcludeConfigRow)
	if !ok {
		err = errors.New("invalid type. not *AttributeExcludeConfigRow")
		logger.ErrorWF("invalid type. not *AttributeExcludeConfigRow", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_exclude"))
		return
	}
	config, ok := container.(*AttributeExcludeConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeExcludeConfig")
		logger.ErrorWF("invalid type. not *AttributeExcludeConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_exclude"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gAttributeExcludeLoader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*AttributeExcludeConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeExcludeConfig")
		logger.ErrorWF("invalid type. not *AttributeExcludeConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_exclude"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gAttributeExcludeLoader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*AttributeExcludeConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeExcludeConfig")
		logger.ErrorWF("invalid type. not *AttributeExcludeConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_exclude"))
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
type gAttributeExcludeParser struct {
}

// New new config row data
func (*gAttributeExcludeParser) New() interface{} {
	return &AttributeExcludeConfigRow{}
}

// Fields get config fields names
func (*gAttributeExcludeParser) Fields() []string {
	return gAttributeExcludeFields
}

// Parse parse raw data to row data
func (*gAttributeExcludeParser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*AttributeExcludeConfigRow)
	if !ok {
		err = errors.New("invalid type. not *AttributeExcludeConfigRow")
		logger.ErrorWF("invalid type. not *AttributeExcludeConfigRow", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_exclude"))
		return
	}
	// compare length
	if len(data) != len(gAttributeExcludeFields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*AttributeExcludeConfigRow",
			zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_exclude"), zap.Int("need_count", len(gAttributeExcludeFields)),
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
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_exclude"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 exclude_source : 来源互斥(此列不使用，使用MAP）
	if data[1] != "" {

		vals := strings.Split(data[1], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field exclude_source 来源互斥(此列不使用，使用MAP） to []int32 failed")
				logger.ErrorWF("parse array field exclude_source 来源互斥(此列不使用，使用MAP） to []int32 failed.",
					zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_exclude"),
					// zap.String("field_data",data[1]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Exclude_source = append(config.Exclude_source, int32(tmp))
		}
	}

	// parse column 2 exclude_select : 互斥时，使用哪一种 1-使用数值大的
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field exclude_select 互斥时，使用哪一种 1-使用数值大的 to int32 failed")
			logger.ErrorWF("parse field exclude_select 互斥时，使用哪一种 1-使用数值大的 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_exclude"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Exclude_select = int32(tmp)
	}

	// parse column 3 source_conflict : 来源互斥(值相同时，value大的优先显示）
	if data[3] != "" {

		config.Source_conflict = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field source_conflict 来源互斥(值相同时，value大的优先显示） to key int32 failed")
				logger.ErrorWF("parse map field source_conflict 来源互斥(值相同时，value大的优先显示） to key int32 failed.",
					zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_exclude"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field source_conflict 来源互斥(值相同时，value大的优先显示） to value int32 failed")
				logger.ErrorWF("parse map field source_conflict 来源互斥(值相同时，value大的优先显示） to value int32 failed.",
					zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_exclude"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Source_conflict[key] = value
		}
	}
	return
}

var gAttributeExcludeFields = []string{
	"id",
	"exclude_source",
	"exclude_select",
	"source_conflict",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gAttributeExcludeParser{}
	loader := &gAttributeExcludeLoader{}
	var data [][]string
	data, err = load("attribute【属性】.xlsx", "attribute_exclude", gAttributeExcludeFields)
	if err != nil {
		logger.ErrorWF("load attribute【属性】.xlsx attribute_exclude data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load attribute【属性】.xlsx attribute_exclude data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gAttributeExcludeFields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse attribute【属性】.xlsx attribute_exclude failed.", zap.Int("row", k), zap.Strings("need", gAttributeExcludeFields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse attribute【属性】.xlsx attribute_exclude row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add attribute【属性】.xlsx attribute_exclude row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check attribute【属性】.xlsx attribute_exclude data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load attribute【属性】.xlsx attribute_exclude data success.")
	return
}
