package GMazeBrushFoeV8Cfg

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

// MazeBrushFoeV8ConfigRow from maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx maze_brush_foe_v8
type MazeBrushFoeV8ConfigRow struct {
	Id                                  int32   `json:"id"`                                  // 序号
	Barries_id                          int32   `json:"barries_id"`                          // 关卡id
	Brush_area_id                       int32   `json:"brush_area_id"`                       // 刷怪区域id
	Group_id                            int32   `json:"group_id"`                            // 波
	In_group_order                      int32   `json:"in_group_order"`                      // 一波内刷新顺序
	Interval_time                       int32   `json:"interval_time"`                       // 每波刷新时间（毫秒）
	Display_wave_id                     int32   `json:"display_wave_id"`                     // 页面展示的波次
	Monsters_id                         []int32 `json:"monsters_id"`                         // 怪物id列表
	Wave_kongfu                         int32   `json:"wave_kongfu"`                         // 每波增加的武力值
	Min_brushtime                       int32   `json:"min_brushtime"`                       // 每波保底刷怪时间（毫秒）
	Iskill_all                          int32   `json:"iskill_all"`                          // 是否需要杀完上波怪再刷本波
	Index                               int32   `json:"index"`                               // 战斗区域id
	Brushing_monsters_coordinate        []int32 `json:"Brushing_monsters_coordinate"`        // 区域刷怪坐标列表（横百分段值,竖百分段值）
	Brushing_monsters_position          []int32 `json:"Brushing_monsters_position"`          // 区域刷怪点坐标（横百分段值,竖百分段值）,范围
	Brushing_monsters_relative_position []int32 `json:"Brushing_monsters_relative_position"` // 相对刷怪点（角度0~360,距离,范围）
	Front_group_order                   []int32 `json:"front_group_order"`                   // 立即刷怪前置条件波次
	Is_def_show                         int32   `json:"is_def_show"`                         // 是否需要默认显示
}

// MazeBrushFoeV8Config from maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx maze_brush_foe_v8
type MazeBrushFoeV8Config struct {
	ConfigRows map[int32]*MazeBrushFoeV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeBrushFoeV8Config {
	ret := &MazeBrushFoeV8Config{ConfigRows: map[int32]*MazeBrushFoeV8ConfigRow{}}
	return ret
}

// GetMazeBrushFoeV8Config get one config by configId
func (c *MazeBrushFoeV8Config) GetMazeBrushFoeV8Config(configId int32) *MazeBrushFoeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeBrushFoeV8Config) Get(configId int32) *MazeBrushFoeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeBrushFoeV8Config get all config slice
func (c *MazeBrushFoeV8Config) GetAllMazeBrushFoeV8Config() (res []*MazeBrushFoeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeBrushFoeV8Config) GetAll() (res []*MazeBrushFoeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeBrushFoeV8Config

// GetMazeBrushFoeV8Config pkg func. get one config by configId
func GetMazeBrushFoeV8Config(configId int32) *MazeBrushFoeV8ConfigRow {
	return gConfigData.GetMazeBrushFoeV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeBrushFoeV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeBrushFoeV8Config pkg func. get all config slice
func GetAllMazeBrushFoeV8Config() []*MazeBrushFoeV8ConfigRow {
	return gConfigData.GetAllMazeBrushFoeV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeBrushFoeV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeBrushFoeV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeBrushFoeV8ConfigRow from maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx maze_brush_foe_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeBrushFoeV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_brush_foe_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_brush_foe_v8.json",
		"maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx", "maze_brush_foe_v8",
		&gMazeBrushFoeV8Parser{}, &gMazeBrushFoeV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeBrushFoeV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeBrushFoeV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeBrushFoeV8Config))(c)
		return true
	})
}

// RegisterMazeBrushFoeV8InitCallBack reg config update func (old func)
var RegisterMazeBrushFoeV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeBrushFoeV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeBrushFoeV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeBrushFoeV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeBrushFoeV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeBrushFoeV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeBrushFoeV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeBrushFoeV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeBrushFoeV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeBrushFoeV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeBrushFoeV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeBrushFoeV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeBrushFoeV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeBrushFoeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushFoeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBrushFoeV8ConfigRow", zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"),
			zap.String("sheet", "maze_brush_foe_v8"))
		return
	}
	config, ok := container.(*MazeBrushFoeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushFoeV8Config")
		logger.ErrorWF("invalid type. not *MazeBrushFoeV8Config", zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"),
			zap.String("sheet", "maze_brush_foe_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeBrushFoeV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeBrushFoeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushFoeV8Config")
		logger.ErrorWF("invalid type. not *MazeBrushFoeV8Config", zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"),
			zap.String("sheet", "maze_brush_foe_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeBrushFoeV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeBrushFoeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushFoeV8Config")
		logger.ErrorWF("invalid type. not *MazeBrushFoeV8Config", zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"),
			zap.String("sheet", "maze_brush_foe_v8"))
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
type gMazeBrushFoeV8Parser struct {
}

// New new config row data
func (*gMazeBrushFoeV8Parser) New() interface{} {
	return &MazeBrushFoeV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeBrushFoeV8Parser) Fields() []string {
	return gMazeBrushFoeV8Fields
}

// Parse parse raw data to row data
func (*gMazeBrushFoeV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeBrushFoeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushFoeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBrushFoeV8ConfigRow", zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"),
			zap.String("sheet", "maze_brush_foe_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeBrushFoeV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeBrushFoeV8ConfigRow",
			zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"),
			zap.String("sheet", "maze_brush_foe_v8"), zap.Int("need_count", len(gMazeBrushFoeV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 序号
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 序号 to int32 failed")
			logger.ErrorWF("parse field id 序号 to int32 failed.",
				zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 barries_id : 关卡id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field barries_id 关卡id to int32 failed")
			logger.ErrorWF("parse field barries_id 关卡id to int32 failed.",
				zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Barries_id = int32(tmp)
	}

	// parse column 2 brush_area_id : 刷怪区域id
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field brush_area_id 刷怪区域id to int32 failed")
			logger.ErrorWF("parse field brush_area_id 刷怪区域id to int32 failed.",
				zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Brush_area_id = int32(tmp)
	}

	// parse column 3 group_id : 波
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field group_id 波 to int32 failed")
			logger.ErrorWF("parse field group_id 波 to int32 failed.",
				zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Group_id = int32(tmp)
	}

	// parse column 4 in_group_order : 一波内刷新顺序
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field in_group_order 一波内刷新顺序 to int32 failed")
			logger.ErrorWF("parse field in_group_order 一波内刷新顺序 to int32 failed.",
				zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.In_group_order = int32(tmp)
	}

	// parse column 5 interval_time : 每波刷新时间（毫秒）
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field interval_time 每波刷新时间（毫秒） to int32 failed")
			logger.ErrorWF("parse field interval_time 每波刷新时间（毫秒） to int32 failed.",
				zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Interval_time = int32(tmp)
	}

	// parse column 6 display_wave_id : 页面展示的波次
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field display_wave_id 页面展示的波次 to int32 failed")
			logger.ErrorWF("parse field display_wave_id 页面展示的波次 to int32 failed.",
				zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Display_wave_id = int32(tmp)
	}

	// parse column 7 monsters_id : 怪物id列表
	if data[7] != "" {

		vals := strings.Split(data[7], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field monsters_id 怪物id列表 to []int32 failed")
				logger.ErrorWF("parse array field monsters_id 怪物id列表 to []int32 failed.",
					zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
					// zap.String("field_data",data[7]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Monsters_id = append(config.Monsters_id, int32(tmp))
		}
	}

	// parse column 8 wave_kongfu : 每波增加的武力值
	if data[8] != "" {
		tmp, err = strconv.ParseInt(data[8], 10, 64)
		if err != nil {
			err = errors.New("parse field wave_kongfu 每波增加的武力值 to int32 failed")
			logger.ErrorWF("parse field wave_kongfu 每波增加的武力值 to int32 failed.",
				zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
				zap.String("parse_data", data[8]),
				zap.Error(err))
			return
		}
		config.Wave_kongfu = int32(tmp)
	}

	// parse column 9 min_brushtime : 每波保底刷怪时间（毫秒）
	if data[9] != "" {
		tmp, err = strconv.ParseInt(data[9], 10, 64)
		if err != nil {
			err = errors.New("parse field min_brushtime 每波保底刷怪时间（毫秒） to int32 failed")
			logger.ErrorWF("parse field min_brushtime 每波保底刷怪时间（毫秒） to int32 failed.",
				zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
				zap.String("parse_data", data[9]),
				zap.Error(err))
			return
		}
		config.Min_brushtime = int32(tmp)
	}

	// parse column 10 iskill_all : 是否需要杀完上波怪再刷本波
	if data[10] != "" {
		tmp, err = strconv.ParseInt(data[10], 10, 64)
		if err != nil {
			err = errors.New("parse field iskill_all 是否需要杀完上波怪再刷本波 to int32 failed")
			logger.ErrorWF("parse field iskill_all 是否需要杀完上波怪再刷本波 to int32 failed.",
				zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
				zap.String("parse_data", data[10]),
				zap.Error(err))
			return
		}
		config.Iskill_all = int32(tmp)
	}

	// parse column 11 index : 战斗区域id
	if data[11] != "" {
		tmp, err = strconv.ParseInt(data[11], 10, 64)
		if err != nil {
			err = errors.New("parse field index 战斗区域id to int32 failed")
			logger.ErrorWF("parse field index 战斗区域id to int32 failed.",
				zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
				zap.String("parse_data", data[11]),
				zap.Error(err))
			return
		}
		config.Index = int32(tmp)
	}

	// parse column 12 Brushing_monsters_coordinate : 区域刷怪坐标列表（横百分段值,竖百分段值）
	if data[12] != "" {

		vals := strings.Split(data[12], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field Brushing_monsters_coordinate 区域刷怪坐标列表（横百分段值,竖百分段值） to []int32 failed")
				logger.ErrorWF("parse array field Brushing_monsters_coordinate 区域刷怪坐标列表（横百分段值,竖百分段值） to []int32 failed.",
					zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
					// zap.String("field_data",data[12]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Brushing_monsters_coordinate = append(config.Brushing_monsters_coordinate, int32(tmp))
		}
	}

	// parse column 13 Brushing_monsters_position : 区域刷怪点坐标（横百分段值,竖百分段值）,范围
	if data[13] != "" {

		vals := strings.Split(data[13], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field Brushing_monsters_position 区域刷怪点坐标（横百分段值,竖百分段值）,范围 to []int32 failed")
				logger.ErrorWF("parse array field Brushing_monsters_position 区域刷怪点坐标（横百分段值,竖百分段值）,范围 to []int32 failed.",
					zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
					// zap.String("field_data",data[13]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Brushing_monsters_position = append(config.Brushing_monsters_position, int32(tmp))
		}
	}

	// parse column 14 Brushing_monsters_relative_position : 相对刷怪点（角度0~360,距离,范围）
	if data[14] != "" {

		vals := strings.Split(data[14], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field Brushing_monsters_relative_position 相对刷怪点（角度0~360,距离,范围） to []int32 failed")
				logger.ErrorWF("parse array field Brushing_monsters_relative_position 相对刷怪点（角度0~360,距离,范围） to []int32 failed.",
					zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
					// zap.String("field_data",data[14]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Brushing_monsters_relative_position = append(config.Brushing_monsters_relative_position, int32(tmp))
		}
	}

	// parse column 15 front_group_order : 立即刷怪前置条件波次
	if data[15] != "" {

		vals := strings.Split(data[15], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field front_group_order 立即刷怪前置条件波次 to []int32 failed")
				logger.ErrorWF("parse array field front_group_order 立即刷怪前置条件波次 to []int32 failed.",
					zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
					// zap.String("field_data",data[15]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Front_group_order = append(config.Front_group_order, int32(tmp))
		}
	}

	// parse column 16 is_def_show : 是否需要默认显示
	if data[16] != "" {
		tmp, err = strconv.ParseInt(data[16], 10, 64)
		if err != nil {
			err = errors.New("parse field is_def_show 是否需要默认显示 to int32 failed")
			logger.ErrorWF("parse field is_def_show 是否需要默认显示 to int32 failed.",
				zap.String("xlsx", "maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx"), zap.String("sheet", "maze_brush_foe_v8"),
				zap.String("parse_data", data[16]),
				zap.Error(err))
			return
		}
		config.Is_def_show = int32(tmp)
	}
	return
}

var gMazeBrushFoeV8Fields = []string{
	"id",
	"barries_id",
	"brush_area_id",
	"group_id",
	"in_group_order",
	"interval_time",
	"display_wave_id",
	"monsters_id",
	"wave_kongfu",
	"min_brushtime",
	"iskill_all",
	"index",
	"Brushing_monsters_coordinate",
	"Brushing_monsters_position",
	"Brushing_monsters_relative_position",
	"front_group_order",
	"is_def_show",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeBrushFoeV8Parser{}
	loader := &gMazeBrushFoeV8Loader{}
	var data [][]string
	data, err = load("maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx", "maze_brush_foe_v8", gMazeBrushFoeV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx maze_brush_foe_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx maze_brush_foe_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeBrushFoeV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx maze_brush_foe_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeBrushFoeV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx maze_brush_foe_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx maze_brush_foe_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx maze_brush_foe_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_brush_foe_v8【迷宫-刷怪相关时间数量类型】.xlsx maze_brush_foe_v8 data success.")
	return
}
