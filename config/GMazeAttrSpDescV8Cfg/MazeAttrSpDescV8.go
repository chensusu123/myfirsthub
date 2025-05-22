package GMazeAttrSpDescV8Cfg


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


// MazeAttrSpDescV8ConfigRow from maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx maze_attr_sp_desc_v8
type MazeAttrSpDescV8ConfigRow struct {
    Id       int32  `json:"id"` // 属性ID
    Equip_affix_desc       string  `json:"equip_affix_desc"` // 装备详情替换描述
    Equip_symbol       string  `json:"equip_symbol"` // 装备详情符号
    Equip_affix_suffix       string  `json:"equip_affix_suffix"` // 装备详情词条描述后缀
    Attr_list_affix_desc       string  `json:"attr_list_affix_desc"` // 属性列表替换描述
    Attr_list_affix_suffix       string  `json:"attr_list_affix_suffix"` // 属性列表后缀单位
}

// MazeAttrSpDescV8Config from maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx maze_attr_sp_desc_v8
type MazeAttrSpDescV8Config struct {
	ConfigRows map[int32]*MazeAttrSpDescV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeAttrSpDescV8Config {
	ret := &MazeAttrSpDescV8Config{ConfigRows: map[int32]*MazeAttrSpDescV8ConfigRow{}}
	return ret
}

// GetMazeAttrSpDescV8Config get one config by configId
func (c *MazeAttrSpDescV8Config) GetMazeAttrSpDescV8Config(configId int32) *MazeAttrSpDescV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeAttrSpDescV8Config) Get(configId int32) *MazeAttrSpDescV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeAttrSpDescV8Config get all config slice
func (c *MazeAttrSpDescV8Config)  GetAllMazeAttrSpDescV8Config () (res []*MazeAttrSpDescV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeAttrSpDescV8Config)  GetAll() (res []*MazeAttrSpDescV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}


// global config pointer 
var gConfigData *MazeAttrSpDescV8Config 

// GetMazeAttrSpDescV8Config pkg func. get one config by configId
func GetMazeAttrSpDescV8Config(configId int32) *MazeAttrSpDescV8ConfigRow {
	return gConfigData.GetMazeAttrSpDescV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeAttrSpDescV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeAttrSpDescV8Config pkg func. get all config slice
func GetAllMazeAttrSpDescV8Config () []*MazeAttrSpDescV8ConfigRow {
	return gConfigData.GetAllMazeAttrSpDescV8Config ()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeAttrSpDescV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeAttrSpDescV8ConfigRow{
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string{
	return "MazeAttrSpDescV8ConfigRow from maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx maze_attr_sp_desc_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeAttrSpDescV8Config{
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_attr_sp_desc_v8"
}

func init() {
	// reg config auto load 
	config_manager.RegAutoConfig("maze_attr_sp_desc_v8.json", 
		"maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx", "maze_attr_sp_desc_v8",
	 	&gMazeAttrSpDescV8Parser{}, &gMazeAttrSpDescV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeAttrSpDescV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeAttrSpDescV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeAttrSpDescV8Config))(c)
		return true
	})
}

// RegisterMazeAttrSpDescV8InitCallBack reg config update func (old func)
var RegisterMazeAttrSpDescV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeAttrSpDescV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeAttrSpDescV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeAttrSpDescV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeAttrSpDescV8Config)error)(c)
		return err == nil
	})
	return 
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeAttrSpDescV8Config)error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeAttrSpDescV8Loader struct {
}
// NewContainer new data container pointer
func (*gMazeAttrSpDescV8Loader) NewContainer() interface{}{
	return newConfig()
}
// Check check new config data ptr
func (*gMazeAttrSpDescV8Loader) Check(newPtr interface{})error{
	// set global ptr
	cfgData := newPtr.(*MazeAttrSpDescV8Config)
	return doConfigCheckCallback(cfgData)
}
// Swap swap global config data ptr
func (*gMazeAttrSpDescV8Loader) Swap(newPtr interface{}){
	// convert pointer
	cache := newPtr.(*MazeAttrSpDescV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}
// Add append config item to map,save result
func (*gMazeAttrSpDescV8Loader) Add(logger fklog.FKLogI,container interface{}, ri interface{})(err error) {
	row,ok := ri.(*MazeAttrSpDescV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSpDescV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAttrSpDescV8ConfigRow", zap.String("xlsx", "maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx"),
			zap.String("sheet", "maze_attr_sp_desc_v8"))
		return 
	}
	config,ok := container.(*MazeAttrSpDescV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSpDescV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrSpDescV8Config", zap.String("xlsx", "maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx"),
			zap.String("sheet", "maze_attr_sp_desc_v8"))
		return 
	}
	config.ConfigRows[row.Id] = row
	return
}
// GetValue get real map value for json parse
func (*gMazeAttrSpDescV8Loader) GetValue(logger fklog.FKLogI,container interface{})(real interface{},err error) {
	config,ok := container.(*MazeAttrSpDescV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSpDescV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrSpDescV8Config", zap.String("xlsx", "maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx"),
			zap.String("sheet", "maze_attr_sp_desc_v8"))
		return 
	}
	real = &config.ConfigRows
	return
}
// Range range all data for json append data parse.
func (*gMazeAttrSpDescV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error){
	config,ok := container.(*MazeAttrSpDescV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSpDescV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrSpDescV8Config", zap.String("xlsx", "maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx"),
			zap.String("sheet", "maze_attr_sp_desc_v8"))
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
type gMazeAttrSpDescV8Parser struct {
}
// New new config row data
func (*gMazeAttrSpDescV8Parser) New() interface{} {
	return &MazeAttrSpDescV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeAttrSpDescV8Parser) Fields() []string {
	return gMazeAttrSpDescV8Fields
}
// Parse parse raw data to row data
func (*gMazeAttrSpDescV8Parser) Parse(logger fklog.FKLogI,data []string, row interface{}) (err error) {
	// convert row type
	config,ok := row.(*MazeAttrSpDescV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSpDescV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAttrSpDescV8ConfigRow", zap.String("xlsx", "maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx"),
			zap.String("sheet", "maze_attr_sp_desc_v8"))
		return 
	}
	// compare length
	if len(data) != len(gMazeAttrSpDescV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeAttrSpDescV8ConfigRow", 
			zap.String("xlsx", "maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx"),
			zap.String("sheet", "maze_attr_sp_desc_v8"), zap.Int("need_count",len(gMazeAttrSpDescV8Fields)), 
			zap.Int("had_count",len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 属性ID 
	if data[0] != "" {
		tmp,err = strconv.ParseInt(data[0],10,64)
		if err != nil {
			err = errors.New("parse field id 属性ID to int32 failed")
			logger.ErrorWF("parse field id 属性ID to int32 failed.", 
				zap.String("xlsx", "maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx"), zap.String("sheet", "maze_attr_sp_desc_v8"), 
				zap.String("parse_data",data[0]), 
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 equip_affix_desc : 装备详情替换描述 
	if data[1] != "" {
		config.Equip_affix_desc = data[1]
	}

	// parse column 2 equip_symbol : 装备详情符号 
	if data[2] != "" {
		config.Equip_symbol = data[2]
	}

	// parse column 3 equip_affix_suffix : 装备详情词条描述后缀 
	if data[3] != "" {
		config.Equip_affix_suffix = data[3]
	}

	// parse column 4 attr_list_affix_desc : 属性列表替换描述 
	if data[4] != "" {
		config.Attr_list_affix_desc = data[4]
	}

	// parse column 5 attr_list_affix_suffix : 属性列表后缀单位 
	if data[5] != "" {
		config.Attr_list_affix_suffix = data[5]
	}
	return
}

var gMazeAttrSpDescV8Fields = []string{
    "id",
    "equip_affix_desc",
    "equip_symbol",
    "equip_affix_suffix",
    "attr_list_affix_desc",
    "attr_list_affix_suffix",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string,error)) (err error) {
	parser := &gMazeAttrSpDescV8Parser{}
	loader := &gMazeAttrSpDescV8Loader{}
	var data [][]string
	data,err = load("maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx", "maze_attr_sp_desc_v8", gMazeAttrSpDescV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx maze_attr_sp_desc_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx maze_attr_sp_desc_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeAttrSpDescV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx maze_attr_sp_desc_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeAttrSpDescV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx maze_attr_sp_desc_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx maze_attr_sp_desc_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx maze_attr_sp_desc_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_attr_sp_desc_v8【迷宫-属性-特殊展示】.xlsx maze_attr_sp_desc_v8 data success.")
	return
}
