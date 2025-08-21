package GMazeAttrListTypeV8Cfg

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

// MazeAttrListTypeV8ConfigRow from maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx maze_attr_list_type_v8
type MazeAttrListTypeV8ConfigRow struct {
	Type      int32  `json:"type"`      // 类型
	Type_name string `json:"type_name"` // 类型名称
}

// MazeAttrListTypeV8Config from maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx maze_attr_list_type_v8
type MazeAttrListTypeV8Config struct {
	ConfigRows map[int32]*MazeAttrListTypeV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeAttrListTypeV8Config {
	ret := &MazeAttrListTypeV8Config{ConfigRows: map[int32]*MazeAttrListTypeV8ConfigRow{}}
	return ret
}

// GetMazeAttrListTypeV8Config get one config by configId
func (c *MazeAttrListTypeV8Config) GetMazeAttrListTypeV8Config(configId int32) *MazeAttrListTypeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeAttrListTypeV8Config) Get(configId int32) *MazeAttrListTypeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeAttrListTypeV8Config get all config slice
func (c *MazeAttrListTypeV8Config) GetAllMazeAttrListTypeV8Config() (res []*MazeAttrListTypeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeAttrListTypeV8Config) GetAll() (res []*MazeAttrListTypeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeAttrListTypeV8Config

// GetMazeAttrListTypeV8Config pkg func. get one config by configId
func GetMazeAttrListTypeV8Config(configId int32) *MazeAttrListTypeV8ConfigRow {
	return gConfigData.GetMazeAttrListTypeV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeAttrListTypeV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeAttrListTypeV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_attr_list_type_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeAttrListTypeV8Config pkg func. get all config slice
func GetAllMazeAttrListTypeV8Config() []*MazeAttrListTypeV8ConfigRow {
	return gConfigData.GetAllMazeAttrListTypeV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeAttrListTypeV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeAttrListTypeV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeAttrListTypeV8ConfigRow from maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx maze_attr_list_type_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeAttrListTypeV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_attr_list_type_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_attr_list_type_v8.json",
		"maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx", "maze_attr_list_type_v8",
		&gMazeAttrListTypeV8Parser{}, &gMazeAttrListTypeV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeAttrListTypeV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeAttrListTypeV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeAttrListTypeV8Config))(c)
		return true
	})
}

// RegisterMazeAttrListTypeV8InitCallBack reg config update func (old func)
var RegisterMazeAttrListTypeV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeAttrListTypeV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeAttrListTypeV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeAttrListTypeV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeAttrListTypeV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeAttrListTypeV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeAttrListTypeV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeAttrListTypeV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeAttrListTypeV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeAttrListTypeV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeAttrListTypeV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeAttrListTypeV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeAttrListTypeV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeAttrListTypeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrListTypeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAttrListTypeV8ConfigRow", zap.String("xlsx", "maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx"),
			zap.String("sheet", "maze_attr_list_type_v8"))
		return
	}
	config, ok := container.(*MazeAttrListTypeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrListTypeV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrListTypeV8Config", zap.String("xlsx", "maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx"),
			zap.String("sheet", "maze_attr_list_type_v8"))
		return
	}
	config.ConfigRows[row.Type] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeAttrListTypeV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeAttrListTypeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrListTypeV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrListTypeV8Config", zap.String("xlsx", "maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx"),
			zap.String("sheet", "maze_attr_list_type_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeAttrListTypeV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeAttrListTypeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrListTypeV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrListTypeV8Config", zap.String("xlsx", "maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx"),
			zap.String("sheet", "maze_attr_list_type_v8"))
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
type gMazeAttrListTypeV8Parser struct {
}

// New new config row data
func (*gMazeAttrListTypeV8Parser) New() interface{} {
	return &MazeAttrListTypeV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeAttrListTypeV8Parser) Fields() []string {
	return gMazeAttrListTypeV8Fields
}

// Parse parse raw data to row data
func (*gMazeAttrListTypeV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeAttrListTypeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrListTypeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAttrListTypeV8ConfigRow", zap.String("xlsx", "maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx"),
			zap.String("sheet", "maze_attr_list_type_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeAttrListTypeV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeAttrListTypeV8ConfigRow",
			zap.String("xlsx", "maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx"),
			zap.String("sheet", "maze_attr_list_type_v8"), zap.Int("need_count", len(gMazeAttrListTypeV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 type : 类型
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field type 类型 to int32 failed")
			logger.ErrorWF("parse field type 类型 to int32 failed.",
				zap.String("xlsx", "maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx"), zap.String("sheet", "maze_attr_list_type_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Type = int32(tmp)
	}

	// parse column 1 type_name : 类型名称
	if data[1] != "" {
		config.Type_name = data[1]
	}
	return
}

var gMazeAttrListTypeV8Fields = []string{
	"type",
	"type_name",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeAttrListTypeV8Parser{}
	loader := &gMazeAttrListTypeV8Loader{}
	var data [][]string
	data, err = load("maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx", "maze_attr_list_type_v8", gMazeAttrListTypeV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx maze_attr_list_type_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx maze_attr_list_type_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeAttrListTypeV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx maze_attr_list_type_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeAttrListTypeV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx maze_attr_list_type_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx maze_attr_list_type_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx maze_attr_list_type_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_attr_list_type_v8【迷宫-属性列表-分类】.xlsx maze_attr_list_type_v8 data success.")
	return
}
