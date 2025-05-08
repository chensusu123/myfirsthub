// @Author: ZhaoXiming 2025/3/21 16:35
// @Desc:

package mail

import (
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/protodef/ItemSvr"
	"gitlab.ifreetalk.com/plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/plate/protodef/MessageType"
	"go.uber.org/zap"
	"math"
	"time"
)

//type FirstAddItem struct {
//	Id         int32 // 类型id
//	limitCount int64 // 上限
//	count      int64 // 当前数量
//	errMsg     error // 出错提示
//	alreadyErr bool  // 是否已经出错
//}
//
//func newFirstAddItem(ctx context.Context, userID uint64, logger fklog.FKLogI, itemId int32) *FirstAddItem {
//	ret := &FirstAddItem{}
//	ret.Id = itemId
//
//	switch itemId {
//	case classItem.ITEM_CLASS_STORE_HOUSE:
//		//仓库道具
//		ret.loadStoreHouseLimit(ctx, userID, logger)
//	}
//
//	return ret
//}
//
//func (f *FirstAddItem) canAdd(count int64) error {
//	if f.alreadyErr {
//		return f.errMsg
//	}
//
//	f.count += count
//
//	if f.limitCount == 0 {
//		f.alreadyErr = true
//		return f.errMsg
//	}
//
//	if f.limitCount < f.count {
//		f.alreadyErr = true
//		return nil
//	}
//	return nil
//}
//
////检查仓库
//func (f *FirstAddItem) canAddStoreHouse(count int64) error {
//	if f.alreadyErr {
//		return f.errMsg
//	}
//
//	f.count += count
//
//	if f.limitCount < f.count {
//		f.alreadyErr = true
//		return f.errMsg
//	}
//	return nil
//}
//
//func (f *FirstAddItem) loadStoreHouseLimit(ctx context.Context, userID uint64, logger fklog.FKLogI) {
//	req := &ItemSvr.QueryItemRQ{}
//	req.UserId = proto.Uint64(userID)
//	req.Items = []*MazeCommon.MazeItem{{ItemId: proto.Int32(classItem.ITEM_BAG_LATTICE)}}
//	res := &ItemSvr.QueryItemRS{}
//	var remain int64
//	err := ItemsNewRpc.QueryItemSvrRQ(logger, req, res)
//	if err != nil {
//		remain = 0
//	} else {
//		if len(res.Items) > 0 {
//			remain = res.Items[0].GetCount()
//		}
//	}
//	f.limitCount = remain
//	f.errMsg = errors.STOREHOUSE_FULL_ERROR
//}

var totalMailCache = simCache.NewCache()

func OnTotalMailRecvRQ(ctx fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req := rqMsg.(*MazeMailCli.TotalMailRecvRQ)
	res := rsMsg.(*MazeMailCli.TotalMailRecvRS)
	res.ErrInfo = errors.NO_ERROR
	userID := shardingID
	logger := ctx.FKLogI
	classType := req.GetClassType()
	res.ClassType = req.ClassType

	defer fkprometheus.DebugPMT("OnTotalMailRecvRQ")()

	// 限制重复请求
	if totalMailCache.SetCache(userID, struct{}{}) {
		res.ErrInfo = errors.NewCodeError(62091, "").ToInfo()
		return
	}
	defer totalMailCache.DelCache(userID)

	defer func() {
		res.Header = req.Header
		var num int32
		num, err = mailunreaddb.UnreadMailNum(logger, userID, classType)
		if err != nil {
			logger.ErrorWF("OnTotalMailRecvRQ get unread mail num", zap.Error(err), zap.Int32("classType", classType))
		} else {
			res.UnreadNum = proto.Int32(num)
		}

		// 获取最新的过期时间
		var latestExpireTime int64
		latestExpireTime, err = mailunreaddb.GetUnreaMailExpireTime(logger, userID, classType)
		if err != nil {
			logger.ErrorWF("OnTotalMailRecvRQ load data error", zap.Error(err))
		} else if latestExpireTime > 0 {
			res.LatestExpireTime = proto.Int64(latestExpireTime)
		}
		logger.WarnWF("OnTotalMailRecvRQ end", zap.Any("res", res), zap.Any("req", req))
	}()

	needDelMailIDS := make([]uint64, 0, 20)
	offset := uint64(math.MaxUint64)
	mailList := make([]*MazeMail.MailInfo, 0, 30)
	readMailId := make([]uint64, 0, 30)
	//var globalLimitItem int32
	//firstNoLimit := make(map[int32]*FirstAddItem, 100)
	//awardItemList := make(map[int][]*MazeCommon.MazeItem, 100)

	var specialAward bool

	for loopIndex := 0; loopIndex < MaxMailLoopCount; loopIndex++ {
		indexes, _ := mailunreaddb.GetMailIndexesFromUnreadList(logger, userID, classType, offset, mailunreaddb.GetMailIndexesLimitCount)
		if len(indexes) <= 0 {
			//结束了
			ctx.InfoWF("OnTotalMailRecvRQ indexes empty", zap.Any("offset", offset), zap.Any("loopIndex", loopIndex))
			break
		}
		searchIndexes := make([]uint64, 0)
		for _, index := range indexes {
			searchIndexes = append(searchIndexes, index.GetId())
			ctx.InfoWF("OnTotalMailRecvRQ all index", zap.Any("index", index))
		}

		mails, _ := MailBoxBodyRedis.GetMailList(logger, searchIndexes)
		now := time.Now().Unix()
		for i, mail := range mails {
			if mail.GetBase() == nil {
				ctx.InfoWF("OnTotalMailRecvRQ invalid need del 1", zap.Any("mail.Base", mail.GetBase()), zap.Uint64("delMailID", searchIndexes[i]))
				needDelMailIDS = append(needDelMailIDS, searchIndexes[i])
				continue
			}
			if mail.Base.GetMailId() == 0 {
				ctx.InfoWF("OnTotalMailRecvRQ invalid need del 2", zap.Any("mail.Base", mail.GetBase()), zap.Uint64("delMailID", searchIndexes[i]))
				needDelMailIDS = append(needDelMailIDS, searchIndexes[i])
				continue
			}

			if mail.GetBase().GetExpires() > 0 && now > mail.GetBase().GetExpires() {
				ctx.InfoWF("OnTotalMailRecvRQ Expires need del", zap.Any("mail.Base", mail.GetBase()), zap.Uint64("delMailID", searchIndexes[i]))
				needDelMailIDS = append(needDelMailIDS, searchIndexes[i])
				continue
			}

			//if (mail.GetBase().GetMailContentType() == int32(MazeMail.MailContentType_CONTENT_TRADE_BOX) && mail.TradeBox.GetIsDiamond() == 0) ||
			//	mail.GetBase().GetMailContentType() == int32(MazeMail.MailContentType_CONTENT_VOTE) {
			//	boxAward = true
			//	continue
			//}

			if mail.GetBase().GetNotTotalRecv() {
				// 不能一键领取
				specialAward = true
				continue
			}

			//mailSpecial := GMailSpecaialCfg.GetMailSpecaialConfig(mail.Base.GetExtType())
			//if mailSpecial != nil {
			//	// 有加成的需要单个领取，不能意见领取
			//	specialAward = true
			//	continue
			//}

			var hasRead, hasOperator int32
			_, hasRead, hasOperator, _, err = getMailFlag(logger, userID, mail)
			if err != nil {
				continue
			}

			if mailutil.IsMailTextOrReport(mail) {
				//文本
				if hasRead != 1 {
					mailList = append(mailList, mail)
				} else {
					readMailId = append(readMailId, mail.Base.GetMailId())
				}
				continue
			} else {
				//奖励
				if hasOperator == 1 {
					readMailId = append(readMailId, mail.Base.GetMailId())
					continue
				}
				//if isMailAwardFailFlag(logger, userID, mail.Base.GetMailId()) {
				//	continue
				//}
			}

			annexItem := make([]*MazeCommon.MazeItem, 0, 10)
			for _, annex := range mail.GetAnnex() {
				if annex.GetHasOpen() == 1 {
					continue
				}

				for _, tmpItem := range annex.Items {
					//// 一键领取不处理宝箱类型
					//if tmpItem.GetItemId() == TreasureBoxId {
					//	tmpItem.ItemId = proto.Int32(MaterialBoxId)
					//}
					annexItem = append(annexItem, tmpItem)
				}
			}
			if len(annexItem) < 1 {
				logger.WarnWF("OnTotalMailRecvRQ mail annex not have award", zap.Any("mail", mail))
				mailList = append(mailList, mail)
				continue
			}

			//opType := mail.Base.GetExtType()

			// TODO check add break最后再定, 等有上限逻辑了再用
			//limitErr := itemrpc.CheckAddItemsRpc(logger, userID, opType, annexItem)
			//if limitErr != nil {
			//	ctx.ErrorWF("OnTotalMailRecvRQ check add items fail", zap.Any("mail", mail), zap.Error(limitErr))
			//	break
			//}

			mailList = append(mailList, mail)
		}

		if len(indexes) < mailunreaddb.GetMailIndexesLimitCount {
			ctx.InfoWF("OnTotalMailRecvRQ is end empty", zap.Any("offset", offset), zap.Any("loopIndex", loopIndex))
			break
		}

		offset = indexes[len(indexes)-1].GetSequence()
		if offset == 0 {
			break
		} else {
			offset--
		}
	}

	ctx.InfoWF("OnTotalMailRecvRQ getMailList end",
		zap.Any("offset", offset),
		zap.Any("needDelMailIDS", needDelMailIDS),
		zap.Int("mail", len(mailList)),
	)

	if len(needDelMailIDS) > 0 {
		syncDelMail(logger, userID, classType, needDelMailIDS, "TotalMailDelRQ")
	}

	if len(mailList) < 1 && len(readMailId) < 1 {
		//if globalLimitItem != 0 {
		//	res.ErrInfo = wrapKnowErr(firstNoLimit[globalLimitItem].errMsg, errors.ERR_AWAD_MAIL_EMPTY)
		//	return nil
		//}
		if specialAward {
			res.ErrInfo = errors.NewCommonCodeError("剩余邮件需要手动领取")
		} else {
			res.ErrInfo = errors.ERR_AWAD_MAIL_EMPTY.ToInfo()
		}
		return nil
	}

	var convertErr error
	now := time.Now().Unix()
	itemIds := make(map[int32]int64, 3*len(mailList))
	awardList := make([]*MazeMail.MailInfo, 0, len(mailList))
	for _, mail := range mailList {
		if mailutil.IsMailHasAward(mail) {

			// 先改为已读 已操作
			mail.HasOperator = proto.Int32(1)
			mail.MailHasRead = proto.Int32(1)

			// 计算奖励 修改领取状态
			items := make([]*MazeCommon.MazeItem, 0, 4)
			for _, annex := range mail.GetAnnex() {
				if annex.GetHasOpen() == 0 {
					annex.HasOpen = proto.Int32(1)
					items = append(items, annex.Items...)
				}
			}

			if mail.Base.GetBroadcastFlag() == 0 {
				// 已读、已领取信封设置过期时间是1天
				mail.Base.Expires = proto.Int64(now + cgkargs.AlreadyReadMailExpireTime)
				err = MailBoxBodyRedis.AddMail(logger, mail)
			} else {
				err = broadcastmailstatdb.SetUseBroadcastMailState(logger, userID, mail.Base.GetMailId(), 1, 1, nil, mail.Base.GetExpires())
			}
			if err != nil {
				logger.ErrorWF("OnTotalMailRecvRQ TotalMailReceive add mail", zap.Any("mail", mail), zap.Error(err))
				return err
			}

			var gatherErr *MessageType.ErrorInfo

			// 只根据上面计算的有无奖励 决定是否领奖
			if len(items) > 0 {
				opType := mail.Base.GetExtType()

				var addRes *ItemSvr.AddItemRS
				tradeNum := mail.GetBase().GetMailId()
				if mail.Base.GetBroadcastId() != 0 {
					tradeNum = tradeno.GetTradeNoMaker().MakeTradeNo()
				}

				//TODO check 阶段的错误待处理
				gatherErr, err = itemrpc.GatherItems(logger, userID, req.Header, opType, tradeNum, items...)
				if err != nil {
					ctx.ErrorWF("OnTotalMailRecvRQ add items err",
						zap.Uint64("tradeNo", tradeNum),
						zap.Uint64("mailId", mail.GetBase().GetMailId()),
						zap.Any("items", items),
						zap.Error(err),
					)
					res.ErrInfo = errors.NewCommonCodeError("add item err")

				}
				if gatherErr != nil {
					ctx.WarnWF("OnTotalMailRecvRQ add items failed",
						zap.Uint64("tradeNo", tradeNum),
						zap.Uint64("mailId", mail.GetBase().GetMailId()),
						zap.Any("items", items),
						zap.Any("errorInfo", gatherErr),
					)
					res.ErrInfo = gatherErr
				}

				//TODO
				aggregationResults(gatherErr, itemIds, items, addRes)

				//流水
				mail_module.PushMailAwardLog(logger, userID, mail, items)
			}

			awardList = append(awardList, mail)
			readMailId = append(readMailId, mail.GetBase().GetMailId())
			if gatherErr != nil {
				err = errors.New(string(gatherErr.GetErrMsg()))
				break
			}
			if err != nil {
				break
			}

		} else {
			mail.HasOperator = proto.Int32(1)
			if mailutil.IsMailTextOrReport(mail) {
				//文本类型需要设置，操作类不修改
				mail.MailHasRead = proto.Int32(1)
			}
			if mail.Base.GetBroadcastFlag() == 0 {
				// 已读、已领取信封设置过期时间是1天
				mail.Base.Expires = proto.Int64(now + cgkargs.AlreadyReadMailExpireTime)
				err = MailBoxBodyRedis.AddMail(logger, mail)
			} else {
				err = broadcastmailstatdb.SetUseBroadcastMailState(logger, userID, mail.Base.GetMailId(), 1, 1, nil, mail.Base.GetExpires())
			}
			awardList = append(awardList, mail)
			readMailId = append(readMailId, mail.GetBase().GetMailId())
			if err != nil {
				logger.ErrorWF("OnTotalMailRecvRQ  add mail", zap.Any("mail", mail), zap.Error(err))
				break
			}
		}
	}

	items := make([]*MazeCommon.MazeItem, 0, len(itemIds))
	for id, count := range itemIds {
		items = append(items, &MazeCommon.MazeItem{ItemId: proto.Int32(id), Count: proto.Int64(count)})
	}
	res.AwardList = items

	if len(awardList) == 0 && convertErr != nil {
		res.ErrInfo = errors.NewCommonCodeError(convertErr.Error())
	}

	//if globalLimitItem != 0 {
	//	res.ErrInfo = wrapKnowErr(firstNoLimit[globalLimitItem].errMsg, errors.ERR_AWAD_MAIL_EMPTY)
	//} else
	if err != nil {
		res.ErrInfo = errors.NewCommonCodeError(err.Error())
	}

	//处理未读
	_ = mailunreaddb.DelUnreadMail(logger, userID, classType, readMailId)
	res.OpenMailList = readMailId

	return nil
}

//func wrapKnowErr(err error, defaultErr *errors.CodeError) (ErrInfo *MessageType.ErrorInfo) {
//	errInfo, ok := err.(*errors.CodeError)
//	if !ok {
//		return defaultErr.ToInfo()
//	}
//
//	switch errInfo {
//	case errors.SILVER_REACH_MAX,
//		errors.ERR_SHIP_EXP_FULL,
//		//航海术
//		errors.ERR_MERITORIOUS_SERVICE_FULL,
//		// 功勋
//		errors.STOREHOUSE_FULL_ERROR,
//		//背包
//		errors.ERR_SHIP_FULL_EQUIP_BAG,
//		//装备
//		errors.ERR_SHIP_LESS_EQUIP_BAG,
//		//装备
//		errors.BAG_FULL_ERROR,
//		// 背包满
//		errors.BAG_ITEM_LIMIT_MAX_ERROR,
//		errors.EACH_ADD_ITEM_ERROR,
//		// 单日领取上限
//		errors.ITEM_ADD_TIME_OUT,
//		errors.ITEM_CHECK_ERROR,
//		errors.ItemCheckFailure:
//		ErrInfo = errInfo.ToInfo()
//	default:
//		ErrInfo = defaultErr.ToInfo()
//	}
//
//	return ErrInfo
//}

//func GetClassifiedItems(logger fklog.FKLogI, items ...*MazeCommon.MazeItem) (classItemMap map[int][]*MazeCommon.MazeItem) {
//	classItemMap = make(map[int][]*MazeCommon.MazeItem)
//
//	var (
//		classItems []*MazeCommon.MazeItem
//		ok         bool
//	)
//
//	for _, item := range items {
//		class := GetItemClass(logger, item.GetItemId())
//		if class != classItem.ITEM_CLASS_INVALID {
//			if classItems, ok = classItemMap[class]; !ok {
//				classItems = make([]*MazeCommon.MazeItem, 0, 1)
//			}
//			classItems = append(classItems, item)
//			classItemMap[class] = classItems
//		}
//	}
//
//	return
//}

//func GetItemClass(logger fklog.FKLogI, itemId int32) int {
//	if itemId == 0 {
//		return classItem.ITEM_CLASS_INVALID
//	}
//
//	cnf := GItemsCfg.GetItemsConfig(itemId)
//	if cnf == nil {
//		logger.ErrorWF("invalid itemId", zap.Any("id", itemId))
//		return classItem.ITEM_CLASS_INVALID
//	}
//
//	// 老背包
//	if cnf.Type == classItem.ITEM_OLD_BAG_PET_ITEM ||
//		cnf.Type == classItem.ITEM_OLD_BAG_SKILL_ITEM ||
//		cnf.Type == classItem.ITEM_OLD_BAG_NORMAL_TYPE ||
//		cnf.Type == classItem.ITEM_OLD_BAG_TASK_ITEM {
//		logger.InfoWF("GetItemClass old bag", zap.Int32("id", itemId))
//		return classItem.ITEM_CLASS_OLD_BAG
//	}
//
//	if cnf.Is_bag == 1 {
//		return classItem.ITEM_CLASS_STORE_HOUSE
//	}
//	return classItem.ITEM_CLASS_INVALID
//}

// 聚合加道具的奖励信息,客户端票动画展示获得东西使用
func aggregationResults(err *MessageType.ErrorInfo, itemIds map[int32]int64, items []*MazeCommon.MazeItem, addRes *ItemSvr.AddItemRS) {
	if addRes == nil {
		return
	}

	// 1、加成功并且没有道具转化
	// 2、加道具失败，展示全部信息
	if err != nil || len(addRes.ConvertItems) <= 0 {
		for _, tmpItem := range items {
			itemIds[tmpItem.GetItemId()] += tmpItem.GetCount()
		}
		return
	}

	// 加道具成功并且发生了道具转化
	for _, cItem := range addRes.SuccItems {
		itemIds[cItem.GetItemId()] += cItem.GetCount()
	}
	if len(addRes.Items) > 0 {
		for _, item := range addRes.Items {
			for _, cItem := range item.ChgItems {
				itemIds[cItem.GetItemId()] += cItem.GetCount()
			}
		}
	}
	return
}
