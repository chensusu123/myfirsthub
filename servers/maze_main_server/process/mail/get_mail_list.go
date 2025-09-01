package mail

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeMail"
	"maze_game_server/services/mailservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

// 邮件列表
func (g *Mail) OnMazeGetMailListRQ_10626_10627(s *session.Session, req *MazeMail.MazeGetMailListRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeGetMailListRQ")()
	ctx := s.Context()
	logger := log.Clone("Frame", uint64(s.UID()), 0)
	res := &MazeMail.MazeGetMailListRS{}

	logger.CtxInfo(ctx, "OnMazeGetMailListRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnMazeGetMailListRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	_, err = mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnMazeGetMailListRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	//首页签分页查询结果
	mailList, err := mailservice.GlobalMailService.GetMailListByLabel(ctx, userId, 0, 0, 30)
	if err != nil {
		logger.CtxError(ctx, "OnMazeGetMailListRQ GetMailList fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return err
	}
	for _, info := range mailList {
		res.MailList = append(res.MailList, PbMailData(info))
	}

	return nil
}
