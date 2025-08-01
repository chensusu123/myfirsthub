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
	Skill_attr_id            int32           `json:"skill_attr_id"`            // 获得技能对应属性id
	Type                     int32           `json:"type"`                     // 技能类型
	Priority                 int32           `json:"priority"`                 // 动作技能释放优先级
	Level                    int32           `json:"level"`                    // 技能等级
	Name                     string          `json:"name"`                     // 技能名
	Desc                     string          `json:"desc"`                     // 技能描述_文本
	Auto_release_time        int32           `json:"auto_release_time"`        // 自动释放时机
	Auto_release_condition   string          `json:"auto_release_condition"`   // 自动释放条件
	Skill_cool_time          map[int32]int32 `json:"skill_cool_time"`          // 技能释放冷却时间（豪秒）
	Distance_min             map[int32]int32 `json:"distance_min"`             // 最近释放距离
	Distance_max             map[int32]int32 `json:"distance_max"`             // 最远释放距离
	Is_break                 int32           `json:"is_break"`                 // 是否打断当前动作（0-不 1-打断
	Is_no_target             int32           `json:"is_no_target"`             // 是否允许无目标释放（0-不允许 1-允许）
	Scope_type               int32           `json:"scope_type"`               // 释放目标类型
	Target_type              int32           `json:"target_type"`              // 目标阵营
	Scope_param1             map[int32]int32 `json:"scope_param1"`             // 技能效果半径
	Target_num               map[int32]int32 `json:"target_num"`               // 目标数量
	Damage_type              int32           `json:"damage_type"`              // 伤害计算类型
	Damage_element           []int32         `json:"damage_element"`           // 参与伤害计算的攻击元素类型（0-物理、1-冰、2-火、3-毒、4-电）
	Damage_element_adjust    []int32         `json:"damage_element_adjust"`    // 元素伤害系数调整关联属性id（按顺序：物理、冰、火、毒、电，万分比）
	Main_target_damage_fix   map[int32]int32 `json:"main_target_damage_fix"`   // 主目标直接伤害固定值
	Main_target_damage       map[int32]int32 `json:"main_target_damage"`       // 主目标直接伤害系数
	Second_target_damage_fix map[int32]int32 `json:"second_target_damage_fix"` // 非主目标直接伤害固定值
	Second_target_damage     map[int32]int32 `json:"second_target_damage"`     // 非主目标直接伤害系数
	Self_effect              []int32         `json:"self_effect"`              // 释放后自身效果
	Target_effect            []int32         `json:"target_effect"`            // 释放后对目标效果
	Summon_id                int32           `json:"summon_id"`                // 召唤物id（默认召唤1个）
	Duration                 map[int32]int32 `json:"duration"`                 // 技能持续时间（毫秒）
	Interval                 map[int32]int32 `json:"interval"`                 // 技能伤害或效果生效间隔（毫秒）
	Damage_adjustment        map[int32]int32 `json:"damage_adjustment"`        // 每次造成伤害后的伤害调整系数
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

	// parse column 1 skill_attr_id : 获得技能对应属性id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field skill_attr_id 获得技能对应属性id to int32 failed")
			logger.ErrorWF("parse field skill_attr_id 获得技能对应属性id to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Skill_attr_id = int32(tmp)
	}

	// parse column 2 type : 技能类型
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field type 技能类型 to int32 failed")
			logger.ErrorWF("parse field type 技能类型 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Type = int32(tmp)
	}

	// parse column 3 priority : 动作技能释放优先级
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field priority 动作技能释放优先级 to int32 failed")
			logger.ErrorWF("parse field priority 动作技能释放优先级 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Priority = int32(tmp)
	}

	// parse column 4 level : 技能等级
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field level 技能等级 to int32 failed")
			logger.ErrorWF("parse field level 技能等级 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Level = int32(tmp)
	}

	// parse column 5 name : 技能名
	if data[5] != "" {
		config.Name = data[5]
	}

	// parse column 6 desc : 技能描述_文本
	if data[6] != "" {
		config.Desc = data[6]
	}

	// parse column 7 auto_release_time : 自动释放时机
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field auto_release_time 自动释放时机 to int32 failed")
			logger.ErrorWF("parse field auto_release_time 自动释放时机 to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Auto_release_time = int32(tmp)
	}

	// parse column 8 auto_release_condition : 自动释放条件
	if data[8] != "" {
		config.Auto_release_condition = data[8]
	}

	// parse column 9 skill_cool_time : 技能释放冷却时间（豪秒）
	if data[9] != "" {

		config.Skill_cool_time = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[9], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field skill_cool_time 技能释放冷却时间（豪秒） to key int32 failed")
				logger.ErrorWF("parse map field skill_cool_time 技能释放冷却时间（豪秒） to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[9]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field skill_cool_time 技能释放冷却时间（豪秒） to value int32 failed")
				logger.ErrorWF("parse map field skill_cool_time 技能释放冷却时间（豪秒） to value int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[9]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Skill_cool_time[key] = value
		}
	}

	// parse column 10 distance_min : 最近释放距离
	if data[10] != "" {

		config.Distance_min = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[10], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field distance_min 最近释放距离 to key int32 failed")
				logger.ErrorWF("parse map field distance_min 最近释放距离 to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[10]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field distance_min 最近释放距离 to value int32 failed")
				logger.ErrorWF("parse map field distance_min 最近释放距离 to value int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[10]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Distance_min[key] = value
		}
	}

	// parse column 11 distance_max : 最远释放距离
	if data[11] != "" {

		config.Distance_max = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[11], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field distance_max 最远释放距离 to key int32 failed")
				logger.ErrorWF("parse map field distance_max 最远释放距离 to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[11]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field distance_max 最远释放距离 to value int32 failed")
				logger.ErrorWF("parse map field distance_max 最远释放距离 to value int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[11]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Distance_max[key] = value
		}
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

		config.Scope_param1 = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[16], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field scope_param1 技能效果半径 to key int32 failed")
				logger.ErrorWF("parse map field scope_param1 技能效果半径 to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[16]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field scope_param1 技能效果半径 to value int32 failed")
				logger.ErrorWF("parse map field scope_param1 技能效果半径 to value int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[16]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Scope_param1[key] = value
		}
	}

	// parse column 17 target_num : 目标数量
	if data[17] != "" {

		config.Target_num = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[17], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field target_num 目标数量 to key int32 failed")
				logger.ErrorWF("parse map field target_num 目标数量 to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[17]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field target_num 目标数量 to value int32 failed")
				logger.ErrorWF("parse map field target_num 目标数量 to value int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[17]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Target_num[key] = value
		}
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

	// parse column 19 damage_element : 参与伤害计算的攻击元素类型（0-物理、1-冰、2-火、3-毒、4-电）
	if data[19] != "" {

		vals := strings.Split(data[19], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field damage_element 参与伤害计算的攻击元素类型（0-物理、1-冰、2-火、3-毒、4-电） to []int32 failed")
				logger.ErrorWF("parse array field damage_element 参与伤害计算的攻击元素类型（0-物理、1-冰、2-火、3-毒、4-电） to []int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[19]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Damage_element = append(config.Damage_element, int32(tmp))
		}
	}

	// parse column 20 damage_element_adjust : 元素伤害系数调整关联属性id（按顺序：物理、冰、火、毒、电，万分比）
	if data[20] != "" {

		vals := strings.Split(data[20], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field damage_element_adjust 元素伤害系数调整关联属性id（按顺序：物理、冰、火、毒、电，万分比） to []int32 failed")
				logger.ErrorWF("parse array field damage_element_adjust 元素伤害系数调整关联属性id（按顺序：物理、冰、火、毒、电，万分比） to []int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[20]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Damage_element_adjust = append(config.Damage_element_adjust, int32(tmp))
		}
	}

	// parse column 21 main_target_damage_fix : 主目标直接伤害固定值
	if data[21] != "" {

		config.Main_target_damage_fix = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[21], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field main_target_damage_fix 主目标直接伤害固定值 to key int32 failed")
				logger.ErrorWF("parse map field main_target_damage_fix 主目标直接伤害固定值 to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[21]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field main_target_damage_fix 主目标直接伤害固定值 to value int32 failed")
				logger.ErrorWF("parse map field main_target_damage_fix 主目标直接伤害固定值 to value int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[21]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Main_target_damage_fix[key] = value
		}
	}

	// parse column 22 main_target_damage : 主目标直接伤害系数
	if data[22] != "" {

		config.Main_target_damage = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[22], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field main_target_damage 主目标直接伤害系数 to key int32 failed")
				logger.ErrorWF("parse map field main_target_damage 主目标直接伤害系数 to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[22]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field main_target_damage 主目标直接伤害系数 to value int32 failed")
				logger.ErrorWF("parse map field main_target_damage 主目标直接伤害系数 to value int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[22]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Main_target_damage[key] = value
		}
	}

	// parse column 23 second_target_damage_fix : 非主目标直接伤害固定值
	if data[23] != "" {

		config.Second_target_damage_fix = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[23], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field second_target_damage_fix 非主目标直接伤害固定值 to key int32 failed")
				logger.ErrorWF("parse map field second_target_damage_fix 非主目标直接伤害固定值 to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[23]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field second_target_damage_fix 非主目标直接伤害固定值 to value int32 failed")
				logger.ErrorWF("parse map field second_target_damage_fix 非主目标直接伤害固定值 to value int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[23]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Second_target_damage_fix[key] = value
		}
	}

	// parse column 24 second_target_damage : 非主目标直接伤害系数
	if data[24] != "" {

		config.Second_target_damage = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[24], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field second_target_damage 非主目标直接伤害系数 to key int32 failed")
				logger.ErrorWF("parse map field second_target_damage 非主目标直接伤害系数 to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[24]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field second_target_damage 非主目标直接伤害系数 to value int32 failed")
				logger.ErrorWF("parse map field second_target_damage 非主目标直接伤害系数 to value int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[24]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Second_target_damage[key] = value
		}
	}

	// parse column 25 self_effect : 释放后自身效果
	if data[25] != "" {

		vals := strings.Split(data[25], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field self_effect 释放后自身效果 to []int32 failed")
				logger.ErrorWF("parse array field self_effect 释放后自身效果 to []int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[25]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Self_effect = append(config.Self_effect, int32(tmp))
		}
	}

	// parse column 26 target_effect : 释放后对目标效果
	if data[26] != "" {

		vals := strings.Split(data[26], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field target_effect 释放后对目标效果 to []int32 failed")
				logger.ErrorWF("parse array field target_effect 释放后对目标效果 to []int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[26]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Target_effect = append(config.Target_effect, int32(tmp))
		}
	}

	// parse column 27 summon_id : 召唤物id（默认召唤1个）
	if data[27] != "" {
		tmp, err = strconv.ParseInt(data[27], 10, 64)
		if err != nil {
			err = errors.New("parse field summon_id 召唤物id（默认召唤1个） to int32 failed")
			logger.ErrorWF("parse field summon_id 召唤物id（默认召唤1个） to int32 failed.",
				zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
				zap.String("parse_data", data[27]),
				zap.Error(err))
			return
		}
		config.Summon_id = int32(tmp)
	}

	// parse column 28 duration : 技能持续时间（毫秒）
	if data[28] != "" {

		config.Duration = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[28], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field duration 技能持续时间（毫秒） to key int32 failed")
				logger.ErrorWF("parse map field duration 技能持续时间（毫秒） to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[28]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field duration 技能持续时间（毫秒） to value int32 failed")
				logger.ErrorWF("parse map field duration 技能持续时间（毫秒） to value int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[28]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Duration[key] = value
		}
	}

	// parse column 29 interval : 技能伤害或效果生效间隔（毫秒）
	if data[29] != "" {

		config.Interval = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[29], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field interval 技能伤害或效果生效间隔（毫秒） to key int32 failed")
				logger.ErrorWF("parse map field interval 技能伤害或效果生效间隔（毫秒） to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[29]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field interval 技能伤害或效果生效间隔（毫秒） to value int32 failed")
				logger.ErrorWF("parse map field interval 技能伤害或效果生效间隔（毫秒） to value int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[29]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Interval[key] = value
		}
	}

	// parse column 30 damage_adjustment : 每次造成伤害后的伤害调整系数
	if data[30] != "" {

		config.Damage_adjustment = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[30], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field damage_adjustment 每次造成伤害后的伤害调整系数 to key int32 failed")
				logger.ErrorWF("parse map field damage_adjustment 每次造成伤害后的伤害调整系数 to key int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[30]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field damage_adjustment 每次造成伤害后的伤害调整系数 to value int32 failed")
				logger.ErrorWF("parse map field damage_adjustment 每次造成伤害后的伤害调整系数 to value int32 failed.",
					zap.String("xlsx", "maze_skill_info_v8【迷宫-技能-技能信息】.xlsx"), zap.String("sheet", "maze_skill_info_v8"),
					// zap.String("field_data",data[30]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Damage_adjustment[key] = value
		}
	}
	return
}

var gMazeSkillInfoV8Fields = []string{
	"id",
	"skill_attr_id",
	"type",
	"priority",
	"level",
	"name",
	"desc",
	"auto_release_time",
	"auto_release_condition",
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
	"damage_element",
	"damage_element_adjust",
	"main_target_damage_fix",
	"main_target_damage",
	"second_target_damage_fix",
	"second_target_damage",
	"self_effect",
	"target_effect",
	"summon_id",
	"duration",
	"interval",
	"damage_adjustment",
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
