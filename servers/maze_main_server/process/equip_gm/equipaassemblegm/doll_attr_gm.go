/*
 * @Author: majian
 * @Date: 2024-07-15 16:40:12
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-01 20:05:41
 */
package equipaassemblegm

import (
	"bytes"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/vardef"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeAttrSpDescV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeAttributeV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/config/GMazeInitialAttrV8Cfg"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebuffinforedis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecalcattrredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeattrformula"
)

func DumpDollCalcAttr(logger fklog.FKLogI, userId uint64) (attrInfo string, err error) {
	attrs, err := mazecalcattrredis.HScanMazeCalcAttr(logger, userId)
	if err != nil {
		return
	}
	var bs bytes.Buffer
	showAttrs := make(map[int32]int64)
	realAttrs := make(map[int32]int64)

	for k, attr := range attrs {
		attrRow := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(k)
		if attrRow != nil {
			if attrRow.Type == constdef.DollAttrTypeShow {
				showAttrs[k] = attr
			} else {
				realAttrs[k] = attr
			}
		}
	}

	bs.WriteString("\n展示属性:\n")
	for k, v := range showAttrs {
		var attrName string
		attrRow := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(k)
		if attrRow != nil {
			attrName = attrRow.Name
		}
		spRow := GMazeAttrSpDescV8Cfg.GetMazeAttrSpDescV8Config(k)
		if spRow != nil {
			attrName = spRow.Attr_list_affix_desc
		}

		key := fmt.Sprintf("%s(%d):\t", attrName, k)
		bs.WriteString(fmt.Sprintf("%s%-10d\n", key, v))
	}

	bs.WriteString("\n生效属性:\n")
	for k, v := range realAttrs {
		var attrName string
		attrRow := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(k)
		if attrRow != nil {
			attrName = attrRow.Name
		}
		spRow := GMazeAttrSpDescV8Cfg.GetMazeAttrSpDescV8Config(k)
		if spRow != nil {
			attrName = spRow.Attr_list_affix_desc
		}
		if constdef.DollFormulaBlood == k {
			attrName = "气血最终值"
		}
		if constdef.DollFormulaAttack == k {
			attrName = "攻击最终值"
		}
		if constdef.DollFormulaDefend == k {
			attrName = "防御最终值"
		}
		key := fmt.Sprintf("%s(%d):\t", attrName, k)
		bs.WriteString(fmt.Sprintf("%s%-10d\n", key, v))
	}

	bs.WriteString("\n初始属性:\n")
	bs.WriteString(fmt.Sprintf("配表:%s\n", GMazeInitialAttrV8Cfg.SheetName()))
	row := GMazeInitialAttrV8Cfg.GetMazeInitialAttrV8Config(constdef.DollInitAttrCfgId)
	if row != nil {
		for k, v := range row.Initial_attr {
			var attrName string
			attrRow := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(k)
			if attrRow != nil {
				attrName = attrRow.Name
			}

			spRow := GMazeAttrSpDescV8Cfg.GetMazeAttrSpDescV8Config(k)
			if spRow != nil {
				attrName = spRow.Attr_list_affix_desc
			}
			key := fmt.Sprintf("%s(%d):\t", attrName, k)
			bs.WriteString(fmt.Sprintf("%s%-10d\n", key, v))
		}
	}
	return bs.String(), nil
}

// func DumpForceAttr(logger fklog.FKLogI, userId uint64) (attrInfo string, err error) {
// 	allAttrs, err := dollassembleattrredis.GetDollAllForce(logger, userId)
// 	if err != nil {
// 		return
// 	}
// 	var bs bytes.Buffer
// 	for k, attrs := range allAttrs {
// 		bs.WriteString(fmt.Sprintf("\n武力加成来源%d(%s):\n", k, vardef.ForceSrcMap[k]))
// 		for k, v := range attrs {
// 			var attrName string
// 			attrRow := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(k)
// 			if attrRow != nil {
// 				attrName = attrRow.Name
// 			}

// 			spRow := GMazeAttrSpDescV8Cfg.GetMazeAttrSpDescV8Config(k)
// 			if spRow != nil {
// 				attrName = spRow.Attr_list_affix_desc
// 			}
// 			key := fmt.Sprintf("%s(%d):\t", attrName, k)
// 			bs.WriteString(fmt.Sprintf("%s%-10d\n", key, v))
// 		}
// 	}
// 	return bs.String(), nil
// }

func DumpNoForceAttr(logger fklog.FKLogI, userId uint64) (attrInfo string, err error) {
	var bs bytes.Buffer
	bs.WriteString("\n迷宫面板初始属性:\n")
	bs.WriteString(fmt.Sprintf("配表:%s\n", GMazeInitialAttrV8Cfg.SheetName()))
	row := GMazeInitialAttrV8Cfg.GetMazeInitialAttrV8Config(constdef.DollInitAttrCfgId2)
	if row != nil {
		for k, v := range row.Initial_attr {
			var attrName string
			attrRow := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(k)
			if attrRow != nil {
				attrName = attrRow.Name
			}

			spRow := GMazeAttrSpDescV8Cfg.GetMazeAttrSpDescV8Config(k)
			if spRow != nil {
				attrName = spRow.Attr_list_affix_desc
			}
			key := fmt.Sprintf("%s(%d):\t", attrName, k)
			bs.WriteString(fmt.Sprintf("%s%-10d\n", key, v))
		}
	}

	allAttrs, err := mazebuffinforedis.GetAllMazeBuffs(logger, userId)
	if err != nil {
		return
	}

	var bcIds []int32
	ids := mazeattrformula.GetGFXFormulaParamAttrs(constdef.DollFormulaAttack)
	bcIds = append(bcIds, ids...)
	ids = mazeattrformula.GetGFXFormulaParamAttrs(constdef.DollFormulaDefend)
	bcIds = append(bcIds, ids...)
	ids = mazeattrformula.GetGFXFormulaParamAttrs(constdef.DollFormulaBlood)
	bcIds = append(bcIds, ids...)
	ids = mazeattrformula.GetGFXFormulaParamAttrs(constdef.MazeForce)
	bcIds = append(bcIds, ids...)
	bs.WriteString("\n攻+防+血+武力值加成来源汇总:\n")
	for src, attrs := range allAttrs {
		bs.WriteString(fmt.Sprintf("\n来源%d(%s):\n", src, vardef.MazeBuffSrcDescMap[src]))
		for _, attrDb := range attrs.GetMazeRealBuffs() {
			var find bool
			for _, fpId := range bcIds {
				if fpId == attrDb.GetAttrId() {
					find = true
					break
				}
			}
			if !find {
				continue
			}
			var attrName string
			attrRow := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(attrDb.GetAttrId())
			if attrRow != nil {
				attrName = attrRow.Name
			}

			spRow := GMazeAttrSpDescV8Cfg.GetMazeAttrSpDescV8Config(attrDb.GetAttrId())
			if spRow != nil {
				attrName = spRow.Attr_list_affix_desc
			}
			key := fmt.Sprintf("%s(%d)(真):\t", attrName, attrDb.GetAttrId())
			bs.WriteString(fmt.Sprintf("%s%-10d\n", key, attrDb.GetAttrVal()))
		}

		for _, attrDb := range attrs.GetMazeShowBuffs() {

			var find bool
			for _, fpId := range bcIds {
				if fpId == attrDb.GetAttrId() {
					find = true
					break
				}
			}
			if !find {
				continue
			}

			var attrName string
			attrRow := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(attrDb.GetAttrId())
			if attrRow != nil {
				attrName = attrRow.Name
			}

			spRow := GMazeAttrSpDescV8Cfg.GetMazeAttrSpDescV8Config(attrDb.GetAttrId())
			if spRow != nil {
				attrName = spRow.Attr_list_affix_desc
			}
			key := fmt.Sprintf("%s(%d)(显):\t", attrName, attrDb.GetAttrId())
			bs.WriteString(fmt.Sprintf("%s%-10d\n", key, attrDb.GetAttrVal()))
		}
	}
	return bs.String(), nil
}

type DumpStrOrder struct {
	Order int32
	Str   string
}
