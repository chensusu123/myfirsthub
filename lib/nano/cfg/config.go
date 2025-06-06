package cfg

type CfgCenter struct {
}

func (c *CfgCenter) Init() error {
	return nil
}

func (c *CfgCenter) GetConfig(name string) (string, error) {
	return "", nil
}

func (c *CfgCenter) LoadConfig(key string, info interface{}) error {
	return nil
}

func (c *CfgCenter) LoadParam(key string, valueRef interface{}) error {
	return nil
}

// 由其他模块实现这个功能
//func (c *CfgCenter) Watch(key string, call func(valueRef interface{})) error {
//	return nil
//}
