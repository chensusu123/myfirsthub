package GMazeBarriesOnHookV8Cfg

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

// MazeBarriesOnHookV8ConfigRow from maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx maze_barries_on_hook_v8
type MazeBarriesOnHookV8ConfigRow struct {
	Order            int32           `json:"order"`            // 关卡id
	Cycle_award      map[int32]int64 `json:"cycle_award"`      // 每周期产出挂机奖励万分比
	Cycle_time       int32           `json:"cycle_time"`       // 产出周期（秒）
	Maxlimit_cycle   int32           `json:"maxlimit_cycle"`   // 不领取可储存的奖励周期数
	Can_receive_time int32           `json:"can_receive_time"` // 允许领取奖励的时间（秒）
	Cycle_award_2    map[int32]int64 `json:"cycle_award_2"`    // 每周期产出挂机稀有奖励万分比
}

// MazeBarriesOnHookV8Config from maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx maze_barries_on_hook_v8
type MazeBarriesOnHookV8Config struct {
	ConfigRows map[int32]*MazeBarriesOnHookV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeBarriesOnHookV8Config {
	ret := &MazeBarriesOnHookV8Config{ConfigRows: map[int32]*MazeBarriesOnHookV8ConfigRow{}}
	return ret
}

// GetMazeBarriesOnHookV8Config get one config by configId
func (c *MazeBarriesOnHookV8Config) GetMazeBarriesOnHookV8Config(configId int32) *MazeBarriesOnHookV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeBarriesOnHookV8Config) Get(configId int32) *MazeBarriesOnHookV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeBarriesOnHookV8Config get all config slice
func (c *MazeBarriesOnHookV8Config) GetAllMazeBarriesOnHookV8Config() (res []*MazeBarriesOnHookV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeBarriesOnHookV8Config) GetAll() (res []*MazeBarriesOnHookV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeBarriesOnHookV8Config

// GetMazeBarriesOnHookV8Config pkg func. get one config by configId
func GetMazeBarriesOnHookV8Config(configId int32) *MazeBarriesOnHookV8ConfigRow {
	return gConfigData.GetMazeBarriesOnHookV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeBarriesOnHookV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeBarriesOnHookV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_barries_on_hook_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeBarriesOnHookV8Config pkg func. get all config slice
func GetAllMazeBarriesOnHookV8Config() []*MazeBarriesOnHookV8ConfigRow {
	return gConfigData.GetAllMazeBarriesOnHookV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeBarriesOnHookV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeBarriesOnHookV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeBarriesOnHookV8ConfigRow from maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx maze_barries_on_hook_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeBarriesOnHookV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_barries_on_hook_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_barries_on_hook_v8.json",
		"maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx", "maze_barries_on_hook_v8",
		&gMazeBarriesOnHookV8Parser{}, &gMazeBarriesOnHookV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeBarriesOnHookV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeBarriesOnHookV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeBarriesOnHookV8Config))(c)
		return true
	})
}

// RegisterMazeBarriesOnHookV8InitCallBack reg config update func (old func)
var RegisterMazeBarriesOnHookV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeBarriesOnHookV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeBarriesOnHookV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeBarriesOnHookV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeBarriesOnHookV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeBarriesOnHookV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeBarriesOnHookV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeBarriesOnHookV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeBarriesOnHookV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeBarriesOnHookV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeBarriesOnHookV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeBarriesOnHookV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeBarriesOnHookV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeBarriesOnHookV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesOnHookV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBarriesOnHookV8ConfigRow", zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"),
			zap.String("sheet", "maze_barries_on_hook_v8"))
		return
	}
	config, ok := container.(*MazeBarriesOnHookV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesOnHookV8Config")
		logger.ErrorWF("invalid type. not *MazeBarriesOnHookV8Config", zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"),
			zap.String("sheet", "maze_barries_on_hook_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeBarriesOnHookV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeBarriesOnHookV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesOnHookV8Config")
		logger.ErrorWF("invalid type. not *MazeBarriesOnHookV8Config", zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"),
			zap.String("sheet", "maze_barries_on_hook_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeBarriesOnHookV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeBarriesOnHookV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesOnHookV8Config")
		logger.ErrorWF("invalid type. not *MazeBarriesOnHookV8Config", zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"),
			zap.String("sheet", "maze_barries_on_hook_v8"))
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
type gMazeBarriesOnHookV8Parser struct {
}

// New new config row data
func (*gMazeBarriesOnHookV8Parser) New() interface{} {
	return &MazeBarriesOnHookV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeBarriesOnHookV8Parser) Fields() []string {
	return gMazeBarriesOnHookV8Fields
}

// Parse parse raw data to row data
func (*gMazeBarriesOnHookV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeBarriesOnHookV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesOnHookV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBarriesOnHookV8ConfigRow", zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"),
			zap.String("sheet", "maze_barries_on_hook_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeBarriesOnHookV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeBarriesOnHookV8ConfigRow",
			zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"),
			zap.String("sheet", "maze_barries_on_hook_v8"), zap.Int("need_count", len(gMazeBarriesOnHookV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 关卡id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 关卡id to int32 failed")
			logger.ErrorWF("parse field order 关卡id to int32 failed.",
				zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"), zap.String("sheet", "maze_barries_on_hook_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 cycle_award : 每周期产出挂机奖励万分比
	if data[1] != "" {

		config.Cycle_award = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[1], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field cycle_award 每周期产出挂机奖励万分比 to key int32 failed")
				logger.ErrorWF("parse map field cycle_award 每周期产出挂机奖励万分比 to key int32 failed.",
					zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"), zap.String("sheet", "maze_barries_on_hook_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field cycle_award 每周期产出挂机奖励万分比 to value int64 failed")
				logger.ErrorWF("parse map field cycle_award 每周期产出挂机奖励万分比 to value int64 failed.",
					zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"), zap.String("sheet", "maze_barries_on_hook_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Cycle_award[key] = value
		}
	}

	// parse column 2 cycle_time : 产出周期（秒）
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field cycle_time 产出周期（秒） to int32 failed")
			logger.ErrorWF("parse field cycle_time 产出周期（秒） to int32 failed.",
				zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"), zap.String("sheet", "maze_barries_on_hook_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Cycle_time = int32(tmp)
	}

	// parse column 3 maxlimit_cycle : 不领取可储存的奖励周期数
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field maxlimit_cycle 不领取可储存的奖励周期数 to int32 failed")
			logger.ErrorWF("parse field maxlimit_cycle 不领取可储存的奖励周期数 to int32 failed.",
				zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"), zap.String("sheet", "maze_barries_on_hook_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Maxlimit_cycle = int32(tmp)
	}

	// parse column 4 can_receive_time : 允许领取奖励的时间（秒）
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field can_receive_time 允许领取奖励的时间（秒） to int32 failed")
			logger.ErrorWF("parse field can_receive_time 允许领取奖励的时间（秒） to int32 failed.",
				zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"), zap.String("sheet", "maze_barries_on_hook_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Can_receive_time = int32(tmp)
	}

	// parse column 5 cycle_award_2 : 每周期产出挂机稀有奖励万分比
	if data[5] != "" {

		config.Cycle_award_2 = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[5], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field cycle_award_2 每周期产出挂机稀有奖励万分比 to key int32 failed")
				logger.ErrorWF("parse map field cycle_award_2 每周期产出挂机稀有奖励万分比 to key int32 failed.",
					zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"), zap.String("sheet", "maze_barries_on_hook_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field cycle_award_2 每周期产出挂机稀有奖励万分比 to value int64 failed")
				logger.ErrorWF("parse map field cycle_award_2 每周期产出挂机稀有奖励万分比 to value int64 failed.",
					zap.String("xlsx", "maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx"), zap.String("sheet", "maze_barries_on_hook_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Cycle_award_2[key] = value
		}
	}
	return
}

var gMazeBarriesOnHookV8Fields = []string{
	"order",
	"cycle_award",
	"cycle_time",
	"maxlimit_cycle",
	"can_receive_time",
	"cycle_award_2",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeBarriesOnHookV8Parser{}
	loader := &gMazeBarriesOnHookV8Loader{}
	var data [][]string
	data, err = load("maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx", "maze_barries_on_hook_v8", gMazeBarriesOnHookV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx maze_barries_on_hook_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx maze_barries_on_hook_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeBarriesOnHookV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx maze_barries_on_hook_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeBarriesOnHookV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx maze_barries_on_hook_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx maze_barries_on_hook_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx maze_barries_on_hook_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_barries_on_hook_v8【迷宫-关卡-挂机奖励】.xlsx maze_barries_on_hook_v8 data success.")
	return
}
