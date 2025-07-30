package friend

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/Friend"
	"maze_game_server/services/friendservice"
)

func (f *FriendComponent) OnFriendList_10560_10561(s *session.Session, req *Friend.FriendListRQ) (err error) {
	userId := uint64(s.UID())
	logger := log.Clone("OnFriendList", userId, 0)
	res := &Friend.FriendListRS{}

	logger.InfoWF("OnFriendList start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnFriendList end", zap.Any("res", res))
	}()

	page := req.GetPage()
	pageSize := req.GetPageSize()
	if page < 0 || pageSize < 0 {
		return errors.COMMON_ERROR_TIPS.WrapMsg("page 和 pageSize 必须为正数")
	}

	friends, codeErr := friendservice.GlobalFriendService.FriendList(logger, userId, page, pageSize)
	if codeErr != nil {
		logger.ErrorWF("OnFriendList FriendList err ", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
		res.ErrInfo = codeErr.ToInfo()
		return codeErr
	}
	res.Page = req.Page
	res.PageSize = req.PageSize
	res.FriendList = make([]*Friend.FriendInfo, 0, len(friends))
	// todo 头像昵称
	for _, friend := range friends {
		res.FriendList = append(res.FriendList, &Friend.FriendInfo{
			UserId:  proto.Uint64(friend.UserId),
			AddTime: proto.Int64(friend.CreateAt),
		})
	}

	return nil
}
