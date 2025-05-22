package GMazeSkillAutoReleaseV8Cfg


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


// MazeSkillAutoReleaseV8ConfigRow from maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8
type MazeSkillAutoReleaseV8ConfigRow struct {
    Order       int32  `json:"order"` // auto_skill_id
    Attr       int32  `json:"attr"` // 逻辑属性id
    Attr_value_1_variable_id       map[int32]int32  `json:"attr_value_1_variable_id"` // 参数1关联的变量id
    Attr_value_1_type       int32  `json:"attr_value_1_type"` // 参数1数值类型
    Attr_value_1       int32  `json:"attr_value_1"` // 参数1
    Attr_value_2_variable_id       map[int32]int32  `json:"attr_value_2_variable_id"` // 触发几率关联的变量id
    Attr_value_2_type       int32  `json:"attr_value_2_type"` // 触发几率数值类型
    Attr_value_2       int32  `json:"attr_value_2"` // 触发几率
    Attr_value_3_variable_id       map[int32]int32  `json:"attr_value_3_variable_id"` // 释放次数关联的变量id
    Attr_value_3_type       int32  `json:"attr_value_3_type"` // 释放次数数值类型
    Attr_value_3       int32  `json:"attr_value_3"` // 释放次数
    Attr_value_4       []int32  `json:"attr_value_4"` // 参数4
    Last_time_variable_id       map[int32]int32  `json:"last_time_variable_id"` // 持续时长（毫秒）关联的变量id
    Last_time       int32  `json:"last_time"` // 持续时长（毫秒）
    Base_hitrate_variable_id       map[int32]int32  `json:"base_hitrate_variable_id"` // 基础命中率变量
    Base_hitrate       int32  `json:"base_hitrate"` // 基础命中率（万分比）
    Attr_value_7_variable_id       map[int32]int32  `json:"attr_value_7_variable_id"` // 参数7关联的变量id
    Attr_value_7_type       int32  `json:"attr_value_7_type"` // 参数7数值类型
    Attr_value_7       int32  `json:"attr_value_7"` // 参数7
    Attr_value_8_variable_id       map[int32]int32  `json:"attr_value_8_variable_id"` // 参数8关联的变量id
    Attr_value_8_type       int32  `json:"attr_value_8_type"` // 参数8数值类型
    Attr_value_8       int32  `json:"attr_value_8"` // 参数8
}

// MazeSkillAutoReleaseV8Config from maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8
type MazeSkillAutoReleaseV8Config struct {
	ConfigRows map[int32]*MazeSkillAutoReleaseV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeSkillAutoReleaseV8Config {
	ret := &MazeSkillAutoReleaseV8Config{ConfigRows: map[int32]*MazeSkillAutoReleaseV8ConfigRow{}}
	return ret
}

// GetMazeSkillAutoReleaseV8Config get one config by configId
func (c *MazeSkillAutoReleaseV8Config) GetMazeSkillAutoReleaseV8Config(configId int32) *MazeSkillAutoReleaseV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeSkillAutoReleaseV8Config) Get(configId int32) *MazeSkillAutoReleaseV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeSkillAutoReleaseV8Config get all config slice
func (c *MazeSkillAutoReleaseV8Config)  GetAllMazeSkillAutoReleaseV8Config () (res []*MazeSkillAutoReleaseV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeSkillAutoReleaseV8Config)  GetAll() (res []*MazeSkillAutoReleaseV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}


// global config pointer 
var gConfigData *MazeSkillAutoReleaseV8Config 

// GetMazeSkillAutoReleaseV8Config pkg func. get one config by configId
func GetMazeSkillAutoReleaseV8Config(configId int32) *MazeSkillAutoReleaseV8ConfigRow {
	return gConfigData.GetMazeSkillAutoReleaseV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeSkillAutoReleaseV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeSkillAutoReleaseV8Config pkg func. get all config slice
func GetAllMazeSkillAutoReleaseV8Config () []*MazeSkillAutoReleaseV8ConfigRow {
	return gConfigData.GetAllMazeSkillAutoReleaseV8Config ()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeSkillAutoReleaseV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeSkillAutoReleaseV8ConfigRow{
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string{
	return "MazeSkillAutoReleaseV8ConfigRow from maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeSkillAutoReleaseV8Config{
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_skill_auto_release_v8"
}

func init() {
	// reg config auto load 
	config_manager.RegAutoConfig("maze_skill_auto_release_v8.json", 
		"maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx", "maze_skill_auto_release_v8",
	 	&gMazeSkillAutoReleaseV8Parser{}, &gMazeSkillAutoReleaseV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeSkillAutoReleaseV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeSkillAutoReleaseV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeSkillAutoReleaseV8Config))(c)
		return true
	})
}

// RegisterMazeSkillAutoReleaseV8InitCallBack reg config update func (old func)
var RegisterMazeSkillAutoReleaseV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeSkillAutoReleaseV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeSkillAutoReleaseV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeSkillAutoReleaseV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeSkillAutoReleaseV8Config)error)(c)
		return err == nil
	})
	return 
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeSkillAutoReleaseV8Config)error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeSkillAutoReleaseV8Loader struct {
}
// NewContainer new data container pointer
func (*gMazeSkillAutoReleaseV8Loader) NewContainer() interface{}{
	return newConfig()
}
// Check check new config data ptr
func (*gMazeSkillAutoReleaseV8Loader) Check(newPtr interface{})error{
	// set global ptr
	cfgData := newPtr.(*MazeSkillAutoReleaseV8Config)
	return doConfigCheckCallback(cfgData)
}
// Swap swap global config data ptr
func (*gMazeSkillAutoReleaseV8Loader) Swap(newPtr interface{}){
	// convert pointer
	cache := newPtr.(*MazeSkillAutoReleaseV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}
// Add append config item to map,save result
func (*gMazeSkillAutoReleaseV8Loader) Add(logger fklog.FKLogI,container interface{}, ri interface{})(err error) {
	row,ok := ri.(*MazeSkillAutoReleaseV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoReleaseV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkillAutoReleaseV8ConfigRow", zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"),
			zap.String("sheet", "maze_skill_auto_release_v8"))
		return 
	}
	config,ok := container.(*MazeSkillAutoReleaseV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoReleaseV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillAutoReleaseV8Config", zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"),
			zap.String("sheet", "maze_skill_auto_release_v8"))
		return 
	}
	config.ConfigRows[row.Order] = row
	return
}
// GetValue get real map value for json parse
func (*gMazeSkillAutoReleaseV8Loader) GetValue(logger fklog.FKLogI,container interface{})(real interface{},err error) {
	config,ok := container.(*MazeSkillAutoReleaseV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoReleaseV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillAutoReleaseV8Config", zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"),
			zap.String("sheet", "maze_skill_auto_release_v8"))
		return 
	}
	real = &config.ConfigRows
	return
}
// Range range all data for json append data parse.
func (*gMazeSkillAutoReleaseV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error){
	config,ok := container.(*MazeSkillAutoReleaseV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoReleaseV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillAutoReleaseV8Config", zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"),
			zap.String("sheet", "maze_skill_auto_release_v8"))
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
type gMazeSkillAutoReleaseV8Parser struct {
}
// New new config row data
func (*gMazeSkillAutoReleaseV8Parser) New() interface{} {
	return &MazeSkillAutoReleaseV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeSkillAutoReleaseV8Parser) Fields() []string {
	return gMazeSkillAutoReleaseV8Fields
}
// Parse parse raw data to row data
func (*gMazeSkillAutoReleaseV8Parser) Parse(logger fklog.FKLogI,data []string, row interface{}) (err error) {
	// convert row type
	config,ok := row.(*MazeSkillAutoReleaseV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoReleaseV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkillAutoReleaseV8ConfigRow", zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"),
			zap.String("sheet", "maze_skill_auto_release_v8"))
		return 
	}
	// compare length
	if len(data) != len(gMazeSkillAutoReleaseV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeSkillAutoReleaseV8ConfigRow", 
			zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"),
			zap.String("sheet", "maze_skill_auto_release_v8"), zap.Int("need_count",len(gMazeSkillAutoReleaseV8Fields)), 
			zap.Int("had_count",len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : auto_skill_id 
	if data[0] != "" {
		tmp,err = strconv.ParseInt(data[0],10,64)
		if err != nil {
			err = errors.New("parse field order auto_skill_id to int32 failed")
			logger.ErrorWF("parse field order auto_skill_id to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[0]), 
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 attr : 逻辑属性id 
	if data[1] != "" {
		tmp,err = strconv.ParseInt(data[1],10,64)
		if err != nil {
			err = errors.New("parse field attr 逻辑属性id to int32 failed")
			logger.ErrorWF("parse field attr 逻辑属性id to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[1]), 
				zap.Error(err))
			return
		}
		config.Attr = int32(tmp)
	}

	// parse column 2 attr_value_1_variable_id : 参数1关联的变量id 
	if data[2] != "" {

		config.Attr_value_1_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[2],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field attr_value_1_variable_id 参数1关联的变量id to key int32 failed")
				logger.ErrorWF("parse map field attr_value_1_variable_id 参数1关联的变量id to key int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[2]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field attr_value_1_variable_id 参数1关联的变量id to value int32 failed")
				logger.ErrorWF("parse map field attr_value_1_variable_id 参数1关联的变量id to value int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[2]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Attr_value_1_variable_id[key] = value
		}
	}

	// parse column 3 attr_value_1_type : 参数1数值类型 
	if data[3] != "" {
		tmp,err = strconv.ParseInt(data[3],10,64)
		if err != nil {
			err = errors.New("parse field attr_value_1_type 参数1数值类型 to int32 failed")
			logger.ErrorWF("parse field attr_value_1_type 参数1数值类型 to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[3]), 
				zap.Error(err))
			return
		}
		config.Attr_value_1_type = int32(tmp)
	}

	// parse column 4 attr_value_1 : 参数1 
	if data[4] != "" {
		tmp,err = strconv.ParseInt(data[4],10,64)
		if err != nil {
			err = errors.New("parse field attr_value_1 参数1 to int32 failed")
			logger.ErrorWF("parse field attr_value_1 参数1 to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[4]), 
				zap.Error(err))
			return
		}
		config.Attr_value_1 = int32(tmp)
	}

	// parse column 5 attr_value_2_variable_id : 触发几率关联的变量id 
	if data[5] != "" {

		config.Attr_value_2_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[5],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field attr_value_2_variable_id 触发几率关联的变量id to key int32 failed")
				logger.ErrorWF("parse map field attr_value_2_variable_id 触发几率关联的变量id to key int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[5]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field attr_value_2_variable_id 触发几率关联的变量id to value int32 failed")
				logger.ErrorWF("parse map field attr_value_2_variable_id 触发几率关联的变量id to value int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[5]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Attr_value_2_variable_id[key] = value
		}
	}

	// parse column 6 attr_value_2_type : 触发几率数值类型 
	if data[6] != "" {
		tmp,err = strconv.ParseInt(data[6],10,64)
		if err != nil {
			err = errors.New("parse field attr_value_2_type 触发几率数值类型 to int32 failed")
			logger.ErrorWF("parse field attr_value_2_type 触发几率数值类型 to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[6]), 
				zap.Error(err))
			return
		}
		config.Attr_value_2_type = int32(tmp)
	}

	// parse column 7 attr_value_2 : 触发几率 
	if data[7] != "" {
		tmp,err = strconv.ParseInt(data[7],10,64)
		if err != nil {
			err = errors.New("parse field attr_value_2 触发几率 to int32 failed")
			logger.ErrorWF("parse field attr_value_2 触发几率 to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[7]), 
				zap.Error(err))
			return
		}
		config.Attr_value_2 = int32(tmp)
	}

	// parse column 8 attr_value_3_variable_id : 释放次数关联的变量id 
	if data[8] != "" {

		config.Attr_value_3_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[8],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field attr_value_3_variable_id 释放次数关联的变量id to key int32 failed")
				logger.ErrorWF("parse map field attr_value_3_variable_id 释放次数关联的变量id to key int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[8]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field attr_value_3_variable_id 释放次数关联的变量id to value int32 failed")
				logger.ErrorWF("parse map field attr_value_3_variable_id 释放次数关联的变量id to value int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[8]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Attr_value_3_variable_id[key] = value
		}
	}

	// parse column 9 attr_value_3_type : 释放次数数值类型 
	if data[9] != "" {
		tmp,err = strconv.ParseInt(data[9],10,64)
		if err != nil {
			err = errors.New("parse field attr_value_3_type 释放次数数值类型 to int32 failed")
			logger.ErrorWF("parse field attr_value_3_type 释放次数数值类型 to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[9]), 
				zap.Error(err))
			return
		}
		config.Attr_value_3_type = int32(tmp)
	}

	// parse column 10 attr_value_3 : 释放次数 
	if data[10] != "" {
		tmp,err = strconv.ParseInt(data[10],10,64)
		if err != nil {
			err = errors.New("parse field attr_value_3 释放次数 to int32 failed")
			logger.ErrorWF("parse field attr_value_3 释放次数 to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[10]), 
				zap.Error(err))
			return
		}
		config.Attr_value_3 = int32(tmp)
	}

	// parse column 11 attr_value_4 : 参数4 
	if data[11] != "" {
    
		vals := strings.Split(data[11],",")
		for k,v := range vals {
			tmp,err = strconv.ParseInt(v,10,64)
			if err != nil {
				err = errors.New("parse array field attr_value_4 参数4 to []int32 failed")
				logger.ErrorWF("parse array field attr_value_4 参数4 to []int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[11]), 
					zap.String("parse_data", v),zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Attr_value_4 = append(config.Attr_value_4, int32(tmp))
		}
	}

	// parse column 12 last_time_variable_id : 持续时长（毫秒）关联的变量id 
	if data[12] != "" {

		config.Last_time_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[12],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field last_time_variable_id 持续时长（毫秒）关联的变量id to key int32 failed")
				logger.ErrorWF("parse map field last_time_variable_id 持续时长（毫秒）关联的变量id to key int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[12]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field last_time_variable_id 持续时长（毫秒）关联的变量id to value int32 failed")
				logger.ErrorWF("parse map field last_time_variable_id 持续时长（毫秒）关联的变量id to value int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[12]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Last_time_variable_id[key] = value
		}
	}

	// parse column 13 last_time : 持续时长（毫秒） 
	if data[13] != "" {
		tmp,err = strconv.ParseInt(data[13],10,64)
		if err != nil {
			err = errors.New("parse field last_time 持续时长（毫秒） to int32 failed")
			logger.ErrorWF("parse field last_time 持续时长（毫秒） to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[13]), 
				zap.Error(err))
			return
		}
		config.Last_time = int32(tmp)
	}

	// parse column 14 base_hitrate_variable_id : 基础命中率变量 
	if data[14] != "" {

		config.Base_hitrate_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[14],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field base_hitrate_variable_id 基础命中率变量 to key int32 failed")
				logger.ErrorWF("parse map field base_hitrate_variable_id 基础命中率变量 to key int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[14]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field base_hitrate_variable_id 基础命中率变量 to value int32 failed")
				logger.ErrorWF("parse map field base_hitrate_variable_id 基础命中率变量 to value int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[14]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Base_hitrate_variable_id[key] = value
		}
	}

	// parse column 15 base_hitrate : 基础命中率（万分比） 
	if data[15] != "" {
		tmp,err = strconv.ParseInt(data[15],10,64)
		if err != nil {
			err = errors.New("parse field base_hitrate 基础命中率（万分比） to int32 failed")
			logger.ErrorWF("parse field base_hitrate 基础命中率（万分比） to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[15]), 
				zap.Error(err))
			return
		}
		config.Base_hitrate = int32(tmp)
	}

	// parse column 16 attr_value_7_variable_id : 参数7关联的变量id 
	if data[16] != "" {

		config.Attr_value_7_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[16],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field attr_value_7_variable_id 参数7关联的变量id to key int32 failed")
				logger.ErrorWF("parse map field attr_value_7_variable_id 参数7关联的变量id to key int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[16]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field attr_value_7_variable_id 参数7关联的变量id to value int32 failed")
				logger.ErrorWF("parse map field attr_value_7_variable_id 参数7关联的变量id to value int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[16]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Attr_value_7_variable_id[key] = value
		}
	}

	// parse column 17 attr_value_7_type : 参数7数值类型 
	if data[17] != "" {
		tmp,err = strconv.ParseInt(data[17],10,64)
		if err != nil {
			err = errors.New("parse field attr_value_7_type 参数7数值类型 to int32 failed")
			logger.ErrorWF("parse field attr_value_7_type 参数7数值类型 to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[17]), 
				zap.Error(err))
			return
		}
		config.Attr_value_7_type = int32(tmp)
	}

	// parse column 18 attr_value_7 : 参数7 
	if data[18] != "" {
		tmp,err = strconv.ParseInt(data[18],10,64)
		if err != nil {
			err = errors.New("parse field attr_value_7 参数7 to int32 failed")
			logger.ErrorWF("parse field attr_value_7 参数7 to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[18]), 
				zap.Error(err))
			return
		}
		config.Attr_value_7 = int32(tmp)
	}

	// parse column 19 attr_value_8_variable_id : 参数8关联的变量id 
	if data[19] != "" {

		config.Attr_value_8_variable_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[19],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field attr_value_8_variable_id 参数8关联的变量id to key int32 failed")
				logger.ErrorWF("parse map field attr_value_8_variable_id 参数8关联的变量id to key int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[19]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field attr_value_8_variable_id 参数8关联的变量id to value int32 failed")
				logger.ErrorWF("parse map field attr_value_8_variable_id 参数8关联的变量id to value int32 failed.", 
					zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
					// zap.String("field_data",data[19]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Attr_value_8_variable_id[key] = value
		}
	}

	// parse column 20 attr_value_8_type : 参数8数值类型 
	if data[20] != "" {
		tmp,err = strconv.ParseInt(data[20],10,64)
		if err != nil {
			err = errors.New("parse field attr_value_8_type 参数8数值类型 to int32 failed")
			logger.ErrorWF("parse field attr_value_8_type 参数8数值类型 to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[20]), 
				zap.Error(err))
			return
		}
		config.Attr_value_8_type = int32(tmp)
	}

	// parse column 21 attr_value_8 : 参数8 
	if data[21] != "" {
		tmp,err = strconv.ParseInt(data[21],10,64)
		if err != nil {
			err = errors.New("parse field attr_value_8 参数8 to int32 failed")
			logger.ErrorWF("parse field attr_value_8 参数8 to int32 failed.", 
				zap.String("xlsx", "maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx"), zap.String("sheet", "maze_skill_auto_release_v8"), 
				zap.String("parse_data",data[21]), 
				zap.Error(err))
			return
		}
		config.Attr_value_8 = int32(tmp)
	}
	return
}

var gMazeSkillAutoReleaseV8Fields = []string{
    "order",
    "attr",
    "attr_value_1_variable_id",
    "attr_value_1_type",
    "attr_value_1",
    "attr_value_2_variable_id",
    "attr_value_2_type",
    "attr_value_2",
    "attr_value_3_variable_id",
    "attr_value_3_type",
    "attr_value_3",
    "attr_value_4",
    "last_time_variable_id",
    "last_time",
    "base_hitrate_variable_id",
    "base_hitrate",
    "attr_value_7_variable_id",
    "attr_value_7_type",
    "attr_value_7",
    "attr_value_8_variable_id",
    "attr_value_8_type",
    "attr_value_8",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string,error)) (err error) {
	parser := &gMazeSkillAutoReleaseV8Parser{}
	loader := &gMazeSkillAutoReleaseV8Loader{}
	var data [][]string
	data,err = load("maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx", "maze_skill_auto_release_v8", gMazeSkillAutoReleaseV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeSkillAutoReleaseV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeSkillAutoReleaseV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_skill_auto_release_v8【迷宫-技能-自动释放技能】.xlsx maze_skill_auto_release_v8 data success.")
	return
}
