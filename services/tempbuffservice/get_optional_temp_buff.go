package tempbuffservice

import (
	"context"
	"fmt"
	"math/rand"
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeEnergyAffixFrontV8Cfg"
	"maze_game_server/config/GMazeEnergyAffixLibraryV8Cfg"
	"maze_game_server/config/GMazeEnergyAffixRandRuleV8Cfg"
	"maze_game_server/config/GMazeEnergyAffixV8Cfg"
	"maze_game_server/config/GMazeEnergyLevelV8Cfg"
	"maze_game_server/excel/mazeconfigv8"
	"maze_game_server/excel/mazeconfigv8config"
	"maze_game_server/excel/mazeenergyaffixlibrarycfgex"
	"maze_game_server/excel/mazeenergyaffixrandrulev8config"
	"maze_game_server/excel/mazeenergylevelv8config"
	"maze_game_server/excel/mazeenergyresetcostv8config"
	"maze_game_server/model/tempbuffmodel"
	"maze_game_server/pb/common/MazeTempBuff"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) GetOptionalTempBuffList(ctx context.Context, userId uint64, barrierId, level, buffType, areaId, areaIndex, attrMask int32) (*OptionalBuffInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
	buffInfo, err := tempbuffmodel.NewTempBuffInfoModel(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "GetOptionalMazeTempBuffListRQ GetMazeTempBuff failed", zap.Error(err))
		return nil, fmt.Errorf("获取用户buff信息失败")
	}

	if buffInfo.BuffSequence == nil {
		// 没有buff 初始化队列
		buffInfo = &tempbuffmodel.TempBuffInfoModel{
			BuffSequence: &tempbuffmodel.BuffSequence{
				Level: 1,
			},
		}
	}

	// 生成可选buff列表
	err = s.genOptionalBuffList(ctx, userId, barrierId, level, buffType, areaId, areaIndex, attrMask, buffInfo)
	if err != nil {
		logger.CtxError(ctx, "GetOptionalMazeTempBuffListRQ getOptionalBuffList", zap.Int32("barrierId", barrierId),
			zap.Any("info", buffInfo), zap.Any("err", err.Error()))
		return nil, err
	}

	optionalBuffInfo := s.packOptionalInfo(ctx, buffInfo)
	if optionalBuffInfo != nil {
		// 能力等级
		optionalBuffInfo.Level = buffInfo.BuffSequence.Level
		return optionalBuffInfo, nil
	}

	if len(buffInfo.BuffSequence.OptionalBuffList) == 0 {
		// 无buff可选
		return nil, fmt.Errorf("buff已全部选择")
	}

	return nil, fmt.Errorf("buff配置异常")
}

// 生成可选buff列表
func (s *service) genOptionalBuffList(ctx context.Context, userId uint64, barrierId, level, buffType, areaId, areaIndex, attrMask int32,
	buffInfo *tempbuffmodel.TempBuffInfoModel) error {
	logger := fklog.ContextAppLogger(ctx)
	if level < buffInfo.BuffSequence.Level {
		logger.CtxWarn(ctx, "genOptionalBuffList level already select", zap.Int32("level", level),
			zap.Int32("needLevel", buffInfo.BuffSequence.Level))
		return fmt.Errorf("当前等级已选择过buff")
	}

	// 校验选择buff数量
	barrierConfig := GMazeBarriesV8Cfg.GetWithCtx(ctx, barrierId)
	if barrierConfig == nil {
		logger.CtxWarn(ctx, "genOptionalBuffList barrier config unknown", zap.Int32("barrierId", barrierId))
		return fmt.Errorf("关卡配置异常")
	}

	energyId, ok := barrierConfig.Energy_id[areaId]
	if !ok || energyId <= 0 {
		logger.CtxError(ctx, "genOptionalBuffList energyId unknown", zap.Bool("findEnergyId", ok), zap.Int32("areaId", areaId))
		return fmt.Errorf("找不到当前区域能力配置")
	}

	if buffType == int32(MazeTempBuff.Type_UP_LEVEL) {
		checkErr := s.checkUpLevelSelectBuff(ctx, logger, level, energyId, buffInfo)
		if checkErr != nil {
			return checkErr
		}
	} else if buffType == int32(MazeTempBuff.Type_USE_ITEM) {
		checkErr := s.checkUseItemLevelSelectBuff(ctx, logger, level, energyId, buffInfo)
		if checkErr != nil {
			return checkErr
		}
	}

	buffInfo.BuffSequence.Level = level
	// 生成可选的buff列表
	buffList, err := s.createOptionalBuffList(ctx, buffInfo, level, areaId, attrMask, barrierConfig)
	if err != nil {
		logger.CtxError(ctx, "genOptionalBuffList createOptionalBuffList failed", zap.Error(err))
		return fmt.Errorf("创建可选buff列表失败")
	}

	buffInfo.BuffSequence.OptionalBuffList = buffList
	buffInfo.BuffSequence.AreaId = areaId
	buffInfo.BuffSequence.AreaIndex = areaIndex
	err = buffInfo.Save(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "genOptionalBuffList SetMazeTempBuff failed", zap.Int32("barrierId", barrierId),
			zap.Any("info", buffInfo), zap.Error(err))
		return fmt.Errorf("保存buff信息失败")
	}

	return nil
}

func (s *service) checkUpLevelSelectBuff(ctx context.Context, logger fklog.FKLogI, level, energyID int32, buffInfo *tempbuffmodel.TempBuffInfoModel) error {
	configId := mazeenergylevelv8config.GetKey(energyID, level)
	config := GMazeEnergyLevelV8Cfg.GetWithCtx(ctx, configId)
	if config == nil {
		logger.CtxError(ctx, "checkUpLevelSelectBuff level config unknown", zap.Int32("level", level), zap.Int32("configId", configId))
		return fmt.Errorf("当前等级无法选择buff")
	}
	var count int32
	for _, info := range buffInfo.SelectedBuff {
		if info.Level == level && info.Type == int32(MazeTempBuff.Type_UP_LEVEL) {
			count++
		}
	}

	if count >= config.Energy_select {
		logger.CtxError(ctx, "checkUpLevelSelectBuff buff count select max", zap.Int32("count", count),
			zap.Int32("maxCount", config.Energy_select))
		return fmt.Errorf("当前等级已选择完buff")
	}
	return nil
}

func (s *service) checkUseItemLevelSelectBuff(ctx context.Context, logger fklog.FKLogI, level, energyID int32, buffInfo *tempbuffmodel.TempBuffInfoModel) error {
	configId := mazeenergylevelv8config.GetKey(energyID, level)
	config := GMazeEnergyLevelV8Cfg.GetWithCtx(ctx, configId)
	if config == nil {
		logger.CtxError(ctx, "checkUseItemLevelSelectBuff level config unknown", zap.Int32("configId", configId))
		return fmt.Errorf("读取buff能力配置失败")
	}
	var count int32
	for _, info := range buffInfo.SelectedBuff {
		if info.Level == level && info.Type == int32(MazeTempBuff.Type_USE_ITEM) {
			count++
		}
	}

	if count >= config.Energy_item_select {
		logger.CtxError(ctx, "checkUseItemLevelSelectBuff buff count select max", zap.Int32("count", count),
			zap.Int32("maxCount", config.Energy_item_select))
		return fmt.Errorf("道具选择buff次数已用完")
	}
	return nil
}

func (s *service) packOptionalInfo(ctx context.Context, info *tempbuffmodel.TempBuffInfoModel) *OptionalBuffInfo {
	optionalInfo := &OptionalBuffInfo{
		SelectBuffList: s.packSelectBuffList(ctx, info.BuffSequence.OptionalBuffList),
		SelectBuffTime: mazeconfigv8config.GetBuffSelectTime(ctx),
	}

	if len(optionalInfo.SelectBuffList) == 0 {
		return nil
	}

	// 是否可以刷新
	config := mazeenergyresetcostv8config.GetEnergyResetCostConfig(info.BuffSequence.RefreshCount + 1)
	if config == nil {
		return optionalInfo
	}

	optionalInfo.IsRefresh = 1
	for itemId, count := range config.Cost {
		optionalInfo.Cost = append(optionalInfo.Cost, &Item{
			ItemId: itemId,
			Count:  count,
		})
	}

	return optionalInfo
}

func (s *service) packSelectBuffList(ctx context.Context, buffList []int32) []*BuffInfo {
	if len(buffList) == 0 {
		return nil
	}

	showBuffList := make([]*BuffInfo, 0, len(buffList))
	for _, buffId := range buffList {
		if buffId == 0 {
			continue
		}

		if showBuff := s.packShowBuff(ctx, buffId, 1); showBuff != nil {
			showBuffList = append(showBuffList, showBuff)
		}
	}

	return showBuffList
}

// 生成可选的buff列表
func (s *service) createOptionalBuffList(ctx context.Context, buffInfo *tempbuffmodel.TempBuffInfoModel,
	level, areaId, attrMask int32, barrierConfig *GMazeBarriesV8Cfg.MazeBarriesV8ConfigRow) ([]int32, error) {
	logger := fklog.ContextAppLogger(ctx)
	ruleId, ok := barrierConfig.Energy_affix_rand_rule[areaId]
	if !ok || ruleId <= 0 {
		logger.CtxError(ctx, "genOptionalBuffList barrierConfig.Energy_affix_rand_rule not found",
			zap.Int32("level", level), zap.Int32("areaId", areaId))
		return nil, fmt.Errorf("找不到当前区域能力随机规则")
	}
	configId := mazeenergyaffixrandrulev8config.GetKey(ruleId, level)

	randConfig := GMazeEnergyAffixRandRuleV8Cfg.GetWithCtx(ctx, configId)
	if randConfig == nil {
		return nil, nil
	}
	// 统计词条已经选择的次数
	selectedBuffMap := make(map[int32]int32)
	for _, info := range buffInfo.SelectedBuff {
		selectedBuffMap[info.BuffId] += 1
	}
	// 统计词条组以及选择的次数
	selectedBuffGroupMap := make(map[int32]int32)
	for _, i := range buffInfo.SelectedBuff {
		buffConfig := GMazeEnergyAffixV8Cfg.GetWithCtx(ctx, i.BuffId)
		if buffConfig == nil {
			logger.CtxError(ctx, "createOptionalBuffList GetAffixConfig is nil", zap.Int32("buffId", i.BuffId))
			continue
		}
		selectedBuffGroupMap[buffConfig.Affix_group_id] += 1
	}
	// 统计随机库的选择次数
	selectAffixLibraryCountMap := map[int32]int32{}
	for _, i := range buffInfo.SelectedBuff {
		libraryList := mazeenergyaffixlibrarycfgex.GetLibrary(i.BuffId)
		for _, libraryId := range libraryList {
			selectAffixLibraryCountMap[libraryId] += 1
		}
	}

	optionalBuffMap := make(map[int32]struct{})
	libraryMap := make(map[int32]int32) // 本次已选择的库
	num := mazeconfigv8config.GetBuffSelectCount()
	for i := int64(1); i <= num; i++ {
		var libraryId int32
		var posLib map[int32]int32
		switch i {
		case 1:
			posLib = randConfig.Pos_1_lib
		case 2:
			posLib = randConfig.Pos_2_lib
		case 3:
			posLib = randConfig.Pos_3_lib
		default:
			logger.CtxError(ctx, "createOptionalBuffList unknown id", zap.Int64("pos", i))
			return nil, nil
		}
		maxRandLibCount := len(posLib)
		for j := 1; j <= maxRandLibCount; j++ {
			// 随机一个库
			libraryId, _ = s.randLibraryId(ctx, posLib, selectAffixLibraryCountMap, libraryMap)
			if libraryId == 0 {
				logger.CtxError(ctx, "randLibraryId libraryId id=0", zap.Any("posLib", posLib), zap.Any("configId", configId))
				break
			}
			libraryConfig := GMazeEnergyAffixLibraryV8Cfg.GetWithCtx(ctx, libraryId)
			if libraryConfig == nil || len(libraryConfig.Affix_id_list) == 0 {
				return nil, errors.New("affixList is nil")
			}

			// 过滤出可选择的词条
			optionalList, totalWeight := s.filterBuffList(ctx, optionalBuffMap, libraryConfig.Affix_id_list, libraryConfig.Certainly_affix_id_list,
				selectedBuffMap, selectedBuffGroupMap, attrMask)
			// 随机选择个词条
			buffId, weight := s.randomId(optionalList, totalWeight)

			logger.CtxDebug(ctx, "createOptionalBuffList random", zap.Int64("pos", i), zap.Int32("libraryId", libraryId),
				zap.Int32s("affixList", libraryConfig.Affix_id_list), zap.Any("optionalList", optionalList),
				zap.Int32("weight", weight), zap.Int32("id", buffId), zap.Int("randLibCount", j))

			if buffId != 0 {
				optionalBuffMap[buffId] = struct{}{}
				libraryMap[libraryConfig.Order] += 1
				break
			} else {
				// 这次没随机到就把这个库删掉重新随机
				tempPosLib := make(map[int32]int32)
				for k, v := range posLib {
					if k == libraryId {
						continue
					}
					tempPosLib[k] = v
				}
				posLib = tempPosLib
			}
		}
	}

	var optionalList []int32
	for optionId := range optionalBuffMap {
		optionalList = append(optionalList, optionId)
	}

	logger.CtxInfo(ctx, "createOptionalBuffList end", zap.Int32s("optionalList", optionalList))
	return optionalList, nil
}

func (s *service) calcLibraryAddWeight(weight int32, weightAdjust1Add []int32, buffCount int32) int32 {
	if len(weightAdjust1Add) == 0 || buffCount <= 0 {
		return weight
	}
	addIndex := buffCount - 1
	addCount := int32(len(weightAdjust1Add))
	if buffCount >= addCount {
		addIndex = addCount - 1
	}
	weightCoefficient := weightAdjust1Add[addIndex]
	weight = weight * (weightCoefficient + 10000) / 10000
	return weight
}

func (s *service) calcLibrarySubWeight(weight int32, weightAdjust2Sub []int32, thisSelectCount int32) int32 {
	if len(weightAdjust2Sub) == 0 || thisSelectCount <= 0 {
		return weight
	}
	subCount := int32(len(weightAdjust2Sub))
	subIndex := thisSelectCount - 1
	if thisSelectCount >= subCount {
		subIndex = subCount - 1
	}
	weightCoefficient := weightAdjust2Sub[subIndex]
	weight = weight * (10000 - weightCoefficient) / 10000
	return weight
}

func (s *service) randLibraryId(ctx context.Context, libraryMap map[int32]int32, selectAffixLibraryCountMap map[int32]int32, thisSelectLibraryMap map[int32]int32) (int32, int32) {
	var (
		weightList  []*WeightInfo
		totalWeight int32
	)
	for id, weight := range libraryMap {
		if weight == 0 {
			continue
		}
		libraryConfig := GMazeEnergyAffixLibraryV8Cfg.GetWithCtx(ctx, id)
		if libraryConfig == nil {
			continue
		}
		if selectAffixLibraryCountMap[id] > 0 {
			// 计算增加的权重系数
			weight = s.calcLibraryAddWeight(weight, libraryConfig.Weight_adjust1_add, selectAffixLibraryCountMap[id])
		}
		if thisSelectLibraryMap[id] > 0 {
			// 计算减少的权重系数
			weight = s.calcLibrarySubWeight(weight, libraryConfig.Weight_adjust2_sub, thisSelectLibraryMap[id])
			if weight <= 0 {
				continue
			}
		}

		weightList = append(weightList, &WeightInfo{
			Id:     id,
			Weight: weight,
		})
		totalWeight += weight
	}

	return s.randomId(weightList, totalWeight)
}

type WeightInfo struct {
	Id     int32
	Weight int32
}

// 过滤本次可选的词条
func (s *service) filterBuffList(ctx context.Context, optionalMap map[int32]struct{}, buffList, certainlyList []int32,
	selectedBuffMap, selectedBuffGroupMap map[int32]int32, attrMask int32) ([]*WeightInfo, int32) {

	var (
		optionalList []*WeightInfo
		totalWeight  int32
	)

	// 先添加必选buff
	for _, buffId := range certainlyList {
		if buffId == 0 {
			continue
		}
		if _, ok := optionalMap[buffId]; ok {
			continue
		}
		buffWeight := s.GetOptionBuffWeightInfo(ctx, buffId, selectedBuffMap, selectedBuffGroupMap, optionalMap, attrMask)
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

		buffWeight := s.GetOptionBuffWeightInfo(ctx, buffId, selectedBuffMap, selectedBuffGroupMap, optionalMap, attrMask)
		if buffWeight == nil {
			continue
		}

		optionalList = append(optionalList, buffWeight)
		totalWeight += buffWeight.Weight
	}

	return optionalList, totalWeight
}

// 检查buff是否满足可选条件， 获取可选buff的权重信息
func (s *service) GetOptionBuffWeightInfo(ctx context.Context, buffId int32, selectedBuffMap, selectedBuffGroupMap map[int32]int32,
	optionalMap map[int32]struct{}, attrMask int32) *WeightInfo {
	logger := fklog.ContextAppLogger(ctx)
	buffConfig := GMazeEnergyAffixV8Cfg.GetWithCtx(ctx, buffId)
	if buffConfig == nil {
		logger.CtxWarn(ctx, "getOptionBuffWeightInfo buffConfig is nil", zap.Int32("buffId", buffId))
		return nil
	}
	if !TestBuffAttrMask(buffConfig.Affix_group_id, attrMask) {
		return nil
	}
	if buffConfig.Weight == 0 {
		logger.CtxWarn(ctx, "getOptionBuffWeightInfo buff weight is 0", zap.Int32("buffId", buffId))
		return nil
	}

	if buffConfig.Affix_group_id != int32(mazeconfigv8.GetSpecialBuffGroupId(ctx)) {
		_, exist := selectedBuffGroupMap[buffConfig.Affix_group_id]
		if !exist {
			maxCount := int(mazeconfigv8.GetMaxBuffGroupCount(ctx))
			if len(selectedBuffGroupMap) > maxCount {
				return nil
			} else if len(selectedBuffGroupMap) == maxCount {
				_, exist = selectedBuffGroupMap[int32(mazeconfigv8.GetSpecialBuffGroupId(ctx))]
				if !exist {
					return nil
				}
			}
		}
	}

	// 检查选择数量
	optionalCount := buffConfig.Use_num_max - selectedBuffMap[buffId]
	if optionalCount <= 0 {
		return nil
	}

	// 检查前置条件
	for _, frontId := range buffConfig.Font_affix_condition {
		if frontId == 0 {
			continue
		}
		isOk := s.checkFrontCondition(ctx, logger, frontId, selectedBuffMap, selectedBuffGroupMap, optionalMap)
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
func (s *service) checkFrontCondition(ctx context.Context, logger fklog.FKLogI, frontId int32, selectedBuffMap, selectedBuffGroupMap map[int32]int32,
	optionalMap map[int32]struct{}) bool {
	frontConfig := GMazeEnergyAffixFrontV8Cfg.GetWithCtx(ctx, frontId)
	if frontConfig == nil {
		return false
	}

	if frontConfig.Exclusive_affix__id != 0 {
		// 检查互斥词条
		if _, ok := optionalMap[frontConfig.Exclusive_affix__id]; ok {
			return false
		}
	}

	// 检查前置词条是否满足
	var count int32
	for _, affixId := range frontConfig.Affix_id_set {
		if _, ok := selectedBuffMap[affixId]; !ok {
			continue
		}

		count++

		if count >= frontConfig.Must_num {
			break
		}
	}

	if count < frontConfig.Must_num {
		//logger.CtxDebug(ctx,"checkFrontCondition affix id set not enough", zap.Int32("frontId", frontId), zap.Int32("count", count))
		return false
	}

	// 检查前置词条组数量是否满足
	for k, v := range frontConfig.Affix_group_num {
		if k == 0 {
			continue
		}
		groupCount, _ := selectedBuffGroupMap[k]
		if groupCount < v {
			logger.CtxDebug(ctx, "checkFrontCondition groupCount not enough", zap.Int32("frontId", frontId), zap.Int32("groupCount", groupCount))
			return false
		}
	}

	return true
}

func (s *service) randomId(optionalList []*WeightInfo, totalWeight int32) (int32, int32) {
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

const (
	IceMask int32 = 1 << iota
	FireMask
	FlashMask
	PoisonMask
)

// 6001 6002 电
// 6003 6004 冰
// 6005 6006 火
// 6007 6008 毒
// 测试用，只选需要的buff
func TestBuffAttrMask(groupId, attrMask int32) bool {
	if attrMask == 0 {
		return true
	}
	if groupId < 6001 || groupId > 6008 {
		return true
	}
	// 测试用属性掩码 0-全部 1-冰 2-火 4-电 8-毒
	if attrMask&IceMask > 0 {
		if groupId == 6003 || groupId == 6004 {
			return true
		}
	}
	if attrMask&FireMask > 0 {
		if groupId == 6005 || groupId == 6006 {
			return true
		}
	}
	if attrMask&FlashMask > 0 {
		if groupId == 6001 || groupId == 6002 {
			return true
		}
	}
	if attrMask&PoisonMask > 0 {
		if groupId == 6007 || groupId == 6008 {
			return true
		}
	}
	return false
}
