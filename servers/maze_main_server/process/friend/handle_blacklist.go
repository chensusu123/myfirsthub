package friend

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
	"maze_game_server/services/userprofileservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (f *FriendComponent) OnBlacklist_10717_10718(s *session.Session, req *Friend.BlacklistRQ) (err error) {
	userId := uint64(s.UID())
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Friend.BlacklistRS{}
	res.Header = req.Header
	res.Page = req.Page

	logger.CtxInfo(ctx, "OnBlacklist start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnBlacklist end", zap.Any("res", res))
	}()

	page := req.GetPage()
	pageSize := int32(10)

	if page < 0 {
		return errors.COMMON_ERROR_TIPS.WrapMsg("page必须为正数")
	}

	friends, isFinsh, err := friendservice.GlobalFriendService.Blacklist(ctx, userId, page, pageSize)
	if err != nil {
		logger.CtxError(ctx, "OnBlacklist Blacklist err ", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return
	}

	res.IsFinish = proto.Bool(isFinsh)
	for _, friend := range friends {
		userProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, friend.UserId)
		if err != nil {
			logger.CtxError(ctx, "OnBlacklist GetUserProfile Fail",
				zap.Any("toID", friend.UserId),
				zap.Error(err),
			)
			return err
		}
		res.BlackUserList = append(res.BlackUserList, &Friend.BlackUserInfo{
			UserInfo: &Friend.User{
				UserId:     proto.Int64(int64(userProfile.UserID)),
				UserName:   proto.String(userProfile.NickName),
				UserGender: proto.Int32(userProfile.Sex),
				AvaterUrl:  proto.String(userProfile.Avatar),
			},
			AddTime: proto.Int64(friend.CreateAt),
		})
	}

	return nil
}
