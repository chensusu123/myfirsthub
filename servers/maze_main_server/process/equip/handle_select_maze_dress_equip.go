/*
 * @Author: majian
 * @Date: 2025-03-17 14:41:29
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 18:01:42
 */
package equip

import (
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/structsdef"
	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/excel/toastmsgtipexcel"
	"maze_game_server/io/kafka/dollequipassmeblekakfa"
	"maze_game_server/io/redis/dollassemblesuitredis"
	"maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/io/redis/mazeuserlevelredis"
	"maze_game_server/lib/nano/session"
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

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (ep *Equip) OnSelectDressMazeEquipRQ_10420_10421(s *session.Session, req *MazeGameEquip.SelectDressMazeEquipRQ) (err error) {
	defer fkprometheus.DebugPMT("OnSelectDressMazeEquipRQ")()

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeGameEquip.SelectDressMazeEquipRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.EquipPos = req.EquipPos

	userId := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnSelectDressMazeEquipRQ end", zap.Any("res", res))
	}()

	logger.CtxInfo(ctx, "OnSelectDressMazeEquipRQ with", zap.Any("req", req))

	var curSuitSeq int32 = 1
	pos := req.GetEquipPos()
	dressEquipGuid := req.GetSelectGuid() // 要穿戴的装备
	replacedGuid := req.GetLoadedGuid()   // 要替换掉的装备

	var (
		dressEquipDb    *MazeEquipCache.MazeEquipInfoDb
		chgEquipPosList []*MazeEquipCache.MazeEquipPosInfo
		updateEquipPos  []*MazeEquipCache.MazeEquipPosInfo
	)

	if pos <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("请指定装备位")
		logger.CtxWarn(ctx, "OnDressMazeEquipRQ equip pos invalid", zap.Int32("pos", pos))
		return nil
	}

	mazeLv, err := mazeuserlevelredis.GetUserLevel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnSelectDressMazeEquipRQ GetUserLevel fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	if dressEquipGuid <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("请选择要穿戴的装备")
		logger.CtxWarn(ctx, "OnSelectDressMazeEquipRQ dress equip invalid", zap.Int64("dressGuid", dressEquipGuid))
		return nil
	}

	if dressEquipGuid == replacedGuid {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数无效")
		logger.CtxWarn(ctx, "OnSelectDressMazeEquipRQ equip not match", zap.Int32("pos", pos),
			zap.Int64("dressGuid", dressEquipGuid), zap.Int64("replacedGuid", replacedGuid))
		return nil
	}

	dressEquipDb, err = effectequip.GetEffectEquipInfo(logger, userId, dressEquipGuid)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		logger.CtxError(ctx, "OnSelectDressMazeEquipRQ get bag equip info fail", zap.Error(err), zap.Int64("dressEquipGuid", dressEquipGuid))
		return err
	}
	if dressEquipDb == nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(toastmsgtipexcel.GetToastMsgTip(2019, "要穿戴装备不存在"))
		logger.CtxWarn(ctx, "OnSelectDressMazeEquipRQ equip not exist", zap.Int64("dressEquipGuid", dressEquipGuid))
		return err
	}

	assembleInfo, oldEffect, err := dollassembleinfo.GetDollAssembleInfoEx(logger, userId)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		logger.CtxError(ctx, "OnSelectDressMazeEquipRQ Get Assemble info fail", zap.Error(err))
		return err
	}
	if assembleInfo.GetCurSuitIndex() == 0 {
		res.ErrInfo = errors.NewErrorInfo(constdef.DE_ERR_EQUIP_SUIT_NOT_INIT, "未设置生效装备套")
		logger.CtxWarn(ctx, "OnSelectDressMazeEquipRQ no set cur suit", zap.Int32("cCurSeq", curSuitSeq))
		return err
	}
	if assembleInfo.GetCurSuitIndex() != curSuitSeq {
		res.ErrInfo = errors.NewErrorInfo(constdef.DE_ERR_EQUIP_SUIT_NOT_MATCH, "套装不匹配")
		logger.CtxWarn(ctx, "OnSelectDressMazeEquipRQ cur suit not match",
			zap.Int32("sCurSeq", assembleInfo.GetCurSuitIndex()), zap.Int32("cCurSeq", curSuitSeq))
		return err
	}

	dressedEquip := module.GetEquipPosInfo(assembleInfo, pos)
	if dressedEquip == nil {
		logger.CtxWarn(ctx, "OnSelectDressMazeEquipRQ pos not unlock", zap.Int32("pos", pos))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("装备位未解锁")
		return
	}
	dressedGuid := module.GetDressedGuid(dressedEquip)
	if dressedGuid > 0 {
		if dressedGuid != replacedGuid {
			logger.CtxWarn(ctx, "OnSelectDressMazeEquipRQ dressed equip not match",
				zap.Int32("pos", pos),
				zap.Int64("dressedGuid", dressedGuid),
				zap.Int64("replacedGuid", replacedGuid))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("已穿戴装备数据不匹配")
			return nil
		}
	} else {
		if replacedGuid > 0 {
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("当前未穿戴装备")
			logger.CtxWarn(ctx, "OnSelectDressMazeEquipRQ no dressed equip",
				zap.Int32("pos", pos),
				zap.Int64("dressedGuid", dressedGuid),
				zap.Int64("replacedGuid", replacedGuid))
			return nil
		}
	}

	euqipInfoCfgRow := GMazeEquipInfoV8Cfg.GetMazeEquipInfoV8Config(dressEquipDb.GetEquipId())
	if euqipInfoCfgRow == nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("未找到装备配置")
		logger.CtxError(ctx, "OnSelectDressMazeEquipRQ 未找到装备配置 ", zap.Int32("equipId", dressEquipDb.GetEquipId()), zap.Int64("equipGuid", dressEquipGuid))
		return err
	}
	if euqipInfoCfgRow.Pos != pos {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("装备位不匹配")
		logger.CtxWarn(ctx, "OnSelectDressMazeEquipRQ equip pos not match",
			zap.Int32("pos", pos),
			zap.Int32("previewPos", euqipInfoCfgRow.Pos),
			zap.Int32("previewEquipId", dressEquipDb.GetEquipId()),
			zap.Int64("previewEquipGuid", dressEquipGuid))
		return nil
	}
	if mazeLv < int64(euqipInfoCfgRow.Level) {
		logger.CtxWarn(ctx, "OnSelectDressMazeEquipRQ maze level no enough",
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

	module.ReplaceEquip(logger, dressedEquip, dressEquipDb)
	chgEquipPosList = append(chgEquipPosList, dressedEquip)
	updateEquipPos = append(updateEquipPos, dressedEquip)
	var effectInfo *calcassembleattr.EquipmentEffectInfo
	effectInfo, err = calcassembleattr.CalcEquipEffect(logger, assembleInfo.MazeEquips)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		logger.CtxError(ctx, "OnSelectDressMazeEquipRQ CalcEquipEffect fail ",
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

	record := StartEquipAssmebleRecord(userId, pos, recordType, dressedEquip,
		oldEquipPos, oldEffect)

	var opCode int32 = 1
	var opMask int32

	defer func() {
		EndEquipAssmebleRecord(ctx, record, opCode, opMask, effectInfo)
	}()
	// 保存装配数据
	err = dollassemblesuitredis.SaveEquipAssembleInfoV2(logger, userId, assembleInfo.GetCurSuitIndex(), updateEquipPos)
	if err != nil {
		res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
		opMask |= demconstdef.DollEquipAssembleOpMaskDbSave
		logger.CtxError(ctx, "OnSelectDressMazeEquipRQ SaveEquipAssembleInfoV2 fail", zap.Error(err), zap.Any("updateEquipPos", updateEquipPos))
		return err
	}
	opCode = 0

	// 通知背包服务
	e := NotifyBagSvrV2(logger, userId, pos, 0, dressEquipGuid, dressedGuid, req.GetAbandonGuids())
	if e != nil {
		opMask |= demconstdef.DollEquipAssembleOpMaskNotifyBag
		logger.CtxError(ctx, "OnDressMazeEquipRQ NotifyBagSvr fail", zap.Error(e),
			zap.Int64("upGuid", dressEquipGuid), zap.Int64("downGuid", dressedGuid))
	}

	var chgMask int32

	// 更新buff中心
	e = mazebuffinforedis.SaveMazeEquipBuff(logger, userId, effectInfo.Other)
	if e != nil {
		opMask |= demconstdef.DollEquipAssembleOpMaskNonForceBuff
	} else {
		// 通知计算属性
		calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{}
		calcAttrNotify.FromServer = demconstdef.MySvr
		calcAttrNotify.UserId = userId
		calcAttrNotify.ChgType = constdef.MazeBuffChgTypeEquipDress
		calcAttrNotify.BuffSrc = constdef.MazeBuffSrcEquip

		calcAttrNotify.Session = req.GetHeader().GetSession()
		e = mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(ctx, calcAttrNotify)
		if e != nil {
			opMask |= demconstdef.DollEquipAssembleOpMaskCalcAttr
			logger.CtxError(ctx, "OnSelectDressMazeEquipRQ SendDollAttrCalcNotify fail", zap.Error(e))
		}

		mazebuffchgrrecordapi.SendMazeBuffChgRecord(ctx, userId,
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
	assembleidpack.SendAssembleChgID(logger, userId, chgAssembleInfo,
		chgMask, int32(int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_LOAD_EQUIP_INFO)),
		constdef.DollAssembleChgTypeReplaceEquip)

	return nil
}

func NotifyBagSvrV2(logger fklog.FKLogI, userId uint64, pos, src int32, upGuid, downGuid int64, loseGuids []int64) error {
	rq := &MazeEquipSvr.SvrMazeEquipAssembleRQ{}
	rq.UserId = proto.Uint64(userId)
	rq.EquipPos = proto.Int32(pos)
	rq.EquipGuid = proto.Int64(upGuid)
	rq.ReplacedEquipGuid = proto.Int64(downGuid)
	rq.DiscardedGuids = loseGuids
	rq.EquipSrc = proto.Int32(1)
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
