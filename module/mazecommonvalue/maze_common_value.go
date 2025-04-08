package mazecommonvalue

import (
	"errors"

	"gitlab.ifreetalk.com/plate/excel/auto/GMazeAttributeFormulaV8Cfg"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/servers/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/servers/maze_game_server/io/redis/mazecalcattrredis"
	"go.uber.org/zap"
)

func GetExtraAdditionForce(logger fklog.FKLogI, userId uint64, userLevel, forceVal int64) (money, equip, exp int64, err error) {
	// allCfg := GMazeKongfuExtraCoinV8Cfg.GetAll()
	// if len(allCfg) == 0 {
	// 	err = errors.New("maze kongfu extra coin cfg empty")
	// 	return
	// }
	// for _, cfg := range allCfg {
	// 	if forceVal <= cfg.Kongfu_max && forceVal >= cfg.Kongfu_min {
	// 		if cfg.Drop_type == 1 {
	// 			money = cfg.Extra_drop_num[int32(userLevel)]
	// 		} else if cfg.Drop_type == 2 {
	// 			equip = cfg.Extra_drop_num[int32(userLevel)]
	// 		} else if cfg.Drop_type == 3 {
	// 			exp = cfg.Extra_drop_num[int32(userLevel)]
	// 		}
	// 	}
	// }

	return
}

func GetMoneyExtraAdditionEquip(logger fklog.FKLogI, userId uint64) (money int64, err error) {

	attrMap, err := mazecalcattrredis.BatchGetMazeCalcAttr(logger, userId, []int32{constdef.MazeMoneyBuff10258})
	if err != nil {
		logger.ErrorWF("GetMoneyExtraAdditionEquip BatchGetDollCalcAttr fail", zap.Error(err), zap.Any("attrId", constdef.MazeMoneyBuff10258))
		return
	}

	// pb := attrMap[constdef.MazeMoneyBuff10258]
	// money = pb.GetAttrVal()
	money = attrMap[constdef.MazeMoneyBuff10258]
	logger.InfoWF("GetMoneyExtraAdditionEquip succ ", zap.Any("money", money), zap.Any("attrMap", attrMap))
	return
}

func GetExpExtraAdditionEquip(logger fklog.FKLogI, userId uint64) (money int64, err error) {

	attrMap, err := mazecalcattrredis.BatchGetMazeCalcAttr(logger, userId, []int32{constdef.MazeExpBuff10261})
	if err != nil {
		logger.ErrorWF("GetExpExtraAdditionEquip BatchGetDollCalcAttr fail", zap.Error(err), zap.Any("attrId", constdef.MazeExpBuff10261))
		return
	}

	// pb := attrMap[constdef.MazeExpBuff10261]
	// money = pb.GetAttrVal()
	money = attrMap[constdef.MazeExpBuff10261]
	logger.InfoWF("GetExpExtraAdditionEquip succ ", zap.Any("money", money), zap.Any("attrMap", attrMap))
	return
}

func GetCalRet(logger fklog.FKLogI, userId uint64, userLevel, forceVal int64, foeExpBase, foeMoneyBase, foeEquipBase int64) (totalExp, totalMoney int64, totalEquip int64, err error) {

	attrs := make([]int32, 0)

	expFormula := GMazeAttributeFormulaV8Cfg.Get(constdef.MazeExp)
	if expFormula == nil {
		err = errors.New("maze exp formula cfg empty")
		return
	}

	attrs = append(attrs, expFormula.Parameter_1, expFormula.Parameter_2, expFormula.Parameter_3, expFormula.Parameter_4, expFormula.Parameter_5, expFormula.Parameter_11, expFormula.Parameter_12)
	attrs = append(attrs, expFormula.Parameter_8...)
	attrs = append(attrs, expFormula.Parameter_9...)
	attrs = append(attrs, expFormula.Parameter_10...)

	moneyFormula := GMazeAttributeFormulaV8Cfg.Get(constdef.MazeMoney)
	if moneyFormula == nil {
		err = errors.New("maze money formula cfg empty")
		return
	}

	attrs = append(attrs, moneyFormula.Parameter_1, moneyFormula.Parameter_2, moneyFormula.Parameter_3, moneyFormula.Parameter_4, moneyFormula.Parameter_5, moneyFormula.Parameter_11, moneyFormula.Parameter_12)
	attrs = append(attrs, moneyFormula.Parameter_8...)
	attrs = append(attrs, moneyFormula.Parameter_9...)
	attrs = append(attrs, moneyFormula.Parameter_10...)

	logger.InfoWF("GetCalRet dump need attrIds", zap.Any("attrs", attrs))

	// attrs = append(attrs, constdef.MazeForce)

	attrDbs, err := mazecalcattrredis.BatchGetMazeCalcAttr(logger, userId, attrs)
	if err != nil {
		logger.ErrorWF("GetCalRet BatchGetDollCalcAttr fail", zap.Error(err), zap.Any("attrs", attrs))
		return
	}

	money, equip, exp, err := GetExtraAdditionForce(logger, userId, userLevel, forceVal)
	if err != nil {
		logger.ErrorWF("GetCalRet GetExtraAdditionForce fail", zap.Error(err), zap.Any("userLevel", userLevel), zap.Any("forceVal", forceVal))
		return
	}

	attrMap := make(map[int32]int64)
	for k, attr := range attrDbs {
		if k == constdef.MazeMoneyBuff10258 {
			attrMap[k] = attr + money + foeMoneyBase
		} else if k == constdef.MazeExpBuff10261 {
			attrMap[k] = attr + exp + foeExpBase
		} else {
			attrMap[k] = attr
		}
	}

	totalExp = calRet(logger, attrMap, expFormula)
	totalMoney = calRet(logger, attrMap, moneyFormula)
	totalEquip = foeEquipBase + equip
	logger.InfoWF("GetCalRet calRet dump", zap.Any("totalExp", totalExp), zap.Any("totalMoney", totalMoney), zap.Any("totalEquip", totalEquip))
	return
}

func calRet(logger fklog.FKLogI, attrMap map[int32]int64, cfg *GMazeAttributeFormulaV8Cfg.MazeAttributeFormulaV8ConfigRow) (ret int64) {
	p1 := attrMap[cfg.Parameter_1]
	p2 := attrMap[cfg.Parameter_2]
	p3 := attrMap[cfg.Parameter_3]
	p4 := attrMap[cfg.Parameter_4]
	p5 := attrMap[cfg.Parameter_5]
	p8 := convertAttr(attrMap, cfg.Parameter_8)
	p9 := convertAttr(attrMap, cfg.Parameter_9)
	p10 := convertAttr(attrMap, cfg.Parameter_10)
	p11 := attrMap[cfg.Parameter_11]
	p12 := attrMap[cfg.Parameter_12]

	//最终展示武力值 = {[参数1 * （1+参数2/10000) * （1+参数8/10000） * （1+参数3/10000) + 参数4]* (1+参数5/10000)*(1+参数9/10000)*(1+参数10/10000) +参数11}*（1+参数12/10000）
	ret1 := float64(p1) * (1 + float64(p2)/10000.0) * p8 * (1 + float64(p3)/10000.0)
	ret2 := (ret1 + float64(p4)) * (1 + float64(p5)/10000.0) * p9 * p10
	ret = int64((ret2 + float64(p11)) * (1 + float64(p12)/10000.0))
	logger.InfoWF("calRet dump", zap.Any("p1", p1), zap.Any("p2", p2), zap.Any("p3", p3), zap.Any("p4", p4), zap.Any("p5", p5),
		zap.Any("p8", p8), zap.Any("p9", p9), zap.Any("p10", p10), zap.Any("p11", p11), zap.Any("p12", p12), zap.Any("ret", ret))
	return
}

func convertAttr(attrMap map[int32]int64, attrIds []int32) (ret float64) {
	ret = 1.0
	for _, attrId := range attrIds {
		if attrId == 0 {
			continue
		}
		ret = ret * (1 + float64(attrMap[attrId])/10000.0)
	}
	return
}
