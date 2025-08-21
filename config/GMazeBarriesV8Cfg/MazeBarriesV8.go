package GMazeBarriesV8Cfg

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

// MazeBarriesV8ConfigRow from maze_barries_v8【迷宫-关卡信息】.xlsx maze_barries_v8
type MazeBarriesV8ConfigRow struct {
	Order                      int32           `json:"order"`                      // 关卡id
	Last_id                    int32           `json:"last_id"`                    // 上一关id
	Next_id                    int32           `json:"next_id"`                    // 下一关id
	Name                       string          `json:"name"`                       // 名称
	Box_id                     int32           `json:"box_id"`                     // 通关奖励（每次通关都有）
	Box_id_first               int32           `json:"box_id_first"`               // 首次通关奖励
	Hp_vial                    int32           `json:"hp_vial"`                    // 血瓶id
	Drop_vial_foe              []int32         `json:"drop_vial_foe"`              // 掉落血瓶的怪物数
	Challenge_cost             int32           `json:"challenge_cost"`             // 挑战消耗次数（废弃）
	Mop_cost                   int32           `json:"mop_cost"`                   // 挑战和扫荡消耗体力
	Drop_equip_lv_min          int32           `json:"drop_equip_lv_min"`          // 掉落装备等级，小
	Drop_equip_lv_max          int32           `json:"drop_equip_lv_max"`          // 掉落装备等级,大
	Energy_list                []int32         `json:"energy_list"`                // 能力等级队列（随机）
	Barries_add_kongfu         int32           `json:"barries_add_kongfu"`         // 关卡通关值增量
	Energy_id                  map[int32]int32 `json:"energy_id"`                  // 当前关卡使用的能力id
	Energy_affix_rand_rule     map[int32]int32 `json:"energy_affix_rand_rule"`     // 能力词条随机规则
	Rare_items_show            []int32         `json:"rare_items_show"`            // 展示为稀有的物品id
	Drop_id                    []int32         `json:"drop_id"`                    // 掉落id
	Attack_action1_need_kongfu map[int32]int64 `json:"attack_action1_need_kongfu"` // 普攻档位1刷怪区域：所需武力值
	Attack_action2_need_kongfu map[int32]int64 `json:"attack_action2_need_kongfu"` // 普攻档位2刷怪区域：所需武力值
	Attack_action3_need_kongfu map[int32]int64 `json:"attack_action3_need_kongfu"` // 普攻档位3刷怪区域：所需武力值
	Monster_max_num            map[int32]int32 `json:"monster_max_num"`            // 区域刷怪数量上限
	Box_ids                    []int32         `json:"box_ids"`                    // 宝箱列表
	Equ_drop                   int32           `json:"equ_drop"`                   // 装备掉落id
	Need_equip_score           int32           `json:"need_equip_score"`           // 掉落装备所需积分
	Need_item1_score           int32           `json:"need_item1_score"`           // 掉落道具1金币所需积分
	Need_item2_score           int32           `json:"need_item2_score"`           // 掉落道具2强化石所需积分
	Item1_nums_per_pile        int32           `json:"item1_nums_per_pile"`        // 掉落道具1金币单堆数量
	Item2_nums_per_pile        int32           `json:"item2_nums_per_pile"`        // 掉落道具2强化石单堆数量
	Demon_num                  int32           `json:"demon_num"`                  // 猪妖值
	Initial_speed              int32           `json:"initial_speed"`              // 初始速度（万分比)
	Demon_distance             []int32         `json:"demon_distance"`             // 阶段猪妖距离
}

// MazeBarriesV8Config from maze_barries_v8【迷宫-关卡信息】.xlsx maze_barries_v8
type MazeBarriesV8Config struct {
	ConfigRows map[int32]*MazeBarriesV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeBarriesV8Config {
	ret := &MazeBarriesV8Config{ConfigRows: map[int32]*MazeBarriesV8ConfigRow{}}
	return ret
}

// GetMazeBarriesV8Config get one config by configId
func (c *MazeBarriesV8Config) GetMazeBarriesV8Config(configId int32) *MazeBarriesV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeBarriesV8Config) Get(configId int32) *MazeBarriesV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeBarriesV8Config get all config slice
func (c *MazeBarriesV8Config) GetAllMazeBarriesV8Config() (res []*MazeBarriesV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeBarriesV8Config) GetAll() (res []*MazeBarriesV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeBarriesV8Config

// GetMazeBarriesV8Config pkg func. get one config by configId
func GetMazeBarriesV8Config(configId int32) *MazeBarriesV8ConfigRow {
	return gConfigData.GetMazeBarriesV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeBarriesV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeBarriesV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_barries_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeBarriesV8Config pkg func. get all config slice
func GetAllMazeBarriesV8Config() []*MazeBarriesV8ConfigRow {
	return gConfigData.GetAllMazeBarriesV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeBarriesV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeBarriesV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeBarriesV8ConfigRow from maze_barries_v8【迷宫-关卡信息】.xlsx maze_barries_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeBarriesV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_barries_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_barries_v8.json",
		"maze_barries_v8【迷宫-关卡信息】.xlsx", "maze_barries_v8",
		&gMazeBarriesV8Parser{}, &gMazeBarriesV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeBarriesV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeBarriesV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeBarriesV8Config))(c)
		return true
	})
}

// RegisterMazeBarriesV8InitCallBack reg config update func (old func)
var RegisterMazeBarriesV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeBarriesV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeBarriesV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeBarriesV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeBarriesV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeBarriesV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeBarriesV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeBarriesV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeBarriesV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeBarriesV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeBarriesV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeBarriesV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeBarriesV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeBarriesV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBarriesV8ConfigRow", zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"),
			zap.String("sheet", "maze_barries_v8"))
		return
	}
	config, ok := container.(*MazeBarriesV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesV8Config")
		logger.ErrorWF("invalid type. not *MazeBarriesV8Config", zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"),
			zap.String("sheet", "maze_barries_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeBarriesV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeBarriesV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesV8Config")
		logger.ErrorWF("invalid type. not *MazeBarriesV8Config", zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"),
			zap.String("sheet", "maze_barries_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeBarriesV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeBarriesV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesV8Config")
		logger.ErrorWF("invalid type. not *MazeBarriesV8Config", zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"),
			zap.String("sheet", "maze_barries_v8"))
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
type gMazeBarriesV8Parser struct {
}

// New new config row data
func (*gMazeBarriesV8Parser) New() interface{} {
	return &MazeBarriesV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeBarriesV8Parser) Fields() []string {
	return gMazeBarriesV8Fields
}

// Parse parse raw data to row data
func (*gMazeBarriesV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeBarriesV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBarriesV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBarriesV8ConfigRow", zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"),
			zap.String("sheet", "maze_barries_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeBarriesV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeBarriesV8ConfigRow",
			zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"),
			zap.String("sheet", "maze_barries_v8"), zap.Int("need_count", len(gMazeBarriesV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 关卡id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 关卡id to int32 failed")
			logger.ErrorWF("parse field order 关卡id to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 last_id : 上一关id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field last_id 上一关id to int32 failed")
			logger.ErrorWF("parse field last_id 上一关id to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Last_id = int32(tmp)
	}

	// parse column 2 next_id : 下一关id
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field next_id 下一关id to int32 failed")
			logger.ErrorWF("parse field next_id 下一关id to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Next_id = int32(tmp)
	}

	// parse column 3 name : 名称
	if data[3] != "" {
		config.Name = data[3]
	}

	// parse column 4 box_id : 通关奖励（每次通关都有）
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field box_id 通关奖励（每次通关都有） to int32 failed")
			logger.ErrorWF("parse field box_id 通关奖励（每次通关都有） to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Box_id = int32(tmp)
	}

	// parse column 5 box_id_first : 首次通关奖励
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field box_id_first 首次通关奖励 to int32 failed")
			logger.ErrorWF("parse field box_id_first 首次通关奖励 to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Box_id_first = int32(tmp)
	}

	// parse column 6 hp_vial : 血瓶id
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field hp_vial 血瓶id to int32 failed")
			logger.ErrorWF("parse field hp_vial 血瓶id to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Hp_vial = int32(tmp)
	}

	// parse column 7 drop_vial_foe : 掉落血瓶的怪物数
	if data[7] != "" {

		vals := strings.Split(data[7], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field drop_vial_foe 掉落血瓶的怪物数 to []int32 failed")
				logger.ErrorWF("parse array field drop_vial_foe 掉落血瓶的怪物数 to []int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[7]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Drop_vial_foe = append(config.Drop_vial_foe, int32(tmp))
		}
	}

	// parse column 8 challenge_cost : 挑战消耗次数（废弃）
	if data[8] != "" {
		tmp, err = strconv.ParseInt(data[8], 10, 64)
		if err != nil {
			err = errors.New("parse field challenge_cost 挑战消耗次数（废弃） to int32 failed")
			logger.ErrorWF("parse field challenge_cost 挑战消耗次数（废弃） to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[8]),
				zap.Error(err))
			return
		}
		config.Challenge_cost = int32(tmp)
	}

	// parse column 9 mop_cost : 挑战和扫荡消耗体力
	if data[9] != "" {
		tmp, err = strconv.ParseInt(data[9], 10, 64)
		if err != nil {
			err = errors.New("parse field mop_cost 挑战和扫荡消耗体力 to int32 failed")
			logger.ErrorWF("parse field mop_cost 挑战和扫荡消耗体力 to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[9]),
				zap.Error(err))
			return
		}
		config.Mop_cost = int32(tmp)
	}

	// parse column 10 drop_equip_lv_min : 掉落装备等级，小
	if data[10] != "" {
		tmp, err = strconv.ParseInt(data[10], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_equip_lv_min 掉落装备等级，小 to int32 failed")
			logger.ErrorWF("parse field drop_equip_lv_min 掉落装备等级，小 to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[10]),
				zap.Error(err))
			return
		}
		config.Drop_equip_lv_min = int32(tmp)
	}

	// parse column 11 drop_equip_lv_max : 掉落装备等级,大
	if data[11] != "" {
		tmp, err = strconv.ParseInt(data[11], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_equip_lv_max 掉落装备等级,大 to int32 failed")
			logger.ErrorWF("parse field drop_equip_lv_max 掉落装备等级,大 to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[11]),
				zap.Error(err))
			return
		}
		config.Drop_equip_lv_max = int32(tmp)
	}

	// parse column 12 energy_list : 能力等级队列（随机）
	if data[12] != "" {

		vals := strings.Split(data[12], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field energy_list 能力等级队列（随机） to []int32 failed")
				logger.ErrorWF("parse array field energy_list 能力等级队列（随机） to []int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[12]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Energy_list = append(config.Energy_list, int32(tmp))
		}
	}

	// parse column 13 barries_add_kongfu : 关卡通关值增量
	if data[13] != "" {
		tmp, err = strconv.ParseInt(data[13], 10, 64)
		if err != nil {
			err = errors.New("parse field barries_add_kongfu 关卡通关值增量 to int32 failed")
			logger.ErrorWF("parse field barries_add_kongfu 关卡通关值增量 to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[13]),
				zap.Error(err))
			return
		}
		config.Barries_add_kongfu = int32(tmp)
	}

	// parse column 14 energy_id : 当前关卡使用的能力id
	if data[14] != "" {

		config.Energy_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[14], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field energy_id 当前关卡使用的能力id to key int32 failed")
				logger.ErrorWF("parse map field energy_id 当前关卡使用的能力id to key int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[14]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field energy_id 当前关卡使用的能力id to value int32 failed")
				logger.ErrorWF("parse map field energy_id 当前关卡使用的能力id to value int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[14]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Energy_id[key] = value
		}
	}

	// parse column 15 energy_affix_rand_rule : 能力词条随机规则
	if data[15] != "" {

		config.Energy_affix_rand_rule = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[15], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field energy_affix_rand_rule 能力词条随机规则 to key int32 failed")
				logger.ErrorWF("parse map field energy_affix_rand_rule 能力词条随机规则 to key int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[15]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field energy_affix_rand_rule 能力词条随机规则 to value int32 failed")
				logger.ErrorWF("parse map field energy_affix_rand_rule 能力词条随机规则 to value int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[15]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Energy_affix_rand_rule[key] = value
		}
	}

	// parse column 16 rare_items_show : 展示为稀有的物品id
	if data[16] != "" {

		vals := strings.Split(data[16], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field rare_items_show 展示为稀有的物品id to []int32 failed")
				logger.ErrorWF("parse array field rare_items_show 展示为稀有的物品id to []int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[16]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Rare_items_show = append(config.Rare_items_show, int32(tmp))
		}
	}

	// parse column 17 drop_id : 掉落id
	if data[17] != "" {

		vals := strings.Split(data[17], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field drop_id 掉落id to []int32 failed")
				logger.ErrorWF("parse array field drop_id 掉落id to []int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[17]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Drop_id = append(config.Drop_id, int32(tmp))
		}
	}

	// parse column 18 attack_action1_need_kongfu : 普攻档位1刷怪区域：所需武力值
	if data[18] != "" {

		config.Attack_action1_need_kongfu = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[18], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field attack_action1_need_kongfu 普攻档位1刷怪区域：所需武力值 to key int32 failed")
				logger.ErrorWF("parse map field attack_action1_need_kongfu 普攻档位1刷怪区域：所需武力值 to key int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[18]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field attack_action1_need_kongfu 普攻档位1刷怪区域：所需武力值 to value int64 failed")
				logger.ErrorWF("parse map field attack_action1_need_kongfu 普攻档位1刷怪区域：所需武力值 to value int64 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[18]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Attack_action1_need_kongfu[key] = value
		}
	}

	// parse column 19 attack_action2_need_kongfu : 普攻档位2刷怪区域：所需武力值
	if data[19] != "" {

		config.Attack_action2_need_kongfu = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[19], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field attack_action2_need_kongfu 普攻档位2刷怪区域：所需武力值 to key int32 failed")
				logger.ErrorWF("parse map field attack_action2_need_kongfu 普攻档位2刷怪区域：所需武力值 to key int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[19]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field attack_action2_need_kongfu 普攻档位2刷怪区域：所需武力值 to value int64 failed")
				logger.ErrorWF("parse map field attack_action2_need_kongfu 普攻档位2刷怪区域：所需武力值 to value int64 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[19]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Attack_action2_need_kongfu[key] = value
		}
	}

	// parse column 20 attack_action3_need_kongfu : 普攻档位3刷怪区域：所需武力值
	if data[20] != "" {

		config.Attack_action3_need_kongfu = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[20], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field attack_action3_need_kongfu 普攻档位3刷怪区域：所需武力值 to key int32 failed")
				logger.ErrorWF("parse map field attack_action3_need_kongfu 普攻档位3刷怪区域：所需武力值 to key int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[20]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field attack_action3_need_kongfu 普攻档位3刷怪区域：所需武力值 to value int64 failed")
				logger.ErrorWF("parse map field attack_action3_need_kongfu 普攻档位3刷怪区域：所需武力值 to value int64 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[20]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Attack_action3_need_kongfu[key] = value
		}
	}

	// parse column 21 monster_max_num : 区域刷怪数量上限
	if data[21] != "" {

		config.Monster_max_num = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[21], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field monster_max_num 区域刷怪数量上限 to key int32 failed")
				logger.ErrorWF("parse map field monster_max_num 区域刷怪数量上限 to key int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[21]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field monster_max_num 区域刷怪数量上限 to value int32 failed")
				logger.ErrorWF("parse map field monster_max_num 区域刷怪数量上限 to value int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[21]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Monster_max_num[key] = value
		}
	}

	// parse column 22 box_ids : 宝箱列表
	if data[22] != "" {

		vals := strings.Split(data[22], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field box_ids 宝箱列表 to []int32 failed")
				logger.ErrorWF("parse array field box_ids 宝箱列表 to []int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[22]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Box_ids = append(config.Box_ids, int32(tmp))
		}
	}

	// parse column 23 equ_drop : 装备掉落id
	if data[23] != "" {
		tmp, err = strconv.ParseInt(data[23], 10, 64)
		if err != nil {
			err = errors.New("parse field equ_drop 装备掉落id to int32 failed")
			logger.ErrorWF("parse field equ_drop 装备掉落id to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[23]),
				zap.Error(err))
			return
		}
		config.Equ_drop = int32(tmp)
	}

	// parse column 24 need_equip_score : 掉落装备所需积分
	if data[24] != "" {
		tmp, err = strconv.ParseInt(data[24], 10, 64)
		if err != nil {
			err = errors.New("parse field need_equip_score 掉落装备所需积分 to int32 failed")
			logger.ErrorWF("parse field need_equip_score 掉落装备所需积分 to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[24]),
				zap.Error(err))
			return
		}
		config.Need_equip_score = int32(tmp)
	}

	// parse column 25 need_item1_score : 掉落道具1金币所需积分
	if data[25] != "" {
		tmp, err = strconv.ParseInt(data[25], 10, 64)
		if err != nil {
			err = errors.New("parse field need_item1_score 掉落道具1金币所需积分 to int32 failed")
			logger.ErrorWF("parse field need_item1_score 掉落道具1金币所需积分 to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[25]),
				zap.Error(err))
			return
		}
		config.Need_item1_score = int32(tmp)
	}

	// parse column 26 need_item2_score : 掉落道具2强化石所需积分
	if data[26] != "" {
		tmp, err = strconv.ParseInt(data[26], 10, 64)
		if err != nil {
			err = errors.New("parse field need_item2_score 掉落道具2强化石所需积分 to int32 failed")
			logger.ErrorWF("parse field need_item2_score 掉落道具2强化石所需积分 to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[26]),
				zap.Error(err))
			return
		}
		config.Need_item2_score = int32(tmp)
	}

	// parse column 27 item1_nums_per_pile : 掉落道具1金币单堆数量
	if data[27] != "" {
		tmp, err = strconv.ParseInt(data[27], 10, 64)
		if err != nil {
			err = errors.New("parse field item1_nums_per_pile 掉落道具1金币单堆数量 to int32 failed")
			logger.ErrorWF("parse field item1_nums_per_pile 掉落道具1金币单堆数量 to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[27]),
				zap.Error(err))
			return
		}
		config.Item1_nums_per_pile = int32(tmp)
	}

	// parse column 28 item2_nums_per_pile : 掉落道具2强化石单堆数量
	if data[28] != "" {
		tmp, err = strconv.ParseInt(data[28], 10, 64)
		if err != nil {
			err = errors.New("parse field item2_nums_per_pile 掉落道具2强化石单堆数量 to int32 failed")
			logger.ErrorWF("parse field item2_nums_per_pile 掉落道具2强化石单堆数量 to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[28]),
				zap.Error(err))
			return
		}
		config.Item2_nums_per_pile = int32(tmp)
	}

	// parse column 29 demon_num : 猪妖值
	if data[29] != "" {
		tmp, err = strconv.ParseInt(data[29], 10, 64)
		if err != nil {
			err = errors.New("parse field demon_num 猪妖值 to int32 failed")
			logger.ErrorWF("parse field demon_num 猪妖值 to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[29]),
				zap.Error(err))
			return
		}
		config.Demon_num = int32(tmp)
	}

	// parse column 30 initial_speed : 初始速度（万分比)
	if data[30] != "" {
		tmp, err = strconv.ParseInt(data[30], 10, 64)
		if err != nil {
			err = errors.New("parse field initial_speed 初始速度（万分比) to int32 failed")
			logger.ErrorWF("parse field initial_speed 初始速度（万分比) to int32 failed.",
				zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
				zap.String("parse_data", data[30]),
				zap.Error(err))
			return
		}
		config.Initial_speed = int32(tmp)
	}

	// parse column 31 demon_distance : 阶段猪妖距离
	if data[31] != "" {

		vals := strings.Split(data[31], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field demon_distance 阶段猪妖距离 to []int32 failed")
				logger.ErrorWF("parse array field demon_distance 阶段猪妖距离 to []int32 failed.",
					zap.String("xlsx", "maze_barries_v8【迷宫-关卡信息】.xlsx"), zap.String("sheet", "maze_barries_v8"),
					// zap.String("field_data",data[31]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Demon_distance = append(config.Demon_distance, int32(tmp))
		}
	}
	return
}

var gMazeBarriesV8Fields = []string{
	"order",
	"last_id",
	"next_id",
	"name",
	"box_id",
	"box_id_first",
	"hp_vial",
	"drop_vial_foe",
	"challenge_cost",
	"mop_cost",
	"drop_equip_lv_min",
	"drop_equip_lv_max",
	"energy_list",
	"barries_add_kongfu",
	"energy_id",
	"energy_affix_rand_rule",
	"rare_items_show",
	"drop_id",
	"attack_action1_need_kongfu",
	"attack_action2_need_kongfu",
	"attack_action3_need_kongfu",
	"monster_max_num",
	"box_ids",
	"equ_drop",
	"need_equip_score",
	"need_item1_score",
	"need_item2_score",
	"item1_nums_per_pile",
	"item2_nums_per_pile",
	"demon_num",
	"initial_speed",
	"demon_distance",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeBarriesV8Parser{}
	loader := &gMazeBarriesV8Loader{}
	var data [][]string
	data, err = load("maze_barries_v8【迷宫-关卡信息】.xlsx", "maze_barries_v8", gMazeBarriesV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_barries_v8【迷宫-关卡信息】.xlsx maze_barries_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_barries_v8【迷宫-关卡信息】.xlsx maze_barries_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeBarriesV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_barries_v8【迷宫-关卡信息】.xlsx maze_barries_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeBarriesV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_barries_v8【迷宫-关卡信息】.xlsx maze_barries_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_barries_v8【迷宫-关卡信息】.xlsx maze_barries_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_barries_v8【迷宫-关卡信息】.xlsx maze_barries_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_barries_v8【迷宫-关卡信息】.xlsx maze_barries_v8 data success.")
	return
}
