package mysql

import (
	"context"
	"fmt"
	"time"

	"maze_game_server/lib/log"

	"gitlab.ifreetalk.com/maze-plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/discovery"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/instanceutil"
	"gitlab.ifreetalk.com/maze-plate/freetk/plateregistry"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type BizCfg struct {
	Addr       string `yaml:"addr"`
	DbUser     string `yaml:"db_user"`
	Pwd        string `yaml:"pwd"`
	DbName     string `yaml:"db_name"`
	TableName  string `yaml:"table_name"`
	IsLocalDev bool   `yaml:"-"`
}
type BizFlow struct {
	bizeName string
}

func NewBizFlow(name string) *BizFlow {
	return &BizFlow{
		bizeName: name,
	}
}

func (flow *BizFlow) Name() string {
	return flow.bizeName
}

var (
	bizCfg *BizCfg
	db     *gorm.DB
)

func (flow *BizFlow) Init(resolver discovery.Resolver) error {
	logger := fklog.AppLogger().Clone("BizFlow")
	bizCfg = &BizCfg{}
	groupResolver := plateregistry.NewGroupResolver(resolver)
	mysqlInfo, err := groupResolver.Resolve(context.TODO(), fkconfig.EnvVal.Namespace+":"+flow.bizeName)
	if err != nil {
		logger.ErrorWF("BizFlow Init Resolve failed", zap.Any("err", err))
		return err
	}
	mysqlCfg, mysqlCfgErr := instanceutil.GetMysqlCfg(mysqlInfo.Instances)
	logger.InfoWF("BizFlow Init  mysqlInfo GetMysqlCfg show ", zap.Any("mysqlCfg", mysqlCfg),
		zap.Any("mysqlCfgErr", mysqlCfgErr), zap.Any("InstancesLen", len(mysqlInfo.Instances)))
	if mysqlCfgErr != nil {
		logger.ErrorWF("BizFlow Init GetMysqlCfg failed", zap.Any("mysqlCfgErr", mysqlCfgErr))
		return mysqlCfgErr
	}

	bizCfg.Addr = mysqlCfg.Addr
	bizCfg.DbUser = mysqlCfg.DbUser
	bizCfg.Pwd = mysqlCfg.Password
	bizCfg.DbName = mysqlCfg.DbName

	bizCfg.IsLocalDev = fkconfig.EnvVal.IsLocalDev
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
	if bizCfg.IsLocalDev {
		return fmt.Sprintf("`t_%s`", baseTable)
	}
	t := time.Now()
	return fmt.Sprintf("`t_%s_%s`", baseTable, t.Format("2"))
}

// 获取分库名字
func getShardingDbName(baseDb string) string {
	if bizCfg.IsLocalDev {
		return fmt.Sprintf("`%s`", bizCfg.DbName)
	}
	t := time.Now()
	return fmt.Sprintf("`db_%s_%s`", baseDb, t.Format("200601"))
}

// 获取完全限定表名
func GetFullyQualifiedTableName(baseName string) string {
	return fmt.Sprintf("%s.%s", getShardingDbName(baseName), getShardingTableName(baseName))
}

func initGorm(cfg *BizCfg) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DbUser, cfg.Pwd, cfg.Addr, cfg.DbName)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fkfmt.Println("init gorm fail ", err)
		return err
	}
	db = gormDB
	return nil
}
