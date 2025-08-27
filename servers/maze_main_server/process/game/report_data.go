package game

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/io/kafka/mazemoneykafka"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/equipdropmodel"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/services/equipdropservice"
	"maze_game_server/services/moneyservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (g *Game) OnReportDataRQ_10453_10454(s *session.Session, req *MazeGame.ReportDataRQ) (err error) {
	defer fkprometheus.InfoPMT("ReportDataRQ")()

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

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

	userInfo, err = mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.ErrorWF("ReportDataRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	if reportInfo.GetReportMask()&1 == 1 {

		// 上报总经验
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

	var dropInfo *equipdropmodel.EquipSpecialDropModel
	//var shopInfo *mazeshopseqredis.MazeShopInfo
	if reportInfo.GetReportMask()&4 == 4 {
		// 上报装备积分
		dropInfo, err = equipdropmodel.NewEquipSpecialDropModel(ctx, userId)
		if err != nil {
			logger.ErrorWF("ReportDataRQ GetEquipSpecialDropModel fail", zap.Error(err), zap.Uint64("userId", userId))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
		dropInfo.EquipPoints = int32(reportInfo.GetEquipPoint())

		//shopInfo, err = calequipsequence.GetMazeShopInfo(logger, userId, int32(userInfo.Level), userInfo.Barrier)
		//if err != nil {
		//	logger.ErrorWF("ReportDataRQ GetMazeShopInfo fail", zap.Error(err))
		//	res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		//	return
		//}
		//shopInfo.EquipPoints = int32(reportInfo.GetEquipPoint())
	}

	// 修改上报数据的存储
	if reportInfo.GetReportMask()&1 == 1 {
		err = mazeuserinfo.SetUserInfoV2(ctx, userId, userInfo)
		if err != nil {
			logger.ErrorWF("ReportDataRQ SetUserInfoV2 fail", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
		mazecommonvalue.HandleUserLevelExpChg(logger, userId, userInfo.Level, userInfo.Exp, req.GetHeader().GetSession())
		if levelRecord.OldLevel != levelRecord.NewLevel {
			mazeuserlevelkafka.PushMazeLevelRecord(ctx, levelRecord)
		}

	}

	if reportInfo.GetReportMask()&2 == 2 {
		var oldCoin int64
		oldCoin, _, err = moneyservice.GlobalMoneyService.GetUserMoney(context.TODO(), userId)
		if err != nil {
			logger.ErrorWF("MazeCommonValueQueryRQ GetUserMoney fail", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		}

		// rpc不支持set 他们也需要加锁 目前先自己直接设置
		err = moneyservice.GlobalMoneyService.SetMoney(context.TODO(), userId, constdef.MazeCommonItemCoin, reportInfo.GetMoneyCount())
		if err != nil {
			logger.ErrorWF("ReportDataRQ SetMoney fail", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}

		// 上报金币
		if oldCoin != reportInfo.GetMoneyCount() {
			record := &mazemoneykafka.MazeMoneyRecord{
				UserId:        userId,
				OldMoneyId:    constdef.MazeCommonItemCoin,
				OldMoneyCount: oldCoin,
				NewMoneyId:    constdef.MazeCommonItemCoin,
				NewMoneyCount: reportInfo.GetMoneyCount(),
				TradeNo:       int64(0),
				ChgReason:     0,
			}
			mazemoneykafka.PushMazeMoneyRecord(ctx, record)
		}
	}

	if reportInfo.GetReportMask()&4 == 4 {
		newLevel := equipdropservice.GlobalEquipDropService.GetMazeBarrierLv(int32(userInfo.Level), userInfo.Barrier)
		err = dropInfo.Save(ctx, userId)
		if err != nil {
			logger.ErrorWF("ReportDataRQ EquipSpecialDropModel save fail", zap.Error(err), zap.Any("level", newLevel), zap.Any("dropInfo", dropInfo))
			return err
		}

		//err = mazeshopseqredis.SetMazeShopInfo(logger, userId, int32(newLevel), shopInfo)
		//if err != nil {
		//	logger.ErrorWF("ReportDataRQ SetMazeShopInfo fail", zap.Error(err), zap.Any("level", newLevel), zap.Any("shopInfo", shopInfo))
		//}
	}

	// if reportInfo.GetReportMask()&8 == 8 {
	// 	// todo 客户端数据 存储

	// }

	return
}
