package friend

// // 添加好友请求包
// func (f *FriendComponent) OnFriendApply_10694_10695(s *session.Session, req *Friend.FriendApplyRQ) (err error) {
// 	userId := uint64(s.UID())
// 	ctx := s.Context()
// 	logger := fklog.ContextAppLogger(ctx)
// 	res := &Friend.FriendApplyRS{}

// 	res.UserId = req.UserId
// 	res.From = req.From

// 	logger.CtxInfo(ctx, "OnFriendApply start", zap.Any("req", req))
// 	defer func() {
// 		err = s.Response(res)
// 		logger.InfoWF("OnFriendApply end", zap.Any("res", res))
// 	}()

// 	toID := req.GetUserId()
// 	if userId == toID {
// 		logger.CtxError(ctx, "OnFriendApply userId args error", zap.Any("req", req))
// 		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("不能加自己为好友")
// 		return err
// 	}

// 	err = friendservice.GlobalFriendService.AddFriendRequest(ctx, userId, toID)
// 	if err != nil {
// 		logger.ErrorWF("OnFriendApply FriendRequest failed ", zap.Error(err), zap.Uint64("userId", userId), zap.Uint64("toID", toID))
// 		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
// 		return
// 	}

// 	userProfiel, err := userprofileservice.GlobalUserProfileService.GetUserProfile(ctx, userId)
// 	if err != nil {
// 		logger.CtxError(ctx, "OnFriendApply GetUserProfile Fail",
// 			zap.Uint64("userID", userId),
// 			zap.Error(err),
// 		)
// 		res.ErrInfo = errors.COMMON_ERROR_TIPS.ToInfo()
// 		return
// 	}

// 	nowTime := time.Now()
// 	pushMsg := &Friend.FriendApplyID{
// 		UserInfo: &Friend.User{
// 			UserId:     proto.Uint64(userId),
// 			UserName:   proto.String(userProfiel.NickName),
// 			UserGender: proto.Int32(userProfiel.Sex),
// 			AvaterUrl:  proto.String(userProfiel.Avatar),
// 		},
// 		CreateTime: proto.Int64(nowTime.UnixMilli()),
// 		// 暂时设置7天过期
// 		ExpireTime: proto.Int64(nowTime.UnixMilli() + int64(friendmodel.ExpireTime)),
// 	}
// 	// 通知对方请求加好友
// 	online.ClusterPush(ctx, toID, 10696, pushMsg)

// 	return nil
// }
