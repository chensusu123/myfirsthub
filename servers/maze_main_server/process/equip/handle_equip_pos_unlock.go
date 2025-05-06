/*
 * @Author: majian
 * @Date: 2024-08-10 11:24:04
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-17 16:15:31
 * @Desc 装备位解锁
 */
package equip

import (
	"fmt"
	"sort"

	"go.uber.org/zap"

	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipPosRankV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/DollEquip"
	"gitlab.ifreetalk.com/plate/protodef/MazeEquipCache"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassembleredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazelevel"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/assembleidpack"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
)

const (
	UnlockSrcInit    string = "登录检查"
	UnlockSrcTask    string = "任务通知"
	UnlockSrcDungoen string = "副本通知"
	UnlockSrcGM      string = "gm检查"
	UnlockSrcDollLv  string = "等级变化"
	UnlockSrcSexChg  string = "性别变化"
)

type UnlockLogic struct {
	Pos       int32
	CondName  string
	NeedCond  string // 需要的条件
	CurCond   string // 当前满足的条件
	UnlockPos *MazeEquipCache.MazeEquipSlotDb
}

func ChkEquipPosUnlock(logger fklog.FKLogI, userId uint64, src string, needNotify bool) error {
	allEquipPos := GMazeEquipPosRankV8Cfg.GetAll()
	equipList, err := dollassembleredis.GetDollEquipPosInfo(logger, userId, len(allEquipPos))
	if err != nil {
		logger.ErrorWF("ChkEquipPosUnlock get pos info fal", zap.Error(err), zap.String("src", src))
		return err
	}
	// 如果存在未解锁的装备
	// var finTaskList map[int32]struct{}
	var mazeLv int64
	if len(allEquipPos) <= len(equipList) {
		logger.InfoWF("ChkEquipPosUnlock already all unlock", zap.Int("len", len(equipList)), zap.String("src", src))
		return nil
	}

	mazeLv, err = mazelevel.GetMazelLevel(logger, userId)
	if err != nil {
		logger.ErrorWF("ChkEquipPosUnlock GetDollLevel fail", zap.Error(err))
		return err
	}

	var unLockLogicList []*UnlockLogic
	var unLockPosList []*MazeEquipCache.MazeEquipSlotDb
	assembleInfo := &MazeEquipCache.MazeAssembleDb{}
	for _, posRow := range allEquipPos {
		// 已解锁
		if equipList[posRow.Pos_id] != nil {
			logger.InfoWF("ChkEquipPosUnlock pos already unlock",
				zap.Int32("pos", posRow.Pos_id),
				zap.String("src", src))
			continue
		}
		unlockCfg := GMazeEquipPosRankV8Cfg.GetMazeEquipPosRankV8Config(posRow.Pos_id)
		if unlockCfg == nil {
			logger.ErrorWF("ChkEquipPosUnlock cannot find cfg",
				zap.Int32("pos", posRow.Pos_id), zap.String("src", src))
			return nil
		}
		// 默认解锁
		if unlockCfg.Is_default == 1 {
			unlockInfo := NewUnlockLogic(unlockCfg)
			unlockInfo.UnlockPos = NewUnlockPosInfo(posRow.Pos_id)
			unLockLogicList = append(unLockLogicList, unlockInfo)
			unLockPosList = append(unLockPosList, unlockInfo.UnlockPos)

			posInfo := &MazeEquipCache.MazeEquipPosInfo{}
			posInfo.EquipPos = unlockInfo.UnlockPos
			assembleInfo.MazeEquips = append(assembleInfo.MazeEquips, posInfo)
			continue
		}
		if unlockCfg.Need_level > 0 {
			unlockInfo := NewUnlockLogic(unlockCfg)
			if mazeLv >= int64(unlockCfg.Need_level) {
				unlockInfo.CurCond = fmt.Sprintf("等级满足 need %d cur %d", unlockCfg.Need_level, mazeLv)
				unlockInfo.UnlockPos = NewUnlockPosInfo(posRow.Pos_id)
				unLockPosList = append(unLockPosList, unlockInfo.UnlockPos)

				posInfo := &MazeEquipCache.MazeEquipPosInfo{}
				posInfo.EquipPos = unlockInfo.UnlockPos
				assembleInfo.MazeEquips = append(assembleInfo.MazeEquips, posInfo)
			} else {
				unlockInfo.CurCond = fmt.Sprintf("等级不满足 need %d cur %d", unlockCfg.Need_level, mazeLv)
			}
			unLockLogicList = append(unLockLogicList, unlockInfo)
			continue
		}
	}
	sort.Slice(unLockLogicList, func(i, j int) bool {
		return unLockLogicList[i].Pos < unLockLogicList[j].Pos
	})

	// 汇总新解锁的装备位
	if len(unLockPosList) == 0 {
		logger.DebugWF("ChkEquipPosUnlock no unlock pos", zap.String("src", src), zap.Any("unLockLogicList", unLockLogicList))
		return nil
	}
	err = dollassembleredis.SetDollEquipPosInfo(logger, userId, unLockPosList)
	if err != nil {
		logger.ErrorWF("ChkEquipPosUnlock unlock fal", zap.Error(err), zap.String("src", src), zap.Any("unLockLogicList", unLockLogicList))
	} else {
		logger.WarnWF("ChkEquipPosUnlock unlock succ", zap.String("src", src), zap.Any("unLockLogicList", unLockLogicList),
			zap.Bool("needNotify", needNotify))
	}
	if needNotify {
		assembleidpack.SendAssembleChgID(logger, userId, assembleInfo,
			int32(DollEquip.ENUM_ASSEMBLE_CHG_TYPE_MASK_ENUM_DOLL_EQUIP_POS_MASK), int32(-1), constdef.DollAssembleChgTypeEquipPosUnlock)
	}

	return err
}

func NewUnlockPosInfo(posId int32) *MazeEquipCache.MazeEquipSlotDb {
	unlockPosInfo := &MazeEquipCache.MazeEquipSlotDb{}
	unlockPosInfo.Pos = proto.Int32(posId)
	unlockPosInfo.Level = proto.Int32(0)
	return unlockPosInfo
}

func NewUnlockLogic(cfg *GMazeEquipPosRankV8Cfg.MazeEquipPosRankV8ConfigRow) *UnlockLogic {
	logic := new(UnlockLogic)
	logic.Pos = cfg.Pos_id
	if cfg.Is_default == 1 {
		logic.CondName = "默认解锁"
		return logic
	}
	if cfg.Need_level > 0 {
		logic.CondName = "等级解锁"
		logic.NeedCond = fmt.Sprintf("需要人物等级:%d", cfg.Need_level)
		return logic
	}

	if cfg.Need_task > 0 {
		logic.CondName = "任务解锁"
		logic.NeedCond = fmt.Sprintf("需要完成任务:%d", cfg.Need_task)
		return logic
	}
	for k, v := range cfg.Need_dungeon_id {
		if k > 0 {
			logic.CondName = "副本解锁"
			logic.NeedCond = fmt.Sprintf("需要通关副本:%d_%d", k, v)
			return logic
		}
	}
	logic.CondName = "未找到解锁条件"
	return logic
}
