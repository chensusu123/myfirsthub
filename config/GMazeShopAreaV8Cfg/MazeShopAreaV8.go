package GMazeShopAreaV8Cfg

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

// MazeShopAreaV8ConfigRow from maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx maze_shop_area_v8
type MazeShopAreaV8ConfigRow struct {
	Order          int32 `json:"order"`          // 区域id
	Shop_level_min int32 `json:"shop_level_min"` // 对应最小商店等级
	Shop_level_max int32 `json:"shop_level_max"` // 对应最大商店等级
}

// MazeShopAreaV8Config from maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx maze_shop_area_v8
type MazeShopAreaV8Config struct {
	ConfigRows map[int32]*MazeShopAreaV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeShopAreaV8Config {
	ret := &MazeShopAreaV8Config{ConfigRows: map[int32]*MazeShopAreaV8ConfigRow{}}
	return ret
}

// GetMazeShopAreaV8Config get one config by configId
func (c *MazeShopAreaV8Config) GetMazeShopAreaV8Config(configId int32) *MazeShopAreaV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeShopAreaV8Config) Get(configId int32) *MazeShopAreaV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeShopAreaV8Config get all config slice
func (c *MazeShopAreaV8Config) GetAllMazeShopAreaV8Config() (res []*MazeShopAreaV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeShopAreaV8Config) GetAll() (res []*MazeShopAreaV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeShopAreaV8Config

// GetMazeShopAreaV8Config pkg func. get one config by configId
func GetMazeShopAreaV8Config(configId int32) *MazeShopAreaV8ConfigRow {
	return gConfigData.GetMazeShopAreaV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeShopAreaV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeShopAreaV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_shop_area_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeShopAreaV8Config pkg func. get all config slice
func GetAllMazeShopAreaV8Config() []*MazeShopAreaV8ConfigRow {
	return gConfigData.GetAllMazeShopAreaV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeShopAreaV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeShopAreaV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeShopAreaV8ConfigRow from maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx maze_shop_area_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeShopAreaV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_shop_area_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_shop_area_v8.json",
		"maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx", "maze_shop_area_v8",
		&gMazeShopAreaV8Parser{}, &gMazeShopAreaV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeShopAreaV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeShopAreaV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeShopAreaV8Config))(c)
		return true
	})
}

// RegisterMazeShopAreaV8InitCallBack reg config update func (old func)
var RegisterMazeShopAreaV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeShopAreaV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeShopAreaV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeShopAreaV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeShopAreaV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeShopAreaV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeShopAreaV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeShopAreaV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeShopAreaV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeShopAreaV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeShopAreaV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeShopAreaV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeShopAreaV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeShopAreaV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeShopAreaV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeShopAreaV8ConfigRow", zap.String("xlsx", "maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx"),
			zap.String("sheet", "maze_shop_area_v8"))
		return
	}
	config, ok := container.(*MazeShopAreaV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeShopAreaV8Config")
		logger.ErrorWF("invalid type. not *MazeShopAreaV8Config", zap.String("xlsx", "maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx"),
			zap.String("sheet", "maze_shop_area_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeShopAreaV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeShopAreaV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeShopAreaV8Config")
		logger.ErrorWF("invalid type. not *MazeShopAreaV8Config", zap.String("xlsx", "maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx"),
			zap.String("sheet", "maze_shop_area_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeShopAreaV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeShopAreaV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeShopAreaV8Config")
		logger.ErrorWF("invalid type. not *MazeShopAreaV8Config", zap.String("xlsx", "maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx"),
			zap.String("sheet", "maze_shop_area_v8"))
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
type gMazeShopAreaV8Parser struct {
}

// New new config row data
func (*gMazeShopAreaV8Parser) New() interface{} {
	return &MazeShopAreaV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeShopAreaV8Parser) Fields() []string {
	return gMazeShopAreaV8Fields
}

// Parse parse raw data to row data
func (*gMazeShopAreaV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeShopAreaV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeShopAreaV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeShopAreaV8ConfigRow", zap.String("xlsx", "maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx"),
			zap.String("sheet", "maze_shop_area_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeShopAreaV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeShopAreaV8ConfigRow",
			zap.String("xlsx", "maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx"),
			zap.String("sheet", "maze_shop_area_v8"), zap.Int("need_count", len(gMazeShopAreaV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 区域id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 区域id to int32 failed")
			logger.ErrorWF("parse field order 区域id to int32 failed.",
				zap.String("xlsx", "maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx"), zap.String("sheet", "maze_shop_area_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 shop_level_min : 对应最小商店等级
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field shop_level_min 对应最小商店等级 to int32 failed")
			logger.ErrorWF("parse field shop_level_min 对应最小商店等级 to int32 failed.",
				zap.String("xlsx", "maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx"), zap.String("sheet", "maze_shop_area_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Shop_level_min = int32(tmp)
	}

	// parse column 2 shop_level_max : 对应最大商店等级
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field shop_level_max 对应最大商店等级 to int32 failed")
			logger.ErrorWF("parse field shop_level_max 对应最大商店等级 to int32 failed.",
				zap.String("xlsx", "maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx"), zap.String("sheet", "maze_shop_area_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Shop_level_max = int32(tmp)
	}
	return
}

var gMazeShopAreaV8Fields = []string{
	"order",
	"shop_level_min",
	"shop_level_max",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeShopAreaV8Parser{}
	loader := &gMazeShopAreaV8Loader{}
	var data [][]string
	data, err = load("maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx", "maze_shop_area_v8", gMazeShopAreaV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx maze_shop_area_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx maze_shop_area_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeShopAreaV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx maze_shop_area_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeShopAreaV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx maze_shop_area_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx maze_shop_area_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx maze_shop_area_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_shop_area_v8【迷宫-商店-刷怪区域id对应商店】.xlsx maze_shop_area_v8 data success.")
	return
}
