/*
 * @Author: majian
 * @Date: 2025-03-11 11:32:50
 * @Last Modified by: majian
 * @Last Modified time: 2025-04-01 14:43:34
 */
package mazeattrlogic

import (
	"fmt"
	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/config/GMazeAttrSpDescV8Cfg"
	"maze_game_server/config/GMazeAttributeV8Cfg"
	"maze_game_server/io/kafka/mazeattrchgrecord"
	"maze_game_server/io/kafka/mazeattrmsg"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/module/mazeattrformula"
	"maze_game_server/pb/server/MazeBuffData"
	"maze_game_server/servers/maze_main_server/process/attr_calc/commonlogic"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type DAC struct {
	fklog.FKLogI
	UserId  uint64
	InParam *DACParam
	//	BuffCenterAttrsIn map[int32]int64                        // buff中心
	DollBuffCenterIn map[int32][]*MazeBuffData.MazeBuffAttr // 人偶buff中心

	AttrResultMap  map[int32]int64 // 结果汇总
	CalcAttrOldMap map[int32]int64 // 旧值

	ErrNoRetry          bool                          // 失败不能重试
	FormulaAttrIds      []int32                       // 需要计算的公式属性ID列表
	CareAttrs           map[int32]bool                // 相关属性列表
	FormulaAttrExtraMap map[int32]string              // 计算信息
	Acr                 *commonlogic.AttrCalcRecorder // 属性计算记录器
	ReplaceData         *ReplaceData                  // 替换的数据
}

type ReplaceData struct {
	ReplaceAttrs map[int32]int64 // 替换的属性Id
	ReplaceSrc   int32           // 替换的来源
	PreviewId    int64           // 预览Id
}

type PreviewParam struct {
	PreviewId    int64           // 预览Id
	CareAttrs    []int32         // 预览的属性
	BuffSrc      int32           // buff来源
	RepalceAttrs map[int32]int64 // 要替换的属性
}

// 输入参数
type DACParam struct {
	ChgType    int32  // 变化类型
	ChgSubType int32  // 变化子类型
	ChgDesc    string // 变化原因
	IsPreview  bool   // 是否预览
	Session    string // session 业务串联用
}

func NewDACParam() *DACParam {
	obj := &DACParam{}
	return obj
}

func NewReplaceData() *ReplaceData {
	obj := new(ReplaceData)
	obj.ReplaceAttrs = make(map[int32]int64)
	return obj
}

func NewDAC(logger fklog.FKLogI, userId uint64) *DAC {
	dac := &DAC{}
	dac.UserId = userId
	dac.FKLogI = logger
	dac.AttrResultMap = make(map[int32]int64)
	// dac.BuffCenterAttrsIn = make(map[int32]int64)
	dac.DollBuffCenterIn = make(map[int32][]*MazeBuffData.MazeBuffAttr)
	dac.CalcAttrOldMap = make(map[int32]int64)
	dac.FormulaAttrExtraMap = make(map[int32]string)
	dac.Acr = commonlogic.NewAttrCalcRecorder()
	dac.ReplaceData = NewReplaceData()
	dac.InParam = NewDACParam()
	dac.CareAttrs = make(map[int32]bool)
	return dac
}

func (m *DAC) CloneData() *DAC {
	cp := NewDAC(m, m.UserId)
	cp.CalcAttrOldMap = m.CalcAttrOldMap
	//	cp.BuffCenterAttrsIn = m.BuffCenterAttrsIn
	cp.ErrNoRetry = m.ErrNoRetry
	//	cp.FormulaAttrIds = m.FormulaAttrIds
	cp.InParam = m.InParam
	for src, attrs := range m.DollBuffCenterIn {
		for _, attr := range attrs {
			attrTmp := &MazeBuffData.MazeBuffAttr{}
			attrTmp.AttrId = proto.Int32(attr.GetAttrId())
			attrTmp.AttrVal = proto.Int64(attr.GetAttrVal())
			cp.DollBuffCenterIn[src] = append(cp.DollBuffCenterIn[src], attrTmp)
		}
	}
	return cp
}

// 初始化
func (m *DAC) InitData(iParam *DACParam) error {
	var err error
	defer func() {
		if err != nil {
			m.ErrorWF("DAC InitData fail", zap.Error(err), zap.Any("param", iParam))
		} else {
			m.InfoWF("DAC InitData succ", zap.Any("param", iParam))
		}
	}()
	m.InParam = iParam
	if !m.InParam.IsPreview { // 非预览计算所有公式属性
		m.FormulaAttrIds = append(m.FormulaAttrIds, constdef.DollFormulaAttack,
			constdef.DollFormulaDefend, constdef.DollFormulaBlood, constdef.MazeForce)
	}
	// 初始化buff中心数据(展示属性+非展示属性)
	// attrIdSet := simpleset.NewSet()
	// attrIdSet.Add(constdef.DollAttrTypeCalc)
	// attrIdSet.Add(constdef.DollAttrTypeShow)
	// attrIdSet.Add(constdef.DollAttrTypeForce)
	// attrSet := commonlogic.GetAllMazeAttrByType(attrIdSet)
	// if len(attrSet) > 0 {
	// 	m.BuffCenterAttrsIn, err = BuffManagerRedis.GetBuffs(context.TODO(), m, m.UserId, 0, attrSet)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// 初始化人偶buff中心
	dollBuffs, e := mazebuffinforedis.GetMazeBuffsV2(m, m.UserId, 3)
	if e != nil {
		err = e
		return err
	}
	m.DollBuffCenterIn = dollBuffs
	// 初始化旧值
	m.CalcAttrOldMap, err = mazecalcattrredis.HScanMazeCalcAttr(m, m.UserId)
	if err != nil {
		return err
	}
	return nil
}

// 数据预处理
func (m *DAC) Prepare() error {
	if !m.NeedReplaceAttr() {
		return nil
	}
	if v, ok := m.DollBuffCenterIn[m.ReplaceData.ReplaceSrc]; ok {
		m.DollBuffCenterIn[m.ReplaceData.ReplaceSrc] = ReplaceAttr(v, m.ReplaceData)
	} else {
		m.DollBuffCenterIn[m.ReplaceData.ReplaceSrc] = ReplaceAttr(nil, m.ReplaceData)
	}
	m.InfoWF("Prepare replace", zap.Any("attrs", m.ReplaceData))
	return nil
}

// 计算
func (m *DAC) Calc() error {

	var err error
	var step string
	defer func() {
		if err != nil {
			m.ErrorWF("DAC Calc fail", zap.Error(err), zap.String("step", step),
				zap.Any("param", m.InParam))
		} else {
			m.InfoWF("DAC Calc succ", zap.Any("param", m.InParam))
		}
	}()

	// 汇总不同来源属性
	for _, cbNode := range buffFuncList {
		srcName := commonlogic.GetSrcName(cbNode.SrcType)
		err = cbNode.Cbf(m, cbNode.SrcType)
		if err != nil {
			step = fmt.Sprintf("Cbf:%s", srcName)
			m.ErrorWF("DAC Calc fail", zap.Error(err),
				zap.Int("srcType", cbNode.SrcType),
				zap.String("srcName", srcName),
				zap.Any("param", m.InParam))
			return err
		}
		m.InfoWF("DAC Cbf succ", zap.Int("srcType", cbNode.SrcType),
			zap.String("srcName", srcName))
	}
	// 计算公式属性，比如攻击 防御 耐久需要按公式计算
	err = m.CalcFormulaAttr()
	if err != nil {
		step = "CalcFormulaAttr"
		return err
	}
	return err
}

// 计算公式属性
func (m *DAC) CalcFormulaAttr() error {
	for _, formulaAttrId := range m.FormulaAttrIds {
		v, ext, e := mazeattrformula.CalcMazeFormulaAttr(m, formulaAttrId, m.AttrResultMap)
		if e != nil {
			return e
		}
		commonlogic.AddtionMazeAttrKv(m.AttrResultMap, formulaAttrId, v)
		m.FormulaAttrExtraMap[formulaAttrId] = ext
	}
	return nil
}

// 保存通知
func (m *DAC) End() error {
	if m.InParam.IsPreview {
		return nil
	}
	// 对比变化和可以删除的Id
	chgAttrs, delIds := m.compareDiff()
	var err error
	defer func() {
		if err != nil {
			m.ErrorWF("DAC End fail", zap.Error(err), zap.Any("param", m.InParam),
				zap.Any("chgAttrs", chgAttrs), zap.Any("delIds", delIds))
		} else {
			m.InfoWF("DAC End succ", zap.Any("param", m.InParam),
				zap.Any("chgAttrs", chgAttrs),
				zap.Any("delIds", delIds))
		}
	}()
	if len(chgAttrs) <= 0 {
		return nil
	}
	err = mazecalcattrredis.SaveMazeCalcAttr(m, m.UserId, chgAttrs)
	if err != nil {
		return err
	}

	// 删除0值的ID，节省空间，放置key过大，失败可以忽略
	if len(delIds) > 0 && len(m.CalcAttrOldMap) > int(commonlogic.CleanIdLen) {
		mazecalcattrredis.HDelMazeCalcAttr(m, m.UserId, delIds)
	}

	// 记录流水
	m.Record(chgAttrs)

	// 变化通知
	m.Notify(chgAttrs)
	return err
}

// 失败是否需要重试
func (m *DAC) ErrNeedRetry() bool {
	return !m.ErrNoRetry
}

func (m *DAC) compareDiff() (chgs map[int32]int64, delIds []int32) {
	chgs = make(map[int32]int64)
	for k, newVal := range m.AttrResultMap {
		if oldVal, ok := m.CalcAttrOldMap[k]; ok {
			if oldVal != newVal {
				// 值有变化
				chgs[k] = newVal
			}
			continue
		}
		// 新增属性
		chgs[k] = newVal
	}
	for k, oldVal := range m.CalcAttrOldMap {
		if _, ok := m.AttrResultMap[k]; ok {
			continue
		}
		if oldVal <= 0 {
			delIds = append(delIds, k)
			continue
		}
		// 旧属性不存在，置成0
		chgs[k] = 0
		delIds = append(delIds, k)
	}
	return chgs, delIds
}

func (m *DAC) Record(chgAttrs map[int32]int64) {
	// 打消息变化
	now := time.Now().UnixNano() / 1000000
	for k, newVal := range chgAttrs {
		var attrName string
		spRow := GMazeAttrSpDescV8Cfg.GetMazeAttrSpDescV8Config(k)
		if spRow != nil {
			if spRow.Equip_affix_desc != "" {
				attrName = spRow.Equip_affix_desc
			}
		}
		if attrName == "" {
			attrCfg := GMazeAttributeV8Cfg.GetMazeAttributeV8Config(k)
			if attrCfg != nil {
				attrName = attrCfg.Name
			}
		}

		oldAttr := m.CalcAttrOldMap[k]
		dacr := &structsdef.MazeGameAttrChgRecord{}
		dacr.UserId = m.UserId
		dacr.AttrId = k
		dacr.ChgType = m.InParam.ChgType
		dacr.ChgSubType = m.InParam.ChgSubType
		dacr.ChgDesc = fmt.Sprintf("%s(%s)", m.InParam.ChgDesc, attrName)
		dacr.OldVal = oldAttr
		dacr.NewVal = newVal
		dacr.CreateTime = now
		// 如果计算属性，会有扩展信息，记录计算过程的日志，记录到流水
		if ext, ok := m.FormulaAttrExtraMap[k]; ok {
			dacr.Extra = ext
		} else {
			dacr.Extra = m.Acr.DumpWeightInfo(k)
		}
		mazeattrchgrecord.SendMazeGameAttrChgRecord(m, dacr)
	}
}

func (m *DAC) Notify(chgAttrs map[int32]int64) {
	// 打消息变化
	now := time.Now().UnixNano() / 1000000
	msg := &structsdef.DollAttrChgNotify{}
	msg.UserId = m.UserId
	msg.ChgType = m.InParam.ChgType
	msg.ChgSubType = m.InParam.ChgSubType
	msg.ChgDesc = m.InParam.ChgDesc
	msg.CreateTime = now
	msg.Session = m.InParam.Session

	for k, newAttr := range chgAttrs {
		chgAttr := &structsdef.AttrChgInfo{}
		chgAttr.AttrId = k
		chgAttr.OldVal = m.CalcAttrOldMap[k]
		chgAttr.CurVal = newAttr
		msg.ChgAttrs = append(msg.ChgAttrs, chgAttr)
	}
	mazeattrmsg.SendMazeAttrChgNotify(m, msg)
	commonlogic.NotifyClientAttrChg(m, m.UserId, msg)
}

func (m *DAC) NeedReplaceAttr() bool {
	if m.InParam.IsPreview {
		if m.ReplaceData != nil && m.ReplaceData.ReplaceSrc > 0 {
			return true
		}
	}
	return false
}

func (m *DAC) GetCareAttrs(careAttrs []int32) map[int32]int64 {
	rs := make(map[int32]int64)
	for _, careId := range careAttrs {
		rs[careId] = m.AttrResultMap[careId]
	}
	return rs
}

func (m *DAC) SetPreviewInfo(param *PreviewParam) {
	m.ReplaceData.ReplaceSrc = param.BuffSrc
	m.ReplaceData.ReplaceAttrs = param.RepalceAttrs
	m.ReplaceData.PreviewId = param.PreviewId

	if len(param.CareAttrs) == 1 && param.CareAttrs[0] == constdef.MazeForce {
		cares := mazeattrformula.GetGFXFormulaParamAttrs(constdef.MazeForce)
		for _, attr := range cares {
			m.CareAttrs[attr] = true
		}
	}

	for _, attrId := range param.CareAttrs {
		if attrId == constdef.DollFormulaAttack ||
			attrId == constdef.DollFormulaDefend ||
			attrId == constdef.DollFormulaBlood ||
			attrId == constdef.MazeForce {
			m.FormulaAttrIds = append(m.FormulaAttrIds, attrId)
		}
	}
	m.InfoWF("SetPreviewInfo", zap.Any("param", param),
		zap.Any("careAttrs", m.CareAttrs),
		zap.Any("formulaAttrs", m.FormulaAttrIds))
}

func ReplaceAttr(in []*MazeBuffData.MazeBuffAttr, repData *ReplaceData) []*MazeBuffData.MazeBuffAttr {
	for k, v := range repData.ReplaceAttrs {
		var has bool
		for _, attr := range in {
			if attr.GetAttrId() == k {
				attr.AttrVal = proto.Int64(v)
				has = true
				break
			}
		}
		if !has {
			in = append(in, &MazeBuffData.MazeBuffAttr{AttrId: proto.Int32(k), AttrVal: proto.Int64(v)})
		}
	}
	return in
}

func (m *DAC) IsNoCareAttr(attrId int32) bool {
	if m.InParam.IsPreview {
		if !m.CareAttrs[attrId] {
			return true
		}
	}
	return false
}
