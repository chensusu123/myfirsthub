/*
 * @Author: majian
 * @Date: 2025-03-17 14:41:29
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-28 21:00:53
 */
package equip

import (
	"context"
	"fmt"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/function/maputil"
	"maze_game_server/common/function/packtopb/packequipostopb"
	"maze_game_server/common/function/pbutil"
	"maze_game_server/common/structsdef"
	"maze_game_server/common/tradeno"
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
	"maze_game_server/services/costumeservice"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (ep *Equip) OnDressMazeEquipRQ_10418_10419(s *session.Session, req *MazeGameEquip.MazeDressEquipRQ) (err error) {
	defer fkprometheus.DebugPMT("OnDressMazeEquipRQ")()

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeGameEquip.MazeDressEquipRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.EquipPos = req.EquipPos
	res.OpSrc = req.OpSrc

	userId := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnDressMazeEquipRQ end", zap.Any("res", res))
	}()

	logger.CtxInfo(ctx, "OnDressMazeEquipRQ with", zap.Any("req", req))

	limitKey := MakeLimiterKey(userId, 16181)
	if GtcpLimiter.IsRateLimit(limitKey) {
		logger.CtxWarn(ctx, "OnDressMazeEquipRQ rate limiter", zap.Any("req", req))
		res.ErrInfo = errors.NewCodeError(constdef.DE_ERR_RQ_RATE_LIMITER,
			toastmsgtipexcel.GetToastMsgTip(2018, "操作太频繁")).ToInfo()
		return
	}
	GtcpLimiter.UpdateTime(ctx, limitKey)

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
		logger.CtxWarn(ctx, "OnDressMazeEquipRQ equip pos invalid", zap.Int32("pos", pos))
		return nil
	}

	if opType != 1 && opType != 2 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效的操作类型")
		logger.CtxWarn(ctx, "OnDressMazeEquipRQ op type invalid", zap.Int32("opType", opType))
		return nil
	}

	mazeLv, err := mazeuserlevelredis.GetUserLevel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnDressMazeEquipRQ GetUserLevel fail", zap.Error(err))
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
		logger.CtxWarn(ctx, "OnDressMazeEquipRQ dress equip invalid", zap.Int64("dressGuid", dressEquipGuid))
		return nil
	}

	if dressEquipGuid == replacedGuid {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数无效")
		logger.CtxWarn(ctx, "OnDressMazeEquipRQ equip not match", zap.Int32("pos", pos),
			zap.Int64("dressGuid", dressEquipGuid), zap.Int64("replacedGuid", replacedGuid))
		return nil
	}

	dressEquipDb, err = effectequip.GetEffectEquipInfo(ctx, userId, dressEquipGuid)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		logger.CtxError(ctx, "OnDressMazeEquipRQ get bag equip info fail", zap.Error(err), zap.Int64("dressEquipGuid", dressEquipGuid))
		return err
	}
	if dressEquipDb == nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(toastmsgtipexcel.GetToastMsgTip(2019, "要穿戴装备不存在"))
		logger.CtxWarn(ctx, "OnDressMazeEquipRQ equip not exist", zap.Int64("dressEquipGuid", dressEquipGuid))
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
	assembleInfo, oldEffect, err := dollassembleinfo.GetDollAssembleInfoEx(ctx, userId)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		logger.CtxError(ctx, "OnDressMazeEquipRQ Get Assemble info fail", zap.Error(err))
		return err
	}
	if assembleInfo.GetCurSuitIndex() == 0 {
		res.ErrInfo = errors.NewErrorInfo(constdef.DE_ERR_EQUIP_SUIT_NOT_INIT, "未设置生效装备套")
		logger.CtxError(ctx, "OnDressMazeEquipRQ no set cur suit", zap.Int32("cCurSeq", curSuitSeq))
		return err
	}
	if assembleInfo.GetCurSuitIndex() != curSuitSeq {
		res.ErrInfo = errors.NewErrorInfo(constdef.DE_ERR_EQUIP_SUIT_NOT_MATCH, "套装不匹配")
		logger.CtxWarn(ctx, "OnDressMazeEquipRQ cur suit not match",
			zap.Int32("sCurSeq", assembleInfo.GetCurSuitIndex()), zap.Int32("cCurSeq", curSuitSeq))
		return err
	}
	defer func() {
		// 如果客户端请求穿戴信息和服务器不一致时，带回正确的数据给客户端
		if noMatch && len(res.MazeEquipPos) == 0 {
			for _, chgEquip := range assembleInfo.GetMazeEquips() {
				cliPb, _ := packequipostopb.PackEquipPosPb(ctx, chgEquip, int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_LOAD_EQUIP_INFO))
				if cliPb != nil {
					res.MazeEquipPos = append(res.MazeEquipPos, cliPb)
				}
			}
		}
	}()
	dressedEquip := module.GetEquipPosInfo(assembleInfo, pos)
	if dressedEquip == nil {
		logger.CtxWarn(ctx, "OnDressMazeEquipRQ pos not unlock", zap.Int32("pos", pos))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("装备位未解锁")
		return
	}
	dressedGuid := module.GetDressedGuid(dressedEquip)
	if dressedGuid > 0 {
		if dressedGuid != replacedGuid {
			noMatch = true
			logger.CtxError(ctx, "OnDressMazeEquipRQ dressed equip not match",
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
			logger.CtxWarn(ctx, "OnDressMazeEquipRQ no dressed equip",
				zap.Int32("pos", pos),
				zap.Int64("dressedGuid", dressedGuid),
				zap.Int64("replacedGuid", replacedGuid))
			return nil
		}
	}

	euqipInfoCfgRow := GMazeEquipInfoV8Cfg.GetMazeEquipInfoV8Config(dressEquipDb.GetEquipId())
	if euqipInfoCfgRow == nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("未找到装备配置")
		logger.CtxError(ctx, "OnDressMazeEquipRQ 未找到装备配置 ", zap.Int32("equipId", dressEquipDb.GetEquipId()), zap.Int64("equipGuid", dressEquipGuid))
		return err
	}
	if euqipInfoCfgRow.Pos != pos {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("装备位不匹配")
		logger.CtxWarn(ctx, "OnDressMazeEquipRQ equip pos not match",
			zap.Int32("pos", pos),
			zap.Int32("previewPos", euqipInfoCfgRow.Pos),
			zap.Int32("previewEquipId", dressEquipDb.GetEquipId()),
			zap.Int64("previewEquipGuid", dressEquipGuid))
		return nil
	}
	if mazeLv < int64(euqipInfoCfgRow.Level) {
		logger.CtxWarn(ctx, "OnDressMazeEquipRQ maze level no enough",
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

	module.ReplaceEquip(ctx, dressedEquip, dressEquipDb)
	chgEquipPosList = append(chgEquipPosList, dressedEquip)
	updateEquipPos = append(updateEquipPos, dressedEquip)
	var effectInfo *calcassembleattr.EquipmentEffectInfo
	effectInfo, err = calcassembleattr.CalcEquipEffect(ctx, assembleInfo.MazeEquips)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		logger.CtxError(ctx, "OnDressMazeEquipRQ CalcEquipEffect fail ",
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

	// 穿戴装备引起的技能信息变化包
	equipSkillInfoChange, changed, err := GetEquipSkillInfoChange(ctx, userId, oldEquipPos, dressedEquip)
	if err != nil {
		logger.CtxError(ctx, "OnDressMazeEquipRQ GetEquipSkillInfoChange fail",
			zap.Error(err),
			zap.Uint64("userId", userId),
			zap.Any("oldEquipPos", oldEquipPos),
			zap.Any("dressedEquip", dressedEquip),
		)
	} else if changed {
		defer func() {
			online.ClusterPush(context.TODO(), userId, 10510, equipSkillInfoChange)
		}()
	}

	record := StartEquipAssmebleRecord(ctx, userId, pos, recordType, dressedEquip,
		oldEquipPos, oldEffect)

	var opCode int32 = 1
	var opMask int32

	defer func() {
		EndEquipAssmebleRecord(ctx, record, opCode, opMask, effectInfo)
	}()
	// 保存装配数据
	err = dollassemblesuitredis.SaveEquipAssembleInfoV2(ctx, userId, assembleInfo.GetCurSuitIndex(), updateEquipPos)
	if err != nil {
		res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
		opMask |= demconstdef.DollEquipAssembleOpMaskDbSave
		logger.CtxError(ctx, "OnDressMazeEquipRQ SaveEquipAssembleInfoV2 fail", zap.Error(err), zap.Any("updateEquipPos", updateEquipPos))
		return err
	}
	opCode = 0

	//是否自动分解
	//if req.GetIsAutoDismantle() {
	//	// 获取装备可以出售获得的材料
	//	award := make(map[int32]int64)
	//	var equipAward map[int32]int64
	//	equipId := dressedEquip.GetEquipInfo().GetEquipId()
	//	equipCfg := GMazeEquipInfoV8Cfg.Get(equipId)
	//	if equipCfg == nil {
	//		logger.CtxError(ctx, "OnDressMazeEquipRQ get equip cfg fail", zap.Any("equipId", equipId))
	//		res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
	//		return
	//	}
	//	equipAward, err = getEquipDismantle(ctx, dressedGuid, int64(equipId), equipCfg)
	//	if err != nil {
	//		logger.CtxError(ctx, "OnDressMazeEquipRQ get dismantle award fail", zap.Any("equipGuid", dressedGuid), zap.Error(err))
	//		res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
	//		return
	//	}
	//
	//	for k, v := range equipAward {
	//		award[k] += v
	//	}
	//	awardItems := itemutil.Map2Common(award)
	//	res.DismantleAward = awardItems
	//
	//	if len(awardItems) > 0 {
	//		// 699	UN_CGK_COMMON_BILL_TYPE_699	迷宫分解装备
	//		tradeNo := tradeno.GetTradeNum()
	//		items := itemutil.Map2ItemInfo(award)
	//		errInfo := itemservice.GlobalItemService.AddItem(context.TODO(), userId, itemservice.ItemOpTypeDismantle, tradeNo, items...)
	//		if errInfo != nil {
	//			logger.CtxError(ctx, "OnDressMazeEquipRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("rq", req))
	//		}
	//	}

	//bagEquipMgr := bagmodule.NewBagEquipMgr(ctx, userId)
	//err = bagEquipMgr.LoadBagFromRedis(ctx)
	//if err != nil {
	//	logger.CtxError(ctx, "OnDressMazeEquipRQ LoadBagFromRedis error!", zap.Error(err))
	//	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	//	return
	//}
	//delEquipList := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
	//delEquipList = append(delEquipList, dressedEquip.EquipInfo)
	//
	//bagEquipMgr.AddBagRem(dressedEquip.EquipInfo)
	//err = bagEquipMgr.SaveBagInfoToRedis(ctx)
	//if err != nil {
	//	tradeNo := tradeno.GetTradeNum()
	//	logger.CtxError(ctx, "OnDressMazeEquipRQ SaveBagInfoToRedis error!", zap.Error(err))
	//	res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
	//	PushMazeEquipBagLogEx(ctx, userId, nil, delEquipList, tradeNo, 5, mazeequipbagrecord.MazeDelEquip, 0, 1)
	//	return
	//}

	//分解了，只通知穿戴的装备
	//	arr := make([]int64, 0)
	//	arr = append(arr, dressedGuid)
	//	e := NotifyBagSvrV2(ctx, userId, pos, 0, dressEquipGuid, 0, arr)
	//	if e != nil {
	//		opMask |= demconstdef.DollEquipAssembleOpMaskNotifyBag
	//		logger.CtxError(ctx, "OnDressMazeEquipRQ NotifyBagSvr fail", zap.Error(e),
	//			zap.Int64("upGuid", dressEquipGuid), zap.Int64("downGuid", dressedGuid))
	//	}
	//} else {
	// 通知背包服务
	e := NotifyBagSvr(ctx, userId, pos, 0, dressEquipGuid, dressedGuid)
	if e != nil {
		opMask |= demconstdef.DollEquipAssembleOpMaskNotifyBag
		logger.CtxError(ctx, "OnDressMazeEquipRQ NotifyBagSvr fail", zap.Error(e),
			zap.Int64("upGuid", dressEquipGuid), zap.Int64("downGuid", dressedGuid))
	}
	//}

	if req.GetIsAutoDismantle() {
		// 分解装备
		guidsBag := make([]int64, 0)
		guidsBag = append(guidsBag, dressedGuid)
		tradeNo := tradeno.GetTradeNum()
		rqSale := &MazeEquipSvr.SvrMazeEquipSaleRQ{
			UserId:     proto.Uint64(userId),
			EquipGuids: guidsBag,
			TradeNum:   proto.Uint64(tradeNo),
		}
		rsSale := &MazeEquipSvr.SvrMazeEquipSaleRS{}
		logger.CtxError(ctx, "OnDollEquipDismantleRQ SvrDollEquipSaleRS dump", zap.Any("rqSale", rqSale), zap.Any("rsSale", rsSale))
		err = OnSvrDollEquipSaleRQ(ctx, int64(userId), rqSale, rsSale)
		if err != nil {
			logger.CtxError(ctx, "OnDollEquipDismantleRQ SvrDollEquipSaleRS fail", zap.Error(err), zap.Any("rq", rqSale))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}

		award := make(map[int32]int64)
		var equipAward map[int32]int64
		equipId := dressedEquip.GetEquipInfo().GetEquipId()
		equipCfg := GMazeEquipInfoV8Cfg.Get(equipId)
		if equipCfg == nil {
			logger.CtxError(ctx, "OnDressMazeEquipRQ get equip cfg fail", zap.Any("equipId", equipId))
			res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
			return
		}
		equipAward, err = getEquipDismantle(ctx, dressedGuid, int64(equipId), equipCfg)
		if err != nil {
			logger.CtxError(ctx, "OnDressMazeEquipRQ get dismantle award fail", zap.Any("equipGuid", dressedGuid), zap.Error(err))
			res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
			return
		}

		for k, v := range equipAward {
			award[k] += v
		}
		awardItems := itemutil.Map2Common(award)
		res.DismantleAward = awardItems
	}

	var chgMask int32

	// 更新buff中心

	e = mazebuffinforedis.SaveMazeEquipBuff(ctx, userId, effectInfo.Other)

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
			logger.CtxError(ctx, "OnDressMazeEquipRQ SendDollAttrCalcNotify fail", zap.Error(e))
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
	assembleidpack.SendAssembleChgID(ctx, userId, chgAssembleInfo,
		chgMask, int32(int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_LOAD_EQUIP_INFO)),
		constdef.DollAssembleChgTypeReplaceEquip)

	// 换装备推送装扮变化id包
	costumeservice.GlobalCostumeService.ChangeCostume(context.TODO(), userId)
	return nil
}

func NotifyBagSvr(ctx context.Context, userId uint64, pos, src int32, upGuid, downGuid int64) error {
	rq := &MazeEquipSvr.SvrMazeEquipAssembleRQ{}
	rq.UserId = proto.Uint64(userId)
	rq.EquipPos = proto.Int32(pos)
	rq.EquipGuid = proto.Int64(upGuid)
	rq.ReplacedEquipGuid = proto.Int64(downGuid)
	rq.EquipSrc = proto.Int32(src)
	rs := &MazeEquipSvr.SvrMazeEquipAssembleRS{}
	// e := dollequipbagrpc.MazeEquipAssembleRQ(ctx, rq, rs)
	e := OnSvrMazeEquipAssembleRQ(ctx, int64(userId), rq, rs)
	if e != nil {
		return e
	}
	if rs.GetErrInfo().GetErrCode() == errors.NO_ERROR_CODE {
		return nil
	}
	return errors.New("请求背包失败")
}

func StartEquipAssmebleRecord(ctx context.Context, userId uint64, pos, op int32, newEquip, oldEquip *MazeEquipCache.MazeEquipPosInfo, oldEffect *calcassembleattr.EquipmentEffectInfo) *dollequipassmeblekakfa.MazeGameEquipAssembleRecord {
	record := new(dollequipassmeblekakfa.MazeGameEquipAssembleRecord)
	record.UserId = userId
	record.EquipPos = pos
	var oldExtra, newExtra string
	if newEquip != nil {
		record.NewEquipId = newEquip.GetEquipLoadInfo().GetEquipId()
		record.NewGuid = uint64(newEquip.GetEquipLoadInfo().GetEquipGuid())
		newResId, newEquipName := pbutil.GetDollEquipName(ctx, newEquip.GetEquipInfo(), 0)
		newExtra = MakeExtra(ctx, newResId, newEquipName)
	}
	if oldEquip != nil {
		record.OldEquipId = oldEquip.GetEquipLoadInfo().GetEquipId()
		record.OldGuid = uint64(oldEquip.GetEquipLoadInfo().GetEquipGuid())
		oldResId, oldEquipName := pbutil.GetDollEquipName(ctx, oldEquip.GetEquipInfo(), 0)
		oldExtra = MakeExtra(ctx, oldResId, oldEquipName)
	}
	record.OpType = op
	hsStr := maputil.MapToString32(oldEffect.GetSuitCalc().GetSuitNumMap())
	//	fiveStr := maputil.MapToString32(oldEffect.FiveStateMap)
	record.OldFElem = fmt.Sprintf("hurtSuit:%s|%s", hsStr, oldExtra)
	record.NewFElem = newExtra
	return record
}

func EndEquipAssmebleRecord(ctx context.Context, record *dollequipassmeblekakfa.MazeGameEquipAssembleRecord,
	opRet, oMask int32, newEffect *calcassembleattr.EquipmentEffectInfo) {
	record.RetCode = opRet
	record.CodeMask = oMask

	// fiveElemStr := maputil.MapToString32(newEffect.FiveStateMap)
	hsStr := maputil.MapToString32(newEffect.GetSuitCalc().GetSuitNumMap())
	record.NewFElem = fmt.Sprintf("hurtSuit:%s|%s", hsStr, record.NewFElem)
	dollequipassmeblekakfa.SendMazeGameEquipAssembleRecord(ctx, record)
}

func MakeExtra(ctx context.Context, resId int32, equipName string) string {
	if equipName == "" {
		logger := fklog.ContextAppLogger(ctx)
		logger.CtxWarn(ctx, "MakeExtra equipName is empty", zap.Int32("resId", resId))
	}
	if resId > 0 {
		return fmt.Sprintf("equipName:%s_resId:%d", equipName, resId)
	} else {
		return fmt.Sprintf("equipName:%s", equipName)
	}
}
