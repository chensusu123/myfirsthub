package equip

import (
	"fmt"
	"strings"

	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeEquipPosLvSuiteV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeEquipPosLvV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/protodef/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/protodef/Common"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipPos"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGameEquip"

	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassembleredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/toastmsgtipexcel"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/equipposexcel"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeuserlevelredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/itemutil"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/excelutil"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/assemble"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/equippossuit"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/uniqueid"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/gentradeno"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/equipposstrengrecordkafka"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/assembleidpack"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/maputil"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebuffinforedis"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structsdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/calcassembleattr"
)

// 装备位强化
func OnEquipPosLvUpRQ(logger fknet.TCPContext, userId uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnEquipPosLvUpRQ")()
	rq := rqMsg.(*MazeEquipPos.MazeEquipPosLvUpRQ)
	rs := rsMsg.(*MazeEquipPos.MazeEquipPosLvUpRS)

	rs.ErrInfo = errors.NO_ERROR
	rs.Header = rq.Header
	rs.PosId = rq.PosId
	rs.OutsideBarrier = rq.OutsideBarrier

	posId := rq.GetPosId()
	rqCurLv := rq.GetLevel()
	rqCost := rq.GetCostItems()

	logger.InfoWF("OnEquipPosLvUpRQ start", zap.Any("rq", rq))
	defer func() {
		logger.InfoWF("OnEquipPosLvUpRQ end", zap.Any("rs", rs))
	}()

	if posId < 1 || posId > constdef.EquipPosNum {
		logger.ErrorWF("OnEquipPosLvUpRQ param invalid", zap.Int32("posId", posId))
		rs.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
		return
	}

	assembleDb, err := dollassembleredis.GetAllAssembleInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("OnEquipPosLvUpRQ GetAllAssembleInfo", zap.Error(err))
		rs.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	curSuitId := assembleDb.GetEpEnSuitId()

	equipPos := assembleDb.GetMazeEquips()
	if len(equipPos) < constdef.EquipPosNum { // 未解锁全部装备位
		logger.WarnWF("OnEquipPosStrengPreviewRQ not unlock all equip pos",
			zap.Int("unlockNum", len(equipPos)), zap.Int("needNum", constdef.EquipPosNum))
		rs.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(toastmsgtipexcel.GetToastMsgTip(2024, "解锁全部装备位后可强化"))
		return
	}

	var curLv, otherMinLv int32
	otherMinLvIsInit := true // 标识otherMinLv是否是初始值，初始值和初始装备位等级冲突
	var curEquipPosDb *MazeEquipCache.MazeEquipSlotDb
	for _, v := range equipPos {
		lv := v.GetEquipPos().GetLevel()
		if otherMinLvIsInit || lv < otherMinLv {
			otherMinLv = lv
			otherMinLvIsInit = false
		}
		if posId == v.GetEquipPos().GetPos() {
			curLv = lv
			curEquipPosDb = v.GetEquipPos()
		}
	}

	if rqCurLv != curLv {
		logger.ErrorWF("OnEquipPosLvUpRQ param invalid", zap.Int32("rqCurLv", rqCurLv), zap.Int32("svrCurLv", curLv))
		rs.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
		return
	}

	if curEquipPosDb == nil {
		logger.ErrorWF("OnEquipPosLvUpRQ param invalid, no current equip position", zap.Int32("posId", posId))
		rs.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
		return
	}

	curCfg := equipposexcel.GetPosStrengthCfg(posId, curLv)
	if curCfg == nil {
		logger.ErrorWF("OnEquipPosLvUpRQ GetPosStrengthCfg", zap.Error(err),
			zap.Int32("posId", posId), zap.Int32("curLv", curLv))
		rs.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		return
	}
	if curCfg.Next_order < 0 {
		logger.WarnWF("OnEquipPosLvUpRQ equip pos is max level, client should filter",
			zap.Int32("curLv", curLv))
		rs.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("装备位已达到最高等级")
		return
	}
	nextCfg := equipposexcel.GetPosStrengthCfgByKey(curCfg.Next_order)
	if nextCfg == nil {
		logger.ErrorWF("OnEquipPosLvUpRQ cannot found next cfg",
			zap.Int32("curLv", curLv), zap.Int32("key", curCfg.Next_order))
		rs.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		return
	}
	if otherMinLv < curCfg.Need_other {
		logger.ErrorWF("OnEquipPosLvUpRQ other equip pos min level limit, client should filter",
			zap.Int32("curLv", curLv),
			zap.Int32("pos", posId),
			zap.Int32("otherMinLv", otherMinLv),
			zap.Int32("needOtherMinLv", curCfg.Need_other))
		rs.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(strings.ReplaceAll(toastmsgtipexcel.GetToastMsgTip(2011, "所有部位达到N级才可继续强化"), "N", fmt.Sprintf("%d", curCfg.Need_other)))
		return
	}

	mazeLv, err := mazeuserlevelredis.GetUserLevel(logger, userId)
	if err != nil {
		logger.ErrorWF("OnEquipPosLvUpRQ GetUserLevel", zap.Error(err))
		rs.ErrInfo = errors.MODULE_ERROR.Wrap("查询迷宫等级失败")
		return nil
	}
	if mazeLv < int64(curCfg.Need_maze_level) {
		logger.WarnWF("OnEquipPosLvUpRQ insufficient doll lv",
			zap.Int64("mazeLv", mazeLv),
			zap.Int32("needLv", curCfg.Need_maze_level))
		rs.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(strings.ReplaceAll(toastmsgtipexcel.GetToastMsgTip(2012, "等级达到N级才可继续强化"), "N", fmt.Sprintf("%d", curCfg.Need_maze_level)))
		return nil
	}
	svrCost := itemutil.Map2Common(curCfg.Cost)
	if !itemutil.CheckItemMatch(rqCost, svrCost) {
		logger.ErrorWF("OnEquipPosLvUpRQ param invalid", zap.Any("rqCost", rqCost),
			zap.Any("svrCost", svrCost))
		rs.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
		return
	}
	var careCost []*MazeCommon.MazeItem
	for _, item := range svrCost {
		if !rq.GetOutsideBarrier() {
			if item.GetItemId() == constdef.MazeCommonItemCoin {
				continue
			}
		}
		careCost = append(careCost, item)
	}
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

	// 设置装备位等级 设置套装id pk段位
	chgAssemDb := &MazeEquipCache.MazeAssembleDb{}
	var mask int32
	var fields []string
	curEquipPosDb.Level = proto.Int32(nextCfg.Level)
	chgAssemDb.MazeEquips = append(chgAssemDb.MazeEquips, &MazeEquipCache.MazeEquipPosInfo{
		EquipPos: curEquipPosDb,
	})
	fields = append(fields, assemble.EnCodeAssemblePosField(posId))
	mask |= int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_POS_MASK)

	curSuit, nextSuit, err := equippossuit.GetCurAndNextSuit(logger, curSuitId)
	if err != nil {
		logger.ErrorWF("OnEquipPosLvUpRQ GetCurAndNextSuit", zap.Error(err))
		rs.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	newSuitId, err := equippossuit.CalcPosSuit(logger, equipPos)
	if err != nil {
		logger.ErrorWF("OnEquipPosLvUpRQ CalcPosSuit", zap.Error(err))
		rs.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if newSuitId != curSuitId {
		chgAssemDb.EpEnSuitId = proto.Int32(newSuitId)
		fields = append(fields, constdef.AssemblePrefixEquipPosEnSuit)
		mask |= int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_ST_SUIT_MASK)
	}

	// 扣物品
	tid := uniqueid.GenUniqueIdUInt64()
	if len(careCost) > 0 {
		// 通用	693	UN_CGK_COMMON_BILL_TYPE_693	迷宫装备位强化		否	马健	2025-03-22 17:28:42
		errInfo := gentradeno.DeductItemsEx(logger, userId, 693, tid, careCost...)
		if errInfo != nil {
			logger.ErrorWF("OnEquipPosLvUpRQ DeductItemsEx",
				zap.Any("svrCost", svrCost),
				zap.Any("careCost", careCost),
				zap.Uint64("tid", tid),
				zap.Any("errInfo", errInfo))

			if errInfo.GetErrCode() == 50049 {
				rs.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("物品不足")
			} else {
				rs.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("扣物品失败")
			}
			return
		}
	}
	logger.InfoWF("OnEquipPosLvUpRQ DeductItemsEx succ",
		zap.Any("svrCost", svrCost),
		zap.Any("careCost", careCost),
		zap.Uint64("tid", tid))
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
		if err := equipposstrengrecordkafka.PushEquipPosStrengRecord(logger, userId,
			posId, curLv, nextCfg.Level, curSuitId, newSuitId, tid, svrCost, result, retMask); err != nil {
			logger.ErrorWF("OnEquipPosLvUpRQ PushEquipPosStrengRecord", zap.Error(err))
		}
	}()

	// 强化
	if err = dollassembleredis.SetAssembleInfoByFields(logger, userId, fields, chgAssemDb); err != nil {
		logger.ErrorWF("OnEquipPosLvUpRQ SetAssembleInfoByFields", zap.Error(err))
		rs.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
		result = 0
		retMask |= MazeEquipPosErrDB
		return
	}
	logger.InfoWF("OnEquipPosLvUpRQ SetAssembleInfoByFields succ", zap.Any("chgAssemDb", chgAssemDb))
	// 计算属性加成
	retMask |= UpdateEquipPosBuff(logger, userId, curSuitId, newSuitId,
		assembleDb.GetMazeEquips(), posCurAttrs, rq.GetHeader().GetSession())

	// 推装配信息变化包
	assembleidpack.SendAssembleChgID(logger, userId, chgAssemDb,
		mask, int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_STRENGTHEN_INFO), constdef.DollAssembleChgTypeEquipPosEnhancement)

	// 封装响应数据
	rsCurCfg := nextCfg
	var rsNextCfg *GMazeEquipPosLvV8Cfg.MazeEquipPosLvV8ConfigRow
	if rsCurCfg.Next_order > 0 {
		rsNextCfg = equipposexcel.GetPosStrengthCfgByKey(rsCurCfg.Next_order)
	}

	rs.StLevel = &Common.AttrChgInfo{
		AttrVal: proto.Int64(int64(rsCurCfg.Level)),
	}
	if nextCfg != nil {
		rs.StLevel.NextVal = proto.Int64(int64(rsNextCfg.Level))
	}

	// 对比属性加成
	for _, id := range rsCurCfg.Show_attr_order {
		if id <= 0 {
			continue
		}
		attrInfo := &Common.AttrChgInfo{
			AttrId:  proto.Int32(id),
			AttrVal: proto.Int64(rsCurCfg.Show_attr[id]),
		}
		if rsNextCfg != nil {
			attrInfo.NextVal = proto.Int64(rsNextCfg.Show_attr[id])
		}
		rs.AttrBuffs = append(rs.AttrBuffs, attrInfo)
	}

	rs.StCond = &MazeGameEquip.EquipPosStCond{
		CostItems:   itemutil.Map2Common(rsCurCfg.Cost),
		NeedOtherLv: proto.Int32(rsCurCfg.Need_other),
		NeedMazeLv:  proto.Int32(rsCurCfg.Need_maze_level),
	}
	if newSuitId == curSuitId { // 套装无变化
		rs.CurSuit, rs.NextSuit = curSuit, nextSuit
	} else {
		var e error
		rs.CurSuit, rs.NextSuit, e = equippossuit.GetCurAndNextSuit(logger, newSuitId)
		if e != nil {
			logger.ErrorWF("OnEquipPosLvUpRQ GetCurAndNextSuit", zap.Error(e),
				zap.Int32("newSuitId", newSuitId))
		}
	}
	return
}

func UpdateEquipPosBuff(logger fklog.FKLogI, userId uint64,
	curSuitId, newSuitId int32, posInfo []*MazeEquipCache.MazeEquipPosInfo, posCurAttrs map[int32]int64, session string) (result int32) {
	reals, shows := CalcEquipPosBuffs(posInfo)
	posNewAttrs := map[int32]int64{} // 装备位升级后属性
	posNewShowAttrs := map[int32]int64{}

	maputil.Int64MapAppend(posNewAttrs, reals)
	maputil.Int64MapAppend(posNewShowAttrs, shows)

	reals, shows = CalcEquipPosSuitBuffs(newSuitId)

	maputil.Int64MapAppend(posNewAttrs, reals)
	maputil.Int64MapAppend(posNewAttrs, shows)

	err := mazebuffinforedis.SaveMazeEquipPosBuff(logger, userId, calcassembleattr.PackMapAttrAll(posNewAttrs, posNewShowAttrs))
	if err != nil {
		result |= MazeEquipPosErrSaveBuff
		logger.ErrorWF("UpdateEquipPosBuff SaveMazeEquipPosBuff err", zap.Error(err))
	} else {
		calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{}
		calcAttrNotify.FromServer = PosMySvr
		calcAttrNotify.UserId = userId
		calcAttrNotify.ChgType = constdef.MazeBuffEquipPosUpgrade
		calcAttrNotify.Session = session
		err = mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, calcAttrNotify)
		if err != nil {
			result |= MazeEquipPosErrBuffNotify
			logger.ErrorWF("UpdateEquipPosBuff SendMazeAttrCalcNotify fail", zap.Error(err))
		}
	}
	return
}

func CalcEquipPosBuffs(equipPos []*MazeEquipCache.MazeEquipPosInfo) (reals, shows map[int32]int64) {
	reals = make(map[int32]int64)
	shows = make(map[int32]int64)

	for _, v := range equipPos {
		id := v.GetEquipPos().GetPos()
		lv := v.GetEquipPos().GetLevel()
		posCfg := GMazeEquipPosLvV8Cfg.Get(excelutil.GetEquipPosEnLevelKey(id, lv))
		if posCfg == nil {
			continue
		}
		for k, v := range posCfg.Add_attr {
			reals[k] += v
		}
		for k, v := range posCfg.Show_attr {
			shows[k] += v
		}
	}
	return
}

func CalcEquipPosSuitBuffs(suitId int32) (reals, shows map[int32]int64) {
	reals = make(map[int32]int64)
	shows = make(map[int32]int64)
	cfg := GMazeEquipPosLvSuiteV8Cfg.Get(suitId)
	if cfg != nil {
		for k, v := range cfg.Add_attr {
			reals[k] += v
		}
		for k, v := range cfg.Show_attr {
			reals[k] += v
		}
	}
	return
}
