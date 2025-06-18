// @Author pangchenyang 2025/6/17 22:06:00
// @Desc: 
package userprofile

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"maze_game_server/pb/common/UserProfile"
	"maze_game_server/common/errors"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"google.golang.org/protobuf/proto"
	"fmt"
)

// OnQueryAvatarToken 查询用户操作头像jwt token
func (p *Profile) OnQueryAvatarToken_10511_10512(ctx fklog.FKLogI, shardingID int64, rqMsg proto.Message, rsMsg proto.Message, opData string) (err error) {
	defer fkprometheus.DebugPMT("OnQueryAvatarToken")()
	req := rqMsg.(*UserProfile.QueryAvatarTokenRQ)
	res := rsMsg.(*UserProfile.QueryAvatarTokenRS)
	res.ErrInfo = errors.NO_ERROR
	ctx.WarnWF("OnQueryAvatarToken with", zap.Any("rq", req))

	defer func() {
		ctx.WarnWF("OnQueryAvatarToken end ", zap.Any("req", req), zap.Any("res", res),
			zap.String("errMsg", string(res.GetErrInfo().GetErrMsg())))
	}()

	// 检查rq
	if shardingID <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	// 生成token
	token, err := GenerateAvatarToken()
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(fmt.Sprintf("生成token失败: %v", err))
		return
	}
	res.Token = proto.String(token)
	return
}
