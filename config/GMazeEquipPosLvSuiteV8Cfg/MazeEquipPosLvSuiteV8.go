package GMazeEquipPosLvSuiteV8Cfg

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

// MazeEquipPosLvSuiteV8ConfigRow from maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx maze_equip_pos_lv_suite_v8
type MazeEquipPosLvSuiteV8ConfigRow struct {
	Level           int32           `json:"level"`           // 套装等级
	Add_attr        map[int32]int64 `json:"add_attr"`        // 增加属性
	Show_attr       map[int32]int64 `json:"show_attr"`       // 展示属性
	Show_attr_order []int32         `json:"show_attr_order"` // 展示属性排序
	Next_level      int32           `json:"next_level"`      // 下一级套装等级
}

// MazeEquipPosLvSuiteV8Config from maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx maze_equip_pos_lv_suite_v8
type MazeEquipPosLvSuiteV8Config struct {
	ConfigRows map[int32]*MazeEquipPosLvSuiteV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipPosLvSuiteV8Config {
	ret := &MazeEquipPosLvSuiteV8Config{ConfigRows: map[int32]*MazeEquipPosLvSuiteV8ConfigRow{}}
	return ret
}

// GetMazeEquipPosLvSuiteV8Config get one config by configId
func (c *MazeEquipPosLvSuiteV8Config) GetMazeEquipPosLvSuiteV8Config(configId int32) *MazeEquipPosLvSuiteV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipPosLvSuiteV8Config) Get(configId int32) *MazeEquipPosLvSuiteV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipPosLvSuiteV8Config get all config slice
func (c *MazeEquipPosLvSuiteV8Config) GetAllMazeEquipPosLvSuiteV8Config() (res []*MazeEquipPosLvSuiteV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipPosLvSuiteV8Config) GetAll() (res []*MazeEquipPosLvSuiteV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipPosLvSuiteV8Config

// GetMazeEquipPosLvSuiteV8Config pkg func. get one config by configId
func GetMazeEquipPosLvSuiteV8Config(configId int32) *MazeEquipPosLvSuiteV8ConfigRow {
	return gConfigData.GetMazeEquipPosLvSuiteV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipPosLvSuiteV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeEquipPosLvSuiteV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_equip_pos_lv_suite_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeEquipPosLvSuiteV8Config pkg func. get all config slice
func GetAllMazeEquipPosLvSuiteV8Config() []*MazeEquipPosLvSuiteV8ConfigRow {
	return gConfigData.GetAllMazeEquipPosLvSuiteV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipPosLvSuiteV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipPosLvSuiteV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipPosLvSuiteV8ConfigRow from maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx maze_equip_pos_lv_suite_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipPosLvSuiteV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_pos_lv_suite_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_pos_lv_suite_v8.json",
		"maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx", "maze_equip_pos_lv_suite_v8",
		&gMazeEquipPosLvSuiteV8Parser{}, &gMazeEquipPosLvSuiteV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipPosLvSuiteV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipPosLvSuiteV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipPosLvSuiteV8Config))(c)
		return true
	})
}

// RegisterMazeEquipPosLvSuiteV8InitCallBack reg config update func (old func)
var RegisterMazeEquipPosLvSuiteV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipPosLvSuiteV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipPosLvSuiteV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipPosLvSuiteV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipPosLvSuiteV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipPosLvSuiteV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipPosLvSuiteV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipPosLvSuiteV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipPosLvSuiteV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipPosLvSuiteV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipPosLvSuiteV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipPosLvSuiteV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipPosLvSuiteV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipPosLvSuiteV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosLvSuiteV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipPosLvSuiteV8ConfigRow", zap.String("xlsx", "maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx"),
			zap.String("sheet", "maze_equip_pos_lv_suite_v8"))
		return
	}
	config, ok := container.(*MazeEquipPosLvSuiteV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosLvSuiteV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipPosLvSuiteV8Config", zap.String("xlsx", "maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx"),
			zap.String("sheet", "maze_equip_pos_lv_suite_v8"))
		return
	}
	config.ConfigRows[row.Level] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipPosLvSuiteV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipPosLvSuiteV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosLvSuiteV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipPosLvSuiteV8Config", zap.String("xlsx", "maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx"),
			zap.String("sheet", "maze_equip_pos_lv_suite_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipPosLvSuiteV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipPosLvSuiteV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosLvSuiteV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipPosLvSuiteV8Config", zap.String("xlsx", "maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx"),
			zap.String("sheet", "maze_equip_pos_lv_suite_v8"))
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
type gMazeEquipPosLvSuiteV8Parser struct {
}

// New new config row data
func (*gMazeEquipPosLvSuiteV8Parser) New() interface{} {
	return &MazeEquipPosLvSuiteV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipPosLvSuiteV8Parser) Fields() []string {
	return gMazeEquipPosLvSuiteV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipPosLvSuiteV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipPosLvSuiteV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosLvSuiteV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipPosLvSuiteV8ConfigRow", zap.String("xlsx", "maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx"),
			zap.String("sheet", "maze_equip_pos_lv_suite_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipPosLvSuiteV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipPosLvSuiteV8ConfigRow",
			zap.String("xlsx", "maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx"),
			zap.String("sheet", "maze_equip_pos_lv_suite_v8"), zap.Int("need_count", len(gMazeEquipPosLvSuiteV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 level : 套装等级
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field level 套装等级 to int32 failed")
			logger.ErrorWF("parse field level 套装等级 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_suite_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Level = int32(tmp)
	}

	// parse column 1 add_attr : 增加属性
	if data[1] != "" {

		config.Add_attr = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[1], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field add_attr 增加属性 to key int32 failed")
				logger.ErrorWF("parse map field add_attr 增加属性 to key int32 failed.",
					zap.String("xlsx", "maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_suite_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field add_attr 增加属性 to value int64 failed")
				logger.ErrorWF("parse map field add_attr 增加属性 to value int64 failed.",
					zap.String("xlsx", "maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_suite_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Add_attr[key] = value
		}
	}

	// parse column 2 show_attr : 展示属性
	if data[2] != "" {

		config.Show_attr = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[2], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field show_attr 展示属性 to key int32 failed")
				logger.ErrorWF("parse map field show_attr 展示属性 to key int32 failed.",
					zap.String("xlsx", "maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_suite_v8"),
					// zap.String("field_data",data[2]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field show_attr 展示属性 to value int64 failed")
				logger.ErrorWF("parse map field show_attr 展示属性 to value int64 failed.",
					zap.String("xlsx", "maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_suite_v8"),
					// zap.String("field_data",data[2]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Show_attr[key] = value
		}
	}

	// parse column 3 show_attr_order : 展示属性排序
	if data[3] != "" {

		vals := strings.Split(data[3], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field show_attr_order 展示属性排序 to []int32 failed")
				logger.ErrorWF("parse array field show_attr_order 展示属性排序 to []int32 failed.",
					zap.String("xlsx", "maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_suite_v8"),
					// zap.String("field_data",data[3]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Show_attr_order = append(config.Show_attr_order, int32(tmp))
		}
	}

	// parse column 4 next_level : 下一级套装等级
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field next_level 下一级套装等级 to int32 failed")
			logger.ErrorWF("parse field next_level 下一级套装等级 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_suite_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Next_level = int32(tmp)
	}
	return
}

var gMazeEquipPosLvSuiteV8Fields = []string{
	"level",
	"add_attr",
	"show_attr",
	"show_attr_order",
	"next_level",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipPosLvSuiteV8Parser{}
	loader := &gMazeEquipPosLvSuiteV8Loader{}
	var data [][]string
	data, err = load("maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx", "maze_equip_pos_lv_suite_v8", gMazeEquipPosLvSuiteV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx maze_equip_pos_lv_suite_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx maze_equip_pos_lv_suite_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipPosLvSuiteV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx maze_equip_pos_lv_suite_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipPosLvSuiteV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx maze_equip_pos_lv_suite_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx maze_equip_pos_lv_suite_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx maze_equip_pos_lv_suite_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_pos_lv_suite_v8【迷宫-装备-部位强化等级套装】.xlsx maze_equip_pos_lv_suite_v8 data success.")
	return
}
