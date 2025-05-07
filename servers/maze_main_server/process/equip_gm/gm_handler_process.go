package equip_gm

import (
	"math/rand"
	"time"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip_gm/copyequipgm"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip_gm/equipaassemblegm"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip_gm/equipbaggm"
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
