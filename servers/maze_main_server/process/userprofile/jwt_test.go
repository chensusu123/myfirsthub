// @Author pangchenyang 2025/6/17 10:59:00
// @Desc: 
package userprofile

import (
	"testing"
	"encoding/base64"
	"math/rand"
	"time"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/pb/common/UserProfile"
	"github.com/stretchr/testify/assert"
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
			res := &UserProfile.QueryAvatarTokenRS{}
			shardingID := int64(tt.alterProfile.GetUserId())
			err := testProfile.OnQueryAvatarToken_10511_10512(logger, shardingID, req, res, "")
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("TestOnQueryAvatarToken failed, error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Logf("TestOnQueryAvatarToken success, res = %v", res)
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
	// t.Logf("GenerateAvatarToken token: %s", token)
	// _, err = VerifyWithCustomClaims("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxNjcyNjh9.7XleGVrAtihKaYlLBXIVMHGjXkZ6P2Bej8CIHyjm7WQ")
	// t.Logf("VerifyWithCustomClaims err: %v", err)
	// validate := ValidateAvatarToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMDAwMDEsInVzZXJuYW1lIjoiIiwiZXhwIjoxNzUwMTY3NTAyLCJpc3MiOiJ1c2VyX2F2YXRhciJ9.0vdyWBE8VUNTN-9ypD5nPWvd2DxanOIJnkxcw2Ww8No")
	// t.Logf("ValidateAvatarToken validate: %v", validate)
	_, err = VerifyWithCustomClaims(token)
	t.Logf("VerifyWithCustomClaims err: %v", err)
	time.Sleep(time.Second * 2)
	_, err = VerifyWithCustomClaims(token)
	t.Logf("VerifyWithCustomClaims err: %v", err)
	// validate := ValidateAvatarToken(token)
	// t.Logf("ValidateAvatarToken validate: %v", validate)
	// time.Sleep(4 * time.Second)
	// _, err = VerifyWithCustomClaims(token)
	// t.Logf("VerifyWithCustomClaims err: %v", err)
	// validate = ValidateAvatarToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTAxNzIzNjF9.7l0UVm82SuENlvLzzZTMkhoBnoJnow4jM5jmbnxgIOo")
	// t.Logf("ValidateAvatarToken validate: %v", validate)
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
