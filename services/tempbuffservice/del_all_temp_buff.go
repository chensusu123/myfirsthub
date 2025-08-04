package tempbuffservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/model/tempbuffmodel"
)

func (s *service) DelTempBuff(logger fklog.FKLogI, userID uint64, stageId int32) error {
	var model = &tempbuffmodel.TempBuffInfoModel{}
	return model.Del(logger, userID, stageId)
}
