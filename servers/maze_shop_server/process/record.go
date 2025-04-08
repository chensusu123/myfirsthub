package process

// func PushDollMazeShopInfoLog(logger fklog.FKLogI, userId uint64, level, itemType int32, mazeShopInfo *mazeshopseqredis.MazeShopInfo, addEquipList []*DollEquipSvr.DollEquipInfoSvr, opType int32, tradeNumber uint64, retCode int32) error {
// 	addEquipGuidStr := make([]string, 0)
// 	for _, equipInfo := range addEquipList {
// 		addEquipGuidStr = append(addEquipGuidStr, fmt.Sprintf("%d:%d", equipInfo.GetEquipGuid(), equipInfo.GetEquipId()))
// 	}
// 	saleMsg := &dollmazeshoprecord.DollMazeShopRecord{
// 		UserId:        userId,
// 		AreaId:        level,
// 		SeqId:         mazeShopInfo.SeqId,
// 		BackId:        mazeShopInfo.BackId,
// 		SlotId:        itemType,
// 		CurSeqIndex:   mazeShopInfo.CurSeqIndex,
// 		BackIndex:     mazeShopInfo.BackIndex,
// 		TotalCount:    mazeShopInfo.TotalCount,
// 		OpType:        opType,
// 		TradeNumber:   tradeNumber,
// 		AddEquipGuids: strings.Join(addEquipGuidStr, ","),
// 		RetCode:       retCode,
// 	}
// 	if err := dollmazeshoprecord.PushDollMazeShopRecord(logger, saleMsg); err != nil {
// 		logger.ErrorWF("PushDollMazeShopInfoLog PushDollMazeShopRecord err", zap.Error(err))
// 	}
// 	return nil
// }

// func PushMazeMoneyChgRecord(logger fklog.FKLogI, userId uint64, moneyId int32, oldCount int64, newCount int64) {
// 	record := &mazemoneykafka.MazeMoneyRecord{
// 		UserId:        userId,
// 		OldMoneyId:    moneyId,
// 		OldMoneyCount: oldCount,
// 		NewMoneyId:    moneyId,
// 		NewMoneyCount: newCount,
// 		ChgReason:     mazemoneykafka.MoneyChgReasonDeath,
// 	}

// 	mazemoneykafka.PushMazeMoneyRecord(logger, record)
// }
