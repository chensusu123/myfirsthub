package gm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"maze_game_server/common/function/gm"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/redis/UnionIDBindRedis"
	"maze_game_server/io/redis/mazefixedbarrierredis"
	"maze_game_server/io/redis/useridredis"
	"maze_game_server/io/redis/usersection"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/servers/maze_main_server/process/gm/cmdbattledata"
	"maze_game_server/usecase/online"

	"github.com/gorilla/schema"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

var form = schema.NewDecoder()

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

		// 更新等级经验
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

	type SetBarrierParams struct {
		UserID    uint64 `schema:"userId,required"`
		BarrierID int32  `schema:"barrierId,required"`
		Lock      int32  `schema:"lock"` // 如果提供这个参数，则设置的关卡会记录下来，重置游戏数据也会继续生效
	}

	gm.SafeHttpRegister(logger, "/SetBarrier", func(writer http.ResponseWriter, request *http.Request) {
		var params SetBarrierParams

		err := form.Decode(&params, request.Form)
		if err != nil {
			writer.Write([]byte("参数不正确"))
			return
		}

		logger.SetLogId(time.Now().UnixNano())

		if params.UserID <= 0 || params.BarrierID <= 0 {
			writer.Write([]byte("参数不正确"))
			return
		}

		userInfo, err := mazeuserinfo.GetUserInfoV2(logger, params.UserID)
		if err != nil {
			logger.ErrorWF("SetBarrier GetUserInfoV2 fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}

		barrierCfg := GMazeBarriesV8Cfg.Get(params.BarrierID)
		if barrierCfg == nil {
			errCfg := errors.New("cant find barrier cfg")
			logger.ErrorWF("SetBarrier Get barrier fail", zap.Error(err))
			writer.Write([]byte(errCfg.Error()))
			return
		}

		// oldBarrier := userInfo.Barrier

		// if params.BarrierID <= oldBarrier {
		// 	fmt.Fprintf(writer, "仅支持跳过关卡，当前第%d关", oldBarrier)
		// 	return
		// }

		userInfo.SetBarrier(params.BarrierID)
		userInfo.SetPassBarrier(params.BarrierID - 1)

		// 更新设置关卡
		err = mazeuserinfo.SetUserInfoV2(logger, params.UserID, userInfo)
		if err != nil {
			logger.ErrorWF("SetBarrier SetUserInfoV2 fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}

		// 锁定设置的关卡
		if params.Lock == 1 {
			mazefixedbarrierredis.SetUserFixedBarrierID(logger, params.UserID, params.BarrierID)
		} else {
			mazefixedbarrierredis.DelUserFixedBarrierID(logger, params.UserID)
		}

		writer.Write([]byte("设置成功，注意尽量不要在迷宫杀怪时使用本gm"))
	})

	gm.SafeHttpRegister(logger, "/generateUser", func(writer http.ResponseWriter, request *http.Request) {
		logger.SetLogId(time.Now().UnixNano())
		AuthId := fkutil.ToUint64(request.Form.Get("AuthId"))

		userID := uint64(0)

		generateUser := &GenerateUser{}
		defer func() {
			generateUser.UserId = userID
			jsonData, err := json.Marshal(generateUser)
			if err != nil {
				writer.Write([]byte(err.Error()))
				return
			}
			writer.Write(jsonData)
		}()

		if AuthId == 0 {
			generateUser.ErrorCode = 1
			generateUser.ErrorMsg = "AuthId is 0"
			return
		}

		users, err := UnionIDBindRedis.GetUsersWithUnionID(logger, uint64(AuthId))
		if err != nil {
			generateUser.ErrorCode = 1
			generateUser.ErrorMsg = err.Error()
			return
		}

		if len(users) == 0 {
			newUserID := useridredis.Generate(logger)
			if newUserID == 0 {
				generateUser.ErrorCode = 1
				generateUser.ErrorMsg = "Generate error"
				return
			}
			err = UnionIDBindRedis.AddUnionID2UserID(logger, uint64(AuthId), newUserID)
			if err != nil {
				generateUser.ErrorCode = 1
				generateUser.ErrorMsg = err.Error()
				return
			}
			err = UnionIDBindRedis.AddUserID2UnionID(logger, newUserID, uint64(AuthId))
			if err != nil {
				generateUser.ErrorCode = 1
				generateUser.ErrorMsg = err.Error()
				return
			}
			err = usersection.Set(context.TODO(), newUserID, appconfig.GlobalConfig().Global.SectionID)
			if err != nil {
				logger.ErrorWF("usersection.Set fail",
					zap.Uint64("userID", userID),
					zap.Error(err))
			}
			userID = newUserID
		} else {
			userID = users[0]
		}
		generateUser.ErrorCode = 0
		generateUser.ErrorMsg = "success"
	})

	gm.SafeHttpRegister(logger, "/showSheet", func(writer http.ResponseWriter, request *http.Request) {
		logger.SetLogId(time.Now().UnixNano())
		showSheet := &ShowSheet{}
		defer func() {
			jsonData, err := json.Marshal(showSheet)
			if err != nil {
				writer.Write([]byte(err.Error()))
				return
			}
			writer.Write(jsonData)
		}()
		showSheet.Data = config_manager.ShowSheet()
		showSheet.ErrorCode = 0
		showSheet.ErrorMsg = "success"
	})

	gm.SafeHttpRegister(logger, "/online", func(writer http.ResponseWriter, request *http.Request) {
		logger.SetLogId(time.Now().UnixNano())
		fmt.Fprintf(writer, "会话ID 用户ID 客户端地址\n")
		online.Scan(func(id int64, s *session.Session) {
			s.RLock()
			defer s.RUnlock()
			if userID := s.UID(); userID <= 0 {
				fmt.Fprintf(writer, "%d 验证中 %s\n", s.ID(), s.RemoteAddr().String())
			} else {
				fmt.Fprintf(writer, "%d %d %s\n", s.ID(), userID, s.RemoteAddr().String())
			}
		})
	})
}

type ShowSheet struct {
	ErrorCode uint64                          `json:"errorCode"`
	ErrorMsg  string                          `json:"errorMsg"`
	Data      []config_manager.ConfigShowItem `json:"data"`
}
type GenerateUser struct {
	ErrorCode uint64 `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	UserId    uint64 `json:"userId"`
}
