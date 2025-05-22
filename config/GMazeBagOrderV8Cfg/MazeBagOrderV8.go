package GMazeBagOrderV8Cfg


import (
	"sync"
	"sync/atomic"
	"unsafe"
	"strconv"
	"errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)


// MazeBagOrderV8ConfigRow from maze_bag_order_v8【迷宫-背包-物品排序】.xlsx maze_bag_order_v8
type MazeBagOrderV8ConfigRow struct {
    Item_id       int32  `json:"item_id"` // 物品id
    Quality       int32  `json:"quality"` // 品质
    Order_type_2       int32  `json:"order_type_2"` // 同品质排序（数字大的在前面）
}

// MazeBagOrderV8Config from maze_bag_order_v8【迷宫-背包-物品排序】.xlsx maze_bag_order_v8
type MazeBagOrderV8Config struct {
	ConfigRows map[int32]*MazeBagOrderV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeBagOrderV8Config {
	ret := &MazeBagOrderV8Config{ConfigRows: map[int32]*MazeBagOrderV8ConfigRow{}}
	return ret
}

// GetMazeBagOrderV8Config get one config by configId
func (c *MazeBagOrderV8Config) GetMazeBagOrderV8Config(configId int32) *MazeBagOrderV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeBagOrderV8Config) Get(configId int32) *MazeBagOrderV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeBagOrderV8Config get all config slice
func (c *MazeBagOrderV8Config)  GetAllMazeBagOrderV8Config () (res []*MazeBagOrderV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeBagOrderV8Config)  GetAll() (res []*MazeBagOrderV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}


// global config pointer 
var gConfigData *MazeBagOrderV8Config 

// GetMazeBagOrderV8Config pkg func. get one config by configId
func GetMazeBagOrderV8Config(configId int32) *MazeBagOrderV8ConfigRow {
	return gConfigData.GetMazeBagOrderV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeBagOrderV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeBagOrderV8Config pkg func. get all config slice
func GetAllMazeBagOrderV8Config () []*MazeBagOrderV8ConfigRow {
	return gConfigData.GetAllMazeBagOrderV8Config ()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeBagOrderV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeBagOrderV8ConfigRow{
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string{
	return "MazeBagOrderV8ConfigRow from maze_bag_order_v8【迷宫-背包-物品排序】.xlsx maze_bag_order_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeBagOrderV8Config{
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_bag_order_v8"
}

func init() {
	// reg config auto load 
	config_manager.RegAutoConfig("maze_bag_order_v8.json", 
		"maze_bag_order_v8【迷宫-背包-物品排序】.xlsx", "maze_bag_order_v8",
	 	&gMazeBagOrderV8Parser{}, &gMazeBagOrderV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeBagOrderV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeBagOrderV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeBagOrderV8Config))(c)
		return true
	})
}

// RegisterMazeBagOrderV8InitCallBack reg config update func (old func)
var RegisterMazeBagOrderV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeBagOrderV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeBagOrderV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeBagOrderV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeBagOrderV8Config)error)(c)
		return err == nil
	})
	return 
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeBagOrderV8Config)error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeBagOrderV8Loader struct {
}
// NewContainer new data container pointer
func (*gMazeBagOrderV8Loader) NewContainer() interface{}{
	return newConfig()
}
// Check check new config data ptr
func (*gMazeBagOrderV8Loader) Check(newPtr interface{})error{
	// set global ptr
	cfgData := newPtr.(*MazeBagOrderV8Config)
	return doConfigCheckCallback(cfgData)
}
// Swap swap global config data ptr
func (*gMazeBagOrderV8Loader) Swap(newPtr interface{}){
	// convert pointer
	cache := newPtr.(*MazeBagOrderV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}
// Add append config item to map,save result
func (*gMazeBagOrderV8Loader) Add(logger fklog.FKLogI,container interface{}, ri interface{})(err error) {
	row,ok := ri.(*MazeBagOrderV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBagOrderV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBagOrderV8ConfigRow", zap.String("xlsx", "maze_bag_order_v8【迷宫-背包-物品排序】.xlsx"),
			zap.String("sheet", "maze_bag_order_v8"))
		return 
	}
	config,ok := container.(*MazeBagOrderV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBagOrderV8Config")
		logger.ErrorWF("invalid type. not *MazeBagOrderV8Config", zap.String("xlsx", "maze_bag_order_v8【迷宫-背包-物品排序】.xlsx"),
			zap.String("sheet", "maze_bag_order_v8"))
		return 
	}
	config.ConfigRows[row.Item_id] = row
	return
}
// GetValue get real map value for json parse
func (*gMazeBagOrderV8Loader) GetValue(logger fklog.FKLogI,container interface{})(real interface{},err error) {
	config,ok := container.(*MazeBagOrderV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBagOrderV8Config")
		logger.ErrorWF("invalid type. not *MazeBagOrderV8Config", zap.String("xlsx", "maze_bag_order_v8【迷宫-背包-物品排序】.xlsx"),
			zap.String("sheet", "maze_bag_order_v8"))
		return 
	}
	real = &config.ConfigRows
	return
}
// Range range all data for json append data parse.
func (*gMazeBagOrderV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error){
	config,ok := container.(*MazeBagOrderV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBagOrderV8Config")
		logger.ErrorWF("invalid type. not *MazeBagOrderV8Config", zap.String("xlsx", "maze_bag_order_v8【迷宫-背包-物品排序】.xlsx"),
			zap.String("sheet", "maze_bag_order_v8"))
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
type gMazeBagOrderV8Parser struct {
}
// New new config row data
func (*gMazeBagOrderV8Parser) New() interface{} {
	return &MazeBagOrderV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeBagOrderV8Parser) Fields() []string {
	return gMazeBagOrderV8Fields
}
// Parse parse raw data to row data
func (*gMazeBagOrderV8Parser) Parse(logger fklog.FKLogI,data []string, row interface{}) (err error) {
	// convert row type
	config,ok := row.(*MazeBagOrderV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBagOrderV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBagOrderV8ConfigRow", zap.String("xlsx", "maze_bag_order_v8【迷宫-背包-物品排序】.xlsx"),
			zap.String("sheet", "maze_bag_order_v8"))
		return 
	}
	// compare length
	if len(data) != len(gMazeBagOrderV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeBagOrderV8ConfigRow", 
			zap.String("xlsx", "maze_bag_order_v8【迷宫-背包-物品排序】.xlsx"),
			zap.String("sheet", "maze_bag_order_v8"), zap.Int("need_count",len(gMazeBagOrderV8Fields)), 
			zap.Int("had_count",len(data)))
		return
	}

	var tmp int64

	// parse column 0 item_id : 物品id 
	if data[0] != "" {
		tmp,err = strconv.ParseInt(data[0],10,64)
		if err != nil {
			err = errors.New("parse field item_id 物品id to int32 failed")
			logger.ErrorWF("parse field item_id 物品id to int32 failed.", 
				zap.String("xlsx", "maze_bag_order_v8【迷宫-背包-物品排序】.xlsx"), zap.String("sheet", "maze_bag_order_v8"), 
				zap.String("parse_data",data[0]), 
				zap.Error(err))
			return
		}
		config.Item_id = int32(tmp)
	}

	// parse column 1 quality : 品质 
	if data[1] != "" {
		tmp,err = strconv.ParseInt(data[1],10,64)
		if err != nil {
			err = errors.New("parse field quality 品质 to int32 failed")
			logger.ErrorWF("parse field quality 品质 to int32 failed.", 
				zap.String("xlsx", "maze_bag_order_v8【迷宫-背包-物品排序】.xlsx"), zap.String("sheet", "maze_bag_order_v8"), 
				zap.String("parse_data",data[1]), 
				zap.Error(err))
			return
		}
		config.Quality = int32(tmp)
	}

	// parse column 2 order_type_2 : 同品质排序（数字大的在前面） 
	if data[2] != "" {
		tmp,err = strconv.ParseInt(data[2],10,64)
		if err != nil {
			err = errors.New("parse field order_type_2 同品质排序（数字大的在前面） to int32 failed")
			logger.ErrorWF("parse field order_type_2 同品质排序（数字大的在前面） to int32 failed.", 
				zap.String("xlsx", "maze_bag_order_v8【迷宫-背包-物品排序】.xlsx"), zap.String("sheet", "maze_bag_order_v8"), 
				zap.String("parse_data",data[2]), 
				zap.Error(err))
			return
		}
		config.Order_type_2 = int32(tmp)
	}
	return
}

var gMazeBagOrderV8Fields = []string{
    "item_id",
    "quality",
    "order_type_2",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string,error)) (err error) {
	parser := &gMazeBagOrderV8Parser{}
	loader := &gMazeBagOrderV8Loader{}
	var data [][]string
	data,err = load("maze_bag_order_v8【迷宫-背包-物品排序】.xlsx", "maze_bag_order_v8", gMazeBagOrderV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_bag_order_v8【迷宫-背包-物品排序】.xlsx maze_bag_order_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_bag_order_v8【迷宫-背包-物品排序】.xlsx maze_bag_order_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeBagOrderV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_bag_order_v8【迷宫-背包-物品排序】.xlsx maze_bag_order_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeBagOrderV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_bag_order_v8【迷宫-背包-物品排序】.xlsx maze_bag_order_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_bag_order_v8【迷宫-背包-物品排序】.xlsx maze_bag_order_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_bag_order_v8【迷宫-背包-物品排序】.xlsx maze_bag_order_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_bag_order_v8【迷宫-背包-物品排序】.xlsx maze_bag_order_v8 data success.")
	return
}
