package GMazeShopV8Cfg

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

// MazeShopV8ConfigRow from maze_shop_v8【迷宫-商店】.xlsx maze_shop_v8
type MazeShopV8ConfigRow struct {
	Order            int32           `json:"order"`            // 冒险等级
	Buy_cost         int32           `json:"buy_cost"`         // 购买需要钱币
	Buy_list         []int32         `json:"buy_list"`         // 购买到装备的队列id
	Buy_list_base2   int32           `json:"buy_list_base2"`   // 队列走完后，使用的备用队列
	Equip_value_max  map[int32]int64 `json:"equip_value_max"`  // 部位：主属性最大值
	Need_equip_score int32           `json:"need_equip_score"` // 掉落装备需要的积分
}

// MazeShopV8Config from maze_shop_v8【迷宫-商店】.xlsx maze_shop_v8
type MazeShopV8Config struct {
	ConfigRows map[int32]*MazeShopV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeShopV8Config {
	ret := &MazeShopV8Config{ConfigRows: map[int32]*MazeShopV8ConfigRow{}}
	return ret
}

// GetMazeShopV8Config get one config by configId
func (c *MazeShopV8Config) GetMazeShopV8Config(configId int32) *MazeShopV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeShopV8Config) Get(configId int32) *MazeShopV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeShopV8Config get all config slice
func (c *MazeShopV8Config) GetAllMazeShopV8Config() (res []*MazeShopV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeShopV8Config) GetAll() (res []*MazeShopV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeShopV8Config

// GetMazeShopV8Config pkg func. get one config by configId
func GetMazeShopV8Config(configId int32) *MazeShopV8ConfigRow {
	return gConfigData.GetMazeShopV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeShopV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeShopV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_shop_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeShopV8Config pkg func. get all config slice
func GetAllMazeShopV8Config() []*MazeShopV8ConfigRow {
	return gConfigData.GetAllMazeShopV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeShopV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeShopV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeShopV8ConfigRow from maze_shop_v8【迷宫-商店】.xlsx maze_shop_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeShopV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_shop_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_shop_v8.json",
		"maze_shop_v8【迷宫-商店】.xlsx", "maze_shop_v8",
		&gMazeShopV8Parser{}, &gMazeShopV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeShopV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeShopV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeShopV8Config))(c)
		return true
	})
}

// RegisterMazeShopV8InitCallBack reg config update func (old func)
var RegisterMazeShopV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeShopV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeShopV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeShopV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeShopV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeShopV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeShopV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeShopV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeShopV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeShopV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeShopV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeShopV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeShopV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeShopV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeShopV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeShopV8ConfigRow", zap.String("xlsx", "maze_shop_v8【迷宫-商店】.xlsx"),
			zap.String("sheet", "maze_shop_v8"))
		return
	}
	config, ok := container.(*MazeShopV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeShopV8Config")
		logger.ErrorWF("invalid type. not *MazeShopV8Config", zap.String("xlsx", "maze_shop_v8【迷宫-商店】.xlsx"),
			zap.String("sheet", "maze_shop_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeShopV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeShopV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeShopV8Config")
		logger.ErrorWF("invalid type. not *MazeShopV8Config", zap.String("xlsx", "maze_shop_v8【迷宫-商店】.xlsx"),
			zap.String("sheet", "maze_shop_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeShopV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeShopV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeShopV8Config")
		logger.ErrorWF("invalid type. not *MazeShopV8Config", zap.String("xlsx", "maze_shop_v8【迷宫-商店】.xlsx"),
			zap.String("sheet", "maze_shop_v8"))
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
type gMazeShopV8Parser struct {
}

// New new config row data
func (*gMazeShopV8Parser) New() interface{} {
	return &MazeShopV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeShopV8Parser) Fields() []string {
	return gMazeShopV8Fields
}

// Parse parse raw data to row data
func (*gMazeShopV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeShopV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeShopV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeShopV8ConfigRow", zap.String("xlsx", "maze_shop_v8【迷宫-商店】.xlsx"),
			zap.String("sheet", "maze_shop_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeShopV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeShopV8ConfigRow",
			zap.String("xlsx", "maze_shop_v8【迷宫-商店】.xlsx"),
			zap.String("sheet", "maze_shop_v8"), zap.Int("need_count", len(gMazeShopV8Fields)),
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
				zap.String("xlsx", "maze_shop_v8【迷宫-商店】.xlsx"), zap.String("sheet", "maze_shop_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 buy_cost : 购买需要钱币
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field buy_cost 购买需要钱币 to int32 failed")
			logger.ErrorWF("parse field buy_cost 购买需要钱币 to int32 failed.",
				zap.String("xlsx", "maze_shop_v8【迷宫-商店】.xlsx"), zap.String("sheet", "maze_shop_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Buy_cost = int32(tmp)
	}

	// parse column 2 buy_list : 购买到装备的队列id
	if data[2] != "" {

		vals := strings.Split(data[2], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field buy_list 购买到装备的队列id to []int32 failed")
				logger.ErrorWF("parse array field buy_list 购买到装备的队列id to []int32 failed.",
					zap.String("xlsx", "maze_shop_v8【迷宫-商店】.xlsx"), zap.String("sheet", "maze_shop_v8"),
					// zap.String("field_data",data[2]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Buy_list = append(config.Buy_list, int32(tmp))
		}
	}

	// parse column 3 buy_list_base2 : 队列走完后，使用的备用队列
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field buy_list_base2 队列走完后，使用的备用队列 to int32 failed")
			logger.ErrorWF("parse field buy_list_base2 队列走完后，使用的备用队列 to int32 failed.",
				zap.String("xlsx", "maze_shop_v8【迷宫-商店】.xlsx"), zap.String("sheet", "maze_shop_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Buy_list_base2 = int32(tmp)
	}

	// parse column 4 equip_value_max : 部位：主属性最大值
	if data[4] != "" {

		config.Equip_value_max = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[4], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field equip_value_max 部位：主属性最大值 to key int32 failed")
				logger.ErrorWF("parse map field equip_value_max 部位：主属性最大值 to key int32 failed.",
					zap.String("xlsx", "maze_shop_v8【迷宫-商店】.xlsx"), zap.String("sheet", "maze_shop_v8"),
					// zap.String("field_data",data[4]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field equip_value_max 部位：主属性最大值 to value int64 failed")
				logger.ErrorWF("parse map field equip_value_max 部位：主属性最大值 to value int64 failed.",
					zap.String("xlsx", "maze_shop_v8【迷宫-商店】.xlsx"), zap.String("sheet", "maze_shop_v8"),
					// zap.String("field_data",data[4]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Equip_value_max[key] = value
		}
	}

	// parse column 5 need_equip_score : 掉落装备需要的积分
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field need_equip_score 掉落装备需要的积分 to int32 failed")
			logger.ErrorWF("parse field need_equip_score 掉落装备需要的积分 to int32 failed.",
				zap.String("xlsx", "maze_shop_v8【迷宫-商店】.xlsx"), zap.String("sheet", "maze_shop_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Need_equip_score = int32(tmp)
	}
	return
}

var gMazeShopV8Fields = []string{
	"order",
	"buy_cost",
	"buy_list",
	"buy_list_base2",
	"equip_value_max",
	"need_equip_score",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeShopV8Parser{}
	loader := &gMazeShopV8Loader{}
	var data [][]string
	data, err = load("maze_shop_v8【迷宫-商店】.xlsx", "maze_shop_v8", gMazeShopV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_shop_v8【迷宫-商店】.xlsx maze_shop_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_shop_v8【迷宫-商店】.xlsx maze_shop_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeShopV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_shop_v8【迷宫-商店】.xlsx maze_shop_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeShopV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_shop_v8【迷宫-商店】.xlsx maze_shop_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_shop_v8【迷宫-商店】.xlsx maze_shop_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_shop_v8【迷宫-商店】.xlsx maze_shop_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_shop_v8【迷宫-商店】.xlsx maze_shop_v8 data success.")
	return
}
