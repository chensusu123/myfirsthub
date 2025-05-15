package common_value

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommonValueSvr"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGame"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/errors"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazemoneykafka"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazecommonvalue"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazemoney"
	"go.uber.org/zap"
)

func MazeCommonValueSubRQ(logger fklog.FKLogI, shardingID int64, req *MazeCommonValueSvr.MazeCommonValueSubRQ, res *MazeCommonValueSvr.MazeCommonValueSubRS) (err error) {
	res.ErrInfo = errors.NO_ERROR
	res.UserId = req.UserId

	logger.WarnWF("MazeCommonValueSubRQ with", zap.Any("rq", req))

	addStartTime := time.Now()
	defer func() {
		costTime := time.Since(addStartTime).Seconds()
		logger.WarnWF("MazeCommonValueSubRQ end ", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
		if costTime >= 0.5 {
			logger.ErrorWF("MazeCommonValueSubRQ timeout", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
		}
	}()

	query := req.GetSubItems()

	oldCoin, oldDiamond, err := mazemoney.GetUserMoney(logger, req.GetUserId())
	if err != nil {
		logger.ErrorWF("MazeCommonValueSubRQ GetUserMoney fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	oldMap := map[int32]int64{constdef.MazeCommonItemCoin: oldCoin, constdef.MazeCommonItemDiamond: oldDiamond}
	moneyMap := make(map[int32]int64)

	for _, item := range query {
		if item.GetItemId() == constdef.MazeCommonItemCoin {
			if oldCoin < item.GetCount() {
				res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("金币不足")
				return
			}
			moneyMap[item.GetItemId()] = oldCoin - item.GetCount()
		} else if item.GetItemId() == constdef.MazeCommonItemDiamond {
			if oldDiamond < item.GetCount() {
				res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("钻石不足")
				return
			}
			moneyMap[item.GetItemId()] = oldDiamond - item.GetCount()
		}
		res.Items = append(res.Items, &MazeCommon.MazeItem{ItemId: proto.Int32(item.GetItemId()), Count: proto.Int64(moneyMap[item.GetItemId()])})
	}

	err = mazemoney.BatchSetUserMoney(logger, req.GetUserId(), moneyMap)
	if err != nil {
		logger.ErrorWF("MazeCommonValueSubRQ BatchSetUserMoney fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	commonList := make([]*mazecommonvalue.CommonValueStruct, 0)
	moneyCommon := mazecommonvalue.MakeCommonValueList(logger, map[int32]int64{
		int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_MONEY):   moneyMap[constdef.MazeCommonItemCoin],
		int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_DIAMOND): moneyMap[constdef.MazeCommonItemDiamond]},
		map[int32]int32{}, map[int32]string{})
	commonList = append(commonList, moneyCommon...)
	mazecommonvalue.SendCommonValueIdPack(logger, req.GetUserId(), commonList)

	for _, item := range query {
		record := &mazemoneykafka.MazeMoneyRecord{
			UserId:        req.GetUserId(),
			OldMoneyId:    item.GetItemId(),
			OldMoneyCount: oldMap[item.GetItemId()],
			NewMoneyId:    item.GetItemId(),
			NewMoneyCount: moneyMap[item.GetItemId()],
			TradeNo:       int64(req.GetTradeNumber()),
			ChgReason:     req.GetOpType(),
		}
		mazemoneykafka.PushMazeMoneyRecord(logger, record)
	}
	return
}
