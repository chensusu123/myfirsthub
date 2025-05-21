package GAttributeFormulaCfg

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

// AttributeFormulaConfigRow from attribute【属性】.xlsx attribute_formula
type AttributeFormulaConfigRow struct {
	Formula_id   int32   `json:"formula_id"`   // 公式id
	Type         int32   `json:"type"`         // 公式类型
	Parameter_1  int32   `json:"parameter_1"`  // 参数1属性id
	Parameter_2  int32   `json:"parameter_2"`  // 参数2属性id
	Parameter_8  []int32 `json:"parameter_8"`  // 参数8属性id
	Parameter_3  int32   `json:"parameter_3"`  // 参数3属性id
	Parameter_4  int32   `json:"parameter_4"`  // 参数4属性id
	Parameter_5  int32   `json:"parameter_5"`  // 参数5属性id
	Parameter_9  []int32 `json:"parameter_9"`  // 参数9属性id
	Parameter_10 []int32 `json:"parameter_10"` // 参数10属性id
	Parameter_11 int32   `json:"parameter_11"` // 参数11属性id
	Parameter_12 int32   `json:"parameter_12"` // 参数12属性id
}

// AttributeFormulaConfig from attribute【属性】.xlsx attribute_formula
type AttributeFormulaConfig struct {
	ConfigRows map[int32]*AttributeFormulaConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *AttributeFormulaConfig {
	ret := &AttributeFormulaConfig{ConfigRows: map[int32]*AttributeFormulaConfigRow{}}
	return ret
}

// GetAttributeFormulaConfig get one config by configId
func (c *AttributeFormulaConfig) GetAttributeFormulaConfig(configId int32) *AttributeFormulaConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *AttributeFormulaConfig) Get(configId int32) *AttributeFormulaConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllAttributeFormulaConfig get all config slice
func (c *AttributeFormulaConfig) GetAllAttributeFormulaConfig() (res []*AttributeFormulaConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *AttributeFormulaConfig) GetAll() (res []*AttributeFormulaConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *AttributeFormulaConfig

// GetAttributeFormulaConfig pkg func. get one config by configId
func GetAttributeFormulaConfig(configId int32) *AttributeFormulaConfigRow {
	return gConfigData.GetAttributeFormulaConfig(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *AttributeFormulaConfigRow {
	return gConfigData.Get(configId)
}

// GetAllAttributeFormulaConfig pkg func. get all config slice
func GetAllAttributeFormulaConfig() []*AttributeFormulaConfigRow {
	return gConfigData.GetAllAttributeFormulaConfig()
}

// GetAll pkg func. get all config slice
func GetAll() []*AttributeFormulaConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*AttributeFormulaConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "AttributeFormulaConfigRow from attribute【属性】.xlsx attribute_formula"
}

// GetRawValue get raw data
func GetRawValue() *AttributeFormulaConfig {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "attribute_formula"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("attribute_formula.json",
		"attribute【属性】.xlsx", "attribute_formula",
		&gAttributeFormulaParser{}, &gAttributeFormulaLoader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*AttributeFormulaConfig)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *AttributeFormulaConfig) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*AttributeFormulaConfig))(c)
		return true
	})
}

// RegisterAttributeFormulaInitCallBack reg config update func (old func)
var RegisterAttributeFormulaInitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*AttributeFormulaConfig)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*AttributeFormulaConfig)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *AttributeFormulaConfig) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*AttributeFormulaConfig) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*AttributeFormulaConfig) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gAttributeFormulaLoader struct {
}

// NewContainer new data container pointer
func (*gAttributeFormulaLoader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gAttributeFormulaLoader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*AttributeFormulaConfig)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gAttributeFormulaLoader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*AttributeFormulaConfig)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gAttributeFormulaLoader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*AttributeFormulaConfigRow)
	if !ok {
		err = errors.New("invalid type. not *AttributeFormulaConfigRow")
		logger.ErrorWF("invalid type. not *AttributeFormulaConfigRow", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_formula"))
		return
	}
	config, ok := container.(*AttributeFormulaConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeFormulaConfig")
		logger.ErrorWF("invalid type. not *AttributeFormulaConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_formula"))
		return
	}
	config.ConfigRows[row.Formula_id] = row
	return
}

// GetValue get real map value for json parse
func (*gAttributeFormulaLoader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*AttributeFormulaConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeFormulaConfig")
		logger.ErrorWF("invalid type. not *AttributeFormulaConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_formula"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gAttributeFormulaLoader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*AttributeFormulaConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeFormulaConfig")
		logger.ErrorWF("invalid type. not *AttributeFormulaConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_formula"))
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
type gAttributeFormulaParser struct {
}

// New new config row data
func (*gAttributeFormulaParser) New() interface{} {
	return &AttributeFormulaConfigRow{}
}

// Fields get config fields names
func (*gAttributeFormulaParser) Fields() []string {
	return gAttributeFormulaFields
}

// Parse parse raw data to row data
func (*gAttributeFormulaParser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*AttributeFormulaConfigRow)
	if !ok {
		err = errors.New("invalid type. not *AttributeFormulaConfigRow")
		logger.ErrorWF("invalid type. not *AttributeFormulaConfigRow", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_formula"))
		return
	}
	// compare length
	if len(data) != len(gAttributeFormulaFields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*AttributeFormulaConfigRow",
			zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute_formula"), zap.Int("need_count", len(gAttributeFormulaFields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 formula_id : 公式id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field formula_id 公式id to int32 failed")
			logger.ErrorWF("parse field formula_id 公式id to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_formula"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Formula_id = int32(tmp)
	}

	// parse column 1 type : 公式类型
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field type 公式类型 to int32 failed")
			logger.ErrorWF("parse field type 公式类型 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_formula"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Type = int32(tmp)
	}

	// parse column 2 parameter_1 : 参数1属性id
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field parameter_1 参数1属性id to int32 failed")
			logger.ErrorWF("parse field parameter_1 参数1属性id to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_formula"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Parameter_1 = int32(tmp)
	}

	// parse column 3 parameter_2 : 参数2属性id
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field parameter_2 参数2属性id to int32 failed")
			logger.ErrorWF("parse field parameter_2 参数2属性id to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_formula"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Parameter_2 = int32(tmp)
	}

	// parse column 4 parameter_8 : 参数8属性id
	if data[4] != "" {

		vals := strings.Split(data[4], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field parameter_8 参数8属性id to []int32 failed")
				logger.ErrorWF("parse array field parameter_8 参数8属性id to []int32 failed.",
					zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_formula"),
					// zap.String("field_data",data[4]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Parameter_8 = append(config.Parameter_8, int32(tmp))
		}
	}

	// parse column 5 parameter_3 : 参数3属性id
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field parameter_3 参数3属性id to int32 failed")
			logger.ErrorWF("parse field parameter_3 参数3属性id to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_formula"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Parameter_3 = int32(tmp)
	}

	// parse column 6 parameter_4 : 参数4属性id
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field parameter_4 参数4属性id to int32 failed")
			logger.ErrorWF("parse field parameter_4 参数4属性id to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_formula"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Parameter_4 = int32(tmp)
	}

	// parse column 7 parameter_5 : 参数5属性id
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field parameter_5 参数5属性id to int32 failed")
			logger.ErrorWF("parse field parameter_5 参数5属性id to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_formula"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Parameter_5 = int32(tmp)
	}

	// parse column 8 parameter_9 : 参数9属性id
	if data[8] != "" {

		vals := strings.Split(data[8], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field parameter_9 参数9属性id to []int32 failed")
				logger.ErrorWF("parse array field parameter_9 参数9属性id to []int32 failed.",
					zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_formula"),
					// zap.String("field_data",data[8]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Parameter_9 = append(config.Parameter_9, int32(tmp))
		}
	}

	// parse column 9 parameter_10 : 参数10属性id
	if data[9] != "" {

		vals := strings.Split(data[9], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field parameter_10 参数10属性id to []int32 failed")
				logger.ErrorWF("parse array field parameter_10 参数10属性id to []int32 failed.",
					zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_formula"),
					// zap.String("field_data",data[9]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Parameter_10 = append(config.Parameter_10, int32(tmp))
		}
	}

	// parse column 10 parameter_11 : 参数11属性id
	if data[10] != "" {
		tmp, err = strconv.ParseInt(data[10], 10, 64)
		if err != nil {
			err = errors.New("parse field parameter_11 参数11属性id to int32 failed")
			logger.ErrorWF("parse field parameter_11 参数11属性id to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_formula"),
				zap.String("parse_data", data[10]),
				zap.Error(err))
			return
		}
		config.Parameter_11 = int32(tmp)
	}

	// parse column 11 parameter_12 : 参数12属性id
	if data[11] != "" {
		tmp, err = strconv.ParseInt(data[11], 10, 64)
		if err != nil {
			err = errors.New("parse field parameter_12 参数12属性id to int32 failed")
			logger.ErrorWF("parse field parameter_12 参数12属性id to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute_formula"),
				zap.String("parse_data", data[11]),
				zap.Error(err))
			return
		}
		config.Parameter_12 = int32(tmp)
	}
	return
}

var gAttributeFormulaFields = []string{
	"formula_id",
	"type",
	"parameter_1",
	"parameter_2",
	"parameter_8",
	"parameter_3",
	"parameter_4",
	"parameter_5",
	"parameter_9",
	"parameter_10",
	"parameter_11",
	"parameter_12",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gAttributeFormulaParser{}
	loader := &gAttributeFormulaLoader{}
	var data [][]string
	data, err = load("attribute【属性】.xlsx", "attribute_formula", gAttributeFormulaFields)
	if err != nil {
		logger.ErrorWF("load attribute【属性】.xlsx attribute_formula data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load attribute【属性】.xlsx attribute_formula data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gAttributeFormulaFields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse attribute【属性】.xlsx attribute_formula failed.", zap.Int("row", k), zap.Strings("need", gAttributeFormulaFields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse attribute【属性】.xlsx attribute_formula row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add attribute【属性】.xlsx attribute_formula row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check attribute【属性】.xlsx attribute_formula data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load attribute【属性】.xlsx attribute_formula data success.")
	return
}
