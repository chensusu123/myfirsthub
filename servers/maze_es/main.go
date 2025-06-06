package main

import (
	"time"

	"maze_game_server/usecase/cmdbconfig"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func initLog() {
	logConfig := fklog.LogConfig{
		LogDir:     "/tmp",
		LogLevel:   "debug",
		LogName:    "maze_es",
		LogType:    "zap",
		WithCaller: true,
	}
	fklog.InitAppFkLog(&logConfig)
}

type Config struct {
	Global struct {
		Namespace     string `yaml:"namespace"`      // Namespace for the configuration.
		EnvName       string `yaml:"env_name"`       // Environment name.
		ContainerName string `yaml:"container_name"` // Container name.
		LocalIP       string `yaml:"local_ip"`       // Local IP address.
		EnableSet     string `yaml:"enable_set"`     // Y/N. Whether to enable Set. Default is N.
		// Full set name with the format: [set name].[set region].[set group name].
		FullSetName string `yaml:"full_set_name"`
		// Size of the read buffer in bytes. <=0 means read buffer disabled. Default value will be used if not set.
		ReadBufferSize *int `yaml:"read_buffer_size,omitempty"`
	}
	Server Server `yaml:"server"`
}
type Server struct {
	App      string `yaml:"app"`       // Application name.
	Server   string `yaml:"server"`    // Server name.
	BinPath  string `yaml:"bin_path"`  // Binary file path.
	DataPath string `yaml:"data_path"` // Data file path.
	ConfPath string `yaml:"conf_path"` // Configuration file path.
}

type BizCfg struct {
	Addr      string `yaml:"addr"`
	DbUser    string `yaml:"db_user"`
	Pwd       string `yaml:"pwd"`
	DbName    string `yaml:"db_name"`
	TableName string `yaml:"table_name"`
}

func main() {
	initLog()

	gTestLogger := fklog.AppLogger().Clone("maze_energy_server_t")
	// c, err := configreader.Load("./conf.d/nano_actor.yaml", configreader.WithCodec("yaml"), configreader.WithProvider("cmdb"))

	// yy := Server{App: "app1"}
	// xx := c.Get("server", &yy)
	// var ttt int
	// xxx := c.Get("axx", ttt)
	// gTestLogger.InfoWF("start", zap.Any("err", err), zap.Any("config", c),
	// 	zap.Any("configxx", xx), zap.Any("configyy", yy), zap.Any("configxxx", xxx))

	// cfg := localconfig.New("../maze_main_server/conf.d/localconfig.yaml")
	cfg := cmdbconfig.New("../maze_main_server/conf.d/polaris.yaml")
	cfgValue, err := cfg.LoadConfig("BizCfg.addr", &BizCfg{})
	gTestLogger.InfoWF("start", zap.Any("err", err), zap.Any("config", cfgValue))
	time.Sleep(time.Second * 2)
}
