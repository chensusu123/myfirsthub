package equip_gm

import (
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
	"gitlab.ifreetalk.com/plate/protodef/MazeGameEquip"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/equip_gm/resetequipcmd"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
)

func RegTcpHandler() {
	// 处理装备命令
	_ = websocket_service.RegProcSimple(16186, &MazeGameEquip.SendMazeEquipCmdRQ{},
		16187, &MazeGameEquip.SendMazeEquipCmdRS{}, OnSendMazeEquipCmdRQ)
}

func OnSendMazeEquipCmdRQ(ctx fknet.TCPContext, shardingID uint64, request proto.Message, response proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnSendMazeEquipCmdRQ")()
	req := request.(*MazeGameEquip.SendMazeEquipCmdRQ)
	res := response.(*MazeGameEquip.SendMazeEquipCmdRS)

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.CmdCode = req.CmdCode
	res.CmdParam = req.CmdParam
	userCtx := fkserver.NewUserContext(ctx.Context, shardingID, ctx.FKLogI)

	defer func() {
		userCtx.InfoWF("OnSendMazeEquipCmdRQ end", zap.Any("res", res))
	}()

	userCtx.InfoWF("OnSendMazeEquipCmdRQ with", zap.Any("req", req))

	codeS := req.GetCmdCode()
	code := fkutil.ToInt32(codeS)
	switch code {
	case 1001:
		err = resetequipcmd.RunCmd1001(userCtx, shardingID, req.GetHeader().GetSession(), req.GetCmdParam())
		if err != nil {
			err = errors.New("执行失败")
		}
	default:
		err = errors.New("未知命令")
	}
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
	}
	return nil
}
