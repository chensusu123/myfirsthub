package GMazeEquipTypeResV8Cfg

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

// MazeEquipTypeResV8ConfigRow from maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_res_v8
type MazeEquipTypeResV8ConfigRow struct {
	Order        int32  `json:"order"`        // 序号（品质*1E+部位*100W+子类型*10000+属性枚举*1000+id后3位）
	Equipment_id int32  `json:"equipment_id"` // 装备id
	Pos_sub_type int32  `json:"pos_sub_type"` // 部位子类型
	Attr_id      int32  `json:"attr_id"`      // 攻击属性id
	Name         string `json:"name"`         // 名称
	IconAtlas    string `json:"iconAtlas"`    // 图集
	Icon         string `json:"icon"`         // 图标
	Weapon_model int32  `json:"weapon_model"` // 武器模型
	Maze_model   int32  `json:"maze_model"`   // 迷宫模型资源
}

// MazeEquipTypeResV8Config from maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_res_v8
type MazeEquipTypeResV8Config struct {
	ConfigRows map[int32]*MazeEquipTypeResV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipTypeResV8Config {
	ret := &MazeEquipTypeResV8Config{ConfigRows: map[int32]*MazeEquipTypeResV8ConfigRow{}}
	return ret
}

// GetMazeEquipTypeResV8Config get one config by configId
func (c *MazeEquipTypeResV8Config) GetMazeEquipTypeResV8Config(configId int32) *MazeEquipTypeResV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipTypeResV8Config) Get(configId int32) *MazeEquipTypeResV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipTypeResV8Config get all config slice
func (c *MazeEquipTypeResV8Config) GetAllMazeEquipTypeResV8Config() (res []*MazeEquipTypeResV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipTypeResV8Config) GetAll() (res []*MazeEquipTypeResV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipTypeResV8Config

// GetMazeEquipTypeResV8Config pkg func. get one config by configId
func GetMazeEquipTypeResV8Config(configId int32) *MazeEquipTypeResV8ConfigRow {
	return gConfigData.GetMazeEquipTypeResV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipTypeResV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeEquipTypeResV8Config pkg func. get all config slice
func GetAllMazeEquipTypeResV8Config() []*MazeEquipTypeResV8ConfigRow {
	return gConfigData.GetAllMazeEquipTypeResV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipTypeResV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipTypeResV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipTypeResV8ConfigRow from maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_res_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipTypeResV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_type_res_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_type_res_v8.json",
		"maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx", "maze_equip_type_res_v8",
		&gMazeEquipTypeResV8Parser{}, &gMazeEquipTypeResV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipTypeResV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipTypeResV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipTypeResV8Config))(c)
		return true
	})
}

// RegisterMazeEquipTypeResV8InitCallBack reg config update func (old func)
var RegisterMazeEquipTypeResV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipTypeResV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipTypeResV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipTypeResV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipTypeResV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipTypeResV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipTypeResV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipTypeResV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipTypeResV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipTypeResV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipTypeResV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipTypeResV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipTypeResV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipTypeResV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipTypeResV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipTypeResV8ConfigRow", zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"),
			zap.String("sheet", "maze_equip_type_res_v8"))
		return
	}
	config, ok := container.(*MazeEquipTypeResV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipTypeResV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipTypeResV8Config", zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"),
			zap.String("sheet", "maze_equip_type_res_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipTypeResV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipTypeResV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipTypeResV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipTypeResV8Config", zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"),
			zap.String("sheet", "maze_equip_type_res_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipTypeResV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipTypeResV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipTypeResV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipTypeResV8Config", zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"),
			zap.String("sheet", "maze_equip_type_res_v8"))
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
type gMazeEquipTypeResV8Parser struct {
}

// New new config row data
func (*gMazeEquipTypeResV8Parser) New() interface{} {
	return &MazeEquipTypeResV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipTypeResV8Parser) Fields() []string {
	return gMazeEquipTypeResV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipTypeResV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipTypeResV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipTypeResV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipTypeResV8ConfigRow", zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"),
			zap.String("sheet", "maze_equip_type_res_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipTypeResV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipTypeResV8ConfigRow",
			zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"),
			zap.String("sheet", "maze_equip_type_res_v8"), zap.Int("need_count", len(gMazeEquipTypeResV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 序号（品质*1E+部位*100W+子类型*10000+属性枚举*1000+id后3位）
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 序号（品质*1E+部位*100W+子类型*10000+属性枚举*1000+id后3位） to int32 failed")
			logger.ErrorWF("parse field order 序号（品质*1E+部位*100W+子类型*10000+属性枚举*1000+id后3位） to int32 failed.",
				zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"), zap.String("sheet", "maze_equip_type_res_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 equipment_id : 装备id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field equipment_id 装备id to int32 failed")
			logger.ErrorWF("parse field equipment_id 装备id to int32 failed.",
				zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"), zap.String("sheet", "maze_equip_type_res_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Equipment_id = int32(tmp)
	}

	// parse column 2 pos_sub_type : 部位子类型
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field pos_sub_type 部位子类型 to int32 failed")
			logger.ErrorWF("parse field pos_sub_type 部位子类型 to int32 failed.",
				zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"), zap.String("sheet", "maze_equip_type_res_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Pos_sub_type = int32(tmp)
	}

	// parse column 3 attr_id : 攻击属性id
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field attr_id 攻击属性id to int32 failed")
			logger.ErrorWF("parse field attr_id 攻击属性id to int32 failed.",
				zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"), zap.String("sheet", "maze_equip_type_res_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Attr_id = int32(tmp)
	}

	// parse column 4 name : 名称
	if data[4] != "" {
		config.Name = data[4]
	}

	// parse column 5 iconAtlas : 图集
	if data[5] != "" {
		config.IconAtlas = data[5]
	}

	// parse column 6 icon : 图标
	if data[6] != "" {
		config.Icon = data[6]
	}

	// parse column 7 weapon_model : 武器模型
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field weapon_model 武器模型 to int32 failed")
			logger.ErrorWF("parse field weapon_model 武器模型 to int32 failed.",
				zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"), zap.String("sheet", "maze_equip_type_res_v8"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Weapon_model = int32(tmp)
	}

	// parse column 8 maze_model : 迷宫模型资源
	if data[8] != "" {
		tmp, err = strconv.ParseInt(data[8], 10, 64)
		if err != nil {
			err = errors.New("parse field maze_model 迷宫模型资源 to int32 failed")
			logger.ErrorWF("parse field maze_model 迷宫模型资源 to int32 failed.",
				zap.String("xlsx", "maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx"), zap.String("sheet", "maze_equip_type_res_v8"),
				zap.String("parse_data", data[8]),
				zap.Error(err))
			return
		}
		config.Maze_model = int32(tmp)
	}
	return
}

var gMazeEquipTypeResV8Fields = []string{
	"order",
	"equipment_id",
	"pos_sub_type",
	"attr_id",
	"name",
	"iconAtlas",
	"icon",
	"weapon_model",
	"maze_model",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipTypeResV8Parser{}
	loader := &gMazeEquipTypeResV8Loader{}
	var data [][]string
	data, err = load("maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx", "maze_equip_type_res_v8", gMazeEquipTypeResV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_res_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_res_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipTypeResV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_res_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipTypeResV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_res_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_res_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_res_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_type_res_v8【迷宫-装备-类型与对应资源】.xlsx maze_equip_type_res_v8 data success.")
	return
}
