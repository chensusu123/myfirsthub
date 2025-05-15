package randfuncs

import (
	"math/rand"
	"sort"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

type ItemWeight struct {
	Id     int32
	Weight int32
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 按权重随机
func RandByWeight(log fklog.FKLogI, input map[int32]int32) (result int32) {
	return RandByWeightV3(log, input, true)
}

// 按权重随机v2
func RandByWeightV2(log fklog.FKLogI, input map[int32]int32, isLog bool) (result int32) {
	return RandByWeightV3(log, input, isLog)
}

// 按权重随机v4 旧
// func RandByWeightV4(log fklog.FKLogI, input map[int32]int32, isLog bool) (result int32) {
// 	weightTotal := int32(0)
// 	weightSlice := make([]*ItemWeight, 0)

// 	for id, weight := range input {
// 		if weight == 0 { //只过滤权重
// 			continue
// 		}

// 		weightTotal += weight
// 		weightSlice = append(weightSlice, &ItemWeight{Id: id, Weight: weight})
// 	}

// 	randWeight := CustomRand(1, int(weightTotal+1))
// 	tmpWeightPos := int32(0)
// 	sort.Slice(weightSlice, func(i, j int) bool {
// 		return weightSlice[i].Weight < weightSlice[j].Weight
// 	})

// 	defer func() {
// 		if isLog {
// 			log.DebugWF("RandByWeightV4 dump",
// 				zap.Int32("result", result),
// 				zap.Any("input", weightSlice),
// 				zap.Int32("weightTotal", weightTotal),
// 				zap.Int("randWeight", randWeight))
// 		}
// 	}()

// 	for _, v := range weightSlice {
// 		tmpWeightPos += v.Weight
// 		if randWeight <= int(tmpWeightPos) {
// 			result = v.Id
// 			return
// 		}
// 	}
// 	return
// }

// 按权重随机v3，输入序列总权重不要为0
func RandByWeightV3(log fklog.FKLogI, input map[int32]int32, isLog bool) (result int32) {
	weightTotal := int32(0)
	weightSlice := make([]*ItemWeight, 0)

	for id, weight := range input {
		if weight == 0 { //只过滤权重
			continue
		}

		weightTotal += weight
		weightSlice = append(weightSlice, &ItemWeight{Id: id, Weight: weight})
	}
	var randWeight int
	defer func() {
		if isLog {
			log.InfoWF("RandByWeightV3 dump",
				zap.Int32("result", result),
				zap.Any("input", weightSlice),
				zap.Int32("weightTotal", weightTotal),
				zap.Int("randWeight", randWeight))
		}
	}()

	if weightTotal == 0 {
		return
	}

	randWeight = CustomRand(1, int(weightTotal+1))
	sort.Slice(weightSlice, func(i, j int) bool {
		return weightSlice[i].Id < weightSlice[j].Id
	})

	tmpWeightPos := int32(0)
	for _, v := range weightSlice {
		tmpWeightPos += v.Weight
		if randWeight <= int(tmpWeightPos) {
			result = v.Id
			return
		}
	}
	return
}

func CustomRand(start, end int) int {
	return fkutil.RandInt32(start, end)
}
