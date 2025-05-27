package equip

import (
	"maze_game_server/common/errors"
	"maze_game_server/common/function/packtopb"
	"maze_game_server/io/redis/mazebagequipredis"
	"maze_game_server/lib/log"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"

	"maze_game_server/lib/nano/session"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (e *Equip) OnQueryMazeEquipDetailRQ_10407_10408(s *session.Session, req *MazeGameEquip.QueryMazeEquipDetailRQ) (err error) {
	defer fkprometheus.DebugPMT("OnQueryMazeEquipDetailRQ")()

	logger := log.Clone("Equip", uint64(s.UID()), 0)
	res := &MazeGameEquip.QueryMazeEquipDetailRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.QueryType = req.QueryType
	res.PreviewForce = req.PreviewForce

	userId := uint64(s.UID())

	logger.InfoWF("OnQueryMazeEquipDetailRQ with", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnQueryMazeEquipDetailRQ end", zap.Any("res", res))
	}()
	// if !BreedVersionFC.IsDollVersion(logger, userId) {
	// 	logger.ErrorWF("OnQueryMazeEquipDetailRQ not doll version", zap.Uint64("userID", userId))
	// 	return
	// }
	if len(req.GetEquipGuids()) <= 0 {
		return
	}
	equipMap, err := mazebagequipredis.GetBatchEquipInfo(logger, userId, req.GetEquipGuids()...)
	bagEquips := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
	for _, equipInfo := range equipMap {
		// newEquipInfo, err := pbutil.ConvertIdentifyEquipDb(logger, equipInfo)
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
			equipCli, err = packtopb.EquipInfoToCliPB(logger, equipInfo)
			if err != nil {
				res.ErrInfo = errors.MODULE_ERROR.ToInfo()
				logger.ErrorWF("OnQueryMazeEquipDetailRQ EquipInfoToCliPB fail", zap.Error(err))
				return err
			}
		}
		equipList = append(equipList, equipCli)
	}
	res.EquipList = equipList
	res.Token = proto.Int64(GetToken())
	return
}
