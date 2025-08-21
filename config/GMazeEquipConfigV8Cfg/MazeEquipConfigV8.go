package GMazeEquipConfigV8Cfg

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

// MazeEquipConfigV8ConfigRow from maze_equip_config_v8【迷宫-装备-通用配置】.xlsx maze_equip_config_v8
type MazeEquipConfigV8ConfigRow struct {
	Order      int32           `json:"order"`      // 序号
	Value_str  string          `json:"value_str"`  // 参数string
	Value_map  map[int32]int64 `json:"value_map"`  // map参数
	Value_int  int64           `json:"value_int"`  // 参数
	Value_list []int64         `json:"value_list"` // 参数list
}

// MazeEquipConfigV8Config from maze_equip_config_v8【迷宫-装备-通用配置】.xlsx maze_equip_config_v8
type MazeEquipConfigV8Config struct {
	ConfigRows map[int32]*MazeEquipConfigV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipConfigV8Config {
	ret := &MazeEquipConfigV8Config{ConfigRows: map[int32]*MazeEquipConfigV8ConfigRow{}}
	return ret
}

// GetMazeEquipConfigV8Config get one config by configId
func (c *MazeEquipConfigV8Config) GetMazeEquipConfigV8Config(configId int32) *MazeEquipConfigV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipConfigV8Config) Get(configId int32) *MazeEquipConfigV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipConfigV8Config get all config slice
func (c *MazeEquipConfigV8Config) GetAllMazeEquipConfigV8Config() (res []*MazeEquipConfigV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipConfigV8Config) GetAll() (res []*MazeEquipConfigV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipConfigV8Config

// GetMazeEquipConfigV8Config pkg func. get one config by configId
func GetMazeEquipConfigV8Config(configId int32) *MazeEquipConfigV8ConfigRow {
	return gConfigData.GetMazeEquipConfigV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipConfigV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeEquipConfigV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_equip_config_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeEquipConfigV8Config pkg func. get all config slice
func GetAllMazeEquipConfigV8Config() []*MazeEquipConfigV8ConfigRow {
	return gConfigData.GetAllMazeEquipConfigV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipConfigV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipConfigV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipConfigV8ConfigRow from maze_equip_config_v8【迷宫-装备-通用配置】.xlsx maze_equip_config_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipConfigV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_config_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_config_v8.json",
		"maze_equip_config_v8【迷宫-装备-通用配置】.xlsx", "maze_equip_config_v8",
		&gMazeEquipConfigV8Parser{}, &gMazeEquipConfigV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipConfigV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipConfigV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipConfigV8Config))(c)
		return true
	})
}

// RegisterMazeEquipConfigV8InitCallBack reg config update func (old func)
var RegisterMazeEquipConfigV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipConfigV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipConfigV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipConfigV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipConfigV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipConfigV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipConfigV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipConfigV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipConfigV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipConfigV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipConfigV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipConfigV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipConfigV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipConfigV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipConfigV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipConfigV8ConfigRow", zap.String("xlsx", "maze_equip_config_v8【迷宫-装备-通用配置】.xlsx"),
			zap.String("sheet", "maze_equip_config_v8"))
		return
	}
	config, ok := container.(*MazeEquipConfigV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipConfigV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipConfigV8Config", zap.String("xlsx", "maze_equip_config_v8【迷宫-装备-通用配置】.xlsx"),
			zap.String("sheet", "maze_equip_config_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipConfigV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipConfigV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipConfigV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipConfigV8Config", zap.String("xlsx", "maze_equip_config_v8【迷宫-装备-通用配置】.xlsx"),
			zap.String("sheet", "maze_equip_config_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipConfigV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipConfigV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipConfigV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipConfigV8Config", zap.String("xlsx", "maze_equip_config_v8【迷宫-装备-通用配置】.xlsx"),
			zap.String("sheet", "maze_equip_config_v8"))
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
type gMazeEquipConfigV8Parser struct {
}

// New new config row data
func (*gMazeEquipConfigV8Parser) New() interface{} {
	return &MazeEquipConfigV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipConfigV8Parser) Fields() []string {
	return gMazeEquipConfigV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipConfigV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipConfigV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipConfigV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipConfigV8ConfigRow", zap.String("xlsx", "maze_equip_config_v8【迷宫-装备-通用配置】.xlsx"),
			zap.String("sheet", "maze_equip_config_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipConfigV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipConfigV8ConfigRow",
			zap.String("xlsx", "maze_equip_config_v8【迷宫-装备-通用配置】.xlsx"),
			zap.String("sheet", "maze_equip_config_v8"), zap.Int("need_count", len(gMazeEquipConfigV8Fields)),
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
				zap.String("xlsx", "maze_equip_config_v8【迷宫-装备-通用配置】.xlsx"), zap.String("sheet", "maze_equip_config_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 value_str : 参数string
	if data[1] != "" {
		config.Value_str = data[1]
	}

	// parse column 2 value_map : map参数
	if data[2] != "" {

		config.Value_map = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[2], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field value_map map参数 to key int32 failed")
				logger.ErrorWF("parse map field value_map map参数 to key int32 failed.",
					zap.String("xlsx", "maze_equip_config_v8【迷宫-装备-通用配置】.xlsx"), zap.String("sheet", "maze_equip_config_v8"),
					// zap.String("field_data",data[2]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field value_map map参数 to value int64 failed")
				logger.ErrorWF("parse map field value_map map参数 to value int64 failed.",
					zap.String("xlsx", "maze_equip_config_v8【迷宫-装备-通用配置】.xlsx"), zap.String("sheet", "maze_equip_config_v8"),
					// zap.String("field_data",data[2]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Value_map[key] = value
		}
	}

	// parse column 3 value_int : 参数
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field value_int 参数 to int64 failed")
			logger.ErrorWF("parse field value_int 参数 to int64 failed.",
				zap.String("xlsx", "maze_equip_config_v8【迷宫-装备-通用配置】.xlsx"), zap.String("sheet", "maze_equip_config_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Value_int = int64(tmp)
	}

	// parse column 4 value_list : 参数list
	if data[4] != "" {

		vals := strings.Split(data[4], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field value_list 参数list to []int64 failed")
				logger.ErrorWF("parse array field value_list 参数list to []int64 failed.",
					zap.String("xlsx", "maze_equip_config_v8【迷宫-装备-通用配置】.xlsx"), zap.String("sheet", "maze_equip_config_v8"),
					// zap.String("field_data",data[4]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Value_list = append(config.Value_list, int64(tmp))
		}
	}
	return
}

var gMazeEquipConfigV8Fields = []string{
	"order",
	"value_str",
	"value_map",
	"value_int",
	"value_list",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipConfigV8Parser{}
	loader := &gMazeEquipConfigV8Loader{}
	var data [][]string
	data, err = load("maze_equip_config_v8【迷宫-装备-通用配置】.xlsx", "maze_equip_config_v8", gMazeEquipConfigV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_config_v8【迷宫-装备-通用配置】.xlsx maze_equip_config_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_config_v8【迷宫-装备-通用配置】.xlsx maze_equip_config_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipConfigV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_config_v8【迷宫-装备-通用配置】.xlsx maze_equip_config_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipConfigV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_config_v8【迷宫-装备-通用配置】.xlsx maze_equip_config_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_config_v8【迷宫-装备-通用配置】.xlsx maze_equip_config_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_config_v8【迷宫-装备-通用配置】.xlsx maze_equip_config_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_config_v8【迷宫-装备-通用配置】.xlsx maze_equip_config_v8 data success.")
	return
}
