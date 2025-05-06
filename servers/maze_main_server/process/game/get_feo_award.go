package game

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserbarrier"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
)

func OnGetFeoAwardRQ(logger fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	fkprometheus.InfoPMT("OnGetFeoAwardRQ")()

	// req := rqMsg.(*DollMazeBarrier.GetFoeAwardRQ)
	// res := rsMsg.(*DollMazeBarrier.GetFoeAwardRS)

	// logger.InfoWF("OnGetFeoAwardRQ start", zap.Any("req", req))
	// defer func() {
	// 	logger.InfoWF("OnGetFeoAwardRQ end", zap.Any("res", res))
	// }()

	// res.Header = req.Header
	// res.ErrInfo = errors.NO_ERROR
	// res.DefeatedFoeList = req.DefeatedFoeList
	// res.BarrierId = req.BarrierId
	// res.AreaId = req.AreaId

	// userId := shardingID

	// mazeAreaFoe := req.GetMazeReportInfo()
	// if mazeAreaFoe == nil || len(mazeAreaFoe.GetAreaMonsterInfo()) == 0 {
	// 	// logger.ErrorWF("OnGetFeoAwardRQ kill monster empty")
	// 	// res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("没有击败任何怪")
	// 	// 兼容假打旧逻辑
	// 	OnGetFeoAwardOldRQ(logger, shardingID, rqMsg, rsMsg)
	// 	return
	// }

	// // areaCfg := GMazeBrushAreaV8Cfg.Get(mazeAreaFoe.GetAreaId())
	// // if areaCfg == nil {
	// // 	logger.ErrorWF("OnGetFeoAwardRQ get maze brush area cfg nil", zap.Any("areaId", req.GetAreaId()))
	// // 	res.ErrInfo = errors.CONFIG_NOT_FOUND.Wrap("找不到对应的区域id")
	// // 	return
	// // }
	// // if areaCfg.Barries != req.GetBarrierId() {
	// // 	logger.ErrorWF("OnGetFeoAwardRQ barrier and area not match", zap.Any("barrierId", req.GetBarrierId()), zap.Any("areaId", req.GetAreaId()))
	// // 	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡与区域不匹配")
	// // 	return
	// // }

	// //校验怪id是否是当前关卡的
	// barrierId, _, _, err := mazebarrier.GetMazeInfo(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("OnGetFeoAwardRQ GetCurrBarrier fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }
	// if barrierId != req.GetBarrierId() {
	// 	logger.ErrorWF("OnGetFeoAwardRQ curr barrier not match", zap.Any("barrierId", barrierId), zap.Any("req", req.GetBarrierId()))
	// 	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id不匹配")
	// 	return
	// }
	// barrierCfg := GMazeBarriesV8Cfg.Get(req.GetBarrierId())
	// if barrierCfg == nil {
	// 	logger.ErrorWF("OnGetFeoAwardRQ get barrier cfg nil", zap.Any("barrierId", req.GetBarrierId()))
	// 	res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
	// 	return
	// }

	// userInfo, err := mazeuserinfo.GetUserInfo(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("OnGetFeoAwardRQ GetUserInfo fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }
	// level := userInfo.Level
	// levelCfg := GMazeLevelV8Cfg.Get(int32(level))
	// if levelCfg == nil {
	// 	res.ErrInfo = errors.CONFIG_NOT_FOUND.Wrap("等级表获取失败")
	// 	return
	// }
	// shopCfg := GMazeShopV8Cfg.Get(int32(level))
	// if shopCfg == nil {
	// 	logger.ErrorWF("OnGetFeoAwardRQ get maze shop cfg fail", zap.Any("level", level))
	// 	res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
	// 	return
	// }

	// force, err := mazecalcattrredis.GetMazeForce(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("OnGetFeoAwardRQ GetMazeForce fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	// userBarrier, err := mazeuserbarrier.GetUserBarrier(logger, userId, barrierId)
	// if err != nil {
	// 	logger.ErrorWF("OnGetFeoAwardRQ GetUserBarrier fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	// // var currMoneyId int32
	// currMoneyId := barrierCfg.Currency_id
	// var addMoney, addExp, addEquipPoints int64
	// foeMoney := make(map[int32]int64, 0)
	// foeExp := make(map[int32]int64, 0)
	// foeEquipPoints := make(map[int32]int64, 0)
	// var needFoeEquip bool
	// var addEquipMap map[int32]int32
	// uesrArea, ok := userBarrier.AreaInfo[req.GetAreaId()]
	// if !ok {
	// 	uesrArea = &mazeuserbarrier.BarrierArea{
	// 		AreaId:          req.GetAreaId(),
	// 		DefeatedFoe:     make(map[int32]int32),
	// 		CollectEquipNum: 0,
	// 	}
	// }

	// var oldFoeCount, newFoeCount int32
	// oldFoeCount = getTotalFoeCount(userBarrier)

	// record := &dollmazefoekafka.DollMazeFoeRecord{
	// 	UserId:  userId,
	// 	Barrier: barrierId,
	// 	Area:    req.GetAreaId(),
	// 	Level:   int32(level),
	// }

	// foeStrSlice := make([]string, 0)

	// // 遍历上报的打怪总计 对比存储计算奖励
	// for _, area := range mazeAreaFoe.GetAreaMonsterInfo() {
	// 	for _, foe := range area.GetKillMonsterInfo() {
	// 		addFoeCount := foe.GetMonsterCount() - uesrArea.DefeatedFoe[foe.GetConfigId()]
	// 		if addFoeCount < 0 {
	// 			logger.ErrorWF("OnGetFeoAwardRQ req KillMonsterInfo foe count lt save data",
	// 				zap.Any("saveDefeatedFoe", uesrArea), zap.Any("area", area))
	// 			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("上报的杀怪数量小于存储的数值")
	// 			return
	// 		}

	// 		foeCfg := GMazeFoeV8Cfg.Get(foe.GetConfigId())
	// 		if foeCfg == nil {
	// 			logger.ErrorWF("OnGetFeoAwardRQ cant get foe cfg", zap.Any("foeId", foe.GetConfigId()))
	// 			res.ErrInfo = errors.CONFIG_NOT_FOUND.Wrap("找不到怪物配置")
	// 			return
	// 		}

	// 		fExp, fMoney, fEquip, err2 := mazecommonvalue.GetCalRet(logger, userId, level, force,
	// 			foeCfg.Drop_exp_num[int32(level)], foeCfg.Drop_coin_num[int32(level)], foeCfg.Drop_equip_score_num[int32(level)])
	// 		if err2 == nil {
	// 			logger.ErrorWF("OnGetFeoAwardRQ GetCalRet fail", zap.Error(err2))
	// 			res.ErrInfo = errors.CONFIG_NOT_FOUND.Wrap("计算奖励失败")
	// 			return
	// 		}
	// 		foeMoney[foe.GetConfigId()] = fMoney * int64(addFoeCount)
	// 		foeExp[foe.GetConfigId()] = fExp * int64(addFoeCount)
	// 		foeEquipPoints[foe.GetConfigId()] = fEquip * int64(addFoeCount)
	// 		addMoney += fMoney * int64(addFoeCount)
	// 		addExp += fExp * int64(addFoeCount)
	// 		addEquipPoints += fEquip * int64(addFoeCount)

	// 		//更新打怪
	// 		uesrArea.DefeatedFoe[foe.GetConfigId()] = foe.GetMonsterCount()

	// 		foeStrSlice = append(foeStrSlice, fmt.Sprintf("%d:%d", foe.GetConfigId(), foe.GetMonsterCount()))
	// 	}

	// }

	// if addEquipPoints > 0 {
	// 	if userBarrier.EquipPoints+addEquipPoints >= int64(shopCfg.Need_equip_score) {
	// 		needFoeEquip = true
	// 		userBarrier.EquipPoints = userBarrier.EquipPoints + addEquipPoints - int64(shopCfg.Need_equip_score)
	// 	} else {
	// 		userBarrier.EquipPoints = userBarrier.EquipPoints + addEquipPoints
	// 	}
	// }

	// var shopInfo *dollmazeshopseqredis.MazeShopInfo
	// if needFoeEquip {
	// 	//取各部位最大值
	// 	equipMax, isFull, err2 := mazeshopmodule.GetEquipMaxMap(logger, userId, int32(level), userInfo.HighArea)
	// 	if err2 != nil {
	// 		logger.ErrorWF("OnGetFeoAwardRQ GetEquipMaxMap fail", zap.Error(err2), zap.Any("level", level))
	// 		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 		return
	// 	}
	// 	logger.InfoWF("OnGetFeoAwardRQ GetEquipMaxMap dump", zap.Any("equipMax", equipMax), zap.Any("isFull", isFull))
	// 	if isFull {
	// 		needFoeEquip = false
	// 	} else {
	// 		shopInfo, err2 = mazeshopmodule.GetMazeShopInfo(logger, userId, int32(level), userInfo.HighArea)
	// 		if err2 != nil {
	// 			logger.ErrorWF("OnGetFeoAwardRQ GetMazeShopInfo fail", zap.Error(err2), zap.Any("level", level))
	// 			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 			return
	// 		}

	// 		addEquipMap, err2 = mazeshopmodule.GetMazeSeqEquipId(logger, userId, int32(level), userInfo.HighArea, shopInfo, 1, 0, equipMax)
	// 		if err2 != nil {
	// 			logger.ErrorWF("OnGetFeoAwardRQ GetMazeSeqEquipId fail", zap.Error(err2),
	// 				zap.Any("shopInfo", shopInfo), zap.Any("equipMax", equipMax))
	// 			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 			return
	// 		}

	// 		err2 = dollmazeshopseqredis.SetMazeShopInfo(logger, userId, int32(level), shopInfo)
	// 		if err2 != nil {
	// 			logger.ErrorWF("OnGetFeoAwardRQ SetMazeShopInfo fail", zap.Error(err2), zap.Any("level", level), zap.Any("shopInfo", shopInfo))
	// 			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 			return
	// 		}
	// 	}
	// }

	// record.FoeList = strings.Join(foeStrSlice, ",")
	// record.AwardList = fmt.Sprintf("1(%d):%d,2:%d,3:%d", currMoneyId, addMoney, addEquipPoints, addExp)
	// record.Equips = maputil.MapToString32(addEquipMap)
	// record.EquipPoints = int32(userBarrier.EquipPoints)

	// defer func() {
	// 	dollmazefoekafka.PushDollMazeFoeRecord(logger, record)
	// }()

	// userBarrier.AreaInfo[req.GetAreaId()] = uesrArea

	// pbData, err := proto.Marshal(mazeAreaFoe)
	// if err != nil {
	// 	logger.ErrorWF("OnGetFeoAwardRQ Marshal fail", zap.Error(err), zap.Any("mazeAreaFoe", mazeAreaFoe))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	// // 更新杀怪存储
	// err = mazeuserbarrier.SetUserBarrier(logger, userId, barrierId, userBarrier, pbData)
	// if err != nil {
	// 	logger.ErrorWF("OnGetFeoAwardRQ SetUserBarrier fail", zap.Error(err), zap.Any("userBarrier", userBarrier))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	// //根据新旧打死怪的数量 读数值表确定血瓶奖励信息
	// newFoeCount = getTotalFoeCount(userBarrier)
	// var addBlood int64
	// for _, num := range barrierCfg.Drop_vial_foe {
	// 	if oldFoeCount < num && newFoeCount >= num {
	// 		addBlood++
	// 	}
	// 	if oldFoeCount < num && newFoeCount < num {
	// 		break
	// 	}
	// }

	// logger.InfoWF("OnGetFeoAwardRQ all award dump", zap.Any("foeList", mazeAreaFoe), zap.Any("level", level), zap.Any("force", force),
	// 	zap.Any("addMoney", addMoney), zap.Any("addEquipPoint", addEquipPoints), zap.Any("addExp", addExp),
	// 	zap.Any("addBlood", addBlood), zap.Any("isHaveEquip", needFoeEquip), zap.Any("addEquip", addEquipMap))

	// //通用数值id包
	// commonList := make([]*CommonValueStruct, 0)

	// //发货币
	// if addMoney > 0 {
	// 	newCount, err := mazebarriermoneyredis.AddMoney(logger, userId, currMoneyId, addMoney)
	// 	if err == nil {
	// 		// SendMoneyChgPack(logger, userId, currMoneyId, newCount, constdef.DollMazeMoneyChgTypeAdd, req.GetHeader().GetSession())

	// 		commonList = append(commonList, &CommonValueStruct{
	// 			DataType:     int32(DollMazeBarrier.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_MONEY),
	// 			DataValueInt: newCount,
	// 			// ChgReason:    int32(1),
	// 			Session: req.GetHeader().GetSession(),
	// 		})

	// 		record2 := &mazemoneykafka.MazeMoneyRecord{
	// 			UserId:        userId,
	// 			OldMoneyId:    currMoneyId,
	// 			OldMoneyCount: newCount - addMoney,
	// 			NewMoneyId:    currMoneyId,
	// 			NewMoneyCount: newCount,
	// 			ChgReason:     mazemoneykafka.MoneyChgReasonFoe,
	// 		}
	// 		mazemoneykafka.PushMazeMoneyRecord(logger, record2)

	// 		res.AwardItems = append(res.AwardItems, &Common.Item{ItemId: proto.Int32(currMoneyId), Count: proto.Int64(addMoney)})
	// 	} else {
	// 		logger.ErrorWF("OnGetFeoAwardRQ AddMoney fail", zap.Error(err), zap.Any("currMoneyId", currMoneyId), zap.Any("addMoney", addMoney))
	// 	}
	// }

	// //发关卡经验
	// err = userInfo.AddExp(addExp)
	// if err == nil {
	// 	err = mazeuserinfo.SetUserInfo(logger, userId, userInfo)
	// 	if err == nil {
	// 		//等级经验id包
	// 		commonList = append(commonList, &CommonValueStruct{
	// 			DataType:     int32(DollMazeBarrier.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_LEVEL),
	// 			DataValueInt: userInfo.Level,
	// 			// ChgReason:    int32(1),
	// 			Session: req.GetHeader().GetSession(),
	// 		})
	// 		commonList = append(commonList, &CommonValueStruct{
	// 			DataType:     int32(DollMazeBarrier.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP),
	// 			DataValueInt: userInfo.Exp,
	// 			// ChgReason:    int32(1),
	// 			Session: req.GetHeader().GetSession(),
	// 		})
	// 		newLevelCfg := GMazeLevelV8Cfg.Get(int32(userInfo.Level))
	// 		if newLevelCfg != nil {
	// 			commonList = append(commonList, &CommonValueStruct{
	// 				DataType:     int32(DollMazeBarrier.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP_MAX),
	// 				DataValueInt: newLevelCfg.Next_level_need_exp,
	// 				// ChgReason:    int32(1),
	// 				Session: req.GetHeader().GetSession(),
	// 			})
	// 		} else {
	// 			logger.ErrorWF("OnGetFeoAwardRQ get new level cfg fail", zap.Any("level", userInfo.Level))
	// 		}

	// 		extra, err3 := MakeCommonValueExtra(logger, userId, userInfo.Level, force)
	// 		if err3 == nil {
	// 			commonList = append(commonList, &CommonValueStruct{
	// 				DataType:     int32(DollMazeBarrier.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_INCOME),
	// 				DataValueInt: extra,
	// 				// ChgReason:    int32(1),
	// 				Session: req.GetHeader().GetSession(),
	// 			})
	// 		} else {
	// 			logger.ErrorWF("OnGetFeoAwardRQ MakeCommonValueExtra fail", zap.Error(err3), zap.Any("level", userInfo.Level), zap.Any("force", force))
	// 		}

	// 		res.AwardItems = append(res.AwardItems, &Common.Item{ItemId: proto.Int32(46500001), Count: proto.Int64(addExp)})

	// 	} else {
	// 		logger.ErrorWF("OnGetFeoAwardRQ SetUserInfo fail", zap.Error(err), zap.Any("userInfo", userInfo))
	// 	}
	// } else {
	// 	logger.ErrorWF("OnGetFeoAwardRQ AddExp fail", zap.Error(err), zap.Any("addExp", addExp))
	// }

	// //发物品
	// tradeNo := gentradeno.GetTradeNum()
	// if addBlood > 0 {
	// 	record.AwardList += fmt.Sprintf(",%d:%d", barrierCfg.Hp_vial, addBlood)

	// 	addItems := itemutil.Map2Common(map[int32]int64{barrierCfg.Hp_vial: addBlood})
	// 	errInfo := gentradeno.AddItemEx(logger, userId, 689, tradeNo, 0, 1, req.GetHeader(), addItems...)
	// 	if errInfo != nil {
	// 		logger.ErrorWF("OnGetFeoAwardRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("addItems", addItems))
	// 	}

	// 	res.AwardItems = append(res.AwardItems, &Common.Item{ItemId: proto.Int32(barrierCfg.Hp_vial), Count: proto.Int64(addBlood)})
	// }

	// //发装备
	// if needFoeEquip {
	// 	// addEquipMap = map[int32]int32{42711001: 1}
	// 	// DOLL_EQUIP_MAZE_FOE = 19;//迷宫怪物掉落
	// 	// DOLL_EQUIP_MAZE_BOX_AWARD = 20;//迷宫宝箱掉落
	// 	rs, err2 := addequip.AddEquipToBag(logger, userId, 19, tradeNo, addEquipMap)
	// 	if err2 != nil {
	// 		logger.ErrorWF("OnGetFeoAwardRQ addEquipToBag fail", zap.Error(err2), zap.Any("optype", 19),
	// 			zap.Any("tradeNo", tradeNo), zap.Any("addEquip", addEquipMap), zap.Any("rs", rs))
	// 	}
	// 	PushDollMazeShopInfoLog(logger, userId, int32(level), shopInfo, rs.EquipList, 19, tradeNo, 0)
	// }

	// SendCommonValueIdPack(logger, userId, commonList)

	return nil
}

// func SendMoneyChgPack(logger fklog.FKLogI, userId uint64, currMoneyId int32, currCount int64, reason int32, session string) (err error) {
// 	moneyPack := &DollMazeBarrier.BarrierMoneyID{
// 		NewMoney: &Common.Item{ItemId: proto.Int32(currMoneyId), Count: proto.Int64(currCount)},
// 		Reason:   proto.Int32(reason),
// 		Session:  proto.String(session),
// 		Token:    proto.Int64(time.Now().UnixNano() / 1e6),
// 	}

// 	logger.InfoWF("SendMoneyChgPack send client with", zap.Any("moneyPack", moneyPack))
// 	return commonmustarriveredis.SendArrivePacketFix(userId, 16139, moneyPack)
// }

func OnGetFeoAwardOldRQ(logger fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	fkprometheus.InfoPMT("OnGetFeoAwardOldRQ")()

	// req := rqMsg.(*DollMazeBarrier.GetFoeAwardRQ)
	// res := rsMsg.(*DollMazeBarrier.GetFoeAwardRS)

	// logger.InfoWF("OnGetFeoAwardOldRQ start", zap.Any("req", req))
	// defer func() {
	// 	logger.InfoWF("OnGetFeoAwardOldRQ end", zap.Any("res", res))
	// }()

	// res.Header = req.Header
	// res.ErrInfo = errors.NO_ERROR
	// res.DefeatedFoeList = req.DefeatedFoeList
	// res.BarrierId = req.BarrierId
	// res.AreaId = req.AreaId

	// userId := shardingID

	// // areaCfg := GMazeBrushAreaV8Cfg.Get(req.GetAreaId())
	// // if areaCfg == nil {
	// // 	logger.ErrorWF("OnGetFeoAwardOldRQ get maze brush area cfg nil", zap.Any("areaId", req.GetAreaId()))
	// // 	res.ErrInfo = errors.CONFIG_NOT_FOUND.Wrap("找不到对应的区域id")
	// // 	return
	// // }
	// // if areaCfg.Barries != req.GetBarrierId() {
	// // 	logger.ErrorWF("OnGetFeoAwardOldRQ barrier and area not match", zap.Any("barrierId", req.GetBarrierId()), zap.Any("areaId", req.GetAreaId()))
	// // 	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡与区域不匹配")
	// // 	return
	// // }

	// //校验怪id是否是当前关卡的
	// barrierId, _, _, err := mazebarrier.GetMazeInfo(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("OnGetFeoAwardOldRQ GetCurrBarrier fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }
	// if barrierId != req.GetBarrierId() {
	// 	logger.ErrorWF("OnGetFeoAwardOldRQ curr barrier not match", zap.Any("barrierId", barrierId), zap.Any("req", req.GetBarrierId()))
	// 	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id不匹配")
	// 	return
	// }
	// barrierCfg := GMazeBarriesV8Cfg.Get(req.GetBarrierId())
	// if barrierCfg == nil {
	// 	logger.ErrorWF("OnGetFeoAwardOldRQ get barrier cfg nil", zap.Any("barrierId", req.GetBarrierId()))
	// 	res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
	// 	return
	// }

	// userInfo, err := mazeuserinfo.GetUserInfo(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("OnGetFeoAwardOldRQ GetUserInfo fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }
	// level := userInfo.Level
	// levelCfg := GMazeLevelV8Cfg.Get(int32(level))
	// if levelCfg == nil {
	// 	res.ErrInfo = errors.CONFIG_NOT_FOUND.Wrap("等级表获取失败")
	// 	return
	// }
	// shopCfg := GMazeShopV8Cfg.Get(int32(level))
	// if shopCfg == nil {
	// 	logger.ErrorWF("OnGetFeoAwardOldRQ get maze shop cfg fail", zap.Any("level", level))
	// 	res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
	// 	return
	// }

	// force, err := mazecalcattrredis.GetMazeForce(logger, userId)
	// if err != nil {
	// 	logger.ErrorWF("OnGetFeoAwardOldRQ GetMazeForce fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	// userBarrier, err := mazeuserbarrier.GetUserBarrier(logger, userId, barrierId)
	// if err != nil {
	// 	logger.ErrorWF("OnGetFeoAwardOldRQ GetUserBarrier fail", zap.Error(err))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	// // var currMoneyId int32
	// currMoneyId := barrierCfg.Currency_id
	// var addMoney int64
	// foeMoney := make(map[int32]int64, 0)
	// var addExp int64
	// foeExp := make(map[int32]int64, 0)
	// var addEquipPoints int64
	// foeEquipPoints := make(map[int32]int64, 0)
	// var needFoeEquip bool
	// var addEquipMap map[int32]int32
	// uesrArea, ok := userBarrier.AreaInfo[req.GetAreaId()]
	// if !ok {
	// 	uesrArea = &mazeuserbarrier.BarrierArea{
	// 		AreaId:          req.GetAreaId(),
	// 		DefeatedFoe:     make(map[int32]int32),
	// 		CollectEquipNum: 0,
	// 	}
	// }

	// var oldFoeCount, newFoeCount int32
	// oldFoeCount = getTotalFoeCount(userBarrier)

	// record := &dollmazefoekafka.DollMazeFoeRecord{
	// 	UserId:  userId,
	// 	Barrier: barrierId,
	// 	Area:    req.GetAreaId(),
	// 	Level:   int32(level),
	// }

	// foeStrSlice := make([]string, 0)

	// // 遍历上报的打怪总计 对比存储计算奖励
	// for _, foe := range req.GetDefeatedFoeList() {

	// 	addFoeCount := foe.GetFoeCount()

	// 	foeCfg := GMazeFoeV8Cfg.Get(foe.GetFoeId())
	// 	if foeCfg == nil {
	// 		logger.ErrorWF("OnGetFeoAwardOldRQ cant get foe cfg", zap.Any("foeId", foe.GetFoeId()))
	// 		res.ErrInfo = errors.CONFIG_NOT_FOUND.Wrap("找不到怪物配置")
	// 		return
	// 	}

	// 	fExp, fMoney, fEquip, err2 := mazecommonvalue.GetCalRet(logger, userId, level, force,
	// 		foeCfg.Drop_exp_num[int32(level)], foeCfg.Drop_coin_num[int32(level)], foeCfg.Drop_equip_score_num[int32(level)])
	// 	if err2 != nil {
	// 		logger.ErrorWF("OnGetFeoAwardOldRQ GetCalRet fail", zap.Error(err2))
	// 		res.ErrInfo = errors.CONFIG_NOT_FOUND.Wrap("计算奖励失败")
	// 		return
	// 	}
	// 	foeMoney[foe.GetFoeId()] = fMoney * int64(addFoeCount)
	// 	foeExp[foe.GetFoeId()] = fExp * int64(addFoeCount)
	// 	foeEquipPoints[foe.GetFoeId()] = fEquip * int64(addFoeCount)
	// 	addMoney += fMoney * int64(addFoeCount)
	// 	addExp += fExp * int64(addFoeCount)
	// 	addEquipPoints += fEquip * int64(addFoeCount)

	// 	//更新打怪
	// 	uesrArea.DefeatedFoe[foe.GetFoeId()] += foe.GetFoeCount()

	// 	foeStrSlice = append(foeStrSlice, fmt.Sprintf("%d:%d", foe.GetFoeId(), foe.GetFoeCount()))
	// }

	// if addEquipPoints > 0 {
	// 	if userBarrier.EquipPoints+addEquipPoints >= int64(shopCfg.Need_equip_score) {
	// 		needFoeEquip = true
	// 		userBarrier.EquipPoints = userBarrier.EquipPoints + addEquipPoints - int64(shopCfg.Need_equip_score)
	// 	} else {
	// 		userBarrier.EquipPoints = userBarrier.EquipPoints + addEquipPoints
	// 	}
	// }

	// var shopInfo *dollmazeshopseqredis.MazeShopInfo
	// if needFoeEquip {
	// 	//取各部位最大值
	// 	equipMax, isFull, err2 := mazeshopmodule.GetEquipMaxMap(logger, userId, int32(level), userInfo.HighArea)
	// 	if err2 != nil {
	// 		logger.ErrorWF("OnGetFeoAwardOldRQ GetEquipMaxMap fail", zap.Error(err2), zap.Any("level", level))
	// 		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 		return
	// 	}
	// 	logger.InfoWF("OnGetFeoAwardOldRQ GetEquipMaxMap dump", zap.Any("equipMax", equipMax), zap.Any("isFull", isFull))
	// 	if isFull {
	// 		needFoeEquip = false
	// 	} else {
	// 		shopInfo, err2 = mazeshopmodule.GetMazeShopInfo(logger, userId, int32(level), userInfo.HighArea)
	// 		if err2 != nil {
	// 			logger.ErrorWF("OnGetFeoAwardOldRQ GetMazeShopInfo fail", zap.Error(err2), zap.Any("level", level))
	// 			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 			return
	// 		}

	// 		addEquipMap, err2 = mazeshopmodule.GetMazeSeqEquipId(logger, userId, int32(level), userInfo.HighArea, shopInfo, 1, 0, equipMax)
	// 		if err2 != nil {
	// 			logger.ErrorWF("OnGetFeoAwardOldRQ GetMazeSeqEquipId fail", zap.Error(err2),
	// 				zap.Any("shopInfo", shopInfo), zap.Any("equipMax", equipMax))
	// 			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 			return
	// 		}

	// 		err2 = dollmazeshopseqredis.SetMazeShopInfo(logger, userId, int32(level), shopInfo)
	// 		if err2 != nil {
	// 			logger.ErrorWF("OnGetFeoAwardOldRQ SetMazeShopInfo fail", zap.Error(err2), zap.Any("level", level), zap.Any("shopInfo", shopInfo))
	// 			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 			return
	// 		}
	// 	}
	// }

	// record.FoeList = strings.Join(foeStrSlice, ",")
	// record.AwardList = fmt.Sprintf("1(%d):%d,2:%d,3:%d", currMoneyId, addMoney, addEquipPoints, addExp)
	// record.Equips = maputil.MapToString32(addEquipMap)
	// record.EquipPoints = int32(userBarrier.EquipPoints)

	// defer func() {
	// 	dollmazefoekafka.PushDollMazeFoeRecord(logger, record)
	// }()

	// userBarrier.AreaInfo[req.GetAreaId()] = uesrArea

	// // 更新杀怪存储
	// err = mazeuserbarrier.SetUserBarrier(logger, userId, barrierId, userBarrier, nil)
	// if err != nil {
	// 	logger.ErrorWF("OnGetFeoAwardOldRQ SetUserBarrier fail", zap.Error(err), zap.Any("userBarrier", userBarrier))
	// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
	// 	return
	// }

	// //根据新旧打死怪的数量 读数值表确定血瓶奖励信息
	// newFoeCount = getTotalFoeCount(userBarrier)
	// var addBlood int64
	// for _, num := range barrierCfg.Drop_vial_foe {
	// 	if oldFoeCount < num && newFoeCount >= num {
	// 		addBlood++
	// 	}
	// 	if oldFoeCount < num && newFoeCount < num {
	// 		break
	// 	}
	// }

	// logger.InfoWF("OnGetFeoAwardRQ all award dump", zap.Any("foeList", req.GetDefeatedFoeList()), zap.Any("level", level), zap.Any("force", force),
	// 	zap.Any("addMoney", addMoney), zap.Any("addEquipPoint", addEquipPoints), zap.Any("addExp", addExp),
	// 	zap.Any("addBlood", addBlood), zap.Any("isHaveEquip", needFoeEquip), zap.Any("addEquip", addEquipMap))

	// //通用数值id包
	// commonList := make([]*CommonValueStruct, 0)

	// //发货币
	// if addMoney > 0 {
	// 	newCount, err := mazebarriermoneyredis.AddMoney(logger, userId, currMoneyId, addMoney)
	// 	if err == nil {
	// 		// SendMoneyChgPack(logger, userId, currMoneyId, newCount, constdef.DollMazeMoneyChgTypeAdd, req.GetHeader().GetSession())

	// 		commonList = append(commonList, &CommonValueStruct{
	// 			DataType:     int32(DollMazeBarrier.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_MONEY),
	// 			DataValueInt: newCount,
	// 			// ChgReason:    int32(1),
	// 			Session: req.GetHeader().GetSession(),
	// 		})

	// 		record2 := &mazemoneykafka.MazeMoneyRecord{
	// 			UserId:        userId,
	// 			OldMoneyId:    currMoneyId,
	// 			OldMoneyCount: newCount - addMoney,
	// 			NewMoneyId:    currMoneyId,
	// 			NewMoneyCount: newCount,
	// 			ChgReason:     mazemoneykafka.MoneyChgReasonFoe,
	// 		}
	// 		mazemoneykafka.PushMazeMoneyRecord(logger, record2)

	// 		res.AwardItems = append(res.AwardItems, &Common.Item{ItemId: proto.Int32(currMoneyId), Count: proto.Int64(addMoney)})
	// 	} else {
	// 		logger.ErrorWF("OnGetFeoAwardOldRQ AddMoney fail", zap.Error(err), zap.Any("currMoneyId", currMoneyId), zap.Any("addMoney", addMoney))
	// 	}

	// }

	// //发关卡经验
	// err = userInfo.AddExp(addExp)
	// if err == nil {
	// 	err = mazeuserinfo.SetUserInfo(logger, userId, userInfo)
	// 	if err == nil {
	// 		//等级经验id包
	// 		commonList = append(commonList, &CommonValueStruct{
	// 			DataType:     int32(DollMazeBarrier.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_LEVEL),
	// 			DataValueInt: userInfo.Level,
	// 			// ChgReason:    int32(1),
	// 			Session: req.GetHeader().GetSession(),
	// 		})
	// 		commonList = append(commonList, &CommonValueStruct{
	// 			DataType:     int32(DollMazeBarrier.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP),
	// 			DataValueInt: userInfo.Exp,
	// 			// ChgReason:    int32(1),
	// 			Session: req.GetHeader().GetSession(),
	// 		})
	// 		newLevelCfg := GMazeLevelV8Cfg.Get(int32(userInfo.Level))
	// 		if newLevelCfg != nil {
	// 			commonList = append(commonList, &CommonValueStruct{
	// 				DataType:     int32(DollMazeBarrier.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP_MAX),
	// 				DataValueInt: newLevelCfg.Next_level_need_exp,
	// 				// ChgReason:    int32(1),
	// 				Session: req.GetHeader().GetSession(),
	// 			})
	// 		} else {
	// 			logger.ErrorWF("OnGetFeoAwardOldRQ get new level cfg fail", zap.Any("level", userInfo.Level))
	// 		}

	// 		extra, err3 := MakeCommonValueExtra(logger, userId, userInfo.Level, force)
	// 		if err3 == nil {
	// 			commonList = append(commonList, &CommonValueStruct{
	// 				DataType:     int32(DollMazeBarrier.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_INCOME),
	// 				DataValueInt: extra,
	// 				// ChgReason:    int32(1),
	// 				Session: req.GetHeader().GetSession(),
	// 			})
	// 		} else {
	// 			logger.ErrorWF("OnGetFeoAwardOldRQ MakeCommonValueExtra fail", zap.Error(err3), zap.Any("level", userInfo.Level), zap.Any("force", force))
	// 		}
	// 		res.AwardItems = append(res.AwardItems, &Common.Item{ItemId: proto.Int32(46500001), Count: proto.Int64(addExp)})
	// 	} else {
	// 		logger.ErrorWF("OnGetFeoAwardOldRQ SetUserInfo fail", zap.Error(err), zap.Any("userInfo", userInfo))
	// 	}
	// } else {
	// 	logger.ErrorWF("OnGetFeoAwardOldRQ AddExp fail", zap.Error(err), zap.Any("addExp", addExp))
	// }

	// //发物品
	// tradeNo := gentradeno.GetTradeNum()
	// if addBlood > 0 {
	// 	record.AwardList += fmt.Sprintf(",%d:%d", barrierCfg.Hp_vial, addBlood)

	// 	addItems := itemutil.Map2Common(map[int32]int64{barrierCfg.Hp_vial: addBlood})
	// 	errInfo := gentradeno.AddItemEx(logger, userId, 689, tradeNo, 0, 1, req.GetHeader(), addItems...)
	// 	if errInfo != nil {
	// 		logger.ErrorWF("OnGetFeoAwardOldRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("addItems", addItems))
	// 	}

	// 	res.AwardItems = append(res.AwardItems, &Common.Item{ItemId: proto.Int32(barrierCfg.Hp_vial), Count: proto.Int64(addBlood)})
	// }

	// //发装备
	// if needFoeEquip {
	// 	// addEquipMap = map[int32]int32{42711001: 1}
	// 	// DOLL_EQUIP_MAZE_FOE = 19;//迷宫怪物掉落
	// 	// DOLL_EQUIP_MAZE_BOX_AWARD = 20;//迷宫宝箱掉落
	// 	rs, err2 := addequip.AddEquipToBag(logger, userId, 19, tradeNo, addEquipMap)
	// 	if err2 != nil {
	// 		logger.ErrorWF("OnGetFeoAwardOldRQ addEquipToBag fail", zap.Error(err2), zap.Any("optype", 19),
	// 			zap.Any("tradeNo", tradeNo), zap.Any("addEquip", addEquipMap), zap.Any("rs", rs))
	// 	}
	// 	PushDollMazeShopInfoLog(logger, userId, int32(level), shopInfo, rs.EquipList, 19, tradeNo, 0)
	// }

	// SendCommonValueIdPack(logger, userId, commonList)

	return
}

func getTotalFoeCount(userInfo *mazeuserbarrier.MazeUserBarrier) (total int32) {
	for _, area := range userInfo.AreaInfo {
		for _, v := range area.DefeatedFoe {
			total += v
		}
	}
	return
}
