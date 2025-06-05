package GMazeAnswerMechanismV8Cfg

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

// MazeAnswerMechanismV8ConfigRow from maze_answer_mechanism_v8【答题机关词库】.xlsx maze_answer_mechanism_v8
type MazeAnswerMechanismV8ConfigRow struct {
	Id               int32  `json:"id"`               // 唯一id
	Question_bank_id int32  `json:"question_bank_id"` // 题库id
	Question_stem    string `json:"question_stem"`    // 题干
	Answer_a         string `json:"answer_a"`         // 选项a
	Answer_b         string `json:"answer_b"`         // 选项b
	Answer_c         string `json:"answer_c"`         // 选项c
	Answer_d         string `json:"answer_d"`         // 选项d
	Right_answer     int32  `json:"right_answer"`     // 正确选项
}

// MazeAnswerMechanismV8Config from maze_answer_mechanism_v8【答题机关词库】.xlsx maze_answer_mechanism_v8
type MazeAnswerMechanismV8Config struct {
	ConfigRows map[int32]*MazeAnswerMechanismV8ConfigRow
	// lock       sync.RWMutex
}

// newConfig new one config
func newConfig() *MazeAnswerMechanismV8Config {
	ret := &MazeAnswerMechanismV8Config{ConfigRows: map[int32]*MazeAnswerMechanismV8ConfigRow{}}
	return ret
}

// GetMazeAnswerMechanismV8Config get one config by configId
func (c *MazeAnswerMechanismV8Config) GetMazeAnswerMechanismV8Config(configId int32) *MazeAnswerMechanismV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// Get get one config by configId
func (c *MazeAnswerMechanismV8Config) Get(configId int32) *MazeAnswerMechanismV8ConfigRow {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	if cnf, ok := c.ConfigRows[configId]; ok {
		return cnf
	}
	return nil
}

// GetAllMazeAnswerMechanismV8Config get all config slice
func (c *MazeAnswerMechanismV8Config) GetAllMazeAnswerMechanismV8Config() (res []*MazeAnswerMechanismV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// GetAll get all config slice
func (c *MazeAnswerMechanismV8Config) GetAll() (res []*MazeAnswerMechanismV8ConfigRow) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	for _, val := range c.ConfigRows {
		res = append(res, val)
	}
	return
}

// global config pointer
var gConfigData *MazeAnswerMechanismV8Config

// GetMazeAnswerMechanismV8Config pkg func. get one config by configId
func GetMazeAnswerMechanismV8Config(configId int32) *MazeAnswerMechanismV8ConfigRow {
	return gConfigData.GetMazeAnswerMechanismV8Config(configId)
}

// Get pkg func. get one config by configId
func Get(configId int32) *MazeAnswerMechanismV8ConfigRow {
	return gConfigData.Get(configId)
}

// GetAllMazeAnswerMechanismV8Config pkg func. get all config slice
func GetAllMazeAnswerMechanismV8Config() []*MazeAnswerMechanismV8ConfigRow {
	return gConfigData.GetAllMazeAnswerMechanismV8Config()
}

// GetAll pkg func. get all config slice
func GetAll() []*MazeAnswerMechanismV8ConfigRow {
	return gConfigData.GetAll()
}

// ConfigRows get raw map data
func ConfigRows() map[int32]*MazeAnswerMechanismV8ConfigRow {
	return gConfigData.ConfigRows
}

// GetConfigDesc get config desc for lod,debug etc.
func GetConfigDesc() string {
	return "MazeAnswerMechanismV8ConfigRow from maze_answer_mechanism_v8【答题机关词库】.xlsx maze_answer_mechanism_v8"
}

// GetRawValue get raw data
func GetRawValue() *MazeAnswerMechanismV8Config {
	return gConfigData
}

// SheetName get sheet name
func SheetName() string {
	return "maze_answer_mechanism_v8"
}

func init() {
	// reg config auto load
	config_manager.RegAutoConfig("maze_answer_mechanism_v8.json",
		"maze_answer_mechanism_v8【答题机关词库】.xlsx", "maze_answer_mechanism_v8",
		&gMazeAnswerMechanismV8Parser{}, &gMazeAnswerMechanismV8Loader{})
}

// data update call back
var cfgUpdateCallBack sync.Map //map[string]func(*MazeAnswerMechanismV8Config)

// doConfigUpdateCallback do config update callback
func doConfigUpdateCallback(c *MazeAnswerMechanismV8Config) {
	cfgUpdateCallBack.Range(func(key, value interface{}) bool {
		value.(func(*MazeAnswerMechanismV8Config))(c)
		return true
	})
}

// RegisterMazeAnswerMechanismV8InitCallBack reg config update func (old func)
var RegisterMazeAnswerMechanismV8InitCallBack = RegisterLoadedCallBack

// RegisterLoadedCallBack reg config update func.
func RegisterLoadedCallBack(key string, f func(*MazeAnswerMechanismV8Config)) {
	cfgUpdateCallBack.Store(key, f)
}

// data check call back
var cfgCheckCallBack sync.Map //map[string]func(*MazeAnswerMechanismV8Config)error

// doConfigCheckCallback do config check callback
func doConfigCheckCallback(c *MazeAnswerMechanismV8Config) (err error) {
	cfgCheckCallBack.Range(func(key, value interface{}) bool {
		err := value.(func(*MazeAnswerMechanismV8Config) error)(c)
		return err == nil
	})
	return
}

// RegisterLoadedCallBack reg config update func.
func RegisterConfigCheck(key string, f func(*MazeAnswerMechanismV8Config) error) {
	cfgCheckCallBack.Store(key, f)
}

// implete ConfigLoader interface
type gMazeAnswerMechanismV8Loader struct {
}

// NewContainer new data container pointer
func (*gMazeAnswerMechanismV8Loader) NewContainer() interface{} {
	return newConfig()
}

// Check check new config data ptr
func (*gMazeAnswerMechanismV8Loader) Check(newPtr interface{}) error {
	// set global ptr
	cfgData := newPtr.(*MazeAnswerMechanismV8Config)
	return doConfigCheckCallback(cfgData)
}

// Swap swap global config data ptr
func (*gMazeAnswerMechanismV8Loader) Swap(newPtr interface{}) {
	// convert pointer
	cache := newPtr.(*MazeAnswerMechanismV8Config)
	// update second edit
	doConfigUpdateCallback(cache)
	// set global ptr
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&gConfigData)), unsafe.Pointer(cache))
}

// Add append config item to map,save result
func (*gMazeAnswerMechanismV8Loader) Add(logger fklog.FKLogI, container interface{}, ri interface{}) (err error) {
	row, ok := ri.(*MazeAnswerMechanismV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAnswerMechanismV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAnswerMechanismV8ConfigRow", zap.String("xlsx", "maze_answer_mechanism_v8【答题机关词库】.xlsx"),
			zap.String("sheet", "maze_answer_mechanism_v8"))
		return
	}
	config, ok := container.(*MazeAnswerMechanismV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAnswerMechanismV8Config")
		logger.ErrorWF("invalid type. not *MazeAnswerMechanismV8Config", zap.String("xlsx", "maze_answer_mechanism_v8【答题机关词库】.xlsx"),
			zap.String("sheet", "maze_answer_mechanism_v8"))
		return
	}
	config.ConfigRows[row.Id] = row
	return
}

// GetValue get real map value for json parse
func (*gMazeAnswerMechanismV8Loader) GetValue(logger fklog.FKLogI, container interface{}) (real interface{}, err error) {
	config, ok := container.(*MazeAnswerMechanismV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAnswerMechanismV8Config")
		logger.ErrorWF("invalid type. not *MazeAnswerMechanismV8Config", zap.String("xlsx", "maze_answer_mechanism_v8【答题机关词库】.xlsx"),
			zap.String("sheet", "maze_answer_mechanism_v8"))
		return
	}
	real = &config.ConfigRows
	return
}

// Range range all data for json append data parse.
func (*gMazeAnswerMechanismV8Loader) Range(logger fklog.FKLogI, container interface{}, rf func(row interface{}) error) (err error) {
	config, ok := container.(*MazeAnswerMechanismV8Config)
	if !ok {
		err = errors.New("invalid type. not *MazeAnswerMechanismV8Config")
		logger.ErrorWF("invalid type. not *MazeAnswerMechanismV8Config", zap.String("xlsx", "maze_answer_mechanism_v8【答题机关词库】.xlsx"),
			zap.String("sheet", "maze_answer_mechanism_v8"))
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
type gMazeAnswerMechanismV8Parser struct {
}

// New new config row data
func (*gMazeAnswerMechanismV8Parser) New() interface{} {
	return &MazeAnswerMechanismV8ConfigRow{}
}

// Fields get config fields names
func (*gMazeAnswerMechanismV8Parser) Fields() []string {
	return gMazeAnswerMechanismV8Fields
}

// Parse parse raw data to row data
func (*gMazeAnswerMechanismV8Parser) Parse(logger fklog.FKLogI, data []string, row interface{}) (err error) {
	// convert row type
	config, ok := row.(*MazeAnswerMechanismV8ConfigRow)
	if !ok {
		err = errors.New("invalid type. not *MazeAnswerMechanismV8ConfigRow")
		logger.ErrorWF("invalid type. not *MazeAnswerMechanismV8ConfigRow", zap.String("xlsx", "maze_answer_mechanism_v8【答题机关词库】.xlsx"),
			zap.String("sheet", "maze_answer_mechanism_v8"))
		return
	}
	// compare length
	if len(data) != len(gMazeAnswerMechanismV8Fields) {
		err = errors.New("fields count not match.")
		logger.ErrorWF("invalid type. not *map[int32]*MazeAnswerMechanismV8ConfigRow",
			zap.String("xlsx", "maze_answer_mechanism_v8【答题机关词库】.xlsx"),
			zap.String("sheet", "maze_answer_mechanism_v8"), zap.Int("need_count", len(gMazeAnswerMechanismV8Fields)),
			zap.Int("had_count", len(data)))
		return
	}

	var tmp int64

	// parse column 0 id : 唯一id
	if data[0] != "" {
		tmp, err = strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			err = errors.New("parse field id 唯一id to int32 failed")
			logger.ErrorWF("parse field id 唯一id to int32 failed.",
				zap.String("xlsx", "maze_answer_mechanism_v8【答题机关词库】.xlsx"), zap.String("sheet", "maze_answer_mechanism_v8"),
				zap.String("parse_data", data[0]),
				zap.Error(err))
			return
		}
		config.Id = int32(tmp)
	}

	// parse column 1 question_bank_id : 题库id
	if data[1] != "" {
		tmp, err = strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			err = errors.New("parse field question_bank_id 题库id to int32 failed")
			logger.ErrorWF("parse field question_bank_id 题库id to int32 failed.",
				zap.String("xlsx", "maze_answer_mechanism_v8【答题机关词库】.xlsx"), zap.String("sheet", "maze_answer_mechanism_v8"),
				zap.String("parse_data", data[1]),
				zap.Error(err))
			return
		}
		config.Question_bank_id = int32(tmp)
	}

	// parse column 2 question_stem : 题干
	if data[2] != "" {
		config.Question_stem = data[2]
	}

	// parse column 3 answer_a : 选项a
	if data[3] != "" {
		config.Answer_a = data[3]
	}

	// parse column 4 answer_b : 选项b
	if data[4] != "" {
		config.Answer_b = data[4]
	}

	// parse column 5 answer_c : 选项c
	if data[5] != "" {
		config.Answer_c = data[5]
	}

	// parse column 6 answer_d : 选项d
	if data[6] != "" {
		config.Answer_d = data[6]
	}

	// parse column 7 right_answer : 正确选项
	if data[7] != "" {
		tmp, err = strconv.ParseInt(data[7], 10, 64)
		if err != nil {
			err = errors.New("parse field right_answer 正确选项 to int32 failed")
			logger.ErrorWF("parse field right_answer 正确选项 to int32 failed.",
				zap.String("xlsx", "maze_answer_mechanism_v8【答题机关词库】.xlsx"), zap.String("sheet", "maze_answer_mechanism_v8"),
				zap.String("parse_data", data[7]),
				zap.Error(err))
			return
		}
		config.Right_answer = int32(tmp)
	}
	return
}

var gMazeAnswerMechanismV8Fields = []string{
	"id",
	"question_bank_id",
	"question_stem",
	"answer_a",
	"answer_b",
	"answer_c",
	"answer_d",
	"right_answer",
}

// LoadDataManual load data for test
func LoadDataManual(logger fklog.FKLogI, load func(file, sheet string, fields []string) ([][]string, error)) (err error) {
	parser := &gMazeAnswerMechanismV8Parser{}
	loader := &gMazeAnswerMechanismV8Loader{}
	var data [][]string
	data, err = load("maze_answer_mechanism_v8【答题机关词库】.xlsx", "maze_answer_mechanism_v8", gMazeAnswerMechanismV8Fields)
	if err != nil {
		logger.ErrorWF("load maze_answer_mechanism_v8【答题机关词库】.xlsx maze_answer_mechanism_v8 data failed.", zap.Error(err))
		return
	}
	if len(data) < 1 {
		logger.WarnWF("load maze_answer_mechanism_v8【答题机关词库】.xlsx maze_answer_mechanism_v8 data empty.")
		return
	}
	container := loader.NewContainer()
	for k, row := range data {
		item := parser.New()
		if len(parser.Fields()) != len(gMazeAnswerMechanismV8Fields) {
			err = errors.New("invalid request.fileds not match")
			logger.ErrorWF("parse maze_answer_mechanism_v8【答题机关词库】.xlsx maze_answer_mechanism_v8 failed.", zap.Int("row", k), zap.Strings("need", gMazeAnswerMechanismV8Fields), zap.Strings("has", row))
			return
		}
		err = parser.Parse(logger, row, item)
		if err != nil {
			logger.ErrorWF("parse maze_answer_mechanism_v8【答题机关词库】.xlsx maze_answer_mechanism_v8 row data failed.", zap.Error(err))
			return
		}
		err = loader.Add(logger, container, item)
		if err != nil {
			logger.ErrorWF("add maze_answer_mechanism_v8【答题机关词库】.xlsx maze_answer_mechanism_v8 row data failed.", zap.Error(err))
			return
		}
	}
	err = loader.Check(container)
	if err != nil {
		logger.ErrorWF("check maze_answer_mechanism_v8【答题机关词库】.xlsx maze_answer_mechanism_v8 data failed.", zap.Error(err))
		return
	}
	loader.Swap(container)
	logger.InfoWF("load maze_answer_mechanism_v8【答题机关词库】.xlsx maze_answer_mechanism_v8 data success.")
	return
}
