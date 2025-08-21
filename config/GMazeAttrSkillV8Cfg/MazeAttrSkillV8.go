package GMazeAttrSkillV8Cfg

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

// MazeAttrSkillV8ConfigRow from maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx maze_attr_skill_v8
type MazeAttrSkillV8ConfigRow struct {
	Attr_id       int32 `json:"attr_id"`       // 属性id
	Skill_id      int32 `json:"skill_id"`      // 技能id
	Auto_skill_id int32 `json:"auto_skill_id"` // 自动技能
}

// MazeAttrSkillV8Config from maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx maze_attr_skill_v8
type MazeAttrSkillV8Config struct {
	ConfigRows map[int32]*MazeAttrSkillV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeAttrSkillV8Config {
	ret := &MazeAttrSkillV8Config{ConfigRows: map[int32]*MazeAttrSkillV8ConfigRow{}}
	return ret
}

// GetMazeAttrSkillV8Config get one config by configId
func (c *MazeAttrSkillV8Config) GetMazeAttrSkillV8Config(configId int32) *MazeAttrSkillV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeAttrSkillV8Config) Get(configId int32) *MazeAttrSkillV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeAttrSkillV8Config get all config slice
func (c *MazeAttrSkillV8Config) GetAllMazeAttrSkillV8Config() (res []*MazeAttrSkillV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeAttrSkillV8Config) GetAll() (res []*MazeAttrSkillV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeAttrSkillV8Config

// GetMazeAttrSkillV8Config pkg func. get one config by configId
func GetMazeAttrSkillV8Config(configId int32) *MazeAttrSkillV8ConfigRow {
	return gConfigData.GetMazeAttrSkillV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeAttrSkillV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeAttrSkillV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_attr_skill_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeAttrSkillV8Config pkg func. get all config slice
func GetAllMazeAttrSkillV8Config() []*MazeAttrSkillV8ConfigRow {
	return gConfigData.GetAllMazeAttrSkillV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeAttrSkillV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeAttrSkillV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeAttrSkillV8ConfigRow from maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx maze_attr_skill_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeAttrSkillV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_attr_skill_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_attr_skill_v8.json",
		"maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx", "maze_attr_skill_v8",
		&gMazeAttrSkillV8Parser{}, &gMazeAttrSkillV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeAttrSkillV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeAttrSkillV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeAttrSkillV8Config))(c)
		return true
	})
}

// RegisterMazeAttrSkillV8InitCallBack reg config update func (old func)
var RegisterMazeAttrSkillV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeAttrSkillV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeAttrSkillV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeAttrSkillV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeAttrSkillV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeAttrSkillV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeAttrSkillV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeAttrSkillV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeAttrSkillV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeAttrSkillV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeAttrSkillV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeAttrSkillV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeAttrSkillV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeAttrSkillV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSkillV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAttrSkillV8ConfigRow", zap.String("xlsx", "maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx"),
			zap.String("sheet", "maze_attr_skill_v8"))
		return
	}
	config, ok := container.(*MazeAttrSkillV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSkillV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrSkillV8Config", zap.String("xlsx", "maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx"),
			zap.String("sheet", "maze_attr_skill_v8"))
		return
	}
	config.ConfigRows[row.Attr_id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeAttrSkillV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeAttrSkillV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSkillV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrSkillV8Config", zap.String("xlsx", "maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx"),
			zap.String("sheet", "maze_attr_skill_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeAttrSkillV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeAttrSkillV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSkillV8Config")
		logger.ErrorWF("invalid type. not *MazeAttrSkillV8Config", zap.String("xlsx", "maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx"),
			zap.String("sheet", "maze_attr_skill_v8"))
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
type gMazeAttrSkillV8Parser struct {
}

// New new config row data
func (*gMazeAttrSkillV8Parser) New() interface{} {
	return &MazeAttrSkillV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeAttrSkillV8Parser) Fields() []string {
	return gMazeAttrSkillV8Fields
}

// Parse parse raw data to row data
func (*gMazeAttrSkillV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeAttrSkillV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAttrSkillV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAttrSkillV8ConfigRow", zap.String("xlsx", "maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx"),
			zap.String("sheet", "maze_attr_skill_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeAttrSkillV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeAttrSkillV8ConfigRow",
			zap.String("xlsx", "maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx"),
			zap.String("sheet", "maze_attr_skill_v8"), zap.Int("need_count", len(gMazeAttrSkillV8Fields)),
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
				zap.String("xlsx", "maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx"), zap.String("sheet", "maze_attr_skill_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Attr_id = int32(tmp)
	}

	// parse column 1 skill_id : 技能id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field skill_id 技能id to int32 failed")
			logger.ErrorWF("parse field skill_id 技能id to int32 failed.",
				zap.String("xlsx", "maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx"), zap.String("sheet", "maze_attr_skill_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Skill_id = int32(tmp)
	}

	// parse column 2 auto_skill_id : 自动技能
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field auto_skill_id 自动技能 to int32 failed")
			logger.ErrorWF("parse field auto_skill_id 自动技能 to int32 failed.",
				zap.String("xlsx", "maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx"), zap.String("sheet", "maze_attr_skill_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Auto_skill_id = int32(tmp)
	}
	return
}

var gMazeAttrSkillV8Fields = []string{
	"attr_id",
	"skill_id",
	"auto_skill_id",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeAttrSkillV8Parser{}
	loader := &gMazeAttrSkillV8Loader{}
	var data [][]string
	data, err = load("maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx", "maze_attr_skill_v8", gMazeAttrSkillV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx maze_attr_skill_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx maze_attr_skill_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeAttrSkillV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx maze_attr_skill_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeAttrSkillV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx maze_attr_skill_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx maze_attr_skill_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx maze_attr_skill_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_attr_skill_v8【迷宫-属性id关联技能id】.xlsx maze_attr_skill_v8 data success.")
	return
}
