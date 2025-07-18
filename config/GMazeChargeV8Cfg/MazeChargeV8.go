package GMazeChargeV8Cfg

import (
	"errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"
)

// MazeChargeV8ConfigRow from maze_charge_v8【迷宫-充值】.xlsx maze_charge_v8
type MazeChargeV8ConfigRow struct {
	Id           int32 `json:"id"`           // 序号
	Unique_id    int32 `json:"unique_id"`    // 商品支付unique_id
	Currency_num int32 `json:"currency_num"` // 货币数目
}

// MazeChargeV8Config from maze_charge_v8【迷宫-充值】.xlsx maze_charge_v8
type MazeChargeV8Config struct {
	ConfigRows map[int32]*MazeChargeV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeChargeV8Config {
	ret := &MazeChargeV8Config{ConfigRows: map[int32]*MazeChargeV8ConfigRow{}}
	return ret
}

// GetMazeChargeV8Config get one config by configId
func (c *MazeChargeV8Config) GetMazeChargeV8Config(configId int32) *MazeChargeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeChargeV8Config) Get(configId int32) *MazeChargeV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeChargeV8Config get all config slice
func (c *MazeChargeV8Config) GetAllMazeChargeV8Config() (res []*MazeChargeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeChargeV8Config) GetAll() (res []*MazeChargeV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeChargeV8Config

// GetMazeChargeV8Config pkg func. get one config by configId
func GetMazeChargeV8Config(configId int32) *MazeChargeV8ConfigRow {
	return gConfigData.GetMazeChargeV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeChargeV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeChargeV8Config pkg func. get all config slice
func GetAllMazeChargeV8Config() []*MazeChargeV8ConfigRow {
	return gConfigData.GetAllMazeChargeV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeChargeV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeChargeV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeChargeV8ConfigRow from maze_charge_v8【迷宫-充值】.xlsx maze_charge_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeChargeV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_charge_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_charge_v8.json",
		"maze_charge_v8【迷宫-充值】.xlsx", "maze_charge_v8",
		&gMazeChargeV8Parser{}, &gMazeChargeV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeChargeV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeChargeV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeChargeV8Config))(c)
		return true
	})
}

// RegisterMazeChargeV8InitCallBack reg config update func (old func)
var RegisterMazeChargeV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeChargeV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeChargeV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeChargeV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeChargeV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeChargeV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeChargeV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeChargeV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeChargeV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeChargeV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeChargeV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeChargeV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeChargeV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeChargeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeChargeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeChargeV8ConfigRow", zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"),
			zap.String("sheet", "maze_charge_v8"))
		return
	}
	config, ok := container.(*MazeChargeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeChargeV8Config")
		logger.ErrorWF("invalid type. not *MazeChargeV8Config", zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"),
			zap.String("sheet", "maze_charge_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeChargeV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeChargeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeChargeV8Config")
		logger.ErrorWF("invalid type. not *MazeChargeV8Config", zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"),
			zap.String("sheet", "maze_charge_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeChargeV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeChargeV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeChargeV8Config")
		logger.ErrorWF("invalid type. not *MazeChargeV8Config", zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"),
			zap.String("sheet", "maze_charge_v8"))
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
type gMazeChargeV8Parser struct {
}

// New new config row data
func (*gMazeChargeV8Parser) New() interface{} {
	return &MazeChargeV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeChargeV8Parser) Fields() []string {
	return gMazeChargeV8Fields
}

// Parse parse raw data to row data
func (*gMazeChargeV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeChargeV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeChargeV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeChargeV8ConfigRow", zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"),
			zap.String("sheet", "maze_charge_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeChargeV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeChargeV8ConfigRow",
			zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"),
			zap.String("sheet", "maze_charge_v8"), zap.Int("need_count", len(gMazeChargeV8Fields)),
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
				zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"), zap.String("sheet", "maze_charge_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 unique_id : 商品支付unique_id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field unique_id 商品支付unique_id to int32 failed")
			logger.ErrorWF("parse field unique_id 商品支付unique_id to int32 failed.",
				zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"), zap.String("sheet", "maze_charge_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Unique_id = int32(tmp)
	}

	// parse column 2 currency_num : 货币数目
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field currency_num 货币数目 to int32 failed")
			logger.ErrorWF("parse field currency_num 货币数目 to int32 failed.",
				zap.String("xlsx", "maze_charge_v8【迷宫-充值】.xlsx"), zap.String("sheet", "maze_charge_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Currency_num = int32(tmp)
	}
	return
}

var gMazeChargeV8Fields = []string{
	"id",
	"unique_id",
	"currency_num",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeChargeV8Parser{}
	loader := &gMazeChargeV8Loader{}
	var data [][]string
	data, err = load("maze_charge_v8【迷宫-充值】.xlsx", "maze_charge_v8", gMazeChargeV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_charge_v8【迷宫-充值】.xlsx maze_charge_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_charge_v8【迷宫-充值】.xlsx maze_charge_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeChargeV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_charge_v8【迷宫-充值】.xlsx maze_charge_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeChargeV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_charge_v8【迷宫-充值】.xlsx maze_charge_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_charge_v8【迷宫-充值】.xlsx maze_charge_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_charge_v8【迷宫-充值】.xlsx maze_charge_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_charge_v8【迷宫-充值】.xlsx maze_charge_v8 data success.")
	return
}
