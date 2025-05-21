package GMazeBarriesEnergyLevelV8Cfg

import (
	"errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"
)

// MazeBarriesEnergyLevelV8ConfigRow from maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx maze_barries_energy_level_v8
type MazeBarriesEnergyLevelV8ConfigRow struct {
	Order             int32 `json:"order"`             // 等级
	Max_energy        int32 `json:"max_energy"`        // 升到下一级需要的能量点数
	Max_energy_all    int32 `json:"max_energy_all"`    // 到当前等级能量点数总数
	Energy_select     int32 `json:"energy_select"`     // 提升到当前等级时可以选择能力的次数
	Energy_select_all int32 `json:"energy_select_all"` // 当前等级选择能力次数总数
}

// MazeBarriesEnergyLevelV8Config from maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx maze_barries_energy_level_v8
type MazeBarriesEnergyLevelV8Config struct {
	ConfigRows map[int32]*MazeBarriesEnergyLevelV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeBarriesEnergyLevelV8Config {
	ret := &MazeBarriesEnergyLevelV8Config{ConfigRows: map[int32]*MazeBarriesEnergyLevelV8ConfigRow{}}
	return ret
}

// GetMazeBarriesEnergyLevelV8Config get one config by configId
func (c *MazeBarriesEnergyLevelV8Config) GetMazeBarriesEnergyLevelV8Config(configId int32) *MazeBarriesEnergyLevelV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeBarriesEnergyLevelV8Config) Get(configId int32) *MazeBarriesEnergyLevelV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeBarriesEnergyLevelV8Config get all config slice
func (c *MazeBarriesEnergyLevelV8Config) GetAllMazeBarriesEnergyLevelV8Config() (res []*MazeBarriesEnergyLevelV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeBarriesEnergyLevelV8Config) GetAll() (res []*MazeBarriesEnergyLevelV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeBarriesEnergyLevelV8Config

// GetMazeBarriesEnergyLevelV8Config pkg func. get one config by configId
func GetMazeBarriesEnergyLevelV8Config(configId int32) *MazeBarriesEnergyLevelV8ConfigRow {
	return gConfigData.GetMazeBarriesEnergyLevelV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeBarriesEnergyLevelV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeBarriesEnergyLevelV8Config pkg func. get all config slice
func GetAllMazeBarriesEnergyLevelV8Config() []*MazeBarriesEnergyLevelV8ConfigRow {
	return gConfigData.GetAllMazeBarriesEnergyLevelV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeBarriesEnergyLevelV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeBarriesEnergyLevelV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeBarriesEnergyLevelV8ConfigRow from maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx maze_barries_energy_level_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeBarriesEnergyLevelV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_barries_energy_level_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_barries_energy_level_v8.json",
		"maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx", "maze_barries_energy_level_v8",
		&gMazeBarriesEnergyLevelV8Parser{}, &gMazeBarriesEnergyLevelV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeBarriesEnergyLevelV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeBarriesEnergyLevelV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeBarriesEnergyLevelV8Config))(c)
		return true
	})
}

// RegisterMazeBarriesEnergyLevelV8InitCallBack reg config update func (old func)
var RegisterMazeBarriesEnergyLevelV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeBarriesEnergyLevelV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeBarriesEnergyLevelV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeBarriesEnergyLevelV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeBarriesEnergyLevelV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeBarriesEnergyLevelV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeBarriesEnergyLevelV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeBarriesEnergyLevelV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeBarriesEnergyLevelV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeBarriesEnergyLevelV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeBarriesEnergyLevelV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeBarriesEnergyLevelV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeBarriesEnergyLevelV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeBarriesEnergyLevelV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesEnergyLevelV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBarriesEnergyLevelV8ConfigRow", zap.String("xlsx", "maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx"),
			zap.String("sheet", "maze_barries_energy_level_v8"))
		return
	}
	config, ok := container.(*MazeBarriesEnergyLevelV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesEnergyLevelV8Config")
		logger.ErrorWF("invalid type. not *MazeBarriesEnergyLevelV8Config", zap.String("xlsx", "maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx"),
			zap.String("sheet", "maze_barries_energy_level_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeBarriesEnergyLevelV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeBarriesEnergyLevelV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesEnergyLevelV8Config")
		logger.ErrorWF("invalid type. not *MazeBarriesEnergyLevelV8Config", zap.String("xlsx", "maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx"),
			zap.String("sheet", "maze_barries_energy_level_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeBarriesEnergyLevelV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeBarriesEnergyLevelV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesEnergyLevelV8Config")
		logger.ErrorWF("invalid type. not *MazeBarriesEnergyLevelV8Config", zap.String("xlsx", "maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx"),
			zap.String("sheet", "maze_barries_energy_level_v8"))
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
type gMazeBarriesEnergyLevelV8Parser struct {
}

// New new config row data
func (*gMazeBarriesEnergyLevelV8Parser) New() interface{} {
	return &MazeBarriesEnergyLevelV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeBarriesEnergyLevelV8Parser) Fields() []string {
	return gMazeBarriesEnergyLevelV8Fields
}

// Parse parse raw data to row data
func (*gMazeBarriesEnergyLevelV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeBarriesEnergyLevelV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesEnergyLevelV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBarriesEnergyLevelV8ConfigRow", zap.String("xlsx", "maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx"),
			zap.String("sheet", "maze_barries_energy_level_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeBarriesEnergyLevelV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeBarriesEnergyLevelV8ConfigRow",
			zap.String("xlsx", "maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx"),
			zap.String("sheet", "maze_barries_energy_level_v8"), zap.Int("need_count", len(gMazeBarriesEnergyLevelV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 等级
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 等级 to int32 failed")
			logger.ErrorWF("parse field order 等级 to int32 failed.",
				zap.String("xlsx", "maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx"), zap.String("sheet", "maze_barries_energy_level_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 max_energy : 升到下一级需要的能量点数
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field max_energy 升到下一级需要的能量点数 to int32 failed")
			logger.ErrorWF("parse field max_energy 升到下一级需要的能量点数 to int32 failed.",
				zap.String("xlsx", "maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx"), zap.String("sheet", "maze_barries_energy_level_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Max_energy = int32(tmp)
	}

	// parse column 2 max_energy_all : 到当前等级能量点数总数
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field max_energy_all 到当前等级能量点数总数 to int32 failed")
			logger.ErrorWF("parse field max_energy_all 到当前等级能量点数总数 to int32 failed.",
				zap.String("xlsx", "maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx"), zap.String("sheet", "maze_barries_energy_level_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Max_energy_all = int32(tmp)
	}

	// parse column 3 energy_select : 提升到当前等级时可以选择能力的次数
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field energy_select 提升到当前等级时可以选择能力的次数 to int32 failed")
			logger.ErrorWF("parse field energy_select 提升到当前等级时可以选择能力的次数 to int32 failed.",
				zap.String("xlsx", "maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx"), zap.String("sheet", "maze_barries_energy_level_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Energy_select = int32(tmp)
	}

	// parse column 4 energy_select_all : 当前等级选择能力次数总数
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field energy_select_all 当前等级选择能力次数总数 to int32 failed")
			logger.ErrorWF("parse field energy_select_all 当前等级选择能力次数总数 to int32 failed.",
				zap.String("xlsx", "maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx"), zap.String("sheet", "maze_barries_energy_level_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Energy_select_all = int32(tmp)
	}
	return
}

var gMazeBarriesEnergyLevelV8Fields = []string{
	"order",
	"max_energy",
	"max_energy_all",
	"energy_select",
	"energy_select_all",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeBarriesEnergyLevelV8Parser{}
	loader := &gMazeBarriesEnergyLevelV8Loader{}
	var data [][]string
	data, err = load("maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx", "maze_barries_energy_level_v8", gMazeBarriesEnergyLevelV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx maze_barries_energy_level_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx maze_barries_energy_level_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeBarriesEnergyLevelV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx maze_barries_energy_level_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeBarriesEnergyLevelV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx maze_barries_energy_level_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx maze_barries_energy_level_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx maze_barries_energy_level_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_barries_energy_level_v8【迷宫-关卡内能力等级】.xlsx maze_barries_energy_level_v8 data success.")
	return
}
