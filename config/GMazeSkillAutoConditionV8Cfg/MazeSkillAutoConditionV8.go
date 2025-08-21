package GMazeSkillAutoConditionV8Cfg

import (
	"context"
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

// MazeSkillAutoConditionV8ConfigRow from maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx maze_skill_auto_condition_v8
type MazeSkillAutoConditionV8ConfigRow struct {
	Order          int32           `json:"order"`          // 条件id(技能*10000+条件数*100+条件类型
	Type           int32           `json:"type"`           // 条件类型
	Symbol         int32           `json:"symbol"`         // 条件符号
	Value_variable map[int32]int32 `json:"value_variable"` // 参数变量id:变化值
	Value          int32           `json:"value"`          // 参数1
}

// MazeSkillAutoConditionV8Config from maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx maze_skill_auto_condition_v8
type MazeSkillAutoConditionV8Config struct {
	ConfigRows map[int32]*MazeSkillAutoConditionV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeSkillAutoConditionV8Config {
	ret := &MazeSkillAutoConditionV8Config{ConfigRows: map[int32]*MazeSkillAutoConditionV8ConfigRow{}}
	return ret
}

// GetMazeSkillAutoConditionV8Config get one config by configId
func (c *MazeSkillAutoConditionV8Config) GetMazeSkillAutoConditionV8Config(configId int32) *MazeSkillAutoConditionV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeSkillAutoConditionV8Config) Get(configId int32) *MazeSkillAutoConditionV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeSkillAutoConditionV8Config get all config slice
func (c *MazeSkillAutoConditionV8Config) GetAllMazeSkillAutoConditionV8Config() (res []*MazeSkillAutoConditionV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeSkillAutoConditionV8Config) GetAll() (res []*MazeSkillAutoConditionV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeSkillAutoConditionV8Config

// GetMazeSkillAutoConditionV8Config pkg func. get one config by configId
func GetMazeSkillAutoConditionV8Config(configId int32) *MazeSkillAutoConditionV8ConfigRow {
	return gConfigData.GetMazeSkillAutoConditionV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeSkillAutoConditionV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeSkillAutoConditionV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_skill_auto_condition_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeSkillAutoConditionV8Config pkg func. get all config slice
func GetAllMazeSkillAutoConditionV8Config() []*MazeSkillAutoConditionV8ConfigRow {
	return gConfigData.GetAllMazeSkillAutoConditionV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeSkillAutoConditionV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeSkillAutoConditionV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeSkillAutoConditionV8ConfigRow from maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx maze_skill_auto_condition_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeSkillAutoConditionV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_skill_auto_condition_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_skill_auto_condition_v8.json",
		"maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx", "maze_skill_auto_condition_v8",
		&gMazeSkillAutoConditionV8Parser{}, &gMazeSkillAutoConditionV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeSkillAutoConditionV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeSkillAutoConditionV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeSkillAutoConditionV8Config))(c)
		return true
	})
}

// RegisterMazeSkillAutoConditionV8InitCallBack reg config update func (old func)
var RegisterMazeSkillAutoConditionV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeSkillAutoConditionV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeSkillAutoConditionV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeSkillAutoConditionV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeSkillAutoConditionV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeSkillAutoConditionV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeSkillAutoConditionV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeSkillAutoConditionV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeSkillAutoConditionV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeSkillAutoConditionV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeSkillAutoConditionV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeSkillAutoConditionV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeSkillAutoConditionV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeSkillAutoConditionV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoConditionV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkillAutoConditionV8ConfigRow", zap.String("xlsx", "maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx"),
			zap.String("sheet", "maze_skill_auto_condition_v8"))
		return
	}
	config, ok := container.(*MazeSkillAutoConditionV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoConditionV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillAutoConditionV8Config", zap.String("xlsx", "maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx"),
			zap.String("sheet", "maze_skill_auto_condition_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeSkillAutoConditionV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeSkillAutoConditionV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoConditionV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillAutoConditionV8Config", zap.String("xlsx", "maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx"),
			zap.String("sheet", "maze_skill_auto_condition_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeSkillAutoConditionV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeSkillAutoConditionV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoConditionV8Config")
		logger.ErrorWF("invalid type. not *MazeSkillAutoConditionV8Config", zap.String("xlsx", "maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx"),
			zap.String("sheet", "maze_skill_auto_condition_v8"))
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
type gMazeSkillAutoConditionV8Parser struct {
}

// New new config row data
func (*gMazeSkillAutoConditionV8Parser) New() interface{} {
	return &MazeSkillAutoConditionV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeSkillAutoConditionV8Parser) Fields() []string {
	return gMazeSkillAutoConditionV8Fields
}

// Parse parse raw data to row data
func (*gMazeSkillAutoConditionV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeSkillAutoConditionV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeSkillAutoConditionV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeSkillAutoConditionV8ConfigRow", zap.String("xlsx", "maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx"),
			zap.String("sheet", "maze_skill_auto_condition_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeSkillAutoConditionV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeSkillAutoConditionV8ConfigRow",
			zap.String("xlsx", "maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx"),
			zap.String("sheet", "maze_skill_auto_condition_v8"), zap.Int("need_count", len(gMazeSkillAutoConditionV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 条件id(技能*10000+条件数*100+条件类型
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 条件id(技能*10000+条件数*100+条件类型 to int32 failed")
			logger.ErrorWF("parse field order 条件id(技能*10000+条件数*100+条件类型 to int32 failed.",
				zap.String("xlsx", "maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx"), zap.String("sheet", "maze_skill_auto_condition_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 type : 条件类型
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field type 条件类型 to int32 failed")
			logger.ErrorWF("parse field type 条件类型 to int32 failed.",
				zap.String("xlsx", "maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx"), zap.String("sheet", "maze_skill_auto_condition_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Type = int32(tmp)
	}

	// parse column 2 symbol : 条件符号
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field symbol 条件符号 to int32 failed")
			logger.ErrorWF("parse field symbol 条件符号 to int32 failed.",
				zap.String("xlsx", "maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx"), zap.String("sheet", "maze_skill_auto_condition_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Symbol = int32(tmp)
	}

	// parse column 3 value_variable : 参数变量id:变化值
	if data[3] != "" {

		config.Value_variable = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field value_variable 参数变量id:变化值 to key int32 failed")
				logger.ErrorWF("parse map field value_variable 参数变量id:变化值 to key int32 failed.",
					zap.String("xlsx", "maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx"), zap.String("sheet", "maze_skill_auto_condition_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field value_variable 参数变量id:变化值 to value int32 failed")
				logger.ErrorWF("parse map field value_variable 参数变量id:变化值 to value int32 failed.",
					zap.String("xlsx", "maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx"), zap.String("sheet", "maze_skill_auto_condition_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Value_variable[key] = value
		}
	}

	// parse column 4 value : 参数1
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field value 参数1 to int32 failed")
			logger.ErrorWF("parse field value 参数1 to int32 failed.",
				zap.String("xlsx", "maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx"), zap.String("sheet", "maze_skill_auto_condition_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Value = int32(tmp)
	}
	return
}

var gMazeSkillAutoConditionV8Fields = []string{
	"order",
	"type",
	"symbol",
	"value_variable",
	"value",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeSkillAutoConditionV8Parser{}
	loader := &gMazeSkillAutoConditionV8Loader{}
	var data [][]string
	data, err = load("maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx", "maze_skill_auto_condition_v8", gMazeSkillAutoConditionV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx maze_skill_auto_condition_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx maze_skill_auto_condition_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeSkillAutoConditionV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx maze_skill_auto_condition_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeSkillAutoConditionV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx maze_skill_auto_condition_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx maze_skill_auto_condition_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx maze_skill_auto_condition_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_skill_auto_condition_v8【迷宫-技能-自动释放条件】.xlsx maze_skill_auto_condition_v8 data success.")
	return
}
