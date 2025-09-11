package allianceservice

import (
	"context"
	"maze_game_server/model/alliancemodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 订阅
func (s *service) SubscribeAllianceChat(ctx context.Context, userID uint64, allianceID int32, group_ids []int64) error {
	//TODO: 订阅联盟聊天室
	logger := fklog.ContextAppLogger(ctx)
	allianceModel := &alliancemodel.AllianceSubscribeModel{
		AllianceID: allianceID,
		UserID:     userID,
		GroupIDs:   group_ids,
	}
	//数组不为空，保存订阅 数组为空，则删除订阅
	if len(group_ids) > 0 {
		err := allianceModel.Save(ctx)
		if err != nil {
			logger.CtxError(ctx, "SubscribeAllianceChat Save error", zap.Any("err", err))
			return err
		}
	} else {
		err := allianceModel.Delete(ctx)
		if err != nil {
			logger.CtxError(ctx, "SubscribeAllianceChat Delete error", zap.Any("err", err))
			return err
		}
	}

	return nil
}
