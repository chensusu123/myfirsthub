/*
 * @Author: majian
 * @Date: 2025-02-19 21:20:43
 * @Last Modified by: majian
 * @Last Modified time: 2025-02-19 21:57:31
 * @Desc 迷宫怪物伤害计算
 */
package mazehurtcalc

import (
	"errors"

	"gitlab.ifreetalk.com/plate/excel/auto/GFightKongfuMazeV8Cfg"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type MazeAttrInfo struct {
	Force int64 // 武力
	Hp    int64 // 血量
}

type MazeLoseHurt struct {
	PlayerLoseHurt int64 // 玩家每回合损血
	FoeLoseHurt    int64 // 怪物每回合损血
}

func CalcLoseHurt(logger fklog.FKLogI, loseHpType int32, role, monster *MazeAttrInfo) (lose *MazeLoseHurt, err error) {
	forceRate := role.Force * 10000 / monster.Force
	cfg := GetMazeLoseHurtCfg(forceRate, loseHpType)
	if cfg == nil {
		logger.ErrorWF("CalcLoseHurt cannot find cfg",
			zap.Int64("roleForce", role.Force),
			zap.Int64("MonsterForce", monster.Force),
			zap.Int32("loseType", loseHpType),
			zap.String("cfg", GFightKongfuMazeV8Cfg.GetConfigDesc()))
		return nil, errors.New("no found cfg")
	}
	lose = &MazeLoseHurt{}
	lose.FoeLoseHurt = LoseHurtFormula(monster.Hp, cfg.Foe_hp_lose, cfg.Round)
	lose.PlayerLoseHurt = LoseHurtFormula(role.Hp, cfg.Player_hp_lose, cfg.Round)
	logger.InfoWF("CalcLoseHurt end", zap.Any("playerParam", role),
		zap.Any("foeParam", monster), zap.Int32("loseType", loseHpType),
		zap.Any("ret", lose))
	return lose, nil
}

func LoseHurtFormula(hp int64, hpLostRate, round int32) int64 {
	return int64(hpLostRate) * hp * 10000 / (int64(round) * 1000000)
}

func GetMazeLoseHurtCfg(forceRate int64, loseType int32) *GFightKongfuMazeV8Cfg.FightKongfuMazeV8ConfigRow {
	allRows := GFightKongfuMazeV8Cfg.GetAllFightKongfuMazeV8Config()
	for _, row := range allRows {
		if row.Hp_lose_type == loseType && (forceRate >= row.Pro__min && forceRate <= row.Pro_max) {
			return row
		}
	}
	return nil
}
