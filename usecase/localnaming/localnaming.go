package localnaming

import (
	"strings"

	"gitlab.ifreetalk.com/maze-plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkini"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	namingI "gitlab.ifreetalk.com/maze-plate/freetk/pkg/naming"
	"go.uber.org/zap"
)

// LocalNaming It is used to assemble multiple associated client's Options
type LocalNaming struct {
	cfgName     string
	cfg         *fkini.IniConfig
	InitFunList []InitFunc
	fklog.FKLogI
}

type InitFunc func(namingI.NamingI) error

func NewLocalNaming(cfgName string, initFuncList []InitFunc) *LocalNaming {
	return &LocalNaming{
		cfgName:     cfgName,
		InitFunList: initFuncList,
	}
}

type mysqlConfig struct {
	Address string `ini:"address"`
	DbUser  string `ini:"db_user"`
	Pwd     string `ini:"pwd"`
	DbName  string `ini:"db_name"`
}
type redisConfig struct {
	Address string `ini:"redis"`
}

// LoadConfig read config from ini
func (c *LocalNaming) LoadConfig(config interface{}, name string, type_id int) error {
	section, err := c.cfg.GetSection(name)
	if err != nil {
		return err
	}
	// c.InfoWF("load section config with", zap.Any("section", section))
	err = section.MapTo(config)
	if err != nil {
		fkfmt.Println("reflect config failed. ", name, err)
		return err
	}
	return err
}

func (n *LocalNaming) GetServer(svrName string, groupID int32) ([]string, error) {
	return []string{"127.0.0.1:6379"}, nil
}

func (n *LocalNaming) GetDB(dbName string, groupID int32) (string, error) {
	if strings.Contains(dbName, "mysql") {
		mysqlCfg := &mysqlConfig{}
		err := n.cfg.LoadConfig(mysqlCfg, "mysql", 0)
		n.InfoWF("GetDB", zap.Any("dbName", dbName), zap.Any("groupID", groupID), zap.Any("mysqlCfg", mysqlCfg))
		if err != nil {
			fkfmt.Println("init failed.", err)
			return "", err
		}
		return mysqlCfg.Address, nil
	}
	if strings.Contains(dbName, "redis") {
		redisCfg := &redisConfig{}
		err := n.cfg.LoadConfig(redisCfg, "CgkCfg", 0)
		n.InfoWF("GetDB", zap.Any("dbName", dbName), zap.Any("groupID", groupID), zap.Any("redisCfg", redisCfg))
		if err != nil {
			fkfmt.Println("init failed.", err)
			return "", err
		}
		return redisCfg.Address, nil
	}
	n.InfoWF("GetDB", zap.Any("dbName", dbName), zap.Any("groupID", groupID))
	return "", nil
}

func (n *LocalNaming) GetDBUserAndPassword(dbName string, groupID int32) (string, string, error) {
	mysqlCfg := &mysqlConfig{}
	err := n.cfg.LoadConfig(mysqlCfg, "mysql", 0)
	n.InfoWF("GetDBUserAndPassword", zap.Any("dbName", dbName), zap.Any("groupID", groupID), zap.Any("mysqlCfg", mysqlCfg))
	if err != nil {
		fkfmt.Println("init failed.", err)
		return "", "", err
	}
	return mysqlCfg.DbUser, mysqlCfg.Pwd, nil
}

func (n *LocalNaming) Init() error {
	fkfmt.Println("LocalNaming Init config file", n.cfgName)
	cfg, err := readConfig(n.cfgName)
	if err != nil {
		return err
	}
	n.cfg = cfg

	logger := fklog.AppLogger().Clone("")
	n.FKLogI = logger
	for _, f := range n.InitFunList {
		err = f(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func readConfig(cfgName string) (*fkini.IniConfig, error) {
	// 打开默认配置文件
	iniCfg, err := fkini.Open(cfgName)
	if err != nil {
		fkfmt.Println("open config file", "conf.d/config.ini", "failed")
		return nil, err
	}
	return iniCfg, nil
}
