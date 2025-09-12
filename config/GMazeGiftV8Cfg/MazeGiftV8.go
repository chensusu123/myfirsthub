package GMazeGiftV8Cfg

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

// MazeGiftV8ConfigRow from maze_exchange_v8【迷宫-商城】.xlsx maze_gift_v8
type MazeGiftV8ConfigRow struct {
	Id        int32           `json:"id"`        // 礼包id
	Item_num  map[int32]int64 `json:"item_num"`  // 礼包道具列表物品：数量
	Item_sort []int32         `json:"item_sort"` // 道具排序
}

// MazeGiftV8Config from maze_exchange_v8【迷宫-商城】.xlsx maze_gift_v8
type MazeGiftV8Config struct {
	ConfigRows map[int32]*MazeGiftV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeGiftV8Config {
	ret := &MazeGiftV8Config{ConfigRows: map[int32]*MazeGiftV8ConfigRow{}}
	return ret
}

// GetMazeGiftV8Config get one config by configId
func (c *MazeGiftV8Config) GetMazeGiftV8Config(configId int32) *MazeGiftV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeGiftV8Config) Get(configId int32) *MazeGiftV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeGiftV8Config get all config slice
func (c *MazeGiftV8Config) GetAllMazeGiftV8Config() (res []*MazeGiftV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeGiftV8Config) GetAll() (res []*MazeGiftV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeGiftV8Config

// GetMazeGiftV8Config pkg func. get one config by configId
func GetMazeGiftV8Config(configId int32) *MazeGiftV8ConfigRow {
	return gConfigData.GetMazeGiftV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeGiftV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeGiftV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_gift_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeGiftV8Config pkg func. get all config slice
func GetAllMazeGiftV8Config() []*MazeGiftV8ConfigRow {
	return gConfigData.GetAllMazeGiftV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeGiftV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeGiftV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeGiftV8ConfigRow from maze_exchange_v8【迷宫-商城】.xlsx maze_gift_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeGiftV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_gift_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_gift_v8.json",
		"maze_exchange_v8【迷宫-商城】.xlsx", "maze_gift_v8",
		&gMazeGiftV8Parser{}, &gMazeGiftV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeGiftV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeGiftV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeGiftV8Config))(c)
		return true
	})
}

// RegisterMazeGiftV8InitCallBack reg config update func (old func)
var RegisterMazeGiftV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeGiftV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeGiftV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeGiftV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeGiftV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeGiftV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeGiftV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeGiftV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeGiftV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeGiftV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeGiftV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeGiftV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeGiftV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeGiftV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeGiftV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeGiftV8ConfigRow", zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"),
			zap.String("sheet", "maze_gift_v8"))
		return
	}
	config, ok := container.(*MazeGiftV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeGiftV8Config")
		logger.ErrorWF("invalid type. not *MazeGiftV8Config", zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"),
			zap.String("sheet", "maze_gift_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeGiftV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeGiftV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeGiftV8Config")
		logger.ErrorWF("invalid type. not *MazeGiftV8Config", zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"),
			zap.String("sheet", "maze_gift_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeGiftV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeGiftV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeGiftV8Config")
		logger.ErrorWF("invalid type. not *MazeGiftV8Config", zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"),
			zap.String("sheet", "maze_gift_v8"))
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
type gMazeGiftV8Parser struct {
}

// New new config row data
func (*gMazeGiftV8Parser) New() interface{} {
	return &MazeGiftV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeGiftV8Parser) Fields() []string {
	return gMazeGiftV8Fields
}

// Parse parse raw data to row data
func (*gMazeGiftV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeGiftV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeGiftV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeGiftV8ConfigRow", zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"),
			zap.String("sheet", "maze_gift_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeGiftV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeGiftV8ConfigRow",
			zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"),
			zap.String("sheet", "maze_gift_v8"), zap.Int("need_count", len(gMazeGiftV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 礼包id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 礼包id to int32 failed")
			logger.ErrorWF("parse field id 礼包id to int32 failed.",
				zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_gift_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 item_num : 礼包道具列表物品：数量
	if data[1] != "" {

		config.Item_num = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[1], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field item_num 礼包道具列表物品：数量 to key int32 failed")
				logger.ErrorWF("parse map field item_num 礼包道具列表物品：数量 to key int32 failed.",
					zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_gift_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field item_num 礼包道具列表物品：数量 to value int64 failed")
				logger.ErrorWF("parse map field item_num 礼包道具列表物品：数量 to value int64 failed.",
					zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_gift_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Item_num[key] = value
		}
	}

	// parse column 2 item_sort : 道具排序
	if data[2] != "" {

		vals := strings.Split(data[2], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field item_sort 道具排序 to []int32 failed")
				logger.ErrorWF("parse array field item_sort 道具排序 to []int32 failed.",
					zap.String("xlsx", "maze_exchange_v8【迷宫-商城】.xlsx"), zap.String("sheet", "maze_gift_v8"),
					// zap.String("field_data",data[2]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Item_sort = append(config.Item_sort, int32(tmp))
		}
	}
	return
}

var gMazeGiftV8Fields = []string{
	"id",
	"item_num",
	"item_sort",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeGiftV8Parser{}
	loader := &gMazeGiftV8Loader{}
	var data [][]string
	data, err = load("maze_exchange_v8【迷宫-商城】.xlsx", "maze_gift_v8", gMazeGiftV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_exchange_v8【迷宫-商城】.xlsx maze_gift_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_exchange_v8【迷宫-商城】.xlsx maze_gift_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeGiftV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_exchange_v8【迷宫-商城】.xlsx maze_gift_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeGiftV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_exchange_v8【迷宫-商城】.xlsx maze_gift_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_exchange_v8【迷宫-商城】.xlsx maze_gift_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_exchange_v8【迷宫-商城】.xlsx maze_gift_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_exchange_v8【迷宫-商城】.xlsx maze_gift_v8 data success.")
	return
}
