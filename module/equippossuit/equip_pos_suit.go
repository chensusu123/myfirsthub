package equippossuit

import (
	"fmt"
	"sort"

	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipPosLvSuiteV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/Common"
	"gitlab.ifreetalk.com/plate/protodef/MazeEquipCache"
	"gitlab.ifreetalk.com/plate/protodef/MazeGameEquip"

	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
)

// 当前套装和下级套装
// 无激活套装:1级套装和1级套装(客户端需要)
// 套装已满级:当前套装和套装ID为-1的下级套装
// 已激活套装并且未满级:当前套装和下级套装
func GetCurAndNextSuit(logger fklog.FKLogI, curSuitId int32) (cur, next *MazeGameEquip.EquipPosSuitInfo, err error) {
	var suitId, nextSuitId int32
	if curSuitId == 0 {
		suitId = findFirstSuit()
	} else {
		suitId = curSuitId
	}

	defer func() {
		logger.InfoWF("GetCurAndNextSuit",
			zap.Int32("curSuitId", curSuitId),
			zap.Any("curSuit", cur), zap.Any("nextSuit", next),
			zap.Bool("is init suit", curSuitId == 0),
			zap.Bool("is max level suit", nextSuitId == -1))
	}()

	cfg := GMazeEquipPosLvSuiteV8Cfg.Get(suitId)
	if cfg == nil {
		err = fmt.Errorf("GetCurAndNextSuit no cfg, sheetName: %s, id: %d", GMazeEquipPosLvSuiteV8Cfg.SheetName(), suitId)
		logger.ErrorWF("GetCurAndNextSuit", zap.Error(err))
		return
	}
	nextSuitId = cfg.Next_level

	cur = &MazeGameEquip.EquipPosSuitInfo{
		SuitId: proto.Int32(cfg.Level),
		SuitLv: proto.Int32(cfg.Level),
	}
	for _, id := range cfg.Show_attr_order {
		if id <= 0 {
			continue
		}
		cur.Attrs = append(cur.Attrs, &Common.Attr{AttrId: proto.Int32(id), AttrValue: proto.Int32(int32(cfg.Show_attr[id]))})
	}

	if curSuitId == 0 { // 无激活套装
		next = cur
		return
	}

	if nextSuitId == -1 { // 已满级
		next = &MazeGameEquip.EquipPosSuitInfo{
			SuitId: proto.Int32(-1),
		}
		return
	}

	nextCfg := GMazeEquipPosLvSuiteV8Cfg.Get(nextSuitId)
	if nextCfg == nil {
		err = fmt.Errorf("GetCurAndNextSuit no cfg, sheetName: %s, id: %d", GMazeEquipPosLvSuiteV8Cfg.SheetName(), nextSuitId)
		logger.ErrorWF("GetCurAndNextSuit", zap.Error(err))
		return
	}

	next = &MazeGameEquip.EquipPosSuitInfo{
		SuitId: proto.Int32(nextCfg.Level),
		SuitLv: proto.Int32(nextCfg.Level),
	}
	for _, id := range nextCfg.Show_attr_order {
		if id <= 0 {
			continue
		}
		next.Attrs = append(next.Attrs, &Common.Attr{AttrId: proto.Int32(id), AttrValue: proto.Int32(int32(nextCfg.Show_attr[id]))})
	}

	return
}

// 装备位套装是否能升级
// func IsEquipPosSuitCanLevelUp(logger fklog.FKLogI, curSuitId int32, poss []*MazeEquipCache.MazeEquipPosInfo) (can bool, err error) {
// 	if len(poss) < constdef.EquipPosNum { // 有装备位未解锁
// 		logger.InfoWF("IsEquipPosSuitCanLevelUp equip pos num not enough",
// 			zap.Int("curNum", len(poss)), zap.Int("needNum", constdef.EquipPosNum))
// 		return
// 	}

// 	for _, pos := range poss {
// 		lv := pos.GetEquipPos().GetLevel()
// 		if lv <= curSuitId { // 装备等级不够
// 			logger.InfoWF("IsEquipPosSuitCanLevelUp equip pos level not enough",
// 				zap.Int32("posId", pos.GetEquipPos().GetPos()),
// 				zap.Int32("curLevel", lv), zap.Int32("needLevel", curSuitId))
// 			return
// 		}
// 	}

// 	cfg := GMazeEquipPosLvSuiteV8Cfg.Get(curSuitId)
// 	if cfg == nil {
// 		err = fmt.Errorf("GetCurAndNextSuit no cfg, sheetName: %s, id: %d", GMazeEquipPosLvSuiteV8Cfg.SheetName(), curSuitId)
// 		logger.ErrorWF("IsEquipPosSuitCanLevelUp", zap.Error(err))
// 		return
// 	}

// 	if cfg.Next_level == -1 { // 套装已满级
// 		logger.InfoWF("IsEquipPosSuitCanLevelUp suit is max level", zap.Int32("curSuitId", curSuitId))
// 		return
// 	}

// 	logger.InfoWF("IsEquipPosSuitCanLevelUp can level up", zap.Int32("curSuitId", curSuitId), zap.Int32("nextSuitId", cfg.Next_level))
// 	can = true

// 	return
// }

// 计算装备套装
func CalcPosSuit(logger fklog.FKLogI, poss []*MazeEquipCache.MazeEquipPosInfo) (suitId int32, err error) {
	if len(poss) < constdef.EquipPosNum { // 有装备位未解锁
		logger.InfoWF("IsEquipPosSuitCanLevelUp equip pos num not enough",
			zap.Int("curNum", len(poss)), zap.Int("needNum", constdef.EquipPosNum))
		return
	}
	var minLv int32 = -1
	for _, pos := range poss {
		lv := pos.GetEquipPos().GetLevel()
		if minLv == -1 || minLv > lv {
			minLv = lv
		}
	}

	allRows := GMazeEquipPosLvSuiteV8Cfg.GetAll()
	sort.Slice(allRows, func(i, j int) bool {
		return allRows[i].Level > allRows[j].Level
	})
	for _, row := range allRows {
		if row.Level <= minLv {
			return row.Level, nil
		}
	}
	return 0, nil
}

func findFirstSuit() int32 {
	allRows := GMazeEquipPosLvSuiteV8Cfg.GetAll()
	var firstSuit int32
	for _, row := range allRows {
		if firstSuit == 0 || row.Level < firstSuit {
			firstSuit = row.Level
		}
	}
	return firstSuit
}
