package process

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_gm_server/process/cmdbattledata"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
)

func RegGm(logger fklog.FKLogI) {
	cmdbattledata.RegBattleDataGm(logger)
}
