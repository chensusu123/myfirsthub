package game

// func PushDollMazeShopInfoLog(ctx context.Context, userId uint64, areaId int32, mazeShopInfo *dollmazeshopseqredis.MazeShopInfo, addEquipList []*MazeEquipSvr.MazeEquipInfoSvr, opType int32, tradeNumber uint64, retCode int32) error {
// 	addEquipGuidStr := make([]string, 0)
// 	for _, equipInfo := range addEquipList {
// 		addEquipGuidStr = append(addEquipGuidStr, fmt.Sprintf("%d:%d", equipInfo.GetEquipGuid(), equipInfo.GetEquipId()))
// 	}
// 	saleMsg := &dollmazeshoprecord.DollMazeShopRecord{
// 		UserId:        userId,
// 		AreaId:        areaId,
// 		SeqId:         mazeShopInfo.SeqId,
// 		BackId:        mazeShopInfo.BackId,
// 		CurSeqIndex:   mazeShopInfo.CurSeqIndex,
// 		BackIndex:     mazeShopInfo.BackIndex,
// 		TotalCount:    mazeShopInfo.TotalCount,
// 		OpType:        opType,
// 		TradeNumber:   tradeNumber,
// 		AddEquipGuids: strings.Join(addEquipGuidStr, ","),
// 		RetCode:       retCode,
// 	}
// 	if err := dollmazeshoprecord.PushDollMazeShopRecord(logger, saleMsg); err != nil {
// 		logger.CtxError(ctx,"PushDollMazeShopInfoLog PushDollMazeShopRecord err", zap.Error(err))
// 	}
// 	return nil
// }
