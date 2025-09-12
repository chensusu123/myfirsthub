// @Author pangchenyang 2025/6/11 21:55:00
// @Desc:
package userprofile

import (
	"maze_game_server/io/redis/userprofileredis"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/UserProfile"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/serverdepend"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"google.golang.org/protobuf/proto"
)

func init() {
	// 设置start.sh中的环境变量
	os.Setenv("mode", "dev")
	os.Setenv("HOSTNAME", "maze-main-server-c0-g4-0")
	os.Setenv("S_GROUP_ID", "4")
	os.Setenv("MAZE_REDIS_ADDR_4", "10.101.110.231:9005")
	os.Setenv("LOCAL_DEV", "true")

	// 建立客户端
	db := userprofileredis.NewRedisDemo("maze_main_server.redis", "redis")
	serverdepend.RegisterDepend(db)

	go func() {
		fkserver.Run()
	}()
	time.Sleep(time.Second * 2)
}

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
			req := &UserProfile.QueryUserProfileRQ{}
			s := &session.Session{}
			s.Bind(tt.shardingID)
			err := testProfile.OnQueryUserProfile_10481_10482(s, req)
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("TestOnQueryUserProfile failed, error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("TestOnQueryUserProfile success")
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
			shardingID := int64(tt.alterProfile.GetUserId())
			session := &session.Session{}
			session.Bind(shardingID)
			err := testProfile.OnAlterUserProfile_10483_10484(session, req)
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("TestOnAlterUserProfile failed, error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("TestOnAlterUserProfile success")
			}
		})
	}
}
