package tempbuffservice

import (
	"context"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/model/passareamodel"
)

func (s *service) DelPassArea(ctx context.Context, userID uint64, barrierId int32) error {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "DelPassArea", zap.Int32("barrierId", barrierId))
	var model = &passareamodel.PassAreaModel{}
	return model.Del(ctx, userID, barrierId)
}
