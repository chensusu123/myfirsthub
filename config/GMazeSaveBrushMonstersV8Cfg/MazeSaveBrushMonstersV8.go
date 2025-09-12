package GMazeSaveBrushMonstersV8Cfg

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

// MazeSaveBrushMonstersV8ConfigRow from maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx maze_save_brush_monsters_v8
type MazeSaveBrushMonstersV8ConfigRow struct {
	Id       int32   `json:"id"`       // 刷怪id
	Monsters []int32 `json:"monsters"` // 刷怪id
	Time     []int32 `json:"time"`     // 刷怪时序(刷怪间隔时间，毫秒)
}

// MazeSaveBrushMonstersV8Config from maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx maze_save_brush_monsters_v8
type MazeSaveBrushMonstersV8Config struct {
	ConfigRows map[int32]*MazeSaveBrushMonstersV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeSaveBrushMonstersV8Config {
	ret := &MazeSaveBrushMonstersV8Config{ConfigRows: map[int32]*MazeSaveBrushMonstersV8ConfigRow{}}
	return ret
}

// GetMazeSaveBrushMonstersV8Config get one config by configId
func (c *MazeSaveBrushMonstersV8Config) GetMazeSaveBrushMonstersV8Config(configId int32) *MazeSaveBrushMonstersV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeSaveBrushMonstersV8Config) Get(configId int32) *MazeSaveBrushMonstersV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeSaveBrushMonstersV8Config get all config slice
func (c *MazeSaveBrushMonstersV8Config) GetAllMazeSaveBrushMonstersV8Config() (res []*MazeSaveBrushMonstersV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeSaveBrushMonstersV8Config) GetAll() (res []*MazeSaveBrushMonstersV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeSaveBrushMonstersV8Config

// GetMazeSaveBrushMonstersV8Config pkg func. get one config by configId
func GetMazeSaveBrushMonstersV8Config(configId int32) *MazeSaveBrushMonstersV8ConfigRow {
	return gConfigData.GetMazeSaveBrushMonstersV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeSaveBrushMonstersV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeSaveBrushMonstersV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_save_brush_monsters_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeSaveBrushMonstersV8Config pkg func. get all config slice
func GetAllMazeSaveBrushMonstersV8Config() []*MazeSaveBrushMonstersV8ConfigRow {
	return gConfigData.GetAllMazeSaveBrushMonstersV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeSaveBrushMonstersV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeSaveBrushMonstersV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeSaveBrushMonstersV8ConfigRow from maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx maze_save_brush_monsters_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeSaveBrushMonstersV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_save_brush_monsters_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_save_brush_monsters_v8.json",
		"maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx", "maze_save_brush_monsters_v8",
		&gMazeSaveBrushMonstersV8Parser{}, &gMazeSaveBrushMonstersV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeSaveBrushMonstersV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeSaveBrushMonstersV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeSaveBrushMonstersV8Config))(c)
		return true
	})
}

// RegisterMazeSaveBrushMonstersV8InitCallBack reg config update func (old func)
var RegisterMazeSaveBrushMonstersV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeSaveBrushMonstersV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeSaveBrushMonstersV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeSaveBrushMonstersV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeSaveBrushMonstersV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeSaveBrushMonstersV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeSaveBrushMonstersV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeSaveBrushMonstersV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeSaveBrushMonstersV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeSaveBrushMonstersV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeSaveBrushMonstersV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeSaveBrushMonstersV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeSaveBrushMonstersV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeSaveBrushMonstersV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveBrushMonstersV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSaveBrushMonstersV8ConfigRow", zap.String("xlsx", "maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx"),
			zap.String("sheet", "maze_save_brush_monsters_v8"))
		return
	}
	config, ok := container.(*MazeSaveBrushMonstersV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveBrushMonstersV8Config")
		logger.ErrorWF("invalid type. not *MazeSaveBrushMonstersV8Config", zap.String("xlsx", "maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx"),
			zap.String("sheet", "maze_save_brush_monsters_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeSaveBrushMonstersV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeSaveBrushMonstersV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveBrushMonstersV8Config")
		logger.ErrorWF("invalid type. not *MazeSaveBrushMonstersV8Config", zap.String("xlsx", "maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx"),
			zap.String("sheet", "maze_save_brush_monsters_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeSaveBrushMonstersV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeSaveBrushMonstersV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveBrushMonstersV8Config")
		logger.ErrorWF("invalid type. not *MazeSaveBrushMonstersV8Config", zap.String("xlsx", "maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx"),
			zap.String("sheet", "maze_save_brush_monsters_v8"))
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
type gMazeSaveBrushMonstersV8Parser struct {
}

// New new config row data
func (*gMazeSaveBrushMonstersV8Parser) New() interface{} {
	return &MazeSaveBrushMonstersV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeSaveBrushMonstersV8Parser) Fields() []string {
	return gMazeSaveBrushMonstersV8Fields
}

// Parse parse raw data to row data
func (*gMazeSaveBrushMonstersV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeSaveBrushMonstersV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveBrushMonstersV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSaveBrushMonstersV8ConfigRow", zap.String("xlsx", "maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx"),
			zap.String("sheet", "maze_save_brush_monsters_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeSaveBrushMonstersV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeSaveBrushMonstersV8ConfigRow",
			zap.String("xlsx", "maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx"),
			zap.String("sheet", "maze_save_brush_monsters_v8"), zap.Int("need_count", len(gMazeSaveBrushMonstersV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 刷怪id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 刷怪id to int32 failed")
			logger.ErrorWF("parse field id 刷怪id to int32 failed.",
				zap.String("xlsx", "maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx"), zap.String("sheet", "maze_save_brush_monsters_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 monsters : 刷怪id
	if data[1] != "" {

		vals := strings.Split(data[1], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field monsters 刷怪id to []int32 failed")
				logger.ErrorWF("parse array field monsters 刷怪id to []int32 failed.",
					zap.String("xlsx", "maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx"), zap.String("sheet", "maze_save_brush_monsters_v8"),
					// zap.String("field_data",data[1]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Monsters = append(config.Monsters, int32(tmp))
		}
	}

	// parse column 2 time : 刷怪时序(刷怪间隔时间，毫秒)
	if data[2] != "" {

		vals := strings.Split(data[2], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field time 刷怪时序(刷怪间隔时间，毫秒) to []int32 failed")
				logger.ErrorWF("parse array field time 刷怪时序(刷怪间隔时间，毫秒) to []int32 failed.",
					zap.String("xlsx", "maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx"), zap.String("sheet", "maze_save_brush_monsters_v8"),
					// zap.String("field_data",data[2]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Time = append(config.Time, int32(tmp))
		}
	}
	return
}

var gMazeSaveBrushMonstersV8Fields = []string{
	"id",
	"monsters",
	"time",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeSaveBrushMonstersV8Parser{}
	loader := &gMazeSaveBrushMonstersV8Loader{}
	var data [][]string
	data, err = load("maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx", "maze_save_brush_monsters_v8", gMazeSaveBrushMonstersV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx maze_save_brush_monsters_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx maze_save_brush_monsters_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeSaveBrushMonstersV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx maze_save_brush_monsters_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeSaveBrushMonstersV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx maze_save_brush_monsters_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx maze_save_brush_monsters_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx maze_save_brush_monsters_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_save_brush_monsters_v8【迷宫-救援敌人-刷怪】.xlsx maze_save_brush_monsters_v8 data success.")
	return
}
