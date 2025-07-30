package buff

import (
	"maze_game_server/excel/mazeenergyaffixlibraryv8config"
	"maze_game_server/pb/server/MazeTempBuffSvr"
	"testing"
	"time"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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
	// 随机库id
	buffInfo := &MazeTempBuffSvr.TempBuffInfo{
		BuffSequence: new(MazeTempBuffSvr.BuffSequence),
		SelectedBuff: []*MazeTempBuffSvr.SelectedBuffInfo{
			{BuffId: proto.Int32(401010101)},
			{BuffId: proto.Int32(401010202)},
			{BuffId: proto.Int32(401010303)},
			{BuffId: proto.Int32(401010404)},
		},
	}
	optionalMap := make(map[int32]struct{})
	affixList, certainly_list := mazeenergyaffixlibraryv8config.GetEnergyLibraryAffixList(1)
	optionalList, totalWeight := filterBuffList(buffInfo, optionalMap, affixList, certainly_list)
	buffMap := make(map[int32]int64)
	for i := 0; i < 1000000000; i++ {
		buffId, _ := randomId(optionalList, totalWeight)
		buffMap[buffId] += 1
	}

	for _, info := range optionalList {
		gTestLogger.InfoWF("random buff", zap.Any("info", info), zap.Int64("count", buffMap[info.Id]),
			zap.Int64("value", buffMap[info.Id]/int64(info.Weight)))
	}

	time.Sleep(time.Second)
}

func TestGetOptionalBuffList(t *testing.T) {
	var userId uint64 = 40000007
	var stageId int32 = 12
	var level int32 = 5
	var buffType int32 = 1
	var areaId int32 = 120001
	buffInfo := &MazeTempBuffSvr.TempBuffInfo{
		BuffSequence: &MazeTempBuffSvr.BuffSequence{
			Index: proto.Int32(5),
		},
		TotalBuff: []*MazeTempBuffSvr.TotalBuffInfo{
			{
				BuffId:    proto.Int32(9201),
				BuffValue: proto.Int64(9000),
			},
			{
				BuffId:    proto.Int32(3060001),
				BuffValue: proto.Int64(22000),
			},
			{
				BuffId:    proto.Int32(3041501),
				BuffValue: proto.Int64(3),
			},
			{
				BuffId:    proto.Int32(3060000),
				BuffValue: proto.Int64(1),
			},
			{
				BuffId:    proto.Int32(7000302),
				BuffValue: proto.Int64(5000),
			},
		},
		SelectedBuff: []*MazeTempBuffSvr.SelectedBuffInfo{
			{BuffId: proto.Int32(700100001), Level: proto.Int32(2), Type: proto.Int32(1)},
			{BuffId: proto.Int32(700100101), Level: proto.Int32(3), Type: proto.Int32(1)},
			{BuffId: proto.Int32(700100102), Level: proto.Int32(4), Type: proto.Int32(1)},
		},
	}
	errInfo := getOptionalBuffList(gTestLogger, userId, stageId, level, buffType, areaId, buffInfo)
	if errInfo != nil {
		gTestLogger.ErrorWF("GetOptionalMazeTempBuffListRQ getOptionalBuffList", zap.Int32("stageId", stageId),
			zap.Any("info", buffInfo), zap.Any("errInfo", errInfo))
	}
}
