package GMazeAttrSkillEffectV8Cfg

import (
	"errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"
)

// MazeAttrSkillEffectV8ConfigRow from maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx maze_attr_skill_effect_v8
type MazeAttrSkillEffectV8ConfigRow struct {
	Attr_id         int32 `json:"attr_id"`         // 属性id
	Skill_effect_id int32 `json:"skill_effect_id"` // 技能效果id
}

// MazeAttrSkillEffectV8Config from maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx maze_attr_skill_effect_v8
type MazeAttrSkillEffectV8Config struct {
	ConfigRows map[int32]*MazeAttrSkillEffectV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeAttrSkillEffectV8Config {
	ret := &MazeAttrSkillEffectV8Config{ConfigRows: map[int32]*MazeAttrSkillEffectV8ConfigRow{}}
	return ret
}

// GetMazeAttrSkillEffectV8Config get one config by configId
func (c *MazeAttrSkillEffectV8Config) GetMazeAttrSkillEffectV8Config(configId int32) *MazeAttrSkillEffectV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeAttrSkillEffectV8Config) Get(configId int32) *MazeAttrSkillEffectV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeAttrSkillEffectV8Config get all config slice
func (c *MazeAttrSkillEffectV8Config) GetAllMazeAttrSkillEffectV8Config() (res []*MazeAttrSkillEffectV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeAttrSkillEffectV8Config) GetAll() (res []*MazeAttrSkillEffectV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeAttrSkillEffectV8Config

// GetMazeAttrSkillEffectV8Config pkg func. get one config by configId
func GetMazeAttrSkillEffectV8Config(configId int32) *MazeAttrSkillEffectV8ConfigRow {
	return gConfigData.GetMazeAttrSkillEffectV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeAttrSkillEffectV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeAttrSkillEffectV8Config pkg func. get all config slice
func GetAllMazeAttrSkillEffectV8Config() []*MazeAttrSkillEffectV8ConfigRow {
	return gConfigData.GetAllMazeAttrSkillEffectV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeAttrSkillEffectV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeAttrSkillEffectV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeAttrSkillEffectV8ConfigRow from maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx maze_attr_skill_effect_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeAttrSkillEffectV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_attr_skill_effect_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_attr_skill_effect_v8.json",
		"maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx", "maze_attr_skill_effect_v8",
		&gMazeAttrSkillEffectV8Parser{}, &gMazeAttrSkillEffectV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeAttrSkillEffectV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeAttrSkillEffectV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeAttrSkillEffectV8Config))(c)
		return true
	})
}

// RegisterMazeAttrSkillEffectV8InitCallBack reg config update func (old func)
var RegisterMazeAttrSkillEffectV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeAttrSkillEffectV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeAttrSkillEffectV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeAttrSkillEffectV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeAttrSkillEffectV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeAttrSkillEffectV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeAttrSkillEffectV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeAttrSkillEffectV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeAttrSkillEffectV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeAttrSkillEffectV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeAttrSkillEffectV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeAttrSkillEffectV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeAttrSkillEffectV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeAttrSkillEffectV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSkillEffectV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAttrSkillEffectV8ConfigRow", zap.String("xlsx", "maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx"),
			zap.String("sheet", "maze_attr_skill_effect_v8"))
		return
	}
	config, ok := container.(*MazeAttrSkillEffectV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSkillEffectV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrSkillEffectV8Config", zap.String("xlsx", "maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx"),
			zap.String("sheet", "maze_attr_skill_effect_v8"))
		return
	}
	config.ConfigRows[row.Attr_id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeAttrSkillEffectV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeAttrSkillEffectV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSkillEffectV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrSkillEffectV8Config", zap.String("xlsx", "maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx"),
			zap.String("sheet", "maze_attr_skill_effect_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeAttrSkillEffectV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeAttrSkillEffectV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSkillEffectV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrSkillEffectV8Config", zap.String("xlsx", "maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx"),
			zap.String("sheet", "maze_attr_skill_effect_v8"))
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
type gMazeAttrSkillEffectV8Parser struct {
}

// New new config row data
func (*gMazeAttrSkillEffectV8Parser) New() interface{} {
	return &MazeAttrSkillEffectV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeAttrSkillEffectV8Parser) Fields() []string {
	return gMazeAttrSkillEffectV8Fields
}

// Parse parse raw data to row data
func (*gMazeAttrSkillEffectV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeAttrSkillEffectV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSkillEffectV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAttrSkillEffectV8ConfigRow", zap.String("xlsx", "maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx"),
			zap.String("sheet", "maze_attr_skill_effect_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeAttrSkillEffectV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeAttrSkillEffectV8ConfigRow",
			zap.String("xlsx", "maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx"),
			zap.String("sheet", "maze_attr_skill_effect_v8"), zap.Int("need_count", len(gMazeAttrSkillEffectV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 attr_id : 属性id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field attr_id 属性id to int32 failed")
			logger.ErrorWF("parse field attr_id 属性id to int32 failed.",
				zap.String("xlsx", "maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx"), zap.String("sheet", "maze_attr_skill_effect_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Attr_id = int32(tmp)
	}

	// parse column 1 skill_effect_id : 技能效果id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field skill_effect_id 技能效果id to int32 failed")
			logger.ErrorWF("parse field skill_effect_id 技能效果id to int32 failed.",
				zap.String("xlsx", "maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx"), zap.String("sheet", "maze_attr_skill_effect_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Skill_effect_id = int32(tmp)
	}
	return
}

var gMazeAttrSkillEffectV8Fields = []string{
	"attr_id",
	"skill_effect_id",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeAttrSkillEffectV8Parser{}
	loader := &gMazeAttrSkillEffectV8Loader{}
	var data [][]string
	data, err = load("maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx", "maze_attr_skill_effect_v8", gMazeAttrSkillEffectV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx maze_attr_skill_effect_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx maze_attr_skill_effect_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeAttrSkillEffectV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx maze_attr_skill_effect_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeAttrSkillEffectV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx maze_attr_skill_effect_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx maze_attr_skill_effect_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx maze_attr_skill_effect_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_attr_skill_effect_v8【迷宫-属性-对应技能效果】.xlsx maze_attr_skill_effect_v8 data success.")
	return
}
