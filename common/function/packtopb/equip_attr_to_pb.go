package packtopb

import (
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipAffixRandPoolV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGameEquip"

	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeequipconfigv8"
)

func EquipBaseAttrToCliPB(logger fklog.FKLogI, equipAttrs []*MazeEquipCache.BaseAttrInfo) (mainAttrs *MazeGameEquip.BaseAttrInfo, baseAttrs []*MazeGameEquip.BaseAttrInfo, attackScope, attackNumber, hurtType int32) {
	attackScopeAttrId := GetEquipAttackScopeCfg()
	attackNumberAttrId := GetEquipAttackNumberCfg()
	for _, attrInfoDb := range equipAttrs {
		if attrInfoDb.ShowAttrList[0].GetAttrId() == attackScopeAttrId {
			attackScope = int32(attrInfoDb.ShowAttrList[0].GetAttrValue())
			continue
		}
		if attrInfoDb.ShowAttrList[0].GetAttrId() == attackNumberAttrId {
			attackNumber = int32(attrInfoDb.ShowAttrList[0].GetAttrValue())
			continue
		}
		if hurtType <= 0 {
			for _, realAttr := range attrInfoDb.RealAttrList {
				hurtType = mazeequipconfigv8.GetSuitEffectType(realAttr.GetAttrId())
				if hurtType > 0 {
					break
				}
			}
			if hurtType > 0 {
				continue
			}
		}
		attrCfg := GMazeEquipAffixRandPoolV8Cfg.Get(attrInfoDb.GetAttrGroup())
		if attrCfg != nil {
			baseAttrInfo := &MazeGameEquip.BaseAttrInfo{}
			for _, showAttrInfo := range attrInfoDb.ShowAttrList {
				attrInfo := EquipAttrToCliPb(showAttrInfo, attrCfg.Show_attr_min, attrCfg.Show_attr_max)
				baseAttrInfo.AttrInfo = attrInfo
				break
			}
			baseAttrInfo.AttrType = proto.Int32(GetEquipAttrType(baseAttrInfo.GetAttrInfo().GetAttrId(), attrInfoDb.GetAttrType()))
			baseAttrInfo.Index = proto.Int32(attrInfoDb.GetIndex())
			if attrInfoDb.GetIndex() == 1 {
				mainAttrs = baseAttrInfo
				continue
			}
			baseAttrs = append(baseAttrs, baseAttrInfo)
		} else {
			logger.ErrorWF("maze_equip_affix_pool_v8【迷宫-装备-词条随机池】.xlsx 随机词条库缺失", zap.Any("池子id:", attrInfoDb.GetAttrGroup()))
		}
	}
	return
}
