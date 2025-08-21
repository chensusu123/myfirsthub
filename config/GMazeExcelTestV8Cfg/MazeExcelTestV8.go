package GMazeExcelTestV8Cfg

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

// MazeExcelTestV8ConfigRow from maze_excel_test_v8【迷宫-导表-测试】.xlsx maze_excel_test_v8
type MazeExcelTestV8ConfigRow struct {
	Order               int32           `json:"order"`               // 等级
	Next_level_need_exp int64           `json:"next_level_need_exp"` // 升到下一级需要的经验值
	All_exp             int64           `json:"all_exp"`             // 总经验
	Attr                map[int32]int64 `json:"attr"`                // 属性id:属性值
}

// MazeExcelTestV8Config from maze_excel_test_v8【迷宫-导表-测试】.xlsx maze_excel_test_v8
type MazeExcelTestV8Config struct {
	ConfigRows map[int32]*MazeExcelTestV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeExcelTestV8Config {
	ret := &MazeExcelTestV8Config{ConfigRows: map[int32]*MazeExcelTestV8ConfigRow{}}
	return ret
}

// GetMazeExcelTestV8Config get one config by configId
func (c *MazeExcelTestV8Config) GetMazeExcelTestV8Config(configId int32) *MazeExcelTestV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeExcelTestV8Config) Get(configId int32) *MazeExcelTestV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeExcelTestV8Config get all config slice
func (c *MazeExcelTestV8Config) GetAllMazeExcelTestV8Config() (res []*MazeExcelTestV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeExcelTestV8Config) GetAll() (res []*MazeExcelTestV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeExcelTestV8Config

// GetMazeExcelTestV8Config pkg func. get one config by configId
func GetMazeExcelTestV8Config(configId int32) *MazeExcelTestV8ConfigRow {
	return gConfigData.GetMazeExcelTestV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeExcelTestV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeExcelTestV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_excel_test_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeExcelTestV8Config pkg func. get all config slice
func GetAllMazeExcelTestV8Config() []*MazeExcelTestV8ConfigRow {
	return gConfigData.GetAllMazeExcelTestV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeExcelTestV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeExcelTestV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeExcelTestV8ConfigRow from maze_excel_test_v8【迷宫-导表-测试】.xlsx maze_excel_test_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeExcelTestV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_excel_test_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_excel_test_v8.json",
		"maze_excel_test_v8【迷宫-导表-测试】.xlsx", "maze_excel_test_v8",
		&gMazeExcelTestV8Parser{}, &gMazeExcelTestV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeExcelTestV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeExcelTestV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeExcelTestV8Config))(c)
		return true
	})
}

// RegisterMazeExcelTestV8InitCallBack reg config update func (old func)
var RegisterMazeExcelTestV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeExcelTestV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeExcelTestV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeExcelTestV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeExcelTestV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeExcelTestV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeExcelTestV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeExcelTestV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeExcelTestV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeExcelTestV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeExcelTestV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeExcelTestV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeExcelTestV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeExcelTestV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeExcelTestV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeExcelTestV8ConfigRow", zap.String("xlsx", "maze_excel_test_v8【迷宫-导表-测试】.xlsx"),
			zap.String("sheet", "maze_excel_test_v8"))
		return
	}
	config, ok := container.(*MazeExcelTestV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeExcelTestV8Config")
		logger.ErrorWF("invalid type. not *MazeExcelTestV8Config", zap.String("xlsx", "maze_excel_test_v8【迷宫-导表-测试】.xlsx"),
			zap.String("sheet", "maze_excel_test_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeExcelTestV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeExcelTestV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeExcelTestV8Config")
		logger.ErrorWF("invalid type. not *MazeExcelTestV8Config", zap.String("xlsx", "maze_excel_test_v8【迷宫-导表-测试】.xlsx"),
			zap.String("sheet", "maze_excel_test_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeExcelTestV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeExcelTestV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeExcelTestV8Config")
		logger.ErrorWF("invalid type. not *MazeExcelTestV8Config", zap.String("xlsx", "maze_excel_test_v8【迷宫-导表-测试】.xlsx"),
			zap.String("sheet", "maze_excel_test_v8"))
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
type gMazeExcelTestV8Parser struct {
}

// New new config row data
func (*gMazeExcelTestV8Parser) New() interface{} {
	return &MazeExcelTestV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeExcelTestV8Parser) Fields() []string {
	return gMazeExcelTestV8Fields
}

// Parse parse raw data to row data
func (*gMazeExcelTestV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeExcelTestV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeExcelTestV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeExcelTestV8ConfigRow", zap.String("xlsx", "maze_excel_test_v8【迷宫-导表-测试】.xlsx"),
			zap.String("sheet", "maze_excel_test_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeExcelTestV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeExcelTestV8ConfigRow",
			zap.String("xlsx", "maze_excel_test_v8【迷宫-导表-测试】.xlsx"),
			zap.String("sheet", "maze_excel_test_v8"), zap.Int("need_count", len(gMazeExcelTestV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 等级
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 等级 to int32 failed")
			logger.ErrorWF("parse field order 等级 to int32 failed.",
				zap.String("xlsx", "maze_excel_test_v8【迷宫-导表-测试】.xlsx"), zap.String("sheet", "maze_excel_test_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 next_level_need_exp : 升到下一级需要的经验值
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field next_level_need_exp 升到下一级需要的经验值 to int64 failed")
			logger.ErrorWF("parse field next_level_need_exp 升到下一级需要的经验值 to int64 failed.",
				zap.String("xlsx", "maze_excel_test_v8【迷宫-导表-测试】.xlsx"), zap.String("sheet", "maze_excel_test_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Next_level_need_exp = int64(tmp)
	}

	// parse column 2 all_exp : 总经验
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field all_exp 总经验 to int64 failed")
			logger.ErrorWF("parse field all_exp 总经验 to int64 failed.",
				zap.String("xlsx", "maze_excel_test_v8【迷宫-导表-测试】.xlsx"), zap.String("sheet", "maze_excel_test_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.All_exp = int64(tmp)
	}

	// parse column 3 attr : 属性id:属性值
	if data[3] != "" {

		config.Attr = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr 属性id:属性值 to key int32 failed")
				logger.ErrorWF("parse map field attr 属性id:属性值 to key int32 failed.",
					zap.String("xlsx", "maze_excel_test_v8【迷宫-导表-测试】.xlsx"), zap.String("sheet", "maze_excel_test_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field attr 属性id:属性值 to value int64 failed")
				logger.ErrorWF("parse map field attr 属性id:属性值 to value int64 failed.",
					zap.String("xlsx", "maze_excel_test_v8【迷宫-导表-测试】.xlsx"), zap.String("sheet", "maze_excel_test_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Attr[key] = value
		}
	}
	return
}

var gMazeExcelTestV8Fields = []string{
	"order",
	"next_level_need_exp",
	"all_exp",
	"attr",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeExcelTestV8Parser{}
	loader := &gMazeExcelTestV8Loader{}
	var data [][]string
	data, err = load("maze_excel_test_v8【迷宫-导表-测试】.xlsx", "maze_excel_test_v8", gMazeExcelTestV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_excel_test_v8【迷宫-导表-测试】.xlsx maze_excel_test_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_excel_test_v8【迷宫-导表-测试】.xlsx maze_excel_test_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeExcelTestV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_excel_test_v8【迷宫-导表-测试】.xlsx maze_excel_test_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeExcelTestV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_excel_test_v8【迷宫-导表-测试】.xlsx maze_excel_test_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_excel_test_v8【迷宫-导表-测试】.xlsx maze_excel_test_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_excel_test_v8【迷宫-导表-测试】.xlsx maze_excel_test_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_excel_test_v8【迷宫-导表-测试】.xlsx maze_excel_test_v8 data success.")
	return
}
