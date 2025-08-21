package GMazeSummonV8Cfg

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

// MazeSummonV8ConfigRow from maze_summon_v8【迷宫-召唤物】.xlsx maze_summon_v8
type MazeSummonV8ConfigRow struct {
	Id                  int32   `json:"id"`                  // 召唤物id
	Name                string  `json:"name"`                // 召唤物名称
	Suffer_damage       int32   `json:"suffer_damage"`       // 是否承伤
	Duration            int32   `json:"duration"`            // 持续时间（毫秒）
	Max_number          int32   `json:"max_number"`          // 同时存在数量上限
	Model_id            int32   `json:"model_id"`            // 资源组id
	Nor_attack_skill_id int32   `json:"nor_attack_skill_id"` // 普通攻击技能id
	Skill_id            []int32 `json:"skill_id"`            // 技能id
	Speed               int32   `json:"speed"`               // 移动速度(万分比）
	Attack_speed_pro    int32   `json:"attack_speed_pro"`    // 攻击速度系数（>10000加速 ,<10000减速
	Inheritance_attack  int32   `json:"inheritance_attack"`  // 继承攻击万分比
	Inheritance_def     int32   `json:"inheritance_def"`     // 继承防御万分比
	Inheritance_hp      int32   `json:"inheritance_hp"`      // 继承生命万分比
	Attack_value        int32   `json:"attack_value"`        // 额外攻击固定值
	Def_value           int32   `json:"def_value"`           // 额外防御固定值
	Hp_value            int32   `json:"hp_value"`            // 额外生命固定值
}

// MazeSummonV8Config from maze_summon_v8【迷宫-召唤物】.xlsx maze_summon_v8
type MazeSummonV8Config struct {
	ConfigRows map[int32]*MazeSummonV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeSummonV8Config {
	ret := &MazeSummonV8Config{ConfigRows: map[int32]*MazeSummonV8ConfigRow{}}
	return ret
}

// GetMazeSummonV8Config get one config by configId
func (c *MazeSummonV8Config) GetMazeSummonV8Config(configId int32) *MazeSummonV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeSummonV8Config) Get(configId int32) *MazeSummonV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeSummonV8Config get all config slice
func (c *MazeSummonV8Config) GetAllMazeSummonV8Config() (res []*MazeSummonV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeSummonV8Config) GetAll() (res []*MazeSummonV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeSummonV8Config

// GetMazeSummonV8Config pkg func. get one config by configId
func GetMazeSummonV8Config(configId int32) *MazeSummonV8ConfigRow {
	return gConfigData.GetMazeSummonV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeSummonV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeSummonV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_summon_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeSummonV8Config pkg func. get all config slice
func GetAllMazeSummonV8Config() []*MazeSummonV8ConfigRow {
	return gConfigData.GetAllMazeSummonV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeSummonV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeSummonV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeSummonV8ConfigRow from maze_summon_v8【迷宫-召唤物】.xlsx maze_summon_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeSummonV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_summon_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_summon_v8.json",
		"maze_summon_v8【迷宫-召唤物】.xlsx", "maze_summon_v8",
		&gMazeSummonV8Parser{}, &gMazeSummonV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeSummonV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeSummonV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeSummonV8Config))(c)
		return true
	})
}

// RegisterMazeSummonV8InitCallBack reg config update func (old func)
var RegisterMazeSummonV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeSummonV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeSummonV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeSummonV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeSummonV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeSummonV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeSummonV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeSummonV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeSummonV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeSummonV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeSummonV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeSummonV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeSummonV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeSummonV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSummonV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSummonV8ConfigRow", zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"),
			zap.String("sheet", "maze_summon_v8"))
		return
	}
	config, ok := container.(*MazeSummonV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSummonV8Config")
		logger.ErrorWF("invalid type. not *MazeSummonV8Config", zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"),
			zap.String("sheet", "maze_summon_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeSummonV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeSummonV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSummonV8Config")
		logger.ErrorWF("invalid type. not *MazeSummonV8Config", zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"),
			zap.String("sheet", "maze_summon_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeSummonV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeSummonV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSummonV8Config")
		logger.ErrorWF("invalid type. not *MazeSummonV8Config", zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"),
			zap.String("sheet", "maze_summon_v8"))
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
type gMazeSummonV8Parser struct {
}

// New new config row data
func (*gMazeSummonV8Parser) New() interface{} {
	return &MazeSummonV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeSummonV8Parser) Fields() []string {
	return gMazeSummonV8Fields
}

// Parse parse raw data to row data
func (*gMazeSummonV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeSummonV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSummonV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSummonV8ConfigRow", zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"),
			zap.String("sheet", "maze_summon_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeSummonV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeSummonV8ConfigRow",
			zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"),
			zap.String("sheet", "maze_summon_v8"), zap.Int("need_count", len(gMazeSummonV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 召唤物id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 召唤物id to int32 failed")
			logger.ErrorWF("parse field id 召唤物id to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 name : 召唤物名称
	if data[1] != "" {
		config.Name = data[1]
	}

	// parse column 2 suffer_damage : 是否承伤
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field suffer_damage 是否承伤 to int32 failed")
			logger.ErrorWF("parse field suffer_damage 是否承伤 to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Suffer_damage = int32(tmp)
	}

	// parse column 3 duration : 持续时间（毫秒）
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field duration 持续时间（毫秒） to int32 failed")
			logger.ErrorWF("parse field duration 持续时间（毫秒） to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Duration = int32(tmp)
	}

	// parse column 4 max_number : 同时存在数量上限
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field max_number 同时存在数量上限 to int32 failed")
			logger.ErrorWF("parse field max_number 同时存在数量上限 to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Max_number = int32(tmp)
	}

	// parse column 5 model_id : 资源组id
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field model_id 资源组id to int32 failed")
			logger.ErrorWF("parse field model_id 资源组id to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Model_id = int32(tmp)
	}

	// parse column 6 nor_attack_skill_id : 普通攻击技能id
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field nor_attack_skill_id 普通攻击技能id to int32 failed")
			logger.ErrorWF("parse field nor_attack_skill_id 普通攻击技能id to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Nor_attack_skill_id = int32(tmp)
	}

	// parse column 7 skill_id : 技能id
	if data[7] != "" {

		vals := strings.Split(data[7], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field skill_id 技能id to []int32 failed")
				logger.ErrorWF("parse array field skill_id 技能id to []int32 failed.",
					zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
					// zap.String("field_data",data[7]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Skill_id = append(config.Skill_id, int32(tmp))
		}
	}

	// parse column 8 speed : 移动速度(万分比）
	if data[8] != "" {
		tmp, err = strconv.ParseInt(data[8], 10, 64)
		if err != nil {
			err = errors.New("parse field speed 移动速度(万分比） to int32 failed")
			logger.ErrorWF("parse field speed 移动速度(万分比） to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[8]),
				zap.Error(err))
			return
		}
		config.Speed = int32(tmp)
	}

	// parse column 9 attack_speed_pro : 攻击速度系数（>10000加速 ,<10000减速
	if data[9] != "" {
		tmp, err = strconv.ParseInt(data[9], 10, 64)
		if err != nil {
			err = errors.New("parse field attack_speed_pro 攻击速度系数（>10000加速 ,<10000减速 to int32 failed")
			logger.ErrorWF("parse field attack_speed_pro 攻击速度系数（>10000加速 ,<10000减速 to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[9]),
				zap.Error(err))
			return
		}
		config.Attack_speed_pro = int32(tmp)
	}

	// parse column 10 inheritance_attack : 继承攻击万分比
	if data[10] != "" {
		tmp, err = strconv.ParseInt(data[10], 10, 64)
		if err != nil {
			err = errors.New("parse field inheritance_attack 继承攻击万分比 to int32 failed")
			logger.ErrorWF("parse field inheritance_attack 继承攻击万分比 to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[10]),
				zap.Error(err))
			return
		}
		config.Inheritance_attack = int32(tmp)
	}

	// parse column 11 inheritance_def : 继承防御万分比
	if data[11] != "" {
		tmp, err = strconv.ParseInt(data[11], 10, 64)
		if err != nil {
			err = errors.New("parse field inheritance_def 继承防御万分比 to int32 failed")
			logger.ErrorWF("parse field inheritance_def 继承防御万分比 to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[11]),
				zap.Error(err))
			return
		}
		config.Inheritance_def = int32(tmp)
	}

	// parse column 12 inheritance_hp : 继承生命万分比
	if data[12] != "" {
		tmp, err = strconv.ParseInt(data[12], 10, 64)
		if err != nil {
			err = errors.New("parse field inheritance_hp 继承生命万分比 to int32 failed")
			logger.ErrorWF("parse field inheritance_hp 继承生命万分比 to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[12]),
				zap.Error(err))
			return
		}
		config.Inheritance_hp = int32(tmp)
	}

	// parse column 13 attack_value : 额外攻击固定值
	if data[13] != "" {
		tmp, err = strconv.ParseInt(data[13], 10, 64)
		if err != nil {
			err = errors.New("parse field attack_value 额外攻击固定值 to int32 failed")
			logger.ErrorWF("parse field attack_value 额外攻击固定值 to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[13]),
				zap.Error(err))
			return
		}
		config.Attack_value = int32(tmp)
	}

	// parse column 14 def_value : 额外防御固定值
	if data[14] != "" {
		tmp, err = strconv.ParseInt(data[14], 10, 64)
		if err != nil {
			err = errors.New("parse field def_value 额外防御固定值 to int32 failed")
			logger.ErrorWF("parse field def_value 额外防御固定值 to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[14]),
				zap.Error(err))
			return
		}
		config.Def_value = int32(tmp)
	}

	// parse column 15 hp_value : 额外生命固定值
	if data[15] != "" {
		tmp, err = strconv.ParseInt(data[15], 10, 64)
		if err != nil {
			err = errors.New("parse field hp_value 额外生命固定值 to int32 failed")
			logger.ErrorWF("parse field hp_value 额外生命固定值 to int32 failed.",
				zap.String("xlsx", "maze_summon_v8【迷宫-召唤物】.xlsx"), zap.String("sheet", "maze_summon_v8"),
				zap.String("parse_data", data[15]),
				zap.Error(err))
			return
		}
		config.Hp_value = int32(tmp)
	}
	return
}

var gMazeSummonV8Fields = []string{
	"id",
	"name",
	"suffer_damage",
	"duration",
	"max_number",
	"model_id",
	"nor_attack_skill_id",
	"skill_id",
	"speed",
	"attack_speed_pro",
	"inheritance_attack",
	"inheritance_def",
	"inheritance_hp",
	"attack_value",
	"def_value",
	"hp_value",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeSummonV8Parser{}
	loader := &gMazeSummonV8Loader{}
	var data [][]string
	data, err = load("maze_summon_v8【迷宫-召唤物】.xlsx", "maze_summon_v8", gMazeSummonV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_summon_v8【迷宫-召唤物】.xlsx maze_summon_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_summon_v8【迷宫-召唤物】.xlsx maze_summon_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeSummonV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_summon_v8【迷宫-召唤物】.xlsx maze_summon_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeSummonV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_summon_v8【迷宫-召唤物】.xlsx maze_summon_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_summon_v8【迷宫-召唤物】.xlsx maze_summon_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_summon_v8【迷宫-召唤物】.xlsx maze_summon_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_summon_v8【迷宫-召唤物】.xlsx maze_summon_v8 data success.")
	return
}
