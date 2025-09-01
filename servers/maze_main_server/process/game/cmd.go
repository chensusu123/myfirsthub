package game

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/io/kafka/mazeenergyrecord"
	"maze_game_server/io/redis/mazeboxredis"
	"maze_game_server/io/redis/mazechallengenumredis"
	"maze_game_server/io/redis/mazecollectredis"
	"maze_game_server/io/redis/mazeequipgetnumredis"

	"maze_game_server/io/redis/mazeuserbarrierredis"
	"maze_game_server/io/redis/mazeuserlevelredis"
	"maze_game_server/io/redis/syncmazestorageinforedis"
	"maze_game_server/lib/nano/session"

	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/services/barrierenergyservice"
	"maze_game_server/services/barrierscorerewardservice"
	"maze_game_server/services/equipdropservice"
	"maze_game_server/services/moneyservice"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

func (g *Game) OnSendDollMazeCmdRQ_10463_10464(s *session.Session, req *MazeGame.SendDollMazeCmdRQ) (err error) {
	defer fkprometheus.DebugPMT("OnSendDollMazeCmdRQ")()

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeGame.SendDollMazeCmdRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.CmdCode = req.CmdCode
	res.CmdParam = req.CmdParam

	userId := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnSendDollMazeCmdRQ end", zap.Any("res", res))
	}()

	logger.CtxInfo(ctx, "OnSendDollMazeCmdRQ with", zap.Any("req", req))
	// if !BreedVersionFC.IsDollVersion(userCtx, shardingID) {
	// 	userCtx.WarnWF("OnSendDollMazeCmdRQ not doll version", zap.Uint64("userID", shardingID))
	// 	return
	// }
	codeS := req.GetCmdCode()
	code := fkutil.ToInt32(codeS)
	args := ParseCmdParam(req.GetCmdParam())
	switch code {
	// case 1001:
	// 	err = resetequipcmd.RunCmd1001(userCtx, shardingID, req.GetHeader().GetSession(), req.GetCmdParam())
	// 	if err != nil {
	// 		err = errors.New("执行失败")
	// 	}
	case 1002:
		err = ParseCmd(ctx, userId, code, req.GetCmdParam(), req.GetHeader().GetSession())
		if err != nil {
			err = errors.New("执行失败")
		}
	case 1003:
		err = ParseCmd(ctx, userId, code, req.GetCmdParam(), req.GetHeader().GetSession())
		if err != nil {
			err = errors.New("执行失败")
		}
	case 1004:
		err = ParseCmd(ctx, userId, code, req.GetCmdParam(), req.GetHeader().GetSession())
		if err != nil {
			err = errors.New("执行失败")
		}
	case 1005:
		err = ParseCmd(ctx, userId, code, req.GetCmdParam(), req.GetHeader().GetSession())
		if err != nil {
			err = errors.New("执行失败")
		}
	case 1006: // 添加体力
		err = CmdAddEnergy(ctx, userId, args)
		if err != nil {
			err = errors.New("执行失败")
		}
	default:
		err = errors.New("未知命令")
	}
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
	}
	return nil
}

func ParseCmd(ctx context.Context, uid uint64, cmdCode int32, cmd string, session string) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	switch cmdCode {
	case 1002:
		cmdParams := strings.Split(cmd, "&")
		if len(cmdParams) != 3 {
			err = errors.New("cmd_param设置错误")
			return
		}
		params := make(map[string]string)
		for _, p := range cmdParams {
			datas := strings.Split(p, "=")
			if len(datas) != 2 {
				err = errors.New("cmd_param设置错误")
				return
			}
			params[datas[0]] = datas[1]
		}
		logger.CtxInfo(ctx, "ParseCmd SetMazeMoney dump params", zap.Any("params", params))
		err = SetMazeMoney(ctx, fkutil.ToUint64(params["user"]), fkutil.ToInt32(params["diamond"]), fkutil.ToInt64(params["money"]), session)
		return
	case 1003:
		cmdParams := strings.Split(cmd, "&")
		if len(cmdParams) != 2 {
			err = errors.New("cmd_param设置错误")
			return
		}
		params := make(map[string]string)
		for _, p := range cmdParams {
			datas := strings.Split(p, "=")
			if len(datas) != 2 {
				err = errors.New("cmd_param设置错误")
				return
			}
			params[datas[0]] = datas[1]
		}
		logger.CtxInfo(ctx, "ParseCmd SetMazeBarrier dump params", zap.Any("params", params))
		err = SetMazeBarrier(ctx, fkutil.ToUint64(params["user"]), fkutil.ToInt32(params["barrier"]))
		return
	case 1004:
		cmdParams := strings.Split(cmd, "&")
		if len(cmdParams) != 1 {
			err = errors.New("cmd_param设置错误")
			return
		}
		params := make(map[string]string)
		for _, p := range cmdParams {
			datas := strings.Split(p, "=")
			if len(datas) != 2 {
				err = errors.New("cmd_param设置错误")
				return
			}
			params[datas[0]] = datas[1]
		}
		logger.CtxInfo(ctx, "ParseCmd ClearBarrier dump params", zap.Any("params", params))
		err = mazeboxredis.ClearOpenBoxTime(ctx, fkutil.ToUint64(params["user"]))
		if err != nil {
			logger.CtxError(ctx, "ParseCmd ClearOpenBoxTime fail", zap.Error(err), zap.Uint64("userID", fkutil.ToUint64(params["user"])))
		}
		err = ClearBarrier(ctx, fkutil.ToUint64(params["user"]))
		if err != nil {
			logger.CtxError(ctx, "ParseCmd End ClearOpenBoxTime failed", zap.Error(err), zap.Uint64("userID", fkutil.ToUint64(params["user"])))
		}
		return
	case 1005:
		cmdParams := strings.Split(cmd, "&")
		if len(cmdParams) != 3 {
			err = errors.New("cmd_param设置错误")
			return
		}
		params := make(map[string]string)
		for _, p := range cmdParams {
			datas := strings.Split(p, "=")
			if len(datas) != 2 {
				err = errors.New("cmd_param设置错误")
				return
			}
			params[datas[0]] = datas[1]
		}
		logger.CtxInfo(ctx, "ParseCmd SetMazeUserInfo dump params", zap.Any("params", params))
		err = SetMazeUserInfo(ctx, fkutil.ToUint64(params["user"]), fkutil.ToInt64(params["level"]), fkutil.ToInt64(params["exp"]))
		return
	default:
		err = errors.New("wrong cmd_code")
	}
	return
}

func SetMazeMoney(ctx context.Context, uid uint64, diamond int32, money int64, session string) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	if uid <= 0 || diamond < 0 || money < 0 {
		logger.CtxError(ctx, "userId不能小于等于0,diamond、moneyCount不能小于0")
		err = errors.New("userId不能小于等于0,diamond、moneyCount不能小于0")
		return
	}
	err = moneyservice.GlobalMoneyService.SetMoney(context.TODO(), uid, constdef.MazeCommonItemCoin, money)
	if err != nil {
		logger.CtxError(ctx, "SetMazeMoney GMSet fail", zap.Error(err))
		return
	}
	err = moneyservice.GlobalMoneyService.SetMoney(context.TODO(), uid, constdef.MazeCommonItemDiamond, int64(diamond))
	if err != nil {
		logger.CtxError(ctx, "SetMazeMoney GMSet fail", zap.Error(err))
		return
	}

	return
}

func SetMazeBarrier(ctx context.Context, userId uint64, barrierId int32) (err error) {

	if userId <= 0 || barrierId <= 0 {
		err = errors.New("userId、barrierId不能小于等于0")
		return
	}

	//设置关卡升级
	// err = mazebarrierredis.SetBarrier(logger, uint64(userId), int32(barrierId))
	// if err != nil {
	// 	return
	// }

	return
}

func ClearBarrier(ctx context.Context, userId uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	if userId <= 0 {
		err = errors.New("userId不能小于等于0")
		return
	}

	//清楚关卡信息
	// err = mazebarrierredis.GMDel(logger, uint64(userId))
	// if err != nil {
	// 	return
	// }

	//清除关卡存档
	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		return
	}
	err = syncmazestorageinforedis.DelSyncMazeStorageInfo(userId, userInfo.Barrier)
	logger.CtxInfo(ctx, "ClearBarrier end", zap.Error(err), zap.Any("userInfo", userInfo), zap.Any("user", userId))
	if err != nil {
		return err
	}

	//清除等级经验通用数值
	err = mazeuserlevelredis.GMDel(ctx, userId)
	if err != nil {
		return
	}

	// maxBarrierId := 7
	// for barrierId := 1; barrierId <= maxBarrierId; barrierId++ {
	err = mazeuserbarrierredis.GMDel(ctx, userId, 0)
	if err != nil {
		return
	}
	// }
	//err = mazeshopseqredis.GMDel(logger, userId)
	//if err != nil {
	//	return
	//}

	//err = mazeequipspecialdropredis.GMDel(ctx, userId)
	//if err != nil {
	//	return err
	//}

	err = equipdropservice.GlobalEquipDropService.GmDelete(ctx, userId)
	if err != nil {
		return err
	}

	err = mazecollectredis.GMDel(ctx, userId)
	if err != nil {
		return
	}

	now := time.Now()
	today := now.Year()*10000 + int(now.Month())*100 + now.Day()
	err = mazechallengenumredis.GMDel(ctx, userId, today)
	if err != nil {
		return
	}

	err = mazeequipgetnumredis.GMDel(ctx, userId)
	if err != nil {
		return
	}

	//清除关卡已获得奖励存档
	err = barrierscorerewardservice.GlobalScoreRewardService.DelBarrierScoreRewardItem(context.TODO(), userId, userInfo.Barrier)
	if err != nil {
		return
	}

	//重置体力
	err = barrierenergyservice.GlobalBarrierEnergyService.ResetEnergy(ctx, userId)
	if err != nil {
		return
	}

	ClearBarriersTempData(ctx, userId, userInfo.Barrier)
	return

}

func SetMazeUserInfo(ctx context.Context, uid uint64, level int64, exp int64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	if uid <= 0 || level <= 0 || exp < 0 {
		logger.CtxError(ctx, "userId、level不能小于等于0,exp不能小于0")
		err = errors.New("userId、level不能小于等于0,exp不能小于0")
		return
	}

	// userInfo, err := mazeuserinfo.GetUserInfoV2(logger, uid)
	// if err != nil {
	// 	return
	// }
	// userInfo.Level = int64(level)
	// userInfo.Exp = exp

	// levelCfg := GMazeLevelV8Cfg.GetWithCtx(ctx,int32(level))
	// if levelCfg == nil {
	// 	logger.CtxError(ctx,"根据level不能读取等级配置表")
	// 	err = errors.New("根据level不能读取等级配置表")
	// 	return
	// }
	// userInfo.TotalExp = levelCfg.All_exp - levelCfg.Next_level_need_exp + exp

	// err = mazeuserinfo.SetUserInfo(logger, uint64(uid), userInfo)
	// if err != nil {
	// 	return
	// }

	return
}

func ParseCmdParam(p string) map[string]string {
	args := make(map[string]string)
	ss := strings.Split(p, "&")
	for i := 0; i < len(ss); i++ {
		if ss[i] == "" {
			continue
		}
		sss := strings.Split(ss[i], "=")
		if len(sss) == 2 {
			args[sss[0]] = sss[1]
		}
	}
	return args
}

func CmdAddEnergy(ctx context.Context, userId uint64, args map[string]string) error {
	logger := fklog.ContextAppLogger(ctx)
	var vInt int32
	if v, ok := args["add_cnt"]; ok {
		vInt = fkutil.ToInt32(v)
	}
	if vInt <= 0 {
		logger.CtxWarn(ctx, "CmdAddEnergy vInt=0", zap.Any("args", args))
		return errors.New("加体力参数错误")
	}
	oldEnergy, _, err := barrierenergyservice.GlobalBarrierEnergyService.GetBarrierEnergy(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "CmdAddEnergy GetBarrierEnergy failed", zap.Error(err))
		return err
	}

	newEnergy, nextUpdateTime, err := barrierenergyservice.GlobalBarrierEnergyService.AddEnergy(ctx, userId, vInt)
	if err != nil {
		logger.CtxError(ctx, "CmdAddEnergy AddEnergy failed", zap.Error(err))
		return err
	}
	barrierenergyservice.GlobalBarrierEnergyService.PushEnergyRecord(ctx, userId, oldEnergy, newEnergy, mazeenergyrecord.GMAdd, nextUpdateTime)

	//rq := &MazeEnergySvr.AddMazeEnergyRQ{}
	//rs := &MazeEnergySvr.AddMazeEnergyRS{}
	//rq.UserId = proto.Uint64(userId)
	//rq.AddVal = proto.Int32(vInt)
	//rq.OpType = proto.Int32(int32(MazeEnergySvr.ENUM_MAZE_ENERGY_OP_TYPE_GMADD))
	//rq.OpDesc = proto.String("CmdGmAdd")
	//rq.TradeNumber = proto.Uint64(uniqueid.GenUniqueIdUInt64())
	// 合并服务，直接访问函数
	// return mazeenergyrpc.AddMazeEnergyRQ(logger, rq, rs)
	//return energy.AddMazeEnergyRQ(logger, userId, rq, rs)
	return err
}
