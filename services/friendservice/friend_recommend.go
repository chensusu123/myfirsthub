package friendservice

import (
	"context"
	"maze_game_server/io/redis/onlineredis"
	"maze_game_server/model/friendmodel"

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
		logger.CtxError(ctx, "FriendRecommend NewSendFriendRequestModel err",
			zap.Uint64("userID", userID),
			zap.Error(err))
		return nil, err
	}

	blacklistModel, err := friendmodel.NewBlacklistModel(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "FriendRecommend NewBlacklistModel err", zap.Error(err), zap.Uint64("userID", userID))
		return nil, err
	}

	// todo 清理过期请求

	res := make([]uint64, 0)
	nowUser := make(map[uint64]struct{})

	onlineUser, err := onlineredis.GetAllOnline(ctx)
	if err != nil {
		logger.CtxError(ctx, "FriendRecommend GetAllOnline fail",
			zap.Uint64("userID", userID),
			zap.Error(err),
		)
		return nil, err
	}

	for _, uid := range onlineUser {
		nowUser[uid] = struct{}{}
	}

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

		for _, black := range blacklistModel.Blacklist {
			if nowID == black.UserId {
				isPass = true
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
		zap.Any("pageSize", pageSize),
	)
	return res, nil
}
