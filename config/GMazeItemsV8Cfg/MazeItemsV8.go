package GMazeItemsV8Cfg

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"go.uber.org/zap"
)

// MazeItemsV8ConfigRow from maze_items_v8【迷宫-道具】.xlsx maze_items_v8
type MazeItemsV8ConfigRow struct {
	Key_id           int32           `json:"key_id"`           // 序号
	Type             int32           `json:"type"`             // 类型
	Id               int32           `json:"id"`               // id
	Prop_name        string          `json:"prop_name"`        // 名称
	Name_desc        string          `json:"name_desc"`        // 名称描述
	Is_native        int32           `json:"is_native"`        // 是否原生端需要信息
	Pic_token        int32           `json:"pic_token"`        // 图片资源版本
	IconAtlas        string          `json:"iconAtlas"`        // icon图集文件夹
	Icon             string          `json:"icon"`             // 图标
	Quality          int32           `json:"quality"`          // 品质
	Is_sell          int32           `json:"is_sell"`          // 能否出售
	Price            int32           `json:"price"`            // 出售价格
	Sale_id          int32           `json:"sale_id"`          // 出售后得到货币的ID
	Price_v2         int32           `json:"price_v2"`         // 培育版出售价格
	Sale_id_v2       int32           `json:"sale_id_v2"`       // 培育版出售后得到货币的ID
	Is_bag           int32           `json:"is_bag"`           // 进入仓库类型
	Fly_type         int32           `json:"fly_type"`         // 飞物品动画类型
	Is_rob           int32           `json:"is_rob"`           // 能否可抢夺
	Diamond_price    map[int32]int32 `json:"diamond_price"`    // 货币：价格
	Diamond_price_v2 map[int32]int32 `json:"diamond_price_v2"` // 货币：价格
	Strength_pay     int32           `json:"strength_pay"`     // 夺取消耗体力
	Protectiontime   int32           `json:"protectiontime"`   // 被夺后保护时间
	Is_use           int32           `json:"is_use"`           // 是否可用
	If_popup         int32           `json:"if_popup"`         // 是否新物品弹窗
	If_binding       int32           `json:"if_binding"`       // 是否绑定可拍卖
	Avaliable_time   int32           `json:"avaliable_time"`   // 有效期（秒）
	JumpType         []int32         `json:"jumpType"`         // 跳转类型
	JumpTarget       string          `json:"jumpTarget"`       // 跳转目标1
	Source_type      int32           `json:"source_type"`      // 来源类型
	Btn_type         int32           `json:"btn_type"`         // 按钮动作
	ItemGettype      int32           `json:"itemGettype"`      // 道具添加类型
	Is_pay           int32           `json:"is_pay"`           // 是否有付费映射
	Price_v8         int32           `json:"price_v8"`         // 人偶版出售价格
	Sale_id_v8       int32           `json:"sale_id_v8"`       // 人偶版出售后得到货币的ID
	Diamond_price_v8 map[int32]int32 `json:"diamond_price_v8"` // 货币：价格
	Overtime_type    int32           `json:"overtime_type"`    // 倒计时类型（0-默认 1-进船背包
}

// MazeItemsV8Config from maze_items_v8【迷宫-道具】.xlsx maze_items_v8
type MazeItemsV8Config struct {
	ConfigRows map[int32]*MazeItemsV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeItemsV8Config {
	ret := &MazeItemsV8Config{ConfigRows: map[int32]*MazeItemsV8ConfigRow{}}
	return ret
}

// GetMazeItemsV8Config get one config by configId
func (c *MazeItemsV8Config) GetMazeItemsV8Config(configId int32) *MazeItemsV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeItemsV8Config) Get(configId int32) *MazeItemsV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeItemsV8Config get all config slice
func (c *MazeItemsV8Config) GetAllMazeItemsV8Config() (res []*MazeItemsV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeItemsV8Config) GetAll() (res []*MazeItemsV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeItemsV8Config

// GetMazeItemsV8Config pkg func. get one config by configId
func GetMazeItemsV8Config(configId int32) *MazeItemsV8ConfigRow {
	return gConfigData.GetMazeItemsV8Config(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeItemsV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeItemsV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_items_v8", configId, otps...)
		logger := fklog.ContextAppLogger(ctx)
		logger.CtxWarn(ctx, "config not found", zap.Int32("config_id", configId), zap.Any("otps", otps))
	}
	return cfg
}

// GetAllMazeItemsV8Config pkg func. get all config slice
func GetAllMazeItemsV8Config() []*MazeItemsV8ConfigRow {
	return gConfigData.GetAllMazeItemsV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeItemsV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeItemsV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeItemsV8ConfigRow from maze_items_v8【迷宫-道具】.xlsx maze_items_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeItemsV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_items_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_items_v8.json",
		"maze_items_v8【迷宫-道具】.xlsx", "maze_items_v8",
		&gMazeItemsV8Parser{}, &gMazeItemsV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeItemsV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeItemsV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeItemsV8Config))(c)
		return true
	})
}

// RegisterMazeItemsV8InitCallBack reg config update func (old func)
var RegisterMazeItemsV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeItemsV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeItemsV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeItemsV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeItemsV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeItemsV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeItemsV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeItemsV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeItemsV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeItemsV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeItemsV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeItemsV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeItemsV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeItemsV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeItemsV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeItemsV8ConfigRow", zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"),
			zap.String("sheet", "maze_items_v8"))
		return
	}
	config, ok := container.(*MazeItemsV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeItemsV8Config")
		logger.ErrorWF("invalid type. not *MazeItemsV8Config", zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"),
			zap.String("sheet", "maze_items_v8"))
		return
	}
	config.ConfigRows[row.Key_id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeItemsV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeItemsV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeItemsV8Config")
		logger.ErrorWF("invalid type. not *MazeItemsV8Config", zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"),
			zap.String("sheet", "maze_items_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeItemsV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeItemsV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeItemsV8Config")
		logger.ErrorWF("invalid type. not *MazeItemsV8Config", zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"),
			zap.String("sheet", "maze_items_v8"))
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
type gMazeItemsV8Parser struct {
}

// New new config row data
func (*gMazeItemsV8Parser) New() interface{} {
	return &MazeItemsV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeItemsV8Parser) Fields() []string {
	return gMazeItemsV8Fields
}

// Parse parse raw data to row data
func (*gMazeItemsV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeItemsV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeItemsV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeItemsV8ConfigRow", zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"),
			zap.String("sheet", "maze_items_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeItemsV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeItemsV8ConfigRow",
			zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"),
			zap.String("sheet", "maze_items_v8"), zap.Int("need_count", len(gMazeItemsV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 key_id : 序号
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field key_id 序号 to int32 failed")
			logger.ErrorWF("parse field key_id 序号 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Key_id = int32(tmp)
	}

	// parse column 1 type : 类型
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field type 类型 to int32 failed")
			logger.ErrorWF("parse field type 类型 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Type = int32(tmp)
	}

	// parse column 2 id : id
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field id id to int32 failed")
			logger.ErrorWF("parse field id id to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 3 prop_name : 名称
	if data[3] != "" {
		config.Prop_name = data[3]
	}

	// parse column 4 name_desc : 名称描述
	if data[4] != "" {
		config.Name_desc = data[4]
	}

	// parse column 5 is_native : 是否原生端需要信息
	if data[5] != "" {
		tmp, err = strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			err = errors.New("parse field is_native 是否原生端需要信息 to int32 failed")
			logger.ErrorWF("parse field is_native 是否原生端需要信息 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[5]),
				zap.Error(err))
			return
		}
		config.Is_native = int32(tmp)
	}

	// parse column 6 pic_token : 图片资源版本
	if data[6] != "" {
		tmp, err = strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			err = errors.New("parse field pic_token 图片资源版本 to int32 failed")
			logger.ErrorWF("parse field pic_token 图片资源版本 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[6]),
				zap.Error(err))
			return
		}
		config.Pic_token = int32(tmp)
	}

	// parse column 7 iconAtlas : icon图集文件夹
	if data[7] != "" {
		config.IconAtlas = data[7]
	}

	// parse column 8 icon : 图标
	if data[8] != "" {
		config.Icon = data[8]
	}

	// parse column 9 quality : 品质
	if data[9] != "" {
		tmp, err = strconv.ParseInt(data[9], 10, 64)
		if err != nil {
			err = errors.New("parse field quality 品质 to int32 failed")
			logger.ErrorWF("parse field quality 品质 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[9]),
				zap.Error(err))
			return
		}
		config.Quality = int32(tmp)
	}

	// parse column 10 is_sell : 能否出售
	if data[10] != "" {
		tmp, err = strconv.ParseInt(data[10], 10, 64)
		if err != nil {
			err = errors.New("parse field is_sell 能否出售 to int32 failed")
			logger.ErrorWF("parse field is_sell 能否出售 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[10]),
				zap.Error(err))
			return
		}
		config.Is_sell = int32(tmp)
	}

	// parse column 11 price : 出售价格
	if data[11] != "" {
		tmp, err = strconv.ParseInt(data[11], 10, 64)
		if err != nil {
			err = errors.New("parse field price 出售价格 to int32 failed")
			logger.ErrorWF("parse field price 出售价格 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[11]),
				zap.Error(err))
			return
		}
		config.Price = int32(tmp)
	}

	// parse column 12 sale_id : 出售后得到货币的ID
	if data[12] != "" {
		tmp, err = strconv.ParseInt(data[12], 10, 64)
		if err != nil {
			err = errors.New("parse field sale_id 出售后得到货币的ID to int32 failed")
			logger.ErrorWF("parse field sale_id 出售后得到货币的ID to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[12]),
				zap.Error(err))
			return
		}
		config.Sale_id = int32(tmp)
	}

	// parse column 13 price_v2 : 培育版出售价格
	if data[13] != "" {
		tmp, err = strconv.ParseInt(data[13], 10, 64)
		if err != nil {
			err = errors.New("parse field price_v2 培育版出售价格 to int32 failed")
			logger.ErrorWF("parse field price_v2 培育版出售价格 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[13]),
				zap.Error(err))
			return
		}
		config.Price_v2 = int32(tmp)
	}

	// parse column 14 sale_id_v2 : 培育版出售后得到货币的ID
	if data[14] != "" {
		tmp, err = strconv.ParseInt(data[14], 10, 64)
		if err != nil {
			err = errors.New("parse field sale_id_v2 培育版出售后得到货币的ID to int32 failed")
			logger.ErrorWF("parse field sale_id_v2 培育版出售后得到货币的ID to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[14]),
				zap.Error(err))
			return
		}
		config.Sale_id_v2 = int32(tmp)
	}

	// parse column 15 is_bag : 进入仓库类型
	if data[15] != "" {
		tmp, err = strconv.ParseInt(data[15], 10, 64)
		if err != nil {
			err = errors.New("parse field is_bag 进入仓库类型 to int32 failed")
			logger.ErrorWF("parse field is_bag 进入仓库类型 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[15]),
				zap.Error(err))
			return
		}
		config.Is_bag = int32(tmp)
	}

	// parse column 16 fly_type : 飞物品动画类型
	if data[16] != "" {
		tmp, err = strconv.ParseInt(data[16], 10, 64)
		if err != nil {
			err = errors.New("parse field fly_type 飞物品动画类型 to int32 failed")
			logger.ErrorWF("parse field fly_type 飞物品动画类型 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[16]),
				zap.Error(err))
			return
		}
		config.Fly_type = int32(tmp)
	}

	// parse column 17 is_rob : 能否可抢夺
	if data[17] != "" {
		tmp, err = strconv.ParseInt(data[17], 10, 64)
		if err != nil {
			err = errors.New("parse field is_rob 能否可抢夺 to int32 failed")
			logger.ErrorWF("parse field is_rob 能否可抢夺 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[17]),
				zap.Error(err))
			return
		}
		config.Is_rob = int32(tmp)
	}

	// parse column 18 diamond_price : 货币：价格
	if data[18] != "" {

		config.Diamond_price = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[18], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field diamond_price 货币：价格 to key int32 failed")
				logger.ErrorWF("parse map field diamond_price 货币：价格 to key int32 failed.",
					zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
					// zap.String("field_data",data[18]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field diamond_price 货币：价格 to value int32 failed")
				logger.ErrorWF("parse map field diamond_price 货币：价格 to value int32 failed.",
					zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
					// zap.String("field_data",data[18]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Diamond_price[key] = value
		}
	}

	// parse column 19 diamond_price_v2 : 货币：价格
	if data[19] != "" {

		config.Diamond_price_v2 = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[19], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field diamond_price_v2 货币：价格 to key int32 failed")
				logger.ErrorWF("parse map field diamond_price_v2 货币：价格 to key int32 failed.",
					zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
					// zap.String("field_data",data[19]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field diamond_price_v2 货币：价格 to value int32 failed")
				logger.ErrorWF("parse map field diamond_price_v2 货币：价格 to value int32 failed.",
					zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
					// zap.String("field_data",data[19]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Diamond_price_v2[key] = value
		}
	}

	// parse column 20 strength_pay : 夺取消耗体力
	if data[20] != "" {
		tmp, err = strconv.ParseInt(data[20], 10, 64)
		if err != nil {
			err = errors.New("parse field strength_pay 夺取消耗体力 to int32 failed")
			logger.ErrorWF("parse field strength_pay 夺取消耗体力 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[20]),
				zap.Error(err))
			return
		}
		config.Strength_pay = int32(tmp)
	}

	// parse column 21 protectiontime : 被夺后保护时间
	if data[21] != "" {
		tmp, err = strconv.ParseInt(data[21], 10, 64)
		if err != nil {
			err = errors.New("parse field protectiontime 被夺后保护时间 to int32 failed")
			logger.ErrorWF("parse field protectiontime 被夺后保护时间 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[21]),
				zap.Error(err))
			return
		}
		config.Protectiontime = int32(tmp)
	}

	// parse column 22 is_use : 是否可用
	if data[22] != "" {
		tmp, err = strconv.ParseInt(data[22], 10, 64)
		if err != nil {
			err = errors.New("parse field is_use 是否可用 to int32 failed")
			logger.ErrorWF("parse field is_use 是否可用 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[22]),
				zap.Error(err))
			return
		}
		config.Is_use = int32(tmp)
	}

	// parse column 23 if_popup : 是否新物品弹窗
	if data[23] != "" {
		tmp, err = strconv.ParseInt(data[23], 10, 64)
		if err != nil {
			err = errors.New("parse field if_popup 是否新物品弹窗 to int32 failed")
			logger.ErrorWF("parse field if_popup 是否新物品弹窗 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[23]),
				zap.Error(err))
			return
		}
		config.If_popup = int32(tmp)
	}

	// parse column 24 if_binding : 是否绑定可拍卖
	if data[24] != "" {
		tmp, err = strconv.ParseInt(data[24], 10, 64)
		if err != nil {
			err = errors.New("parse field if_binding 是否绑定可拍卖 to int32 failed")
			logger.ErrorWF("parse field if_binding 是否绑定可拍卖 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[24]),
				zap.Error(err))
			return
		}
		config.If_binding = int32(tmp)
	}

	// parse column 25 avaliable_time : 有效期（秒）
	if data[25] != "" {
		tmp, err = strconv.ParseInt(data[25], 10, 64)
		if err != nil {
			err = errors.New("parse field avaliable_time 有效期（秒） to int32 failed")
			logger.ErrorWF("parse field avaliable_time 有效期（秒） to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[25]),
				zap.Error(err))
			return
		}
		config.Avaliable_time = int32(tmp)
	}

	// parse column 26 jumpType : 跳转类型
	if data[26] != "" {

		vals := strings.Split(data[26], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field jumpType 跳转类型 to []int32 failed")
				logger.ErrorWF("parse array field jumpType 跳转类型 to []int32 failed.",
					zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
					// zap.String("field_data",data[26]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.JumpType = append(config.JumpType, int32(tmp))
		}
	}

	// parse column 27 jumpTarget : 跳转目标1
	if data[27] != "" {
		config.JumpTarget = data[27]
	}

	// parse column 28 source_type : 来源类型
	if data[28] != "" {
		tmp, err = strconv.ParseInt(data[28], 10, 64)
		if err != nil {
			err = errors.New("parse field source_type 来源类型 to int32 failed")
			logger.ErrorWF("parse field source_type 来源类型 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[28]),
				zap.Error(err))
			return
		}
		config.Source_type = int32(tmp)
	}

	// parse column 29 btn_type : 按钮动作
	if data[29] != "" {
		tmp, err = strconv.ParseInt(data[29], 10, 64)
		if err != nil {
			err = errors.New("parse field btn_type 按钮动作 to int32 failed")
			logger.ErrorWF("parse field btn_type 按钮动作 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[29]),
				zap.Error(err))
			return
		}
		config.Btn_type = int32(tmp)
	}

	// parse column 30 itemGettype : 道具添加类型
	if data[30] != "" {
		tmp, err = strconv.ParseInt(data[30], 10, 64)
		if err != nil {
			err = errors.New("parse field itemGettype 道具添加类型 to int32 failed")
			logger.ErrorWF("parse field itemGettype 道具添加类型 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[30]),
				zap.Error(err))
			return
		}
		config.ItemGettype = int32(tmp)
	}

	// parse column 31 is_pay : 是否有付费映射
	if data[31] != "" {
		tmp, err = strconv.ParseInt(data[31], 10, 64)
		if err != nil {
			err = errors.New("parse field is_pay 是否有付费映射 to int32 failed")
			logger.ErrorWF("parse field is_pay 是否有付费映射 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[31]),
				zap.Error(err))
			return
		}
		config.Is_pay = int32(tmp)
	}

	// parse column 32 price_v8 : 人偶版出售价格
	if data[32] != "" {
		tmp, err = strconv.ParseInt(data[32], 10, 64)
		if err != nil {
			err = errors.New("parse field price_v8 人偶版出售价格 to int32 failed")
			logger.ErrorWF("parse field price_v8 人偶版出售价格 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[32]),
				zap.Error(err))
			return
		}
		config.Price_v8 = int32(tmp)
	}

	// parse column 33 sale_id_v8 : 人偶版出售后得到货币的ID
	if data[33] != "" {
		tmp, err = strconv.ParseInt(data[33], 10, 64)
		if err != nil {
			err = errors.New("parse field sale_id_v8 人偶版出售后得到货币的ID to int32 failed")
			logger.ErrorWF("parse field sale_id_v8 人偶版出售后得到货币的ID to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[33]),
				zap.Error(err))
			return
		}
		config.Sale_id_v8 = int32(tmp)
	}

	// parse column 34 diamond_price_v8 : 货币：价格
	if data[34] != "" {

		config.Diamond_price_v8 = make(map[int32]int32)
		var key int32
		var value int32
		vals := strings.Split(data[34], "_")
		for k, val := range vals {
			items := strings.Split(val, ":")
			tmp, err = strconv.ParseInt(items[0], 10, 64)
			if err != nil {
				err = errors.New("parse map field diamond_price_v8 货币：价格 to key int32 failed")
				logger.ErrorWF("parse map field diamond_price_v8 货币：价格 to key int32 failed.",
					zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
					// zap.String("field_data",data[34]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[0]),
					zap.Error(err))
				return
			}
			key = int32(tmp)
			tmp, err = strconv.ParseInt(items[1], 10, 64)
			if err != nil {
				err = errors.New("parse map field diamond_price_v8 货币：价格 to value int32 failed")
				logger.ErrorWF("parse map field diamond_price_v8 货币：价格 to value int32 failed.",
					zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
					// zap.String("field_data",data[34]),
					zap.String("item_data", val), zap.Int("index", k),
					zap.String("parse_data", items[1]),
					zap.Error(err))
				return
			}
			value = int32(tmp)
			config.Diamond_price_v8[key] = value
		}
	}

	// parse column 35 overtime_type : 倒计时类型（0-默认 1-进船背包
	if data[35] != "" {
		tmp, err = strconv.ParseInt(data[35], 10, 64)
		if err != nil {
			err = errors.New("parse field overtime_type 倒计时类型（0-默认 1-进船背包 to int32 failed")
			logger.ErrorWF("parse field overtime_type 倒计时类型（0-默认 1-进船背包 to int32 failed.",
				zap.String("xlsx", "maze_items_v8【迷宫-道具】.xlsx"), zap.String("sheet", "maze_items_v8"),
				zap.String("parse_data", data[35]),
				zap.Error(err))
			return
		}
		config.Overtime_type = int32(tmp)
	}
	return
}

var gMazeItemsV8Fields = []string{
	"key_id",
	"type",
	"id",
	"prop_name",
	"name_desc",
	"is_native",
	"pic_token",
	"iconAtlas",
	"icon",
	"quality",
	"is_sell",
	"price",
	"sale_id",
	"price_v2",
	"sale_id_v2",
	"is_bag",
	"fly_type",
	"is_rob",
	"diamond_price",
	"diamond_price_v2",
	"strength_pay",
	"protectiontime",
	"is_use",
	"if_popup",
	"if_binding",
	"avaliable_time",
	"jumpType",
	"jumpTarget",
	"source_type",
	"btn_type",
	"itemGettype",
	"is_pay",
	"price_v8",
	"sale_id_v8",
	"diamond_price_v8",
	"overtime_type",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeItemsV8Parser{}
	loader := &gMazeItemsV8Loader{}
	var data [][]string
	data, err = load("maze_items_v8【迷宫-道具】.xlsx", "maze_items_v8", gMazeItemsV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_items_v8【迷宫-道具】.xlsx maze_items_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_items_v8【迷宫-道具】.xlsx maze_items_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeItemsV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_items_v8【迷宫-道具】.xlsx maze_items_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeItemsV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_items_v8【迷宫-道具】.xlsx maze_items_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_items_v8【迷宫-道具】.xlsx maze_items_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_items_v8【迷宫-道具】.xlsx maze_items_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_items_v8【迷宫-道具】.xlsx maze_items_v8 data success.")
	return
}
