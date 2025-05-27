/*
 * @Author: majian
 * @Date: 2025-03-17 10:34:49
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-26 11:43:50
 */
package equip

import (
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/assemble"
	"maze_game_server/common/function/packtopb/asequipsuittopb"
	"maze_game_server/common/function/packtopb/packequipostopb"
	"maze_game_server/common/structsdef"
	"maze_game_server/config/GMazeEquipPosRankV8Cfg"
	"maze_game_server/io/kafka/dollequipassmeblekakfa"
	"maze_game_server/io/redis/dollassemblesuitredis"
	"maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/assembleidpack"
	"maze_game_server/module/calcassembleattr"
	"maze_game_server/module/dollassembleinfo"
	"maze_game_server/module/equippossuit"
	"maze_game_server/module/mazebuffchgrrecordapi"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"
	"maze_game_server/servers/maze_main_server/process/equip/demconstdef"
	"maze_game_server/servers/maze_main_server/process/equip/module"
	"sort"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (e *Equip) OnGetMazeAssembleRQ_10414_10415(s *session.Session, req *MazeGameEquip.GetMazeGameAssembleInfoRQ) (err error) {
	defer fkprometheus.DebugPMT("OnGetDollAssembleRQ")()

	logger := log.Clone("Equip", uint64(s.UID()), 0)
	res := &MazeGameEquip.GetMazeGameAssembleInfoRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	userId := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnGetMazeAssembleRQ end", zap.Any("res", res))
	}()

	logger.InfoWF("OnGetMazeAssembleRQ with", zap.Any("req", req))
	// 检查装备位解锁
	ChkEquipPosUnlock(logger, userId, UnlockSrcInit, false)

	// 初始装备套检查
	InitDollEquipSuitSeq(logger, userId)

	// 处理初始化装备
	HandleDollEquipInit(logger, userId, false)

	// 人偶属性初始化
	HandleDollAttrInit(logger, userId, req.GetHeader().GetSession())

	assembleInfo, effect, err := dollassembleinfo.GetDollAssembleInfoEx(logger, userId)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		logger.ErrorWF("OnGetMazeAssembleRQ Get Assemble info fail", zap.Error(err))
		return err
	}
	err = checkAssembleEquipConsistent(logger, userId, assembleInfo)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		logger.ErrorWF("OnGetMazeAssembleRQ checkAssembleEquipConsistent fail", zap.Error(err))
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
				cliEquip, e := packequipostopb.PackEquipPosPb(logger, aEquip, -1)
				if e != nil {
					res.ErrInfo = errors.MODULE_ERROR.ToInfo()
					logger.ErrorWF("OnGetMazeAssembleRQ AssembleEquipToCliPb fail", zap.Error(e), zap.Int32("pos", posCfg.Pos_id))
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
		dai.EquipPosNextStSuit, e1 = equippossuit.GetCurAndNextSuit(logger, assembleInfo.GetEpEnSuitId())
	if e1 != nil {
		logger.ErrorWF("OnGetMazeAssembleRQ GetCurAndNextSuit fail", zap.Error(e1), zap.Int32("enSuitId", assembleInfo.GetEpEnSuitId()))
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

// var NeedFuncCondMap = map[int32]int32{
// 	int32(DollEquip.ENUM_FUNC_OPEN_ID_DollMount):                constdef.FuncOpenMount,
// 	int32(DollEquip.ENUM_FUNC_OPEN_ID_PosStrengthen):            constdef.FuncOpenEquipPosStrengthen,
// 	int32(DollEquip.ENUM_FUNC_OPEN_ID_Knife):                    constdef.FuncOpenKnife,
// 	int32(DollEquip.ENUM_FUNC_OPEN_ID_EquipMagicOpen):           constdef.FuncOpenFaBao,
// 	int32(DollEquip.ENUM_FUNC_OPEN_ID_EquipMagicStrengthenOpen): constdef.FuncOpenFaBaoStrenth,
// int32(DollEquip.ENUM_FUNC_OPEN_ID_EquipMagicStageOpen):      constdef.FuncOpenFaBaoStage,
// }

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
