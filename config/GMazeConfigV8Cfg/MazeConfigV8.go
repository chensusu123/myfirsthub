package GMazeConfigV8Cfg


import (
	"sync"
	"sync/atomic"
	"unsafe"
	"strconv"
	"errors"
	"strings"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)


// MazeConfigV8ConfigRow from maze_config_v8【迷宫-通用配置】.xlsx maze_config_v8
type MazeConfigV8ConfigRow struct {
    Order       int32  `json:"order"` // 序号
    Value_int       int64  `json:"value_int"` // 参数
    Value_map       map[int32]int64  `json:"value_map"` // 参数
    Value_list       []int64  `json:"value_list"` // 参数
}

// MazeConfigV8Config from maze_config_v8【迷宫-通用配置】.xlsx maze_config_v8
type MazeConfigV8Config struct {
	ConfigRows map[int32]*MazeConfigV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeConfigV8Config {
	ret := &MazeConfigV8Config{ConfigRows: map[int32]*MazeConfigV8ConfigRow{}}
	return ret
}

// GetMazeConfigV8Config get one config by configId
func (c *MazeConfigV8Config) GetMazeConfigV8Config(configId int32) *MazeConfigV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeConfigV8Config) Get(configId int32) *MazeConfigV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeConfigV8Config get all config slice
func (c *MazeConfigV8Config)  GetAllMazeConfigV8Config () (res []*MazeConfigV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeConfigV8Config)  GetAll() (res []*MazeConfigV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}


// global config pointer 
var gConfigData *MazeConfigV8Config 

// GetMazeConfigV8Config pkg func. get one config by configId
func GetMazeConfigV8Config(configId int32) *MazeConfigV8ConfigRow {
	return gConfigData.GetMazeConfigV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeConfigV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeConfigV8Config pkg func. get all config slice
func GetAllMazeConfigV8Config () []*MazeConfigV8ConfigRow {
	return gConfigData.GetAllMazeConfigV8Config ()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeConfigV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeConfigV8ConfigRow{
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string{
	return "MazeConfigV8ConfigRow from maze_config_v8【迷宫-通用配置】.xlsx maze_config_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeConfigV8Config{
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_config_v8"
}

func init() {
	// reg config auto load 
	config_manager.RegAutoConfig("maze_config_v8.json", 
		"maze_config_v8【迷宫-通用配置】.xlsx", "maze_config_v8",
	 	&gMazeConfigV8Parser{}, &gMazeConfigV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeConfigV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeConfigV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeConfigV8Config))(c)
		return true
	})
}

// RegisterMazeConfigV8InitCallBack reg config update func (old func)
var RegisterMazeConfigV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeConfigV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeConfigV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeConfigV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeConfigV8Config)error)(c)
		return err == nil
	})
	return 
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeConfigV8Config)error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeConfigV8Loader struct {
}
// NewContainer new data container pointer
func (*gMazeConfigV8Loader) NewContainer() interface{}{
	return newConfig()
}
// Check check new config data ptr
func (*gMazeConfigV8Loader) Check(newPtr interface{})error{
	// set global ptr
	cfgData := newPtr.(*MazeConfigV8Config)
	return doConfigCheckCallback(cfgData)
}
// Swap swap global config data ptr
func (*gMazeConfigV8Loader) Swap(newPtr interface{}){
	// convert pointer
	cache := newPtr.(*MazeConfigV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}
// Add append config item to map,save result
func (*gMazeConfigV8Loader) Add(logger fklog.FKLogI,container interface{}, ri interface{})(err error) {
	row,ok := ri.(*MazeConfigV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeConfigV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeConfigV8ConfigRow", zap.String("xlsx", "maze_config_v8【迷宫-通用配置】.xlsx"),
			zap.String("sheet", "maze_config_v8"))
		return 
	}
	config,ok := container.(*MazeConfigV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeConfigV8Config")
		logger.ErrorWF("invalid type. not *MazeConfigV8Config", zap.String("xlsx", "maze_config_v8【迷宫-通用配置】.xlsx"),
			zap.String("sheet", "maze_config_v8"))
		return 
	}
	config.ConfigRows[row.Order] = row
	return
}
// GetValue get real map value for json parse
func (*gMazeConfigV8Loader) GetValue(logger fklog.FKLogI,container interface{})(real interface{},err error) {
	config,ok := container.(*MazeConfigV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeConfigV8Config")
		logger.ErrorWF("invalid type. not *MazeConfigV8Config", zap.String("xlsx", "maze_config_v8【迷宫-通用配置】.xlsx"),
			zap.String("sheet", "maze_config_v8"))
		return 
	}
	real = &config.ConfigRows
	return
}
// Range range all data for json append data parse.
func (*gMazeConfigV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error){
	config,ok := container.(*MazeConfigV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeConfigV8Config")
		logger.ErrorWF("invalid type. not *MazeConfigV8Config", zap.String("xlsx", "maze_config_v8【迷宫-通用配置】.xlsx"),
			zap.String("sheet", "maze_config_v8"))
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
type gMazeConfigV8Parser struct {
}
// New new config row data
func (*gMazeConfigV8Parser) New() interface{} {
	return &MazeConfigV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeConfigV8Parser) Fields() []string {
	return gMazeConfigV8Fields
}
// Parse parse raw data to row data
func (*gMazeConfigV8Parser) Parse(logger fklog.FKLogI,data []string, row interface{}) (err error) {
	// convert row type
	config,ok := row.(*MazeConfigV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeConfigV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeConfigV8ConfigRow", zap.String("xlsx", "maze_config_v8【迷宫-通用配置】.xlsx"),
			zap.String("sheet", "maze_config_v8"))
		return 
	}
	// compare length
	if len(data) != len(gMazeConfigV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeConfigV8ConfigRow", 
			zap.String("xlsx", "maze_config_v8【迷宫-通用配置】.xlsx"),
			zap.String("sheet", "maze_config_v8"), zap.Int("need_count",len(gMazeConfigV8Fields)), 
			zap.Int("had_count",len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 序号 
	if data[0] != "" {
		tmp,err = strconv.ParseInt(data[0],10,64)
		if err != nil {
			err = errors.New("parse field order 序号 to int32 failed")
			logger.ErrorWF("parse field order 序号 to int32 failed.", 
				zap.String("xlsx", "maze_config_v8【迷宫-通用配置】.xlsx"), zap.String("sheet", "maze_config_v8"), 
				zap.String("parse_data",data[0]), 
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 value_int : 参数 
	if data[1] != "" {
		tmp,err = strconv.ParseInt(data[1],10,64)
		if err != nil {
			err = errors.New("parse field value_int 参数 to int64 failed")
			logger.ErrorWF("parse field value_int 参数 to int64 failed.", 
				zap.String("xlsx", "maze_config_v8【迷宫-通用配置】.xlsx"), zap.String("sheet", "maze_config_v8"), 
				zap.String("parse_data",data[1]), 
				zap.Error(err))
			return
		}
		config.Value_int = int64(tmp)
	}

	// parse column 2 value_map : 参数 
	if data[2] != "" {

		config.Value_map = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[2],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field value_map 参数 to key int32 failed")
				logger.ErrorWF("parse map field value_map 参数 to key int32 failed.", 
					zap.String("xlsx", "maze_config_v8【迷宫-通用配置】.xlsx"), zap.String("sheet", "maze_config_v8"), 
					// zap.String("field_data",data[2]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field value_map 参数 to value int64 failed")
				logger.ErrorWF("parse map field value_map 参数 to value int64 failed.", 
					zap.String("xlsx", "maze_config_v8【迷宫-通用配置】.xlsx"), zap.String("sheet", "maze_config_v8"), 
					// zap.String("field_data",data[2]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Value_map[key] = value
		}
	}

	// parse column 3 value_list : 参数 
	if data[3] != "" {
    
		vals := strings.Split(data[3],",")
		for k,v := range vals {
			tmp,err = strconv.ParseInt(v,10,64)
			if err != nil {
				err = errors.New("parse array field value_list 参数 to []int64 failed")
				logger.ErrorWF("parse array field value_list 参数 to []int64 failed.", 
					zap.String("xlsx", "maze_config_v8【迷宫-通用配置】.xlsx"), zap.String("sheet", "maze_config_v8"), 
					// zap.String("field_data",data[3]), 
					zap.String("parse_data", v),zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Value_list = append(config.Value_list, int64(tmp))
		}
	}
	return
}

var gMazeConfigV8Fields = []string{
    "order",
    "value_int",
    "value_map",
    "value_list",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string,error)) (err error) {
	parser := &gMazeConfigV8Parser{}
	loader := &gMazeConfigV8Loader{}
	var data [][]string
	data,err = load("maze_config_v8【迷宫-通用配置】.xlsx", "maze_config_v8", gMazeConfigV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_config_v8【迷宫-通用配置】.xlsx maze_config_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_config_v8【迷宫-通用配置】.xlsx maze_config_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeConfigV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_config_v8【迷宫-通用配置】.xlsx maze_config_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeConfigV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_config_v8【迷宫-通用配置】.xlsx maze_config_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_config_v8【迷宫-通用配置】.xlsx maze_config_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_config_v8【迷宫-通用配置】.xlsx maze_config_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_config_v8【迷宫-通用配置】.xlsx maze_config_v8 data success.")
	return
}
