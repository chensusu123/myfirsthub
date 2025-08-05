package GMazeEquipPosQuaLvV8Cfg

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

// MazeEquipPosQuaLvV8ConfigRow from maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx maze_equip_pos_qua_lv_v8
type MazeEquipPosQuaLvV8ConfigRow struct {
	Order   int32 `json:"order"`   // ID
	Level   int32 `json:"level"`   // 等级
	Quality int32 `json:"quality"` // 品质
	Pos1    int32 `json:"pos1"`    // 部位1装备id武器
	Pos2    int32 `json:"pos2"`    // 部位2装备id帽子
	Pos3    int32 `json:"pos3"`    // 部位3装备id衣服
	Pos4    int32 `json:"pos4"`    // 部位4装备id裤子
	Pos5    int32 `json:"pos5"`    // 部位5装备id戒指
	Pos6    int32 `json:"pos6"`    // 部位6装备id鞋子
	Pos7    int32 `json:"pos7"`    // 部位7装备id护腕
	Pos8    int32 `json:"pos8"`    // 部位8装备id项链
}

// MazeEquipPosQuaLvV8Config from maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx maze_equip_pos_qua_lv_v8
type MazeEquipPosQuaLvV8Config struct {
	ConfigRows map[int32]*MazeEquipPosQuaLvV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipPosQuaLvV8Config {
	ret := &MazeEquipPosQuaLvV8Config{ConfigRows: map[int32]*MazeEquipPosQuaLvV8ConfigRow{}}
	return ret
}

// GetMazeEquipPosQuaLvV8Config get one config by configId
func (c *MazeEquipPosQuaLvV8Config) GetMazeEquipPosQuaLvV8Config(configId int32) *MazeEquipPosQuaLvV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipPosQuaLvV8Config) Get(configId int32) *MazeEquipPosQuaLvV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipPosQuaLvV8Config get all config slice
func (c *MazeEquipPosQuaLvV8Config) GetAllMazeEquipPosQuaLvV8Config() (res []*MazeEquipPosQuaLvV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipPosQuaLvV8Config) GetAll() (res []*MazeEquipPosQuaLvV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipPosQuaLvV8Config

// GetMazeEquipPosQuaLvV8Config pkg func. get one config by configId
func GetMazeEquipPosQuaLvV8Config(configId int32) *MazeEquipPosQuaLvV8ConfigRow {
	return gConfigData.GetMazeEquipPosQuaLvV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipPosQuaLvV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeEquipPosQuaLvV8Config pkg func. get all config slice
func GetAllMazeEquipPosQuaLvV8Config() []*MazeEquipPosQuaLvV8ConfigRow {
	return gConfigData.GetAllMazeEquipPosQuaLvV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipPosQuaLvV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipPosQuaLvV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipPosQuaLvV8ConfigRow from maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx maze_equip_pos_qua_lv_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipPosQuaLvV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_pos_qua_lv_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_pos_qua_lv_v8.json",
		"maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx", "maze_equip_pos_qua_lv_v8",
		&gMazeEquipPosQuaLvV8Parser{}, &gMazeEquipPosQuaLvV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipPosQuaLvV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipPosQuaLvV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipPosQuaLvV8Config))(c)
		return true
	})
}

// RegisterMazeEquipPosQuaLvV8InitCallBack reg config update func (old func)
var RegisterMazeEquipPosQuaLvV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipPosQuaLvV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipPosQuaLvV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipPosQuaLvV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipPosQuaLvV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipPosQuaLvV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipPosQuaLvV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipPosQuaLvV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipPosQuaLvV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipPosQuaLvV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipPosQuaLvV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipPosQuaLvV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipPosQuaLvV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipPosQuaLvV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosQuaLvV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipPosQuaLvV8ConfigRow", zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"),
			zap.String("sheet", "maze_equip_pos_qua_lv_v8"))
		return
	}
	config, ok := container.(*MazeEquipPosQuaLvV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosQuaLvV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipPosQuaLvV8Config", zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"),
			zap.String("sheet", "maze_equip_pos_qua_lv_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipPosQuaLvV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipPosQuaLvV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosQuaLvV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipPosQuaLvV8Config", zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"),
			zap.String("sheet", "maze_equip_pos_qua_lv_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipPosQuaLvV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipPosQuaLvV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosQuaLvV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipPosQuaLvV8Config", zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"),
			zap.String("sheet", "maze_equip_pos_qua_lv_v8"))
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
type gMazeEquipPosQuaLvV8Parser struct {
}

// New new config row data
func (*gMazeEquipPosQuaLvV8Parser) New() interface{} {
	return &MazeEquipPosQuaLvV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipPosQuaLvV8Parser) Fields() []string {
	return gMazeEquipPosQuaLvV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipPosQuaLvV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipPosQuaLvV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosQuaLvV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipPosQuaLvV8ConfigRow", zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"),
			zap.String("sheet", "maze_equip_pos_qua_lv_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipPosQuaLvV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipPosQuaLvV8ConfigRow",
			zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"),
			zap.String("sheet", "maze_equip_pos_qua_lv_v8"), zap.Int("need_count", len(gMazeEquipPosQuaLvV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : ID
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order ID to int32 failed")
			logger.ErrorWF("parse field order ID to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"), zap.String("sheet", "maze_equip_pos_qua_lv_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 level : 等级
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field level 等级 to int32 failed")
			logger.ErrorWF("parse field level 等级 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"), zap.String("sheet", "maze_equip_pos_qua_lv_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Level = int32(tmp)
	}

	// parse column 2 quality : 品质
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field quality 品质 to int32 failed")
			logger.ErrorWF("parse field quality 品质 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"), zap.String("sheet", "maze_equip_pos_qua_lv_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Quality = int32(tmp)
	}

	// parse column 3 pos1 : 部位1装备id武器
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field pos1 部位1装备id武器 to int32 failed")
			logger.ErrorWF("parse field pos1 部位1装备id武器 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"), zap.String("sheet", "maze_equip_pos_qua_lv_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Pos1 = int32(tmp)
	}

	// parse column 4 pos2 : 部位2装备id帽子
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field pos2 部位2装备id帽子 to int32 failed")
			logger.ErrorWF("parse field pos2 部位2装备id帽子 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"), zap.String("sheet", "maze_equip_pos_qua_lv_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Pos2 = int32(tmp)
	}

	// parse column 5 pos3 : 部位3装备id衣服
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field pos3 部位3装备id衣服 to int32 failed")
			logger.ErrorWF("parse field pos3 部位3装备id衣服 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"), zap.String("sheet", "maze_equip_pos_qua_lv_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Pos3 = int32(tmp)
	}

	// parse column 6 pos4 : 部位4装备id裤子
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field pos4 部位4装备id裤子 to int32 failed")
			logger.ErrorWF("parse field pos4 部位4装备id裤子 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"), zap.String("sheet", "maze_equip_pos_qua_lv_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Pos4 = int32(tmp)
	}

	// parse column 7 pos5 : 部位5装备id戒指
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field pos5 部位5装备id戒指 to int32 failed")
			logger.ErrorWF("parse field pos5 部位5装备id戒指 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"), zap.String("sheet", "maze_equip_pos_qua_lv_v8"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Pos5 = int32(tmp)
	}

	// parse column 8 pos6 : 部位6装备id鞋子
	if data[8] != "" {
		tmp, err = strconv.ParseInt(data[8], 10, 64)
		if err != nil {
			err = errors.New("parse field pos6 部位6装备id鞋子 to int32 failed")
			logger.ErrorWF("parse field pos6 部位6装备id鞋子 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"), zap.String("sheet", "maze_equip_pos_qua_lv_v8"),
				zap.String("parse_data", data[8]),
				zap.Error(err))
			return
		}
		config.Pos6 = int32(tmp)
	}

	// parse column 9 pos7 : 部位7装备id护腕
	if data[9] != "" {
		tmp, err = strconv.ParseInt(data[9], 10, 64)
		if err != nil {
			err = errors.New("parse field pos7 部位7装备id护腕 to int32 failed")
			logger.ErrorWF("parse field pos7 部位7装备id护腕 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"), zap.String("sheet", "maze_equip_pos_qua_lv_v8"),
				zap.String("parse_data", data[9]),
				zap.Error(err))
			return
		}
		config.Pos7 = int32(tmp)
	}

	// parse column 10 pos8 : 部位8装备id项链
	if data[10] != "" {
		tmp, err = strconv.ParseInt(data[10], 10, 64)
		if err != nil {
			err = errors.New("parse field pos8 部位8装备id项链 to int32 failed")
			logger.ErrorWF("parse field pos8 部位8装备id项链 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx"), zap.String("sheet", "maze_equip_pos_qua_lv_v8"),
				zap.String("parse_data", data[10]),
				zap.Error(err))
			return
		}
		config.Pos8 = int32(tmp)
	}
	return
}

var gMazeEquipPosQuaLvV8Fields = []string{
	"order",
	"level",
	"quality",
	"pos1",
	"pos2",
	"pos3",
	"pos4",
	"pos5",
	"pos6",
	"pos7",
	"pos8",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipPosQuaLvV8Parser{}
	loader := &gMazeEquipPosQuaLvV8Loader{}
	var data [][]string
	data, err = load("maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx", "maze_equip_pos_qua_lv_v8", gMazeEquipPosQuaLvV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx maze_equip_pos_qua_lv_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx maze_equip_pos_qua_lv_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipPosQuaLvV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx maze_equip_pos_qua_lv_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipPosQuaLvV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx maze_equip_pos_qua_lv_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx maze_equip_pos_qua_lv_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx maze_equip_pos_qua_lv_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_pos_qua_lv_v8【迷宫-装备-部位品质等级对应表】.xlsx maze_equip_pos_qua_lv_v8 data success.")
	return
}
