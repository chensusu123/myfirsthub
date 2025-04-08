package process

import (
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/servers/maze_game_server/servers/maze_gm_server/process/cmdbattledata"
)

func RegGm(logger fklog.FKLogI) {
	cmdbattledata.RegBattleDataGm(logger)
}
