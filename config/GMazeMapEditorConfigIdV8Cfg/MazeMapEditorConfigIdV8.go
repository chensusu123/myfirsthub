package GMazeMapEditorConfigIdV8Cfg

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

// MazeMapEditorConfigIdV8ConfigRow from maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx maze_map_editor_config_id_v8
type MazeMapEditorConfigIdV8ConfigRow struct {
	ID            int32 `json:"ID"`            // 唯一ID
	Level_id      int32 `json:"level_id"`      // 关卡ID
	Config_id     int32 `json:"config_id"`     // 地编配置id
	Add_kungfu    int32 `json:"add_kungfu"`    // 增加通关值
	Add_save_item int32 `json:"add_save_item"` // 增加道具id
	Save_item_pos int32 `json:"save_item_pos"` // 掉落道具位置（万分比）
}

// MazeMapEditorConfigIdV8Config from maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx maze_map_editor_config_id_v8
type MazeMapEditorConfigIdV8Config struct {
	ConfigRows map[int32]*MazeMapEditorConfigIdV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeMapEditorConfigIdV8Config {
	ret := &MazeMapEditorConfigIdV8Config{ConfigRows: map[int32]*MazeMapEditorConfigIdV8ConfigRow{}}
	return ret
}

// GetMazeMapEditorConfigIdV8Config get one config by configId
func (c *MazeMapEditorConfigIdV8Config) GetMazeMapEditorConfigIdV8Config(configId int32) *MazeMapEditorConfigIdV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeMapEditorConfigIdV8Config) Get(configId int32) *MazeMapEditorConfigIdV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeMapEditorConfigIdV8Config get all config slice
func (c *MazeMapEditorConfigIdV8Config) GetAllMazeMapEditorConfigIdV8Config() (res []*MazeMapEditorConfigIdV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeMapEditorConfigIdV8Config) GetAll() (res []*MazeMapEditorConfigIdV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeMapEditorConfigIdV8Config

// GetMazeMapEditorConfigIdV8Config pkg func. get one config by configId
func GetMazeMapEditorConfigIdV8Config(configId int32) *MazeMapEditorConfigIdV8ConfigRow {
	return gConfigData.GetMazeMapEditorConfigIdV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeMapEditorConfigIdV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeMapEditorConfigIdV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_map_editor_config_id_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeMapEditorConfigIdV8Config pkg func. get all config slice
func GetAllMazeMapEditorConfigIdV8Config() []*MazeMapEditorConfigIdV8ConfigRow {
	return gConfigData.GetAllMazeMapEditorConfigIdV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeMapEditorConfigIdV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeMapEditorConfigIdV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeMapEditorConfigIdV8ConfigRow from maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx maze_map_editor_config_id_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeMapEditorConfigIdV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_map_editor_config_id_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_map_editor_config_id_v8.json",
		"maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx", "maze_map_editor_config_id_v8",
		&gMazeMapEditorConfigIdV8Parser{}, &gMazeMapEditorConfigIdV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeMapEditorConfigIdV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeMapEditorConfigIdV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeMapEditorConfigIdV8Config))(c)
		return true
	})
}

// RegisterMazeMapEditorConfigIdV8InitCallBack reg config update func (old func)
var RegisterMazeMapEditorConfigIdV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeMapEditorConfigIdV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeMapEditorConfigIdV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeMapEditorConfigIdV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeMapEditorConfigIdV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeMapEditorConfigIdV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeMapEditorConfigIdV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeMapEditorConfigIdV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeMapEditorConfigIdV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeMapEditorConfigIdV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeMapEditorConfigIdV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeMapEditorConfigIdV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeMapEditorConfigIdV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeMapEditorConfigIdV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeMapEditorConfigIdV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeMapEditorConfigIdV8ConfigRow", zap.String("xlsx", "maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx"),
			zap.String("sheet", "maze_map_editor_config_id_v8"))
		return
	}
	config, ok := container.(*MazeMapEditorConfigIdV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeMapEditorConfigIdV8Config")
		logger.ErrorWF("invalid type. not *MazeMapEditorConfigIdV8Config", zap.String("xlsx", "maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx"),
			zap.String("sheet", "maze_map_editor_config_id_v8"))
		return
	}
	config.ConfigRows[row.ID] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeMapEditorConfigIdV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeMapEditorConfigIdV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeMapEditorConfigIdV8Config")
		logger.ErrorWF("invalid type. not *MazeMapEditorConfigIdV8Config", zap.String("xlsx", "maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx"),
			zap.String("sheet", "maze_map_editor_config_id_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeMapEditorConfigIdV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeMapEditorConfigIdV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeMapEditorConfigIdV8Config")
		logger.ErrorWF("invalid type. not *MazeMapEditorConfigIdV8Config", zap.String("xlsx", "maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx"),
			zap.String("sheet", "maze_map_editor_config_id_v8"))
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
type gMazeMapEditorConfigIdV8Parser struct {
}

// New new config row data
func (*gMazeMapEditorConfigIdV8Parser) New() interface{} {
	return &MazeMapEditorConfigIdV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeMapEditorConfigIdV8Parser) Fields() []string {
	return gMazeMapEditorConfigIdV8Fields
}

// Parse parse raw data to row data
func (*gMazeMapEditorConfigIdV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeMapEditorConfigIdV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeMapEditorConfigIdV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeMapEditorConfigIdV8ConfigRow", zap.String("xlsx", "maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx"),
			zap.String("sheet", "maze_map_editor_config_id_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeMapEditorConfigIdV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeMapEditorConfigIdV8ConfigRow",
			zap.String("xlsx", "maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx"),
			zap.String("sheet", "maze_map_editor_config_id_v8"), zap.Int("need_count", len(gMazeMapEditorConfigIdV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 ID : 唯一ID
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field ID 唯一ID to int32 failed")
			logger.ErrorWF("parse field ID 唯一ID to int32 failed.",
				zap.String("xlsx", "maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx"), zap.String("sheet", "maze_map_editor_config_id_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.ID = int32(tmp)
	}

	// parse column 1 level_id : 关卡ID
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field level_id 关卡ID to int32 failed")
			logger.ErrorWF("parse field level_id 关卡ID to int32 failed.",
				zap.String("xlsx", "maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx"), zap.String("sheet", "maze_map_editor_config_id_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Level_id = int32(tmp)
	}

	// parse column 2 config_id : 地编配置id
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field config_id 地编配置id to int32 failed")
			logger.ErrorWF("parse field config_id 地编配置id to int32 failed.",
				zap.String("xlsx", "maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx"), zap.String("sheet", "maze_map_editor_config_id_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Config_id = int32(tmp)
	}

	// parse column 3 add_kungfu : 增加通关值
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field add_kungfu 增加通关值 to int32 failed")
			logger.ErrorWF("parse field add_kungfu 增加通关值 to int32 failed.",
				zap.String("xlsx", "maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx"), zap.String("sheet", "maze_map_editor_config_id_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Add_kungfu = int32(tmp)
	}

	// parse column 4 add_save_item : 增加道具id
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field add_save_item 增加道具id to int32 failed")
			logger.ErrorWF("parse field add_save_item 增加道具id to int32 failed.",
				zap.String("xlsx", "maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx"), zap.String("sheet", "maze_map_editor_config_id_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Add_save_item = int32(tmp)
	}

	// parse column 5 save_item_pos : 掉落道具位置（万分比）
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field save_item_pos 掉落道具位置（万分比） to int32 failed")
			logger.ErrorWF("parse field save_item_pos 掉落道具位置（万分比） to int32 failed.",
				zap.String("xlsx", "maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx"), zap.String("sheet", "maze_map_editor_config_id_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Save_item_pos = int32(tmp)
	}
	return
}

var gMazeMapEditorConfigIdV8Fields = []string{
	"ID",
	"level_id",
	"config_id",
	"add_kungfu",
	"add_save_item",
	"save_item_pos",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeMapEditorConfigIdV8Parser{}
	loader := &gMazeMapEditorConfigIdV8Loader{}
	var data [][]string
	data, err = load("maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx", "maze_map_editor_config_id_v8", gMazeMapEditorConfigIdV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx maze_map_editor_config_id_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx maze_map_editor_config_id_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeMapEditorConfigIdV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx maze_map_editor_config_id_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeMapEditorConfigIdV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx maze_map_editor_config_id_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx maze_map_editor_config_id_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx maze_map_editor_config_id_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_map_editor_config_id_v8【迷宫-地编配置id】.xlsx maze_map_editor_config_id_v8 data success.")
	return
}
