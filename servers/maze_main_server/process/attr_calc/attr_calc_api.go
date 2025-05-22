/*
 * @Author: majian
 * @Date: 2024-07-10 14:01:38
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-01 14:50:14
 */
package attr_calc

import (
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeBuffData"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/dollattr"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/vardef"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeAttributeV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeEquipConfigV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/attr_calc/commonlogic"
	"go.uber.org/zap"
)

func AddtionMazeAttr(in map[int32]int64, attr *MazeBuffData.MazeBuffAttr) {
	in[attr.GetAttrId()] += attr.GetAttrVal()
}

func AddtionMazeAttrKv(in map[int32]int64, attrId int32, val int64) {
	in[attrId] += val
}

func SetMazeAttrKv(in map[int32]int64, attrId int32, val int64) {
	in[attrId] = val
}

// 判断是否是人偶计算属性
func IsDollCalcAttr(attrId int32) bool {
	row := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(attrId)
	if row != nil {
		if dollattr.IsDollCalcAttr(row.Type) {
			return true
		}
	}
	return false
}

func DumpAttrBySrc(logger fklog.FKLogI, src int, k int32, v int64, extra string) {
	logger.InfoWF("DumpAttrBySrc", zap.Int("src", src), zap.String("srcName", GetSrcName(src)),
		zap.Int32("k", k), zap.Int64("v", v), zap.String("extra", extra))
}

// 抗性属性转换
func ResistanceAttrConvert(logger fklog.FKLogI, attrsMap map[int32]int64) {
	rowCfg := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(constdef.DollResistanceConvert)
	if rowCfg == nil {
		logger.InfoWF("ResistanceAttrConvert no cfg")
		return
	}
	var hasConvert bool
	var logInfo string
	if rowCfg.Value_int > 0 && len(rowCfg.Value_map) > 0 {
		allKangXingVal := attrsMap[int32(rowCfg.Value_int)]
		if allKangXingVal > 0 {
			for k, v := range rowCfg.Value_map {
				if k <= 0 || v <= 0 {
					continue
				}
				convertVal := allKangXingVal * v / 10000
				oldVal := attrsMap[k]
				attrsMap[k] += convertVal
				newVal := attrsMap[k]
				hasConvert = true
				logInfo += fmt.Sprintf("attrId:%d old:%d new:%d rate=%d,", k, oldVal, newVal, v)
			}
		}
	}
	if hasConvert {
		logger.InfoWF("ResistanceAttrConvert", zap.String("result", logInfo))
	} else {
		logger.InfoWF("ResistanceAttrConvert no chg")
	}
}

func GetSrcName(src int) string {
	if desc, ok := vardef.MazeBuffSrcDescMap[int32(src)]; ok {
		return desc
	}
	if desc, ok := commonlogic.DollAttrSrcMapCfg.Load(int32(src)); ok {
		return desc.(string)
	}
	return ""
}
