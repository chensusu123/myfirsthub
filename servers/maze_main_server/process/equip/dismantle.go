package equip

import (
	"context"
	"fmt"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/tradeno"
	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/excel/toastmsgtipexcel"
	"maze_game_server/io/kafka/dollequipdismantlekafka"
	"maze_game_server/io/redis/dollassemblesuitredis"
	"maze_game_server/io/redis/mazebagequipredis"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/services/itemservice"
	"strings"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (e *Equip) OnDollEquipDismantleRQ_10410_10411(s *session.Session, req *MazeGameEquip.MazeEquipDismantleRQ) (err error) {
	defer fkprometheus.DebugPMT("OnDollEquipDismantleRQ")()

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeGameEquip.MazeEquipDismantleRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.SelectQuality = req.SelectQuality
	res.SelectGuid = req.SelectGuid
	res.DismantleFrom = req.DismantleFrom
	// res.DismantleAward = req.DismantleAward

	userId := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnDollEquipDismantleRQ end", zap.Any("res", res))
	}()

	logger.CtxInfo(ctx, "OnDollEquipDismantleRQ with", zap.Any("req", req))

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
	assembleInfoMap, err := dollassemblesuitredis.GetAllDollAssembleSuit(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnDollEquipDismantleRQ GetAllDollAssembleSuit fail", zap.Error(err))
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

	equipInfoMap, err := mazebagequipredis.GetAllEquipInfo(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnDollEquipDismantleRQ GetAllEquipInfo fail", zap.Error(err))
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
			logger.CtxError(ctx, "OnDollEquipDismantleRQ not found equip", zap.Any("guid", equipGuid))
			res.ErrInfo = errors.NewCodeError(constdef.DE_ERR_EQUIP_DISMANTLE_EQUIP_NOT_EXISTS, toastmsgtipexcel.GetToastMsgTip(2031, "选中的部分装备已经被分解")).ToInfo()
			// res.ErrInfo = errors.NewCodeError(constdef.DE_ERR_EQUIP_DISMANTLE_EQUIP_NOT_EXISTS, "选中的部分装备已经被分解").ToInfo()
			return
		}

		if _, ok := assembleGuid[equip.GetEquipGuid()]; ok { // 身上的装备不能分解
			logger.CtxError(ctx, "OnDollEquipDismantleRQ assemble equip", zap.Any("guid", equipGuid))
			res.ErrInfo = errors.NewCodeError(constdef.DE_ERR_EQUIP_DISMANTLE_QUANLITY_WRONG, "").ToInfo()
			return
		}

		logger.CtxInfo(ctx, "OnDollEquipDismantleRQ start dismantle equip", zap.Any("equipGuid", equip.GetEquipGuid()))

		// 获取装备可以出售获得的材料
		var equipAward map[int32]int64

		equipCfg := GMazeEquipInfoV8Cfg.GetWithCtx(ctx, equip.GetEquipId())
		if equipCfg == nil {
			logger.CtxError(ctx, "OnDollEquipDismantleRQ get equip cfg fail", zap.Any("equipId", equip.GetEquipId))
			res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
			return
		}
		equipAward, err = getEquipDismantle(ctx, equip.GetEquipGuid(), int64(equip.GetEquipId()), equipCfg)
		if err != nil {
			logger.CtxError(ctx, "OnDollEquipDismantleRQ get dismantle award fail", zap.Any("equipGuid", equip.GetEquipGuid()), zap.Error(err))
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

	logger.CtxInfo(ctx, "OnDollEquipDismantleRQ end dismantle equip", zap.Any("total", award))

	if len(guidsBag) == 0 && req.GetDismantleFrom() != 4 {
		logger.CtxError(ctx, "OnDollEquipDismantleRQ guids empty")
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	awardItems := itemutil.Map2Common(award)
	if req.GetDismantleFrom() != 4 {
		res.DismantleAward = awardItems
	}

	tradeNo := tradeno.GetTradeNum()
	// 分解装备
	rqSale := &MazeEquipSvr.SvrMazeEquipSaleRQ{
		UserId:     proto.Uint64(userId),
		EquipGuids: guidsBag,
		TradeNum:   proto.Uint64(tradeNo),
	}
	rsSale := &MazeEquipSvr.SvrMazeEquipSaleRS{}
	logger.CtxError(ctx, "OnDollEquipDismantleRQ SvrDollEquipSaleRS dump", zap.Any("rqSale", rqSale), zap.Any("rsSale", rsSale))
	// err = dollequipbagrpc.MazeEquipSaleRQ(logger, rqSale, rsSale)
	err = OnSvrDollEquipSaleRQ(ctx, int64(userId), rqSale, rsSale)
	if err != nil {
		logger.CtxError(ctx, "OnDollEquipDismantleRQ SvrDollEquipSaleRS fail", zap.Error(err), zap.Any("rq", rqSale))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if rsSale.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		if rsSale.GetErrInfo().GetErrCode() == constdef.DE_ERR_EQUIP_DISMANTLE_EQUIP_NOT_EXISTS {
			res.ErrInfo = errors.NewCodeError(constdef.DE_ERR_EQUIP_DISMANTLE_EQUIP_NOT_EXISTS, "装备不存在").ToInfo()
			return
		}
		logger.CtxError(ctx, "OnDollEquipDismantleRQ SvrDollEquipSaleRS checkRs fail", zap.Any("rs", rsSale), zap.Any("rq", rqSale))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(string(rsSale.GetErrInfo().GetErrMsg()))
		return
	}

	record := &dollequipdismantlekafka.MazeGameEquipDismantleRecord{
		UserId:   userId,
		TradeNum: tradeNo,
		IsFail:   0,
	}

	if len(awardItems) > 0 && req.GetDismantleFrom() != 4 {
		// 699	UN_CGK_COMMON_BILL_TYPE_699	迷宫分解装备
		items := itemutil.Map2ItemInfo(award)
		errInfo := itemservice.GlobalItemService.AddItem(ctx, userId, itemservice.ItemOpTypeDismantle, tradeNo, items...)
		if errInfo != nil {
			logger.CtxError(ctx, "OnDollEquipDismantleRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("rq", rqSale))
		}
	}

	DismantleRecordPush(ctx, record, guidsBag, []int64{}, awardBag, map[int32]int64{}, equipGuid2EquipId)

	return
}

func getEquipDismantle(ctx context.Context, equipGuid, equipId int64, cfg *GMazeEquipInfoV8Cfg.MazeEquipInfoV8ConfigRow) (award map[int32]int64, err error) {
	award = make(map[int32]int64)
	logger := fklog.ContextAppLogger(ctx)
	for k, v := range cfg.Sell {
		if k > 0 && v > 0 {
			award[k] += v
		}
	}

	logger.CtxDebug(ctx, "getEquipDismantle one equip succ", zap.Any("equipGuid", equipGuid),
		zap.Int64("equipId", equipId), zap.Any("equip_award", award))

	return
}

func DismantleRecordPush(ctx context.Context, record *dollequipdismantlekafka.MazeGameEquipDismantleRecord,
	guidsBag, guidsTempBag []int64, awardBag, awardTempBag map[int32]int64, equipGuid2EquipId map[int64]int32) {

	equipGuidStr := make([]string, 0)

	if len(guidsBag) > 0 {
		record.Award = itemutil.CommonItemsToString(itemutil.Map2Common(awardBag))
		record.OpType = dollequipdismantlekafka.OpTypeFromBagEquipInfo
		for i := 0; i < len(guidsBag); i++ {
			equipGuidStr = append(equipGuidStr, fmt.Sprintf("%d:%d", guidsBag[i], equipGuid2EquipId[guidsBag[i]]))
			if len(equipGuidStr) >= 80 {
				record.EquipGuids = strings.Join(equipGuidStr, ",")
				dollequipdismantlekafka.PushDollEquipDismantleRecord(ctx, record)
				equipGuidStr = make([]string, 0)
			}
		}

		if len(equipGuidStr) > 0 {
			record.EquipGuids = strings.Join(equipGuidStr, ",")
			dollequipdismantlekafka.PushDollEquipDismantleRecord(ctx, record)
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
				dollequipdismantlekafka.PushDollEquipDismantleRecord(ctx, record)
				equipGuidStr = make([]string, 0)
			}
		}
		if len(equipGuidStr) > 0 {
			record.EquipGuids = strings.Join(equipGuidStr, ",")
			dollequipdismantlekafka.PushDollEquipDismantleRecord(ctx, record)
		}
	}
}
