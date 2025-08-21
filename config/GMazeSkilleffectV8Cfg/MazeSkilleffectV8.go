package GMazeSkilleffectV8Cfg

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

// MazeSkilleffectV8ConfigRow from maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx maze_skilleffect_v8
type MazeSkilleffectV8ConfigRow struct {
	Effect__id                    int32           `json:"effect__id"`                    // effect_id
	Effect_group                  int32           `json:"effect_group"`                  // effect组id
	In_group_weight               int32           `json:"in_group_weight"`               // 同组优先级
	Cool_down                     int32           `json:"cool_down"`                     // 重复中buff组的cd时间（毫毫秒）
	Attr                          int32           `json:"attr"`                          // 逻辑属性id
	Attr_value_variable_id        map[int32]int32 `json:"attr_value_variable_id"`        // 参数1关联的变量id
	Attr_value_type               int32           `json:"attr_value_type"`               // 参数1数值类型
	Attr_value                    int32           `json:"attr_value"`                    // 参数1
	Attr_value_2_variable_id      map[int32]int32 `json:"attr_value_2_variable_id"`      // 参数2关联的变量id
	Attr_value_2_type             int32           `json:"attr_value_2_type"`             // 参数2数值类型
	Attr_value_2                  int32           `json:"attr_value_2"`                  // 参数2
	Attr_value_3_variable_id      map[int32]int32 `json:"attr_value_3_variable_id"`      // 参数3关联的变量id
	Attr_value_3_type             int32           `json:"attr_value_3_type"`             // 参数3数值类型
	Attr_value_3                  int32           `json:"attr_value_3"`                  // 参数3
	Attr_value_4                  []int32         `json:"attr_value_4"`                  // 参数4
	Last_time_variable_id         map[int32]int32 `json:"last_time_variable_id"`         // 持续时长（毫秒）关联的变量id
	Last_time                     int32           `json:"last_time"`                     // 持续时长（毫秒）
	Base_hitrate_variable_id      map[int32]int32 `json:"base_hitrate_variable_id"`      // 基础命中率变量
	Base_hitrate                  int32           `json:"base_hitrate"`                  // 基础命中率（万分比）
	Attr_value_7_variable_id      map[int32]int32 `json:"attr_value_7_variable_id"`      // 参数7关联的变量id
	Attr_value_7_type             int32           `json:"attr_value_7_type"`             // 参数7数值类型
	Attr_value_7                  int32           `json:"attr_value_7"`                  // 参数7
	Attr_value_8_variable_id      map[int32]int32 `json:"attr_value_8_variable_id"`      // 参数8关联的变量id
	Attr_value_8_type             int32           `json:"attr_value_8_type"`             // 参数8数值类型
	Attr_value_8                  int32           `json:"attr_value_8"`                  // 参数8结算间隔时间
	Modify_attr_value_variable_id map[int32]int32 `json:"modify_attr_value_variable_id"` // 修改属性关联变量id
	Modify_attr_value_type        map[int32]int32 `json:"modify_attr_value_type"`        // 修改属性数值类型
	Modify_attr_value_attr_id     map[int32]int32 `json:"modify_attr_value_attr_id"`     // 修改属性影响属性id
	Modify_attr_value             map[int32]int32 `json:"modify_attr_value"`             // 修改属性影响数值
}

// MazeSkilleffectV8Config from maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx maze_skilleffect_v8
type MazeSkilleffectV8Config struct {
	ConfigRows map[int32]*MazeSkilleffectV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeSkilleffectV8Config {
	ret := &MazeSkilleffectV8Config{ConfigRows: map[int32]*MazeSkilleffectV8ConfigRow{}}
	return ret
}

// GetMazeSkilleffectV8Config get one config by configId
func (c *MazeSkilleffectV8Config) GetMazeSkilleffectV8Config(configId int32) *MazeSkilleffectV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeSkilleffectV8Config) Get(configId int32) *MazeSkilleffectV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeSkilleffectV8Config get all config slice
func (c *MazeSkilleffectV8Config) GetAllMazeSkilleffectV8Config() (res []*MazeSkilleffectV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeSkilleffectV8Config) GetAll() (res []*MazeSkilleffectV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeSkilleffectV8Config

// GetMazeSkilleffectV8Config pkg func. get one config by configId
func GetMazeSkilleffectV8Config(configId int32) *MazeSkilleffectV8ConfigRow {
	return gConfigData.GetMazeSkilleffectV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeSkilleffectV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeSkilleffectV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_skilleffect_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeSkilleffectV8Config pkg func. get all config slice
func GetAllMazeSkilleffectV8Config() []*MazeSkilleffectV8ConfigRow {
	return gConfigData.GetAllMazeSkilleffectV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeSkilleffectV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeSkilleffectV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeSkilleffectV8ConfigRow from maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx maze_skilleffect_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeSkilleffectV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_skilleffect_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_skilleffect_v8.json",
		"maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx", "maze_skilleffect_v8",
		&gMazeSkilleffectV8Parser{}, &gMazeSkilleffectV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeSkilleffectV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeSkilleffectV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeSkilleffectV8Config))(c)
		return true
	})
}

// RegisterMazeSkilleffectV8InitCallBack reg config update func (old func)
var RegisterMazeSkilleffectV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeSkilleffectV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeSkilleffectV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeSkilleffectV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeSkilleffectV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeSkilleffectV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeSkilleffectV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeSkilleffectV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeSkilleffectV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeSkilleffectV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeSkilleffectV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeSkilleffectV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeSkilleffectV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeSkilleffectV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkilleffectV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkilleffectV8ConfigRow", zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"),
			zap.String("sheet", "maze_skilleffect_v8"))
		return
	}
	config, ok := container.(*MazeSkilleffectV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkilleffectV8Config")
		logger.ErrorWF("invalid type. not *MazeSkilleffectV8Config", zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"),
			zap.String("sheet", "maze_skilleffect_v8"))
		return
	}
	config.ConfigRows[row.Effect__id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeSkilleffectV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeSkilleffectV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkilleffectV8Config")
		logger.ErrorWF("invalid type. not *MazeSkilleffectV8Config", zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"),
			zap.String("sheet", "maze_skilleffect_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeSkilleffectV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeSkilleffectV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkilleffectV8Config")
		logger.ErrorWF("invalid type. not *MazeSkilleffectV8Config", zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"),
			zap.String("sheet", "maze_skilleffect_v8"))
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
type gMazeSkilleffectV8Parser struct {
}

// New new config row data
func (*gMazeSkilleffectV8Parser) New() interface{} {
	return &MazeSkilleffectV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeSkilleffectV8Parser) Fields() []string {
	return gMazeSkilleffectV8Fields
}

// Parse parse raw data to row data
func (*gMazeSkilleffectV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeSkilleffectV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkilleffectV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkilleffectV8ConfigRow", zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"),
			zap.String("sheet", "maze_skilleffect_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeSkilleffectV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeSkilleffectV8ConfigRow",
			zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"),
			zap.String("sheet", "maze_skilleffect_v8"), zap.Int("need_count", len(gMazeSkilleffectV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 effect__id : effect_id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field effect__id effect_id to int32 failed")
			logger.ErrorWF("parse field effect__id effect_id to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Effect__id = int32(tmp)
	}

	// parse column 1 effect_group : effect组id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field effect_group effect组id to int32 failed")
			logger.ErrorWF("parse field effect_group effect组id to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Effect_group = int32(tmp)
	}

	// parse column 2 in_group_weight : 同组优先级
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field in_group_weight 同组优先级 to int32 failed")
			logger.ErrorWF("parse field in_group_weight 同组优先级 to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.In_group_weight = int32(tmp)
	}

	// parse column 3 cool_down : 重复中buff组的cd时间（毫毫秒）
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field cool_down 重复中buff组的cd时间（毫毫秒） to int32 failed")
			logger.ErrorWF("parse field cool_down 重复中buff组的cd时间（毫毫秒） to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Cool_down = int32(tmp)
	}

	// parse column 4 attr : 逻辑属性id
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field attr 逻辑属性id to int32 failed")
			logger.ErrorWF("parse field attr 逻辑属性id to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Attr = int32(tmp)
	}

	// parse column 5 attr_value_variable_id : 参数1关联的变量id
	if data[5] != "" {

		config.Attr_value_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[5], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr_value_variable_id 参数1关联的变量id to key int32 failed")
				logger.ErrorWF("parse map field attr_value_variable_id 参数1关联的变量id to key int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr_value_variable_id 参数1关联的变量id to value int32 failed")
				logger.ErrorWF("parse map field attr_value_variable_id 参数1关联的变量id to value int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Attr_value_variable_id[key] = value
		}
	}

	// parse column 6 attr_value_type : 参数1数值类型
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field attr_value_type 参数1数值类型 to int32 failed")
			logger.ErrorWF("parse field attr_value_type 参数1数值类型 to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Attr_value_type = int32(tmp)
	}

	// parse column 7 attr_value : 参数1
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field attr_value 参数1 to int32 failed")
			logger.ErrorWF("parse field attr_value 参数1 to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Attr_value = int32(tmp)
	}

	// parse column 8 attr_value_2_variable_id : 参数2关联的变量id
	if data[8] != "" {

		config.Attr_value_2_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[8], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr_value_2_variable_id 参数2关联的变量id to key int32 failed")
				logger.ErrorWF("parse map field attr_value_2_variable_id 参数2关联的变量id to key int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[8]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr_value_2_variable_id 参数2关联的变量id to value int32 failed")
				logger.ErrorWF("parse map field attr_value_2_variable_id 参数2关联的变量id to value int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[8]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Attr_value_2_variable_id[key] = value
		}
	}

	// parse column 9 attr_value_2_type : 参数2数值类型
	if data[9] != "" {
		tmp, err = strconv.ParseInt(data[9], 10, 64)
		if err != nil {
			err = errors.New("parse field attr_value_2_type 参数2数值类型 to int32 failed")
			logger.ErrorWF("parse field attr_value_2_type 参数2数值类型 to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[9]),
				zap.Error(err))
			return
		}
		config.Attr_value_2_type = int32(tmp)
	}

	// parse column 10 attr_value_2 : 参数2
	if data[10] != "" {
		tmp, err = strconv.ParseInt(data[10], 10, 64)
		if err != nil {
			err = errors.New("parse field attr_value_2 参数2 to int32 failed")
			logger.ErrorWF("parse field attr_value_2 参数2 to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[10]),
				zap.Error(err))
			return
		}
		config.Attr_value_2 = int32(tmp)
	}

	// parse column 11 attr_value_3_variable_id : 参数3关联的变量id
	if data[11] != "" {

		config.Attr_value_3_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[11], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr_value_3_variable_id 参数3关联的变量id to key int32 failed")
				logger.ErrorWF("parse map field attr_value_3_variable_id 参数3关联的变量id to key int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[11]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr_value_3_variable_id 参数3关联的变量id to value int32 failed")
				logger.ErrorWF("parse map field attr_value_3_variable_id 参数3关联的变量id to value int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[11]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Attr_value_3_variable_id[key] = value
		}
	}

	// parse column 12 attr_value_3_type : 参数3数值类型
	if data[12] != "" {
		tmp, err = strconv.ParseInt(data[12], 10, 64)
		if err != nil {
			err = errors.New("parse field attr_value_3_type 参数3数值类型 to int32 failed")
			logger.ErrorWF("parse field attr_value_3_type 参数3数值类型 to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[12]),
				zap.Error(err))
			return
		}
		config.Attr_value_3_type = int32(tmp)
	}

	// parse column 13 attr_value_3 : 参数3
	if data[13] != "" {
		tmp, err = strconv.ParseInt(data[13], 10, 64)
		if err != nil {
			err = errors.New("parse field attr_value_3 参数3 to int32 failed")
			logger.ErrorWF("parse field attr_value_3 参数3 to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[13]),
				zap.Error(err))
			return
		}
		config.Attr_value_3 = int32(tmp)
	}

	// parse column 14 attr_value_4 : 参数4
	if data[14] != "" {

		vals := strings.Split(data[14], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field attr_value_4 参数4 to []int32 failed")
				logger.ErrorWF("parse array field attr_value_4 参数4 to []int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[14]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Attr_value_4 = append(config.Attr_value_4, int32(tmp))
		}
	}

	// parse column 15 last_time_variable_id : 持续时长（毫秒）关联的变量id
	if data[15] != "" {

		config.Last_time_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[15], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field last_time_variable_id 持续时长（毫秒）关联的变量id to key int32 failed")
				logger.ErrorWF("parse map field last_time_variable_id 持续时长（毫秒）关联的变量id to key int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[15]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field last_time_variable_id 持续时长（毫秒）关联的变量id to value int32 failed")
				logger.ErrorWF("parse map field last_time_variable_id 持续时长（毫秒）关联的变量id to value int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[15]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Last_time_variable_id[key] = value
		}
	}

	// parse column 16 last_time : 持续时长（毫秒）
	if data[16] != "" {
		tmp, err = strconv.ParseInt(data[16], 10, 64)
		if err != nil {
			err = errors.New("parse field last_time 持续时长（毫秒） to int32 failed")
			logger.ErrorWF("parse field last_time 持续时长（毫秒） to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[16]),
				zap.Error(err))
			return
		}
		config.Last_time = int32(tmp)
	}

	// parse column 17 base_hitrate_variable_id : 基础命中率变量
	if data[17] != "" {

		config.Base_hitrate_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[17], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field base_hitrate_variable_id 基础命中率变量 to key int32 failed")
				logger.ErrorWF("parse map field base_hitrate_variable_id 基础命中率变量 to key int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[17]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field base_hitrate_variable_id 基础命中率变量 to value int32 failed")
				logger.ErrorWF("parse map field base_hitrate_variable_id 基础命中率变量 to value int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[17]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Base_hitrate_variable_id[key] = value
		}
	}

	// parse column 18 base_hitrate : 基础命中率（万分比）
	if data[18] != "" {
		tmp, err = strconv.ParseInt(data[18], 10, 64)
		if err != nil {
			err = errors.New("parse field base_hitrate 基础命中率（万分比） to int32 failed")
			logger.ErrorWF("parse field base_hitrate 基础命中率（万分比） to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[18]),
				zap.Error(err))
			return
		}
		config.Base_hitrate = int32(tmp)
	}

	// parse column 19 attr_value_7_variable_id : 参数7关联的变量id
	if data[19] != "" {

		config.Attr_value_7_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[19], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr_value_7_variable_id 参数7关联的变量id to key int32 failed")
				logger.ErrorWF("parse map field attr_value_7_variable_id 参数7关联的变量id to key int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[19]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr_value_7_variable_id 参数7关联的变量id to value int32 failed")
				logger.ErrorWF("parse map field attr_value_7_variable_id 参数7关联的变量id to value int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[19]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Attr_value_7_variable_id[key] = value
		}
	}

	// parse column 20 attr_value_7_type : 参数7数值类型
	if data[20] != "" {
		tmp, err = strconv.ParseInt(data[20], 10, 64)
		if err != nil {
			err = errors.New("parse field attr_value_7_type 参数7数值类型 to int32 failed")
			logger.ErrorWF("parse field attr_value_7_type 参数7数值类型 to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[20]),
				zap.Error(err))
			return
		}
		config.Attr_value_7_type = int32(tmp)
	}

	// parse column 21 attr_value_7 : 参数7
	if data[21] != "" {
		tmp, err = strconv.ParseInt(data[21], 10, 64)
		if err != nil {
			err = errors.New("parse field attr_value_7 参数7 to int32 failed")
			logger.ErrorWF("parse field attr_value_7 参数7 to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[21]),
				zap.Error(err))
			return
		}
		config.Attr_value_7 = int32(tmp)
	}

	// parse column 22 attr_value_8_variable_id : 参数8关联的变量id
	if data[22] != "" {

		config.Attr_value_8_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[22], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr_value_8_variable_id 参数8关联的变量id to key int32 failed")
				logger.ErrorWF("parse map field attr_value_8_variable_id 参数8关联的变量id to key int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[22]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr_value_8_variable_id 参数8关联的变量id to value int32 failed")
				logger.ErrorWF("parse map field attr_value_8_variable_id 参数8关联的变量id to value int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[22]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Attr_value_8_variable_id[key] = value
		}
	}

	// parse column 23 attr_value_8_type : 参数8数值类型
	if data[23] != "" {
		tmp, err = strconv.ParseInt(data[23], 10, 64)
		if err != nil {
			err = errors.New("parse field attr_value_8_type 参数8数值类型 to int32 failed")
			logger.ErrorWF("parse field attr_value_8_type 参数8数值类型 to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[23]),
				zap.Error(err))
			return
		}
		config.Attr_value_8_type = int32(tmp)
	}

	// parse column 24 attr_value_8 : 参数8结算间隔时间
	if data[24] != "" {
		tmp, err = strconv.ParseInt(data[24], 10, 64)
		if err != nil {
			err = errors.New("parse field attr_value_8 参数8结算间隔时间 to int32 failed")
			logger.ErrorWF("parse field attr_value_8 参数8结算间隔时间 to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
				zap.String("parse_data", data[24]),
				zap.Error(err))
			return
		}
		config.Attr_value_8 = int32(tmp)
	}

	// parse column 25 modify_attr_value_variable_id : 修改属性关联变量id
	if data[25] != "" {

		config.Modify_attr_value_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[25], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field modify_attr_value_variable_id 修改属性关联变量id to key int32 failed")
				logger.ErrorWF("parse map field modify_attr_value_variable_id 修改属性关联变量id to key int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[25]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field modify_attr_value_variable_id 修改属性关联变量id to value int32 failed")
				logger.ErrorWF("parse map field modify_attr_value_variable_id 修改属性关联变量id to value int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[25]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Modify_attr_value_variable_id[key] = value
		}
	}

	// parse column 26 modify_attr_value_type : 修改属性数值类型
	if data[26] != "" {

		config.Modify_attr_value_type = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[26], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field modify_attr_value_type 修改属性数值类型 to key int32 failed")
				logger.ErrorWF("parse map field modify_attr_value_type 修改属性数值类型 to key int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[26]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field modify_attr_value_type 修改属性数值类型 to value int32 failed")
				logger.ErrorWF("parse map field modify_attr_value_type 修改属性数值类型 to value int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[26]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Modify_attr_value_type[key] = value
		}
	}

	// parse column 27 modify_attr_value_attr_id : 修改属性影响属性id
	if data[27] != "" {

		config.Modify_attr_value_attr_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[27], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field modify_attr_value_attr_id 修改属性影响属性id to key int32 failed")
				logger.ErrorWF("parse map field modify_attr_value_attr_id 修改属性影响属性id to key int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[27]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field modify_attr_value_attr_id 修改属性影响属性id to value int32 failed")
				logger.ErrorWF("parse map field modify_attr_value_attr_id 修改属性影响属性id to value int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[27]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Modify_attr_value_attr_id[key] = value
		}
	}

	// parse column 28 modify_attr_value : 修改属性影响数值
	if data[28] != "" {

		config.Modify_attr_value = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[28], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field modify_attr_value 修改属性影响数值 to key int32 failed")
				logger.ErrorWF("parse map field modify_attr_value 修改属性影响数值 to key int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[28]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field modify_attr_value 修改属性影响数值 to value int32 failed")
				logger.ErrorWF("parse map field modify_attr_value 修改属性影响数值 to value int32 failed.",
					zap.String("xlsx", "maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx"), zap.String("sheet", "maze_skilleffect_v8"),
					// zap.String("field_data",data[28]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Modify_attr_value[key] = value
		}
	}
	return
}

var gMazeSkilleffectV8Fields = []string{
	"effect__id",
	"effect_group",
	"in_group_weight",
	"cool_down",
	"attr",
	"attr_value_variable_id",
	"attr_value_type",
	"attr_value",
	"attr_value_2_variable_id",
	"attr_value_2_type",
	"attr_value_2",
	"attr_value_3_variable_id",
	"attr_value_3_type",
	"attr_value_3",
	"attr_value_4",
	"last_time_variable_id",
	"last_time",
	"base_hitrate_variable_id",
	"base_hitrate",
	"attr_value_7_variable_id",
	"attr_value_7_type",
	"attr_value_7",
	"attr_value_8_variable_id",
	"attr_value_8_type",
	"attr_value_8",
	"modify_attr_value_variable_id",
	"modify_attr_value_type",
	"modify_attr_value_attr_id",
	"modify_attr_value",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeSkilleffectV8Parser{}
	loader := &gMazeSkilleffectV8Loader{}
	var data [][]string
	data, err = load("maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx", "maze_skilleffect_v8", gMazeSkilleffectV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx maze_skilleffect_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx maze_skilleffect_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeSkilleffectV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx maze_skilleffect_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeSkilleffectV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx maze_skilleffect_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx maze_skilleffect_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx maze_skilleffect_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_skill_effect_v8【迷宫-技能-技能效果】.xlsx maze_skilleffect_v8 data success.")
	return
}
