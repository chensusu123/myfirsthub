package friend

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/friendmodel"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
	"maze_game_server/services/userprofileservice"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (f *FriendComponent) OnFriendApplySendList_10702_10703(s *session.Session, req *Friend.FriendApplySendListRQ) (err error) {
	userId := uint64(s.UID())
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Friend.FriendApplySendListRS{}
	res.Header = req.Header
	res.Page = req.Page

	logger.CtxInfo(ctx, "OnFriendApplySendList start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnFriendApplySendList end", zap.Any("res", res))
	}()

	page := req.GetPage()

	if page < 0 {
		return errors.COMMON_ERROR_TIPS.WrapMsg("page必须为正数")
	}

	pageSize := int32(10)
	friends, isFinish, err := friendservice.GlobalFriendService.SendFriendRequestList(ctx, userId, page, pageSize)
	if err != nil {
		logger.CtxError(ctx, "OnFriendApplySendList SendFriendRequestList err ", zap.Error(err), zap.Int32("page", page))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return
	}

	res.IsFinish = proto.Bool(isFinish)
	for _, friend := range friends {
		nowTime := time.Now()
		userProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, friend.ToUserId)
		if err != nil {
			logger.CtxError(ctx, "OnFriendApplySendList GetUserProfile err ", zap.Error(err), zap.Int32("page", page))
			res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
			return err
		}
		res.SendList = append(res.SendList, &Friend.SendInfo{
			UserInfo: &Friend.User{
				UserId:     proto.Int64(int64(friend.ToUserId)),
				UserName:   proto.String(userProfile.NickName),
				UserGender: proto.Int32(userProfile.Sex),
				AvaterUrl:  proto.String(userProfile.Avatar),
			},
			SendTime:   proto.Int64(nowTime.UnixMilli()),
			ExpireTime: proto.Int64(nowTime.UnixMilli() + int64(friendmodel.ExpireTime)),
		})
	}

	return nil
}
