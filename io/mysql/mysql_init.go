package mysql

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/naming"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/cfg"
	"os"
	"time"
)

type BizCfg struct {
	Addr      string `yaml:"addr"`
	DbUser    string `yaml:"db_user"`
	Pwd       string `yaml:"pwd"`
	DbName    string `yaml:"db_name"`
	TableName string `yaml:"table_name"`
	Mode      string `yaml:"mode"`
}
type BizFlow struct{}

var (
	bizCfg *BizCfg
	db     *gorm.DB
)

func (flow *BizFlow) Init(cfg cfg.CfgSvr) error {
	bizCfg = &BizCfg{}
	mp, err := cfg.LoadConfig("BizCfg", bizCfg)
	fkfmt.Println("load config", "BizCfg", mp, err)
	if err != nil {
		return err
	}
	c, ok := mp.(map[string]interface{})
	if !ok {
		fkfmt.Println("load config", "BizCfg", mp, err)
	}
	bizCfg.Addr = c["addr"].(string)
	bizCfg.DbUser = c["db_user"].(string)
	bizCfg.Pwd = c["pwd"].(string)
	bizCfg.DbName = c["db_name"].(string)
	bizCfg.TableName = c["table_name"].(string)
	bizCfg.Mode = os.Getenv("mode")
	err = initGorm(bizCfg)
	if err != nil {
		return err
	}
	return nil
}

func (flow *BizFlow) InitWithBiz(biz *BizCfg) error {
	return nil
}

// 获取gorm数据库实例
func GetMysqlDb() (*gorm.DB, error) {
	if _, ok := db.Logger.(*GormLogger); !ok {
		// gorm日志在这里初始化是因为initGorm的执行时机在log初始化之前，所以gorm的日志要晚初始化
		logg := log.Clone("gorm", 0, 0)
		db.Logger = NewGormLogger(logg, logger.Info)
	}

	return db, nil
}

// 获取分表名字
func getShardingTableName(baseTable string) string {
	if bizCfg.Mode == "dev" || bizCfg.Mode == "docker" {
		return fmt.Sprintf("`t_%s`", baseTable)
	}
	t := time.Now()
	return fmt.Sprintf("`t_%s_%s`", baseTable, t.Format("2"))
}

// 获取分库名字
func getShardingDbName(baseDb string) string {
	if bizCfg.Mode == "dev" || bizCfg.Mode == "docker" {
		return fmt.Sprintf("`%s`", bizCfg.DbName)
	}
	t := time.Now()
	return fmt.Sprintf("`db_%s_%s`", baseDb, t.Format("200601"))
}

// 获取完全限定表名
func GetFullyQualifiedTableName(baseName string) string {
	return fmt.Sprintf("%s.%s", getShardingDbName(baseName), getShardingTableName(baseName))
}

// Deprecated:初始化mysql
func InitMysqlEx(namingSvr naming.NamingI) error {
	return nil
}

func initGorm(cfg *BizCfg) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", cfg.DbUser, cfg.Pwd, cfg.Addr)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fkfmt.Println("init gorm fail ", err)
		return err
	}
	db = gormDB
	return nil
}
