package buff

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeTempBuff"
	"maze_game_server/services/tempbuffservice"
	"time"
)

func (b *Buff) GetTempBuffGroupListRQ_10640_10641(s *session.Session, req *MazeTempBuff.GetTempBuffGroupListRQ) (err error) {
	defer fkprometheus.InfoPMT("GetTempBuffGroupListRQ")()
	start := time.Now()

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	logger.InfoWF("GetTempBuffGroupListRQ start", zap.Any("req", req))
	res := &MazeTempBuff.GetTempBuffGroupListRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.BarrierId = req.BarrierId

	defer func() {
		err = s.Response(res)
		logger.InfoWF("GetTempBuffGroupListRQ end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId, barrier := uint64(s.UID()), req.GetBarrierId()
	if userId == 0 || barrier == 0 {
		logger.WarnWF("GetTempBuffGroupListRQ args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
		return
	}

	groupList, err := tempbuffservice.GlobalTempBuffService.GetTempBuffGroupList(ctx, userId, barrier)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return nil
	}

	res.BuffGroupList = make([]*MazeTempBuff.TempBuffGroupInfo, 0, len(groupList))
	for k, v := range groupList {
		res.BuffGroupList = append(res.BuffGroupList, &MazeTempBuff.TempBuffGroupInfo{
			BuffId: proto.Int32(k),
			Count:  proto.Int64(int64(v)),
		})
	}

	return
}
