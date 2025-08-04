package buff

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeTempBuff"
	"maze_game_server/services/tempbuffservice"
	"time"
)

func (b *Buff) GetMazeTempBuffListRQ_10433_10434(s *session.Session, req *MazeTempBuff.GetMazeTempBuffListRQ) (err error) {
	defer fkprometheus.InfoPMT("GetMazeTempBuffListRQ")()

	start := time.Now()

	logger := log.Clone("Buff", uint64(s.UID()), 0)
	res := &MazeTempBuff.GetMazeTempBuffListRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.StageId = req.StageId

	defer func() {
		err = s.Response(res)
		logger.InfoWF("GetMazeTempBuffListRQ end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId, stageId := uint64(s.UID()), req.GetStageId()
	if userId == 0 || stageId == 0 {
		logger.WarnWF("GetMazeTempBuffListRQ args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
		return
	}

	buffList, err := tempbuffservice.GlobalTempBuffService.GetMazeTempBuffList(logger, userId, stageId)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return nil
	}
	res.BuffList = buffInfo2MazeBuffInfo(buffList)

	return
}

func buffInfo2MazeBuffInfo(buffList []*tempbuffservice.BuffInfo) []*MazeTempBuff.MazeBuffInfo {
	mazeBuffList := make([]*MazeTempBuff.MazeBuffInfo, 0, len(buffList))
	for _, buff := range buffList {
		mazeBuffList = append(mazeBuffList, &MazeTempBuff.MazeBuffInfo{
			BuffId:      proto.Int32(buff.BuffId),
			Value:       proto.Int64(buff.Value),
			Name:        proto.String(buff.Name),
			Decs:        proto.String(buff.Decs),
			IsRecommend: proto.Int32(buff.IsRecommend),
		})
	}
	return mazeBuffList
}
