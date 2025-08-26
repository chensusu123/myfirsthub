package tempbuffservice

import (
	"context"
	"fmt"

	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/config/GMazeAttributeV8Cfg"
	"maze_game_server/config/GMazeEnergyAffixV8Cfg"
	"maze_game_server/io/kafka/mazetempbuffchgmsg"
	"maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/model/tempbuffmodel"
	"maze_game_server/pb/common/MazeTempBuff"
	"maze_game_server/pb/server/MazeBuffData"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *service) SelectMazeTempBuff(ctx context.Context, userId uint64, barrierId, level, buffId, buffType int32) ([]*BuffInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
	buffInfo, err := tempbuffmodel.NewTempBuffInfoModel(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "SelectMazeTempBuffRQ GetMazeTempBuff", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}

	if buffInfo == nil {
		logger.CtxError(ctx, "SelectMazeTempBuffRQ buff is nil", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}

	// 检查是否可选，选的buff是否为当次可选buff
	if !s.checkSelectBuff(ctx, logger, level, buffId, buffInfo) {
		return nil, fmt.Errorf("buff信息异常")
	}

	err = s.updateBuffInfo(ctx, userId, barrierId, level, buffId, buffType, buffInfo)
	if err != nil {
		logger.CtxError(ctx, "SelectMazeTempBuffRQ updateBuffInfo failed", zap.Error(err))
		return nil, err
	}

	buffList := s.packShowBuffList(ctx, buffInfo)

	s.pushGroupChange(ctx, logger, userId, buffInfo)

	return buffList, nil
}

func (s *service) checkSelectBuff(ctx context.Context, logger fklog.FKLogI, level, buffId int32, buffInfo *tempbuffmodel.TempBuffInfoModel) bool {
	// 检查是否可选，选的buff是否为当次可选buff
	if level != buffInfo.BuffSequence.Level {
		logger.CtxError(ctx, "checkSelectBuff level unknown", zap.Int32("level", level),
			zap.Int32("needLevel", buffInfo.BuffSequence.Level))
		return false
	}

	for _, id := range buffInfo.BuffSequence.OptionalBuffList {
		if buffId == id {
			return true
		}
	}

	logger.CtxError(ctx, "checkSelectBuff buffId unknown", zap.Int32("buffId", buffId),
		zap.Int32s("buffList", buffInfo.BuffSequence.OptionalBuffList))
	return false
}

func (s *service) updateBuffInfo(ctx context.Context, userId uint64, barrierId, level, buffId, buffType int32,
	buffInfo *tempbuffmodel.TempBuffInfoModel,
) error {
	logger := fklog.ContextAppLogger(ctx)
	areaId := buffInfo.BuffSequence.AreaId
	areaIndex := buffInfo.BuffSequence.AreaIndex
	buffInfo.BuffSequence = &tempbuffmodel.BuffSequence{
		Level: level,
	}

	buffInfo.SelectedBuff = append(buffInfo.SelectedBuff, &tempbuffmodel.SelectedBuffInfo{
		BuffId:    buffId,
		Level:     level,
		Type:      buffType,
		AreaId:    areaId,
		AreaIndex: areaIndex,
	})

	var totalMap map[int32]int64
	totalMap, buffInfo.TotalBuff = s.GetTotalBuff(ctx, buffInfo.SelectedBuff)
	// 更新buff信息
	err := buffInfo.Save(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "updateBuffInfo SetMazeTempBuff failed", zap.Any("info", buffInfo), zap.Error(err))
		return err
	}

	// 推送buff变化信息
	msg := &mazetempbuffchgmsg.MazeTempBuffChangeMsg{
		UserId:  userId,
		StageId: barrierId,
		ChgType: 1,
		ChgDesc: "选择buff",
	}

	// 计算buff变化
	config := GMazeEnergyAffixV8Cfg.GetWithCtx(ctx, buffId)
	if config != nil {
		chgAttrs := make([]*mazetempbuffchgmsg.AttrChgInfo, 0, len(config.Add_attr))
		for id, value := range config.Add_attr {
			chgAttrs = append(chgAttrs, &mazetempbuffchgmsg.AttrChgInfo{
				AttrId: id,
				OldVal: totalMap[id] - value,
				CurVal: totalMap[id],
			})
		}

		msg.ChgAttrs = chgAttrs
	}

	_ = mazetempbuffchgmsg.PushTempBuffChangeMsg(ctx, msg)

	// 同步到buff中心
	s.TempBuffChangeSync(ctx, logger, userId, buffInfo)

	return nil
}

// 同步到buff中心
func (s *service) TempBuffChangeSync(ctx context.Context, logger fklog.FKLogI, userId uint64, buffInfo *tempbuffmodel.TempBuffInfoModel) error {
	forceAttr, err := s.getSelectBuffForceAttr(ctx, buffInfo.TotalBuff)
	if err != nil {
		logger.CtxError(ctx, "updateBuffInfo GetSelectBuffForceAttr failed", zap.Error(err))
		return err
	}
	attrDb := &MazeBuffData.MazeBuffDb{
		MazeRealBuffs: s.packMazeBuff(forceAttr),
		// MazeShowBuffs: PackMazeBuff(showBuff), // todo 现在暂时没有展示武力值
	}

	err = mazebuffinforedis.SaveMazeBuffInfo(logger, userId, constdef.MazeBuffSrcSelectBuffForce, attrDb)
	if err != nil {
		logger.ErrorWF("AddMazeCard SaveMazeBuffInfo failed", zap.Uint64("userId", userId), zap.Error(err))
		return err
	}

	// 推送属性计算消息
	calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{
		UserId: userId,
		// FromServer: fmt.Sprintf("%d %s", fkconfig.EnvVal.ServerType, fkconfig.EnvVal.AppName),
		ChgType: constdef.MazeBuffChgForceValue,
		Session: "buff",
		BuffSrc: constdef.MazeBuffSrcSelectBuffForce,
	}
	mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(ctx, calcAttrNotify)
	return nil
}

func (s *service) GetTotalBuff(ctx context.Context, buffList []*tempbuffmodel.SelectedBuffInfo) (
	map[int32]int64, []*tempbuffmodel.TotalBuffInfo,
) {
	totalMap := make(map[int32]int64)
	for _, info := range buffList {
		// 获取buff实际加成
		config := GMazeEnergyAffixV8Cfg.GetWithCtx(ctx, info.BuffId)
		if config != nil {
			for id, value := range config.Add_attr {
				totalMap[id] += value
			}
		}
	}

	totalList := make([]*tempbuffmodel.TotalBuffInfo, 0, len(totalMap))
	for id, value := range totalMap {
		totalList = append(totalList, &tempbuffmodel.TotalBuffInfo{
			BuffId:    id,
			BuffValue: value,
		})
	}

	return totalMap, totalList
}

func (s *service) getSelectBuffForceAttr(ctx context.Context, buffInfo []*tempbuffmodel.TotalBuffInfo) (forceAttr map[int32]int64, err error) {
	if buffInfo == nil {
		return
	}
	forceAttr = make(map[int32]int64)
	for _, v := range buffInfo {
		// 读属性表
		cfg := GMazeAttributeV8Cfg.GetWithCtx(ctx, v.BuffId)
		if cfg == nil {
			err = fmt.Errorf("GetSelectBuffForceAttr GetMazeAttributeFormulaV8Cfg nil, attrID: %d", v.BuffId)
			return
		}
		if cfg.Type == forceAttrType {
			forceAttr[cfg.Id] += v.BuffValue
		}
	}
	return
}

func (s *service) packMazeBuff(buffMap map[int32]int64) []*MazeBuffData.MazeBuffAttr {
	if len(buffMap) == 0 {
		return nil
	}

	buffList := make([]*MazeBuffData.MazeBuffAttr, 0, len(buffMap))
	for id, value := range buffMap {
		buffList = append(buffList, &MazeBuffData.MazeBuffAttr{
			AttrId:  proto.Int32(id),
			AttrVal: proto.Int64(value),
		})
	}

	return buffList
}

const forceAttrType = 5

// 推送组变化包
func (s *service) pushGroupChange(ctx context.Context, logger fklog.FKLogI, userId uint64, buffModel *tempbuffmodel.TempBuffInfoModel) {
	groupList, err := s.getGroupList(ctx, logger, buffModel)
	if err != nil {
		logger.CtxError(ctx, "getGroupList failed", zap.Error(err))
		return
	}
	res := &MazeTempBuff.TempBuffGroupChangeID{}
	res.BuffGroupList = make([]*MazeTempBuff.TempBuffGroupInfo, 0, len(groupList))
	for _, i := range groupList {
		res.BuffGroupList = append(res.BuffGroupList, &MazeTempBuff.TempBuffGroupInfo{
			BuffId: proto.Int32(i.BuffId),
			Count:  proto.Int64(int64(i.Count)),
		})
	}

	online.PushWithContext(ctx, userId, 10642, res)
}
