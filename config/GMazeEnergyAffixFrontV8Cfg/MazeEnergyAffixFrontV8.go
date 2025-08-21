package GMazeEnergyAffixFrontV8Cfg

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

// MazeEnergyAffixFrontV8ConfigRow from maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx maze_energy_affix_front_v8
type MazeEnergyAffixFrontV8ConfigRow struct {
	Order                int32           `json:"order"`                // 词条前置组id
	Affix_id_set         []int32         `json:"affix_id_set"`         // 词条组id
	Must_num             int32           `json:"must_num"`             // 必须拥有的词条数
	Affix_group_num      map[int32]int32 `json:"affix_group_num"`      // 词条前置所需词条组id:总数量
	Exclusive_affix__id  int32           `json:"exclusive_affix__id"`  // 互斥词条id
	Extra_affix_group_id int32           `json:"extra_affix_group_id"` // 额外词条组id
}

// MazeEnergyAffixFrontV8Config from maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx maze_energy_affix_front_v8
type MazeEnergyAffixFrontV8Config struct {
	ConfigRows map[int32]*MazeEnergyAffixFrontV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEnergyAffixFrontV8Config {
	ret := &MazeEnergyAffixFrontV8Config{ConfigRows: map[int32]*MazeEnergyAffixFrontV8ConfigRow{}}
	return ret
}

// GetMazeEnergyAffixFrontV8Config get one config by configId
func (c *MazeEnergyAffixFrontV8Config) GetMazeEnergyAffixFrontV8Config(configId int32) *MazeEnergyAffixFrontV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEnergyAffixFrontV8Config) Get(configId int32) *MazeEnergyAffixFrontV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEnergyAffixFrontV8Config get all config slice
func (c *MazeEnergyAffixFrontV8Config) GetAllMazeEnergyAffixFrontV8Config() (res []*MazeEnergyAffixFrontV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEnergyAffixFrontV8Config) GetAll() (res []*MazeEnergyAffixFrontV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEnergyAffixFrontV8Config

// GetMazeEnergyAffixFrontV8Config pkg func. get one config by configId
func GetMazeEnergyAffixFrontV8Config(configId int32) *MazeEnergyAffixFrontV8ConfigRow {
	return gConfigData.GetMazeEnergyAffixFrontV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEnergyAffixFrontV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeEnergyAffixFrontV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_energy_affix_front_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeEnergyAffixFrontV8Config pkg func. get all config slice
func GetAllMazeEnergyAffixFrontV8Config() []*MazeEnergyAffixFrontV8ConfigRow {
	return gConfigData.GetAllMazeEnergyAffixFrontV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEnergyAffixFrontV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEnergyAffixFrontV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEnergyAffixFrontV8ConfigRow from maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx maze_energy_affix_front_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEnergyAffixFrontV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_energy_affix_front_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_energy_affix_front_v8.json",
		"maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx", "maze_energy_affix_front_v8",
		&gMazeEnergyAffixFrontV8Parser{}, &gMazeEnergyAffixFrontV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEnergyAffixFrontV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEnergyAffixFrontV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEnergyAffixFrontV8Config))(c)
		return true
	})
}

// RegisterMazeEnergyAffixFrontV8InitCallBack reg config update func (old func)
var RegisterMazeEnergyAffixFrontV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEnergyAffixFrontV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEnergyAffixFrontV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEnergyAffixFrontV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEnergyAffixFrontV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEnergyAffixFrontV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEnergyAffixFrontV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEnergyAffixFrontV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEnergyAffixFrontV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEnergyAffixFrontV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEnergyAffixFrontV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEnergyAffixFrontV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEnergyAffixFrontV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEnergyAffixFrontV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixFrontV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixFrontV8ConfigRow", zap.String("xlsx", "maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx"),
			zap.String("sheet", "maze_energy_affix_front_v8"))
		return
	}
	config, ok := container.(*MazeEnergyAffixFrontV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixFrontV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixFrontV8Config", zap.String("xlsx", "maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx"),
			zap.String("sheet", "maze_energy_affix_front_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEnergyAffixFrontV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEnergyAffixFrontV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixFrontV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixFrontV8Config", zap.String("xlsx", "maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx"),
			zap.String("sheet", "maze_energy_affix_front_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEnergyAffixFrontV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEnergyAffixFrontV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixFrontV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixFrontV8Config", zap.String("xlsx", "maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx"),
			zap.String("sheet", "maze_energy_affix_front_v8"))
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
type gMazeEnergyAffixFrontV8Parser struct {
}

// New new config row data
func (*gMazeEnergyAffixFrontV8Parser) New() interface{} {
	return &MazeEnergyAffixFrontV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEnergyAffixFrontV8Parser) Fields() []string {
	return gMazeEnergyAffixFrontV8Fields
}

// Parse parse raw data to row data
func (*gMazeEnergyAffixFrontV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEnergyAffixFrontV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixFrontV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixFrontV8ConfigRow", zap.String("xlsx", "maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx"),
			zap.String("sheet", "maze_energy_affix_front_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEnergyAffixFrontV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEnergyAffixFrontV8ConfigRow",
			zap.String("xlsx", "maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx"),
			zap.String("sheet", "maze_energy_affix_front_v8"), zap.Int("need_count", len(gMazeEnergyAffixFrontV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 词条前置组id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 词条前置组id to int32 failed")
			logger.ErrorWF("parse field order 词条前置组id to int32 failed.",
				zap.String("xlsx", "maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx"), zap.String("sheet", "maze_energy_affix_front_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 affix_id_set : 词条组id
	if data[1] != "" {

		vals := strings.Split(data[1], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field affix_id_set 词条组id to []int32 failed")
				logger.ErrorWF("parse array field affix_id_set 词条组id to []int32 failed.",
					zap.String("xlsx", "maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx"), zap.String("sheet", "maze_energy_affix_front_v8"),
					// zap.String("field_data",data[1]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Affix_id_set = append(config.Affix_id_set, int32(tmp))
		}
	}

	// parse column 2 must_num : 必须拥有的词条数
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field must_num 必须拥有的词条数 to int32 failed")
			logger.ErrorWF("parse field must_num 必须拥有的词条数 to int32 failed.",
				zap.String("xlsx", "maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx"), zap.String("sheet", "maze_energy_affix_front_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Must_num = int32(tmp)
	}

	// parse column 3 affix_group_num : 词条前置所需词条组id:总数量
	if data[3] != "" {

		config.Affix_group_num = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_group_num 词条前置所需词条组id:总数量 to key int32 failed")
				logger.ErrorWF("parse map field affix_group_num 词条前置所需词条组id:总数量 to key int32 failed.",
					zap.String("xlsx", "maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx"), zap.String("sheet", "maze_energy_affix_front_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field affix_group_num 词条前置所需词条组id:总数量 to value int32 failed")
				logger.ErrorWF("parse map field affix_group_num 词条前置所需词条组id:总数量 to value int32 failed.",
					zap.String("xlsx", "maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx"), zap.String("sheet", "maze_energy_affix_front_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Affix_group_num[key] = value
		}
	}

	// parse column 4 exclusive_affix__id : 互斥词条id
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field exclusive_affix__id 互斥词条id to int32 failed")
			logger.ErrorWF("parse field exclusive_affix__id 互斥词条id to int32 failed.",
				zap.String("xlsx", "maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx"), zap.String("sheet", "maze_energy_affix_front_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Exclusive_affix__id = int32(tmp)
	}

	// parse column 5 extra_affix_group_id : 额外词条组id
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field extra_affix_group_id 额外词条组id to int32 failed")
			logger.ErrorWF("parse field extra_affix_group_id 额外词条组id to int32 failed.",
				zap.String("xlsx", "maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx"), zap.String("sheet", "maze_energy_affix_front_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Extra_affix_group_id = int32(tmp)
	}
	return
}

var gMazeEnergyAffixFrontV8Fields = []string{
	"order",
	"affix_id_set",
	"must_num",
	"affix_group_num",
	"exclusive_affix__id",
	"extra_affix_group_id",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEnergyAffixFrontV8Parser{}
	loader := &gMazeEnergyAffixFrontV8Loader{}
	var data [][]string
	data, err = load("maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx", "maze_energy_affix_front_v8", gMazeEnergyAffixFrontV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx maze_energy_affix_front_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx maze_energy_affix_front_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEnergyAffixFrontV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx maze_energy_affix_front_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEnergyAffixFrontV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx maze_energy_affix_front_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx maze_energy_affix_front_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx maze_energy_affix_front_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_energy_affix_front_v8【迷宫-能力词条-前置词条组】.xlsx maze_energy_affix_front_v8 data success.")
	return
}
