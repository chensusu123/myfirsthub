package GMazeInitialAttrV8Cfg


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


// MazeInitialAttrV8ConfigRow from maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx maze_initial_attr_v8
type MazeInitialAttrV8ConfigRow struct {
    Order       int32  `json:"order"` // 序号
    Initial_attr       map[int32]int64  `json:"initial_attr"` // 初始化属性
}

// MazeInitialAttrV8Config from maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx maze_initial_attr_v8
type MazeInitialAttrV8Config struct {
	ConfigRows map[int32]*MazeInitialAttrV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeInitialAttrV8Config {
	ret := &MazeInitialAttrV8Config{ConfigRows: map[int32]*MazeInitialAttrV8ConfigRow{}}
	return ret
}

// GetMazeInitialAttrV8Config get one config by configId
func (c *MazeInitialAttrV8Config) GetMazeInitialAttrV8Config(configId int32) *MazeInitialAttrV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeInitialAttrV8Config) Get(configId int32) *MazeInitialAttrV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeInitialAttrV8Config get all config slice
func (c *MazeInitialAttrV8Config)  GetAllMazeInitialAttrV8Config () (res []*MazeInitialAttrV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeInitialAttrV8Config)  GetAll() (res []*MazeInitialAttrV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}


// global config pointer 
var gConfigData *MazeInitialAttrV8Config 

// GetMazeInitialAttrV8Config pkg func. get one config by configId
func GetMazeInitialAttrV8Config(configId int32) *MazeInitialAttrV8ConfigRow {
	return gConfigData.GetMazeInitialAttrV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeInitialAttrV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeInitialAttrV8Config pkg func. get all config slice
func GetAllMazeInitialAttrV8Config () []*MazeInitialAttrV8ConfigRow {
	return gConfigData.GetAllMazeInitialAttrV8Config ()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeInitialAttrV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeInitialAttrV8ConfigRow{
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string{
	return "MazeInitialAttrV8ConfigRow from maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx maze_initial_attr_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeInitialAttrV8Config{
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_initial_attr_v8"
}

func init() {
	// reg config auto load 
	config_manager.RegAutoConfig("maze_initial_attr_v8.json", 
		"maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx", "maze_initial_attr_v8",
	 	&gMazeInitialAttrV8Parser{}, &gMazeInitialAttrV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeInitialAttrV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeInitialAttrV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeInitialAttrV8Config))(c)
		return true
	})
}

// RegisterMazeInitialAttrV8InitCallBack reg config update func (old func)
var RegisterMazeInitialAttrV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeInitialAttrV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeInitialAttrV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeInitialAttrV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeInitialAttrV8Config)error)(c)
		return err == nil
	})
	return 
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeInitialAttrV8Config)error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeInitialAttrV8Loader struct {
}
// NewContainer new data container pointer
func (*gMazeInitialAttrV8Loader) NewContainer() interface{}{
	return newConfig()
}
// Check check new config data ptr
func (*gMazeInitialAttrV8Loader) Check(newPtr interface{})error{
	// set global ptr
	cfgData := newPtr.(*MazeInitialAttrV8Config)
	return doConfigCheckCallback(cfgData)
}
// Swap swap global config data ptr
func (*gMazeInitialAttrV8Loader) Swap(newPtr interface{}){
	// convert pointer
	cache := newPtr.(*MazeInitialAttrV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}
// Add append config item to map,save result
func (*gMazeInitialAttrV8Loader) Add(logger fklog.FKLogI,container interface{}, ri interface{})(err error) {
	row,ok := ri.(*MazeInitialAttrV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeInitialAttrV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeInitialAttrV8ConfigRow", zap.String("xlsx", "maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx"),
			zap.String("sheet", "maze_initial_attr_v8"))
		return 
	}
	config,ok := container.(*MazeInitialAttrV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeInitialAttrV8Config")
		logger.ErrorWF("invalid type. not *MazeInitialAttrV8Config", zap.String("xlsx", "maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx"),
			zap.String("sheet", "maze_initial_attr_v8"))
		return 
	}
	config.ConfigRows[row.Order] = row
	return
}
// GetValue get real map value for json parse
func (*gMazeInitialAttrV8Loader) GetValue(logger fklog.FKLogI,container interface{})(real interface{},err error) {
	config,ok := container.(*MazeInitialAttrV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeInitialAttrV8Config")
		logger.ErrorWF("invalid type. not *MazeInitialAttrV8Config", zap.String("xlsx", "maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx"),
			zap.String("sheet", "maze_initial_attr_v8"))
		return 
	}
	real = &config.ConfigRows
	return
}
// Range range all data for json append data parse.
func (*gMazeInitialAttrV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error){
	config,ok := container.(*MazeInitialAttrV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeInitialAttrV8Config")
		logger.ErrorWF("invalid type. not *MazeInitialAttrV8Config", zap.String("xlsx", "maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx"),
			zap.String("sheet", "maze_initial_attr_v8"))
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
type gMazeInitialAttrV8Parser struct {
}
// New new config row data
func (*gMazeInitialAttrV8Parser) New() interface{} {
	return &MazeInitialAttrV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeInitialAttrV8Parser) Fields() []string {
	return gMazeInitialAttrV8Fields
}
// Parse parse raw data to row data
func (*gMazeInitialAttrV8Parser) Parse(logger fklog.FKLogI,data []string, row interface{}) (err error) {
	// convert row type
	config,ok := row.(*MazeInitialAttrV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeInitialAttrV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeInitialAttrV8ConfigRow", zap.String("xlsx", "maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx"),
			zap.String("sheet", "maze_initial_attr_v8"))
		return 
	}
	// compare length
	if len(data) != len(gMazeInitialAttrV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeInitialAttrV8ConfigRow", 
			zap.String("xlsx", "maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx"),
			zap.String("sheet", "maze_initial_attr_v8"), zap.Int("need_count",len(gMazeInitialAttrV8Fields)), 
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
				zap.String("xlsx", "maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx"), zap.String("sheet", "maze_initial_attr_v8"), 
				zap.String("parse_data",data[0]), 
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 initial_attr : 初始化属性 
	if data[1] != "" {

		config.Initial_attr = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[1],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field initial_attr 初始化属性 to key int32 failed")
				logger.ErrorWF("parse map field initial_attr 初始化属性 to key int32 failed.", 
					zap.String("xlsx", "maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx"), zap.String("sheet", "maze_initial_attr_v8"), 
					// zap.String("field_data",data[1]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field initial_attr 初始化属性 to value int64 failed")
				logger.ErrorWF("parse map field initial_attr 初始化属性 to value int64 failed.", 
					zap.String("xlsx", "maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx"), zap.String("sheet", "maze_initial_attr_v8"), 
					// zap.String("field_data",data[1]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Initial_attr[key] = value
		}
	}
	return
}

var gMazeInitialAttrV8Fields = []string{
    "order",
    "initial_attr",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string,error)) (err error) {
	parser := &gMazeInitialAttrV8Parser{}
	loader := &gMazeInitialAttrV8Loader{}
	var data [][]string
	data,err = load("maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx", "maze_initial_attr_v8", gMazeInitialAttrV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx maze_initial_attr_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx maze_initial_attr_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeInitialAttrV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx maze_initial_attr_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeInitialAttrV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx maze_initial_attr_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx maze_initial_attr_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx maze_initial_attr_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_initial_attr_v8【迷宫-角色初始化属性】.xlsx maze_initial_attr_v8 data success.")
	return
}
