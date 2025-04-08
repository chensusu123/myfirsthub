package process

import (
	"time"

	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebarriermoneyredis"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/common/commonmustarriveredis"
	"gitlab.ifreetalk.com/plate/protodef/Common"
	"gitlab.ifreetalk.com/plate/protodef/DollMazeBarrier"
	"go.uber.org/zap"
)

func subMoney(logger fklog.FKLogI, uid uint64, moneyId int32, moneyCount int64, session string) error {
	curCount, err := mazebarriermoneyredis.GetMoneyCount(logger, uid, moneyId)
	if err != nil {
		logger.ErrorWF("GetMoneyCount cfgId fail", zap.Error(err))
		return err
	}
	if curCount < moneyCount {
		logger.ErrorWF("GetMoneyCount cfgId fail", zap.Int32("moneyId", moneyId), zap.Int64("curCount", curCount), zap.Int64("moneyCount", moneyCount))
		return errors.New("货币数量不足，购买失败")
	}
	err = mazebarriermoneyredis.SubMoney(logger, uid, moneyId, moneyCount)
	if err != nil {
		return err
	}

	SendMoneyChgPack(logger, uid, moneyId, curCount-moneyCount, constdef.DollMazeMoneyChgTypeBuy, session)
	// PushMazeMoneyChgRecord(logger, uid, moneyId, curCount, curCount-moneyCount)
	return nil
}

func SendMoneyChgPack(logger fklog.FKLogI, userId uint64, currMoneyId int32, currCount int64, reason int32, session string) (err error) {
	moneyPack := &DollMazeBarrier.BarrierMoneyID{
		NewMoney: &Common.Item{ItemId: proto.Int32(currMoneyId), Count: proto.Int64(currCount)},
		Reason:   proto.Int32(reason),
		Session:  proto.String(session),
		Token:    proto.Int64(time.Now().UnixNano() / 1e6),
	}

	logger.InfoWF("SendMoneyChgPack send client with", zap.Any("moneyPack", moneyPack))
	return commonmustarriveredis.SendArrivePacketFix(userId, 16139, moneyPack)
}
