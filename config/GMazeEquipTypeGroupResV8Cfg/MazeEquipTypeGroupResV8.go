package GMazeEquipTypeGroupResV8Cfg

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

// MazeEquipTypeGroupResV8ConfigRow from maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_group_res_v8
type MazeEquipTypeGroupResV8ConfigRow struct {
	Order            int32           `json:"order"`            // 序号
	Maze_model_group map[int32]int32 `json:"maze_model_group"` // 迷宫模型资源列表
}

// MazeEquipTypeGroupResV8Config from maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_group_res_v8
type MazeEquipTypeGroupResV8Config struct {
	ConfigRows map[int32]*MazeEquipTypeGroupResV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipTypeGroupResV8Config {
	ret := &MazeEquipTypeGroupResV8Config{ConfigRows: map[int32]*MazeEquipTypeGroupResV8ConfigRow{}}
	return ret
}

// GetMazeEquipTypeGroupResV8Config get one config by configId
func (c *MazeEquipTypeGroupResV8Config) GetMazeEquipTypeGroupResV8Config(configId int32) *MazeEquipTypeGroupResV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipTypeGroupResV8Config) Get(configId int32) *MazeEquipTypeGroupResV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipTypeGroupResV8Config get all config slice
func (c *MazeEquipTypeGroupResV8Config) GetAllMazeEquipTypeGroupResV8Config() (res []*MazeEquipTypeGroupResV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipTypeGroupResV8Config) GetAll() (res []*MazeEquipTypeGroupResV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipTypeGroupResV8Config

// GetMazeEquipTypeGroupResV8Config pkg func. get one config by configId
func GetMazeEquipTypeGroupResV8Config(configId int32) *MazeEquipTypeGroupResV8ConfigRow {
	return gConfigData.GetMazeEquipTypeGroupResV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipTypeGroupResV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeEquipTypeGroupResV8Config pkg func. get all config slice
func GetAllMazeEquipTypeGroupResV8Config() []*MazeEquipTypeGroupResV8ConfigRow {
	return gConfigData.GetAllMazeEquipTypeGroupResV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipTypeGroupResV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipTypeGroupResV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipTypeGroupResV8ConfigRow from maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_group_res_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipTypeGroupResV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_type_group_res_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_type_group_res_v8.json",
		"maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx", "maze_equip_type_group_res_v8",
		&gMazeEquipTypeGroupResV8Parser{}, &gMazeEquipTypeGroupResV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipTypeGroupResV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipTypeGroupResV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipTypeGroupResV8Config))(c)
		return true
	})
}

// RegisterMazeEquipTypeGroupResV8InitCallBack reg config update func (old func)
var RegisterMazeEquipTypeGroupResV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipTypeGroupResV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipTypeGroupResV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipTypeGroupResV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipTypeGroupResV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipTypeGroupResV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipTypeGroupResV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipTypeGroupResV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipTypeGroupResV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipTypeGroupResV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipTypeGroupResV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipTypeGroupResV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipTypeGroupResV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipTypeGroupResV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipTypeGroupResV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipTypeGroupResV8ConfigRow", zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"),
			zap.String("sheet", "maze_equip_type_group_res_v8"))
		return
	}
	config, ok := container.(*MazeEquipTypeGroupResV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipTypeGroupResV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipTypeGroupResV8Config", zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"),
			zap.String("sheet", "maze_equip_type_group_res_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipTypeGroupResV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipTypeGroupResV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipTypeGroupResV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipTypeGroupResV8Config", zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"),
			zap.String("sheet", "maze_equip_type_group_res_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipTypeGroupResV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipTypeGroupResV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipTypeGroupResV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipTypeGroupResV8Config", zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"),
			zap.String("sheet", "maze_equip_type_group_res_v8"))
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
type gMazeEquipTypeGroupResV8Parser struct {
}

// New new config row data
func (*gMazeEquipTypeGroupResV8Parser) New() interface{} {
	return &MazeEquipTypeGroupResV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipTypeGroupResV8Parser) Fields() []string {
	return gMazeEquipTypeGroupResV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipTypeGroupResV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipTypeGroupResV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipTypeGroupResV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipTypeGroupResV8ConfigRow", zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"),
			zap.String("sheet", "maze_equip_type_group_res_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipTypeGroupResV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipTypeGroupResV8ConfigRow",
			zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"),
			zap.String("sheet", "maze_equip_type_group_res_v8"), zap.Int("need_count", len(gMazeEquipTypeGroupResV8Fields)),
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
				zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"), zap.String("sheet", "maze_equip_type_group_res_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 maze_model_group : 迷宫模型资源列表
	if data[1] != "" {

		config.Maze_model_group = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[1], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field maze_model_group 迷宫模型资源列表 to key int32 failed")
				logger.ErrorWF("parse map field maze_model_group 迷宫模型资源列表 to key int32 failed.",
					zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"), zap.String("sheet", "maze_equip_type_group_res_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field maze_model_group 迷宫模型资源列表 to value int32 failed")
				logger.ErrorWF("parse map field maze_model_group 迷宫模型资源列表 to value int32 failed.",
					zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"), zap.String("sheet", "maze_equip_type_group_res_v8"),
					// zap.String("field_data",data[1]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Maze_model_group[key] = value
		}
	}
	return
}

var gMazeEquipTypeGroupResV8Fields = []string{
	"order",
	"maze_model_group",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipTypeGroupResV8Parser{}
	loader := &gMazeEquipTypeGroupResV8Loader{}
	var data [][]string
	data, err = load("maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx", "maze_equip_type_group_res_v8", gMazeEquipTypeGroupResV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_group_res_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_group_res_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipTypeGroupResV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_group_res_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipTypeGroupResV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_group_res_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_group_res_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_group_res_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_group_res_v8 data success.")
	return
}
