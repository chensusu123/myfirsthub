package familyservice

import (
	"context"
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 订阅家族聊天室
func (r *service) SubscribeFamilyChat(ctx context.Context, familyID int32, userID uint64, group_ids []int64) error {
	//TODO: 订阅家族聊天室
	logger := fklog.ContextAppLogger(ctx)
	familyModel := &familymodel.FamilySubscribeModel{
		FamilyId: familyID,
		UserId:   userID,
		GroupIds: group_ids,
	}
	//数组不为空，保存订阅 数组为空，则删除订阅
	if len(group_ids) > 0 {
		err := familyModel.Save(ctx)
		if err != nil {
			logger.CtxError(ctx, "SubscribeFamilyChat Save error", zap.Any("err", err))
			return err
		}
	} else {
		err := familyModel.Delete(ctx)
		if err != nil {
			logger.CtxError(ctx, "SubscribeFamilyChat Delete error", zap.Any("err", err))
			return err
		}
	}
	return nil
}
