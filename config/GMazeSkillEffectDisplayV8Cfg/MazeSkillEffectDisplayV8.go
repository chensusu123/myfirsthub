package GMazeSkillEffectDisplayV8Cfg

import (
	"context"
	"errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"
)

// MazeSkillEffectDisplayV8ConfigRow from maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx maze_skill_effect_display_v8
type MazeSkillEffectDisplayV8ConfigRow struct {
	Effect_group   int32  `json:"effect_group"`   // effect组id
	Buff_list_name string `json:"buff_list_name"` // buff列表里的效果名称
}

// MazeSkillEffectDisplayV8Config from maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx maze_skill_effect_display_v8
type MazeSkillEffectDisplayV8Config struct {
	ConfigRows map[int32]*MazeSkillEffectDisplayV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeSkillEffectDisplayV8Config {
	ret := &MazeSkillEffectDisplayV8Config{ConfigRows: map[int32]*MazeSkillEffectDisplayV8ConfigRow{}}
	return ret
}

// GetMazeSkillEffectDisplayV8Config get one config by configId
func (c *MazeSkillEffectDisplayV8Config) GetMazeSkillEffectDisplayV8Config(configId int32) *MazeSkillEffectDisplayV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeSkillEffectDisplayV8Config) Get(configId int32) *MazeSkillEffectDisplayV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeSkillEffectDisplayV8Config get all config slice
func (c *MazeSkillEffectDisplayV8Config) GetAllMazeSkillEffectDisplayV8Config() (res []*MazeSkillEffectDisplayV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeSkillEffectDisplayV8Config) GetAll() (res []*MazeSkillEffectDisplayV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeSkillEffectDisplayV8Config

// GetMazeSkillEffectDisplayV8Config pkg func. get one config by configId
func GetMazeSkillEffectDisplayV8Config(configId int32) *MazeSkillEffectDisplayV8ConfigRow {
	return gConfigData.GetMazeSkillEffectDisplayV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeSkillEffectDisplayV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeSkillEffectDisplayV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_skill_effect_display_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeSkillEffectDisplayV8Config pkg func. get all config slice
func GetAllMazeSkillEffectDisplayV8Config() []*MazeSkillEffectDisplayV8ConfigRow {
	return gConfigData.GetAllMazeSkillEffectDisplayV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeSkillEffectDisplayV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeSkillEffectDisplayV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeSkillEffectDisplayV8ConfigRow from maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx maze_skill_effect_display_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeSkillEffectDisplayV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_skill_effect_display_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_skill_effect_display_v8.json",
		"maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx", "maze_skill_effect_display_v8",
		&gMazeSkillEffectDisplayV8Parser{}, &gMazeSkillEffectDisplayV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeSkillEffectDisplayV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeSkillEffectDisplayV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeSkillEffectDisplayV8Config))(c)
		return true
	})
}

// RegisterMazeSkillEffectDisplayV8InitCallBack reg config update func (old func)
var RegisterMazeSkillEffectDisplayV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeSkillEffectDisplayV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeSkillEffectDisplayV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeSkillEffectDisplayV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeSkillEffectDisplayV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeSkillEffectDisplayV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeSkillEffectDisplayV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeSkillEffectDisplayV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeSkillEffectDisplayV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeSkillEffectDisplayV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeSkillEffectDisplayV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeSkillEffectDisplayV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeSkillEffectDisplayV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeSkillEffectDisplayV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillEffectDisplayV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkillEffectDisplayV8ConfigRow", zap.String("xlsx", "maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx"),
			zap.String("sheet", "maze_skill_effect_display_v8"))
		return
	}
	config, ok := container.(*MazeSkillEffectDisplayV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillEffectDisplayV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillEffectDisplayV8Config", zap.String("xlsx", "maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx"),
			zap.String("sheet", "maze_skill_effect_display_v8"))
		return
	}
	config.ConfigRows[row.Effect_group] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeSkillEffectDisplayV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeSkillEffectDisplayV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillEffectDisplayV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillEffectDisplayV8Config", zap.String("xlsx", "maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx"),
			zap.String("sheet", "maze_skill_effect_display_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeSkillEffectDisplayV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeSkillEffectDisplayV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillEffectDisplayV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillEffectDisplayV8Config", zap.String("xlsx", "maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx"),
			zap.String("sheet", "maze_skill_effect_display_v8"))
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
type gMazeSkillEffectDisplayV8Parser struct {
}

// New new config row data
func (*gMazeSkillEffectDisplayV8Parser) New() interface{} {
	return &MazeSkillEffectDisplayV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeSkillEffectDisplayV8Parser) Fields() []string {
	return gMazeSkillEffectDisplayV8Fields
}

// Parse parse raw data to row data
func (*gMazeSkillEffectDisplayV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeSkillEffectDisplayV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillEffectDisplayV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkillEffectDisplayV8ConfigRow", zap.String("xlsx", "maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx"),
			zap.String("sheet", "maze_skill_effect_display_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeSkillEffectDisplayV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeSkillEffectDisplayV8ConfigRow",
			zap.String("xlsx", "maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx"),
			zap.String("sheet", "maze_skill_effect_display_v8"), zap.Int("need_count", len(gMazeSkillEffectDisplayV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 effect_group : effect组id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field effect_group effect组id to int32 failed")
			logger.ErrorWF("parse field effect_group effect组id to int32 failed.",
				zap.String("xlsx", "maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx"), zap.String("sheet", "maze_skill_effect_display_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Effect_group = int32(tmp)
	}

	// parse column 1 buff_list_name : buff列表里的效果名称
	if data[1] != "" {
		config.Buff_list_name = data[1]
	}
	return
}

var gMazeSkillEffectDisplayV8Fields = []string{
	"effect_group",
	"buff_list_name",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeSkillEffectDisplayV8Parser{}
	loader := &gMazeSkillEffectDisplayV8Loader{}
	var data [][]string
	data, err = load("maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx", "maze_skill_effect_display_v8", gMazeSkillEffectDisplayV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx maze_skill_effect_display_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx maze_skill_effect_display_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeSkillEffectDisplayV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx maze_skill_effect_display_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeSkillEffectDisplayV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx maze_skill_effect_display_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx maze_skill_effect_display_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx maze_skill_effect_display_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_skill_effect_display_v8【迷宫-技能-效果展示信息】.xlsx maze_skill_effect_display_v8 data success.")
	return
}
