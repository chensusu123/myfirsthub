/*
 * @Author: majian
 * @Date: 2025-03-20 11:04:04
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-28 16:22:04
 */
package process

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structsdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeconfigv8"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeuserlevelredis"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/common/commonmustarriveredis"
	"gitlab.ifreetalk.com/plate/protodef/Common"
	"gitlab.ifreetalk.com/plate/protodef/MazeGame"

	"go.uber.org/zap"
)

// 处理迷宫升级属性变化Id包
func HandleMazeLvUpgradeAttrChgId(logger fklog.FKLogI, userId uint64, msg *structsdef.DollAttrChgNotify) {
	if msg.ChgType != 301 {
		return
	}
	mazeLv, e := mazeuserlevelredis.GetUserLevel(logger, userId)
	if e != nil {
		logger.ErrorWF("HandleMazeLvUpgradeAttrChgId GetUserLevel fail", zap.Error(e),
			zap.Uint64("uid", userId))
		return
	}
	if mazeLv <= 1 {
		logger.InfoWF("HandleMazeLvUpgradeAttrChgId level1 ignore",
			zap.Uint64("uid", userId))
		return
	}
	careIds := mazeconfigv8.GetClientCareAttrIds()

	mazeLvChgIDMsg := &MazeGame.MazeLvUpgradeID{}
	mazeLvChgIDMsg.CurLevel = proto.Int32(int32(mazeLv))
	mazeLvChgIDMsg.Session = proto.String(msg.Session)
	for _, attr := range msg.ChgAttrs {
		// 是否需要关心的属性
		if _, ok := careIds[attr.AttrId]; !ok {
			continue
		}
		chgIdInfo := &Common.AttrChgInfo{}
		chgIdInfo.AttrId = proto.Int32(attr.AttrId)
		chgIdInfo.AttrVal = proto.Int64(attr.OldVal)
		chgIdInfo.NextVal = proto.Int64(attr.CurVal)
		mazeLvChgIDMsg.ChgAttrs = append(mazeLvChgIDMsg.ChgAttrs, chgIdInfo)
	}

	logger.InfoWF("HandleMazeLvUpgradeAttrChgId send client with",
		zap.Any("mazeLvChgIDMsg", mazeLvChgIDMsg), zap.Uint64("userId", userId))
	commonmustarriveredis.SendArrivePacketWithLogFix(logger, userId, 16239, mazeLvChgIDMsg)
}

// func IsMazeUpgradeCareAttr(attrId int32) bool {
// 	if attrId == constdef.DollFormulaAttack ||
// 		attrId == constdef.DollFormulaDefend ||
// 		attrId == constdef.DollFormulaBlood ||
// 		attrId == constdef.MazeForce {
// 		return true
// 	}
// 	return false
// }
