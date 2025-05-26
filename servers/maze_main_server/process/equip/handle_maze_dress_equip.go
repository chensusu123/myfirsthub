/*
 * @Author: majian
 * @Date: 2025-03-17 14:41:29
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-28 21:00:53
 */
package equip

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"go.uber.org/zap"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/maputil"
	"maze_game_server/common/function/packtopb/packequipostopb"
	"maze_game_server/common/function/pbutil"
	"maze_game_server/common/structsdef"
	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/excel/toastmsgtipexcel"
	"maze_game_server/io/kafka/dollequipassmeblekakfa"
	"maze_game_server/io/redis/dollassemblesuitredis"
	"maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/io/redis/mazeuserlevelredis"
	"maze_game_server/module/assembleidpack"
	"maze_game_server/module/calcassembleattr"
	"maze_game_server/module/dollassembleinfo"
	"maze_game_server/module/effectequip"
	"maze_game_server/module/mazebuffchgrrecordapi"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/servers/maze_main_server/process/equip/demconstdef"
	"maze_game_server/servers/maze_main_server/process/equip/module"
)

func OnDressMazeEquipRQ(ctx fknet.TCPContext, shardingID uint64, request proto.Message, response proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnDressMazeEquipRQ")()
	req := request.(*MazeGameEquip.MazeDressEquipRQ)
	res := response.(*MazeGameEquip.MazeDressEquipRS)

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.EquipPos = req.EquipPos
	res.OpSrc = req.OpSrc
	userCtx := fkserver.NewUserContext(ctx.Context, shardingID, ctx.FKLogI)

	defer func() {
		userCtx.InfoWF("OnDressMazeEquipRQ end", zap.Any("res", res))
	}()

	userCtx.InfoWF("OnDressMazeEquipRQ with", zap.Any("req", req))

	limitKey := MakeLimiterKey(shardingID, 16181)
	if GtcpLimiter.IsRateLimit(limitKey) {
		userCtx.WarnWF("OnDressMazeEquipRQ rate limiter", zap.Any("req", req))
		res.ErrInfo = errors.NewCodeError(constdef.DE_ERR_RQ_RATE_LIMITER,
			toastmsgtipexcel.GetToastMsgTip(2018, "操作太频繁")).ToInfo()
		return
	}
	GtcpLimiter.UpdateTime(userCtx, limitKey)

	var curSuitSeq int32 = 1
	pos := req.GetEquipPos()
	dressEquipGuid := req.GetEquipGuid()       // 要穿戴的装备
	replacedGuid := req.GetReplacedEquipGuid() // 要替换掉的装备
	opType := req.GetOpType()                  // 操作类型
	var (
		dressEquipDb    *MazeEquipCache.MazeEquipInfoDb
		chgEquipPosList []*MazeEquipCache.MazeEquipPosInfo
		updateEquipPos  []*MazeEquipCache.MazeEquipPosInfo
		noMatch         bool // 客户端数据是否与服务器数据不一致
	)

	if pos <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("请指定装备位")
		userCtx.WarnWF("OnDressMazeEquipRQ equip pos invalid", zap.Int32("pos", pos))
		return nil
	}

	if opType != 1 && opType != 2 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效的操作类型")
		userCtx.WarnWF("OnDressMazeEquipRQ op type invalid", zap.Int32("opType", opType))
		return nil
	}

	mazeLv, err := mazeuserlevelredis.GetUserLevel(userCtx, shardingID)
	if err != nil {
		userCtx.ErrorWF("OnDressMazeEquipRQ GetUserLevel fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	// 检查装备功能开启
	// fc, e := funcopencheck.IsFuncOpen(constdef.FuncOpenEquip, int32(mazeLv))
	// if e != nil {
	// 	userCtx.ErrorWF("OnDressMazeEquipRQ IsFuncOpen fail", zap.Error(e))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }
	// if !fc.IsOpen {
	// 	userCtx.WarnWF("OnDressMazeEquipRQ func no open", zap.Int32("needLv", fc.NeedLv), zap.Int64("mazeLv", mazeLv))
	// 	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(fc.LoseDesc)
	// 	return
	// }
	// if opType == 2 {
	// 	return OnDownEquipRQ(userCtx, shardingID, req, res)
	// }

	if dressEquipGuid <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("请选择要穿戴的装备")
		userCtx.WarnWF("OnDressMazeEquipRQ dress equip invalid", zap.Int64("dressGuid", dressEquipGuid))
		return nil
	}

	if dressEquipGuid == replacedGuid {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数无效")
		userCtx.WarnWF("OnDressMazeEquipRQ equip not match", zap.Int32("pos", pos),
			zap.Int64("dressGuid", dressEquipGuid), zap.Int64("replacedGuid", replacedGuid))
		return nil
	}

	dressEquipDb, err = effectequip.GetEffectEquipInfo(userCtx, shardingID, dressEquipGuid)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		userCtx.ErrorWF("OnDressMazeEquipRQ get bag equip info fail", zap.Error(err), zap.Int64("dressEquipGuid", dressEquipGuid))
		return err
	}
	if dressEquipDb == nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(toastmsgtipexcel.GetToastMsgTip(2019, "要穿戴装备不存在"))
		userCtx.WarnWF("OnDressMazeEquipRQ equip not exist", zap.Int64("dressEquipGuid", dressEquipGuid))
		return err
	}
	// 临时背包装备
	// if req.GetOpSrc()&1 == 1 {
	// 	return OnDressTmpBagEquipRQ(userCtx, shardingID, int32(dollLv), req, res, dressEquipDb)
	// }
	// if dressEquipDb.GetEnterTime() > 0 {
	// 	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("请求类型与数据不匹配")
	// 	ctx.WarnWF("OnDressMazeEquipRQ equip at tmp bag", zap.Int64("dressEquipGuid", dressEquipGuid),
	// 		zap.Int64("enterTime", dressEquipDb.GetEnterTime()))
	// 	return nil
	// }
	assembleInfo, oldEffect, err := dollassembleinfo.GetDollAssembleInfoEx(userCtx, shardingID)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		userCtx.ErrorWF("OnDressMazeEquipRQ Get Assemble info fail", zap.Error(err))
		return err
	}
	if assembleInfo.GetCurSuitIndex() == 0 {
		res.ErrInfo = errors.NewErrorInfo(constdef.DE_ERR_EQUIP_SUIT_NOT_INIT, "未设置生效装备套")
		userCtx.WarnWF("OnDressMazeEquipRQ no set cur suit", zap.Int32("cCurSeq", curSuitSeq))
		return err
	}
	if assembleInfo.GetCurSuitIndex() != curSuitSeq {
		res.ErrInfo = errors.NewErrorInfo(constdef.DE_ERR_EQUIP_SUIT_NOT_MATCH, "套装不匹配")
		userCtx.WarnWF("OnDressMazeEquipRQ cur suit not match",
			zap.Int32("sCurSeq", assembleInfo.GetCurSuitIndex()), zap.Int32("cCurSeq", curSuitSeq))
		return err
	}
	defer func() {
		// 如果客户端请求穿戴信息和服务器不一致时，带回正确的数据给客户端
		if noMatch && len(res.MazeEquipPos) == 0 {
			for _, chgEquip := range assembleInfo.GetMazeEquips() {
				cliPb, _ := packequipostopb.PackEquipPosPb(userCtx, chgEquip, int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_LOAD_EQUIP_INFO))
				if cliPb != nil {
					res.MazeEquipPos = append(res.MazeEquipPos, cliPb)
				}
			}
		}
	}()
	dressedEquip := module.GetEquipPosInfo(assembleInfo, pos)
	if dressedEquip == nil {
		userCtx.WarnWF("OnDressMazeEquipRQ pos not unlock", zap.Int32("pos", pos))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("装备位未解锁")
		return
	}
	dressedGuid := module.GetDressedGuid(dressedEquip)
	if dressedGuid > 0 {
		if dressedGuid != replacedGuid {
			noMatch = true
			userCtx.WarnWF("OnDressMazeEquipRQ dressed equip not match",
				zap.Int32("pos", pos),
				zap.Int64("dressedGuid", dressedGuid),
				zap.Int64("replacedGuid", replacedGuid))
			if req.GetOpSrc()&2 == 0 { // 快捷穿戴不检查
				res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("已穿戴装备数据不匹配")
				return nil
			}
		}
	} else {
		if replacedGuid > 0 {
			noMatch = true
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("当前未穿戴装备")
			userCtx.WarnWF("OnDressMazeEquipRQ no dressed equip",
				zap.Int32("pos", pos),
				zap.Int64("dressedGuid", dressedGuid),
				zap.Int64("replacedGuid", replacedGuid))
			return nil
		}
	}

	euqipInfoCfgRow := GMazeEquipInfoV8Cfg.GetMazeEquipInfoV8Config(dressEquipDb.GetEquipId())
	if euqipInfoCfgRow == nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("未找到装备配置")
		userCtx.ErrorWF("OnDressMazeEquipRQ 未找到装备配置 ", zap.Int32("equipId", dressEquipDb.GetEquipId()), zap.Int64("equipGuid", dressEquipGuid))
		return err
	}
	if euqipInfoCfgRow.Pos != pos {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("装备位不匹配")
		userCtx.WarnWF("OnDressMazeEquipRQ equip pos not match",
			zap.Int32("pos", pos),
			zap.Int32("previewPos", euqipInfoCfgRow.Pos),
			zap.Int32("previewEquipId", dressEquipDb.GetEquipId()),
			zap.Int64("previewEquipGuid", dressEquipGuid))
		return nil
	}
	if mazeLv < int64(euqipInfoCfgRow.Level) {
		userCtx.WarnWF("OnDressMazeEquipRQ maze level no enough",
			zap.Int64("equipGuid", dressEquipGuid),
			zap.Int64("mazeLv", mazeLv),
			zap.Int32("equipId", dressEquipDb.GetEquipId()),
			zap.Int32("equipLevel", euqipInfoCfgRow.Level))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(toastmsgtipexcel.GetToastMsgTip(2021, "等级不足，无法穿戴"))
		return
	}
	// 记录旧值
	var oldEquipPos *MazeEquipCache.MazeEquipPosInfo
	if dressedEquip.GetEquipLoadInfo() != nil {
		oldEquipPos = &MazeEquipCache.MazeEquipPosInfo{}
		oldLoadInfo := &MazeEquipCache.MazeEquipPosDb{
			EquipGuid: proto.Int64(dressedEquip.GetEquipLoadInfo().GetEquipGuid()),
			EquipId:   proto.Int32(dressedEquip.GetEquipLoadInfo().GetEquipId())}
		oldEquipPos.EquipLoadInfo = oldLoadInfo
		oldEquipPos.EquipInfo = dressedEquip.GetEquipInfo()
	}

	module.ReplaceEquip(userCtx, dressedEquip, dressEquipDb)
	chgEquipPosList = append(chgEquipPosList, dressedEquip)
	updateEquipPos = append(updateEquipPos, dressedEquip)
	var effectInfo *calcassembleattr.EquipmentEffectInfo
	effectInfo, err = calcassembleattr.CalcEquipEffect(userCtx, assembleInfo.MazeEquips)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		userCtx.ErrorWF("OnDressMazeEquipRQ CalcEquipEffect fail ",
			zap.Error(err),
			zap.Int32("equipId", dressEquipDb.GetEquipId()),
			zap.Int64("equipGuid", dressEquipGuid))
		return err
	}

	for _, equipPos := range assembleInfo.MazeEquips {
		// 去重，不用重复添加
		if equipPos.GetEquipPos().GetPos() == pos {
			continue
		}
		if effectInfo.ChkPosChg(pos) {
			chgEquipPosList = append(chgEquipPosList, equipPos)
		}
	}

	// 记录流水
	recordType := dollequipassmeblekakfa.DollEquipAssembleOpDress
	if dressedGuid > 0 {
		recordType = dollequipassmeblekakfa.DollEquipAssembleOpReplace
	}

	record := StartEquipAssmebleRecord(shardingID, pos, recordType, dressedEquip,
		oldEquipPos, oldEffect)

	var opCode int32 = 1
	var opMask int32

	defer func() {
		EndEquipAssmebleRecord(userCtx, record, opCode, opMask, effectInfo)
	}()
	// 保存装配数据
	err = dollassemblesuitredis.SaveEquipAssembleInfoV2(userCtx, shardingID, assembleInfo.GetCurSuitIndex(), updateEquipPos)
	if err != nil {
		res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
		opMask |= demconstdef.DollEquipAssembleOpMaskDbSave
		userCtx.ErrorWF("OnDressMazeEquipRQ SaveEquipAssembleInfoV2 fail", zap.Error(err), zap.Any("updateEquipPos", updateEquipPos))
		return err
	}
	opCode = 0

	// 通知背包服务
	e := NotifyBagSvr(userCtx, shardingID, pos, 0, dressEquipGuid, dressedGuid)
	if e != nil {
		opMask |= demconstdef.DollEquipAssembleOpMaskNotifyBag
		userCtx.ErrorWF("OnDressMazeEquipRQ NotifyBagSvr fail", zap.Error(e),
			zap.Int64("upGuid", dressEquipGuid), zap.Int64("downGuid", dressedGuid))
	}

	var chgMask int32

	// 更新buff中心
	e = mazebuffinforedis.SaveMazeEquipBuff(userCtx, shardingID, effectInfo.Other)
	if e != nil {
		opMask |= demconstdef.DollEquipAssembleOpMaskNonForceBuff
	} else {
		// 通知计算属性
		calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{}
		calcAttrNotify.FromServer = demconstdef.MySvr
		calcAttrNotify.UserId = shardingID
		calcAttrNotify.ChgType = constdef.MazeBuffChgTypeEquipDress
		calcAttrNotify.BuffSrc = constdef.MazeBuffSrcEquip

		calcAttrNotify.Session = req.GetHeader().GetSession()
		e = mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(userCtx, calcAttrNotify)
		if e != nil {
			opMask |= demconstdef.DollEquipAssembleOpMaskCalcAttr
			userCtx.ErrorWF("OnDressMazeEquipRQ SendDollAttrCalcNotify fail", zap.Error(e))
		}

		mazebuffchgrrecordapi.SendMazeBuffChgRecord(userCtx, shardingID,
			constdef.MazeBuffSrcEquip,
			constdef.MazeBuffChgTypeEquipDress,
			oldEffect.Other, effectInfo.Other)
	}

	// 打包rs
	chgAssembleInfo := &MazeEquipCache.MazeAssembleDb{}
	if assembleInfo.GetEpSuitId() != effectInfo.SuitId {
		chgAssembleInfo.EpSuitId = proto.Int32(effectInfo.SuitId)
		chgMask |= int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_SUIT_MASK)
	}

	chgAssembleInfo.MazeEquips = chgEquipPosList
	chgMask |= int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_POS_MASK)
	assembleidpack.SendAssembleChgID(userCtx, shardingID, chgAssembleInfo,
		chgMask, int32(int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_LOAD_EQUIP_INFO)),
		constdef.DollAssembleChgTypeReplaceEquip)

	return nil
}

func NotifyBagSvr(logger fklog.FKLogI, userId uint64, pos, src int32, upGuid, downGuid int64) error {
	rq := &MazeEquipSvr.SvrMazeEquipAssembleRQ{}
	rq.UserId = proto.Uint64(userId)
	rq.EquipPos = proto.Int32(pos)
	rq.EquipGuid = proto.Int64(upGuid)
	rq.ReplacedEquipGuid = proto.Int64(downGuid)
	rq.EquipSrc = proto.Int32(src)
	rs := &MazeEquipSvr.SvrMazeEquipAssembleRS{}
	// e := dollequipbagrpc.MazeEquipAssembleRQ(logger, rq, rs)
	e := OnSvrMazeEquipAssembleRQ(logger, int64(userId), rq, rs)
	if e != nil {
		return e
	}
	if rs.GetErrInfo().GetErrCode() == errors.NO_ERROR_CODE {
		return nil
	}
	return errors.New("请求背包失败")
}

func StartEquipAssmebleRecord(userId uint64, pos, op int32, newEquip, oldEquip *MazeEquipCache.MazeEquipPosInfo, oldEffect *calcassembleattr.EquipmentEffectInfo) *dollequipassmeblekakfa.MazeGameEquipAssembleRecord {
	record := new(dollequipassmeblekakfa.MazeGameEquipAssembleRecord)
	record.UserId = userId
	record.EquipPos = pos
	var oldExtra, newExtra string
	if newEquip != nil {
		record.NewEquipId = newEquip.GetEquipLoadInfo().GetEquipId()
		record.NewGuid = uint64(newEquip.GetEquipLoadInfo().GetEquipGuid())
		newExtra = MakeExtra(pbutil.GetDollEquipName(newEquip.GetEquipInfo(), 0))
	}
	if oldEquip != nil {
		record.OldEquipId = oldEquip.GetEquipLoadInfo().GetEquipId()
		record.OldGuid = uint64(oldEquip.GetEquipLoadInfo().GetEquipGuid())
		oldExtra = MakeExtra(pbutil.GetDollEquipName(oldEquip.GetEquipInfo(), 0))
	}
	record.OpType = op
	hsStr := maputil.MapToString32(oldEffect.GetSuitCalc().GetSuitNumMap())
	//	fiveStr := maputil.MapToString32(oldEffect.FiveStateMap)
	record.OldFElem = fmt.Sprintf("hurtSuit:%s|%s", hsStr, oldExtra)
	record.NewFElem = newExtra
	return record
}

func EndEquipAssmebleRecord(logger fklog.FKLogI, record *dollequipassmeblekakfa.MazeGameEquipAssembleRecord,
	opRet, oMask int32, newEffect *calcassembleattr.EquipmentEffectInfo) {
	record.RetCode = opRet
	record.CodeMask = oMask

	// fiveElemStr := maputil.MapToString32(newEffect.FiveStateMap)
	hsStr := maputil.MapToString32(newEffect.GetSuitCalc().GetSuitNumMap())
	record.NewFElem = fmt.Sprintf("hurtSuit:%s|%s", hsStr, record.NewFElem)
	dollequipassmeblekakfa.SendMazeGameEquipAssembleRecord(logger, record)
}

func MakeExtra(resId int32, equipName string) string {
	if resId > 0 {
		return fmt.Sprintf("equipName:%s_resId:%d", equipName, resId)
	} else {
		return fmt.Sprintf("equipName:%s", equipName)
	}
}
