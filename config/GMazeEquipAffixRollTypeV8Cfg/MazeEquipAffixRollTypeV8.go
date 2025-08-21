package GMazeEquipAffixRollTypeV8Cfg

import (
	"context"
	"errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"
)

// MazeEquipAffixRollTypeV8ConfigRow from maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx maze_equip_affix_roll_type_v8
type MazeEquipAffixRollTypeV8ConfigRow struct {
	Order          int32 `json:"order"`          // 序号
	Roll_type      int32 `json:"roll_type"`      // roll值类型
	Get_score_min  int32 `json:"get_score_min"`  // 装备获得分数范围最小值
	Get_score_max  int32 `json:"get_score_max"`  // 装备获得分数范围最大值
	Weight         int32 `json:"weight"`         // 随机范围权重
	Roll_range_min int32 `json:"roll_range_min"` // 范围万分比，最小值
	Roll_range_max int32 `json:"roll_range_max"` // 范围万分比，最大值
	Round_value    int32 `json:"round_value"`    // 取整数，（随机数/取整数）向上取整*取整数
}

// MazeEquipAffixRollTypeV8Config from maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx maze_equip_affix_roll_type_v8
type MazeEquipAffixRollTypeV8Config struct {
	ConfigRows map[int32]*MazeEquipAffixRollTypeV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipAffixRollTypeV8Config {
	ret := &MazeEquipAffixRollTypeV8Config{ConfigRows: map[int32]*MazeEquipAffixRollTypeV8ConfigRow{}}
	return ret
}

// GetMazeEquipAffixRollTypeV8Config get one config by configId
func (c *MazeEquipAffixRollTypeV8Config) GetMazeEquipAffixRollTypeV8Config(configId int32) *MazeEquipAffixRollTypeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipAffixRollTypeV8Config) Get(configId int32) *MazeEquipAffixRollTypeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipAffixRollTypeV8Config get all config slice
func (c *MazeEquipAffixRollTypeV8Config) GetAllMazeEquipAffixRollTypeV8Config() (res []*MazeEquipAffixRollTypeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipAffixRollTypeV8Config) GetAll() (res []*MazeEquipAffixRollTypeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipAffixRollTypeV8Config

// GetMazeEquipAffixRollTypeV8Config pkg func. get one config by configId
func GetMazeEquipAffixRollTypeV8Config(configId int32) *MazeEquipAffixRollTypeV8ConfigRow {
	return gConfigData.GetMazeEquipAffixRollTypeV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipAffixRollTypeV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeEquipAffixRollTypeV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_equip_affix_roll_type_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeEquipAffixRollTypeV8Config pkg func. get all config slice
func GetAllMazeEquipAffixRollTypeV8Config() []*MazeEquipAffixRollTypeV8ConfigRow {
	return gConfigData.GetAllMazeEquipAffixRollTypeV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipAffixRollTypeV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipAffixRollTypeV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipAffixRollTypeV8ConfigRow from maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx maze_equip_affix_roll_type_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipAffixRollTypeV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_affix_roll_type_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_affix_roll_type_v8.json",
		"maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx", "maze_equip_affix_roll_type_v8",
		&gMazeEquipAffixRollTypeV8Parser{}, &gMazeEquipAffixRollTypeV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipAffixRollTypeV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipAffixRollTypeV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipAffixRollTypeV8Config))(c)
		return true
	})
}

// RegisterMazeEquipAffixRollTypeV8InitCallBack reg config update func (old func)
var RegisterMazeEquipAffixRollTypeV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipAffixRollTypeV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipAffixRollTypeV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipAffixRollTypeV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipAffixRollTypeV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipAffixRollTypeV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipAffixRollTypeV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipAffixRollTypeV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipAffixRollTypeV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipAffixRollTypeV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipAffixRollTypeV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipAffixRollTypeV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipAffixRollTypeV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipAffixRollTypeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixRollTypeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipAffixRollTypeV8ConfigRow", zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"),
			zap.String("sheet", "maze_equip_affix_roll_type_v8"))
		return
	}
	config, ok := container.(*MazeEquipAffixRollTypeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixRollTypeV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixRollTypeV8Config", zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"),
			zap.String("sheet", "maze_equip_affix_roll_type_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipAffixRollTypeV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipAffixRollTypeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixRollTypeV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixRollTypeV8Config", zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"),
			zap.String("sheet", "maze_equip_affix_roll_type_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipAffixRollTypeV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipAffixRollTypeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixRollTypeV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixRollTypeV8Config", zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"),
			zap.String("sheet", "maze_equip_affix_roll_type_v8"))
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
type gMazeEquipAffixRollTypeV8Parser struct {
}

// New new config row data
func (*gMazeEquipAffixRollTypeV8Parser) New() interface{} {
	return &MazeEquipAffixRollTypeV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipAffixRollTypeV8Parser) Fields() []string {
	return gMazeEquipAffixRollTypeV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipAffixRollTypeV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipAffixRollTypeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixRollTypeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipAffixRollTypeV8ConfigRow", zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"),
			zap.String("sheet", "maze_equip_affix_roll_type_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipAffixRollTypeV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipAffixRollTypeV8ConfigRow",
			zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"),
			zap.String("sheet", "maze_equip_affix_roll_type_v8"), zap.Int("need_count", len(gMazeEquipAffixRollTypeV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 序号
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 序号 to int32 failed")
			logger.ErrorWF("parse field order 序号 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"), zap.String("sheet", "maze_equip_affix_roll_type_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 roll_type : roll值类型
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field roll_type roll值类型 to int32 failed")
			logger.ErrorWF("parse field roll_type roll值类型 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"), zap.String("sheet", "maze_equip_affix_roll_type_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Roll_type = int32(tmp)
	}

	// parse column 2 get_score_min : 装备获得分数范围最小值
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field get_score_min 装备获得分数范围最小值 to int32 failed")
			logger.ErrorWF("parse field get_score_min 装备获得分数范围最小值 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"), zap.String("sheet", "maze_equip_affix_roll_type_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Get_score_min = int32(tmp)
	}

	// parse column 3 get_score_max : 装备获得分数范围最大值
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field get_score_max 装备获得分数范围最大值 to int32 failed")
			logger.ErrorWF("parse field get_score_max 装备获得分数范围最大值 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"), zap.String("sheet", "maze_equip_affix_roll_type_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Get_score_max = int32(tmp)
	}

	// parse column 4 weight : 随机范围权重
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field weight 随机范围权重 to int32 failed")
			logger.ErrorWF("parse field weight 随机范围权重 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"), zap.String("sheet", "maze_equip_affix_roll_type_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Weight = int32(tmp)
	}

	// parse column 5 roll_range_min : 范围万分比，最小值
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field roll_range_min 范围万分比，最小值 to int32 failed")
			logger.ErrorWF("parse field roll_range_min 范围万分比，最小值 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"), zap.String("sheet", "maze_equip_affix_roll_type_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Roll_range_min = int32(tmp)
	}

	// parse column 6 roll_range_max : 范围万分比，最大值
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field roll_range_max 范围万分比，最大值 to int32 failed")
			logger.ErrorWF("parse field roll_range_max 范围万分比，最大值 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"), zap.String("sheet", "maze_equip_affix_roll_type_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Roll_range_max = int32(tmp)
	}

	// parse column 7 round_value : 取整数，（随机数/取整数）向上取整*取整数
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field round_value 取整数，（随机数/取整数）向上取整*取整数 to int32 failed")
			logger.ErrorWF("parse field round_value 取整数，（随机数/取整数）向上取整*取整数 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx"), zap.String("sheet", "maze_equip_affix_roll_type_v8"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Round_value = int32(tmp)
	}
	return
}

var gMazeEquipAffixRollTypeV8Fields = []string{
	"order",
	"roll_type",
	"get_score_min",
	"get_score_max",
	"weight",
	"roll_range_min",
	"roll_range_max",
	"round_value",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipAffixRollTypeV8Parser{}
	loader := &gMazeEquipAffixRollTypeV8Loader{}
	var data [][]string
	data, err = load("maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx", "maze_equip_affix_roll_type_v8", gMazeEquipAffixRollTypeV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx maze_equip_affix_roll_type_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx maze_equip_affix_roll_type_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipAffixRollTypeV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx maze_equip_affix_roll_type_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipAffixRollTypeV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx maze_equip_affix_roll_type_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx maze_equip_affix_roll_type_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx maze_equip_affix_roll_type_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_affix_roll_type_v8【迷宫-装备-词条值随机范围】.xlsx maze_equip_affix_roll_type_v8 data success.")
	return
}
