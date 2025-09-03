package GMazeBariresDropConditionV8Cfg

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

// MazeBariresDropConditionV8ConfigRow from maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx maze_barires_drop_condition_v8
type MazeBariresDropConditionV8ConfigRow struct {
	Order          int32 `json:"order"`          // 条件id
	Condition_type int32 `json:"condition_type"` // 条件逻辑id
	Value          int32 `json:"value"`          // 条件参数
}

// MazeBariresDropConditionV8Config from maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx maze_barires_drop_condition_v8
type MazeBariresDropConditionV8Config struct {
	ConfigRows map[int32]*MazeBariresDropConditionV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeBariresDropConditionV8Config {
	ret := &MazeBariresDropConditionV8Config{ConfigRows: map[int32]*MazeBariresDropConditionV8ConfigRow{}}
	return ret
}

// GetMazeBariresDropConditionV8Config get one config by configId
func (c *MazeBariresDropConditionV8Config) GetMazeBariresDropConditionV8Config(configId int32) *MazeBariresDropConditionV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeBariresDropConditionV8Config) Get(configId int32) *MazeBariresDropConditionV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeBariresDropConditionV8Config get all config slice
func (c *MazeBariresDropConditionV8Config) GetAllMazeBariresDropConditionV8Config() (res []*MazeBariresDropConditionV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeBariresDropConditionV8Config) GetAll() (res []*MazeBariresDropConditionV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeBariresDropConditionV8Config

// GetMazeBariresDropConditionV8Config pkg func. get one config by configId
func GetMazeBariresDropConditionV8Config(configId int32) *MazeBariresDropConditionV8ConfigRow {
	return gConfigData.GetMazeBariresDropConditionV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeBariresDropConditionV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeBariresDropConditionV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_barires_drop_condition_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeBariresDropConditionV8Config pkg func. get all config slice
func GetAllMazeBariresDropConditionV8Config() []*MazeBariresDropConditionV8ConfigRow {
	return gConfigData.GetAllMazeBariresDropConditionV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeBariresDropConditionV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeBariresDropConditionV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeBariresDropConditionV8ConfigRow from maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx maze_barires_drop_condition_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeBariresDropConditionV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_barires_drop_condition_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_barires_drop_condition_v8.json",
		"maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx", "maze_barires_drop_condition_v8",
		&gMazeBariresDropConditionV8Parser{}, &gMazeBariresDropConditionV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeBariresDropConditionV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeBariresDropConditionV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeBariresDropConditionV8Config))(c)
		return true
	})
}

// RegisterMazeBariresDropConditionV8InitCallBack reg config update func (old func)
var RegisterMazeBariresDropConditionV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeBariresDropConditionV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeBariresDropConditionV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeBariresDropConditionV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeBariresDropConditionV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeBariresDropConditionV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeBariresDropConditionV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeBariresDropConditionV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeBariresDropConditionV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeBariresDropConditionV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeBariresDropConditionV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeBariresDropConditionV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeBariresDropConditionV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeBariresDropConditionV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBariresDropConditionV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBariresDropConditionV8ConfigRow", zap.String("xlsx", "maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx"),
			zap.String("sheet", "maze_barires_drop_condition_v8"))
		return
	}
	config, ok := container.(*MazeBariresDropConditionV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBariresDropConditionV8Config")
		logger.ErrorWF("invalid type. not *MazeBariresDropConditionV8Config", zap.String("xlsx", "maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx"),
			zap.String("sheet", "maze_barires_drop_condition_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeBariresDropConditionV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeBariresDropConditionV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBariresDropConditionV8Config")
		logger.ErrorWF("invalid type. not *MazeBariresDropConditionV8Config", zap.String("xlsx", "maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx"),
			zap.String("sheet", "maze_barires_drop_condition_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeBariresDropConditionV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeBariresDropConditionV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBariresDropConditionV8Config")
		logger.ErrorWF("invalid type. not *MazeBariresDropConditionV8Config", zap.String("xlsx", "maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx"),
			zap.String("sheet", "maze_barires_drop_condition_v8"))
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
type gMazeBariresDropConditionV8Parser struct {
}

// New new config row data
func (*gMazeBariresDropConditionV8Parser) New() interface{} {
	return &MazeBariresDropConditionV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeBariresDropConditionV8Parser) Fields() []string {
	return gMazeBariresDropConditionV8Fields
}

// Parse parse raw data to row data
func (*gMazeBariresDropConditionV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeBariresDropConditionV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBariresDropConditionV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBariresDropConditionV8ConfigRow", zap.String("xlsx", "maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx"),
			zap.String("sheet", "maze_barires_drop_condition_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeBariresDropConditionV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeBariresDropConditionV8ConfigRow",
			zap.String("xlsx", "maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx"),
			zap.String("sheet", "maze_barires_drop_condition_v8"), zap.Int("need_count", len(gMazeBariresDropConditionV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 条件id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 条件id to int32 failed")
			logger.ErrorWF("parse field order 条件id to int32 failed.",
				zap.String("xlsx", "maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx"), zap.String("sheet", "maze_barires_drop_condition_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 condition_type : 条件逻辑id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field condition_type 条件逻辑id to int32 failed")
			logger.ErrorWF("parse field condition_type 条件逻辑id to int32 failed.",
				zap.String("xlsx", "maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx"), zap.String("sheet", "maze_barires_drop_condition_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Condition_type = int32(tmp)
	}

	// parse column 2 value : 条件参数
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field value 条件参数 to int32 failed")
			logger.ErrorWF("parse field value 条件参数 to int32 failed.",
				zap.String("xlsx", "maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx"), zap.String("sheet", "maze_barires_drop_condition_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Value = int32(tmp)
	}
	return
}

var gMazeBariresDropConditionV8Fields = []string{
	"order",
	"condition_type",
	"value",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeBariresDropConditionV8Parser{}
	loader := &gMazeBariresDropConditionV8Loader{}
	var data [][]string
	data, err = load("maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx", "maze_barires_drop_condition_v8", gMazeBariresDropConditionV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx maze_barires_drop_condition_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx maze_barires_drop_condition_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeBariresDropConditionV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx maze_barires_drop_condition_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeBariresDropConditionV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx maze_barires_drop_condition_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx maze_barires_drop_condition_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx maze_barires_drop_condition_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_barires_drop_condition_v8【迷宫-关卡-道具掉落条件】.xlsx maze_barires_drop_condition_v8 data success.")
	return
}
