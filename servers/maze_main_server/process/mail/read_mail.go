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

// 邮件阅读
func (g *Mail) OnMazeReadMailRQ_10628_10629(s *session.Session, req *MazeMail.MazeReadMailRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeReadMailRQ")()
	ctx := s.Context()
	logger := log.Clone("Frame", uint64(s.UID()), 0)
	res := &MazeMail.MazeReadMailRS{}

	logger.CtxInfo(ctx, "OnMazeReadMailRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnMazeReadMailRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	_, err = mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnMazeReadMailRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	if req.GetIsAll() {
		_, err := mailservice.GlobalMailService.ReadAllMail(ctx, userId, req.GetLabel())
		if err != nil {
			logger.CtxError(ctx, "OnMazeReadMailRQ ReadMail fail", zap.Error(err), zap.Uint64("userId", userId), zap.Int32("Label", req.GetLabel()))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
			return err
		}
	} else {
		_, err = mailservice.GlobalMailService.ReadMail(ctx, userId, req.GetMailId(), req.GetLabel())
		if err != nil {
			logger.CtxError(ctx, "OnMazeReadMailRQ ReadMail fail", zap.Error(err), zap.Uint64("mailId", req.GetMailId()))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
			return err
		}
	}

	list, err := mailservice.GlobalMailService.GetMailListByLabel(ctx, userId, req.GetLabel(), 0, 30)
	if err != nil {
		logger.CtxError(ctx, "OnMazeGetMailListRQ GetMailList fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return err
	}

	for _, info := range list {
		res.MailList = append(res.MailList, PbMailData(info))
	}

	return nil
}
