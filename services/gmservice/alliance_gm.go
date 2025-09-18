package gmservice

import (
	"encoding/json"
	"fmt"
	"maze_game_server/model/alliancemodel"
	"maze_game_server/model/gmmodel"
	"maze_game_server/services/allianceservice"
	"net/http"
	"sort"
	"strconv"
	"strings"

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

	pageSizeString := request.Form.Get("page_size")
	pageString := request.Form.Get("page")

	pageSize, err := strconv.ParseInt(pageSizeString, 10, 64)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	page, err := strconv.ParseInt(pageString, 10, 64)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	allianceList, err := allianceservice.GlobalAllianceService.QueryAllianceList(ctx)
	if err != nil {
		logger.CtxError(ctx, "GetAllianceListGM QueryAllianceList Fail")
		writer.Write([]byte(err.Error()))
		return
	}

	allianceListInfo := make([]*gmmodel.AllianceInfo, 0)
	for _, allianceID := range allianceList.Alliances {
		allianceInfoModel, err := alliancemodel.LoadAllianceInfoModel(ctx, allianceID)
		if err != nil {
			logger.CtxError(ctx, "GetAllianceListGM LoadAllianceInfoModel err",
				zap.Int32("allianceID", allianceID), zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}

		allianceListInfo = append(allianceListInfo, &gmmodel.AllianceInfo{
			AllianceID:   int(allianceInfoModel.AllianceGroupID),
			AllianceName: allianceInfoModel.AllianceName,
			CreateAt:     allianceInfoModel.CreateTime,
		})
	}
	total := len(allianceListInfo)

	sort.Slice(allianceListInfo, func(i, j int) bool {
		return allianceListInfo[i].CreateAt < allianceListInfo[j].CreateAt
	})

	rtAllianceListInfo := make([]*gmmodel.AllianceInfo, 0)
	for i := page * pageSize; i < (page+1)*pageSize; i++ {
		if i >= int64(total) {
			break
		}
		rtAllianceListInfo = append(rtAllianceListInfo, allianceListInfo[i])
	}

	userList := &gmmodel.ListInfo{
		List:     rtAllianceListInfo,
		Total:    total,
		Count:    len(rtAllianceListInfo),
		Page:     int(page),
		PageSize: int(pageSize),
	}

	ret := gmmodel.MakeSuccessReturnMsg(userList)

	retJson, err := json.Marshal(ret)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	writer.Write(retJson)
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
