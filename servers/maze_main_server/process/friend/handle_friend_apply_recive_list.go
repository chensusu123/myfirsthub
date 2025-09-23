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

func (f *FriendComponent) OnFriendApplyReceiveList_10700_10701(s *session.Session, req *Friend.FriendApplyReceiveListRQ) (err error) {
	userId := uint64(s.UID())
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Friend.FriendApplyReceiveListRS{}

	logger.CtxInfo(ctx, "OnFriendApplyReceiveList start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnFriendApplyReceiveList end", zap.Any("res", res))
	}()

	page := req.GetPage()
	if page < 0 {
		return errors.COMMON_ERROR_TIPS.WrapMsg("page必须为正数")
	}

	receives, isFinsh, codeErr := friendservice.GlobalFriendService.ReceiveFriendRequestList(ctx, userId, page, int32(friendmodel.ReceivePage))
	if codeErr != nil {
		logger.CtxError(ctx, "OnFriendApplyReceiveList FriendList err ", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", int32(friendmodel.ReceivePage)))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return
	}

	logger.CtxInfo(ctx, "OnFriendApplyReceiveList GetReceives Successful",
		zap.Any("receives", receives),
		zap.Any("isFinsh", isFinsh),
	)

	res.Page = req.Page
	res.IsFinish = proto.Bool(isFinsh)

	res.ReceiveList = make([]*Friend.ReceiveInfo, 0, len(receives))

	for _, i := range receives {
		nowUserProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, i.FromUserId)
		if err != nil {
			logger.CtxError(ctx, "OnFriendApplyReceiveList GetUserProfile Fail",
				zap.Uint64("FromUserId", i.FromUserId),
				zap.Error(err),
			)
			res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
			return err
		}
		nowTime := time.Now().UnixMilli()
		res.ReceiveList = append(res.ReceiveList, &Friend.ReceiveInfo{
			UserInfo: &Friend.User{
				UserId:     proto.Int64(int64(i.FromUserId)),
				UserName:   proto.String(nowUserProfile.NickName),
				UserGender: proto.Int32(nowUserProfile.Sex),
				AvaterUrl:  proto.String(nowUserProfile.Avatar),
			},
			ReceiveTime: proto.Int64(nowTime),
			ExpireTime:  proto.Int64(nowTime + int64(friendmodel.ExpireTime)),
		})
	}

	return nil
}
