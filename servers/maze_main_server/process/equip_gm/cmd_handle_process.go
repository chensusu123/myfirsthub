package equip_gm

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"maze_game_server/lib/net/websocket_service"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/errors"
	"maze_game_server/servers/maze_main_server/process/equip_gm/resetequipcmd"
)

func RegTcpHandler() {
	// 处理装备命令
	_ = websocket_service.RegProcSimple(10412, &MazeGameEquip.SendMazeEquipCmdRQ{},
		10413, &MazeGameEquip.SendMazeEquipCmdRS{}, OnSendMazeEquipCmdRQ)
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
