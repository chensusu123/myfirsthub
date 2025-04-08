package process

import (
	"time"

	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/gentradeno"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebarriertempbuffredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazechallengenumredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeuserbarrierredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/rpc/mazeenergyrpc"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/calequipsequence"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazecommonvalue"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeActionCountV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeBarriesV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeLevelV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/protodef/MazeAIBattle"
	"gitlab.ifreetalk.com/plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/plate/protodef/MazeEnergySvr"
	"gitlab.ifreetalk.com/plate/protodef/MazeGame"

	"go.uber.org/zap"
)

func OnMazeBarrierEnterRQ(logger fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	fkprometheus.InfoPMT("OnMazeBarrierEnterRQ")()

	req := rqMsg.(*MazeGame.MazeBarrierEnterRQ)
	res := rsMsg.(*MazeGame.MazeBarrierEnterRS)

	logger.InfoWF("OnMazeBarrierEnterRQ start", zap.Any("req", req))
	defer func() {
		logger.InfoWF("OnMazeBarrierEnterRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := shardingID

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

	if req.GetBarrierId() < userInfo.Barrier {
		logger.ErrorWF("OnMazeBarrierEnterRQ req barrier lt pass barrier", zap.Any("req", req), zap.Int32("save", userInfo.Barrier))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("该关卡id小于存储的关卡id")
		return
	}

	//	res.Energy = proto.Int32(userInfo.Energy)
	var isNewBarrier bool

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
	if req.GetBarrierId() > userInfo.Barrier {
		// 	isNewBarrier = true
		userInfo.SetBarrier(req.GetBarrierId())
	}

	shopInfo, err := calequipsequence.GetMazeShopInfo(logger, userId, int32(userInfo.Level), req.GetBarrierId())
	if err != nil {
		logger.ErrorWF("OnMazeBarrierEnterRQ GetMazeShopInfo fail", zap.Error(err))
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

		//扣次数
		var maxNum int32
		maxNumCfg := GMazeActionCountV8Cfg.Get(101)
		if maxNumCfg == nil {
			logger.ErrorWF("OnMazeBarrierEnterRQ GMazeActionCountV8Cfg fail", zap.Error(err))
			res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
			return
		}
		maxNum = maxNumCfg.Day_count_v8
		now := time.Now()
		today := now.Year()*10000 + int(now.Month())*100 + now.Day()

		useNumToday, err2 := mazechallengenumredis.GetUserChallengeNum(logger, userId, today)
		if err2 != nil {
			logger.ErrorWF("OnMazeBarrierEnterRQ GetUserChallengeNum fail", zap.Error(err2))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
		if int32(useNumToday)+barrierCfg.Challenge_cost > maxNum {
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("次数不足")
			return
		}
		err = mazechallengenumredis.AddUserChallengeNum(logger, userId, today, barrierCfg.Challenge_cost)
		if err != nil {
			logger.ErrorWF("OnMazeBarrierEnterRQ AddUserChallengeNum fail", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}

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

		//3. 首次进入清临时buff
		mazebarriertempbuffredis.ClearBarrierTempBuff(logger, userId, req.GetBarrierId())
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
		EquipPoint: proto.Int64(int64(shopInfo.EquipPoints)),
		Energy:     proto.Int32(curEnergy),
	}

	return nil
}

func SubUserEnergy(logger fklog.FKLogI, uid uint64, subEnergy int32) (isSucc bool, newEnergy int32, err error) {
	req := &MazeEnergySvr.SubMazeEnergyRQ{
		UserId:      proto.Uint64(uid),
		SubVal:      proto.Int32(subEnergy),
		OpType:      proto.Int32(1), //NUM_MAZE_ENERGY_OP_TYPE_CHALLLENGE
		OpDesc:      proto.String("maze_barrier_enter"),
		TradeNumber: proto.Uint64(gentradeno.GetTradeNum()),
	}
	res := &MazeEnergySvr.SubMazeEnergyRS{}
	err = mazeenergyrpc.SubMazeEnergyRQ(logger, req, res)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierEnterRQ SubMazeEnergyRQ fail", zap.Error(err), zap.Any("req", req), zap.Any("res", res))
		return
	}
	// 体力不足，返回错误码(80000 // 体力不足),并带回剩余的体力值
	// 扣体力成功，返回剩余的体力值
	if res.GetErrInfo().GetErrCode() == errors.NO_ERROR_CODE {
		return true, res.GetRemainVal(), nil
	} else if res.GetErrInfo().GetErrCode() == 80000 {
		return false, res.GetRemainVal(), nil
	} else {
		err = errors.New(string(res.GetErrInfo().GetErrMsg()))
		return false, res.GetRemainVal(), err
	}
}

func GetUserMoney(logger fklog.FKLogI, uid uint64) {

}
