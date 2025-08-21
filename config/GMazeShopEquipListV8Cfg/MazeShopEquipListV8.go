package GMazeShopEquipListV8Cfg

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

// MazeShopEquipListV8ConfigRow from maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx maze_shop_equip_list_v8
type MazeShopEquipListV8ConfigRow struct {
	Order    int32   `json:"order"`    // 序号
	List_id  int32   `json:"list_id"`  // 队列id
	Pos_id   int32   `json:"pos_id"`   // 部位id
	Equip_id []int32 `json:"equip_id"` // 购买到装备的id
}

// MazeShopEquipListV8Config from maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx maze_shop_equip_list_v8
type MazeShopEquipListV8Config struct {
	ConfigRows map[int32]*MazeShopEquipListV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeShopEquipListV8Config {
	ret := &MazeShopEquipListV8Config{ConfigRows: map[int32]*MazeShopEquipListV8ConfigRow{}}
	return ret
}

// GetMazeShopEquipListV8Config get one config by configId
func (c *MazeShopEquipListV8Config) GetMazeShopEquipListV8Config(configId int32) *MazeShopEquipListV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeShopEquipListV8Config) Get(configId int32) *MazeShopEquipListV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeShopEquipListV8Config get all config slice
func (c *MazeShopEquipListV8Config) GetAllMazeShopEquipListV8Config() (res []*MazeShopEquipListV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeShopEquipListV8Config) GetAll() (res []*MazeShopEquipListV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeShopEquipListV8Config

// GetMazeShopEquipListV8Config pkg func. get one config by configId
func GetMazeShopEquipListV8Config(configId int32) *MazeShopEquipListV8ConfigRow {
	return gConfigData.GetMazeShopEquipListV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeShopEquipListV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeShopEquipListV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_shop_equip_list_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeShopEquipListV8Config pkg func. get all config slice
func GetAllMazeShopEquipListV8Config() []*MazeShopEquipListV8ConfigRow {
	return gConfigData.GetAllMazeShopEquipListV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeShopEquipListV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeShopEquipListV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeShopEquipListV8ConfigRow from maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx maze_shop_equip_list_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeShopEquipListV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_shop_equip_list_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_shop_equip_list_v8.json",
		"maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx", "maze_shop_equip_list_v8",
		&gMazeShopEquipListV8Parser{}, &gMazeShopEquipListV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeShopEquipListV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeShopEquipListV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeShopEquipListV8Config))(c)
		return true
	})
}

// RegisterMazeShopEquipListV8InitCallBack reg config update func (old func)
var RegisterMazeShopEquipListV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeShopEquipListV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeShopEquipListV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeShopEquipListV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeShopEquipListV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeShopEquipListV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeShopEquipListV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeShopEquipListV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeShopEquipListV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeShopEquipListV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeShopEquipListV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeShopEquipListV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeShopEquipListV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeShopEquipListV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeShopEquipListV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeShopEquipListV8ConfigRow", zap.String("xlsx", "maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx"),
			zap.String("sheet", "maze_shop_equip_list_v8"))
		return
	}
	config, ok := container.(*MazeShopEquipListV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeShopEquipListV8Config")
		logger.ErrorWF("invalid type. not *MazeShopEquipListV8Config", zap.String("xlsx", "maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx"),
			zap.String("sheet", "maze_shop_equip_list_v8"))
		return
	}
	config.ConfigRows[row.Order] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeShopEquipListV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeShopEquipListV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeShopEquipListV8Config")
		logger.ErrorWF("invalid type. not *MazeShopEquipListV8Config", zap.String("xlsx", "maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx"),
			zap.String("sheet", "maze_shop_equip_list_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeShopEquipListV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeShopEquipListV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeShopEquipListV8Config")
		logger.ErrorWF("invalid type. not *MazeShopEquipListV8Config", zap.String("xlsx", "maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx"),
			zap.String("sheet", "maze_shop_equip_list_v8"))
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
type gMazeShopEquipListV8Parser struct {
}

// New new config row data
func (*gMazeShopEquipListV8Parser) New() interface{} {
	return &MazeShopEquipListV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeShopEquipListV8Parser) Fields() []string {
	return gMazeShopEquipListV8Fields
}

// Parse parse raw data to row data
func (*gMazeShopEquipListV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeShopEquipListV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeShopEquipListV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeShopEquipListV8ConfigRow", zap.String("xlsx", "maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx"),
			zap.String("sheet", "maze_shop_equip_list_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeShopEquipListV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeShopEquipListV8ConfigRow",
			zap.String("xlsx", "maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx"),
			zap.String("sheet", "maze_shop_equip_list_v8"), zap.Int("need_count", len(gMazeShopEquipListV8Fields)),
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
				zap.String("xlsx", "maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx"), zap.String("sheet", "maze_shop_equip_list_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Order = int32(tmp)
	}

	// parse column 1 list_id : 队列id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field list_id 队列id to int32 failed")
			logger.ErrorWF("parse field list_id 队列id to int32 failed.",
				zap.String("xlsx", "maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx"), zap.String("sheet", "maze_shop_equip_list_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.List_id = int32(tmp)
	}

	// parse column 2 pos_id : 部位id
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field pos_id 部位id to int32 failed")
			logger.ErrorWF("parse field pos_id 部位id to int32 failed.",
				zap.String("xlsx", "maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx"), zap.String("sheet", "maze_shop_equip_list_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Pos_id = int32(tmp)
	}

	// parse column 3 equip_id : 购买到装备的id
	if data[3] != "" {

		vals := strings.Split(data[3], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field equip_id 购买到装备的id to []int32 failed")
				logger.ErrorWF("parse array field equip_id 购买到装备的id to []int32 failed.",
					zap.String("xlsx", "maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx"), zap.String("sheet", "maze_shop_equip_list_v8"),
					// zap.String("field_data",data[3]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Equip_id = append(config.Equip_id, int32(tmp))
		}
	}
	return
}

var gMazeShopEquipListV8Fields = []string{
	"order",
	"list_id",
	"pos_id",
	"equip_id",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeShopEquipListV8Parser{}
	loader := &gMazeShopEquipListV8Loader{}
	var data [][]string
	data, err = load("maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx", "maze_shop_equip_list_v8", gMazeShopEquipListV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx maze_shop_equip_list_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx maze_shop_equip_list_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeShopEquipListV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx maze_shop_equip_list_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeShopEquipListV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx maze_shop_equip_list_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx maze_shop_equip_list_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx maze_shop_equip_list_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_shop_equip_list_v8【迷宫-商店购买装备队列】.xlsx maze_shop_equip_list_v8 data success.")
	return
}
