package GMazeSkillActV8Cfg

import (
	"errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

// MazeSkillActV8ConfigRow from maze_skill_act_v8【迷宫-技能-技能动作】.xlsx maze_skill_act_v8
type MazeSkillActV8ConfigRow struct {
	Skill_id  int32   `json:"skill_id"`  // 技能id
	Act_id    []int32 `json:"act_id"`    // 技能对应的动作id
	Effect_id int32   `json:"effect_id"` // 技能对应的特效组
}

// MazeSkillActV8Config from maze_skill_act_v8【迷宫-技能-技能动作】.xlsx maze_skill_act_v8
type MazeSkillActV8Config struct {
	ConfigRows map[int32]*MazeSkillActV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeSkillActV8Config {
	ret := &MazeSkillActV8Config{ConfigRows: map[int32]*MazeSkillActV8ConfigRow{}}
	return ret
}

// GetMazeSkillActV8Config get one config by configId
func (c *MazeSkillActV8Config) GetMazeSkillActV8Config(configId int32) *MazeSkillActV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeSkillActV8Config) Get(configId int32) *MazeSkillActV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeSkillActV8Config get all config slice
func (c *MazeSkillActV8Config) GetAllMazeSkillActV8Config() (res []*MazeSkillActV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeSkillActV8Config) GetAll() (res []*MazeSkillActV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeSkillActV8Config

// GetMazeSkillActV8Config pkg func. get one config by configId
func GetMazeSkillActV8Config(configId int32) *MazeSkillActV8ConfigRow {
	return gConfigData.GetMazeSkillActV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeSkillActV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeSkillActV8Config pkg func. get all config slice
func GetAllMazeSkillActV8Config() []*MazeSkillActV8ConfigRow {
	return gConfigData.GetAllMazeSkillActV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeSkillActV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeSkillActV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeSkillActV8ConfigRow from maze_skill_act_v8【迷宫-技能-技能动作】.xlsx maze_skill_act_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeSkillActV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_skill_act_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_skill_act_v8.json",
		"maze_skill_act_v8【迷宫-技能-技能动作】.xlsx", "maze_skill_act_v8",
		&gMazeSkillActV8Parser{}, &gMazeSkillActV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeSkillActV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeSkillActV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeSkillActV8Config))(c)
		return true
	})
}

// RegisterMazeSkillActV8InitCallBack reg config update func (old func)
var RegisterMazeSkillActV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeSkillActV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeSkillActV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeSkillActV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeSkillActV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeSkillActV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeSkillActV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeSkillActV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeSkillActV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeSkillActV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeSkillActV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeSkillActV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeSkillActV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeSkillActV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillActV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkillActV8ConfigRow", zap.String("xlsx", "maze_skill_act_v8【迷宫-技能-技能动作】.xlsx"),
			zap.String("sheet", "maze_skill_act_v8"))
		return
	}
	config, ok := container.(*MazeSkillActV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillActV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillActV8Config", zap.String("xlsx", "maze_skill_act_v8【迷宫-技能-技能动作】.xlsx"),
			zap.String("sheet", "maze_skill_act_v8"))
		return
	}
	config.ConfigRows[row.Skill_id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeSkillActV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeSkillActV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillActV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillActV8Config", zap.String("xlsx", "maze_skill_act_v8【迷宫-技能-技能动作】.xlsx"),
			zap.String("sheet", "maze_skill_act_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeSkillActV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeSkillActV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillActV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillActV8Config", zap.String("xlsx", "maze_skill_act_v8【迷宫-技能-技能动作】.xlsx"),
			zap.String("sheet", "maze_skill_act_v8"))
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
type gMazeSkillActV8Parser struct {
}

// New new config row data
func (*gMazeSkillActV8Parser) New() interface{} {
	return &MazeSkillActV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeSkillActV8Parser) Fields() []string {
	return gMazeSkillActV8Fields
}

// Parse parse raw data to row data
func (*gMazeSkillActV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeSkillActV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillActV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkillActV8ConfigRow", zap.String("xlsx", "maze_skill_act_v8【迷宫-技能-技能动作】.xlsx"),
			zap.String("sheet", "maze_skill_act_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeSkillActV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeSkillActV8ConfigRow",
			zap.String("xlsx", "maze_skill_act_v8【迷宫-技能-技能动作】.xlsx"),
			zap.String("sheet", "maze_skill_act_v8"), zap.Int("need_count", len(gMazeSkillActV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 skill_id : 技能id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field skill_id 技能id to int32 failed")
			logger.ErrorWF("parse field skill_id 技能id to int32 failed.",
				zap.String("xlsx", "maze_skill_act_v8【迷宫-技能-技能动作】.xlsx"), zap.String("sheet", "maze_skill_act_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Skill_id = int32(tmp)
	}

	// parse column 1 act_id : 技能对应的动作id
	if data[1] != "" {

		vals := strings.Split(data[1], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field act_id 技能对应的动作id to []int32 failed")
				logger.ErrorWF("parse array field act_id 技能对应的动作id to []int32 failed.",
					zap.String("xlsx", "maze_skill_act_v8【迷宫-技能-技能动作】.xlsx"), zap.String("sheet", "maze_skill_act_v8"),
					// zap.String("field_data",data[1]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Act_id = append(config.Act_id, int32(tmp))
		}
	}

	// parse column 2 effect_id : 技能对应的特效组
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field effect_id 技能对应的特效组 to int32 failed")
			logger.ErrorWF("parse field effect_id 技能对应的特效组 to int32 failed.",
				zap.String("xlsx", "maze_skill_act_v8【迷宫-技能-技能动作】.xlsx"), zap.String("sheet", "maze_skill_act_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Effect_id = int32(tmp)
	}
	return
}

var gMazeSkillActV8Fields = []string{
	"skill_id",
	"act_id",
	"effect_id",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeSkillActV8Parser{}
	loader := &gMazeSkillActV8Loader{}
	var data [][]string
	data, err = load("maze_skill_act_v8【迷宫-技能-技能动作】.xlsx", "maze_skill_act_v8", gMazeSkillActV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_skill_act_v8【迷宫-技能-技能动作】.xlsx maze_skill_act_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_skill_act_v8【迷宫-技能-技能动作】.xlsx maze_skill_act_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeSkillActV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_skill_act_v8【迷宫-技能-技能动作】.xlsx maze_skill_act_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeSkillActV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_skill_act_v8【迷宫-技能-技能动作】.xlsx maze_skill_act_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_skill_act_v8【迷宫-技能-技能动作】.xlsx maze_skill_act_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_skill_act_v8【迷宫-技能-技能动作】.xlsx maze_skill_act_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_skill_act_v8【迷宫-技能-技能动作】.xlsx maze_skill_act_v8 data success.")
	return
}
