package mazemoney

import (
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/servers/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/servers/maze_game_server/io/redis/mazebarriermoneyredis"
	"go.uber.org/zap"
)

func GetUserMoney(logger fklog.FKLogI, uid uint64) (coin, diamond int64, err error) {
	moneyMap, err := mazebarriermoneyredis.GetAllMoney(logger, uid)
	if err != nil {
		logger.ErrorWF("GetUserMoney GetMoneyCount fail", zap.Error(err))
		return
	}

	for k, v := range moneyMap {
		if k == constdef.MazeCommonItemCoin {
			coin = v
		} else if k == constdef.MazeCommonItemDiamond {
			diamond = v
		}
	}
	return
}

func SetUserDiamond(logger fklog.FKLogI, uid uint64, moneyCount int64) (err error) {
	err = mazebarriermoneyredis.SetMoney(logger, uid, constdef.MazeCommonItemDiamond, moneyCount)
	if err != nil {
		logger.ErrorWF("SetUserDiamond SetMoney fail", zap.Error(err))
		return
	}
	return
}

// AddUserDiamond 加钻石
// addDiamondCount 需要增加的钻石数量(正数)
func AddUserDiamond(logger fklog.FKLogI, uid uint64, addDiamondCount int64, reason int32) (moneyCount int64, err error) {

	moneyCount, err = mazebarriermoneyredis.AddMoney(logger, uid, constdef.MazeCommonItemDiamond, addDiamondCount)
	if err != nil {
		logger.ErrorWF("AddUserDiamond AddMoney fail", zap.Error(err))
		return
	}

	// record := &mazemoneykafka.MazeMoneyRecord{
	// 	UserId:        uid,
	// 	OldMoneyId:    MONEY_ID,
	// 	OldMoneyCount: moneyCount - addMoneyCount,
	// 	NewMoneyId:    MONEY_ID,
	// 	NewMoneyCount: moneyCount,
	// 	ChgReason:     reason,
	// }
	// mazemoneykafka.PushMazeMoneyRecord(logger, record)
	return
}

// SubUserDiamond 扣钻石
// subDiamondCount 需要扣除的钻石数量(正数)
func SubUserDiamond(logger fklog.FKLogI, uid uint64, subDiamondCount int64, reason int32) (moneyCount int64, err error) {

	moneyCount, err = mazebarriermoneyredis.AddMoney(logger, uid, constdef.MazeCommonItemDiamond, 0-subDiamondCount)
	if err != nil {
		logger.ErrorWF("SubUserDiamond AddMoney fail", zap.Error(err))
		return
	}

	// record := &mazemoneykafka.MazeMoneyRecord{
	// 	UserId:        uid,
	// 	OldMoneyId:    MONEY_ID,
	// 	OldMoneyCount: moneyCount + subMoneyCount,
	// 	NewMoneyId:    MONEY_ID,
	// 	NewMoneyCount: moneyCount,
	// 	ChgReason:     reason,
	// }
	// mazemoneykafka.PushMazeMoneyRecord(logger, record)
	return
}

func SetUserMoney(logger fklog.FKLogI, uid uint64, moneyCount int64) (err error) {
	err = mazebarriermoneyredis.SetMoney(logger, uid, constdef.MazeCommonItemCoin, moneyCount)
	if err != nil {
		logger.ErrorWF("GetUserMoney SetMoney fail", zap.Error(err))
		return
	}
	return
}

// 加扣物品并发问题 目前只有怪加 强化扣
// AddUserMoney 加钱
// addMoneyCount 需要增加的银子数量(正数)
func AddUserMoney(logger fklog.FKLogI, uid uint64, addMoneyCount int64, reason int32) (moneyCount int64, err error) {

	moneyCount, err = mazebarriermoneyredis.AddMoney(logger, uid, constdef.MazeCommonItemCoin, addMoneyCount)
	if err != nil {
		logger.ErrorWF("AddUserMoney AddMoney fail", zap.Error(err))
		return
	}

	// record := &mazemoneykafka.MazeMoneyRecord{
	// 	UserId:        uid,
	// 	OldMoneyId:    MONEY_ID,
	// 	OldMoneyCount: moneyCount - addMoneyCount,
	// 	NewMoneyId:    MONEY_ID,
	// 	NewMoneyCount: moneyCount,
	// 	ChgReason:     reason,
	// }
	// mazemoneykafka.PushMazeMoneyRecord(logger, record)
	return
}

// SubUserMoney 扣钱
// subMoneyCount 需要扣除的银子数量(正数)
func SubUserMoney(logger fklog.FKLogI, uid uint64, subMoneyCount int64, reason int32) (moneyCount int64, err error) {

	moneyCount, err = mazebarriermoneyredis.AddMoney(logger, uid, constdef.MazeCommonItemCoin, 0-subMoneyCount)
	if err != nil {
		logger.ErrorWF("SubUserMoney AddMoney fail", zap.Error(err))
		return
	}

	// record := &mazemoneykafka.MazeMoneyRecord{
	// 	UserId:        uid,
	// 	OldMoneyId:    MONEY_ID,
	// 	OldMoneyCount: moneyCount + subMoneyCount,
	// 	NewMoneyId:    MONEY_ID,
	// 	NewMoneyCount: moneyCount,
	// 	ChgReason:     reason,
	// }
	// mazemoneykafka.PushMazeMoneyRecord(logger, record)
	return
}

func BatchSetUserMoney(logger fklog.FKLogI, uid uint64, moneyMap map[int32]int64) (err error) {
	return mazebarriermoneyredis.HMSetMoney(logger, uid, moneyMap)
}
