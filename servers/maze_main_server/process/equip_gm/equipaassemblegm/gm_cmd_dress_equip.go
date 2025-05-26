/*
 * @Author: majian
 * @Date: 2025-03-20 14:03:38
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-20 20:56:24
 */
package equipaassemblegm

import (
	"context"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"go.uber.org/zap"
	"maze_game_server/common/constdef"
	"maze_game_server/common/function/assemble"
	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/io/redis/mazebagequipredis"
	"maze_game_server/module/dollassembleinfo"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"
	"maze_game_server/servers/maze_main_server/process/equip"
)

func DressEquipGm(logger fklog.FKLogI, userId uint64, pos int32, equipGuid int64) error {
	assembleInfo, _, err := dollassembleinfo.GetDollAssembleInfoEx(logger, userId)
	if err != nil {
		logger.ErrorWF("DressEquipGm Get Assemble info fail", zap.Error(err))
		return err
	}
	allEquips, err := mazebagequipredis.GetAllEquipInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("DressEquipGm GetAllEquipInfo info fail", zap.Error(err))
		return err
	}
	for i := 1; i <= constdef.EquipPosNum; i++ {
		if pos > 0 && pos != int32(i) {
			continue
		}
		var dressInfo *MazeEquipCache.MazeEquipPosInfo
		for _, dressEquip := range assembleInfo.GetMazeEquips() {
			if dressEquip.GetEquipPos().GetPos() == int32(i) {
				dressInfo = dressEquip
				break
			}
		}
		if dressInfo != nil {
			if !assemble.IsAssembleEquip(dressInfo) {
				if equipGuid == 0 {
					for _, bagEquip := range allEquips {
						cfg := GMazeEquipInfoV8Cfg.GetMazeEquipInfoV8Config(bagEquip.GetEquipId())
						if cfg != nil && cfg.Pos == int32(i) {
							equipGuid = bagEquip.GetEquipGuid()
							break
						}
					}
				}
				if equipGuid == 0 {
					logger.WarnWF("DressEquipGm cannot find equip", zap.Int32("pos", int32(i)))
				} else {
					equipDetail := allEquips[equipGuid]
					if equipDetail == nil {
						logger.WarnWF("DressEquipGm bag equip no exist",
							zap.Int32("pos", int32(i)),
							zap.Int64("equipGuid", equipGuid))
					} else {
						gmDressOneEquip(logger, userId, int32(i), equipDetail.GetEquipId(), equipGuid)
					}
				}
			}
		}
	}
	return nil
}

func gmDressOneEquip(logger fklog.FKLogI, userId uint64, pos, equipId int32, guid int64) error {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: logger}
	rq := &MazeGameEquip.MazeDressEquipRQ{}
	rs := &MazeGameEquip.MazeDressEquipRS{}
	rq.EquipPos = proto.Int32(pos)
	rq.EquipGuid = proto.Int64(guid)
	rq.OpSrc = proto.Int32(0)
	rq.OpType = proto.Int32(1)
	return equip.OnDressMazeEquipRQ(ctx, userId, rq, rs)
}
