/*
 * @Author: majian
 * @Date: 2024-08-03 15:18:45
 * @Last Modified by: majian
 * @Last Modified time: 2024-12-26 21:52:36
 */
package equipaassemblegm

import (
	"bytes"

	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipInfoV8Cfg"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebagequipredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/effectequip"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassemblesuitredis"
)

type BagCond struct {
	EquipId    int32
	Pos        int32
	MinLv      int32
	MaxLv      int32
	Quality    int32
	Guid       int64
	Lock       int32
	MadeTime   int64
	RuleId     int32
	QueryInUse int32
	TmpBag     int32
	ShowSeal   int32
}

func CondHelp() string {
	var bs bytes.Buffer
	bs.WriteString("-----使用说明begin-----\n")
	bs.WriteString("equipId:指定装备配置ID\n")
	bs.WriteString("guid:指定装备guid\n")
	bs.WriteString("lock:1=只返回锁定的装备 0=没有限制\n")
	bs.WriteString("pos:指定装备位(1=武器 2=帽子 3=上衣 4=裤子 5=护肩 6=鞋子 7=护手 8=饰品)\n")
	bs.WriteString("minLv:指定穿戴等级下限\n")
	bs.WriteString("maxLv:指定穿戴等级上限\n")
	bs.WriteString("quality:指定装备品质(1=白(普通) 3=蓝(精良) 5=紫(史诗) 6=紫(史诗) 7=橙(传奇) 8=绿(套装))\n")
	bs.WriteString("ruleId:指定特殊规则ID(数值表里配置)\n")
	bs.WriteString("madetime:指定装备出厂时间,指定后只会返回指定时间及之前掉落的装备\n")
	bs.WriteString("queryInUse:0-查询不包含已装配的 1-查询包含已装配的\n")
	bs.WriteString("showseal:1 显示封印装备隐藏属性\n")
	bs.WriteString("不指定任何条件会返回背包里所有装备\n")
	bs.WriteString("不指定任何条件会返回背包里所有装备\n")
	bs.WriteString("-----使用说明end-----\n")
	return bs.String()
}

func QueryBagByCond(logger fklog.FKLogI, userId uint64, cond BagCond) (equips []*MazeEquipCache.MazeEquipInfoDb, err error) {
	var guidEquip map[int64]*MazeEquipCache.MazeEquipInfoDb
	if cond.ShowSeal == 1 {
		guidEquip, err = mazebagequipredis.GetAllEquipInfo(logger, userId)

	} else {
		guidEquip, err = effectequip.GetAllEffectEquipInfo(logger, userId)

	}
	if err != nil {
		return nil, err
	}
	// 过滤掉已装配的
	if cond.QueryInUse == 0 {
		suitEquips, err1 := dollassemblesuitredis.GetAllDollAssembleSuit(logger, userId)
		if err1 != nil {
			err = err1
			return
		}
		for _, equips := range suitEquips {
			for _, equip := range equips {
				guid := equip.GetEquipGuid()
				if _, ok := guidEquip[guid]; ok {
					delete(guidEquip, guid)
					logger.InfoWF("QueryBagByCond ignore equip", zap.Int64("guid", guid))
				}
			}
		}
	}

	for _, equip := range guidEquip {
		if cond.EquipId > 0 && equip.GetEquipId() != cond.EquipId {
			continue
		}

		if cond.Guid > 0 && equip.GetEquipGuid() != cond.Guid {
			continue
		}
		if cond.MadeTime > 0 && equip.GetMakeTime() > cond.MadeTime {
			continue
		}
		if cond.RuleId > 0 && equip.GetRuleId() != cond.RuleId {
			continue
		}
		cfg := GMazeEquipInfoV8Cfg.GetMazeEquipInfoV8Config(equip.GetEquipId())
		if cfg == nil {
			continue
		}
		if cond.MinLv > 0 && cfg.Level < cond.MinLv {
			continue
		}

		if cond.MaxLv > 0 && cfg.Level > cond.MaxLv {
			continue
		}
		if cond.Pos > 0 && cfg.Pos != cond.Pos {
			continue
		}
		if cond.Quality > 0 && cfg.Quality != cond.Quality {
			continue
		}
		// if cond.Lock > 0 && equip.GetLock() != cond.Lock {
		// 	continue
		// }
		// if cond.TmpBag == 1 { // 临时被阿伯
		// 	if equip.GetEnterTime() <= 0 {
		// 		continue
		// 	}
		// } else if cond.TmpBag == 2 { //正式背包
		// 	if equip.GetEnterTime() > 0 {
		// 		continue
		// 	}
		// }
		equips = append(equips, equip)
	}
	// if cond.TmpBag == 1 {
	// 	sort.Slice(equips, func(i, j int) bool {
	// 		return equips[i].GetEnterTime() <= equips[j].GetEnterTime()
	// 	})
	// } else if cond.TmpBag == 2 {
	// 	sort.Slice(equips, func(i, j int) bool {
	// 		return equips[i].GetMakeTime() <= equips[j].GetMakeTime()
	// 	})
	// }
	return
}
