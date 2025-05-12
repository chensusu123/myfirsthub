/*
 * @Author: majian
 * @Date: 2025-03-07 22:13:29
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-20 13:47:21
 */
package game

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structsdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecalcattrredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazecommonvalue"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/MazeGame"
	"go.uber.org/zap"
)

func HandleUserAttrMsg(logger fklog.FKLogI, msg *structsdef.DollAttrChgNotify) {
	// msg := &structsdef.DollAttrChgNotify{}
	// err = json.Unmarshal(data, msg)
	// if err != nil {
	// 	logger.ErrorWF("HandleUserAttrMsg Unmarshal", zap.Error(err),
	// 		zap.Int("msg's len", len(data)), zap.Uint64("userId", msg.UserId))
	// 	return
	// }

	userId := msg.UserId
	logger.InfoWF("HandleUserAttrMsg start", zap.Any("msg", msg),
		zap.Uint64("userId", userId))
	if userId <= 0 || len(msg.ChgAttrs) == 0 {
		return
	}
	handleMazeBattleNotify(logger, userId, msg)
	HandleMazeLvUpgradeAttrChgId(logger, userId, msg)

	// _ = attrsMap
	// TODO 处理属性变化包（武力，银子和经验加成属性等）
	handleMazeCommonValueChg(logger, msg.UserId, msg)

	// TODO 处理战斗数据变化(武力 血量 技能属性等)
	return
}

func HasForceAttr(chgAttrs []*structsdef.AttrChgInfo) int32 {
	for _, attr := range chgAttrs {
		if attr.AttrId == constdef.MazeForce {
			return attr.AttrId
		}
	}
	return 0
}

func HasMoneyAttr(chgAttrs []*structsdef.AttrChgInfo) int32 {
	for _, attr := range chgAttrs {
		if attr.AttrId == constdef.MazeMoneyBuff10258 {
			return attr.AttrId
		}
	}
	return 0
}

func HasExpAttr(chgAttrs []*structsdef.AttrChgInfo) int32 {
	for _, attr := range chgAttrs {
		if attr.AttrId == constdef.MazeExpBuff10261 {
			return attr.AttrId
		}
	}
	return 0
}

func handleMazeCommonValueChg(logger fklog.FKLogI, userId uint64, msg *structsdef.DollAttrChgNotify) (err error) {
	var needAttrs []int32
	// 是否有武力属性
	mazeForceId := HasForceAttr(msg.ChgAttrs)
	if mazeForceId > 0 {
		// needAttrs = append(needAttrs, mazeForceId)
		needAttrs = append(needAttrs, constdef.MazeMoneyBuff10258)
		needAttrs = append(needAttrs, constdef.MazeExpBuff10261)
	} else {
		// 是否有银子加成
		mazeMoneyId := HasMoneyAttr(msg.ChgAttrs)
		if mazeMoneyId > 0 {
			needAttrs = append(needAttrs, mazeMoneyId)
			// needAttrs = append(needAttrs, constdef.MazeForce)
		}

		// 是否有经验加成
		mazeExpId := HasExpAttr(msg.ChgAttrs)
		if mazeExpId > 0 {
			needAttrs = append(needAttrs, mazeExpId)
			// needAttrs = append(needAttrs, constdef.MazeForce)
		}
	}

	if len(needAttrs) == 0 {
		return
	}
	attrsMap, err := mazecalcattrredis.BatchGetMazeCalcAttr(logger, msg.UserId, needAttrs)
	if err != nil {
		logger.ErrorWF("handleMazeCommonValueChg get attrs fail",
			zap.Error(err),
			zap.Any("userId", msg.UserId),
			zap.Any("careId", needAttrs))
		return nil
	}

	//武力及属性变化影响的数值有   武力、额外加成
	// force := attrsMap[constdef.MazeForce]

	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("handleMazeCommonValueChg GetUserInfoV2 fail", zap.Error(err))
		return
	}
	level := userInfo.Level

	addMoneyForce, _, addExpForce, err := mazecommonvalue.GetExtraAdditionForce(logger, userId, level, 0)
	if err != nil {
		logger.ErrorWF("handleMazeCommonValueChg GetExtraAdditionForce fail", zap.Error(err))
		return
	}
	addMoneyEquip, okMoney := attrsMap[constdef.MazeMoneyBuff10258]
	addExpEquip, okExp := attrsMap[constdef.MazeExpBuff10261]

	moneyIncome := addMoneyForce + addMoneyEquip
	expIncome := addExpForce + addExpEquip

	commonList := make([]*mazecommonvalue.CommonValueStruct, 0)
	if okMoney {
		moneyCommon := mazecommonvalue.MakeCommonValueList(logger, map[int32]int64{
			// int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_FORCE):      force,
			// int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP_INCOME): expIncome,
			int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_INCOME): moneyIncome},
			map[int32]int32{}, map[int32]string{})
		commonList = append(commonList, moneyCommon...)
	}
	if okExp {
		expCommon := mazecommonvalue.MakeCommonValueList(logger, map[int32]int64{
			// int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_FORCE):      force,
			int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP_INCOME): expIncome},
			// int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_INCOME):     moneyIncome},
			map[int32]int32{}, map[int32]string{})
		commonList = append(commonList, expCommon...)
	}

	if len(commonList) > 0 {
		// forceCommon := mazecommonvalue.MakeCommonValueList(logger, map[int32]int64{
		// 	int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_FORCE): force},
		// 	// int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_EXP_INCOME): expIncome},
		// 	// int32(MazeGame.MAZE_DATA_TYPE_ENUM_MAZE_DATA_TYPE_INCOME):     moneyIncome},
		// 	map[int32]int32{}, map[int32]string{})
		// commonList = append(commonList, forceCommon...)

		mazecommonvalue.SendCommonValueIdPack(logger, userId, commonList)
	}

	return
}

func handleMazeBattleNotify(logger fklog.FKLogI, userId uint64, msg *structsdef.DollAttrChgNotify) {
	if HasBattleAttr(msg.ChgAttrs) {
		userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
		if err != nil {
			logger.ErrorWF("handleMazeBattleNotify GetUserInfoV2 fail", zap.Error(err))
			return
		}
		mazeBattleInfo, err := GetMazeBattleData(logger, userId, userInfo.Barrier)
		if err != nil {
			logger.ErrorWF("handleMazeBattleNotify GetMazeBattleData fail", zap.Error(err))
			return
		}
		SendMazeBarrierChgPack(logger, userId, mazeBattleInfo)
	}
	return
}
