package GMazeSaveItemsV8Cfg

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

// MazeSaveItemsV8ConfigRow from maze_save_items_v8【迷宫-救援道具】.xlsx maze_save_items_v8
type MazeSaveItemsV8ConfigRow struct {
	Id        int32  `json:"id"`        // 道具id
	IconAtlas string `json:"iconAtlas"` // icon图集文件夹
	Icon      string `json:"icon"`      // 图标
	Type      int32  `json:"type"`      // 道具类型
	Argument  int32  `json:"argument"`  // 效果参数
}

// MazeSaveItemsV8Config from maze_save_items_v8【迷宫-救援道具】.xlsx maze_save_items_v8
type MazeSaveItemsV8Config struct {
	ConfigRows map[int32]*MazeSaveItemsV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeSaveItemsV8Config {
	ret := &MazeSaveItemsV8Config{ConfigRows: map[int32]*MazeSaveItemsV8ConfigRow{}}
	return ret
}

// GetMazeSaveItemsV8Config get one config by configId
func (c *MazeSaveItemsV8Config) GetMazeSaveItemsV8Config(configId int32) *MazeSaveItemsV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeSaveItemsV8Config) Get(configId int32) *MazeSaveItemsV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeSaveItemsV8Config get all config slice
func (c *MazeSaveItemsV8Config) GetAllMazeSaveItemsV8Config() (res []*MazeSaveItemsV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeSaveItemsV8Config) GetAll() (res []*MazeSaveItemsV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeSaveItemsV8Config

// GetMazeSaveItemsV8Config pkg func. get one config by configId
func GetMazeSaveItemsV8Config(configId int32) *MazeSaveItemsV8ConfigRow {
	return gConfigData.GetMazeSaveItemsV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeSaveItemsV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeSaveItemsV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_save_items_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeSaveItemsV8Config pkg func. get all config slice
func GetAllMazeSaveItemsV8Config() []*MazeSaveItemsV8ConfigRow {
	return gConfigData.GetAllMazeSaveItemsV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeSaveItemsV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeSaveItemsV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeSaveItemsV8ConfigRow from maze_save_items_v8【迷宫-救援道具】.xlsx maze_save_items_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeSaveItemsV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_save_items_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_save_items_v8.json",
		"maze_save_items_v8【迷宫-救援道具】.xlsx", "maze_save_items_v8",
		&gMazeSaveItemsV8Parser{}, &gMazeSaveItemsV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeSaveItemsV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeSaveItemsV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeSaveItemsV8Config))(c)
		return true
	})
}

// RegisterMazeSaveItemsV8InitCallBack reg config update func (old func)
var RegisterMazeSaveItemsV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeSaveItemsV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeSaveItemsV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeSaveItemsV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeSaveItemsV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeSaveItemsV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeSaveItemsV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeSaveItemsV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeSaveItemsV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeSaveItemsV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeSaveItemsV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeSaveItemsV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeSaveItemsV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeSaveItemsV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveItemsV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSaveItemsV8ConfigRow", zap.String("xlsx", "maze_save_items_v8【迷宫-救援道具】.xlsx"),
			zap.String("sheet", "maze_save_items_v8"))
		return
	}
	config, ok := container.(*MazeSaveItemsV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveItemsV8Config")
		logger.ErrorWF("invalid type. not *MazeSaveItemsV8Config", zap.String("xlsx", "maze_save_items_v8【迷宫-救援道具】.xlsx"),
			zap.String("sheet", "maze_save_items_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeSaveItemsV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeSaveItemsV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveItemsV8Config")
		logger.ErrorWF("invalid type. not *MazeSaveItemsV8Config", zap.String("xlsx", "maze_save_items_v8【迷宫-救援道具】.xlsx"),
			zap.String("sheet", "maze_save_items_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeSaveItemsV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeSaveItemsV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveItemsV8Config")
		logger.ErrorWF("invalid type. not *MazeSaveItemsV8Config", zap.String("xlsx", "maze_save_items_v8【迷宫-救援道具】.xlsx"),
			zap.String("sheet", "maze_save_items_v8"))
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
type gMazeSaveItemsV8Parser struct {
}

// New new config row data
func (*gMazeSaveItemsV8Parser) New() interface{} {
	return &MazeSaveItemsV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeSaveItemsV8Parser) Fields() []string {
	return gMazeSaveItemsV8Fields
}

// Parse parse raw data to row data
func (*gMazeSaveItemsV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeSaveItemsV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSaveItemsV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSaveItemsV8ConfigRow", zap.String("xlsx", "maze_save_items_v8【迷宫-救援道具】.xlsx"),
			zap.String("sheet", "maze_save_items_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeSaveItemsV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeSaveItemsV8ConfigRow",
			zap.String("xlsx", "maze_save_items_v8【迷宫-救援道具】.xlsx"),
			zap.String("sheet", "maze_save_items_v8"), zap.Int("need_count", len(gMazeSaveItemsV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 道具id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 道具id to int32 failed")
			logger.ErrorWF("parse field id 道具id to int32 failed.",
				zap.String("xlsx", "maze_save_items_v8【迷宫-救援道具】.xlsx"), zap.String("sheet", "maze_save_items_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 iconAtlas : icon图集文件夹
	if data[1] != "" {
		config.IconAtlas = data[1]
	}

	// parse column 2 icon : 图标
	if data[2] != "" {
		config.Icon = data[2]
	}

	// parse column 3 type : 道具类型
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field type 道具类型 to int32 failed")
			logger.ErrorWF("parse field type 道具类型 to int32 failed.",
				zap.String("xlsx", "maze_save_items_v8【迷宫-救援道具】.xlsx"), zap.String("sheet", "maze_save_items_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Type = int32(tmp)
	}

	// parse column 4 argument : 效果参数
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field argument 效果参数 to int32 failed")
			logger.ErrorWF("parse field argument 效果参数 to int32 failed.",
				zap.String("xlsx", "maze_save_items_v8【迷宫-救援道具】.xlsx"), zap.String("sheet", "maze_save_items_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Argument = int32(tmp)
	}
	return
}

var gMazeSaveItemsV8Fields = []string{
	"id",
	"iconAtlas",
	"icon",
	"type",
	"argument",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeSaveItemsV8Parser{}
	loader := &gMazeSaveItemsV8Loader{}
	var data [][]string
	data, err = load("maze_save_items_v8【迷宫-救援道具】.xlsx", "maze_save_items_v8", gMazeSaveItemsV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_save_items_v8【迷宫-救援道具】.xlsx maze_save_items_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_save_items_v8【迷宫-救援道具】.xlsx maze_save_items_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeSaveItemsV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_save_items_v8【迷宫-救援道具】.xlsx maze_save_items_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeSaveItemsV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_save_items_v8【迷宫-救援道具】.xlsx maze_save_items_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_save_items_v8【迷宫-救援道具】.xlsx maze_save_items_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_save_items_v8【迷宫-救援道具】.xlsx maze_save_items_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_save_items_v8【迷宫-救援道具】.xlsx maze_save_items_v8 data success.")
	return
}
