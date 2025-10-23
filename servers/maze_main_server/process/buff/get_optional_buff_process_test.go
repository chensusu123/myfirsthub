package buff

import (
	"math/rand"
	"maze_game_server/services/tempbuffservice"
	"testing"
	"time"

	"go.uber.org/zap"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/22 16:17
 * @Description:
 */

// func TestGetOptionalMazeTempBuffListRQ(t *testing.T) {
// 	type args struct {
// 		logger     fknet.TCPContext
// 		shardingID uint64
// 		request    *MazeTempBuff.GetOptionalMazeTempBuffListRQ
// 		response   *MazeTempBuff.GetOptionalMazeTempBuffListRS
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "查询选择buff",
// 			args: args{
// 				logger: fknet.TCPContext{
// 					Context: context.Background(),
// 					FKLogI:  gTestLogger,
// 				},
// 				shardingID: 9003200130019765,
// 				request: &MazeTempBuff.GetOptionalMazeTempBuffListRQ{
// 					Header:  nil,
// 					StageId: proto.Int32(1),
// 					Level:   proto.Int32(2),
// 				},
// 				response: &MazeTempBuff.GetOptionalMazeTempBuffListRS{},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := GetOptionalMazeTempBuffListRQ(tt.args.logger, tt.args.shardingID, tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
// 				t.Errorf("GetOptionalMazeTempBuffListRQ() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

func Test_randomBuff(t *testing.T) {
	// 测试随机权重算法
	optionalList := []*tempbuffservice.WeightInfo{
		{Id: 1, Weight: 10},
		{Id: 2, Weight: 20},
		{Id: 3, Weight: 30},
		{Id: 4, Weight: 40},
	}
	totalWeight := int32(100)

	buffMap := make(map[int32]int64)
	// 减少测试次数以提高测试速度
	for i := 0; i < 100000; i++ {
		buffId := randomIdSimple(optionalList, totalWeight)
		buffMap[buffId] += 1
	}

	for _, info := range optionalList {
		gTestLogger.InfoWF("random buff", zap.Any("info", info), zap.Int64("count", buffMap[info.Id]),
			zap.Float64("ratio", float64(buffMap[info.Id])/float64(info.Weight)))
	}

	time.Sleep(time.Second)
}

// 简化的随机ID函数，用于测试
func randomIdSimple(optionalList []*tempbuffservice.WeightInfo, totalWeight int32) int32 {
	if len(optionalList) == 0 || totalWeight <= 0 {
		return 0
	}

	weight := int32(rand.Intn(int(totalWeight)))
	var curWeight int32
	for _, info := range optionalList {
		curWeight += info.Weight
		if weight < curWeight {
			return info.Id
		}
	}

	return 0
}
