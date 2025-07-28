package GMazeBoxV8Cfg

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

// MazeBoxV8ConfigRow from maze_box_v8【迷宫-宝箱】.xlsx maze_box_v8
type MazeBoxV8ConfigRow struct {
	Id                 int32           `json:"id"`                 // 宝箱id
	Award_equip        []int32         `json:"award_equip"`        // 装备奖励（非首次）
	Drop_exp_num       map[int32]int64 `json:"drop_exp_num"`       // 冒险等级:掉落经验数量（非首次）
	Drop_item          map[int32]int64 `json:"drop_item"`          // 宝箱掉落物品id：数量（非首次）
	Add_kongfu         int32           `json:"add_kongfu"`         // 增加通关值
	Award_equip_first  []int32         `json:"award_equip_first"`  // 装备奖励（首次）
	Drop_exp_num_first map[int32]int64 `json:"drop_exp_num_first"` // 冒险等级:掉落经验数量（首次）
	Drop_item_first    map[int32]int64 `json:"drop_item_first"`    // 宝箱掉落物品id：数量（首次）
	Res_type           int32           `json:"res_type"`           // 资源类型
	Res_type_first     int32           `json:"res_type_first"`     // 资源类型（首次）
	Res_id             int32           `json:"res_id"`             // 资源id
	Res_id_first       int32           `json:"res_id_first"`       // 资源id（首次）
}

// MazeBoxV8Config from maze_box_v8【迷宫-宝箱】.xlsx maze_box_v8
type MazeBoxV8Config struct {
	ConfigRows map[int32]*MazeBoxV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeBoxV8Config {
	ret := &MazeBoxV8Config{ConfigRows: map[int32]*MazeBoxV8ConfigRow{}}
	return ret
}

// GetMazeBoxV8Config get one config by configId
func (c *MazeBoxV8Config) GetMazeBoxV8Config(configId int32) *MazeBoxV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeBoxV8Config) Get(configId int32) *MazeBoxV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeBoxV8Config get all config slice
func (c *MazeBoxV8Config) GetAllMazeBoxV8Config() (res []*MazeBoxV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeBoxV8Config) GetAll() (res []*MazeBoxV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeBoxV8Config

// GetMazeBoxV8Config pkg func. get one config by configId
func GetMazeBoxV8Config(configId int32) *MazeBoxV8ConfigRow {
	return gConfigData.GetMazeBoxV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeBoxV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeBoxV8Config pkg func. get all config slice
func GetAllMazeBoxV8Config() []*MazeBoxV8ConfigRow {
	return gConfigData.GetAllMazeBoxV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeBoxV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeBoxV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeBoxV8ConfigRow from maze_box_v8【迷宫-宝箱】.xlsx maze_box_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeBoxV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_box_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_box_v8.json",
		"maze_box_v8【迷宫-宝箱】.xlsx", "maze_box_v8",
		&gMazeBoxV8Parser{}, &gMazeBoxV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeBoxV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeBoxV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeBoxV8Config))(c)
		return true
	})
}

// RegisterMazeBoxV8InitCallBack reg config update func (old func)
var RegisterMazeBoxV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeBoxV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeBoxV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeBoxV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeBoxV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeBoxV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeBoxV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeBoxV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeBoxV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeBoxV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeBoxV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeBoxV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeBoxV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeBoxV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBoxV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBoxV8ConfigRow", zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"),
			zap.String("sheet", "maze_box_v8"))
		return
	}
	config, ok := container.(*MazeBoxV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBoxV8Config")
		logger.ErrorWF("invalid type. not *MazeBoxV8Config", zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"),
			zap.String("sheet", "maze_box_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeBoxV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeBoxV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBoxV8Config")
		logger.ErrorWF("invalid type. not *MazeBoxV8Config", zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"),
			zap.String("sheet", "maze_box_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeBoxV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeBoxV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeBoxV8Config")
		logger.ErrorWF("invalid type. not *MazeBoxV8Config", zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"),
			zap.String("sheet", "maze_box_v8"))
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
type gMazeBoxV8Parser struct {
}

// New new config row data
func (*gMazeBoxV8Parser) New() interface{} {
	return &MazeBoxV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeBoxV8Parser) Fields() []string {
	return gMazeBoxV8Fields
}

// Parse parse raw data to row data
func (*gMazeBoxV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeBoxV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeBoxV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeBoxV8ConfigRow", zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"),
			zap.String("sheet", "maze_box_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeBoxV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeBoxV8ConfigRow",
			zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"),
			zap.String("sheet", "maze_box_v8"), zap.Int("need_count", len(gMazeBoxV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 宝箱id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 宝箱id to int32 failed")
			logger.ErrorWF("parse field id 宝箱id to int32 failed.",
				zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 award_equip : 装备奖励（非首次）
	if data[1] != "" {

		vals := strings.Split(data[1], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field award_equip 装备奖励（非首次） to []int32 failed")
				logger.ErrorWF("parse array field award_equip 装备奖励（非首次） to []int32 failed.",
					zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
					// zap.String("field_data",data[1]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Award_equip = append(config.Award_equip, int32(tmp))
		}
	}

	// parse column 2 drop_exp_num : 冒险等级:掉落经验数量（非首次）
	if data[2] != "" {

		config.Drop_exp_num = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[2], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_exp_num 冒险等级:掉落经验数量（非首次） to key int32 failed")
				logger.ErrorWF("parse map field drop_exp_num 冒险等级:掉落经验数量（非首次） to key int32 failed.",
					zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
					// zap.String("field_data",data[2]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_exp_num 冒险等级:掉落经验数量（非首次） to value int64 failed")
				logger.ErrorWF("parse map field drop_exp_num 冒险等级:掉落经验数量（非首次） to value int64 failed.",
					zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
					// zap.String("field_data",data[2]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Drop_exp_num[key] = value
		}
	}

	// parse column 3 drop_item : 宝箱掉落物品id：数量（非首次）
	if data[3] != "" {

		config.Drop_item = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_item 宝箱掉落物品id：数量（非首次） to key int32 failed")
				logger.ErrorWF("parse map field drop_item 宝箱掉落物品id：数量（非首次） to key int32 failed.",
					zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_item 宝箱掉落物品id：数量（非首次） to value int64 failed")
				logger.ErrorWF("parse map field drop_item 宝箱掉落物品id：数量（非首次） to value int64 failed.",
					zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Drop_item[key] = value
		}
	}

	// parse column 4 add_kongfu : 增加通关值
	if data[4] != "" {
		tmp, err = strconv.ParseInt(data[4], 10, 64)
		if err != nil {
			err = errors.New("parse field add_kongfu 增加通关值 to int32 failed")
			logger.ErrorWF("parse field add_kongfu 增加通关值 to int32 failed.",
				zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
				zap.String("parse_data", data[4]),
				zap.Error(err))
			return
		}
		config.Add_kongfu = int32(tmp)
	}

	// parse column 5 award_equip_first : 装备奖励（首次）
	if data[5] != "" {

		vals := strings.Split(data[5], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field award_equip_first 装备奖励（首次） to []int32 failed")
				logger.ErrorWF("parse array field award_equip_first 装备奖励（首次） to []int32 failed.",
					zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
					// zap.String("field_data",data[5]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Award_equip_first = append(config.Award_equip_first, int32(tmp))
		}
	}

	// parse column 6 drop_exp_num_first : 冒险等级:掉落经验数量（首次）
	if data[6] != "" {

		config.Drop_exp_num_first = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[6], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_exp_num_first 冒险等级:掉落经验数量（首次） to key int32 failed")
				logger.ErrorWF("parse map field drop_exp_num_first 冒险等级:掉落经验数量（首次） to key int32 failed.",
					zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
					// zap.String("field_data",data[6]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_exp_num_first 冒险等级:掉落经验数量（首次） to value int64 failed")
				logger.ErrorWF("parse map field drop_exp_num_first 冒险等级:掉落经验数量（首次） to value int64 failed.",
					zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
					// zap.String("field_data",data[6]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Drop_exp_num_first[key] = value
		}
	}

	// parse column 7 drop_item_first : 宝箱掉落物品id：数量（首次）
	if data[7] != "" {

		config.Drop_item_first = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[7], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_item_first 宝箱掉落物品id：数量（首次） to key int32 failed")
				logger.ErrorWF("parse map field drop_item_first 宝箱掉落物品id：数量（首次） to key int32 failed.",
					zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field drop_item_first 宝箱掉落物品id：数量（首次） to value int64 failed")
				logger.ErrorWF("parse map field drop_item_first 宝箱掉落物品id：数量（首次） to value int64 failed.",
					zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
					// zap.String("field_data",data[7]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Drop_item_first[key] = value
		}
	}

	// parse column 8 res_type : 资源类型
	if data[8] != "" {
		tmp, err = strconv.ParseInt(data[8], 10, 64)
		if err != nil {
			err = errors.New("parse field res_type 资源类型 to int32 failed")
			logger.ErrorWF("parse field res_type 资源类型 to int32 failed.",
				zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
				zap.String("parse_data", data[8]),
				zap.Error(err))
			return
		}
		config.Res_type = int32(tmp)
	}

	// parse column 9 res_type_first : 资源类型（首次）
	if data[9] != "" {
		tmp, err = strconv.ParseInt(data[9], 10, 64)
		if err != nil {
			err = errors.New("parse field res_type_first 资源类型（首次） to int32 failed")
			logger.ErrorWF("parse field res_type_first 资源类型（首次） to int32 failed.",
				zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
				zap.String("parse_data", data[9]),
				zap.Error(err))
			return
		}
		config.Res_type_first = int32(tmp)
	}

	// parse column 10 res_id : 资源id
	if data[10] != "" {
		tmp, err = strconv.ParseInt(data[10], 10, 64)
		if err != nil {
			err = errors.New("parse field res_id 资源id to int32 failed")
			logger.ErrorWF("parse field res_id 资源id to int32 failed.",
				zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
				zap.String("parse_data", data[10]),
				zap.Error(err))
			return
		}
		config.Res_id = int32(tmp)
	}

	// parse column 11 res_id_first : 资源id（首次）
	if data[11] != "" {
		tmp, err = strconv.ParseInt(data[11], 10, 64)
		if err != nil {
			err = errors.New("parse field res_id_first 资源id（首次） to int32 failed")
			logger.ErrorWF("parse field res_id_first 资源id（首次） to int32 failed.",
				zap.String("xlsx", "maze_box_v8【迷宫-宝箱】.xlsx"), zap.String("sheet", "maze_box_v8"),
				zap.String("parse_data", data[11]),
				zap.Error(err))
			return
		}
		config.Res_id_first = int32(tmp)
	}
	return
}

var gMazeBoxV8Fields = []string{
	"id",
	"award_equip",
	"drop_exp_num",
	"drop_item",
	"add_kongfu",
	"award_equip_first",
	"drop_exp_num_first",
	"drop_item_first",
	"res_type",
	"res_type_first",
	"res_id",
	"res_id_first",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeBoxV8Parser{}
	loader := &gMazeBoxV8Loader{}
	var data [][]string
	data, err = load("maze_box_v8【迷宫-宝箱】.xlsx", "maze_box_v8", gMazeBoxV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_box_v8【迷宫-宝箱】.xlsx maze_box_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_box_v8【迷宫-宝箱】.xlsx maze_box_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeBoxV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_box_v8【迷宫-宝箱】.xlsx maze_box_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeBoxV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_box_v8【迷宫-宝箱】.xlsx maze_box_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_box_v8【迷宫-宝箱】.xlsx maze_box_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_box_v8【迷宫-宝箱】.xlsx maze_box_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_box_v8【迷宫-宝箱】.xlsx maze_box_v8 data success.")
	return
}
