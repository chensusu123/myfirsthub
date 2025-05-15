package equip

import (
	"fmt"
	"strings"

	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeEquipInfoV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/protodef/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipSvr"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGameEquip"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassemblesuitredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebagequipredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/itemutil"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/uniqueid"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/gentradeno"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/toastmsgtipexcel"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/dollequipdismantlekafka"
)

func OnDollEquipDismantleRQ(ctx fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnDollEquipDismantleRQ")()
	req := rqMsg.(*MazeGameEquip.MazeEquipDismantleRQ)
	res := rsMsg.(*MazeGameEquip.MazeEquipDismantleRS)

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.SelectQuality = req.SelectQuality
	res.SelectGuid = req.SelectGuid
	res.DismantleFrom = req.DismantleFrom
	// res.DismantleAward = req.DismantleAward
	userCtx := fkserver.NewUserContext(ctx.Context, shardingID, ctx.FKLogI)

	defer func() {
		userCtx.InfoWF("OnDollEquipDismantleRQ end", zap.Any("res", res))
	}()

	userCtx.InfoWF("OnDollEquipDismantleRQ with", zap.Any("req", req))

	if len(req.GetSelectGuid()) == 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("缺少选定的装备")
		return
	}

	// if req.GetDismantleFrom() != 0 && req.GetDismantleFrom() != 1 && req.GetDismantleFrom() != 2 {
	// 	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误请重试")
	// 	return
	// }

	selectQuality := make(map[int32]struct{})
	for _, v := range req.GetSelectQuality() {
		if v > 0 {
			selectQuality[v] = struct{}{}
		}
	}

	// 获取身上的装备信息
	assembleInfoMap, err := dollassemblesuitredis.GetAllDollAssembleSuit(userCtx, shardingID)
	if err != nil {
		userCtx.ErrorWF("OnDollEquipDismantleRQ GetAllDollAssembleSuit fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	// 身上的装备
	assembleGuid := make(map[int64]struct{})
	if len(assembleInfoMap) > 0 {
		for _, equipList := range assembleInfoMap {
			if len(equipList) > 0 {
				for _, v := range equipList {
					assembleGuid[v.GetEquipGuid()] = struct{}{}
				}
			}
		}
	}

	equipInfoMap, err := mazebagequipredis.GetAllEquipInfo(userCtx, shardingID)
	if err != nil {
		userCtx.ErrorWF("OnDollEquipDismantleRQ GetAllEquipInfo fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	award := make(map[int32]int64)
	equipGuid2EquipId := make(map[int64]int32)

	awardBag := make(map[int32]int64)
	// awardTempBag := make(map[int32]int64)
	guidsBag := make([]int64, 0)
	// guidsTempBag := make([]int64, 0)

	for _, equipGuid := range req.GetSelectGuid() {

		equip, ok := equipInfoMap[equipGuid]
		if !ok && req.GetDismantleFrom() == 4 {
			continue
		}
		if !ok {
			userCtx.WarnWF("OnDollEquipDismantleRQ not found equip", zap.Any("guid", equipGuid))
			res.ErrInfo = errors.NewCodeError(constdef.DE_ERR_EQUIP_DISMANTLE_EQUIP_NOT_EXISTS, toastmsgtipexcel.GetToastMsgTip(2031, "选中的部分装备已经被分解")).ToInfo()
			// res.ErrInfo = errors.NewCodeError(constdef.DE_ERR_EQUIP_DISMANTLE_EQUIP_NOT_EXISTS, "选中的部分装备已经被分解").ToInfo()
			return
		}

		if _, ok := assembleGuid[equip.GetEquipGuid()]; ok { // 身上的装备不能分解
			userCtx.ErrorWF("OnDollEquipDismantleRQ assemble equip", zap.Any("guid", equipGuid))
			res.ErrInfo = errors.NewCodeError(constdef.DE_ERR_EQUIP_DISMANTLE_QUANLITY_WRONG, "").ToInfo()
			return
		}

		userCtx.DebugWF("OnDollEquipDismantleRQ start dismantle equip", zap.Any("equipGuid", equip.GetEquipGuid()))

		// 获取装备可以出售获得的材料
		var equipAward map[int32]int64

		equipCfg := GMazeEquipInfoV8Cfg.Get(equip.GetEquipId())
		if equipCfg == nil {
			userCtx.ErrorWF("OnDollEquipDismantleRQ get equip cfg fail", zap.Any("equipId", equip.GetEquipId))
			res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
			return
		}
		equipAward, err = getEquipDismantle(userCtx, equip.GetEquipGuid(), int64(equip.GetEquipId()), equipCfg)
		if err != nil {
			userCtx.ErrorWF("OnDollEquipDismantleRQ get dismantle award fail", zap.Any("equipGuid", equip.GetEquipGuid()), zap.Error(err))
			res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
			return
		}

		for k, v := range equipAward {
			award[k] += v
			awardBag[k] += v
		}

		equipGuid2EquipId[equip.GetEquipGuid()] = equip.GetEquipId()

		// 背包装备
		guidsBag = append(guidsBag, equip.GetEquipGuid())
	}

	userCtx.DebugWF("OnDollEquipDismantleRQ end dismantle equip", zap.Any("total", award))

	if len(guidsBag) == 0 && req.GetDismantleFrom() != 4 {
		userCtx.ErrorWF("OnDollEquipDismantleRQ guids empty")
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	awardItems := itemutil.Map2Common(award)
	if req.GetDismantleFrom() != 4 {
		res.DismantleAward = awardItems
	}

	tradeNo := uniqueid.GenUniqueIdUInt64()
	// 分解装备
	rqSale := &MazeEquipSvr.SvrMazeEquipSaleRQ{
		UserId:     proto.Uint64(shardingID),
		EquipGuids: guidsBag,
		TradeNum:   proto.Uint64(tradeNo),
	}
	rsSale := &MazeEquipSvr.SvrMazeEquipSaleRS{}
	userCtx.DebugWF("OnDollEquipDismantleRQ SvrDollEquipSaleRS dump", zap.Any("rqSale", rqSale), zap.Any("rsSale", rsSale))
	// err = dollequipbagrpc.MazeEquipSaleRQ(userCtx, rqSale, rsSale)
	err = OnSvrDollEquipSaleRQ(userCtx, int64(shardingID), rqSale, rsSale)
	if err != nil {
		userCtx.ErrorWF("OnDollEquipDismantleRQ SvrDollEquipSaleRS fail", zap.Error(err), zap.Any("rq", rqSale))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if rsSale.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		if rsSale.GetErrInfo().GetErrCode() == constdef.DE_ERR_EQUIP_DISMANTLE_EQUIP_NOT_EXISTS {
			res.ErrInfo = errors.NewCodeError(constdef.DE_ERR_EQUIP_DISMANTLE_EQUIP_NOT_EXISTS, "装备不存在").ToInfo()
			return
		}
		userCtx.ErrorWF("OnDollEquipDismantleRQ SvrDollEquipSaleRS checkRs fail", zap.Any("rs", rsSale), zap.Any("rq", rqSale))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(string(rsSale.GetErrInfo().GetErrMsg()))
		return
	}

	// record := &dollequipdismantlekafka.MazeGameEquipDismantleRecord{
	// 	UserId:   shardingID,
	// 	TradeNum: tradeNo,
	// 	IsFail:   0,
	// }

	if len(awardItems) > 0 && req.GetDismantleFrom() != 4 {
		// 699	UN_CGK_COMMON_BILL_TYPE_699	迷宫分解装备
		errInfo := gentradeno.AddItemEx(userCtx, shardingID, 699, tradeNo, req.GetHeader(), awardItems...)
		if errInfo != nil {
			userCtx.ErrorWF("OnDollEquipDismantleRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("rq", rqSale))
		}
	}

	// DismantleRecordPush(userCtx, record, guidsBag, []int64{}, awardBag, map[int32]int64{}, equipGuid2EquipId)

	return
}

func getEquipDismantle(logger fklog.FKLogI, equipGuid, equipId int64, cfg *GMazeEquipInfoV8Cfg.MazeEquipInfoV8ConfigRow) (award map[int32]int64, err error) {
	award = make(map[int32]int64)

	for k, v := range cfg.Sell {
		if k > 0 && v > 0 {
			award[k] += v
		}
	}

	logger.DebugWF("getEquipDismantle one equip succ", zap.Any("equipGuid", equipGuid),
		zap.Int64("equipId", equipId), zap.Any("equip_award", award))

	return
}

func DismantleRecordPush(logger fklog.FKLogI, record *dollequipdismantlekafka.MazeGameEquipDismantleRecord,
	guidsBag, guidsTempBag []int64, awardBag, awardTempBag map[int32]int64, equipGuid2EquipId map[int64]int32) {

	equipGuidStr := make([]string, 0)

	if len(guidsBag) > 0 {
		record.Award = itemutil.CommonItemsToString(itemutil.Map2Common(awardBag))
		record.OpType = dollequipdismantlekafka.OpTypeFromBagEquipInfo
		for i := 0; i < len(guidsBag); i++ {
			equipGuidStr = append(equipGuidStr, fmt.Sprintf("%d:%d", guidsBag[i], equipGuid2EquipId[guidsBag[i]]))
			if len(equipGuidStr) >= 80 {
				record.EquipGuids = strings.Join(equipGuidStr, ",")
				dollequipdismantlekafka.PushDollEquipDismantleRecord(logger, record)
				equipGuidStr = make([]string, 0)
			}
		}

		if len(equipGuidStr) > 0 {
			record.EquipGuids = strings.Join(equipGuidStr, ",")
			dollequipdismantlekafka.PushDollEquipDismantleRecord(logger, record)
		}
	}

	if len(guidsTempBag) > 0 {
		equipGuidStr = make([]string, 0)
		record.Award = itemutil.CommonItemsToString(itemutil.Map2Common(awardTempBag))
		record.OpType = dollequipdismantlekafka.OpTypeFromTempBagEquipInfo

		for i := 0; i < len(guidsTempBag); i++ {
			equipGuidStr = append(equipGuidStr, fmt.Sprintf("%d:%d", guidsTempBag[i], equipGuid2EquipId[guidsTempBag[i]]))
			if len(equipGuidStr) >= 80 {
				record.EquipGuids = strings.Join(equipGuidStr, ",")
				dollequipdismantlekafka.PushDollEquipDismantleRecord(logger, record)
				equipGuidStr = make([]string, 0)
			}
		}
		if len(equipGuidStr) > 0 {
			record.EquipGuids = strings.Join(equipGuidStr, ",")
			dollequipdismantlekafka.PushDollEquipDismantleRecord(logger, record)
		}
	}
}
