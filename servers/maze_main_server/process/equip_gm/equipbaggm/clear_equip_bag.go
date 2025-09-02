package equipbaggm

import (
	"context"
	"maze_game_server/io/redis/mazebagequipredis"
)

func ClearEquipBag(ctx context.Context, userId uint64) error {
	// allEquip, err := mazebagequipredis.GetAllEquipInfo(ctx, logger, userId)
	// if err != nil {
	// 	return err
	// }
	err := mazebagequipredis.GMDelEquip(ctx, userId)
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

func ClearEquipBagBatch(ctx context.Context, userId uint64, equipGuids []int64) error {
	err := mazebagequipredis.BatchDelEquip(ctx, userId, equipGuids...)
	return err
}
