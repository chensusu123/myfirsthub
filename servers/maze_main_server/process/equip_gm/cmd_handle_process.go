package equip_gm

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/servers/maze_main_server/process/equip_gm/resetequipcmd"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

type EquipGM struct {
	component.Base
}

func NewEquipGM() *EquipGM {
	return &EquipGM{}
}

func RegTcpHandler() {
	// // 处理装备命令
	// _ = websocket_service.RegProcSimple(10412, &MazeGameEquip.SendMazeEquipCmdRQ{},
	// 	10413, &MazeGameEquip.SendMazeEquipCmdRS{}, OnSendMazeEquipCmdRQ)
}

func (eg *EquipGM) OnSendMazeEquipCmdRQ_10412_10413(s *session.Session, req *MazeGameEquip.SendMazeEquipCmdRQ) (err error) {
	defer fkprometheus.DebugPMT("OnSendMazeEquipCmdRQ")()

	logger := log.Clone("EquipGM", uint64(s.UID()), 0)
	res := &MazeGameEquip.SendMazeEquipCmdRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.CmdCode = req.CmdCode
	res.CmdParam = req.CmdParam

	userId := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnSendMazeEquipCmdRQ end", zap.Any("res", res))
	}()

	logger.InfoWF("OnSendMazeEquipCmdRQ with", zap.Any("req", req))

	codeS := req.GetCmdCode()
	code := fkutil.ToInt32(codeS)
	switch code {
	case 1001:
		err = resetequipcmd.RunCmd1001(logger, userId, req.GetHeader().GetSession(), req.GetCmdParam())
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
