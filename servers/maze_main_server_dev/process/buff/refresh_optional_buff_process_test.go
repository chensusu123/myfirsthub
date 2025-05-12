package buff

import (
	"context"
	"testing"

	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/plate/protodef/MazeTempBuff"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/24 17:18
 * @Description:
 */

func TestRefreshOptionalMazeTempBuffListRQ(t *testing.T) {
	type args struct {
		logger     fknet.TCPContext
		shardingID uint64
		request    *MazeTempBuff.RefreshOptionalMazeTempBuffListRQ
		response   *MazeTempBuff.RefreshOptionalMazeTempBuffListRS
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "刷新buff",
			args: args{
				logger: fknet.TCPContext{
					Context: context.Background(),
					FKLogI:  gTestLogger,
				},
				shardingID: 9003200130019765,
				request: &MazeTempBuff.RefreshOptionalMazeTempBuffListRQ{
					StageId: proto.Int32(1),
					Level:   proto.Int32(2),
					Cost: []*MazeCommon.MazeItem{
						{
							ItemId: proto.Int32(46900001),
							Count:  proto.Int64(50),
						},
					},
				},
				response: &MazeTempBuff.RefreshOptionalMazeTempBuffListRS{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RefreshOptionalMazeTempBuffListRQ(tt.args.logger, tt.args.shardingID, tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("RefreshOptionalMazeTempBuffListRQ() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
