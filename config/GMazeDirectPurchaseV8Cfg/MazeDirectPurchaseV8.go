package GMazeDirectPurchaseV8Cfg

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

// MazeDirectPurchaseV8ConfigRow from maze_charge_v8【迷宫-充值】.xlsx maze_direct_purchase_v8
type MazeDirectPurchaseV8ConfigRow struct {
	Id        int32           `json:"id"`        // 序号
	Charge_id int32           `json:"charge_id"` // 充值id
	Sort      int32           `json:"sort"`      // 排序
	Items     map[int32]int64 `json:"items"`     // 礼包内容道具：数量
}

// MazeDirectPurchaseV8Config from maze_charge_v8【迷宫-充值】.xlsx maze_direct_purchase_v8
type MazeDirectPurchaseV8Config struct {
	ConfigRows map[int32]*MazeDirectPurchaseV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeDirectPurchaseV8Config {
	ret := &MazeDirectPurchaseV8Config{ConfigRows: map[int32]*MazeDirectPurchaseV8ConfigRow{}}
	return ret
}

// GetMazeDirectPurchaseV8Config get one config by configId
func (c *MazeDirectPurchaseV8Config) GetMazeDirectPurchaseV8Config(configId int32) *MazeDirectPurchaseV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeDirectPurchaseV8Config) Get(configId int32) *MazeDirectPurchaseV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeDirectPurchaseV8Config get all config slice
func (c *MazeDirectPurchaseV8Config) GetAllMazeDirectPurchaseV8Config() (res []*MazeDirectPurchaseV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeDirectPurchaseV8Config) GetAll() (res []*MazeDirectPurchaseV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeDirectPurchaseV8Config

// GetMazeDirectPurchaseV8Config pkg func. get one config by configId
func GetMazeDirectPurchaseV8Config(configId int32) *MazeDirectPurchaseV8ConfigRow {
	return gConfigData.GetMazeDirectPurchaseV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeDirectPurchaseV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeDirectPurchaseV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_direct_purchase_v8", configId, otps...)
	}
	return cfg
}

// GetAllMazeDirectPurchaseV8Config pkg func. get all config slice
func GetAllMazeDirectPurchaseV8Config() []*MazeDirectPurchaseV8ConfigRow {
	return gConfigData.GetAllMazeDirectPurchaseV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeDirectPurchaseV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeDirectPurchaseV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeDirectPurchaseV8ConfigRow from maze_charge_v8【迷宫-充值】.xlsx maze_direct_purchase_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeDirectPurchaseV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_direct_purchase_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_direct_purchase_v8.json",
		"maze_charge_v8【迷宫-充值】.xlsx", "maze_direct_purchase_v8",
		&gMazeDirectPurchaseV8Parser{}, &gMazeDirectPurchaseV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeDirectPurchaseV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeDirectPurchaseV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeDirectPurchaseV8Config))(c)
		return true
	})
}

// RegisterMazeDirectPurchaseV8InitCallBack reg config update func (old func)
var RegisterMazeDirectPurchaseV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeDirectPurchaseV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeDirectPurchaseV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeDirectPurchaseV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeDirectPurchaseV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeDirectPurchaseV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeDirectPurchaseV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeDirectPurchaseV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeDirectPurchaseV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeDirectPurchaseV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeDirectPurchaseV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeDirectPurchaseV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeDirectPurchaseV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeDirectPurchaseV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeDirectPurchaseV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeDirectPurchaseV8ConfigRow", zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"),
			zap.String("sheet", "maze_direct_purchase_v8"))
		return
	}
	config, ok := container.(*MazeDirectPurchaseV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeDirectPurchaseV8Config")
		logger.ErrorWF("invalid type. not *MazeDirectPurchaseV8Config", zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"),
			zap.String("sheet", "maze_direct_purchase_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeDirectPurchaseV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeDirectPurchaseV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeDirectPurchaseV8Config")
		logger.ErrorWF("invalid type. not *MazeDirectPurchaseV8Config", zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"),
			zap.String("sheet", "maze_direct_purchase_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeDirectPurchaseV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeDirectPurchaseV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeDirectPurchaseV8Config")
		logger.ErrorWF("invalid type. not *MazeDirectPurchaseV8Config", zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"),
			zap.String("sheet", "maze_direct_purchase_v8"))
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
type gMazeDirectPurchaseV8Parser struct {
}

// New new config row data
func (*gMazeDirectPurchaseV8Parser) New() interface{} {
	return &MazeDirectPurchaseV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeDirectPurchaseV8Parser) Fields() []string {
	return gMazeDirectPurchaseV8Fields
}

// Parse parse raw data to row data
func (*gMazeDirectPurchaseV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeDirectPurchaseV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeDirectPurchaseV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeDirectPurchaseV8ConfigRow", zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"),
			zap.String("sheet", "maze_direct_purchase_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeDirectPurchaseV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeDirectPurchaseV8ConfigRow",
			zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"),
			zap.String("sheet", "maze_direct_purchase_v8"), zap.Int("need_count", len(gMazeDirectPurchaseV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 序号
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 序号 to int32 failed")
			logger.ErrorWF("parse field id 序号 to int32 failed.",
				zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"), zap.String("sheet", "maze_direct_purchase_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 charge_id : 充值id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field charge_id 充值id to int32 failed")
			logger.ErrorWF("parse field charge_id 充值id to int32 failed.",
				zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"), zap.String("sheet", "maze_direct_purchase_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Charge_id = int32(tmp)
	}

	// parse column 2 sort : 排序
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field sort 排序 to int32 failed")
			logger.ErrorWF("parse field sort 排序 to int32 failed.",
				zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"), zap.String("sheet", "maze_direct_purchase_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Sort = int32(tmp)
	}

	// parse column 3 items : 礼包内容道具：数量
	if data[3] != "" {

		config.Items = make(map[int32]int64)
		var key int32
		var value int64
		vals := strings.Split(data[3], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field items 礼包内容道具：数量 to key int32 failed")
				logger.ErrorWF("parse map field items 礼包内容道具：数量 to key int32 failed.",
					zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"), zap.String("sheet", "maze_direct_purchase_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field items 礼包内容道具：数量 to value int64 failed")
				logger.ErrorWF("parse map field items 礼包内容道具：数量 to value int64 failed.",
					zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"), zap.String("sheet", "maze_direct_purchase_v8"),
					// zap.String("field_data",data[3]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int64(tmp)
			config.Items[key] = value
		}
	}
	return
}

var gMazeDirectPurchaseV8Fields = []string{
	"id",
	"charge_id",
	"sort",
	"items",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeDirectPurchaseV8Parser{}
	loader := &gMazeDirectPurchaseV8Loader{}
	var data [][]string
	data, err = load("maze_charge_v8【迷宫-充值】.xlsx", "maze_direct_purchase_v8", gMazeDirectPurchaseV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_charge_v8【迷宫-充值】.xlsx maze_direct_purchase_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_charge_v8【迷宫-充值】.xlsx maze_direct_purchase_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeDirectPurchaseV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_charge_v8【迷宫-充值】.xlsx maze_direct_purchase_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeDirectPurchaseV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_charge_v8【迷宫-充值】.xlsx maze_direct_purchase_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_charge_v8【迷宫-充值】.xlsx maze_direct_purchase_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_charge_v8【迷宫-充值】.xlsx maze_direct_purchase_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_charge_v8【迷宫-充值】.xlsx maze_direct_purchase_v8 data success.")
	return
}
