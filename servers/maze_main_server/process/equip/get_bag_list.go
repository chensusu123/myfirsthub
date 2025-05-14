package equip

import (
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGameEquip"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/bagmodule"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/packtopb"
)

func OnGetMazeBagEquipListRQ(ctx fknet.TCPContext, shardingID uint64, request proto.Message, response proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnGetMazeBagEquipListRQ")()
	req := request.(*MazeGameEquip.GetMazeBagEquipListRQ)
	res := response.(*MazeGameEquip.GetMazeBagEquipListRS)

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.QueryType = req.QueryType
	userCtx := fkserver.NewUserContext(ctx.Context, shardingID, ctx.FKLogI)

	defer func() {
		userCtx.InfoWF("OnGetMazeBagEquipListRQ end", zap.Any("res", res))
	}()

	userCtx.InfoWF("OnGetMazeBagEquipListRQ with", zap.Any("req", req))

	// if !BreedVersionFC.IsDollVersion(userCtx, shardingID) {
	// 	userCtx.ErrorWF("OnGetMazeBagEquipListRQ not doll version", zap.Uint64("userID", shardingID))
	// 	return
	// }
	bagEquipMgr := bagmodule.NewBagEquipMgr(ctx, shardingID)
	err = bagEquipMgr.LoadBagFromRedis()
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		ctx.ErrorWF("OnGetMazeBagEquipListRQ LoadBagFromRedis fail", zap.Error(err))
		return err
	}
	equipMap := bagEquipMgr.MainBagEquips.GetMainEquip()
	equipList := make([]*MazeGameEquip.MazeEquipInfo, 0)
	bagEquips := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
	for _, equipInfo := range equipMap {
		// newEquipInfo, err := pbutil.ConvertIdentifyEquipDb(userCtx, equipInfo)
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
			equipCli, err = packtopb.EquipSimplifyToCliPB(userCtx, equipInfo)
			if err != nil {
				res.ErrInfo = errors.MODULE_ERROR.ToInfo()
				ctx.ErrorWF("OnGetMazeBagEquipListRQ EquipSimplifyToCliPB fail", zap.Error(err))
				return err
			}
		}
		equipList = append(equipList, equipCli)
	}
	res.EquipList = equipList
	res.Token = proto.Int64(GetToken())
	return nil
}
