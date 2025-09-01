package equip

import (
	"maze_game_server/common/errors"
	"maze_game_server/common/function/packtopb"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/bagmodule"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (e *Equip) OnGetMazeBagEquipListRQ_10405_10406(s *session.Session, req *MazeGameEquip.GetMazeBagEquipListRQ) (err error) {
	defer fkprometheus.DebugPMT("OnGetMazeBagEquipListRQ")()

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	//logger := log.Clone("Equip", uint64(s.UID()), 0)
	res := &MazeGameEquip.GetMazeBagEquipListRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.QueryType = req.QueryType

	userId := uint64(s.UID())

	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnGetMazeBagEquipListRQ end", zap.Any("res", res))
	}()

	logger.CtxInfo(ctx, "OnGetMazeBagEquipListRQ with", zap.Any("req", req))

	// if !BreedVersionFC.IsDollVersion(logger, userId) {
	// 	logger.CtxError(ctx,"OnGetMazeBagEquipListRQ not doll version", zap.Uint64("userID", userId))
	// 	return
	// }
	bagEquipMgr := bagmodule.NewBagEquipMgr(ctx, userId)
	err = bagEquipMgr.LoadBagFromRedis(ctx)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		logger.CtxError(ctx, "OnGetMazeBagEquipListRQ LoadBagFromRedis fail", zap.Error(err))
		return err
	}
	equipMap := bagEquipMgr.MainBagEquips.GetMainEquip()
	equipList := make([]*MazeGameEquip.MazeEquipInfo, 0)
	bagEquips := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
	for _, equipInfo := range equipMap {
		// newEquipInfo, err := pbutil.ConvertIdentifyEquipDb(logger, equipInfo)
		// if err != nil {
		//	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		//	ctx.ErrorWF("OnGetMazeBagEquipListRQ ConvertIdentifyEquipDb fail", zap.Error(err))
		//	return err
		// }
		bagEquips = append(bagEquips, equipInfo)
	}
	for _, equipInfo := range bagEquips {
		equipCli := &MazeGameEquip.MazeEquipInfo{
			EquipGuid: proto.Int64(equipInfo.GetEquipGuid()),
		}
		if req.GetQueryType() == 0 {
			equipCli, err = packtopb.EquipSimplifyToCliPB(ctx, equipInfo)
			if err != nil {
				res.ErrInfo = errors.MODULE_ERROR.ToInfo()
				logger.CtxError(ctx, "OnGetMazeBagEquipListRQ EquipSimplifyToCliPB fail", zap.Error(err))
				return err
			}
		}
		equipList = append(equipList, equipCli)
	}
	res.EquipList = equipList
	res.Token = proto.Int64(GetToken())
	return nil
}
