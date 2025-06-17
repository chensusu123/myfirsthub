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
	resolver discovery.Resolver
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

// 监控实例变化必须实现的接口
func (flow *BizFlow) WatcherName() string {
	return flow.bizeName
}

// 监控实例变化必须实现的接口
func (flow *BizFlow) OnInstancesUpdate(change *discovery.Change) {
	logger := fklog.AppLogger().Clone("BizFlow")
	if change == nil {
		logger.ErrorWF("OnInstancesUpdate change is nil")
		return
	}

	if len(change.Result.Instances) > 0 {
		for _, xx := range change.Result.Instances {
			logger.InfoWF("BizFlow OnInstancesUpdate ",
				zap.Any("Address", xx.Address().String()),
				zap.Any("Tags", xx.Tags()), zap.Any("Vsersion", change.Result.Vsersion))
		}
	}
}

var _ discovery.InstancesListener = (*BizFlow)(nil)

func (flow *BizFlow) Init(resolver discovery.Resolver) error {
	logger := fklog.AppLogger().Clone("BizFlow")
	bizCfg = &BizCfg{}
	groupResolver := plateregistry.NewGroupResolver(resolver)
	flow.resolver = groupResolver
	resolveName := fkconfig.EnvVal.Namespace + ":" + flow.bizeName

	mysqlInfo, err := flow.resolver.Resolve(context.TODO(), resolveName)
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

	// 监听实例变化
	err = flow.resolver.Watcher(context.Background(), resolveName, flow)
	if err != nil {
		logger.ErrorWF("BizFlow Init Watcher failed", zap.Any("err", err))
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", cfg.DbUser, cfg.Pwd, cfg.Addr)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fkfmt.Println("init gorm fail ", err)
		return err
	}
	db = gormDB
	return nil
}
