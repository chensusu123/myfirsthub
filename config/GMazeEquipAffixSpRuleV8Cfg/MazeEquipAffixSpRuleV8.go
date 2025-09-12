package GMazeEquipAffixSpRuleV8Cfg

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

// MazeEquipAffixSpRuleV8ConfigRow from maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx maze_equip_affix_sp_rule_v8
type MazeEquipAffixSpRuleV8ConfigRow struct {
	Rule_id              int32           `json:"rule_id"`              // 规则id
	Sub_type_random      map[int32]int32 `json:"sub_type_random"`      // 装备子类型随机权重
	Affix_base_num       map[int32]int32 `json:"affix_base_num"`       // 初始基础词条数量
	Affix_base_pool      map[int32]int32 `json:"affix_base_pool"`      // 基础词条池:权重（doll_equip_affix_pool_v8）
	Affix_rand_num       map[int32]int32 `json:"affix_rand_num"`       // 初始随机词条数量
	Affix_rand_pool      map[int32]int32 `json:"affix_rand_pool"`      // 随机词条:权重（doll_equip_affix_pool_v8）
	Affix_mod_num        map[int32]int32 `json:"affix_mod_num"`        // 初始附魔词条数量
	Affix_mod_pool       map[int32]int32 `json:"affix_mod_pool"`       // 附魔词条:权重
	Affix_extra_pool     map[int32]int32 `json:"affix_extra_pool"`     // 额外词条池子id:权重（sp_pool
	Suite_id             map[int32]int32 `json:"suite_id"`             // 装备套装id:随机权重
	Font_soul            map[int32]int32 `json:"font_soul"`            // 升华头属性：随机权重
	Tail_soul            map[int32]int32 `json:"tail_soul"`            // 升华尾属性：随机权重
	Soul_affix_base_pool map[int32]int32 `json:"soul_affix_base_pool"` // 升华属性（走基础词条池子id:词条位置(affix_pool
}

// MazeEquipAffixSpRuleV8Config from maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx maze_equip_affix_sp_rule_v8
type MazeEquipAffixSpRuleV8Config struct {
	ConfigRows map[int32]*MazeEquipAffixSpRuleV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipAffixSpRuleV8Config {
	ret := &MazeEquipAffixSpRuleV8Config{ConfigRows: map[int32]*MazeEquipAffixSpRuleV8ConfigRow{}}
	return ret
}

// GetMazeEquipAffixSpRuleV8Config get one config by configId
func (c *MazeEquipAffixSpRuleV8Config) GetMazeEquipAffixSpRuleV8Config(configId int32) *MazeEquipAffixSpRuleV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipAffixSpRuleV8Config) Get(configId int32) *MazeEquipAffixSpRuleV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipAffixSpRuleV8Config get all config slice
func (c *MazeEquipAffixSpRuleV8Config) GetAllMazeEquipAffixSpRuleV8Config() (res []*MazeEquipAffixSpRuleV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipAffixSpRuleV8Config) GetAll() (res []*MazeEquipAffixSpRuleV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipAffixSpRuleV8Config

// GetMazeEquipAffixSpRuleV8Config pkg func. get one config by configId
func GetMazeEquipAffixSpRuleV8Config(configId int32) *MazeEquipAffixSpRuleV8ConfigRow {
	return gConfigData.GetMazeEquipAffixSpRuleV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipAffixSpRuleV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeEquipAffixSpRuleV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_equip_affix_sp_rule_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeEquipAffixSpRuleV8Config pkg func. get all config slice
func GetAllMazeEquipAffixSpRuleV8Config() []*MazeEquipAffixSpRuleV8ConfigRow {
	return gConfigData.GetAllMazeEquipAffixSpRuleV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipAffixSpRuleV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipAffixSpRuleV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipAffixSpRuleV8ConfigRow from maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx maze_equip_affix_sp_rule_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipAffixSpRuleV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_affix_sp_rule_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_affix_sp_rule_v8.json",
		"maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx", "maze_equip_affix_sp_rule_v8",
		&gMazeEquipAffixSpRuleV8Parser{}, &gMazeEquipAffixSpRuleV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipAffixSpRuleV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipAffixSpRuleV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipAffixSpRuleV8Config))(c)
		return true
	})
}

// RegisterMazeEquipAffixSpRuleV8InitCallBack reg config update func (old func)
var RegisterMazeEquipAffixSpRuleV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipAffixSpRuleV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipAffixSpRuleV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipAffixSpRuleV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipAffixSpRuleV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipAffixSpRuleV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipAffixSpRuleV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipAffixSpRuleV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipAffixSpRuleV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipAffixSpRuleV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipAffixSpRuleV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipAffixSpRuleV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipAffixSpRuleV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipAffixSpRuleV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixSpRuleV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipAffixSpRuleV8ConfigRow", zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"),
			zap.String("sheet", "maze_equip_affix_sp_rule_v8"))
		return
	}
	config, ok := container.(*MazeEquipAffixSpRuleV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixSpRuleV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixSpRuleV8Config", zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"),
			zap.String("sheet", "maze_equip_affix_sp_rule_v8"))
		return
	}
	config.ConfigRows[row.Rule_id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipAffixSpRuleV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipAffixSpRuleV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixSpRuleV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixSpRuleV8Config", zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"),
			zap.String("sheet", "maze_equip_affix_sp_rule_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipAffixSpRuleV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipAffixSpRuleV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixSpRuleV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixSpRuleV8Config", zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"),
			zap.String("sheet", "maze_equip_affix_sp_rule_v8"))
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
type gMazeEquipAffixSpRuleV8Parser struct {
}

// New new config row data
func (*gMazeEquipAffixSpRuleV8Parser) New() interface{} {
	return &MazeEquipAffixSpRuleV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipAffixSpRuleV8Parser) Fields() []string {
	return gMazeEquipAffixSpRuleV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipAffixSpRuleV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipAffixSpRuleV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixSpRuleV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipAffixSpRuleV8ConfigRow", zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"),
			zap.String("sheet", "maze_equip_affix_sp_rule_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipAffixSpRuleV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipAffixSpRuleV8ConfigRow",
			zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"),
			zap.String("sheet", "maze_equip_affix_sp_rule_v8"), zap.Int("need_count", len(gMazeEquipAffixSpRuleV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 rule_id : 规则id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field rule_id 规则id to int32 failed")
			logger.ErrorWF("parse field rule_id 规则id to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Rule_id = int32(tmp)
	}

	// parse column 1 sub_type_random : 装备子类型随机权重
	if data[1] != "" {

		config.Sub_type_random = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[1], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field sub_type_random 装备子类型随机权重 to key int32 failed")
				logger.ErrorWF("parse map field sub_type_random 装备子类型随机权重 to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[1]),
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
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Sub_type_random[key] = value
		}
	}

	// parse column 2 affix_base_num : 初始基础词条数量
	if data[2] != "" {

		config.Affix_base_num = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[2], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_base_num 初始基础词条数量 to key int32 failed")
				logger.ErrorWF("parse map field affix_base_num 初始基础词条数量 to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[2]),
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
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[2]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_base_num[key] = value
		}
	}

	// parse column 3 affix_base_pool : 基础词条池:权重（doll_equip_affix_pool_v8）
	if data[3] != "" {

		config.Affix_base_pool = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_base_pool 基础词条池:权重（doll_equip_affix_pool_v8） to key int32 failed")
				logger.ErrorWF("parse map field affix_base_pool 基础词条池:权重（doll_equip_affix_pool_v8） to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_base_pool 基础词条池:权重（doll_equip_affix_pool_v8） to value int32 failed")
				logger.ErrorWF("parse map field affix_base_pool 基础词条池:权重（doll_equip_affix_pool_v8） to value int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_base_pool[key] = value
		}
	}

	// parse column 4 affix_rand_num : 初始随机词条数量
	if data[4] != "" {

		config.Affix_rand_num = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[4], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_rand_num 初始随机词条数量 to key int32 failed")
				logger.ErrorWF("parse map field affix_rand_num 初始随机词条数量 to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[4]),
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
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[4]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_rand_num[key] = value
		}
	}

	// parse column 5 affix_rand_pool : 随机词条:权重（doll_equip_affix_pool_v8）
	if data[5] != "" {

		config.Affix_rand_pool = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[5], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_rand_pool 随机词条:权重（doll_equip_affix_pool_v8） to key int32 failed")
				logger.ErrorWF("parse map field affix_rand_pool 随机词条:权重（doll_equip_affix_pool_v8） to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_rand_pool 随机词条:权重（doll_equip_affix_pool_v8） to value int32 failed")
				logger.ErrorWF("parse map field affix_rand_pool 随机词条:权重（doll_equip_affix_pool_v8） to value int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_rand_pool[key] = value
		}
	}

	// parse column 6 affix_mod_num : 初始附魔词条数量
	if data[6] != "" {

		config.Affix_mod_num = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[6], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_mod_num 初始附魔词条数量 to key int32 failed")
				logger.ErrorWF("parse map field affix_mod_num 初始附魔词条数量 to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[6]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_mod_num 初始附魔词条数量 to value int32 failed")
				logger.ErrorWF("parse map field affix_mod_num 初始附魔词条数量 to value int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[6]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_mod_num[key] = value
		}
	}

	// parse column 7 affix_mod_pool : 附魔词条:权重
	if data[7] != "" {

		config.Affix_mod_pool = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[7], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_mod_pool 附魔词条:权重 to key int32 failed")
				logger.ErrorWF("parse map field affix_mod_pool 附魔词条:权重 to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_mod_pool 附魔词条:权重 to value int32 failed")
				logger.ErrorWF("parse map field affix_mod_pool 附魔词条:权重 to value int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_mod_pool[key] = value
		}
	}

	// parse column 8 affix_extra_pool : 额外词条池子id:权重（sp_pool
	if data[8] != "" {

		config.Affix_extra_pool = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[8], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_extra_pool 额外词条池子id:权重（sp_pool to key int32 failed")
				logger.ErrorWF("parse map field affix_extra_pool 额外词条池子id:权重（sp_pool to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[8]),
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
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[8]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_extra_pool[key] = value
		}
	}

	// parse column 9 suite_id : 装备套装id:随机权重
	if data[9] != "" {

		config.Suite_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[9], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field suite_id 装备套装id:随机权重 to key int32 failed")
				logger.ErrorWF("parse map field suite_id 装备套装id:随机权重 to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[9]),
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
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[9]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Suite_id[key] = value
		}
	}

	// parse column 10 font_soul : 升华头属性：随机权重
	if data[10] != "" {

		config.Font_soul = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[10], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field font_soul 升华头属性：随机权重 to key int32 failed")
				logger.ErrorWF("parse map field font_soul 升华头属性：随机权重 to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[10]),
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
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[10]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Font_soul[key] = value
		}
	}

	// parse column 11 tail_soul : 升华尾属性：随机权重
	if data[11] != "" {

		config.Tail_soul = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[11], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field tail_soul 升华尾属性：随机权重 to key int32 failed")
				logger.ErrorWF("parse map field tail_soul 升华尾属性：随机权重 to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[11]),
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
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[11]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Tail_soul[key] = value
		}
	}

	// parse column 12 soul_affix_base_pool : 升华属性（走基础词条池子id:词条位置(affix_pool
	if data[12] != "" {

		config.Soul_affix_base_pool = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[12], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field soul_affix_base_pool 升华属性（走基础词条池子id:词条位置(affix_pool to key int32 failed")
				logger.ErrorWF("parse map field soul_affix_base_pool 升华属性（走基础词条池子id:词条位置(affix_pool to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[12]),
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
					zap.String("xlsx", "maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx"), zap.String("sheet", "maze_equip_affix_sp_rule_v8"),
					// zap.String("field_data",data[12]),
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

var gMazeEquipAffixSpRuleV8Fields = []string{
	"rule_id",
	"sub_type_random",
	"affix_base_num",
	"affix_base_pool",
	"affix_rand_num",
	"affix_rand_pool",
	"affix_mod_num",
	"affix_mod_pool",
	"affix_extra_pool",
	"suite_id",
	"font_soul",
	"tail_soul",
	"soul_affix_base_pool",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipAffixSpRuleV8Parser{}
	loader := &gMazeEquipAffixSpRuleV8Loader{}
	var data [][]string
	data, err = load("maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx", "maze_equip_affix_sp_rule_v8", gMazeEquipAffixSpRuleV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx maze_equip_affix_sp_rule_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx maze_equip_affix_sp_rule_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipAffixSpRuleV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx maze_equip_affix_sp_rule_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipAffixSpRuleV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx maze_equip_affix_sp_rule_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx maze_equip_affix_sp_rule_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx maze_equip_affix_sp_rule_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_affix_sp_rule_v8【迷宫-装备-生成特殊词条规则】.xlsx maze_equip_affix_sp_rule_v8 data success.")
	return
}
