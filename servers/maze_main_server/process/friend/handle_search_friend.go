package friend

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/userinfomodel"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
	"maze_game_server/services/userprofileservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (f *FriendComponent) OnSearchUserRQ_10709_10710(s *session.Session, req *Friend.SearchUserRQ) (err error) {
	userID := s.UID()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Friend.SearchUserRS{}
	res.Header = req.Header

	logger.CtxInfo(ctx, "OnSearchUserRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnSearchUserRQ end", zap.Any("res", res))
	}()

	if req.GetUserId() == userID {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("搜索id为自身id")
		return
	}

	// 检测用户是否存在
	toInfo, err := userinfomodel.NewUserInfoModel(ctx, uint64(req.GetUserId()))
	if err != nil {
		logger.CtxError(ctx, "OnSearchUserRQ NewUserInfoModel Fail",
			zap.Any("userID", userID),
			zap.Any("toID", req.GetUserId()),
		)
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return
	}

	if toInfo.Level < 1 {
		logger.CtxWarn(ctx, "OnSearchUserRQ Not Found User",
			zap.Any("userID", userID),
			zap.Any("toID", req.GetUserId()),
		)
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("用户不存在")
		return
	}

	nowUserProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, uint64(req.GetUserId()))
	if err != nil {
		logger.CtxError(ctx, "OnSearchUserRQ GetUserProfile err ", zap.Error(err), zap.Uint64("searchUserID", uint64(req.GetUserId())))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return
	}

	res.UserInfo = &Friend.User{
		UserId:     proto.Int64(req.GetUserId()),
		UserName:   proto.String(nowUserProfile.NickName),
		UserGender: proto.Int32(nowUserProfile.Sex),
		AvaterUrl:  proto.String(nowUserProfile.Avatar),
	}

	// todo check 是否为好友
	isFriend, err := friendservice.GlobalFriendService.CheckFriend(ctx, uint64(userID), uint64(req.GetUserId()))
	if err != nil {
		logger.CtxError(ctx, "OnSearchUserRQ CheckFriend Fail",
			zap.Error(err),
			zap.Any("userID", userID),
			zap.Any("toID", req.GetUserId()),
		)
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return
	}

	res.IsFriend = proto.Bool(isFriend)
	return
}
