package equip

import (
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/config/GMazeEquipPosLvV8Cfg"
	"maze_game_server/excel/equipposexcel"
	"maze_game_server/io/redis/dollassembleredis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/equippossuit"
	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeEquipPos"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 装备位强化预览
func (e *Equip) OnEquipPosLvUpPreviewRQ_10423_10424(s *session.Session, rq *MazeEquipPos.MazeEquipPosLvUpPreviewRQ) (err error) {
	defer fkprometheus.DebugPMT("OnEquipPosLvUpPreviewRQ")()

	logger := log.Clone("Equip", uint64(s.UID()), 0)
	rs := &MazeEquipPos.MazeEquipPosLvUpPreviewRS{}

	rs.ErrInfo = errors.NO_ERROR
	rs.Header = rq.Header
	rs.PosId = rq.PosId
	posId := rq.GetPosId()

	userId := uint64(s.UID())

	logger.InfoWF("OnEquipPosLvUpPreviewRQ start", zap.Any("rq", rq))
	defer func() {
		err = s.Response(rs)
		logger.InfoWF("OnEquipPosLvUpPreviewRQ end", zap.Any("rs", rs))
	}()

	if posId < 1 || posId > constdef.EquipPosNum {
		logger.ErrorWF("OnEquipPosLvUpPreviewRQ param invalid", zap.Int32("posId", posId))
		rs.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
		return
	}

	assembleDb, err := dollassembleredis.GetAllAssembleInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("OnEquipPosLvUpPreviewRQ GetAllAssembleInfo", zap.Error(err))
		rs.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	equipPos := assembleDb.GetMazeEquips()
	if len(equipPos) == 0 {
		logger.WarnWF("OnEquipPosLvUpPreviewRQ equip pos num is 0")
		rs.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无装备位数据")
		return
	}

	equipPosSuitId := assembleDb.GetEpEnSuitId()
	rs.CurSuit, rs.NextSuit, err = equippossuit.GetCurAndNextSuit(logger, equipPosSuitId)
	if err != nil {
		logger.ErrorWF("OnEquipPosLvUpPreviewRQ GetCurAndNextSuit", zap.Error(err))
		rs.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	var curLv int32
	var curEquipPosDb *MazeEquipCache.MazeEquipSlotDb
	for _, v := range equipPos {
		if posId == v.GetEquipPos().GetPos() {
			curLv = v.GetEquipPos().GetLevel()
			curEquipPosDb = v.GetEquipPos()
		}
	}

	if curEquipPosDb == nil {
		logger.ErrorWF("OnEquipPosLvUpPreviewRQ param invalid, no current equip position", zap.Int32("posId", posId))
		rs.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
		return
	}
	curCfg := equipposexcel.GetPosStrengthCfg(posId, curLv)
	if curCfg == nil {
		logger.ErrorWF("OnEquipPosLvUpPreviewRQ no found strength cfg",
			zap.Int32("posId", posId), zap.Int32("lv", curLv))
		rs.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		return
	}
	var nextCfg *GMazeEquipPosLvV8Cfg.MazeEquipPosLvV8ConfigRow
	if curCfg.Next_order > 0 {
		nextCfg = equipposexcel.GetPosStrengthCfgByKey(curCfg.Next_order)
	}

	rs.StLevel = &Common.AttrChgInfo{
		AttrVal: proto.Int64(int64(curLv)),
	}
	if nextCfg != nil {
		rs.StLevel.NextVal = proto.Int64(int64(nextCfg.Level))
	}

	// 对比属性加成
	for _, id := range curCfg.Show_attr_order {
		if id <= 0 {
			continue
		}
		attrInfo := &Common.AttrChgInfo{
			AttrId:  proto.Int32(id),
			AttrVal: proto.Int64(curCfg.Show_attr[id]),
		}
		if nextCfg != nil {
			attrInfo.NextVal = proto.Int64(nextCfg.Show_attr[id])
		}
		rs.AttrBuffs = append(rs.AttrBuffs, attrInfo)
	}

	rs.StCond = &MazeGameEquip.EquipPosStCond{
		CostItems:   itemutil.Map2Common(curCfg.Cost),
		NeedOtherLv: proto.Int32(curCfg.Need_other),
		NeedMazeLv:  proto.Int32(curCfg.Need_maze_level),
	}

	return
}
