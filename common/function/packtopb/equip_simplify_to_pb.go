package packtopb

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/pbutil"
	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func EquipSimplifyToCliPB(ctx context.Context, equipInfo *MazeEquipCache.MazeEquipInfoDb) (*MazeGameEquip.MazeEquipInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeGameEquip.MazeEquipInfo{}
	res.EquipGuid = proto.Int64(equipInfo.GetEquipGuid())
	equipCfg := GMazeEquipInfoV8Cfg.GetWithCtx(ctx, equipInfo.GetEquipId())
	if equipCfg == nil {
		logger.CtxError(ctx, "EquipSimplifyToCliPB get doll equip info cfg nil", zap.Int32("equipId", equipInfo.GetEquipId()))
		return nil, errors.New("装备详情配置不存在")
	}
	res.Pos = proto.Int32(equipCfg.Pos)
	res.EquipQuality = proto.Int32(equipCfg.Quality)
	res.EquipLevel = proto.Int32(equipCfg.Level)
	res.EquipId = proto.Int32(equipCfg.Equipment_id)

	if equipInfo.GetSuitId() > 0 {
		res.SuitInfo = &MazeGameEquip.EquipSuitInfo{
			SuitId: proto.Int32(equipInfo.GetSuitId()),
		}
	}
	_, equipName, mazeModel, icon, iconAtlas := pbutil.GetDollEquipNameEx(ctx, equipInfo, equipInfo.GetEquipSubType())
	//	res.EquipResId = proto.Int32(equipResId)
	res.EquipName = proto.String(equipName)
	res.MazeModel = proto.Int32(mazeModel)
	res.Icon = proto.String(icon)
	res.IconAtlas = proto.String(iconAtlas)

	mainAttrs, _, _, _, _ := EquipBaseAttrToCliPB(ctx, equipInfo.BaseAttrs)
	res.MainAttrs = mainAttrs
	return res, nil
}

// func BagEquipChgIDToCliPBEx(ctx context.Context, equip *MazeEquipCache.MazeEquipInfoDb) (*MazeGameEquip.MazeEquipInfo, error) {
//	newEquipInfo, err := pbutil.ConvertIdentifyEquipDb(logger, equip)
//	if err != nil {
//		logger.CtxError(ctx,"BagEquipChgIDToCliPBEx ConvertIdentifyEquipDb error", zap.Any("equip", equip), zap.Error(err))
//		return nil, err
//	}
//	res := &MazeGameEquip.MazeEquipInfo{}
//	res.EquipGuid = proto.Int64(newEquipInfo.GetEquipGuid())
//	equipCfg := GMazeEquipInfoV8Cfg.GetWithCtx(ctx,newEquipInfo.GetEquipId())
//	if equipCfg == nil {
//		logger.CtxError(ctx,"EquipSimplifyToCliPB get doll equip info cfg nil", zap.Int32("equipId", newEquipInfo.GetEquipId()))
//		return nil, errors.New("装备详情配置不存在")
//	}
//	res.Pos = proto.Int32(equipCfg.Pos)
//	res.EquipQuality = proto.Int32(equipCfg.Quality)
//	res.EquipLevel = proto.Int32(equipCfg.Level)
//	res.EquipId = proto.Int32(equipCfg.Equipment_id)
//
//	if newEquipInfo.GetSuitId() > 0 {
//		res.SuitInfo = &MazeGameEquip.EquipSuitInfo{
//			SuitId: proto.Int32(newEquipInfo.GetSuitId()),
//		}
//	}
//	_, equipName, mazeModel, icon, iconAtlas := pbutil.GetDollEquipNameEx(newEquipInfo, equipCfg.Pos_sub_type)
//	//	res.EquipResId = proto.Int32(equipResId)
//	res.EquipName = proto.String(equipName)
//	res.MazeModel = proto.Int32(mazeModel)
//	res.Icon = proto.String(icon)
//	res.IconAtlas = proto.String(iconAtlas)
//	return res, nil
// }
