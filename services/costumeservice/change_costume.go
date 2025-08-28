package costumeservice

import (
	"context"

	"maze_game_server/pb/common/Costume"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 换装
func (s *service) ChangeCostume(ctx context.Context, userId uint64) {
	logger := fklog.ContextAppLogger(ctx)
	// todo: 获取装备的逻辑比较复杂，就在这重新获取了
	res, err := s.GetUserCostume(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "ChangeCostume GetUserCostume err", zap.Error(err))
		return
	}
	push := &Costume.ChangeCostumeID{}
	push.Costume = make([]*Costume.CostumeInfo, 0, len(res))
	for k, v := range res {
		push.Costume = append(push.Costume, &Costume.CostumeInfo{
			Pos:     proto.Int32(k),
			ModelId: proto.Int32(v),
		})
	}
	online.PushWithContext(ctx, userId, 10615, push)

	return
}
