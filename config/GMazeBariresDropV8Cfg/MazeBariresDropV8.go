package GMazeBariresDropV8Cfg

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

// MazeBariresDropV8ConfigRow from maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx maze_barires_drop_v8
type MazeBariresDropV8ConfigRow struct {
	Order               int32           `json:"order"`               // 序号
	Drop_id             int32           `json:"drop_id"`             // 掉落id
	Drop_cd             int32           `json:"drop_cd"`             // 掉落cd（毫秒）
	In_barries_kill_min int32           `json:"in_barries_kill_min"` // 杀怪数小
	In_barries_kill_max int32           `json:"in_barries_kill_max"` // 杀怪数大
	Drop_ratio_min      int32           `json:"drop_ratio_min"`      // 掉落几率万分比
	Drop_ratio_max      int32           `json:"drop_ratio_max"`      // 掉落几率万分比
	Drop_items          map[int32]int64 `json:"drop_items"`          // 掉落物品：数量
	Drop_condition      []int32         `json:"drop_condition"`      // 掉落条件id
}

// MazeBariresDropV8Config from maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx maze_barires_drop_v8
type MazeBariresDropV8Config struct {
	ConfigRows map[int32]*MazeBariresDropV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeBariresDropV8Config {
	ret := &MazeBariresDropV8Config{ConfigRows: map[int32]*MazeBariresDropV8ConfigRow{}}
	return ret
}

// GetMazeBariresDropV8Config get one config by configId
func (c *MazeBariresDropV8Config) GetMazeBariresDropV8Config(configId int32) *MazeBariresDropV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeBariresDropV8Config) Get(configId int32) *MazeBariresDropV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeBariresDropV8Config get all config slice
func (c *MazeBariresDropV8Config) GetAllMazeBariresDropV8Config() (res []*MazeBariresDropV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeBariresDropV8Config) GetAll() (res []*MazeBariresDropV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeBariresDropV8Config

// GetMazeBariresDropV8Config pkg func. get one config by configId
func GetMazeBariresDropV8Config(configId int32) *MazeBariresDropV8ConfigRow {
	return gConfigData.GetMazeBariresDropV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeBariresDropV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeBariresDropV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_barires_drop_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeBariresDropV8Config pkg func. get all config slice
func GetAllMazeBariresDropV8Config() []*MazeBariresDropV8ConfigRow {
	return gConfigData.GetAllMazeBariresDropV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeBariresDropV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeBariresDropV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeBariresDropV8ConfigRow from maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx maze_barires_drop_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeBariresDropV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_barires_drop_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_barires_drop_v8.json",
		"maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx", "maze_barires_drop_v8",
		&gMazeBariresDropV8Parser{}, &gMazeBariresDropV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeBariresDropV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeBariresDropV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeBariresDropV8Config))(c)
		return true
	})
}

// RegisterMazeBariresDropV8InitCallBack reg config update func (old func)
var RegisterMazeBariresDropV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeBariresDropV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeBariresDropV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeBariresDropV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeBariresDropV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeBariresDropV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeBariresDropV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeBariresDropV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeBariresDropV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeBariresDropV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeBariresDropV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeBariresDropV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeBariresDropV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeBariresDropV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBariresDropV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBariresDropV8ConfigRow", zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"),
			zap.String("sheet", "maze_barires_drop_v8"))
		return
	}
	config, ok := container.(*MazeBariresDropV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBariresDropV8Config")
		logger.ErrorWF("invalid type. not *MazeBariresDropV8Config", zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"),
			zap.String("sheet", "maze_barires_drop_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeBariresDropV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeBariresDropV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBariresDropV8Config")
		logger.ErrorWF("invalid type. not *MazeBariresDropV8Config", zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"),
			zap.String("sheet", "maze_barires_drop_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeBariresDropV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeBariresDropV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBariresDropV8Config")
		logger.ErrorWF("invalid type. not *MazeBariresDropV8Config", zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"),
			zap.String("sheet", "maze_barires_drop_v8"))
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
type gMazeBariresDropV8Parser struct {
}

// New new config row data
func (*gMazeBariresDropV8Parser) New() interface{} {
	return &MazeBariresDropV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeBariresDropV8Parser) Fields() []string {
	return gMazeBariresDropV8Fields
}

// Parse parse raw data to row data
func (*gMazeBariresDropV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeBariresDropV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBariresDropV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBariresDropV8ConfigRow", zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"),
			zap.String("sheet", "maze_barires_drop_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeBariresDropV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeBariresDropV8ConfigRow",
			zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"),
			zap.String("sheet", "maze_barires_drop_v8"), zap.Int("need_count", len(gMazeBariresDropV8Fields)),
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
				zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"), zap.String("sheet", "maze_barires_drop_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 drop_id : 掉落id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_id 掉落id to int32 failed")
			logger.ErrorWF("parse field drop_id 掉落id to int32 failed.",
				zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"), zap.String("sheet", "maze_barires_drop_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Drop_id = int32(tmp)
	}

	// parse column 2 drop_cd : 掉落cd（毫秒）
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_cd 掉落cd（毫秒） to int32 failed")
			logger.ErrorWF("parse field drop_cd 掉落cd（毫秒） to int32 failed.",
				zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"), zap.String("sheet", "maze_barires_drop_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Drop_cd = int32(tmp)
	}

	// parse column 3 in_barries_kill_min : 杀怪数小
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field in_barries_kill_min 杀怪数小 to int32 failed")
			logger.ErrorWF("parse field in_barries_kill_min 杀怪数小 to int32 failed.",
				zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"), zap.String("sheet", "maze_barires_drop_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.In_barries_kill_min = int32(tmp)
	}

	// parse column 4 in_barries_kill_max : 杀怪数大
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field in_barries_kill_max 杀怪数大 to int32 failed")
			logger.ErrorWF("parse field in_barries_kill_max 杀怪数大 to int32 failed.",
				zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"), zap.String("sheet", "maze_barires_drop_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.In_barries_kill_max = int32(tmp)
	}

	// parse column 5 drop_ratio_min : 掉落几率万分比
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_ratio_min 掉落几率万分比 to int32 failed")
			logger.ErrorWF("parse field drop_ratio_min 掉落几率万分比 to int32 failed.",
				zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"), zap.String("sheet", "maze_barires_drop_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Drop_ratio_min = int32(tmp)
	}

	// parse column 6 drop_ratio_max : 掉落几率万分比
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_ratio_max 掉落几率万分比 to int32 failed")
			logger.ErrorWF("parse field drop_ratio_max 掉落几率万分比 to int32 failed.",
				zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"), zap.String("sheet", "maze_barires_drop_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Drop_ratio_max = int32(tmp)
	}

	// parse column 7 drop_items : 掉落物品：数量
	if data[7] != "" {

		config.Drop_items = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[7], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_items 掉落物品：数量 to key int32 failed")
				logger.ErrorWF("parse map field drop_items 掉落物品：数量 to key int32 failed.",
					zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"), zap.String("sheet", "maze_barires_drop_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_items 掉落物品：数量 to value int64 failed")
				logger.ErrorWF("parse map field drop_items 掉落物品：数量 to value int64 failed.",
					zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"), zap.String("sheet", "maze_barires_drop_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Drop_items[key] = value
		}
	}

	// parse column 8 drop_condition : 掉落条件id
	if data[8] != "" {

		vals := strings.Split(data[8], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field drop_condition 掉落条件id to []int32 failed")
				logger.ErrorWF("parse array field drop_condition 掉落条件id to []int32 failed.",
					zap.String("xlsx", "maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx"), zap.String("sheet", "maze_barires_drop_v8"),
					// zap.String("field_data",data[8]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Drop_condition = append(config.Drop_condition, int32(tmp))
		}
	}
	return
}

var gMazeBariresDropV8Fields = []string{
	"order",
	"drop_id",
	"drop_cd",
	"in_barries_kill_min",
	"in_barries_kill_max",
	"drop_ratio_min",
	"drop_ratio_max",
	"drop_items",
	"drop_condition",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeBariresDropV8Parser{}
	loader := &gMazeBariresDropV8Loader{}
	var data [][]string
	data, err = load("maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx", "maze_barires_drop_v8", gMazeBariresDropV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx maze_barires_drop_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx maze_barires_drop_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeBariresDropV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx maze_barires_drop_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeBariresDropV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx maze_barires_drop_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx maze_barires_drop_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx maze_barires_drop_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_barires_drop_v8【迷宫-关卡-掉落物品】.xlsx maze_barires_drop_v8 data success.")
	return
}
