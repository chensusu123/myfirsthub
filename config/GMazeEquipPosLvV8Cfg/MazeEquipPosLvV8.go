package GMazeEquipPosLvV8Cfg

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

// MazeEquipPosLvV8ConfigRow from maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx maze_equip_pos_lv_v8
type MazeEquipPosLvV8ConfigRow struct {
	Order           int32           `json:"order"`           // 序号
	Pos_id          int32           `json:"pos_id"`          // 部位id
	Level           int32           `json:"level"`           // 当前等级
	Next_order      int32           `json:"next_order"`      // 下一级序号
	Need_other      int32           `json:"need_other"`      // 需要其他部位最低等级
	Cost            map[int32]int64 `json:"cost"`            // 升到下一级消耗
	Need_maze_level int32           `json:"need_maze_level"` // 升到下一级需要探险等级
	Add_attr        map[int32]int64 `json:"add_attr"`        // 增加属性
	Show_attr       map[int32]int64 `json:"show_attr"`       // 展示属性
	Show_attr_order []int32         `json:"show_attr_order"` // 展示属性排序
}

// MazeEquipPosLvV8Config from maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx maze_equip_pos_lv_v8
type MazeEquipPosLvV8Config struct {
	ConfigRows map[int32]*MazeEquipPosLvV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipPosLvV8Config {
	ret := &MazeEquipPosLvV8Config{ConfigRows: map[int32]*MazeEquipPosLvV8ConfigRow{}}
	return ret
}

// GetMazeEquipPosLvV8Config get one config by configId
func (c *MazeEquipPosLvV8Config) GetMazeEquipPosLvV8Config(configId int32) *MazeEquipPosLvV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipPosLvV8Config) Get(configId int32) *MazeEquipPosLvV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipPosLvV8Config get all config slice
func (c *MazeEquipPosLvV8Config) GetAllMazeEquipPosLvV8Config() (res []*MazeEquipPosLvV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipPosLvV8Config) GetAll() (res []*MazeEquipPosLvV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipPosLvV8Config

// GetMazeEquipPosLvV8Config pkg func. get one config by configId
func GetMazeEquipPosLvV8Config(configId int32) *MazeEquipPosLvV8ConfigRow {
	return gConfigData.GetMazeEquipPosLvV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipPosLvV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeEquipPosLvV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_equip_pos_lv_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeEquipPosLvV8Config pkg func. get all config slice
func GetAllMazeEquipPosLvV8Config() []*MazeEquipPosLvV8ConfigRow {
	return gConfigData.GetAllMazeEquipPosLvV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipPosLvV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipPosLvV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipPosLvV8ConfigRow from maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx maze_equip_pos_lv_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipPosLvV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_pos_lv_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_pos_lv_v8.json",
		"maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx", "maze_equip_pos_lv_v8",
		&gMazeEquipPosLvV8Parser{}, &gMazeEquipPosLvV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipPosLvV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipPosLvV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipPosLvV8Config))(c)
		return true
	})
}

// RegisterMazeEquipPosLvV8InitCallBack reg config update func (old func)
var RegisterMazeEquipPosLvV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipPosLvV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipPosLvV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipPosLvV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipPosLvV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipPosLvV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipPosLvV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipPosLvV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipPosLvV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipPosLvV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipPosLvV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipPosLvV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipPosLvV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipPosLvV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosLvV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipPosLvV8ConfigRow", zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"),
			zap.String("sheet", "maze_equip_pos_lv_v8"))
		return
	}
	config, ok := container.(*MazeEquipPosLvV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosLvV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipPosLvV8Config", zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"),
			zap.String("sheet", "maze_equip_pos_lv_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipPosLvV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipPosLvV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosLvV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipPosLvV8Config", zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"),
			zap.String("sheet", "maze_equip_pos_lv_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipPosLvV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipPosLvV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosLvV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipPosLvV8Config", zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"),
			zap.String("sheet", "maze_equip_pos_lv_v8"))
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
type gMazeEquipPosLvV8Parser struct {
}

// New new config row data
func (*gMazeEquipPosLvV8Parser) New() interface{} {
	return &MazeEquipPosLvV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipPosLvV8Parser) Fields() []string {
	return gMazeEquipPosLvV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipPosLvV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipPosLvV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosLvV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipPosLvV8ConfigRow", zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"),
			zap.String("sheet", "maze_equip_pos_lv_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipPosLvV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipPosLvV8ConfigRow",
			zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"),
			zap.String("sheet", "maze_equip_pos_lv_v8"), zap.Int("need_count", len(gMazeEquipPosLvV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 序号
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 序号 to int32 failed")
			logger.ErrorWF("parse field order 序号 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 pos_id : 部位id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field pos_id 部位id to int32 failed")
			logger.ErrorWF("parse field pos_id 部位id to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Pos_id = int32(tmp)
	}

	// parse column 2 level : 当前等级
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field level 当前等级 to int32 failed")
			logger.ErrorWF("parse field level 当前等级 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Level = int32(tmp)
	}

	// parse column 3 next_order : 下一级序号
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field next_order 下一级序号 to int32 failed")
			logger.ErrorWF("parse field next_order 下一级序号 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Next_order = int32(tmp)
	}

	// parse column 4 need_other : 需要其他部位最低等级
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field need_other 需要其他部位最低等级 to int32 failed")
			logger.ErrorWF("parse field need_other 需要其他部位最低等级 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Need_other = int32(tmp)
	}

	// parse column 5 cost : 升到下一级消耗
	if data[5] != "" {

		config.Cost = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[5], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field cost 升到下一级消耗 to key int32 failed")
				logger.ErrorWF("parse map field cost 升到下一级消耗 to key int32 failed.",
					zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field cost 升到下一级消耗 to value int64 failed")
				logger.ErrorWF("parse map field cost 升到下一级消耗 to value int64 failed.",
					zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Cost[key] = value
		}
	}

	// parse column 6 need_maze_level : 升到下一级需要探险等级
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field need_maze_level 升到下一级需要探险等级 to int32 failed")
			logger.ErrorWF("parse field need_maze_level 升到下一级需要探险等级 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Need_maze_level = int32(tmp)
	}

	// parse column 7 add_attr : 增加属性
	if data[7] != "" {

		config.Add_attr = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[7], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field add_attr 增加属性 to key int32 failed")
				logger.ErrorWF("parse map field add_attr 增加属性 to key int32 failed.",
					zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_v8"),
					// zap.String("field_data",data[7]),
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
					zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_v8"),
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

	// parse column 8 show_attr : 展示属性
	if data[8] != "" {

		config.Show_attr = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[8], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field show_attr 展示属性 to key int32 failed")
				logger.ErrorWF("parse map field show_attr 展示属性 to key int32 failed.",
					zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_v8"),
					// zap.String("field_data",data[8]),
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
					zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_v8"),
					// zap.String("field_data",data[8]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Show_attr[key] = value
		}
	}

	// parse column 9 show_attr_order : 展示属性排序
	if data[9] != "" {

		vals := strings.Split(data[9], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field show_attr_order 展示属性排序 to []int32 failed")
				logger.ErrorWF("parse array field show_attr_order 展示属性排序 to []int32 failed.",
					zap.String("xlsx", "maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx"), zap.String("sheet", "maze_equip_pos_lv_v8"),
					// zap.String("field_data",data[9]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Show_attr_order = append(config.Show_attr_order, int32(tmp))
		}
	}
	return
}

var gMazeEquipPosLvV8Fields = []string{
	"order",
	"pos_id",
	"level",
	"next_order",
	"need_other",
	"cost",
	"need_maze_level",
	"add_attr",
	"show_attr",
	"show_attr_order",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipPosLvV8Parser{}
	loader := &gMazeEquipPosLvV8Loader{}
	var data [][]string
	data, err = load("maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx", "maze_equip_pos_lv_v8", gMazeEquipPosLvV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx maze_equip_pos_lv_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx maze_equip_pos_lv_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipPosLvV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx maze_equip_pos_lv_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipPosLvV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx maze_equip_pos_lv_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx maze_equip_pos_lv_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx maze_equip_pos_lv_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_pos_lv_v8【迷宫-装备-部位强化等级】.xlsx maze_equip_pos_lv_v8 data success.")
	return
}
