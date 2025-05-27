package game

import (
	"fmt"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/addequip"
	"maze_game_server/common/function/gentradeno"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/function/packtopb/equiptoitem"
	"maze_game_server/config/GMazeActionCountV8Cfg"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/redis/mazechallengenumredis"
	"maze_game_server/io/redis/mazeuserbarrierredis"
	"maze_game_server/lib/log"
	"maze_game_server/module/mazebarrier"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/server/MazeEquipSvr"
	"strings"
	"time"

	"github.com/lonng/nano/session"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (g *Game) OnMazeBarrierPassRQ_10459_10460(s *session.Session, req *MazeGame.MazeBarrierPassRQ) (err error) {
	defer fkprometheus.InfoPMT("OnMazeBarrierPassRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)

	res := &MazeGame.MazeBarrierPassRS{}

	logger.InfoWF("OnMazeBarrierPassRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeBarrierPassRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	if req.GetBarrierId() <= 0 {
		logger.ErrorWF("OnMazeBarrierPassRQ req barrier invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id未设置")
		return
	}
	if req.GetFoeExp() < 0 {
		logger.ErrorWF("OnMazeBarrierPassRQ req barrier invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("经验参数错误")
		return
	}

	cfg := GMazeBarriesV8Cfg.Get(req.GetBarrierId())
	if cfg == nil {
		logger.ErrorWF("OnMazeBarrierPassRQ get barrier cfg fail", zap.Any("barrier", req.GetBarrierId()))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡配置数据获取失败")
		return
	}
	rareMap := make(map[int32]struct{})
	for _, v := range cfg.Rare_items_show {
		rareMap[v] = struct{}{}
	}
	rareMap = map[int32]struct{}{}

	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierPassRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	if userInfo.PassBarrier >= req.GetBarrierId() {
		logger.ErrorWF("OnMazeBarrierPassRQ req barrier lt pass barrier", zap.Any("req", req), zap.Int32("pass", userInfo.PassBarrier))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("该关卡已上报过通关")
		return
	}

	if userInfo.Barrier != req.GetBarrierId() {
		logger.ErrorWF("OnMazeBarrierPassRQ userinfo barrier not match", zap.Any("req", req), zap.Int32("save", userInfo.Barrier))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("记录用户正在打的关卡与上报通关id不匹配")
		return
	}

	userBarrier, err := mazeuserbarrierredis.GetUserBarrierInfo(logger, userId, req.GetBarrierId())
	if err != nil {
		logger.ErrorWF("OnMazeBarrierPassRQ GetUserBarrierInfo fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	userInfo.SetPassBarrier(req.GetBarrierId())
	// userInfo.SetBarrier(cfg.Next_id)
	oldLevel := userInfo.Level
	oldExp := userInfo.TotalExp
	err = userInfo.AddExp(int64(req.GetFoeExp()))
	if err != nil {
		logger.ErrorWF("OnMazeBarrierPassRQ AddExp fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	newLevel := userInfo.Level
	err = mazeuserinfo.SetUserInfoV2(logger, userId, userInfo)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierPassRQ SetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	mazecommonvalue.HandleUserLevelExpChg(logger, userId, userInfo.Level, userInfo.Exp, req.GetHeader().GetSession())

	defer func() {
		if oldLevel != newLevel {
			levelRecord := &mazeuserlevelkafka.MazeUserLevelRecord{
				UserId:      userId,
				OldLevel:    int32(oldLevel),
				OldTotalExp: oldExp,
				NewLevel:    int32(newLevel),
				NewTotalExp: int32(userInfo.TotalExp),
			}
			mazeuserlevelkafka.PushMazeLevelRecord(logger, levelRecord)
		}
	}()

	if req.GetFoeExp() > 0 {
		expItem := &MazeCommon.MazeItem{ItemId: proto.Int32(constdef.MazeCommonItemExp), Count: proto.Int64(int64(req.GetFoeExp()))}
		_, ok := rareMap[constdef.MazeCommonItemExp]
		if ok {
			res.BarrierRareAward = append(res.BarrierRareAward, expItem)
		} else {
			res.BarrierAward = append(res.BarrierAward, expItem)
		}
	}

	//成功通关需要 (1和2通过kafka清理)
	//1. 清除临时buff
	//2. 影响挂机生产
	//3. 清除客户端透传数据(通过切换关卡id 切换不同的key 目前没清)

	// 死亡之后是否需要清空复活次数
	userBarrier.RebornCount = proto.Int32(0)
	userBarrier.BarrierStatus = proto.Int32(3)
	userBarrier.EndTime = proto.Int64(time.Now().Unix())
	err = mazeuserbarrierredis.SetUserBarrierInfo(logger, userId, 0, userBarrier)
	if err != nil {
		logger.ErrorWF("OnMazeBarrierPassRQ SetUserBarrierInfo fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	now := time.Now()
	today := now.Year()*10000 + int(now.Month())*100 + now.Day()

	var awardBarrierNum int32
	awardBarrierNumCfg := GMazeActionCountV8Cfg.Get(102)
	if awardBarrierNumCfg != nil {
		awardBarrierNum = awardBarrierNumCfg.Day_count_v8
	}

	useNumToday, err2 := mazechallengenumredis.GetUserChallengeNum(logger, userId, today)
	if err2 != nil {
		logger.ErrorWF("OnMazeBarrierPassRQ GetUserChallengeNum fail", zap.Error(err2))
	} else if useNumToday > 0 {
		if awardBarrierNum > int32(useNumToday) {
			awardBarrierNum = int32(useNumToday)
		}
		err = mazechallengenumredis.AddUserChallengeNum(logger, userId, today, 0-awardBarrierNum)
		if err != nil {
			logger.ErrorWF("OnMazeBarrierPassRQ AddUserChallengeNum fail", zap.Error(err))
		}
	}

	awardMap, equipMap, err := mazebarrier.GetBarrierPassAwardWithFirst(logger, req.GetBarrierId())
	if err != nil {
		logger.ErrorWF("OnMazeBarrierPassRQ GetBarrierPassAwardWithFirst fail", zap.Error(err), zap.Any("barrier", req.GetBarrierId()))
	} else {
		//696	UN_CGK_COMMON_BILL_TYPE_696	迷宫通关
		tradeNo := gentradeno.GetTradeNum()
		if len(awardMap) > 0 {
			awardItems := itemutil.Map2Common(awardMap)
			errInfo := gentradeno.AddItemEx(logger, userId, 696, tradeNo, req.GetHeader(), awardItems...)
			if errInfo != nil {
				logger.ErrorWF("OnMazeBarrierPassRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("awardItems", awardItems))
			}

			for _, item := range awardItems {
				_, ok := rareMap[item.GetItemId()]
				if ok {
					res.BarrierRareAward = append(res.BarrierRareAward, item)
				} else {
					res.BarrierAward = append(res.BarrierAward, item)
				}
			}
		}

		if len(equipMap) > 0 {
			//MAZE_EQUIP_PASS_AWARD = 9;//迷宫通关奖励 张登元
			rs, err := addequip.AddEquipToBag(logger, userId, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_PASS_AWARD), tradeNo, equipMap)
			if err != nil {
				logger.ErrorWF("OnMazeBarrierPassRQ addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_BOX_AWARD)),
					zap.Any("tradeNo", tradeNo), zap.Any("addEquip", equipMap))
			}

			for _, equip := range rs.GetEquipList() {
				itemEquip, err := equiptoitem.PackEquipToItem(equip)
				if err != nil {
					logger.ErrorWF("OnMazeBarrierPassRQ PackEquipToItem fail", zap.Error(err), zap.Any("equip", equip))
					continue
				}
				_, ok := rareMap[itemEquip.GetItemId()]
				if ok {
					res.BarrierRareAward = append(res.BarrierRareAward, itemEquip)
				} else {
					res.BarrierAward = append(res.BarrierAward, itemEquip)
				}
			}
		}
	}

	logger.InfoWF("OnMazeBarrierPassRQ award dump", zap.Any("exp", req.GetFoeExp()), zap.Any("awardItem", awardMap), zap.Any("awardEquip", equipMap))

	passRecord := &mazebarrieruserkafka.MazeBarrierUserGameRecord{
		UserId:  userId,
		Barrier: req.GetBarrierId(),
		GameRet: mazebarrieruserkafka.GameRetSucc,
		Awards:  getAwards(req.GetFoeExp(), awardMap, equipMap),
	}

	mazebarrieruserkafka.PushMazeBarrierUserRecord(logger, passRecord)

	return nil
}

func getAwards(exp int32, awardMap map[int32]int64, awardEquip map[int32]int32) string {
	awardStr := make([]string, 0)

	awardStr = append(awardStr, fmt.Sprintf("%d:%d", constdef.MazeCommonItemExp, exp))

	for k, v := range awardMap {
		awardStr = append(awardStr, fmt.Sprintf("%d:%d", k, v))
	}

	for k, v := range awardEquip {
		awardStr = append(awardStr, fmt.Sprintf("%d:%d", k, v))
	}
	return strings.Join(awardStr, "_")
}
