package GMazeEnergyAffixLibraryV8Cfg

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

// MazeEnergyAffixLibraryV8ConfigRow from maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx maze_energy_affix_library_v8
type MazeEnergyAffixLibraryV8ConfigRow struct {
	Order                   int32   `json:"order"`                   // 词条库id
	Certainly_affix_id_list []int32 `json:"certainly_affix_id_list"` // 必出库包括的词条id
	Affix_id_list           []int32 `json:"affix_id_list"`           // 库包括的词条id
	Weight                  int32   `json:"weight"`                  // 库随机权重
}

// MazeEnergyAffixLibraryV8Config from maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx maze_energy_affix_library_v8
type MazeEnergyAffixLibraryV8Config struct {
	ConfigRows map[int32]*MazeEnergyAffixLibraryV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEnergyAffixLibraryV8Config {
	ret := &MazeEnergyAffixLibraryV8Config{ConfigRows: map[int32]*MazeEnergyAffixLibraryV8ConfigRow{}}
	return ret
}

// GetMazeEnergyAffixLibraryV8Config get one config by configId
func (c *MazeEnergyAffixLibraryV8Config) GetMazeEnergyAffixLibraryV8Config(configId int32) *MazeEnergyAffixLibraryV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEnergyAffixLibraryV8Config) Get(configId int32) *MazeEnergyAffixLibraryV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEnergyAffixLibraryV8Config get all config slice
func (c *MazeEnergyAffixLibraryV8Config) GetAllMazeEnergyAffixLibraryV8Config() (res []*MazeEnergyAffixLibraryV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEnergyAffixLibraryV8Config) GetAll() (res []*MazeEnergyAffixLibraryV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEnergyAffixLibraryV8Config

// GetMazeEnergyAffixLibraryV8Config pkg func. get one config by configId
func GetMazeEnergyAffixLibraryV8Config(configId int32) *MazeEnergyAffixLibraryV8ConfigRow {
	return gConfigData.GetMazeEnergyAffixLibraryV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeEnergyAffixLibraryV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeEnergyAffixLibraryV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_energy_affix_library_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeEnergyAffixLibraryV8Config pkg func. get all config slice
func GetAllMazeEnergyAffixLibraryV8Config() []*MazeEnergyAffixLibraryV8ConfigRow {
	return gConfigData.GetAllMazeEnergyAffixLibraryV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEnergyAffixLibraryV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEnergyAffixLibraryV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEnergyAffixLibraryV8ConfigRow from maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx maze_energy_affix_library_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEnergyAffixLibraryV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_energy_affix_library_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_energy_affix_library_v8.json",
		"maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx", "maze_energy_affix_library_v8",
		&gMazeEnergyAffixLibraryV8Parser{}, &gMazeEnergyAffixLibraryV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEnergyAffixLibraryV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEnergyAffixLibraryV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEnergyAffixLibraryV8Config))(c)
		return true
	})
}

// RegisterMazeEnergyAffixLibraryV8InitCallBack reg config update func (old func)
var RegisterMazeEnergyAffixLibraryV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEnergyAffixLibraryV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEnergyAffixLibraryV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEnergyAffixLibraryV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEnergyAffixLibraryV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEnergyAffixLibraryV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEnergyAffixLibraryV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEnergyAffixLibraryV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEnergyAffixLibraryV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEnergyAffixLibraryV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEnergyAffixLibraryV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEnergyAffixLibraryV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEnergyAffixLibraryV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEnergyAffixLibraryV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixLibraryV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixLibraryV8ConfigRow", zap.String("xlsx", "maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx"),
			zap.String("sheet", "maze_energy_affix_library_v8"))
		return
	}
	config, ok := container.(*MazeEnergyAffixLibraryV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixLibraryV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixLibraryV8Config", zap.String("xlsx", "maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx"),
			zap.String("sheet", "maze_energy_affix_library_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEnergyAffixLibraryV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEnergyAffixLibraryV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixLibraryV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixLibraryV8Config", zap.String("xlsx", "maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx"),
			zap.String("sheet", "maze_energy_affix_library_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEnergyAffixLibraryV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEnergyAffixLibraryV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixLibraryV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixLibraryV8Config", zap.String("xlsx", "maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx"),
			zap.String("sheet", "maze_energy_affix_library_v8"))
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
type gMazeEnergyAffixLibraryV8Parser struct {
}

// New new config row data
func (*gMazeEnergyAffixLibraryV8Parser) New() interface{} {
	return &MazeEnergyAffixLibraryV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEnergyAffixLibraryV8Parser) Fields() []string {
	return gMazeEnergyAffixLibraryV8Fields
}

// Parse parse raw data to row data
func (*gMazeEnergyAffixLibraryV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEnergyAffixLibraryV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixLibraryV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixLibraryV8ConfigRow", zap.String("xlsx", "maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx"),
			zap.String("sheet", "maze_energy_affix_library_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEnergyAffixLibraryV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEnergyAffixLibraryV8ConfigRow",
			zap.String("xlsx", "maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx"),
			zap.String("sheet", "maze_energy_affix_library_v8"), zap.Int("need_count", len(gMazeEnergyAffixLibraryV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 词条库id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 词条库id to int32 failed")
			logger.ErrorWF("parse field order 词条库id to int32 failed.",
				zap.String("xlsx", "maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx"), zap.String("sheet", "maze_energy_affix_library_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 certainly_affix_id_list : 必出库包括的词条id
	if data[1] != "" {

		vals := strings.Split(data[1], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field certainly_affix_id_list 必出库包括的词条id to []int32 failed")
				logger.ErrorWF("parse array field certainly_affix_id_list 必出库包括的词条id to []int32 failed.",
					zap.String("xlsx", "maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx"), zap.String("sheet", "maze_energy_affix_library_v8"),
					// zap.String("field_data",data[1]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Certainly_affix_id_list = append(config.Certainly_affix_id_list, int32(tmp))
		}
	}

	// parse column 2 affix_id_list : 库包括的词条id
	if data[2] != "" {

		vals := strings.Split(data[2], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field affix_id_list 库包括的词条id to []int32 failed")
				logger.ErrorWF("parse array field affix_id_list 库包括的词条id to []int32 failed.",
					zap.String("xlsx", "maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx"), zap.String("sheet", "maze_energy_affix_library_v8"),
					// zap.String("field_data",data[2]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Affix_id_list = append(config.Affix_id_list, int32(tmp))
		}
	}

	// parse column 3 weight : 库随机权重
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field weight 库随机权重 to int32 failed")
			logger.ErrorWF("parse field weight 库随机权重 to int32 failed.",
				zap.String("xlsx", "maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx"), zap.String("sheet", "maze_energy_affix_library_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Weight = int32(tmp)
	}
	return
}

var gMazeEnergyAffixLibraryV8Fields = []string{
	"order",
	"certainly_affix_id_list",
	"affix_id_list",
	"weight",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEnergyAffixLibraryV8Parser{}
	loader := &gMazeEnergyAffixLibraryV8Loader{}
	var data [][]string
	data, err = load("maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx", "maze_energy_affix_library_v8", gMazeEnergyAffixLibraryV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx maze_energy_affix_library_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx maze_energy_affix_library_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEnergyAffixLibraryV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx maze_energy_affix_library_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEnergyAffixLibraryV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx maze_energy_affix_library_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx maze_energy_affix_library_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx maze_energy_affix_library_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_energy_affix_library_v8【迷宫-能量词条库权重】.xlsx maze_energy_affix_library_v8 data success.")
	return
}
