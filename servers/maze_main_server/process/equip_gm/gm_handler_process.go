package equip_gm

import (
	"math/rand"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/servers/maze_main_server/process/equip_gm/copyequipgm"
	"maze_game_server/servers/maze_main_server/process/equip_gm/equipaassemblegm"
	"maze_game_server/servers/maze_main_server/process/equip_gm/equipbaggm"
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
