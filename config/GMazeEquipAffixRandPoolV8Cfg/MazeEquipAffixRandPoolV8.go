package GMazeEquipAffixRandPoolV8Cfg

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

// MazeEquipAffixRandPoolV8ConfigRow from maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx maze_equip_affix_rand_pool_v8
type MazeEquipAffixRandPoolV8ConfigRow struct {
	Affix_id      int32           `json:"affix_id"`      // 词条id
	Pool_id       int32           `json:"pool_id"`       // 池子id
	Pool_rank     int32           `json:"pool_rank"`     // 池子内编号
	Weight        int32           `json:"weight"`        // 随机权重
	Group         int32           `json:"group"`         // 排重组（同id去重）
	Add_attr_min  map[int32]int64 `json:"add_attr_min"`  // 实际增加属性
	Add_attr_max  map[int32]int64 `json:"add_attr_max"`  // 实际增加属性
	Show_attr_min map[int32]int64 `json:"show_attr_min"` // 展示属性id
	Show_attr_max map[int32]int64 `json:"show_attr_max"` // 展示属性id
	Roll_type     int32           `json:"roll_type"`     // roll值类型
}

// MazeEquipAffixRandPoolV8Config from maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx maze_equip_affix_rand_pool_v8
type MazeEquipAffixRandPoolV8Config struct {
	ConfigRows map[int32]*MazeEquipAffixRandPoolV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipAffixRandPoolV8Config {
	ret := &MazeEquipAffixRandPoolV8Config{ConfigRows: map[int32]*MazeEquipAffixRandPoolV8ConfigRow{}}
	return ret
}

// GetMazeEquipAffixRandPoolV8Config get one config by configId
func (c *MazeEquipAffixRandPoolV8Config) GetMazeEquipAffixRandPoolV8Config(configId int32) *MazeEquipAffixRandPoolV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipAffixRandPoolV8Config) Get(configId int32) *MazeEquipAffixRandPoolV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipAffixRandPoolV8Config get all config slice
func (c *MazeEquipAffixRandPoolV8Config) GetAllMazeEquipAffixRandPoolV8Config() (res []*MazeEquipAffixRandPoolV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipAffixRandPoolV8Config) GetAll() (res []*MazeEquipAffixRandPoolV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipAffixRandPoolV8Config

// GetMazeEquipAffixRandPoolV8Config pkg func. get one config by configId
func GetMazeEquipAffixRandPoolV8Config(configId int32) *MazeEquipAffixRandPoolV8ConfigRow {
	return gConfigData.GetMazeEquipAffixRandPoolV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipAffixRandPoolV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeEquipAffixRandPoolV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_equip_affix_rand_pool_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeEquipAffixRandPoolV8Config pkg func. get all config slice
func GetAllMazeEquipAffixRandPoolV8Config() []*MazeEquipAffixRandPoolV8ConfigRow {
	return gConfigData.GetAllMazeEquipAffixRandPoolV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipAffixRandPoolV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipAffixRandPoolV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipAffixRandPoolV8ConfigRow from maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx maze_equip_affix_rand_pool_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipAffixRandPoolV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_affix_rand_pool_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_affix_rand_pool_v8.json",
		"maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx", "maze_equip_affix_rand_pool_v8",
		&gMazeEquipAffixRandPoolV8Parser{}, &gMazeEquipAffixRandPoolV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipAffixRandPoolV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipAffixRandPoolV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipAffixRandPoolV8Config))(c)
		return true
	})
}

// RegisterMazeEquipAffixRandPoolV8InitCallBack reg config update func (old func)
var RegisterMazeEquipAffixRandPoolV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipAffixRandPoolV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipAffixRandPoolV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipAffixRandPoolV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipAffixRandPoolV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipAffixRandPoolV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipAffixRandPoolV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipAffixRandPoolV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipAffixRandPoolV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipAffixRandPoolV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipAffixRandPoolV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipAffixRandPoolV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipAffixRandPoolV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipAffixRandPoolV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixRandPoolV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipAffixRandPoolV8ConfigRow", zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"),
			zap.String("sheet", "maze_equip_affix_rand_pool_v8"))
		return
	}
	config, ok := container.(*MazeEquipAffixRandPoolV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixRandPoolV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixRandPoolV8Config", zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"),
			zap.String("sheet", "maze_equip_affix_rand_pool_v8"))
		return
	}
	config.ConfigRows[row.Affix_id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipAffixRandPoolV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipAffixRandPoolV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixRandPoolV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixRandPoolV8Config", zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"),
			zap.String("sheet", "maze_equip_affix_rand_pool_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipAffixRandPoolV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipAffixRandPoolV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixRandPoolV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixRandPoolV8Config", zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"),
			zap.String("sheet", "maze_equip_affix_rand_pool_v8"))
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
type gMazeEquipAffixRandPoolV8Parser struct {
}

// New new config row data
func (*gMazeEquipAffixRandPoolV8Parser) New() interface{} {
	return &MazeEquipAffixRandPoolV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipAffixRandPoolV8Parser) Fields() []string {
	return gMazeEquipAffixRandPoolV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipAffixRandPoolV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipAffixRandPoolV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixRandPoolV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipAffixRandPoolV8ConfigRow", zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"),
			zap.String("sheet", "maze_equip_affix_rand_pool_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipAffixRandPoolV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipAffixRandPoolV8ConfigRow",
			zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"),
			zap.String("sheet", "maze_equip_affix_rand_pool_v8"), zap.Int("need_count", len(gMazeEquipAffixRandPoolV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 affix_id : 词条id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field affix_id 词条id to int32 failed")
			logger.ErrorWF("parse field affix_id 词条id to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Affix_id = int32(tmp)
	}

	// parse column 1 pool_id : 池子id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field pool_id 池子id to int32 failed")
			logger.ErrorWF("parse field pool_id 池子id to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Pool_id = int32(tmp)
	}

	// parse column 2 pool_rank : 池子内编号
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field pool_rank 池子内编号 to int32 failed")
			logger.ErrorWF("parse field pool_rank 池子内编号 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Pool_rank = int32(tmp)
	}

	// parse column 3 weight : 随机权重
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field weight 随机权重 to int32 failed")
			logger.ErrorWF("parse field weight 随机权重 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Weight = int32(tmp)
	}

	// parse column 4 group : 排重组（同id去重）
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field group 排重组（同id去重） to int32 failed")
			logger.ErrorWF("parse field group 排重组（同id去重） to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Group = int32(tmp)
	}

	// parse column 5 add_attr_min : 实际增加属性
	if data[5] != "" {

		config.Add_attr_min = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[5], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field add_attr_min 实际增加属性 to key int32 failed")
				logger.ErrorWF("parse map field add_attr_min 实际增加属性 to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field add_attr_min 实际增加属性 to value int64 failed")
				logger.ErrorWF("parse map field add_attr_min 实际增加属性 to value int64 failed.",
					zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
					// zap.String("field_data",data[5]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Add_attr_min[key] = value
		}
	}

	// parse column 6 add_attr_max : 实际增加属性
	if data[6] != "" {

		config.Add_attr_max = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[6], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field add_attr_max 实际增加属性 to key int32 failed")
				logger.ErrorWF("parse map field add_attr_max 实际增加属性 to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
					// zap.String("field_data",data[6]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field add_attr_max 实际增加属性 to value int64 failed")
				logger.ErrorWF("parse map field add_attr_max 实际增加属性 to value int64 failed.",
					zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
					// zap.String("field_data",data[6]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Add_attr_max[key] = value
		}
	}

	// parse column 7 show_attr_min : 展示属性id
	if data[7] != "" {

		config.Show_attr_min = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[7], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field show_attr_min 展示属性id to key int32 failed")
				logger.ErrorWF("parse map field show_attr_min 展示属性id to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field show_attr_min 展示属性id to value int64 failed")
				logger.ErrorWF("parse map field show_attr_min 展示属性id to value int64 failed.",
					zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Show_attr_min[key] = value
		}
	}

	// parse column 8 show_attr_max : 展示属性id
	if data[8] != "" {

		config.Show_attr_max = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[8], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field show_attr_max 展示属性id to key int32 failed")
				logger.ErrorWF("parse map field show_attr_max 展示属性id to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
					// zap.String("field_data",data[8]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field show_attr_max 展示属性id to value int64 failed")
				logger.ErrorWF("parse map field show_attr_max 展示属性id to value int64 failed.",
					zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
					// zap.String("field_data",data[8]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Show_attr_max[key] = value
		}
	}

	// parse column 9 roll_type : roll值类型
	if data[9] != "" {
		tmp, err = strconv.ParseInt(data[9], 10, 64)
		if err != nil {
			err = errors.New("parse field roll_type roll值类型 to int32 failed")
			logger.ErrorWF("parse field roll_type roll值类型 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx"), zap.String("sheet", "maze_equip_affix_rand_pool_v8"),
				zap.String("parse_data", data[9]),
				zap.Error(err))
			return
		}
		config.Roll_type = int32(tmp)
	}
	return
}

var gMazeEquipAffixRandPoolV8Fields = []string{
	"affix_id",
	"pool_id",
	"pool_rank",
	"weight",
	"group",
	"add_attr_min",
	"add_attr_max",
	"show_attr_min",
	"show_attr_max",
	"roll_type",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipAffixRandPoolV8Parser{}
	loader := &gMazeEquipAffixRandPoolV8Loader{}
	var data [][]string
	data, err = load("maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx", "maze_equip_affix_rand_pool_v8", gMazeEquipAffixRandPoolV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx maze_equip_affix_rand_pool_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx maze_equip_affix_rand_pool_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipAffixRandPoolV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx maze_equip_affix_rand_pool_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipAffixRandPoolV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx maze_equip_affix_rand_pool_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx maze_equip_affix_rand_pool_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx maze_equip_affix_rand_pool_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_affix_rand_pool_v8【迷宫-装备-词条随机池】.xlsx maze_equip_affix_rand_pool_v8 data success.")
	return
}
