/*
 * @Author: majian
 * @Date: 2024-04-03 16:26:41
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-13 22:07:28
 */
package equip

import (
	"fmt"
	"strings"
	"time"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"maze_game_server/common/function/randfuncs"
	"maze_game_server/config/GMazeEquipAffixLimitV8Cfg"
	"maze_game_server/config/GMazeEquipAffixRandPoolV8Cfg"
	"maze_game_server/config/GMazeEquipAffixRollTypeV8Cfg"
	"maze_game_server/config/GMazeEquipAffixSpRuleV8Cfg"
	"maze_game_server/config/GMazeEquipAttrStageV8Cfg"
	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/excel/mazeequipaffixrandpoolv8"
	"maze_game_server/excel/mazeequipaffixrolltypev8"
	"maze_game_server/io/kafka/mazeequipinstancerecord"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/errors"
	"maze_game_server/pb/server/MazeEquipCache"
	"maze_game_server/pb/server/MazeEquipSvr"
)

// 用户参数外部
type DEIUWParam struct {
	ChargeStage       int32              `json:"charge_stage"`
	OpType            int32              `json:"op_type"`
	TradeNum          uint64             `json:"trade_num"`
	PoolLimitMap      map[int32]struct{} `json:"-"`
	PoolLimitBase2Map map[int32]struct{} `json:"-"`
	// DollLv            int64              `json:"doll_lv"`
	// BarrierMap map[int32]int64 `json:"barrier_map"`
}

// 装备参数外部
type DEIEWParam struct {
	EquipId            int32  // 装备 ID
	RuleId             int32  // 规则id
	HasSuit            int32  // 拥有套装
	AssignSuit         int32  // 指定套装
	AssignEquipType    int32  // 指定装备类型
	Guid               int64  // 分配的guid
	Score              int32  // 分值
	Condition          string // 条件
	AssignFiveElemHead int32  // 指定五行头类型
	AssignFiveElemTail int32  // 指定五行尾类型
}

// 实例化结果返回
type EDIResult struct {
	EquipInfo    *MazeEquipCache.MazeEquipInfoDb // 实例化的装备
	DropCntScore int32                           // 掉落次数分值
	TotalScore   int32                           // 装备累计分数
}

// 用户参数内部
type DEIUParam struct {
	ChargeStage    int32 // 充值档位
	StageDropScore int32 // 档位掉落分值
	StageFactor    int32 // 档位系数
	BaseRandoms    []int64
	ElementsRandom int64
}

// 人偶版装备配置
type DollEquipInfoCfg struct {
	EquipId           int32              `json:"equip_id"`
	AffixBaseNum      map[int32]int32    `json:"affix_base_num"`
	AffixBasePool     map[int32]int32    `json:"affix_base_pool"`
	AffixRandNum      map[int32]int32    `json:"affix_rand_num"`
	AffixRandPool     map[int32]int32    `json:"affix_rand_pool"`
	SuiteId           map[int32]int32    `json:"suite_id"`
	LimitNum          int32              `json:"limit_num"`
	LimitBase2Num     int32              `json:"limit_base_2_num"`
	SubTypeRandom     map[int32]int32    `json:"sub_type_random"`
	ScoreGroup        int32              `json:"score_group"`
	PoolLimitMap      map[int32]struct{} `json:"-"`
	PoolLimitBase2Map map[int32]struct{} `json:"-"`
	FontSoul          map[int32]int32    `json:"font_soul"`
	TailSoul          map[int32]int32    `json:"tail_soul"`
}

// 装备参数内部
type DEIEParam struct {
	EquipId         int32             // 装备 ID
	TotalScore      int32             // 累计分值
	EquipGuid       int64             // 装备guid
	EquipInfoCfg    *DollEquipInfoCfg // 装备库配置
	RuleId          int32             // 规则id
	HasSuit         int32             // 拥有套装
	AssignSuit      int32             // 指定套装
	AssignEquipType int32             // 指定装备类型
	// IsSeal             int32             //是否需要鉴定
	AssignFiveElemHead int32 // 指定五行头类型
	AssignFiveElemTail int32 // 指定五行尾类型
}

// 人偶装备实例化接口
type DEInstance struct {
	fklog.FKLogI
	UserID uint64 // 用户Id
	// 用户参数
	UParam *DEIUParam
	// 装备参数
	EParam *DEIEParam
	// 实例化信息
	EquipInfo *MazeEquipCache.MazeEquipInfoDb
	Record    *MazeGameEquipInstanceRecord
}

func NewDEInstance(logger fklog.FKLogI, userID uint64) *DEInstance {
	ins := &DEInstance{}
	ins.FKLogI = logger
	ins.UserID = userID
	return ins
}

func NewDEIEWParam(equipInfo *MazeEquipSvr.SvrEquipInfo) *DEIEWParam {
	ep := &DEIEWParam{}
	ep.EquipId = equipInfo.GetEquipId()
	for _, condition := range equipInfo.GetConditions() {
		if condition.GetValue() <= 0 {
			continue
		}
		if condition.GetId() == int32(MazeEquipSvr.EQUIP_ADD_CONDITION_SPECIAL_EQUIP) {
			ep.RuleId = int32(condition.GetValue())
		}
		if condition.GetId() == int32(MazeEquipSvr.EQUIP_ADD_CONDITION_BE_SURE_SUIT) {
			ep.HasSuit = int32(condition.GetValue())
		}
		if condition.GetId() == int32(MazeEquipSvr.EQUIP_ADD_CONDITION_ASSIGN_SUIT) {
			ep.AssignSuit = int32(condition.GetValue())
		}
		if condition.GetId() == int32(MazeEquipSvr.EQUIP_ADD_CONDITION_ASSIGN_EQUIP_SUB_TYPE) {
			ep.AssignEquipType = int32(condition.GetValue())
		}
	}
	return ep
}

// 初始化用户参数和装备参数
func (dei *DEInstance) DeiInit(uwparam *DEIUWParam, ewparam *DEIEWParam) error {
	e := dei.InitUParam(uwparam)
	if e != nil {
		return e
	}
	e = dei.InitEParam(uwparam, ewparam)
	if e != nil {
		return e
	}
	dei.InitEquipInstanceLog(uwparam, ewparam)
	return e
}

// 初始化用户参数
func (dei *DEInstance) InitUParam(uwparam *DEIUWParam) error {
	dei.UParam = &DEIUParam{}
	dei.UParam.ChargeStage = uwparam.ChargeStage
	e := initStageInfo(dei, dei.UParam)
	if e != nil {
		dei.ErrorWF("DEInstance InitUParam initStageInfo err", zap.Any("uwp", uwparam), zap.Error(e))
		return e
	}
	dei.InfoWF("DEInstance InitUParam dump", zap.Any("uwp", uwparam), zap.Any("up", dei.UParam))
	return nil
}

// 初始化装备参数
func (dei *DEInstance) InitEParam(uwparam *DEIUWParam, ewparam *DEIEWParam) error {
	dei.EParam = &DEIEParam{}
	var (
		e error
		// newGuid      int64
		equipInfoCfg *DollEquipInfoCfg
	)

	defer func() {
		if e != nil {
			dei.ErrorWF("DEInstance InitEParam dump err", zap.Error(e), zap.Any("ewp", ewparam), zap.Any("ep", dei.EParam))
		} else {
			dei.InfoWF("DEInstance InitEParam dump", zap.Any("ewp", ewparam), zap.Any("ep", dei.EParam))
		}
	}()
	equipCfg := GMazeEquipInfoV8Cfg.Get(ewparam.EquipId)
	if equipCfg == nil {
		e = errors.New("装备配置不存在")
		return e
	}
	if ewparam.RuleId > 0 {
		equipRuleCfg := GMazeEquipAffixSpRuleV8Cfg.Get(ewparam.RuleId)
		if equipRuleCfg == nil {
			e = errors.New("特殊装备配置不存在")
			return e
		}
		equipInfoCfg = InitEquipInfoCfg(equipCfg.Id, equipRuleCfg.Affix_base_num, equipRuleCfg.Affix_base_pool,
			equipRuleCfg.Affix_rand_num, equipRuleCfg.Affix_rand_pool, equipRuleCfg.Suite_id, equipRuleCfg.Sub_type_random, equipCfg.Score_group)
	} else {
		equipInfoCfg = InitEquipInfoCfg(equipCfg.Id, equipCfg.Affix_base_num, equipCfg.Affix_base_pool,
			equipCfg.Affix_rand_num, equipCfg.Affix_rand_pool, equipCfg.Suite_id, equipCfg.Sub_type_random, equipCfg.Score_group)
	}
	equipAffixLimitCfg := GMazeEquipAffixLimitV8Cfg.Get(equipCfg.Quality*10000 + equipCfg.Pos*1000)
	if equipAffixLimitCfg != nil {
		equipInfoCfg.LimitNum = equipAffixLimitCfg.Limit_num
		equipInfoCfg.LimitBase2Num = equipAffixLimitCfg.Limit_base2_num
	}
	equipInfoCfg.PoolLimitMap = uwparam.PoolLimitMap
	equipInfoCfg.PoolLimitBase2Map = uwparam.PoolLimitBase2Map

	dei.EParam.EquipId = ewparam.EquipId
	dei.EParam.RuleId = ewparam.RuleId
	dei.EParam.HasSuit = ewparam.HasSuit
	dei.EParam.AssignSuit = ewparam.AssignSuit
	dei.EParam.AssignEquipType = ewparam.AssignEquipType
	dei.EParam.AssignFiveElemHead = ewparam.AssignFiveElemHead
	dei.EParam.AssignFiveElemTail = ewparam.AssignFiveElemTail

	dei.EParam.TotalScore = ewparam.Score
	dei.EParam.EquipInfoCfg = equipInfoCfg

	dei.EParam.EquipGuid = ewparam.Guid

	dei.InfoWF("equip groupId and score",
		zap.Any("equipId", ewparam.EquipId),
		zap.Int64("guid", ewparam.Guid),
		zap.Any("groupId", equipCfg.Score_group),
		zap.Any("addScore", dei.UParam.StageDropScore),
		zap.Any("totalScore", dei.EParam.TotalScore))
	return nil
}

func (dei *DEInstance) InitEquipInstanceLog(uwparam *DEIUWParam, ewparam *DEIEWParam) *MazeGameEquipInstanceRecord {
	obj := &MazeGameEquipInstanceRecord{}
	obj.MazeGameEquipInstanceRecord = new(mazeequipinstancerecord.MazeGameEquipInstanceRecord)
	obj.UserId = dei.UserID
	obj.OpType = uwparam.OpType
	obj.TradeNum = uwparam.TradeNum
	obj.StageFactor = dei.UParam.StageFactor
	obj.Conditions = ewparam.Condition
	dei.Record = obj
	return obj
}

// 实例化
func (dei *DEInstance) DeiInstance() error {
	// 填充装备信息
	dei.EquipInfo = &MazeEquipCache.MazeEquipInfoDb{}
	dei.EquipInfo.EquipGuid = proto.Int64(dei.EParam.EquipGuid)
	dei.EquipInfo.EquipId = proto.Int32(dei.EParam.EquipId)
	if dei.EParam.RuleId > 0 {
		dei.EquipInfo.RuleId = proto.Int32(dei.EParam.RuleId)
	}
	dei.EquipInfo.MakeTime = proto.Int64(time.Now().Unix())

	dei.EquipInfo.EquipScore = proto.Int32(dei.EParam.TotalScore)
	// 生成套装属性
	e := dei.GenSuiteId()
	if e != nil {
		return e
	}
	// 生成基础属性
	e = dei.GenBaseAttr()
	if e != nil {
		return e
	}
	// 生成初始属性属性
	e = dei.GenRandAttr()
	if e != nil {
		return e
	}

	e = dei.GenEquipSubType()
	if e != nil {
		return e
	}

	dei.BeforeRecord()
	dei.DebugWF("DEInstance DeiInstance dump", zap.Any("equipInfo", dei.EquipInfo))
	return nil
}

// 获取实例化结果
func (dei *DEInstance) GetResult() *EDIResult {
	deiR := &EDIResult{}
	deiR.EquipInfo = dei.EquipInfo
	deiR.DropCntScore = dei.UParam.StageDropScore
	deiR.TotalScore = dei.EParam.TotalScore
	return deiR
}

// 生成基础属性
func (dei *DEInstance) GenBaseAttr() error {
	dei.DebugWF("开始实例化Base属性",
		zap.Int32("equipId", dei.EParam.EquipId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
	equipInfoCfg := dei.EParam.EquipInfoCfg
	dei.DebugWF("开始随机Base属性条数", zap.Int64("equipGuid", dei.EParam.EquipGuid))
	// fkfmt.Println("GenBaseAttr0")
	baseNum := randfuncs.RandByWeightV2(dei, equipInfoCfg.AffixBaseNum, true)
	dei.DebugWF("Base属性条数随机结果", zap.Any("baseNum", baseNum), zap.Int64("equipGuid", dei.EParam.EquipGuid))
	if baseNum == 0 {
		return nil
	}
	//	fkfmt.Println("GenBaseAttr1")
	// poolLimitMap := GMazeEquipAffixRandPoolV8CfgEx.GetPoolLimitCfg(dei.EParam.EquipInfoCfg.LimitAttr)
	// poolLimitBase2Map := GMazeEquipAffixRandPoolV8CfgEx.GetPoolLimitCfg(dei.EParam.EquipInfoCfg.LimitBase2Attr)
	suitRemGroupMap, err := GetEquipSuitRemGroupMap(dei.EquipInfo.GetSuitId())
	if err != nil {
		dei.ErrorWF("GenBaseAttr GetEquipSuitRemGroupMap error",
			zap.Int32("suitId", dei.EquipInfo.GetSuitId()), zap.Int32("equipId", dei.EParam.EquipId), zap.Error(err))
		return err
	}
	// fkfmt.Println("GenBaseAttr2")
	var limitNum, limitBase2Num int32
	groupBasePoolMap := make(map[int32]map[int32]int32, 0)
	for poolId, index := range equipInfoCfg.AffixBasePool {
		if index == 0 {
			continue
		}
		if groupBasePoolMap[index] == nil {
			groupBasePoolMap[index] = make(map[int32]int32, 0)
		}
		groupBasePoolMap[index][poolId] = 1000
	}
	// fkfmt.Println("GenBaseAttr3")
	attrGroupList := make([]int32, 0)
	baseAttrs := make([]*MazeEquipCache.BaseAttrInfo, 0)
	remGroupMap := make(map[int32]struct{}, 0)
	for index := int32(1); index <= baseNum; index++ {
		basePoolMap, ok := groupBasePoolMap[index]
		if !ok {
			dei.ErrorWF("GenBaseAttr get GMazeEquipAffixRandPoolV8Cfg weight fail",
				zap.Int32("index", index), zap.Int32("equipId", dei.EParam.EquipId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
			return errors.New("配置不存在")
		}
		dei.DebugWF("开始随机Base属性池", zap.Any("条数", index), zap.Int64("equipGuid", dei.EParam.EquipGuid))
		poolId := randfuncs.RandByWeightV2(dei, basePoolMap, true)
		dei.DebugWF("随机Base属性池结果", zap.Any("poolId", poolId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
		poolMap := mazeequipaffixrandpoolv8.GetEquipPoolWeightCfg(poolId)
		if len(poolMap) == 0 {
			dei.ErrorWF("GenBaseAttr GetEquipPoolWeightCfg weight fail",
				zap.Int32("poolId", poolId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
			return errors.New("属性配置不存在")
		}
		newPoolWeightMap := make(map[int32]int32, 0)
		for k, v := range poolMap {
			if _, ok := remGroupMap[v.Group]; ok {
				dei.DebugWF("随机Base属性词条去重", zap.Any("去重词条id", k), zap.Any("去重组id", v.Group), zap.Int64("equipGuid", dei.EParam.EquipGuid))
				continue
			}
			if limitNum >= dei.EParam.EquipInfoCfg.LimitNum {
				if _, ok := equipInfoCfg.PoolLimitMap[v.Affix_id]; ok {
					dei.DebugWF("随机Base属性抗性词条去重", zap.Any("去重词条id", k), zap.Any("去重组id", v.Group), zap.Int64("equipGuid", dei.EParam.EquipGuid))
					continue
				}
			}
			if limitBase2Num >= dei.EParam.EquipInfoCfg.LimitBase2Num {
				if _, ok := equipInfoCfg.PoolLimitBase2Map[v.Affix_id]; ok {
					dei.DebugWF("随机Base属性垃圾词条去重", zap.Any("去重词条id", k), zap.Any("去重组id", v.Group), zap.Int64("equipGuid", dei.EParam.EquipGuid))
					continue
				}
			}
			if _, ok := suitRemGroupMap[v.Affix_id]; ok {
				dei.DebugWF("套装属性词条去重", zap.Any("去重词条id", k), zap.Any("去重组id", v.Group), zap.Int64("equipGuid", dei.EParam.EquipGuid))
				continue
			}
			newPoolWeightMap[k] = v.Weight
		}
		dei.DebugWF("随机Base属性词条id", zap.Int64("equipGuid", dei.EParam.EquipGuid))
		attrGroupId := randfuncs.RandByWeightV2(dei, newPoolWeightMap, true)
		dei.DebugWF("随机Base属性词条id结果", zap.Any("attrGroupId", attrGroupId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
		attrGroupList = append(attrGroupList, attrGroupId)
		equipAffixCfg := GMazeEquipAffixRandPoolV8Cfg.Get(attrGroupId)
		if equipAffixCfg == nil {
			dei.ErrorWF("GenBaseAttr get GMazeEquipAffixRandPoolV8Cfg fail",
				zap.Int32("poolId", poolId),
				zap.Any("attrGroupId", attrGroupId))
			return errors.New("属性配置不存在")
		}

		baseAttr := &MazeEquipCache.BaseAttrInfo{
			AttrGroup: proto.Int32(attrGroupId),
		}
		showAttrList, realAttrList, randWeight, err := calcEquipAttrValueEx(dei, equipAffixCfg, dei.EParam.TotalScore, dei.EParam.EquipGuid)
		if err != nil {
			dei.ErrorWF("GenBaseAttr calcEquipAttrValueEx fail", zap.Error(err), zap.Any("attrGroupId", attrGroupId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
			return err
		}
		var attrType int32
		if index == 1 {
			attrType = int32(MazeGameEquip.ENUM_EQUIP_ATTR_TYPE_ATTR_MAIN)
		} else {
			attrType = int32(MazeGameEquip.ENUM_EQUIP_ATTR_TYPE_ATTR_BASE)
		}
		baseAttr.ShowAttrList = showAttrList
		baseAttr.RealAttrList = realAttrList
		baseAttr.RandWeight = proto.Int32(int32(randWeight))
		baseAttr.AttrType = proto.Int32(attrType)
		baseAttr.Index = proto.Int32(index)
		baseAttrs = append(baseAttrs, baseAttr)
		dei.UParam.BaseRandoms = append(dei.UParam.BaseRandoms, randWeight)
		remGroupMap[equipAffixCfg.Group] = struct{}{}

		if _, ok := equipInfoCfg.PoolLimitMap[attrGroupId]; ok {
			limitNum++
		}
		if _, ok := equipInfoCfg.PoolLimitBase2Map[attrGroupId]; ok {
			limitBase2Num++
		}
	}
	// fkfmt.Println("GenBaseAttr4")
	dei.DebugWF("GenBaseAttr Instance result",
		zap.Any("baseAttrs", baseAttrs), zap.Int64("equipGuid", dei.EParam.EquipGuid))
	dei.EquipInfo.BaseAttrs = append(dei.EquipInfo.BaseAttrs, baseAttrs...)
	// fkfmt.Println("GenBaseAttr5")
	return nil
}

// 生成初始属性
func (dei *DEInstance) GenRandAttr() error {
	dei.DebugWF("开始实例化Rand属性",
		zap.Int32("equipId", dei.EParam.EquipId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
	equipInfoCfg := dei.EParam.EquipInfoCfg
	dei.DebugWF("开始随机Rand属性条数", zap.Int64("equipGuid", dei.EParam.EquipGuid))
	randNum := randfuncs.RandByWeightV2(dei, equipInfoCfg.AffixRandNum, true)
	dei.DebugWF("Rand属性条数随机结果", zap.Any("randNum", randNum), zap.Int64("equipGuid", dei.EParam.EquipGuid))
	if randNum == 0 {
		return nil
	}
	// poolLimitMap := GMazeEquipAffixRandPoolV8CfgEx.GetPoolLimitCfg(dei.EParam.EquipInfoCfg.LimitAttr)
	// poolLimitBase2Map := GMazeEquipAffixRandPoolV8CfgEx.GetPoolLimitCfg(dei.EParam.EquipInfoCfg.LimitBase2Attr)
	var limitNum, limitBase2Num int32
	for _, attrInfo := range dei.EquipInfo.BaseAttrs {
		if _, ok := equipInfoCfg.PoolLimitMap[attrInfo.GetAttrGroup()]; ok {
			limitNum++
		}
		if _, ok := equipInfoCfg.PoolLimitBase2Map[attrInfo.GetAttrGroup()]; ok {
			limitBase2Num++
		}
	}

	randPoolMap := make(map[int32]int32)
	for k, v := range equipInfoCfg.AffixRandPool {
		randPoolMap[k] = v
	}
	randAttrs := make([]*MazeEquipCache.BaseAttrInfo, 0)
	dei.DebugWF("开始随机Rand属性池", zap.Int64("equipGuid", dei.EParam.EquipGuid))
	poolId := randfuncs.RandByWeightV2(dei, randPoolMap, true)
	dei.DebugWF("随机Rand属性池结果", zap.Any("poolId", poolId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
	if poolId == 0 {
		dei.ErrorWF("GenRandAttr get AffixRandPool weight fail",
			zap.Any("randPoolMap", randPoolMap),
			zap.Int32("equipId", dei.EParam.EquipId))
		return errors.New("属性池不能为0")
	}
	poolWeightMap := mazeequipaffixrandpoolv8.GetEquipPoolWeightCfg(poolId)
	remGroupMap := make(map[int32]struct{}, 0)
	for index := int32(1); index <= randNum; index++ {
		newPoolWeightMap := make(map[int32]int32, 0)
		for k, v := range poolWeightMap {
			if _, ok := remGroupMap[v.Group]; ok {
				dei.DebugWF("随机Base属性词条去重", zap.Any("去重词条id", k), zap.Any("去重组id", v.Group), zap.Int64("equipGuid", dei.EParam.EquipGuid))
				continue
			}
			if limitNum >= dei.EParam.EquipInfoCfg.LimitNum {
				if _, ok := equipInfoCfg.PoolLimitMap[v.Affix_id]; ok {
					dei.DebugWF("随机Base属性抗性词条去重", zap.Any("去重词条id", k), zap.Any("去重组id", v.Group), zap.Int64("equipGuid", dei.EParam.EquipGuid))
					continue
				}
			}
			if limitBase2Num >= dei.EParam.EquipInfoCfg.LimitBase2Num {
				if _, ok := equipInfoCfg.PoolLimitBase2Map[v.Affix_id]; ok {
					dei.DebugWF("随机Base属性垃圾词条去重", zap.Any("去重词条id", k), zap.Any("去重组id", v.Group), zap.Int64("equipGuid", dei.EParam.EquipGuid))
					continue
				}
			}
			newPoolWeightMap[k] = v.Weight
		}

		if len(newPoolWeightMap) == 0 {
			dei.ErrorWF("GenRandAttr get GMazeEquipAffixRandPoolV8Cfg weight fail",
				zap.Int32("poolId", poolId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
			return errors.New("属性配置不存在")
		}

		dei.DebugWF("开始随机Rand属性词条id", zap.Any("条数", index), zap.Int64("equipGuid", dei.EParam.EquipGuid))
		attrGroupId := randfuncs.RandByWeightV2(dei, newPoolWeightMap, true)
		dei.DebugWF("随机Rand属性词条id结果", zap.Any("attrGroupId", attrGroupId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
		equipAffixCfg := GMazeEquipAffixRandPoolV8Cfg.Get(attrGroupId)
		if equipAffixCfg == nil {
			dei.ErrorWF("GenRandAttr GMazeEquipAffixRandPoolV8Cfg fail",
				zap.Int32("poolId", poolId),
				zap.Any("attrGroupId", attrGroupId),
				zap.Int64("equipGuid", dei.EParam.EquipGuid))
			return errors.New("属性配置不存在")
		}

		randAttr := &MazeEquipCache.BaseAttrInfo{
			AttrGroup: proto.Int32(attrGroupId),
		}
		showAttrList, realAttrList, randWeight, err := calcEquipAttrValueEx(dei, equipAffixCfg, dei.EParam.TotalScore, dei.EParam.EquipGuid)
		if err != nil {
			dei.ErrorWF("GenRandAttr calcEquipAttrValueEx fail", zap.Error(err), zap.Any("attrGroupId", attrGroupId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
			return err
		}
		randAttr.ShowAttrList = showAttrList
		randAttr.RealAttrList = realAttrList
		randAttr.RandWeight = proto.Int32(int32(randWeight))
		randAttr.AttrType = proto.Int32(int32(MazeGameEquip.ENUM_EQUIP_ATTR_TYPE_ATTR_RAND))
		randAttr.Index = proto.Int32(index + int32(len(dei.EquipInfo.BaseAttrs)))
		randAttrs = append(randAttrs, randAttr)
		dei.UParam.BaseRandoms = append(dei.UParam.BaseRandoms, randWeight)
		remGroupMap[equipAffixCfg.Group] = struct{}{}
		if _, ok := equipInfoCfg.PoolLimitMap[attrGroupId]; ok {
			limitNum++
		}
		if _, ok := equipInfoCfg.PoolLimitBase2Map[attrGroupId]; ok {
			limitBase2Num++
		}
	}
	dei.DebugWF("GenRandAttr Instance result",
		zap.Int32("randNum", randNum),
		zap.Any("randAttr", randAttrs),
		zap.Int64("equipGuid", dei.EParam.EquipGuid))
	dei.EquipInfo.BaseAttrs = append(dei.EquipInfo.BaseAttrs, randAttrs...)
	return nil
}

// 随机武器类型
func (dei *DEInstance) GenEquipSubType() error {
	dei.DebugWF("开始随机装备类型",
		zap.Int32("equipId", dei.EParam.EquipId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
	equipInfoCfg := dei.EParam.EquipInfoCfg
	equipTypeWeight := make(map[int32]int32, 0)
	// subTypeRemMap, err := GetEquipSubTypeRemMap(dei.EquipInfo.GetSuitId())
	for equipType, weight := range equipInfoCfg.SubTypeRandom {
		if weight == 0 {
			continue
		}
		if dei.EParam.AssignEquipType > 0 {
			if equipType != dei.EParam.AssignEquipType {
				continue
			}
		}
		equipTypeWeight[equipType] = weight
	}
	if len(equipTypeWeight) <= 0 {
		dei.ErrorWF("GenEquipType get GMazeEquipInfoV8Cfg fail",
			zap.Int32("equipId", dei.EParam.EquipId),
			zap.Int32("AssignEquipType", dei.EParam.AssignEquipType),
			zap.Any("SubTypeRandom", equipInfoCfg.SubTypeRandom),
			zap.Int64("equipGuid", dei.EParam.EquipGuid))
		return errors.New("套装类型配置不存在")
	}
	equipType := randfuncs.RandByWeightV2(dei, equipTypeWeight, true)
	dei.DebugWF("随机装备类型结果", zap.Any("equipType", equipType), zap.Int64("equipGuid", dei.EParam.EquipGuid))
	if equipType == 0 {
		return nil
	}
	dei.EquipInfo.EquipSubType = proto.Int32(equipType)
	dei.DebugWF("GenEquipType Instance result",
		zap.Any("equipType", equipType), zap.Int64("equipGuid", dei.EParam.EquipGuid))
	return nil
}

// 生成套装Id
func (dei *DEInstance) GenSuiteId() error {
	dei.DebugWF("开始随机套装Id",
		zap.Int32("equipId", dei.EParam.EquipId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
	equipInfoCfg := dei.EParam.EquipInfoCfg
	suiteWeight := make(map[int32]int32, 0)
	for suitId, weight := range equipInfoCfg.SuiteId {
		if weight == 0 {
			continue
		}
		if dei.EParam.HasSuit > 0 {
			if suitId == 0 {
				continue
			}
			if dei.EParam.AssignSuit > 0 {
				if suitId != dei.EParam.AssignSuit {
					continue
				}
			}
		}
		suiteWeight[suitId] = weight
	}
	if len(suiteWeight) <= 0 {
		if dei.EParam.HasSuit > 0 {
			dei.ErrorWF("GenSuiteId get GMazeEquipInfoV8Cfg fail",
				zap.Int32("equipId", dei.EParam.EquipId),
				zap.Int64("equipGuid", dei.EParam.EquipGuid),
				zap.Int32("AssignSuit", dei.EParam.AssignSuit),
				zap.Int32("HasSuit", dei.EParam.HasSuit))
			return errors.New("套装属性配置不存在")
		}
		return nil
	}
	suitId := randfuncs.RandByWeightV2(dei, suiteWeight, true)
	dei.DebugWF("随机套装Id结果", zap.Any("suitId", suitId))
	if suitId == 0 {
		return nil
	}
	dei.EquipInfo.SuitId = proto.Int32(suitId)
	dei.InfoWF("GenSuiteId Instance result",
		zap.Any("suitId", suitId), zap.Int64("equipGuid", dei.EParam.EquipGuid))
	return nil
}

func (dei *DEInstance) BeforeRecord() {
	dei.Record.EquipId = dei.EquipInfo.GetEquipId()
	dei.Record.EquipGuid = dei.EquipInfo.GetEquipGuid()
	dei.Record.TotalScore = dei.EParam.TotalScore
	dei.Record.SentencesLibId = dei.EParam.EquipInfoCfg.ScoreGroup
	dei.Record.RuleId = dei.EquipInfo.GetRuleId()
	dei.Record.SuitId = dei.EquipInfo.GetSuitId()
	dei.Record.EquipSubType = dei.EquipInfo.GetEquipSubType()
	baseAttrs := []string{}
	for _, baseAttr := range dei.EquipInfo.BaseAttrs {
		key := fmt.Sprintf("%d:%d:%d:%d", baseAttr.GetIndex(), baseAttr.GetAttrType(), baseAttr.GetAttrGroup(), baseAttr.GetRandWeight())
		keyAttr := "<"
		for _, attrInfo := range baseAttr.ShowAttrList {
			keyAttr += fmt.Sprintf("%d:%d", attrInfo.GetAttrId(), attrInfo.GetAttrValue())
		}
		keyAttr += ">"
		keyAttr += "<"
		for i, attrInfo := range baseAttr.RealAttrList {
			keyAttr += fmt.Sprintf("%d:%d", attrInfo.GetAttrId(), attrInfo.GetAttrValue())
			if len(baseAttr.RealAttrList) != i+1 {
				keyAttr += ","
			}
		}
		keyAttr += ">"
		key += keyAttr
		baseAttrs = append(baseAttrs, key)
	}
	dei.Record.BaseAttrs = strings.Join(baseAttrs, "|")
}

func initStageInfo(logger fklog.FKLogI, uparam *DEIUParam) error {
	row := GMazeEquipAttrStageV8Cfg.Get(uparam.ChargeStage)
	if row == nil {
		logger.ErrorWF("initStageInfo  GDollEquipAttrStageV8Cfg error", zap.Int32("stageId", uparam.ChargeStage))
		return errors.New("配置不存在")
	}
	uparam.StageDropScore = row.Score
	uparam.StageFactor = row.Ratio
	uparam.BaseRandoms = make([]int64, 0)
	return nil
}

func calcEquipAttrValueEx(logger fklog.FKLogI, cfg *GMazeEquipAffixRandPoolV8Cfg.MazeEquipAffixRandPoolV8ConfigRow, totalScore int32, equipGuid int64) (showAttrList, realAttrList []*MazeEquipCache.EquipAttrInfo, randWeight int64, err error) {
	rollWeightMap := mazeequipaffixrolltypev8.GetMazeEquipRollCfgByRollTypeAndScore(cfg.Roll_type, totalScore)

	if len(rollWeightMap) <= 0 {
		logger.ErrorWF("calcEquipAttrValueEx GetDollEquipRollCfgByRollTypeAndScore err",
			zap.Int32("rollType", cfg.Roll_type),
			zap.Any("totalScore", totalScore))
		return nil, nil, 0, errors.New("配置不存在")
	}
	logger.DebugWF("开始随机roll值", zap.Any("roll", cfg.Roll_type), zap.Any("分数", totalScore), zap.Int64("equipGuid", equipGuid))
	rollId := randfuncs.RandByWeightV2(logger, rollWeightMap, true)
	logger.DebugWF("随机roll类型结果", zap.Any("rollId", rollId), zap.Int64("equipGuid", equipGuid))
	rollTypeCfg := GMazeEquipAffixRollTypeV8Cfg.Get(rollId)
	if rollTypeCfg == nil {
		logger.ErrorWF("calcEquipAttrValueEx GMazeEquipAffixRollTypeV8Cfg err",
			zap.Int32("rollId", rollId),
			zap.Any("rollWeightMap", rollWeightMap))
		return nil, nil, 0, errors.New("配置不存在")
	}
	rollCount := int64(fkutil.RandInt32(int(rollTypeCfg.Roll_range_min), int(rollTypeCfg.Roll_range_max+1)))
	showAttrList = make([]*MazeEquipCache.EquipAttrInfo, 0)
	calcAddCount, calcRangeCount := CalcAttrRealRoll(logger, cfg.Add_attr_min, cfg.Add_attr_max, cfg.Show_attr_min, cfg.Show_attr_max, rollCount, rollTypeCfg.Round_value)
	for attrId, attrMax := range cfg.Show_attr_max {
		if attrId == 0 {
			logger.ErrorWF("calcEquipAttrValueEx GMazeEquipAffixRandPoolV8Cfg err",
				zap.Int32("attrGroupId", cfg.Affix_id),
				zap.Any("ShowAttrMax", cfg.Show_attr_max))
			return nil, nil, 0, errors.New("配置不存在")
		}
		attrMin := cfg.Show_attr_min[attrId]
		// todo 如果max小于min，配置错误报错
		if attrMax < attrMin {
			logger.ErrorWF("calcEquipAttrValueEx GMazeEquipAffixRandPoolV8Cfg err",
				zap.Int32("attrGroupId", cfg.Affix_id),
				zap.Any("ShowAttrMax", cfg.Show_attr_max),
				zap.Any("ShowAttrMax", cfg.Show_attr_min))
			return nil, nil, 0, errors.New("属性配置异常")
		}
		showAttrInfo := packAttrInfoDb(attrId, attrMin, attrMax, rollTypeCfg.Round_value, calcAddCount, calcRangeCount)
		showAttrList = append(showAttrList, showAttrInfo)
		break
	}
	realAttrList = make([]*MazeEquipCache.EquipAttrInfo, 0)
	for attrId, attrMax := range cfg.Add_attr_max {
		if attrId == 0 {
			logger.ErrorWF("calcEquipAttrValueEx GMazeEquipAffixRandPoolV8Cfg err",
				zap.Int32("attrGroupId", cfg.Affix_id),
				zap.Any("AddAttrMax", cfg.Add_attr_max))
			return nil, nil, 0, errors.New("配置不存在")
		}
		attrMin := cfg.Add_attr_min[attrId]
		if attrMax < attrMin {
			logger.ErrorWF("calcEquipAttrValueEx GMazeEquipAffixRandPoolV8Cfg err",
				zap.Int32("attrGroupId", cfg.Affix_id),
				zap.Any("ShowAttrMax", cfg.Add_attr_max),
				zap.Any("ShowAttrMax", cfg.Add_attr_min))
			return nil, nil, 0, errors.New("属性配置异常")
		}
		realAttrList = append(realAttrList, packAttrInfoDb(attrId, attrMin, attrMax, rollTypeCfg.Round_value, calcAddCount, calcRangeCount))
	}
	logger.InfoWF("calcEquipAttrValueEx result", zap.Int64("equipGuid", equipGuid), zap.Any("affixId", cfg.Affix_id), zap.Any("showAttr", showAttrList), zap.Any("realAttr", realAttrList), zap.Any("rollCount", rollCount), zap.Any("calcAdd", calcAddCount), zap.Any("calcRange", calcRangeCount))
	return showAttrList, realAttrList, rollCount, nil
}

func InitEquipInfoCfg(equipId int32, affixBaseNum, affixBasePool, affixRandNum, affixRandPool, suiteId, subTypeRandom map[int32]int32, scoreGroup int32) *DollEquipInfoCfg {
	equipInfoCfg := &DollEquipInfoCfg{
		EquipId:       equipId,
		AffixBaseNum:  make(map[int32]int32),
		AffixBasePool: make(map[int32]int32),
		AffixRandNum:  make(map[int32]int32),
		AffixRandPool: make(map[int32]int32),
		SuiteId:       make(map[int32]int32),
		SubTypeRandom: make(map[int32]int32),
		ScoreGroup:    scoreGroup,
	}
	for k, v := range affixBaseNum {
		equipInfoCfg.AffixBaseNum[k] = v
	}
	for k, v := range affixBasePool {
		equipInfoCfg.AffixBasePool[k] = v
	}
	for k, v := range affixRandNum {
		equipInfoCfg.AffixRandNum[k] = v
	}
	for k, v := range affixRandPool {
		equipInfoCfg.AffixRandPool[k] = v
	}
	for k, v := range suiteId {
		equipInfoCfg.SuiteId[k] = v
	}
	for k, v := range subTypeRandom {
		equipInfoCfg.SubTypeRandom[k] = v
	}
	return equipInfoCfg
}

func CalcAttrValByRoll(attrMin, attrMax int64, round int32, calcAddCount, calcRangeCount int64) int64 {
	addVal := (attrMax - attrMin) * calcAddCount / calcRangeCount
	remainValue := addVal % int64(round)
	addVal = addVal - remainValue
	return attrMin + addVal
}

func packAttrInfoDb(attrId int32, attrMin, attrMax int64, round int32, calcAddCount, calcRangeCount int64) *MazeEquipCache.EquipAttrInfo {
	attrInfo := &MazeEquipCache.EquipAttrInfo{}
	attrInfo.AttrId = proto.Int32(attrId)
	if attrMax-attrMin <= 0 {
		attrInfo.AttrValue = proto.Int64(attrMax)
		return attrInfo
	}
	attrInfo.AttrValue = proto.Int64(CalcAttrValByRoll(attrMin, attrMax, round, calcAddCount, calcRangeCount))
	return attrInfo
}

func CalcAttrRealRoll(logger fklog.FKLogI, addAttrMin, addAttrMax, showAttrMin, showAttrMax map[int32]int64, rollCount int64, round int32) (int64, int64) {
	var minAttrId int32
	var minRangeCount, minAttrCount, maxAttrCount int64
	for attrId, attrMaxVal := range addAttrMax {
		rangeCount := attrMaxVal - addAttrMin[attrId]
		if rangeCount <= 0 {
			continue
		}
		if minAttrId == 0 {
			minAttrId = attrId
			minRangeCount = rangeCount
			minAttrCount = addAttrMin[attrId]
			maxAttrCount = attrMaxVal
			continue
		}
		if minRangeCount > rangeCount {
			minAttrId = attrId
			minRangeCount = rangeCount
			minAttrCount = addAttrMin[attrId]
			maxAttrCount = attrMaxVal
		}
	}
	for attrId, attrMaxVal := range showAttrMax {
		rangeCount := attrMaxVal - showAttrMin[attrId]
		if rangeCount <= 0 {
			continue
		}
		if minAttrId == 0 {
			minAttrId = attrId
			minRangeCount = rangeCount
			minAttrCount = showAttrMin[attrId]
			maxAttrCount = attrMaxVal
			continue
		}
		if minRangeCount > rangeCount {
			minAttrId = attrId
			minRangeCount = rangeCount
			minAttrCount = showAttrMin[attrId]
			maxAttrCount = attrMaxVal
		}
	}
	addAttrCount := CalcRealAttrVal(minAttrCount, maxAttrCount, rollCount, round)
	// logger.DebugWF("CalcAttrRealRoll end", zap.Any("minAttrId", minAttrId), zap.Any("minAttrCount", minAttrCount), zap.Any("maxAttrCount", maxAttrCount), zap.Any("addAttrCount", addAttrCount))
	return addAttrCount - minAttrCount, maxAttrCount - minAttrCount
}

func CalcRealAttrVal(attrMin, attrMax, rollCount int64, round int32) int64 {
	addVal := (attrMax - attrMin) * rollCount / 10000
	remainValue := addVal % int64(round)
	addVal = addVal - remainValue
	return attrMin + addVal
}
