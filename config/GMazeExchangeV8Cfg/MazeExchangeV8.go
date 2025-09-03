package GMazeExchangeV8Cfg

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

// MazeExchangeV8ConfigRow from maze_exchange_v8【迷宫-商城】.xlsx maze_exchange_v8
type MazeExchangeV8ConfigRow struct {
	Id         int32           `json:"id"`         // 商品id
	Shop_id    int32           `json:"shop_id"`    // 商城id
	Sort       int32           `json:"sort"`       // 排序
	Item_num   map[int32]int64 `json:"item_num"`   // 单次出售物品：数量
	Buy_cost   map[int32]int64 `json:"buy_cost"`   // 购买需要消耗物品：数量
	Limit_type int32           `json:"limit_type"` // 限购类型（0-不限购，1~4日周月终身)
	Limit_num  int32           `json:"limit_num"`  // 限购数量
	Discount   int32           `json:"discount"`   // 显示折扣（万分比）
}

// MazeExchangeV8Config from maze_exchange_v8【迷宫-商城】.xlsx maze_exchange_v8
type MazeExchangeV8Config struct {
	ConfigRows map[int32]*MazeExchangeV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeExchangeV8Config {
	ret := &MazeExchangeV8Config{ConfigRows: map[int32]*MazeExchangeV8ConfigRow{}}
	return ret
}

// GetMazeExchangeV8Config get one config by configId
func (c *MazeExchangeV8Config) GetMazeExchangeV8Config(configId int32) *MazeExchangeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeExchangeV8Config) Get(configId int32) *MazeExchangeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeExchangeV8Config get all config slice
func (c *MazeExchangeV8Config) GetAllMazeExchangeV8Config() (res []*MazeExchangeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeExchangeV8Config) GetAll() (res []*MazeExchangeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeExchangeV8Config

// GetMazeExchangeV8Config pkg func. get one config by configId
func GetMazeExchangeV8Config(configId int32) *MazeExchangeV8ConfigRow {
	return gConfigData.GetMazeExchangeV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeExchangeV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeExchangeV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_exchange_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeExchangeV8Config pkg func. get all config slice
func GetAllMazeExchangeV8Config() []*MazeExchangeV8ConfigRow {
	return gConfigData.GetAllMazeExchangeV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeExchangeV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeExchangeV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeExchangeV8ConfigRow from maze_exchange_v8【迷宫-商城】.xlsx maze_exchange_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeExchangeV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_exchange_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_exchange_v8.json",
		"maze_exchange_v8【迷宫-商城】.xlsx", "maze_exchange_v8",
		&gMazeExchangeV8Parser{}, &gMazeExchangeV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeExchangeV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeExchangeV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeExchangeV8Config))(c)
		return true
	})
}

// RegisterMazeExchangeV8InitCallBack reg config update func (old func)
var RegisterMazeExchangeV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeExchangeV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeExchangeV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeExchangeV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeExchangeV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeExchangeV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeExchangeV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeExchangeV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeExchangeV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeExchangeV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeExchangeV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeExchangeV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeExchangeV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeExchangeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeExchangeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeExchangeV8ConfigRow", zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"),
			zap.String("sheet", "maze_exchange_v8"))
		return
	}
	config, ok := container.(*MazeExchangeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeExchangeV8Config")
		logger.ErrorWF("invalid type. not *MazeExchangeV8Config", zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"),
			zap.String("sheet", "maze_exchange_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeExchangeV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeExchangeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeExchangeV8Config")
		logger.ErrorWF("invalid type. not *MazeExchangeV8Config", zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"),
			zap.String("sheet", "maze_exchange_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeExchangeV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeExchangeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeExchangeV8Config")
		logger.ErrorWF("invalid type. not *MazeExchangeV8Config", zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"),
			zap.String("sheet", "maze_exchange_v8"))
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
type gMazeExchangeV8Parser struct {
}

// New new config row data
func (*gMazeExchangeV8Parser) New() interface{} {
	return &MazeExchangeV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeExchangeV8Parser) Fields() []string {
	return gMazeExchangeV8Fields
}

// Parse parse raw data to row data
func (*gMazeExchangeV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeExchangeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeExchangeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeExchangeV8ConfigRow", zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"),
			zap.String("sheet", "maze_exchange_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeExchangeV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeExchangeV8ConfigRow",
			zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"),
			zap.String("sheet", "maze_exchange_v8"), zap.Int("need_count", len(gMazeExchangeV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 商品id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 商品id to int32 failed")
			logger.ErrorWF("parse field id 商品id to int32 failed.",
				zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_exchange_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 shop_id : 商城id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field shop_id 商城id to int32 failed")
			logger.ErrorWF("parse field shop_id 商城id to int32 failed.",
				zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_exchange_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Shop_id = int32(tmp)
	}

	// parse column 2 sort : 排序
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field sort 排序 to int32 failed")
			logger.ErrorWF("parse field sort 排序 to int32 failed.",
				zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_exchange_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Sort = int32(tmp)
	}

	// parse column 3 item_num : 单次出售物品：数量
	if data[3] != "" {

		config.Item_num = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field item_num 单次出售物品：数量 to key int32 failed")
				logger.ErrorWF("parse map field item_num 单次出售物品：数量 to key int32 failed.",
					zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_exchange_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field item_num 单次出售物品：数量 to value int64 failed")
				logger.ErrorWF("parse map field item_num 单次出售物品：数量 to value int64 failed.",
					zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_exchange_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Item_num[key] = value
		}
	}

	// parse column 4 buy_cost : 购买需要消耗物品：数量
	if data[4] != "" {

		config.Buy_cost = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[4], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field buy_cost 购买需要消耗物品：数量 to key int32 failed")
				logger.ErrorWF("parse map field buy_cost 购买需要消耗物品：数量 to key int32 failed.",
					zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_exchange_v8"),
					// zap.String("field_data",data[4]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field buy_cost 购买需要消耗物品：数量 to value int64 failed")
				logger.ErrorWF("parse map field buy_cost 购买需要消耗物品：数量 to value int64 failed.",
					zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_exchange_v8"),
					// zap.String("field_data",data[4]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Buy_cost[key] = value
		}
	}

	// parse column 5 limit_type : 限购类型（0-不限购，1~4日周月终身)
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field limit_type 限购类型（0-不限购，1~4日周月终身) to int32 failed")
			logger.ErrorWF("parse field limit_type 限购类型（0-不限购，1~4日周月终身) to int32 failed.",
				zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_exchange_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Limit_type = int32(tmp)
	}

	// parse column 6 limit_num : 限购数量
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field limit_num 限购数量 to int32 failed")
			logger.ErrorWF("parse field limit_num 限购数量 to int32 failed.",
				zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_exchange_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Limit_num = int32(tmp)
	}

	// parse column 7 discount : 显示折扣（万分比）
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field discount 显示折扣（万分比） to int32 failed")
			logger.ErrorWF("parse field discount 显示折扣（万分比） to int32 failed.",
				zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_exchange_v8"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Discount = int32(tmp)
	}
	return
}

var gMazeExchangeV8Fields = []string{
	"id",
	"shop_id",
	"sort",
	"item_num",
	"buy_cost",
	"limit_type",
	"limit_num",
	"discount",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeExchangeV8Parser{}
	loader := &gMazeExchangeV8Loader{}
	var data [][]string
	data, err = load("maze_exchange_v8【迷宫-商城】.xlsx", "maze_exchange_v8", gMazeExchangeV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_exchange_v8【迷宫-商城】.xlsx maze_exchange_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_exchange_v8【迷宫-商城】.xlsx maze_exchange_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeExchangeV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_exchange_v8【迷宫-商城】.xlsx maze_exchange_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeExchangeV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_exchange_v8【迷宫-商城】.xlsx maze_exchange_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_exchange_v8【迷宫-商城】.xlsx maze_exchange_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_exchange_v8【迷宫-商城】.xlsx maze_exchange_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_exchange_v8【迷宫-商城】.xlsx maze_exchange_v8 data success.")
	return
}
