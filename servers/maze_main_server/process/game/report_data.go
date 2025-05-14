package game

import (
	"github.com/lonng/nano/session"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazemoneykafka"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazeuserlevelkafka"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeshopseqredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/calequipsequence"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazecommonvalue"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazemoney"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGame"

	"go.uber.org/zap"
)

func (g *Game) OnReportDataRQ(s *session.Session, req *MazeGame.ReportDataRQ) (err error) {
	fkprometheus.InfoPMT("ReportDataRQ")()

	logger := fklog.AppLogger().Clone("game")
	res := &MazeGame.ReportDataRS{}

	logger.InfoWF("ReportDataRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("ReportDataRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.UserInfo = req.UserInfo

	userId := uint64(s.UID())
	reportInfo := req.GetUserInfo()

	if reportInfo.GetReportMask() <= 0 {
		return
	}

	var userInfo *mazeuserinfo.UserInfo
	var levelRecord *mazeuserlevelkafka.MazeUserLevelRecord

	userInfo, err = mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("ReportDataRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	if reportInfo.GetReportMask()&1 == 1 {

		//上报总经验
		levelRecord = &mazeuserlevelkafka.MazeUserLevelRecord{
			UserId:   userId,
			OldLevel: int32(userInfo.Level),
		}

		userInfo.SetTotalExp(reportInfo.GetExpTotal())
		err = userInfo.CalExp()
		if err != nil {
			logger.ErrorWF("ReportDataRQ CalTotalExp fail", zap.Error(err), zap.Any("totalExp", reportInfo.GetExpTotal()))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
		levelRecord.NewLevel = int32(userInfo.Level)

	}

	var shopInfo *mazeshopseqredis.MazeShopInfo
	if reportInfo.GetReportMask()&4 == 4 {
		// 上报装备积分
		shopInfo, err = calequipsequence.GetMazeShopInfo(logger, userId, int32(userInfo.Level), userInfo.Barrier)
		if err != nil {
			logger.ErrorWF("ReportDataRQ GetMazeShopInfo fail", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}

		shopInfo.EquipPoints = int32(reportInfo.GetEquipPoint())
	}

	// 修改上报数据的存储
	if reportInfo.GetReportMask()&1 == 1 {
		err = mazeuserinfo.SetUserInfoV2(logger, userId, userInfo)
		if err != nil {
			logger.ErrorWF("ReportDataRQ SetUserInfoV2 fail", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
		mazecommonvalue.HandleUserLevelExpChg(logger, userId, userInfo.Level, userInfo.Exp, req.GetHeader().GetSession())
		if levelRecord.OldLevel != levelRecord.NewLevel {
			mazeuserlevelkafka.PushMazeLevelRecord(logger, levelRecord)
		}

	}

	if reportInfo.GetReportMask()&2 == 2 {
		//上报金币

		// oldCount, err2 := mazemoney.GetUserMoney(logger, userId)
		// if err2 != nil {
		// 	logger.ErrorWF("ReportDataRQ GetUserMoney fail", zap.Error(err2))
		// 	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		// 	return
		// }

		// rpc不支持set 他们也需要加锁 目前先自己直接设置
		err = mazemoney.SetUserMoney(logger, userId, reportInfo.GetMoneyCount())
		if err != nil {
			logger.ErrorWF("ReportDataRQ SetMoney fail", zap.Error(err))
			// res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			// return
		}

		record := &mazemoneykafka.MazeMoneyRecord{
			UserId:        userId,
			OldMoneyId:    constdef.MazeCommonItemCoin,
			OldMoneyCount: 0,
			NewMoneyId:    constdef.MazeCommonItemCoin,
			NewMoneyCount: reportInfo.GetMoneyCount(),
			TradeNo:       int64(0),
			ChgReason:     0,
		}
		mazemoneykafka.PushMazeMoneyRecord(logger, record)
	}

	if reportInfo.GetReportMask()&4 == 4 {
		newLevel := calequipsequence.GetMazeBarrierLv(int32(userInfo.Level), userInfo.Barrier)
		err = mazeshopseqredis.SetMazeShopInfo(logger, userId, int32(newLevel), shopInfo)
		if err != nil {
			logger.ErrorWF("ReportDataRQ SetMazeShopInfo fail", zap.Error(err), zap.Any("level", newLevel), zap.Any("shopInfo", shopInfo))
		}
	}

	// if reportInfo.GetReportMask()&8 == 8 {
	// 	// todo 客户端数据 存储

	// }

	return
}
