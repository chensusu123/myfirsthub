package GMazeEnergyAffixRandRuleV8Cfg


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


// MazeEnergyAffixRandRuleV8ConfigRow from maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx maze_energy_affix_rand_rule_v8
type MazeEnergyAffixRandRuleV8ConfigRow struct {
    Order       int32  `json:"order"` // 序号
    Rule_id       int32  `json:"rule_id"` // 规则id
    Energy_lv       int32  `json:"energy_lv"` // 能力等级
    Pos_1_lib       map[int32]int32  `json:"pos_1_lib"` // 位置1指定词条库id:随机权重
    Pos_2_lib       map[int32]int32  `json:"pos_2_lib"` // 位置2指定词条库id:随机权重
    Pos_3_lib       map[int32]int32  `json:"pos_3_lib"` // 位置3指定词条库id:随机权重
    Pos_4_lib       map[int32]int32  `json:"pos_4_lib"` // 位置4指定词条库id:随机权重
    Pos_5_lib       map[int32]int32  `json:"pos_5_lib"` // 位置5指定词条库id:随机权重
    Pos_6_lib       map[int32]int32  `json:"pos_6_lib"` // 位置6指定词条库id:随机权重
}

// MazeEnergyAffixRandRuleV8Config from maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx maze_energy_affix_rand_rule_v8
type MazeEnergyAffixRandRuleV8Config struct {
	ConfigRows map[int32]*MazeEnergyAffixRandRuleV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEnergyAffixRandRuleV8Config {
	ret := &MazeEnergyAffixRandRuleV8Config{ConfigRows: map[int32]*MazeEnergyAffixRandRuleV8ConfigRow{}}
	return ret
}

// GetMazeEnergyAffixRandRuleV8Config get one config by configId
func (c *MazeEnergyAffixRandRuleV8Config) GetMazeEnergyAffixRandRuleV8Config(configId int32) *MazeEnergyAffixRandRuleV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEnergyAffixRandRuleV8Config) Get(configId int32) *MazeEnergyAffixRandRuleV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEnergyAffixRandRuleV8Config get all config slice
func (c *MazeEnergyAffixRandRuleV8Config)  GetAllMazeEnergyAffixRandRuleV8Config () (res []*MazeEnergyAffixRandRuleV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEnergyAffixRandRuleV8Config)  GetAll() (res []*MazeEnergyAffixRandRuleV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}


// global config pointer 
var gConfigData *MazeEnergyAffixRandRuleV8Config 

// GetMazeEnergyAffixRandRuleV8Config pkg func. get one config by configId
func GetMazeEnergyAffixRandRuleV8Config(configId int32) *MazeEnergyAffixRandRuleV8ConfigRow {
	return gConfigData.GetMazeEnergyAffixRandRuleV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEnergyAffixRandRuleV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeEnergyAffixRandRuleV8Config pkg func. get all config slice
func GetAllMazeEnergyAffixRandRuleV8Config () []*MazeEnergyAffixRandRuleV8ConfigRow {
	return gConfigData.GetAllMazeEnergyAffixRandRuleV8Config ()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEnergyAffixRandRuleV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEnergyAffixRandRuleV8ConfigRow{
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string{
	return "MazeEnergyAffixRandRuleV8ConfigRow from maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx maze_energy_affix_rand_rule_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEnergyAffixRandRuleV8Config{
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_energy_affix_rand_rule_v8"
}

func init() {
	// reg config auto load 
	config_manager.RegAutoConfig("maze_energy_affix_rand_rule_v8.json", 
		"maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx", "maze_energy_affix_rand_rule_v8",
	 	&gMazeEnergyAffixRandRuleV8Parser{}, &gMazeEnergyAffixRandRuleV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEnergyAffixRandRuleV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEnergyAffixRandRuleV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEnergyAffixRandRuleV8Config))(c)
		return true
	})
}

// RegisterMazeEnergyAffixRandRuleV8InitCallBack reg config update func (old func)
var RegisterMazeEnergyAffixRandRuleV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEnergyAffixRandRuleV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEnergyAffixRandRuleV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEnergyAffixRandRuleV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEnergyAffixRandRuleV8Config)error)(c)
		return err == nil
	})
	return 
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEnergyAffixRandRuleV8Config)error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEnergyAffixRandRuleV8Loader struct {
}
// NewContainer new data container pointer
func (*gMazeEnergyAffixRandRuleV8Loader) NewContainer() interface{}{
	return newConfig()
}
// Check check new config data ptr
func (*gMazeEnergyAffixRandRuleV8Loader) Check(newPtr interface{})error{
	// set global ptr
	cfgData := newPtr.(*MazeEnergyAffixRandRuleV8Config)
	return doConfigCheckCallback(cfgData)
}
// Swap swap global config data ptr
func (*gMazeEnergyAffixRandRuleV8Loader) Swap(newPtr interface{}){
	// convert pointer
	cache := newPtr.(*MazeEnergyAffixRandRuleV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}
// Add append config item to map,save result
func (*gMazeEnergyAffixRandRuleV8Loader) Add(logger fklog.FKLogI,container interface{}, ri interface{})(err error) {
	row,ok := ri.(*MazeEnergyAffixRandRuleV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixRandRuleV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixRandRuleV8ConfigRow", zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"),
			zap.String("sheet", "maze_energy_affix_rand_rule_v8"))
		return 
	}
	config,ok := container.(*MazeEnergyAffixRandRuleV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixRandRuleV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixRandRuleV8Config", zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"),
			zap.String("sheet", "maze_energy_affix_rand_rule_v8"))
		return 
	}
	config.ConfigRows[row.Order] = row
	return
}
// GetValue get real map value for json parse
func (*gMazeEnergyAffixRandRuleV8Loader) GetValue(logger fklog.FKLogI,container interface{})(real interface{},err error) {
	config,ok := container.(*MazeEnergyAffixRandRuleV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixRandRuleV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixRandRuleV8Config", zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"),
			zap.String("sheet", "maze_energy_affix_rand_rule_v8"))
		return 
	}
	real = &config.ConfigRows
	return
}
// Range range all data for json append data parse.
func (*gMazeEnergyAffixRandRuleV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error){
	config,ok := container.(*MazeEnergyAffixRandRuleV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixRandRuleV8Config")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixRandRuleV8Config", zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"),
			zap.String("sheet", "maze_energy_affix_rand_rule_v8"))
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
type gMazeEnergyAffixRandRuleV8Parser struct {
}
// New new config row data
func (*gMazeEnergyAffixRandRuleV8Parser) New() interface{} {
	return &MazeEnergyAffixRandRuleV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEnergyAffixRandRuleV8Parser) Fields() []string {
	return gMazeEnergyAffixRandRuleV8Fields
}
// Parse parse raw data to row data
func (*gMazeEnergyAffixRandRuleV8Parser) Parse(logger fklog.FKLogI,data []string, row interface{}) (err error) {
	// convert row type
	config,ok := row.(*MazeEnergyAffixRandRuleV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEnergyAffixRandRuleV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEnergyAffixRandRuleV8ConfigRow", zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"),
			zap.String("sheet", "maze_energy_affix_rand_rule_v8"))
		return 
	}
	// compare length
	if len(data) != len(gMazeEnergyAffixRandRuleV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEnergyAffixRandRuleV8ConfigRow", 
			zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"),
			zap.String("sheet", "maze_energy_affix_rand_rule_v8"), zap.Int("need_count",len(gMazeEnergyAffixRandRuleV8Fields)), 
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
				zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
				zap.String("parse_data",data[0]), 
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 rule_id : 规则id 
	if data[1] != "" {
		tmp,err = strconv.ParseInt(data[1],10,64)
		if err != nil {
			err = errors.New("parse field rule_id 规则id to int32 failed")
			logger.ErrorWF("parse field rule_id 规则id to int32 failed.", 
				zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
				zap.String("parse_data",data[1]), 
				zap.Error(err))
			return
		}
		config.Rule_id = int32(tmp)
	}

	// parse column 2 energy_lv : 能力等级 
	if data[2] != "" {
		tmp,err = strconv.ParseInt(data[2],10,64)
		if err != nil {
			err = errors.New("parse field energy_lv 能力等级 to int32 failed")
			logger.ErrorWF("parse field energy_lv 能力等级 to int32 failed.", 
				zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
				zap.String("parse_data",data[2]), 
				zap.Error(err))
			return
		}
		config.Energy_lv = int32(tmp)
	}

	// parse column 3 pos_1_lib : 位置1指定词条库id:随机权重 
	if data[3] != "" {

		config.Pos_1_lib = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[3],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field pos_1_lib 位置1指定词条库id:随机权重 to key int32 failed")
				logger.ErrorWF("parse map field pos_1_lib 位置1指定词条库id:随机权重 to key int32 failed.", 
					zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
					// zap.String("field_data",data[3]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field pos_1_lib 位置1指定词条库id:随机权重 to value int32 failed")
				logger.ErrorWF("parse map field pos_1_lib 位置1指定词条库id:随机权重 to value int32 failed.", 
					zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
					// zap.String("field_data",data[3]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Pos_1_lib[key] = value
		}
	}

	// parse column 4 pos_2_lib : 位置2指定词条库id:随机权重 
	if data[4] != "" {

		config.Pos_2_lib = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[4],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field pos_2_lib 位置2指定词条库id:随机权重 to key int32 failed")
				logger.ErrorWF("parse map field pos_2_lib 位置2指定词条库id:随机权重 to key int32 failed.", 
					zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
					// zap.String("field_data",data[4]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field pos_2_lib 位置2指定词条库id:随机权重 to value int32 failed")
				logger.ErrorWF("parse map field pos_2_lib 位置2指定词条库id:随机权重 to value int32 failed.", 
					zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
					// zap.String("field_data",data[4]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Pos_2_lib[key] = value
		}
	}

	// parse column 5 pos_3_lib : 位置3指定词条库id:随机权重 
	if data[5] != "" {

		config.Pos_3_lib = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[5],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field pos_3_lib 位置3指定词条库id:随机权重 to key int32 failed")
				logger.ErrorWF("parse map field pos_3_lib 位置3指定词条库id:随机权重 to key int32 failed.", 
					zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
					// zap.String("field_data",data[5]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field pos_3_lib 位置3指定词条库id:随机权重 to value int32 failed")
				logger.ErrorWF("parse map field pos_3_lib 位置3指定词条库id:随机权重 to value int32 failed.", 
					zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
					// zap.String("field_data",data[5]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Pos_3_lib[key] = value
		}
	}

	// parse column 6 pos_4_lib : 位置4指定词条库id:随机权重 
	if data[6] != "" {

		config.Pos_4_lib = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[6],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field pos_4_lib 位置4指定词条库id:随机权重 to key int32 failed")
				logger.ErrorWF("parse map field pos_4_lib 位置4指定词条库id:随机权重 to key int32 failed.", 
					zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
					// zap.String("field_data",data[6]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field pos_4_lib 位置4指定词条库id:随机权重 to value int32 failed")
				logger.ErrorWF("parse map field pos_4_lib 位置4指定词条库id:随机权重 to value int32 failed.", 
					zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
					// zap.String("field_data",data[6]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Pos_4_lib[key] = value
		}
	}

	// parse column 7 pos_5_lib : 位置5指定词条库id:随机权重 
	if data[7] != "" {

		config.Pos_5_lib = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[7],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field pos_5_lib 位置5指定词条库id:随机权重 to key int32 failed")
				logger.ErrorWF("parse map field pos_5_lib 位置5指定词条库id:随机权重 to key int32 failed.", 
					zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
					// zap.String("field_data",data[7]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field pos_5_lib 位置5指定词条库id:随机权重 to value int32 failed")
				logger.ErrorWF("parse map field pos_5_lib 位置5指定词条库id:随机权重 to value int32 failed.", 
					zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
					// zap.String("field_data",data[7]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Pos_5_lib[key] = value
		}
	}

	// parse column 8 pos_6_lib : 位置6指定词条库id:随机权重 
	if data[8] != "" {

		config.Pos_6_lib = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[8],"_")
		for k,val := range vals {
			items := strings.Split(val,":")
			tmp,err = strconv.ParseInt(items[0],10,64)
			if err != nil {
				err = errors.New("parse map field pos_6_lib 位置6指定词条库id:随机权重 to key int32 failed")
				logger.ErrorWF("parse map field pos_6_lib 位置6指定词条库id:随机权重 to key int32 failed.", 
					zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
					// zap.String("field_data",data[8]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp,err = strconv.ParseInt(items[1],10,64)
			if err != nil {
				err = errors.New("parse map field pos_6_lib 位置6指定词条库id:随机权重 to value int32 failed")
				logger.ErrorWF("parse map field pos_6_lib 位置6指定词条库id:随机权重 to value int32 failed.", 
					zap.String("xlsx", "maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx"), zap.String("sheet", "maze_energy_affix_rand_rule_v8"), 
					// zap.String("field_data",data[8]), 
					zap.String("item_data",val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Pos_6_lib[key] = value
		}
	}
	return
}

var gMazeEnergyAffixRandRuleV8Fields = []string{
    "order",
    "rule_id",
    "energy_lv",
    "pos_1_lib",
    "pos_2_lib",
    "pos_3_lib",
    "pos_4_lib",
    "pos_5_lib",
    "pos_6_lib",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string,error)) (err error) {
	parser := &gMazeEnergyAffixRandRuleV8Parser{}
	loader := &gMazeEnergyAffixRandRuleV8Loader{}
	var data [][]string
	data,err = load("maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx", "maze_energy_affix_rand_rule_v8", gMazeEnergyAffixRandRuleV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx maze_energy_affix_rand_rule_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx maze_energy_affix_rand_rule_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEnergyAffixRandRuleV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx maze_energy_affix_rand_rule_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEnergyAffixRandRuleV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx maze_energy_affix_rand_rule_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx maze_energy_affix_rand_rule_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx maze_energy_affix_rand_rule_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_energy_affix_rand_rule_v8【迷宫-能力词条-随机规则】.xlsx maze_energy_affix_rand_rule_v8 data success.")
	return
}
