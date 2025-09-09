package GMazeFoeV8Cfg

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

// MazeFoeV8ConfigRow from maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_v8
type MazeFoeV8ConfigRow struct {
	Order                        int32           `json:"order"`                        // 怪物id
	In_barries_id                int32           `json:"in_barries_id"`                // 所属关卡id
	Foe_type                     int32           `json:"foe_type"`                     // 怪物类型（1-小怪 2-守关boss 3-巡逻守卫 4-精英怪9-木桩怪
	Drop_item                    map[int32]int64 `json:"drop_item"`                    // 怪物掉落物品
	Drop_equip                   []int32         `json:"drop_equip"`                   // 怪物掉落装备
	Name                         string          `json:"name"`                         // 怪物名称
	Model_id                     int32           `json:"model_id"`                     // 资源组id
	Level                        int32           `json:"level"`                        // 怪物等级
	Kongfu                       int32           `json:"kongfu"`                       // 怪物武力值
	Hp_lose_type                 int32           `json:"hp_lose_type"`                 // 损血类型
	Attack_max                   int32           `json:"attack_max"`                   // 怪物攻击
	Def_max                      int32           `json:"def_max"`                      // 怪物防御
	Hp_max                       int32           `json:"hp_max"`                       // 怪物血量
	Fire_res                     int32           `json:"fire_res"`                     // 火元素抗性
	Ice_res                      int32           `json:"ice_res"`                      // 冰元素抗性
	Poi_res                      int32           `json:"poi_res"`                      // 毒元素抗性
	Ele_res                      int32           `json:"ele_res"`                      // 电元素抗性
	Speed                        int32           `json:"speed"`                        // 移动速度(万分比）
	Hp_num                       int32           `json:"hp_num"`                       // 血条数量
	Tough_max                    int32           `json:"tough_max"`                    // 怪物韧性上限
	Drop_exp_num                 int32           `json:"drop_exp_num"`                 // 掉落经验数量
	Drop_coin_num                int32           `json:"drop_coin_num"`                // 掉落钱币数量
	Drop_equip_score_num         int32           `json:"drop_equip_score_num"`         // 掉落装备分数量
	Drop_item1_score_num         int32           `json:"drop_item1_score_num"`         // 掉落道具1（金币）分
	Drop_item2_score_num         int32           `json:"drop_item2_score_num"`         // 掉落道具2（强化石）分
	Nor_attack_skill_id          int32           `json:"nor_attack_skill_id"`          // 普通攻击技能id
	Passive_skill_id             []int32         `json:"passive_skill_id"`             // 被动技能id
	Drop_energy_num              int32           `json:"drop_energy_num"`              // 死亡后掉落的能量点数
	Tough_borke_raitio           map[int64]int32 `json:"tough_borke_raitio"`           // 武力对应削韧倍率（万分比）
	Attack_speed_pro             int32           `json:"attack_speed_pro"`             // 怪物攻击速度系数（>10000加速 ,<10000减速
	Be_attack_recovery_speed_pro int32           `json:"be_attack_recovery_speed_pro"` // 受击回复动作播放速度系数（>10000加速 ,<10000减速
	Tough_deplete                int32           `json:"tough_deplete"`                // 韧性被打空时释放技能
	Search_for_scope             int32           `json:"search_for_scope"`             // 寻敌范围调整值（默认10米）
	Attacked_back_range          int32           `json:"attacked_back_range"`          // 被击退距离系数（万分比）
	Attacked_back_range_after    int32           `json:"attacked_back_range_after"`    // 被击退距离系数（破除韧性后）（万分比）
	Threat_value                 int32           `json:"threat_value"`                 // 威胁值
	Arrow_threat_value           int32           `json:"arrow_threat_value"`           // 远程威胁值
}

// MazeFoeV8Config from maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_v8
type MazeFoeV8Config struct {
	ConfigRows map[int32]*MazeFoeV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeFoeV8Config {
	ret := &MazeFoeV8Config{ConfigRows: map[int32]*MazeFoeV8ConfigRow{}}
	return ret
}

// GetMazeFoeV8Config get one config by configId
func (c *MazeFoeV8Config) GetMazeFoeV8Config(configId int32) *MazeFoeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeFoeV8Config) Get(configId int32) *MazeFoeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeFoeV8Config get all config slice
func (c *MazeFoeV8Config) GetAllMazeFoeV8Config() (res []*MazeFoeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeFoeV8Config) GetAll() (res []*MazeFoeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeFoeV8Config

// GetMazeFoeV8Config pkg func. get one config by configId
func GetMazeFoeV8Config(configId int32) *MazeFoeV8ConfigRow {
	return gConfigData.GetMazeFoeV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeFoeV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeFoeV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_foe_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeFoeV8Config pkg func. get all config slice
func GetAllMazeFoeV8Config() []*MazeFoeV8ConfigRow {
	return gConfigData.GetAllMazeFoeV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeFoeV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeFoeV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeFoeV8ConfigRow from maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeFoeV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_foe_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_foe_v8.json",
		"maze_foe_v8【迷宫-敌人信息】.xlsx", "maze_foe_v8",
		&gMazeFoeV8Parser{}, &gMazeFoeV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeFoeV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeFoeV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeFoeV8Config))(c)
		return true
	})
}

// RegisterMazeFoeV8InitCallBack reg config update func (old func)
var RegisterMazeFoeV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeFoeV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeFoeV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeFoeV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeFoeV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeFoeV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeFoeV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeFoeV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeFoeV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeFoeV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeFoeV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeFoeV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeFoeV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeFoeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeFoeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeFoeV8ConfigRow", zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "maze_foe_v8"))
		return
	}
	config, ok := container.(*MazeFoeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeFoeV8Config")
		logger.ErrorWF("invalid type. not *MazeFoeV8Config", zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "maze_foe_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeFoeV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeFoeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeFoeV8Config")
		logger.ErrorWF("invalid type. not *MazeFoeV8Config", zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "maze_foe_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeFoeV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeFoeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeFoeV8Config")
		logger.ErrorWF("invalid type. not *MazeFoeV8Config", zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "maze_foe_v8"))
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
type gMazeFoeV8Parser struct {
}

// New new config row data
func (*gMazeFoeV8Parser) New() interface{} {
	return &MazeFoeV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeFoeV8Parser) Fields() []string {
	return gMazeFoeV8Fields
}

// Parse parse raw data to row data
func (*gMazeFoeV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeFoeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeFoeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeFoeV8ConfigRow", zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "maze_foe_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeFoeV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeFoeV8ConfigRow",
			zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "maze_foe_v8"), zap.Int("need_count", len(gMazeFoeV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 怪物id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 怪物id to int32 failed")
			logger.ErrorWF("parse field order 怪物id to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 in_barries_id : 所属关卡id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field in_barries_id 所属关卡id to int32 failed")
			logger.ErrorWF("parse field in_barries_id 所属关卡id to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.In_barries_id = int32(tmp)
	}

	// parse column 2 foe_type : 怪物类型（1-小怪 2-守关boss 3-巡逻守卫 4-精英怪9-木桩怪
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field foe_type 怪物类型（1-小怪 2-守关boss 3-巡逻守卫 4-精英怪9-木桩怪 to int32 failed")
			logger.ErrorWF("parse field foe_type 怪物类型（1-小怪 2-守关boss 3-巡逻守卫 4-精英怪9-木桩怪 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Foe_type = int32(tmp)
	}

	// parse column 3 drop_item : 怪物掉落物品
	if data[3] != "" {

		config.Drop_item = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_item 怪物掉落物品 to key int32 failed")
				logger.ErrorWF("parse map field drop_item 怪物掉落物品 to key int32 failed.",
					zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_item 怪物掉落物品 to value int64 failed")
				logger.ErrorWF("parse map field drop_item 怪物掉落物品 to value int64 failed.",
					zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Drop_item[key] = value
		}
	}

	// parse column 4 drop_equip : 怪物掉落装备
	if data[4] != "" {

		vals := strings.Split(data[4], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field drop_equip 怪物掉落装备 to []int32 failed")
				logger.ErrorWF("parse array field drop_equip 怪物掉落装备 to []int32 failed.",
					zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
					// zap.String("field_data",data[4]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Drop_equip = append(config.Drop_equip, int32(tmp))
		}
	}

	// parse column 5 name : 怪物名称
	if data[5] != "" {
		config.Name = data[5]
	}

	// parse column 6 model_id : 资源组id
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field model_id 资源组id to int32 failed")
			logger.ErrorWF("parse field model_id 资源组id to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Model_id = int32(tmp)
	}

	// parse column 7 level : 怪物等级
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field level 怪物等级 to int32 failed")
			logger.ErrorWF("parse field level 怪物等级 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Level = int32(tmp)
	}

	// parse column 8 kongfu : 怪物武力值
	if data[8] != "" {
		tmp, err = strconv.ParseInt(data[8], 10, 64)
		if err != nil {
			err = errors.New("parse field kongfu 怪物武力值 to int32 failed")
			logger.ErrorWF("parse field kongfu 怪物武力值 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[8]),
				zap.Error(err))
			return
		}
		config.Kongfu = int32(tmp)
	}

	// parse column 9 hp_lose_type : 损血类型
	if data[9] != "" {
		tmp, err = strconv.ParseInt(data[9], 10, 64)
		if err != nil {
			err = errors.New("parse field hp_lose_type 损血类型 to int32 failed")
			logger.ErrorWF("parse field hp_lose_type 损血类型 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[9]),
				zap.Error(err))
			return
		}
		config.Hp_lose_type = int32(tmp)
	}

	// parse column 10 attack_max : 怪物攻击
	if data[10] != "" {
		tmp, err = strconv.ParseInt(data[10], 10, 64)
		if err != nil {
			err = errors.New("parse field attack_max 怪物攻击 to int32 failed")
			logger.ErrorWF("parse field attack_max 怪物攻击 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[10]),
				zap.Error(err))
			return
		}
		config.Attack_max = int32(tmp)
	}

	// parse column 11 def_max : 怪物防御
	if data[11] != "" {
		tmp, err = strconv.ParseInt(data[11], 10, 64)
		if err != nil {
			err = errors.New("parse field def_max 怪物防御 to int32 failed")
			logger.ErrorWF("parse field def_max 怪物防御 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[11]),
				zap.Error(err))
			return
		}
		config.Def_max = int32(tmp)
	}

	// parse column 12 hp_max : 怪物血量
	if data[12] != "" {
		tmp, err = strconv.ParseInt(data[12], 10, 64)
		if err != nil {
			err = errors.New("parse field hp_max 怪物血量 to int32 failed")
			logger.ErrorWF("parse field hp_max 怪物血量 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[12]),
				zap.Error(err))
			return
		}
		config.Hp_max = int32(tmp)
	}

	// parse column 13 fire_res : 火元素抗性
	if data[13] != "" {
		tmp, err = strconv.ParseInt(data[13], 10, 64)
		if err != nil {
			err = errors.New("parse field fire_res 火元素抗性 to int32 failed")
			logger.ErrorWF("parse field fire_res 火元素抗性 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[13]),
				zap.Error(err))
			return
		}
		config.Fire_res = int32(tmp)
	}

	// parse column 14 ice_res : 冰元素抗性
	if data[14] != "" {
		tmp, err = strconv.ParseInt(data[14], 10, 64)
		if err != nil {
			err = errors.New("parse field ice_res 冰元素抗性 to int32 failed")
			logger.ErrorWF("parse field ice_res 冰元素抗性 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[14]),
				zap.Error(err))
			return
		}
		config.Ice_res = int32(tmp)
	}

	// parse column 15 poi_res : 毒元素抗性
	if data[15] != "" {
		tmp, err = strconv.ParseInt(data[15], 10, 64)
		if err != nil {
			err = errors.New("parse field poi_res 毒元素抗性 to int32 failed")
			logger.ErrorWF("parse field poi_res 毒元素抗性 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[15]),
				zap.Error(err))
			return
		}
		config.Poi_res = int32(tmp)
	}

	// parse column 16 ele_res : 电元素抗性
	if data[16] != "" {
		tmp, err = strconv.ParseInt(data[16], 10, 64)
		if err != nil {
			err = errors.New("parse field ele_res 电元素抗性 to int32 failed")
			logger.ErrorWF("parse field ele_res 电元素抗性 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[16]),
				zap.Error(err))
			return
		}
		config.Ele_res = int32(tmp)
	}

	// parse column 17 speed : 移动速度(万分比）
	if data[17] != "" {
		tmp, err = strconv.ParseInt(data[17], 10, 64)
		if err != nil {
			err = errors.New("parse field speed 移动速度(万分比） to int32 failed")
			logger.ErrorWF("parse field speed 移动速度(万分比） to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[17]),
				zap.Error(err))
			return
		}
		config.Speed = int32(tmp)
	}

	// parse column 18 hp_num : 血条数量
	if data[18] != "" {
		tmp, err = strconv.ParseInt(data[18], 10, 64)
		if err != nil {
			err = errors.New("parse field hp_num 血条数量 to int32 failed")
			logger.ErrorWF("parse field hp_num 血条数量 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[18]),
				zap.Error(err))
			return
		}
		config.Hp_num = int32(tmp)
	}

	// parse column 19 tough_max : 怪物韧性上限
	if data[19] != "" {
		tmp, err = strconv.ParseInt(data[19], 10, 64)
		if err != nil {
			err = errors.New("parse field tough_max 怪物韧性上限 to int32 failed")
			logger.ErrorWF("parse field tough_max 怪物韧性上限 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[19]),
				zap.Error(err))
			return
		}
		config.Tough_max = int32(tmp)
	}

	// parse column 20 drop_exp_num : 掉落经验数量
	if data[20] != "" {
		tmp, err = strconv.ParseInt(data[20], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_exp_num 掉落经验数量 to int32 failed")
			logger.ErrorWF("parse field drop_exp_num 掉落经验数量 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[20]),
				zap.Error(err))
			return
		}
		config.Drop_exp_num = int32(tmp)
	}

	// parse column 21 drop_coin_num : 掉落钱币数量
	if data[21] != "" {
		tmp, err = strconv.ParseInt(data[21], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_coin_num 掉落钱币数量 to int32 failed")
			logger.ErrorWF("parse field drop_coin_num 掉落钱币数量 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[21]),
				zap.Error(err))
			return
		}
		config.Drop_coin_num = int32(tmp)
	}

	// parse column 22 drop_equip_score_num : 掉落装备分数量
	if data[22] != "" {
		tmp, err = strconv.ParseInt(data[22], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_equip_score_num 掉落装备分数量 to int32 failed")
			logger.ErrorWF("parse field drop_equip_score_num 掉落装备分数量 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[22]),
				zap.Error(err))
			return
		}
		config.Drop_equip_score_num = int32(tmp)
	}

	// parse column 23 drop_item1_score_num : 掉落道具1（金币）分
	if data[23] != "" {
		tmp, err = strconv.ParseInt(data[23], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_item1_score_num 掉落道具1（金币）分 to int32 failed")
			logger.ErrorWF("parse field drop_item1_score_num 掉落道具1（金币）分 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[23]),
				zap.Error(err))
			return
		}
		config.Drop_item1_score_num = int32(tmp)
	}

	// parse column 24 drop_item2_score_num : 掉落道具2（强化石）分
	if data[24] != "" {
		tmp, err = strconv.ParseInt(data[24], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_item2_score_num 掉落道具2（强化石）分 to int32 failed")
			logger.ErrorWF("parse field drop_item2_score_num 掉落道具2（强化石）分 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[24]),
				zap.Error(err))
			return
		}
		config.Drop_item2_score_num = int32(tmp)
	}

	// parse column 25 nor_attack_skill_id : 普通攻击技能id
	if data[25] != "" {
		tmp, err = strconv.ParseInt(data[25], 10, 64)
		if err != nil {
			err = errors.New("parse field nor_attack_skill_id 普通攻击技能id to int32 failed")
			logger.ErrorWF("parse field nor_attack_skill_id 普通攻击技能id to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[25]),
				zap.Error(err))
			return
		}
		config.Nor_attack_skill_id = int32(tmp)
	}

	// parse column 26 passive_skill_id : 被动技能id
	if data[26] != "" {

		vals := strings.Split(data[26], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field passive_skill_id 被动技能id to []int32 failed")
				logger.ErrorWF("parse array field passive_skill_id 被动技能id to []int32 failed.",
					zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
					// zap.String("field_data",data[26]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Passive_skill_id = append(config.Passive_skill_id, int32(tmp))
		}
	}

	// parse column 27 drop_energy_num : 死亡后掉落的能量点数
	if data[27] != "" {
		tmp, err = strconv.ParseInt(data[27], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_energy_num 死亡后掉落的能量点数 to int32 failed")
			logger.ErrorWF("parse field drop_energy_num 死亡后掉落的能量点数 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[27]),
				zap.Error(err))
			return
		}
		config.Drop_energy_num = int32(tmp)
	}

	// parse column 28 tough_borke_raitio : 武力对应削韧倍率（万分比）
	if data[28] != "" {

		config.Tough_borke_raitio = make(map[int64]int32)
		var key int64
		var value int32
		vals := strings.Split(data[28], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field tough_borke_raitio 武力对应削韧倍率（万分比） to key int64 failed")
				logger.ErrorWF("parse map field tough_borke_raitio 武力对应削韧倍率（万分比） to key int64 failed.",
					zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
					// zap.String("field_data",data[28]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int64(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field tough_borke_raitio 武力对应削韧倍率（万分比） to value int32 failed")
				logger.ErrorWF("parse map field tough_borke_raitio 武力对应削韧倍率（万分比） to value int32 failed.",
					zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
					// zap.String("field_data",data[28]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Tough_borke_raitio[key] = value
		}
	}

	// parse column 29 attack_speed_pro : 怪物攻击速度系数（>10000加速 ,<10000减速
	if data[29] != "" {
		tmp, err = strconv.ParseInt(data[29], 10, 64)
		if err != nil {
			err = errors.New("parse field attack_speed_pro 怪物攻击速度系数（>10000加速 ,<10000减速 to int32 failed")
			logger.ErrorWF("parse field attack_speed_pro 怪物攻击速度系数（>10000加速 ,<10000减速 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[29]),
				zap.Error(err))
			return
		}
		config.Attack_speed_pro = int32(tmp)
	}

	// parse column 30 be_attack_recovery_speed_pro : 受击回复动作播放速度系数（>10000加速 ,<10000减速
	if data[30] != "" {
		tmp, err = strconv.ParseInt(data[30], 10, 64)
		if err != nil {
			err = errors.New("parse field be_attack_recovery_speed_pro 受击回复动作播放速度系数（>10000加速 ,<10000减速 to int32 failed")
			logger.ErrorWF("parse field be_attack_recovery_speed_pro 受击回复动作播放速度系数（>10000加速 ,<10000减速 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[30]),
				zap.Error(err))
			return
		}
		config.Be_attack_recovery_speed_pro = int32(tmp)
	}

	// parse column 31 tough_deplete : 韧性被打空时释放技能
	if data[31] != "" {
		tmp, err = strconv.ParseInt(data[31], 10, 64)
		if err != nil {
			err = errors.New("parse field tough_deplete 韧性被打空时释放技能 to int32 failed")
			logger.ErrorWF("parse field tough_deplete 韧性被打空时释放技能 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[31]),
				zap.Error(err))
			return
		}
		config.Tough_deplete = int32(tmp)
	}

	// parse column 32 search_for_scope : 寻敌范围调整值（默认10米）
	if data[32] != "" {
		tmp, err = strconv.ParseInt(data[32], 10, 64)
		if err != nil {
			err = errors.New("parse field search_for_scope 寻敌范围调整值（默认10米） to int32 failed")
			logger.ErrorWF("parse field search_for_scope 寻敌范围调整值（默认10米） to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[32]),
				zap.Error(err))
			return
		}
		config.Search_for_scope = int32(tmp)
	}

	// parse column 33 attacked_back_range : 被击退距离系数（万分比）
	if data[33] != "" {
		tmp, err = strconv.ParseInt(data[33], 10, 64)
		if err != nil {
			err = errors.New("parse field attacked_back_range 被击退距离系数（万分比） to int32 failed")
			logger.ErrorWF("parse field attacked_back_range 被击退距离系数（万分比） to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[33]),
				zap.Error(err))
			return
		}
		config.Attacked_back_range = int32(tmp)
	}

	// parse column 34 attacked_back_range_after : 被击退距离系数（破除韧性后）（万分比）
	if data[34] != "" {
		tmp, err = strconv.ParseInt(data[34], 10, 64)
		if err != nil {
			err = errors.New("parse field attacked_back_range_after 被击退距离系数（破除韧性后）（万分比） to int32 failed")
			logger.ErrorWF("parse field attacked_back_range_after 被击退距离系数（破除韧性后）（万分比） to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[34]),
				zap.Error(err))
			return
		}
		config.Attacked_back_range_after = int32(tmp)
	}

	// parse column 35 threat_value : 威胁值
	if data[35] != "" {
		tmp, err = strconv.ParseInt(data[35], 10, 64)
		if err != nil {
			err = errors.New("parse field threat_value 威胁值 to int32 failed")
			logger.ErrorWF("parse field threat_value 威胁值 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[35]),
				zap.Error(err))
			return
		}
		config.Threat_value = int32(tmp)
	}

	// parse column 36 arrow_threat_value : 远程威胁值
	if data[36] != "" {
		tmp, err = strconv.ParseInt(data[36], 10, 64)
		if err != nil {
			err = errors.New("parse field arrow_threat_value 远程威胁值 to int32 failed")
			logger.ErrorWF("parse field arrow_threat_value 远程威胁值 to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_v8"),
				zap.String("parse_data", data[36]),
				zap.Error(err))
			return
		}
		config.Arrow_threat_value = int32(tmp)
	}
	return
}

var gMazeFoeV8Fields = []string{
	"order",
	"in_barries_id",
	"foe_type",
	"drop_item",
	"drop_equip",
	"name",
	"model_id",
	"level",
	"kongfu",
	"hp_lose_type",
	"attack_max",
	"def_max",
	"hp_max",
	"fire_res",
	"ice_res",
	"poi_res",
	"ele_res",
	"speed",
	"hp_num",
	"tough_max",
	"drop_exp_num",
	"drop_coin_num",
	"drop_equip_score_num",
	"drop_item1_score_num",
	"drop_item2_score_num",
	"nor_attack_skill_id",
	"passive_skill_id",
	"drop_energy_num",
	"tough_borke_raitio",
	"attack_speed_pro",
	"be_attack_recovery_speed_pro",
	"tough_deplete",
	"search_for_scope",
	"attacked_back_range",
	"attacked_back_range_after",
	"threat_value",
	"arrow_threat_value",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeFoeV8Parser{}
	loader := &gMazeFoeV8Loader{}
	var data [][]string
	data, err = load("maze_foe_v8【迷宫-敌人信息】.xlsx", "maze_foe_v8", gMazeFoeV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeFoeV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeFoeV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_v8 data success.")
	return
}
