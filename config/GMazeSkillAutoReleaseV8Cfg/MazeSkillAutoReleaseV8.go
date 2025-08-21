package GMazeSkillAutoReleaseV8Cfg

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

// MazeSkillAutoReleaseV8ConfigRow from maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8
type MazeSkillAutoReleaseV8ConfigRow struct {
	Order                            int32           `json:"order"`                            // 自动
	Release_time                     int32           `json:"release_time"`                     // 自动释放时机
	Release_condition                string          `json:"release_condition"`                // 自动释放条件
	Max_release_limit_variable       map[int32]int32 `json:"max_release_limit_variable"`       // 变量属性id：变化方式（1-基础值+value,2=基础值-value）
	Max_release_limit                int32           `json:"max_release_limit"`                // 每次挑战最大触发次数（0表示无限制）
	Auto_release_protect_cd_variable map[int32]int32 `json:"auto_release_protect_cd_variable"` // 变量属性id：变化方式（1-基础值+value,2=基础值-value）
	Auto_release_protect_cd          int32           `json:"auto_release_protect_cd"`          // 参数-重复触发保护-自动释放冷却时间（毫秒）
	Release_ratio_variable           map[int32]int32 `json:"release_ratio_variable"`           // 变量属性id：变化方式（1-基础值+value,2=基础值-value）
	Release_ratio                    int32           `json:"release_ratio"`                    // 参数-释放几率
	Release_num_variable             map[int32]int32 `json:"release_num_variable"`             // 变量属性id：变化方式（1-基础值+value,2=基础值-value）
	Release_num                      int32           `json:"release_num"`                      // 参数-释放次数
	Skill_id                         []int32         `json:"skill_id"`                         // 参数-释放的技能id
}

// MazeSkillAutoReleaseV8Config from maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8
type MazeSkillAutoReleaseV8Config struct {
	ConfigRows map[int32]*MazeSkillAutoReleaseV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeSkillAutoReleaseV8Config {
	ret := &MazeSkillAutoReleaseV8Config{ConfigRows: map[int32]*MazeSkillAutoReleaseV8ConfigRow{}}
	return ret
}

// GetMazeSkillAutoReleaseV8Config get one config by configId
func (c *MazeSkillAutoReleaseV8Config) GetMazeSkillAutoReleaseV8Config(configId int32) *MazeSkillAutoReleaseV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeSkillAutoReleaseV8Config) Get(configId int32) *MazeSkillAutoReleaseV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeSkillAutoReleaseV8Config get all config slice
func (c *MazeSkillAutoReleaseV8Config) GetAllMazeSkillAutoReleaseV8Config() (res []*MazeSkillAutoReleaseV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeSkillAutoReleaseV8Config) GetAll() (res []*MazeSkillAutoReleaseV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeSkillAutoReleaseV8Config

// GetMazeSkillAutoReleaseV8Config pkg func. get one config by configId
func GetMazeSkillAutoReleaseV8Config(configId int32) *MazeSkillAutoReleaseV8ConfigRow {
	return gConfigData.GetMazeSkillAutoReleaseV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeSkillAutoReleaseV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeSkillAutoReleaseV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_skill_auto_release_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeSkillAutoReleaseV8Config pkg func. get all config slice
func GetAllMazeSkillAutoReleaseV8Config() []*MazeSkillAutoReleaseV8ConfigRow {
	return gConfigData.GetAllMazeSkillAutoReleaseV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeSkillAutoReleaseV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeSkillAutoReleaseV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeSkillAutoReleaseV8ConfigRow from maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeSkillAutoReleaseV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_skill_auto_release_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_skill_auto_release_v8.json",
		"maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx", "maze_skill_auto_release_v8",
		&gMazeSkillAutoReleaseV8Parser{}, &gMazeSkillAutoReleaseV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeSkillAutoReleaseV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeSkillAutoReleaseV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeSkillAutoReleaseV8Config))(c)
		return true
	})
}

// RegisterMazeSkillAutoReleaseV8InitCallBack reg config update func (old func)
var RegisterMazeSkillAutoReleaseV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeSkillAutoReleaseV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeSkillAutoReleaseV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeSkillAutoReleaseV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeSkillAutoReleaseV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeSkillAutoReleaseV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeSkillAutoReleaseV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeSkillAutoReleaseV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeSkillAutoReleaseV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeSkillAutoReleaseV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeSkillAutoReleaseV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeSkillAutoReleaseV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeSkillAutoReleaseV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeSkillAutoReleaseV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoReleaseV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkillAutoReleaseV8ConfigRow", zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"),
			zap.String("sheet", "maze_skill_auto_release_v8"))
		return
	}
	config, ok := container.(*MazeSkillAutoReleaseV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoReleaseV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillAutoReleaseV8Config", zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"),
			zap.String("sheet", "maze_skill_auto_release_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeSkillAutoReleaseV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeSkillAutoReleaseV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoReleaseV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillAutoReleaseV8Config", zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"),
			zap.String("sheet", "maze_skill_auto_release_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeSkillAutoReleaseV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeSkillAutoReleaseV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoReleaseV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillAutoReleaseV8Config", zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"),
			zap.String("sheet", "maze_skill_auto_release_v8"))
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
type gMazeSkillAutoReleaseV8Parser struct {
}

// New new config row data
func (*gMazeSkillAutoReleaseV8Parser) New() interface{} {
	return &MazeSkillAutoReleaseV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeSkillAutoReleaseV8Parser) Fields() []string {
	return gMazeSkillAutoReleaseV8Fields
}

// Parse parse raw data to row data
func (*gMazeSkillAutoReleaseV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeSkillAutoReleaseV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoReleaseV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkillAutoReleaseV8ConfigRow", zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"),
			zap.String("sheet", "maze_skill_auto_release_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeSkillAutoReleaseV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeSkillAutoReleaseV8ConfigRow",
			zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"),
			zap.String("sheet", "maze_skill_auto_release_v8"), zap.Int("need_count", len(gMazeSkillAutoReleaseV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 自动
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 自动 to int32 failed")
			logger.ErrorWF("parse field order 自动 to int32 failed.",
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 release_time : 自动释放时机
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field release_time 自动释放时机 to int32 failed")
			logger.ErrorWF("parse field release_time 自动释放时机 to int32 failed.",
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Release_time = int32(tmp)
	}

	// parse column 2 release_condition : 自动释放条件
	if data[2] != "" {
		config.Release_condition = data[2]
	}

	// parse column 3 max_release_limit_variable : 变量属性id：变化方式（1-基础值+value,2=基础值-value）
	if data[3] != "" {

		config.Max_release_limit_variable = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field max_release_limit_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to key int32 failed")
				logger.ErrorWF("parse map field max_release_limit_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to key int32 failed.",
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field max_release_limit_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to value int32 failed")
				logger.ErrorWF("parse map field max_release_limit_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to value int32 failed.",
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Max_release_limit_variable[key] = value
		}
	}

	// parse column 4 max_release_limit : 每次挑战最大触发次数（0表示无限制）
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field max_release_limit 每次挑战最大触发次数（0表示无限制） to int32 failed")
			logger.ErrorWF("parse field max_release_limit 每次挑战最大触发次数（0表示无限制） to int32 failed.",
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Max_release_limit = int32(tmp)
	}

	// parse column 5 auto_release_protect_cd_variable : 变量属性id：变化方式（1-基础值+value,2=基础值-value）
	if data[5] != "" {

		config.Auto_release_protect_cd_variable = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[5], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field auto_release_protect_cd_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to key int32 failed")
				logger.ErrorWF("parse map field auto_release_protect_cd_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to key int32 failed.",
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field auto_release_protect_cd_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to value int32 failed")
				logger.ErrorWF("parse map field auto_release_protect_cd_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to value int32 failed.",
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Auto_release_protect_cd_variable[key] = value
		}
	}

	// parse column 6 auto_release_protect_cd : 参数-重复触发保护-自动释放冷却时间（毫秒）
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field auto_release_protect_cd 参数-重复触发保护-自动释放冷却时间（毫秒） to int32 failed")
			logger.ErrorWF("parse field auto_release_protect_cd 参数-重复触发保护-自动释放冷却时间（毫秒） to int32 failed.",
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Auto_release_protect_cd = int32(tmp)
	}

	// parse column 7 release_ratio_variable : 变量属性id：变化方式（1-基础值+value,2=基础值-value）
	if data[7] != "" {

		config.Release_ratio_variable = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[7], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field release_ratio_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to key int32 failed")
				logger.ErrorWF("parse map field release_ratio_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to key int32 failed.",
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field release_ratio_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to value int32 failed")
				logger.ErrorWF("parse map field release_ratio_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to value int32 failed.",
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Release_ratio_variable[key] = value
		}
	}

	// parse column 8 release_ratio : 参数-释放几率
	if data[8] != "" {
		tmp, err = strconv.ParseInt(data[8], 10, 64)
		if err != nil {
			err = errors.New("parse field release_ratio 参数-释放几率 to int32 failed")
			logger.ErrorWF("parse field release_ratio 参数-释放几率 to int32 failed.",
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
				zap.String("parse_data", data[8]),
				zap.Error(err))
			return
		}
		config.Release_ratio = int32(tmp)
	}

	// parse column 9 release_num_variable : 变量属性id：变化方式（1-基础值+value,2=基础值-value）
	if data[9] != "" {

		config.Release_num_variable = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[9], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field release_num_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to key int32 failed")
				logger.ErrorWF("parse map field release_num_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to key int32 failed.",
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
					// zap.String("field_data",data[9]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field release_num_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to value int32 failed")
				logger.ErrorWF("parse map field release_num_variable 变量属性id：变化方式（1-基础值+value,2=基础值-value） to value int32 failed.",
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
					// zap.String("field_data",data[9]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Release_num_variable[key] = value
		}
	}

	// parse column 10 release_num : 参数-释放次数
	if data[10] != "" {
		tmp, err = strconv.ParseInt(data[10], 10, 64)
		if err != nil {
			err = errors.New("parse field release_num 参数-释放次数 to int32 failed")
			logger.ErrorWF("parse field release_num 参数-释放次数 to int32 failed.",
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
				zap.String("parse_data", data[10]),
				zap.Error(err))
			return
		}
		config.Release_num = int32(tmp)
	}

	// parse column 11 skill_id : 参数-释放的技能id
	if data[11] != "" {

		vals := strings.Split(data[11], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field skill_id 参数-释放的技能id to []int32 failed")
				logger.ErrorWF("parse array field skill_id 参数-释放的技能id to []int32 failed.",
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"),
					// zap.String("field_data",data[11]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Skill_id = append(config.Skill_id, int32(tmp))
		}
	}
	return
}

var gMazeSkillAutoReleaseV8Fields = []string{
	"order",
	"release_time",
	"release_condition",
	"max_release_limit_variable",
	"max_release_limit",
	"auto_release_protect_cd_variable",
	"auto_release_protect_cd",
	"release_ratio_variable",
	"release_ratio",
	"release_num_variable",
	"release_num",
	"skill_id",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeSkillAutoReleaseV8Parser{}
	loader := &gMazeSkillAutoReleaseV8Loader{}
	var data [][]string
	data, err = load("maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx", "maze_skill_auto_release_v8", gMazeSkillAutoReleaseV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeSkillAutoReleaseV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeSkillAutoReleaseV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 data success.")
	return
}
