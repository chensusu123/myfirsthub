package costume

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Costume"
	"maze_game_server/services/costumeservice"
	"time"
)

func (c *CostumeComponent) GetUserCostume_10613_10614(s *session.Session, req *Costume.GetUserCostumeRQ) (err error) {
	defer fkprometheus.InfoPMT("GetUserCostume")()

	start := time.Now()

	logger := log.Clone("Costume", uint64(s.UID()), 0)
	res := &Costume.GetUserCostumeRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	defer func() {
		err = s.Response(res)
		logger.InfoWF("GetUserCostume end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId := uint64(s.UID())
	if userId == 0 {
		logger.WarnWF("GetUserCostume args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
		return
	}

	costumeMap, err := costumeservice.GlobalCostumeService.GetUserCostume(logger, userId)
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
