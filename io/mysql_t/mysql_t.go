package mysql_t

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkmysql"
)

func init() {
	fkconfig.RegisterNameNode("AccountToPaipaiUnionIDMysql", 2221, a)
}

var a = &fkmysql.MysqlDB{}
