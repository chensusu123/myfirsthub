package GMazePatrolV8Cfg


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


// MazePatrolV8ConfigRow from maze_patrol_v8【迷宫-巡逻守卫】.xlsx maze_patrol_v8
type MazePatrolV8ConfigRow struct {
    Order       int32  `json:"order"` // 序号
    Level_id       int32  `json:"level_id"` // 关卡id
    Monster_id       []int32  `json:"monster_id"` // 怪物id
    Brush_time       int32  `json:"brush_time"` // 刷新cd
}

// MazePatrolV8Config from maze_patrol_v8【迷宫-巡逻守卫】.xlsx maze_patrol_v8
type MazePatrolV8Config struct {
	ConfigRows map[int32]*MazePatrolV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazePatrolV8Config {
	ret := &MazePatrolV8Config{ConfigRows: map[int32]*MazePatrolV8ConfigRow{}}
	return ret
}

// GetMazePatrolV8Config get one config by configId
func (c *MazePatrolV8Config) GetMazePatrolV8Config(configId int32) *MazePatrolV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazePatrolV8Config) Get(configId int32) *MazePatrolV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazePatrolV8Config get all config slice
func (c *MazePatrolV8Config)  GetAllMazePatrolV8Config () (res []*MazePatrolV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazePatrolV8Config)  GetAll() (res []*MazePatrolV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}


// global config pointer 
var gConfigData *MazePatrolV8Config 

// GetMazePatrolV8Config pkg func. get one config by configId
func GetMazePatrolV8Config(configId int32) *MazePatrolV8ConfigRow {
	return gConfigData.GetMazePatrolV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazePatrolV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazePatrolV8Config pkg func. get all config slice
func GetAllMazePatrolV8Config () []*MazePatrolV8ConfigRow {
	return gConfigData.GetAllMazePatrolV8Config ()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazePatrolV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazePatrolV8ConfigRow{
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string{
	return "MazePatrolV8ConfigRow from maze_patrol_v8【迷宫-巡逻守卫】.xlsx maze_patrol_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazePatrolV8Config{
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_patrol_v8"
}

func init() {
	// reg config auto load 
	config_manager.RegAutoConfig("maze_patrol_v8.json", 
		"maze_patrol_v8【迷宫-巡逻守卫】.xlsx", "maze_patrol_v8",
	 	&gMazePatrolV8Parser{}, &gMazePatrolV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazePatrolV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazePatrolV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazePatrolV8Config))(c)
		return true
	})
}

// RegisterMazePatrolV8InitCallBack reg config update func (old func)
var RegisterMazePatrolV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazePatrolV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazePatrolV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazePatrolV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazePatrolV8Config)error)(c)
		return err == nil
	})
	return 
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazePatrolV8Config)error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazePatrolV8Loader struct {
}
// NewContainer new data container pointer
func (*gMazePatrolV8Loader) NewContainer() interface{}{
	return newConfig()
}
// Check check new config data ptr
func (*gMazePatrolV8Loader) Check(newPtr interface{})error{
	// set global ptr
	cfgData := newPtr.(*MazePatrolV8Config)
	return doConfigCheckCallback(cfgData)
}
// Swap swap global config data ptr
func (*gMazePatrolV8Loader) Swap(newPtr interface{}){
	// convert pointer
	cache := newPtr.(*MazePatrolV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}
// Add append config item to map,save result
func (*gMazePatrolV8Loader) Add(logger fklog.FKLogI,container interface{}, ri interface{})(err error) {
	row,ok := ri.(*MazePatrolV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazePatrolV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazePatrolV8ConfigRow", zap.String("xlsx", "maze_patrol_v8【迷宫-巡逻守卫】.xlsx"),
			zap.String("sheet", "maze_patrol_v8"))
		return 
	}
	config,ok := container.(*MazePatrolV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazePatrolV8Config")
		logger.ErrorWF("invalid type. not *MazePatrolV8Config", zap.String("xlsx", "maze_patrol_v8【迷宫-巡逻守卫】.xlsx"),
			zap.String("sheet", "maze_patrol_v8"))
		return 
	}
	config.ConfigRows[row.Order] = row
	return
}
// GetValue get real map value for json parse
func (*gMazePatrolV8Loader) GetValue(logger fklog.FKLogI,container interface{})(real interface{},err error) {
	config,ok := container.(*MazePatrolV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazePatrolV8Config")
		logger.ErrorWF("invalid type. not *MazePatrolV8Config", zap.String("xlsx", "maze_patrol_v8【迷宫-巡逻守卫】.xlsx"),
			zap.String("sheet", "maze_patrol_v8"))
		return 
	}
	real = &config.ConfigRows
	return
}
// Range range all data for json append data parse.
func (*gMazePatrolV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error){
	config,ok := container.(*MazePatrolV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazePatrolV8Config")
		logger.ErrorWF("invalid type. not *MazePatrolV8Config", zap.String("xlsx", "maze_patrol_v8【迷宫-巡逻守卫】.xlsx"),
			zap.String("sheet", "maze_patrol_v8"))
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
type gMazePatrolV8Parser struct {
}
// New new config row data
func (*gMazePatrolV8Parser) New() interface{} {
	return &MazePatrolV8ConfigRow{}
}

// Fields get config fields names
func (*gMazePatrolV8Parser) Fields() []string {
	return gMazePatrolV8Fields
}
// Parse parse raw data to row data
func (*gMazePatrolV8Parser) Parse(logger fklog.FKLogI,data []string, row interface{}) (err error) {
	// convert row type
	config,ok := row.(*MazePatrolV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazePatrolV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazePatrolV8ConfigRow", zap.String("xlsx", "maze_patrol_v8【迷宫-巡逻守卫】.xlsx"),
			zap.String("sheet", "maze_patrol_v8"))
		return 
	}
	// compare length
	if len(data) != len(gMazePatrolV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazePatrolV8ConfigRow", 
			zap.String("xlsx", "maze_patrol_v8【迷宫-巡逻守卫】.xlsx"),
			zap.String("sheet", "maze_patrol_v8"), zap.Int("need_count",len(gMazePatrolV8Fields)), 
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
				zap.String("xlsx", "maze_patrol_v8【迷宫-巡逻守卫】.xlsx"), zap.String("sheet", "maze_patrol_v8"), 
				zap.String("parse_data",data[0]), 
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 level_id : 关卡id 
	if data[1] != "" {
		tmp,err = strconv.ParseInt(data[1],10,64)
		if err != nil {
			err = errors.New("parse field level_id 关卡id to int32 failed")
			logger.ErrorWF("parse field level_id 关卡id to int32 failed.", 
				zap.String("xlsx", "maze_patrol_v8【迷宫-巡逻守卫】.xlsx"), zap.String("sheet", "maze_patrol_v8"), 
				zap.String("parse_data",data[1]), 
				zap.Error(err))
			return
		}
		config.Level_id = int32(tmp)
	}

	// parse column 2 monster_id : 怪物id 
	if data[2] != "" {
    
		vals := strings.Split(data[2],",")
		for k,v := range vals {
			tmp,err = strconv.ParseInt(v,10,64)
			if err != nil {
				err = errors.New("parse array field monster_id 怪物id to []int32 failed")
				logger.ErrorWF("parse array field monster_id 怪物id to []int32 failed.", 
					zap.String("xlsx", "maze_patrol_v8【迷宫-巡逻守卫】.xlsx"), zap.String("sheet", "maze_patrol_v8"), 
					// zap.String("field_data",data[2]), 
					zap.String("parse_data", v),zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Monster_id = append(config.Monster_id, int32(tmp))
		}
	}

	// parse column 3 brush_time : 刷新cd 
	if data[3] != "" {
		tmp,err = strconv.ParseInt(data[3],10,64)
		if err != nil {
			err = errors.New("parse field brush_time 刷新cd to int32 failed")
			logger.ErrorWF("parse field brush_time 刷新cd to int32 failed.", 
				zap.String("xlsx", "maze_patrol_v8【迷宫-巡逻守卫】.xlsx"), zap.String("sheet", "maze_patrol_v8"), 
				zap.String("parse_data",data[3]), 
				zap.Error(err))
			return
		}
		config.Brush_time = int32(tmp)
	}
	return
}

var gMazePatrolV8Fields = []string{
    "order",
    "level_id",
    "monster_id",
    "brush_time",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string,error)) (err error) {
	parser := &gMazePatrolV8Parser{}
	loader := &gMazePatrolV8Loader{}
	var data [][]string
	data,err = load("maze_patrol_v8【迷宫-巡逻守卫】.xlsx", "maze_patrol_v8", gMazePatrolV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_patrol_v8【迷宫-巡逻守卫】.xlsx maze_patrol_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_patrol_v8【迷宫-巡逻守卫】.xlsx maze_patrol_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazePatrolV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_patrol_v8【迷宫-巡逻守卫】.xlsx maze_patrol_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazePatrolV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_patrol_v8【迷宫-巡逻守卫】.xlsx maze_patrol_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_patrol_v8【迷宫-巡逻守卫】.xlsx maze_patrol_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_patrol_v8【迷宫-巡逻守卫】.xlsx maze_patrol_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_patrol_v8【迷宫-巡逻守卫】.xlsx maze_patrol_v8 data success.")
	return
}
