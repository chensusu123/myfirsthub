package gm

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/gm"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeBarriesV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazeuserlevelkafka"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazecommonvalue"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/gm/cmdbattledata"
	"go.uber.org/zap"
)

func RegGm(logger fklog.FKLogI) {
	cmdbattledata.RegBattleDataGm(logger)

	gm.SafeHttpRegister(logger, "/AddExp", func(writer http.ResponseWriter, request *http.Request) {
		userId := fkutil.ToUint64(request.Form.Get("userId"))
		exp := fkutil.ToInt64(request.Form.Get("exp"))
		logger.SetLogId(time.Now().UnixNano())

		userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
		if err != nil {
			logger.ErrorWF("AddExp GetUserInfoV2 fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}

		oldLevel := userInfo.Level
		oldExp := userInfo.TotalExp

		//更新等级经验
		err = userInfo.AddExp(exp)
		if err != nil {
			logger.ErrorWF("AddExp CalExp fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}
		err = mazeuserinfo.SetUserInfoV2(logger, userId, userInfo)
		if err != nil {
			logger.ErrorWF("AddExp SetUserInfoV2 fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}
		mazecommonvalue.HandleUserLevelExpChg(logger, userId, userInfo.Level, userInfo.Exp, "")

		defer func() {
			if exp != 0 {
				levelRecord := &mazeuserlevelkafka.MazeUserLevelRecord{
					UserId:      userId,
					OldLevel:    int32(oldLevel),
					OldTotalExp: oldExp,
					NewLevel:    int32(userInfo.Level),
					NewTotalExp: int32(userInfo.TotalExp),
				}
				mazeuserlevelkafka.PushMazeLevelRecord(logger, levelRecord)
			}
		}()

		writer.Write([]byte("设置成功，注意尽量不要在迷宫杀怪时使用本gm"))
	})

	gm.SafeHttpRegister(logger, "/SetBarrier", func(writer http.ResponseWriter, request *http.Request) {
		userId := fkutil.ToUint64(request.Form.Get("userId"))
		barrierID := fkutil.ToInt32(request.Form.Get("barrierId"))
		logger.SetLogId(time.Now().UnixNano())

		if userId <= 0 || barrierID <= 0 {
			writer.Write([]byte("参数不正确"))
			return
		}

		userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
		if err != nil {
			logger.ErrorWF("SetBarrier GetUserInfoV2 fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}

		barrierCfg := GMazeBarriesV8Cfg.Get(barrierID)
		if barrierCfg == nil {
			errCfg := errors.New("cant find barrier cfg")
			logger.ErrorWF("SetBarrier Get barrier fail", zap.Error(err))
			writer.Write([]byte(errCfg.Error()))
			return
		}

		var (
			oldBarrier = userInfo.Barrier
		)

		if barrierID <= oldBarrier {
			fmt.Fprintf(writer, "仅支持跳过关卡，当前第%d关", oldBarrier)
			return
		}

		userInfo.SetBarrier(barrierID)
		userInfo.SetPassBarrier(barrierID - 1)

		//更新设置关卡
		err = mazeuserinfo.SetUserInfoV2(logger, userId, userInfo)
		if err != nil {
			logger.ErrorWF("SetBarrier SetUserInfoV2 fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}

		writer.Write([]byte("设置成功，注意尽量不要在迷宫杀怪时使用本gm"))
	})

}
