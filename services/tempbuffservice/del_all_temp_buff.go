package tempbuffservice

import (
	"context"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/model/tempbuffmodel"
)

func (s *service) DelTempBuff(ctx context.Context, userID uint64, barrierId int32) error {
	logger := fklog.ContextAppLogger(ctx)
	logger.InfoWF("DelTempBuff", zap.Int32("barrierId", barrierId))
	var model = &tempbuffmodel.TempBuffInfoModel{}
	return model.Del(ctx, userID, barrierId)
}
