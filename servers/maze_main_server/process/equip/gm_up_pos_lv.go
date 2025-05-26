package equip

import (
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/common/constdef"
	"maze_game_server/common/function/assemble"
	"maze_game_server/common/function/excelutil"
	"maze_game_server/config/GMazeEquipPosLvV8Cfg"
	"maze_game_server/excel/equipposexcel"
	"maze_game_server/io/kafka/equipposstrengrecordkafka"
	"maze_game_server/io/redis/dollassembleredis"
	"maze_game_server/io/redis/mazeuserlevelredis"
	"maze_game_server/module/assembleidpack"
	"maze_game_server/module/equippossuit"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"
)

// OnGmEquipPosLvUp 不扣道具，直接将用户装备提升到某个等级 ！！！谨慎使用
func OnGmEquipPosLvUp(logger fklog.FKLogI, userId uint64, targetLv int32) (err error) {
	defer fkprometheus.DebugPMT("OnGmEquipPosLvUp")()
	logger.InfoWF("OnGmEquipPosLvUp start", zap.Any("userId", userId),
		zap.Int32("targetLv", targetLv))
	defer func() {
		logger.InfoWF("OnGmEquipPosLvUp end", zap.Any("err", err))
	}()

	if userId == 0 {
		err = fmt.Errorf("userId is 0")
		logger.ErrorWF("OnGmEquipPosLvUp userId is 0")
		return
	}

	if targetLv == 0 {
		err = fmt.Errorf("targetLv is 0")
		logger.ErrorWF("OnGmEquipPosLvUp targetLv is 0")
		return
	}

	assembleDb, err := dollassembleredis.GetAllAssembleInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("OnGmEquipPosLvUp GetAllAssembleInfo", zap.Error(err))
		return
	}

	curSuitId := assembleDb.GetEpEnSuitId()

	equipPos := assembleDb.GetMazeEquips()
	if len(equipPos) < constdef.EquipPosNum { // 未解锁全部装备位
		err = fmt.Errorf("OnGmEquipPosLvUp not unlock all equip pos")
		logger.WarnWF("OnGmEquipPosLvUp not unlock all equip pos",
			zap.Int("unlockNum", len(equipPos)), zap.Int("needNum", constdef.EquipPosNum))
		return
	}

	// var curLv, otherMinLv int32
	// otherMinLvIsInit := true // 标识otherMinLv是否是初始值，初始值和初始装备位等级冲突
	// var curEquipPosDb *MazeEquipCache.MazeEquipSlotDb
	// for _, v := range equipPos {
	// 	lv := v.GetEquipPos().GetLevel()
	// 	if otherMinLvIsInit || lv < otherMinLv {
	// 		otherMinLv = lv
	// 		otherMinLvIsInit = false
	// 	}
	// 	if posId == v.GetEquipPos().GetPos() {
	// 		curLv = lv
	// 		curEquipPosDb = v.GetEquipPos()
	// 	}
	// }

	// if rqCurLv != curLv {
	// 	logger.ErrorWF("OnEquipPosLvUpRQ param invalid", zap.Int32("rqCurLv", rqCurLv), zap.Int32("svrCurLv", curLv))
	// 	rs.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
	// 	return
	// }

	// if curEquipPosDb == nil {
	// 	logger.ErrorWF("OnEquipPosLvUpRQ param invalid, no current equip position", zap.Int32("posId", posId))
	// 	rs.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
	// 	return
	// }

	mazeLv, err := mazeuserlevelredis.GetUserLevel(logger, userId)
	if err != nil {
		err = fmt.Errorf("OnGmEquipPosLvUp GetUserLevel")
		logger.ErrorWF("OnEquipPosLvUpRQ GetUserLevel", zap.Error(err))
		return nil
	}
	// 所有部位全部检查一遍
	for _, v := range equipPos {
		posId := v.GetEquipPos().GetPos()
		posLv := v.GetEquipPos().GetLevel()
		cfg := equipposexcel.GetPosStrengthCfg(posId, posLv)
		if cfg == nil {
			err = fmt.Errorf("cannot found pos cfg")
			logger.ErrorWF("OnGmEquipPosLvUp GetPosStrengthCfg", zap.Error(err),
				zap.Int32("posId", posId), zap.Int32("posLv", posLv))
			return
		}
		// if cfg.Next_order < 0 {
		// 	logger.WarnWF("OnGmEquipPosLvUp equip pos is max level, client should filter",
		// 		zap.Any("posId", posId), zap.Int32("posLv", posLv))
		// 	continue
		// }
		if mazeLv < int64(cfg.Need_maze_level) {
			err = fmt.Errorf("OnGmEquipPosLvUp insufficient maze lv, curMazeLv:%d, needMazeLv:%d", mazeLv, cfg.Need_maze_level)
			logger.WarnWF("OnGmEquipPosLvUp insufficient maze lv",
				zap.Int64("mazeLv", mazeLv),
				zap.Int32("needLv", cfg.Need_maze_level))
			return
		}
	}

	// nextCfg := equipposexcel.GetPosStrengthCfgByKey(curCfg.Next_order)
	// if nextCfg == nil {
	// 	logger.ErrorWF("OnEquipPosLvUpRQ cannot found next cfg",
	// 		zap.Int32("curLv", curLv), zap.Int32("key", curCfg.Next_order))
	// 	rs.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
	// 	return
	// }
	// if otherMinLv < curCfg.Need_other {
	// 	logger.ErrorWF("OnEquipPosLvUpRQ other equip pos min level limit, client should filter",
	// 		zap.Int32("curLv", curLv),
	// 		zap.Int32("pos", posId),
	// 		zap.Int32("otherMinLv", otherMinLv),
	// 		zap.Int32("needOtherMinLv", curCfg.Need_other))
	// 	rs.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(strings.ReplaceAll(toastmsgtipexcel.GetToastMsgTip(2011, "所有部位达到N级才可继续强化"), "N", fmt.Sprintf("%d", curCfg.Need_other)))
	// 	return
	// }

	// svrCost := itemutil.Map2Common(curCfg.Cost)
	// if !itemutil.CheckItemMatch(rqCost, svrCost) {
	// 	logger.ErrorWF("OnEquipPosLvUpRQ param invalid", zap.Any("rqCost", rqCost),
	// 		zap.Any("svrCost", svrCost))
	// 	rs.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
	// 	return
	// }
	// var careCost []*MazeCommon.MazeItem
	// for _, item := range svrCost {
	// 	if !rq.GetOutsideBarrier() {
	// 		if item.GetItemId() == constdef.MazeCommonItemCoin {
	// 			continue
	// 		}
	// 	}
	// 	careCost = append(careCost, item)
	// }
	posCurAttrs := map[int32]int64{} // 装备位升级前属性
	for _, v := range assembleDb.GetMazeEquips() {
		id := v.GetEquipPos().GetPos()
		lv := v.GetEquipPos().GetLevel()
		posCfg := GMazeEquipPosLvV8Cfg.Get(excelutil.GetEquipPosEnLevelKey(id, lv))
		if posCfg == nil {
			continue
		}
		for k, v := range posCfg.Add_attr {
			posCurAttrs[k] += v
		}
	}

	var (
		mask       int32
		fields     = make([]string, 0)
		chgAssemDb = &MazeEquipCache.MazeAssembleDb{}
	)

	for _, v := range equipPos {
		posId := v.GetEquipPos().GetPos()
		posLv := v.GetEquipPos().GetLevel()
		cfg := equipposexcel.GetPosStrengthCfg(posId, posLv)
		if cfg == nil {
			err = fmt.Errorf("cannot found pos cfg")
			logger.ErrorWF("OnGmEquipPosLvUp GetPosStrengthCfg", zap.Error(err),
				zap.Int32("posId", posId), zap.Int32("posLv", posLv))
			return
		}
		if cfg.Next_order < 0 {
			logger.WarnWF("OnGmEquipPosLvUp equip pos is max level, client should filter",
				zap.Any("posId", posId), zap.Int32("posLv", posLv))
			continue
		}
		v.EquipPos.Level = proto.Int32(targetLv)
		chgAssemDb.MazeEquips = append(chgAssemDb.MazeEquips, &MazeEquipCache.MazeEquipPosInfo{
			EquipPos: v.GetEquipPos(),
		})
		fields = append(fields, assemble.EnCodeAssemblePosField(posId))
		mask |= int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_POS_MASK)
	}

	// curSuit, nextSuit, err := equippossuit.GetCurAndNextSuit(logger, curSuitId)
	// if err != nil {
	// 	logger.ErrorWF("OnEquipPosLvUpRQ GetCurAndNextSuit", zap.Error(err))
	// 	rs.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	newSuitId, err := equippossuit.CalcPosSuit(logger, equipPos)
	if err != nil {
		err = fmt.Errorf("OnGmEquipPosLvUp CalcPosSuit")
		logger.ErrorWF("OnEquipPosLvUpRQ CalcPosSuit", zap.Error(err))
		return
	}
	if newSuitId != curSuitId {
		chgAssemDb.EpEnSuitId = proto.Int32(newSuitId)
		fields = append(fields, constdef.AssemblePrefixEquipPosEnSuit)
		mask |= int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_ST_SUIT_MASK)
	}

	// // 扣物品
	// tid := uniqueid.GenUniqueIdUInt64()
	// if len(careCost) > 0 {
	// 	// 通用	693	UN_CGK_COMMON_BILL_TYPE_693	迷宫装备位强化		否	马健	2025-03-22 17:28:42
	// 	errInfo := gentradeno.DeductItemsEx(logger, userId, 693, tid, careCost...)
	// 	if errInfo != nil {
	// 		logger.ErrorWF("OnEquipPosLvUpRQ DeductItemsEx",
	// 			zap.Any("svrCost", svrCost),
	// 			zap.Any("careCost", careCost),
	// 			zap.Uint64("tid", tid),
	// 			zap.Any("errInfo", errInfo))

	// 		if errInfo.GetErrCode() == 50049 {
	// 			rs.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("物品不足")
	// 		} else {
	// 			rs.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("扣物品失败")
	// 		}
	// 		return
	// 	}
	// }
	// logger.InfoWF("OnEquipPosLvUpRQ DeductItemsEx succ",
	// 	zap.Any("svrCost", svrCost),
	// 	zap.Any("careCost", careCost),
	// 	zap.Uint64("tid", tid))
	// remainMoney, err := mazemoney.GetUserMoney(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("OnEquipPosLvUpRQ GetUserMoney",
	// 		zap.Error(err))
	// 	rs.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("检查货币失败")
	// 	return
	// }
	// if remainMoney < curCfg.Cost[mazemoney.MONEY_ID] {
	// 	logger.WarnWF("OnEquipPosLvUpRQ money less",
	// 		zap.Int32("pos", posId),
	// 		zap.Int32("curLv", curLv),
	// 		zap.Any("needCost", curCfg.Cost),
	// 		zap.Int64("has", remainMoney))
	// 	rs.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("货币不足")
	// 	return
	// }
	// mazemoney.SubUserMoney(logger, userId, curCfg.Cost[mazemoney.MONEY_ID],
	// 	mazemoneykafka.MoneyChgReasonEquipStrenth)
	// if err != nil {
	// 	logger.ErrorWF("OnEquipPosLvUpRQ SubUserMoney",
	// 		zap.Error(err), zap.Any("Cost", curCfg.Cost))
	// 	rs.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("扣货币失败")
	// 	return
	// }

	// 记录流水
	var result int32 = 1
	var retMask int32 = 0
	defer func() {
		// 整合为一条流水(不扣物品，忽略流水号)
		if err := equipposstrengrecordkafka.PushEquipPosStrengRecord(logger, userId,
			0, 0, targetLv, curSuitId, newSuitId, 0, nil, result, retMask); err != nil {
			logger.ErrorWF("OnEquipPosLvUpRQ PushEquipPosStrengRecord", zap.Error(err))
		}
	}()

	// 强化
	if err = dollassembleredis.SetAssembleInfoByFields(logger, userId, fields, chgAssemDb); err != nil {
		logger.ErrorWF("OnEquipPosLvUpRQ SetAssembleInfoByFields", zap.Error(err))
		// result = 0
		retMask |= MazeEquipPosErrDB
		return
	}
	logger.InfoWF("OnEquipPosLvUpRQ SetAssembleInfoByFields succ", zap.Any("chgAssemDb", chgAssemDb))
	// 计算属性加成
	retMask |= UpdateEquipPosBuff(logger, userId, curSuitId, newSuitId,
		assembleDb.GetMazeEquips(), posCurAttrs, "")

	// 推装配信息变化包
	assembleidpack.SendAssembleChgID(logger, userId, chgAssemDb,
		mask, int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_STRENGTHEN_INFO), constdef.DollAssembleChgTypeEquipPosEnhancement)

	// 封装响应数据
	// rsCurCfg := nextCfg
	// var rsNextCfg *GMazeEquipPosLvV8Cfg.MazeEquipPosLvV8ConfigRow
	// if rsCurCfg.Next_order > 0 {
	// 	rsNextCfg = equipposexcel.GetPosStrengthCfgByKey(rsCurCfg.Next_order)
	// }

	// rs.StLevel = &Common.AttrChgInfo{
	// 	AttrVal: proto.Int64(int64(rsCurCfg.Level)),
	// }
	// if nextCfg != nil {
	// 	rs.StLevel.NextVal = proto.Int64(int64(rsNextCfg.Level))
	// }

	// // 对比属性加成
	// for _, id := range rsCurCfg.Show_attr_order {
	// 	if id <= 0 {
	// 		continue
	// 	}
	// 	attrInfo := &Common.AttrChgInfo{
	// 		AttrId:  proto.Int32(id),
	// 		AttrVal: proto.Int64(rsCurCfg.Show_attr[id]),
	// 	}
	// 	if rsNextCfg != nil {
	// 		attrInfo.NextVal = proto.Int64(rsNextCfg.Show_attr[id])
	// 	}
	// 	rs.AttrBuffs = append(rs.AttrBuffs, attrInfo)
	// }

	// rs.StCond = &MazeGameEquip.EquipPosStCond{
	// 	CostItems:   itemutil.Map2Common(rsCurCfg.Cost),
	// 	NeedOtherLv: proto.Int32(rsCurCfg.Need_other),
	// 	NeedMazeLv:  proto.Int32(rsCurCfg.Need_maze_level),
	// }
	// if newSuitId == curSuitId { // 套装无变化
	// 	rs.CurSuit, rs.NextSuit = curSuit, nextSuit
	// } else {
	// 	var e error
	// 	rs.CurSuit, rs.NextSuit, e = equippossuit.GetCurAndNextSuit(logger, newSuitId)
	// 	if e != nil {
	// 		logger.ErrorWF("OnEquipPosLvUpRQ GetCurAndNextSuit", zap.Error(e),
	// 			zap.Int32("newSuitId", newSuitId))
	// 	}
	// }
	return
}
