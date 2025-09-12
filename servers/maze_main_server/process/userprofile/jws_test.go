// @Author pangchenyang 2025/6/17 10:59:00
// @Desc:
package userprofile

import (
	"encoding/base64"
	"math/rand"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/UserProfile"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"google.golang.org/protobuf/proto"
)

func TestOnQueryAvatarToken(t *testing.T) {
	logger = fklog.AppLogger().Clone("query_avatar_token_t")
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
			req := &UserProfile.QueryAvatarTokenRQ{}
			shardingID := int64(tt.alterProfile.GetUserId())
			session := &session.Session{}
			session.Bind(shardingID)
			err := testProfile.OnQueryAvatarToken_10511_10512(session, req)
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("TestOnQueryAvatarToken failed, error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("TestOnQueryAvatarToken success")
			}
		})
	}
}

func TestGenerateAvatarToken(t *testing.T) {
	token, err := GenerateAvatarToken()
	if err != nil {
		t.Errorf("GenerateAvatarToken error: %v", err)
		return
	}
	_, err = VerifyWithCustomClaims(token)
	t.Logf("VerifyWithCustomClaims err: %v", err)
	time.Sleep(time.Second * 2)
	_, err = VerifyWithCustomClaims(token)
	t.Logf("VerifyWithCustomClaims err: %v", err)
}

func TestGenerateHMACSecretKey(t *testing.T) {
	key, err := generateHMACSecretKey(16)
	if err != nil {
		t.Errorf("generateHMACSecretKey error: %v", err)
		return
	}
	t.Logf("generateHMACSecretKey key: %s", key)
}

// 生成指定长度的随机JWT秘钥
func generateHMACSecretKey(length int) (string, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(secret), nil
}
