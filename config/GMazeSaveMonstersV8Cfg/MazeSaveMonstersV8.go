package GMazeSaveMonstersV8Cfg

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

// MazeSaveMonstersV8ConfigRow from maze_save_monsters_v8【迷宫-救援敌人】.xlsx maze_save_monsters_v8
type MazeSaveMonstersV8ConfigRow struct {
	Id       int32  `json:"id"`       // 怪物id
	Model_id int32  `json:"model_id"` // 资源组id
	Speed    string `json:"speed"`    // 移速
	Hp       int32  `json:"hp"`       // 生命值
}

// MazeSaveMonstersV8Config from maze_save_monsters_v8【迷宫-救援敌人】.xlsx maze_save_monsters_v8
type MazeSaveMonstersV8Config struct {
	ConfigRows map[int32]*MazeSaveMonstersV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeSaveMonstersV8Config {
	ret := &MazeSaveMonstersV8Config{ConfigRows: map[int32]*MazeSaveMonstersV8ConfigRow{}}
	return ret
}

// GetMazeSaveMonstersV8Config get one config by configId
func (c *MazeSaveMonstersV8Config) GetMazeSaveMonstersV8Config(configId int32) *MazeSaveMonstersV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeSaveMonstersV8Config) Get(configId int32) *MazeSaveMonstersV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeSaveMonstersV8Config get all config slice
func (c *MazeSaveMonstersV8Config) GetAllMazeSaveMonstersV8Config() (res []*MazeSaveMonstersV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeSaveMonstersV8Config) GetAll() (res []*MazeSaveMonstersV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeSaveMonstersV8Config

// GetMazeSaveMonstersV8Config pkg func. get one config by configId
func GetMazeSaveMonstersV8Config(configId int32) *MazeSaveMonstersV8ConfigRow {
	return gConfigData.GetMazeSaveMonstersV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeSaveMonstersV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeSaveMonstersV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_save_monsters_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeSaveMonstersV8Config pkg func. get all config slice
func GetAllMazeSaveMonstersV8Config() []*MazeSaveMonstersV8ConfigRow {
	return gConfigData.GetAllMazeSaveMonstersV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeSaveMonstersV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeSaveMonstersV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeSaveMonstersV8ConfigRow from maze_save_monsters_v8【迷宫-救援敌人】.xlsx maze_save_monsters_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeSaveMonstersV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_save_monsters_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_save_monsters_v8.json",
		"maze_save_monsters_v8【迷宫-救援敌人】.xlsx", "maze_save_monsters_v8",
		&gMazeSaveMonstersV8Parser{}, &gMazeSaveMonstersV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeSaveMonstersV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeSaveMonstersV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeSaveMonstersV8Config))(c)
		return true
	})
}

// RegisterMazeSaveMonstersV8InitCallBack reg config update func (old func)
var RegisterMazeSaveMonstersV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeSaveMonstersV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeSaveMonstersV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeSaveMonstersV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeSaveMonstersV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeSaveMonstersV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeSaveMonstersV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeSaveMonstersV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeSaveMonstersV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeSaveMonstersV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeSaveMonstersV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeSaveMonstersV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeSaveMonstersV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeSaveMonstersV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveMonstersV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSaveMonstersV8ConfigRow", zap.String("xlsx", "maze_save_monsters_v8【迷宫-救援敌人】.xlsx"),
			zap.String("sheet", "maze_save_monsters_v8"))
		return
	}
	config, ok := container.(*MazeSaveMonstersV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveMonstersV8Config")
		logger.ErrorWF("invalid type. not *MazeSaveMonstersV8Config", zap.String("xlsx", "maze_save_monsters_v8【迷宫-救援敌人】.xlsx"),
			zap.String("sheet", "maze_save_monsters_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeSaveMonstersV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeSaveMonstersV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveMonstersV8Config")
		logger.ErrorWF("invalid type. not *MazeSaveMonstersV8Config", zap.String("xlsx", "maze_save_monsters_v8【迷宫-救援敌人】.xlsx"),
			zap.String("sheet", "maze_save_monsters_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeSaveMonstersV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeSaveMonstersV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveMonstersV8Config")
		logger.ErrorWF("invalid type. not *MazeSaveMonstersV8Config", zap.String("xlsx", "maze_save_monsters_v8【迷宫-救援敌人】.xlsx"),
			zap.String("sheet", "maze_save_monsters_v8"))
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
type gMazeSaveMonstersV8Parser struct {
}

// New new config row data
func (*gMazeSaveMonstersV8Parser) New() interface{} {
	return &MazeSaveMonstersV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeSaveMonstersV8Parser) Fields() []string {
	return gMazeSaveMonstersV8Fields
}

// Parse parse raw data to row data
func (*gMazeSaveMonstersV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeSaveMonstersV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveMonstersV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSaveMonstersV8ConfigRow", zap.String("xlsx", "maze_save_monsters_v8【迷宫-救援敌人】.xlsx"),
			zap.String("sheet", "maze_save_monsters_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeSaveMonstersV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeSaveMonstersV8ConfigRow",
			zap.String("xlsx", "maze_save_monsters_v8【迷宫-救援敌人】.xlsx"),
			zap.String("sheet", "maze_save_monsters_v8"), zap.Int("need_count", len(gMazeSaveMonstersV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 怪物id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 怪物id to int32 failed")
			logger.ErrorWF("parse field id 怪物id to int32 failed.",
				zap.String("xlsx", "maze_save_monsters_v8【迷宫-救援敌人】.xlsx"), zap.String("sheet", "maze_save_monsters_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 model_id : 资源组id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field model_id 资源组id to int32 failed")
			logger.ErrorWF("parse field model_id 资源组id to int32 failed.",
				zap.String("xlsx", "maze_save_monsters_v8【迷宫-救援敌人】.xlsx"), zap.String("sheet", "maze_save_monsters_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Model_id = int32(tmp)
	}

	// parse column 2 speed : 移速
	if data[2] != "" {
		config.Speed = data[2]
	}

	// parse column 3 hp : 生命值
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field hp 生命值 to int32 failed")
			logger.ErrorWF("parse field hp 生命值 to int32 failed.",
				zap.String("xlsx", "maze_save_monsters_v8【迷宫-救援敌人】.xlsx"), zap.String("sheet", "maze_save_monsters_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Hp = int32(tmp)
	}
	return
}

var gMazeSaveMonstersV8Fields = []string{
	"id",
	"model_id",
	"speed",
	"hp",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeSaveMonstersV8Parser{}
	loader := &gMazeSaveMonstersV8Loader{}
	var data [][]string
	data, err = load("maze_save_monsters_v8【迷宫-救援敌人】.xlsx", "maze_save_monsters_v8", gMazeSaveMonstersV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_save_monsters_v8【迷宫-救援敌人】.xlsx maze_save_monsters_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_save_monsters_v8【迷宫-救援敌人】.xlsx maze_save_monsters_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeSaveMonstersV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_save_monsters_v8【迷宫-救援敌人】.xlsx maze_save_monsters_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeSaveMonstersV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_save_monsters_v8【迷宫-救援敌人】.xlsx maze_save_monsters_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_save_monsters_v8【迷宫-救援敌人】.xlsx maze_save_monsters_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_save_monsters_v8【迷宫-救援敌人】.xlsx maze_save_monsters_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_save_monsters_v8【迷宫-救援敌人】.xlsx maze_save_monsters_v8 data success.")
	return
}
