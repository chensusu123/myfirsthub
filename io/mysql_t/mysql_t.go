package mysql_t

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkmysql"
)

func init() {
	fkconfig.RegisterNameNode("designateNameMysql", 2221, designateName)
	fkconfig.RegisterNameNode("defaultNameMysql", 2221, defaultName)
}

var designateName = &fkmysql.MysqlDB{
	ServiceName: "test.mysql",
}

var defaultName = &fkmysql.MysqlDB{}
