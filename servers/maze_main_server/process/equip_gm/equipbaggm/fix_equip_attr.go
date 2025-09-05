package equipbaggm

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeEquipAffixRandPoolV8Cfg"
	"maze_game_server/config/GMazeEquipAffixRollTypeV8Cfg"
	"maze_game_server/io/redis/mazebagequipredis"
	"maze_game_server/pb/server/MazeEquipCache"
	"maze_game_server/servers/maze_main_server/process/equip"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func FixAllEquipAttrLimit(ctx context.Context, userId uint64) error {
	logger := fklog.ContextAppLogger(ctx)
	equipInfoMap, err := mazebagequipredis.GetAllEquipInfo(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "fixAllEquipAttrLimit GetAllEquipInfo fail", zap.Error(err))
		return err
	}
	chgEquipList := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
	for _, equipInfo := range equipInfoMap {
		newEquipInfo, chgEquip, err := fixEquipAttr(ctx, userId, equipInfo)
		if err != nil {
			logger.CtxError(ctx, "fixAllEquipAttrLimit fixEquipAttr fail", zap.Error(err))
			return err
		}
		if chgEquip {
			chgEquipList = append(chgEquipList, newEquipInfo)
		}
	}
	if len(chgEquipList) > 0 {
		mazebagequipredis.BatchSaveEquipInfo(ctx, userId, chgEquipList)
	}
	return err
}

func fixEquipAttr(ctx context.Context, userId uint64, equipInfo *MazeEquipCache.MazeEquipInfoDb) (*MazeEquipCache.MazeEquipInfoDb, bool, error) {
	logger := fklog.ContextAppLogger(ctx)
	chgEquip := false
	for _, attrInfo := range equipInfo.BaseAttrs {
		cfg := GMazeEquipAffixRandPoolV8Cfg.GetWithCtx(ctx, attrInfo.GetAttrGroup())
		if cfg == nil {
			logger.CtxError(ctx, "FixEquipAttr GMazeEquipAffixRandPoolV8Cfg fail",
				zap.Any("attrGroupId", attrInfo.GetAttrGroup()))
			return nil, false, errors.New("属性配置不存在")
		}
		rollTypeCfg := GMazeEquipAffixRollTypeV8Cfg.GetWithCtx(ctx, cfg.Roll_type)
		if rollTypeCfg == nil {
			logger.CtxError(ctx, "FixEquipAttr GMazeEquipAffixRandPoolV8Cfg err",
				zap.Int32("rollId", cfg.Roll_type))
			return nil, false, errors.New("配置不存在")
		}
		calcAddCount, calcRangeCount := equip.CalcAttrRealRoll(ctx, cfg.Add_attr_min, cfg.Add_attr_max, cfg.Show_attr_min, cfg.Show_attr_max, int64(attrInfo.GetRandWeight()), rollTypeCfg.Round_value)
		for index := 0; index <= len(attrInfo.ShowAttrList); index++ {
			showAttrInfo := attrInfo.ShowAttrList[index]
			attrValue := equip.CalcAttrValByRoll(cfg.Show_attr_min[showAttrInfo.GetAttrId()], cfg.Show_attr_max[showAttrInfo.GetAttrId()], rollTypeCfg.Round_value, calcAddCount, calcRangeCount)
			if showAttrInfo.GetAttrValue() != attrValue {
				showAttrInfo.AttrValue = proto.Int64(attrValue)
				chgEquip = true
			}
		}
		for index := 0; index <= len(attrInfo.RealAttrList); index++ {
			realAttrInfo := attrInfo.RealAttrList[index]
			attrValue := equip.CalcAttrValByRoll(cfg.Add_attr_min[realAttrInfo.GetAttrId()], cfg.Add_attr_max[realAttrInfo.GetAttrId()], rollTypeCfg.Round_value, calcAddCount, calcRangeCount)
			if realAttrInfo.GetAttrValue() != attrValue {
				realAttrInfo.AttrValue = proto.Int64(attrValue)
				chgEquip = true
			}
		}
	}

	return equipInfo, chgEquip, nil
}

// func resetInsAllEquip(ctx context.Context, userId uint64) error {
// 	equipInfoMap, err := mazebagequipredis.GetAllEquipInfo(logger, userId)
// 	if err != nil {
// 		logger.CtxError(ctx,"fixAllEquipAttrLimit GetAllEquipInfo fail", zap.Error(err))
// 		return err
// 	}
// 	chgEquipList := make([]*MazeEquipCache.MazeEquipInfoDb, 0)
// 	for _, equipInfo := range equipInfoMap {
// 		newEquipInfo, err := resetInsEquip(logger, userId, equipInfo)
// 		if err != nil {
// 			logger.CtxError(ctx,"fixAllEquipAttrLimit fixEquipAttr fail", zap.Error(err))
// 			return err
// 		}
// 		chgEquipList = append(chgEquipList, newEquipInfo)
// 	}
// 	if len(chgEquipList) > 0 {
// 		mazebagequipredis.BatchSaveEquipInfo(logger, userId, chgEquipList)
// 	}
// 	return err
// }

// func resetInsEquip(ctx context.Context, userId uint64, equipInfo *MazeEquipCache.MazeEquipInfoDb) (*MazeEquipCache.MazeEquipInfoDb, error) {
// 	// 装备实例化开始
// 	up := &process.DEIUWParam{}
// 	up.ChargeStage = 1
// 	up.OpType = 0
// 	up.TradeNum = 0
// 	// up.BarrierMap = make(map[int32]int64)
// 	up.PoolLimitMap = process.GetPoolLimitMap(logger, mazeequipconfigv8.GetEquipGoodAttrIds())
// 	up.PoolLimitBase2Map = process.GetPoolLimitMap(logger, mazeequipconfigv8.GetEquipBadAttrIds())
// 	equipSvr := &MazeEquipSvr.SvrEquipInfo{
// 		EquipId: proto.Int32(equipInfo.GetEquipId()),
// 	}
// 	ep := process.NewDEIEWParam(equipSvr)
// 	ep.Guid = equipInfo.GetEquipGuid()
// 	ep.Score = equipInfo.GetEquipScore()
// 	r, e := process.DoInsEquip(logger, userId, up, ep)
// 	if e != nil {
// 		logger.CtxError(ctx,"ResetEquipAttr DoInsEquip fail", zap.Error(e))
// 		return nil, e
// 	}
// 	equipInfo.BaseAttrs = r.EquipInfo.BaseAttrs
// 	equipInfo.ModAttrs = r.EquipInfo.ModAttrs
// 	equipInfo.SpecialAttrs = r.EquipInfo.SpecialAttrs
// 	return equipInfo, nil
// }
