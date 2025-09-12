package mazeuserbarrier

import "context"

type MazeUserBarrier struct {
	BarrierId   int32
	AreaInfo    map[int32]*BarrierArea
	BoxState    map[int32]int32
	EquipPoints int64
}

type BarrierArea struct {
	AreaId      int32
	DefeatedFoe map[int32]int32
	// FoePoints       int64
	CollectEquipNum int32
}

// func GetUserBarrier(ctx context.Context, userId uint64, barrierId int32) (userInfo *MazeUserBarrier, err error) {
// 	areaList, boxInfo, points, err := mazeuserbarrierredis.GetBarrierInfo(logger, userId, barrierId)
// 	if err != nil {
// 		logger.CtxError(ctx,"GetUserBarrier GetBarrierInfo fail", zap.Error(err), zap.Any("barrierId", barrierId))
// 		return
// 	}
// 	userInfo = &MazeUserBarrier{
// 		BarrierId:   barrierId,
// 		AreaInfo:    make(map[int32]*BarrierArea),
// 		BoxState:    make(map[int32]int32),
// 		EquipPoints: points,
// 	}

// 	for _, v := range boxInfo.GetBoxList() {
// 		userInfo.BoxState[v.GetBoxId()] = v.GetBoxState()
// 	}

// 	for _, v := range areaList {
// 		userInfo.AreaInfo[v.GetAreaId()] = &BarrierArea{
// 			AreaId:      v.GetAreaId(),
// 			DefeatedFoe: make(map[int32]int32),
// 			// FoePoints:       0,
// 			CollectEquipNum: 0,
// 		}

// 		for _, foe := range v.GetDefeatedFoe() {
// 			userInfo.AreaInfo[v.GetAreaId()].DefeatedFoe[foe.GetFeoId()] = foe.GetFeoCount()
// 		}

// 		// userInfo.AreaInfo[v.GetAreaId()].FoePoints = v.GetFoePoints()
// 		userInfo.AreaInfo[v.GetAreaId()].CollectEquipNum = v.GetCollectEquipNum()
// 	}

// 	return
// }

func SetUserBarrierBox(ctx context.Context, userId uint64, barrierId int32, userInfo *MazeUserBarrier) (err error) {

	// boxInfo := &DollMazeBarrierCache.DollMazeBoxDb{
	// 	BoxList: make([]*DollMazeBarrierCache.MazeBox, 0, len(userInfo.BoxState)),
	// }

	// for k, v := range userInfo.BoxState {
	// 	boxInfo.BoxList = append(boxInfo.BoxList, &DollMazeBarrierCache.MazeBox{
	// 		BoxId:    proto.Int32(k),
	// 		BoxState: proto.Int32(v),
	// 	})
	// }

	// err = mazeuserbarrierredis.SetBox(logger, userId, barrierId, boxInfo)
	// if err != nil {
	// 	logger.CtxError(ctx,"SetUserBarrierBox SetBox fail", zap.Error(err),
	// 		zap.Any("barrierId", barrierId), zap.Any("boxInfo", boxInfo))
	// 	return
	// }

	return
}

func SetUserBarrier(ctx context.Context, userId uint64, barrierId int32, userInfo *MazeUserBarrier, reportPb []byte) (err error) {
	// areaInfo := make([]*DollMazeBarrierCache.DollMazeAreaDb, 0)

	// for _, v := range userInfo.AreaInfo {
	// 	pb := &DollMazeBarrierCache.DollMazeAreaDb{
	// 		AreaId:      proto.Int32(v.AreaId),
	// 		DefeatedFoe: make([]*DollMazeBarrierCache.DollFoeInfoDb, 0),
	// 		// FoePoints:       proto.Int64(v.FoePoints),
	// 		CollectEquipNum: proto.Int32(v.CollectEquipNum),
	// 	}
	// 	for foe, count := range v.DefeatedFoe {
	// 		pb.DefeatedFoe = append(pb.DefeatedFoe, &DollMazeBarrierCache.DollFoeInfoDb{
	// 			FeoId:    proto.Int32(foe),
	// 			FeoCount: proto.Int32(count),
	// 		})
	// 	}

	// 	areaInfo = append(areaInfo, pb)
	// }

	// err = mazeuserbarrierredis.SetUserBarrier(logger, userId, barrierId, areaInfo, userInfo.EquipPoints, reportPb)
	// if err != nil {
	// 	logger.CtxError(ctx,"SetUserBarrierArea SetUserBarrier fail", zap.Error(err),
	// 		zap.Any("barrierId", barrierId), zap.Any("areaInfo", areaInfo))
	// 	return
	// }
	return
}
