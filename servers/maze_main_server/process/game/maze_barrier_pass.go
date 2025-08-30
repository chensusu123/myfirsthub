package game

import (
	"context"
	"fmt"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/structsdef"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"maze_game_server/io/redis/mazebarriereventredis"
	"maze_game_server/io/redis/mazebarrieropstatusredis"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/io/redis/syncmazestorageinforedis"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/servers/maze_main_server/process/game/events"
	"maze_game_server/services/barriersavedataservice"
	"maze_game_server/services/barrierscorerewardservice"
	"maze_game_server/services/barrierservice"
	"maze_game_server/services/barrierstagecounterservice"
	"maze_game_server/services/tempbuffservice"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (g *Game) OnMazeBarrierPassRQ_10459_10460(s *session.Session, req *MazeGame.MazeBarrierPassRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeBarrierPassRQ")()

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	res := &MazeGame.MazeBarrierPassRS{}

	logger.CtxInfo(ctx, "OnMazeBarrierPassRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnMazeBarrierPassRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.BarrierId = req.BarrierId

	userId := uint64(s.UID())

	if req.GetBarrierId() <= 0 {
		logger.CtxError(ctx, "OnMazeBarrierPassRQ req barrier invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id未设置")
		return
	}
	//if req.GetFoeExp() < 0 {
	//	logger.ErrorWF("OnMazeBarrierPassRQ req barrier invalid", zap.Any("req", req))
	//	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("经验参数错误")
	//	return
	//}

	killMonsterNum, totalDamage, awards, rareAwards, errinfo := barrierservice.Global.BarrierPass(ctx, req.GetHeader(), userId, req.GetBarrierId(), req.GetFoeExp())
	if errinfo != nil {
		res.ErrInfo = errinfo
		return
	} else {
		res.KillMonsterNum = proto.Int32(killMonsterNum)
		res.TotalDamage = proto.Int64(totalDamage)
		res.BarrierAward = awards
		res.BarrierRareAward = rareAwards
	}

	err = mazebarriereventredis.LeaveBarrier(logger, userId, req.GetBarrierId(), true)
	if err != nil {
		logger.CtxError(ctx, "OnMazeBarrierPassRQ LeaveBarrier fail", zap.Error(err))
	}

	// 触发离开关卡事件
	events.OnLeaveBarrier(logger, userId, 0, time.Now().UnixMilli(), &MazeGame.BattleEventLeaveBarrier{BarrierId: proto.Int32(req.GetBarrierId()), Result: MazeGame.BarrierResult_PASS.Enum()})

	// 清除关卡的临时数据
	ClearBarriersTempData(ctx, userId, req.GetBarrierId())

	logger.CtxInfo(ctx, "OnMazeBarrierPassRQ award dump", zap.Any("exp", req.GetFoeExp()), zap.Any("awards", awards), zap.Any("rareAwards", rareAwards))

	passRecord := &mazebarrieruserkafka.MazeBarrierUserGameRecord{
		UserId:         userId,
		Barrier:        req.GetBarrierId(),
		GameRet:        mazebarrieruserkafka.GameRetSucc,
		Awards:         getAwards(logger, awards, rareAwards),
		KillMonsterNum: int64(killMonsterNum),
	}

	mazebarrieruserkafka.PushMazeBarrierUserRecord(ctx, passRecord)

	return nil
}

func getAwards(logger fklog.FKLogI, awards ...[]*MazeCommon.MazeItem) string {
	awardStr := make([]string, 0)
	for _, award := range awards {
		for _, v := range award {
			awardStr = append(awardStr, fmt.Sprintf("%d:%d", v.GetItemId(), v.GetCount()))
		}
	}
	return strings.Join(awardStr, "_")
}

// 清除关卡的临时数据
func ClearBarriersTempData(ctx context.Context, userId uint64, barrierId int32) {
	logger := fklog.ContextAppLogger(ctx)
	// 删除关卡存档
	syncmazestorageinforedis.DelSyncMazeStorageInfo(userId, barrierId)
	// 清理关卡操作状态
	mazebarrieropstatusredis.ClearOpStatus(ctx, userId, barrierId)
	//清除关卡已获得奖励存档
	barrierscorerewardservice.GlobalScoreRewardService.DelBarrierScoreRewardItem(context.TODO(), userId, barrierId)
	// 删除关卡存档 new
	barriersavedataservice.GlobalBarrierSaveDataService.DelBarrierSaveData(context.TODO(), userId, barrierId)
	// 删除临时buff
	tempbuffservice.GlobalTempBuffService.DelTempBuff(context.TODO(), userId, barrierId)
	// 删除通过的区域
	tempbuffservice.GlobalTempBuffService.DelPassArea(context.TODO(), userId, barrierId)
	//删除关卡计数
	barrierstagecounterservice.GlobalBarrierStageCounterService.DelBarrierStageCounterOnPass(ctx, userId, barrierId)

	mazebuffinforedis.DelMazeBuffBySrc(logger, userId, constdef.MazeBuffSrcSelectBuffForce)
	// 推送属性计算消息
	calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{
		UserId:  userId,
		ChgType: constdef.MazeBuffChgForceValue,
		Session: "buff",
		BuffSrc: constdef.MazeBuffSrcSelectBuffForce,
	}
	mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(ctx, calcAttrNotify)
}
