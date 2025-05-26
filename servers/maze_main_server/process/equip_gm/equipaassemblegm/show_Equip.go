/*
 * @Author: majian
 * @Date: 2024-08-21 13:40:25
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-01 19:22:21
 */
package equipaassemblegm

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/common/function/assemble"
	"maze_game_server/common/function/mazeutil"
	"maze_game_server/common/function/packtopb"
	"maze_game_server/common/function/pbutil"
	"maze_game_server/common/vardef"
	"maze_game_server/config/GMazeAttrSpDescV8Cfg"
	"maze_game_server/config/GMazeAttributeV8Cfg"
	"maze_game_server/config/GMazeEquipAffixRandPoolV8Cfg"
	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/config/GMazeEquipPosRankV8Cfg"
	"maze_game_server/config/GMazeEquipSuiteInfoV8Cfg"
	"maze_game_server/module/effectequip"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"
)

func DumpEquipPos(logger fklog.FKLogI, uid uint64, w *bytes.Buffer, pos int32, posInfo *MazeEquipCache.MazeEquipPosInfo, suitId int32) {
	var myName string
	posNameRow := GMazeEquipPosRankV8Cfg.GetMazeEquipPosRankV8Config(pos)
	if posNameRow != nil {
		myName = posNameRow.Name
	}
	if assemble.IsAssembleEquip(posInfo) {
		w.WriteString(fmt.Sprintf("\n\n装备位: %d(%s) 已装配\n", pos, myName))
		ew, e := ShowEquip(logger, uid, posInfo.GetEquipInfo(), posInfo.GetForce(), suitId, ShowEquipParam{
			IsActvie: posInfo.GetEquipLoadInfo().GetActivateMask(),
			PosLevel: posInfo.GetEquipPos().GetLevel()})
		if e != nil {
			return
		}
		w.WriteString(ew)
	} else {
		w.WriteString(fmt.Sprintf("\n\n装备位: %d(%s) 未装配\n", pos, myName))
	}
}

func GetEquipInfoByCfgId(logger fklog.FKLogI, userId uint64, cond BagCond) (rs string, err error) {
	equips, err := QueryBagByCond(logger, userId, cond)
	if err != nil {
		return
	}
	var strs []string
	totalMap := make(map[int32]int32)
	var totalCnt int32
	for _, equip := range equips {
		var str string
		str, err = ShowEquip(logger, userId, equip, 0, 0, ShowEquipParam{
			EquipSubType: equip.GetEquipSubType(), ResID: pbutil.GetEquipResId(logger, equip)})
		if err != nil {
			return
		}
		strs = append(strs, str)
		row := GMazeEquipInfoV8Cfg.Get(equip.GetEquipId())
		if row != nil {
			totalMap[row.Pos] += 1
			totalCnt++
		}
	}
	var bs bytes.Buffer
	bs.WriteString(fmt.Sprintf("查询结果:总共%d个\n", totalCnt))
	for k, v := range totalMap {
		posNameRow := GMazeEquipPosRankV8Cfg.GetMazeEquipPosRankV8Config(k)
		if posNameRow != nil {
			bs.WriteString(fmt.Sprintf("%s:%d个\n", posNameRow.Name, v))
		}
	}
	rs = strings.Join(strs, "\n---------------------\n")
	bs.WriteString(rs)
	return bs.String(), nil
}

func GetEquipInfoByGuid(logger fklog.FKLogI, userId uint64, guid int64) (rs string, err error) {
	equip, err := effectequip.GetEffectEquipInfo(logger, userId, guid)
	if err != nil {
		return
	}
	if equip == nil {
		rs = "装备不存在"
		return
	}
	return ShowEquip(logger, userId, equip, 0, 0, ShowEquipParam{
		EquipSubType: equip.GetEquipSubType(),
		ResID:        pbutil.GetEquipResId(logger, equip)})
}

type ShowEquipParam struct {
	EquipSubType int32 // 武器子类型
	ResID        int32 // 装备资源ID
	SuitID       int32 // 套装Id
	IsActvie     int32 // 套装是否激活
	PosLevel     int32 // 装备位强化等级
}

func ShowEquip(logger fklog.FKLogI, userId uint64, equipDb *MazeEquipCache.MazeEquipInfoDb, force int64, suitId int32, p ShowEquipParam) (rs string, err error) {
	cliEquip, err := packtopb.EquipInfoToCliPB(logger, equipDb)
	if err != nil {
		return
	}
	// if force <= 0 {
	// 	var equipForce map[int64]int64
	// 	equipForce, _, err = equipforcepreview.DollBagEquipForcePreview2(logger, userId,
	// 		[]*MazeEquipCache.MazeEquipInfoDb{equipDb}, "maze_equip_gm_server", false)
	// 	if err != nil {
	// 		return
	// 	}
	// 	cliEquip.ForceValue = proto.Int64(equipForce[equipDb.GetEquipGuid()])
	// } else {
	// 	cliEquip.ForceValue = proto.Int64(force)
	// }
	var w bytes.Buffer
	var posName, suitName, subTypeName string
	posNameRow := GMazeEquipPosRankV8Cfg.GetMazeEquipPosRankV8Config(cliEquip.GetPos())
	if posNameRow != nil {
		posName = posNameRow.Name
		subTypeName = posNameRow.Sub_type_name[p.EquipSubType]
	}

	w.WriteString(fmt.Sprintf("实例Id: %d\n", cliEquip.GetEquipGuid()))
	w.WriteString(fmt.Sprintf("配置ID: %d\n", cliEquip.GetEquipId()))
	w.WriteString(fmt.Sprintf("部位: (%d)%s\n", cliEquip.GetPos(), posName))
	if cliEquip.GetPos() == 1 {
		w.WriteString(fmt.Sprintf("武器子类型: (%d)%s\n", p.EquipSubType, subTypeName))
	}
	isMaze := mazeutil.CheckMazeEquip(cliEquip.GetEquipId())
	if isMaze {
		w.WriteString("是否迷宫装备:Yes\n")
	} else {
		w.WriteString("是否迷宫装备:No\n")
	}

	w.WriteString(fmt.Sprintf("部位强化等级:%d\n", p.PosLevel))
	w.WriteString(fmt.Sprintf("掉落时间:%s\n", time.Unix(equipDb.GetMakeTime(), 0).Format("2006-01-02 15:04:05")))
	w.WriteString(fmt.Sprintf("武力值: %d\n", cliEquip.GetForceValue()))
	w.WriteString(fmt.Sprintf("品质: %d(%s)\n", cliEquip.GetEquipQuality(),
		vardef.DollEquipQualityMap[cliEquip.GetEquipQuality()]))
	w.WriteString(fmt.Sprintf("穿戴等级: %d\n", cliEquip.GetEquipLevel()))
	w.WriteString(fmt.Sprintf("资源配置ID: %d\n", p.ResID))
	resId, _ := pbutil.GetDollEquipName(equipDb, equipDb.GetEquipSubType())
	w.WriteString(fmt.Sprintf("展示资源配置ID: %d\n", resId))

	_, equipName := pbutil.GetDollEquipName(equipDb, 0) // 子类型无用，内部会自己取
	w.WriteString(fmt.Sprintf("装备名: %s\n", equipName))

	suitRow := GMazeEquipSuiteInfoV8Cfg.GetMazeEquipSuiteInfoV8Config(equipDb.GetSuitId())
	if suitRow != nil {
		suitName = suitRow.Suite_name
	}
	w.WriteString(fmt.Sprintf("所属套装: %d(%s)\n", equipDb.GetSuitId(), suitName))
	w.WriteString(fmt.Sprintf("装配信息里套装ID: %d\n", p.SuitID))
	if p.SuitID > 0 {
		if assemble.IsHurtSuitActivate(p.IsActvie) {
			w.WriteString("套装效果:已激活\n")
		} else {
			w.WriteString("套装效果:未激活\n")
		}
	}
	// if equipDb.GetEnterTime() > 0 {
	// 	w.WriteString(fmt.Sprintf("进临时背包时间: %s\n", time.Unix(equipDb.GetEnterTime(), 0).Format("2006-01-02 15:04:05")))
	// }
	// if cliEquip.GetLock()&1 > 0 {
	// 	w.WriteString("是否锁定: 是")
	// } else {
	// 	w.WriteString("是否锁定: 否")
	// }
	// if equipDb.GetNewFlag() > 0 {
	// 	w.WriteString("\n是否有new标记: 是")
	// } else {
	// 	w.WriteString("\n是否有new标记: 否")
	// }

	// if equipDb.GetEquipStatus() > 0 {
	// 	if equipDb.GetEquipStatus()&1 > 0 {
	// 		w.WriteString("\n状态:待解封")
	// 		if equipDb.GetRuleId() > 0 {
	// 			w.WriteString(fmt.Sprintf("\n源装备ID:%d", equipDb.GetRuleId()))
	// 		}
	// 	} else {
	// 		if equipDb.GetEquipStatus()&2 > 0 {
	// 			w.WriteString("\n状态:已解封")
	// 		}
	// 	}
	// } else {
	// 	w.WriteString("\n状态:无")
	// }
	// w.WriteString(fmt.Sprintf("\n槽位数量:%d", equipDb.GetSkillSlotNum()))
	w.WriteString("\n基础属性\n")
	baseAttrs := equipDb.GetBaseAttrs()
	for _, attrInfo := range baseAttrs {
		var entryInfo string
		entryInfo += fmt.Sprintf("词条Id:%d\t", attrInfo.GetAttrGroup())
		entryInfo += fmt.Sprintf("Roll值:%d\t", attrInfo.GetRandWeight())
		if attrInfo.GetAttrType() == 1 {
			entryInfo += "(主)\t"
		}
		attrCfg := GMazeEquipAffixRandPoolV8Cfg.Get(attrInfo.GetAttrGroup())
		for _, showAttr := range attrInfo.GetShowAttrList() {
			row := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(showAttr.GetAttrId())
			if row == nil {
				continue
			}
			attrName := row.Name
			spRow := GMazeAttrSpDescV8Cfg.GetMazeAttrSpDescV8Config(showAttr.GetAttrId())
			if spRow != nil {
				if spRow.Equip_affix_desc != "" {
					attrName = spRow.Equip_affix_desc
				}
			}
			var min, max int64
			if attrCfg != nil {
				min = attrCfg.Show_attr_min[showAttr.GetAttrId()]
				max = attrCfg.Show_attr_max[showAttr.GetAttrId()]
			}
			dw := GetColorLevel(cliEquip, attrInfo.GetIndex(), attrInfo.GetAttrType(), showAttr.GetAttrId())
			entryInfo += fmt.Sprintf("%s(%d)(显)(档位%d):%d(minmax%d-%d)\t|\t", attrName, showAttr.GetAttrId(), dw, showAttr.GetAttrValue(), min, max)
		}

		for _, realAttr := range attrInfo.GetRealAttrList() {
			row := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(realAttr.GetAttrId())
			if row == nil {
				continue
			}
			attrName := row.Name
			spRow := GMazeAttrSpDescV8Cfg.GetMazeAttrSpDescV8Config(realAttr.GetAttrId())
			if spRow != nil {
				if spRow.Equip_affix_desc != "" {
					attrName = spRow.Equip_affix_desc
				}
			}

			var min, max int64
			if attrCfg != nil {
				min = attrCfg.Add_attr_min[realAttr.GetAttrId()]
				max = attrCfg.Add_attr_max[realAttr.GetAttrId()]
			}
			entryInfo += fmt.Sprintf("%s(%d)(真):%d(minmax%d-%d)\t|\t", attrName, realAttr.GetAttrId(), realAttr.GetAttrValue(), min, max)
		}
		entryInfo += "\n"
		w.WriteString(entryInfo)
	}

	w.WriteString("\n")
	return w.String(), nil
}

func GetColorLevel(in *MazeGameEquip.MazeEquipInfo, index, typ, attrId int32) int32 {
	for _, attr := range in.GetBaseAttrs() {
		if attr.GetAttrType() == typ && attr.GetIndex() == index {
			if attr.GetAttrInfo().GetAttrId() == attrId {
				return attr.AttrInfo.GetScoreTap()
			}
		}
	}
	// for _, attr := range in.GetSpecialAttrs() {
	// 	if attr.GetAttrType() == typ && attr.GetIndex() == index {
	// 		if attr.GetAttrInfo().GetAttrId() == attrId {
	// 			return attr.AttrInfo.GetScoreTap()
	// 		}
	// 	}
	// }
	return 0
}
