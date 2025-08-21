package GMazeEquDropV8Cfg

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

// MazeEquDropV8ConfigRow from maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx maze_equ_drop_v8
type MazeEquDropV8ConfigRow struct {
	Order           int32           `json:"order"`           // 装备掉落id
	Special_drop    []int32         `json:"special_drop"`    // 特殊掉落组（品质、部位）
	Regularity_drop map[int32]int32 `json:"regularity_drop"` // 常规掉落组（品质：权重）
}

// MazeEquDropV8Config from maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx maze_equ_drop_v8
type MazeEquDropV8Config struct {
	ConfigRows map[int32]*MazeEquDropV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquDropV8Config {
	ret := &MazeEquDropV8Config{ConfigRows: map[int32]*MazeEquDropV8ConfigRow{}}
	return ret
}

// GetMazeEquDropV8Config get one config by configId
func (c *MazeEquDropV8Config) GetMazeEquDropV8Config(configId int32) *MazeEquDropV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquDropV8Config) Get(configId int32) *MazeEquDropV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquDropV8Config get all config slice
func (c *MazeEquDropV8Config) GetAllMazeEquDropV8Config() (res []*MazeEquDropV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquDropV8Config) GetAll() (res []*MazeEquDropV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquDropV8Config

// GetMazeEquDropV8Config pkg func. get one config by configId
func GetMazeEquDropV8Config(configId int32) *MazeEquDropV8ConfigRow {
	return gConfigData.GetMazeEquDropV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquDropV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeEquDropV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_equ_drop_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeEquDropV8Config pkg func. get all config slice
func GetAllMazeEquDropV8Config() []*MazeEquDropV8ConfigRow {
	return gConfigData.GetAllMazeEquDropV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquDropV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquDropV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquDropV8ConfigRow from maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx maze_equ_drop_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquDropV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equ_drop_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equ_drop_v8.json",
		"maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx", "maze_equ_drop_v8",
		&gMazeEquDropV8Parser{}, &gMazeEquDropV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquDropV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquDropV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquDropV8Config))(c)
		return true
	})
}

// RegisterMazeEquDropV8InitCallBack reg config update func (old func)
var RegisterMazeEquDropV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquDropV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquDropV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquDropV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquDropV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquDropV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquDropV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquDropV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquDropV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquDropV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquDropV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquDropV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquDropV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquDropV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquDropV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquDropV8ConfigRow", zap.String("xlsx", "maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx"),
			zap.String("sheet", "maze_equ_drop_v8"))
		return
	}
	config, ok := container.(*MazeEquDropV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquDropV8Config")
		logger.ErrorWF("invalid type. not *MazeEquDropV8Config", zap.String("xlsx", "maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx"),
			zap.String("sheet", "maze_equ_drop_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquDropV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquDropV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquDropV8Config")
		logger.ErrorWF("invalid type. not *MazeEquDropV8Config", zap.String("xlsx", "maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx"),
			zap.String("sheet", "maze_equ_drop_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquDropV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquDropV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquDropV8Config")
		logger.ErrorWF("invalid type. not *MazeEquDropV8Config", zap.String("xlsx", "maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx"),
			zap.String("sheet", "maze_equ_drop_v8"))
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
type gMazeEquDropV8Parser struct {
}

// New new config row data
func (*gMazeEquDropV8Parser) New() interface{} {
	return &MazeEquDropV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquDropV8Parser) Fields() []string {
	return gMazeEquDropV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquDropV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquDropV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquDropV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquDropV8ConfigRow", zap.String("xlsx", "maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx"),
			zap.String("sheet", "maze_equ_drop_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquDropV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquDropV8ConfigRow",
			zap.String("xlsx", "maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx"),
			zap.String("sheet", "maze_equ_drop_v8"), zap.Int("need_count", len(gMazeEquDropV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 装备掉落id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 装备掉落id to int32 failed")
			logger.ErrorWF("parse field order 装备掉落id to int32 failed.",
				zap.String("xlsx", "maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx"), zap.String("sheet", "maze_equ_drop_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 special_drop : 特殊掉落组（品质、部位）
	if data[1] != "" {

		vals := strings.Split(data[1], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field special_drop 特殊掉落组（品质、部位） to []int32 failed")
				logger.ErrorWF("parse array field special_drop 特殊掉落组（品质、部位） to []int32 failed.",
					zap.String("xlsx", "maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx"), zap.String("sheet", "maze_equ_drop_v8"),
					// zap.String("field_data",data[1]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Special_drop = append(config.Special_drop, int32(tmp))
		}
	}

	// parse column 2 regularity_drop : 常规掉落组（品质：权重）
	if data[2] != "" {

		config.Regularity_drop = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[2], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field regularity_drop 常规掉落组（品质：权重） to key int32 failed")
				logger.ErrorWF("parse map field regularity_drop 常规掉落组（品质：权重） to key int32 failed.",
					zap.String("xlsx", "maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx"), zap.String("sheet", "maze_equ_drop_v8"),
					// zap.String("field_data",data[2]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field regularity_drop 常规掉落组（品质：权重） to value int32 failed")
				logger.ErrorWF("parse map field regularity_drop 常规掉落组（品质：权重） to value int32 failed.",
					zap.String("xlsx", "maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx"), zap.String("sheet", "maze_equ_drop_v8"),
					// zap.String("field_data",data[2]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Regularity_drop[key] = value
		}
	}
	return
}

var gMazeEquDropV8Fields = []string{
	"order",
	"special_drop",
	"regularity_drop",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquDropV8Parser{}
	loader := &gMazeEquDropV8Loader{}
	var data [][]string
	data, err = load("maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx", "maze_equ_drop_v8", gMazeEquDropV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx maze_equ_drop_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx maze_equ_drop_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquDropV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx maze_equ_drop_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquDropV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx maze_equ_drop_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx maze_equ_drop_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx maze_equ_drop_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equ_drop_v8【迷宫-装备掉落概率】.xlsx maze_equ_drop_v8 data success.")
	return
}
