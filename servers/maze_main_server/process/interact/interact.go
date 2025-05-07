package interact

import (
	"fmt"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/cache/simCache"

	"gitlab.ifreetalk.com/maze/maze_game_server/excel/equipmixcostcfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeuserlevelredis"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipMixListV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkserver/tcp_service"

	"gitlab.ifreetalk.com/plate/freetk/fkutil/saferand"
	"gitlab.ifreetalk.com/plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/plate/protodef/MazeEquipMix"
	"gitlab.ifreetalk.com/plate/protodef/MazeEquipSvr"
	"gitlab.ifreetalk.com/plate/protodef/MessageType"
	"go.uber.org/zap"
)

func RegTcpHandler() {

	// 迷宫装备合成消耗
	tcp_service.RegProcSimple(16265, &MazeEquipMix.MazeEquipMixCostRQ{},
		16266, &MazeEquipMix.MazeEquipMixCostRS{}, OnMazeEquipMixCostRQ)

	// 迷宫装备合成
	tcp_service.RegProcSimple(16267, &MazeEquipMix.MazeEquipMixRQ{},
		16268, &MazeEquipMix.MazeEquipMixRS{}, OnMazeEquipMixRQ)

}

func OnMazeEquipMixCostRQ(ctx fknet.TCPContext, uid uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req := rqMsg.(*MazeEquipMix.MazeEquipMixCostRQ)
	res := rsMsg.(*MazeEquipMix.MazeEquipMixCostRS)
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	defer fkprometheus.DebugPMT("OnMazeEquipMixCostRQ")()
	defer func() {
		ctx.InfoWF("OnMazeEquipMixCostRQ end",
			zap.Any("req", req),
			zap.Any("res", res),
		)
	}()

	// 查等级
	lv, err := mazeuserlevelredis.GetUserLevel(ctx, uid)
	if err != nil {
		ctx.ErrorWF("OnMazeEquipMixCostRQ get user level error", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	// 读表返回
	cfg := equipmixcostcfg.GetEquipCost(int32(lv))
	if cfg == nil || len(cfg.Cost) == 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("配表错误")
		return
	}

	res.Cost = cfg.Cost

	return
}

var equipMixCache = simCache.NewCache()

func OnMazeEquipMixRQ(ctx fknet.TCPContext, uid uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req := rqMsg.(*MazeEquipMix.MazeEquipMixRQ)
	res := rsMsg.(*MazeEquipMix.MazeEquipMixRS)
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	defer fkprometheus.DebugPMT("OnMazeEquipMixRQ")()
	defer func() {
		ctx.InfoWF("OnMazeEquipMixRQ end",
			zap.Any("req", req),
			zap.Any("res", res),
		)
	}()

	// 限制重复请求
	if equipMixCache.SetCache(uid, struct{}{}) {
		res.ErrInfo = errors.NewCodeError(62091, "").ToInfo()
		return
	}
	defer equipMixCache.DelCache(uid)

	// 查等级
	lv, err := mazeuserlevelredis.GetUserLevel(ctx, uid)
	if err != nil {
		ctx.ErrorWF("OnMazeEquipMixRQ get user level error", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	// 读表 取消耗
	cfg := equipmixcostcfg.GetEquipCost(int32(lv))
	if cfg == nil || cfg.Cfg == nil || len(cfg.Cost) == 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("配表错误")
		return
	}

	// 校验参数消耗
	if !checkReqCost(req.GetCost(), cfg.Cost) {
		ctx.WarnWF("OnMazeEquipMixRQ cost check invalid", zap.Any("reqCost", req.GetCost()), zap.Any("cfgCost", cfg.Cost))
		res.ErrInfo = errors.COST_NOT_MATCH.ToInfo()
		return
	}

	// 读取合成信息. 取上次等级、配置、索引。 有变化重新随
	lastData, err := mazeequipmixdb.GetEquipMixData(ctx, uid)
	if err != nil {
		ctx.ErrorWF("OnMazeEquipMixRQ get user equip mix data error", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if lastData == nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	// 检查合成信息
	useOld, desc, listCfg, err := checkEquipData(ctx, lastData, int32(lv))
	if err != nil {
		ctx.ErrorWF("OnMazeEquipMixRQ check equip data err", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("配表错误")
		return
	}

	if useOld {
		ctx.DebugWF("OnMazeEquipMixRQ equip data use old",
			zap.Any("lastData", lastData),
		)
	} else {
		newData := &equipmix.MixData{}

		// 重新随
		newData.Lv = int32(lv)
		newData.Cfg = cfg.Cfg.Result_equip_list[saferand.Intn(len(cfg.Cfg.Result_equip_list))]
		newData.Idx = 0

		// 使用配置 取装备
		listCfg = GMazeEquipMixListV8Cfg.Get(newData.Cfg)
		if listCfg == nil || len(listCfg.Equip_id) == 0 {
			ctx.ErrorWF("OnMazeEquipMixRQ GMazeEquipMixListV8Cfg nil", zap.Int32("id", newData.Cfg))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("配表错误")
			return
		}

		ctx.DebugWF("OnMazeEquipMixRQ equip data use new",
			zap.Bool("useOld", useOld),
			zap.String("desc", desc),
			zap.Any("lastData", lastData),
			zap.Any("newData", newData),
		)

		lastData = newData
	}

	// 找装备
	if lastData.Idx >= len(listCfg.Equip_id) {
		ctx.ErrorWF("OnMazeEquipMixRQ idx more than equip len",
			zap.Any("listCfg", listCfg),
			zap.Any("lastData", lastData),
		)
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("配表错误")
		return
	}

	equip := listCfg.Equip_id[lastData.Idx]

	// 取完配置 idx++
	lastData.Idx++

	// 保存 装备合成配置存储
	err = mazeequipmixdb.SetEquipMixData(ctx, uid, lastData)
	if err != nil {
		ctx.ErrorWF("OnMazeEquipMixRQ set user equip mix data error", zap.Error(err))
		res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
		return
	}

	tradeNo := tradeno.GetTradeNum()

	// 扣物品
	if len(cfg.Cost) != 0 {
		cost := cfg.Cost

		var errorInfo *MessageType.ErrorInfo
		errorInfo, err = itemrpc.DeductItems(ctx, uid, 695, tradeNo, cost...)
		if err != nil {
			ctx.ErrorWF("OnMazeEquipMixRQ DeductItems err", zap.Uint64("tradeNo", tradeNo), zap.Any("cost", cost),
				zap.Any("errorInfo", errorInfo), zap.Error(err),
			)
			res.ErrInfo = errors.NewCommonCodeError("sub item err")
			return
		}
		if errorInfo != nil {
			ctx.WarnWF("OnMazeEquipMixRQ DeductItems invalid", zap.Uint64("tradeNo", tradeNo), zap.Any("cost", cost),
				zap.Any("errorInfo", errorInfo), zap.Error(err),
			)
			res.ErrInfo = errorInfo
			return
		}
	}

	// 加装备
	equipReq := &MazeEquipSvr.SvrAddMazeEquipRQ{
		UserId:      proto.Uint64(uid),
		OpType:      proto.Int32(8),
		TradeNumber: proto.Uint64(tradeNo),
	}

	equipReq.EquipList = append(equipReq.EquipList, &MazeEquipSvr.SvrEquipInfo{
		EquipId: proto.Int32(equip),
	})

	equipRes := &MazeEquipSvr.SvrAddMazeEquipRS{}

	errInfo, err := equiprpc.AddEquipRQ(ctx, equipReq, equipRes)
	if err != nil {
		ctx.ErrorWF("OnMazeEquipMixRQ add equip err",
			zap.Uint64("tradeNo", tradeNo),
			zap.Error(err),
		)
		res.ErrInfo = errors.NewCommonCodeError("add equip err")
		return
	}
	if errInfo != nil {
		ctx.WarnWF("OnMazeEquipMixRQ add equip err",
			zap.Uint64("tradeNo", tradeNo),
			zap.Any("errorInfo", errInfo),
		)
		res.ErrInfo = errInfo
		return
	}

	res.Equip = []*MazeCommon.MazeItem{
		{
			ItemId: proto.Int32(equip),
			Count:  proto.Int64(1),
		},
	}
	return
}

func checkReqCost(cost []*MazeCommon.MazeItem, cfgCost []*MazeCommon.MazeItem) bool {
	if len(cost) != len(cfgCost) {
		return false
	}

	costMap := make(map[int32]int64)
	for _, item := range cost {
		if item.GetItemId() <= 0 || item.GetCount() <= 0 {
			continue
		}
		costMap[item.GetItemId()] += item.GetCount()
	}

	for _, item := range cfgCost {
		if item.GetItemId() <= 0 || item.GetCount() <= 0 {
			continue
		}

		if item.GetCount() != costMap[item.GetItemId()] {
			return false
		}
	}

	return true
}

func checkEquipData(logger fklog.FKLogI, lastData *equipmix.MixData,
	lv int32) (useOld bool, desc string, listCft *GMazeEquipMixListV8Cfg.MazeEquipMixListV8ConfigRow, err error) {
	/*
		配置>0 	校验等级是否变化、配置、索引是否有效。
												有效 直接用
												无效 重新随
		配置<=0	重新随
	*/

	if lastData == nil {
		desc = "last nil"
		return
	}

	if lastData.Cfg <= 0 {
		desc = "last cfg <=0"
		return
	}

	if lastData.Lv != lv {
		desc = fmt.Sprintf("lv change last=%d new=%d", lastData.Lv, lv)
		return
	}

	listCft = GMazeEquipMixListV8Cfg.Get(lastData.Cfg)
	if listCft == nil {
		desc = fmt.Sprintf("last cfg nil, id=%d", lastData.Cfg)
		err = errors.New(desc)
		return
	}

	if lastData.Idx >= len(listCft.Equip_id) {
		desc = fmt.Sprintf("last idx over flow, id=%d lastIdx=%d cfgLen=%d", lastData.Cfg, lastData.Idx, len(listCft.Equip_id))
		return
	}

	useOld = true
	return
}
