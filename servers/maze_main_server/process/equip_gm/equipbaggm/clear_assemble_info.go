package equipbaggm

import (
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassemblesuitredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/assemble"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassembleredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebuffinforedis"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structsdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeattrcalcnotifyqueue"
)

// 删除人偶装配信息
func ClearDollAssembleInfo(logger fklog.FKLogI, userId uint64) error {
	var step int
	var e error
	defer func() {
		logger.InfoWF("ClearDollAssembleInfo result", zap.Error(e),
			zap.Int("step", step), zap.Uint64("uid", userId))
	}()

	e = dollassemblesuitredis.DelEquipSuitInfo(logger, userId)
	if e != nil {
		return e
	}
	step = 1
	var delList []string
	delList = append(delList, constdef.AssemblePrefixInitEquip,
		constdef.AssemblePrefixCurAssembleSuitIndex, constdef.AssemblePrefixSwitchSuitTime)
	delList = append(delList, constdef.AssemblePrefixEquipPosEnSuit)
	for i := 1; i <= constdef.EquipPosNum; i++ {
		delList = append(delList, assemble.EnCodeAssemblePosField(int32(i)))
	}
	e = dollassembleredis.BatchDelAssmebleInfo(logger, userId, delList...)
	if e != nil {
		return e
	}
	step = 2

	e = mazebuffinforedis.DelMazeBuffBySrc(logger, userId, constdef.MazeBuffSrcEquip)
	if e == nil {
		step = 4
		calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{}
		calcAttrNotify.FromServer = "maze_equip_gm_server"
		calcAttrNotify.UserId = userId
		calcAttrNotify.ChgType = constdef.MazeBuffChgTypeEquipGm
		calcAttrNotify.Session = "gm"
		calcAttrNotify.BuffSrc = constdef.MazeBuffSrcEquip
		mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, calcAttrNotify)
	}

	e = mazebuffinforedis.DelMazeBuffBySrc(logger, userId, constdef.MazeBuffSrcEquipPos)
	if e == nil {
		step = 5
		calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{}
		calcAttrNotify.FromServer = "maze_equip_gm_server"
		calcAttrNotify.UserId = userId
		calcAttrNotify.ChgType = constdef.MazeBuffChgTypeEquipGm
		calcAttrNotify.Session = "gm"
		calcAttrNotify.BuffSrc = constdef.MazeBuffSrcEquipPos
		mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, calcAttrNotify)
	}
	return e
}
