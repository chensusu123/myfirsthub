/*
 * @Author: majian
 * @Date: 2024-04-26 11:10:18
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-17 18:00:19
 */
package equip

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/common/errors"
	"maze_game_server/common/structsdef"
	"maze_game_server/common/tradeno"
	"maze_game_server/config/GMazeConfigV8Cfg"
	"maze_game_server/config/GMazeEquipConfigV8Cfg"
	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/io/kafka/dollequipassmeblekakfa"
	"maze_game_server/io/redis/dollassembleredis"
	"maze_game_server/io/redis/dollassemblesuitredis"
	"maze_game_server/io/redis/mazeattrcalcnotifyqueue"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/module/assembleidpack"
	"maze_game_server/module/calcassembleattr"
	"maze_game_server/module/dollassembleinfo"
	"maze_game_server/module/effectequip"
	"maze_game_server/module/mazebuffchgrrecordapi"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/server/MazeEquipCache"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/servers/maze_main_server/process/equip/demconstdef"
	"maze_game_server/servers/maze_main_server/process/equip/module"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// redis + 内存缓存
func GetEquipInitState(logger fklog.FKLogI, userId uint64) (state int64, err error) {
	return dollassembleredis.GetDollEquipInitState(logger, userId)
}

func SaveEquipInitState(logger fklog.FKLogI, userId uint64, state int64) error {
	return dollassembleredis.SetDollEquipInitState(logger, userId, state)
}

// 从背包找
func getInitquipFromBag(logger fklog.FKLogI, userId uint64, equipIds map[int32]int64) (equipInfoMap map[int32]*MazeEquipCache.MazeEquipInfoDb, err error) {
	equipMap, err := effectequip.GetAllEffectEquipInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("getInitquipFromBag GetAllEquipInfo error", zap.Error(err))
		return nil, err
	}
	// 找到一个即可,不用guid最小，包括穿在身上的
	equipInfoMap = make(map[int32]*MazeEquipCache.MazeEquipInfoDb, 0)
	for _, v := range equipMap {
		// 去除临时背包装备
		// if v.GetEnterTime() > 0 {
		// 	continue
		// }
		if equipIds[v.GetEquipId()] > 0 {
			equipInfoMap[v.GetEquipId()] = v
		}
	}

	return equipInfoMap, nil
}

func HandleDollEquipInit(ctx context.Context, userId uint64, needNotify bool) error {
	logger := fklog.ContextAppLogger(ctx)
	state, e := GetEquipInitState(logger, userId)
	if e != nil {
		logger.ErrorWF("HandleDollEquipInit get init state fail", zap.Error(e))
		return e
	}
	// 初始化完成
	if state&(constdef.DollEquipInitStateDress|constdef.DollEquipInitStateBag) ==
		constdef.DollEquipInitStateDress|constdef.DollEquipInitStateBag {
		logger.InfoWF("HandleDollEquipInit already init")
		return nil
	}
	// 初始化进行中
	if state&constdef.DollEquipInitDoing > 0 {
		logger.WarnWF("HandleDollEquipInit init doing", zap.Int64("state", state))
		return nil
	}
	cfg := GMazeConfigV8Cfg.GetMazeConfigV8Config(constdef.DollCfgId3301)
	if cfg == nil {
		logger.ErrorWF("HandleDollEquipInit no init equip cfg", zap.Int32("cfgId", constdef.DollCfgId3301))
		return errors.New("初始装备配置不存在")
	}
	initEquips := make(map[int32]int64)
	for k, v := range cfg.Value_map {
		if k != 0 && v != 0 {
			initEquips[k] = v
		}
	}
	if len(initEquips) == 0 {
		return NoEquipInit(logger, userId)
	}
	return doInitDollEquip(ctx, userId, state, initEquips, needNotify)
}

func NoEquipInit(logger fklog.FKLogI, userId uint64) error {
	state := constdef.DollEquipInitDoing | constdef.DollEquipInitStateDress | constdef.DollEquipInitStateBag
	err := SaveEquipInitState(logger, userId, int64(state))
	if err != nil {
		logger.ErrorWF("NoEquipInit init fail", zap.Error(err), zap.Int("state", state))
		return err
	}
	logger.WarnWF("NoEquipInit init ok", zap.Int("state", state))
	return err
}

func doInitDollEquip(ctx context.Context, userId uint64, state int64, equips map[int32]int64, needNotify bool) error {
	logger := fklog.ContextAppLogger(ctx)
	// 先加入背包
	var err error

	initEquipInfoMap, err := getInitquipFromBag(logger, userId, equips)
	if err != nil {
		logger.ErrorWF("doInitDollEquip getInitquipFromBag fail", zap.Error(err), zap.Any("equips", equips),
			zap.Int64("state", state))
		return err
	}
	// 设置开始初始化状态
	state |= constdef.DollEquipInitDoing
	err = SaveEquipInitState(logger, userId, state)
	if err != nil {
		logger.ErrorWF("doInitDollEquip set init doing state fail", zap.Error(err), zap.Any("equips", equips),
			zap.Int64("state", state))
		return err
	}
	var oldState = state

	defer func() {
		if oldState == state {
			logger.WarnWF("doInitDollEquip SaveEquipInitState no chg", zap.Any("equips", equips),
				zap.Any("initEquipInfoMap", initEquipInfoMap),
				zap.Int64("oldState", oldState),
				zap.Int64("state", state))
			return
		}
		err = SaveEquipInitState(logger, userId, state)
		if err != nil {
			logger.ErrorWF("doInitDollEquip SaveEquipInitState fail", zap.Error(err), zap.Any("equips", equips),
				zap.Any("initEquipInfoMap", initEquipInfoMap),
				zap.Int64("oldState", oldState),
				zap.Int64("state", state))
		} else {
			logger.InfoWF("doInitDollEquip SaveEquipInitState succ", zap.Any("equips", equips),
				zap.Any("initEquipInfoMap", initEquipInfoMap),
				zap.Int64("oldState", oldState),
				zap.Int64("state", state))
		}
	}()

	if len(initEquipInfoMap) == 0 {
		if state&constdef.DollEquipInitStateBag == 0 {
			// 添加到背包后，再次查询装备
			addInitEquipToBag(ctx, userId, equips)
			initEquipInfoMap, err = getInitquipFromBag(logger, userId, equips)
		}
	}
	if err != nil {
		logger.ErrorWF("doInitDollEquip add init equip fail", zap.Any("equips", equips),
			zap.Int64("oldState", oldState),
			zap.Int64("state", state))
		return err
	}
	if len(initEquipInfoMap) > 0 {
		state |= constdef.DollEquipInitStateBag
	} else {
		logger.WarnWF("doInitDollEquip no get init equip", zap.Any("equips", equips),
			zap.Int64("oldState", oldState),
			zap.Int64("state", state))
		return nil
	}

	// 穿戴到身上
	err = dressInitEquip(ctx, userId, initEquipInfoMap, needNotify)
	if err != nil {
		return err
	}
	state |= constdef.DollEquipInitStateDress
	return nil
}

// 添加初始化装备到背包
func addInitEquipToBag(ctx context.Context, userId uint64, equips map[int32]int64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	tradeNo := tradeno.GetTradeNum()
	rqAdd := &MazeEquipSvr.SvrAddMazeEquipRQ{
		UserId:      proto.Uint64(userId),
		OpType:      proto.Int32(int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_INIT_EQUIP)),
		TradeNumber: proto.Uint64(tradeNo),
	}
	for k := range equips {
		equipCond := &MazeEquipSvr.SvrEquipInfo{EquipId: proto.Int32(k)}
		// 初始化武器子类型固定是1 刀
		equipCfg := GMazeEquipInfoV8Cfg.GetMazeEquipInfoV8Config(k)
		if equipCfg != nil && equipCfg.Pos == 1 {
			equipCond.Conditions = append(equipCond.Conditions, &MazeEquipSvr.ConditionInfo{
				Id: proto.Int32(int32(MazeEquipSvr.EQUIP_ADD_CONDITION_ASSIGN_EQUIP_SUB_TYPE)), Value: proto.Int64(1)})
		}
		rqAdd.EquipList = append(rqAdd.EquipList, equipCond)
	}

	rsAdd := &MazeEquipSvr.SvrAddMazeEquipRS{}
	// err = dollequipbagrpc.MazeBagAddRQ(logger, rqAdd, rsAdd)
	err = OnSvrAddMazeEquipRQ(ctx, int64(userId), rqAdd, rsAdd, "")
	if err != nil {
		logger.ErrorWF("addInitEquipToBag MazeBagAddRQ fail", zap.Error(err), zap.Any("req", rqAdd), zap.Any("rs", rsAdd))
	} else {
		if rsAdd.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
			logger.ErrorWF("addInitEquipToBag rs fail", zap.Any("req", rqAdd), zap.Any("rs", rsAdd))
		}
	}
	return
}

// 穿戴初始化装备
func dressInitEquip(ctx context.Context, userId uint64, equipInfoMap map[int32]*MazeEquipCache.MazeEquipInfoDb, needNotify bool) error {
	logger := fklog.ContextAppLogger(ctx)
	var (
		chgEquipPosList []*MazeEquipCache.MazeEquipPosInfo
		recordList      []*dollequipassmeblekakfa.MazeGameEquipAssembleRecord
	)
	qualityMap := make(map[int32]int32, 0)
	constCfgRow := GMazeEquipConfigV8Cfg.GetMazeEquipConfigV8Config(1)
	if constCfgRow == nil {
		logger.ErrorWF("dressInitEquip no found Getmazeequipconfigv8Config")
		return errors.New("未找到装备通用配置")
	}

	assembleInfo, oldEffect, err := dollassembleinfo.GetDollAssembleInfoEx(logger, userId)
	if err != nil {
		logger.ErrorWF("dressInitEquip Get Assemble info fail", zap.Error(err))
		return err
	}
	if assembleInfo.GetCurSuitIndex() == 0 {
		logger.WarnWF("dressInitEquip no set cur suit")
		err = errors.New("未设置生效装备套")
		return err
	}

	for _, equipInfo := range equipInfoMap {
		row := GMazeEquipInfoV8Cfg.GetMazeEquipInfoV8Config(equipInfo.GetEquipId())
		if row == nil {
			logger.ErrorWF("dressInitEquip 未找到装备配置 ", zap.Int32("equipId", equipInfo.GetEquipId()),
				zap.Int64("equipGuid", equipInfo.GetEquipGuid()))
			return errors.New("未找到装备配置")
		}
		equipPosInfo := module.GetEquipPosInfo(assembleInfo, row.Pos)
		if equipPosInfo == nil {
			logger.ErrorWF("dressInitEquip 装备位未解锁", zap.Int32("pos", row.Pos))
			return errors.New("装备位未解锁")
		}
		guid := module.GetDressedGuid(equipPosInfo)
		if guid > 0 {
			if equipPosInfo.GetEquipLoadInfo().GetEquipId() == equipInfo.GetEquipId() {
				logger.WarnWF("dressInitEquip has dressed equip", zap.Int64("guid", guid), zap.Int32("equipId", equipPosInfo.GetEquipLoadInfo().GetEquipId()))
			} else {
				logger.WarnWF("dressInitEquip has dressed other equip", zap.Int64("guid", guid), zap.Int32("equipId", equipPosInfo.GetEquipLoadInfo().GetEquipId()))
			}
			continue
		}

		// 记录旧值
		oldEquipPos := &MazeEquipCache.MazeEquipPosInfo{}
		warpEquipPosOld := &MazeEquipCache.MazeEquipPosDb{
			EquipGuid: proto.Int64(equipPosInfo.GetEquipLoadInfo().GetEquipGuid()),
			EquipId:   proto.Int32(equipPosInfo.GetEquipLoadInfo().GetEquipId())}
		oldEquipPos.EquipLoadInfo = warpEquipPosOld
		oldEquipPos.EquipInfo = equipPosInfo.GetEquipInfo()

		module.ReplaceEquip(logger, equipPosInfo, equipInfo)

		// 记录流水
		recordType := dollequipassmeblekakfa.DollEquipAssembleOpInit
		record := StartEquipAssmebleRecord(userId, row.Pos, recordType, equipPosInfo, oldEquipPos, oldEffect)
		recordList = append(recordList, record)
		chgEquipPosList = append(chgEquipPosList, equipPosInfo)

		qualityMap[row.Quality] = 0
	}

	var opCode int32 = 1
	var opMask int32
	var effect *calcassembleattr.EquipmentEffectInfo
	defer func() {
		for _, record := range recordList {
			tmpMask := opMask
			if opCode == 0 {
				// 通知背包服务
				e := NotifyBagSvr(logger, userId, record.EquipPos, 0, int64(record.NewGuid), 0)
				if e != nil {
					tmpMask |= demconstdef.DollEquipAssembleOpMaskNotifyBag
				}
			}
			// 记录流水
			EndEquipAssmebleRecord(ctx, record, opCode, tmpMask, effect)
		}
	}()
	// 保存装配数据
	err = dollassemblesuitredis.SaveEquipAssembleInfoV2(logger, userId, assembleInfo.GetCurSuitIndex(), chgEquipPosList)
	if err != nil {
		opMask |= demconstdef.DollEquipAssembleOpMaskDbSave
		logger.ErrorWF("dressInitEquip SaveEquipAssembleInfoV2 fail", zap.Error(err), zap.Any("chgEquipPosList", chgEquipPosList))
		return err
	}
	opCode = 0
	var e error
	effect, e = calcassembleattr.CalcEquipEffect(logger, assembleInfo.MazeEquips)
	if e != nil {
		opMask |= demconstdef.DollEquipAssembleOpMaskCalcBuff
	} else {
		e = mazebuffinforedis.SaveMazeEquipBuff(logger, userId, effect.Other)
		if e != nil {
			opMask |= demconstdef.DollEquipAssembleOpMaskNonForceBuff
		} else {
			// 通知计算属性
			calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{}
			calcAttrNotify.FromServer = demconstdef.MySvr
			calcAttrNotify.UserId = userId
			calcAttrNotify.ChgType = constdef.MazeBuffChgTypeEquipInit
			calcAttrNotify.BuffSrc = constdef.MazeBuffSrcEquip
			e = mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(ctx, calcAttrNotify)
			if e != nil {
				opMask |= demconstdef.DollEquipAssembleOpMaskCalcAttr
				logger.ErrorWF("dressInitEquip SendDollAttrCalcNotify fail", zap.Error(e))
			}

			mazebuffchgrrecordapi.SendMazeBuffChgRecord(ctx, userId,
				constdef.MazeBuffSrcEquip,
				constdef.MazeBuffChgTypeEquipInit,
				oldEffect.Other, effect.Other)
		}
	}
	if needNotify {
		chgAssembleInfo := &MazeEquipCache.MazeAssembleDb{}
		var chgMask int32
		if assembleInfo.GetEpSuitId() != effect.SuitId {
			chgAssembleInfo.EpSuitId = proto.Int32(effect.SuitId)
			chgMask |= int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_SUIT_MASK)
		}

		chgAssembleInfo.MazeEquips = chgEquipPosList
		chgMask |= int32(MazeGameEquip.ENUM_MAZE_ASSEMBLE_CHG_TYPE_MASK_EQUIP_POS_MASK)
		assembleidpack.SendAssembleChgID(logger, userId, chgAssembleInfo,
			chgMask, int32(int32(MazeGameEquip.ENUM_MAZE_EQUIP_POS_MASK_LOAD_EQUIP_INFO)),
			constdef.DollAssembleChgTypeReplaceEquip)
	}

	return nil
}

// 设置初始装备套序号
func InitDollEquipSuitSeq(logger fklog.FKLogI, userId uint64) error {
	// 初始装备套检查
	r, err := dollassembleredis.GetDollAssembleMetaInfo(logger, userId, constdef.AssemblePrefixCurAssembleSuitIndex)
	if err != nil {
		return err
	}
	if r.GetCurSuitIndex() > 0 {
	} else {
		dm := &dollassembleredis.DollAssembleMetaSt{}
		initSuitIndex := constdef.CurSuitDef
		switchTime := time.Now().Unix()
		dm.CurAssmebleSuitIndex = &initSuitIndex
		dm.LastSwitchSuitTime = &switchTime
		err = dollassembleredis.SetDollAssmebleMetaInfo(logger, userId, dm)
		if err != nil {
			return err
		}
	}
	return err
}

// 人偶属性初始化
func HandleDollAttrInit(ctx context.Context, userId uint64, session string) {
	logger := fklog.ContextAppLogger(ctx)
	slen, err := mazecalcattrredis.HlenMazeCalcAttr(logger, userId)
	if err != nil {
		return
	}
	if slen > 0 {
		// 已经初始化
		return
	}

	// 通知计算属性
	calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{}
	calcAttrNotify.FromServer = demconstdef.MySvr
	calcAttrNotify.UserId = userId
	calcAttrNotify.ChgType = constdef.MazeBuffChgTypeEquipInit
	calcAttrNotify.Session = session
	_ = mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(ctx, calcAttrNotify)
}
