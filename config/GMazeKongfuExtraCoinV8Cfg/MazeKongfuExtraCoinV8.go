package GMazeKongfuExtraCoinV8Cfg

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

// MazeKongfuExtraCoinV8ConfigRow from maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx maze_kongfu_extra_coin_v8
type MazeKongfuExtraCoinV8ConfigRow struct {
	Order          int32           `json:"order"`          // 序号
	Drop_type      int32           `json:"drop_type"`      // 产出类型（1-银币 2-装备分 3-经验）
	Kongfu_min     int64           `json:"kongfu_min"`     // 武力值区间，下限
	Kongfu_max     int64           `json:"kongfu_max"`     // 武力值区间，上限
	Extra_drop_num map[int32]int64 `json:"extra_drop_num"` // 冒险等级:额外掉落数量
}

// MazeKongfuExtraCoinV8Config from maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx maze_kongfu_extra_coin_v8
type MazeKongfuExtraCoinV8Config struct {
	ConfigRows map[int32]*MazeKongfuExtraCoinV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeKongfuExtraCoinV8Config {
	ret := &MazeKongfuExtraCoinV8Config{ConfigRows: map[int32]*MazeKongfuExtraCoinV8ConfigRow{}}
	return ret
}

// GetMazeKongfuExtraCoinV8Config get one config by configId
func (c *MazeKongfuExtraCoinV8Config) GetMazeKongfuExtraCoinV8Config(configId int32) *MazeKongfuExtraCoinV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeKongfuExtraCoinV8Config) Get(configId int32) *MazeKongfuExtraCoinV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeKongfuExtraCoinV8Config get all config slice
func (c *MazeKongfuExtraCoinV8Config) GetAllMazeKongfuExtraCoinV8Config() (res []*MazeKongfuExtraCoinV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeKongfuExtraCoinV8Config) GetAll() (res []*MazeKongfuExtraCoinV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeKongfuExtraCoinV8Config

// GetMazeKongfuExtraCoinV8Config pkg func. get one config by configId
func GetMazeKongfuExtraCoinV8Config(configId int32) *MazeKongfuExtraCoinV8ConfigRow {
	return gConfigData.GetMazeKongfuExtraCoinV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeKongfuExtraCoinV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeKongfuExtraCoinV8Config pkg func. get all config slice
func GetAllMazeKongfuExtraCoinV8Config() []*MazeKongfuExtraCoinV8ConfigRow {
	return gConfigData.GetAllMazeKongfuExtraCoinV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeKongfuExtraCoinV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeKongfuExtraCoinV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeKongfuExtraCoinV8ConfigRow from maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx maze_kongfu_extra_coin_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeKongfuExtraCoinV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_kongfu_extra_coin_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_kongfu_extra_coin_v8.json",
		"maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx", "maze_kongfu_extra_coin_v8",
		&gMazeKongfuExtraCoinV8Parser{}, &gMazeKongfuExtraCoinV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeKongfuExtraCoinV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeKongfuExtraCoinV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeKongfuExtraCoinV8Config))(c)
		return true
	})
}

// RegisterMazeKongfuExtraCoinV8InitCallBack reg config update func (old func)
var RegisterMazeKongfuExtraCoinV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeKongfuExtraCoinV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeKongfuExtraCoinV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeKongfuExtraCoinV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeKongfuExtraCoinV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeKongfuExtraCoinV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeKongfuExtraCoinV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeKongfuExtraCoinV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeKongfuExtraCoinV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeKongfuExtraCoinV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeKongfuExtraCoinV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeKongfuExtraCoinV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeKongfuExtraCoinV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeKongfuExtraCoinV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuExtraCoinV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeKongfuExtraCoinV8ConfigRow", zap.String("xlsx", "maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx"),
			zap.String("sheet", "maze_kongfu_extra_coin_v8"))
		return
	}
	config, ok := container.(*MazeKongfuExtraCoinV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuExtraCoinV8Config")
		logger.ErrorWF("invalid type. not *MazeKongfuExtraCoinV8Config", zap.String("xlsx", "maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx"),
			zap.String("sheet", "maze_kongfu_extra_coin_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeKongfuExtraCoinV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeKongfuExtraCoinV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuExtraCoinV8Config")
		logger.ErrorWF("invalid type. not *MazeKongfuExtraCoinV8Config", zap.String("xlsx", "maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx"),
			zap.String("sheet", "maze_kongfu_extra_coin_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeKongfuExtraCoinV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeKongfuExtraCoinV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuExtraCoinV8Config")
		logger.ErrorWF("invalid type. not *MazeKongfuExtraCoinV8Config", zap.String("xlsx", "maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx"),
			zap.String("sheet", "maze_kongfu_extra_coin_v8"))
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
type gMazeKongfuExtraCoinV8Parser struct {
}

// New new config row data
func (*gMazeKongfuExtraCoinV8Parser) New() interface{} {
	return &MazeKongfuExtraCoinV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeKongfuExtraCoinV8Parser) Fields() []string {
	return gMazeKongfuExtraCoinV8Fields
}

// Parse parse raw data to row data
func (*gMazeKongfuExtraCoinV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeKongfuExtraCoinV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeKongfuExtraCoinV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeKongfuExtraCoinV8ConfigRow", zap.String("xlsx", "maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx"),
			zap.String("sheet", "maze_kongfu_extra_coin_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeKongfuExtraCoinV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeKongfuExtraCoinV8ConfigRow",
			zap.String("xlsx", "maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx"),
			zap.String("sheet", "maze_kongfu_extra_coin_v8"), zap.Int("need_count", len(gMazeKongfuExtraCoinV8Fields)),
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
				zap.String("xlsx", "maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx"), zap.String("sheet", "maze_kongfu_extra_coin_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 drop_type : 产出类型（1-银币 2-装备分 3-经验）
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field drop_type 产出类型（1-银币 2-装备分 3-经验） to int32 failed")
			logger.ErrorWF("parse field drop_type 产出类型（1-银币 2-装备分 3-经验） to int32 failed.",
				zap.String("xlsx", "maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx"), zap.String("sheet", "maze_kongfu_extra_coin_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Drop_type = int32(tmp)
	}

	// parse column 2 kongfu_min : 武力值区间，下限
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field kongfu_min 武力值区间，下限 to int64 failed")
			logger.ErrorWF("parse field kongfu_min 武力值区间，下限 to int64 failed.",
				zap.String("xlsx", "maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx"), zap.String("sheet", "maze_kongfu_extra_coin_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Kongfu_min = int64(tmp)
	}

	// parse column 3 kongfu_max : 武力值区间，上限
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field kongfu_max 武力值区间，上限 to int64 failed")
			logger.ErrorWF("parse field kongfu_max 武力值区间，上限 to int64 failed.",
				zap.String("xlsx", "maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx"), zap.String("sheet", "maze_kongfu_extra_coin_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Kongfu_max = int64(tmp)
	}

	// parse column 4 extra_drop_num : 冒险等级:额外掉落数量
	if data[4] != "" {

		config.Extra_drop_num = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[4], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field extra_drop_num 冒险等级:额外掉落数量 to key int32 failed")
				logger.ErrorWF("parse map field extra_drop_num 冒险等级:额外掉落数量 to key int32 failed.",
					zap.String("xlsx", "maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx"), zap.String("sheet", "maze_kongfu_extra_coin_v8"),
					// zap.String("field_data",data[4]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field extra_drop_num 冒险等级:额外掉落数量 to value int64 failed")
				logger.ErrorWF("parse map field extra_drop_num 冒险等级:额外掉落数量 to value int64 failed.",
					zap.String("xlsx", "maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx"), zap.String("sheet", "maze_kongfu_extra_coin_v8"),
					// zap.String("field_data",data[4]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Extra_drop_num[key] = value
		}
	}
	return
}

var gMazeKongfuExtraCoinV8Fields = []string{
	"order",
	"drop_type",
	"kongfu_min",
	"kongfu_max",
	"extra_drop_num",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeKongfuExtraCoinV8Parser{}
	loader := &gMazeKongfuExtraCoinV8Loader{}
	var data [][]string
	data, err = load("maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx", "maze_kongfu_extra_coin_v8", gMazeKongfuExtraCoinV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx maze_kongfu_extra_coin_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx maze_kongfu_extra_coin_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeKongfuExtraCoinV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx maze_kongfu_extra_coin_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeKongfuExtraCoinV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx maze_kongfu_extra_coin_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx maze_kongfu_extra_coin_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx maze_kongfu_extra_coin_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_kongfu_extra_coin_v8【迷宫-武力值对应额外钱币】.xlsx maze_kongfu_extra_coin_v8 data success.")
	return
}
