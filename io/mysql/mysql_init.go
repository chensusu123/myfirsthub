package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"maze_game_server/lib/nano/cfg"
	"sync"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkini"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkmysql"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/naming"
)

type mysqlConfig struct {
	Address string `ini:"address"`
	DbUser  string `ini:"db_user"`
	Pwd     string `ini:"pwd"`
	DbName  string `ini:"db_name"`
}

var mysqlCfg = &mysqlConfig{}

const (
	retryInterval = 10 * time.Second
)

// todo 框架沒有提供获取配置的方法，暂时只能自己读取一遍了
func readConfig() (*fkini.IniConfig, error) {
	cfgName := "conf.d/config.ini"
	// 打开默认配置文件
	iniCfg, err := fkini.Open(cfgName)
	if err != nil {
		fkfmt.Println("open config file", "conf.d/config.ini", "failed")
		return nil, err
	}
	return iniCfg, nil
}

type BizCfg struct {
	Addr   string `yaml:"addr"`
	DbUser string `yaml:"db_user"`
	Pwd    string `yaml:"pwd"`
	DbName string `yaml:"db_name"`
}
type BizFlow struct {
}

func (flow *BizFlow) Init(cfg cfg.CfgSvr) error {
	bizCfg := &BizCfg{}
	err := cfg.LoadConfig("BizCfg", bizCfg)
	if err != nil {
		return err
	}
	return nil
}

func (flow *BizFlow) InitWithBiz(biz *BizCfg) error {
	return nil
}

// 初始化mysql
func InitMysqlBack() {
	if fkconfig.EnvVal.IsLocalDev == true {
		c, err := readConfig()
		if err != nil {
			fkfmt.Println("init mysql failed.", err)
			return
		}

		err = c.LoadConfig(mysqlCfg, "mysql", 0)
		if err != nil {
			fkfmt.Println("init failed.", err)
			return
		}
		return
	}

	var err error

	mysqlCfg.Address, err = fkconfig.GetEnv("MYSQL_ADDRESS")
	if err != nil {
		fkfmt.Println("MYSQL_ADDRESS env not set ")
	}

	mysqlCfg.DbUser, err = fkconfig.GetEnv("MYSQL_USER")
	if err != nil {
		fkfmt.Println("MYSQL_USER env not set ")
	}
	mysqlCfg.Pwd, err = fkconfig.GetEnv("MYSQL_PWD")
	if err != nil {
		fkfmt.Println("MYSQL_PWD env not set ")
	}
}

// 获取分表名字
func GetShardingTableName(baseTable string) string {
	if fkconfig.EnvVal.IsLocalDev == true {
		return fmt.Sprintf("t_%s", baseTable)
	}
	t := time.Now()
	return fmt.Sprintf("t_%s_%s", baseTable, t.Format("2"))
}

// 获取分库名字
func GetShardingDbName(baseDb string) string {
	if fkconfig.EnvVal.IsLocalDev == true {
		return mysqlCfg.DbName
	}
	t := time.Now()
	return fmt.Sprintf("db_%s_%s", baseDb, t.Format("200601"))
}

type DBEntry struct {
	db      *sql.DB
	err     error
	mu      sync.Mutex
	lastTry time.Time
}

var dbPool sync.Map // map[dbName]*DBEntry

func GetMysqlDb(baseDbName string) (*sql.DB, error) {
	dbName := GetShardingDbName(baseDbName)
	val, _ := dbPool.LoadOrStore(dbName, &DBEntry{})
	entry := val.(*DBEntry)

	entry.mu.Lock()
	defer entry.mu.Unlock()

	if entry.db != nil && entry.err == nil {
		return entry.db, nil
	}

	if time.Since(entry.lastTry) < retryInterval {
		return nil, fmt.Errorf("recent connection attempt to '%s' failed: %v", dbName, entry.err)
	}

	entry.lastTry = time.Now()

	cfg := &fkconfig.CGKConfigNode{}
	cfg.NodeList = make(map[uint32]*fkconfig.CGKNodeItem)
	node := &fkconfig.CGKNodeItem{}
	cfg.NodeList[1] = node
	node.Address = mysqlCfg.Address
	node.DbUser = mysqlCfg.DbUser
	node.DbPassword = mysqlCfg.Pwd
	node.DbName = dbName
	gMysql := &fkmysql.MysqlDB{}
	var err error
	var mysqlDb *sql.DB
	err = gMysql.Open(fklog.AppLogger(), cfg)
	if err != nil {
		entry.err = err
		fkfmt.Println("open mysql fail", err)
		return nil, err
	}
	mysqlDb, err = gMysql.GetDB(0)
	if err != nil {
		entry.err = err
		fkfmt.Println("get db fail", err)
		return nil, err
	}
	err = mysqlDb.Ping()
	if err != nil {
		entry.err = err
		fkfmt.Println("ping db fail", err)
		return nil, err
	}

	entry.db = mysqlDb
	entry.err = err

	closeOldMysqlDb(baseDbName)

	return mysqlDb, err
}

// 关闭可能的旧连接
func closeOldMysqlDb(baseDbName string) {
	t := time.Now().AddDate(0, -1, 0)
	dbName := fmt.Sprintf("db_%s_%s", baseDbName, t.Format("200601"))
	val, ok := dbPool.Load(dbName)
	if !ok {
		return // 没有这个连接池，不需要关闭
	}
	entry := val.(*DBEntry)
	// entry.mu.Lock()
	// defer entry.mu.Unlock()
	if entry.db == nil {
		return
	}
	entry.db.Close()
	entry.db = nil
	dbPool.Delete(dbName)
}

// 初始化mysql
func InitMysqlEx(namingSvr naming.NamingI) error {
	const (
		DBServiceName = "maze_main_server.mysql"
	)
	if namingSvr == nil {
		fkfmt.Println("naming is nil")
		return errors.New("naming is nil")
	}
	var err error
	mysqlCfg.Address, err = namingSvr.GetDB(DBServiceName, int32(fkconfig.EnvVal.GroupID))
	if err != nil {
		fkfmt.Println("MYSQL_ADDRESS env not set ")
		return err
	}
	dbUser, dbPwd, err := namingSvr.GetDBUserAndPassword(DBServiceName, int32(fkconfig.EnvVal.GroupID))
	// mysqlCfg.DbUser, err = fkconfig.GetEnv("MYSQL_USER")
	if err != nil {
		fkfmt.Println("MYSQL_USER env not set ")
		return err
	}
	mysqlCfg.DbUser = dbUser
	mysqlCfg.Pwd = dbPwd
	return nil
}
