package alliance

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeFamily"
	"maze_game_server/services/allianceservice"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Alliance struct {
	component.Base
}

func NewAlliance() *Alliance {
	return &Alliance{}
}

func (a *Alliance) OnGetAllianceListRQ_10606_10607(s *session.Session, req *MazeFamily.GetAllianceListRQ) (err error) {
	logger := log.Clone("OnGetAllianceListRQ", uint64(s.UID()), 0)
	res := &MazeFamily.GetAllianceListRS{}

	res.ErrInfo = errors.NO_ERROR

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnGetAllianceListRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	allianceList, err := allianceservice.GlobalAllianceService.QueryAllianceList(logger)
	if err != nil {
		logger.ErrorWF("OnGetAllianceListRQ allianceService.GetAllianceList err", zap.Error(err))
		return
	}

	res.AllianceInfo = allianceList.DataToAllianceListPb(logger)
	return nil
}

func (a *Alliance) OnGetAllianceInfoRQ_10604_10605(s *session.Session, req *MazeFamily.GetUserAllianceRQ) (err error) {
	logger := log.Clone("OnGetAllianceInfoRQ", uint64(s.UID()), 0)
	res := &MazeFamily.GetUserAllianceRS{}

	res.ErrInfo = errors.NO_ERROR

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnGetAllianceInfoRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	allianceInfo, err := allianceservice.GlobalAllianceService.QueryAllianceInfo(logger, req.GetAllianceId())
	if err != nil {
		logger.ErrorWF("OnGetAllianceInfoRQ allianceService.GetAllianceInfo err", zap.Error(err))
		return
	}

	res.AllianceInfo = allianceInfo.DataToAllianceInfoPb()
	return nil
}

func (a *Alliance) OnQueryUserAllianceRQ_10608_10609(s *session.Session, req *MazeFamily.QueryUserAllianceRQ) (err error) {
	logger := log.Clone("OnQueryUserAllianceRQ", uint64(s.UID()), 0)
	res := &MazeFamily.QueryUserAllianceRS{}

	res.ErrInfo = errors.NO_ERROR

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnQueryUserAllianceRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	allianceID, err := allianceservice.GlobalAllianceService.QueryUserAlliance(logger, req.GetUserId())
	if err != nil {
		logger.ErrorWF("OnQueryUserAllianceRQ allianceService.GetAllianceInfo err", zap.Error(err))
		return err
	}

	res.AllianceId = proto.Int32(allianceID)
	return nil
}

func (a *Alliance) OnCreateAllianceRQ_10610_10611(s *session.Session, allianceName string) (err error) {
	logger := log.Clone("OnCreateAllianceRQ", uint64(s.UID()), 0)

	err = allianceservice.GlobalAllianceService.AddAlliance(logger, allianceName)
	if err != nil {
		logger.ErrorWF("OnCreateAllianceRQ allianceService.AddAlliance err", zap.Error(err))
		return err
	}

	return nil
}
