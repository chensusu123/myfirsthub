package GMazeActionCountV8Cfg


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


// MazeActionCountV8ConfigRow from maze_action_count_v8【行为次数】.xlsx maze_action_count_v8
type MazeActionCountV8ConfigRow struct {
    Id       int32  `json:"id"` // 唯一id
    Action_name       string  `json:"action_name"` // 行为名称
    Day_count_v8       int32  `json:"day_count_v8"` // 每天次数
    Count_level_v8       map[int32]int32  `json:"count_level_v8"` // 需要迷宫冒险等级
    Lose_desc_v8       string  `json:"lose_desc_v8"` // 未开启提示文本
}

// MazeActionCountV8Config from maze_action_count_v8【行为次数】.xlsx maze_action_count_v8
type MazeActionCountV8Config struct {
	ConfigRows map[int32]*MazeActionCountV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeActionCountV8Config {
	ret := &MazeActionCountV8Config{ConfigRows: map[int32]*MazeActionCountV8ConfigRow{}}
	return ret
}

// GetMazeActionCountV8Config get one config by configId
func (c *MazeActionCountV8Config) GetMazeActionCountV8Config(configId int32) *MazeActionCountV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeActionCountV8Config) Get(configId int32) *MazeActionCountV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeActionCountV8Config get all config slice
func (c *MazeActionCountV8Config)  GetAllMazeActionCountV8Config () (res []*MazeActionCountV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeActionCountV8Config)  GetAll() (res []*MazeActionCountV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}


// global config pointer 
var gConfigData *MazeActionCountV8Config 

// GetMazeActionCountV8Config pkg func. get one config by configId
func GetMazeActionCountV8Config(configId int32) *MazeActionCountV8ConfigRow {
	return gConfigData.GetMazeActionCountV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeActionCountV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeActionCountV8Config pkg func. get all config slice
func GetAllMazeActionCountV8Config () []*MazeActionCountV8ConfigRow {
	return gConfigData.GetAllMazeActionCountV8Config ()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeActionCountV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeActionCountV8ConfigRow{
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string{
	return "MazeActionCountV8ConfigRow from maze_action_count_v8【行为次数】.xlsx maze_action_count_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeActionCountV8Config{
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_action_count_v8"
}

func init() {
	// reg config auto load 
	config_manager.RegAutoConfig("maze_action_count_v8.json", 
		"maze_action_count_v8【行为次数】.xlsx", "maze_action_count_v8",
	 	&gMazeActionCountV8Parser{}, &gMazeActionCountV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeActionCountV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeActionCountV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeActionCountV8Config))(c)
		return true
	})
}

// RegisterMazeActionCountV8InitCallBack reg config update func (old func)
var RegisterMazeActionCountV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeActionCountV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeActionCountV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeActionCountV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeActionCountV8Config)error)(c)
		return err == nil
	})
	return 
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeActionCountV8Config)error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeActionCountV8Loader struct {
}
// NewContainer new data container pointer
func (*gMazeActionCountV8Loader) NewContainer() interface{}{
	return newConfig()
}
// Check check new config data ptr
func (*gMazeActionCountV8Loader) Check(newPtr interface{})error{
	// set global ptr
	cfgData := newPtr.(*MazeActionCountV8Config)
	return doConfigCheckCallback(cfgData)
}
// Swap swap global config data ptr
func (*gMazeActionCountV8Loader) Swap(newPtr interface{}){
	// convert pointer
	cache := newPtr.(*MazeActionCountV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}
// Add append config item to map,save result
func (*gMazeActionCountV8Loader) Add(logger fklog.FKLogI,container interface{}, ri interface{})(err error) {
	row,ok := ri.(*MazeActionCountV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeActionCountV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeActionCountV8ConfigRow", zap.String("xlsx", "maze_action_count_v8【行为次数】.xlsx"),
			zap.String("sheet", "maze_action_count_v8"))
		return 
	}
	config,ok := container.(*MazeActionCountV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeActionCountV8Config")
		logger.ErrorWF("invalid type. not *MazeActionCountV8Config", zap.String("xlsx", "maze_action_count_v8【行为次数】.xlsx"),
			zap.String("sheet", "maze_action_count_v8"))
		return 
	}
	config.ConfigRows[row.Id] = row
	return
}
// GetValue get real map value for json parse
func (*gMazeActionCountV8Loader) GetValue(logger fklog.FKLogI,container interface{})(real interface{},err error) {
	config,ok := container.(*MazeActionCountV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeActionCountV8Config")
		logger.ErrorWF("invalid type. not *MazeActionCountV8Config", zap.String("xlsx", "maze_action_count_v8【行为次数】.xlsx"),
			zap.String("sheet", "maze_action_count_v8"))
		return 
	}
	real = &config.ConfigRows
	return
}
// Range range all data for json append data parse.
func (*gMazeActionCountV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error){
	config,ok := container.(*MazeActionCountV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeActionCountV8Config")
		logger.ErrorWF("invalid type. not *MazeActionCountV8Config", zap.String("xlsx", "maze_action_count_v8【行为次数】.xlsx"),
			zap.String("sheet", "maze_action_count_v8"))
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
type gMazeActionCountV8Parser struct {
}
// New new config row data
func (*gMazeActionCountV8Parser) New() interface{} {
	return &MazeActionCountV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeActionCountV8Parser) Fields() []string {
	return gMazeActionCountV8Fields
}
// Parse parse raw data to row data
func (*gMazeActionCountV8Parser) Parse(logger fklog.FKLogI,data []string, row interface{}) (err error) {
	// convert row type
	config,ok := row.(*MazeActionCountV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeActionCountV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeActionCountV8ConfigRow", zap.String("xlsx", "maze_action_count_v8【行为次数】.xlsx"),
			zap.String("sheet", "maze_action_count_v8"))
		return 
	}
	// compare length
	if len(data) != len(gMazeActionCountV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeActionCountV8ConfigRow", 
			zap.String("xlsx", "maze_action_count_v8【行为次数】.xlsx"),
			zap.String("sheet", "maze_action_count_v8"), zap.Int("need_count",len(gMazeActionCountV8Fields)), 
			zap.Int("had_count",len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 唯一id 
	if data[0] != "" {
		tmp,err = strconv.ParseInt(data[0],10,64)
		if err != nil {
			err = errors.New("parse field id 唯一id to int32 failed")
			logger.ErrorWF("parse field id 唯一id to int32 failed.", 
				zap.String("xlsx", "maze_action_count_v8【行为次数】.xlsx"), zap.String("sheet", "maze_action_count_v8"), 
				zap.String("parse_data",data[0]), 
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 action_name : 行为名称 
	if data[1] != "" {
		config.Action_name = data[1]
	}

	// parse column 2 day_count_v8 : 每天次数 
	if data[2] != "" {
		tmp,err = strconv.ParseInt(data[2],10,64)
		if err != nil {
			err = errors.New("parse field day_count_v8 每天次数 to int32 failed")
			logger.ErrorWF("parse field day_count_v8 每天次数 to int32 failed.", 
				zap.String("xlsx", "maze_action_count_v8【行为次数】.xlsx"), zap.String("sheet", "maze_action_count_v8"), 
				zap.String("parse_data",data[2]), 
				zap.Error(err))
			return
		}
		config.Day_count_v8 = int32(tmp)
	}

	// parse column 3 count_level_v8 : 需要迷宫冒险等级 
	if data[3] != "" {

		config.Count_level_v8 = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[3],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field count_level_v8 需要迷宫冒险等级 to key int32 failed")
				logger.ErrorWF("parse map field count_level_v8 需要迷宫冒险等级 to key int32 failed.", 
					zap.String("xlsx", "maze_action_count_v8【行为次数】.xlsx"), zap.String("sheet", "maze_action_count_v8"), 
					// zap.String("field_data",data[3]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field count_level_v8 需要迷宫冒险等级 to value int32 failed")
				logger.ErrorWF("parse map field count_level_v8 需要迷宫冒险等级 to value int32 failed.", 
					zap.String("xlsx", "maze_action_count_v8【行为次数】.xlsx"), zap.String("sheet", "maze_action_count_v8"), 
					// zap.String("field_data",data[3]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Count_level_v8[key] = value
		}
	}

	// parse column 4 lose_desc_v8 : 未开启提示文本 
	if data[4] != "" {
		config.Lose_desc_v8 = data[4]
	}
	return
}

var gMazeActionCountV8Fields = []string{
    "id",
    "action_name",
    "day_count_v8",
    "count_level_v8",
    "lose_desc_v8",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string,error)) (err error) {
	parser := &gMazeActionCountV8Parser{}
	loader := &gMazeActionCountV8Loader{}
	var data [][]string
	data,err = load("maze_action_count_v8【行为次数】.xlsx", "maze_action_count_v8", gMazeActionCountV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_action_count_v8【行为次数】.xlsx maze_action_count_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_action_count_v8【行为次数】.xlsx maze_action_count_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeActionCountV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_action_count_v8【行为次数】.xlsx maze_action_count_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeActionCountV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_action_count_v8【行为次数】.xlsx maze_action_count_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_action_count_v8【行为次数】.xlsx maze_action_count_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_action_count_v8【行为次数】.xlsx maze_action_count_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_action_count_v8【行为次数】.xlsx maze_action_count_v8 data success.")
	return
}
