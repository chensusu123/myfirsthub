package GCombatEffectivenessCfg

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

// CombatEffectivenessConfigRow from attribute【属性】.xlsx combat_effectiveness
type CombatEffectivenessConfigRow struct {
	Id                       int32 `json:"id"`                       // 评级ID
	Combat_effectiveness_min int32 `json:"combat_effectiveness_min"` // 战力比范围下限
	Combat_effectiveness_max int32 `json:"combat_effectiveness_max"` // 战力比范围上限
}

// CombatEffectivenessConfig from attribute【属性】.xlsx combat_effectiveness
type CombatEffectivenessConfig struct {
	ConfigRows map[int32]*CombatEffectivenessConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *CombatEffectivenessConfig {
	ret := &CombatEffectivenessConfig{ConfigRows: map[int32]*CombatEffectivenessConfigRow{}}
	return ret
}

// GetCombatEffectivenessConfig get one config by configId
func (c *CombatEffectivenessConfig) GetCombatEffectivenessConfig(configId int32) *CombatEffectivenessConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *CombatEffectivenessConfig) Get(configId int32) *CombatEffectivenessConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllCombatEffectivenessConfig get all config slice
func (c *CombatEffectivenessConfig) GetAllCombatEffectivenessConfig() (res []*CombatEffectivenessConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *CombatEffectivenessConfig) GetAll() (res []*CombatEffectivenessConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *CombatEffectivenessConfig

// GetCombatEffectivenessConfig pkg func. get one config by configId
func GetCombatEffectivenessConfig(configId int32) *CombatEffectivenessConfigRow {
	return gConfigData.GetCombatEffectivenessConfig(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *CombatEffectivenessConfigRow {
	return gConfigData.Get(configId)
}

// GetAllCombatEffectivenessConfig pkg func. get all config slice
func GetAllCombatEffectivenessConfig() []*CombatEffectivenessConfigRow {
	return gConfigData.GetAllCombatEffectivenessConfig()
}

// GetAll pkg func. get all config slice
func GetAll() []*CombatEffectivenessConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*CombatEffectivenessConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "CombatEffectivenessConfigRow from attribute【属性】.xlsx combat_effectiveness"
}

// GetRawValue get raw data
func GetRawValue() *CombatEffectivenessConfig {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "combat_effectiveness"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("combat_effectiveness.json",
		"attribute【属性】.xlsx", "combat_effectiveness",
		&gCombatEffectivenessParser{}, &gCombatEffectivenessLoader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*CombatEffectivenessConfig)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *CombatEffectivenessConfig) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*CombatEffectivenessConfig))(c)
		return true
	})
}

// RegisterCombatEffectivenessInitCallBack reg config update func (old func)
var RegisterCombatEffectivenessInitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*CombatEffectivenessConfig)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*CombatEffectivenessConfig)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *CombatEffectivenessConfig) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*CombatEffectivenessConfig) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*CombatEffectivenessConfig) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gCombatEffectivenessLoader struct {
}

// NewContainer new data container pointer
func (*gCombatEffectivenessLoader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gCombatEffectivenessLoader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*CombatEffectivenessConfig)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gCombatEffectivenessLoader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*CombatEffectivenessConfig)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gCombatEffectivenessLoader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*CombatEffectivenessConfigRow)
	if !ok {
		err = errors.New("invalid type. not *CombatEffectivenessConfigRow")
		logger.ErrorWF("invalid type. not *CombatEffectivenessConfigRow", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "combat_effectiveness"))
		return
	}
	config, ok := container.(*CombatEffectivenessConfig)
	if !ok {
		err = errors.New("invalid type. not *CombatEffectivenessConfig")
		logger.ErrorWF("invalid type. not *CombatEffectivenessConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "combat_effectiveness"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gCombatEffectivenessLoader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*CombatEffectivenessConfig)
	if !ok {
		err = errors.New("invalid type. not *CombatEffectivenessConfig")
		logger.ErrorWF("invalid type. not *CombatEffectivenessConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "combat_effectiveness"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gCombatEffectivenessLoader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*CombatEffectivenessConfig)
	if !ok {
		err = errors.New("invalid type. not *CombatEffectivenessConfig")
		logger.ErrorWF("invalid type. not *CombatEffectivenessConfig", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "combat_effectiveness"))
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
type gCombatEffectivenessParser struct {
}

// New new config row data
func (*gCombatEffectivenessParser) New() interface{} {
	return &CombatEffectivenessConfigRow{}
}

// Fields get config fields names
func (*gCombatEffectivenessParser) Fields() []string {
	return gCombatEffectivenessFields
}

// Parse parse raw data to row data
func (*gCombatEffectivenessParser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*CombatEffectivenessConfigRow)
	if !ok {
		err = errors.New("invalid type. not *CombatEffectivenessConfigRow")
		logger.ErrorWF("invalid type. not *CombatEffectivenessConfigRow", zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "combat_effectiveness"))
		return
	}
	// compare length
	if len(data) != len(gCombatEffectivenessFields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*CombatEffectivenessConfigRow",
			zap.String("xlsx", "attribute【属性】.xlsx"),
			zap.String("sheet", "combat_effectiveness"), zap.Int("need_count", len(gCombatEffectivenessFields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 评级ID
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 评级ID to int32 failed")
			logger.ErrorWF("parse field id 评级ID to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "combat_effectiveness"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 combat_effectiveness_min : 战力比范围下限
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field combat_effectiveness_min 战力比范围下限 to int32 failed")
			logger.ErrorWF("parse field combat_effectiveness_min 战力比范围下限 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "combat_effectiveness"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Combat_effectiveness_min = int32(tmp)
	}

	// parse column 2 combat_effectiveness_max : 战力比范围上限
	if data[2] != "" {
		tmp, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			err = errors.New("parse field combat_effectiveness_max 战力比范围上限 to int32 failed")
			logger.ErrorWF("parse field combat_effectiveness_max 战力比范围上限 to int32 failed.",
				zap.String("xlsx", "attribute【属性】.xlsx"), zap.String("sheet", "combat_effectiveness"),
				zap.String("parse_data", data[2]),
				zap.Error(err))
			return
		}
		config.Combat_effectiveness_max = int32(tmp)
	}
	return
}

var gCombatEffectivenessFields = []string{
	"id",
	"combat_effectiveness_min",
	"combat_effectiveness_max",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gCombatEffectivenessParser{}
	loader := &gCombatEffectivenessLoader{}
	var data [][]string
	data, err = load("attribute【属性】.xlsx", "combat_effectiveness", gCombatEffectivenessFields)
	if err != nil {
		logger.ErrorWF("load attribute【属性】.xlsx combat_effectiveness data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load attribute【属性】.xlsx combat_effectiveness data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gCombatEffectivenessFields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse attribute【属性】.xlsx combat_effectiveness failed.", zap.Int("row", k), zap.Strings("need", gCombatEffectivenessFields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse attribute【属性】.xlsx combat_effectiveness row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add attribute【属性】.xlsx combat_effectiveness row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check attribute【属性】.xlsx combat_effectiveness data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load attribute【属性】.xlsx combat_effectiveness data success.")
	return
}
