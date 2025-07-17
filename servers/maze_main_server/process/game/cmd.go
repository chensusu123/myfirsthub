package game

import (
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/uniqueid"
	"maze_game_server/io/redis/mazebarriermoneyredis"
	"maze_game_server/io/redis/mazeboxredis"
	"maze_game_server/io/redis/mazechallengenumredis"
	"maze_game_server/io/redis/mazecollectredis"
	"maze_game_server/io/redis/mazeequipgetnumredis"
	"maze_game_server/io/redis/mazeshopseqredis"
	"maze_game_server/io/redis/mazeuserbarrierredis"
	"maze_game_server/io/redis/mazeuserlevelredis"
	"maze_game_server/io/redis/syncmazestorageinforedis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/server/MazeEnergySvr"
	"maze_game_server/servers/maze_main_server/process/game/energy"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (g *Game) OnSendDollMazeCmdRQ_10463_10464(s *session.Session, req *MazeGame.SendDollMazeCmdRQ) (err error) {
	defer fkprometheus.DebugPMT("OnSendDollMazeCmdRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.SendDollMazeCmdRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.CmdCode = req.CmdCode
	res.CmdParam = req.CmdParam

	userId := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnSendDollMazeCmdRQ end", zap.Any("res", res))
	}()

	logger.InfoWF("OnSendDollMazeCmdRQ with", zap.Any("req", req))
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
		err = ParseCmd(logger, userId, code, req.GetCmdParam(), req.GetHeader().GetSession())
		if err != nil {
			err = errors.New("执行失败")
		}
	case 1003:
		err = ParseCmd(logger, userId, code, req.GetCmdParam(), req.GetHeader().GetSession())
		if err != nil {
			err = errors.New("执行失败")
		}
	case 1004:
		err = ParseCmd(logger, userId, code, req.GetCmdParam(), req.GetHeader().GetSession())
		if err != nil {
			err = errors.New("执行失败")
		}
	case 1005:
		err = ParseCmd(logger, userId, code, req.GetCmdParam(), req.GetHeader().GetSession())
		if err != nil {
			err = errors.New("执行失败")
		}
	case 1006: // 添加体力
		err = CmdAddEnergy(logger, userId, args)
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

func ParseCmd(logger fklog.FKLogI, uid uint64, cmdCode int32, cmd string, session string) (err error) {

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
		logger.InfoWF("ParseCmd SetMazeMoney dump params", zap.Any("params", params))
		err = SetMazeMoney(logger, fkutil.ToUint64(params["user"]), fkutil.ToInt32(params["diamond"]), fkutil.ToInt64(params["money"]), session)
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
		logger.InfoWF("ParseCmd SetMazeBarrier dump params", zap.Any("params", params))
		err = SetMazeBarrier(logger, fkutil.ToUint64(params["user"]), fkutil.ToInt32(params["barrier"]))
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
		logger.InfoWF("ParseCmd ClearBarrier dump params", zap.Any("params", params))
		err = mazeboxredis.ClearOpenBoxTime(logger, fkutil.ToUint64(params["user"]))
		if err != nil {
			logger.ErrorWF("ParseCmd ClearOpenBoxTime fail", zap.Error(err), zap.Uint64("userID", fkutil.ToUint64(params["user"])))
		}
		err = ClearBarrier(logger, fkutil.ToUint64(params["user"]))

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
		logger.InfoWF("ParseCmd SetMazeUserInfo dump params", zap.Any("params", params))
		err = SetMazeUserInfo(logger, fkutil.ToUint64(params["user"]), fkutil.ToInt64(params["level"]), fkutil.ToInt64(params["exp"]))
		return
	default:
		err = errors.New("wrong cmd_code")
	}
	return
}

func SetMazeMoney(logger fklog.FKLogI, uid uint64, diamond int32, money int64, session string) (err error) {

	if uid <= 0 || diamond < 0 || money < 0 {
		logger.ErrorWF("userId不能小于等于0,diamond、moneyCount不能小于0")
		err = errors.New("userId不能小于等于0,diamond、moneyCount不能小于0")
		return
	}

	err = mazebarriermoneyredis.GMSet(logger, uint64(uid), int32(constdef.MazeCommonItemCoin), money)
	if err != nil {
		logger.ErrorWF("SetMazeMoney GMSet fail", zap.Error(err))
		return
	}

	err = mazebarriermoneyredis.GMSet(logger, uint64(uid), int32(constdef.MazeCommonItemDiamond), money)
	if err != nil {
		logger.ErrorWF("SetMazeMoney GMSet fail", zap.Error(err))
		return
	}

	return
}

func SetMazeBarrier(logger fklog.FKLogI, userId uint64, barrierId int32) (err error) {

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

func ClearBarrier(logger fklog.FKLogI, userId uint64) (err error) {
	if userId <= 0 {
		err = errors.New("userId不能小于等于0")
		return
	}

	//清楚关卡信息
	// err = mazebarrierredis.GMDel(logger, uint64(userId))
	// if err != nil {
	// 	return
	// }

	//清除等级经验通用数值
	err = mazeuserlevelredis.GMDel(logger, userId)
	if err != nil {
		return
	}

	// maxBarrierId := 7
	// for barrierId := 1; barrierId <= maxBarrierId; barrierId++ {
	err = mazeuserbarrierredis.GMDel(logger, userId, 0)
	if err != nil {
		return
	}
	// }
	err = mazeshopseqredis.GMDel(logger, userId)
	if err != nil {
		return
	}

	err = mazecollectredis.GMDel(logger, userId)
	if err != nil {
		return
	}

	now := time.Now()
	today := now.Year()*10000 + int(now.Month())*100 + now.Day()
	err = mazechallengenumredis.GMDel(logger, userId, today)
	if err != nil {
		return
	}

	err = mazeequipgetnumredis.GMDel(logger, userId)
	if err != nil {
		return
	}

	//清除关卡存档
	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		return
	}
	err = syncmazestorageinforedis.DelSyncMazeStorageInfo(userId, userInfo.Barrier)
	if err != nil {
		return err
	}

	return

}

func SetMazeUserInfo(logger fklog.FKLogI, uid uint64, level int64, exp int64) (err error) {

	if uid <= 0 || level <= 0 || exp < 0 {
		logger.ErrorWF("userId、level不能小于等于0,exp不能小于0")
		err = errors.New("userId、level不能小于等于0,exp不能小于0")
		return
	}

	// userInfo, err := mazeuserinfo.GetUserInfoV2(logger, uid)
	// if err != nil {
	// 	return
	// }
	// userInfo.Level = int64(level)
	// userInfo.Exp = exp

	// levelCfg := GMazeLevelV8Cfg.Get(int32(level))
	// if levelCfg == nil {
	// 	logger.ErrorWF("根据level不能读取等级配置表")
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

func CmdAddEnergy(logger fklog.FKLogI, userId uint64, args map[string]string) error {
	var vInt int32
	if v, ok := args["add_cnt"]; ok {
		vInt = fkutil.ToInt32(v)
	}
	if vInt <= 0 {
		logger.WarnWF("CmdAddEnergy vInt=0", zap.Any("args", args))
		return errors.New("加体力参数错误")
	}
	rq := &MazeEnergySvr.AddMazeEnergyRQ{}
	rs := &MazeEnergySvr.AddMazeEnergyRS{}
	rq.UserId = proto.Uint64(userId)
	rq.AddVal = proto.Int32(vInt)
	rq.OpType = proto.Int32(int32(MazeEnergySvr.ENUM_MAZE_ENERGY_OP_TYPE_GMADD))
	rq.OpDesc = proto.String("CmdGmAdd")
	rq.TradeNumber = proto.Uint64(uniqueid.GenUniqueIdUInt64())
	// 合并服务，直接访问函数
	// return mazeenergyrpc.AddMazeEnergyRQ(logger, rq, rs)
	return energy.AddMazeEnergyRQ(logger, userId, rq, rs)
}
