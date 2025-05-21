package GMazeEquipAffixOrderV8Cfg

import (
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

// MazeEquipAffixOrderV8ConfigRow from maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx maze_equip_affix_order_v8
type MazeEquipAffixOrderV8ConfigRow struct {
	Pos        int32           `json:"pos"`        // 装备部位类型
	Attr_order map[int32]int32 `json:"attr_order"` // 属性id:排序
}

// MazeEquipAffixOrderV8Config from maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx maze_equip_affix_order_v8
type MazeEquipAffixOrderV8Config struct {
	ConfigRows map[int32]*MazeEquipAffixOrderV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipAffixOrderV8Config {
	ret := &MazeEquipAffixOrderV8Config{ConfigRows: map[int32]*MazeEquipAffixOrderV8ConfigRow{}}
	return ret
}

// GetMazeEquipAffixOrderV8Config get one config by configId
func (c *MazeEquipAffixOrderV8Config) GetMazeEquipAffixOrderV8Config(configId int32) *MazeEquipAffixOrderV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipAffixOrderV8Config) Get(configId int32) *MazeEquipAffixOrderV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipAffixOrderV8Config get all config slice
func (c *MazeEquipAffixOrderV8Config) GetAllMazeEquipAffixOrderV8Config() (res []*MazeEquipAffixOrderV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipAffixOrderV8Config) GetAll() (res []*MazeEquipAffixOrderV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipAffixOrderV8Config

// GetMazeEquipAffixOrderV8Config pkg func. get one config by configId
func GetMazeEquipAffixOrderV8Config(configId int32) *MazeEquipAffixOrderV8ConfigRow {
	return gConfigData.GetMazeEquipAffixOrderV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipAffixOrderV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeEquipAffixOrderV8Config pkg func. get all config slice
func GetAllMazeEquipAffixOrderV8Config() []*MazeEquipAffixOrderV8ConfigRow {
	return gConfigData.GetAllMazeEquipAffixOrderV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipAffixOrderV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipAffixOrderV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipAffixOrderV8ConfigRow from maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx maze_equip_affix_order_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipAffixOrderV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_affix_order_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_affix_order_v8.json",
		"maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx", "maze_equip_affix_order_v8",
		&gMazeEquipAffixOrderV8Parser{}, &gMazeEquipAffixOrderV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipAffixOrderV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipAffixOrderV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipAffixOrderV8Config))(c)
		return true
	})
}

// RegisterMazeEquipAffixOrderV8InitCallBack reg config update func (old func)
var RegisterMazeEquipAffixOrderV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipAffixOrderV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipAffixOrderV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipAffixOrderV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipAffixOrderV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipAffixOrderV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipAffixOrderV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipAffixOrderV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipAffixOrderV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipAffixOrderV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipAffixOrderV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipAffixOrderV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipAffixOrderV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipAffixOrderV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixOrderV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipAffixOrderV8ConfigRow", zap.String("xlsx", "maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx"),
			zap.String("sheet", "maze_equip_affix_order_v8"))
		return
	}
	config, ok := container.(*MazeEquipAffixOrderV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixOrderV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixOrderV8Config", zap.String("xlsx", "maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx"),
			zap.String("sheet", "maze_equip_affix_order_v8"))
		return
	}
	config.ConfigRows[row.Pos] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipAffixOrderV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipAffixOrderV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixOrderV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixOrderV8Config", zap.String("xlsx", "maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx"),
			zap.String("sheet", "maze_equip_affix_order_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipAffixOrderV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipAffixOrderV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixOrderV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixOrderV8Config", zap.String("xlsx", "maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx"),
			zap.String("sheet", "maze_equip_affix_order_v8"))
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
type gMazeEquipAffixOrderV8Parser struct {
}

// New new config row data
func (*gMazeEquipAffixOrderV8Parser) New() interface{} {
	return &MazeEquipAffixOrderV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipAffixOrderV8Parser) Fields() []string {
	return gMazeEquipAffixOrderV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipAffixOrderV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipAffixOrderV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixOrderV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipAffixOrderV8ConfigRow", zap.String("xlsx", "maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx"),
			zap.String("sheet", "maze_equip_affix_order_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipAffixOrderV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipAffixOrderV8ConfigRow",
			zap.String("xlsx", "maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx"),
			zap.String("sheet", "maze_equip_affix_order_v8"), zap.Int("need_count", len(gMazeEquipAffixOrderV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 pos : 装备部位类型
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field pos 装备部位类型 to int32 failed")
			logger.ErrorWF("parse field pos 装备部位类型 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx"), zap.String("sheet", "maze_equip_affix_order_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Pos = int32(tmp)
	}

	// parse column 1 attr_order : 属性id:排序
	if data[1] != "" {

		config.Attr_order = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[1], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr_order 属性id:排序 to key int32 failed")
				logger.ErrorWF("parse map field attr_order 属性id:排序 to key int32 failed.",
					zap.String("xlsx", "maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx"), zap.String("sheet", "maze_equip_affix_order_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr_order 属性id:排序 to value int32 failed")
				logger.ErrorWF("parse map field attr_order 属性id:排序 to value int32 failed.",
					zap.String("xlsx", "maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx"), zap.String("sheet", "maze_equip_affix_order_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Attr_order[key] = value
		}
	}
	return
}

var gMazeEquipAffixOrderV8Fields = []string{
	"pos",
	"attr_order",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipAffixOrderV8Parser{}
	loader := &gMazeEquipAffixOrderV8Loader{}
	var data [][]string
	data, err = load("maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx", "maze_equip_affix_order_v8", gMazeEquipAffixOrderV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx maze_equip_affix_order_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx maze_equip_affix_order_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipAffixOrderV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx maze_equip_affix_order_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipAffixOrderV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx maze_equip_affix_order_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx maze_equip_affix_order_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx maze_equip_affix_order_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_affix_order_v8【迷宫-装备-词条排序】.xlsx maze_equip_affix_order_v8 data success.")
	return
}
