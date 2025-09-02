/*
 * @Author: majian
 * @Date: 2025-01-08 21:13:42
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-01 15:18:53
 */
package calcassembleattr

import (
	"context"
	"errors"

	"maze_game_server/common/constdef"
	"maze_game_server/common/function/assemble"
	"maze_game_server/config/GMazeAttributeV8Cfg"
	"maze_game_server/pb/server/MazeBuffData"
	"maze_game_server/pb/server/MazeEquipCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type EquipmentEffectInfo struct {
	ForceAttrs   map[int32]int64          // 武力属性
	Other        *MazeBuffData.MazeBuffDb // 所有属性
	SuitId       int32                    // 激活的套装Id
	ChgPoss      map[int32]int32          // 有变化的格子
	SuitCalc     *EquipSuitMgr            // 套装计算信息
	SkillSlotCnt int32                    // 技能槽位数量
	// FiveStateMap map[int32]int32                 // 五行属性状态
	fklog.FKLogI
}

func NewEquipmentEffectInfo(ctx context.Context) *EquipmentEffectInfo {
	r := new(EquipmentEffectInfo)
	r.ChgPoss = make(map[int32]int32)
	r.ForceAttrs = make(map[int32]int64)
	//	r.FiveStateMap = make(map[int32]int32)
	r.Other = &MazeBuffData.MazeBuffDb{} // 其他属性
	r.FKLogI = fklog.ContextAppLogger(ctx)
	return r
}

func (m *EquipmentEffectInfo) GetSuitCalc() *EquipSuitMgr {
	if m != nil && m.SuitCalc != nil {
		return m.SuitCalc
	}
	return nil
}

func (m *EquipmentEffectInfo) PackEquipRealAttr(eAttr *MazeEquipCache.EquipAttrInfo) error {
	attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(eAttr.GetAttrId())
	if attrCfg == nil {
		m.ErrorWF("PackEquipRealAttr cannot find attr", zap.Int32("attrId", eAttr.GetAttrId()))
		return errors.New("no found attr cfg")
	}
	if attrCfg.Type == constdef.DollAttrTypeForce {
		m.ForceAttrs[eAttr.GetAttrId()] += eAttr.GetAttrValue() // 武力值属性
	}
	AppendDollAttr(m.Other, PackCommonAttr(eAttr))
	return nil
}

func (m *EquipmentEffectInfo) PackEquipShowAttr(eAttr *MazeEquipCache.EquipAttrInfo) error {
	attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(eAttr.GetAttrId())
	if attrCfg == nil {
		m.ErrorWF("PackEquipShowAttr cannot find attr", zap.Int32("attrId", eAttr.GetAttrId()))
		return errors.New("no found attr cfg")
	}
	AppendDollPanelAttr(m.Other, PackCommonAttr(eAttr))
	return nil
}

func (m *EquipmentEffectInfo) PackEquipRealAttrKv(k int32, v int64) error {
	attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(k)
	if attrCfg == nil {
		m.ErrorWF("PackEquipRealAttrKv cannot find buff attr", zap.Int32("attrId", k))
		return errors.New("no found attr cfg")
	}
	if attrCfg.Type == constdef.DollAttrTypeForce {
		m.ForceAttrs[k] += v // 武力值属性
	}
	AppendDollAttrKv(m.Other, k, v)
	return nil
}

func (m *EquipmentEffectInfo) PackEquipShowAttrKv(k int32, v int64) error {
	attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(k)
	if attrCfg == nil {
		m.ErrorWF("PackEquipShowAttrKv cannot find show buff attr", zap.Int32("attrId", k))
		return errors.New("no found attr cfg")
	}
	AppendDollPanelAttrKv(m.Other, k, v)
	return nil
}

// func (m *EquipmentEffectInfo) CalcFElemAttr(equip *MazeEquipCache.MazeEquipInfoDb) error {
// 	fiveAttrs := equip.GetFiveElemInfo().GetElemAttrs()
// 	for _, attr := range fiveAttrs {
// 		for _, realAttr := range attr.GetRealAttrList() {
// 			e := m.PackEquipRealAttr(realAttr)
// 			if e != nil {
// 				return e
// 			}
// 		}

// 		for _, realAttr := range attr.GetShowAttrList() {
// 			e := m.PackEquipShowAttr(realAttr)
// 			if e != nil {
// 				return e
// 			}
// 		}
// 	}
// 	return nil
// }

func (m *EquipmentEffectInfo) CalcBaseAttr(equip *MazeEquipCache.MazeEquipInfoDb) error {
	baseAttrs := equip.GetBaseAttrs()
	for _, attr := range baseAttrs {
		for _, realAttr := range attr.GetRealAttrList() {
			e := m.PackEquipRealAttr(realAttr)
			if e != nil {
				return e
			}
		}

		for _, realAttr := range attr.GetShowAttrList() {
			e := m.PackEquipShowAttr(realAttr)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

// func (m *EquipmentEffectInfo) CalcModAttr(equip *MazeEquipCache.MazeEquipInfoDb) error {
// 	modAttrs := equip.GetModAttrs()
// 	for _, attr := range modAttrs {
// 		for _, realAttr := range attr.GetRealAttrList() {
// 			e := m.PackEquipRealAttr(realAttr)
// 			if e != nil {
// 				return e
// 			}
// 		}

// 		for _, realAttr := range attr.GetShowAttrList() {
// 			e := m.PackEquipShowAttr(realAttr)
// 			if e != nil {
// 				return e
// 			}
// 		}
// 	}
// 	return nil
// }

// func (m *EquipmentEffectInfo) CalcSpAttr(equip *MazeEquipCache.MazeEquipInfoDb) error {
// 	spAttrs := equip.GetSpecialAttrs()
// 	for _, attr := range spAttrs {
// 		for _, realAttr := range attr.GetRealAttrList() {
// 			e := m.PackEquipRealAttr(realAttr)
// 			if e != nil {
// 				return e
// 			}
// 		}

// 		for _, realAttr := range attr.GetShowAttrList() {
// 			e := m.PackEquipShowAttr(realAttr)
// 			if e != nil {
// 				return e
// 			}
// 		}
// 	}
// 	return nil
// }

// 计算伤害套装
func (m *EquipmentEffectInfo) CalcHurtSuit(ctx context.Context, equips []*MazeEquipCache.MazeEquipPosInfo) error {
	// 计算套装
	suitMgr, err := CalcEquipSuit(ctx, equips)
	if err != nil {
		m.FKLogI.CtxError(ctx, "EquipmentEffectInfo CalcEquipSuit err", zap.Error(err))
		return err
	}
	m.SuitCalc = suitMgr
	// // 检查套装是否激活
	m.SuitId = suitMgr.GetActiveSuitId()

	// 计算传说套装加成
	suitCntMap := suitMgr.GetSuitNumMap()
	buffs, showBuffs := CalcLegendSuitBuff(suitCntMap)
	for k, v := range buffs {
		if k <= 0 || v <= 0 {
			continue
		}
		e := m.PackEquipRealAttrKv(k, v)
		if e != nil {
			m.FKLogI.CtxError(ctx, "EquipmentEffectInfo buffs PackEquipRealAttrKv err", zap.Error(e), zap.Int32("attrId", k), zap.Int64("attrValue", v))
			return e
		}
	}

	for k, v := range showBuffs {
		if k <= 0 || v <= 0 {
			continue
		}
		e := m.PackEquipRealAttrKv(k, v)
		if e != nil {
			m.FKLogI.CtxError(ctx, "EquipmentEffectInfo showBuffs PackEquipRealAttrKv err", zap.Error(e), zap.Int32("attrId", k), zap.Int64("attrValue", v))
			return e
		}
	}

	return nil
}

type EffectCalcInParam struct {
	IsLog   bool // 是否打开日志
	IsForce bool // 是否仅需要武力属性
	// IsFiveOnlyRead bool // 五行激活是否只读，不用计算
}

// 计算装备的效果
func CalcEquipEffectAll(ctx context.Context, equips []*MazeEquipCache.MazeEquipPosInfo, ep EffectCalcInParam) (effect *EquipmentEffectInfo, err error) {
	effect = NewEquipmentEffectInfo(ctx)
	var step string
	defer func() {
		if err != nil {
			effect.FKLogI.CtxError(ctx, "CalcEquipEffectAll dump err", zap.Error(err),
				zap.Any("equips", equips), zap.Any("effect", effect),
				zap.String("step", step),
				zap.Any("suitMgr", effect.SuitCalc))
		} else {
			if ep.IsLog {
				effect.FKLogI.CtxInfo(ctx, "CalcEquipEffectAll dump",
					zap.Any("equips", equips), zap.Any("effect", effect),
					zap.Any("suitMgr", effect.SuitCalc))
			}
		}
	}()

	// 计算套装
	err = effect.CalcHurtSuit(ctx, equips)
	if err != nil {
		effect.FKLogI.CtxError(ctx, "CalcHurtSuit err", zap.Error(err))
		return
	}

	// equipPosMap := make(map[int32]*MazeEquipCache.MazeEquipPosInfo) // <位置:装配信息>
	// for _, equipPos := range equips {
	// 	if !assemble.IsAssembleEquip(equipPos) {
	// 		continue
	// 	}
	// 	equipPosMap[equipPos.GetEquipPos().GetPos()] = equipPos
	// }

	suitMap := effect.SuitCalc.GetSuitNumMap()
	for _, equipPos := range equips {
		if !assemble.IsAssembleEquip(equipPos) {
			continue
		}
		pos := equipPos.GetEquipLoadInfo().GetPos()
		// 检查单个装备套装效果是否激活
		suitId := equipPos.GetEquipInfo().GetSuitId()
		oldVal := equipPos.GetEquipLoadInfo().GetActivateMask()
		activateMask := int32(0)
		row := GetLeastRow(suitId, suitMap[suitId])
		if row != nil {
			activateMask |= constdef.HurtSuitActivate
		}

		// var activate bool
		// activate, err = CalcFiveElemActive(logger, equipPos, equipPosMap)
		// if err != nil {
		// 	step = "calc five elem active"
		// 	return nil, err
		// }
		// if activate {
		// 	activateMask |= constdef.FiveElemActivate
		// }

		if oldVal != activateMask {
			effect.ChgPoss[pos] += 1 // 变化次数
			equipPos.EquipLoadInfo.ActivateMask = proto.Int32(activateMask)
		}
		if pos == 1 {
			txAttr := CalcElementEffect(ctx, equipPos, suitMap, ep.IsLog)
			if txAttr != nil {
				AppendDollAttr(effect.Other, txAttr)
			}
		}

		equip := equipPos.GetEquipInfo()

		err = effect.CalcBaseAttr(equip)
		if err != nil {
			step = "calc base attr"
			return
		}
		// err = effect.CalcModAttr(equip)
		// if err != nil {
		// 	step = "calc mod attr"
		// 	return nil, err
		// }
		// err = effect.CalcSpAttr(equip)
		// if err != nil {
		// 	step = "calc sp attr"
		// 	return nil, err
		// }
		// 添加五行属性buff
		// if assemble.IsFiveElemActivate(equipPos.GetEquipLoadInfo().GetActivateMask()) {
		// 	effect.FiveStateMap[pos] = 1 // 记录流水用
		// 	err = effect.CalcFElemAttr(equip)
		// 	if err != nil {
		// 		step = "calc five elem attr"
		// 		return nil, err
		// 	}
		// }
	}
	return
}

// 是否有手势技能加成
// func (m *EquipmentEffectInfo) HasHandSkillBuff() bool {
// 	rs := dollhandskillcfgv8.GetAttrIdToSkilIdMap()
// 	for _, realAttr := range m.Other.GetDollEquipAttrList() {
// 		if _, ok := rs[realAttr.GetAttrId()]; ok {
// 			return true
// 		}
// 	}
// 	return false
// }

// func HasHandSkillChg(m, n *EquipmentEffectInfo) bool {
// 	return m.SkillSlotCnt != n.SkillSlotCnt || m.HasHandSkillBuff() || n.HasHandSkillBuff()
// }

func (m *EquipmentEffectInfo) ChkPosChg(pos int32) bool {
	return m.ChgPoss[pos] > 0
}

// func ChkPosChg(m, n map[int32]*PosChgInfo, pos int32) bool {
// 	if m == nil && n == nil {
// 		return false
// 	}
// 	if m == nil || n == nil {
// 		return true
// 	}

// 	if m[pos].ActiveMask != n[pos].ActiveMask {
// 		return true
// 	}
// 	return false
// }
