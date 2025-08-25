package costume

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Costume"
	"maze_game_server/services/costumeservice"
)

func (c *CostumeComponent) GetUserCostume_10613_10614(s *session.Session, req *Costume.GetUserCostumeRQ) (err error) {
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Costume.GetUserCostumeRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "GetUserCostume end", zap.Any("req", req), zap.Any("res", res))
	}()

	userId := uint64(s.UID())

	costumeMap, err := costumeservice.GlobalCostumeService.GetUserCostume(ctx, userId)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return nil
	}
	res.CostumeList = make([]*Costume.CostumeInfo, 0, len(costumeMap))
	for k, v := range costumeMap {
		res.CostumeList = append(res.CostumeList, &Costume.CostumeInfo{
			Pos:     proto.Int32(k),
			ModelId: proto.Int32(v),
		})
	}

	return
}
