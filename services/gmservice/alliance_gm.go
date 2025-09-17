package gmservice

import (
	"encoding/json"
	"fmt"
	"maze_game_server/model/alliancemodel"
	"maze_game_server/services/allianceservice"
	"net/http"
	"strings"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

func (s *service) CreateAlliance(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	name := request.Form.Get("name")
	allianceID, err := allianceservice.GlobalAllianceService.CreateAlliance(ctx, name)
	if err != nil {
		logger.CtxError(ctx, "CreateAlliance Gm Fail",
			zap.Error(err),
		)
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write([]byte(fmt.Sprintf("ok 联盟id:%d", allianceID)))
}

func (s *service) GetAllianceInfo(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)
	allianceID := fkutil.ToInt32(request.Form.Get("allianceID"))

	allianceInfoModel, err := alliancemodel.LoadAllianceInfoModel(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "GetAllianceInfo LoadAllianceInfoModel err",
			zap.Int32("allianceID", allianceID), zap.Error(err))
		writer.Write([]byte(err.Error()))
		return
	}

	jsondata, err := json.Marshal(allianceInfoModel)
	if err != nil {
		logger.CtxError(ctx, "GetAllianceInfo Marshal err",
			zap.Int32("allianceID", allianceID), zap.Error(err))
		writer.Write([]byte(err.Error()))
		return
	}

	writer.Write(jsondata)
}

func (s *service) ChangeAllianceInfo(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)
	name := request.Form.Get("name")
	allianceID := fkutil.ToInt32(request.Form.Get("alliance_id"))

	err := allianceservice.GlobalAllianceService.ChangeAllianceInfo(ctx, allianceID, name)
	if err != nil {
		logger.CtxError(ctx, "ChangeAllianceInfoGM  ChangeAllianceInfo Fail",
			zap.Error(err),
			zap.Any("name", name),
			zap.Any("allianceID", allianceID),
		)
		writer.Write([]byte(err.Error()))
		return
	}

	writer.Write([]byte("ok"))
}

func (s *service) BatchChangeAlliance(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	userIDs := request.Form.Get("user_ids")
	allianceID := fkutil.ToInt32(request.Form.Get("alliance_id"))

	var userId []uint64
	for _, str := range strings.Split(userIDs, ",") {
		if userID := fkutil.ToUint64(str); userID != 0 {
			userId = append(userId, userID)
		}
	}

	err := allianceservice.GlobalAllianceService.BatchChangeAlliance(ctx, userId, allianceID)
	if err != nil {
		logger.CtxError(ctx, "BatchChangeAllianceGM  BatchChangeAlliance Fail",
			zap.Error(err),
			zap.Any("userId", userId),
			zap.Any("userIDs", userIDs),
			zap.Any("allianceID", allianceID),
		)
		writer.Write([]byte(err.Error()))
		return
	}

	writer.Write([]byte("ok"))
}

func (s *service) GetAllianceList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)
	allianceList, err := allianceservice.GlobalAllianceService.QueryAllianceList(ctx)
	if err != nil {
		logger.CtxError(ctx, "GetAllianceListGM QueryAllianceList Fail")
		writer.Write([]byte(err.Error()))
		return
	}

	for _, alallianceID := range allianceList.Alliances {
		writer.Write([]byte(fmt.Sprintf("联盟id:%d\n", alallianceID)))
	}
}

func (s *service) GetUserAlliance(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	userID := fkutil.ToUint64(request.Form.Get("user_id"))
	allianceID, err := allianceservice.GlobalAllianceService.QueryUserAlliance(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "GetUserAlliance QueryUserAlliance Fail",
			zap.Error(err),
		)
		writer.Write([]byte(err.Error()))
	}

	writer.Write([]byte(fmt.Sprintf("%d", allianceID)))
}
