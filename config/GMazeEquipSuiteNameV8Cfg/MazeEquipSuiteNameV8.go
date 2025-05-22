package GMazeEquipSuiteNameV8Cfg


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


// MazeEquipSuiteNameV8ConfigRow from maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx maze_equip_suite_name_v8
type MazeEquipSuiteNameV8ConfigRow struct {
    Equipment_id       int32  `json:"equipment_id"` // 装备id
    Suite_equip_name       map[int32]string  `json:"suite_equip_name"` // 套装装备名称
}

// MazeEquipSuiteNameV8Config from maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx maze_equip_suite_name_v8
type MazeEquipSuiteNameV8Config struct {
	ConfigRows map[int32]*MazeEquipSuiteNameV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipSuiteNameV8Config {
	ret := &MazeEquipSuiteNameV8Config{ConfigRows: map[int32]*MazeEquipSuiteNameV8ConfigRow{}}
	return ret
}

// GetMazeEquipSuiteNameV8Config get one config by configId
func (c *MazeEquipSuiteNameV8Config) GetMazeEquipSuiteNameV8Config(configId int32) *MazeEquipSuiteNameV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipSuiteNameV8Config) Get(configId int32) *MazeEquipSuiteNameV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipSuiteNameV8Config get all config slice
func (c *MazeEquipSuiteNameV8Config)  GetAllMazeEquipSuiteNameV8Config () (res []*MazeEquipSuiteNameV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipSuiteNameV8Config)  GetAll() (res []*MazeEquipSuiteNameV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}


// global config pointer 
var gConfigData *MazeEquipSuiteNameV8Config 

// GetMazeEquipSuiteNameV8Config pkg func. get one config by configId
func GetMazeEquipSuiteNameV8Config(configId int32) *MazeEquipSuiteNameV8ConfigRow {
	return gConfigData.GetMazeEquipSuiteNameV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipSuiteNameV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeEquipSuiteNameV8Config pkg func. get all config slice
func GetAllMazeEquipSuiteNameV8Config () []*MazeEquipSuiteNameV8ConfigRow {
	return gConfigData.GetAllMazeEquipSuiteNameV8Config ()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipSuiteNameV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipSuiteNameV8ConfigRow{
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string{
	return "MazeEquipSuiteNameV8ConfigRow from maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx maze_equip_suite_name_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipSuiteNameV8Config{
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_suite_name_v8"
}

func init() {
	// reg config auto load 
	config_manager.RegAutoConfig("maze_equip_suite_name_v8.json", 
		"maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx", "maze_equip_suite_name_v8",
	 	&gMazeEquipSuiteNameV8Parser{}, &gMazeEquipSuiteNameV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipSuiteNameV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipSuiteNameV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipSuiteNameV8Config))(c)
		return true
	})
}

// RegisterMazeEquipSuiteNameV8InitCallBack reg config update func (old func)
var RegisterMazeEquipSuiteNameV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipSuiteNameV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipSuiteNameV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipSuiteNameV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipSuiteNameV8Config)error)(c)
		return err == nil
	})
	return 
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipSuiteNameV8Config)error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipSuiteNameV8Loader struct {
}
// NewContainer new data container pointer
func (*gMazeEquipSuiteNameV8Loader) NewContainer() interface{}{
	return newConfig()
}
// Check check new config data ptr
func (*gMazeEquipSuiteNameV8Loader) Check(newPtr interface{})error{
	// set global ptr
	cfgData := newPtr.(*MazeEquipSuiteNameV8Config)
	return doConfigCheckCallback(cfgData)
}
// Swap swap global config data ptr
func (*gMazeEquipSuiteNameV8Loader) Swap(newPtr interface{}){
	// convert pointer
	cache := newPtr.(*MazeEquipSuiteNameV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}
// Add append config item to map,save result
func (*gMazeEquipSuiteNameV8Loader) Add(logger fklog.FKLogI,container interface{}, ri interface{})(err error) {
	row,ok := ri.(*MazeEquipSuiteNameV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteNameV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteNameV8ConfigRow", zap.String("xlsx", "maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx"),
			zap.String("sheet", "maze_equip_suite_name_v8"))
		return 
	}
	config,ok := container.(*MazeEquipSuiteNameV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteNameV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteNameV8Config", zap.String("xlsx", "maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx"),
			zap.String("sheet", "maze_equip_suite_name_v8"))
		return 
	}
	config.ConfigRows[row.Equipment_id] = row
	return
}
// GetValue get real map value for json parse
func (*gMazeEquipSuiteNameV8Loader) GetValue(logger fklog.FKLogI,container interface{})(real interface{},err error) {
	config,ok := container.(*MazeEquipSuiteNameV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteNameV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteNameV8Config", zap.String("xlsx", "maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx"),
			zap.String("sheet", "maze_equip_suite_name_v8"))
		return 
	}
	real = &config.ConfigRows
	return
}
// Range range all data for json append data parse.
func (*gMazeEquipSuiteNameV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error){
	config,ok := container.(*MazeEquipSuiteNameV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteNameV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteNameV8Config", zap.String("xlsx", "maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx"),
			zap.String("sheet", "maze_equip_suite_name_v8"))
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
type gMazeEquipSuiteNameV8Parser struct {
}
// New new config row data
func (*gMazeEquipSuiteNameV8Parser) New() interface{} {
	return &MazeEquipSuiteNameV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipSuiteNameV8Parser) Fields() []string {
	return gMazeEquipSuiteNameV8Fields
}
// Parse parse raw data to row data
func (*gMazeEquipSuiteNameV8Parser) Parse(logger fklog.FKLogI,data []string, row interface{}) (err error) {
	// convert row type
	config,ok := row.(*MazeEquipSuiteNameV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteNameV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteNameV8ConfigRow", zap.String("xlsx", "maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx"),
			zap.String("sheet", "maze_equip_suite_name_v8"))
		return 
	}
	// compare length
	if len(data) != len(gMazeEquipSuiteNameV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipSuiteNameV8ConfigRow", 
			zap.String("xlsx", "maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx"),
			zap.String("sheet", "maze_equip_suite_name_v8"), zap.Int("need_count",len(gMazeEquipSuiteNameV8Fields)), 
			zap.Int("had_count",len(data)))
		return
	}

	var tmp int64

	// parse column 0 equipment_id : 装备id 
	if data[0] != "" {
		tmp,err = strconv.ParseInt(data[0],10,64)
		if err != nil {
			err = errors.New("parse field equipment_id 装备id to int32 failed")
			logger.ErrorWF("parse field equipment_id 装备id to int32 failed.", 
				zap.String("xlsx", "maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx"), zap.String("sheet", "maze_equip_suite_name_v8"), 
				zap.String("parse_data",data[0]), 
				zap.Error(err))
			return
		}
		config.Equipment_id = int32(tmp)
	}

	// parse column 1 suite_equip_name : 套装装备名称 
	if data[1] != "" {

		config.Suite_equip_name = make(map[int32]string)
		var key int32
		var value string
		vals := strings.Split(data[1],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field suite_equip_name 套装装备名称 to key int32 failed")
				logger.ErrorWF("parse map field suite_equip_name 套装装备名称 to key int32 failed.", 
					zap.String("xlsx", "maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx"), zap.String("sheet", "maze_equip_suite_name_v8"), 
					// zap.String("field_data",data[1]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			value = items[1]
			config.Suite_equip_name[key] = value
		}
	}
	return
}

var gMazeEquipSuiteNameV8Fields = []string{
    "equipment_id",
    "suite_equip_name",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string,error)) (err error) {
	parser := &gMazeEquipSuiteNameV8Parser{}
	loader := &gMazeEquipSuiteNameV8Loader{}
	var data [][]string
	data,err = load("maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx", "maze_equip_suite_name_v8", gMazeEquipSuiteNameV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx maze_equip_suite_name_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx maze_equip_suite_name_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipSuiteNameV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx maze_equip_suite_name_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipSuiteNameV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx maze_equip_suite_name_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx maze_equip_suite_name_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx maze_equip_suite_name_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_suite_name_v8【迷宫-装备-套装装备名称】.xlsx maze_equip_suite_name_v8 data success.")
	return
}
