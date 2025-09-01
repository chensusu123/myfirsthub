package equiptoitem

import (
	"context"
	"errors"
	"fmt"
	"maze_game_server/excel/mazeequiptyperesv8"

	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/config/GMazeEquipTypeResV8Cfg"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/server/MazeEquipSvr"

	"google.golang.org/protobuf/proto"
)

func PackMazeEquipInfoSvrToItem(ctx context.Context, equipId int32) (item *MazeCommon.MazeItem, err error) {
	equipTypeResCfg := mazeequiptyperesv8.GetEquipTypeResCfg(equipId, 0, 0)
	if equipTypeResCfg == nil {
		err = fmt.Errorf("mazeequiptyperesv8 get cfg fail, equipId:%d", equipId)
		return
	}
	equip := &MazeEquipSvr.MazeEquipInfoSvr{
		EquipId:    proto.Int32(equipId),
		EquipResId: proto.Int32(equipTypeResCfg.Order),
		EquipName:  proto.String(equipTypeResCfg.Name),
	}
	itemEquip, err := PackEquipToItem(ctx, equip)
	return itemEquip, err
}

func PackEquipToItem(ctx context.Context, equip *MazeEquipSvr.MazeEquipInfoSvr) (item *MazeCommon.MazeItem, err error) {
	if equip == nil {
		err = errors.New("equip nil")
		return
	}

	equipInfoCfg := GMazeEquipInfoV8Cfg.GetWithCtx(ctx, equip.GetEquipId())
	if equipInfoCfg == nil {
		err = fmt.Errorf("get equip info cfg fail, equipId:%d", equip.GetEquipId())
		return
	}
	equipResCfg := GMazeEquipTypeResV8Cfg.GetWithCtx(ctx, equip.GetEquipResId())
	if equipResCfg == nil {
		err = fmt.Errorf("get equip res cfg fail, resId:%d", equip.GetEquipResId())
		return
	}
	item = &MazeCommon.MazeItem{
		ItemId:    proto.Int32(equip.GetEquipId()),
		Count:     proto.Int64(1),
		ItemName:  proto.String(equip.GetEquipName()),
		Guid:      proto.Int64(equip.GetEquipGuid()),
		IconName:  proto.String(equipResCfg.Icon),
		AtlasName: proto.String(equipResCfg.IconAtlas),
		Quality:   proto.Int32(equipInfoCfg.Quality),
	}
	return
}
