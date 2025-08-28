package GMazeToastMsgInfoCfg

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

// MazeToastMsgInfoConfigRow from maze_toast_msg_info【迷宫-全局消息通知】.xlsx maze_toast_msg_info
type MazeToastMsgInfoConfigRow struct {
	Msg_id        int32   `json:"msg_id"`        // 消息id
	Msg_text      string  `json:"msg_text"`      // 消息内容（没有格式化内容）
	Display_time  int32   `json:"display_time"`  // 展示时长（秒）时间=0表示不会展示信息
	Msg_order     int32   `json:"msg_order"`     // 展示优先级
	Text_wildcard []int32 `json:"text_wildcard"` // 内容通配符id（废字段）
}

// MazeToastMsgInfoConfig from maze_toast_msg_info【迷宫-全局消息通知】.xlsx maze_toast_msg_info
type MazeToastMsgInfoConfig struct {
	ConfigRows map[int32]*MazeToastMsgInfoConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeToastMsgInfoConfig {
	ret := &MazeToastMsgInfoConfig{ConfigRows: map[int32]*MazeToastMsgInfoConfigRow{}}
	return ret
}

// GetMazeToastMsgInfoConfig get one config by configId
func (c *MazeToastMsgInfoConfig) GetMazeToastMsgInfoConfig(configId int32) *MazeToastMsgInfoConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeToastMsgInfoConfig) Get(configId int32) *MazeToastMsgInfoConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeToastMsgInfoConfig get all config slice
func (c *MazeToastMsgInfoConfig) GetAllMazeToastMsgInfoConfig() (res []*MazeToastMsgInfoConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeToastMsgInfoConfig) GetAll() (res []*MazeToastMsgInfoConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeToastMsgInfoConfig

// GetMazeToastMsgInfoConfig pkg func. get one config by configId
func GetMazeToastMsgInfoConfig(configId int32) *MazeToastMsgInfoConfigRow {
	return gConfigData.GetMazeToastMsgInfoConfig(configId)
}

// Deprecated: 链路追踪信息缺失。推荐使用GetWithCtx
// Get pkg func. get one config by configId
func Get(configId int32) *MazeToastMsgInfoConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *MazeToastMsgInfoConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "maze_toast_msg_info", configId, otps...)
	}
	return cfg
}

// GetAllMazeToastMsgInfoConfig pkg func. get all config slice
func GetAllMazeToastMsgInfoConfig() []*MazeToastMsgInfoConfigRow {
	return gConfigData.GetAllMazeToastMsgInfoConfig()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeToastMsgInfoConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeToastMsgInfoConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeToastMsgInfoConfigRow from maze_toast_msg_info【迷宫-全局消息通知】.xlsx maze_toast_msg_info"
}

// GetRawValue get raw data
func GetRawValue() *MazeToastMsgInfoConfig {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_toast_msg_info"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_toast_msg_info.json",
		"maze_toast_msg_info【迷宫-全局消息通知】.xlsx", "maze_toast_msg_info",
		&gMazeToastMsgInfoParser{}, &gMazeToastMsgInfoLoader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeToastMsgInfoConfig)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeToastMsgInfoConfig) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeToastMsgInfoConfig))(c)
		return true
	})
}

// RegisterMazeToastMsgInfoInitCallBack reg config update func (old func)
var RegisterMazeToastMsgInfoInitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeToastMsgInfoConfig)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeToastMsgInfoConfig)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeToastMsgInfoConfig) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeToastMsgInfoConfig) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeToastMsgInfoConfig) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeToastMsgInfoLoader struct {
}

// NewContainer new data container pointer
func (*gMazeToastMsgInfoLoader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeToastMsgInfoLoader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeToastMsgInfoConfig)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeToastMsgInfoLoader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeToastMsgInfoConfig)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeToastMsgInfoLoader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeToastMsgInfoConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeToastMsgInfoConfigRow")
		logger.ErrorWF("invalid type. not *MazeToastMsgInfoConfigRow", zap.String("xlsx", "maze_toast_msg_info【迷宫-全局消息通知】.xlsx"),
			zap.String("sheet", "maze_toast_msg_info"))
		return
	}
	config, ok := container.(*MazeToastMsgInfoConfig)
	if !ok {
		err = errors.New("invalid type. not *MazeToastMsgInfoConfig")
		logger.ErrorWF("invalid type. not *MazeToastMsgInfoConfig", zap.String("xlsx", "maze_toast_msg_info【迷宫-全局消息通知】.xlsx"),
			zap.String("sheet", "maze_toast_msg_info"))
		return
	}
	config.ConfigRows[row.Msg_id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeToastMsgInfoLoader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeToastMsgInfoConfig)
	if !ok {
		err = errors.New("invalid type. not *MazeToastMsgInfoConfig")
		logger.ErrorWF("invalid type. not *MazeToastMsgInfoConfig", zap.String("xlsx", "maze_toast_msg_info【迷宫-全局消息通知】.xlsx"),
			zap.String("sheet", "maze_toast_msg_info"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeToastMsgInfoLoader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeToastMsgInfoConfig)
	if !ok {
		err = errors.New("invalid type. not *MazeToastMsgInfoConfig")
		logger.ErrorWF("invalid type. not *MazeToastMsgInfoConfig", zap.String("xlsx", "maze_toast_msg_info【迷宫-全局消息通知】.xlsx"),
			zap.String("sheet", "maze_toast_msg_info"))
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
type gMazeToastMsgInfoParser struct {
}

// New new config row data
func (*gMazeToastMsgInfoParser) New() interface{} {
	return &MazeToastMsgInfoConfigRow{}
}

// Fields get config fields names
func (*gMazeToastMsgInfoParser) Fields() []string {
	return gMazeToastMsgInfoFields
}

// Parse parse raw data to row data
func (*gMazeToastMsgInfoParser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeToastMsgInfoConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeToastMsgInfoConfigRow")
		logger.ErrorWF("invalid type. not *MazeToastMsgInfoConfigRow", zap.String("xlsx", "maze_toast_msg_info【迷宫-全局消息通知】.xlsx"),
			zap.String("sheet", "maze_toast_msg_info"))
		return
	}
	// compare length
	if len(data) != len(gMazeToastMsgInfoFields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeToastMsgInfoConfigRow",
			zap.String("xlsx", "maze_toast_msg_info【迷宫-全局消息通知】.xlsx"),
			zap.String("sheet", "maze_toast_msg_info"), zap.Int("need_count", len(gMazeToastMsgInfoFields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 msg_id : 消息id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field msg_id 消息id to int32 failed")
			logger.ErrorWF("parse field msg_id 消息id to int32 failed.",
				zap.String("xlsx", "maze_toast_msg_info【迷宫-全局消息通知】.xlsx"), zap.String("sheet", "maze_toast_msg_info"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Msg_id = int32(tmp)
	}

	// parse column 1 msg_text : 消息内容（没有格式化内容）
	if data[1] != "" {
		config.Msg_text = data[1]
	}

	// parse column 2 display_time : 展示时长（秒）时间=0表示不会展示信息
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field display_time 展示时长（秒）时间=0表示不会展示信息 to int32 failed")
			logger.ErrorWF("parse field display_time 展示时长（秒）时间=0表示不会展示信息 to int32 failed.",
				zap.String("xlsx", "maze_toast_msg_info【迷宫-全局消息通知】.xlsx"), zap.String("sheet", "maze_toast_msg_info"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Display_time = int32(tmp)
	}

	// parse column 3 msg_order : 展示优先级
	if data[3] != "" {
		tmp, err = strconv.ParseInt(data[3], 10, 64)
		if err != nil {
			err = errors.New("parse field msg_order 展示优先级 to int32 failed")
			logger.ErrorWF("parse field msg_order 展示优先级 to int32 failed.",
				zap.String("xlsx", "maze_toast_msg_info【迷宫-全局消息通知】.xlsx"), zap.String("sheet", "maze_toast_msg_info"),
				zap.String("parse_data", data[3]),
				zap.Error(err))
			return
		}
		config.Msg_order = int32(tmp)
	}

	// parse column 4 text_wildcard : 内容通配符id（废字段）
	if data[4] != "" {

		vals := strings.Split(data[4], ",")
		for k, v := range vals {
			tmp, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				err = errors.New("parse array field text_wildcard 内容通配符id（废字段） to []int32 failed")
				logger.ErrorWF("parse array field text_wildcard 内容通配符id（废字段） to []int32 failed.",
					zap.String("xlsx", "maze_toast_msg_info【迷宫-全局消息通知】.xlsx"), zap.String("sheet", "maze_toast_msg_info"),
					// zap.String("field_data",data[4]),
					zap.String("parse_data", v), zap.Int("index", k),
					zap.Error(err))
				return
			}
			config.Text_wildcard = append(config.Text_wildcard, int32(tmp))
		}
	}
	return
}

var gMazeToastMsgInfoFields = []string{
	"msg_id",
	"msg_text",
	"display_time",
	"msg_order",
	"text_wildcard",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeToastMsgInfoParser{}
	loader := &gMazeToastMsgInfoLoader{}
	var data [][]string
	data, err = load("maze_toast_msg_info【迷宫-全局消息通知】.xlsx", "maze_toast_msg_info", gMazeToastMsgInfoFields)
	if err != nil {
		logger.ErrorWF("load maze_toast_msg_info【迷宫-全局消息通知】.xlsx maze_toast_msg_info data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_toast_msg_info【迷宫-全局消息通知】.xlsx maze_toast_msg_info data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeToastMsgInfoFields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_toast_msg_info【迷宫-全局消息通知】.xlsx maze_toast_msg_info failed.", zap.Int("row", k), zap.Strings("need", gMazeToastMsgInfoFields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_toast_msg_info【迷宫-全局消息通知】.xlsx maze_toast_msg_info row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_toast_msg_info【迷宫-全局消息通知】.xlsx maze_toast_msg_info row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_toast_msg_info【迷宫-全局消息通知】.xlsx maze_toast_msg_info data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_toast_msg_info【迷宫-全局消息通知】.xlsx maze_toast_msg_info data success.")
	return
}
