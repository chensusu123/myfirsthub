package friendservice

import (
	"context"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/friendmodel"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) FriendRecommend(ctx context.Context, userID uint64, pageSize int32) ([]uint64, error) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "FriendRecommend start",
		zap.Uint64("userID", userID),
	)
	defer func() {
		logger.CtxInfo(ctx, "FriendRecommend end",
			zap.Uint64("userID", userID),
		)
	}()

	friendModel, err := friendmodel.NewFriendModel(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "FriendRecommend  NewFriendModel err", zap.Error(err), zap.Int32("pageSize", pageSize))
		return nil, err
	}

	sends, err := friendmodel.NewSendFriendRequestModel(ctx, userID)
	if err != nil {
		logger.ErrorWF("FriendRecommend NewSendFriendRequestModel err",
			zap.Uint64("userID", userID),
			zap.Error(err))
		return nil, err
	}

	res := make([]uint64, 0)
	nowUser := make(map[uint64]struct{})
	// 暂时做成推荐在线玩家
	online.Scan(func(id int64, s *session.Session) {
		s.RLock()
		defer s.RUnlock()
		if userID := s.UID(); userID <= 0 {
			return
		} else {
			nowUser[uint64(s.UID())] = struct{}{}
		}
	})

	for nowID := range nowUser {
		if len(res) == int(pageSize) {
			break
		}

		isPass := false

		if nowID == userID {
			isPass = true
		}

		for _, friendInfo := range friendModel.FriendList {
			if nowID == friendInfo.UserId {
				isPass = true
				break
			}
		}

		for _, send := range sends.SendList {
			if nowID == send.ToUserId {
				isPass = true
				break
			}
		}

		if isPass {
			continue
		}

		res = append(res, nowID)
	}

	logger.CtxInfo(ctx, "FriendRecommend GetUser",
		zap.Any("nowUser", nowUser),
		zap.Any("res", res),
	)
	return res, nil
}
