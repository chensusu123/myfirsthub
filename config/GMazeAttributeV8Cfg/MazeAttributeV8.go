package GMazeAttributeV8Cfg

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

// MazeAttributeV8ConfigRow from maze_attribute_v8【迷宫-属性】.xlsx maze_attribute_v8
type MazeAttributeV8ConfigRow struct {
	Id                   int32  `json:"id"`                   // 属性ID
	Is_into_buff         int32  `json:"is_into_buff"`         // 是否进入buff中心
	Formula_parameter_id int32  `json:"formula_parameter_id"` // 公式属性id
	Name                 string `json:"name"`                 // 属性名称
	Type                 int32  `json:"type"`                 // 属性类型
	Figure               int32  `json:"figure"`               // 数值类型
	Initial_value        int32  `json:"initial_value"`        // 初始值
}

// MazeAttributeV8Config from maze_attribute_v8【迷宫-属性】.xlsx maze_attribute_v8
type MazeAttributeV8Config struct {
	ConfigRows map[int32]*MazeAttributeV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeAttributeV8Config {
	ret := &MazeAttributeV8Config{ConfigRows: map[int32]*MazeAttributeV8ConfigRow{}}
	return ret
}

// GetMazeAttributeV8Config get one config by configId
func (c *MazeAttributeV8Config) GetMazeAttributeV8Config(configId int32) *MazeAttributeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeAttributeV8Config) Get(configId int32) *MazeAttributeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeAttributeV8Config get all config slice
func (c *MazeAttributeV8Config) GetAllMazeAttributeV8Config() (res []*MazeAttributeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeAttributeV8Config) GetAll() (res []*MazeAttributeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeAttributeV8Config

// GetMazeAttributeV8Config pkg func. get one config by configId
func GetMazeAttributeV8Config(configId int32) *MazeAttributeV8ConfigRow {
	return gConfigData.GetMazeAttributeV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeAttributeV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeAttributeV8Config pkg func. get all config slice
func GetAllMazeAttributeV8Config() []*MazeAttributeV8ConfigRow {
	return gConfigData.GetAllMazeAttributeV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeAttributeV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeAttributeV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeAttributeV8ConfigRow from maze_attribute_v8【迷宫-属性】.xlsx maze_attribute_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeAttributeV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_attribute_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_attribute_v8.json",
		"maze_attribute_v8【迷宫-属性】.xlsx", "maze_attribute_v8",
		&gMazeAttributeV8Parser{}, &gMazeAttributeV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeAttributeV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeAttributeV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeAttributeV8Config))(c)
		return true
	})
}

// RegisterMazeAttributeV8InitCallBack reg config update func (old func)
var RegisterMazeAttributeV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeAttributeV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeAttributeV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeAttributeV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeAttributeV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeAttributeV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeAttributeV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeAttributeV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeAttributeV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeAttributeV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeAttributeV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeAttributeV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeAttributeV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeAttributeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAttributeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAttributeV8ConfigRow", zap.String("xlsx", "maze_attribute_v8【迷宫-属性】.xlsx"),
			zap.String("sheet", "maze_attribute_v8"))
		return
	}
	config, ok := container.(*MazeAttributeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttributeV8Config")
		logger.ErrorWF("invalid type. not *MazeAttributeV8Config", zap.String("xlsx", "maze_attribute_v8【迷宫-属性】.xlsx"),
			zap.String("sheet", "maze_attribute_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeAttributeV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeAttributeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttributeV8Config")
		logger.ErrorWF("invalid type. not *MazeAttributeV8Config", zap.String("xlsx", "maze_attribute_v8【迷宫-属性】.xlsx"),
			zap.String("sheet", "maze_attribute_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeAttributeV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeAttributeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttributeV8Config")
		logger.ErrorWF("invalid type. not *MazeAttributeV8Config", zap.String("xlsx", "maze_attribute_v8【迷宫-属性】.xlsx"),
			zap.String("sheet", "maze_attribute_v8"))
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
type gMazeAttributeV8Parser struct {
}

// New new config row data
func (*gMazeAttributeV8Parser) New() interface{} {
	return &MazeAttributeV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeAttributeV8Parser) Fields() []string {
	return gMazeAttributeV8Fields
}

// Parse parse raw data to row data
func (*gMazeAttributeV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeAttributeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAttributeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAttributeV8ConfigRow", zap.String("xlsx", "maze_attribute_v8【迷宫-属性】.xlsx"),
			zap.String("sheet", "maze_attribute_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeAttributeV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeAttributeV8ConfigRow",
			zap.String("xlsx", "maze_attribute_v8【迷宫-属性】.xlsx"),
			zap.String("sheet", "maze_attribute_v8"), zap.Int("need_count", len(gMazeAttributeV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 属性ID
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 属性ID to int32 failed")
			logger.ErrorWF("parse field id 属性ID to int32 failed.",
				zap.String("xlsx", "maze_attribute_v8【迷宫-属性】.xlsx"), zap.String("sheet", "maze_attribute_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 is_into_buff : 是否进入buff中心
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field is_into_buff 是否进入buff中心 to int32 failed")
			logger.ErrorWF("parse field is_into_buff 是否进入buff中心 to int32 failed.",
				zap.String("xlsx", "maze_attribute_v8【迷宫-属性】.xlsx"), zap.String("sheet", "maze_attribute_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Is_into_buff = int32(tmp)
	}

	// parse column 2 formula_parameter_id : 公式属性id
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field formula_parameter_id 公式属性id to int32 failed")
			logger.ErrorWF("parse field formula_parameter_id 公式属性id to int32 failed.",
				zap.String("xlsx", "maze_attribute_v8【迷宫-属性】.xlsx"), zap.String("sheet", "maze_attribute_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Formula_parameter_id = int32(tmp)
	}

	// parse column 3 name : 属性名称
	if data[3] != "" {
		config.Name = data[3]
	}

	// parse column 4 type : 属性类型
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field type 属性类型 to int32 failed")
			logger.ErrorWF("parse field type 属性类型 to int32 failed.",
				zap.String("xlsx", "maze_attribute_v8【迷宫-属性】.xlsx"), zap.String("sheet", "maze_attribute_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Type = int32(tmp)
	}

	// parse column 5 figure : 数值类型
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field figure 数值类型 to int32 failed")
			logger.ErrorWF("parse field figure 数值类型 to int32 failed.",
				zap.String("xlsx", "maze_attribute_v8【迷宫-属性】.xlsx"), zap.String("sheet", "maze_attribute_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Figure = int32(tmp)
	}

	// parse column 6 initial_value : 初始值
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field initial_value 初始值 to int32 failed")
			logger.ErrorWF("parse field initial_value 初始值 to int32 failed.",
				zap.String("xlsx", "maze_attribute_v8【迷宫-属性】.xlsx"), zap.String("sheet", "maze_attribute_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Initial_value = int32(tmp)
	}
	return
}

var gMazeAttributeV8Fields = []string{
	"id",
	"is_into_buff",
	"formula_parameter_id",
	"name",
	"type",
	"figure",
	"initial_value",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeAttributeV8Parser{}
	loader := &gMazeAttributeV8Loader{}
	var data [][]string
	data, err = load("maze_attribute_v8【迷宫-属性】.xlsx", "maze_attribute_v8", gMazeAttributeV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_attribute_v8【迷宫-属性】.xlsx maze_attribute_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_attribute_v8【迷宫-属性】.xlsx maze_attribute_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeAttributeV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_attribute_v8【迷宫-属性】.xlsx maze_attribute_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeAttributeV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_attribute_v8【迷宫-属性】.xlsx maze_attribute_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_attribute_v8【迷宫-属性】.xlsx maze_attribute_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_attribute_v8【迷宫-属性】.xlsx maze_attribute_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_attribute_v8【迷宫-属性】.xlsx maze_attribute_v8 data success.")
	return
}
