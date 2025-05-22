package GMazeEquipSuiteInfoV8Cfg


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


// MazeEquipSuiteInfoV8ConfigRow from maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx maze_equip_suite_info_v8
type MazeEquipSuiteInfoV8ConfigRow struct {
    Suite_id       int32  `json:"suite_id"` // 套装id
    Suite_name       string  `json:"suite_name"` // 套装名称
    Max_num       int32  `json:"max_num"` // 套装数量
    Pos_list       []int32  `json:"pos_list"` // 套装部位
    Exclude_attr       []int32  `json:"exclude_attr"` // 套装的排除
    Counter_weapons_damage_attr       int32  `json:"counter_weapons_damage_attr"` // 套装对应的武器属性id
    Counter_weapons_damage_type       int32  `json:"counter_weapons_damage_type"` // 套装对应的武器伤害类型
    Exclude_sub_type       []int32  `json:"exclude_sub_type"` // 套装排除的武器子类型
    Prefix       string  `json:"prefix"` // 装备名前缀
}

// MazeEquipSuiteInfoV8Config from maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx maze_equip_suite_info_v8
type MazeEquipSuiteInfoV8Config struct {
	ConfigRows map[int32]*MazeEquipSuiteInfoV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipSuiteInfoV8Config {
	ret := &MazeEquipSuiteInfoV8Config{ConfigRows: map[int32]*MazeEquipSuiteInfoV8ConfigRow{}}
	return ret
}

// GetMazeEquipSuiteInfoV8Config get one config by configId
func (c *MazeEquipSuiteInfoV8Config) GetMazeEquipSuiteInfoV8Config(configId int32) *MazeEquipSuiteInfoV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipSuiteInfoV8Config) Get(configId int32) *MazeEquipSuiteInfoV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipSuiteInfoV8Config get all config slice
func (c *MazeEquipSuiteInfoV8Config)  GetAllMazeEquipSuiteInfoV8Config () (res []*MazeEquipSuiteInfoV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipSuiteInfoV8Config)  GetAll() (res []*MazeEquipSuiteInfoV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}


// global config pointer 
var gConfigData *MazeEquipSuiteInfoV8Config 

// GetMazeEquipSuiteInfoV8Config pkg func. get one config by configId
func GetMazeEquipSuiteInfoV8Config(configId int32) *MazeEquipSuiteInfoV8ConfigRow {
	return gConfigData.GetMazeEquipSuiteInfoV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipSuiteInfoV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeEquipSuiteInfoV8Config pkg func. get all config slice
func GetAllMazeEquipSuiteInfoV8Config () []*MazeEquipSuiteInfoV8ConfigRow {
	return gConfigData.GetAllMazeEquipSuiteInfoV8Config ()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipSuiteInfoV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipSuiteInfoV8ConfigRow{
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string{
	return "MazeEquipSuiteInfoV8ConfigRow from maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx maze_equip_suite_info_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipSuiteInfoV8Config{
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_suite_info_v8"
}

func init() {
	// reg config auto load 
	config_manager.RegAutoConfig("maze_equip_suite_info_v8.json", 
		"maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx", "maze_equip_suite_info_v8",
	 	&gMazeEquipSuiteInfoV8Parser{}, &gMazeEquipSuiteInfoV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipSuiteInfoV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipSuiteInfoV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipSuiteInfoV8Config))(c)
		return true
	})
}

// RegisterMazeEquipSuiteInfoV8InitCallBack reg config update func (old func)
var RegisterMazeEquipSuiteInfoV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipSuiteInfoV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipSuiteInfoV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipSuiteInfoV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipSuiteInfoV8Config)error)(c)
		return err == nil
	})
	return 
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipSuiteInfoV8Config)error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipSuiteInfoV8Loader struct {
}
// NewContainer new data container pointer
func (*gMazeEquipSuiteInfoV8Loader) NewContainer() interface{}{
	return newConfig()
}
// Check check new config data ptr
func (*gMazeEquipSuiteInfoV8Loader) Check(newPtr interface{})error{
	// set global ptr
	cfgData := newPtr.(*MazeEquipSuiteInfoV8Config)
	return doConfigCheckCallback(cfgData)
}
// Swap swap global config data ptr
func (*gMazeEquipSuiteInfoV8Loader) Swap(newPtr interface{}){
	// convert pointer
	cache := newPtr.(*MazeEquipSuiteInfoV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}
// Add append config item to map,save result
func (*gMazeEquipSuiteInfoV8Loader) Add(logger fklog.FKLogI,container interface{}, ri interface{})(err error) {
	row,ok := ri.(*MazeEquipSuiteInfoV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteInfoV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteInfoV8ConfigRow", zap.String("xlsx", "maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx"),
			zap.String("sheet", "maze_equip_suite_info_v8"))
		return 
	}
	config,ok := container.(*MazeEquipSuiteInfoV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteInfoV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteInfoV8Config", zap.String("xlsx", "maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx"),
			zap.String("sheet", "maze_equip_suite_info_v8"))
		return 
	}
	config.ConfigRows[row.Suite_id] = row
	return
}
// GetValue get real map value for json parse
func (*gMazeEquipSuiteInfoV8Loader) GetValue(logger fklog.FKLogI,container interface{})(real interface{},err error) {
	config,ok := container.(*MazeEquipSuiteInfoV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteInfoV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteInfoV8Config", zap.String("xlsx", "maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx"),
			zap.String("sheet", "maze_equip_suite_info_v8"))
		return 
	}
	real = &config.ConfigRows
	return
}
// Range range all data for json append data parse.
func (*gMazeEquipSuiteInfoV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error){
	config,ok := container.(*MazeEquipSuiteInfoV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteInfoV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteInfoV8Config", zap.String("xlsx", "maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx"),
			zap.String("sheet", "maze_equip_suite_info_v8"))
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
type gMazeEquipSuiteInfoV8Parser struct {
}
// New new config row data
func (*gMazeEquipSuiteInfoV8Parser) New() interface{} {
	return &MazeEquipSuiteInfoV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipSuiteInfoV8Parser) Fields() []string {
	return gMazeEquipSuiteInfoV8Fields
}
// Parse parse raw data to row data
func (*gMazeEquipSuiteInfoV8Parser) Parse(logger fklog.FKLogI,data []string, row interface{}) (err error) {
	// convert row type
	config,ok := row.(*MazeEquipSuiteInfoV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteInfoV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteInfoV8ConfigRow", zap.String("xlsx", "maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx"),
			zap.String("sheet", "maze_equip_suite_info_v8"))
		return 
	}
	// compare length
	if len(data) != len(gMazeEquipSuiteInfoV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipSuiteInfoV8ConfigRow", 
			zap.String("xlsx", "maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx"),
			zap.String("sheet", "maze_equip_suite_info_v8"), zap.Int("need_count",len(gMazeEquipSuiteInfoV8Fields)), 
			zap.Int("had_count",len(data)))
		return
	}

	var tmp int64

	// parse column 0 suite_id : 套装id 
	if data[0] != "" {
		tmp,err = strconv.ParseInt(data[0],10,64)
		if err != nil {
			err = errors.New("parse field suite_id 套装id to int32 failed")
			logger.ErrorWF("parse field suite_id 套装id to int32 failed.", 
				zap.String("xlsx", "maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx"), zap.String("sheet", "maze_equip_suite_info_v8"), 
				zap.String("parse_data",data[0]), 
				zap.Error(err))
			return
		}
		config.Suite_id = int32(tmp)
	}

	// parse column 1 suite_name : 套装名称 
	if data[1] != "" {
		config.Suite_name = data[1]
	}

	// parse column 2 max_num : 套装数量 
	if data[2] != "" {
		tmp,err = strconv.ParseInt(data[2],10,64)
		if err != nil {
			err = errors.New("parse field max_num 套装数量 to int32 failed")
			logger.ErrorWF("parse field max_num 套装数量 to int32 failed.", 
				zap.String("xlsx", "maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx"), zap.String("sheet", "maze_equip_suite_info_v8"), 
				zap.String("parse_data",data[2]), 
				zap.Error(err))
			return
		}
		config.Max_num = int32(tmp)
	}

	// parse column 3 pos_list : 套装部位 
	if data[3] != "" {
    
		vals := strings.Split(data[3],",")
		for k,v := range vals {
			tmp,err = strconv.ParseInt(v,10,64)
			if err != nil {
				err = errors.New("parse array field pos_list 套装部位 to []int32 failed")
				logger.ErrorWF("parse array field pos_list 套装部位 to []int32 failed.", 
					zap.String("xlsx", "maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx"), zap.String("sheet", "maze_equip_suite_info_v8"), 
					// zap.String("field_data",data[3]), 
					zap.String("parse_data", v),zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Pos_list = append(config.Pos_list, int32(tmp))
		}
	}

	// parse column 4 exclude_attr : 套装的排除 
	if data[4] != "" {
    
		vals := strings.Split(data[4],",")
		for k,v := range vals {
			tmp,err = strconv.ParseInt(v,10,64)
			if err != nil {
				err = errors.New("parse array field exclude_attr 套装的排除 to []int32 failed")
				logger.ErrorWF("parse array field exclude_attr 套装的排除 to []int32 failed.", 
					zap.String("xlsx", "maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx"), zap.String("sheet", "maze_equip_suite_info_v8"), 
					// zap.String("field_data",data[4]), 
					zap.String("parse_data", v),zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Exclude_attr = append(config.Exclude_attr, int32(tmp))
		}
	}

	// parse column 5 counter_weapons_damage_attr : 套装对应的武器属性id 
	if data[5] != "" {
		tmp,err = strconv.ParseInt(data[5],10,64)
		if err != nil {
			err = errors.New("parse field counter_weapons_damage_attr 套装对应的武器属性id to int32 failed")
			logger.ErrorWF("parse field counter_weapons_damage_attr 套装对应的武器属性id to int32 failed.", 
				zap.String("xlsx", "maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx"), zap.String("sheet", "maze_equip_suite_info_v8"), 
				zap.String("parse_data",data[5]), 
				zap.Error(err))
			return
		}
		config.Counter_weapons_damage_attr = int32(tmp)
	}

	// parse column 6 counter_weapons_damage_type : 套装对应的武器伤害类型 
	if data[6] != "" {
		tmp,err = strconv.ParseInt(data[6],10,64)
		if err != nil {
			err = errors.New("parse field counter_weapons_damage_type 套装对应的武器伤害类型 to int32 failed")
			logger.ErrorWF("parse field counter_weapons_damage_type 套装对应的武器伤害类型 to int32 failed.", 
				zap.String("xlsx", "maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx"), zap.String("sheet", "maze_equip_suite_info_v8"), 
				zap.String("parse_data",data[6]), 
				zap.Error(err))
			return
		}
		config.Counter_weapons_damage_type = int32(tmp)
	}

	// parse column 7 exclude_sub_type : 套装排除的武器子类型 
	if data[7] != "" {
    
		vals := strings.Split(data[7],",")
		for k,v := range vals {
			tmp,err = strconv.ParseInt(v,10,64)
			if err != nil {
				err = errors.New("parse array field exclude_sub_type 套装排除的武器子类型 to []int32 failed")
				logger.ErrorWF("parse array field exclude_sub_type 套装排除的武器子类型 to []int32 failed.", 
					zap.String("xlsx", "maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx"), zap.String("sheet", "maze_equip_suite_info_v8"), 
					// zap.String("field_data",data[7]), 
					zap.String("parse_data", v),zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Exclude_sub_type = append(config.Exclude_sub_type, int32(tmp))
		}
	}

	// parse column 8 prefix : 装备名前缀 
	if data[8] != "" {
		config.Prefix = data[8]
	}
	return
}

var gMazeEquipSuiteInfoV8Fields = []string{
    "suite_id",
    "suite_name",
    "max_num",
    "pos_list",
    "exclude_attr",
    "counter_weapons_damage_attr",
    "counter_weapons_damage_type",
    "exclude_sub_type",
    "prefix",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string,error)) (err error) {
	parser := &gMazeEquipSuiteInfoV8Parser{}
	loader := &gMazeEquipSuiteInfoV8Loader{}
	var data [][]string
	data,err = load("maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx", "maze_equip_suite_info_v8", gMazeEquipSuiteInfoV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx maze_equip_suite_info_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx maze_equip_suite_info_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipSuiteInfoV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx maze_equip_suite_info_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipSuiteInfoV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx maze_equip_suite_info_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx maze_equip_suite_info_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx maze_equip_suite_info_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_suite_info_v8【迷宫-装备-套装信息】.xlsx maze_equip_suite_info_v8 data success.")
	return
}
