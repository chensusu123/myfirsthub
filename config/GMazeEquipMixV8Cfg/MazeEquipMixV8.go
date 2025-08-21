package GMazeEquipMixV8Cfg

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

// MazeEquipMixV8ConfigRow from maze_equip_mix_v8【迷宫-装备-合成】.xlsx maze_equip_mix_v8
type MazeEquipMixV8ConfigRow struct {
	Order             int32           `json:"order"`             // 冒险等级
	Cost              map[int32]int64 `json:"cost"`              // 合成消耗
	Result_equip_list []int32         `json:"result_equip_list"` // 合成结果队列id
}

// MazeEquipMixV8Config from maze_equip_mix_v8【迷宫-装备-合成】.xlsx maze_equip_mix_v8
type MazeEquipMixV8Config struct {
	ConfigRows map[int32]*MazeEquipMixV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipMixV8Config {
	ret := &MazeEquipMixV8Config{ConfigRows: map[int32]*MazeEquipMixV8ConfigRow{}}
	return ret
}

// GetMazeEquipMixV8Config get one config by configId
func (c *MazeEquipMixV8Config) GetMazeEquipMixV8Config(configId int32) *MazeEquipMixV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipMixV8Config) Get(configId int32) *MazeEquipMixV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipMixV8Config get all config slice
func (c *MazeEquipMixV8Config) GetAllMazeEquipMixV8Config() (res []*MazeEquipMixV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipMixV8Config) GetAll() (res []*MazeEquipMixV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipMixV8Config

// GetMazeEquipMixV8Config pkg func. get one config by configId
func GetMazeEquipMixV8Config(configId int32) *MazeEquipMixV8ConfigRow {
	return gConfigData.GetMazeEquipMixV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipMixV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeEquipMixV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_equip_mix_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeEquipMixV8Config pkg func. get all config slice
func GetAllMazeEquipMixV8Config() []*MazeEquipMixV8ConfigRow {
	return gConfigData.GetAllMazeEquipMixV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipMixV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipMixV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipMixV8ConfigRow from maze_equip_mix_v8【迷宫-装备-合成】.xlsx maze_equip_mix_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipMixV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_mix_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_mix_v8.json",
		"maze_equip_mix_v8【迷宫-装备-合成】.xlsx", "maze_equip_mix_v8",
		&gMazeEquipMixV8Parser{}, &gMazeEquipMixV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipMixV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipMixV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipMixV8Config))(c)
		return true
	})
}

// RegisterMazeEquipMixV8InitCallBack reg config update func (old func)
var RegisterMazeEquipMixV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipMixV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipMixV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipMixV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipMixV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipMixV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipMixV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipMixV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipMixV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipMixV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipMixV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipMixV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipMixV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipMixV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipMixV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipMixV8ConfigRow", zap.String("xlsx", "maze_equip_mix_v8【迷宫-装备-合成】.xlsx"),
			zap.String("sheet", "maze_equip_mix_v8"))
		return
	}
	config, ok := container.(*MazeEquipMixV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipMixV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipMixV8Config", zap.String("xlsx", "maze_equip_mix_v8【迷宫-装备-合成】.xlsx"),
			zap.String("sheet", "maze_equip_mix_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipMixV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipMixV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipMixV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipMixV8Config", zap.String("xlsx", "maze_equip_mix_v8【迷宫-装备-合成】.xlsx"),
			zap.String("sheet", "maze_equip_mix_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipMixV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipMixV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipMixV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipMixV8Config", zap.String("xlsx", "maze_equip_mix_v8【迷宫-装备-合成】.xlsx"),
			zap.String("sheet", "maze_equip_mix_v8"))
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
type gMazeEquipMixV8Parser struct {
}

// New new config row data
func (*gMazeEquipMixV8Parser) New() interface{} {
	return &MazeEquipMixV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipMixV8Parser) Fields() []string {
	return gMazeEquipMixV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipMixV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipMixV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipMixV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipMixV8ConfigRow", zap.String("xlsx", "maze_equip_mix_v8【迷宫-装备-合成】.xlsx"),
			zap.String("sheet", "maze_equip_mix_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipMixV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipMixV8ConfigRow",
			zap.String("xlsx", "maze_equip_mix_v8【迷宫-装备-合成】.xlsx"),
			zap.String("sheet", "maze_equip_mix_v8"), zap.Int("need_count", len(gMazeEquipMixV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 冒险等级
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 冒险等级 to int32 failed")
			logger.ErrorWF("parse field order 冒险等级 to int32 failed.",
				zap.String("xlsx", "maze_equip_mix_v8【迷宫-装备-合成】.xlsx"), zap.String("sheet", "maze_equip_mix_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 cost : 合成消耗
	if data[1] != "" {

		config.Cost = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[1], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field cost 合成消耗 to key int32 failed")
				logger.ErrorWF("parse map field cost 合成消耗 to key int32 failed.",
					zap.String("xlsx", "maze_equip_mix_v8【迷宫-装备-合成】.xlsx"), zap.String("sheet", "maze_equip_mix_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field cost 合成消耗 to value int64 failed")
				logger.ErrorWF("parse map field cost 合成消耗 to value int64 failed.",
					zap.String("xlsx", "maze_equip_mix_v8【迷宫-装备-合成】.xlsx"), zap.String("sheet", "maze_equip_mix_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Cost[key] = value
		}
	}

	// parse column 2 result_equip_list : 合成结果队列id
	if data[2] != "" {

		vals := strings.Split(data[2], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field result_equip_list 合成结果队列id to []int32 failed")
				logger.ErrorWF("parse array field result_equip_list 合成结果队列id to []int32 failed.",
					zap.String("xlsx", "maze_equip_mix_v8【迷宫-装备-合成】.xlsx"), zap.String("sheet", "maze_equip_mix_v8"),
					// zap.String("field_data",data[2]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Result_equip_list = append(config.Result_equip_list, int32(tmp))
		}
	}
	return
}

var gMazeEquipMixV8Fields = []string{
	"order",
	"cost",
	"result_equip_list",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipMixV8Parser{}
	loader := &gMazeEquipMixV8Loader{}
	var data [][]string
	data, err = load("maze_equip_mix_v8【迷宫-装备-合成】.xlsx", "maze_equip_mix_v8", gMazeEquipMixV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_mix_v8【迷宫-装备-合成】.xlsx maze_equip_mix_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_mix_v8【迷宫-装备-合成】.xlsx maze_equip_mix_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipMixV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_mix_v8【迷宫-装备-合成】.xlsx maze_equip_mix_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipMixV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_mix_v8【迷宫-装备-合成】.xlsx maze_equip_mix_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_mix_v8【迷宫-装备-合成】.xlsx maze_equip_mix_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_mix_v8【迷宫-装备-合成】.xlsx maze_equip_mix_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_mix_v8【迷宫-装备-合成】.xlsx maze_equip_mix_v8 data success.")
	return
}
