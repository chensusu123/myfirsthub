package friend

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Friend"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (f *FriendComponent) OnFriendListRecommend_10711_10712(s *session.Session, req *Friend.FriendListRecommendRQ) (err error) {
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Friend.FriendListRecommendRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	logger.CtxInfo(ctx, "OnFriendListRecommend start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnFriendListRecommend end", zap.Any("res", res))
	}()

	// 临时测试使用
	// 固定推荐
	nowUsers := []uint64{10003384, 10003386, 10003388, 10003390, 10003392, 10003394}
	for _, nowUser := range nowUsers {
		res.UserInfo = append(res.UserInfo, &Friend.User{
			UserId:     proto.Int64(int64(nowUser)),
			UserName:   proto.String("爱地牢"),
			UserGender: proto.Int32(1),
			AvaterUrl:  proto.String(""),
		})
	}

	// recommendUsers, err := friendservice.GlobalFriendService.FriendRecommend(ctx, uint64(s.UID()), friendmodel.RecommendSize)
	// if err != nil {
	// 	logger.CtxError(ctx, "OnFriendListRecommend FriendRecommend fail",
	// 		zap.Any("userID", s.UID()),
	// 		zap.Error(err),
	// 	)
	// 	res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
	// 	return
	// }

	// logger.CtxInfo(ctx, "OnFriendListRecommend FriendRecommend Successful",
	// 	zap.Any("recommendUsers", recommendUsers),
	// )

	// for _, recommendUser := range recommendUsers {
	// 	nowUserProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, recommendUser)
	// 	if err != nil {
	// 		logger.CtxError(ctx, "OnFriendListRecommend GetUserProfile fail",
	// 			zap.Any("userID", s.UID()),
	// 			zap.Any("recommendUser", recommendUser),
	// 			zap.Error(err),
	// 		)
	// 		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
	// 		return err
	// 	}
	// 	res.UserInfo = append(res.UserInfo, &Friend.User{
	// 		UserId:     proto.Int64(int64(recommendUser)),
	// 		UserName:   proto.String(nowUserProfile.NickName),
	// 		UserGender: proto.Int32(nowUserProfile.Sex),
	// 		AvaterUrl:  proto.String(nowUserProfile.Avatar),
	// 	})
	// 	logger.CtxInfo(ctx, "OnFriendListRecommend res append data",
	// 		zap.Any("res", res),
	// 	)
	// }

	return nil
}
