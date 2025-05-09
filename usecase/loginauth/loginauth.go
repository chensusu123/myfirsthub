package loginauth

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/io/mysql/PaipaiUnionIDAuthMysql"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/mysql/UnionIDBindMysql"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func AuthJudge(logger fklog.FKLogI, userID uint64, authToken uint64) (ok bool, err error) {
	paipaiUnionID, err := UnionIDBindMysql.GetUserID2UnionID(logger, userID)
	if err != nil {
		logger.ErrorWF("AuthJudge call mysql.GetUserID2UnionID",
			zap.Uint64("userID", userID),
			zap.Any("err", err),
		)
		return false, err
	}

	if paipaiUnionID == 0 {
		logger.WarnWF("AuthJudge user not bind paipaiUnionID",
			zap.Uint64("userID", userID),
			zap.Any("err", err),
		)
		return false, nil
	}

	authTokenTmp, _, _, _, errTmp := PaipaiUnionIDAuthMysql.GetAuth(logger, paipaiUnionID)
	if errTmp != nil {
		logger.ErrorWF("AuthJudge call mysql.PaipaiUnionIDAuthMysql",
			zap.Uint64("userID", userID),
			zap.Uint64("paipaiUnionID", paipaiUnionID),
			zap.Any("err", err),
		)
		return false, errTmp
	}

	if authTokenTmp != int64(authToken) {
		logger.WarnWF("AuthJudge authToken not match",
			zap.Uint64("userID", userID),
			zap.Uint64("paipaiUnionID", paipaiUnionID),
			zap.Any("err", err),
		)
	}
	return true, nil
}
