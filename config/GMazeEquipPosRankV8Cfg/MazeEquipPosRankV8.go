package GMazeEquipPosRankV8Cfg

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

// MazeEquipPosRankV8ConfigRow from maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx maze_equip_pos_rank_v8
type MazeEquipPosRankV8ConfigRow struct {
	Pos_id          int32            `json:"pos_id"`          // 部位id
	Name            string           `json:"name"`            // 备注
	Sub_type_name   map[int32]string `json:"sub_type_name"`   // 子类型
	Rank            int32            `json:"rank"`            // 展示顺序
	Is_default      int32            `json:"is_default"`      // 是否默认解锁
	Need_level      int32            `json:"need_level"`      // 需要人偶等级
	Need_task       int32            `json:"need_task"`       // 需要完成任务id
	Need_dungeon_id map[int32]int32  `json:"need_dungeon_id"` // 需要通关章节id（=）：层数id(>=)
	Unlock_desc     string           `json:"unlock_desc"`     // 未解锁描述
}

// MazeEquipPosRankV8Config from maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx maze_equip_pos_rank_v8
type MazeEquipPosRankV8Config struct {
	ConfigRows map[int32]*MazeEquipPosRankV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipPosRankV8Config {
	ret := &MazeEquipPosRankV8Config{ConfigRows: map[int32]*MazeEquipPosRankV8ConfigRow{}}
	return ret
}

// GetMazeEquipPosRankV8Config get one config by configId
func (c *MazeEquipPosRankV8Config) GetMazeEquipPosRankV8Config(configId int32) *MazeEquipPosRankV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipPosRankV8Config) Get(configId int32) *MazeEquipPosRankV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipPosRankV8Config get all config slice
func (c *MazeEquipPosRankV8Config) GetAllMazeEquipPosRankV8Config() (res []*MazeEquipPosRankV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipPosRankV8Config) GetAll() (res []*MazeEquipPosRankV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipPosRankV8Config

// GetMazeEquipPosRankV8Config pkg func. get one config by configId
func GetMazeEquipPosRankV8Config(configId int32) *MazeEquipPosRankV8ConfigRow {
	return gConfigData.GetMazeEquipPosRankV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipPosRankV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeEquipPosRankV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_equip_pos_rank_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeEquipPosRankV8Config pkg func. get all config slice
func GetAllMazeEquipPosRankV8Config() []*MazeEquipPosRankV8ConfigRow {
	return gConfigData.GetAllMazeEquipPosRankV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipPosRankV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipPosRankV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipPosRankV8ConfigRow from maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx maze_equip_pos_rank_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipPosRankV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_pos_rank_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_pos_rank_v8.json",
		"maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx", "maze_equip_pos_rank_v8",
		&gMazeEquipPosRankV8Parser{}, &gMazeEquipPosRankV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipPosRankV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipPosRankV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipPosRankV8Config))(c)
		return true
	})
}

// RegisterMazeEquipPosRankV8InitCallBack reg config update func (old func)
var RegisterMazeEquipPosRankV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipPosRankV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipPosRankV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipPosRankV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipPosRankV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipPosRankV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipPosRankV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipPosRankV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipPosRankV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipPosRankV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipPosRankV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipPosRankV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipPosRankV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipPosRankV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosRankV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipPosRankV8ConfigRow", zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"),
			zap.String("sheet", "maze_equip_pos_rank_v8"))
		return
	}
	config, ok := container.(*MazeEquipPosRankV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosRankV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipPosRankV8Config", zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"),
			zap.String("sheet", "maze_equip_pos_rank_v8"))
		return
	}
	config.ConfigRows[row.Pos_id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipPosRankV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipPosRankV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosRankV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipPosRankV8Config", zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"),
			zap.String("sheet", "maze_equip_pos_rank_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipPosRankV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipPosRankV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosRankV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipPosRankV8Config", zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"),
			zap.String("sheet", "maze_equip_pos_rank_v8"))
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
type gMazeEquipPosRankV8Parser struct {
}

// New new config row data
func (*gMazeEquipPosRankV8Parser) New() interface{} {
	return &MazeEquipPosRankV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipPosRankV8Parser) Fields() []string {
	return gMazeEquipPosRankV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipPosRankV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipPosRankV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipPosRankV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipPosRankV8ConfigRow", zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"),
			zap.String("sheet", "maze_equip_pos_rank_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipPosRankV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipPosRankV8ConfigRow",
			zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"),
			zap.String("sheet", "maze_equip_pos_rank_v8"), zap.Int("need_count", len(gMazeEquipPosRankV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 pos_id : 部位id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field pos_id 部位id to int32 failed")
			logger.ErrorWF("parse field pos_id 部位id to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"), zap.String("sheet", "maze_equip_pos_rank_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Pos_id = int32(tmp)
	}

	// parse column 1 name : 备注
	if data[1] != "" {
		config.Name = data[1]
	}

	// parse column 2 sub_type_name : 子类型
	if data[2] != "" {

		config.Sub_type_name = make(map[int32]string)
		var key int32
		var value string
		vals := strings.Split(data[2], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field sub_type_name 子类型 to key int32 failed")
				logger.ErrorWF("parse map field sub_type_name 子类型 to key int32 failed.",
					zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"), zap.String("sheet", "maze_equip_pos_rank_v8"),
					// zap.String("field_data",data[2]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			value = items[1]
			config.Sub_type_name[key] = value
		}
	}

	// parse column 3 rank : 展示顺序
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field rank 展示顺序 to int32 failed")
			logger.ErrorWF("parse field rank 展示顺序 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"), zap.String("sheet", "maze_equip_pos_rank_v8"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Rank = int32(tmp)
	}

	// parse column 4 is_default : 是否默认解锁
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field is_default 是否默认解锁 to int32 failed")
			logger.ErrorWF("parse field is_default 是否默认解锁 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"), zap.String("sheet", "maze_equip_pos_rank_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Is_default = int32(tmp)
	}

	// parse column 5 need_level : 需要人偶等级
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field need_level 需要人偶等级 to int32 failed")
			logger.ErrorWF("parse field need_level 需要人偶等级 to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"), zap.String("sheet", "maze_equip_pos_rank_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Need_level = int32(tmp)
	}

	// parse column 6 need_task : 需要完成任务id
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field need_task 需要完成任务id to int32 failed")
			logger.ErrorWF("parse field need_task 需要完成任务id to int32 failed.",
				zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"), zap.String("sheet", "maze_equip_pos_rank_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Need_task = int32(tmp)
	}

	// parse column 7 need_dungeon_id : 需要通关章节id（=）：层数id(>=)
	if data[7] != "" {

		config.Need_dungeon_id = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[7], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field need_dungeon_id 需要通关章节id（=）：层数id(>=) to key int32 failed")
				logger.ErrorWF("parse map field need_dungeon_id 需要通关章节id（=）：层数id(>=) to key int32 failed.",
					zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"), zap.String("sheet", "maze_equip_pos_rank_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field need_dungeon_id 需要通关章节id（=）：层数id(>=) to value int32 failed")
				logger.ErrorWF("parse map field need_dungeon_id 需要通关章节id（=）：层数id(>=) to value int32 failed.",
					zap.String("xlsx", "maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx"), zap.String("sheet", "maze_equip_pos_rank_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Need_dungeon_id[key] = value
		}
	}

	// parse column 8 unlock_desc : 未解锁描述
	if data[8] != "" {
		config.Unlock_desc = data[8]
	}
	return
}

var gMazeEquipPosRankV8Fields = []string{
	"pos_id",
	"name",
	"sub_type_name",
	"rank",
	"is_default",
	"need_level",
	"need_task",
	"need_dungeon_id",
	"unlock_desc",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipPosRankV8Parser{}
	loader := &gMazeEquipPosRankV8Loader{}
	var data [][]string
	data, err = load("maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx", "maze_equip_pos_rank_v8", gMazeEquipPosRankV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx maze_equip_pos_rank_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx maze_equip_pos_rank_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipPosRankV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx maze_equip_pos_rank_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipPosRankV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx maze_equip_pos_rank_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx maze_equip_pos_rank_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx maze_equip_pos_rank_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_pos_rank_v8【迷宫-装备-部位排序】.xlsx maze_equip_pos_rank_v8 data success.")
	return
}
