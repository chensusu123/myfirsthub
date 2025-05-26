package equip

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/packtopb"
	"maze_game_server/io/redis/mazebagequipredis"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"
)

func OnQueryMazeEquipDetailRQ(ctx fknet.TCPContext, shardingID uint64, request proto.Message, response proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnQueryMazeEquipDetailRQ")()
	req := request.(*MazeGameEquip.QueryMazeEquipDetailRQ)
	res := response.(*MazeGameEquip.QueryMazeEquipDetailRS)

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.QueryType = req.QueryType
	res.PreviewForce = req.PreviewForce
	userCtx := fkserver.NewUserContext(ctx.Context, shardingID, ctx.FKLogI)
	userCtx.InfoWF("OnQueryMazeEquipDetailRQ with", zap.Any("req", req))
	defer func() {
		userCtx.InfoWF("OnQueryMazeEquipDetailRQ end", zap.Any("res", res))
	}()
	// if !BreedVersionFC.IsDollVersion(userCtx, shardingID) {
	// 	userCtx.ErrorWF("OnQueryMazeEquipDetailRQ not doll version", zap.Uint64("userID", shardingID))
	// 	return
	// }
	if len(req.GetEquipGuids()) <= 0 {
		return
	}
	equipMap, err := mazebagequipredis.GetBatchEquipInfo(userCtx, userCtx.UserID, req.GetEquipGuids()...)
	bagEquips := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
	for _, equipInfo := range equipMap {
		// newEquipInfo, err := pbutil.ConvertIdentifyEquipDb(userCtx, equipInfo)
		// if err != nil {
		//	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		//	ctx.ErrorWF("OnQueryMazeEquipDetailRQ ConvertIdentifyEquipDb fail", zap.Error(err))
		//	return err
		// }
		bagEquips = append(bagEquips, equipInfo)
	}

	equipList := make([]*MazeGameEquip.MazeEquipInfo, 0)
	for _, equipInfo := range bagEquips {
		equipCli := &MazeGameEquip.MazeEquipInfo{
			EquipGuid: proto.Int64(equipInfo.GetEquipGuid()),
		}
		if req.GetQueryType() == 0 {
			equipCli, err = packtopb.EquipInfoToCliPB(userCtx, equipInfo)
			if err != nil {
				res.ErrInfo = errors.MODULE_ERROR.ToInfo()
				ctx.ErrorWF("OnQueryMazeEquipDetailRQ EquipInfoToCliPB fail", zap.Error(err))
				return err
			}
		}
		equipList = append(equipList, equipCli)
	}
	res.EquipList = equipList
	res.Token = proto.Int64(GetToken())
	return
}
