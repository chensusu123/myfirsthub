package GGitVersionV8Cfg

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

// GitVersionV8ConfigRow from git_version_v8【配表版本号】.xlsx git_version_v8
type GitVersionV8ConfigRow struct {
	ID          int32  `json:"ID"`          // 序号
	Git_version string `json:"git_version"` // 版本
	Update_time string `json:"update_time"` // 更新时间
}

// GitVersionV8Config from git_version_v8【配表版本号】.xlsx git_version_v8
type GitVersionV8Config struct {
	ConfigRows map[int32]*GitVersionV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *GitVersionV8Config {
	ret := &GitVersionV8Config{ConfigRows: map[int32]*GitVersionV8ConfigRow{}}
	return ret
}

// GetGitVersionV8Config get one config by configId
func (c *GitVersionV8Config) GetGitVersionV8Config(configId int32) *GitVersionV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *GitVersionV8Config) Get(configId int32) *GitVersionV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllGitVersionV8Config get all config slice
func (c *GitVersionV8Config) GetAllGitVersionV8Config() (res []*GitVersionV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *GitVersionV8Config) GetAll() (res []*GitVersionV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *GitVersionV8Config

// GetGitVersionV8Config pkg func. get one config by configId
func GetGitVersionV8Config(configId int32) *GitVersionV8ConfigRow {
	return gConfigData.GetGitVersionV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *GitVersionV8ConfigRow {
	return GetWithCtx(context.Background(), configId)
}

// GetWithCtx pkg func. get one config by configId
func GetWithCtx(ctx context.Context, configId int32, otps ...config_manager.QueryOption) *GitVersionV8ConfigRow {
	cfg := gConfigData.Get(configId)
	if cfg == nil {
		config_manager.MissRecord(ctx, "git_version_v8", configId, otps...)
	}
	return cfg
}

// GetAllGitVersionV8Config pkg func. get all config slice
func GetAllGitVersionV8Config() []*GitVersionV8ConfigRow {
	return gConfigData.GetAllGitVersionV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*GitVersionV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*GitVersionV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "GitVersionV8ConfigRow from git_version_v8【配表版本号】.xlsx git_version_v8"
}

// GetRawValue get raw data
func GetRawValue() *GitVersionV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "git_version_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("git_version_v8.json",
		"git_version_v8【配表版本号】.xlsx", "git_version_v8",
		&gGitVersionV8Parser{}, &gGitVersionV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*GitVersionV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *GitVersionV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*GitVersionV8Config))(c)
		return true
	})
}

// RegisterGitVersionV8InitCallBack reg config update func (old func)
var RegisterGitVersionV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*GitVersionV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*GitVersionV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *GitVersionV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*GitVersionV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*GitVersionV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gGitVersionV8Loader struct {
}

// NewContainer new data container pointer
func (*gGitVersionV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gGitVersionV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*GitVersionV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gGitVersionV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*GitVersionV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gGitVersionV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*GitVersionV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *GitVersionV8ConfigRow")
		logger.ErrorWF("invalid type. not *GitVersionV8ConfigRow", zap.String("xlsx", "git_version_v8【配表版本号】.xlsx"),
			zap.String("sheet", "git_version_v8"))
		return
	}
	config, ok := container.(*GitVersionV8Config)
	if !ok {
		err = errors.New("invalid type. not *GitVersionV8Config")
		logger.ErrorWF("invalid type. not *GitVersionV8Config", zap.String("xlsx", "git_version_v8【配表版本号】.xlsx"),
			zap.String("sheet", "git_version_v8"))
		return
	}
	config.ConfigRows[row.ID] = row
	return
}

// GetValue get real map value for json parse
func (*gGitVersionV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*GitVersionV8Config)
	if !ok {
		err = errors.New("invalid type. not *GitVersionV8Config")
		logger.ErrorWF("invalid type. not *GitVersionV8Config", zap.String("xlsx", "git_version_v8【配表版本号】.xlsx"),
			zap.String("sheet", "git_version_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gGitVersionV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*GitVersionV8Config)
	if !ok {
		err = errors.New("invalid type. not *GitVersionV8Config")
		logger.ErrorWF("invalid type. not *GitVersionV8Config", zap.String("xlsx", "git_version_v8【配表版本号】.xlsx"),
			zap.String("sheet", "git_version_v8"))
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
type gGitVersionV8Parser struct {
}

// New new config row data
func (*gGitVersionV8Parser) New() interface{} {
	return &GitVersionV8ConfigRow{}
}

// Fields get config fields names
func (*gGitVersionV8Parser) Fields() []string {
	return gGitVersionV8Fields
}

// Parse parse raw data to row data
func (*gGitVersionV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*GitVersionV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *GitVersionV8ConfigRow")
		logger.ErrorWF("invalid type. not *GitVersionV8ConfigRow", zap.String("xlsx", "git_version_v8【配表版本号】.xlsx"),
			zap.String("sheet", "git_version_v8"))
		return
	}
	// compare length
	if len(data) != len(gGitVersionV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*GitVersionV8ConfigRow",
			zap.String("xlsx", "git_version_v8【配表版本号】.xlsx"),
			zap.String("sheet", "git_version_v8"), zap.Int("need_count", len(gGitVersionV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 ID : 序号
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field ID 序号 to int32 failed")
			logger.ErrorWF("parse field ID 序号 to int32 failed.",
				zap.String("xlsx", "git_version_v8【配表版本号】.xlsx"), zap.String("sheet", "git_version_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.ID = int32(tmp)
	}

	// parse column 1 git_version : 版本
	if data[1] != "" {
		config.Git_version = data[1]
	}

	// parse column 2 update_time : 更新时间
	if data[2] != "" {
		config.Update_time = data[2]
	}
	return
}

var gGitVersionV8Fields = []string{
	"ID",
	"git_version",
	"update_time",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gGitVersionV8Parser{}
	loader := &gGitVersionV8Loader{}
	var data [][]string
	data, err = load("git_version_v8【配表版本号】.xlsx", "git_version_v8", gGitVersionV8Fields)
	if err != nil {
		logger.ErrorWF("load git_version_v8【配表版本号】.xlsx git_version_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load git_version_v8【配表版本号】.xlsx git_version_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gGitVersionV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse git_version_v8【配表版本号】.xlsx git_version_v8 failed.", zap.Int("row", k), zap.Strings("need", gGitVersionV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse git_version_v8【配表版本号】.xlsx git_version_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add git_version_v8【配表版本号】.xlsx git_version_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check git_version_v8【配表版本号】.xlsx git_version_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load git_version_v8【配表版本号】.xlsx git_version_v8 data success.")
	return
}
