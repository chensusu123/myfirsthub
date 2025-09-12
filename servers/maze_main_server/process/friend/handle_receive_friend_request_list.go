package friend

// func (f *FriendComponent) OnReceiveFriendRequestList_10570_10571(s *session.Session, req *Friend.ReceiveFriendRequestListRQ) (err error) {
// 	userId := uint64(s.UID())
// 	logger := log.Clone("OnReceiveFriendRequestList", userId, 0)
// 	res := &Friend.ReceiveFriendRequestListRS{}

// 	logger.InfoWF("OnReceiveFriendRequestList start", zap.Any("req", req))
// 	defer func() {
// 		err = s.Response(res)
// 		logger.InfoWF("OnReceiveFriendRequestList end", zap.Any("res", res))
// 	}()

// 	page := req.GetPage()
// 	pageSize := req.GetPageSize()
// 	if page < 0 || pageSize < 0 {
// 		return errors.COMMON_ERROR_TIPS.WrapMsg("page 和 pageSize 必须为正数")
// 	}

// 	receives, codeErr := friendservice.GlobalFriendService.ReceiveFriendRequestList(logger, userId, page, pageSize)
// 	if codeErr != nil {
// 		logger.ErrorWF("OnReceiveFriendRequestList FriendList err ", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
// 		res.ErrInfo = codeErr.ToInfo()
// 		return
// 	}
// 	res.Page = req.Page
// 	res.PageSize = req.PageSize
// 	res.ReceiveList = make([]*Friend.ReceiveInfo, 0, len(receives))
// 	// todo 头像昵称
// 	for _, i := range receives {
// 		res.ReceiveList = append(res.ReceiveList, &Friend.ReceiveInfo{
// 			FromUserId:  proto.Uint64(i.FromUserId),
// 			ReceiveTime: proto.Int64(i.CreateAt),
// 			Status:      proto.Int32(i.Status),
// 		})
// 	}

// 	return nil
// }
