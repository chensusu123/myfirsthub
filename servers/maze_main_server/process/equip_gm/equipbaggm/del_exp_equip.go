/*
 * @Author: majian
 * @Date: 2024-08-29 15:37:44
 * @Last Modified by: majian
 * @Last Modified time: 2024-08-29 16:23:51
 */
package equipbaggm

// func GmDelExpModAttrBag(ctx context.Context, userId uint64, run bool) (need bool, err error) {
// 	allEquips, err := mazebagequipredis.GetAllEquipInfo(logger, userId)
// 	if err != nil {
// 		return false, err
// 	}

// 	var needDelEquips []int64
// 	for _, equip := range allEquips {
// 		var needDel bool
// 		mods := equip.GetModAttrs()
// 		for _, mod := range mods {
// 			for _, real := range mod.GetRealAttrList() {
// 				if real.GetAttrId() == 0 {
// 					needDel = true
// 					break
// 				}
// 			}
// 			for _, real := range mod.GetShowAttrList() {
// 				if real.GetAttrId() == 0 {
// 					needDel = true
// 					break
// 				}
// 			}
// 		}
// 		if needDel {
// 			needDelEquips = append(needDelEquips, equip.GetEquipGuid())
// 		}
// 	}
// 	if len(needDelEquips) > 0 {
// 		need = true
// 		equips, err := dollassemblesuitredis.GetDollAssembleSuit(logger, userId, 1, 8)
// 		if err != nil {
// 			return false, err
// 		}
// 		var asEquips []*MazeEquipCache.MazeEquipPosDb
// 		for _, equip := range equips {
// 			for _, bagEqiup := range needDelEquips {
// 				if bagEqiup == equip.GetEquipGuid() {
// 					equip.EquipGuid = proto.Int64(0)
// 					equip.EquipId = proto.Int32(0)
// 					asEquips = append(asEquips, equip)
// 					break
// 				}
// 			}
// 		}
// 		if len(asEquips) > 0 {
// 			if run {
// 				dollassemblesuitredis.SetDollAssembleSuit(logger, userId, 1, asEquips)
// 			}
// 			logger.CtxInfo(ctx,"GmDelExpModAttrBag del as", zap.Any("asEquips", asEquips))
// 		}
// 		if run {
// 			mazebagequipredis.BatchDelEquip(logger, userId, needDelEquips...)
// 		}

// 		logger.CtxInfo(ctx,"GmDelExpModAttrBag del", zap.Any("bagEquips", needDelEquips))
// 	}
// 	return need, nil
// }
