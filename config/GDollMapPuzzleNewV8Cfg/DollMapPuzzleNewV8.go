package GDollMapPuzzleNewV8Cfg

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

// DollMapPuzzleNewV8ConfigRow from doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx doll_map_puzzle_new_v8
type DollMapPuzzleNewV8ConfigRow struct {
	Order        int32  `json:"order"`        // 序号
	Uid          int32  `json:"uid"`          // 格子id
	Show_type    int32  `json:"show_type"`    // 道具大类型（比如门）
	Sub_type     int32  `json:"sub_type"`     // 分类小类型（比如红门、黄门等）
	Level        int32  `json:"level"`        // 关卡
	Monster_area string `json:"monster_area"` // 区域id_战区
	Config_id    string `json:"config_id"`    // type*1000+index
	Stage        int32  `json:"stage"`        // 区域阶段属性
}

// DollMapPuzzleNewV8Config from doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx doll_map_puzzle_new_v8
type DollMapPuzzleNewV8Config struct {
	ConfigRows map[int32]*DollMapPuzzleNewV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *DollMapPuzzleNewV8Config {
	ret := &DollMapPuzzleNewV8Config{ConfigRows: map[int32]*DollMapPuzzleNewV8ConfigRow{}}
	return ret
}

// GetDollMapPuzzleNewV8Config get one config by configId
func (c *DollMapPuzzleNewV8Config) GetDollMapPuzzleNewV8Config(configId int32) *DollMapPuzzleNewV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *DollMapPuzzleNewV8Config) Get(configId int32) *DollMapPuzzleNewV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllDollMapPuzzleNewV8Config get all config slice
func (c *DollMapPuzzleNewV8Config) GetAllDollMapPuzzleNewV8Config() (res []*DollMapPuzzleNewV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *DollMapPuzzleNewV8Config) GetAll() (res []*DollMapPuzzleNewV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *DollMapPuzzleNewV8Config

// GetDollMapPuzzleNewV8Config pkg func. get one config by configId
func GetDollMapPuzzleNewV8Config(configId int32) *DollMapPuzzleNewV8ConfigRow {
	return gConfigData.GetDollMapPuzzleNewV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *DollMapPuzzleNewV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *DollMapPuzzleNewV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "doll_map_puzzle_new_v8", configId, otps...)
	}
	return cfg
}

// GetAllDollMapPuzzleNewV8Config pkg func. get all config slice
func GetAllDollMapPuzzleNewV8Config() []*DollMapPuzzleNewV8ConfigRow {
	return gConfigData.GetAllDollMapPuzzleNewV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*DollMapPuzzleNewV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*DollMapPuzzleNewV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "DollMapPuzzleNewV8ConfigRow from doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx doll_map_puzzle_new_v8"
}

// GetRawValue get raw data
func GetRawValue() *DollMapPuzzleNewV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "doll_map_puzzle_new_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("doll_map_puzzle_new_v8.json",
		"doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx", "doll_map_puzzle_new_v8",
		&gDollMapPuzzleNewV8Parser{}, &gDollMapPuzzleNewV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*DollMapPuzzleNewV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *DollMapPuzzleNewV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*DollMapPuzzleNewV8Config))(c)
		return true
	})
}

// RegisterDollMapPuzzleNewV8InitCallBack reg config update func (old func)
var RegisterDollMapPuzzleNewV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*DollMapPuzzleNewV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*DollMapPuzzleNewV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *DollMapPuzzleNewV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*DollMapPuzzleNewV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*DollMapPuzzleNewV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gDollMapPuzzleNewV8Loader struct {
}

// NewContainer new data container pointer
func (*gDollMapPuzzleNewV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gDollMapPuzzleNewV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*DollMapPuzzleNewV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gDollMapPuzzleNewV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*DollMapPuzzleNewV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gDollMapPuzzleNewV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*DollMapPuzzleNewV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *DollMapPuzzleNewV8ConfigRow")
		logger.ErrorWF("invalid type. not *DollMapPuzzleNewV8ConfigRow", zap.String("xlsx", "doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx"),
			zap.String("sheet", "doll_map_puzzle_new_v8"))
		return
	}
	config, ok := container.(*DollMapPuzzleNewV8Config)
	if !ok {
		err = errors.New("invalid type. not *DollMapPuzzleNewV8Config")
		logger.ErrorWF("invalid type. not *DollMapPuzzleNewV8Config", zap.String("xlsx", "doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx"),
			zap.String("sheet", "doll_map_puzzle_new_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gDollMapPuzzleNewV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*DollMapPuzzleNewV8Config)
	if !ok {
		err = errors.New("invalid type. not *DollMapPuzzleNewV8Config")
		logger.ErrorWF("invalid type. not *DollMapPuzzleNewV8Config", zap.String("xlsx", "doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx"),
			zap.String("sheet", "doll_map_puzzle_new_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gDollMapPuzzleNewV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*DollMapPuzzleNewV8Config)
	if !ok {
		err = errors.New("invalid type. not *DollMapPuzzleNewV8Config")
		logger.ErrorWF("invalid type. not *DollMapPuzzleNewV8Config", zap.String("xlsx", "doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx"),
			zap.String("sheet", "doll_map_puzzle_new_v8"))
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
type gDollMapPuzzleNewV8Parser struct {
}

// New new config row data
func (*gDollMapPuzzleNewV8Parser) New() interface{} {
	return &DollMapPuzzleNewV8ConfigRow{}
}

// Fields get config fields names
func (*gDollMapPuzzleNewV8Parser) Fields() []string {
	return gDollMapPuzzleNewV8Fields
}

// Parse parse raw data to row data
func (*gDollMapPuzzleNewV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*DollMapPuzzleNewV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *DollMapPuzzleNewV8ConfigRow")
		logger.ErrorWF("invalid type. not *DollMapPuzzleNewV8ConfigRow", zap.String("xlsx", "doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx"),
			zap.String("sheet", "doll_map_puzzle_new_v8"))
		return
	}
	// compare length
	if len(data) != len(gDollMapPuzzleNewV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*DollMapPuzzleNewV8ConfigRow",
			zap.String("xlsx", "doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx"),
			zap.String("sheet", "doll_map_puzzle_new_v8"), zap.Int("need_count", len(gDollMapPuzzleNewV8Fields)),
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
				zap.String("xlsx", "doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx"), zap.String("sheet", "doll_map_puzzle_new_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 uid : 格子id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field uid 格子id to int32 failed")
			logger.ErrorWF("parse field uid 格子id to int32 failed.",
				zap.String("xlsx", "doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx"), zap.String("sheet", "doll_map_puzzle_new_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Uid = int32(tmp)
	}

	// parse column 2 show_type : 道具大类型（比如门）
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field show_type 道具大类型（比如门） to int32 failed")
			logger.ErrorWF("parse field show_type 道具大类型（比如门） to int32 failed.",
				zap.String("xlsx", "doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx"), zap.String("sheet", "doll_map_puzzle_new_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Show_type = int32(tmp)
	}

	// parse column 3 sub_type : 分类小类型（比如红门、黄门等）
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field sub_type 分类小类型（比如红门、黄门等） to int32 failed")
			logger.ErrorWF("parse field sub_type 分类小类型（比如红门、黄门等） to int32 failed.",
				zap.String("xlsx", "doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx"), zap.String("sheet", "doll_map_puzzle_new_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Sub_type = int32(tmp)
	}

	// parse column 4 level : 关卡
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field level 关卡 to int32 failed")
			logger.ErrorWF("parse field level 关卡 to int32 failed.",
				zap.String("xlsx", "doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx"), zap.String("sheet", "doll_map_puzzle_new_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Level = int32(tmp)
	}

	// parse column 5 monster_area : 区域id_战区
	if data[5] != "" {
		config.Monster_area = data[5]
	}

	// parse column 6 config_id : type*1000+index
	if data[6] != "" {
		config.Config_id = data[6]
	}

	// parse column 7 stage : 区域阶段属性
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field stage 区域阶段属性 to int32 failed")
			logger.ErrorWF("parse field stage 区域阶段属性 to int32 failed.",
				zap.String("xlsx", "doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx"), zap.String("sheet", "doll_map_puzzle_new_v8"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Stage = int32(tmp)
	}
	return
}

var gDollMapPuzzleNewV8Fields = []string{
	"order",
	"uid",
	"show_type",
	"sub_type",
	"level",
	"monster_area",
	"config_id",
	"stage",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gDollMapPuzzleNewV8Parser{}
	loader := &gDollMapPuzzleNewV8Loader{}
	var data [][]string
	data, err = load("doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx", "doll_map_puzzle_new_v8", gDollMapPuzzleNewV8Fields)
	if err != nil {
		logger.ErrorWF("load doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx doll_map_puzzle_new_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx doll_map_puzzle_new_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gDollMapPuzzleNewV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx doll_map_puzzle_new_v8 failed.", zap.Int("row", k), zap.Strings("need", gDollMapPuzzleNewV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx doll_map_puzzle_new_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx doll_map_puzzle_new_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx doll_map_puzzle_new_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load doll_map_puzzle_new_v8【人偶-地图数据-新解谜】.xlsx doll_map_puzzle_new_v8 data success.")
	return
}
