package game

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"google.golang.org/protobuf/proto"
)

func OnBarrierAreaHeartRQ(logger fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	fkprometheus.InfoPMT("OnBarrierAreaHeartRQ")()

	// req := rqMsg.(*DollMazeBarrier.BarrierAreaHeartRQ)
	// res := rsMsg.(*DollMazeBarrier.BarrierAreaHeartRS)

	// logger.InfoWF("OnBarrierAreaHeartRQ start", zap.Any("req", req))
	// defer func() {
	// 	logger.InfoWF("OnBarrierAreaHeartRQ end", zap.Any("res", res))
	// }()

	// res.Header = req.Header
	// res.ErrInfo = errors.NO_ERROR
	// res.BarrierId = req.BarrierId
	// res.AreaId = req.AreaId

	// userId := shardingID

	// res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("接口废弃")
	return

	// barrierId, areaId, highArea, err := dollmazebarrier.GetMazeInfo(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("OnBarrierAreaHeartRQ GetMazeInfo fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }
	// if barrierId != req.GetBarrierId() {
	// 	logger.ErrorWF("OnBarrierAreaHeartRQ barrier not match", zap.Any("barrierId", barrierId),
	// 		zap.Any("barrierReq", req.GetBarrierId()))
	// 	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡数据不匹配")
	// 	return
	// }
	// if areaId != req.GetAreaId() {
	// 	logger.ErrorWF("OnBarrierAreaHeartRQ area not match", zap.Any("areaId", areaId),
	// 		zap.Any("areaReq", req.GetAreaId()))
	// 	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("区域数据不匹配")
	// 	return
	// }

	// cfg := GMazeBarriesV8Cfg.Get(req.GetBarrierId())
	// if cfg == nil {
	// 	logger.ErrorWF("OnBarrierAreaHeartRQ curr barrier not found cfg", zap.Any("barrierId", req.GetBarrierId()))
	// 	res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
	// 	return
	// }

	// if req.GetAreaId() != highArea {
	// 	//不是最高区域不更新生产信息
	// 	logger.WarnWF("OnBarrierAreaHeartRQ not high area", zap.Any("barrierId", req.GetBarrierId()),
	// 		zap.Any("areaId", req.GetAreaId()), zap.Any("highArea", highArea))
	// 	res.ErrInfo = errors.NewErrorInfo(ERROR_CODE_NOT_IN_HIGH_AREA, errMsg[ERROR_CODE_NOT_IN_HIGH_AREA])
	// 	return
	// }

	// moneyId := cfg.Currency_id

	// // // 心跳时间得记录 不然不知道是否心跳过期退出生产
	// // produce, err := dollmazeproduceredis.GetUserProduce(logger, userId)
	// // if err != nil {
	// // 	logger.ErrorWF("OnBarrierAreaHeartRQ GetUserProduce fail", zap.Error(err))
	// // 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// // 	return
	// // }
	// // if produce == nil {
	// // 	logger.ErrorWF("OnBarrierAreaHeartRQ GetUserProduce nil")
	// // 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// // 	return
	// // }
	// // produce.LastHeartTime = proto.Int64(time.Now().Unix())
	// // err = dollmazeproduceredis.SetUserProduce(logger, userId, produce)
	// // if err != nil {
	// // 	logger.ErrorWF("OnBarrierAreaHeartRQ SetUserProduce update heartTime fail", zap.Error(err))
	// // 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// // 	return
	// // }

	// err = StartFixMazeProduce(logger, userId, highArea)
	// if err != nil {
	// 	logger.ErrorWF("OnBarrierAreaHeartRQ StartFixMazeProduce fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	// rate, round, nextLevelForce, _, err := GetProduceRate(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("OnBarrierAreaHeartRQ GetProduceRate fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	// res.AreaInfo = &DollMazeBarrier.BarrierAreaInfo{
	// 	AreaId: proto.Int32(req.GetAreaId()),
	// 	MoneyProduce: &DollMazeBarrier.BarrierMoneyProduce{
	// 		Money:          &Common.Item{ItemId: proto.Int32(moneyId)},
	// 		Rate:           proto.Int64(rate),
	// 		Round:          proto.Int64(round),
	// 		NextLevelForce: proto.Int64(nextLevelForce),
	// 	},
	// }

	return nil
}

func GetProduceRate(logger fklog.FKLogI, userId uint64) (rate int64, round int64, nextLevelForce int64, limitVal int64, err error) {
	// forceVal, err := dollforceredis.GetShowDollForce(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("GetProduceRate GetShowDollForce fail")
	// 	return
	// }

	// allProduceCfg := GMazeKongfuExtraCoinV8Cfg.GetAll()

	// for _, cfg := range allProduceCfg {
	// 	if forceVal <= cfg.Kongfu_max && forceVal >= cfg.Kongfu_min {
	// 		// rate = cfg.Cyecle_add_coin
	// 		rate = cfg.Cyecle_add[1]
	// 		round = int64(cfg.Cyecle_time) / 1000
	// 		limitVal = cfg.Short_time[1]
	// 		nextLevelForce = cfg.Kongfu_max + 1
	// 		break
	// 	}
	// }

	// if limitVal == 0 {
	// 	err = errors.New("maze kongfu coin empty")
	// }
	return
}
