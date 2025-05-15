/*
 * @Author: majian
 * @Date: 2025-03-22 16:23:09
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-24 14:36:00
 */
package commonlogic

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/Common"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazePropertyPanel"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structsdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeconfigv8"
	"gitlab.ifreetalk.com/maze-plate/io/redisio/commonmustarriveredis"
)

func NotifyClientAttrChg(logger fklog.FKLogI, userId uint64, msg *structsdef.DollAttrChgNotify) {
	if msg.ChgType == constdef.MazeBuffLvChg {
		logger.InfoWF("NotifyClientAttrChg maze level chg ignore", zap.Any("msg", msg))
		return
	}
	mazePanelChgIDMsg := &MazePropertyPanel.MazePanelAttrChgID{}
	mazePanelChgIDMsg.Session = proto.String(msg.Session)
	mazePanelChgIDMsg.ChgType = proto.Int32(msg.ChgType)
	careIds := mazeconfigv8.GetClientCareAttrIds()
	for _, attr := range msg.ChgAttrs {
		// 是否需要关心的属性
		if _, ok := careIds[attr.AttrId]; !ok {
			continue
		}
		chgIdInfo := &Common.AttrChgInfo{}
		chgIdInfo.AttrId = proto.Int32(attr.AttrId)
		chgIdInfo.AttrVal = proto.Int64(attr.OldVal)
		chgIdInfo.NextVal = proto.Int64(attr.CurVal)
		mazePanelChgIDMsg.ChgAttrs = append(mazePanelChgIDMsg.ChgAttrs, chgIdInfo)
	}
	if len(mazePanelChgIDMsg.ChgAttrs) > 0 {
		logger.InfoWF("NotifyClientAttrChg send client with", zap.Any("mazePanelChgIDMsg", mazePanelChgIDMsg))
		commonmustarriveredis.SendArrivePacketWithLog(logger, userId, 16262, mazePanelChgIDMsg)
	} else {
		logger.InfoWF("NotifyClientAttrChg no care attrs", zap.Any("mazePanelChgIDMsg", mazePanelChgIDMsg))
	}
}
