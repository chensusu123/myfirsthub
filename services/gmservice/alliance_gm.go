package gmservice

import (
	"encoding/json"
	"maze_game_server/services/allianceservice"
	"net/http"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

func (s *service) CreateAlliance(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	name := request.Form.Get("name")
	err := allianceservice.GlobalAllianceService.CreateAlliance(ctx, name)
	if err != nil {
		logger.CtxError(ctx, "CreateAlliance Gm Fail",
			zap.Error(err),
		)
		writer.Write([]byte(err.Error()))
	}
	writer.Write([]byte("ok"))
}

func (s *service) GetAllianceInfo(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)
	allianceID := fkutil.ToInt32(request.Form.Get("allianceID"))

	info, err := allianceservice.GlobalAllianceService.QueryAllianceInfo(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "GetAllianceInfo Gm Fail",
			zap.Error(err),
		)
		writer.Write([]byte(err.Error()))
	}

	jsonInfo, err := json.Marshal(info)
	if err != nil {
		logger.CtxError(ctx, "GetAllianceInfo Marshal Fail",
			zap.Int32("allianceID", allianceID),
			zap.Error(err),
		)
		writer.Write([]byte(err.Error()))
		return
	}

	writer.Write(jsonInfo)
}
