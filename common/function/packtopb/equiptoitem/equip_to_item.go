package equiptoitem

import (
	"errors"
	"fmt"

	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/config/GMazeEquipTypeResV8Cfg"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/server/MazeEquipSvr"

	"google.golang.org/protobuf/proto"
)

func PackEquipToItem(equip *MazeEquipSvr.MazeEquipInfoSvr) (item *MazeCommon.MazeItem, err error) {
	if equip == nil {
		err = errors.New("equip nil")
		return
	}

	equipInfoCfg := GMazeEquipInfoV8Cfg.Get(equip.GetEquipId())
	if equipInfoCfg == nil {
		err = fmt.Errorf("get equip info cfg fail, equipId:%d", equip.GetEquipId())
		return
	}
	equipResCfg := GMazeEquipTypeResV8Cfg.Get(equip.GetEquipResId())
	if equipResCfg == nil {
		err = fmt.Errorf("get equip res cfg fail, resId:%d", equip.GetEquipResId())
		return
	}
	item = &MazeCommon.MazeItem{
		ItemId:    proto.Int32(equip.GetEquipId()),
		ItemName:  proto.String(equip.GetEquipName()),
		Guid:      proto.Int64(equip.GetEquipGuid()),
		IconName:  proto.String(equipResCfg.Icon),
		AtlasName: proto.String(equipResCfg.IconAtlas),
		Quality:   proto.Int32(equipInfoCfg.Quality),
	}
	return
}
