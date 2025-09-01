package mazecommonvalue

import (
	"context"
	"maze_game_server/config/GMazeLevelV8Cfg"
	"maze_game_server/pb/common/MazeGame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 上报数据后 只更新额外加成
func HandleUserLevelChg(ctx context.Context, userId uint64, level int64, session string) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	commonList := make([]*CommonValueStruct, 0)

	addMoney, err := MakeMoneyExtra(ctx, userId, int64(level), 0)
	if err != nil {
		logger.CtxError(ctx, "handleUserLevelChg MakeMoneyExtra fail", zap.Error(err))
		return
	}
	// addExp, err := MakeExpExtra(logger, userId, int64(level), force)
	// if err != nil {
	// 	logger.CtxError(ctx,"handleUserLevelChg MakeExpExtra fail", zap.Error(err))
	// 	return
	// }

	extraMoney := &CommonValueStruct{
		DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_INCOME),
		DataValueInt: addMoney,
		// ChgReason:    int32(1),
		Session: session,
	}

	// extraExp := &CommonValueStruct{
	// 	DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP_INCOME),
	// 	DataValueInt: addExp,
	// 	// ChgReason:    int32(1),
	// 	Session: session,
	// }

	commonList = append(commonList, extraMoney)

	return SendCommonValueIdPack(ctx, userId, commonList)
}

// 等级经验变化后 都推送 用于通关、死亡、扫荡
func HandleUserLevelExpChg(ctx context.Context, userId uint64, level int64, exp int64, session string) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	commonList := make([]*CommonValueStruct, 0)

	levelCfg := GMazeLevelV8Cfg.GetWithCtx(ctx,int32(level))
	if levelCfg != nil {
		levelStruct := &CommonValueStruct{
			DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_LEVEL),
			DataValueInt: level,
			// ChgReason:    int32(1),
			Session: session,
		}
		expStruct := &CommonValueStruct{
			DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP),
			DataValueInt: exp,
			// ChgReason:    int32(1),
			Session: session,
		}
		expMaxStruct := &CommonValueStruct{
			DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP_MAX),
			DataValueInt: levelCfg.Next_level_need_exp,
			// ChgReason:    int32(1),
			Session: session,
		}
		commonList = append(commonList, levelStruct, expMaxStruct, expStruct)
	}

	addMoney, err := MakeMoneyExtra(ctx, userId, int64(level), 0)
	if err != nil {
		logger.CtxError(ctx, "handleUserLevelChg MakeMoneyExtra fail", zap.Error(err))
		return
	}
	// addExp, err := MakeExpExtra(logger, userId, int64(level), force)
	// if err != nil {
	// 	logger.CtxError(ctx,"handleUserLevelChg MakeExpExtra fail", zap.Error(err))
	// 	return
	// }

	extraMoney := &CommonValueStruct{
		DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_INCOME),
		DataValueInt: addMoney,
		// ChgReason:    int32(1),
		Session: session,
	}

	// extraExp := &CommonValueStruct{
	// 	DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP_INCOME),
	// 	DataValueInt: addExp,
	// 	// ChgReason:    int32(1),
	// 	Session: session,
	// }

	commonList = append(commonList, extraMoney)

	return SendCommonValueIdPack(ctx, userId, commonList)
}

func MakeMoneyExtra(ctx context.Context, userId uint64, level int64, force int64) (extra int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	moneyAddEquip, err := GetMoneyExtraAdditionEquip(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "MakeMoneyExtra GetMoneyExtraAdditionEquip fail", zap.Error(err))
		return
	}

	// moneyAddForce, _, _, err := GetExtraAdditionForce(logger, userId, level, force)
	// if err != nil {
	// 	logger.CtxError(ctx,"MakeMoneyExtra GetExtraAdditionForce fail", zap.Error(err))
	// 	return
	// }
	moneyAddForce := int64(0)

	extra = moneyAddEquip + moneyAddForce
	return
}

func MakeExpExtra(ctx context.Context, userId uint64, level int64, force int64) (extra int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	expAddEquip, err := GetExpExtraAdditionEquip(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "MakeExpExtra GetExpExtraAdditionEquip fail", zap.Error(err))
		return
	}

	// _, _, expAddForce, err := GetExtraAdditionForce(logger, userId, level, force)
	// if err != nil {
	// 	logger.CtxError(ctx,"MakeExpExtra GetExtraAdditionForce fail", zap.Error(err))
	// 	return
	// }

	expAddForce := int64(0)

	extra = expAddEquip + expAddForce
	return
}
