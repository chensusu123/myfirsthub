package cfg

type CfgSvr interface {
	Init() error
	// Unmarshal deserializes the config into input param.
	Unmarshal(interface{}) error
	// Get returns config by key.
	LoadConfig(string, interface{}) (interface{}, error)

	// IsSet returns if the config specified by key exists.
	IsSet(string) bool

	// GetInt returns int value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetInt(string, int) int

	// GetInt32 returns int32 value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetInt32(string, int32) int32

	// GetInt64 returns int64 value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetInt64(string, int64) int64

	// GetUint returns uint value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetUint(string, uint) uint

	// GetUint32 returns uint32 value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetUint32(string, uint32) uint32

	// GetUint64 returns uint64 value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetUint64(string, uint64) uint64

	// GetFloat32 returns float32 value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetFloat32(string, float32) float32

	// GetFloat64 returns float64 value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetFloat64(string, float64) float64

	// GetString returns string value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetString(string, string) string

	// GetBool returns bool value by key, the second parameter
	// is default value when key is absent or type conversion fails.
	GetBool(string, bool) bool

	// Bytes returns config data as bytes.
	Bytes() []byte
}

type emptySvr struct{}

func (*emptySvr) Init() error                                         { return nil }
func (*emptySvr) Unmarshal(interface{}) error                         { return nil }
func (*emptySvr) LoadConfig(string, interface{}) (interface{}, error) { return nil, nil }

func (*emptySvr) IsSet(string) bool { return false }

// Bytes returns original config data as bytes.
func (c *emptySvr) Bytes() []byte {
	return nil
}

// GetInt returns int value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *emptySvr) GetInt(key string, defaultValue int) int {
	return 0
}

// GetInt32 returns int32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *emptySvr) GetInt32(key string, defaultValue int32) int32 {
	return 0
}

// GetInt64 returns int64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *emptySvr) GetInt64(key string, defaultValue int64) int64 {
	return 0
}

// GetUint returns uint value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *emptySvr) GetUint(key string, defaultValue uint) uint {
	return 0
}

// GetUint32 returns uint32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *emptySvr) GetUint32(key string, defaultValue uint32) uint32 {
	return 0
}

// GetUint64 returns uint64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *emptySvr) GetUint64(key string, defaultValue uint64) uint64 {
	return 0
}

// GetFloat64 returns float64 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *emptySvr) GetFloat64(key string, defaultValue float64) float64 {
	return 0
}

// GetFloat32 returns float32 value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *emptySvr) GetFloat32(key string, defaultValue float32) float32 {
	return 0
}

// GetBool returns bool value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *emptySvr) GetBool(key string, defaultValue bool) bool {
	return false
}

// GetString returns string value by key, the second parameter
// is default value when key is absent or type conversion fails.
func (c *emptySvr) GetString(key string, defaultValue string) string {
	return ""
}

type CfgCenter struct{}

func (c *CfgCenter) GetConfigSvr() (CfgSvr, error) {
	return &emptySvr{}, nil
}

var Cfg *CfgCenter
