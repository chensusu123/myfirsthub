package cfg

type CfgSvr interface {
	Init() error
	GetConfig() interface{}
	LoadConfig(string, interface{}) error
	LoadParam(string, interface{}) error
}

type emptySvr struct{}

func (*emptySvr) Init() error                          { return nil }
func (*emptySvr) GetConfig() interface{}               { return nil }
func (*emptySvr) LoadConfig(string, interface{}) error { return nil }
func (*emptySvr) LoadParam(string, interface{}) error  { return nil }

type CfgCenter struct{}

func (c *CfgCenter) GetConfigSvr() (CfgSvr, error) {
	return &emptySvr{}, nil
}

var Cfg *CfgCenter
