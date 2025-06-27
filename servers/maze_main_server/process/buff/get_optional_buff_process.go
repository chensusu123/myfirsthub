package buff

import (
	"fmt"
	"math/rand"
	"maze_game_server/common/errors"
	"maze_game_server/excel/mazebarriesv8config"
	"maze_game_server/excel/mazeconfigv8config"
	"maze_game_server/excel/mazeenergyaffixfrontv8config"
	"maze_game_server/excel/mazeenergyaffixlibraryv8config"
	"maze_game_server/excel/mazeenergyaffixlvv8config"
	"maze_game_server/excel/mazeenergyaffixrandrulev8config"
	"maze_game_server/excel/mazeenergylevelv8config"
	"maze_game_server/excel/mazeenergyresetcostv8config"
	"maze_game_server/io/redis/mazetempbuffredis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeTempBuff"
	"maze_game_server/pb/common/MessageType"
	"maze_game_server/pb/server/MazeTempBuffSvr"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/21 15:27
 * @Description: 查询迷宫可选buff列表
 */

func (b *Buff) GetOptionalMazeTempBuffListRQ_10435_10436(s *session.Session, req *MazeTempBuff.GetOptionalMazeTempBuffListRQ) (err error) {
	defer fkprometheus.InfoPMT("GetOptionalMazeTempBuffListRQ")()

	start := time.Now()

	logger := log.Clone("Buff", uint64(s.UID()), 0)
	res := &MazeTempBuff.GetOptionalMazeTempBuffListRS{}
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.StageId = req.StageId
	res.Level = req.Level
	res.Type = req.Type
	res.AreaId = req.AreaId
	defer func() {
		err = s.Response(res)
		logger.InfoWF("GetOptionalMazeTempBuffListRQ end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId, stageId, level, buffType, areaId := uint64(s.UID()), req.GetStageId(), req.GetLevel(), int32(req.GetType()), req.GetAreaId()
	if userId == 0 || stageId == 0 || level == 0 || areaId == 0 || buffType == 0 {
		logger.WarnWF("GetOptionalMazeTempBuffListRQ args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
		return nil
	}
	if buffType != int32(MazeTempBuff.Type_UP_LEVEL) && buffType != int32(MazeTempBuff.Type_USE_ITEM) {
		logger.ErrorWF("GetOptionalMazeTempBuffListRQ buffType args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("buff类型参数错误")
		return nil
	}

	// todo 检查用户是不是小程序用户

	buffInfo, err := mazetempbuffredis.GetMazeTempBuff(logger, userId, stageId)
	if err != nil {
		logger.ErrorWF("GetOptionalMazeTempBuffListRQ GetMazeTempBuff failed", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("获取用户buff信息失败")
		return err
	}

	if buffInfo == nil {
		// 初始化队列
		buffInfo = &MazeTempBuffSvr.TempBuffInfo{
			BuffSequence: &MazeTempBuffSvr.BuffSequence{
				Index: proto.Int32(1),
			},
		}
	}

	if len(buffInfo.GetBuffSequence().GetSelectBuffList()) == 0 {
		// 生成可选buff列表
		errInfo := getOptionalBuffList(logger, userId, stageId, level, buffType, areaId, buffInfo)
		if errInfo != nil {
			logger.ErrorWF("GetOptionalMazeTempBuffListRQ getOptionalBuffList", zap.Int32("stageId", stageId),
				zap.Any("info", buffInfo), zap.Any("errInfo", errInfo))
			res.ErrInfo = errInfo
			return err
		}
	}

	if buffInfo.GetBuffSequence().GetIndex() != level {
		logger.WarnWF("GetOptionalMazeTempBuffListRQ level is not need", zap.Int32("level", level),
			zap.Int32("needLevel", buffInfo.GetBuffSequence().GetIndex()))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("能力等级异常")
		return nil
	}

	res.OptionalBuffInfo = packOptionalInfo(logger, buffInfo)
	if res.GetOptionalBuffInfo() != nil {
		return nil
	}

	if len(buffInfo.GetBuffSequence().GetSelectBuffList()) == 0 {
		// 无buff可选
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("buff已全部选择")
		return nil
	}

	res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("buff配置异常")
	return nil
}

// 获取可选buff列表
func getOptionalBuffList(logger fklog.FKLogI, userId uint64, stageId, level, buffType, areaId int32,
	buffInfo *MazeTempBuffSvr.TempBuffInfo) *MessageType.ErrorInfo {
	if level < buffInfo.GetBuffSequence().GetIndex() {
		logger.WarnWF("getOptionalBuffList level already select", zap.Int32("level", level),
			zap.Int32("needLevel", buffInfo.GetBuffSequence().GetIndex()))
		return errors.COMMON_ERROR_TIPS.Wrap("当前等级已选择过buff")
	}

	// 校验选择buff数量
	stageConfig := mazebarriesv8config.GetStageConfig(stageId)
	if stageConfig == nil {
		logger.WarnWF("getOptionalBuffList stage config unknown", zap.Int32("stageId", stageId))
		return errors.COMMON_ERROR_TIPS.Wrap("关卡配置异常")
	}

	// 校验选择buff数量
	// configId := mazeenergylevelv8config.GetKey(stageConfig.Energy_id, level)
	// config := mazeenergylevelv8config.GetEnergyLevelConfig(configId)
	// if config == nil {
	// 	logger.WarnWF("getOptionalBuffList level config unknown", zap.Int32("level", level))
	// 	return errors.COMMON_ERROR_TIPS.Wrap("当前等级无法选择buff")
	// }

	energyId, ok := stageConfig.Energy_id[areaId]
	if !ok || energyId <= 0 {
		logger.WarnWF("getOptionalBuffList energyId unknown", zap.Bool("findEnergyId", ok), zap.Int32("areaId", areaId))
		return errors.COMMON_ERROR_TIPS.Wrap("找不到当前区域能力配置")
	}

	if buffType == int32(MazeTempBuff.Type_UP_LEVEL) {
		checkErr := checkUpLevelSelectBuff(logger, level, energyId, buffInfo)
		if checkErr != nil && checkErr.GetErrCode() != errors.NO_ERROR_CODE {
			return checkErr
		}
	} else if buffType == int32(MazeTempBuff.Type_USE_ITEM) {
		checkErr := checkUseItemLevelSelectBuff(logger, level, energyId, buffInfo)
		if checkErr != nil && checkErr.GetErrCode() != errors.NO_ERROR_CODE {
			return checkErr
		}
	}

	// var count int32
	// for _, info := range buffInfo.GetSelectedBuff() {
	// 	if info.GetLevel() == level && info.GetType() == int32(MazeTempBuff.Type_UP_LEVEL) {
	// 		count++
	// 	}
	// }
	//
	// if count >= config.Energy_select {
	// 	logger.WarnWF("getOptionalBuffList buff count select max", zap.Int32("count", count),
	// 		zap.Int32("maxCount", config.Energy_select))
	// 	return errors.COMMON_ERROR_TIPS.Wrap("当前等级已选择完buff")
	// }

	buffInfo.BuffSequence.Index = proto.Int32(level)
	// 生成可选的buff列表
	ruleId, ok := stageConfig.Energy_affix_rand_rule[areaId]
	if !ok || ruleId <= 0 {
		logger.ErrorWF("getOptionalBuffList stageConfig.Energy_affix_rand_rule not found",
			zap.Int32("stageId", stageId), zap.Int32("level", level), zap.Int32("areaId", areaId))
		return errors.COMMON_ERROR_TIPS.Wrap("找不到当前区域能力随机规则")
	}
	configId := mazeenergyaffixrandrulev8config.GetKey(ruleId, level)
	buffList, err := createOptionalBuffList(logger, buffInfo, configId)
	if err != nil {
		logger.ErrorWF("getOptionalBuffList createOptionalBuffList failed", zap.Error(err))
		return errors.COMMON_ERROR_TIPS.Wrap("获取可选buff失败")
	}

	buffInfo.BuffSequence.SelectBuffList = buffList
	err = mazetempbuffredis.SetMazeTempBuff(logger, userId, stageId, buffInfo)
	if err != nil {
		logger.ErrorWF("getOptionalBuffList SetMazeTempBuff failed", zap.Int32("stageId", stageId),
			zap.Any("info", buffInfo), zap.Error(err))
		return errors.COMMON_ERROR_TIPS.Wrap("保存buff信息失败")
	}

	return nil
}

func checkUpLevelSelectBuff(logger fklog.FKLogI, level, energyID int32, buffInfo *MazeTempBuffSvr.TempBuffInfo) *MessageType.ErrorInfo {
	configId := mazeenergylevelv8config.GetKey(energyID, level)
	config := mazeenergylevelv8config.GetEnergyLevelConfig(configId)
	if config == nil {
		logger.WarnWF("checkUpLevelSelectBuff level config unknown", zap.Int32("level", level), zap.Int32("configId", configId))
		return errors.COMMON_ERROR_TIPS.Wrap("当前等级无法选择buff")
	}
	var count int32
	for _, info := range buffInfo.GetSelectedBuff() {
		if info.GetLevel() == level && info.GetType() == int32(MazeTempBuff.Type_UP_LEVEL) {
			count++
		}
	}

	if count >= config.Energy_select {
		logger.WarnWF("checkUpLevelSelectBuff buff count select max", zap.Int32("count", count),
			zap.Int32("maxCount", config.Energy_select))
		return errors.COMMON_ERROR_TIPS.Wrap("当前等级已选择完buff")
	}
	return nil
}

func checkUseItemLevelSelectBuff(logger fklog.FKLogI, level, energyID int32, buffInfo *MazeTempBuffSvr.TempBuffInfo) *MessageType.ErrorInfo {
	configId := mazeenergylevelv8config.GetKey(energyID, level)
	config := mazeenergylevelv8config.GetEnergyLevelConfig(configId)
	if config == nil {
		logger.WarnWF("checkUseItemLevelSelectBuff level config unknown", zap.Int32("configId", configId))
		return errors.COMMON_ERROR_TIPS.Wrap("读取buff能力配置失败")
	}
	var count int32
	for _, info := range buffInfo.GetSelectedBuff() {
		if info.GetLevel() == level && info.GetType() == int32(MazeTempBuff.Type_USE_ITEM) {
			count++
		}
	}

	if count >= config.Energy_item_select {
		logger.WarnWF("checkUseItemLevelSelectBuff buff count select max", zap.Int32("count", count),
			zap.Int32("maxCount", config.Energy_item_select))
		return errors.COMMON_ERROR_TIPS.Wrap("道具选择buff次数已用完")
	}
	return nil
}

func packOptionalInfo(logger fklog.FKLogI, info *MazeTempBuffSvr.TempBuffInfo) *MazeTempBuff.OptionalBuffInfo {
	optionalInfo := &MazeTempBuff.OptionalBuffInfo{
		SelectBuffList: packSelectBuffList(logger, info.GetBuffSequence().GetSelectBuffList()),
		SelectBuffTime: proto.Int32(mazeconfigv8config.GetBuffSelectTime()),
	}

	if len(optionalInfo.GetSelectBuffList()) == 0 {
		return nil
	}

	// 是否可以刷新
	config := mazeenergyresetcostv8config.GetEnergyResetCostConfig(info.GetBuffSequence().GetRefreshCount() + 1)
	if config == nil {
		return optionalInfo
	}

	optionalInfo.IsRefresh = proto.Int32(1)
	for itemId, count := range config.Cost {
		optionalInfo.Cost = append(optionalInfo.Cost, &MazeCommon.MazeItem{
			ItemId: proto.Int32(itemId),
			Count:  proto.Int64(count),
		})
	}

	return optionalInfo
}

func packSelectBuffList(logger fklog.FKLogI, buffList []int32) []*MazeTempBuff.MazeBuffInfo {
	if len(buffList) == 0 {
		return nil
	}

	showBuffList := make([]*MazeTempBuff.MazeBuffInfo, 0, len(buffList))
	for _, buffId := range buffList {
		if buffId == 0 {
			continue
		}

		if showBuff := packShowBuff(logger, buffId, 1); showBuff != nil {
			showBuffList = append(showBuffList, showBuff)
		}
	}

	return showBuffList
}

// 生成可选的buff列表
func createOptionalBuffList(logger fklog.FKLogI, buffInfo *MazeTempBuffSvr.TempBuffInfo, configId int32) ([]int32, error) {
	randConfig := mazeenergyaffixrandrulev8config.GetEnergyAffixRandRuleConfig(configId)
	if randConfig == nil {
		return nil, nil
	}

	optionalMap := make(map[int32]struct{})
	num := mazeconfigv8config.GetBuffSelectCount()
	for i := int64(1); i <= num; i++ {
		var libraryId int32
		switch i {
		case 1:
			libraryId, _ = randLibraryId(randConfig.Pos_1_lib)
		case 2:
			libraryId, _ = randLibraryId(randConfig.Pos_2_lib)
		case 3:
			libraryId, _ = randLibraryId(randConfig.Pos_3_lib)
		case 4:
			libraryId, _ = randLibraryId(randConfig.Pos_4_lib)
		case 5:
			libraryId, _ = randLibraryId(randConfig.Pos_5_lib)
		case 6:
			libraryId, _ = randLibraryId(randConfig.Pos_6_lib)
		default:
			logger.WarnWF("createOptionalBuffList unknown id", zap.Int64("num", i))
			return nil, nil
		}

		// 随机库id
		affixList, certainly_list := mazeenergyaffixlibraryv8config.GetEnergyLibraryAffixList(libraryId)
		if len(affixList) == 0 {
			return nil, errors.New("affixList is nil")
		}

		// 过滤不可选择的词条
		optionalList, totalWeight := filterBuffList(buffInfo, optionalMap, affixList, certainly_list)
		// 随机选择个词条
		buffId, weight := randomId(optionalList, totalWeight)
		if buffId != 0 {
			optionalMap[buffId] = struct{}{}
		}

		logger.DebugWF("createOptionalBuffList random", zap.Int64("i", i), zap.Int32("libraryId", libraryId),
			zap.Int32s("affixList", affixList), zap.Any("optionalList", optionalList),
			zap.Int32("weight", weight), zap.Int32("id", buffId))
	}

	var optionalList []int32
	for optionId := range optionalMap {
		optionalList = append(optionalList, optionId)
	}

	logger.InfoWF("createOptionalBuffList end", zap.Int32s("optionalList", optionalList))
	return optionalList, nil
}

func randLibraryId(libraryMap map[int32]int32) (int32, int32) {
	var (
		weightList  []*WeightInfo
		totalWeight int32
	)
	for id, weight := range libraryMap {
		if weight == 0 {
			continue
		}

		weightList = append(weightList, &WeightInfo{
			Id:     id,
			Weight: weight,
		})
		totalWeight += weight
	}

	return randomId(weightList, totalWeight)
}

type WeightInfo struct {
	Id     int32
	Weight int32
}

func filterBuffList(buffInfo *MazeTempBuffSvr.TempBuffInfo, optionalMap map[int32]struct{}, buffList []int32, ce_buffList []int32) ([]*WeightInfo, int32) {
	buffMap := make(map[int32]int32)
	for _, info := range buffInfo.GetSelectedBuff() {
		buffMap[info.GetBuffId()] += 1
	}

	var (
		optionalList []*WeightInfo
		totalWeight  int32
	)

	// 先添加必选buff
	for _, buffId := range ce_buffList {
		if _, ok := optionalMap[buffId]; ok {
			continue
		}

		buffWeight, _ := getOptionBuffWeightInfo(buffId, buffMap)
		if buffWeight == nil {
			continue
		}

		optionalList = append(optionalList, buffWeight)
		totalWeight += buffWeight.Weight
	}

	if len(optionalList) > 0 {
		return optionalList, totalWeight
	}

	for _, buffId := range buffList {
		if _, ok := optionalMap[buffId]; ok {
			continue
		}

		buffWeight, _ := getOptionBuffWeightInfo(buffId, buffMap)
		if buffWeight == nil {
			continue
		}

		optionalList = append(optionalList, buffWeight)
		totalWeight += buffWeight.Weight
	}

	return optionalList, totalWeight
}

// 获取可选buff的权重信息
func getOptionBuffWeightInfo(buffId int32, buffMap map[int32]int32) (*WeightInfo, error) {
	buffConfig := mazeenergyaffixlvv8config.GetAffixConfig(buffId)
	if buffConfig == nil {
		return nil, errors.New("能力词条配置不存在")
	}

	if buffConfig.Weight == 0 {
		return nil, errors.New("词条配置权重为0")
	}

	// 检查选择数量
	optionalCount := buffConfig.Use_num_max - buffMap[buffId]
	if optionalCount <= 0 {
		return nil, fmt.Errorf("最大选择数量: %d", buffConfig.Use_num_max)
	}

	// 检查条件
	for _, frontId := range buffConfig.Font_affix_condition {
		if frontId == 0 {
			continue
		}

		isOk := checkFrontCondition(frontId, buffMap)
		if !isOk {
			return nil, errors.New("前置条件不满足")
		}
	}

	return &WeightInfo{
		Id:     buffId,
		Weight: optionalCount * buffConfig.Weight,
	}, nil
}

// 检查前置条件
func checkFrontCondition(frontId int32, buffMap map[int32]int32) bool {
	frontConfig := mazeenergyaffixfrontv8config.GetMazeEnergyAffixFrontConfig(frontId)
	if frontConfig == nil {
		return false
	}

	var count int32
	for _, affixId := range frontConfig.Affix_id_set {
		if _, ok := buffMap[affixId]; !ok {
			continue
		}

		count++

		if count >= frontConfig.Must_num {
			break
		}
	}

	if count < frontConfig.Must_num {
		return false
	}
	return true
}

func randomId(optionalList []*WeightInfo, totalWeight int32) (int32, int32) {
	if len(optionalList) == 0 || totalWeight <= 0 {
		return 0, 0
	}

	weight := int32(rand.Intn(int(totalWeight)))
	var curWeight int32
	for _, info := range optionalList {
		curWeight += info.Weight
		if weight < curWeight {
			return info.Id, weight
		}
	}

	return 0, weight
}
