package buff

import (
	"math/rand"
	"time"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/common/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeTempBuff"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeTempBuffSvr"
	"gitlab.ifreetalk.com/maze-plate/protodef/MessageType"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazebarriesv8config"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeconfigv8config"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeenergyaffixfrontv8config"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeenergyaffixlibraryv8config"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeenergyaffixlvv8config"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeenergyaffixrandrulev8config"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeenergylevelv8config"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeenergyresetcostv8config"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazetempbuffredis"
	"go.uber.org/zap"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/21 15:27
 * @Description: 查询迷宫可选buff列表
 */

func GetOptionalMazeTempBuffListRQ(logger fknet.TCPContext, shardingID uint64, request, response proto.Message) error {
	defer fkprometheus.InfoPMT("GetOptionalMazeTempBuffListRQ")()
	start := time.Now()
	req := request.(*MazeTempBuff.GetOptionalMazeTempBuffListRQ)
	res := response.(*MazeTempBuff.GetOptionalMazeTempBuffListRS)
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.StageId = req.StageId
	res.Level = req.Level
	defer func() {
		logger.InfoWF("GetOptionalMazeTempBuffListRQ end", zap.Any("req", req), zap.Any("res", res),
			zap.Duration("costTime", time.Now().Sub(start)))
	}()

	userId, stageId, level := shardingID, req.GetStageId(), req.GetLevel()
	if userId == 0 || stageId == 0 || level == 0 {
		logger.WarnWF("GetOptionalMazeTempBuffListRQ args error", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("参数错误")
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
		errInfo := getOptionalBuffList(logger, userId, stageId, level, buffInfo)
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
func getOptionalBuffList(logger fklog.FKLogI, userId uint64, stageId, level int32,
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
	configId := mazeenergylevelv8config.GetKey(stageConfig.Energy_id, level)
	config := mazeenergylevelv8config.GetEnergyLevelConfig(configId)
	if config == nil {
		logger.WarnWF("getOptionalBuffList level config unknown", zap.Int32("level", level))
		return errors.COMMON_ERROR_TIPS.Wrap("当前等级无法选择buff")
	}

	var count int32
	for _, info := range buffInfo.GetSelectedBuff() {
		if info.GetLevel() == level {
			count++
		}
	}

	if count >= config.Energy_select {
		logger.WarnWF("getOptionalBuffList buff count select max", zap.Int32("count", count),
			zap.Int32("maxCount", config.Energy_select))
		return errors.COMMON_ERROR_TIPS.Wrap("当前等级已选择完buff")
	}

	buffInfo.BuffSequence.Index = proto.Int32(level)
	// 生成可选的buff列表
	configId = mazeenergyaffixrandrulev8config.GetKey(stageConfig.Energy_affix_rand_rule, level)
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

		buffWeight := getOptionBuffWeightInfo(buffId, buffMap)
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

		buffWeight := getOptionBuffWeightInfo(buffId, buffMap)
		if buffWeight == nil {
			continue
		}

		optionalList = append(optionalList, buffWeight)
		totalWeight += buffWeight.Weight
	}

	return optionalList, totalWeight
}

// 获取可选buff的权重信息
func getOptionBuffWeightInfo(buffId int32, buffMap map[int32]int32) *WeightInfo {
	buffConfig := mazeenergyaffixlvv8config.GetAffixConfig(buffId)
	if buffConfig == nil {
		return nil
	}

	if buffConfig.Weight == 0 {
		return nil
	}

	// 检查选择数量
	optionalCount := buffConfig.Use_num_max - buffMap[buffId]
	if optionalCount <= 0 {
		return nil
	}

	// 检查条件
	for _, frontId := range buffConfig.Font_affix_condition {
		if frontId == 0 {
			continue
		}

		isOk := checkFrontCondition(frontId, buffMap)
		if !isOk {
			return nil
		}
	}

	return &WeightInfo{
		Id:     buffId,
		Weight: optionalCount * buffConfig.Weight,
	}
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
