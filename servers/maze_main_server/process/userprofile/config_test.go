// @Author pangchenyang 2025/6/17 14:40:00
// @Desc: 
package userprofile

import (
	"testing"
	"fmt"
	_ "gitlab.ifreetalk.com/maze-plate/freetk/fktestutil/testlogger" // 初始化日志
	"time"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"os"
	"gopkg.in/yaml.v3"
	"maze_game_server/usecase/redisconfig"
	"maze_game_server/io/mysql"
)

var (
	logger      = fklog.AppLogger().Clone("user_profile_t")
	testProfile *Profile
)

type fileResolver struct {
	cfg Config
}

type ServiceInfo struct {
	Name string            `yaml:"name"`
	Addr string            `yaml:"addr"`
	Tags map[string]string `yaml:"tags"` // Transport type.
}

type Config struct {
	Service []ServiceInfo `yaml:"service"`
}

func (fr *fileResolver) GetMysqlConfig() (*MysqlConfig, error) {
	for _, service := range fr.cfg.Service {
		if service.Name == "Mysql" || service.Name == "BizCfg" || service.Name == "maze-main-server-dev.gm" {
			tags := service.Tags
			return &MysqlConfig{
				Addr:   service.Addr,
				User:   tags["db_user"],
				Pd:     tags["password"],
				DBName: tags["db_name"],
			}, nil
		}
	}
	return nil, fmt.Errorf("mysql config not found")
}

func (fr *fileResolver) GetRedisConfig() (*RedisConfig, error) {
	for _, service := range fr.cfg.Service {
		if service.Name == "Redis" || service.Name == "maze_main_server.redis" {
			tags := service.Tags
			return &RedisConfig{
				Addr: service.Addr,
				Pd:   tags["password"],
			}, nil
		}
	}
	return nil, fmt.Errorf("redis config not found")
}

type MysqlConfig struct {
	Addr   string
	User   string
	Pd     string
	DBName string
}

type RedisConfig struct {
	Addr string
	Pd   string
}

func New(name string) *fileResolver {
	buf, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}
	fileResolver := &fileResolver{}
	if err := yaml.Unmarshal(buf, &fileResolver.cfg); err != nil {
		panic(err)
	}
	return fileResolver
}

func TestMain(m *testing.M) {
	fmt.Println("TestMain begin")
	cfgSvr := New("../../conf.d/service.yaml")
	fmt.Println("fileCfg", cfgSvr)
	// redis
	redisCfg, err := cfgSvr.GetRedisConfig()
	if err != nil {
		panic("GetRedisConfig panic:" + err.Error())
	}
	err = redisconfig.GetRedisService().Init(redisCfg.Addr, redisCfg.Pd)
	if err != nil {
		panic("Init redis panic:" + err.Error())
	}
	err = redisconfig.GetRedisService().Start()
	if err != nil {
		panic("Start redis panic:" + err.Error())
	}
	defer redisconfig.GetRedisService().Stop()
	// mysql
	mysqlCfg, err := cfgSvr.GetMysqlConfig()
	if err != nil {
		panic("GetRedisConfig panic:" + err.Error())
	}
	bizFlow := mysql.NewBizFlow("test_mysql")
	err = bizFlow.InitWithBiz(&mysql.BizCfg{
		Addr:   mysqlCfg.Addr,
		DbUser: mysqlCfg.User,
		Pwd:    mysqlCfg.Pd,
		DbName: mysqlCfg.DBName,
	})
	if err != nil {
		panic("mysql init err:" + err.Error())
	}
	// local cache
	initCache()
	m.Run()
	fmt.Println("TestMain end")
	time.Sleep(time.Second * 2)
}

func TestGetTestConfig(t *testing.T) {
	time.Sleep(time.Second * 1)
}
