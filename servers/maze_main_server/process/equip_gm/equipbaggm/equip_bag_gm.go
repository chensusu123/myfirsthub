package equipbaggm

import (
	"fmt"
	"net/http"
	"strings"

	"gitlab.ifreetalk.com/plate/excel/auto/GMazeEquipInfoV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/common/new_map_db/FamilyAllocUserRedis"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/common/new_map_db/LeagueFamilyRedis"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/common/new_map_db/WorldLeagueRedis"
	"gitlab.ifreetalk.com/plate/protodef/MazeEquipSvr"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/gm"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/uniqueid"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/rpc/dollequipbagrpc"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/fileio"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeequipgetnumredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/dollassemblesuitredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebagequipredis"
)

func Reg(logger fklog.FKLogI) {
	gm.SafeHttpRegister(logger, "/AddEquip", func(writer http.ResponseWriter, request *http.Request) {

		// 外网线上环境不允许使用GM
		request.ParseForm()

		userId := fkutil.ToUint64(request.Form.Get("userId"))
		equipId := fkutil.ToInt32(request.Form.Get("equipId"))
		ruleId := fkutil.ToInt32(request.Form.Get("ruleId"))
		subType := fkutil.ToInt32(request.Form.Get("subType"))
		headType := fkutil.ToInt32(request.Form.Get("headType"))
		tailType := fkutil.ToInt32(request.Form.Get("tailType"))

		if subType > 6 {
			writer.Write([]byte("subType 子类型无效"))
			return
		}

		req := &MazeEquipSvr.SvrAddMazeEquipRQ{
			UserId: proto.Uint64(userId),
			// ShipNumber: proto.Int32(1),
		}
		equipInfo := &MazeEquipSvr.SvrEquipInfo{
			EquipId: proto.Int32(equipId),
		}
		if ruleId > 0 {
			equipInfo.Conditions = append(equipInfo.Conditions, &MazeEquipSvr.ConditionInfo{
				Id:    proto.Int32(8),
				Value: proto.Int64(int64(ruleId)),
			})
		}
		if subType > 0 {
			equipInfo.Conditions = append(equipInfo.Conditions, &MazeEquipSvr.ConditionInfo{
				Id:    proto.Int32(4),
				Value: proto.Int64(int64(subType)),
			})
		}
		if headType > 0 {
			equipInfo.Conditions = append(equipInfo.Conditions, &MazeEquipSvr.ConditionInfo{
				Id:    proto.Int32(32),
				Value: proto.Int64(int64(headType)),
			})
		}
		if tailType > 0 {
			equipInfo.Conditions = append(equipInfo.Conditions, &MazeEquipSvr.ConditionInfo{
				Id:    proto.Int32(64),
				Value: proto.Int64(int64(tailType)),
			})
		}
		req.EquipList = append(req.EquipList, equipInfo)
		res := &MazeEquipSvr.SvrAddMazeEquipRS{}
		req.OpType = proto.Int32(int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_FOE))
		req.TradeNumber = proto.Uint64(uniqueid.GenUniqueIdUInt64())
		err := dollequipbagrpc.MazeBagAddRQ(logger, req, res)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write([]byte("ok"))
		return
	})

	gm.SafeHttpRegister(logger, "/ClearBag", func(writer http.ResponseWriter, request *http.Request) {
		// 外网线上环境不允许使用GM
		request.ParseForm()

		userId := fkutil.ToUint64(request.Form.Get("userId"))

		err := ClearUserBag(logger, userId)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}

		writer.Write([]byte("ok"))

		return
	})

	gm.SafeHttpRegister(logger, "/BatchClearBag", func(writer http.ResponseWriter, request *http.Request) {
		// 外网线上环境不允许使用GM
		request.ParseForm()

		filename := request.Form.Get("file")

		fr := fileio.NewDefFReaderEx(logger, ",")
		err := fr.Open(filename)
		if err != nil {
			return
		}
		defer fr.Close()

		fr.SetDumpRow(5000)
		// fr.SetSleep(int32(waitLine), int32(sleep))
		fkfmt.Println("open file", filename, "succ")
		defer fkutil.CaptureException()

		fr.InitAsync(int(8), int(100))

		fr.Range(func(logger fklog.FKLogI, line []uint64) bool {
			if len(line) != 1 {
				logger.ErrorWF("file line not match")
				return false
			}
			userId := line[0]

			err := ClearUserBag(logger, userId)
			if err != nil {
				logger.ErrorWF("BatchClearBag ClearUserBag fail", zap.Error(err), zap.Uint64("uid", userId))
				return false
			}

			return true
		})

		writer.Write([]byte("ok"))

		return
	})

	gm.SafeHttpRegister(logger, "/ClearBagByMap", func(writer http.ResponseWriter, request *http.Request) {
		// 外网线上环境不允许使用GM
		request.ParseForm()

		mapId := fkutil.ToUint64(request.Form.Get("mapId"))

		leagueIDMap, err := WorldLeagueRedis.GetWorldLeagueInfo(logger, mapId)
		if err != nil {
			logger.ErrorWF("ClearBagByMap load world leagueInfo fail",
				zap.Uint64("mapID", mapId),
				zap.Error(err))
			return
		}
		logger.InfoWF("ClearBagByMap map info",
			zap.Uint64("map", mapId),
			zap.Int("leagueLen", len(leagueIDMap)),
		)

		var familyID uint64
		for leagueID := range leagueIDMap {
			// 取联盟下的散人家族
			familyIDs, err := LeagueFamilyRedis.GetAllLeagueFamilyIDs(logger, leagueID)
			if err != nil {
				logger.ErrorWF("ClearBagByMap get league familyIDs fail", zap.Any("leagueID", leagueID), zap.Error(err))
				continue
			}

			for _, family := range familyIDs {
				familyID = fkutil.ToUint64(family)
				// 取家族下所有人
				users, err := FamilyAllocUserRedis.GetAllFamilyUIDSliceFix(logger, familyID)
				if err != nil {
					logger.ErrorWF("ClearBagByMap get family users fail", zap.Error(err))
					continue
				}

				if len(users) == 0 {
					continue
				}

				for _, uid := range users {
					err := ClearUserBag(logger, uid)
					if err != nil {
						logger.ErrorWF("ClearBagByMap ClearUserBag fail", zap.Error(err), zap.Uint64("uid", uid))
						continue
					}
				}
			}
		}

		writer.Write([]byte("ok"))

		return
	})

	gm.SafeHttpRegister(logger, "/CheckDollEquipCfg", func(writer http.ResponseWriter, request *http.Request) {
		// 外网线上环境不允许使用GM
		request.ParseForm()

		// for _, cfg := range GDollEquipInfoV8Cfg.GetAll() {
		// 	//校验基础属性
		// 	var maxBaseNum int32
		// 	for baseNum := range cfg.Affix_base_num {
		// 		if baseNum > maxBaseNum {
		// 			maxBaseNum = baseNum
		// 		}
		// 	}
		// 	groupBasePoolMap := make(map[int32]map[int32]int32, 0)
		// 	for poolId, index := range cfg.Affix_base_pool {
		// 		if index == 0 {
		// 			continue
		// 		}
		// 		if groupBasePoolMap[index] == nil {
		// 			groupBasePoolMap[index] = make(map[int32]int32, 0)
		// 		}
		// 		groupBasePoolMap[index][poolId] = 1000
		// 	}
		// 	for index := int32(1); index <= maxBaseNum; index++ {
		// 		basePoolMap, ok := groupBasePoolMap[index]
		// 		if !ok {
		// 			logger.ErrorWF("doll_equip_info_v8【人偶-装备-信息】.xlsx 基础属性条数缺失 ", zap.Any("装备id：", cfg.Equipment_id), zap.Any("缺失条数位置：", index))
		// 			continue
		// 		}
		// 		for poolId := range basePoolMap {
		// 			poolMap := mazeequipaffixrandpoolv8.GetEquipPoolWeightCfg(poolId)
		// 			if len(poolMap) == 0 {
		// 				logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 基础词条库缺失 ", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId))
		// 				continue
		// 			}
		// 		}
		// 	}

		// 	//校验随机属性
		// 	var maxRandNum int32
		// 	for randNum := range cfg.Affix_rand_num {
		// 		if randNum > maxRandNum {
		// 			maxRandNum = randNum
		// 		}
		// 	}
		// 	for poolId := range cfg.Affix_rand_pool {
		// 		if poolId == 0 {
		// 			continue
		// 		}
		// 		poolMap := mazeequipaffixrandpoolv8.GetEquipPoolWeightCfg(poolId)
		// 		if len(poolMap) == 0 {
		// 			logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 随机词条库缺失 ", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId))
		// 			continue
		// 		}
		// 		if len(poolMap) < int(maxRandNum) {
		// 			logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 随机词条库条数不足", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId),
		// 				zap.Any("当前词条数量:", len(poolMap)), zap.Any("最大词条数量", maxRandNum))
		// 			continue
		// 		}
		// 		poolGroupMap := mazeequipaffixrandpoolv8.GetPoolGroupCfg(poolId)
		// 		if len(poolGroupMap) < int(maxRandNum) {
		// 			logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 随机词条库条去重组不足", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId),
		// 				zap.Any("当前词条去重组数量:", len(poolGroupMap)), zap.Any("最大词条数量", maxRandNum))
		// 			continue
		// 		}
		// 	}

		// 	//校验天赋属性
		// 	var maxModNum int32
		// 	for modNum := range cfg.Affix_mod_num {
		// 		if modNum > maxModNum {
		// 			maxModNum = modNum
		// 		}
		// 	}
		// 	for poolId := range cfg.Affix_mod_pool {
		// 		if poolId == 0 {
		// 			continue
		// 		}
		// 		poolWeightMap := GDollEquipAffixModPoolV8CfgEx.GetEnchantPoolWeightCfg(poolId)
		// 		if len(poolWeightMap) == 0 {
		// 			logger.ErrorWF("doll_equip_affix_mod_pool_v8【人偶-装备-附魔词条随机池】.xlsx 随机词条库缺失 ", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId))
		// 			continue
		// 		}
		// 		if len(poolWeightMap) < int(maxModNum) {
		// 			logger.ErrorWF("doll_equip_affix_mod_pool_v8【人偶-装备-附魔词条随机池】.xlsx 随机词条库条数不足", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId),
		// 				zap.Any("当前词条数量:", len(poolWeightMap)), zap.Any("最大词条数量", maxModNum))
		// 			continue
		// 		}
		// 		poolGroupMap := GDollEquipAffixModPoolV8CfgEx.GetEnchantPoolGroupCfg(poolId)
		// 		if len(poolGroupMap) < int(maxModNum) {
		// 			logger.ErrorWF("doll_equip_affix_mod_pool_v8【人偶-装备-附魔词条随机池】.xlsx 随机词条库条去重组不足", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId),
		// 				zap.Any("当前词条去重组数量:", len(poolGroupMap)), zap.Any("最大词条数量", maxModNum))
		// 			continue
		// 		}
		// 	}

		// 	//校验特殊属性
		// 	for poolId := range cfg.Affix_mod_pool {
		// 		if poolId == 0 {
		// 			continue
		// 		}
		// 		poolMap := GDollEquipAffixSpPoolV8CfgEx.GetSpPoolWeightCfg(poolId)
		// 		if len(poolMap) == 0 {
		// 			logger.ErrorWF("doll_equip_affix_sp_pool_v8【人偶-装备-特殊词条随机库】.xlsx 随机词条库缺失 ", zap.Any("装备id:", cfg.Equipment_id), zap.Any("池子id:", poolId))
		// 			continue
		// 		}
		// 	}

		// }

		// for _, cfg := range GDollEquipAffixSpRuleV8Cfg.GetAll() {
		// 	//校验基础属性
		// 	var maxBaseNum int32
		// 	for baseNum := range cfg.Affix_base_num {
		// 		if baseNum > maxBaseNum {
		// 			maxBaseNum = baseNum
		// 		}
		// 	}
		// 	groupBasePoolMap := make(map[int32]map[int32]int32, 0)
		// 	for poolId, index := range cfg.Affix_base_pool {
		// 		if index == 0 {
		// 			continue
		// 		}
		// 		if groupBasePoolMap[index] == nil {
		// 			groupBasePoolMap[index] = make(map[int32]int32, 0)
		// 		}
		// 		groupBasePoolMap[index][poolId] = 1000
		// 	}
		// 	for index := int32(1); index <= maxBaseNum; index++ {
		// 		basePoolMap, ok := groupBasePoolMap[index]
		// 		if !ok {
		// 			logger.ErrorWF("doll_equip_affix_sp_rule_v8【人偶-装备-生成特殊词条规则】.xlsx 基础属性条数缺失 ", zap.Any("装备规则id：", cfg.Rule_id), zap.Any("缺失条数位置：", index))
		// 			continue
		// 		}
		// 		for poolId := range basePoolMap {
		// 			poolMap := mazeequipaffixrandpoolv8.GetEquipPoolWeightCfg(poolId)
		// 			if len(poolMap) == 0 {
		// 				logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 基础词条库缺失 ", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId))
		// 				continue
		// 			}
		// 		}
		// 	}

		// 	//校验随机属性
		// 	var maxRandNum int32
		// 	for randNum := range cfg.Affix_rand_num {
		// 		if randNum > maxRandNum {
		// 			maxRandNum = randNum
		// 		}
		// 	}
		// 	for poolId := range cfg.Affix_rand_pool {
		// 		if poolId == 0 {
		// 			continue
		// 		}
		// 		poolMap := mazeequipaffixrandpoolv8.GetEquipPoolWeightCfg(poolId)
		// 		if len(poolMap) == 0 {
		// 			logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 随机词条库缺失 ", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId))
		// 			continue
		// 		}
		// 		if len(poolMap) < int(maxRandNum) {
		// 			logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 随机词条库条数不足", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId),
		// 				zap.Any("当前词条数量:", len(poolMap)), zap.Any("最大词条数量", maxRandNum))
		// 			continue
		// 		}
		// 		poolGroupMap := mazeequipaffixrandpoolv8.GetPoolGroupCfg(poolId)
		// 		if len(poolGroupMap) < int(maxRandNum) {
		// 			logger.ErrorWF("doll_equip_affix_pool_v8【人偶-装备-词条随机池】.xlsx 随机词条库条去重组不足", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId),
		// 				zap.Any("当前词条去重组数量:", len(poolGroupMap)), zap.Any("最大词条数量", maxRandNum))
		// 			continue
		// 		}
		// 	}

		// 	//校验天赋属性
		// 	var maxModNum int32
		// 	for modNum := range cfg.Affix_mod_num {
		// 		if modNum > maxModNum {
		// 			maxModNum = modNum
		// 		}
		// 	}
		// 	for poolId := range cfg.Affix_mod_pool {
		// 		if poolId == 0 {
		// 			continue
		// 		}
		// 		poolWeightMap := GDollEquipAffixModPoolV8CfgEx.GetEnchantPoolWeightCfg(poolId)
		// 		if len(poolWeightMap) == 0 {
		// 			logger.ErrorWF("doll_equip_affix_mod_pool_v8【人偶-装备-附魔词条随机池】.xlsx 随机词条库缺失 ", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId))
		// 			continue
		// 		}
		// 		if len(poolWeightMap) < int(maxModNum) {
		// 			logger.ErrorWF("doll_equip_affix_mod_pool_v8【人偶-装备-附魔词条随机池】.xlsx 随机词条库条数不足", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId),
		// 				zap.Any("当前词条数量:", len(poolWeightMap)), zap.Any("最大词条数量", maxModNum))
		// 			continue
		// 		}
		// 		poolGroupMap := GDollEquipAffixModPoolV8CfgEx.GetEnchantPoolGroupCfg(poolId)
		// 		if len(poolGroupMap) < int(maxModNum) {
		// 			logger.ErrorWF("doll_equip_affix_mod_pool_v8【人偶-装备-附魔词条随机池】.xlsx 随机词条库条去重组不足", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId),
		// 				zap.Any("当前词条去重组数量:", len(poolGroupMap)), zap.Any("最大词条数量", maxModNum))
		// 			continue
		// 		}
		// 	}

		// 	//校验特殊属性
		// 	for poolId := range cfg.Affix_extra_pool {
		// 		if poolId == 0 {
		// 			continue
		// 		}
		// 		poolMap := GDollEquipAffixSpPoolV8CfgEx.GetSpPoolWeightCfg(poolId)
		// 		if len(poolMap) == 0 {
		// 			logger.ErrorWF("doll_equip_affix_sp_pool_v8【人偶-装备-特殊词条随机库】.xlsx 随机词条库缺失 ", zap.Any("装备规则id:", cfg.Rule_id), zap.Any("池子id:", poolId))
		// 			continue
		// 		}
		// 	}

		// }
		writer.Write([]byte("ok"))

		return
	})

	// gm.SafeHttpRegister(logger, "/FixDelModExpEquipByMap", func(writer http.ResponseWriter, request *http.Request) {
	// 	request.ParseForm()
	// 	mapId := fkutil.ToUint64(request.Form.Get("mapId"))
	// 	run := fkutil.ToUint64(request.Form.Get("run"))
	// 	leagueIDMap, err := WorldLeagueRedis.GetWorldLeagueInfo(logger, mapId)
	// 	if err != nil {
	// 		logger.ErrorWF("FixDelModExpEquipByMap load world leagueInfo fail",
	// 			zap.Uint64("mapID", mapId),
	// 			zap.Error(err))
	// 		return
	// 	}
	// 	logger.InfoWF("FixDelModExpEquipByMap map info",
	// 		zap.Uint64("map", mapId),
	// 		zap.Int("leagueLen", len(leagueIDMap)),
	// 	)

	// 	var familyID uint64
	// 	var cnt int
	// 	for leagueID := range leagueIDMap {
	// 		// 取联盟下的散人家族
	// 		familyIDs, err := LeagueFamilyRedis.GetAllLeagueFamilyIDs(logger, leagueID)
	// 		if err != nil {
	// 			logger.ErrorWF("FixDelModExpEquipByMap get league familyIDs fail", zap.Any("leagueID", leagueID), zap.Error(err))
	// 			continue
	// 		}

	// 		for _, family := range familyIDs {
	// 			familyID = fkutil.ToUint64(family)
	// 			// 取家族下所有人
	// 			users, err := FamilyAllocUserRedis.GetAllFamilyUIDSliceFix(logger, familyID)
	// 			if err != nil {
	// 				logger.ErrorWF("FixDelModExpEquipByMap get family users fail", zap.Error(err))
	// 				continue
	// 			}

	// 			if len(users) == 0 {
	// 				continue
	// 			}

	// 			for _, uid := range users {
	// 				logger.SetUid(uid)
	// 				n, _ := GmDelExpModAttrBag(logger, uid, run == 1)
	// 				if n {
	// 					cnt++
	// 				}
	// 			}
	// 		}
	// 	}

	// 	writer.Write([]byte(fmt.Sprintf("需要修复数量:%d", cnt)))
	// })

	// gm.SafeHttpRegister(logger, "/FixDelModExpEquip", func(writer http.ResponseWriter, request *http.Request) {
	// 	userId := fkutil.ToUint64(request.Form.Get("userId"))
	// 	logger.SetUid(userId)
	// 	_, e := GmDelExpModAttrBag(logger, userId, true)
	// 	if e == nil {
	// 		writer.Write([]byte("ok"))
	// 	} else {
	// 		writer.Write([]byte("fail"))
	// 	}
	// })

	gm.SafeHttpRegister(logger, "/SetEquipRollScore", func(writer http.ResponseWriter, request *http.Request) {
		// 外网线上环境不允许使用GM
		request.ParseForm()

		uid := fkutil.ToUint64(request.Form.Get("uid"))
		equipId := fkutil.ToInt32(request.Form.Get("equipId"))
		score := fkutil.ToInt32(request.Form.Get("score"))

		if uid == 0 || equipId == 10 || score < 0 {
			writer.Write([]byte("uid 不能为0 , equipId 不能为0, score 不能小于0"))
			return
		}

		cfg := GMazeEquipInfoV8Cfg.Get(equipId)
		if cfg == nil {
			writer.Write([]byte("equipId找不到对应的装备配置"))
			return
		}

		err := mazeequipgetnumredis.SetEquipGetNum(logger, uid, cfg.Score_group, score)
		if err != nil {
			logger.ErrorWF("SetEquipRollScore  SetEquipGetNum fail", zap.Error(err), zap.Uint64("uid", uid),
				zap.Int32("equipId", equipId), zap.Int32("score", score))
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write([]byte("ok"))
		return
	})

	gm.SafeHttpRegister(logger, "/ClearBagNotAssemble", func(writer http.ResponseWriter, request *http.Request) {
		// 外网线上环境不允许使用GM
		request.ParseForm()

		uid := fkutil.ToUint64(request.Form.Get("uid"))

		if uid == 0 {
			writer.Write([]byte("uid 不能为0"))
			return
		}

		// 获取身上的装备信息
		assembleInfoMap, err := dollassemblesuitredis.GetAllDollAssembleSuit(logger, uid)
		if err != nil {
			logger.ErrorWF("ClearBagNotAssemble GetAllDollAssembleSuit fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}
		// 身上的装备
		assembleGuid := make(map[int64]struct{})
		if len(assembleInfoMap) > 0 {
			for _, equipList := range assembleInfoMap {
				if len(equipList) > 0 {
					for _, v := range equipList {
						assembleGuid[v.GetEquipGuid()] = struct{}{}
					}
				}
			}
		}

		equipInfoMap, err := mazebagequipredis.GetAllEquipInfo(logger, uid)
		if err != nil {
			logger.ErrorWF("ClearBagNotAssemble GetAllEquipInfo fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}

		equipGuids := make([]int64, 0)
		for k := range equipInfoMap {
			_, ok := assembleGuid[k]
			if !ok {
				equipGuids = append(equipGuids, k)
			}
		}

		err = ClearEquipBagBatch(logger, uid, equipGuids)
		if err != nil {
			logger.ErrorWF("ClearBagNotAssemble ClearEquipBagBatch fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write([]byte("ok"))
		return
	})

	gm.SafeHttpRegister(logger, "/BatchAddEquip", func(writer http.ResponseWriter, request *http.Request) {
		// 外网线上环境不允许使用GM
		request.ParseForm()

		uid := fkutil.ToUint64(request.Form.Get("uid"))
		param := request.Form.Get("equips")
		rp := request.Form.Get("rules")
		if uid <= 0 {
			writer.Write([]byte("uid 不能为0"))
			return
		}

		if param == "" {
			writer.Write([]byte("请指定装备参数"))
			return
		}
		rulesMap := ParseRules(rp)
		var equipConds []*MazeEquipSvr.SvrEquipInfo
		pairs := strings.Split(param, "_")
		req := &MazeEquipSvr.SvrAddMazeEquipRQ{
			UserId: proto.Uint64(uid),
		}
		for _, kv := range pairs {
			elems := strings.Split(kv, ":")
			if len(elems) == 2 {
				equipId := fkutil.ToInt32(elems[0])
				cfg := GMazeEquipInfoV8Cfg.Get(equipId)
				if cfg == nil {
					writer.Write([]byte(fmt.Sprintf("equipId(%d)找不到对应的装备配置", equipId)))
					return
				}

				cnt := fkutil.ToInt32(elems[1])
				for i := 1; i <= int(cnt); i++ {
					cond := MazeEquipSvr.SvrEquipInfo{}
					cond.EquipId = proto.Int32(equipId)
					if rulesMap[equipId] > 0 {
						cond.Conditions = append(cond.Conditions, &MazeEquipSvr.ConditionInfo{
							Id:    proto.Int32(8),
							Value: proto.Int64(int64(rulesMap[equipId]))})
					}
					equipConds = append(equipConds, &cond)
				}
			}
		}
		if len(equipConds) > 100 {
			writer.Write([]byte("一次添加装备太多,最多100件"))
			return
		}

		req.EquipList = append(req.EquipList, equipConds...)
		res := &MazeEquipSvr.SvrAddMazeEquipRS{}
		req.OpType = proto.Int32(int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_FOE))
		req.TradeNumber = proto.Uint64(uniqueid.GenUniqueIdUInt64())
		err := dollequipbagrpc.MazeBagAddRQ(logger, req, res)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write([]byte("ok"))
	})

	gm.SafeHttpRegister(logger, "/fixAllEquipAttrLimit", func(writer http.ResponseWriter, request *http.Request) {
		userId := fkutil.ToUint64(request.Form.Get("userId"))
		logger.SetUid(userId)
		e := fixAllEquipAttrLimit(logger, userId)
		if e == nil {
			writer.Write([]byte("ok"))
		} else {
			writer.Write([]byte("fail"))
		}
	})
	// gm.SafeHttpRegister(logger, "/resetInsAllEquip", func(writer http.ResponseWriter, request *http.Request) {
	// 	userId := fkutil.ToUint64(request.Form.Get("userId"))
	// 	logger.SetUid(userId)
	// 	e := resetInsAllEquip(logger, userId)
	// 	if e == nil {
	// 		writer.Write([]byte("ok"))
	// 	} else {
	// 		writer.Write([]byte("fail"))
	// 	}
	// })
}

func ClearUserBag(logger fklog.FKLogI, userId uint64) (err error) {
	// todo 找装配的删除
	err = ClearDollAssembleInfo(logger, userId)
	if err != nil {
		return
	}
	// todo 删除背包里的
	err = ClearEquipBag(logger, userId)
	if err != nil {
		return
	}

	return
}
