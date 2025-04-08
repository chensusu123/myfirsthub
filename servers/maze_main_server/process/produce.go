package process

// // 进出区域、心跳、打怪等生产变化掉用
// func StopMazeProduce(logger fklog.FKLogI, userId uint64) (err error) {
// 	produce, err := mazeproduceredis.GetUserProduce(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("StopMazeProduce GetUserProduce fail", zap.Error(err))
// 		return
// 	}

// 	if produce == nil {
// 		return
// 	}

// 	cfgMap, err := GetCurrProduceCfg(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("StopMazeProduce GetCurrProduceCfg fail", zap.Error(err))
// 		return
// 	}

// 	now := time.Now().Unix()
// 	// if produce.GetLastSettleTime() != now {
// 	// 校验是否超时, 按最大心跳数结算生产数量
// 	timeGap := now - produce.GetLastSettleTime()
// 	if timeGap > HeartTime {
// 		timeGap = HeartTime
// 	}

// 	cycleRoundNum := timeGap / produce.GetCycleTime()

// 	for _, v := range produce.GetProduceItems() {
// 		cfg := cfgMap[v.GetItemId()]
// 		if cfg == nil {
// 			logger.ErrorWF("StopMazeProduce cant get cfg", zap.Any("itemId", v.GetItemId()))
// 			continue
// 		}
// 		// 计算最大心跳可以获得的生产数量
// 		addCount := cycleRoundNum * v.GetCycleCount()

// 		// 判断生产加暂存是否超过上限
// 		if addCount+v.GetProduceCount() > cfg.TempMax {
// 			v.ProduceCount = proto.Int64(cfg.TempMax)
// 		} else {
// 			v.ProduceCount = proto.Int64(v.GetProduceCount() + addCount)
// 		}

// 		v.CycleCount = proto.Int64(0)
// 	}

// 	produce.LastSettleTime = proto.Int64(now)
// 	// }

// 	err = mazeproduceredis.SetUserProduce(logger, userId, produce)
// 	if err != nil {
// 		logger.ErrorWF("StopMazeProduce SetUserProduce fail", zap.Error(err))
// 		return
// 	}

// 	logger.InfoWF("StopMazeProduce produce succ", zap.Any("produce", produce))
// 	return
// }

// // 进出区域、心跳、打怪等生产变化掉用
// func StartFixMazeProduce(logger fklog.FKLogI, userId uint64, areaId int32) (err error) {
// 	produce, err := mazeproduceredis.GetUserProduce(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("StartFixMazeProduce GetUserProduce fail", zap.Error(err))
// 		return
// 	}

// 	if produce == nil || produce.ProduceItems == nil || len(produce.GetProduceItems()) == 0 {
// 		return InitMazeProduce(logger, userId, areaId)
// 	}

// 	//修复生产进度 更新生产存储 开始生产
// 	return CheckRealProduceTime(logger, userId, produce, areaId)
// }

// func InitMazeProduce(logger fklog.FKLogI, userId uint64, areaId int32) (err error) {

// 	cfgMap, err := GetCurrProduceCfg(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("InitMazeProduce GetCurrProduceCfg fail", zap.Error(err))
// 		return
// 	}
// 	now := time.Now().Unix()
// 	produce := &DollMazeBarrierCache.DollMazeMoneyDb{
// 		AreaId: proto.Int32(areaId),
// 		// ProduceCount:     proto.Int64(0),
// 		// CycleAddCount:    proto.Int64(cycleCount),
// 		LastSettleTime:   proto.Int64(now),
// 		StartProduceTime: proto.Int64(now),
// 		CycleTime:        proto.Int64(0),
// 		LastHeartTime:    proto.Int64(now),
// 		ProduceItems:     make([]*DollMazeBarrierCache.DollMazeProduceItemDb, 0),
// 	}

// 	for _, v := range cfgMap {
// 		items := &DollMazeBarrierCache.DollMazeProduceItemDb{
// 			ItemId:       proto.Int32(v.ItemId),
// 			ProduceCount: proto.Int64(0),
// 			CycleCount:   proto.Int64(v.CycleAdd),
// 		}
// 		if produce.GetCycleTime() == 0 {
// 			produce.CycleTime = proto.Int64(v.CycleTime / 1000)
// 		}

// 		if v.ItemId == 2 {
// 			canProduce, err2 := CheckEquipPointsProduce(logger, userId, areaId)
// 			if err2 != nil {
// 				logger.ErrorWF("InitMazeProduce CheckEquipPointsProduce fail", zap.Error(err2))
// 				return
// 			}

// 			if !canProduce {
// 				items.CycleCount = proto.Int64(0)
// 			}
// 		}

// 		produce.ProduceItems = append(produce.ProduceItems, items)
// 	}

// 	if produce.GetCycleTime() == 0 {
// 		err = errors.New("cant get cycleTime")
// 		logger.ErrorWF("InitMazeProduce cant get cycleTime")
// 		return
// 	}

// 	err = mazeproduceredis.SetUserProduce(logger, userId, produce)
// 	if err != nil {
// 		logger.ErrorWF("InitMazeProduce SetUserProduce fail", zap.Error(err))
// 		return
// 	}

// 	logger.InfoWF("InitMazeProduce init produce succ", zap.Any("produce", produce), zap.Any("areaId", areaId))

// 	return
// }

// func CheckRealProduceTime(logger fklog.FKLogI, userId uint64, produce *DollMazeBarrierCache.DollMazeMoneyDb, areaId int32) (err error) {
// 	now := time.Now().Unix()

// 	cfgMap, err := GetCurrProduceCfg(logger, userId)
// 	if err != nil {
// 		logger.ErrorWF("CheckRealProduceTime GetCurrProduceCfg fail", zap.Error(err))
// 		return
// 	}

// 	if produce.GetLastSettleTime() == now {
// 		//结算时间一致 不需要结算 直接返回
// 		logger.InfoWF("CheckRealProduceTime no need check", zap.Any("settleTime", produce.GetLastSettleTime()), zap.Any("now", now))
// 		return
// 	}

// 	// 按最大心跳数结算生产数量
// 	timeGap := now - produce.GetLastSettleTime()
// 	if timeGap > HeartTime {
// 		timeGap = HeartTime
// 	}

// 	// 计算最大心跳可以获得的生产数量
// 	cycleRoundNum := timeGap / produce.GetCycleTime()
// 	var isProduceFull bool = true
// 	newProduceItem := make([]*DollMazeBarrierCache.DollMazeProduceItemDb, 0)
// 	var cycleTime int64
// 	//判断是否生产已满 装备积分额外校验未拾取装备数量
// 	for _, item := range produce.GetProduceItems() {
// 		cfg := cfgMap[item.GetItemId()]
// 		if cfg == nil {
// 			logger.ErrorWF("CheckRealProduceTime get produce cfg fail", zap.Any("item", item.GetItemId()))
// 			continue
// 		}
// 		if item.GetProduceCount() < cfg.TempMax {
// 			isProduceFull = false

// 			addCount := cycleRoundNum * item.GetCycleCount()

// 			// 判断生产加暂存是否超过上限
// 			if addCount+item.GetProduceCount() > cfg.TempMax {
// 				item.ProduceCount = proto.Int64(cfg.TempMax)
// 			} else {
// 				item.ProduceCount = proto.Int64(item.GetProduceCount() + addCount)
// 			}
// 			item.CycleCount = proto.Int64(cfg.CycleAdd)
// 		}

// 		if item.GetItemId() == 2 {
// 			canProduce, err2 := CheckEquipPointsProduce(logger, userId, areaId)
// 			if err2 != nil {
// 				logger.ErrorWF("CheckRealProduceTime CheckEquipPointsProduce fail", zap.Error(err))
// 				return err2
// 			}

// 			if !canProduce {
// 				logger.WarnWF("CheckRealProduceTime CheckEquipPointsProduce cant produce")
// 				item.CycleCount = proto.Int64(0)
// 			}
// 		}

// 		cycleTime = cfg.CycleTime / 1000

// 		newProduceItem = append(newProduceItem, item)
// 	}

// 	if isProduceFull {
// 		logger.WarnWF("CheckRealProduceTime produceFull", zap.Any("produce", produce))
// 		return
// 	}

// 	produce.AreaId = proto.Int32(areaId)
// 	produce.ProduceItems = newProduceItem
// 	produce.LastSettleTime = proto.Int64(now) //如果不是一秒一生产 这里不应该直接更新结算时间
// 	produce.StartProduceTime = proto.Int64(now)
// 	produce.LastHeartTime = proto.Int64(now)
// 	produce.CycleTime = proto.Int64(cycleTime)

// 	// if isProduceFull {
// 	// 	//生产已满
// 	// 	logger.InfoWF("CheckRealProduceTime produce full", zap.Any("produce", produce))
// 	// 	return
// 	// }

// 	// if produce.GetCycleAddCount() == 0 {
// 	// 	//如果上次是正常退出生产 周期数量是0,则正确重新生产

// 	// 	produce.CycleAddCount = proto.Int64(cycleCount)

// 	// 	produce.CycleTime = proto.Int64(cycleTime)
// 	// 	produce.LastHeartTime = proto.Int64(now)
// 	// } else {

// 	// 	addCount := cycleRoundNum * produce.GetCycleAddCount()

// 	// 	// 判断生产加暂存是否超过上限
// 	// 	if addCount+produce.GetProduceCount() > tempMax {
// 	// 		produce.ProduceCount = proto.Int64(tempMax)
// 	// 	} else {
// 	// 		produce.ProduceCount = proto.Int64(produce.GetProduceCount() + addCount)
// 	// 	}
// 	// 	produce.CycleAddCount = proto.Int64(cycleCount)
// 	// 	produce.StartProduceTime = proto.Int64(now)
// 	// 	produce.CycleTime = proto.Int64(cycleTime)
// 	// 	produce.LastHeartTime = proto.Int64(now)
// 	// 	produce.LastSettleTime = proto.Int64(now)
// 	// }

// 	err = mazeproduceredis.SetUserProduce(logger, userId, produce)
// 	if err != nil {
// 		logger.ErrorWF("CheckRealProduceTime SetUserProduce fail", zap.Error(err))
// 		return
// 	}

// 	logger.InfoWF("CheckRealProduceTime check produce succ", zap.Any("produce", produce), zap.Any("areaId", areaId))
// 	return
// }

// func CheckEquipPointsProduce(logger fklog.FKLogI, userId uint64, areaId int32) (canProduce bool, err error) {
// 	barrierId := areaId / 10000
// 	userInfo, err := mazeuserbarrier.GetUserBarrier(logger, userId, barrierId)
// 	if err != nil {
// 		logger.ErrorWF("CheckEquipPointsProduce GetUserBarrier fail", zap.Error(err),
// 			zap.Any("barrierid", barrierId))
// 		return
// 	}
// 	maxCfg := GDollMazeConfigV8Cfg.Get(1)
// 	if maxCfg == nil {
// 		logger.ErrorWF("CheckEquipPointsProduce get maze config order 1 fail")
// 		err = errors.New("CheckEquipPointsProduce get maze config order 1 fail")
// 		return
// 	}

// 	areaInfo, ok := userInfo.AreaInfo[areaId]
// 	if !ok {
// 		logger.WarnWF("CheckEquipPointsProduce empty area info", zap.Any("userInfo", userInfo))
// 		canProduce = true
// 		return
// 	}

// 	if areaInfo.CollectEquipNum < int32(maxCfg.Value_map[areaId]) {
// 		canProduce = true
// 	}
// 	logger.InfoWF("CheckEquipPointsProduce info", zap.Any("areaInfo.CollectEquipNum", areaInfo.CollectEquipNum),
// 		zap.Any("maxCfg.Value_map", maxCfg.Value_map), zap.Any("areaId", areaId))
// 	return
// }
