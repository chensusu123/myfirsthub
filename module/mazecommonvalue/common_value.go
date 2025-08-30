package mazecommonvalue

import (
	"context"
	"time"

	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var commonMap = map[int32]struct{}{
	int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_LEVEL):       {},
	int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP):         {},
	int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP_MAX):     {},
	int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_FORCE):       {},
	int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_MONEY):       {},
	int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_INCOME):      {},
	int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP_INCOME):  {},
	int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EQUIP_POINT): {},
	int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_DIAMOND):     {},
	int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_PASS_VALUE):  {},
}

type CommonValueStruct struct {
	DataType     int32
	DataValueInt int64
	ChgReason    int32
	Session      string
}

func SendCommonValueIdPack(ctx context.Context, userId uint64, commonList []*CommonValueStruct) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	commonValuePack := &MazeGame.MazeCommonValueChgID{
		CommonValueList: make([]*MazeGame.MazeCommonValueChg, 0),
	}

	token := time.Now().UnixNano() / 1e6

	for _, v := range commonList {
		_, ok := commonMap[v.DataType]
		if !ok {
			continue
		}
		cv := &MazeGame.MazeCommonValueChg{
			NewValue: &MazeGame.MazeCommonValue{
				ValueInt64: proto.Int64(v.DataValueInt),
			},
			DataType:  proto.Int32(v.DataType),
			ChgReason: proto.Int32(v.ChgReason),
			Session:   proto.String(v.Session),
			Token:     proto.Int64(token),
		}

		commonValuePack.CommonValueList = append(commonValuePack.CommonValueList, cv)
	}

	logger.CtxInfo(ctx, "sendCommonValueIdPack send client with", zap.Any("commonList", commonList), zap.Any("commonValuePack", commonValuePack))
	return online.ClusterPush(ctx, uint64(userId), 10478, commonValuePack)
}

func MakeCommonValueList(ctx context.Context, commonValue map[int32]int64, commonReason map[int32]int32, commonSession map[int32]string) (commonList []*CommonValueStruct) {
	commonList = make([]*CommonValueStruct, 0)
	for k, v := range commonValue {
		commonList = append(commonList, &CommonValueStruct{
			DataType:     k,
			DataValueInt: v,
			ChgReason:    commonReason[k],
			Session:      commonSession[k],
		})
	}
	return
}

func MakeCommonValueExtra(ctx context.Context, userId uint64, level int64, force int64) (extra int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	moneyAddEquip, err := GetMoneyExtraAdditionEquip(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "MakeCommonValueExtra GetExtraAdditionEquip fail", zap.Error(err))
		return
	}

	// moneyAddForce, _, _, err := GetExtraAdditionForce(logger, userId, level, force)
	// if err != nil {
	// 	logger.CtxError(ctx,"MakeCommonValueExtra GetExtraAdditionForce fail", zap.Error(err))
	// 	return
	// }

	moneyAddForce := int64(0)
	extra = moneyAddEquip + moneyAddForce
	return
}

func MakeCommonValueExtraExp(ctx context.Context, userId uint64, level int64, force int64) (extra int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	expAddEquip, err := GetExpExtraAdditionEquip(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "MakeCommonValueExtraExp GetExpExtraAdditionEquip fail", zap.Error(err))
		return
	}

	// expAddForce, _, _, err := GetExtraAdditionForce(logger, userId, level, force)
	// if err != nil {
	// 	logger.CtxError(ctx,"MakeCommonValueExtra GetExtraAdditionForce fail", zap.Error(err))
	// 	return
	// }

	expAddForce := int64(0)
	extra = expAddEquip + expAddForce
	return
}

func MakeAllCommonValue(logger fklog.FKLogI, userId uint64, level, exp, expMax, force, money, extra, extraExp, diamond, passValue int64, session string) (commonList []*CommonValueStruct) {
	commonList = make([]*CommonValueStruct, 0)

	lvStruct := &CommonValueStruct{
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
		DataValueInt: expMax,
		// ChgReason:    int32(1),
		Session: session,
	}

	forceStruct := &CommonValueStruct{
		DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_FORCE),
		DataValueInt: force,
		// ChgReason:    int32(1),
		Session: session,
	}

	moneyStruct := &CommonValueStruct{
		DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_MONEY),
		DataValueInt: money,
		// ChgReason:    int32(1),
		Session: session,
	}

	extraStruct := &CommonValueStruct{
		DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_INCOME),
		DataValueInt: extra,
		// ChgReason:    int32(1),
		Session: session,
	}

	// extraExpStruct := &CommonValueStruct{
	// 	DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP_INCOME),
	// 	DataValueInt: extraExp,
	// 	// ChgReason:    int32(1),
	// 	Session: session,
	// }

	diamondStruct := &CommonValueStruct{
		DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_DIAMOND),
		DataValueInt: diamond,
		// ChgReason:    int32(1),
		Session: session,
	}

	passValueStruct := &CommonValueStruct{
		DataType:     int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_PASS_VALUE),
		DataValueInt: passValue,
		// ChgReason:    int32(1),
		Session: session,
	}

	// commonList = append(commonList, lvStruct, expStruct, expMaxStruct, forceStruct, moneyStruct, extraStruct)
	commonList = append(commonList, lvStruct, expStruct, extraStruct, expMaxStruct, moneyStruct, diamondStruct, forceStruct, passValueStruct)
	return
}
