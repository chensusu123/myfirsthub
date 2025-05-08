// @Author: ZhaoXiming 2025/3/21 16:31
// @Desc:

package mail

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mailboxsetdb"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mailunreaddb"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/protodef/MazeMailCli"
	"go.uber.org/zap"
)

func OnMailUnReadRQ(ctx fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req := rqMsg.(*MazeMailCli.MailUnReadRQ)
	res := rsMsg.(*MazeMailCli.MailUnReadRS)
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	userID := shardingID
	logger := ctx.FKLogI
	logger.InfoWF("OnMailUnReadRQ req", zap.Any("req", req))
	defer fkprometheus.DebugPMT("OnMailUnReadRQ")()

	res.TokenInfo = make([]*MazeMailCli.MailTokenInfo, 0, 2)
	for classType := int32(0); classType < 2; classType++ {
		tokenInfo := &MazeMailCli.MailTokenInfo{}
		tokenInfo.ClassType = proto.Int32(classType)
		res.TokenInfo = append(res.TokenInfo, tokenInfo)
		// 获取最新的过期时间
		expireMailIDList, latestExpireTime, loadErr := mailunreaddb.GetLatestUnreadMail(logger, userID, classType)
		if loadErr != nil {
			logger.WarnWF("OnMailUnReadRQ load data error", zap.Error(loadErr), zap.Int32("classType", classType))
		} else {
			tokenInfo.LatestExpireTime = proto.Int64(latestExpireTime)

			if len(expireMailIDList) > 0 {
				// 删除过期邮件id
				syncDelMail(logger, userID, classType, expireMailIDList, "OnMailUnReadRQ")
			}
		}
		var num int32
		num, err = mailunreaddb.UnreadMailNum(logger, userID, classType)
		if err != nil {
			logger.ErrorWF("OnMailUnReadRQ get unread mail num ", zap.Any("req", req), zap.Error(err))
		} else {
			tokenInfo.UnreadNum = proto.Int32(num)
		}
	}
	return nil
}

func syncDelMail(logger fklog.FKLogI, userID uint64, classType int32, mailIDS []uint64, callFuncName string) {
	if len(mailIDS) < 1 {
		return
	}
	logger.InfoWF("syncDelMail req", zap.String("callFuncName", callFuncName), zap.Uint64s("mailIDS", mailIDS))
	err := mailboxsetdb.DelMailBox(logger, userID, classType, mailIDS)
	if err != nil {
		logger.ErrorWF("syncDelMail end", zap.String("callFuncName", callFuncName), zap.Uint64s("mailIDS", mailIDS), zap.Error(err))
	}
	err = mailunreaddb.DelUnreadMail(logger, userID, classType, mailIDS)
	if err != nil {
		logger.ErrorWF("syncDelMail DelUnreadMail error", zap.String("callFuncName", callFuncName), zap.Uint64s("mailIDS", mailIDS), zap.Error(err))
	}
}
