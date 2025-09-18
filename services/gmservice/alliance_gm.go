package gmservice

import (
	"encoding/json"
	"fmt"
	"maze_game_server/model/alliancemodel"
	"maze_game_server/model/gmmodel"
	"maze_game_server/services/allianceservice"
	"net/http"
	"strings"

	"github.com/iancoleman/orderedmap"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

func (s *service) CreateAlliance(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	var outPut gmmodel.Output
	defer func() {
		jsonOut, err := json.Marshal(outPut)
		if err != nil {
			logger.CtxError(ctx, "Post: /CreateAlliance  Marshal Fail",
				zap.Any("request", request),
				zap.Any("ouput", outPut),
				zap.Error(err),
			)
		}
		writer.Write(jsonOut)
	}()

	name := request.Form.Get("name")
	allianceID, err := allianceservice.GlobalAllianceService.CreateAlliance(ctx, name)
	if err != nil {
		logger.CtxError(ctx, "CreateAlliance Gm Fail",
			zap.Error(err),
		)
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, err.Error(), gmmodel.DynamicData{})
		return
	}
	outPut = *gmmodel.NewOutPut(http.StatusOK, fmt.Sprintf("ok 联盟id:%d", allianceID), gmmodel.DynamicData{})
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

	var outPut gmmodel.Output
	defer func() {
		jsonOut, err := json.Marshal(outPut)
		if err != nil {
			logger.CtxError(ctx, "Post: /ChangeAllianceInfo  Marshal Fail",
				zap.Any("request", request),
				zap.Any("ouput", outPut),
				zap.Error(err),
			)
		}
		writer.Write(jsonOut)
	}()

	err := allianceservice.GlobalAllianceService.ChangeAllianceInfo(ctx, allianceID, name)
	if err != nil {
		logger.CtxError(ctx, "ChangeAllianceInfoGM  ChangeAllianceInfo Fail",
			zap.Error(err),
			zap.Any("name", name),
			zap.Any("allianceID", allianceID),
		)
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, err.Error(), gmmodel.DynamicData{})
		return
	}

	outPut = *gmmodel.NewOutPut(http.StatusOK, "ok", gmmodel.DynamicData{})
}

func (s *service) BatchChangeAlliance(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	page := request.Form.Get("page")
	count := request.Form.Get("count")

	_ = page
	_ = count

	var outPut gmmodel.Output
	defer func() {
		jsonOut, err := json.Marshal(outPut)
		if err != nil {
			logger.CtxError(ctx, "Post: /ChangeAllianceInfo  Marshal Fail",
				zap.Any("request", request),
				zap.Any("ouput", outPut),
				zap.Error(err),
			)
		}
		writer.Write(jsonOut)
	}()

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
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, err.Error(), gmmodel.DynamicData{})
		return
	}

	outPut = *gmmodel.NewOutPut(http.StatusOK, "ok", gmmodel.DynamicData{})
}

func (s *service) GetAllianceList(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	var outPut gmmodel.Output
	defer func() {
		jsonOut, err := json.Marshal(outPut)
		if err != nil {
			logger.CtxError(ctx, "Post: /GetAllianceList  Marshal Fail",
				zap.Any("request", request),
				zap.Any("ouput", outPut),
				zap.Error(err),
			)
		}
		writer.Write(jsonOut)
	}()

	allianceList, err := allianceservice.GlobalAllianceService.QueryAllianceList(ctx)
	if err != nil {
		logger.CtxError(ctx, "GetAllianceListGM QueryAllianceList Fail")
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, err.Error(), gmmodel.DynamicData{})
		return
	}

	var records []*orderedmap.OrderedMap
	for _, allianceID := range allianceList.Alliances {

		allianceInfoModel, err := alliancemodel.LoadAllianceInfoModel(ctx, allianceID)
		if err != nil {
			logger.CtxError(ctx, "GetAllianceListGM LoadAllianceInfoModel err",
				zap.Int32("allianceID", allianceID), zap.Error(err))
			outPut = *gmmodel.NewOutPut(http.StatusBadGateway, err.Error(), gmmodel.DynamicData{})
			return
		}

		descRecord := orderedmap.New()
		descRecord.Set("alliance_id", allianceInfoModel.AllianceID)
		descRecord.Set("alliance_name", allianceInfoModel.AllianceName)
		descRecord.Set("create_at", allianceInfoModel.CreateTime)
		records = append(records, descRecord)
	}

	outPut = *gmmodel.NewOutPut(http.StatusOK, "操作成功", *gmmodel.NewDynamicData(records, len(records)))
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
