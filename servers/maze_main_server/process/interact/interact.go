package interact

import (
	"fmt"

	"maze_game_server/common/cache/simCache"
	"maze_game_server/common/equipmix"
	"maze_game_server/common/errors"
	"maze_game_server/common/tradeno"
	"maze_game_server/config/GMazeEquipMixListV8Cfg"
	"maze_game_server/excel/equipmixcostcfg"
	"maze_game_server/io/redis/mazeequipmixdb"
	"maze_game_server/io/redis/mazeuserlevelredis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeEquipMix"
	"maze_game_server/pb/common/MessageType"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/pb/server/MazeItemSvr"
	equiprpc "maze_game_server/servers/maze_main_server/process/equip"
	itemrpc "maze_game_server/servers/maze_main_server/process/item"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil/saferand"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Interact struct {
	component.Base
}

func NewInteract() *Interact {
	return &Interact{}
}

func RegTcpHandler() {

	// // 迷宫装备合成消耗
	// websocket_service.RegProcSimple(10441, &MazeEquipMix.MazeEquipMixCostRQ{},
	// 	10442, &MazeEquipMix.MazeEquipMixCostRS{}, OnMazeEquipMixCostRQ)

	// // 迷宫装备合成
	// websocket_service.RegProcSimple(10443, &MazeEquipMix.MazeEquipMixRQ{},
	// 	10444, &MazeEquipMix.MazeEquipMixRS{}, OnMazeEquipMixRQ)

}

func (*Interact) OnMazeEquipMixCostRQ_10441_10442(s *session.Session, req *MazeEquipMix.MazeEquipMixCostRQ) (err error) {

	logger := log.Clone("Interact", uint64(s.UID()), 0)
	res := &MazeEquipMix.MazeEquipMixCostRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	defer fkprometheus.DebugPMT("OnMazeEquipMixCostRQ")()
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeEquipMixCostRQ end",
			zap.Any("req", req),
			zap.Any("res", res),
		)
	}()

	uid := uint64(s.UID())

	// 查等级
	lv, err := mazeuserlevelredis.GetUserLevel(logger, uid)
	if err != nil {
		logger.ErrorWF("OnMazeEquipMixCostRQ get user level error", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	// 读表返回
	cfg := equipmixcostcfg.GetEquipCost(int32(lv))
	if cfg == nil || len(cfg.Cost) == 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("配表错误")
		return
	}

	// todo： temp fix pb error
	// res.Cost = cfg.Cost

	return
}

var equipMixCache = simCache.NewCache()

func (*Interact) OnMazeEquipMixRQ_10443_10444(s *session.Session, req *MazeEquipMix.MazeEquipMixRQ) (err error) {

	logger := log.Clone("Interact", uint64(s.UID()), 0)
	res := &MazeEquipMix.MazeEquipMixRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	defer fkprometheus.DebugPMT("OnMazeEquipMixRQ")()
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeEquipMixRQ end",
			zap.Any("req", req),
			zap.Any("res", res),
		)
	}()

	uid := uint64(s.UID())

	// 限制重复请求
	if equipMixCache.SetCache(uid, struct{}{}) {
		res.ErrInfo = errors.NewCodeError(62091, "").ToInfo()
		return
	}
	defer equipMixCache.DelCache(uid)

	// 查等级
	lv, err := mazeuserlevelredis.GetUserLevel(logger, uid)
	if err != nil {
		logger.ErrorWF("OnMazeEquipMixRQ get user level error", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	// 读表 取消耗
	cfg := equipmixcostcfg.GetEquipCost(int32(lv))
	if cfg == nil || cfg.Cfg == nil || len(cfg.Cost) == 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("配表错误")
		return
	}

	// 校验参数消耗 //todo： temp fix pb error
	// if !checkReqCost(req.GetCost(), cfg.Cost) {
	//	ctx.WarnWF("OnMazeEquipMixRQ cost check invalid", zap.Any("reqCost", req.GetCost()), zap.Any("cfgCost", cfg.Cost))
	//	res.ErrInfo = errors.COST_NOT_MATCH.ToInfo()
	//	return
	// }

	// 读取合成信息. 取上次等级、配置、索引。 有变化重新随
	lastData, err := mazeequipmixdb.GetEquipMixData(logger, uid)
	if err != nil {
		logger.ErrorWF("OnMazeEquipMixRQ get user equip mix data error", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if lastData == nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	// 检查合成信息
	useOld, desc, listCfg, err := checkEquipData(logger, lastData, int32(lv))
	if err != nil {
		logger.ErrorWF("OnMazeEquipMixRQ check equip data err", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("配表错误")
		return
	}

	if useOld {
		logger.DebugWF("OnMazeEquipMixRQ equip data use old",
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
			logger.ErrorWF("OnMazeEquipMixRQ GMazeEquipMixListV8Cfg nil", zap.Int32("id", newData.Cfg))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("配表错误")
			return
		}

		logger.DebugWF("OnMazeEquipMixRQ equip data use new",
			zap.Bool("useOld", useOld),
			zap.String("desc", desc),
			zap.Any("lastData", lastData),
			zap.Any("newData", newData),
		)

		lastData = newData
	}

	// 找装备
	if lastData.Idx >= len(listCfg.Equip_id) {
		logger.ErrorWF("OnMazeEquipMixRQ idx more than equip len",
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
	err = mazeequipmixdb.SetEquipMixData(logger, uid, lastData)
	if err != nil {
		logger.ErrorWF("OnMazeEquipMixRQ set user equip mix data error", zap.Error(err))
		res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
		return
	}

	tradeNo := tradeno.GetTradeNum()

	// 扣物品
	if len(cfg.Cost) != 0 {
		cost := cfg.Cost

		var errorInfo *MessageType.ErrorInfo

		rpcreq := &MazeItemSvr.ConsumeItemRQ{
			UserId:      proto.Uint64(uid),
			Items:       cost,
			OpType:      proto.Int32(695),
			TradeNumber: proto.Uint64(tradeNo),
		}
		rpcres := &MazeItemSvr.ConsumeItemRS{}

		err = itemrpc.OnAddItemRQ(logger, rpcreq, rpcres)
		if err != nil {
			logger.ErrorWF("OnMazeEquipMixRQ DeductItems err", zap.Uint64("tradeNo", tradeNo), zap.Any("cost", cost),
				zap.Any("errorInfo", errorInfo), zap.Error(err),
			)
			res.ErrInfo = errors.NewCommonCodeError("sub item err")
			return
		}
		if rpcres.ErrInfo != nil {
			logger.WarnWF("OnMazeEquipMixRQ DeductItems invalid", zap.Uint64("tradeNo", tradeNo), zap.Any("cost", cost),
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
	err = equiprpc.OnSvrAddMazeEquipRQ(logger, int64(uid), equipReq, equipRes, "")
	if err != nil {
		logger.ErrorWF("OnMazeEquipMixRQ add equip err",
			zap.Uint64("tradeNo", tradeNo),
			zap.Error(err),
		)
		res.ErrInfo = errors.NewCommonCodeError("add equip err")
		return
	}
	// if errInfo != nil {
	//	ctx.WarnWF("OnMazeEquipMixRQ add equip err",
	//		zap.Uint64("tradeNo", tradeNo),
	//		zap.Any("errorInfo", errInfo),
	//	)
	//	res.ErrInfo = errInfo
	//	return
	// }

	// todo： temp fix pb error
	// res.Equip = []*MazeCommon.MazeItem{
	//	{
	//		ItemId: proto.Int32(equip),
	//		Count:  proto.Int64(1),
	//	},
	// }
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
