package localconfig

import (
	"errors"

	"maze_game_server/lib/nano/configreader"
)

type LocalConfig struct {
	cfg configreader.Config
}

func New(name string) *LocalConfig {
	c, err := configreader.Load(name, configreader.WithCodec("yaml"), configreader.WithProvider("file"))
	if err != nil {
		panic(err)
	}
	return &LocalConfig{
		cfg: c,
	}
}

func (cfg *LocalConfig) Init() error                      { return nil }
func (cfg *LocalConfig) Unmarshal(data interface{}) error { return cfg.cfg.Unmarshal(data) }
func (cfg *LocalConfig) LoadConfig(key string, defaultValue interface{}) (interface{}, error) {
	if !cfg.cfg.IsSet(key) {
		return nil, errors.New("key not found")
	}
	return cfg.cfg.Get(key, defaultValue), nil
}

func (cfg *LocalConfig) LoadParam(key string, defaultValue interface{}) (interface{}, error) {
	if !cfg.cfg.IsSet(key) {
		return nil, errors.New("key not found")
	}
	return cfg.cfg.Get(key, defaultValue), nil
}

func (cfg *LocalConfig) IsSet(key string) bool { return cfg.cfg.IsSet(key) }

// Bytes returns original config data as bytes.
func (c *LocalConfig) Bytes() []byte {
	return c.cfg.Bytes()
}

// GetInt returns int value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *LocalConfig) GetInt(key string, defaultValue int) int {
	return c.cfg.GetInt(key, defaultValue)
}

// GetInt32 returns int32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *LocalConfig) GetInt32(key string, defaultValue int32) int32 {
	return c.cfg.GetInt32(key, defaultValue)
}

// GetInt64 returns int64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *LocalConfig) GetInt64(key string, defaultValue int64) int64 {
	return c.cfg.GetInt64(key, defaultValue)
}

// GetUint returns uint value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *LocalConfig) GetUint(key string, defaultValue uint) uint {
	return c.cfg.GetUint(key, defaultValue)
}

// GetUint32 returns uint32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *LocalConfig) GetUint32(key string, defaultValue uint32) uint32 {
	return c.cfg.GetUint32(key, defaultValue)
}

// GetUint64 returns uint64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *LocalConfig) GetUint64(key string, defaultValue uint64) uint64 {
	return c.cfg.GetUint64(key, defaultValue)
}

// GetFloat64 returns float64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *LocalConfig) GetFloat64(key string, defaultValue float64) float64 {
	return c.cfg.GetFloat64(key, defaultValue)
}

// GetFloat32 returns float32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *LocalConfig) GetFloat32(key string, defaultValue float32) float32 {
	return c.cfg.GetFloat32(key, defaultValue)
}

// GetBool returns bool value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *LocalConfig) GetBool(key string, defaultValue bool) bool {
	return c.cfg.GetBool(key, defaultValue)
}

// GetString returns string value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *LocalConfig) GetString(key string, defaultValue string) string {
	return c.cfg.GetString(key, defaultValue)
}
