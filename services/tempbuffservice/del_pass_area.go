package tempbuffservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/model/passareamodel"
)

func (s *service) DelPassArea(logger fklog.FKLogI, userID uint64, stageId int32) error {
	var model = &passareamodel.PassAreaModel{}
	return model.Del(logger, userID, stageId)
}
