package friend

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/friendmodel"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
	"maze_game_server/services/userprofileservice"

	"github.com/gogo/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (f *FriendComponent) OnFriendListRecommend_10711_10712(s *session.Session, req *Friend.FriendListRecommendRQ) (err error) {
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Friend.FriendListRecommendRS{}

	logger.CtxInfo(ctx, "OnFriendListRecommend start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnFriendListRecommend end", zap.Any("res", res))
	}()

	recommendUsers, err := friendservice.GlobalFriendService.FriendRecommend(ctx, uint64(s.UID()), friendmodel.RecommendSize)
	if err != nil {
		logger.CtxError(ctx, "OnFriendListRecommend FriendRecommend fail",
			zap.Any("userID", s.UID()),
			zap.Error(err),
		)
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return
	}

	for _, recommendUser := range recommendUsers {
		nowUserProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, recommendUser)
		if err != nil {
			logger.CtxError(ctx, "OnFriendListRecommend GetUserProfile fail",
				zap.Any("userID", s.UID()),
				zap.Any("recommendUser", recommendUser),
				zap.Error(err),
			)
			res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
			return err
		}
		res.UserInfo = append(res.UserInfo, &Friend.User{
			UserId:     proto.Int64(int64(recommendUser)),
			UserName:   proto.String(nowUserProfile.NickName),
			UserGender: proto.Int32(nowUserProfile.Sex),
			AvaterUrl:  proto.String(nowUserProfile.Avatar),
		})
	}

	return
}
