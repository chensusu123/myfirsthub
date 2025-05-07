package equip

import (
	"math/rand"
	"time"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip/equipbaggm"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip/equipaassemblegm"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip/copyequipgm"
)

func RegEquipGm(logger fklog.FKLogI) {
	logger.WarnWF("RegEquipGm Begin")
	equipbaggm.Reg(logger)
}

func RegGm(logger fklog.FKLogI) {
	RegEquipGm(logger)
	equipaassemblegm.RegGm(logger)
	copyequipgm.RegGm(logger)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
