package dollassembleinfo

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/assemble"
	"maze_game_server/config/GMazeEquipPosRankV8Cfg"
	"maze_game_server/io/redis/dollassembleredis"
	"maze_game_server/io/redis/dollassemblesuitredis"
	"maze_game_server/module/calcassembleattr"
	"maze_game_server/module/effectequip"
	"maze_game_server/pb/server/MazeEquipCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 打包装配信息
func GetDollAssembleInfo(ctx context.Context, userId uint64) (assembleInfo *MazeEquipCache.MazeAssembleDb, err error) {
	assembleInfo, _, err = GetDollAssembleInfoEx(ctx, userId)
	return
}

func GetDollAssembleInfoEx(ctx context.Context, userId uint64) (assembleInfo *MazeEquipCache.MazeAssembleDb, effectInfo *calcassembleattr.EquipmentEffectInfo, err error) {
	return GetDollAssembleInfoV2(ctx, userId, nil)
}

func CheckIsAssemble(ctx context.Context, userId uint64, equipGuid int64, equipPos int32) (isAssemble bool, err error) {
	logger := fklog.ContextAppLogger(ctx)
	dollAssembleMetaSt, err := dollassembleredis.GetDollAssembleMetaInfo(ctx, userId, constdef.AssemblePrefixCurAssembleSuitIndex)
	if err != nil {
		logger.CtxError(ctx, "CheckIsAssemble GetDollAssembleMetaInfo fail", zap.Error(err), zap.Any("equipPos", equipPos))
		return
	}
	equipPosDb, err := dollassemblesuitredis.GetDollAssembleByPos(ctx, userId, dollAssembleMetaSt.GetCurSuitIndex(), equipPos)
	if err != nil {
		logger.CtxError(ctx, "CheckIsAssemble GetDollAssembleByPos fail", zap.Error(err),
			zap.Any("equipPos", equipPos), zap.Any("suitIndex", dollAssembleMetaSt.GetCurSuitIndex()))
		return
	}

	return equipPosDb.GetEquipGuid() == int64(equipGuid), nil
}

type AssParam struct {
	ReplaeEquipDb *MazeEquipCache.MazeEquipInfoDb
}

// 装备信息加上每个pos的装备信息和装备效果
func GetDollAssembleInfoV2(ctx context.Context, userId uint64, p *AssParam) (assembleInfo *MazeEquipCache.MazeAssembleDb, effectInfo *calcassembleattr.EquipmentEffectInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	// 查询装配信息
	assembleInfo, err = dollassembleredis.GetAllAssembleInfo(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "GetDollAssembleInfo get fail", zap.Error(err))
		return
	}
	effectInfo = calcassembleattr.NewEquipmentEffectInfo(ctx)

	allEquipPos := GMazeEquipPosRankV8Cfg.GetAll()
	posNum := len(allEquipPos)
	var equips []*MazeEquipCache.MazeEquipPosDb
	suitIndex := constdef.CurSuitDef
	if assembleInfo.GetCurSuitIndex() > 0 {
		suitIndex = assembleInfo.GetCurSuitIndex()
	}
	//获取每个位置的装备信息
	if suitIndex > 0 {
		equips, err = dollassemblesuitredis.GetDollAssembleSuit(ctx, userId,
			assembleInfo.GetCurSuitIndex(), posNum)
		if err != nil {
			logger.CtxError(ctx, "GetDollAssembleInfo GetDollAssembleSuit fail", zap.Error(err),
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
	realEquips, err := effectequip.BatchGetEffectEquipInfo(ctx, userId, equipGuids...)
	if err != nil {
		logger.CtxError(ctx, "GetDollAssembleInfo GetBatchEquipInfo fail", zap.Error(err),
			zap.Any("equipGuids", equipGuids))
		return
	}

	if p != nil && p.ReplaeEquipDb != nil {
		if _, ok := realEquips[p.ReplaeEquipDb.GetEquipGuid()]; ok {
			realEquips[p.ReplaeEquipDb.GetEquipGuid()] = p.ReplaeEquipDb
		} else {
			logger.CtxWarn(ctx, "GetDollAssembleInfo no find old equip",
				zap.Any("oldEquips", p.ReplaeEquipDb))
		}
	}

	for _, assemblePos := range assembleInfo.MazeEquips {
		if !assemble.IsAssembleEquip(assemblePos) {
			continue
		}
		equipDetail := realEquips[assemblePos.GetEquipLoadInfo().GetEquipGuid()]
		if equipDetail == nil || equipDetail.GetEquipGuid() == 0 {
			logger.CtxError(ctx, "GetDollAssembleInfo equip info no exist",
				zap.Int64("guid", assemblePos.GetEquipLoadInfo().GetEquipGuid()))
			err = errors.New("dressed equip no exist")
			return
		}
		assemblePos.EquipInfo = equipDetail
	}

	cond := calcassembleattr.EffectCalcInParam{IsLog: true, IsForce: false}
	// cond.IsFiveOnlyRead = true
	effectInfo, err = calcassembleattr.CalcEquipEffectAll(ctx, assembleInfo.MazeEquips, cond)
	if err != nil {
		return
	}

	// 计算套装
	if effectInfo.SuitId > 0 {
		assembleInfo.EpSuitId = proto.Int32(effectInfo.SuitId)
	}
	return
}
