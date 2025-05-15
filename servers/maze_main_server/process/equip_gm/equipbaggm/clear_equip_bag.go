package equipbaggm

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebagequipredis"
)

func ClearEquipBag(logger fklog.FKLogI, userId uint64) error {
	// allEquip, err := mazebagequipredis.GetAllEquipInfo(context.TODO(), logger, userId)
	// if err != nil {
	// 	return err
	// }
	err := mazebagequipredis.GMDelEquip(logger, userId)
	if err != nil {
		return err
	}
	// equipCliList := make([]*MazeGameEquip.MazeEquipInfo, 0)
	// for _, equipInfo := range allEquip {
	// 	equipCliList = append(equipCliList, packtopb.EquipInfoToCliPB(equipInfo))
	// }
	// process.SendDollBagEquipChgIDEx(logger, userId, nil, equipCliList, nil, int32(DollEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_DOLL_EQUIP_BG))

	return err
}

func ClearEquipBagBatch(logger fklog.FKLogI, userId uint64, equipGuids []int64) error {
	err := mazebagequipredis.BatchDelEquip(logger, userId, equipGuids...)
	return err
}
