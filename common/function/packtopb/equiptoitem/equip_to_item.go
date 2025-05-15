package equiptoitem

import (
	"errors"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeEquipInfoV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeEquipTypeResV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipSvr"
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
