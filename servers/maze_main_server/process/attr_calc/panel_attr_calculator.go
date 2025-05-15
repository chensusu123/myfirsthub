/*
 * @Author: majian
 * @Date: 2024-07-10 15:45:44
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-01 14:52:02
 * @Desc 面板属性计算器
 */
package attr_calc

import (
	"sort"

	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeAttrListOrderV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeAttrListTypeV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeAttrSpDescV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/excel/auto/GMazeAttributeV8Cfg"
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeBuffData"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGameEquip"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazePropertyPanel"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecalcattrredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeattrorderv8"
)

type DPAC struct {
	UserId        uint64
	InParam       DPACParam
	panelAttrCfgs []*GMazeAttrListOrderV8Cfg.MazeAttrListOrderV8ConfigRow
	panelAttrMap  map[int32]int64
	Panel         *MazePropertyPanel.MazePropertyPanel
	DungeonRTime  int64                        // 地下城复活时间 10531初始值
	AttrResult    []*MazeBuffData.MazeBuffAttr // 计算结果
}

// 输入参数
type DPACParam struct {
	Force int64 // 武力值
}

func NewDPAC(userId uint64) *DPAC {
	dac := &DPAC{}
	dac.UserId = userId
	dac.Panel = new(MazePropertyPanel.MazePropertyPanel)
	dac.panelAttrMap = make(map[int32]int64)
	return dac
}

// 初始化
func (m *DPAC) Init(logger fklog.FKLogI, iParam DPACParam) error {
	m.InParam = iParam
	var err error
	var step string
	defer func() {
		if err != nil {
			logger.ErrorWF("DPAC Init fail", zap.Error(err), zap.Any("param", m.InParam),
				zap.Uint64("userId", m.UserId),
				zap.Int64("dungenonTime", m.DungeonRTime),
				zap.String("step", step))
		} else {
			logger.InfoWF("DPAC Init succ", zap.Uint64("userId", m.UserId),
				zap.Int64("dungenonTime", m.DungeonRTime),
				zap.Any("param", m.InParam))
		}
	}()

	// 初始化地下城复活时间
	// row := GMazeInitialAttrV8Cfg.GetMazeInitialAttrV8Config(constdef.DollInitAttrCfgId2)
	// if row != nil {
	// 	if val, ok := row.Initial_attr[constdef.AttrId10531]; ok {
	// 		m.DungeonRTime = val
	// 	}
	// }

	m.panelAttrCfgs = GMazeAttrListOrderV8Cfg.GetAll()

	// 查询汇总的属性
	m.panelAttrMap, err = mazecalcattrredis.HScanMazeCalcAttr(logger, m.UserId)
	if err != nil {
		step = "HScanMazeCalcAttr"
		return err
	}
	return nil
}

func (m *DPAC) Calc(logger fklog.FKLogI) error {
	var err error
	var step int32
	defer func() {
		if err != nil {
			logger.ErrorWF("DPAC Calc fail", zap.Error(err),
				zap.Int32("step", step),
				zap.Uint64("userId", m.UserId),
				zap.Any("param", m.InParam))
		} else {
			logger.InfoWF("DPAC Calc succ", zap.Uint64("userId", m.UserId),
				zap.Any("param", m.InParam))
		}
	}()

	m.AttrTransform(logger)

	allGroup := GMazeAttrListTypeV8Cfg.GetAllMazeAttrListTypeV8Config()
	sort.Slice(allGroup, func(i, j int) bool {
		return allGroup[i].Type <= allGroup[j].Type
	})

	for _, gRow := range allGroup {
		attrGroup := &MazePropertyPanel.MazePropertyGroup{}
		attrGroup.GroupId = proto.Int32(gRow.Type)
		attrGroup.GroupName = proto.String(gRow.Type_name)
		rows := mazeattrorderv8.GetPanelAttrByGroup(gRow.Type)
		for _, row := range rows {
			attrCliPb := AttrDbToPanelAttr(logger, row, m.panelAttrMap[row.Attr_id])
			attrGroup.Attrs = append(attrGroup.Attrs, attrCliPb)
		}
		m.Panel.AttrGroup = append(m.Panel.AttrGroup, attrGroup)
	}
	return nil
}

func AttrDbToPanelAttr(logger fklog.FKLogI, row *GMazeAttrListOrderV8Cfg.MazeAttrListOrderV8ConfigRow, v int64) *MazeGameEquip.EquipAttrInfo {
	attrInfo := &MazeGameEquip.EquipAttrInfo{}
	attrInfo.AttrId = proto.Int32(row.Attr_id)
	attrRow := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(row.Attr_id)
	if attrRow != nil {
		attrInfo.AttrDesc = proto.String(attrRow.Name)
		attrInfo.Figure = proto.Int32(attrRow.Figure)
		attrInfo.Comment = proto.String(row.Attr_desc)
	} else {
		logger.ErrorWF("AttrDbToPanelAttr no cfg", zap.Int32("attrId", row.Attr_id))
	}
	spRow := GMazeAttrSpDescV8Cfg.GetMazeAttrSpDescV8Config(row.Attr_id)
	if spRow != nil {
		if spRow.Attr_list_affix_suffix != "" {
			attrInfo.AttrUnits = proto.String(spRow.Attr_list_affix_suffix)
		}
		if spRow.Attr_list_affix_desc != "" {
			attrInfo.AttrDesc = proto.String(spRow.Attr_list_affix_desc)
		}
	}
	attrInfo.Value = proto.Int64(v)
	return attrInfo
}

// 面板属性转换处理
func (m *DPAC) AttrTransform(logger fklog.FKLogI) {
	// 抗性转换
	ResistanceAttrConvert(logger, m.panelAttrMap)

	// 计算特殊属性:攻防血
	// 计算攻防血
	// TODO 不再做转换
	// m.calcGFX(logger, constdef.DollAttrAttack, constdef.DollFormulaAttack)
	// m.calcGFX(logger, constdef.DollAttrDefend, constdef.DollFormulaDefend)
	// m.calcGFX(logger, constdef.DollAttrBlood, constdef.DollFormulaBlood)

	// transfer10531 地下城复活时间
	// m.transfer10531(logger)
}

func (m *DPAC) transfer10531(logger fklog.FKLogI) {
	var dgt int64
	var oldVal int64
	if oldVal, ok := m.panelAttrMap[constdef.AttrId10531]; ok {
		dgt = m.DungeonRTime - oldVal
		if dgt < 0 {
			dgt = 0
		}
	} else {
		dgt = m.DungeonRTime
	}
	SetMazeAttrKv(m.panelAttrMap, constdef.AttrId10531, dgt)
	logger.InfoWF("transfer10531",
		zap.Int64("new", dgt),
		zap.Int64("old", oldVal),
		zap.Int64("rt", m.DungeonRTime))
}

func (m *DPAC) calcGFX(logger fklog.FKLogI, attrId, formulaId int32) {
	SetMazeAttrKv(m.panelAttrMap, attrId, m.panelAttrMap[formulaId])
	logger.InfoWF("calcGFX",
		zap.Uint64("userId", m.UserId),
		zap.Int32("attrId", attrId),
		zap.Int32("formulaId", formulaId))
}
