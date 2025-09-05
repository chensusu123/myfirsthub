package mysql

import (
	"errors"
	"fmt"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database/nanogorm"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/discovery"
	"gorm.io/gorm"
)

type BizGorm struct {
	*nanogorm.NanoGorm
}

func NewBizGorm(serviceName string, name string) *BizGorm {
	rt := &BizGorm{}
	rt.NanoGorm = nanogorm.NewNanoGorm(serviceName, name)
	return rt
}

var gBizGorm *BizGorm

func (g *BizGorm) Init(resolver discovery.Resolver) error {
	appConfig := appconfig.GlobalConfig()
	err := g.NanoGorm.Init(resolver)
	if err != nil {
		return err
	}
	gIsLocalDev = appConfig.Global.IsLocalDev
	gBizGorm = g
	return nil
}

// type BizCfg struct {
// 	Addr       string `yaml:"addr"`
// 	DbUser     string `yaml:"db_user"`
// 	Pwd        string `yaml:"pwd"`
// 	DbName     string `yaml:"db_name"`
// 	TableName  string `yaml:"table_name"`
// 	IsLocalDev bool   `yaml:"-"`
// }
// type BizFlow struct {
// 	bizeName string
// 	resolver discovery.Resolver
// }

var gIsLocalDev bool

// func NewBizFlow(name string) *BizFlow {
// 	return &BizFlow{
// 		bizeName: name,
// 	}
// }

// func (flow *BizFlow) Name() string {
// 	return flow.bizeName
// }

// var (
// 	bizCfg *BizCfg
// 	db     *gorm.DB
// )

// // 监控实例变化必须实现的接口
// func (flow *BizFlow) WatcherName() string {
// 	return flow.bizeName
// }

// // 监控实例变化必须实现的接口
// func (flow *BizFlow) OnInstancesUpdate(change *discovery.Change) {
// 	logger := fklog.AppLogger().Clone("BizFlow")
// 	if change == nil {
// 		logger.CtxError(ctx,"OnInstancesUpdate change is nil")
// 		return
// 	}

// 	if len(change.Result.Instances) > 0 {
// 		for _, xx := range change.Result.Instances {
// 			logger.CtxInfo(ctx,"BizFlow OnInstancesUpdate ",
// 				zap.Any("Address", xx.Address().String()),
// 				zap.Any("Tags", xx.Tags()), zap.Any("Vsersion", change.Result.Vsersion))
// 		}
// 	}
// }

// var _ discovery.InstancesListener = (*BizFlow)(nil)

// func (flow *BizFlow) Init(resolver discovery.Resolver) error {
// 	appConfig := appconfig.GlobalConfig()
// 	logger := fklog.AppLogger().Clone("BizFlow")
// 	bizCfg = &BizCfg{}
// 	sectionResolver := plateregistry.NewSectionResolver(resolver, appConfig.Global.SectionID)
// 	flow.resolver = sectionResolver
// 	resolveName := appConfig.Global.Namespace + ":" + flow.bizeName

// 	mysqlInfo, err := flow.resolver.Resolve(context.TODO(), resolveName)
// 	if err != nil {
// 		logger.CtxError(ctx,"BizFlow Init Resolve failed", zap.Any("err", err))
// 		return err
// 	}
// 	if len(mysqlInfo.Instances) == 0 {
// 		logger.CtxError(ctx,"BizFlow Init Resolve failed, Instances is empty")
// 		return errors.New("BizFlow Init Resolve failed, Instances is empty")
// 	}

// 	addr := mysqlInfo.Instances[0].Address().String()

// 	dataBaseInfo := configuration.GetDatabase(flow.bizeName)
// 	if dataBaseInfo == nil {
// 		logger.CtxError(ctx,"BizFlow  GetDatabase failed")
// 		return errors.New("BizFlow GetDatabase failed")
// 	}

// 	dbCfg, ok := dataBaseInfo.Get().(*datastruct.MysqlDBCfg)
// 	if !ok {
// 		logger.CtxError(ctx,"BizFlow  GetDatabase failed")
// 		return errors.New("BizFlow GetDatabase failed")
// 	}

// 	logger.CtxInfo(ctx,"BizFlow Init  mysqlInfo GetMysqlCfg show ",
// 		zap.Any("addr", addr), zap.Any("InstancesLen",
// 			len(mysqlInfo.Instances)), zap.Any("dbCfg", dbCfg))

// 	// 监听实例变化
// 	err = flow.resolver.Watcher(context.Background(), resolveName, flow)
// 	if err != nil {
// 		logger.CtxError(ctx,"BizFlow Init Watcher failed", zap.Any("err", err))
// 	}

// 	bizCfg.Addr = addr
// 	bizCfg.DbUser = dbCfg.User
// 	bizCfg.Pwd = dbCfg.Password
// 	bizCfg.DbName = dbCfg.DB

// 	bizCfg.IsLocalDev = appConfig.Global.IsLocalDev

// 	err = initGorm(bizCfg)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (flow *BizFlow) InitWithBiz(biz *BizCfg) error {
// 	return nil
// }

// 获取gorm数据库实例
func GetMysqlDb() (*gorm.DB, error) {
	// if _, ok := db.Logger.(*GormLogger); !ok {
	// 	// gorm日志在这里初始化是因为initGorm的执行时机在log初始化之前，所以gorm的日志要晚初始化
	// 	logg := log.Clone("gorm", 0, 0)
	// 	db.Logger = NewGormLogger(logg, logger.Info)
	// }
	if gBizGorm == nil {
		return nil, errors.New("gBizGorm is nil")
	}
	return gBizGorm.GetGormDB()
}

// 获取分表名字
func getShardingTableName(baseTable string) string {
	if gIsLocalDev {
		return fmt.Sprintf("`t_%s`", baseTable)
	}
	t := time.Now()
	return fmt.Sprintf("`t_%s_%s`", baseTable, t.Format("2"))
}

// 获取分库名字
func getShardingDbName(baseDb string) string {
	if gIsLocalDev {
		return fmt.Sprintf("`%s`", "maze")
	}
	t := time.Now()
	return fmt.Sprintf("`db_%s_%s`", baseDb, t.Format("200601"))
}

// 获取完全限定表名
func GetFullyQualifiedTableName(baseName string) string {
	return fmt.Sprintf("%s.%s", getShardingDbName(baseName), getShardingTableName(baseName))
}

// func initGorm(cfg *BizCfg) error {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", cfg.DbUser, cfg.Pwd, cfg.Addr)
// 	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		fkfmt.Println("init gorm fail ", err)
// 		return err
// 	}
// 	db = gormDB
// 	return nil
// }
