/*
 * @Author: majian
 * @Date: 2025-03-11 11:31:26
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-29 10:37:21
 */
package mazeattrlogic

import (
	"maze_game_server/config/GMazeAttributeV8Cfg"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"

	"maze_game_server/common/constdef"
	"maze_game_server/servers/maze_main_server/process/attr_calc/commonlogic"

	"go.uber.org/zap"
)

type DollAttrCalcCBF func(m *DAC, src int) (err error)

type CBFNode struct {
	Cbf     DollAttrCalcCBF
	SrcType int
}

var buffFuncList []*CBFNode

func init() {
	RegDollAttrCb(constdef.MazePartBc, HandleMazeBc)
	//	RegDollAttrCb(constdef.MazePartOldBc, HandleBc)
	RegDollAttrCb(constdef.MazePartInitAttr, HandleMazeInit)
	// RegDollAttrCb(constdef.AttrSrcInitShow, HandleDollInitShow)
}

func RegDollAttrCb(st int, f DollAttrCalcCBF) {
	buffFuncList = append(buffFuncList, &CBFNode{Cbf: f, SrcType: st})
}

func HandleMazeBc(m *DAC, src int) (err error) {
	for dollSrc, dollAttrSrcMap := range m.DollBuffCenterIn {
		for _, dollAttr := range dollAttrSrcMap {
			if m.IsNoCareAttr(dollAttr.GetAttrId()) {
				continue
			}
			commonlogic.AddtionMazeAttr(m.AttrResultMap, dollAttr)
			DumpAttrBySrc(m, int(dollSrc), dollAttr.GetAttrId(), dollAttr.GetAttrVal(), "迷宫buff中心")
			if !m.InParam.IsPreview {
				m.Acr.AddWeight(dollAttr.GetAttrId(), dollSrc, dollAttr.GetAttrVal())
			}
		}
	}
	return
}

// func HandleBc(m *DAC, src int) (err error) {
// 	for id, val := range m.BuffCenterAttrsIn {
// 		if id <= 0 {
// 			continue
// 		}
// 		if val <= 0 { // 没有负值
// 			continue
// 		}
// 		if m.IsNoCareAttr(id) {
// 			continue
// 		}
// 		commonlogic.AddtionMazeAttrKv(m.AttrResultMap, id, val)
// 		DumpAttrBySrc(m, int(constdef.MazeBuffSrcOldBC), id, val, "")
// 		if !m.InParam.IsPreview {
// 			m.Acr.AddWeight(id, constdef.MazeBuffSrcOldBC, val)
// 		}
// 	}
// 	return
// }

func HandleMazeInit(m *DAC, src int) (err error) {
	var (
		initialAttrs = make(map[int32]int64)
	)

	for _, attr := range GMazeAttributeV8Cfg.GetAll() {
		// 初始化属性值非0则为初始属性
		if attr.Initial_value != 0 {
			initialAttrs[attr.Id] = int64(attr.Initial_value)
		}
	}

	// row := GMazeInitialAttrV8Cfg.GetMazeInitialAttrV8Config(constdef.MazeInitAttrCfgId)
	// if row == nil {
	// 	return
	// }
	//	excludeAttrs := dollinitattrv8.GetDollInitExcludeAttrs()

	for id, val := range initialAttrs {
		if id <= 0 {
			continue
		}
		if val <= 0 { // 没有负值
			continue
		}
		// if _, ok := excludeAttrs[id]; ok {
		// 	m.WarnWF("HandleDollInit excludeAttrs ignore", zap.Int32("id", id))
		// 	continue
		// }

		if m.IsNoCareAttr(id) {
			continue
		}

		commonlogic.AddtionMazeAttrKv(m.AttrResultMap, id, val)
		DumpAttrBySrc(m, int(constdef.MazeBuffSrcInit), id, val, "")
		if !m.InParam.IsPreview {
			m.Acr.AddWeight(id, constdef.MazeBuffSrcInit, val)
		}
	}
	return
}

// 人偶初始展示属性
// func HandleDollInitShow(m *DAC, src int) (err error) {
// 	row := GMazeInitialAttrV8Cfg.GetMazeInitialAttrV8Config(constdef.DollInitAttrCfgId2)
// 	if row == nil {
// 		return
// 	}
// 	excludeAttrs := dollinitattrv8.GetDollInitExcludeAttrs()
// 	for id, val := range row.Initial_attr {
// 		if id <= 0 {
// 			continue
// 		}
// 		if val <= 0 { // 没有负值
// 			continue
// 		}
// 		if id == constdef.AttrId10531 {
// 			m.WarnWF("HandleDollInitShow 10531 ignore", zap.Int32("id", id))
// 			continue
// 		}
// 		if _, ok := excludeAttrs[id]; ok {
// 			m.WarnWF("HandleDollInitShow excludeAttrs ignore", zap.Int32("id", id))
// 			continue
// 		}
// 		if m.IsNoCareAttr(id) {
// 			continue
// 		}

// 		commonlogic.AddtionMazeAttrKv(m.AttrResultMap, id, val)
// 		DumpAttrBySrc(m, src, id, val, "")
// 		if !m.InParam.IsPreview {
// 			m.Acr.AddWeight(id, int32(src), val)
// 		}

// 	}
// 	return
// }

func DumpAttrBySrc(logger fklog.FKLogI, src int, k int32, v int64, extra string) {
	logger.InfoWF("DumpAttrBySrc", zap.Int("src", src), zap.String("srcName", commonlogic.GetSrcName(src)),
		zap.Int32("k", k), zap.Int64("v", v), zap.String("extra", extra))
}
