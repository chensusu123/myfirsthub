package friend

// func (f *FriendComponent) OnBlacklist_10572_10573(s *session.Session, req *Friend.BlacklistRQ) (err error) {
// 	userId := uint64(s.UID())
// 	logger := log.Clone("OnBlacklist", userId, 0)
// 	res := &Friend.BlacklistRS{}

// 	logger.InfoWF("OnBlacklist start", zap.Any("req", req))
// 	defer func() {
// 		err = s.Response(res)
// 		logger.InfoWF("OnBlacklist end", zap.Any("res", res))
// 	}()

// 	page := req.GetPage()
// 	pageSize := req.GetPageSize()
// 	if page < 0 || pageSize < 0 {
// 		return errors.COMMON_ERROR_TIPS.WrapMsg("page 和 pageSize 必须为正数")
// 	}

// 	friends, codeErr := friendservice.GlobalFriendService.Blacklist(logger, userId, page, pageSize)
// 	if codeErr != nil {
// 		logger.ErrorWF("OnBlacklist Blacklist err ", zap.Error(err), zap.Int32("page", page), zap.Int32("pageSize", pageSize))
// 		res.ErrInfo = codeErr.ToInfo()
// 		return codeErr
// 	}
// 	res.Page = req.Page
// 	res.PageSize = req.PageSize
// 	res.BlackUserList = make([]*Friend.BlackUserInfo, 0, len(friends))
// 	// todo 头像昵称
// 	for _, friend := range friends {
// 		res.BlackUserList = append(res.BlackUserList, &Friend.BlackUserInfo{
// 			UserId:  proto.Uint64(friend.UserId),
// 			AddTime: proto.Int64(friend.CreateAt),
// 		})
// 	}

// 	return nil
// }
