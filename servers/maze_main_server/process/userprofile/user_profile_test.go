// @Author pangchenyang 2025/6/11 21:55:00
// @Desc: 
package userprofile

import (
	"testing"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"time"
	"maze_game_server/pb/common/UserProfile"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestOnQueryUserProfile(t *testing.T) {
	logger = fklog.AppLogger().Clone("query_user_profile_t")
	logger.SetLogId(time.Now().UnixNano())
	defer func() {
		time.Sleep(time.Second * 2)
	}()
	cases := []struct {
		name       string
		shardingID int64
		userID     []uint64
		wantErr    error
	}{
		{
			name:       "case1",
			shardingID: 1,
			userID:     []uint64{1},
			wantErr:    nil,
		},
		{
			name:       "case1",
			shardingID: 2,
			userID:     []uint64{1, 2},
			wantErr:    nil,
		},
		{
			name:       "case2",
			shardingID: 1,
			userID:     []uint64{1, 2},
			wantErr:    nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := &UserProfile.QueryUserProfileRQ{
				UserId: tt.userID,
			}
			res := &UserProfile.QueryUserProfileRS{}
			err := testProfile.OnQueryUserProfile_10481_10482(logger, tt.shardingID, req, res, "")
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("TestOnQueryUserProfile failed, error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("TestOnQueryUserProfile success, res = %v", res)
			}

		})
	}
}

func TestOnAlterUserProfile(t *testing.T) {
	logger = fklog.AppLogger().Clone("alter_user_profile_t")
	logger.SetLogId(time.Now().UnixNano())
	defer func() {
		time.Sleep(time.Second * 2)
	}()
	cases := []struct {
		name         string
		alterProfile UserProfile.UserProfile
		wantErr      error
	}{
		{
			name: "case1",
			alterProfile: UserProfile.UserProfile{
				UserId:   proto.Uint64(1),
				NickName: proto.String("alterTest"),
			},
			wantErr: nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := &UserProfile.AlterUserProfileRQ{
				AlterProfile: &tt.alterProfile,
			}
			res := &UserProfile.AlterUserProfileRS{}
			shardingID := int64(tt.alterProfile.GetUserId())
			err := testProfile.OnAlterUserProfile_10483_10484(logger, shardingID, req, res, "")
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("TestOnAlterUserProfile failed, error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("TestOnAlterUserProfile success, res = %v", res)
			}
		})
	}
}
