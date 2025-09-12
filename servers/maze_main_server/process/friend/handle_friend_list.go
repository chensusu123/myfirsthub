package friend

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
	"maze_game_server/services/userprofileservice"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (f *FriendComponent) OnFriendList_10704_10705(s *session.Session, req *Friend.FriendListRQ) (err error) {
	userId := uint64(s.UID())
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Friend.FriendListRS{}
	res.Page = req.Page

	logger.CtxInfo(ctx, "OnFriendList start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnFriendList end", zap.Any("res", res))
	}()

	page := req.GetPage()
	if page < 0 {
		return errors.COMMON_ERROR_TIPS.WrapMsg("page 和 pageSize 必须为正数")
	}

	pageSize := int32(10)

	friends, finish, err := friendservice.GlobalFriendService.FriendList(ctx, userId, page, pageSize)
	if err != nil {
		logger.CtxError(ctx, "OnFriendList FriendList err ", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return err
	}

	res.IsFinish = proto.Bool(finish)
	res.FriendList = make([]*Friend.FriendInfo, 0, len(friends))

	for _, friend := range friends {
		nowUserProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, friend.UserId)
		if err != nil {
			logger.CtxError(ctx, "OnFriendList GetUserProfile Fail",
				zap.Uint64("friendID", friend.UserId),
				zap.Error(err),
			)
			res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
			return err
		}
		res.FriendList = append(res.FriendList, &Friend.FriendInfo{
			UserInfo: &Friend.User{
				UserId:     proto.Int64(int64(friend.UserId)),
				UserName:   proto.String(nowUserProfile.NickName),
				UserGender: proto.Int32(nowUserProfile.Sex),
				AvaterUrl:  proto.String(nowUserProfile.Avatar),
			},
			FriendType: proto.Int32(friend.FriendType),
			AddTime:    proto.Int64(friend.CreateAt),
		})
	}

	if len(res.FriendList) == 0 {
		for i := 10010; i <= 10020; i++ {
			res.FriendList = append(res.FriendList, &Friend.FriendInfo{
				UserInfo: &Friend.User{
					UserId:     proto.Int64(int64(i)),
					UserName:   proto.String("爱地牢"),
					UserGender: proto.Int32(1),
					AvaterUrl:  proto.String(""),
				},
				FriendType: proto.Int32(1),
				AddTime:    proto.Int64(time.Now().UnixMilli()),
			})
		}
	}

	return nil
}
