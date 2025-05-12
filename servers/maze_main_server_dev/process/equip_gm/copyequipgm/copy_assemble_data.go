/*
 * @Author: majian
 * @Date: 2024-12-16 14:39:31
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 19:19:20
 */
package copyequipgm

import (
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/MazeEquipCache"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/dollassembleinfo"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassemblesuitredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/assemble"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassembleredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/calcassembleattr"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebuffinforedis"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structsdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazebuffchgrrecordapi"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip_gm/copyequipgm/copyinterface"
)

func CopyAssembleData(logger fklog.FKLogI, srcUserId uint64, dstUsers []uint64, param copyinterface.CopyParam) error {

	var err error
	var assembleInfo *MazeEquipCache.MazeAssembleDb
	var curUser uint64
	defer func() {
		if err != nil {
			logger.ErrorWF("CopyAssembleData fail", zap.Error(err),
				zap.Uint64("src", srcUserId),
				zap.Uint64("dst", curUser))
		} else {
			logger.InfoWF("CopyAssembleData succ",
				zap.Uint64("src", srcUserId),
				zap.Int("dstLen", len(dstUsers)))
		}
	}()
	assembleInfo, ef, err := dollassembleinfo.GetDollAssembleInfoEx(logger, srcUserId)
	if err != nil {
		return err
	}

	for _, dstId := range dstUsers {
		curUser = dstId
		// 删除旧的装配数据
		err = dollassemblesuitredis.DelEquipSuitInfo(logger, dstId)
		if err != nil {
			return err
		}
		// 写装备位数据
		var srcFields []string
		for i := 1; i <= int(constdef.EquipPosNum); i++ {
			field := assemble.EnCodeAssemblePosField(int32(i))
			srcFields = append(srcFields, field)
		}

		// 删除旧装备位数据
		err = dollassembleredis.BatchDelAssmebleInfo(logger, dstId, srcFields...)
		if err != nil {
			return err
		}
		var dstFields []string
		for _, pos := range assembleInfo.MazeEquips {
			dstFields = append(dstFields, assemble.EnCodeAssemblePosField(int32(pos.GetEquipPos().GetPos())))
		}

		// 写装配数据
		err = dollassemblesuitredis.SaveEquipAssembleInfo(logger, dstId, assembleInfo.GetCurSuitIndex(), assembleInfo.MazeEquips)
		if err != nil {
			return err
		}

		dstAssmebleInfo := &MazeEquipCache.MazeAssembleDb{}
		dstAssmebleInfo.CurSuitIndex = proto.Int32(assembleInfo.GetCurSuitIndex())
		dstAssmebleInfo.SwitchSuitTime = proto.Int64(assembleInfo.GetSwitchSuitTime())
		dstAssmebleInfo.MazeEquips = assembleInfo.MazeEquips
		dstFields = append(dstFields, constdef.AssemblePrefixCurAssembleSuitIndex, constdef.AssemblePrefixSwitchSuitTime)
		err = dollassembleredis.SetAssembleInfoByFields(logger, dstId, dstFields, dstAssmebleInfo)
		if err != nil {
			return err
		}
		CalcDollAttr(logger, dstId, assembleInfo, ef)
		logger.InfoWF("CopyAssembleData user succ",
			zap.Uint64("src", srcUserId),
			zap.Int("dst", int(dstId)))
	}
	return nil
}

func CalcDollAttr(logger fklog.FKLogI, userId uint64, assembleInfo *MazeEquipCache.MazeAssembleDb, effectInfo *calcassembleattr.EquipmentEffectInfo) error {
	_, otherAttrs := effectInfo.ForceAttrs, effectInfo.Other
	// 更新buff中心
	e := mazebuffinforedis.SaveMazeEquipBuff(logger, userId, otherAttrs)
	if e != nil {
		return e
	}

	calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{}
	calcAttrNotify.FromServer = "maze_equip_gm_server"
	calcAttrNotify.UserId = userId
	calcAttrNotify.ChgType = constdef.MazeBuffChgTypeGm
	calcAttrNotify.Session = ""
	calcAttrNotify.BuffSrc = constdef.MazeBuffSrcEquip

	mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, calcAttrNotify)

	mazebuffchgrrecordapi.SendMazeBuffChgRecord(logger, userId,
		constdef.MazeBuffSrcEquip,
		constdef.MazeBuffChgTypeGm,
		nil, effectInfo.Other)

	return nil
}
