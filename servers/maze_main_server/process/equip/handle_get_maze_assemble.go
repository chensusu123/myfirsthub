/*
 * @Author: majian
 * @Date: 2025-03-17 10:34:49
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-26 11:43:50
 */
package equip

import (
	"sort"

	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeEquipPosRankV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/protodef/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/protodef/DollEquip"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGameEquip"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/dollassembleinfo"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/packtopb/packequipostopb"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/packtopb/asequipsuittopb"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/assembleidpack"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/assemble"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip/module"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structsdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/equippossuit"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassemblesuitredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/calcassembleattr"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/dollequipassmeblekakfa"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebuffinforedis"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip/demconstdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazebuffchgrrecordapi"
)

func OnGetMazeAssembleRQ(ctx fknet.TCPContext, shardingID uint64, request proto.Message, response proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnGetDollAssembleRQ")()
	req := request.(*MazeGameEquip.GetMazeGameAssembleInfoRQ)
	res := response.(*MazeGameEquip.GetMazeGameAssembleInfoRS)

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	userCtx := fkserver.NewUserContext(ctx.Context, shardingID, ctx.FKLogI)

	defer func() {
		userCtx.InfoWF("OnGetMazeAssembleRQ end", zap.Any("res", res))
	}()

	userCtx.InfoWF("OnGetMazeAssembleRQ with", zap.Any("req", req))
	// 检查装备位解锁
	ChkEquipPosUnlock(userCtx, shardingID, UnlockSrcInit, false)

	// 初始装备套检查
	InitDollEquipSuitSeq(userCtx, shardingID)

	// 处理初始化装备
	HandleDollEquipInit(userCtx, shardingID, false)

	// 人偶属性初始化
	HandleDollAttrInit(userCtx, shardingID, req.GetHeader().GetSession())

	assembleInfo, effect, err := dollassembleinfo.GetDollAssembleInfoEx(userCtx, shardingID)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		userCtx.ErrorWF("OnGetMazeAssembleRQ Get Assemble info fail", zap.Error(err))
		return err
	}
	err = checkAssembleEquipConsistent(userCtx, shardingID, assembleInfo)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		userCtx.ErrorWF("OnGetMazeAssembleRQ checkAssembleEquipConsistent fail", zap.Error(err))
		return err
	}
	_ = effect

	// mazeforcepreview.AssembleMazeForcePreview(userCtx, shardingID, assembleInfo, "maze_equip_main_server")

	dai := &MazeGameEquip.MazeAssembleInfo{}
	//	var posUnlockNum int32
	allRows := GMazeEquipPosRankV8Cfg.GetAllMazeEquipPosRankV8Config()
	for _, posCfg := range allRows {
		// 打包装配的装备信息
		var unlock int32
		for _, aEquip := range assembleInfo.MazeEquips {
			if posCfg.Pos_id == aEquip.GetEquipPos().GetPos() {
				cliEquip, e := packequipostopb.PackEquipPosPb(userCtx, aEquip, -1)
				if e != nil {
					res.ErrInfo = errors.MODULE_ERROR.ToInfo()
					userCtx.ErrorWF("OnGetMazeAssembleRQ AssembleEquipToCliPb fail", zap.Error(e), zap.Int32("pos", posCfg.Pos_id))
					return e
				}
				dai.EquipPosList = append(dai.EquipPosList, cliEquip)
				unlock = 1
				break
			}
		}
		if unlock == 1 {
			//	posUnlockNum++
			continue
		}
		// 未解锁打包
		unlockPos := &MazeGameEquip.EquipPosInfo{}
		unlockPos.EquipPos = proto.Int32(posCfg.Pos_id)
		unlockPos.IsUnlock = proto.Int32(0)
		unlockPos.UnlockDesc = proto.String(posCfg.Unlock_desc)
		dai.EquipPosList = append(dai.EquipPosList, unlockPos)
	}
	sort.Slice(dai.EquipPosList, func(i, j int) bool {
		return dai.EquipPosList[i].GetEquipPos() < dai.EquipPosList[j].GetEquipPos()
	})

	dai.AsEquipSuitInfo = asequipsuittopb.PackAsEquipSuitInfo(assembleInfo.GetEpSuitId())
	var e1 error
	dai.EquipPosStSuit,
		dai.EquipPosNextStSuit, e1 = equippossuit.GetCurAndNextSuit(userCtx, assembleInfo.GetEpEnSuitId())
	if e1 != nil {
		userCtx.ErrorWF("OnGetMazeAssembleRQ GetCurAndNextSuit fail", zap.Error(e1), zap.Int32("enSuitId", assembleInfo.GetEpEnSuitId()))
	}
	res.MazeAssembleInfo = dai

	res.Token = proto.Int64(assembleidpack.GetAssembleToken())
	return nil
}

// func getFuncOpenMask(logger fklog.FKLogI, funcId, dolllv, maskId int32) int32 {
// 	fc, e := funcopencheck.IsFuncOpen(funcId, dolllv)
// 	if e != nil {
// 		logger.ErrorWF("getFuncOpenMask IsFuncOpen fail", zap.Error(e), zap.Int32("funcId", funcId))

// 	} else {
// 		if fc.IsOpen {
// 			return maskId
// 		} else {
// 			logger.InfoWF("getFuncOpenMask IsFuncOpen", zap.Int32("funcId", funcId), zap.Any("fc", fc), zap.Int32("dollLv", dolllv))
// 		}
// 	}
// 	return 0
// }

var NeedFuncCondMap = map[int32]int32{
	int32(DollEquip.ENUM_FUNC_OPEN_ID_DollMount):                constdef.FuncOpenMount,
	int32(DollEquip.ENUM_FUNC_OPEN_ID_PosStrengthen):            constdef.FuncOpenEquipPosStrengthen,
	int32(DollEquip.ENUM_FUNC_OPEN_ID_Knife):                    constdef.FuncOpenKnife,
	int32(DollEquip.ENUM_FUNC_OPEN_ID_EquipMagicOpen):           constdef.FuncOpenFaBao,
	int32(DollEquip.ENUM_FUNC_OPEN_ID_EquipMagicStrengthenOpen): constdef.FuncOpenFaBaoStrenth,
	// int32(DollEquip.ENUM_FUNC_OPEN_ID_EquipMagicStageOpen):      constdef.FuncOpenFaBaoStage,
}

// func getFuncOpenCond() (conds []*DollEquip.DollFuncOpenCond) {
// 	for k, v := range NeedFuncCondMap {
// 		row := GActionCountCfg.GetActionCountConfig(v)
// 		if row != nil {
// 			cond := &DollEquip.DollFuncOpenCond{}
// 			cond.FuncId = proto.Int32(k)
// 			cond.DollLevel = proto.Int32(row.Count_level_v8[7])
// 			cond.CondDesc = proto.String(row.Lose_desc_v8)
// 			conds = append(conds, cond)
// 		}
// 	}
// 	return conds
// }

// 检查装备数据
func checkAssembleEquipConsistent(logger fklog.FKLogI, userId uint64, assembleInfo *MazeEquipCache.MazeAssembleDb) error {
	var needUpdateEquipPos []*MazeEquipCache.MazeEquipPosInfo
	for _, equipPos := range assembleInfo.GetMazeEquips() {
		if !assemble.IsAssembleEquip(equipPos) {
			continue
		}
		if equipPos.GetEquipLoadInfo().GetEquipGuid() != equipPos.GetEquipInfo().GetEquipGuid() {
			continue
		}
		var chg bool
		if equipPos.GetEquipLoadInfo().GetEquipId() != equipPos.GetEquipInfo().GetEquipId() {
			equipPos.GetEquipLoadInfo().EquipId = equipPos.GetEquipInfo().EquipId
			chg = true
		}
		// if equipPos.GetEquipLoadInfo().GetEquipSubType() != equipPos.GetEquipInfo().GetEquipSubType() {
		// 	equipPos.GetEquipLoadInfo().EquipSubType = equipPos.GetEquipInfo().EquipSubType
		// 	chg = true
		// }
		if chg {
			needUpdateEquipPos = append(needUpdateEquipPos, equipPos)
		}
	}
	if len(needUpdateEquipPos) > 0 {
		err := dollassemblesuitredis.SaveEquipAssembleInfoV2(logger, userId, assembleInfo.GetCurSuitIndex(), needUpdateEquipPos)
		if err != nil {
			logger.ErrorWF("checkAssembleEquipConsistent SaveEquipAssembleInfoV2 fail", zap.Error(err), zap.Any("needUpdateEquipPos", needUpdateEquipPos))
		} else {
			logger.WarnWF("checkAssembleEquipConsistent SaveEquipAssembleInfoV2 succ", zap.Any("needUpdateEquipPos", needUpdateEquipPos))
		}
		return err
	}
	return nil
}

// 检查装配的装备是否存在，已测ok
func checkAssembleEquipLose(logger fklog.FKLogI, userId uint64, assembleInfo *MazeEquipCache.MazeAssembleDb, effectOld *calcassembleattr.EquipmentEffectInfo) error {
	var needUpdateEquipPos []*MazeEquipCache.MazeEquipPosInfo
	var recordList []*dollequipassmeblekakfa.MazeGameEquipAssembleRecord
	for _, equipPos := range assembleInfo.GetMazeEquips() {
		if !assemble.IsAssembleEquip(equipPos) {
			continue
		}
		if equipPos.GetEquipLoadInfo().GetEquipGuid() > 0 && equipPos.GetEquipInfo().GetEquipGuid() == 0 {

			record := StartEquipAssmebleRecord(userId, equipPos.EquipPos.GetPos(), dollequipassmeblekakfa.DollEquipAssembleOpDown, nil, equipPos,
				effectOld)
			recordList = append(recordList, record)
			module.ResetEquipPos(equipPos)
			needUpdateEquipPos = append(needUpdateEquipPos, equipPos)
		}
	}
	if len(needUpdateEquipPos) > 0 {
		err := dollassemblesuitredis.SaveEquipAssembleInfoV2(logger, userId, assembleInfo.GetCurSuitIndex(), needUpdateEquipPos)
		if err != nil {
			logger.ErrorWF("checkAssembleEquipLose SaveEquipAssembleInfoV2 fail", zap.Error(err), zap.Any("needUpdateEquipPos", needUpdateEquipPos))
			return err
		}
		logger.WarnWF("checkAssembleEquipLose SaveEquipAssembleInfoV2 succ", zap.Any("needUpdateEquipPos", needUpdateEquipPos))

		effectInfo, e := calcassembleattr.CalcEquipEffect(logger, assembleInfo.MazeEquips)
		if e == nil {
			assembleInfo.EpSuitId = proto.Int32(effectInfo.SuitId)
			// 更新人偶buff
			e = mazebuffinforedis.SaveMazeEquipBuff(logger, userId, effectInfo.Other)
			if e == nil {
				// 通知计算属性
				calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{}
				calcAttrNotify.FromServer = demconstdef.MySvr
				calcAttrNotify.UserId = userId
				calcAttrNotify.ChgType = constdef.MazeBuffEquipFix
				calcAttrNotify.BuffSrc = constdef.MazeBuffSrcEquip
				mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, calcAttrNotify)

				mazebuffchgrrecordapi.SendMazeBuffChgRecord(logger, userId,
					constdef.MazeBuffSrcEquip,
					constdef.MazeBuffEquipFix,
					effectOld.Other, effectInfo.Other)
			}
		}
		// 记录流水
		for _, record := range recordList {
			EndEquipAssmebleRecord(logger, record, 0, 0, effectInfo)
		}
	}
	return nil
}
