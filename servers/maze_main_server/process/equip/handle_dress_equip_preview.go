package equip

import (
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipInfoV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGameEquip"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeuserlevelredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip/module"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeequipconfigv8"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/effectequip"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/packtopb"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/dollassembleinfo"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/packtopb/equipsuittopb"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
)

func OnDressEquipPreviewRQ(ctx fknet.TCPContext, shardingID uint64, request proto.Message, response proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnDressEquipPreviewRQ")()
	req := request.(*MazeGameEquip.MazeDressEquipPreviewRQ)
	res := response.(*MazeGameEquip.MazeDressEquipPreviewRS)

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.EquipPos = req.EquipPos
	res.OpSrc = req.OpSrc
	userCtx := fkserver.NewUserContext(ctx.Context, shardingID, ctx.FKLogI)

	defer func() {
		userCtx.InfoWF("OnDressEquipPreviewRQ end", zap.Any("res", res))
	}()

	userCtx.InfoWF("OnDressEquipPreviewRQ with", zap.Any("req", req))

	pos := req.GetEquipPos()
	replaceGuid := req.GetPreviewGuid()
	selfGuid := req.GetSelfGuid()
	var curSuitSeq int32 = 1

	if curSuitSeq <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("请指定套装")
		userCtx.WarnWF("OnDressEquipPreviewRQ no suit seq", zap.Int32("suitSeq", curSuitSeq))
		return nil
	}
	maxSuitSeq := mazeequipconfigv8.GetMaxEquipSuitNum()
	if curSuitSeq > maxSuitSeq {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("套装序号无效")
		userCtx.WarnWF("OnDressEquipPreviewRQ suit seq invalid", zap.Int32("suitSeq", curSuitSeq),
			zap.Int32("maxSuitSeq", maxSuitSeq))
		return nil
	}

	if pos <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("请选择装备位")
		userCtx.WarnWF("OnDressEquipPreviewRQ euqip pos invalid", zap.Int32("pos", pos))
		return nil
	}

	if selfGuid <= 0 && replaceGuid <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("请选择装备")
		userCtx.WarnWF("OnDressEquipPreviewRQ no equip", zap.Int32("pos", pos),
			zap.Int64("selfGuid", selfGuid), zap.Int64("replaceGuid", replaceGuid))
		return nil
	}

	mazeLv, err := mazeuserlevelredis.GetUserLevel(userCtx, shardingID)
	if err != nil {
		userCtx.ErrorWF("OnDressEquipPreviewRQ GetDollLevel fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	var replaceEquip *MazeEquipCache.MazeEquipInfoDb
	var replaceEquipCli, selfEquipCli *MazeGameEquip.MazeEquipInfo

	if replaceGuid > 0 {
		replaceEquip, err = effectequip.GetEffectEquipInfo(userCtx, shardingID, replaceGuid)
		//	replaceEquip, isIns, err = FindEquip(userCtx, shardingID, replaceGuid, req.GetOpSrc())
		if err != nil {
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
			userCtx.ErrorWF("OnDressEquipPreviewRQ get bag equip info fail", zap.Error(err),
				zap.Int64("bagEquip", replaceGuid))
			return err
		}

		row := GMazeEquipInfoV8Cfg.GetMazeEquipInfoV8Config(replaceEquip.GetEquipId())
		if row == nil {
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("未找到装备配置")
			userCtx.ErrorWF("OnDressEquipPreviewRQ 未找到装备配置 ", zap.Int32("equipId", replaceEquip.GetEquipId()), zap.Int64("equipGuid", replaceGuid))
			return err
		}
		if row.Pos != pos {
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("装备位不匹配")
			userCtx.WarnWF("OnDressEquipPreviewRQ equip pos not match",
				zap.Int32("pos", pos),
				zap.Int32("previewPos", row.Pos),
				zap.Int32("previewEquipId", replaceEquip.GetEquipId()),
				zap.Int64("previewEquipGuid", replaceGuid))
			return nil
		}

		var e1 error
		replaceEquipCli, e1 = packtopb.EquipInfoToCliPB(userCtx, replaceEquip)
		if e1 != nil {
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("预览失败")
			userCtx.ErrorWF("OnDressEquipPreviewRQ EquipInfoToCliPB fail",
				zap.Int32("equipId", replaceEquip.GetEquipId()),
				zap.Int64("equipGuid", replaceGuid))
			return e1
		}
	}
	assembleInfo, err := dollassembleinfo.GetDollAssembleInfo(userCtx, shardingID)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		userCtx.ErrorWF("OnDressEquipPreviewRQ Get Assemble info fail", zap.Error(err))
		return err
	}
	if assembleInfo.GetCurSuitIndex() == 0 {
		res.ErrInfo = errors.NewErrorInfo(constdef.DE_ERR_EQUIP_SUIT_NOT_INIT, "未设置生效装备套")
		userCtx.WarnWF("OnDressEquipPreviewRQ no set cur suit", zap.Int32("cCurSeq", curSuitSeq))
		return err
	}
	if assembleInfo.GetCurSuitIndex() != curSuitSeq {
		res.ErrInfo = errors.NewErrorInfo(constdef.DE_ERR_EQUIP_SUIT_NOT_MATCH, "套装不匹配")
		userCtx.WarnWF("OnDressEquipPreviewRQ cur suit not match",
			zap.Int32("sCurSeq", assembleInfo.GetCurSuitIndex()), zap.Int32("cCurSeq", curSuitSeq))
		return err
	}

	dressedEquipPos := module.GetEquipPosInfo(assembleInfo, pos)
	dressedGuid := module.GetDressedGuid(dressedEquipPos)
	if selfGuid != dressedGuid {
		// 如果是快捷预览武力值，不校验
		userCtx.WarnWF("OnDressEquipPreviewRQ dressed equip not match", zap.Int32("pos", pos),
			zap.Int64("dressedGuid", dressedGuid),
			zap.Int64("cliDressedGuid", selfGuid))
		res.ErrInfo = errors.NewCodeError(constdef.DE_ERR_DRESS_EQUIP_DATA_NO_MATCH, "").ToInfo()
		return nil
	}
	if dressedGuid > 0 {
		equipSuitPb, _ := equipsuittopb.PackEquipSuitCliPb(userCtx, pos, assembleInfo.GetMazeEquips(), dressedEquipPos.GetEquipInfo().GetSuitId(), int32(mazeLv))
		selfEquipCli, _ = packtopb.EquipInfoToCliPB(userCtx, dressedEquipPos.GetEquipInfo())
		// if assemble.IsFiveElemActivate(dressedEquipPos.GetEquipLoadInfo().GetActivateMask()) {
		// 	fe := selfEquipCli.GetFiveElemInfo()
		// 	if fe != nil {
		// 		fe.IsActive = proto.Bool(true)
		// 	}
		// }
		if equipSuitPb != nil {
			selfEquipCli.SuitInfo = equipSuitPb
		}
	}

	if replaceEquipCli != nil && replaceEquip != nil {
		equipSuitPb, _ := equipsuittopb.PackEquipSuitCliPb(userCtx, pos, assembleInfo.GetMazeEquips(), replaceEquip.GetSuitId(), int32(mazeLv))
		if equipSuitPb != nil {
			replaceEquipCli.SuitInfo = equipSuitPb
		}
	}
	// selfForce, replaceForce, e := mazeforcepreview.ReplaceEquipForcePreview(userCtx, shardingID, assembleInfo, replaceEquip, pos, "maze_equip_main_server")
	// if e == nil {
	// 	if replaceEquipCli != nil {
	// 		replaceEquipCli.ForceValue = proto.Int64(replaceForce)
	// 	}
	// 	if selfEquipCli != nil {
	// 		selfEquipCli.ForceValue = proto.Int64(selfForce)
	// 	}
	// }

	// 计算迷宫武力
	if replaceEquipCli != nil {
		res.DollEquipInfo = append(res.DollEquipInfo, replaceEquipCli)
	}
	if selfEquipCli != nil {
		res.DollEquipInfo = append(res.DollEquipInfo, selfEquipCli)
	}

	return nil
}

// func clearEquipNew(logger fklog.FKLogI, userID uint64, bagEquips, loadEquips []int64) {
// 	rq := &DollEquipSvr.SvrEquipCleanNewFlagRQ{}
// 	rq.UserId = proto.Uint64(userID)
// 	if len(bagEquips) > 0 {
// 		rq.EquipGuids = append(rq.EquipGuids, bagEquips...)
// 	}
// 	if len(loadEquips) > 0 {
// 		rq.WearEquipGuids = append(rq.WearEquipGuids, loadEquips...)
// 	}

// 	rs := &DollEquipSvr.SvrEquipCleanNewFlagRS{}
// 	dollequipbagrpc.ClearEquipNewFlagRQ(logger, rq, rs)
// }

// func FindEquip(logger fklog.FKLogI, userId uint64, equipGuid int64, src int32) (effectEquip *MazeEquipCache.MazeEquipInfoDb, ins bool, err error) {
// 	defer func() {
// 		if err == effectequip.EquipNoExist {
// 			err = errors.New(toastmsgtipexcel.GetToastMsgTip(2027, "预览装备不存在"))
// 		}
// 	}()
// 	if src == 4 { // 预览实例化装备
// 		effectEquip, err = effectequip.GetInstanceEffectEquipInfo(logger, userId, equipGuid)
// 		if err == nil {
// 			ins = true
// 		}
// 		if err == effectequip.EquipNoExist {
// 			effectEquip, err = effectequip.GetEffectEquipInfo(logger, userId, equipGuid)
// 		}
// 	} else {
// 		effectEquip, err = effectequip.GetEffectEquipInfo(logger, userId, equipGuid)
// 		if err == effectequip.EquipNoExist {
// 			effectEquip, err = effectequip.GetInstanceEffectEquipInfo(logger, userId, equipGuid)
// 			if err == nil {
// 				ins = true
// 			}
// 		}
// 	}
// 	return effectEquip, ins, err
// }
