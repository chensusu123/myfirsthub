package equip

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"maze_game_server/pb/server/MazeEquipCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"maze_game_server/config/GMazeEquipAttrStageV8Cfg"
	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/pb/server/MazeEquipSvr"

	"context"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"go.uber.org/zap"
	"maze_game_server/common/function/packtopb"
	"maze_game_server/excel/mazeequipaffixrandpoolv8"
	"maze_game_server/excel/mazeequipconfigv8"
	"maze_game_server/io/kafka/mazeequipbagrecord"
	"maze_game_server/io/kafka/mazeequipinstancerecord"
	"maze_game_server/io/redis/mazeequipgetnumredis"
	"maze_game_server/io/redis/mazeequipguidredis"
	"maze_game_server/module/bagmodule"
	"maze_game_server/pb/common/MazeGameEquip"
	"maze_game_server/pb/errors"
)

func OnSvrAddMazeEquipRQ(ctx fklog.FKLogI, shardingID int64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	defer fkprometheus.DebugPMT("OnSvrAddMazeEquipRQ")()
	userCtx := fkserver.NewUserContext(context.TODO(), uint64(shardingID), ctx)
	req := rqMsg.(*MazeEquipSvr.SvrAddMazeEquipRQ)
	res := rsMsg.(*MazeEquipSvr.SvrAddMazeEquipRS)
	res.ErrInfo = errors.NO_ERROR
	userCtx.WarnWF("OnSvrAddMazeEquipRQ with", zap.Any("rq", req))

	addStartTime := time.Now()
	defer func() {
		costTime := time.Since(addStartTime).Seconds()
		userCtx.WarnWF("OnSvrAddMazeEquipRQ end ", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
		if costTime >= 0.5 {
			ctx.ErrorWF("OnSvrAddMazeEquipRQ timeout", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
		}
	}()

	lock := globalLock.LockOp(uint64(shardingID))
	defer lock.Unlock()

	if len(req.EquipList) <= 0 {
		userCtx.ErrorWF("OnSvrAddMazeEquipRQ none equip", zap.Any("req", req))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("没有要添加的装备信息")
		return err
	}
	if req.GetOpType() <= 0 || req.GetTradeNumber() <= 0 {
		userCtx.ErrorWF("OnSvrAddMazeEquipRQ input invalid args", zap.Any("req", req))
		res.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
		return err
	}

	bagEquipMgr := bagmodule.NewBagEquipMgr(ctx, uint64(shardingID))
	err = bagEquipMgr.LoadBagFromRedis()
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		ctx.ErrorWF("OnSvrAddMazeEquipRQ LoadBagFromRedis fail", zap.Error(err))
		return err
	}

	stageAddition, err := GetUserEquipAddition(userCtx, userCtx.UserID)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		ctx.ErrorWF("OnSvrAddMazeEquipRQ GetUserEquipAddition fail", zap.Error(err))
		return
	}

	attrStageRow := GMazeEquipAttrStageV8Cfg.Get(stageAddition)
	if attrStageRow == nil {
		res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		ctx.ErrorWF("OnSvrAddMazeEquipRQ GMazeEquipAttrStageV8Cfg fail", zap.Int32("tap", stageAddition))
		return
	}
	totalScoreMap, err := GetTotalScoreAndBarrierMap(userCtx, userCtx.UserID, req.EquipList)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		ctx.ErrorWF("OnSvrAddMazeEquipRQ GetTotalScoreMap fail", zap.Error(err))
		return
	}

	allotGuids, err := AddEquipAllotGuid(userCtx, userCtx.UserID, int32(len(req.EquipList)))
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		ctx.ErrorWF("OnSvrAddMazeEquipRQ AddEquipAllotGuid fail", zap.Error(err))
		return
	}
	// 检查guid分配数量是否充足
	if len(allotGuids) != len(req.EquipList) {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("alloc guid fail")
		ctx.ErrorWF("OnSvrAddMazeEquipRQ alloc guid less",
			zap.Int("need", len(req.EquipList)), zap.Int("alloc", len(allotGuids)))
		return
	}
	// 装备实例化开始
	up := &DEIUWParam{}
	up.ChargeStage = stageAddition
	up.OpType = req.GetOpType()
	up.TradeNum = req.GetTradeNumber()
	// 装备获得次数分值
	up.PoolLimitMap = GetPoolLimitMap(userCtx, mazeequipconfigv8.GetEquipGoodAttrIds())
	up.PoolLimitBase2Map = GetPoolLimitMap(userCtx, mazeequipconfigv8.GetEquipBadAttrIds())

	equipInstanceRecordMap := make(map[int64]*MazeGameEquipInstanceRecord, 0)
	startTime := time.Now().Unix()
	addInstanceEquips := make([]*MazeEquipCache.MazeEquipInfoDb, 0) // 实例化的装备
	var wg sync.WaitGroup
	result := make([]*DEInstance, 0)
	now := time.Now()
	equipLock := &sync.Mutex{}
	for i, equipInfo := range req.EquipList {
		equipCfg := GMazeEquipInfoV8Cfg.Get(equipInfo.GetEquipId())
		if equipCfg != nil {
			totalScoreMap[equipCfg.Score_group] += attrStageRow.Score
		} else {
			res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
			ctx.ErrorWF("OnSvrAddMazeEquipRQ less equip cfg",
				zap.Int32("equipId", equipInfo.GetEquipId()))
			return
		}
		var condParam string
		for _, condition := range equipInfo.GetConditions() {
			if condition.GetValue() <= 0 {
				continue
			}
			condParam += fmt.Sprintf("%d:%d", condition.GetId(), condition.GetValue())
		}
		// 需要先分配好 guid、分值
		// 把condtion 拼接好，记录流水用，在内部赋值了
		// 另一种是按分值组走，不同分值组并行，同分值组串行,这种是每件装备都并行
		ep := NewDEIEWParam(equipInfo)
		ep.Guid = allotGuids[i]
		ep.Score = totalScoreMap[equipCfg.Score_group]
		ep.Condition = condParam

		wg.Add(1)
		go func(logger fklog.FKLogI, userId uint64, uwparam *DEIUWParam, ewparam *DEIEWParam) {
			defer wg.Done()
			r, e := DoInsEquip(logger, userId, uwparam, ewparam)
			if e == nil {
				equipLock.Lock()
				result = append(result, r)
				equipLock.Unlock()
			}
		}(userCtx, userCtx.UserID, up, ep)
	}
	wg.Wait()
	userCtx.InfoWF("OnSvrAddMazeEquipRQ total", zap.Duration("cost", time.Since(now)))
	if len(result) != len(req.EquipList) {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("insEquip fail")
		ctx.ErrorWF("OnSvrAddMazeEquipRQ insEquip fail",
			zap.Any("result", len(result)), zap.Any("equipList", len(req.EquipList)))
		return
	}

	for _, dei := range result {
		bagEquipMgr.AddBagAdd(dei.GetResult().EquipInfo)
		addInstanceEquips = append(addInstanceEquips, dei.GetResult().EquipInfo)
		equipInstanceRecordMap[dei.Record.EquipGuid] = dei.Record
	}

	err = bagEquipMgr.SaveBagInfoToRedis()
	if err != nil {
		userCtx.ErrorWF("OnSvrAddMazeEquipRQ BatchSaveEquipInfo error!", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		PushMazeEquipInstanceLog(userCtx, equipInstanceRecordMap, mazeequipinstancerecord.MazeAddEquip, 1)
		PushMazeEquipBagLogEx(userCtx, userCtx.UserID, addInstanceEquips, nil, req.GetTradeNumber(), req.GetOpType(), mazeequipbagrecord.MazeAddEquip, startTime, 1)
		return
	}
	PushMazeEquipBagLogEx(userCtx, userCtx.UserID, addInstanceEquips, nil, req.GetTradeNumber(), req.GetOpType(), mazeequipbagrecord.MazeAddEquip, startTime, 0)

	var isFail int32
	if len(totalScoreMap) > 0 {
		err1 := mazeequipgetnumredis.BatchSetEquipGetNum(userCtx, userCtx.UserID, totalScoreMap)
		if err1 != nil {
			userCtx.ErrorWF("OnSvrAddMazeEquipRQ BatchSetEquipGetNum error!",
				zap.Error(err1),
				zap.Any("totalScoreMap", totalScoreMap))
			isFail += 100
		}
	}
	if len(equipInstanceRecordMap) > 0 {
		PushMazeEquipInstanceLog(userCtx, equipInstanceRecordMap, mazeequipinstancerecord.MazeAddEquip, isFail)
	}

	//	equipCliList := make([]*MazeGameEquip.MazeEquipInfo, 0)
	bagEquipCliList := make([]*MazeGameEquip.MazeEquipInfo, 0)
	equipSvrList := make([]*MazeEquipSvr.MazeEquipInfoSvr, 0)
	for _, equipInfo := range addInstanceEquips {
		equipCli, equipResId, _ := packtopb.EquipInfoToCliPBEx(userCtx, equipInfo)
		// equipCli.ForceValue = proto.Int64(forceMap[equipInfo.GetEquipGuid()])
		equipSvrList = append(equipSvrList, &MazeEquipSvr.MazeEquipInfoSvr{
			EquipId:    proto.Int32(equipCli.GetEquipId()),
			EquipGuid:  proto.Int64(equipCli.GetEquipGuid()),
			EquipName:  proto.String(equipCli.GetEquipName()),
			EquipResId: proto.Int32(equipResId),
			SuitId:     proto.Int32(equipCli.GetSuitInfo().GetSuitId()),
			Pos:        proto.Int32(equipCli.GetPos()),
		})
		// equipCliList = append(equipCliList, equipCli)
		bagEquipCliList = append(bagEquipCliList, equipCli)
	}
	sort.Slice(equipSvrList, func(i, j int) bool {
		return equipSvrList[i].GetEquipGuid() < equipSvrList[j].GetEquipGuid()
	})
	res.EquipList = equipSvrList
	if len(bagEquipCliList) > 0 {
		SendMazeBagEquipChgIDEx(userCtx, userCtx.UserID, bagEquipCliList, nil, nil, req.GetOpType())
	}
	return
}

func GetUserEquipAddition(logger fkserver.UserContext, uid uint64) (int32, error) {
	var stageAddition int32 = 1
	// ok, err := UserBlackDiamondFC.IsBlackDiamondUserFC(logger, uid, 50)
	// if err != nil {
	// 	logger.ErrorWF("GetUserEquipAddition get black diamond user fail", zap.Error(err))
	// 	return 0, err
	// }
	// if ok {
	// 	stageAddition = 3
	// } else {
	// 	e, monthCardEnabled := MonthlyCardRedis.CheckMonthlyCardEx(logger, uid, 800001)
	// 	if e != nil {
	// 		logger.ErrorWF("GetUserEquipAddition get month card fail", zap.Error(err))
	// 		return 0, e
	// 	}
	// 	if monthCardEnabled {
	// 		stageAddition = 2
	// 	}
	// }
	// return stageAddition, err
	return stageAddition, nil
}

func AddEquipAllotGuid(logger fklog.FKLogI, userId uint64, addCount int32) ([]int64, error) {
	if addCount == 0 {
		return make([]int64, 0), nil
	}
	newGuid, e := mazeequipguidredis.GetNewGuid(logger, userId, addCount)
	if e != nil {
		logger.ErrorWF("AddEquipAllotGuid GetNewGuid error", zap.Error(e))
		return nil, e
	}
	allotGuids := make([]int64, 0, addCount)
	for i := int32(1); i <= addCount; i++ {
		guid := newGuid - int64(addCount) + int64(i)
		allotGuids = append(allotGuids, guid)
	}
	logger.InfoWF("AddEquipAllotGuid success", zap.Int64s("allotGuids", allotGuids))
	return allotGuids, nil
}

func GetTotalScoreAndBarrierMap(logger fklog.FKLogI, userId uint64, equipList []*MazeEquipSvr.SvrEquipInfo) (map[int32]int32, error) {
	if len(equipList) <= 0 {
		return make(map[int32]int32), nil
	}
	scoreGroupMap := make(map[int32]struct{}, 0)
	for _, equipId := range equipList {
		equipCfg := GMazeEquipInfoV8Cfg.Get(equipId.GetEquipId())
		if equipCfg == nil {
			logger.ErrorWF("GetTotalScoreAndBarrierMap GMazeEquipInfoV8Cfg error", zap.Int32("equipId", equipId.GetEquipId()))
			return nil, errors.New("配置不存在")
		}
		scoreGroupMap[equipCfg.Score_group] = struct{}{}
	}
	equipIds := make([]int32, 0)
	for equipId := range scoreGroupMap {
		equipIds = append(equipIds, equipId)
	}
	equipScoreMap, err := mazeequipgetnumredis.GetBatchEquipGetNum(logger, userId, equipIds)
	if err != nil {
		logger.ErrorWF("GetTotalScoreAndBarrierMap GetBatchEquipGetNum error", zap.Error(err))
		return nil, err
	}
	return equipScoreMap, nil
}

func DoInsEquip(logger fklog.FKLogI, userId uint64, uwparam *DEIUWParam, ewparam *DEIEWParam) (ins *DEInstance, err error) {
	dei := NewDEInstance(logger, userId)
	err = dei.DeiInit(uwparam, ewparam)
	if err != nil {
		logger.ErrorWF("DoInsEquip DeiInit error!", zap.Error(err))
		return nil, err
	}
	err = dei.DeiInstance()
	if err != nil {
		logger.ErrorWF("DoInsEquip DeiInstance error!", zap.Error(err))
		return nil, err
	}
	return dei, nil
}

func GetPoolLimitMap(logger fklog.FKLogI, attrLimits []int64) map[int32]struct{} {
	affixLimits := make([]int32, 0, len(attrLimits))
	for _, attrId := range attrLimits {
		affixLimits = append(affixLimits, int32(attrId))
	}
	poolLimitMap := mazeequipaffixrandpoolv8.GetPoolLimitCfg(affixLimits)
	logger.InfoWF("GetPoolLimitMap end", zap.Int("attrLimits", len(attrLimits)), zap.Int("poolLimitMap", len(poolLimitMap)))
	return poolLimitMap
}
