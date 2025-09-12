package friend

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/userprofileservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (f *FriendComponent) OnSearchUserRQ_10709_10710(s *session.Session, req *Friend.SearchUserRQ) (err error) {
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &Friend.SearchUserRS{}

	logger.InfoWF("OnSearchUserRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnSearchUserRQ end", zap.Any("res", res))
	}()

	nowUserProfile, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, uint64(req.GetUserId()))
	if err != nil {
		logger.ErrorWF("OnSearchUserRQ GetUserProfile err ", zap.Error(err), zap.Uint64("searchUserID", uint64(req.GetUserId())))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
		return
	}

	if nowUserProfile == nil {
		res.UserInfo = &Friend.User{
			UserId:     proto.Uint64(uint64(req.GetUserId())),
			UserName:   proto.String("爱地牢"),
			UserGender: proto.Int32(1),
			AvaterUrl:  proto.String(""),
		}
		res.IsFriend = proto.Bool(false)
	}

	// todo check 是否为好友

	return
}
