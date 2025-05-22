package GMazeEquipAttrStageV8Cfg


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


// MazeEquipAttrStageV8ConfigRow from maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx maze_equip_attr_stage_v8
type MazeEquipAttrStageV8ConfigRow struct {
    Stage_id       int32  `json:"stage_id"` // 所属档位
    Ratio       int32  `json:"ratio"` // 属性系数
    Score       int32  `json:"score"` // 掉落分数
}

// MazeEquipAttrStageV8Config from maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx maze_equip_attr_stage_v8
type MazeEquipAttrStageV8Config struct {
	ConfigRows map[int32]*MazeEquipAttrStageV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipAttrStageV8Config {
	ret := &MazeEquipAttrStageV8Config{ConfigRows: map[int32]*MazeEquipAttrStageV8ConfigRow{}}
	return ret
}

// GetMazeEquipAttrStageV8Config get one config by configId
func (c *MazeEquipAttrStageV8Config) GetMazeEquipAttrStageV8Config(configId int32) *MazeEquipAttrStageV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipAttrStageV8Config) Get(configId int32) *MazeEquipAttrStageV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipAttrStageV8Config get all config slice
func (c *MazeEquipAttrStageV8Config)  GetAllMazeEquipAttrStageV8Config () (res []*MazeEquipAttrStageV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipAttrStageV8Config)  GetAll() (res []*MazeEquipAttrStageV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}


// global config pointer 
var gConfigData *MazeEquipAttrStageV8Config 

// GetMazeEquipAttrStageV8Config pkg func. get one config by configId
func GetMazeEquipAttrStageV8Config(configId int32) *MazeEquipAttrStageV8ConfigRow {
	return gConfigData.GetMazeEquipAttrStageV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipAttrStageV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeEquipAttrStageV8Config pkg func. get all config slice
func GetAllMazeEquipAttrStageV8Config () []*MazeEquipAttrStageV8ConfigRow {
	return gConfigData.GetAllMazeEquipAttrStageV8Config ()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipAttrStageV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipAttrStageV8ConfigRow{
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string{
	return "MazeEquipAttrStageV8ConfigRow from maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx maze_equip_attr_stage_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipAttrStageV8Config{
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_attr_stage_v8"
}

func init() {
	// reg config auto load 
	config_manager.RegAutoConfig("maze_equip_attr_stage_v8.json", 
		"maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx", "maze_equip_attr_stage_v8",
	 	&gMazeEquipAttrStageV8Parser{}, &gMazeEquipAttrStageV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipAttrStageV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipAttrStageV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipAttrStageV8Config))(c)
		return true
	})
}

// RegisterMazeEquipAttrStageV8InitCallBack reg config update func (old func)
var RegisterMazeEquipAttrStageV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipAttrStageV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipAttrStageV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipAttrStageV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipAttrStageV8Config)error)(c)
		return err == nil
	})
	return 
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipAttrStageV8Config)error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipAttrStageV8Loader struct {
}
// NewContainer new data container pointer
func (*gMazeEquipAttrStageV8Loader) NewContainer() interface{}{
	return newConfig()
}
// Check check new config data ptr
func (*gMazeEquipAttrStageV8Loader) Check(newPtr interface{})error{
	// set global ptr
	cfgData := newPtr.(*MazeEquipAttrStageV8Config)
	return doConfigCheckCallback(cfgData)
}
// Swap swap global config data ptr
func (*gMazeEquipAttrStageV8Loader) Swap(newPtr interface{}){
	// convert pointer
	cache := newPtr.(*MazeEquipAttrStageV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}
// Add append config item to map,save result
func (*gMazeEquipAttrStageV8Loader) Add(logger fklog.FKLogI,container interface{}, ri interface{})(err error) {
	row,ok := ri.(*MazeEquipAttrStageV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAttrStageV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipAttrStageV8ConfigRow", zap.String("xlsx", "maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx"),
			zap.String("sheet", "maze_equip_attr_stage_v8"))
		return 
	}
	config,ok := container.(*MazeEquipAttrStageV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAttrStageV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAttrStageV8Config", zap.String("xlsx", "maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx"),
			zap.String("sheet", "maze_equip_attr_stage_v8"))
		return 
	}
	config.ConfigRows[row.Stage_id] = row
	return
}
// GetValue get real map value for json parse
func (*gMazeEquipAttrStageV8Loader) GetValue(logger fklog.FKLogI,container interface{})(real interface{},err error) {
	config,ok := container.(*MazeEquipAttrStageV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAttrStageV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAttrStageV8Config", zap.String("xlsx", "maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx"),
			zap.String("sheet", "maze_equip_attr_stage_v8"))
		return 
	}
	real = &config.ConfigRows
	return
}
// Range range all data for json append data parse.
func (*gMazeEquipAttrStageV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error){
	config,ok := container.(*MazeEquipAttrStageV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAttrStageV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipAttrStageV8Config", zap.String("xlsx", "maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx"),
			zap.String("sheet", "maze_equip_attr_stage_v8"))
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
type gMazeEquipAttrStageV8Parser struct {
}
// New new config row data
func (*gMazeEquipAttrStageV8Parser) New() interface{} {
	return &MazeEquipAttrStageV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipAttrStageV8Parser) Fields() []string {
	return gMazeEquipAttrStageV8Fields
}
// Parse parse raw data to row data
func (*gMazeEquipAttrStageV8Parser) Parse(logger fklog.FKLogI,data []string, row interface{}) (err error) {
	// convert row type
	config,ok := row.(*MazeEquipAttrStageV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipAttrStageV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipAttrStageV8ConfigRow", zap.String("xlsx", "maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx"),
			zap.String("sheet", "maze_equip_attr_stage_v8"))
		return 
	}
	// compare length
	if len(data) != len(gMazeEquipAttrStageV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipAttrStageV8ConfigRow", 
			zap.String("xlsx", "maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx"),
			zap.String("sheet", "maze_equip_attr_stage_v8"), zap.Int("need_count",len(gMazeEquipAttrStageV8Fields)), 
			zap.Int("had_count",len(data)))
		return
	}

	var tmp int64

	// parse column 0 stage_id : 所属档位 
	if data[0] != "" {
		tmp,err = strconv.ParseInt(data[0],10,64)
		if err != nil {
			err = errors.New("parse field stage_id 所属档位 to int32 failed")
			logger.ErrorWF("parse field stage_id 所属档位 to int32 failed.", 
				zap.String("xlsx", "maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx"), zap.String("sheet", "maze_equip_attr_stage_v8"), 
				zap.String("parse_data",data[0]), 
				zap.Error(err))
			return
		}
		config.Stage_id = int32(tmp)
	}

	// parse column 1 ratio : 属性系数 
	if data[1] != "" {
		tmp,err = strconv.ParseInt(data[1],10,64)
		if err != nil {
			err = errors.New("parse field ratio 属性系数 to int32 failed")
			logger.ErrorWF("parse field ratio 属性系数 to int32 failed.", 
				zap.String("xlsx", "maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx"), zap.String("sheet", "maze_equip_attr_stage_v8"), 
				zap.String("parse_data",data[1]), 
				zap.Error(err))
			return
		}
		config.Ratio = int32(tmp)
	}

	// parse column 2 score : 掉落分数 
	if data[2] != "" {
		tmp,err = strconv.ParseInt(data[2],10,64)
		if err != nil {
			err = errors.New("parse field score 掉落分数 to int32 failed")
			logger.ErrorWF("parse field score 掉落分数 to int32 failed.", 
				zap.String("xlsx", "maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx"), zap.String("sheet", "maze_equip_attr_stage_v8"), 
				zap.String("parse_data",data[2]), 
				zap.Error(err))
			return
		}
		config.Score = int32(tmp)
	}
	return
}

var gMazeEquipAttrStageV8Fields = []string{
    "stage_id",
    "ratio",
    "score",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string,error)) (err error) {
	parser := &gMazeEquipAttrStageV8Parser{}
	loader := &gMazeEquipAttrStageV8Loader{}
	var data [][]string
	data,err = load("maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx", "maze_equip_attr_stage_v8", gMazeEquipAttrStageV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx maze_equip_attr_stage_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx maze_equip_attr_stage_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipAttrStageV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx maze_equip_attr_stage_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipAttrStageV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx maze_equip_attr_stage_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx maze_equip_attr_stage_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx maze_equip_attr_stage_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_attr_stage_v8【迷宫-装备-词条属性档位】.xlsx maze_equip_attr_stage_v8 data success.")
	return
}
