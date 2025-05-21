package GDollMazeFoeV8Cfg

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

// DollMazeFoeV8ConfigRow from doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx doll_maze_foe_v8
type DollMazeFoeV8ConfigRow struct {
	Order            int32           `json:"order"`            // 怪物id
	Monster_type     int32           `json:"monster_type"`     // 怪物类型1:小怪，2:守卫
	Drop_item        map[int32]int64 `json:"drop_item"`        // 怪物掉落物品
	Drop_equip       map[int32]int64 `json:"drop_equip"`       // 怪物掉落装备
	Name             string          `json:"name"`             // 怪物名称
	Model_id         int32           `json:"model_id"`         // 资源组id
	Kongfu           int32           `json:"kongfu"`           // 怪物武力值
	Hp_max           int32           `json:"hp_max"`           // 怪物展示血量
	Hp_lose_type     int32           `json:"hp_lose_type"`     // 损血类型
	Drop_equip_score int32           `json:"drop_equip_score"` // 掉落装备分数
}

// DollMazeFoeV8Config from doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx doll_maze_foe_v8
type DollMazeFoeV8Config struct {
	ConfigRows map[int32]*DollMazeFoeV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *DollMazeFoeV8Config {
	ret := &DollMazeFoeV8Config{ConfigRows: map[int32]*DollMazeFoeV8ConfigRow{}}
	return ret
}

// GetDollMazeFoeV8Config get one config by configId
func (c *DollMazeFoeV8Config) GetDollMazeFoeV8Config(configId int32) *DollMazeFoeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *DollMazeFoeV8Config) Get(configId int32) *DollMazeFoeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllDollMazeFoeV8Config get all config slice
func (c *DollMazeFoeV8Config) GetAllDollMazeFoeV8Config() (res []*DollMazeFoeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *DollMazeFoeV8Config) GetAll() (res []*DollMazeFoeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *DollMazeFoeV8Config

// GetDollMazeFoeV8Config pkg func. get one config by configId
func GetDollMazeFoeV8Config(configId int32) *DollMazeFoeV8ConfigRow {
	return gConfigData.GetDollMazeFoeV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *DollMazeFoeV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllDollMazeFoeV8Config pkg func. get all config slice
func GetAllDollMazeFoeV8Config() []*DollMazeFoeV8ConfigRow {
	return gConfigData.GetAllDollMazeFoeV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*DollMazeFoeV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*DollMazeFoeV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "DollMazeFoeV8ConfigRow from doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx doll_maze_foe_v8"
}

// GetRawValue get raw data
func GetRawValue() *DollMazeFoeV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "doll_maze_foe_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("doll_maze_foe_v8.json",
		"doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx", "doll_maze_foe_v8",
		&gDollMazeFoeV8Parser{}, &gDollMazeFoeV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*DollMazeFoeV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *DollMazeFoeV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*DollMazeFoeV8Config))(c)
		return true
	})
}

// RegisterDollMazeFoeV8InitCallBack reg config update func (old func)
var RegisterDollMazeFoeV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*DollMazeFoeV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*DollMazeFoeV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *DollMazeFoeV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*DollMazeFoeV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*DollMazeFoeV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gDollMazeFoeV8Loader struct {
}

// NewContainer new data container pointer
func (*gDollMazeFoeV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gDollMazeFoeV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*DollMazeFoeV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gDollMazeFoeV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*DollMazeFoeV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gDollMazeFoeV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*DollMazeFoeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *DollMazeFoeV8ConfigRow")
		logger.ErrorWF("invalid type. not *DollMazeFoeV8ConfigRow", zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "doll_maze_foe_v8"))
		return
	}
	config, ok := container.(*DollMazeFoeV8Config)
	if !ok {
		err = errors.New("invalid type. not *DollMazeFoeV8Config")
		logger.ErrorWF("invalid type. not *DollMazeFoeV8Config", zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "doll_maze_foe_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gDollMazeFoeV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*DollMazeFoeV8Config)
	if !ok {
		err = errors.New("invalid type. not *DollMazeFoeV8Config")
		logger.ErrorWF("invalid type. not *DollMazeFoeV8Config", zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "doll_maze_foe_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gDollMazeFoeV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*DollMazeFoeV8Config)
	if !ok {
		err = errors.New("invalid type. not *DollMazeFoeV8Config")
		logger.ErrorWF("invalid type. not *DollMazeFoeV8Config", zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "doll_maze_foe_v8"))
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
type gDollMazeFoeV8Parser struct {
}

// New new config row data
func (*gDollMazeFoeV8Parser) New() interface{} {
	return &DollMazeFoeV8ConfigRow{}
}

// Fields get config fields names
func (*gDollMazeFoeV8Parser) Fields() []string {
	return gDollMazeFoeV8Fields
}

// Parse parse raw data to row data
func (*gDollMazeFoeV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*DollMazeFoeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *DollMazeFoeV8ConfigRow")
		logger.ErrorWF("invalid type. not *DollMazeFoeV8ConfigRow", zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "doll_maze_foe_v8"))
		return
	}
	// compare length
	if len(data) != len(gDollMazeFoeV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*DollMazeFoeV8ConfigRow",
			zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "doll_maze_foe_v8"), zap.Int("need_count", len(gDollMazeFoeV8Fields)),
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
				zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"), zap.String("sheet", "doll_maze_foe_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 monster_type : 怪物类型1:小怪，2:守卫
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field monster_type 怪物类型1:小怪，2:守卫 to int32 failed")
			logger.ErrorWF("parse field monster_type 怪物类型1:小怪，2:守卫 to int32 failed.",
				zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"), zap.String("sheet", "doll_maze_foe_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Monster_type = int32(tmp)
	}

	// parse column 2 drop_item : 怪物掉落物品
	if data[2] != "" {

		config.Drop_item = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[2], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_item 怪物掉落物品 to key int32 failed")
				logger.ErrorWF("parse map field drop_item 怪物掉落物品 to key int32 failed.",
					zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"), zap.String("sheet", "doll_maze_foe_v8"),
					// zap.String("field_data",data[2]),
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
					zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"), zap.String("sheet", "doll_maze_foe_v8"),
					// zap.String("field_data",data[2]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Drop_item[key] = value
		}
	}

	// parse column 3 drop_equip : 怪物掉落装备
	if data[3] != "" {

		config.Drop_equip = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_equip 怪物掉落装备 to key int32 failed")
				logger.ErrorWF("parse map field drop_equip 怪物掉落装备 to key int32 failed.",
					zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"), zap.String("sheet", "doll_maze_foe_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_equip 怪物掉落装备 to value int64 failed")
				logger.ErrorWF("parse map field drop_equip 怪物掉落装备 to value int64 failed.",
					zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"), zap.String("sheet", "doll_maze_foe_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Drop_equip[key] = value
		}
	}

	// parse column 4 name : 怪物名称
	if data[4] != "" {
		config.Name = data[4]
	}

	// parse column 5 model_id : 资源组id
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field model_id 资源组id to int32 failed")
			logger.ErrorWF("parse field model_id 资源组id to int32 failed.",
				zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"), zap.String("sheet", "doll_maze_foe_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Model_id = int32(tmp)
	}

	// parse column 6 kongfu : 怪物武力值
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field kongfu 怪物武力值 to int32 failed")
			logger.ErrorWF("parse field kongfu 怪物武力值 to int32 failed.",
				zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"), zap.String("sheet", "doll_maze_foe_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Kongfu = int32(tmp)
	}

	// parse column 7 hp_max : 怪物展示血量
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field hp_max 怪物展示血量 to int32 failed")
			logger.ErrorWF("parse field hp_max 怪物展示血量 to int32 failed.",
				zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"), zap.String("sheet", "doll_maze_foe_v8"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Hp_max = int32(tmp)
	}

	// parse column 8 hp_lose_type : 损血类型
	if data[8] != "" {
		tmp, err = strconv.ParseInt(data[8], 10, 64)
		if err != nil {
			err = errors.New("parse field hp_lose_type 损血类型 to int32 failed")
			logger.ErrorWF("parse field hp_lose_type 损血类型 to int32 failed.",
				zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"), zap.String("sheet", "doll_maze_foe_v8"),
				zap.String("parse_data", data[8]),
				zap.Error(err))
			return
		}
		config.Hp_lose_type = int32(tmp)
	}

	// parse column 9 drop_equip_score : 掉落装备分数
	if data[9] != "" {
		tmp, err = strconv.ParseInt(data[9], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_equip_score 掉落装备分数 to int32 failed")
			logger.ErrorWF("parse field drop_equip_score 掉落装备分数 to int32 failed.",
				zap.String("xlsx", "doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx"), zap.String("sheet", "doll_maze_foe_v8"),
				zap.String("parse_data", data[9]),
				zap.Error(err))
			return
		}
		config.Drop_equip_score = int32(tmp)
	}
	return
}

var gDollMazeFoeV8Fields = []string{
	"order",
	"monster_type",
	"drop_item",
	"drop_equip",
	"name",
	"model_id",
	"kongfu",
	"hp_max",
	"hp_lose_type",
	"drop_equip_score",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gDollMazeFoeV8Parser{}
	loader := &gDollMazeFoeV8Loader{}
	var data [][]string
	data, err = load("doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx", "doll_maze_foe_v8", gDollMazeFoeV8Fields)
	if err != nil {
		logger.ErrorWF("load doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx doll_maze_foe_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx doll_maze_foe_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gDollMazeFoeV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx doll_maze_foe_v8 failed.", zap.Int("row", k), zap.Strings("need", gDollMazeFoeV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx doll_maze_foe_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx doll_maze_foe_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx doll_maze_foe_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load doll_maze_foe_v8【人偶-迷宫-敌人信息】.xlsx doll_maze_foe_v8 data success.")
	return
}
