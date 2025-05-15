package buff

import (
	"context"
	"testing"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeTempBuff"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/22 16:24
 * @Description:
 */

func TestSelectMazeTempBuffRQ(t *testing.T) {
	type args struct {
		logger     fknet.TCPContext
		shardingID uint64
		request    *MazeTempBuff.SelectMazeTempBuffRQ
		response   *MazeTempBuff.SelectMazeTempBuffRS
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "选择buff",
			args: args{
				logger: fknet.TCPContext{
					Context: context.Background(),
					FKLogI:  gTestLogger,
				},
				shardingID: 9003200130019765,
				request: &MazeTempBuff.SelectMazeTempBuffRQ{
					Header:  nil,
					StageId: proto.Int32(1),
					Level:   proto.Int32(2),
					BuffId:  proto.Int32(20007),
				},
				response: &MazeTempBuff.SelectMazeTempBuffRS{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SelectMazeTempBuffRQ(tt.args.logger, tt.args.shardingID, tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("SelectMazeTempBuffRQ() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
