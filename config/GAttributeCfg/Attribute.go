package GAttributeCfg

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

// AttributeConfigRow from attribute【属性】.xlsx attribute
type AttributeConfigRow struct {
	Id                    int32  `json:"id"`                    // 属性ID
	Is_into_buff          int32  `json:"is_into_buff"`          // 是否进入buff中心
	Formula_parameter_id  int32  `json:"formula_parameter_id"`  // 公式属性id
	Name                  string `json:"name"`                  // 属性名称
	Type                  int32  `json:"type"`                  // 属性类型
	Figure                int32  `json:"figure"`                // 数值类型
	ShowType              int32  `json:"showType"`              // 是否显示在属性面板
	Quotiety              int32  `json:"quotiety"`              // 战力系数
	Atk_coefficient       int32  `json:"atk_coefficient"`       // 攻击系数
	Def_coefficient       int32  `json:"def_coefficient"`       // 防御系数
	Durable_coefficient   int32  `json:"durable_coefficient"`   // 耐久系数
	Movespeed_coefficient int32  `json:"movespeed_coefficient"` // 移动系数
	Load_coefficient      int32  `json:"load_coefficient"`      // 载重系数
	Ishide                int32  `json:"ishide"`                // 是否隐藏
	Enemy_type            int32  `json:"enemy_type"`            // 是否是npc怪
	To_attr               int32  `json:"to_attr"`               // 映射属性id
	Action_range          int32  `json:"action_range"`          // 作用范围（1-旗舰 2-非旗舰）
	Origin                int32  `json:"origin"`                // 本源属性id
	Use_type              int32  `json:"use_type"`              // 属性使用类型
}

// AttributeConfig from attribute【属性】.xlsx attribute
type AttributeConfig struct {
	ConfigRows map[int32]*AttributeConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *AttributeConfig {
	ret := &AttributeConfig{ConfigRows: map[int32]*AttributeConfigRow{}}
	return ret
}

// GetAttributeConfig get one config by configId
func (c *AttributeConfig) GetAttributeConfig(configId int32) *AttributeConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *AttributeConfig) Get(configId int32) *AttributeConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllAttributeConfig get all config slice
func (c *AttributeConfig) GetAllAttributeConfig() (res []*AttributeConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *AttributeConfig) GetAll() (res []*AttributeConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *AttributeConfig

// GetAttributeConfig pkg func. get one config by configId
func GetAttributeConfig(configId int32) *AttributeConfigRow {
	return gConfigData.GetAttributeConfig(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *AttributeConfigRow {
	return gConfigData.Get(configId)
}

// GetAllAttributeConfig pkg func. get all config slice
func GetAllAttributeConfig() []*AttributeConfigRow {
	return gConfigData.GetAllAttributeConfig()
}

// GetAll pkg func. get all config slice
func GetAll() []*AttributeConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*AttributeConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "AttributeConfigRow from attribute【属性】.xlsx attribute"
}

// GetRawValue get raw data
func GetRawValue() *AttributeConfig {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "attribute"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("attribute.json",
		"attribute【属性】.xlsx", "attribute",
		&gAttributeParser{}, &gAttributeLoader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*AttributeConfig)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *AttributeConfig) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*AttributeConfig))(c)
		return true
	})
}

// RegisterAttributeInitCallBack reg config update func (old func)
var RegisterAttributeInitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*AttributeConfig)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*AttributeConfig)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *AttributeConfig) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*AttributeConfig) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*AttributeConfig) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gAttributeLoader struct {
}

// NewContainer new data container pointer
func (*gAttributeLoader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gAttributeLoader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*AttributeConfig)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gAttributeLoader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*AttributeConfig)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gAttributeLoader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*AttributeConfigRow)
	if !ok {
		err = errors.New("invalid type. not *AttributeConfigRow")
		logger.ErrorWF("invalid type. not *AttributeConfigRow", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute"))
		return
	}
	config, ok := container.(*AttributeConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeConfig")
		logger.ErrorWF("invalid type. not *AttributeConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gAttributeLoader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*AttributeConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeConfig")
		logger.ErrorWF("invalid type. not *AttributeConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gAttributeLoader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*AttributeConfig)
	if !ok {
		err = errors.New("invalid type. not *AttributeConfig")
		logger.ErrorWF("invalid type. not *AttributeConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute"))
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
type gAttributeParser struct {
}

// New new config row data
func (*gAttributeParser) New() interface{} {
	return &AttributeConfigRow{}
}

// Fields get config fields names
func (*gAttributeParser) Fields() []string {
	return gAttributeFields
}

// Parse parse raw data to row data
func (*gAttributeParser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*AttributeConfigRow)
	if !ok {
		err = errors.New("invalid type. not *AttributeConfigRow")
		logger.ErrorWF("invalid type. not *AttributeConfigRow", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute"))
		return
	}
	// compare length
	if len(data) != len(gAttributeFields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*AttributeConfigRow",
			zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "attribute"), zap.Int("need_count", len(gAttributeFields)),
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
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 is_into_buff : 是否进入buff中心
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field is_into_buff 是否进入buff中心 to int32 failed")
			logger.ErrorWF("parse field is_into_buff 是否进入buff中心 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Is_into_buff = int32(tmp)
	}

	// parse column 2 formula_parameter_id : 公式属性id
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field formula_parameter_id 公式属性id to int32 failed")
			logger.ErrorWF("parse field formula_parameter_id 公式属性id to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Formula_parameter_id = int32(tmp)
	}

	// parse column 3 name : 属性名称
	if data[3] != "" {
		config.Name = data[3]
	}

	// parse column 4 type : 属性类型
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field type 属性类型 to int32 failed")
			logger.ErrorWF("parse field type 属性类型 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Type = int32(tmp)
	}

	// parse column 5 figure : 数值类型
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field figure 数值类型 to int32 failed")
			logger.ErrorWF("parse field figure 数值类型 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Figure = int32(tmp)
	}

	// parse column 6 showType : 是否显示在属性面板
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field showType 是否显示在属性面板 to int32 failed")
			logger.ErrorWF("parse field showType 是否显示在属性面板 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.ShowType = int32(tmp)
	}

	// parse column 7 quotiety : 战力系数
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field quotiety 战力系数 to int32 failed")
			logger.ErrorWF("parse field quotiety 战力系数 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Quotiety = int32(tmp)
	}

	// parse column 8 atk_coefficient : 攻击系数
	if data[8] != "" {
		tmp, err = strconv.ParseInt(data[8], 10, 64)
		if err != nil {
			err = errors.New("parse field atk_coefficient 攻击系数 to int32 failed")
			logger.ErrorWF("parse field atk_coefficient 攻击系数 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[8]),
				zap.Error(err))
			return
		}
		config.Atk_coefficient = int32(tmp)
	}

	// parse column 9 def_coefficient : 防御系数
	if data[9] != "" {
		tmp, err = strconv.ParseInt(data[9], 10, 64)
		if err != nil {
			err = errors.New("parse field def_coefficient 防御系数 to int32 failed")
			logger.ErrorWF("parse field def_coefficient 防御系数 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[9]),
				zap.Error(err))
			return
		}
		config.Def_coefficient = int32(tmp)
	}

	// parse column 10 durable_coefficient : 耐久系数
	if data[10] != "" {
		tmp, err = strconv.ParseInt(data[10], 10, 64)
		if err != nil {
			err = errors.New("parse field durable_coefficient 耐久系数 to int32 failed")
			logger.ErrorWF("parse field durable_coefficient 耐久系数 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[10]),
				zap.Error(err))
			return
		}
		config.Durable_coefficient = int32(tmp)
	}

	// parse column 11 movespeed_coefficient : 移动系数
	if data[11] != "" {
		tmp, err = strconv.ParseInt(data[11], 10, 64)
		if err != nil {
			err = errors.New("parse field movespeed_coefficient 移动系数 to int32 failed")
			logger.ErrorWF("parse field movespeed_coefficient 移动系数 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[11]),
				zap.Error(err))
			return
		}
		config.Movespeed_coefficient = int32(tmp)
	}

	// parse column 12 load_coefficient : 载重系数
	if data[12] != "" {
		tmp, err = strconv.ParseInt(data[12], 10, 64)
		if err != nil {
			err = errors.New("parse field load_coefficient 载重系数 to int32 failed")
			logger.ErrorWF("parse field load_coefficient 载重系数 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[12]),
				zap.Error(err))
			return
		}
		config.Load_coefficient = int32(tmp)
	}

	// parse column 13 ishide : 是否隐藏
	if data[13] != "" {
		tmp, err = strconv.ParseInt(data[13], 10, 64)
		if err != nil {
			err = errors.New("parse field ishide 是否隐藏 to int32 failed")
			logger.ErrorWF("parse field ishide 是否隐藏 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[13]),
				zap.Error(err))
			return
		}
		config.Ishide = int32(tmp)
	}

	// parse column 14 enemy_type : 是否是npc怪
	if data[14] != "" {
		tmp, err = strconv.ParseInt(data[14], 10, 64)
		if err != nil {
			err = errors.New("parse field enemy_type 是否是npc怪 to int32 failed")
			logger.ErrorWF("parse field enemy_type 是否是npc怪 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[14]),
				zap.Error(err))
			return
		}
		config.Enemy_type = int32(tmp)
	}

	// parse column 15 to_attr : 映射属性id
	if data[15] != "" {
		tmp, err = strconv.ParseInt(data[15], 10, 64)
		if err != nil {
			err = errors.New("parse field to_attr 映射属性id to int32 failed")
			logger.ErrorWF("parse field to_attr 映射属性id to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[15]),
				zap.Error(err))
			return
		}
		config.To_attr = int32(tmp)
	}

	// parse column 16 action_range : 作用范围（1-旗舰 2-非旗舰）
	if data[16] != "" {
		tmp, err = strconv.ParseInt(data[16], 10, 64)
		if err != nil {
			err = errors.New("parse field action_range 作用范围（1-旗舰 2-非旗舰） to int32 failed")
			logger.ErrorWF("parse field action_range 作用范围（1-旗舰 2-非旗舰） to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[16]),
				zap.Error(err))
			return
		}
		config.Action_range = int32(tmp)
	}

	// parse column 17 origin : 本源属性id
	if data[17] != "" {
		tmp, err = strconv.ParseInt(data[17], 10, 64)
		if err != nil {
			err = errors.New("parse field origin 本源属性id to int32 failed")
			logger.ErrorWF("parse field origin 本源属性id to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[17]),
				zap.Error(err))
			return
		}
		config.Origin = int32(tmp)
	}

	// parse column 18 use_type : 属性使用类型
	if data[18] != "" {
		tmp, err = strconv.ParseInt(data[18], 10, 64)
		if err != nil {
			err = errors.New("parse field use_type 属性使用类型 to int32 failed")
			logger.ErrorWF("parse field use_type 属性使用类型 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "attribute"),
				zap.String("parse_data", data[18]),
				zap.Error(err))
			return
		}
		config.Use_type = int32(tmp)
	}
	return
}

var gAttributeFields = []string{
	"id",
	"is_into_buff",
	"formula_parameter_id",
	"name",
	"type",
	"figure",
	"showType",
	"quotiety",
	"atk_coefficient",
	"def_coefficient",
	"durable_coefficient",
	"movespeed_coefficient",
	"load_coefficient",
	"ishide",
	"enemy_type",
	"to_attr",
	"action_range",
	"origin",
	"use_type",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gAttributeParser{}
	loader := &gAttributeLoader{}
	var data [][]string
	data, err = load("attribute【属性】.xlsx", "attribute", gAttributeFields)
	if err != nil {
		logger.ErrorWF("load attribute【属性】.xlsx attribute data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load attribute【属性】.xlsx attribute data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gAttributeFields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse attribute【属性】.xlsx attribute failed.", zap.Int("row", k), zap.Strings("need", gAttributeFields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse attribute【属性】.xlsx attribute row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add attribute【属性】.xlsx attribute row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check attribute【属性】.xlsx attribute data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load attribute【属性】.xlsx attribute data success.")
	return
}
