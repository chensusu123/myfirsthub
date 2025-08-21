package GMazeLevelV8Cfg

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

// MazeLevelV8ConfigRow from maze_level_v8【迷宫-探险等级】.xlsx maze_level_v8
type MazeLevelV8ConfigRow struct {
	Order               int32           `json:"order"`               // 等级
	Next_level_need_exp int64           `json:"next_level_need_exp"` // 升到下一级需要的经验值
	All_exp             int64           `json:"all_exp"`             // 总经验
	Attr                map[int32]int64 `json:"attr"`                // 属性id:属性值
}

// MazeLevelV8Config from maze_level_v8【迷宫-探险等级】.xlsx maze_level_v8
type MazeLevelV8Config struct {
	ConfigRows map[int32]*MazeLevelV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeLevelV8Config {
	ret := &MazeLevelV8Config{ConfigRows: map[int32]*MazeLevelV8ConfigRow{}}
	return ret
}

// GetMazeLevelV8Config get one config by configId
func (c *MazeLevelV8Config) GetMazeLevelV8Config(configId int32) *MazeLevelV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeLevelV8Config) Get(configId int32) *MazeLevelV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeLevelV8Config get all config slice
func (c *MazeLevelV8Config) GetAllMazeLevelV8Config() (res []*MazeLevelV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeLevelV8Config) GetAll() (res []*MazeLevelV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeLevelV8Config

// GetMazeLevelV8Config pkg func. get one config by configId
func GetMazeLevelV8Config(configId int32) *MazeLevelV8ConfigRow {
	return gConfigData.GetMazeLevelV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeLevelV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeLevelV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_level_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeLevelV8Config pkg func. get all config slice
func GetAllMazeLevelV8Config() []*MazeLevelV8ConfigRow {
	return gConfigData.GetAllMazeLevelV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeLevelV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeLevelV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeLevelV8ConfigRow from maze_level_v8【迷宫-探险等级】.xlsx maze_level_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeLevelV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_level_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_level_v8.json",
		"maze_level_v8【迷宫-探险等级】.xlsx", "maze_level_v8",
		&gMazeLevelV8Parser{}, &gMazeLevelV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeLevelV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeLevelV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeLevelV8Config))(c)
		return true
	})
}

// RegisterMazeLevelV8InitCallBack reg config update func (old func)
var RegisterMazeLevelV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeLevelV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeLevelV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeLevelV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeLevelV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeLevelV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeLevelV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeLevelV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeLevelV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeLevelV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeLevelV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeLevelV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeLevelV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeLevelV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeLevelV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeLevelV8ConfigRow", zap.String("xlsx", "maze_level_v8【迷宫-探险等级】.xlsx"),
			zap.String("sheet", "maze_level_v8"))
		return
	}
	config, ok := container.(*MazeLevelV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeLevelV8Config")
		logger.ErrorWF("invalid type. not *MazeLevelV8Config", zap.String("xlsx", "maze_level_v8【迷宫-探险等级】.xlsx"),
			zap.String("sheet", "maze_level_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeLevelV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeLevelV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeLevelV8Config")
		logger.ErrorWF("invalid type. not *MazeLevelV8Config", zap.String("xlsx", "maze_level_v8【迷宫-探险等级】.xlsx"),
			zap.String("sheet", "maze_level_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeLevelV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeLevelV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeLevelV8Config")
		logger.ErrorWF("invalid type. not *MazeLevelV8Config", zap.String("xlsx", "maze_level_v8【迷宫-探险等级】.xlsx"),
			zap.String("sheet", "maze_level_v8"))
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
type gMazeLevelV8Parser struct {
}

// New new config row data
func (*gMazeLevelV8Parser) New() interface{} {
	return &MazeLevelV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeLevelV8Parser) Fields() []string {
	return gMazeLevelV8Fields
}

// Parse parse raw data to row data
func (*gMazeLevelV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeLevelV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeLevelV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeLevelV8ConfigRow", zap.String("xlsx", "maze_level_v8【迷宫-探险等级】.xlsx"),
			zap.String("sheet", "maze_level_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeLevelV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeLevelV8ConfigRow",
			zap.String("xlsx", "maze_level_v8【迷宫-探险等级】.xlsx"),
			zap.String("sheet", "maze_level_v8"), zap.Int("need_count", len(gMazeLevelV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 等级
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 等级 to int32 failed")
			logger.ErrorWF("parse field order 等级 to int32 failed.",
				zap.String("xlsx", "maze_level_v8【迷宫-探险等级】.xlsx"), zap.String("sheet", "maze_level_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 next_level_need_exp : 升到下一级需要的经验值
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field next_level_need_exp 升到下一级需要的经验值 to int64 failed")
			logger.ErrorWF("parse field next_level_need_exp 升到下一级需要的经验值 to int64 failed.",
				zap.String("xlsx", "maze_level_v8【迷宫-探险等级】.xlsx"), zap.String("sheet", "maze_level_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Next_level_need_exp = int64(tmp)
	}

	// parse column 2 all_exp : 总经验
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field all_exp 总经验 to int64 failed")
			logger.ErrorWF("parse field all_exp 总经验 to int64 failed.",
				zap.String("xlsx", "maze_level_v8【迷宫-探险等级】.xlsx"), zap.String("sheet", "maze_level_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.All_exp = int64(tmp)
	}

	// parse column 3 attr : 属性id:属性值
	if data[3] != "" {

		config.Attr = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr 属性id:属性值 to key int32 failed")
				logger.ErrorWF("parse map field attr 属性id:属性值 to key int32 failed.",
					zap.String("xlsx", "maze_level_v8【迷宫-探险等级】.xlsx"), zap.String("sheet", "maze_level_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr 属性id:属性值 to value int64 failed")
				logger.ErrorWF("parse map field attr 属性id:属性值 to value int64 failed.",
					zap.String("xlsx", "maze_level_v8【迷宫-探险等级】.xlsx"), zap.String("sheet", "maze_level_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Attr[key] = value
		}
	}
	return
}

var gMazeLevelV8Fields = []string{
	"order",
	"next_level_need_exp",
	"all_exp",
	"attr",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeLevelV8Parser{}
	loader := &gMazeLevelV8Loader{}
	var data [][]string
	data, err = load("maze_level_v8【迷宫-探险等级】.xlsx", "maze_level_v8", gMazeLevelV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_level_v8【迷宫-探险等级】.xlsx maze_level_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_level_v8【迷宫-探险等级】.xlsx maze_level_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeLevelV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_level_v8【迷宫-探险等级】.xlsx maze_level_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeLevelV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_level_v8【迷宫-探险等级】.xlsx maze_level_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_level_v8【迷宫-探险等级】.xlsx maze_level_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_level_v8【迷宫-探险等级】.xlsx maze_level_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_level_v8【迷宫-探险等级】.xlsx maze_level_v8 data success.")
	return
}
