package dollassembleinfo

import (
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipPosRankV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassembleredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/calcassembleattr"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassemblesuitredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/assemble"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/effectequip"
)

// 打包装配信息
func GetDollAssembleInfo(logger fklog.FKLogI, userId uint64) (assembleInfo *MazeEquipCache.MazeAssembleDb, err error) {
	assembleInfo, _, err = GetDollAssembleInfoEx(logger, userId)
	return
}

func GetDollAssembleInfoEx(logger fklog.FKLogI, userId uint64) (assembleInfo *MazeEquipCache.MazeAssembleDb, effectInfo *calcassembleattr.EquipmentEffectInfo, err error) {
	return GetDollAssembleInfoV2(logger, userId, nil)
}

func CheckIsAssemble(logger fklog.FKLogI, userId uint64, equipGuid int64, equipPos int32) (isAssemble bool, err error) {
	dollAssembleMetaSt, err := dollassembleredis.GetDollAssembleMetaInfo(logger, userId, constdef.AssemblePrefixCurAssembleSuitIndex)
	if err != nil {
		logger.ErrorWF("CheckIsAssemble GetDollAssembleMetaInfo fail", zap.Error(err), zap.Any("equipPos", equipPos))
		return
	}
	equipPosDb, err := dollassemblesuitredis.GetDollAssembleByPos(logger, userId, dollAssembleMetaSt.GetCurSuitIndex(), equipPos)
	if err != nil {
		logger.ErrorWF("CheckIsAssemble GetDollAssembleByPos fail", zap.Error(err),
			zap.Any("equipPos", equipPos), zap.Any("suitIndex", dollAssembleMetaSt.GetCurSuitIndex()))
		return
	}

	return equipPosDb.GetEquipGuid() == int64(equipGuid), nil
}

type AssParam struct {
	ReplaeEquipDb *MazeEquipCache.MazeEquipInfoDb
}

func GetDollAssembleInfoV2(logger fklog.FKLogI, userId uint64, p *AssParam) (assembleInfo *MazeEquipCache.MazeAssembleDb, effectInfo *calcassembleattr.EquipmentEffectInfo, err error) {
	// 查询装配信息
	assembleInfo, err = dollassembleredis.GetAllAssembleInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("GetDollAssembleInfo get fail", zap.Error(err))
		return
	}
	effectInfo = calcassembleattr.NewEquipmentEffectInfo(logger)

	allEquipPos := GMazeEquipPosRankV8Cfg.GetAll()
	posNum := len(allEquipPos)
	var equips []*MazeEquipCache.MazeEquipPosDb
	suitIndex := constdef.CurSuitDef
	if assembleInfo.GetCurSuitIndex() > 0 {
		suitIndex = assembleInfo.GetCurSuitIndex()
	}
	if suitIndex > 0 {
		equips, err = dollassemblesuitredis.GetDollAssembleSuit(logger, userId,
			assembleInfo.GetCurSuitIndex(), posNum)
		if err != nil {
			logger.ErrorWF("GetDollAssembleInfo GetDollAssembleSuit fail", zap.Error(err),
				zap.Int32("curSuit", assembleInfo.GetCurSuitIndex()))
			return
		}
		// 填充装配信息
		for _, pos := range assembleInfo.GetMazeEquips() {
			if pos.GetEquipPos() == nil {
				continue
			}
			for _, loadInfo := range equips {
				if loadInfo.GetPos() == pos.GetEquipPos().GetPos() {
					pos.EquipLoadInfo = loadInfo
					break
				}
			}
		}
	}

	equipGuids := make([]int64, len(assembleInfo.MazeEquips))
	for _, assemblePos := range assembleInfo.MazeEquips {
		if !assemble.IsAssembleEquip(assemblePos) {
			continue
		}
		equipGuids = append(equipGuids, assemblePos.GetEquipLoadInfo().GetEquipGuid())
	}
	if len(equipGuids) == 0 {
		return
	}
	// 从背包查询装备信息
	realEquips, err := effectequip.BatchGetEffectEquipInfo(logger, userId, equipGuids...)
	if err != nil {
		logger.ErrorWF("GetDollAssembleInfo GetBatchEquipInfo fail", zap.Error(err),
			zap.Any("equipGuids", equipGuids))
		return
	}

	if p != nil && p.ReplaeEquipDb != nil {
		if _, ok := realEquips[p.ReplaeEquipDb.GetEquipGuid()]; ok {
			realEquips[p.ReplaeEquipDb.GetEquipGuid()] = p.ReplaeEquipDb
		} else {
			logger.WarnWF("GetDollAssembleInfo no find old equip",
				zap.Any("oldEquips", p.ReplaeEquipDb))
		}
	}

	for _, assemblePos := range assembleInfo.MazeEquips {
		if !assemble.IsAssembleEquip(assemblePos) {
			continue
		}
		equipDetail := realEquips[assemblePos.GetEquipLoadInfo().GetEquipGuid()]
		if equipDetail == nil || equipDetail.GetEquipGuid() == 0 {
			logger.ErrorWF("GetDollAssembleInfo equip info no exist",
				zap.Int64("guid", assemblePos.GetEquipLoadInfo().GetEquipGuid()))
			err = errors.New("dressed equip no exist")
			return
		}
		assemblePos.EquipInfo = equipDetail
	}

	cond := calcassembleattr.EffectCalcInParam{IsLog: true, IsForce: false}
	// cond.IsFiveOnlyRead = true
	effectInfo, err = calcassembleattr.CalcEquipEffectAll(logger, assembleInfo.MazeEquips, cond)
	if err != nil {
		return
	}

	// 计算套装
	if effectInfo.SuitId > 0 {
		assembleInfo.EpSuitId = proto.Int32(effectInfo.SuitId)
	}
	return
}
