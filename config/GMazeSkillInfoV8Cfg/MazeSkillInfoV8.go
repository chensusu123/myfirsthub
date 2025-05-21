package GMazeSkillInfoV8Cfg

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

// MazeSkillInfoV8ConfigRow from maze_skill_info_v8【迷宫-技能-技能信息】.xlsx maze_skill_info_v8
type MazeSkillInfoV8ConfigRow struct {
	Id                       int32           `json:"id"`                       // 技能id
	Type                     int32           `json:"type"`                     // 技能类型
	Group                    int32           `json:"group"`                    // 技能组id
	Level                    int32           `json:"level"`                    // 技能等级
	Name                     string          `json:"name"`                     // 技能名
	Cost                     map[int32]int64 `json:"cost"`                     // 消耗物品：数量
	Desc                     string          `json:"desc"`                     // 技能描述_文本
	Initial_cool_time        int32           `json:"initial_cool_time"`        // 技能初始冷却时间（豪秒）
	Public_cool_time         int32           `json:"public_cool_time"`         // 技能公共冷却时间（豪秒）
	Skill_cool_time          int32           `json:"skill_cool_time"`          // 技能释放冷却时间（豪秒）
	Distance_min             int32           `json:"distance_min"`             // 最近释放距离
	Distance_max             int32           `json:"distance_max"`             // 最远释放距离
	Is_break                 int32           `json:"is_break"`                 // 是否打断当前动作（0-不 1-打断
	Is_no_target             int32           `json:"is_no_target"`             // 是否允许无目标释放（0-不允许 1-允许）
	Scope_type               int32           `json:"scope_type"`               // 释放目标类型
	Target_type              int32           `json:"target_type"`              // 目标阵营
	Scope_param1             int32           `json:"scope_param1"`             // 技能效果半径
	Target_num               int32           `json:"target_num"`               // 目标数量
	Damage_type              int32           `json:"damage_type"`              // 伤害计算类型
	Main_target_damage_fix   int32           `json:"main_target_damage_fix"`   // 主目标直接伤害固定值
	Main_target_damage       int32           `json:"main_target_damage"`       // 主目标直接伤害系数
	Second_target_damage_fix int32           `json:"second_target_damage_fix"` // 非主目标直接伤害固定值
	Second_target_damage     int32           `json:"second_target_damage"`     // 非主目标直接伤害系数
	Self_effect              []int32         `json:"self_effect"`              // 释放后自身效果
	Target_effect            []int32         `json:"target_effect"`            // 释放后对目标效果
	Target_effect_pro        map[int32]int32 `json:"target_effect_pro"`        // 效果同组几率
	Is_allow                 []int32         `json:"is_allow"`                 // 允许释放状态
	Is_target                []int32         `json:"is_target"`                // 能被选中的目标状态
}

// MazeSkillInfoV8Config from maze_skill_info_v8【迷宫-技能-技能信息】.xlsx maze_skill_info_v8
type MazeSkillInfoV8Config struct {
	ConfigRows map[int32]*MazeSkillInfoV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeSkillInfoV8Config {
	ret := &MazeSkillInfoV8Config{ConfigRows: map[int32]*MazeSkillInfoV8ConfigRow{}}
	return ret
}

// GetMazeSkillInfoV8Config get one config by configId
func (c *MazeSkillInfoV8Config) GetMazeSkillInfoV8Config(configId int32) *MazeSkillInfoV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeSkillInfoV8Config) Get(configId int32) *MazeSkillInfoV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeSkillInfoV8Config get all config slice
func (c *MazeSkillInfoV8Config) GetAllMazeSkillInfoV8Config() (res []*MazeSkillInfoV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeSkillInfoV8Config) GetAll() (res []*MazeSkillInfoV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeSkillInfoV8Config

// GetMazeSkillInfoV8Config pkg func. get one config by configId
func GetMazeSkillInfoV8Config(configId int32) *MazeSkillInfoV8ConfigRow {
	return gConfigData.GetMazeSkillInfoV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeSkillInfoV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeSkillInfoV8Config pkg func. get all config slice
func GetAllMazeSkillInfoV8Config() []*MazeSkillInfoV8ConfigRow {
	return gConfigData.GetAllMazeSkillInfoV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeSkillInfoV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeSkillInfoV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeSkillInfoV8ConfigRow from maze_skill_info_v8【迷宫-技能-技能信息】.xlsx maze_skill_info_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeSkillInfoV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_skill_info_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_skill_info_v8.json",
		"maze_skill_info_v8【迷宫-技能-技能信息】.xlsx", "maze_skill_info_v8",
		&gMazeSkillInfoV8Parser{}, &gMazeSkillInfoV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeSkillInfoV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeSkillInfoV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeSkillInfoV8Config))(c)
		return true
	})
}

// RegisterMazeSkillInfoV8InitCallBack reg config update func (old func)
var RegisterMazeSkillInfoV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeSkillInfoV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeSkillInfoV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeSkillInfoV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeSkillInfoV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeSkillInfoV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeSkillInfoV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeSkillInfoV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeSkillInfoV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeSkillInfoV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeSkillInfoV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeSkillInfoV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeSkillInfoV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeSkillInfoV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillInfoV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkillInfoV8ConfigRow", zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"),
			zap.String("sheet", "maze_skill_info_v8"))
		return
	}
	config, ok := container.(*MazeSkillInfoV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillInfoV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillInfoV8Config", zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"),
			zap.String("sheet", "maze_skill_info_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeSkillInfoV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeSkillInfoV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillInfoV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillInfoV8Config", zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"),
			zap.String("sheet", "maze_skill_info_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeSkillInfoV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeSkillInfoV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillInfoV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillInfoV8Config", zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"),
			zap.String("sheet", "maze_skill_info_v8"))
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
type gMazeSkillInfoV8Parser struct {
}

// New new config row data
func (*gMazeSkillInfoV8Parser) New() interface{} {
	return &MazeSkillInfoV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeSkillInfoV8Parser) Fields() []string {
	return gMazeSkillInfoV8Fields
}

// Parse parse raw data to row data
func (*gMazeSkillInfoV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeSkillInfoV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillInfoV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkillInfoV8ConfigRow", zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"),
			zap.String("sheet", "maze_skill_info_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeSkillInfoV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeSkillInfoV8ConfigRow",
			zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"),
			zap.String("sheet", "maze_skill_info_v8"), zap.Int("need_count", len(gMazeSkillInfoV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 技能id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 技能id to int32 failed")
			logger.ErrorWF("parse field id 技能id to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 type : 技能类型
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field type 技能类型 to int32 failed")
			logger.ErrorWF("parse field type 技能类型 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Type = int32(tmp)
	}

	// parse column 2 group : 技能组id
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field group 技能组id to int32 failed")
			logger.ErrorWF("parse field group 技能组id to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Group = int32(tmp)
	}

	// parse column 3 level : 技能等级
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field level 技能等级 to int32 failed")
			logger.ErrorWF("parse field level 技能等级 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Level = int32(tmp)
	}

	// parse column 4 name : 技能名
	if data[4] != "" {
		config.Name = data[4]
	}

	// parse column 5 cost : 消耗物品：数量
	if data[5] != "" {

		config.Cost = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[5], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field cost 消耗物品：数量 to key int32 failed")
				logger.ErrorWF("parse map field cost 消耗物品：数量 to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field cost 消耗物品：数量 to value int64 failed")
				logger.ErrorWF("parse map field cost 消耗物品：数量 to value int64 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Cost[key] = value
		}
	}

	// parse column 6 desc : 技能描述_文本
	if data[6] != "" {
		config.Desc = data[6]
	}

	// parse column 7 initial_cool_time : 技能初始冷却时间（豪秒）
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field initial_cool_time 技能初始冷却时间（豪秒） to int32 failed")
			logger.ErrorWF("parse field initial_cool_time 技能初始冷却时间（豪秒） to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Initial_cool_time = int32(tmp)
	}

	// parse column 8 public_cool_time : 技能公共冷却时间（豪秒）
	if data[8] != "" {
		tmp, err = strconv.ParseInt(data[8], 10, 64)
		if err != nil {
			err = errors.New("parse field public_cool_time 技能公共冷却时间（豪秒） to int32 failed")
			logger.ErrorWF("parse field public_cool_time 技能公共冷却时间（豪秒） to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[8]),
				zap.Error(err))
			return
		}
		config.Public_cool_time = int32(tmp)
	}

	// parse column 9 skill_cool_time : 技能释放冷却时间（豪秒）
	if data[9] != "" {
		tmp, err = strconv.ParseInt(data[9], 10, 64)
		if err != nil {
			err = errors.New("parse field skill_cool_time 技能释放冷却时间（豪秒） to int32 failed")
			logger.ErrorWF("parse field skill_cool_time 技能释放冷却时间（豪秒） to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[9]),
				zap.Error(err))
			return
		}
		config.Skill_cool_time = int32(tmp)
	}

	// parse column 10 distance_min : 最近释放距离
	if data[10] != "" {
		tmp, err = strconv.ParseInt(data[10], 10, 64)
		if err != nil {
			err = errors.New("parse field distance_min 最近释放距离 to int32 failed")
			logger.ErrorWF("parse field distance_min 最近释放距离 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[10]),
				zap.Error(err))
			return
		}
		config.Distance_min = int32(tmp)
	}

	// parse column 11 distance_max : 最远释放距离
	if data[11] != "" {
		tmp, err = strconv.ParseInt(data[11], 10, 64)
		if err != nil {
			err = errors.New("parse field distance_max 最远释放距离 to int32 failed")
			logger.ErrorWF("parse field distance_max 最远释放距离 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[11]),
				zap.Error(err))
			return
		}
		config.Distance_max = int32(tmp)
	}

	// parse column 12 is_break : 是否打断当前动作（0-不 1-打断
	if data[12] != "" {
		tmp, err = strconv.ParseInt(data[12], 10, 64)
		if err != nil {
			err = errors.New("parse field is_break 是否打断当前动作（0-不 1-打断 to int32 failed")
			logger.ErrorWF("parse field is_break 是否打断当前动作（0-不 1-打断 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[12]),
				zap.Error(err))
			return
		}
		config.Is_break = int32(tmp)
	}

	// parse column 13 is_no_target : 是否允许无目标释放（0-不允许 1-允许）
	if data[13] != "" {
		tmp, err = strconv.ParseInt(data[13], 10, 64)
		if err != nil {
			err = errors.New("parse field is_no_target 是否允许无目标释放（0-不允许 1-允许） to int32 failed")
			logger.ErrorWF("parse field is_no_target 是否允许无目标释放（0-不允许 1-允许） to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[13]),
				zap.Error(err))
			return
		}
		config.Is_no_target = int32(tmp)
	}

	// parse column 14 scope_type : 释放目标类型
	if data[14] != "" {
		tmp, err = strconv.ParseInt(data[14], 10, 64)
		if err != nil {
			err = errors.New("parse field scope_type 释放目标类型 to int32 failed")
			logger.ErrorWF("parse field scope_type 释放目标类型 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[14]),
				zap.Error(err))
			return
		}
		config.Scope_type = int32(tmp)
	}

	// parse column 15 target_type : 目标阵营
	if data[15] != "" {
		tmp, err = strconv.ParseInt(data[15], 10, 64)
		if err != nil {
			err = errors.New("parse field target_type 目标阵营 to int32 failed")
			logger.ErrorWF("parse field target_type 目标阵营 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[15]),
				zap.Error(err))
			return
		}
		config.Target_type = int32(tmp)
	}

	// parse column 16 scope_param1 : 技能效果半径
	if data[16] != "" {
		tmp, err = strconv.ParseInt(data[16], 10, 64)
		if err != nil {
			err = errors.New("parse field scope_param1 技能效果半径 to int32 failed")
			logger.ErrorWF("parse field scope_param1 技能效果半径 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[16]),
				zap.Error(err))
			return
		}
		config.Scope_param1 = int32(tmp)
	}

	// parse column 17 target_num : 目标数量
	if data[17] != "" {
		tmp, err = strconv.ParseInt(data[17], 10, 64)
		if err != nil {
			err = errors.New("parse field target_num 目标数量 to int32 failed")
			logger.ErrorWF("parse field target_num 目标数量 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[17]),
				zap.Error(err))
			return
		}
		config.Target_num = int32(tmp)
	}

	// parse column 18 damage_type : 伤害计算类型
	if data[18] != "" {
		tmp, err = strconv.ParseInt(data[18], 10, 64)
		if err != nil {
			err = errors.New("parse field damage_type 伤害计算类型 to int32 failed")
			logger.ErrorWF("parse field damage_type 伤害计算类型 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[18]),
				zap.Error(err))
			return
		}
		config.Damage_type = int32(tmp)
	}

	// parse column 19 main_target_damage_fix : 主目标直接伤害固定值
	if data[19] != "" {
		tmp, err = strconv.ParseInt(data[19], 10, 64)
		if err != nil {
			err = errors.New("parse field main_target_damage_fix 主目标直接伤害固定值 to int32 failed")
			logger.ErrorWF("parse field main_target_damage_fix 主目标直接伤害固定值 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[19]),
				zap.Error(err))
			return
		}
		config.Main_target_damage_fix = int32(tmp)
	}

	// parse column 20 main_target_damage : 主目标直接伤害系数
	if data[20] != "" {
		tmp, err = strconv.ParseInt(data[20], 10, 64)
		if err != nil {
			err = errors.New("parse field main_target_damage 主目标直接伤害系数 to int32 failed")
			logger.ErrorWF("parse field main_target_damage 主目标直接伤害系数 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[20]),
				zap.Error(err))
			return
		}
		config.Main_target_damage = int32(tmp)
	}

	// parse column 21 second_target_damage_fix : 非主目标直接伤害固定值
	if data[21] != "" {
		tmp, err = strconv.ParseInt(data[21], 10, 64)
		if err != nil {
			err = errors.New("parse field second_target_damage_fix 非主目标直接伤害固定值 to int32 failed")
			logger.ErrorWF("parse field second_target_damage_fix 非主目标直接伤害固定值 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[21]),
				zap.Error(err))
			return
		}
		config.Second_target_damage_fix = int32(tmp)
	}

	// parse column 22 second_target_damage : 非主目标直接伤害系数
	if data[22] != "" {
		tmp, err = strconv.ParseInt(data[22], 10, 64)
		if err != nil {
			err = errors.New("parse field second_target_damage 非主目标直接伤害系数 to int32 failed")
			logger.ErrorWF("parse field second_target_damage 非主目标直接伤害系数 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[22]),
				zap.Error(err))
			return
		}
		config.Second_target_damage = int32(tmp)
	}

	// parse column 23 self_effect : 释放后自身效果
	if data[23] != "" {

		vals := strings.Split(data[23], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field self_effect 释放后自身效果 to []int32 failed")
				logger.ErrorWF("parse array field self_effect 释放后自身效果 to []int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[23]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Self_effect = append(config.Self_effect, int32(tmp))
		}
	}

	// parse column 24 target_effect : 释放后对目标效果
	if data[24] != "" {

		vals := strings.Split(data[24], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field target_effect 释放后对目标效果 to []int32 failed")
				logger.ErrorWF("parse array field target_effect 释放后对目标效果 to []int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[24]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Target_effect = append(config.Target_effect, int32(tmp))
		}
	}

	// parse column 25 target_effect_pro : 效果同组几率
	if data[25] != "" {

		config.Target_effect_pro = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[25], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field target_effect_pro 效果同组几率 to key int32 failed")
				logger.ErrorWF("parse map field target_effect_pro 效果同组几率 to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[25]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field target_effect_pro 效果同组几率 to value int32 failed")
				logger.ErrorWF("parse map field target_effect_pro 效果同组几率 to value int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[25]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Target_effect_pro[key] = value
		}
	}

	// parse column 26 is_allow : 允许释放状态
	if data[26] != "" {

		vals := strings.Split(data[26], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field is_allow 允许释放状态 to []int32 failed")
				logger.ErrorWF("parse array field is_allow 允许释放状态 to []int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[26]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Is_allow = append(config.Is_allow, int32(tmp))
		}
	}

	// parse column 27 is_target : 能被选中的目标状态
	if data[27] != "" {

		vals := strings.Split(data[27], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field is_target 能被选中的目标状态 to []int32 failed")
				logger.ErrorWF("parse array field is_target 能被选中的目标状态 to []int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[27]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Is_target = append(config.Is_target, int32(tmp))
		}
	}
	return
}

var gMazeSkillInfoV8Fields = []string{
	"id",
	"type",
	"group",
	"level",
	"name",
	"cost",
	"desc",
	"initial_cool_time",
	"public_cool_time",
	"skill_cool_time",
	"distance_min",
	"distance_max",
	"is_break",
	"is_no_target",
	"scope_type",
	"target_type",
	"scope_param1",
	"target_num",
	"damage_type",
	"main_target_damage_fix",
	"main_target_damage",
	"second_target_damage_fix",
	"second_target_damage",
	"self_effect",
	"target_effect",
	"target_effect_pro",
	"is_allow",
	"is_target",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeSkillInfoV8Parser{}
	loader := &gMazeSkillInfoV8Loader{}
	var data [][]string
	data, err = load("maze_skill_info_v8【迷宫-技能-技能信息】.xlsx", "maze_skill_info_v8", gMazeSkillInfoV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_skill_info_v8【迷宫-技能-技能信息】.xlsx maze_skill_info_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_skill_info_v8【迷宫-技能-技能信息】.xlsx maze_skill_info_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeSkillInfoV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_skill_info_v8【迷宫-技能-技能信息】.xlsx maze_skill_info_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeSkillInfoV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_skill_info_v8【迷宫-技能-技能信息】.xlsx maze_skill_info_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_skill_info_v8【迷宫-技能-技能信息】.xlsx maze_skill_info_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_skill_info_v8【迷宫-技能-技能信息】.xlsx maze_skill_info_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_skill_info_v8【迷宫-技能-技能信息】.xlsx maze_skill_info_v8 data success.")
	return
}
