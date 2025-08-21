package GMazeKongfuDeathV8Cfg

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

// MazeKongfuDeathV8ConfigRow from maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx maze_kongfu_death_v8
type MazeKongfuDeathV8ConfigRow struct {
	Order          int32 `json:"order"`          // 序号
	Kongfu_min     int64 `json:"kongfu_min"`     // 武力值区间，下限
	Kongfu_max     int64 `json:"kongfu_max"`     // 武力值区间，上限
	Death_less_pro int32 `json:"death_less_pro"` // 死亡扣除金币万分比
	Death_less_min int64 `json:"death_less_min"` // 死亡扣除金币最小数
	Death_less_max int64 `json:"death_less_max"` // 死亡扣除金币最大数
}

// MazeKongfuDeathV8Config from maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx maze_kongfu_death_v8
type MazeKongfuDeathV8Config struct {
	ConfigRows map[int32]*MazeKongfuDeathV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeKongfuDeathV8Config {
	ret := &MazeKongfuDeathV8Config{ConfigRows: map[int32]*MazeKongfuDeathV8ConfigRow{}}
	return ret
}

// GetMazeKongfuDeathV8Config get one config by configId
func (c *MazeKongfuDeathV8Config) GetMazeKongfuDeathV8Config(configId int32) *MazeKongfuDeathV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeKongfuDeathV8Config) Get(configId int32) *MazeKongfuDeathV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeKongfuDeathV8Config get all config slice
func (c *MazeKongfuDeathV8Config) GetAllMazeKongfuDeathV8Config() (res []*MazeKongfuDeathV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeKongfuDeathV8Config) GetAll() (res []*MazeKongfuDeathV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeKongfuDeathV8Config

// GetMazeKongfuDeathV8Config pkg func. get one config by configId
func GetMazeKongfuDeathV8Config(configId int32) *MazeKongfuDeathV8ConfigRow {
	return gConfigData.GetMazeKongfuDeathV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeKongfuDeathV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeKongfuDeathV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_kongfu_death_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeKongfuDeathV8Config pkg func. get all config slice
func GetAllMazeKongfuDeathV8Config() []*MazeKongfuDeathV8ConfigRow {
	return gConfigData.GetAllMazeKongfuDeathV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeKongfuDeathV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeKongfuDeathV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeKongfuDeathV8ConfigRow from maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx maze_kongfu_death_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeKongfuDeathV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_kongfu_death_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_kongfu_death_v8.json",
		"maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx", "maze_kongfu_death_v8",
		&gMazeKongfuDeathV8Parser{}, &gMazeKongfuDeathV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeKongfuDeathV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeKongfuDeathV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeKongfuDeathV8Config))(c)
		return true
	})
}

// RegisterMazeKongfuDeathV8InitCallBack reg config update func (old func)
var RegisterMazeKongfuDeathV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeKongfuDeathV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeKongfuDeathV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeKongfuDeathV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeKongfuDeathV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeKongfuDeathV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeKongfuDeathV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeKongfuDeathV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeKongfuDeathV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeKongfuDeathV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeKongfuDeathV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeKongfuDeathV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeKongfuDeathV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeKongfuDeathV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuDeathV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeKongfuDeathV8ConfigRow", zap.String("xlsx", "maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx"),
			zap.String("sheet", "maze_kongfu_death_v8"))
		return
	}
	config, ok := container.(*MazeKongfuDeathV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuDeathV8Config")
		logger.ErrorWF("invalid type. not *MazeKongfuDeathV8Config", zap.String("xlsx", "maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx"),
			zap.String("sheet", "maze_kongfu_death_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeKongfuDeathV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeKongfuDeathV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuDeathV8Config")
		logger.ErrorWF("invalid type. not *MazeKongfuDeathV8Config", zap.String("xlsx", "maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx"),
			zap.String("sheet", "maze_kongfu_death_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeKongfuDeathV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeKongfuDeathV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuDeathV8Config")
		logger.ErrorWF("invalid type. not *MazeKongfuDeathV8Config", zap.String("xlsx", "maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx"),
			zap.String("sheet", "maze_kongfu_death_v8"))
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
type gMazeKongfuDeathV8Parser struct {
}

// New new config row data
func (*gMazeKongfuDeathV8Parser) New() interface{} {
	return &MazeKongfuDeathV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeKongfuDeathV8Parser) Fields() []string {
	return gMazeKongfuDeathV8Fields
}

// Parse parse raw data to row data
func (*gMazeKongfuDeathV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeKongfuDeathV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuDeathV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeKongfuDeathV8ConfigRow", zap.String("xlsx", "maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx"),
			zap.String("sheet", "maze_kongfu_death_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeKongfuDeathV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeKongfuDeathV8ConfigRow",
			zap.String("xlsx", "maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx"),
			zap.String("sheet", "maze_kongfu_death_v8"), zap.Int("need_count", len(gMazeKongfuDeathV8Fields)),
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
				zap.String("xlsx", "maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx"), zap.String("sheet", "maze_kongfu_death_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 kongfu_min : 武力值区间，下限
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field kongfu_min 武力值区间，下限 to int64 failed")
			logger.ErrorWF("parse field kongfu_min 武力值区间，下限 to int64 failed.",
				zap.String("xlsx", "maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx"), zap.String("sheet", "maze_kongfu_death_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Kongfu_min = int64(tmp)
	}

	// parse column 2 kongfu_max : 武力值区间，上限
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field kongfu_max 武力值区间，上限 to int64 failed")
			logger.ErrorWF("parse field kongfu_max 武力值区间，上限 to int64 failed.",
				zap.String("xlsx", "maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx"), zap.String("sheet", "maze_kongfu_death_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Kongfu_max = int64(tmp)
	}

	// parse column 3 death_less_pro : 死亡扣除金币万分比
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field death_less_pro 死亡扣除金币万分比 to int32 failed")
			logger.ErrorWF("parse field death_less_pro 死亡扣除金币万分比 to int32 failed.",
				zap.String("xlsx", "maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx"), zap.String("sheet", "maze_kongfu_death_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Death_less_pro = int32(tmp)
	}

	// parse column 4 death_less_min : 死亡扣除金币最小数
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field death_less_min 死亡扣除金币最小数 to int64 failed")
			logger.ErrorWF("parse field death_less_min 死亡扣除金币最小数 to int64 failed.",
				zap.String("xlsx", "maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx"), zap.String("sheet", "maze_kongfu_death_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Death_less_min = int64(tmp)
	}

	// parse column 5 death_less_max : 死亡扣除金币最大数
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field death_less_max 死亡扣除金币最大数 to int64 failed")
			logger.ErrorWF("parse field death_less_max 死亡扣除金币最大数 to int64 failed.",
				zap.String("xlsx", "maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx"), zap.String("sheet", "maze_kongfu_death_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Death_less_max = int64(tmp)
	}
	return
}

var gMazeKongfuDeathV8Fields = []string{
	"order",
	"kongfu_min",
	"kongfu_max",
	"death_less_pro",
	"death_less_min",
	"death_less_max",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeKongfuDeathV8Parser{}
	loader := &gMazeKongfuDeathV8Loader{}
	var data [][]string
	data, err = load("maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx", "maze_kongfu_death_v8", gMazeKongfuDeathV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx maze_kongfu_death_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx maze_kongfu_death_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeKongfuDeathV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx maze_kongfu_death_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeKongfuDeathV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx maze_kongfu_death_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx maze_kongfu_death_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx maze_kongfu_death_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_kongfu_death_v8【迷宫-武力值对应死亡扣除】.xlsx maze_kongfu_death_v8 data success.")
	return
}
