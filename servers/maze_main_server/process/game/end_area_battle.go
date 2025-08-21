package game

import (
	"context"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeSkillInfoV8Cfg"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/passareamodel"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeAIBattle"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/server/MazeTempBuffSvr"
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
	userId := uint64(s.UID())

	if req.GetStageId() <= 0 || req.GetAreaId() <= 0 {
		logger.ErrorWF("OnEndAreaBattleRQ req invalid", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id未设置")
		return
	}

	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnEndAreaBattleRQ GetUserInfo fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if userInfo.Barrier != req.GetStageId() {
		logger.ErrorWF("OnEndAreaBattleRQ barrier err", zap.Any("req", req), zap.Any("barrier", userInfo.Barrier))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡id错误")
		return
	}

	passAreaModel, err := passareamodel.NewPassAreaModel(context.TODO(), userId, req.GetStageId())
	if err != nil {
		logger.ErrorWF("OnEndAreaBattleRQ GetBarrierPassArea fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	exist := false
	for _, area := range passAreaModel.PassAreaList {
		if area.AreaId == req.GetAreaId() && area.AreaIndex == req.GetAreaIndex() {
			exist = true
		}
	}
	if exist {
		logger.InfoWF("OnEndAreaBattleRQ area already passed", zap.Int32("stageId", req.GetStageId()),
			zap.Int32("areaId", req.GetAreaId()), zap.Int32("areaIndex", req.GetAreaIndex()))
		return
	}
	passAreaModel.PassAreaList = append(passAreaModel.PassAreaList, &passareamodel.PassAreaInfo{
		AreaId:    req.GetAreaId(),
		AreaIndex: req.GetAreaIndex(),
	})

	err = passAreaModel.Save(context.TODO(), userId, req.GetStageId())
	if err != nil {
		logger.ErrorWF("OnEndAreaBattleRQ SetBarrierPassArea fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
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
