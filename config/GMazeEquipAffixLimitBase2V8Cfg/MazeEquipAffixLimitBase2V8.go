package GMazeEquipAffixLimitBase2V8Cfg

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

// MazeEquipAffixLimitBase2V8ConfigRow from maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx maze_equip_affix_limit_base2_v8
type MazeEquipAffixLimitBase2V8ConfigRow struct {
	Order           int32 `json:"order"`           // 序号
	Quality         int32 `json:"quality"`         // 品质
	Pos             int32 `json:"pos"`             // 部位
	Limit_num       int32 `json:"limit_num"`       // 初始化和洗练时，受限制的抗性词条数量
	Limit_base2_num int32 `json:"limit_base2_num"` // 初始化时，受到数量限制的垃圾词条id数量
}

// MazeEquipAffixLimitBase2V8Config from maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx maze_equip_affix_limit_base2_v8
type MazeEquipAffixLimitBase2V8Config struct {
	ConfigRows map[int32]*MazeEquipAffixLimitBase2V8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipAffixLimitBase2V8Config {
	ret := &MazeEquipAffixLimitBase2V8Config{ConfigRows: map[int32]*MazeEquipAffixLimitBase2V8ConfigRow{}}
	return ret
}

// GetMazeEquipAffixLimitBase2V8Config get one config by configId
func (c *MazeEquipAffixLimitBase2V8Config) GetMazeEquipAffixLimitBase2V8Config(configId int32) *MazeEquipAffixLimitBase2V8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipAffixLimitBase2V8Config) Get(configId int32) *MazeEquipAffixLimitBase2V8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipAffixLimitBase2V8Config get all config slice
func (c *MazeEquipAffixLimitBase2V8Config) GetAllMazeEquipAffixLimitBase2V8Config() (res []*MazeEquipAffixLimitBase2V8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipAffixLimitBase2V8Config) GetAll() (res []*MazeEquipAffixLimitBase2V8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipAffixLimitBase2V8Config

// GetMazeEquipAffixLimitBase2V8Config pkg func. get one config by configId
func GetMazeEquipAffixLimitBase2V8Config(configId int32) *MazeEquipAffixLimitBase2V8ConfigRow {
	return gConfigData.GetMazeEquipAffixLimitBase2V8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipAffixLimitBase2V8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeEquipAffixLimitBase2V8Config pkg func. get all config slice
func GetAllMazeEquipAffixLimitBase2V8Config() []*MazeEquipAffixLimitBase2V8ConfigRow {
	return gConfigData.GetAllMazeEquipAffixLimitBase2V8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipAffixLimitBase2V8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipAffixLimitBase2V8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipAffixLimitBase2V8ConfigRow from maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx maze_equip_affix_limit_base2_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipAffixLimitBase2V8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_affix_limit_base2_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_affix_limit_base2_v8.json",
		"maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx", "maze_equip_affix_limit_base2_v8",
		&gMazeEquipAffixLimitBase2V8Parser{}, &gMazeEquipAffixLimitBase2V8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipAffixLimitBase2V8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipAffixLimitBase2V8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipAffixLimitBase2V8Config))(c)
		return true
	})
}

// RegisterMazeEquipAffixLimitBase2V8InitCallBack reg config update func (old func)
var RegisterMazeEquipAffixLimitBase2V8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipAffixLimitBase2V8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipAffixLimitBase2V8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipAffixLimitBase2V8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipAffixLimitBase2V8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipAffixLimitBase2V8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipAffixLimitBase2V8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipAffixLimitBase2V8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipAffixLimitBase2V8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipAffixLimitBase2V8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipAffixLimitBase2V8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipAffixLimitBase2V8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipAffixLimitBase2V8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipAffixLimitBase2V8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixLimitBase2V8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipAffixLimitBase2V8ConfigRow", zap.String("xlsx", "maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx"),
			zap.String("sheet", "maze_equip_affix_limit_base2_v8"))
		return
	}
	config, ok := container.(*MazeEquipAffixLimitBase2V8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixLimitBase2V8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixLimitBase2V8Config", zap.String("xlsx", "maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx"),
			zap.String("sheet", "maze_equip_affix_limit_base2_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipAffixLimitBase2V8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipAffixLimitBase2V8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixLimitBase2V8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixLimitBase2V8Config", zap.String("xlsx", "maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx"),
			zap.String("sheet", "maze_equip_affix_limit_base2_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipAffixLimitBase2V8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipAffixLimitBase2V8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixLimitBase2V8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAffixLimitBase2V8Config", zap.String("xlsx", "maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx"),
			zap.String("sheet", "maze_equip_affix_limit_base2_v8"))
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
type gMazeEquipAffixLimitBase2V8Parser struct {
}

// New new config row data
func (*gMazeEquipAffixLimitBase2V8Parser) New() interface{} {
	return &MazeEquipAffixLimitBase2V8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipAffixLimitBase2V8Parser) Fields() []string {
	return gMazeEquipAffixLimitBase2V8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipAffixLimitBase2V8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipAffixLimitBase2V8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAffixLimitBase2V8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipAffixLimitBase2V8ConfigRow", zap.String("xlsx", "maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx"),
			zap.String("sheet", "maze_equip_affix_limit_base2_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipAffixLimitBase2V8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipAffixLimitBase2V8ConfigRow",
			zap.String("xlsx", "maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx"),
			zap.String("sheet", "maze_equip_affix_limit_base2_v8"), zap.Int("need_count", len(gMazeEquipAffixLimitBase2V8Fields)),
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
				zap.String("xlsx", "maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx"), zap.String("sheet", "maze_equip_affix_limit_base2_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 quality : 品质
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field quality 品质 to int32 failed")
			logger.ErrorWF("parse field quality 品质 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx"), zap.String("sheet", "maze_equip_affix_limit_base2_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Quality = int32(tmp)
	}

	// parse column 2 pos : 部位
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field pos 部位 to int32 failed")
			logger.ErrorWF("parse field pos 部位 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx"), zap.String("sheet", "maze_equip_affix_limit_base2_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Pos = int32(tmp)
	}

	// parse column 3 limit_num : 初始化和洗练时，受限制的抗性词条数量
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field limit_num 初始化和洗练时，受限制的抗性词条数量 to int32 failed")
			logger.ErrorWF("parse field limit_num 初始化和洗练时，受限制的抗性词条数量 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx"), zap.String("sheet", "maze_equip_affix_limit_base2_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Limit_num = int32(tmp)
	}

	// parse column 4 limit_base2_num : 初始化时，受到数量限制的垃圾词条id数量
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field limit_base2_num 初始化时，受到数量限制的垃圾词条id数量 to int32 failed")
			logger.ErrorWF("parse field limit_base2_num 初始化时，受到数量限制的垃圾词条id数量 to int32 failed.",
				zap.String("xlsx", "maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx"), zap.String("sheet", "maze_equip_affix_limit_base2_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Limit_base2_num = int32(tmp)
	}
	return
}

var gMazeEquipAffixLimitBase2V8Fields = []string{
	"order",
	"quality",
	"pos",
	"limit_num",
	"limit_base2_num",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipAffixLimitBase2V8Parser{}
	loader := &gMazeEquipAffixLimitBase2V8Loader{}
	var data [][]string
	data, err = load("maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx", "maze_equip_affix_limit_base2_v8", gMazeEquipAffixLimitBase2V8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx maze_equip_affix_limit_base2_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx maze_equip_affix_limit_base2_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipAffixLimitBase2V8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx maze_equip_affix_limit_base2_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipAffixLimitBase2V8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx maze_equip_affix_limit_base2_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx maze_equip_affix_limit_base2_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx maze_equip_affix_limit_base2_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_affix_limit_base2_v8【迷宫-装备-受数量限制的词条id】.xlsx maze_equip_affix_limit_base2_v8 data success.")
	return
}
