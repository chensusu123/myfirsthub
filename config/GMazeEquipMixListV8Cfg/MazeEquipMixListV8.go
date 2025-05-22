package GMazeEquipMixListV8Cfg


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


// MazeEquipMixListV8ConfigRow from maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx maze_equip_mix_list_v8
type MazeEquipMixListV8ConfigRow struct {
    Order       int32  `json:"order"` // 队列id
    Equip_id       []int32  `json:"equip_id"` // 装备id
}

// MazeEquipMixListV8Config from maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx maze_equip_mix_list_v8
type MazeEquipMixListV8Config struct {
	ConfigRows map[int32]*MazeEquipMixListV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipMixListV8Config {
	ret := &MazeEquipMixListV8Config{ConfigRows: map[int32]*MazeEquipMixListV8ConfigRow{}}
	return ret
}

// GetMazeEquipMixListV8Config get one config by configId
func (c *MazeEquipMixListV8Config) GetMazeEquipMixListV8Config(configId int32) *MazeEquipMixListV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipMixListV8Config) Get(configId int32) *MazeEquipMixListV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipMixListV8Config get all config slice
func (c *MazeEquipMixListV8Config)  GetAllMazeEquipMixListV8Config () (res []*MazeEquipMixListV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipMixListV8Config)  GetAll() (res []*MazeEquipMixListV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}


// global config pointer 
var gConfigData *MazeEquipMixListV8Config 

// GetMazeEquipMixListV8Config pkg func. get one config by configId
func GetMazeEquipMixListV8Config(configId int32) *MazeEquipMixListV8ConfigRow {
	return gConfigData.GetMazeEquipMixListV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipMixListV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeEquipMixListV8Config pkg func. get all config slice
func GetAllMazeEquipMixListV8Config () []*MazeEquipMixListV8ConfigRow {
	return gConfigData.GetAllMazeEquipMixListV8Config ()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipMixListV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipMixListV8ConfigRow{
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string{
	return "MazeEquipMixListV8ConfigRow from maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx maze_equip_mix_list_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipMixListV8Config{
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_mix_list_v8"
}

func init() {
	// reg config auto load 
	config_manager.RegAutoConfig("maze_equip_mix_list_v8.json", 
		"maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx", "maze_equip_mix_list_v8",
	 	&gMazeEquipMixListV8Parser{}, &gMazeEquipMixListV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipMixListV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipMixListV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipMixListV8Config))(c)
		return true
	})
}

// RegisterMazeEquipMixListV8InitCallBack reg config update func (old func)
var RegisterMazeEquipMixListV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipMixListV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipMixListV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipMixListV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipMixListV8Config)error)(c)
		return err == nil
	})
	return 
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipMixListV8Config)error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipMixListV8Loader struct {
}
// NewContainer new data container pointer
func (*gMazeEquipMixListV8Loader) NewContainer() interface{}{
	return newConfig()
}
// Check check new config data ptr
func (*gMazeEquipMixListV8Loader) Check(newPtr interface{})error{
	// set global ptr
	cfgData := newPtr.(*MazeEquipMixListV8Config)
	return doConfigCheckCallback(cfgData)
}
// Swap swap global config data ptr
func (*gMazeEquipMixListV8Loader) Swap(newPtr interface{}){
	// convert pointer
	cache := newPtr.(*MazeEquipMixListV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}
// Add append config item to map,save result
func (*gMazeEquipMixListV8Loader) Add(logger fklog.FKLogI,container interface{}, ri interface{})(err error) {
	row,ok := ri.(*MazeEquipMixListV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipMixListV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipMixListV8ConfigRow", zap.String("xlsx", "maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx"),
			zap.String("sheet", "maze_equip_mix_list_v8"))
		return 
	}
	config,ok := container.(*MazeEquipMixListV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipMixListV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipMixListV8Config", zap.String("xlsx", "maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx"),
			zap.String("sheet", "maze_equip_mix_list_v8"))
		return 
	}
	config.ConfigRows[row.Order] = row
	return
}
// GetValue get real map value for json parse
func (*gMazeEquipMixListV8Loader) GetValue(logger fklog.FKLogI,container interface{})(real interface{},err error) {
	config,ok := container.(*MazeEquipMixListV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipMixListV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipMixListV8Config", zap.String("xlsx", "maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx"),
			zap.String("sheet", "maze_equip_mix_list_v8"))
		return 
	}
	real = &config.ConfigRows
	return
}
// Range range all data for json append data parse.
func (*gMazeEquipMixListV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error){
	config,ok := container.(*MazeEquipMixListV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipMixListV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipMixListV8Config", zap.String("xlsx", "maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx"),
			zap.String("sheet", "maze_equip_mix_list_v8"))
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
type gMazeEquipMixListV8Parser struct {
}
// New new config row data
func (*gMazeEquipMixListV8Parser) New() interface{} {
	return &MazeEquipMixListV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipMixListV8Parser) Fields() []string {
	return gMazeEquipMixListV8Fields
}
// Parse parse raw data to row data
func (*gMazeEquipMixListV8Parser) Parse(logger fklog.FKLogI,data []string, row interface{}) (err error) {
	// convert row type
	config,ok := row.(*MazeEquipMixListV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipMixListV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipMixListV8ConfigRow", zap.String("xlsx", "maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx"),
			zap.String("sheet", "maze_equip_mix_list_v8"))
		return 
	}
	// compare length
	if len(data) != len(gMazeEquipMixListV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipMixListV8ConfigRow", 
			zap.String("xlsx", "maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx"),
			zap.String("sheet", "maze_equip_mix_list_v8"), zap.Int("need_count",len(gMazeEquipMixListV8Fields)), 
			zap.Int("had_count",len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 队列id 
	if data[0] != "" {
		tmp,err = strconv.ParseInt(data[0],10,64)
		if err != nil {
			err = errors.New("parse field order 队列id to int32 failed")
			logger.ErrorWF("parse field order 队列id to int32 failed.", 
				zap.String("xlsx", "maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx"), zap.String("sheet", "maze_equip_mix_list_v8"), 
				zap.String("parse_data",data[0]), 
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 equip_id : 装备id 
	if data[1] != "" {
    
		vals := strings.Split(data[1],",")
		for k,v := range vals {
			tmp,err = strconv.ParseInt(v,10,64)
			if err != nil {
				err = errors.New("parse array field equip_id 装备id to []int32 failed")
				logger.ErrorWF("parse array field equip_id 装备id to []int32 failed.", 
					zap.String("xlsx", "maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx"), zap.String("sheet", "maze_equip_mix_list_v8"), 
					// zap.String("field_data",data[1]), 
					zap.String("parse_data", v),zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Equip_id = append(config.Equip_id, int32(tmp))
		}
	}
	return
}

var gMazeEquipMixListV8Fields = []string{
    "order",
    "equip_id",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string,error)) (err error) {
	parser := &gMazeEquipMixListV8Parser{}
	loader := &gMazeEquipMixListV8Loader{}
	var data [][]string
	data,err = load("maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx", "maze_equip_mix_list_v8", gMazeEquipMixListV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx maze_equip_mix_list_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx maze_equip_mix_list_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipMixListV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx maze_equip_mix_list_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipMixListV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx maze_equip_mix_list_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx maze_equip_mix_list_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx maze_equip_mix_list_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_mix_list_v8【迷宫-装备-合成结果队列】.xlsx maze_equip_mix_list_v8 data success.")
	return
}
