package cmdbconfig

import (
	"errors"
	"fmt"

	"maze_game_server/lib/nano/configreader"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"
	"gitlab.ifreetalk.com/maze-plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
)

type CmdbConfig struct {
	cfg configreader.Config
}

type CmdbProvider struct {
	Data []byte
}

func (p *CmdbProvider) Name() string {
	return "cmdb"
}

func (p *CmdbProvider) Read(key string) ([]byte, error) {
	return p.Data, nil
}

func (p *CmdbProvider) Watch(cb configreader.ProviderCallback) {
	cb("cmdb", p.Data)
}

func NewProvider(data []byte) *CmdbProvider {
	return &CmdbProvider{
		Data: data,
	}
}

func New(name string) *CmdbConfig {
	configAPI, err := polaris.NewConfigAPIByFile(name)
	if err != nil {
		panic(err)
	}
	pp := polaris.GetConfigFileRequest{
		&model.GetConfigFileRequest{
			Namespace: "minigame-test",
		},
	}

	pp.FileName = "cmdbconf.yaml"
	pp.FileGroup = fmt.Sprintf("group%d", fkconfig.EnvVal.GroupID)
	fkfmt.Println(pp)
	configFile, err := configAPI.FetchConfigFile(&pp)
	if err != nil {
		// xlog.AppLogger().Error("configuration.loadFile GetConfigFile failed",
		// 	zap.String("namespace", namespace),
		// 	zap.String("fileGroup", fileGroup),
		// 	zap.String("fileName", fileName),
		// 	zap.Error(err),
		// )
		panic(err)
	}
	fkfmt.Println(configFile.GetContent())
	kk := NewProvider([]byte(configFile.GetContent()))
	configreader.RegisterProvider(kk)

	c, err := configreader.Load(name, configreader.WithCodec("yaml"), configreader.WithProvider(kk.Name()))
	if err != nil {
		panic(err)
	}
	return &CmdbConfig{
		cfg: c,
	}
}

func (cfg *CmdbConfig) Init() error                      { return nil }
func (cfg *CmdbConfig) Unmarshal(data interface{}) error { return cfg.cfg.Unmarshal(data) }
func (cfg *CmdbConfig) LoadConfig(key string, defaultValue interface{}) (interface{}, error) {
	if !cfg.cfg.IsSet(key) {
		return nil, errors.New("key not found")
	}
	return cfg.cfg.Get(key, defaultValue), nil
}

func (cfg *CmdbConfig) LoadParam(key string, defaultValue interface{}) (interface{}, error) {
	if !cfg.cfg.IsSet(key) {
		return nil, errors.New("key not found")
	}
	return cfg.cfg.Get(key, defaultValue), nil
}

func (cfg *CmdbConfig) IsSet(key string) bool { return cfg.cfg.IsSet(key) }

// Bytes returns original config data as bytes.
func (c *CmdbConfig) Bytes() []byte {
	return c.cfg.Bytes()
}

// GetInt returns int value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *CmdbConfig) GetInt(key string, defaultValue int) int {
	return c.cfg.GetInt(key, defaultValue)
}

// GetInt32 returns int32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *CmdbConfig) GetInt32(key string, defaultValue int32) int32 {
	return c.cfg.GetInt32(key, defaultValue)
}

// GetInt64 returns int64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *CmdbConfig) GetInt64(key string, defaultValue int64) int64 {
	return c.cfg.GetInt64(key, defaultValue)
}

// GetUint returns uint value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *CmdbConfig) GetUint(key string, defaultValue uint) uint {
	return c.cfg.GetUint(key, defaultValue)
}

// GetUint32 returns uint32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *CmdbConfig) GetUint32(key string, defaultValue uint32) uint32 {
	return c.cfg.GetUint32(key, defaultValue)
}

// GetUint64 returns uint64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *CmdbConfig) GetUint64(key string, defaultValue uint64) uint64 {
	return c.cfg.GetUint64(key, defaultValue)
}

// GetFloat64 returns float64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *CmdbConfig) GetFloat64(key string, defaultValue float64) float64 {
	return c.cfg.GetFloat64(key, defaultValue)
}

// GetFloat32 returns float32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *CmdbConfig) GetFloat32(key string, defaultValue float32) float32 {
	return c.cfg.GetFloat32(key, defaultValue)
}

// GetBool returns bool value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *CmdbConfig) GetBool(key string, defaultValue bool) bool {
	return c.cfg.GetBool(key, defaultValue)
}

// GetString returns string value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *CmdbConfig) GetString(key string, defaultValue string) string {
	return c.cfg.GetString(key, defaultValue)
}
