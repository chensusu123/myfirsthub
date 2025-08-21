package GMazeRebornCostV8Cfg

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

// MazeRebornCostV8ConfigRow from maze_reborn_cost_v8【迷宫-复活消耗】.xlsx maze_reborn_cost_v8
type MazeRebornCostV8ConfigRow struct {
	Order      int32           `json:"order"`      // 序号
	Reborn_min int32           `json:"reborn_min"` // 复活次数
	Reborn_max int32           `json:"reborn_max"` // 复活次数
	Cost       map[int32]int64 `json:"cost"`       // 复活消耗
}

// MazeRebornCostV8Config from maze_reborn_cost_v8【迷宫-复活消耗】.xlsx maze_reborn_cost_v8
type MazeRebornCostV8Config struct {
	ConfigRows map[int32]*MazeRebornCostV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeRebornCostV8Config {
	ret := &MazeRebornCostV8Config{ConfigRows: map[int32]*MazeRebornCostV8ConfigRow{}}
	return ret
}

// GetMazeRebornCostV8Config get one config by configId
func (c *MazeRebornCostV8Config) GetMazeRebornCostV8Config(configId int32) *MazeRebornCostV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeRebornCostV8Config) Get(configId int32) *MazeRebornCostV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeRebornCostV8Config get all config slice
func (c *MazeRebornCostV8Config) GetAllMazeRebornCostV8Config() (res []*MazeRebornCostV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeRebornCostV8Config) GetAll() (res []*MazeRebornCostV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeRebornCostV8Config

// GetMazeRebornCostV8Config pkg func. get one config by configId
func GetMazeRebornCostV8Config(configId int32) *MazeRebornCostV8ConfigRow {
	return gConfigData.GetMazeRebornCostV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeRebornCostV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeRebornCostV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_reborn_cost_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeRebornCostV8Config pkg func. get all config slice
func GetAllMazeRebornCostV8Config() []*MazeRebornCostV8ConfigRow {
	return gConfigData.GetAllMazeRebornCostV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeRebornCostV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeRebornCostV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeRebornCostV8ConfigRow from maze_reborn_cost_v8【迷宫-复活消耗】.xlsx maze_reborn_cost_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeRebornCostV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_reborn_cost_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_reborn_cost_v8.json",
		"maze_reborn_cost_v8【迷宫-复活消耗】.xlsx", "maze_reborn_cost_v8",
		&gMazeRebornCostV8Parser{}, &gMazeRebornCostV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeRebornCostV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeRebornCostV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeRebornCostV8Config))(c)
		return true
	})
}

// RegisterMazeRebornCostV8InitCallBack reg config update func (old func)
var RegisterMazeRebornCostV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeRebornCostV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeRebornCostV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeRebornCostV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeRebornCostV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeRebornCostV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeRebornCostV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeRebornCostV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeRebornCostV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeRebornCostV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeRebornCostV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeRebornCostV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeRebornCostV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeRebornCostV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeRebornCostV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeRebornCostV8ConfigRow", zap.String("xlsx", "maze_reborn_cost_v8【迷宫-复活消耗】.xlsx"),
			zap.String("sheet", "maze_reborn_cost_v8"))
		return
	}
	config, ok := container.(*MazeRebornCostV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeRebornCostV8Config")
		logger.ErrorWF("invalid type. not *MazeRebornCostV8Config", zap.String("xlsx", "maze_reborn_cost_v8【迷宫-复活消耗】.xlsx"),
			zap.String("sheet", "maze_reborn_cost_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeRebornCostV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeRebornCostV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeRebornCostV8Config")
		logger.ErrorWF("invalid type. not *MazeRebornCostV8Config", zap.String("xlsx", "maze_reborn_cost_v8【迷宫-复活消耗】.xlsx"),
			zap.String("sheet", "maze_reborn_cost_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeRebornCostV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeRebornCostV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeRebornCostV8Config")
		logger.ErrorWF("invalid type. not *MazeRebornCostV8Config", zap.String("xlsx", "maze_reborn_cost_v8【迷宫-复活消耗】.xlsx"),
			zap.String("sheet", "maze_reborn_cost_v8"))
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
type gMazeRebornCostV8Parser struct {
}

// New new config row data
func (*gMazeRebornCostV8Parser) New() interface{} {
	return &MazeRebornCostV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeRebornCostV8Parser) Fields() []string {
	return gMazeRebornCostV8Fields
}

// Parse parse raw data to row data
func (*gMazeRebornCostV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeRebornCostV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeRebornCostV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeRebornCostV8ConfigRow", zap.String("xlsx", "maze_reborn_cost_v8【迷宫-复活消耗】.xlsx"),
			zap.String("sheet", "maze_reborn_cost_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeRebornCostV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeRebornCostV8ConfigRow",
			zap.String("xlsx", "maze_reborn_cost_v8【迷宫-复活消耗】.xlsx"),
			zap.String("sheet", "maze_reborn_cost_v8"), zap.Int("need_count", len(gMazeRebornCostV8Fields)),
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
				zap.String("xlsx", "maze_reborn_cost_v8【迷宫-复活消耗】.xlsx"), zap.String("sheet", "maze_reborn_cost_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 reborn_min : 复活次数
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field reborn_min 复活次数 to int32 failed")
			logger.ErrorWF("parse field reborn_min 复活次数 to int32 failed.",
				zap.String("xlsx", "maze_reborn_cost_v8【迷宫-复活消耗】.xlsx"), zap.String("sheet", "maze_reborn_cost_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Reborn_min = int32(tmp)
	}

	// parse column 2 reborn_max : 复活次数
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field reborn_max 复活次数 to int32 failed")
			logger.ErrorWF("parse field reborn_max 复活次数 to int32 failed.",
				zap.String("xlsx", "maze_reborn_cost_v8【迷宫-复活消耗】.xlsx"), zap.String("sheet", "maze_reborn_cost_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Reborn_max = int32(tmp)
	}

	// parse column 3 cost : 复活消耗
	if data[3] != "" {

		config.Cost = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field cost 复活消耗 to key int32 failed")
				logger.ErrorWF("parse map field cost 复活消耗 to key int32 failed.",
					zap.String("xlsx", "maze_reborn_cost_v8【迷宫-复活消耗】.xlsx"), zap.String("sheet", "maze_reborn_cost_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field cost 复活消耗 to value int64 failed")
				logger.ErrorWF("parse map field cost 复活消耗 to value int64 failed.",
					zap.String("xlsx", "maze_reborn_cost_v8【迷宫-复活消耗】.xlsx"), zap.String("sheet", "maze_reborn_cost_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Cost[key] = value
		}
	}
	return
}

var gMazeRebornCostV8Fields = []string{
	"order",
	"reborn_min",
	"reborn_max",
	"cost",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeRebornCostV8Parser{}
	loader := &gMazeRebornCostV8Loader{}
	var data [][]string
	data, err = load("maze_reborn_cost_v8【迷宫-复活消耗】.xlsx", "maze_reborn_cost_v8", gMazeRebornCostV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_reborn_cost_v8【迷宫-复活消耗】.xlsx maze_reborn_cost_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_reborn_cost_v8【迷宫-复活消耗】.xlsx maze_reborn_cost_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeRebornCostV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_reborn_cost_v8【迷宫-复活消耗】.xlsx maze_reborn_cost_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeRebornCostV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_reborn_cost_v8【迷宫-复活消耗】.xlsx maze_reborn_cost_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_reborn_cost_v8【迷宫-复活消耗】.xlsx maze_reborn_cost_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_reborn_cost_v8【迷宫-复活消耗】.xlsx maze_reborn_cost_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_reborn_cost_v8【迷宫-复活消耗】.xlsx maze_reborn_cost_v8 data success.")
	return
}
