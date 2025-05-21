package GDollEquipInfoV8Cfg

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

// DollEquipInfoV8ConfigRow from doll_equip_info_v8【人偶-装备-信息】.xlsx doll_equip_info_v8
type DollEquipInfoV8ConfigRow struct {
	Equipment_id             int32           `json:"equipment_id"`             // 序号
	Score_group              int32           `json:"score_group"`              // 记录装备分数的组id
	Quality                  int32           `json:"quality"`                  // 品质
	Pos                      int32           `json:"pos"`                      // 部位
	Skill_pos_num            map[int32]int32 `json:"skill_pos_num"`            // 随机技能孔数量:权重
	Sub_type_random          map[int32]int32 `json:"sub_type_random"`          // 装备子类型随机权重
	Level                    int32           `json:"level"`                    // 穿戴等级
	Level_show               int32           `json:"level_show"`               // 显示等级
	Need_dungeon             map[int32]int32 `json:"need_dungeon"`             // 允许穿戴需要通关的副本
	Sell                     map[int32]int64 `json:"sell"`                     // 出售价格
	Split_num                map[int32]int32 `json:"split_num"`                // 分解后获得材料次数：权重
	Split_random_id          int32           `json:"split_random_id"`          // 分解随机库id
	Affix_base_num           map[int32]int32 `json:"affix_base_num"`           // 初始基础词条数量
	Affix_base_pool          map[int32]int32 `json:"affix_base_pool"`          // 基础词条池子id:词条位置(affix_pool
	Affix_rand_num           map[int32]int32 `json:"affix_rand_num"`           // 初始随机词条数量
	Affix_rand_pool          map[int32]int32 `json:"affix_rand_pool"`          // 随机词条池子id:权重(affix_pool
	Affix_mod_num            map[int32]int32 `json:"affix_mod_num"`            // 初始传奇词条数量
	Affix_mod_pool           map[int32]int32 `json:"affix_mod_pool"`           // 传奇词条池子id:权重（mod_pool
	Affix_extra_pool         map[int32]int32 `json:"affix_extra_pool"`         // 额外词条池子id:权重（sp_pool
	Suite_id                 map[int32]int32 `json:"suite_id"`                 // 装备套装id:随机权重
	Unknow_view_equipment_id int32           `json:"unknow_view_equipment_id"` // 未鉴定的装备id（0-表示不需要鉴定
	View_cost                map[int32]int64 `json:"view_cost"`                // 鉴定消耗
	View_need_dungeon        map[int32]int32 `json:"view_need_dungeon"`        // 允许鉴定需要通关的副本
	Pos_sub_type             int32           `json:"pos_sub_type"`             // 部位子类型
	Weapon_model             int32           `json:"weapon_model"`             // 武器模型（废弃）
	Id                       int32           `json:"id"`                       // 记录装备分数的组id
	Preview_attr             map[int32]int64 `json:"preview_attr"`             // 预览属性
	Attack_type              int32           `json:"attack_type"`              // 攻击方式
	Damage_type              int32           `json:"damage_type"`              // 伤害类型
	Font_soul                map[int32]int32 `json:"font_soul"`                // 升华头属性：随机权重
	Tail_soul                map[int32]int32 `json:"tail_soul"`                // 升华尾属性：随机权重
	Recast_cost              map[int32]int64 `json:"recast_cost"`              // 升华属性重铸消耗
	Soul_affix_base_pool     map[int32]int32 `json:"soul_affix_base_pool"`     // 升华属性（走基础词条池子id:词条位置(affix_pool
}

// DollEquipInfoV8Config from doll_equip_info_v8【人偶-装备-信息】.xlsx doll_equip_info_v8
type DollEquipInfoV8Config struct {
	ConfigRows map[int32]*DollEquipInfoV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *DollEquipInfoV8Config {
	ret := &DollEquipInfoV8Config{ConfigRows: map[int32]*DollEquipInfoV8ConfigRow{}}
	return ret
}

// GetDollEquipInfoV8Config get one config by configId
func (c *DollEquipInfoV8Config) GetDollEquipInfoV8Config(configId int32) *DollEquipInfoV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *DollEquipInfoV8Config) Get(configId int32) *DollEquipInfoV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllDollEquipInfoV8Config get all config slice
func (c *DollEquipInfoV8Config) GetAllDollEquipInfoV8Config() (res []*DollEquipInfoV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *DollEquipInfoV8Config) GetAll() (res []*DollEquipInfoV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *DollEquipInfoV8Config

// GetDollEquipInfoV8Config pkg func. get one config by configId
func GetDollEquipInfoV8Config(configId int32) *DollEquipInfoV8ConfigRow {
	return gConfigData.GetDollEquipInfoV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *DollEquipInfoV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllDollEquipInfoV8Config pkg func. get all config slice
func GetAllDollEquipInfoV8Config() []*DollEquipInfoV8ConfigRow {
	return gConfigData.GetAllDollEquipInfoV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*DollEquipInfoV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*DollEquipInfoV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "DollEquipInfoV8ConfigRow from doll_equip_info_v8【人偶-装备-信息】.xlsx doll_equip_info_v8"
}

// GetRawValue get raw data
func GetRawValue() *DollEquipInfoV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "doll_equip_info_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("doll_equip_info_v8.json",
		"doll_equip_info_v8【人偶-装备-信息】.xlsx", "doll_equip_info_v8",
		&gDollEquipInfoV8Parser{}, &gDollEquipInfoV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*DollEquipInfoV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *DollEquipInfoV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*DollEquipInfoV8Config))(c)
		return true
	})
}

// RegisterDollEquipInfoV8InitCallBack reg config update func (old func)
var RegisterDollEquipInfoV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*DollEquipInfoV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*DollEquipInfoV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *DollEquipInfoV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*DollEquipInfoV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*DollEquipInfoV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gDollEquipInfoV8Loader struct {
}

// NewContainer new data container pointer
func (*gDollEquipInfoV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gDollEquipInfoV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*DollEquipInfoV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gDollEquipInfoV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*DollEquipInfoV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gDollEquipInfoV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*DollEquipInfoV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *DollEquipInfoV8ConfigRow")
		logger.ErrorWF("invalid type. not *DollEquipInfoV8ConfigRow", zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"),
			zap.String("sheet", "doll_equip_info_v8"))
		return
	}
	config, ok := container.(*DollEquipInfoV8Config)
	if !ok {
		err = errors.New("invalid type. not *DollEquipInfoV8Config")
		logger.ErrorWF("invalid type. not *DollEquipInfoV8Config", zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"),
			zap.String("sheet", "doll_equip_info_v8"))
		return
	}
	config.ConfigRows[row.Equipment_id] = row
	return
}

// GetValue get real map value for json parse
func (*gDollEquipInfoV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*DollEquipInfoV8Config)
	if !ok {
		err = errors.New("invalid type. not *DollEquipInfoV8Config")
		logger.ErrorWF("invalid type. not *DollEquipInfoV8Config", zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"),
			zap.String("sheet", "doll_equip_info_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gDollEquipInfoV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*DollEquipInfoV8Config)
	if !ok {
		err = errors.New("invalid type. not *DollEquipInfoV8Config")
		logger.ErrorWF("invalid type. not *DollEquipInfoV8Config", zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"),
			zap.String("sheet", "doll_equip_info_v8"))
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
type gDollEquipInfoV8Parser struct {
}

// New new config row data
func (*gDollEquipInfoV8Parser) New() interface{} {
	return &DollEquipInfoV8ConfigRow{}
}

// Fields get config fields names
func (*gDollEquipInfoV8Parser) Fields() []string {
	return gDollEquipInfoV8Fields
}

// Parse parse raw data to row data
func (*gDollEquipInfoV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*DollEquipInfoV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *DollEquipInfoV8ConfigRow")
		logger.ErrorWF("invalid type. not *DollEquipInfoV8ConfigRow", zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"),
			zap.String("sheet", "doll_equip_info_v8"))
		return
	}
	// compare length
	if len(data) != len(gDollEquipInfoV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*DollEquipInfoV8ConfigRow",
			zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"),
			zap.String("sheet", "doll_equip_info_v8"), zap.Int("need_count", len(gDollEquipInfoV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 equipment_id : 序号
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field equipment_id 序号 to int32 failed")
			logger.ErrorWF("parse field equipment_id 序号 to int32 failed.",
				zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Equipment_id = int32(tmp)
	}

	// parse column 1 score_group : 记录装备分数的组id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field score_group 记录装备分数的组id to int32 failed")
			logger.ErrorWF("parse field score_group 记录装备分数的组id to int32 failed.",
				zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Score_group = int32(tmp)
	}

	// parse column 2 quality : 品质
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field quality 品质 to int32 failed")
			logger.ErrorWF("parse field quality 品质 to int32 failed.",
				zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Quality = int32(tmp)
	}

	// parse column 3 pos : 部位
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field pos 部位 to int32 failed")
			logger.ErrorWF("parse field pos 部位 to int32 failed.",
				zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Pos = int32(tmp)
	}

	// parse column 4 skill_pos_num : 随机技能孔数量:权重
	if data[4] != "" {

		config.Skill_pos_num = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[4], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field skill_pos_num 随机技能孔数量:权重 to key int32 failed")
				logger.ErrorWF("parse map field skill_pos_num 随机技能孔数量:权重 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[4]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field skill_pos_num 随机技能孔数量:权重 to value int32 failed")
				logger.ErrorWF("parse map field skill_pos_num 随机技能孔数量:权重 to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[4]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Skill_pos_num[key] = value
		}
	}

	// parse column 5 sub_type_random : 装备子类型随机权重
	if data[5] != "" {

		config.Sub_type_random = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[5], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field sub_type_random 装备子类型随机权重 to key int32 failed")
				logger.ErrorWF("parse map field sub_type_random 装备子类型随机权重 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field sub_type_random 装备子类型随机权重 to value int32 failed")
				logger.ErrorWF("parse map field sub_type_random 装备子类型随机权重 to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Sub_type_random[key] = value
		}
	}

	// parse column 6 level : 穿戴等级
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field level 穿戴等级 to int32 failed")
			logger.ErrorWF("parse field level 穿戴等级 to int32 failed.",
				zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Level = int32(tmp)
	}

	// parse column 7 level_show : 显示等级
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field level_show 显示等级 to int32 failed")
			logger.ErrorWF("parse field level_show 显示等级 to int32 failed.",
				zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Level_show = int32(tmp)
	}

	// parse column 8 need_dungeon : 允许穿戴需要通关的副本
	if data[8] != "" {

		config.Need_dungeon = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[8], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field need_dungeon 允许穿戴需要通关的副本 to key int32 failed")
				logger.ErrorWF("parse map field need_dungeon 允许穿戴需要通关的副本 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[8]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field need_dungeon 允许穿戴需要通关的副本 to value int32 failed")
				logger.ErrorWF("parse map field need_dungeon 允许穿戴需要通关的副本 to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[8]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Need_dungeon[key] = value
		}
	}

	// parse column 9 sell : 出售价格
	if data[9] != "" {

		config.Sell = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[9], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field sell 出售价格 to key int32 failed")
				logger.ErrorWF("parse map field sell 出售价格 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[9]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field sell 出售价格 to value int64 failed")
				logger.ErrorWF("parse map field sell 出售价格 to value int64 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[9]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Sell[key] = value
		}
	}

	// parse column 10 split_num : 分解后获得材料次数：权重
	if data[10] != "" {

		config.Split_num = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[10], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field split_num 分解后获得材料次数：权重 to key int32 failed")
				logger.ErrorWF("parse map field split_num 分解后获得材料次数：权重 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[10]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field split_num 分解后获得材料次数：权重 to value int32 failed")
				logger.ErrorWF("parse map field split_num 分解后获得材料次数：权重 to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[10]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Split_num[key] = value
		}
	}

	// parse column 11 split_random_id : 分解随机库id
	if data[11] != "" {
		tmp, err = strconv.ParseInt(data[11], 10, 64)
		if err != nil {
			err = errors.New("parse field split_random_id 分解随机库id to int32 failed")
			logger.ErrorWF("parse field split_random_id 分解随机库id to int32 failed.",
				zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
				zap.String("parse_data", data[11]),
				zap.Error(err))
			return
		}
		config.Split_random_id = int32(tmp)
	}

	// parse column 12 affix_base_num : 初始基础词条数量
	if data[12] != "" {

		config.Affix_base_num = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[12], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_base_num 初始基础词条数量 to key int32 failed")
				logger.ErrorWF("parse map field affix_base_num 初始基础词条数量 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[12]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_base_num 初始基础词条数量 to value int32 failed")
				logger.ErrorWF("parse map field affix_base_num 初始基础词条数量 to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[12]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_base_num[key] = value
		}
	}

	// parse column 13 affix_base_pool : 基础词条池子id:词条位置(affix_pool
	if data[13] != "" {

		config.Affix_base_pool = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[13], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_base_pool 基础词条池子id:词条位置(affix_pool to key int32 failed")
				logger.ErrorWF("parse map field affix_base_pool 基础词条池子id:词条位置(affix_pool to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[13]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_base_pool 基础词条池子id:词条位置(affix_pool to value int32 failed")
				logger.ErrorWF("parse map field affix_base_pool 基础词条池子id:词条位置(affix_pool to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[13]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_base_pool[key] = value
		}
	}

	// parse column 14 affix_rand_num : 初始随机词条数量
	if data[14] != "" {

		config.Affix_rand_num = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[14], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_rand_num 初始随机词条数量 to key int32 failed")
				logger.ErrorWF("parse map field affix_rand_num 初始随机词条数量 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[14]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_rand_num 初始随机词条数量 to value int32 failed")
				logger.ErrorWF("parse map field affix_rand_num 初始随机词条数量 to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[14]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_rand_num[key] = value
		}
	}

	// parse column 15 affix_rand_pool : 随机词条池子id:权重(affix_pool
	if data[15] != "" {

		config.Affix_rand_pool = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[15], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_rand_pool 随机词条池子id:权重(affix_pool to key int32 failed")
				logger.ErrorWF("parse map field affix_rand_pool 随机词条池子id:权重(affix_pool to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[15]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_rand_pool 随机词条池子id:权重(affix_pool to value int32 failed")
				logger.ErrorWF("parse map field affix_rand_pool 随机词条池子id:权重(affix_pool to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[15]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_rand_pool[key] = value
		}
	}

	// parse column 16 affix_mod_num : 初始传奇词条数量
	if data[16] != "" {

		config.Affix_mod_num = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[16], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_mod_num 初始传奇词条数量 to key int32 failed")
				logger.ErrorWF("parse map field affix_mod_num 初始传奇词条数量 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[16]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_mod_num 初始传奇词条数量 to value int32 failed")
				logger.ErrorWF("parse map field affix_mod_num 初始传奇词条数量 to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[16]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_mod_num[key] = value
		}
	}

	// parse column 17 affix_mod_pool : 传奇词条池子id:权重（mod_pool
	if data[17] != "" {

		config.Affix_mod_pool = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[17], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_mod_pool 传奇词条池子id:权重（mod_pool to key int32 failed")
				logger.ErrorWF("parse map field affix_mod_pool 传奇词条池子id:权重（mod_pool to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[17]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_mod_pool 传奇词条池子id:权重（mod_pool to value int32 failed")
				logger.ErrorWF("parse map field affix_mod_pool 传奇词条池子id:权重（mod_pool to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[17]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_mod_pool[key] = value
		}
	}

	// parse column 18 affix_extra_pool : 额外词条池子id:权重（sp_pool
	if data[18] != "" {

		config.Affix_extra_pool = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[18], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_extra_pool 额外词条池子id:权重（sp_pool to key int32 failed")
				logger.ErrorWF("parse map field affix_extra_pool 额外词条池子id:权重（sp_pool to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[18]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_extra_pool 额外词条池子id:权重（sp_pool to value int32 failed")
				logger.ErrorWF("parse map field affix_extra_pool 额外词条池子id:权重（sp_pool to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[18]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_extra_pool[key] = value
		}
	}

	// parse column 19 suite_id : 装备套装id:随机权重
	if data[19] != "" {

		config.Suite_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[19], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field suite_id 装备套装id:随机权重 to key int32 failed")
				logger.ErrorWF("parse map field suite_id 装备套装id:随机权重 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[19]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field suite_id 装备套装id:随机权重 to value int32 failed")
				logger.ErrorWF("parse map field suite_id 装备套装id:随机权重 to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[19]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Suite_id[key] = value
		}
	}

	// parse column 20 unknow_view_equipment_id : 未鉴定的装备id（0-表示不需要鉴定
	if data[20] != "" {
		tmp, err = strconv.ParseInt(data[20], 10, 64)
		if err != nil {
			err = errors.New("parse field unknow_view_equipment_id 未鉴定的装备id（0-表示不需要鉴定 to int32 failed")
			logger.ErrorWF("parse field unknow_view_equipment_id 未鉴定的装备id（0-表示不需要鉴定 to int32 failed.",
				zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
				zap.String("parse_data", data[20]),
				zap.Error(err))
			return
		}
		config.Unknow_view_equipment_id = int32(tmp)
	}

	// parse column 21 view_cost : 鉴定消耗
	if data[21] != "" {

		config.View_cost = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[21], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field view_cost 鉴定消耗 to key int32 failed")
				logger.ErrorWF("parse map field view_cost 鉴定消耗 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[21]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field view_cost 鉴定消耗 to value int64 failed")
				logger.ErrorWF("parse map field view_cost 鉴定消耗 to value int64 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[21]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.View_cost[key] = value
		}
	}

	// parse column 22 view_need_dungeon : 允许鉴定需要通关的副本
	if data[22] != "" {

		config.View_need_dungeon = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[22], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field view_need_dungeon 允许鉴定需要通关的副本 to key int32 failed")
				logger.ErrorWF("parse map field view_need_dungeon 允许鉴定需要通关的副本 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[22]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field view_need_dungeon 允许鉴定需要通关的副本 to value int32 failed")
				logger.ErrorWF("parse map field view_need_dungeon 允许鉴定需要通关的副本 to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[22]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.View_need_dungeon[key] = value
		}
	}

	// parse column 23 pos_sub_type : 部位子类型
	if data[23] != "" {
		tmp, err = strconv.ParseInt(data[23], 10, 64)
		if err != nil {
			err = errors.New("parse field pos_sub_type 部位子类型 to int32 failed")
			logger.ErrorWF("parse field pos_sub_type 部位子类型 to int32 failed.",
				zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
				zap.String("parse_data", data[23]),
				zap.Error(err))
			return
		}
		config.Pos_sub_type = int32(tmp)
	}

	// parse column 24 weapon_model : 武器模型（废弃）
	if data[24] != "" {
		tmp, err = strconv.ParseInt(data[24], 10, 64)
		if err != nil {
			err = errors.New("parse field weapon_model 武器模型（废弃） to int32 failed")
			logger.ErrorWF("parse field weapon_model 武器模型（废弃） to int32 failed.",
				zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
				zap.String("parse_data", data[24]),
				zap.Error(err))
			return
		}
		config.Weapon_model = int32(tmp)
	}

	// parse column 25 id : 记录装备分数的组id
	if data[25] != "" {
		tmp, err = strconv.ParseInt(data[25], 10, 64)
		if err != nil {
			err = errors.New("parse field id 记录装备分数的组id to int32 failed")
			logger.ErrorWF("parse field id 记录装备分数的组id to int32 failed.",
				zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
				zap.String("parse_data", data[25]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 26 preview_attr : 预览属性
	if data[26] != "" {

		config.Preview_attr = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[26], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field preview_attr 预览属性 to key int32 failed")
				logger.ErrorWF("parse map field preview_attr 预览属性 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[26]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field preview_attr 预览属性 to value int64 failed")
				logger.ErrorWF("parse map field preview_attr 预览属性 to value int64 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[26]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Preview_attr[key] = value
		}
	}

	// parse column 27 attack_type : 攻击方式
	if data[27] != "" {
		tmp, err = strconv.ParseInt(data[27], 10, 64)
		if err != nil {
			err = errors.New("parse field attack_type 攻击方式 to int32 failed")
			logger.ErrorWF("parse field attack_type 攻击方式 to int32 failed.",
				zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
				zap.String("parse_data", data[27]),
				zap.Error(err))
			return
		}
		config.Attack_type = int32(tmp)
	}

	// parse column 28 damage_type : 伤害类型
	if data[28] != "" {
		tmp, err = strconv.ParseInt(data[28], 10, 64)
		if err != nil {
			err = errors.New("parse field damage_type 伤害类型 to int32 failed")
			logger.ErrorWF("parse field damage_type 伤害类型 to int32 failed.",
				zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
				zap.String("parse_data", data[28]),
				zap.Error(err))
			return
		}
		config.Damage_type = int32(tmp)
	}

	// parse column 29 font_soul : 升华头属性：随机权重
	if data[29] != "" {

		config.Font_soul = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[29], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field font_soul 升华头属性：随机权重 to key int32 failed")
				logger.ErrorWF("parse map field font_soul 升华头属性：随机权重 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[29]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field font_soul 升华头属性：随机权重 to value int32 failed")
				logger.ErrorWF("parse map field font_soul 升华头属性：随机权重 to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[29]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Font_soul[key] = value
		}
	}

	// parse column 30 tail_soul : 升华尾属性：随机权重
	if data[30] != "" {

		config.Tail_soul = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[30], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field tail_soul 升华尾属性：随机权重 to key int32 failed")
				logger.ErrorWF("parse map field tail_soul 升华尾属性：随机权重 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[30]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field tail_soul 升华尾属性：随机权重 to value int32 failed")
				logger.ErrorWF("parse map field tail_soul 升华尾属性：随机权重 to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[30]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Tail_soul[key] = value
		}
	}

	// parse column 31 recast_cost : 升华属性重铸消耗
	if data[31] != "" {

		config.Recast_cost = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[31], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field recast_cost 升华属性重铸消耗 to key int32 failed")
				logger.ErrorWF("parse map field recast_cost 升华属性重铸消耗 to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[31]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field recast_cost 升华属性重铸消耗 to value int64 failed")
				logger.ErrorWF("parse map field recast_cost 升华属性重铸消耗 to value int64 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[31]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Recast_cost[key] = value
		}
	}

	// parse column 32 soul_affix_base_pool : 升华属性（走基础词条池子id:词条位置(affix_pool
	if data[32] != "" {

		config.Soul_affix_base_pool = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[32], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field soul_affix_base_pool 升华属性（走基础词条池子id:词条位置(affix_pool to key int32 failed")
				logger.ErrorWF("parse map field soul_affix_base_pool 升华属性（走基础词条池子id:词条位置(affix_pool to key int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[32]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field soul_affix_base_pool 升华属性（走基础词条池子id:词条位置(affix_pool to value int32 failed")
				logger.ErrorWF("parse map field soul_affix_base_pool 升华属性（走基础词条池子id:词条位置(affix_pool to value int32 failed.",
					zap.String("xlsx", "doll_equip_info_v8【人偶-装备-信息】.xlsx"), zap.String("sheet", "doll_equip_info_v8"),
					// zap.String("field_data",data[32]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Soul_affix_base_pool[key] = value
		}
	}
	return
}

var gDollEquipInfoV8Fields = []string{
	"equipment_id",
	"score_group",
	"quality",
	"pos",
	"skill_pos_num",
	"sub_type_random",
	"level",
	"level_show",
	"need_dungeon",
	"sell",
	"split_num",
	"split_random_id",
	"affix_base_num",
	"affix_base_pool",
	"affix_rand_num",
	"affix_rand_pool",
	"affix_mod_num",
	"affix_mod_pool",
	"affix_extra_pool",
	"suite_id",
	"unknow_view_equipment_id",
	"view_cost",
	"view_need_dungeon",
	"pos_sub_type",
	"weapon_model",
	"id",
	"preview_attr",
	"attack_type",
	"damage_type",
	"font_soul",
	"tail_soul",
	"recast_cost",
	"soul_affix_base_pool",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gDollEquipInfoV8Parser{}
	loader := &gDollEquipInfoV8Loader{}
	var data [][]string
	data, err = load("doll_equip_info_v8【人偶-装备-信息】.xlsx", "doll_equip_info_v8", gDollEquipInfoV8Fields)
	if err != nil {
		logger.ErrorWF("load doll_equip_info_v8【人偶-装备-信息】.xlsx doll_equip_info_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load doll_equip_info_v8【人偶-装备-信息】.xlsx doll_equip_info_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gDollEquipInfoV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse doll_equip_info_v8【人偶-装备-信息】.xlsx doll_equip_info_v8 failed.", zap.Int("row", k), zap.Strings("need", gDollEquipInfoV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse doll_equip_info_v8【人偶-装备-信息】.xlsx doll_equip_info_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add doll_equip_info_v8【人偶-装备-信息】.xlsx doll_equip_info_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check doll_equip_info_v8【人偶-装备-信息】.xlsx doll_equip_info_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load doll_equip_info_v8【人偶-装备-信息】.xlsx doll_equip_info_v8 data success.")
	return
}
