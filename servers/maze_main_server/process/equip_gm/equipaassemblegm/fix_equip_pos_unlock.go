/*
 * @Author: majian
 * @Date: 2024-08-21 11:20:33
 * @Last Modified by: majian
 * @Last Modified time: 2024-08-21 12:09:28
 */
package equipaassemblegm

import (
	"context"
	"maze_game_server/io/redis/dollassembleredis"
	"maze_game_server/io/redis/dollassemblesuitredis"
	"maze_game_server/pb/server/MazeEquipCache"

	"google.golang.org/protobuf/proto"
)

func UnlockPosByEquip(ctx context.Context, userId uint64) (cnt int32, err error) {
	// logger := fklog.ContextAppLogger(ctx)
	suitInfo, e := dollassemblesuitredis.GetDollAssembleSuit(ctx, userId, 1, 8)
	if e != nil {
		return 0, e
	}
	posMap, e := dollassembleredis.GetDollEquipPosInfo(ctx, userId, 8)
	if e != nil {
		return 0, e
	}
	var unlockList []*MazeEquipCache.MazeEquipSlotDb
	for _, equip := range suitInfo {
		if equip.GetEquipGuid() <= 0 {
			continue
		}
		if equip.GetPos() <= 0 {
			continue
		}
		posInfo := posMap[equip.GetPos()]
		if posInfo != nil && posInfo.GetPos() > 0 {
			continue
		}
		unlockList = append(unlockList, &MazeEquipCache.MazeEquipSlotDb{
			Pos:   proto.Int32(equip.GetPos()),
			Level: proto.Int32(0)})
	}
	e = dollassembleredis.SetDollEquipPosInfo(ctx, userId, unlockList)
	return int32(len(unlockList)), e
}
