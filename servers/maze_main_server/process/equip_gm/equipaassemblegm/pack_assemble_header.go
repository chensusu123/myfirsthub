/*
 * @Author: majian
 * @Date: 2024-08-28 17:52:09
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-01 20:15:44
 */
package equipaassemblegm

import (
	"bytes"
	"fmt"
	"time"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeEquipCache"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeuserlevelredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecalcattrredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassembleredis"
)

func PackAssembleHeader(logger fklog.FKLogI, userId uint64, as *MazeEquipCache.MazeAssembleDb) (header string, err error) {
	var headerBs bytes.Buffer
	headerBs.WriteString("基本信息:\n")
	dollLv, err := mazeuserlevelredis.GetUserLevel(logger, userId)
	if err != nil {
		return
	}
	headerBs.WriteString(fmt.Sprintf("迷宫等级:%d\n", dollLv))
	force, err := mazecalcattrredis.GetMazeForce(logger, userId)
	if err != nil {
		return
	}
	headerBs.WriteString(fmt.Sprintf("武力值:%d\n", force))
	headerBs.WriteString(EndLine)

	// 初始化状态
	var fields []string
	fields = append(fields, constdef.AssemblePrefixCurAssembleSuitIndex,
		constdef.AssemblePrefixSwitchSuitTime, constdef.AssemblePrefixInitEquip)
	s, e := dollassembleredis.GetDollAssembleMetaInfo(logger, userId, fields...)
	if e == nil {
		initTime := time.Unix(s.GetSwitchTime(), 0).Format("2006-01-02 15:04:05")
		var stateName string = "未初始化"
		if s.GetInitState() == 7 {
			stateName = "已初始化"
		} else if s.GetInitState() == 4 {
			stateName = "初始化失败"
		}
		headerBs.WriteString(fmt.Sprintf("初始化状态:%d(%s) 套装初始化时间:%s\n", s.GetInitState(), stateName, initTime))
		headerBs.WriteString(EndLine)
	}

	// userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	// if err != nil {
	// 	return
	// }
	// var mazeInfo string
	// mazeInfo += "用户解谜数据:\n"
	// mazeInfo += fmt.Sprintf("通关ID:%d\n", userInfo.PassBarrier)
	// mazeInfo += fmt.Sprintf("当前关卡:%d\n", userInfo.Barrier)
	// mazeInfo += fmt.Sprintf("体力:%d\n", userInfo.Energy)
	// mazeInfo += fmt.Sprintf("上次恢复时间:%s\n", time.Unix(userInfo.EnergyLastTime, 0).Format("2006-01-02 15:04:05"))
	// mazeInfo += fmt.Sprintf("装备积分:%d\n", userInfo.EquipPoint)
	// mazeInfo += fmt.Sprintf("总经验:%d\n", userInfo.TotalExp)
	// headerBs.WriteString(mazeInfo)
	// headerBs.WriteString(EndLine)
	header = headerBs.String()
	return
}
