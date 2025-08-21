package GMazeBrushPointConditionV8Cfg

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

// MazeBrushPointConditionV8ConfigRow from maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx maze_brush_point_condition_v8
type MazeBrushPointConditionV8ConfigRow struct {
	Brush_point_id       int32   `json:"brush_point_id"`       // 出怪点id
	Open_condition_type  int32   `json:"open_condition_type"`  // 开始刷怪的条件类型
	Open_condition_value []int32 `json:"open_condition_value"` // 条件参数
}

// MazeBrushPointConditionV8Config from maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx maze_brush_point_condition_v8
type MazeBrushPointConditionV8Config struct {
	ConfigRows map[int32]*MazeBrushPointConditionV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeBrushPointConditionV8Config {
	ret := &MazeBrushPointConditionV8Config{ConfigRows: map[int32]*MazeBrushPointConditionV8ConfigRow{}}
	return ret
}

// GetMazeBrushPointConditionV8Config get one config by configId
func (c *MazeBrushPointConditionV8Config) GetMazeBrushPointConditionV8Config(configId int32) *MazeBrushPointConditionV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeBrushPointConditionV8Config) Get(configId int32) *MazeBrushPointConditionV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeBrushPointConditionV8Config get all config slice
func (c *MazeBrushPointConditionV8Config) GetAllMazeBrushPointConditionV8Config() (res []*MazeBrushPointConditionV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeBrushPointConditionV8Config) GetAll() (res []*MazeBrushPointConditionV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeBrushPointConditionV8Config

// GetMazeBrushPointConditionV8Config pkg func. get one config by configId
func GetMazeBrushPointConditionV8Config(configId int32) *MazeBrushPointConditionV8ConfigRow {
	return gConfigData.GetMazeBrushPointConditionV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeBrushPointConditionV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeBrushPointConditionV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_brush_point_condition_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeBrushPointConditionV8Config pkg func. get all config slice
func GetAllMazeBrushPointConditionV8Config() []*MazeBrushPointConditionV8ConfigRow {
	return gConfigData.GetAllMazeBrushPointConditionV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeBrushPointConditionV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeBrushPointConditionV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeBrushPointConditionV8ConfigRow from maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx maze_brush_point_condition_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeBrushPointConditionV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_brush_point_condition_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_brush_point_condition_v8.json",
		"maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx", "maze_brush_point_condition_v8",
		&gMazeBrushPointConditionV8Parser{}, &gMazeBrushPointConditionV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeBrushPointConditionV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeBrushPointConditionV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeBrushPointConditionV8Config))(c)
		return true
	})
}

// RegisterMazeBrushPointConditionV8InitCallBack reg config update func (old func)
var RegisterMazeBrushPointConditionV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeBrushPointConditionV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeBrushPointConditionV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeBrushPointConditionV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeBrushPointConditionV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeBrushPointConditionV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeBrushPointConditionV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeBrushPointConditionV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeBrushPointConditionV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeBrushPointConditionV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeBrushPointConditionV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeBrushPointConditionV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeBrushPointConditionV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeBrushPointConditionV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushPointConditionV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBrushPointConditionV8ConfigRow", zap.String("xlsx", "maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx"),
			zap.String("sheet", "maze_brush_point_condition_v8"))
		return
	}
	config, ok := container.(*MazeBrushPointConditionV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushPointConditionV8Config")
		logger.ErrorWF("invalid type. not *MazeBrushPointConditionV8Config", zap.String("xlsx", "maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx"),
			zap.String("sheet", "maze_brush_point_condition_v8"))
		return
	}
	config.ConfigRows[row.Brush_point_id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeBrushPointConditionV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeBrushPointConditionV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushPointConditionV8Config")
		logger.ErrorWF("invalid type. not *MazeBrushPointConditionV8Config", zap.String("xlsx", "maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx"),
			zap.String("sheet", "maze_brush_point_condition_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeBrushPointConditionV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeBrushPointConditionV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushPointConditionV8Config")
		logger.ErrorWF("invalid type. not *MazeBrushPointConditionV8Config", zap.String("xlsx", "maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx"),
			zap.String("sheet", "maze_brush_point_condition_v8"))
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
type gMazeBrushPointConditionV8Parser struct {
}

// New new config row data
func (*gMazeBrushPointConditionV8Parser) New() interface{} {
	return &MazeBrushPointConditionV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeBrushPointConditionV8Parser) Fields() []string {
	return gMazeBrushPointConditionV8Fields
}

// Parse parse raw data to row data
func (*gMazeBrushPointConditionV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeBrushPointConditionV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBrushPointConditionV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBrushPointConditionV8ConfigRow", zap.String("xlsx", "maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx"),
			zap.String("sheet", "maze_brush_point_condition_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeBrushPointConditionV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeBrushPointConditionV8ConfigRow",
			zap.String("xlsx", "maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx"),
			zap.String("sheet", "maze_brush_point_condition_v8"), zap.Int("need_count", len(gMazeBrushPointConditionV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 brush_point_id : 出怪点id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field brush_point_id 出怪点id to int32 failed")
			logger.ErrorWF("parse field brush_point_id 出怪点id to int32 failed.",
				zap.String("xlsx", "maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx"), zap.String("sheet", "maze_brush_point_condition_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Brush_point_id = int32(tmp)
	}

	// parse column 1 open_condition_type : 开始刷怪的条件类型
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field open_condition_type 开始刷怪的条件类型 to int32 failed")
			logger.ErrorWF("parse field open_condition_type 开始刷怪的条件类型 to int32 failed.",
				zap.String("xlsx", "maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx"), zap.String("sheet", "maze_brush_point_condition_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Open_condition_type = int32(tmp)
	}

	// parse column 2 open_condition_value : 条件参数
	if data[2] != "" {

		vals := strings.Split(data[2], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field open_condition_value 条件参数 to []int32 failed")
				logger.ErrorWF("parse array field open_condition_value 条件参数 to []int32 failed.",
					zap.String("xlsx", "maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx"), zap.String("sheet", "maze_brush_point_condition_v8"),
					// zap.String("field_data",data[2]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Open_condition_value = append(config.Open_condition_value, int32(tmp))
		}
	}
	return
}

var gMazeBrushPointConditionV8Fields = []string{
	"brush_point_id",
	"open_condition_type",
	"open_condition_value",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeBrushPointConditionV8Parser{}
	loader := &gMazeBrushPointConditionV8Loader{}
	var data [][]string
	data, err = load("maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx", "maze_brush_point_condition_v8", gMazeBrushPointConditionV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx maze_brush_point_condition_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx maze_brush_point_condition_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeBrushPointConditionV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx maze_brush_point_condition_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeBrushPointConditionV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx maze_brush_point_condition_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx maze_brush_point_condition_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx maze_brush_point_condition_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_brush_point_condition_v8【迷宫-刷怪点刷怪条件】.xlsx maze_brush_point_condition_v8 data success.")
	return
}
