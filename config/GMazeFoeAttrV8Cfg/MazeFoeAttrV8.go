package GMazeFoeAttrV8Cfg

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

// MazeFoeAttrV8ConfigRow from maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_attr_v8
type MazeFoeAttrV8ConfigRow struct {
	Order      int32           `json:"order"`      // 属性id
	Extra_attr map[int32]int32 `json:"extra_attr"` // 额外属性列表
}

// MazeFoeAttrV8Config from maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_attr_v8
type MazeFoeAttrV8Config struct {
	ConfigRows map[int32]*MazeFoeAttrV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeFoeAttrV8Config {
	ret := &MazeFoeAttrV8Config{ConfigRows: map[int32]*MazeFoeAttrV8ConfigRow{}}
	return ret
}

// GetMazeFoeAttrV8Config get one config by configId
func (c *MazeFoeAttrV8Config) GetMazeFoeAttrV8Config(configId int32) *MazeFoeAttrV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeFoeAttrV8Config) Get(configId int32) *MazeFoeAttrV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeFoeAttrV8Config get all config slice
func (c *MazeFoeAttrV8Config) GetAllMazeFoeAttrV8Config() (res []*MazeFoeAttrV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeFoeAttrV8Config) GetAll() (res []*MazeFoeAttrV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeFoeAttrV8Config

// GetMazeFoeAttrV8Config pkg func. get one config by configId
func GetMazeFoeAttrV8Config(configId int32) *MazeFoeAttrV8ConfigRow {
	return gConfigData.GetMazeFoeAttrV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeFoeAttrV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeFoeAttrV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_foe_attr_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeFoeAttrV8Config pkg func. get all config slice
func GetAllMazeFoeAttrV8Config() []*MazeFoeAttrV8ConfigRow {
	return gConfigData.GetAllMazeFoeAttrV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeFoeAttrV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeFoeAttrV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeFoeAttrV8ConfigRow from maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_attr_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeFoeAttrV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_foe_attr_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_foe_attr_v8.json",
		"maze_foe_v8【迷宫-敌人信息】.xlsx", "maze_foe_attr_v8",
		&gMazeFoeAttrV8Parser{}, &gMazeFoeAttrV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeFoeAttrV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeFoeAttrV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeFoeAttrV8Config))(c)
		return true
	})
}

// RegisterMazeFoeAttrV8InitCallBack reg config update func (old func)
var RegisterMazeFoeAttrV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeFoeAttrV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeFoeAttrV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeFoeAttrV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeFoeAttrV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeFoeAttrV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeFoeAttrV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeFoeAttrV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeFoeAttrV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeFoeAttrV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeFoeAttrV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeFoeAttrV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeFoeAttrV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeFoeAttrV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeFoeAttrV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeFoeAttrV8ConfigRow", zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "maze_foe_attr_v8"))
		return
	}
	config, ok := container.(*MazeFoeAttrV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeFoeAttrV8Config")
		logger.ErrorWF("invalid type. not *MazeFoeAttrV8Config", zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "maze_foe_attr_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeFoeAttrV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeFoeAttrV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeFoeAttrV8Config")
		logger.ErrorWF("invalid type. not *MazeFoeAttrV8Config", zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "maze_foe_attr_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeFoeAttrV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeFoeAttrV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeFoeAttrV8Config")
		logger.ErrorWF("invalid type. not *MazeFoeAttrV8Config", zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "maze_foe_attr_v8"))
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
type gMazeFoeAttrV8Parser struct {
}

// New new config row data
func (*gMazeFoeAttrV8Parser) New() interface{} {
	return &MazeFoeAttrV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeFoeAttrV8Parser) Fields() []string {
	return gMazeFoeAttrV8Fields
}

// Parse parse raw data to row data
func (*gMazeFoeAttrV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeFoeAttrV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeFoeAttrV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeFoeAttrV8ConfigRow", zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "maze_foe_attr_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeFoeAttrV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeFoeAttrV8ConfigRow",
			zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"),
			zap.String("sheet", "maze_foe_attr_v8"), zap.Int("need_count", len(gMazeFoeAttrV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 属性id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 属性id to int32 failed")
			logger.ErrorWF("parse field order 属性id to int32 failed.",
				zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_attr_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 extra_attr : 额外属性列表
	if data[1] != "" {

		config.Extra_attr = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[1], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field extra_attr 额外属性列表 to key int32 failed")
				logger.ErrorWF("parse map field extra_attr 额外属性列表 to key int32 failed.",
					zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_attr_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field extra_attr 额外属性列表 to value int32 failed")
				logger.ErrorWF("parse map field extra_attr 额外属性列表 to value int32 failed.",
					zap.String("xlsx", "maze_foe_v8【迷宫-敌人信息】.xlsx"), zap.String("sheet", "maze_foe_attr_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Extra_attr[key] = value
		}
	}
	return
}

var gMazeFoeAttrV8Fields = []string{
	"order",
	"extra_attr",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeFoeAttrV8Parser{}
	loader := &gMazeFoeAttrV8Loader{}
	var data [][]string
	data, err = load("maze_foe_v8【迷宫-敌人信息】.xlsx", "maze_foe_attr_v8", gMazeFoeAttrV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_attr_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_attr_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeFoeAttrV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_attr_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeFoeAttrV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_attr_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_attr_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_attr_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_foe_v8【迷宫-敌人信息】.xlsx maze_foe_attr_v8 data success.")
	return
}
