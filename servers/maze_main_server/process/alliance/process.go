package alliance

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeFamily"
	"maze_game_server/services/allianceservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type Alliance struct {
	component.Base
}

func NewAlliance() *Alliance {
	return &Alliance{}
}

func (a *Alliance) OnGetAllianceListRQ_10606_10607(s *session.Session, req *MazeFamily.GetAllianceListRQ) (err error) {
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeFamily.GetAllianceListRS{}

	res.ErrInfo = errors.NO_ERROR

	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnGetAllianceListRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	allianceList, err := allianceservice.GlobalAllianceService.QueryAllianceList(ctx)
	if err != nil {
		logger.CtxError(ctx, "OnGetAllianceListRQ allianceService.GetAllianceList err", zap.Error(err))
		return
	}

	res.AllianceInfo = allianceList.DataToAllianceListPb(ctx)
	return nil
}

func (a *Alliance) OnGetAllianceInfoRQ_10604_10605(s *session.Session, req *MazeFamily.GetUserAllianceRQ) (err error) {
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeFamily.GetUserAllianceRS{}

	res.ErrInfo = errors.NO_ERROR

	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnGetAllianceInfoRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	allianceInfo, err := allianceservice.GlobalAllianceService.QueryAllianceInfo(ctx, req.GetAllianceId())
	if err != nil {
		logger.CtxError(ctx, "OnGetAllianceInfoRQ allianceService.GetAllianceInfo err", zap.Error(err))
		return
	}

	res.AllianceInfo = allianceInfo.DataToAllianceInfoPb()
	return nil
}

func (a *Alliance) OnQueryUserAllianceRQ_10608_10609(s *session.Session, req *MazeFamily.QueryUserAllianceRQ) (err error) {
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeFamily.QueryUserAllianceRS{}

	res.ErrInfo = errors.NO_ERROR

	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnQueryUserAllianceRQ end", zap.Any("req", req), zap.Any("res", res))
	}()

	allianceID, err := allianceservice.GlobalAllianceService.QueryUserAlliance(ctx, req.GetUserId())
	if err != nil {
		logger.CtxError(ctx, "OnQueryUserAllianceRQ allianceService.GetAllianceInfo err", zap.Error(err))
		return err
	}
	_ = allianceID
	// res.AllianceId = proto.Int32(allianceID)
	return nil
}

func (a *Alliance) OnCreateAllianceRQ_10610_10611(s *session.Session, allianceName string) (err error) {
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	_, err = allianceservice.GlobalAllianceService.CreateAlliance(ctx, allianceName)
	if err != nil {
		logger.CtxError(ctx, "OnCreateAllianceRQ allianceService.AddAlliance err", zap.Error(err))
		return err
	}

	return nil
}

// 订阅联盟群聊
func (a *Alliance) OnSubscribeAllianceChatRQ_10691_10692(s *session.Session, req *MazeFamily.SubscribeAllianceChatRQ) (err error) {
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeFamily.SubscribeAllianceChatRS{}

	res.ErrInfo = errors.NO_ERROR
	uid := uint64(s.UID())
	group_ids := req.GetGroupIds()
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnSubscribeAllianceChatRQ end", zap.Any("req", req), zap.Any("res", res))
	}()
	// 检查参数
	if uid <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("联盟uid参数错误")
		return
	}
	allianceID, err := allianceservice.GlobalAllianceService.QueryUserAlliance(ctx, uid)
	if err != nil {
		logger.CtxError(ctx, "OnSubscribeAllianceChatRQ allianceService.GetAllianceInfo err", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("查询用户所属联盟失败")
		return
	}
	if allianceID <= 0 {
		logger.CtxError(ctx, "OnSubscribeAllianceChatRQ allianceService.GetAllianceInfo err", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("用户不属于任何联盟")
		return
	}

	// 订阅家族群聊
	err = allianceservice.GlobalAllianceService.SubscribeAllianceChat(ctx, uid, allianceID, group_ids)
	if err != nil {
		logger.CtxError(ctx, "OnSubscribeAllianceChatRQ allianceService.SubscribeAllianceChat err", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.Wrap("订阅联盟群聊失败")
		return
	}
	return
}
