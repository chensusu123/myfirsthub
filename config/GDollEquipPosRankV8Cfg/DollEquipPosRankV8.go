package GDollEquipPosRankV8Cfg

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

// DollEquipPosRankV8ConfigRow from doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx doll_equip_pos_rank_v8
type DollEquipPosRankV8ConfigRow struct {
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

// DollEquipPosRankV8Config from doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx doll_equip_pos_rank_v8
type DollEquipPosRankV8Config struct {
	ConfigRows map[int32]*DollEquipPosRankV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *DollEquipPosRankV8Config {
	ret := &DollEquipPosRankV8Config{ConfigRows: map[int32]*DollEquipPosRankV8ConfigRow{}}
	return ret
}

// GetDollEquipPosRankV8Config get one config by configId
func (c *DollEquipPosRankV8Config) GetDollEquipPosRankV8Config(configId int32) *DollEquipPosRankV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *DollEquipPosRankV8Config) Get(configId int32) *DollEquipPosRankV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllDollEquipPosRankV8Config get all config slice
func (c *DollEquipPosRankV8Config) GetAllDollEquipPosRankV8Config() (res []*DollEquipPosRankV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *DollEquipPosRankV8Config) GetAll() (res []*DollEquipPosRankV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *DollEquipPosRankV8Config

// GetDollEquipPosRankV8Config pkg func. get one config by configId
func GetDollEquipPosRankV8Config(configId int32) *DollEquipPosRankV8ConfigRow {
	return gConfigData.GetDollEquipPosRankV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *DollEquipPosRankV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllDollEquipPosRankV8Config pkg func. get all config slice
func GetAllDollEquipPosRankV8Config() []*DollEquipPosRankV8ConfigRow {
	return gConfigData.GetAllDollEquipPosRankV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*DollEquipPosRankV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*DollEquipPosRankV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "DollEquipPosRankV8ConfigRow from doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx doll_equip_pos_rank_v8"
}

// GetRawValue get raw data
func GetRawValue() *DollEquipPosRankV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "doll_equip_pos_rank_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("doll_equip_pos_rank_v8.json",
		"doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx", "doll_equip_pos_rank_v8",
		&gDollEquipPosRankV8Parser{}, &gDollEquipPosRankV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*DollEquipPosRankV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *DollEquipPosRankV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*DollEquipPosRankV8Config))(c)
		return true
	})
}

// RegisterDollEquipPosRankV8InitCallBack reg config update func (old func)
var RegisterDollEquipPosRankV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*DollEquipPosRankV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*DollEquipPosRankV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *DollEquipPosRankV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*DollEquipPosRankV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*DollEquipPosRankV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gDollEquipPosRankV8Loader struct {
}

// NewContainer new data container pointer
func (*gDollEquipPosRankV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gDollEquipPosRankV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*DollEquipPosRankV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gDollEquipPosRankV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*DollEquipPosRankV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gDollEquipPosRankV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*DollEquipPosRankV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *DollEquipPosRankV8ConfigRow")
		logger.ErrorWF("invalid type. not *DollEquipPosRankV8ConfigRow", zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"),
			zap.String("sheet", "doll_equip_pos_rank_v8"))
		return
	}
	config, ok := container.(*DollEquipPosRankV8Config)
	if !ok {
		err = errors.New("invalid type. not *DollEquipPosRankV8Config")
		logger.ErrorWF("invalid type. not *DollEquipPosRankV8Config", zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"),
			zap.String("sheet", "doll_equip_pos_rank_v8"))
		return
	}
	config.ConfigRows[row.Pos_id] = row
	return
}

// GetValue get real map value for json parse
func (*gDollEquipPosRankV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*DollEquipPosRankV8Config)
	if !ok {
		err = errors.New("invalid type. not *DollEquipPosRankV8Config")
		logger.ErrorWF("invalid type. not *DollEquipPosRankV8Config", zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"),
			zap.String("sheet", "doll_equip_pos_rank_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gDollEquipPosRankV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*DollEquipPosRankV8Config)
	if !ok {
		err = errors.New("invalid type. not *DollEquipPosRankV8Config")
		logger.ErrorWF("invalid type. not *DollEquipPosRankV8Config", zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"),
			zap.String("sheet", "doll_equip_pos_rank_v8"))
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
type gDollEquipPosRankV8Parser struct {
}

// New new config row data
func (*gDollEquipPosRankV8Parser) New() interface{} {
	return &DollEquipPosRankV8ConfigRow{}
}

// Fields get config fields names
func (*gDollEquipPosRankV8Parser) Fields() []string {
	return gDollEquipPosRankV8Fields
}

// Parse parse raw data to row data
func (*gDollEquipPosRankV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*DollEquipPosRankV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *DollEquipPosRankV8ConfigRow")
		logger.ErrorWF("invalid type. not *DollEquipPosRankV8ConfigRow", zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"),
			zap.String("sheet", "doll_equip_pos_rank_v8"))
		return
	}
	// compare length
	if len(data) != len(gDollEquipPosRankV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*DollEquipPosRankV8ConfigRow",
			zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"),
			zap.String("sheet", "doll_equip_pos_rank_v8"), zap.Int("need_count", len(gDollEquipPosRankV8Fields)),
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
				zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"), zap.String("sheet", "doll_equip_pos_rank_v8"),
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
					zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"), zap.String("sheet", "doll_equip_pos_rank_v8"),
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
				zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"), zap.String("sheet", "doll_equip_pos_rank_v8"),
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
				zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"), zap.String("sheet", "doll_equip_pos_rank_v8"),
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
				zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"), zap.String("sheet", "doll_equip_pos_rank_v8"),
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
				zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"), zap.String("sheet", "doll_equip_pos_rank_v8"),
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
					zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"), zap.String("sheet", "doll_equip_pos_rank_v8"),
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
					zap.String("xlsx", "doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx"), zap.String("sheet", "doll_equip_pos_rank_v8"),
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

var gDollEquipPosRankV8Fields = []string{
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
	parser := &gDollEquipPosRankV8Parser{}
	loader := &gDollEquipPosRankV8Loader{}
	var data [][]string
	data, err = load("doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx", "doll_equip_pos_rank_v8", gDollEquipPosRankV8Fields)
	if err != nil {
		logger.ErrorWF("load doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx doll_equip_pos_rank_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx doll_equip_pos_rank_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gDollEquipPosRankV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx doll_equip_pos_rank_v8 failed.", zap.Int("row", k), zap.Strings("need", gDollEquipPosRankV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx doll_equip_pos_rank_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx doll_equip_pos_rank_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx doll_equip_pos_rank_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load doll_equip_pos_rank_v8【人偶-装备-部位排序】.xlsx doll_equip_pos_rank_v8 data success.")
	return
}
