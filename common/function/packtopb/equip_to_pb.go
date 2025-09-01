package packtopb

import (
	"context"
	"sort"

	"maze_game_server/config/GMazeEquipAffixOrderV8Cfg"
	"maze_game_server/pb/common/MazeGameEquip"

	"maze_game_server/config/GMazeAttributeV8Cfg"
	"maze_game_server/config/GMazeEquipConfigV8Cfg"

	"maze_game_server/common/errors"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"go.uber.org/zap"

	"maze_game_server/common/constdef"
	"maze_game_server/common/function/pbutil"
	"maze_game_server/common/function/randfuncs"
	"maze_game_server/config/GMazeAttrSpDescV8Cfg"
	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/pb/server/MazeEquipCache"

	"google.golang.org/protobuf/proto"
)

func EquipInfoToCliPB(ctx context.Context, equipInfo *MazeEquipCache.MazeEquipInfoDb) (*MazeGameEquip.MazeEquipInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeGameEquip.MazeEquipInfo{}
	res.EquipGuid = proto.Int64(equipInfo.GetEquipGuid())
	equipCfg := GMazeEquipInfoV8Cfg.GetWithCtx(ctx, equipInfo.GetEquipId())
	if equipCfg == nil {
		logger.CtxError(ctx, "EquipInfoToCliPB get doll equip info cfg nil", zap.Int32("equipId", equipInfo.GetEquipId()))
		return nil, errors.New("装备详情配置不存在")
	}
	equipSubType := equipCfg.Pos_sub_type
	if equipInfo.GetEquipSubType() > 0 {
		equipSubType = equipInfo.GetEquipSubType()
	}
	res.Pos = proto.Int32(equipCfg.Pos)
	res.EquipQuality = proto.Int32(equipCfg.Quality)
	res.EquipLevel = proto.Int32(equipCfg.Level)
	res.EquipId = proto.Int32(equipCfg.Equipment_id)
	_, equipName, mazeModel, icon, iconAtlas := pbutil.GetDollEquipNameEx(ctx, equipInfo, equipSubType)
	//	res.EquipResId = proto.Int32(equipResId)
	res.EquipName = proto.String(equipName)
	res.MazeModel = proto.Int32(mazeModel)
	res.Icon = proto.String(icon)
	res.IconAtlas = proto.String(iconAtlas)
	mainAttrs, baseAttrs, _, _, _ := EquipBaseAttrToCliPB(ctx, equipInfo.BaseAttrs)
	EquipBaseAttrSort(ctx, baseAttrs, equipCfg.Pos)
	res.MainAttrs = mainAttrs
	res.BaseAttrs = baseAttrs
	if equipInfo.GetSuitId() > 0 {
		res.SuitInfo = &MazeGameEquip.EquipSuitInfo{
			SuitId: proto.Int32(equipInfo.GetSuitId()),
		}
	}
	return res, nil
}

func EquipInfoToCliPBEx(ctx context.Context, equip *MazeEquipCache.MazeEquipInfoDb) (*MazeGameEquip.MazeEquipInfo, int32, error) {
	// newEquipInfo, err := pbutil.ConvertIdentifyEquipDb(logger, equip)
	// if err != nil {
	//	logger.CtxError(ctx,"EquipInfoToCliPBEx ConvertIdentifyEquipDb error", zap.Any("equip", equip), zap.Error(err))
	//	return nil, 0, err
	// }
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeGameEquip.MazeEquipInfo{}
	res.EquipGuid = proto.Int64(equip.GetEquipGuid())
	equipCfg := GMazeEquipInfoV8Cfg.GetWithCtx(ctx, equip.GetEquipId())
	if equipCfg == nil {
		logger.CtxError(ctx, "EquipInfoToCliPB get doll equip info cfg nil", zap.Int32("equipId", equip.GetEquipId()))
		return nil, 0, errors.New("装备详情配置不存在")
	}
	equipSubType := equipCfg.Pos_sub_type
	if equip.GetEquipSubType() > 0 {
		equipSubType = equip.GetEquipSubType()
	}
	res.Pos = proto.Int32(equipCfg.Pos)
	res.EquipQuality = proto.Int32(equipCfg.Quality)
	res.EquipLevel = proto.Int32(equipCfg.Level)
	res.EquipId = proto.Int32(equipCfg.Equipment_id)
	equipResId, equipName, mazeModel, icon, iconAtlas := pbutil.GetDollEquipNameEx(ctx, equip, equipSubType)
	//	res.EquipResId = proto.Int32(equipResId)
	res.EquipName = proto.String(equipName)
	res.MazeModel = proto.Int32(mazeModel)
	res.Icon = proto.String(icon)
	res.IconAtlas = proto.String(iconAtlas)
	mainAttrs, baseAttrs, _, _, _ := EquipBaseAttrToCliPB(ctx, equip.BaseAttrs)
	EquipBaseAttrSort(ctx, baseAttrs, equipCfg.Pos)
	res.MainAttrs = mainAttrs
	res.BaseAttrs = baseAttrs
	if equip.GetSuitId() > 0 {
		res.SuitInfo = &MazeGameEquip.EquipSuitInfo{
			SuitId: proto.Int32(equip.GetSuitId()),
		}
	}
	return res, equipResId, nil
}

func EquipBaseAttrSort(ctx context.Context, equipAttrs []*MazeGameEquip.BaseAttrInfo, pos int32) {
	attrSortMap := make(map[int32]int32, 0)
	attrSortCfg := GMazeEquipAffixOrderV8Cfg.GetWithCtx(ctx, pos)
	if attrSortCfg != nil {
		attrSortMap = attrSortCfg.Attr_order
	}
	// 先按类型排序 再按属性权重排序 最后按属性id排序
	sort.Slice(equipAttrs, func(i, j int) bool {
		k1 := equipAttrs[i].GetAttrInfo().GetAttrId()
		k2 := equipAttrs[j].GetAttrInfo().GetAttrId()
		if attrSortMap[k1] != attrSortMap[k2] {
			return attrSortMap[k1] > attrSortMap[k2]
		}
		if equipAttrs[i].GetAttrType() != equipAttrs[j].GetAttrType() {
			return equipAttrs[i].GetAttrType() < equipAttrs[j].GetAttrType()
		}
		if equipAttrs[i].GetIndex() != equipAttrs[j].GetIndex() {
			return equipAttrs[i].GetIndex() < equipAttrs[j].GetIndex()
		}
		return k1 > k2
	})
	return
}

// 只打包装备guid
func EquipInfoToCliPbGuid(equipInfo *MazeEquipCache.MazeEquipInfoDb) *MazeGameEquip.MazeEquipInfo {
	res := &MazeGameEquip.MazeEquipInfo{}
	res.EquipGuid = proto.Int64(equipInfo.GetEquipGuid())
	return res
}

func EquipAttrToCliPb(ctx context.Context, showAttrInfo *MazeEquipCache.EquipAttrInfo, showAttrMin, showAttrMax map[int32]int64) *MazeGameEquip.EquipAttrInfo {
	attrInfo := &MazeGameEquip.EquipAttrInfo{}
	attrInfo.AttrId = proto.Int32(showAttrInfo.GetAttrId())
	attrInfo.Value = proto.Int64(showAttrInfo.GetAttrValue())
	// attrInfo.MinValue = proto.Int64(showAttrMin[showAttrInfo.GetAttrId()])
	// attrInfo.MaxValue = proto.Int64(showAttrMax[showAttrInfo.GetAttrId()])
	var randWeight int32
	if showAttrMin[showAttrInfo.GetAttrId()] == showAttrMax[showAttrInfo.GetAttrId()] {
		randWeight = 10000
	} else {
		randWeight = int32((showAttrInfo.GetAttrValue() - showAttrMin[showAttrInfo.GetAttrId()]) * 10000 / (showAttrMax[showAttrInfo.GetAttrId()] - showAttrMin[showAttrInfo.GetAttrId()]))
	}
	attrInfo.ScoreTap = proto.Int32(GetEquipAttrScoreTap(ctx, randWeight))
	attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(attrInfo.GetAttrId())
	if attrCfg == nil {
		return attrInfo
	}
	attrInfo.Figure = proto.Int32(attrCfg.Figure)
	attrSpDescCfg := GMazeAttrSpDescV8Cfg.GetWithCtx(context.Background(), showAttrInfo.GetAttrId(), config_manager.QueryNullable())
	if attrSpDescCfg != nil {
		attrInfo.AttrDesc = proto.String(attrSpDescCfg.Equip_affix_desc)
		attrInfo.AttrUnits = proto.String(attrSpDescCfg.Equip_affix_suffix)
		attrInfo.SpecialSymbol = proto.String(attrSpDescCfg.Equip_symbol)
	}
	if attrInfo.GetAttrDesc() == "" {
		attrInfo.AttrDesc = proto.String(attrCfg.Name)
	}
	return attrInfo
}

func GetEquipAttrType(attrId int32, attrType int32) int32 {
	// for _, cfg := range GDollEquipModTypeV8Cfg.GetAll() {
	//	if cfg.Show_id == attrId {
	//		return int32(DollEquip.ENUM_EQUIP_ATTR_TYPE_ATTR_RESIST)
	//	}
	// }
	return attrType
}

func GetEquipAttackScopeCfg(ctx context.Context) int32 {
	equipConfig := GMazeEquipConfigV8Cfg.GetWithCtx(ctx, constdef.DollEquipCfg501)
	if equipConfig != nil {
		for _, v := range equipConfig.Value_map {
			return int32(v)
		}
	}
	return 0
}
func GetEquipAttackNumberCfg(ctx context.Context) int32 {
	equipConfig := GMazeEquipConfigV8Cfg.GetWithCtx(ctx, constdef.DollEquipCfg502)
	if equipConfig != nil {
		for _, v := range equipConfig.Value_map {
			return int32(v)
		}
	}
	return 0
}

func GetEquipAttrScoreTap(ctx context.Context, randWeight int32) int32 {
	equipConfig := GMazeEquipConfigV8Cfg.GetWithCtx(ctx, constdef.DollEquipCfg701)
	if equipConfig == nil {
		return 0
	}
	weightSlice := make([]*randfuncs.ItemWeight, 0)
	for weight, id := range equipConfig.Value_map {
		if weight == 0 { // 只过滤权重
			continue
		}
		weightSlice = append(weightSlice, &randfuncs.ItemWeight{Id: int32(id), Weight: weight})
	}
	sort.Slice(weightSlice, func(i, j int) bool {
		return weightSlice[i].Weight < weightSlice[j].Weight
	})

	for _, v := range weightSlice {
		if randWeight <= v.Weight {
			return v.Id
		}
	}
	return 0
}
