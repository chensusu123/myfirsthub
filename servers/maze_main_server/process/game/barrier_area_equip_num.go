package game

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"google.golang.org/protobuf/proto"
)

func OnMazeAreaEquipNumRQ(logger fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	fkprometheus.InfoPMT("OnMazeAreaEquipNumRQ")()

	// req := rqMsg.(*DollMazeBarrier.MazeAreaEquipNumRQ)
	// res := rsMsg.(*DollMazeBarrier.MazeAreaEquipNumRS)

	// logger.CtxInfo(ctx,"OnMazeAreaEquipNumRQ start", zap.Any("req", req))
	// defer func() {
	// 	logger.CtxInfo(ctx,"OnMazeAreaEquipNumRQ end", zap.Any("res", res))
	// }()

	// res.Header = req.Header
	// res.ErrInfo = errors.NO_ERROR
	// res.BarrierId = req.BarrierId
	// res.AreaId = req.AreaId
	// res.CollectEquipNum = req.CollectEquipNum

	// userId := shardingID

	// areaCfg := GMazeBrushAreaV8Cfg.GetWithCtx(ctx,req.GetAreaId())
	// if areaCfg == nil {
	// 	logger.CtxError(ctx,"OnMazeAreaEquipNumRQ get maze brush area cfg nil", zap.Any("areaId", req.GetAreaId()))
	// 	res.ErrInfo = errors.CONFIG_NOT_FOUND.Wrap("区域id未取到配表数据")
	// 	return
	// }
	// if areaCfg.Barries != req.GetBarrierId() {
	// 	logger.CtxError(ctx,"OnMazeAreaEquipNumRQ barrier and area not match", zap.Any("barrierId", req.GetBarrierId()), zap.Any("areaId", req.GetAreaId()))
	// 	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡与区域不匹配")
	// 	return
	// }

	// maxCfg := GMazeConfigV8Cfg.GetWithCtx(ctx,1)
	// if maxCfg == nil {
	// 	logger.CtxError(ctx,"OnMazeAreaEquipNumRQ get maze maxEquip cfg nil")
	// 	res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
	// 	return
	// }

	// maxVal := maxCfg.Value_map[req.GetAreaId()]
	// if maxVal == -1 {
	// 	//-1 表示不限制
	// 	return
	// }

	// userBarrier, err := mazeuserbarrier.GetUserBarrier(logger, userId, req.GetBarrierId())
	// if err != nil {
	// 	logger.CtxError(ctx,"OnMazeAreaEquipNumRQ GetUserBarrier fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	// areaInfo, ok := userBarrier.AreaInfo[req.GetAreaId()]
	// if !ok {
	// 	areaInfo = &mazeuserbarrier.BarrierArea{
	// 		AreaId:          req.GetAreaId(),
	// 		DefeatedFoe:     make(map[int32]int32),
	// 		CollectEquipNum: 0,
	// 	}
	// 	userBarrier.AreaInfo[req.GetAreaId()] = areaInfo
	// }

	// oldEquipNum := userBarrier.AreaInfo[req.GetAreaId()].CollectEquipNum
	// userBarrier.AreaInfo[req.GetAreaId()].CollectEquipNum = req.GetCollectEquipNum()

	// err = mazeuserbarrier.SetUserBarrier(logger, userId, req.GetBarrierId(), userBarrier, nil)
	// if err != nil {
	// 	logger.CtxError(ctx,"OnMazeAreaEquipNumRQ SetUserBarrierArea fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	// if oldEquipNum < int32(maxVal) && req.GetCollectEquipNum() >= int32(maxVal) {
	// 	err = StartFixMazeProduce(logger, userId, req.GetAreaId())
	// 	if err != nil {
	// 		logger.CtxError(ctx,"OnMazeAreaEquipNumRQ StartFixMazeProduce fail", zap.Error(err))
	// 		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 		return
	// 	}
	// }

	// // if oldEquipNum >= int32(maxVal) && req.GetCollectEquipNum() < int32(maxVal) {
	// // 	err = StartFixMazeProduce(logger, userId, req.GetAreaId())
	// // 	if err != nil {
	// // 		logger.CtxError(ctx,"OnMazeAreaEquipNumRQ StartFixMazeProduce fail", zap.Error(err))
	// // 		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// // 		return
	// // 	}
	// // }

	return
}
