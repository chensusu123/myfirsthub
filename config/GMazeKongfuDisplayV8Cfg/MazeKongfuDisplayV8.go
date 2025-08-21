package GMazeKongfuDisplayV8Cfg

import (
	"context"
	"errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"
)

// MazeKongfuDisplayV8ConfigRow from maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx maze_kongfu_display_v8
type MazeKongfuDisplayV8ConfigRow struct {
	Order         int32 `json:"order"`         // 序号
	Kongfu_min    int64 `json:"kongfu_min"`    // 玩家武力/主目标武力万分比，下限
	Kongfu_max    int64 `json:"kongfu_max"`    // 玩家武力/主目标武力万分比，上限
	Display_level int32 `json:"display_level"` // 刷怪表现等级
}

// MazeKongfuDisplayV8Config from maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx maze_kongfu_display_v8
type MazeKongfuDisplayV8Config struct {
	ConfigRows map[int32]*MazeKongfuDisplayV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeKongfuDisplayV8Config {
	ret := &MazeKongfuDisplayV8Config{ConfigRows: map[int32]*MazeKongfuDisplayV8ConfigRow{}}
	return ret
}

// GetMazeKongfuDisplayV8Config get one config by configId
func (c *MazeKongfuDisplayV8Config) GetMazeKongfuDisplayV8Config(configId int32) *MazeKongfuDisplayV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeKongfuDisplayV8Config) Get(configId int32) *MazeKongfuDisplayV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeKongfuDisplayV8Config get all config slice
func (c *MazeKongfuDisplayV8Config) GetAllMazeKongfuDisplayV8Config() (res []*MazeKongfuDisplayV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeKongfuDisplayV8Config) GetAll() (res []*MazeKongfuDisplayV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeKongfuDisplayV8Config

// GetMazeKongfuDisplayV8Config pkg func. get one config by configId
func GetMazeKongfuDisplayV8Config(configId int32) *MazeKongfuDisplayV8ConfigRow {
	return gConfigData.GetMazeKongfuDisplayV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeKongfuDisplayV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeKongfuDisplayV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_kongfu_display_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeKongfuDisplayV8Config pkg func. get all config slice
func GetAllMazeKongfuDisplayV8Config() []*MazeKongfuDisplayV8ConfigRow {
	return gConfigData.GetAllMazeKongfuDisplayV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeKongfuDisplayV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeKongfuDisplayV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeKongfuDisplayV8ConfigRow from maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx maze_kongfu_display_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeKongfuDisplayV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_kongfu_display_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_kongfu_display_v8.json",
		"maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx", "maze_kongfu_display_v8",
		&gMazeKongfuDisplayV8Parser{}, &gMazeKongfuDisplayV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeKongfuDisplayV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeKongfuDisplayV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeKongfuDisplayV8Config))(c)
		return true
	})
}

// RegisterMazeKongfuDisplayV8InitCallBack reg config update func (old func)
var RegisterMazeKongfuDisplayV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeKongfuDisplayV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeKongfuDisplayV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeKongfuDisplayV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeKongfuDisplayV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeKongfuDisplayV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeKongfuDisplayV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeKongfuDisplayV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeKongfuDisplayV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeKongfuDisplayV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeKongfuDisplayV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeKongfuDisplayV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeKongfuDisplayV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeKongfuDisplayV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuDisplayV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeKongfuDisplayV8ConfigRow", zap.String("xlsx", "maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx"),
			zap.String("sheet", "maze_kongfu_display_v8"))
		return
	}
	config, ok := container.(*MazeKongfuDisplayV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuDisplayV8Config")
		logger.ErrorWF("invalid type. not *MazeKongfuDisplayV8Config", zap.String("xlsx", "maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx"),
			zap.String("sheet", "maze_kongfu_display_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeKongfuDisplayV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeKongfuDisplayV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuDisplayV8Config")
		logger.ErrorWF("invalid type. not *MazeKongfuDisplayV8Config", zap.String("xlsx", "maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx"),
			zap.String("sheet", "maze_kongfu_display_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeKongfuDisplayV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeKongfuDisplayV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuDisplayV8Config")
		logger.ErrorWF("invalid type. not *MazeKongfuDisplayV8Config", zap.String("xlsx", "maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx"),
			zap.String("sheet", "maze_kongfu_display_v8"))
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
type gMazeKongfuDisplayV8Parser struct {
}

// New new config row data
func (*gMazeKongfuDisplayV8Parser) New() interface{} {
	return &MazeKongfuDisplayV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeKongfuDisplayV8Parser) Fields() []string {
	return gMazeKongfuDisplayV8Fields
}

// Parse parse raw data to row data
func (*gMazeKongfuDisplayV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeKongfuDisplayV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuDisplayV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeKongfuDisplayV8ConfigRow", zap.String("xlsx", "maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx"),
			zap.String("sheet", "maze_kongfu_display_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeKongfuDisplayV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeKongfuDisplayV8ConfigRow",
			zap.String("xlsx", "maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx"),
			zap.String("sheet", "maze_kongfu_display_v8"), zap.Int("need_count", len(gMazeKongfuDisplayV8Fields)),
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
				zap.String("xlsx", "maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx"), zap.String("sheet", "maze_kongfu_display_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 kongfu_min : 玩家武力/主目标武力万分比，下限
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field kongfu_min 玩家武力/主目标武力万分比，下限 to int64 failed")
			logger.ErrorWF("parse field kongfu_min 玩家武力/主目标武力万分比，下限 to int64 failed.",
				zap.String("xlsx", "maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx"), zap.String("sheet", "maze_kongfu_display_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Kongfu_min = int64(tmp)
	}

	// parse column 2 kongfu_max : 玩家武力/主目标武力万分比，上限
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field kongfu_max 玩家武力/主目标武力万分比，上限 to int64 failed")
			logger.ErrorWF("parse field kongfu_max 玩家武力/主目标武力万分比，上限 to int64 failed.",
				zap.String("xlsx", "maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx"), zap.String("sheet", "maze_kongfu_display_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Kongfu_max = int64(tmp)
	}

	// parse column 3 display_level : 刷怪表现等级
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field display_level 刷怪表现等级 to int32 failed")
			logger.ErrorWF("parse field display_level 刷怪表现等级 to int32 failed.",
				zap.String("xlsx", "maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx"), zap.String("sheet", "maze_kongfu_display_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Display_level = int32(tmp)
	}
	return
}

var gMazeKongfuDisplayV8Fields = []string{
	"order",
	"kongfu_min",
	"kongfu_max",
	"display_level",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeKongfuDisplayV8Parser{}
	loader := &gMazeKongfuDisplayV8Loader{}
	var data [][]string
	data, err = load("maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx", "maze_kongfu_display_v8", gMazeKongfuDisplayV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx maze_kongfu_display_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx maze_kongfu_display_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeKongfuDisplayV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx maze_kongfu_display_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeKongfuDisplayV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx maze_kongfu_display_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx maze_kongfu_display_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx maze_kongfu_display_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_kongfu_display_v8【迷宫-武力值对应表现等级】.xlsx maze_kongfu_display_v8 data success.")
	return
}
