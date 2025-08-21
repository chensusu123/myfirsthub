package GMazeEnergyLevelV8Cfg

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

// MazeEnergyLevelV8ConfigRow from maze_energy_level_v8【迷宫-能力等级】.xlsx maze_energy_level_v8
type MazeEnergyLevelV8ConfigRow struct {
	Order              int32 `json:"order"`              // 序号（能力id*10000+等级
	Energy_id          int32 `json:"energy_id"`          // 能力id
	Energy_level       int32 `json:"energy_level"`       // 等级
	Max_energy         int32 `json:"max_energy"`         // 升到下一级需要的能量点数
	Energy_select      int32 `json:"energy_select"`      // 提升到当前等级时可以选择能力的次数
	Energy_item_select int32 `json:"energy_item_select"` // 提升到当前等级时，可以使用道具触发的选择能力次数
}

// MazeEnergyLevelV8Config from maze_energy_level_v8【迷宫-能力等级】.xlsx maze_energy_level_v8
type MazeEnergyLevelV8Config struct {
	ConfigRows map[int32]*MazeEnergyLevelV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEnergyLevelV8Config {
	ret := &MazeEnergyLevelV8Config{ConfigRows: map[int32]*MazeEnergyLevelV8ConfigRow{}}
	return ret
}

// GetMazeEnergyLevelV8Config get one config by configId
func (c *MazeEnergyLevelV8Config) GetMazeEnergyLevelV8Config(configId int32) *MazeEnergyLevelV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEnergyLevelV8Config) Get(configId int32) *MazeEnergyLevelV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEnergyLevelV8Config get all config slice
func (c *MazeEnergyLevelV8Config) GetAllMazeEnergyLevelV8Config() (res []*MazeEnergyLevelV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEnergyLevelV8Config) GetAll() (res []*MazeEnergyLevelV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEnergyLevelV8Config

// GetMazeEnergyLevelV8Config pkg func. get one config by configId
func GetMazeEnergyLevelV8Config(configId int32) *MazeEnergyLevelV8ConfigRow {
	return gConfigData.GetMazeEnergyLevelV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeEnergyLevelV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeEnergyLevelV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_energy_level_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeEnergyLevelV8Config pkg func. get all config slice
func GetAllMazeEnergyLevelV8Config() []*MazeEnergyLevelV8ConfigRow {
	return gConfigData.GetAllMazeEnergyLevelV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEnergyLevelV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEnergyLevelV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEnergyLevelV8ConfigRow from maze_energy_level_v8【迷宫-能力等级】.xlsx maze_energy_level_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEnergyLevelV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_energy_level_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_energy_level_v8.json",
		"maze_energy_level_v8【迷宫-能力等级】.xlsx", "maze_energy_level_v8",
		&gMazeEnergyLevelV8Parser{}, &gMazeEnergyLevelV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEnergyLevelV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEnergyLevelV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEnergyLevelV8Config))(c)
		return true
	})
}

// RegisterMazeEnergyLevelV8InitCallBack reg config update func (old func)
var RegisterMazeEnergyLevelV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEnergyLevelV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEnergyLevelV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEnergyLevelV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEnergyLevelV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEnergyLevelV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEnergyLevelV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEnergyLevelV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEnergyLevelV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEnergyLevelV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEnergyLevelV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEnergyLevelV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEnergyLevelV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEnergyLevelV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyLevelV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEnergyLevelV8ConfigRow", zap.String("xlsx", "maze_energy_level_v8【迷宫-能力等级】.xlsx"),
			zap.String("sheet", "maze_energy_level_v8"))
		return
	}
	config, ok := container.(*MazeEnergyLevelV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyLevelV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyLevelV8Config", zap.String("xlsx", "maze_energy_level_v8【迷宫-能力等级】.xlsx"),
			zap.String("sheet", "maze_energy_level_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEnergyLevelV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEnergyLevelV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyLevelV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyLevelV8Config", zap.String("xlsx", "maze_energy_level_v8【迷宫-能力等级】.xlsx"),
			zap.String("sheet", "maze_energy_level_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEnergyLevelV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEnergyLevelV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyLevelV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyLevelV8Config", zap.String("xlsx", "maze_energy_level_v8【迷宫-能力等级】.xlsx"),
			zap.String("sheet", "maze_energy_level_v8"))
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
type gMazeEnergyLevelV8Parser struct {
}

// New new config row data
func (*gMazeEnergyLevelV8Parser) New() interface{} {
	return &MazeEnergyLevelV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEnergyLevelV8Parser) Fields() []string {
	return gMazeEnergyLevelV8Fields
}

// Parse parse raw data to row data
func (*gMazeEnergyLevelV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEnergyLevelV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyLevelV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEnergyLevelV8ConfigRow", zap.String("xlsx", "maze_energy_level_v8【迷宫-能力等级】.xlsx"),
			zap.String("sheet", "maze_energy_level_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEnergyLevelV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEnergyLevelV8ConfigRow",
			zap.String("xlsx", "maze_energy_level_v8【迷宫-能力等级】.xlsx"),
			zap.String("sheet", "maze_energy_level_v8"), zap.Int("need_count", len(gMazeEnergyLevelV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 序号（能力id*10000+等级
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 序号（能力id*10000+等级 to int32 failed")
			logger.ErrorWF("parse field order 序号（能力id*10000+等级 to int32 failed.",
				zap.String("xlsx", "maze_energy_level_v8【迷宫-能力等级】.xlsx"), zap.String("sheet", "maze_energy_level_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 energy_id : 能力id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field energy_id 能力id to int32 failed")
			logger.ErrorWF("parse field energy_id 能力id to int32 failed.",
				zap.String("xlsx", "maze_energy_level_v8【迷宫-能力等级】.xlsx"), zap.String("sheet", "maze_energy_level_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Energy_id = int32(tmp)
	}

	// parse column 2 energy_level : 等级
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field energy_level 等级 to int32 failed")
			logger.ErrorWF("parse field energy_level 等级 to int32 failed.",
				zap.String("xlsx", "maze_energy_level_v8【迷宫-能力等级】.xlsx"), zap.String("sheet", "maze_energy_level_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Energy_level = int32(tmp)
	}

	// parse column 3 max_energy : 升到下一级需要的能量点数
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field max_energy 升到下一级需要的能量点数 to int32 failed")
			logger.ErrorWF("parse field max_energy 升到下一级需要的能量点数 to int32 failed.",
				zap.String("xlsx", "maze_energy_level_v8【迷宫-能力等级】.xlsx"), zap.String("sheet", "maze_energy_level_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Max_energy = int32(tmp)
	}

	// parse column 4 energy_select : 提升到当前等级时可以选择能力的次数
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field energy_select 提升到当前等级时可以选择能力的次数 to int32 failed")
			logger.ErrorWF("parse field energy_select 提升到当前等级时可以选择能力的次数 to int32 failed.",
				zap.String("xlsx", "maze_energy_level_v8【迷宫-能力等级】.xlsx"), zap.String("sheet", "maze_energy_level_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Energy_select = int32(tmp)
	}

	// parse column 5 energy_item_select : 提升到当前等级时，可以使用道具触发的选择能力次数
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field energy_item_select 提升到当前等级时，可以使用道具触发的选择能力次数 to int32 failed")
			logger.ErrorWF("parse field energy_item_select 提升到当前等级时，可以使用道具触发的选择能力次数 to int32 failed.",
				zap.String("xlsx", "maze_energy_level_v8【迷宫-能力等级】.xlsx"), zap.String("sheet", "maze_energy_level_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Energy_item_select = int32(tmp)
	}
	return
}

var gMazeEnergyLevelV8Fields = []string{
	"order",
	"energy_id",
	"energy_level",
	"max_energy",
	"energy_select",
	"energy_item_select",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEnergyLevelV8Parser{}
	loader := &gMazeEnergyLevelV8Loader{}
	var data [][]string
	data, err = load("maze_energy_level_v8【迷宫-能力等级】.xlsx", "maze_energy_level_v8", gMazeEnergyLevelV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_energy_level_v8【迷宫-能力等级】.xlsx maze_energy_level_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_energy_level_v8【迷宫-能力等级】.xlsx maze_energy_level_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEnergyLevelV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_energy_level_v8【迷宫-能力等级】.xlsx maze_energy_level_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEnergyLevelV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_energy_level_v8【迷宫-能力等级】.xlsx maze_energy_level_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_energy_level_v8【迷宫-能力等级】.xlsx maze_energy_level_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_energy_level_v8【迷宫-能力等级】.xlsx maze_energy_level_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_energy_level_v8【迷宫-能力等级】.xlsx maze_energy_level_v8 data success.")
	return
}
