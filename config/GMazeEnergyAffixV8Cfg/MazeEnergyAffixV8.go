package GMazeEnergyAffixV8Cfg

import (
	"context"
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

// MazeEnergyAffixV8ConfigRow from maze_energy_affix_v8【迷宫-能力词条】.xlsx maze_energy_affix_v8
type MazeEnergyAffixV8ConfigRow struct {
	Order                int32           `json:"order"`                // 词条id
	Affix_group_id       int32           `json:"affix_group_id"`       // 词条组id
	Affix_group_order    int32           `json:"affix_group_order"`    // 词条组内编号
	Affix_lv             int32           `json:"affix_lv"`             // 词条等级
	Font_affix_condition []int32         `json:"font_affix_condition"` // 前置词条id
	Use_num_max          int32           `json:"use_num_max"`          // 可使用最大次数
	Weight               int32           `json:"weight"`               // 随机权重
	Add_attr             map[int32]int64 `json:"add_attr"`             // 实际增加的属性
	Affix_name           string          `json:"affix_name"`           // 词条名字
	Affix_desc           string          `json:"affix_desc"`           // 词条的描述
	Affix_wildcard       []int32         `json:"affix_wildcard"`       // 描述通配符id（属性id）
}

// MazeEnergyAffixV8Config from maze_energy_affix_v8【迷宫-能力词条】.xlsx maze_energy_affix_v8
type MazeEnergyAffixV8Config struct {
	ConfigRows map[int32]*MazeEnergyAffixV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEnergyAffixV8Config {
	ret := &MazeEnergyAffixV8Config{ConfigRows: map[int32]*MazeEnergyAffixV8ConfigRow{}}
	return ret
}

// GetMazeEnergyAffixV8Config get one config by configId
func (c *MazeEnergyAffixV8Config) GetMazeEnergyAffixV8Config(configId int32) *MazeEnergyAffixV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEnergyAffixV8Config) Get(configId int32) *MazeEnergyAffixV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEnergyAffixV8Config get all config slice
func (c *MazeEnergyAffixV8Config) GetAllMazeEnergyAffixV8Config() (res []*MazeEnergyAffixV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEnergyAffixV8Config) GetAll() (res []*MazeEnergyAffixV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEnergyAffixV8Config

// GetMazeEnergyAffixV8Config pkg func. get one config by configId
func GetMazeEnergyAffixV8Config(configId int32) *MazeEnergyAffixV8ConfigRow {
	return gConfigData.GetMazeEnergyAffixV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeEnergyAffixV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeEnergyAffixV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_energy_affix_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeEnergyAffixV8Config pkg func. get all config slice
func GetAllMazeEnergyAffixV8Config() []*MazeEnergyAffixV8ConfigRow {
	return gConfigData.GetAllMazeEnergyAffixV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEnergyAffixV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEnergyAffixV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEnergyAffixV8ConfigRow from maze_energy_affix_v8【迷宫-能力词条】.xlsx maze_energy_affix_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEnergyAffixV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_energy_affix_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_energy_affix_v8.json",
		"maze_energy_affix_v8【迷宫-能力词条】.xlsx", "maze_energy_affix_v8",
		&gMazeEnergyAffixV8Parser{}, &gMazeEnergyAffixV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEnergyAffixV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEnergyAffixV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEnergyAffixV8Config))(c)
		return true
	})
}

// RegisterMazeEnergyAffixV8InitCallBack reg config update func (old func)
var RegisterMazeEnergyAffixV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEnergyAffixV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEnergyAffixV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEnergyAffixV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEnergyAffixV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEnergyAffixV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEnergyAffixV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEnergyAffixV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEnergyAffixV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEnergyAffixV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEnergyAffixV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEnergyAffixV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEnergyAffixV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEnergyAffixV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixV8ConfigRow", zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"),
			zap.String("sheet", "maze_energy_affix_v8"))
		return
	}
	config, ok := container.(*MazeEnergyAffixV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixV8Config", zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"),
			zap.String("sheet", "maze_energy_affix_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEnergyAffixV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEnergyAffixV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixV8Config", zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"),
			zap.String("sheet", "maze_energy_affix_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEnergyAffixV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEnergyAffixV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixV8Config", zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"),
			zap.String("sheet", "maze_energy_affix_v8"))
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
type gMazeEnergyAffixV8Parser struct {
}

// New new config row data
func (*gMazeEnergyAffixV8Parser) New() interface{} {
	return &MazeEnergyAffixV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEnergyAffixV8Parser) Fields() []string {
	return gMazeEnergyAffixV8Fields
}

// Parse parse raw data to row data
func (*gMazeEnergyAffixV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEnergyAffixV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixV8ConfigRow", zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"),
			zap.String("sheet", "maze_energy_affix_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEnergyAffixV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEnergyAffixV8ConfigRow",
			zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"),
			zap.String("sheet", "maze_energy_affix_v8"), zap.Int("need_count", len(gMazeEnergyAffixV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 词条id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 词条id to int32 failed")
			logger.ErrorWF("parse field order 词条id to int32 failed.",
				zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"), zap.String("sheet", "maze_energy_affix_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 affix_group_id : 词条组id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field affix_group_id 词条组id to int32 failed")
			logger.ErrorWF("parse field affix_group_id 词条组id to int32 failed.",
				zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"), zap.String("sheet", "maze_energy_affix_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Affix_group_id = int32(tmp)
	}

	// parse column 2 affix_group_order : 词条组内编号
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field affix_group_order 词条组内编号 to int32 failed")
			logger.ErrorWF("parse field affix_group_order 词条组内编号 to int32 failed.",
				zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"), zap.String("sheet", "maze_energy_affix_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Affix_group_order = int32(tmp)
	}

	// parse column 3 affix_lv : 词条等级
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field affix_lv 词条等级 to int32 failed")
			logger.ErrorWF("parse field affix_lv 词条等级 to int32 failed.",
				zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"), zap.String("sheet", "maze_energy_affix_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Affix_lv = int32(tmp)
	}

	// parse column 4 font_affix_condition : 前置词条id
	if data[4] != "" {

		vals := strings.Split(data[4], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field font_affix_condition 前置词条id to []int32 failed")
				logger.ErrorWF("parse array field font_affix_condition 前置词条id to []int32 failed.",
					zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"), zap.String("sheet", "maze_energy_affix_v8"),
					// zap.String("field_data",data[4]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Font_affix_condition = append(config.Font_affix_condition, int32(tmp))
		}
	}

	// parse column 5 use_num_max : 可使用最大次数
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field use_num_max 可使用最大次数 to int32 failed")
			logger.ErrorWF("parse field use_num_max 可使用最大次数 to int32 failed.",
				zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"), zap.String("sheet", "maze_energy_affix_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Use_num_max = int32(tmp)
	}

	// parse column 6 weight : 随机权重
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field weight 随机权重 to int32 failed")
			logger.ErrorWF("parse field weight 随机权重 to int32 failed.",
				zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"), zap.String("sheet", "maze_energy_affix_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Weight = int32(tmp)
	}

	// parse column 7 add_attr : 实际增加的属性
	if data[7] != "" {

		config.Add_attr = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[7], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field add_attr 实际增加的属性 to key int32 failed")
				logger.ErrorWF("parse map field add_attr 实际增加的属性 to key int32 failed.",
					zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"), zap.String("sheet", "maze_energy_affix_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field add_attr 实际增加的属性 to value int64 failed")
				logger.ErrorWF("parse map field add_attr 实际增加的属性 to value int64 failed.",
					zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"), zap.String("sheet", "maze_energy_affix_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Add_attr[key] = value
		}
	}

	// parse column 8 affix_name : 词条名字
	if data[8] != "" {
		config.Affix_name = data[8]
	}

	// parse column 9 affix_desc : 词条的描述
	if data[9] != "" {
		config.Affix_desc = data[9]
	}

	// parse column 10 affix_wildcard : 描述通配符id（属性id）
	if data[10] != "" {

		vals := strings.Split(data[10], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field affix_wildcard 描述通配符id（属性id） to []int32 failed")
				logger.ErrorWF("parse array field affix_wildcard 描述通配符id（属性id） to []int32 failed.",
					zap.String("xlsx", "maze_energy_affix_v8【迷宫-能力词条】.xlsx"), zap.String("sheet", "maze_energy_affix_v8"),
					// zap.String("field_data",data[10]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Affix_wildcard = append(config.Affix_wildcard, int32(tmp))
		}
	}
	return
}

var gMazeEnergyAffixV8Fields = []string{
	"order",
	"affix_group_id",
	"affix_group_order",
	"affix_lv",
	"font_affix_condition",
	"use_num_max",
	"weight",
	"add_attr",
	"affix_name",
	"affix_desc",
	"affix_wildcard",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEnergyAffixV8Parser{}
	loader := &gMazeEnergyAffixV8Loader{}
	var data [][]string
	data, err = load("maze_energy_affix_v8【迷宫-能力词条】.xlsx", "maze_energy_affix_v8", gMazeEnergyAffixV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_energy_affix_v8【迷宫-能力词条】.xlsx maze_energy_affix_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_energy_affix_v8【迷宫-能力词条】.xlsx maze_energy_affix_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEnergyAffixV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_energy_affix_v8【迷宫-能力词条】.xlsx maze_energy_affix_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEnergyAffixV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_energy_affix_v8【迷宫-能力词条】.xlsx maze_energy_affix_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_energy_affix_v8【迷宫-能力词条】.xlsx maze_energy_affix_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_energy_affix_v8【迷宫-能力词条】.xlsx maze_energy_affix_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_energy_affix_v8【迷宫-能力词条】.xlsx maze_energy_affix_v8 data success.")
	return
}
