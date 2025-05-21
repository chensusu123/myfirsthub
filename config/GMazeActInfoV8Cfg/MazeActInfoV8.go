package GMazeActInfoV8Cfg

import (
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

// MazeActInfoV8ConfigRow from maze_act_info_v8【迷宫-动作配置】.xlsx maze_act_info_v8
type MazeActInfoV8ConfigRow struct {
	Order                     int32           `json:"order"`                     // 动作id
	Attack_point_damage_ratio map[int32]int32 `json:"attack_point_damage_ratio"` // 打击点:伤害系数
	Tough_broke_value         int32           `json:"tough_broke_value"`         // 削韧值
	Temp_tough                int32           `json:"temp_tough"`                // 动作临时韧性
}

// MazeActInfoV8Config from maze_act_info_v8【迷宫-动作配置】.xlsx maze_act_info_v8
type MazeActInfoV8Config struct {
	ConfigRows map[int32]*MazeActInfoV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeActInfoV8Config {
	ret := &MazeActInfoV8Config{ConfigRows: map[int32]*MazeActInfoV8ConfigRow{}}
	return ret
}

// GetMazeActInfoV8Config get one config by configId
func (c *MazeActInfoV8Config) GetMazeActInfoV8Config(configId int32) *MazeActInfoV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeActInfoV8Config) Get(configId int32) *MazeActInfoV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeActInfoV8Config get all config slice
func (c *MazeActInfoV8Config) GetAllMazeActInfoV8Config() (res []*MazeActInfoV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeActInfoV8Config) GetAll() (res []*MazeActInfoV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeActInfoV8Config

// GetMazeActInfoV8Config pkg func. get one config by configId
func GetMazeActInfoV8Config(configId int32) *MazeActInfoV8ConfigRow {
	return gConfigData.GetMazeActInfoV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeActInfoV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeActInfoV8Config pkg func. get all config slice
func GetAllMazeActInfoV8Config() []*MazeActInfoV8ConfigRow {
	return gConfigData.GetAllMazeActInfoV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeActInfoV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeActInfoV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeActInfoV8ConfigRow from maze_act_info_v8【迷宫-动作配置】.xlsx maze_act_info_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeActInfoV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_act_info_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_act_info_v8.json",
		"maze_act_info_v8【迷宫-动作配置】.xlsx", "maze_act_info_v8",
		&gMazeActInfoV8Parser{}, &gMazeActInfoV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeActInfoV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeActInfoV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeActInfoV8Config))(c)
		return true
	})
}

// RegisterMazeActInfoV8InitCallBack reg config update func (old func)
var RegisterMazeActInfoV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeActInfoV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeActInfoV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeActInfoV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeActInfoV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeActInfoV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeActInfoV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeActInfoV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeActInfoV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeActInfoV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeActInfoV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeActInfoV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeActInfoV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeActInfoV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeActInfoV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeActInfoV8ConfigRow", zap.String("xlsx", "maze_act_info_v8【迷宫-动作配置】.xlsx"),
			zap.String("sheet", "maze_act_info_v8"))
		return
	}
	config, ok := container.(*MazeActInfoV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeActInfoV8Config")
		logger.ErrorWF("invalid type. not *MazeActInfoV8Config", zap.String("xlsx", "maze_act_info_v8【迷宫-动作配置】.xlsx"),
			zap.String("sheet", "maze_act_info_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeActInfoV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeActInfoV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeActInfoV8Config")
		logger.ErrorWF("invalid type. not *MazeActInfoV8Config", zap.String("xlsx", "maze_act_info_v8【迷宫-动作配置】.xlsx"),
			zap.String("sheet", "maze_act_info_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeActInfoV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeActInfoV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeActInfoV8Config")
		logger.ErrorWF("invalid type. not *MazeActInfoV8Config", zap.String("xlsx", "maze_act_info_v8【迷宫-动作配置】.xlsx"),
			zap.String("sheet", "maze_act_info_v8"))
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
type gMazeActInfoV8Parser struct {
}

// New new config row data
func (*gMazeActInfoV8Parser) New() interface{} {
	return &MazeActInfoV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeActInfoV8Parser) Fields() []string {
	return gMazeActInfoV8Fields
}

// Parse parse raw data to row data
func (*gMazeActInfoV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeActInfoV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeActInfoV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeActInfoV8ConfigRow", zap.String("xlsx", "maze_act_info_v8【迷宫-动作配置】.xlsx"),
			zap.String("sheet", "maze_act_info_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeActInfoV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeActInfoV8ConfigRow",
			zap.String("xlsx", "maze_act_info_v8【迷宫-动作配置】.xlsx"),
			zap.String("sheet", "maze_act_info_v8"), zap.Int("need_count", len(gMazeActInfoV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 动作id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 动作id to int32 failed")
			logger.ErrorWF("parse field order 动作id to int32 failed.",
				zap.String("xlsx", "maze_act_info_v8【迷宫-动作配置】.xlsx"), zap.String("sheet", "maze_act_info_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 attack_point_damage_ratio : 打击点:伤害系数
	if data[1] != "" {

		config.Attack_point_damage_ratio = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[1], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field attack_point_damage_ratio 打击点:伤害系数 to key int32 failed")
				logger.ErrorWF("parse map field attack_point_damage_ratio 打击点:伤害系数 to key int32 failed.",
					zap.String("xlsx", "maze_act_info_v8【迷宫-动作配置】.xlsx"), zap.String("sheet", "maze_act_info_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field attack_point_damage_ratio 打击点:伤害系数 to value int32 failed")
				logger.ErrorWF("parse map field attack_point_damage_ratio 打击点:伤害系数 to value int32 failed.",
					zap.String("xlsx", "maze_act_info_v8【迷宫-动作配置】.xlsx"), zap.String("sheet", "maze_act_info_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Attack_point_damage_ratio[key] = value
		}
	}

	// parse column 2 tough_broke_value : 削韧值
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field tough_broke_value 削韧值 to int32 failed")
			logger.ErrorWF("parse field tough_broke_value 削韧值 to int32 failed.",
				zap.String("xlsx", "maze_act_info_v8【迷宫-动作配置】.xlsx"), zap.String("sheet", "maze_act_info_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Tough_broke_value = int32(tmp)
	}

	// parse column 3 temp_tough : 动作临时韧性
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field temp_tough 动作临时韧性 to int32 failed")
			logger.ErrorWF("parse field temp_tough 动作临时韧性 to int32 failed.",
				zap.String("xlsx", "maze_act_info_v8【迷宫-动作配置】.xlsx"), zap.String("sheet", "maze_act_info_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Temp_tough = int32(tmp)
	}
	return
}

var gMazeActInfoV8Fields = []string{
	"order",
	"attack_point_damage_ratio",
	"tough_broke_value",
	"temp_tough",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeActInfoV8Parser{}
	loader := &gMazeActInfoV8Loader{}
	var data [][]string
	data, err = load("maze_act_info_v8【迷宫-动作配置】.xlsx", "maze_act_info_v8", gMazeActInfoV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_act_info_v8【迷宫-动作配置】.xlsx maze_act_info_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_act_info_v8【迷宫-动作配置】.xlsx maze_act_info_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeActInfoV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_act_info_v8【迷宫-动作配置】.xlsx maze_act_info_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeActInfoV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_act_info_v8【迷宫-动作配置】.xlsx maze_act_info_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_act_info_v8【迷宫-动作配置】.xlsx maze_act_info_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_act_info_v8【迷宫-动作配置】.xlsx maze_act_info_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_act_info_v8【迷宫-动作配置】.xlsx maze_act_info_v8 data success.")
	return
}
