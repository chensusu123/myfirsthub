/*
 * @Author: majian
 * @Date: 2024-08-27 19:38:34
 * @Last Modified by: majian
 * @Last Modified time: 2024-09-13 14:39:29
 */
package equipaassemblegm

import (
	"bytes"
	"context"
	"fmt"

	"maze_game_server/config/GMazeEquipPosLvSuiteV8Cfg"
	"maze_game_server/module/calcassembleattr"
	"maze_game_server/pb/server/MazeEquipCache"
)

func PackEquipSuitInfo(ctx context.Context, userId uint64, assembleInfo *MazeEquipCache.MazeAssembleDb, effectInfo *calcassembleattr.EquipmentEffectInfo) (s string, e error) {
	// 打包装备套装属性
	var suitString bytes.Buffer
	suitString.WriteString("装备套装信息:\n")
	if effectInfo.SuitCalc != nil {
		suitMap := effectInfo.SuitCalc.GetSuitNumMap()

		for k, v := range suitMap {
			if k > 0 && v > 0 {
				row := calcassembleattr.GetLeastRow(k, v)
				//	row := GMazeEquipSuiteAttrV8Cfg.GetDollEquipSuiteAttrV8Config(calcassembleattr.LegendSuitKey(k, v))
				if row == nil {
					continue
				}
				suitString.WriteString(fmt.Sprintf("套装Id:%d(%s) 套装数量:%d 激活数量:%d\t", k, row.Add_attr_desc, row.Affix_num, v))

				for id, val := range row.Add_attr {
					suitString.WriteString(fmt.Sprintf("属性Id(真):%d 属性值:%d|", id, val))
				}
				suitString.WriteString("\t")
				for id, val := range row.Add_attr_show {
					suitString.WriteString(fmt.Sprintf("属性Id(显):%d 属性值:%d|", id, val))
				}
				suitString.WriteString("\n")
			}
		}
	}
	enSuitId := assembleInfo.GetEpEnSuitId()
	suitString.WriteString("装备位强化套装信息:\n")
	if enSuitId > 0 {
		suitString.WriteString(fmt.Sprintf("套装Id:%d\t", enSuitId))
		enSuitCfg := GMazeEquipPosLvSuiteV8Cfg.GetWithCtx(ctx, enSuitId)
		if enSuitCfg != nil {
			for id, val := range enSuitCfg.Add_attr {
				suitString.WriteString(fmt.Sprintf("属性Id(真):%d 属性值:%d|", id, val))
			}

			for id, val := range enSuitCfg.Show_attr {
				suitString.WriteString(fmt.Sprintf("属性Id(显):%d 属性值:%d|", id, val))
			}
			suitString.WriteString("\n")
		}
	}
	return suitString.String(), nil
}
