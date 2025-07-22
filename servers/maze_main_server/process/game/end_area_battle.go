package game

import (
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/structsdef"
	"maze_game_server/config/GMazeSkillInfoV8Cfg"
	"maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/io/redis/mazetempbuffredis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeAIBattle"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/server/MazeTempBuffSvr"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (g *Game) OnEndAreaBattleRQ_10525_10526(s *session.Session, req *MazeGame.EndAreaBattleRQ) (err error) {
	defer fkprometheus.InfoPMT("OnEndAreaBattleRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)
	res := &MazeGame.EndAreaBattleRS{}

	logger.InfoWF("OnEndAreaBattleRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnEndAreaBattleRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	return err
	//temp block logic error

	userId := uint64(s.UID())

	if req.GetStageId() <= 0 {
		logger.ErrorWF("OnEndAreaBattleRQ req invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id未设置")
		return
	}
	tempBuffInfo, err := mazetempbuffredis.GetMazeTempBuff(logger, userId, req.GetStageId())
	if err != nil {
		logger.ErrorWF("OnEndAreaBattleRQ GetMazeTempBuff failed", zap.Error(err), zap.Any("req", req))
		return
	}

	if tempBuffInfo == nil {
		logger.InfoWF("OnEndAreaBattleRQ user tempBuffInfo is nil", zap.Any("req", req))
		return
	}

	err = mazetempbuffredis.DelMazeTempBuff(logger, userId, req.GetStageId())
	if err != nil {
		logger.ErrorWF("OnEndAreaBattleRQ DelMazeTempBuff failed", zap.Error(err), zap.Any("req", req))
		return
	}

	// 通知
	// 删除临时buff武力属性
	err = mazebuffinforedis.DelMazeBuffBySrc(logger, userId, constdef.MazeBuffSrcSelectBuffForce)
	if err != nil {
		logger.ErrorWF("OnEndAreaBattleRQ DelMazeBuffBySrc failed", zap.Uint64("userId", userId), zap.Error(err))
		return
	}
	// 推送属性计算消息
	calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{
		UserId: userId,
		// FromServer: fmt.Sprintf("%d %s", fkconfig.EnvVal.ServerType, fkconfig.EnvVal.AppName),
		ChgType: constdef.MazeBuffChgForceValue,
		Session: "buff",
		BuffSrc: constdef.MazeBuffSrcSelectBuffForce,
	}
	mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, calcAttrNotify)

	// 通知删除技能
	equipSkillInfoChange, changed, err := GetTempBuffSkillInfoChange(logger, userId, tempBuffInfo)
	if err != nil {
		logger.ErrorWF("OnEndAreaBattleRQ GetTempBuffSkillInfoChange fail",
			zap.Error(err),
			zap.Uint64("userId", userId),
			zap.Any("req", req),
		)
	} else if changed {
		defer func() {
			online.Push(logger, userId, 10510, equipSkillInfoChange)
		}()
	}

	return nil
}

// GetTempBuffSkillInfoChange 获取临时buff变化引起的技能变化
func GetTempBuffSkillInfoChange(logger fklog.FKLogI, userID uint64, tempBuffInfo *MazeTempBuffSvr.TempBuffInfo) (ret *MazeAIBattle.MazeUserSkillInfoChangeID, changed bool, err error) {
	userAttrMap, err := GetUserAttrMap(logger, userID)
	if err != nil {
		logger.ErrorWF("OnEndAreaBattleRQ GetUserAttrMap err", zap.Error(err))
		return nil, false, err
	}

	for _, buffInfo := range tempBuffInfo.TotalBuff {
		userAttrMap[buffInfo.GetBuffId()] += buffInfo.GetBuffValue()
	}

	ret = &MazeAIBattle.MazeUserSkillInfoChangeID{}
	// 删除临时buff会删除技能
	oldChanged := false
	ret.DelSkillInfoList, oldChanged, err = GetUserSkillTotalInfo(logger, userAttrMap, tempBuffInfo)
	if oldChanged {
		changed = true
	}
	return ret, changed, nil
}

func GetUserSkillTotalInfo(logger fklog.FKLogI, userAttrMap map[int32]int64, tempBuff *MazeTempBuffSvr.TempBuffInfo) (skillTotalInfo *MazeAIBattle.MazeAISkillTotalInfo, changed bool, err error) {
	skillTotalInfo = &MazeAIBattle.MazeAISkillTotalInfo{}
	attrMap := make(map[int32]int64)
	for _, buff := range tempBuff.GetTotalBuff() {
		attrMap[buff.GetBuffId()] += buff.GetBuffValue()
	}
	for _, cfg := range GMazeSkillInfoV8Cfg.GetAll() {
		if cfg.Skill_attr_id > 0 {
			_, ok := attrMap[cfg.Skill_attr_id]
			// 判断是否激活技能
			if ok {
				skillInfo, _, err := GetUserBattleSkillInfo(logger, cfg.Id, userAttrMap)
				if err != nil {
					return nil, false, err
				}
				skillTotalInfo.SkillInfoList = append(skillTotalInfo.SkillInfoList, skillInfo)
				changed = true
			}
		}
	}
	return skillTotalInfo, changed, nil
}
