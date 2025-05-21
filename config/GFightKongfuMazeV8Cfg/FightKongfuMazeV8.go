package GFightKongfuMazeV8Cfg

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

// FightKongfuMazeV8ConfigRow from fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx fight_kongfu_maze_v8
type FightKongfuMazeV8ConfigRow struct {
	Id             int32 `json:"id"`             // id
	Hp_lose_type   int32 `json:"hp_lose_type"`   // 战斗损血类型
	Pro__min       int64 `json:"pro__min"`       // 玩家武力/怪物武力结果万分比范围，最小值
	Pro_max        int64 `json:"pro_max"`        // 玩家武力/怪物武力结果万分比范围，最大值
	Player_hp_lose int32 `json:"player_hp_lose"` // 玩家掉血比例范围，最小值（百万分比）
	Foe_hp_lose    int32 `json:"foe_hp_lose"`    // 怪物掉血比例范围，最小值（百万分比）
	Round          int32 `json:"round"`          // 满血战斗回合数(万分比)
}

// FightKongfuMazeV8Config from fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx fight_kongfu_maze_v8
type FightKongfuMazeV8Config struct {
	ConfigRows map[int32]*FightKongfuMazeV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *FightKongfuMazeV8Config {
	ret := &FightKongfuMazeV8Config{ConfigRows: map[int32]*FightKongfuMazeV8ConfigRow{}}
	return ret
}

// GetFightKongfuMazeV8Config get one config by configId
func (c *FightKongfuMazeV8Config) GetFightKongfuMazeV8Config(configId int32) *FightKongfuMazeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *FightKongfuMazeV8Config) Get(configId int32) *FightKongfuMazeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllFightKongfuMazeV8Config get all config slice
func (c *FightKongfuMazeV8Config) GetAllFightKongfuMazeV8Config() (res []*FightKongfuMazeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *FightKongfuMazeV8Config) GetAll() (res []*FightKongfuMazeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *FightKongfuMazeV8Config

// GetFightKongfuMazeV8Config pkg func. get one config by configId
func GetFightKongfuMazeV8Config(configId int32) *FightKongfuMazeV8ConfigRow {
	return gConfigData.GetFightKongfuMazeV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *FightKongfuMazeV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllFightKongfuMazeV8Config pkg func. get all config slice
func GetAllFightKongfuMazeV8Config() []*FightKongfuMazeV8ConfigRow {
	return gConfigData.GetAllFightKongfuMazeV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*FightKongfuMazeV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*FightKongfuMazeV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "FightKongfuMazeV8ConfigRow from fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx fight_kongfu_maze_v8"
}

// GetRawValue get raw data
func GetRawValue() *FightKongfuMazeV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "fight_kongfu_maze_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("fight_kongfu_maze_v8.json",
		"fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx", "fight_kongfu_maze_v8",
		&gFightKongfuMazeV8Parser{}, &gFightKongfuMazeV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*FightKongfuMazeV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *FightKongfuMazeV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*FightKongfuMazeV8Config))(c)
		return true
	})
}

// RegisterFightKongfuMazeV8InitCallBack reg config update func (old func)
var RegisterFightKongfuMazeV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*FightKongfuMazeV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*FightKongfuMazeV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *FightKongfuMazeV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*FightKongfuMazeV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*FightKongfuMazeV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gFightKongfuMazeV8Loader struct {
}

// NewContainer new data container pointer
func (*gFightKongfuMazeV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gFightKongfuMazeV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*FightKongfuMazeV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gFightKongfuMazeV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*FightKongfuMazeV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gFightKongfuMazeV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*FightKongfuMazeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *FightKongfuMazeV8ConfigRow")
		logger.ErrorWF("invalid type. not *FightKongfuMazeV8ConfigRow", zap.String("xlsx", "fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx"),
			zap.String("sheet", "fight_kongfu_maze_v8"))
		return
	}
	config, ok := container.(*FightKongfuMazeV8Config)
	if !ok {
		err = errors.New("invalid type. not *FightKongfuMazeV8Config")
		logger.ErrorWF("invalid type. not *FightKongfuMazeV8Config", zap.String("xlsx", "fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx"),
			zap.String("sheet", "fight_kongfu_maze_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gFightKongfuMazeV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*FightKongfuMazeV8Config)
	if !ok {
		err = errors.New("invalid type. not *FightKongfuMazeV8Config")
		logger.ErrorWF("invalid type. not *FightKongfuMazeV8Config", zap.String("xlsx", "fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx"),
			zap.String("sheet", "fight_kongfu_maze_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gFightKongfuMazeV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*FightKongfuMazeV8Config)
	if !ok {
		err = errors.New("invalid type. not *FightKongfuMazeV8Config")
		logger.ErrorWF("invalid type. not *FightKongfuMazeV8Config", zap.String("xlsx", "fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx"),
			zap.String("sheet", "fight_kongfu_maze_v8"))
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
type gFightKongfuMazeV8Parser struct {
}

// New new config row data
func (*gFightKongfuMazeV8Parser) New() interface{} {
	return &FightKongfuMazeV8ConfigRow{}
}

// Fields get config fields names
func (*gFightKongfuMazeV8Parser) Fields() []string {
	return gFightKongfuMazeV8Fields
}

// Parse parse raw data to row data
func (*gFightKongfuMazeV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*FightKongfuMazeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *FightKongfuMazeV8ConfigRow")
		logger.ErrorWF("invalid type. not *FightKongfuMazeV8ConfigRow", zap.String("xlsx", "fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx"),
			zap.String("sheet", "fight_kongfu_maze_v8"))
		return
	}
	// compare length
	if len(data) != len(gFightKongfuMazeV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*FightKongfuMazeV8ConfigRow",
			zap.String("xlsx", "fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx"),
			zap.String("sheet", "fight_kongfu_maze_v8"), zap.Int("need_count", len(gFightKongfuMazeV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id id to int32 failed")
			logger.ErrorWF("parse field id id to int32 failed.",
				zap.String("xlsx", "fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx"), zap.String("sheet", "fight_kongfu_maze_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 hp_lose_type : 战斗损血类型
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field hp_lose_type 战斗损血类型 to int32 failed")
			logger.ErrorWF("parse field hp_lose_type 战斗损血类型 to int32 failed.",
				zap.String("xlsx", "fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx"), zap.String("sheet", "fight_kongfu_maze_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Hp_lose_type = int32(tmp)
	}

	// parse column 2 pro__min : 玩家武力/怪物武力结果万分比范围，最小值
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field pro__min 玩家武力/怪物武力结果万分比范围，最小值 to int64 failed")
			logger.ErrorWF("parse field pro__min 玩家武力/怪物武力结果万分比范围，最小值 to int64 failed.",
				zap.String("xlsx", "fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx"), zap.String("sheet", "fight_kongfu_maze_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Pro__min = int64(tmp)
	}

	// parse column 3 pro_max : 玩家武力/怪物武力结果万分比范围，最大值
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field pro_max 玩家武力/怪物武力结果万分比范围，最大值 to int64 failed")
			logger.ErrorWF("parse field pro_max 玩家武力/怪物武力结果万分比范围，最大值 to int64 failed.",
				zap.String("xlsx", "fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx"), zap.String("sheet", "fight_kongfu_maze_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Pro_max = int64(tmp)
	}

	// parse column 4 player_hp_lose : 玩家掉血比例范围，最小值（百万分比）
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field player_hp_lose 玩家掉血比例范围，最小值（百万分比） to int32 failed")
			logger.ErrorWF("parse field player_hp_lose 玩家掉血比例范围，最小值（百万分比） to int32 failed.",
				zap.String("xlsx", "fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx"), zap.String("sheet", "fight_kongfu_maze_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Player_hp_lose = int32(tmp)
	}

	// parse column 5 foe_hp_lose : 怪物掉血比例范围，最小值（百万分比）
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field foe_hp_lose 怪物掉血比例范围，最小值（百万分比） to int32 failed")
			logger.ErrorWF("parse field foe_hp_lose 怪物掉血比例范围，最小值（百万分比） to int32 failed.",
				zap.String("xlsx", "fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx"), zap.String("sheet", "fight_kongfu_maze_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Foe_hp_lose = int32(tmp)
	}

	// parse column 6 round : 满血战斗回合数(万分比)
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field round 满血战斗回合数(万分比) to int32 failed")
			logger.ErrorWF("parse field round 满血战斗回合数(万分比) to int32 failed.",
				zap.String("xlsx", "fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx"), zap.String("sheet", "fight_kongfu_maze_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Round = int32(tmp)
	}
	return
}

var gFightKongfuMazeV8Fields = []string{
	"id",
	"hp_lose_type",
	"pro__min",
	"pro_max",
	"player_hp_lose",
	"foe_hp_lose",
	"round",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gFightKongfuMazeV8Parser{}
	loader := &gFightKongfuMazeV8Loader{}
	var data [][]string
	data, err = load("fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx", "fight_kongfu_maze_v8", gFightKongfuMazeV8Fields)
	if err != nil {
		logger.ErrorWF("load fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx fight_kongfu_maze_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx fight_kongfu_maze_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gFightKongfuMazeV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx fight_kongfu_maze_v8 failed.", zap.Int("row", k), zap.Strings("need", gFightKongfuMazeV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx fight_kongfu_maze_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx fight_kongfu_maze_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx fight_kongfu_maze_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load fight_kongfu_maze_v8【人偶-迷宫-武力比损血表】.xlsx fight_kongfu_maze_v8 data success.")
	return
}
