package buff

import (
	"context"
	"testing"
	"time"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/22 16:29
 * @Description:
 */

func TestMazeBarrierNotifyProcess(t *testing.T) {
	type args struct {
		c      context.Context
		logger fklog.FKLogI
		index  int
		key    []byte
		data   *MazeBarrierUserGameRecord
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "关卡通关",
			args: args{
				c:      context.Background(),
				logger: gTestLogger,
				index:  0,
				data: &MazeBarrierUserGameRecord{
					UserId:     9003200130019765,
					Barrier:    1,
					GameRet:    1,
					GroupID:    12,
					CreateTime: time.Now().Unix(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// msg, err := json.Marshal(tt.args.data)
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			// }

			MazeBarrierNotifyProcess(tt.args.logger, tt.args.data)
		})
	}
}
