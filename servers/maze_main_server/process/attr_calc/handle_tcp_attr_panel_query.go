/*
 * @Author: majian
 * @Date: 2024-07-10 11:40:04
 * @Last Modified by: majian
 * @Last Modified time: 2024-12-13 21:49:08
 */
package attr_calc

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazePropertyPanel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (p *Property) OnQueryPropertyPanelRQ_10427_10428(s *session.Session, req *MazePropertyPanel.QueryMazePropertyPanelRQ) (err error) {
	defer fkprometheus.DebugPMT("OnQueryPropertyPanelRQ")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazePropertyPanel.QueryMazePropertyPanelRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	userId := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnQueryPropertyPanelRQ end", zap.Any("res", res))
	}()

	logger.CtxInfo(ctx, "OnQueryPropertyPanelRQ with", zap.Any("req", req))

	panelCalc := NewDPAC(userId)
	err = panelCalc.Init(ctx, DPACParam{Force: 0})
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		logger.ErrorWF("OnQueryPropertyPanelRQ Init fail", zap.Error(err))
		return err
	}
	err = panelCalc.Calc(ctx)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		logger.ErrorWF("OnQueryPropertyPanelRQ Calc fail", zap.Error(err))
		return err
	}
	res.PropertyPanel = panelCalc.Panel
	return
}
