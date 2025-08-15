package GMazeAttrItemAttrV8Cfg

import (
	"errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"
)

// MazeAttrItemAttrV8ConfigRow from maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx maze_attr_item_attr_v8
type MazeAttrItemAttrV8ConfigRow struct {
	Order    int32 `json:"order"`    // 道具id
	Add_attr int32 `json:"add_attr"` // 技能属性id
}

// MazeAttrItemAttrV8Config from maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx maze_attr_item_attr_v8
type MazeAttrItemAttrV8Config struct {
	ConfigRows map[int32]*MazeAttrItemAttrV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeAttrItemAttrV8Config {
	ret := &MazeAttrItemAttrV8Config{ConfigRows: map[int32]*MazeAttrItemAttrV8ConfigRow{}}
	return ret
}

// GetMazeAttrItemAttrV8Config get one config by configId
func (c *MazeAttrItemAttrV8Config) GetMazeAttrItemAttrV8Config(configId int32) *MazeAttrItemAttrV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeAttrItemAttrV8Config) Get(configId int32) *MazeAttrItemAttrV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeAttrItemAttrV8Config get all config slice
func (c *MazeAttrItemAttrV8Config) GetAllMazeAttrItemAttrV8Config() (res []*MazeAttrItemAttrV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeAttrItemAttrV8Config) GetAll() (res []*MazeAttrItemAttrV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeAttrItemAttrV8Config

// GetMazeAttrItemAttrV8Config pkg func. get one config by configId
func GetMazeAttrItemAttrV8Config(configId int32) *MazeAttrItemAttrV8ConfigRow {
	return gConfigData.GetMazeAttrItemAttrV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeAttrItemAttrV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeAttrItemAttrV8Config pkg func. get all config slice
func GetAllMazeAttrItemAttrV8Config() []*MazeAttrItemAttrV8ConfigRow {
	return gConfigData.GetAllMazeAttrItemAttrV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeAttrItemAttrV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeAttrItemAttrV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeAttrItemAttrV8ConfigRow from maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx maze_attr_item_attr_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeAttrItemAttrV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_attr_item_attr_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_attr_item_attr_v8.json",
		"maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx", "maze_attr_item_attr_v8",
		&gMazeAttrItemAttrV8Parser{}, &gMazeAttrItemAttrV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeAttrItemAttrV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeAttrItemAttrV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeAttrItemAttrV8Config))(c)
		return true
	})
}

// RegisterMazeAttrItemAttrV8InitCallBack reg config update func (old func)
var RegisterMazeAttrItemAttrV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeAttrItemAttrV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeAttrItemAttrV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeAttrItemAttrV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeAttrItemAttrV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeAttrItemAttrV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeAttrItemAttrV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeAttrItemAttrV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeAttrItemAttrV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeAttrItemAttrV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeAttrItemAttrV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeAttrItemAttrV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeAttrItemAttrV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeAttrItemAttrV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrItemAttrV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAttrItemAttrV8ConfigRow", zap.String("xlsx", "maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx"),
			zap.String("sheet", "maze_attr_item_attr_v8"))
		return
	}
	config, ok := container.(*MazeAttrItemAttrV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrItemAttrV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrItemAttrV8Config", zap.String("xlsx", "maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx"),
			zap.String("sheet", "maze_attr_item_attr_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeAttrItemAttrV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeAttrItemAttrV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrItemAttrV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrItemAttrV8Config", zap.String("xlsx", "maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx"),
			zap.String("sheet", "maze_attr_item_attr_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeAttrItemAttrV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeAttrItemAttrV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrItemAttrV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrItemAttrV8Config", zap.String("xlsx", "maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx"),
			zap.String("sheet", "maze_attr_item_attr_v8"))
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
type gMazeAttrItemAttrV8Parser struct {
}

// New new config row data
func (*gMazeAttrItemAttrV8Parser) New() interface{} {
	return &MazeAttrItemAttrV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeAttrItemAttrV8Parser) Fields() []string {
	return gMazeAttrItemAttrV8Fields
}

// Parse parse raw data to row data
func (*gMazeAttrItemAttrV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeAttrItemAttrV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrItemAttrV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAttrItemAttrV8ConfigRow", zap.String("xlsx", "maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx"),
			zap.String("sheet", "maze_attr_item_attr_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeAttrItemAttrV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeAttrItemAttrV8ConfigRow",
			zap.String("xlsx", "maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx"),
			zap.String("sheet", "maze_attr_item_attr_v8"), zap.Int("need_count", len(gMazeAttrItemAttrV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 道具id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 道具id to int32 failed")
			logger.ErrorWF("parse field order 道具id to int32 failed.",
				zap.String("xlsx", "maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx"), zap.String("sheet", "maze_attr_item_attr_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 add_attr : 技能属性id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field add_attr 技能属性id to int32 failed")
			logger.ErrorWF("parse field add_attr 技能属性id to int32 failed.",
				zap.String("xlsx", "maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx"), zap.String("sheet", "maze_attr_item_attr_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Add_attr = int32(tmp)
	}
	return
}

var gMazeAttrItemAttrV8Fields = []string{
	"order",
	"add_attr",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeAttrItemAttrV8Parser{}
	loader := &gMazeAttrItemAttrV8Loader{}
	var data [][]string
	data, err = load("maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx", "maze_attr_item_attr_v8", gMazeAttrItemAttrV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx maze_attr_item_attr_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx maze_attr_item_attr_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeAttrItemAttrV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx maze_attr_item_attr_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeAttrItemAttrV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx maze_attr_item_attr_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx maze_attr_item_attr_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx maze_attr_item_attr_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_item_attr_v8【迷宫-道具-道具id对应属性id】.xlsx maze_attr_item_attr_v8 data success.")
	return
}
