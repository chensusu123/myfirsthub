package GMazeEnergyResetCostV8Cfg


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


// MazeEnergyResetCostV8ConfigRow from maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx maze_energy_reset_cost_v8
type MazeEnergyResetCostV8ConfigRow struct {
    Order       int32  `json:"order"` // 序号
    Mun_min       int32  `json:"mun_min"` // 刷新次数小
    Mun_max       int32  `json:"mun_max"` // 刷新次数大
    Cost       map[int32]int64  `json:"cost"` // 刷新消耗
}

// MazeEnergyResetCostV8Config from maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx maze_energy_reset_cost_v8
type MazeEnergyResetCostV8Config struct {
	ConfigRows map[int32]*MazeEnergyResetCostV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEnergyResetCostV8Config {
	ret := &MazeEnergyResetCostV8Config{ConfigRows: map[int32]*MazeEnergyResetCostV8ConfigRow{}}
	return ret
}

// GetMazeEnergyResetCostV8Config get one config by configId
func (c *MazeEnergyResetCostV8Config) GetMazeEnergyResetCostV8Config(configId int32) *MazeEnergyResetCostV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEnergyResetCostV8Config) Get(configId int32) *MazeEnergyResetCostV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEnergyResetCostV8Config get all config slice
func (c *MazeEnergyResetCostV8Config)  GetAllMazeEnergyResetCostV8Config () (res []*MazeEnergyResetCostV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEnergyResetCostV8Config)  GetAll() (res []*MazeEnergyResetCostV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}


// global config pointer 
var gConfigData *MazeEnergyResetCostV8Config 

// GetMazeEnergyResetCostV8Config pkg func. get one config by configId
func GetMazeEnergyResetCostV8Config(configId int32) *MazeEnergyResetCostV8ConfigRow {
	return gConfigData.GetMazeEnergyResetCostV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEnergyResetCostV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeEnergyResetCostV8Config pkg func. get all config slice
func GetAllMazeEnergyResetCostV8Config () []*MazeEnergyResetCostV8ConfigRow {
	return gConfigData.GetAllMazeEnergyResetCostV8Config ()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEnergyResetCostV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEnergyResetCostV8ConfigRow{
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string{
	return "MazeEnergyResetCostV8ConfigRow from maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx maze_energy_reset_cost_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEnergyResetCostV8Config{
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_energy_reset_cost_v8"
}

func init() {
	// reg config auto load 
	config_manager.RegAutoConfig("maze_energy_reset_cost_v8.json", 
		"maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx", "maze_energy_reset_cost_v8",
	 	&gMazeEnergyResetCostV8Parser{}, &gMazeEnergyResetCostV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEnergyResetCostV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEnergyResetCostV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEnergyResetCostV8Config))(c)
		return true
	})
}

// RegisterMazeEnergyResetCostV8InitCallBack reg config update func (old func)
var RegisterMazeEnergyResetCostV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEnergyResetCostV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEnergyResetCostV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEnergyResetCostV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEnergyResetCostV8Config)error)(c)
		return err == nil
	})
	return 
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEnergyResetCostV8Config)error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEnergyResetCostV8Loader struct {
}
// NewContainer new data container pointer
func (*gMazeEnergyResetCostV8Loader) NewContainer() interface{}{
	return newConfig()
}
// Check check new config data ptr
func (*gMazeEnergyResetCostV8Loader) Check(newPtr interface{})error{
	// set global ptr
	cfgData := newPtr.(*MazeEnergyResetCostV8Config)
	return doConfigCheckCallback(cfgData)
}
// Swap swap global config data ptr
func (*gMazeEnergyResetCostV8Loader) Swap(newPtr interface{}){
	// convert pointer
	cache := newPtr.(*MazeEnergyResetCostV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}
// Add append config item to map,save result
func (*gMazeEnergyResetCostV8Loader) Add(logger fklog.FKLogI,container interface{}, ri interface{})(err error) {
	row,ok := ri.(*MazeEnergyResetCostV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyResetCostV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEnergyResetCostV8ConfigRow", zap.String("xlsx", "maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx"),
			zap.String("sheet", "maze_energy_reset_cost_v8"))
		return 
	}
	config,ok := container.(*MazeEnergyResetCostV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyResetCostV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyResetCostV8Config", zap.String("xlsx", "maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx"),
			zap.String("sheet", "maze_energy_reset_cost_v8"))
		return 
	}
	config.ConfigRows[row.Order] = row
	return
}
// GetValue get real map value for json parse
func (*gMazeEnergyResetCostV8Loader) GetValue(logger fklog.FKLogI,container interface{})(real interface{},err error) {
	config,ok := container.(*MazeEnergyResetCostV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyResetCostV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyResetCostV8Config", zap.String("xlsx", "maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx"),
			zap.String("sheet", "maze_energy_reset_cost_v8"))
		return 
	}
	real = &config.ConfigRows
	return
}
// Range range all data for json append data parse.
func (*gMazeEnergyResetCostV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error){
	config,ok := container.(*MazeEnergyResetCostV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyResetCostV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyResetCostV8Config", zap.String("xlsx", "maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx"),
			zap.String("sheet", "maze_energy_reset_cost_v8"))
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
type gMazeEnergyResetCostV8Parser struct {
}
// New new config row data
func (*gMazeEnergyResetCostV8Parser) New() interface{} {
	return &MazeEnergyResetCostV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEnergyResetCostV8Parser) Fields() []string {
	return gMazeEnergyResetCostV8Fields
}
// Parse parse raw data to row data
func (*gMazeEnergyResetCostV8Parser) Parse(logger fklog.FKLogI,data []string, row interface{}) (err error) {
	// convert row type
	config,ok := row.(*MazeEnergyResetCostV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyResetCostV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEnergyResetCostV8ConfigRow", zap.String("xlsx", "maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx"),
			zap.String("sheet", "maze_energy_reset_cost_v8"))
		return 
	}
	// compare length
	if len(data) != len(gMazeEnergyResetCostV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEnergyResetCostV8ConfigRow", 
			zap.String("xlsx", "maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx"),
			zap.String("sheet", "maze_energy_reset_cost_v8"), zap.Int("need_count",len(gMazeEnergyResetCostV8Fields)), 
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
				zap.String("xlsx", "maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx"), zap.String("sheet", "maze_energy_reset_cost_v8"), 
				zap.String("parse_data",data[0]), 
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 mun_min : 刷新次数小 
	if data[1] != "" {
		tmp,err = strconv.ParseInt(data[1],10,64)
		if err != nil {
			err = errors.New("parse field mun_min 刷新次数小 to int32 failed")
			logger.ErrorWF("parse field mun_min 刷新次数小 to int32 failed.", 
				zap.String("xlsx", "maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx"), zap.String("sheet", "maze_energy_reset_cost_v8"), 
				zap.String("parse_data",data[1]), 
				zap.Error(err))
			return
		}
		config.Mun_min = int32(tmp)
	}

	// parse column 2 mun_max : 刷新次数大 
	if data[2] != "" {
		tmp,err = strconv.ParseInt(data[2],10,64)
		if err != nil {
			err = errors.New("parse field mun_max 刷新次数大 to int32 failed")
			logger.ErrorWF("parse field mun_max 刷新次数大 to int32 failed.", 
				zap.String("xlsx", "maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx"), zap.String("sheet", "maze_energy_reset_cost_v8"), 
				zap.String("parse_data",data[2]), 
				zap.Error(err))
			return
		}
		config.Mun_max = int32(tmp)
	}

	// parse column 3 cost : 刷新消耗 
	if data[3] != "" {

		config.Cost = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[3],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field cost 刷新消耗 to key int32 failed")
				logger.ErrorWF("parse map field cost 刷新消耗 to key int32 failed.", 
					zap.String("xlsx", "maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx"), zap.String("sheet", "maze_energy_reset_cost_v8"), 
					// zap.String("field_data",data[3]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field cost 刷新消耗 to value int64 failed")
				logger.ErrorWF("parse map field cost 刷新消耗 to value int64 failed.", 
					zap.String("xlsx", "maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx"), zap.String("sheet", "maze_energy_reset_cost_v8"), 
					// zap.String("field_data",data[3]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Cost[key] = value
		}
	}
	return
}

var gMazeEnergyResetCostV8Fields = []string{
    "order",
    "mun_min",
    "mun_max",
    "cost",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string,error)) (err error) {
	parser := &gMazeEnergyResetCostV8Parser{}
	loader := &gMazeEnergyResetCostV8Loader{}
	var data [][]string
	data,err = load("maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx", "maze_energy_reset_cost_v8", gMazeEnergyResetCostV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx maze_energy_reset_cost_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx maze_energy_reset_cost_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEnergyResetCostV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx maze_energy_reset_cost_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEnergyResetCostV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx maze_energy_reset_cost_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx maze_energy_reset_cost_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx maze_energy_reset_cost_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_energy_reset_cost_v8【迷宫-能力-刷新次数与消耗】.xlsx maze_energy_reset_cost_v8 data success.")
	return
}
