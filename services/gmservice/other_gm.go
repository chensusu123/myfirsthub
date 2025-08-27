package gmservice

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/common/function/fileio"
	"maze_game_server/io/redis/UnionIDBindRedis"
	"maze_game_server/io/redis/dollassemblesuitredis"
	"maze_game_server/io/redis/mazebagequipredis"
	"maze_game_server/io/redis/useridredis"
	"maze_game_server/io/redis/usersection"
	"maze_game_server/lib/nano/session"
	"maze_game_server/servers/maze_main_server/process/equip_gm/equipbaggm"
	"maze_game_server/usecase/online"
	"net/http"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"gitlab.ifreetalk.com/nano-ecosystem/nlog/fkfmt"
	"go.uber.org/zap"
)

func (s *service) ClearBag(writer http.ResponseWriter, request *http.Request) {
	// 外网线上环境不允许使用GM
	ctx := request.Context()
	request.ParseForm()

	userId := fkutil.ToUint64(request.Form.Get("user_id"))

	err := equipbaggm.ClearUserBag(ctx, userId)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	writer.Write([]byte("ok"))

	return
}

func (s *service) BatchClearBag(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)
	// 外网线上环境不允许使用GM
	request.ParseForm()

	filename := request.Form.Get("file")

	fr := fileio.NewDefFReaderEx(logger, ",")
	err := fr.Open(filename)
	if err != nil {
		return
	}
	defer fr.Close()

	fr.SetDumpRow(5000)
	// fr.SetSleep(int32(waitLine), int32(sleep))
	fkfmt.Println("open file", filename, "succ")
	defer fkutil.CaptureException()

	fr.InitAsync(int(8), int(100))

	fr.Range(func(logger fklog.FKLogI, line []uint64) bool {
		if len(line) != 1 {
			logger.ErrorWF("file line not match")
			return false
		}
		userId := line[0]

		err := equipbaggm.ClearUserBag(context.TODO(), userId)
		if err != nil {
			logger.ErrorWF("BatchClearBag ClearUserBag fail", zap.Error(err), zap.Uint64("uid", userId))
			return false
		}

		return true
	})

	writer.Write([]byte("ok"))

	return
}

func (s *service) ClearBagNotAssemble(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)
	// 外网线上环境不允许使用GM
	request.ParseForm()

	uid := fkutil.ToUint64(request.Form.Get("user_id"))

	if uid == 0 {
		writer.Write([]byte("uid 不能为0"))
		return
	}

	// 获取身上的装备信息
	assembleInfoMap, err := dollassemblesuitredis.GetAllDollAssembleSuit(logger, uid)
	if err != nil {
		logger.ErrorWF("ClearBagNotAssemble GetAllDollAssembleSuit fail", zap.Error(err))
		writer.Write([]byte(err.Error()))
		return
	}
	// 身上的装备
	assembleGuid := make(map[int64]struct{})
	if len(assembleInfoMap) > 0 {
		for _, equipList := range assembleInfoMap {
			if len(equipList) > 0 {
				for _, v := range equipList {
					assembleGuid[v.GetEquipGuid()] = struct{}{}
				}
			}
		}
	}

	equipInfoMap, err := mazebagequipredis.GetAllEquipInfo(logger, uid)
	if err != nil {
		logger.ErrorWF("ClearBagNotAssemble GetAllEquipInfo fail", zap.Error(err))
		writer.Write([]byte(err.Error()))
		return
	}

	equipGuids := make([]int64, 0)
	for k := range equipInfoMap {
		_, ok := assembleGuid[k]
		if !ok {
			equipGuids = append(equipGuids, k)
		}
	}

	err = equipbaggm.ClearEquipBagBatch(logger, uid, equipGuids)
	if err != nil {
		logger.ErrorWF("ClearBagNotAssemble ClearEquipBagBatch fail", zap.Error(err))
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write([]byte("ok"))
	return
}

type GenerateUser struct {
	ErrorCode uint64 `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	UserId    uint64 `json:"userId"`
}

func (s *service) GenerateUser(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	AuthId := fkutil.ToUint64(request.Form.Get("user_id"))
	userID := uint64(0)
	callLogger := fklog.ContextAppLogger(ctx)
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

	users, err := UnionIDBindRedis.GetUsersWithUnionID(ctx, callLogger, uint64(AuthId))
	if err != nil {
		generateUser.ErrorCode = 1
		generateUser.ErrorMsg = err.Error()
		return
	}

	if len(users) == 0 {
		newUserID := useridredis.Generate(ctx, callLogger)
		if newUserID == 0 {
			generateUser.ErrorCode = 1
			generateUser.ErrorMsg = "Generate error"
			return
		}
		err = UnionIDBindRedis.AddUnionID2UserID(ctx, callLogger, uint64(AuthId), newUserID)
		if err != nil {
			generateUser.ErrorCode = 1
			generateUser.ErrorMsg = err.Error()
			return
		}
		err = UnionIDBindRedis.AddUserID2UnionID(ctx, callLogger, newUserID, uint64(AuthId))
		if err != nil {
			generateUser.ErrorCode = 1
			generateUser.ErrorMsg = err.Error()
			return
		}
		err = usersection.Set(ctx, newUserID, appconfig.GlobalConfig().Global.SectionID)
		if err != nil {
			callLogger.CtxError(ctx, "usersection.Set fail",
				zap.Uint64("userID", userID),
				zap.Error(err))
		}
		userID = newUserID
	} else {
		userID = users[0]
	}
	generateUser.ErrorCode = 0
	generateUser.ErrorMsg = "success"
}

func (s *service) Online(writer http.ResponseWriter, request *http.Request) {

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
}
