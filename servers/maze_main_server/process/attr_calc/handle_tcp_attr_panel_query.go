/*
 * @Author: majian
 * @Date: 2024-07-10 11:40:04
 * @Last Modified by: majian
 * @Last Modified time: 2024-12-13 21:49:08
 */
package attr_calc

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/common/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazePropertyPanel"
	"go.uber.org/zap"
)

func OnQueryPropertyPanelRQ(ctx fknet.TCPContext, shardingID uint64, request proto.Message, response proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnQueryPropertyPanelRQ")()
	req := request.(*MazePropertyPanel.QueryMazePropertyPanelRQ)
	res := response.(*MazePropertyPanel.QueryMazePropertyPanelRS)

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	userCtx := fkserver.NewUserContext(ctx.Context, shardingID, ctx.FKLogI)

	defer func() {
		userCtx.InfoWF("OnQueryPropertyPanelRQ end", zap.Any("res", res))
	}()

	userCtx.InfoWF("OnQueryPropertyPanelRQ with", zap.Any("req", req))

	panelCalc := NewDPAC(shardingID)
	err = panelCalc.Init(userCtx, DPACParam{Force: 0})
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		userCtx.ErrorWF("OnQueryPropertyPanelRQ Init fail", zap.Error(err))
		return err
	}
	err = panelCalc.Calc(userCtx)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		userCtx.ErrorWF("OnQueryPropertyPanelRQ Calc fail", zap.Error(err))
		return err
	}
	res.PropertyPanel = panelCalc.Panel
	return
}
