package GMazeEquipSuiteAttrV8Cfg

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

// MazeEquipSuiteAttrV8ConfigRow from maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx maze_equip_suite_attr_v8
type MazeEquipSuiteAttrV8ConfigRow struct {
	Order         int32           `json:"order"`         // 序号
	Suite_id      int32           `json:"suite_id"`      // 套装id
	Affix_num     int32           `json:"affix_num"`     // 套装数量
	Add_attr      map[int32]int64 `json:"add_attr"`      // 凑齐对应数量后增加的属性:属性值(直接替换已有,目前实现的是只取当前套的属性）
	Add_attr_show map[int32]int64 `json:"add_attr_show"` // 属性面板展示值
	Add_attr_desc string          `json:"add_attr_desc"` // 套装效果描述
}

// MazeEquipSuiteAttrV8Config from maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx maze_equip_suite_attr_v8
type MazeEquipSuiteAttrV8Config struct {
	ConfigRows map[int32]*MazeEquipSuiteAttrV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeEquipSuiteAttrV8Config {
	ret := &MazeEquipSuiteAttrV8Config{ConfigRows: map[int32]*MazeEquipSuiteAttrV8ConfigRow{}}
	return ret
}

// GetMazeEquipSuiteAttrV8Config get one config by configId
func (c *MazeEquipSuiteAttrV8Config) GetMazeEquipSuiteAttrV8Config(configId int32) *MazeEquipSuiteAttrV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeEquipSuiteAttrV8Config) Get(configId int32) *MazeEquipSuiteAttrV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeEquipSuiteAttrV8Config get all config slice
func (c *MazeEquipSuiteAttrV8Config) GetAllMazeEquipSuiteAttrV8Config() (res []*MazeEquipSuiteAttrV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeEquipSuiteAttrV8Config) GetAll() (res []*MazeEquipSuiteAttrV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeEquipSuiteAttrV8Config

// GetMazeEquipSuiteAttrV8Config pkg func. get one config by configId
func GetMazeEquipSuiteAttrV8Config(configId int32) *MazeEquipSuiteAttrV8ConfigRow {
	return gConfigData.GetMazeEquipSuiteAttrV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeEquipSuiteAttrV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeEquipSuiteAttrV8Config pkg func. get all config slice
func GetAllMazeEquipSuiteAttrV8Config() []*MazeEquipSuiteAttrV8ConfigRow {
	return gConfigData.GetAllMazeEquipSuiteAttrV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeEquipSuiteAttrV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeEquipSuiteAttrV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeEquipSuiteAttrV8ConfigRow from maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx maze_equip_suite_attr_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeEquipSuiteAttrV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_equip_suite_attr_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_equip_suite_attr_v8.json",
		"maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx", "maze_equip_suite_attr_v8",
		&gMazeEquipSuiteAttrV8Parser{}, &gMazeEquipSuiteAttrV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeEquipSuiteAttrV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeEquipSuiteAttrV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeEquipSuiteAttrV8Config))(c)
		return true
	})
}

// RegisterMazeEquipSuiteAttrV8InitCallBack reg config update func (old func)
var RegisterMazeEquipSuiteAttrV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeEquipSuiteAttrV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeEquipSuiteAttrV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeEquipSuiteAttrV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeEquipSuiteAttrV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeEquipSuiteAttrV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeEquipSuiteAttrV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeEquipSuiteAttrV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeEquipSuiteAttrV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeEquipSuiteAttrV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeEquipSuiteAttrV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeEquipSuiteAttrV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeEquipSuiteAttrV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeEquipSuiteAttrV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteAttrV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteAttrV8ConfigRow", zap.String("xlsx", "maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx"),
			zap.String("sheet", "maze_equip_suite_attr_v8"))
		return
	}
	config, ok := container.(*MazeEquipSuiteAttrV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteAttrV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteAttrV8Config", zap.String("xlsx", "maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx"),
			zap.String("sheet", "maze_equip_suite_attr_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeEquipSuiteAttrV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeEquipSuiteAttrV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteAttrV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteAttrV8Config", zap.String("xlsx", "maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx"),
			zap.String("sheet", "maze_equip_suite_attr_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeEquipSuiteAttrV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeEquipSuiteAttrV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteAttrV8Config")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteAttrV8Config", zap.String("xlsx", "maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx"),
			zap.String("sheet", "maze_equip_suite_attr_v8"))
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
type gMazeEquipSuiteAttrV8Parser struct {
}

// New new config row data
func (*gMazeEquipSuiteAttrV8Parser) New() interface{} {
	return &MazeEquipSuiteAttrV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeEquipSuiteAttrV8Parser) Fields() []string {
	return gMazeEquipSuiteAttrV8Fields
}

// Parse parse raw data to row data
func (*gMazeEquipSuiteAttrV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeEquipSuiteAttrV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeEquipSuiteAttrV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeEquipSuiteAttrV8ConfigRow", zap.String("xlsx", "maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx"),
			zap.String("sheet", "maze_equip_suite_attr_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeEquipSuiteAttrV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeEquipSuiteAttrV8ConfigRow",
			zap.String("xlsx", "maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx"),
			zap.String("sheet", "maze_equip_suite_attr_v8"), zap.Int("need_count", len(gMazeEquipSuiteAttrV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 order : 序号
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field order 序号 to int32 failed")
			logger.ErrorWF("parse field order 序号 to int32 failed.",
				zap.String("xlsx", "maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx"), zap.String("sheet", "maze_equip_suite_attr_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 suite_id : 套装id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field suite_id 套装id to int32 failed")
			logger.ErrorWF("parse field suite_id 套装id to int32 failed.",
				zap.String("xlsx", "maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx"), zap.String("sheet", "maze_equip_suite_attr_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Suite_id = int32(tmp)
	}

	// parse column 2 affix_num : 套装数量
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field affix_num 套装数量 to int32 failed")
			logger.ErrorWF("parse field affix_num 套装数量 to int32 failed.",
				zap.String("xlsx", "maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx"), zap.String("sheet", "maze_equip_suite_attr_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Affix_num = int32(tmp)
	}

	// parse column 3 add_attr : 凑齐对应数量后增加的属性:属性值(直接替换已有,目前实现的是只取当前套的属性）
	if data[3] != "" {

		config.Add_attr = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field add_attr 凑齐对应数量后增加的属性:属性值(直接替换已有,目前实现的是只取当前套的属性） to key int32 failed")
				logger.ErrorWF("parse map field add_attr 凑齐对应数量后增加的属性:属性值(直接替换已有,目前实现的是只取当前套的属性） to key int32 failed.",
					zap.String("xlsx", "maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx"), zap.String("sheet", "maze_equip_suite_attr_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field add_attr 凑齐对应数量后增加的属性:属性值(直接替换已有,目前实现的是只取当前套的属性） to value int64 failed")
				logger.ErrorWF("parse map field add_attr 凑齐对应数量后增加的属性:属性值(直接替换已有,目前实现的是只取当前套的属性） to value int64 failed.",
					zap.String("xlsx", "maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx"), zap.String("sheet", "maze_equip_suite_attr_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Add_attr[key] = value
		}
	}

	// parse column 4 add_attr_show : 属性面板展示值
	if data[4] != "" {

		config.Add_attr_show = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[4], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field add_attr_show 属性面板展示值 to key int32 failed")
				logger.ErrorWF("parse map field add_attr_show 属性面板展示值 to key int32 failed.",
					zap.String("xlsx", "maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx"), zap.String("sheet", "maze_equip_suite_attr_v8"),
					// zap.String("field_data",data[4]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field add_attr_show 属性面板展示值 to value int64 failed")
				logger.ErrorWF("parse map field add_attr_show 属性面板展示值 to value int64 failed.",
					zap.String("xlsx", "maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx"), zap.String("sheet", "maze_equip_suite_attr_v8"),
					// zap.String("field_data",data[4]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Add_attr_show[key] = value
		}
	}

	// parse column 5 add_attr_desc : 套装效果描述
	if data[5] != "" {
		config.Add_attr_desc = data[5]
	}
	return
}

var gMazeEquipSuiteAttrV8Fields = []string{
	"order",
	"suite_id",
	"affix_num",
	"add_attr",
	"add_attr_show",
	"add_attr_desc",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeEquipSuiteAttrV8Parser{}
	loader := &gMazeEquipSuiteAttrV8Loader{}
	var data [][]string
	data, err = load("maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx", "maze_equip_suite_attr_v8", gMazeEquipSuiteAttrV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx maze_equip_suite_attr_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx maze_equip_suite_attr_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeEquipSuiteAttrV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx maze_equip_suite_attr_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeEquipSuiteAttrV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx maze_equip_suite_attr_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx maze_equip_suite_attr_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx maze_equip_suite_attr_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_equip_suite_attr_v8【迷宫-装备-套装部位激活属性】.xlsx maze_equip_suite_attr_v8 data success.")
	return
}
