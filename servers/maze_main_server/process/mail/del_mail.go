package mail

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeMail"
	"maze_game_server/services/mailservice"
)

// 邮件删除
func (g *Mail) OnMazeDelMailRQ_10632_10633(s *session.Session, req *MazeMail.MazeDelMailRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeDelMailRQ")()

	logger := log.Clone("Frame", uint64(s.UID()), 0)
	res := &MazeMail.MazeDelMailRS{}

	logger.InfoWF("OnMazeDelMailRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeDelMailRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	_, err = mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeDelMailRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	if req.GetIsAll() {
		err = mailservice.GlobalMailService.DelAllMail(logger, userId, req.GetLabel())
		if err != nil {
			logger.ErrorWF("OnMazeDelMailRQ ReadMail fail", zap.Error(err), zap.Int32("label", req.GetLabel()))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
			return
		}
	} else {
		err = mailservice.GlobalMailService.DelMail(logger, userId, req.GetMailId(), req.GetLabel())
		if err != nil {
			logger.ErrorWF("OnMazeDelMailRQ ReadMail fail", zap.Error(err), zap.Uint64("mailId", req.GetMailId()), zap.Int32("label", req.GetLabel()))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
			return
		}
	}

	list, err := mailservice.GlobalMailService.GetMailListByLabel(logger, userId, req.GetLabel(), 0, 30)
	if err != nil {
		logger.ErrorWF("OnMazeGetMailListRQ GetMailList fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return err
	}

	for _, info := range list {
		res.MailList = append(res.MailList, PbMailData(info))
	}

	return nil
}
