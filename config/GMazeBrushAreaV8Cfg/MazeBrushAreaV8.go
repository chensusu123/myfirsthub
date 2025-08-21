package GMazeBrushAreaV8Cfg

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

// MazeBrushAreaV8ConfigRow from maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx maze_brush_area_v8
type MazeBrushAreaV8ConfigRow struct {
	Order         int32           `json:"order"`         // 序号（区域id）
	Area_id       int32           `json:"area_id"`       // 区域id
	Barries       int32           `json:"barries"`       // 区域所属关卡id
	First_foe     int32           `json:"first_foe"`     // 初始怪物数量
	Min_foe       int32           `json:"min_foe"`       // 最少怪物数量（少于这个数就刷新）
	Max_foe       int32           `json:"max_foe"`       // 区域内最大怪物数量(>=这个数就不刷新
	Brush_cd      int32           `json:"brush_cd"`      // 怪物刷新cd（毫秒）
	Foe_pool      map[int32]int32 `json:"foe_pool"`      // 怪物id:总数量
	Brusharea_min int32           `json:"brusharea_min"` // 动态刷怪区域（以人为原点的最小半径）
	Brusharea_max int32           `json:"brusharea_max"` // 动态刷怪区域（以人为原点的最大半径）
}

// MazeBrushAreaV8Config from maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx maze_brush_area_v8
type MazeBrushAreaV8Config struct {
	ConfigRows map[int32]*MazeBrushAreaV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeBrushAreaV8Config {
	ret := &MazeBrushAreaV8Config{ConfigRows: map[int32]*MazeBrushAreaV8ConfigRow{}}
	return ret
}

// GetMazeBrushAreaV8Config get one config by configId
func (c *MazeBrushAreaV8Config) GetMazeBrushAreaV8Config(configId int32) *MazeBrushAreaV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeBrushAreaV8Config) Get(configId int32) *MazeBrushAreaV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeBrushAreaV8Config get all config slice
func (c *MazeBrushAreaV8Config) GetAllMazeBrushAreaV8Config() (res []*MazeBrushAreaV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeBrushAreaV8Config) GetAll() (res []*MazeBrushAreaV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeBrushAreaV8Config

// GetMazeBrushAreaV8Config pkg func. get one config by configId
func GetMazeBrushAreaV8Config(configId int32) *MazeBrushAreaV8ConfigRow {
	return gConfigData.GetMazeBrushAreaV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeBrushAreaV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeBrushAreaV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_brush_area_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeBrushAreaV8Config pkg func. get all config slice
func GetAllMazeBrushAreaV8Config() []*MazeBrushAreaV8ConfigRow {
	return gConfigData.GetAllMazeBrushAreaV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeBrushAreaV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeBrushAreaV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeBrushAreaV8ConfigRow from maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx maze_brush_area_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeBrushAreaV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_brush_area_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_brush_area_v8.json",
		"maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx", "maze_brush_area_v8",
		&gMazeBrushAreaV8Parser{}, &gMazeBrushAreaV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeBrushAreaV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeBrushAreaV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeBrushAreaV8Config))(c)
		return true
	})
}

// RegisterMazeBrushAreaV8InitCallBack reg config update func (old func)
var RegisterMazeBrushAreaV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeBrushAreaV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeBrushAreaV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeBrushAreaV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeBrushAreaV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeBrushAreaV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeBrushAreaV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeBrushAreaV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeBrushAreaV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeBrushAreaV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeBrushAreaV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeBrushAreaV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeBrushAreaV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeBrushAreaV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushAreaV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBrushAreaV8ConfigRow", zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"),
			zap.String("sheet", "maze_brush_area_v8"))
		return
	}
	config, ok := container.(*MazeBrushAreaV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushAreaV8Config")
		logger.ErrorWF("invalid type. not *MazeBrushAreaV8Config", zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"),
			zap.String("sheet", "maze_brush_area_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeBrushAreaV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeBrushAreaV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushAreaV8Config")
		logger.ErrorWF("invalid type. not *MazeBrushAreaV8Config", zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"),
			zap.String("sheet", "maze_brush_area_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeBrushAreaV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeBrushAreaV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushAreaV8Config")
		logger.ErrorWF("invalid type. not *MazeBrushAreaV8Config", zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"),
			zap.String("sheet", "maze_brush_area_v8"))
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
type gMazeBrushAreaV8Parser struct {
}

// New new config row data
func (*gMazeBrushAreaV8Parser) New() interface{} {
	return &MazeBrushAreaV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeBrushAreaV8Parser) Fields() []string {
	return gMazeBrushAreaV8Fields
}

// Parse parse raw data to row data
func (*gMazeBrushAreaV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeBrushAreaV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushAreaV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBrushAreaV8ConfigRow", zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"),
			zap.String("sheet", "maze_brush_area_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeBrushAreaV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeBrushAreaV8ConfigRow",
			zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"),
			zap.String("sheet", "maze_brush_area_v8"), zap.Int("need_count", len(gMazeBrushAreaV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 序号（区域id）
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 序号（区域id） to int32 failed")
			logger.ErrorWF("parse field order 序号（区域id） to int32 failed.",
				zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"), zap.String("sheet", "maze_brush_area_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 area_id : 区域id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field area_id 区域id to int32 failed")
			logger.ErrorWF("parse field area_id 区域id to int32 failed.",
				zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"), zap.String("sheet", "maze_brush_area_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Area_id = int32(tmp)
	}

	// parse column 2 barries : 区域所属关卡id
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field barries 区域所属关卡id to int32 failed")
			logger.ErrorWF("parse field barries 区域所属关卡id to int32 failed.",
				zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"), zap.String("sheet", "maze_brush_area_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Barries = int32(tmp)
	}

	// parse column 3 first_foe : 初始怪物数量
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field first_foe 初始怪物数量 to int32 failed")
			logger.ErrorWF("parse field first_foe 初始怪物数量 to int32 failed.",
				zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"), zap.String("sheet", "maze_brush_area_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.First_foe = int32(tmp)
	}

	// parse column 4 min_foe : 最少怪物数量（少于这个数就刷新）
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field min_foe 最少怪物数量（少于这个数就刷新） to int32 failed")
			logger.ErrorWF("parse field min_foe 最少怪物数量（少于这个数就刷新） to int32 failed.",
				zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"), zap.String("sheet", "maze_brush_area_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Min_foe = int32(tmp)
	}

	// parse column 5 max_foe : 区域内最大怪物数量(>=这个数就不刷新
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field max_foe 区域内最大怪物数量(>=这个数就不刷新 to int32 failed")
			logger.ErrorWF("parse field max_foe 区域内最大怪物数量(>=这个数就不刷新 to int32 failed.",
				zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"), zap.String("sheet", "maze_brush_area_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Max_foe = int32(tmp)
	}

	// parse column 6 brush_cd : 怪物刷新cd（毫秒）
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field brush_cd 怪物刷新cd（毫秒） to int32 failed")
			logger.ErrorWF("parse field brush_cd 怪物刷新cd（毫秒） to int32 failed.",
				zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"), zap.String("sheet", "maze_brush_area_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Brush_cd = int32(tmp)
	}

	// parse column 7 foe_pool : 怪物id:总数量
	if data[7] != "" {

		config.Foe_pool = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[7], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field foe_pool 怪物id:总数量 to key int32 failed")
				logger.ErrorWF("parse map field foe_pool 怪物id:总数量 to key int32 failed.",
					zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"), zap.String("sheet", "maze_brush_area_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field foe_pool 怪物id:总数量 to value int32 failed")
				logger.ErrorWF("parse map field foe_pool 怪物id:总数量 to value int32 failed.",
					zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"), zap.String("sheet", "maze_brush_area_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Foe_pool[key] = value
		}
	}

	// parse column 8 brusharea_min : 动态刷怪区域（以人为原点的最小半径）
	if data[8] != "" {
		tmp, err = strconv.ParseInt(data[8], 10, 64)
		if err != nil {
			err = errors.New("parse field brusharea_min 动态刷怪区域（以人为原点的最小半径） to int32 failed")
			logger.ErrorWF("parse field brusharea_min 动态刷怪区域（以人为原点的最小半径） to int32 failed.",
				zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"), zap.String("sheet", "maze_brush_area_v8"),
				zap.String("parse_data", data[8]),
				zap.Error(err))
			return
		}
		config.Brusharea_min = int32(tmp)
	}

	// parse column 9 brusharea_max : 动态刷怪区域（以人为原点的最大半径）
	if data[9] != "" {
		tmp, err = strconv.ParseInt(data[9], 10, 64)
		if err != nil {
			err = errors.New("parse field brusharea_max 动态刷怪区域（以人为原点的最大半径） to int32 failed")
			logger.ErrorWF("parse field brusharea_max 动态刷怪区域（以人为原点的最大半径） to int32 failed.",
				zap.String("xlsx", "maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx"), zap.String("sheet", "maze_brush_area_v8"),
				zap.String("parse_data", data[9]),
				zap.Error(err))
			return
		}
		config.Brusharea_max = int32(tmp)
	}
	return
}

var gMazeBrushAreaV8Fields = []string{
	"order",
	"area_id",
	"barries",
	"first_foe",
	"min_foe",
	"max_foe",
	"brush_cd",
	"foe_pool",
	"brusharea_min",
	"brusharea_max",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeBrushAreaV8Parser{}
	loader := &gMazeBrushAreaV8Loader{}
	var data [][]string
	data, err = load("maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx", "maze_brush_area_v8", gMazeBrushAreaV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx maze_brush_area_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx maze_brush_area_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeBrushAreaV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx maze_brush_area_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeBrushAreaV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx maze_brush_area_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx maze_brush_area_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx maze_brush_area_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_brush_area_v8【迷宫-区域怪物刷新】.xlsx maze_brush_area_v8 data success.")
	return
}
