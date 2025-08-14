package game

import (
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/gentradeno"
	"maze_game_server/common/structsdef"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeLevelV8Cfg"
	"maze_game_server/io/kafka/mazeenergyrecord"
	"maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"maze_game_server/io/redis/mazebarriereventredis"
	"maze_game_server/io/redis/mazebarrieropstatusredis"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/io/redis/mazeuserbarrierredis"
	"maze_game_server/io/redis/syncmazestorageinforedis"
	"maze_game_server/lib/codec"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/equipdropmodel"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeAIBattle"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeEnergy"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/servers/maze_main_server/process/game/events"
	"maze_game_server/services/barrierarearecordservice"
	"maze_game_server/services/barrierenergyservice"
	"maze_game_server/services/barriersavedataservice"
	"maze_game_server/services/tempbuffservice"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (g *Game) OnMazeBarrierEnterRQ_10447_10448(s *session.Session, req *MazeGame.MazeBarrierEnterRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeBarrierEnterRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.MazeBarrierEnterRS{}
	energyID := &MazeEnergy.EnergyChangeID{} //defer时多补一个体力ID包

	logger.InfoWF("OnMazeBarrierEnterRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeBarrierEnterRQ end", zap.Any("res", res))

		err = s.ResponseMID(codec.ToMessageID(uint32(time.Now().Unix()), 0, 10610), energyID)
		logger.InfoWF("OnMazeBarrierEnterRQ end send EnergyChangeID", zap.Any("energyID", energyID))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	if req.GetBarrierId() <= 0 {
		logger.ErrorWF("OnMazeBarrierEnterRQ req barrier invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id未设置")
		return
	}
	barrierCfg := GMazeBarriesV8Cfg.Get(req.GetBarrierId())
	if barrierCfg == nil {
		logger.ErrorWF("OnMazeBarrierEnterRQ get barrier cfg fail", zap.Any("barrier", req.GetBarrierId()))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("找不到该关卡配置")
		return
	}

	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierEnterRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	energy, _, err := barrierenergyservice.GlobalBarrierEnergyService.GetBarrierEnergy(logger, userId)
	if err != nil {
		return err
	}
	energyID.EnergyInfo = &MazeEnergy.EnergyInfo{
		CurVal:           proto.Int32(energy),
		MaxVal:           proto.Int32(barrierenergyservice.GlobalBarrierEnergyService.GetEnergyMaxValue()),
		NextRecoveryTime: proto.Int64(userInfo.EnergyLastTime),
	}
	oldEnergy := energy

	// TODO 客户端需要进入任意关卡
	// if req.GetBarrierId() < userInfo.Barrier {
	// 	logger.ErrorWF("OnMazeBarrierEnterRQ req barrier lt pass barrier", zap.Any("req", req), zap.Int32("save", userInfo.Barrier))
	// 	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("该关卡id小于存储的关卡id")
	// 	return
	// }

	//	res.Energy = proto.Int32(userInfo.Energy)
	var isNewBarrier bool
	storageInfo, _ := syncmazestorageinforedis.GetSyncMazeStorageInfo(userId, req.GetBarrierId())

	// 获取存档数据 new
	saveData, err := barriersavedataservice.GlobalBarrierSaveDataService.GetBarrierSaveData(logger, userId, req.GetBarrierId())
	if err != nil {
		logger.ErrorWF("OnMazeBarrierEnterRQ GetBarrierSaveData err", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return err
	}
	if storageInfo == nil || saveData.StageId == 0 {
		mazebuffinforedis.DelMazeBuffBySrc(logger, userId, constdef.MazeBuffSrcSelectBuffForce)
		// 推送属性计算消息
		calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{
			UserId:  userId,
			ChgType: constdef.MazeBuffChgForceValue,
			Session: "buff",
			BuffSrc: constdef.MazeBuffSrcSelectBuffForce,
		}
		mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, calcAttrNotify)
	} else {
		// 刷一半的情况需要检查三选一是否有问题
		tempBuff, err := tempbuffservice.GlobalTempBuffService.CheckTempBuff(logger, userId, req.GetBarrierId(), saveData.StageId)
		if err != nil {
			logger.ErrorWF("OnMazeBarrierEnterRQ checkTempBuff", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return err
		}
		if tempBuff != nil && tempBuff.BuffSequence != nil {
			res.EnergyLevel = proto.Int32(tempBuff.BuffSequence.Level)
		}

		//刷一半的情况需要把未通过的区域杀怪记录删除
		err = barrierarearecordservice.GlobalBarrierAreaRecordService.DelBarrierAreaRecord(logger, userId, req.GetBarrierId())
		if err != nil {
			logger.ErrorWF("OnMazeBarrierEnterRQ DelBarrierAreaRecord fail", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return err
		}
	}

	res.SaveData = &MazeGame.BarrierSaveData{
		StageId:     proto.Int32(saveData.StageId),
		RescueValue: proto.Int32(saveData.RescueValue),
		BossPower:   proto.Int32(saveData.BossPower),
	}
	initPassValue, err := mazecommonvalue.CalcInitPassValue(logger, req.GetBarrierId())
	if err != nil {
		logger.ErrorWF("OnMazeBarrierEnterRQ CalcInitPassvalue fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	res.InitPassValue = proto.Int32(int32(initPassValue))
	// 推送通关值
	mazecommonvalue.SendPassValueIdPack(logger, userId, req.GetBarrierId(), saveData.StageId)
	// 清理关卡操作状态
	mazebarrieropstatusredis.ClearOpStatus(logger, userId, req.GetBarrierId())

	mazeBattleInfo, err3 := GetMazeBattleData(logger, userId, req.GetBarrierId())
	if err3 != nil {
		logger.ErrorWF("OnMazeBarrierEnterRQ GetMazeBattleData fail", zap.Error(err3))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	res.MazeBarrierInfo = mazeBattleInfo

	userBarrier, err := mazeuserbarrierredis.GetUserBarrierInfo(logger, userId, 0)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierEnterRQ GetUserBarrierInfo fail", zap.Error(err3))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if userInfo.Barrier == req.GetBarrierId() && userBarrier.GetBarrierStatus() == 2 {
		isNewBarrier = false
	} else {
		isNewBarrier = true
	}
	// if req.GetBarrierId() > userInfo.Barrier {
	isNewBarrier = true
	userInfo.SetBarrier(req.GetBarrierId())
	// }

	//shopInfo, err := calequipsequence.GetMazeShopInfo(logger, userId, int32(userInfo.Level), req.GetBarrierId())
	//if err != nil {
	//	logger.ErrorWF("OnMazeBarrierEnterRQ GetMazeShopInfo fail", zap.Error(err))
	//	return
	//}
	dropInfo, err := equipdropmodel.NewEquipSpecialDropModel(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierEnterRQ GetEquipSpecialDropModel fail", zap.Error(err))
		return
	}

	var curEnergy int32
	//进入关卡需要

	//首次进入新关还额外需要
	//0. 扣次数
	//1. 更新记录的关卡id
	//2. 判断是否切换装备序列 清空装备积分 (不需要清 旧关卡积分保留 扫荡会继续加
	//3. 清临时buff
	if isNewBarrier {
		//首次进入判断体力是否足够 直接扣根据错误码判断
		// isEnergyEnough, remainVal, err2 := SubUserEnergy(logger, userId, barrierCfg.Mop_cost)
		// if err2 != nil {
		// 	logger.ErrorWF("OnMazeBarrierEnterRQ SubUserEnergy fail", zap.Error(err2))
		// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		// 	return
		// }
		// if !isEnergyEnough {
		// 	logger.WarnWF("OnMazeBarrierEnterRQ user energy not enough", zap.Any("user", remainVal), zap.Any("need", barrierCfg.Mop_cost))
		// 	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("体力不足")
		// 	return
		// }
		// curEnergy = remainVal

		//扣体力
		curEnergy, err = barrierenergyservice.GlobalBarrierEnergyService.SubEnergy(logger, userId, barrierCfg.Mop_cost)
		if err != nil {
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("体力不足")
			logger.ErrorWF("OnMazeBarrierEnterRQ SubEnergy fail", zap.Error(err))
			return err
		}

		energyID.EnergyInfo = &MazeEnergy.EnergyInfo{
			CurVal:           proto.Int32(curEnergy),
			MaxVal:           proto.Int32(barrierenergyservice.GlobalBarrierEnergyService.GetEnergyMaxValue()),
			NextRecoveryTime: proto.Int64(userInfo.EnergyLastTime),
		}

		defer func() {
			barrierenergyservice.GlobalBarrierEnergyService.PushEnergyRecord(logger, userId, oldEnergy, curEnergy, mazeenergyrecord.EnterBarrier, userInfo.EnergyLastTime)
		}()
		//扣次数
		//var maxNum int32
		//maxNumCfg := GMazeActionCountV8Cfg.Get(101)
		//if maxNumCfg == nil {
		//	logger.ErrorWF("OnMazeBarrierEnterRQ GMazeActionCountV8Cfg fail", zap.Error(err))
		//	res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		//	return
		//}
		//maxNum = maxNumCfg.Day_count_v8
		//now := time.Now()
		//today := now.Year()*10000 + int(now.Month())*100 + now.Day()
		//
		//useNumToday, err2 := mazechallengenumredis.GetUserChallengeNum(logger, userId, today)
		//if err2 != nil {
		//	logger.ErrorWF("OnMazeBarrierEnterRQ GetUserChallengeNum fail", zap.Error(err2))
		//	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		//	return
		//}
		//if int32(useNumToday)+barrierCfg.Challenge_cost > maxNum {
		//	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("次数不足")
		//	return
		//}
		//err = mazechallengenumredis.AddUserChallengeNum(logger, userId, today, barrierCfg.Challenge_cost)
		//if err != nil {
		//	logger.ErrorWF("OnMazeBarrierEnterRQ AddUserChallengeNum fail", zap.Error(err))
		//	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		//	return
		//}
		//
		err = mazeuserinfo.SetUserInfoV2(logger, userId, userInfo)
		if err != nil {
			logger.ErrorWF("OnMazeBarrierEnterRQ SetUserInfoV2 fail", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}

		// 记录用户关卡状态 清除客户端上报数据
		userBarrier.BarrierId = proto.Int32(req.GetBarrierId())
		userBarrier.BarrierStatus = proto.Int32(2)
		userBarrier.RebornCount = proto.Int32(0)
		userBarrier.StartTime = proto.Int64(time.Now().Unix())
		userBarrier.EndTime = proto.Int64(0)
		err = mazeuserbarrierredis.SetUserBarrierInfo(logger, userId, 0, userBarrier)
		if err != nil {
			logger.ErrorWF("OnMazeBarrierEnterRQ SetUserBarrierInfo fail", zap.Error(err))
		}

		//2. 判断是否切换装备序列 清空装备积分 (不需要清 旧关卡积分保留 扫荡会继续加)

		// //3. 首次进入清临时buff
		// mazebarriertempbuffredis.ClearBarrierTempBuff(logger, userId, req.GetBarrierId())
	}

	_, rebornMax, _ := GetReviveCost(userBarrier.GetRebornCount() + 1)
	res.MazeReportInfo = &MazeAIBattle.MazeAIReportInfo{
		RebornInfo: &MazeCommon.MazeCount{
			CurCount:   proto.Int32(userBarrier.GetRebornCount()),
			TotalCount: proto.Int32(rebornMax),
		},
	}

	// 同步给客户端当前服务器记录的通用数值
	var expMax, money, diamond int64
	levelCfg := GMazeLevelV8Cfg.Get(int32(userInfo.Level))
	if levelCfg != nil {
		expMax = levelCfg.Next_level_need_exp
	}

	items := []*MazeCommon.MazeItem{
		&MazeCommon.MazeItem{ItemId: proto.Int32(constdef.MazeCommonItemCoin)},
		&MazeCommon.MazeItem{ItemId: proto.Int32(constdef.MazeCommonItemDiamond)}}
	queryItems, errInfo := gentradeno.QueryItems(logger, userId, items...)
	if errInfo == nil {
		for _, v := range queryItems {
			if v.GetItemId() == constdef.MazeCommonItemCoin {
				money = v.GetCount()
			} else if v.GetItemId() == constdef.MazeCommonItemDiamond {
				diamond = v.GetCount()
			}
		}
	}

	income, err := mazecommonvalue.MakeCommonValueExtra(logger, userId, userInfo.Level, 0)

	res.SyncData = &MazeGame.MazeCommonValueSync{
		Level:      proto.Int64(userInfo.Level),
		Exp:        proto.Int64(userInfo.Exp),
		ExpMax:     proto.Int64(expMax),
		Money:      proto.Int64(money),
		Income:     proto.Int64(income),
		Diamond:    proto.Int64(diamond),
		EquipPoint: proto.Int64(int64(dropInfo.EquipPoints)),
		Energy:     proto.Int32(curEnergy),
	}

	err = mazebarriereventredis.EnterBarrier(logger, userId, req.GetBarrierId())
	if err != nil {
		logger.ErrorWF("OnMazeBarrierEnterRQ EnterBarrier fail", zap.Error(err))
	}

	// 触发进入关卡事件
	events.OnEnterBarrier(logger, userId, 0, time.Now().UnixMilli(), &MazeGame.BattleEventEnterBarrier{BarrierId: proto.Int32(req.GetBarrierId())})
	return nil
}

//func SubUserEnergy(logger fklog.FKLogI, uid uint64, subEnergy int32) (isSucc bool, newEnergy int32, err error) {
//	req := &MazeEnergySvr.SubMazeEnergyRQ{
//		UserId:      proto.Uint64(uid),
//		SubVal:      proto.Int32(subEnergy),
//		OpType:      proto.Int32(1), //NUM_MAZE_ENERGY_OP_TYPE_CHALLLENGE
//		OpDesc:      proto.String("maze_barrier_enter"),
//		TradeNumber: proto.Uint64(gentradeno.GetTradeNum()),
//	}
//	res := &MazeEnergySvr.SubMazeEnergyRS{}
//	// 合并服务，内聚接口
//	// err = mazeenergyrpc.SubMazeEnergyRQ(logger, req, res)
//	err = energy.SubMazeEnergyRQ(logger, uid, req, res)
//	if err != nil {
//		logger.ErrorWF("OnMazeBarrierEnterRQ SubMazeEnergyRQ fail", zap.Error(err), zap.Any("req", req), zap.Any("res", res))
//		return
//	}
//	// 体力不足，返回错误码(80000 // 体力不足),并带回剩余的体力值
//	// 扣体力成功，返回剩余的体力值
//	if res.GetErrInfo().GetErrCode() == errors.NO_ERROR_CODE {
//		return true, res.GetRemainVal(), nil
//	} else if res.GetErrInfo().GetErrCode() == 80000 {
//		return false, res.GetRemainVal(), nil
//	} else {
//		err = errors.New(string(res.GetErrInfo().GetErrMsg()))
//		return false, res.GetRemainVal(), err
//	}
//}

func (g *Game) OnGetStorageInfoRQ_10529_10530(s *session.Session, req *MazeGame.MazeBarrierEnterRQ) (err error) {
	defer fkprometheus.InfoPMT("OnGetStorageInfoRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.GetStorageInfoRS{}

	logger.InfoWF("OnGetStorageInfoRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnGetStorageInfoRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.BarrierId = req.BarrierId

	userId := uint64(s.UID())

	if req.GetBarrierId() <= 0 {
		logger.ErrorWF("OnGetStorageInfoRQ req barrier invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id未设置")
		return
	}

	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnGetStorageInfoRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	// 默认是从存档进入
	storageInfo, err := syncmazestorageinforedis.GetSyncMazeStorageInfo(userId, userInfo.Barrier)
	if err != nil {
		logger.ErrorWF("OnGetStorageInfoRQ GetSyncMazeStorageInfo fail", zap.Error(err))
		return
	}
	res.StorageInfo = storageInfo
	//if storageInfo == nil {
	//	// 进入清临时buff
	//	mazebarriertempbuffredis.ClearBarrierTempBuff(logger, userId, req.GetBarrierId())
	//	mazebuffinforedis.DelMazeBuffBySrc(logger, userId, constdef.MazeBuffSrcSelectBuffForce)
	//	// 推送属性计算消息
	//	calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{
	//		UserId:  userId,
	//		ChgType: constdef.MazeBuffChgForceValue,
	//		Session: "buff",
	//		BuffSrc: constdef.MazeBuffSrcSelectBuffForce,
	//	}
	//	mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, calcAttrNotify)
	//	// 清理关卡操作状态
	//	mazebarrieropstatusredis.ClearOpStatus(logger, userId, req.GetBarrierId())
	//}

	return nil
}

func GetUserMoney(logger fklog.FKLogI, uid uint64) {

}
