package friend

// func (f *FriendComponent) OnSendFriendRequestList_10568_10569(s *session.Session, req *Friend.SendFriendRequestListRQ) (err error) {
// 	userId := uint64(s.UID())
// 	logger := log.Clone("OnSendFriendRequestList", userId, 0)
// 	res := &Friend.SendFriendRequestListRS{}

// 	logger.InfoWF("OnSendFriendRequestList start", zap.Any("req", req))
// 	defer func() {
// 		err = s.Response(res)
// 		logger.InfoWF("OnSendFriendRequestList end", zap.Any("res", res))
// 	}()

// 	page := req.GetPage()
// 	pageSize := req.GetPageSize()
// 	if page < 0 || pageSize < 0 {
// 		return errors.COMMON_ERROR_TIPS.WrapMsg("page 和 pageSize 必须为正数")
// 	}

// 	friends, codeErr := friendservice.GlobalFriendService.SendFriendRequestList(logger, userId, page, pageSize)
// 	if codeErr != nil {
// 		logger.ErrorWF("OnSendFriendRequestList SendFriendRequestList err ", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
// 		res.ErrInfo = codeErr.ToInfo()
// 		return codeErr
// 	}
// 	res.Page = req.Page
// 	res.PageSize = req.PageSize
// 	res.SendList = make([]*Friend.SendInfo, 0, len(friends))
// 	// todo 头像昵称
// 	for _, friend := range friends {
// 		res.SendList = append(res.SendList, &Friend.SendInfo{
// 			ToUserId: proto.Uint64(friend.ToUserId),
// 			SendTime: proto.Int64(friend.CreateAt),
// 		})
// 	}

// 	return nil
// }
