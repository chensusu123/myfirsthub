package equipaassemblegm

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync/atomic"
	"time"

	"maze_game_server/common/function/fileio"
	"maze_game_server/common/function/gm"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/io/redis/mazeuserlevelredis"
	"maze_game_server/module/calcassembleattr"
	"maze_game_server/module/dollassembleinfo"
	"maze_game_server/servers/maze_main_server/process/equip"
	"maze_game_server/servers/maze_main_server/process/equip_gm/equipbaggm"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

var EndLine = "-----------------------------------------------------------\n"

func RegGm(logger fklog.FKLogI) {
	gm.SafeHttpRegister(logger, "/LookAssembleInfo", func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userId := fkutil.ToUint64(request.Form.Get("user_id"))

		logger.SetLogId(time.Now().UnixNano())
		logger.SetUid(userId)
		logger.InfoWF("LookAssembleInfo begin")

		assembleInfo, effect, err := dollassembleinfo.GetDollAssembleInfoEx(logger, userId)
		if err != nil {
			logger.ErrorWF("LookAssembleInfo Get Assemble info fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}

		var showBuff bytes.Buffer
		header, err := PackAssembleHeader(ctx, userId, assembleInfo)
		if err != nil {
			logger.ErrorWF("LookAssembleInfo PackAssembleHeader fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}
		showBuff.WriteString(header)

		showBuff.WriteString("装备位信息:\n")
		sort.Slice(assembleInfo.MazeEquips, func(i, j int) bool {
			return assembleInfo.MazeEquips[i].GetEquipPos().GetPos() <= assembleInfo.MazeEquips[j].GetEquipPos().GetPos()
		})
		for _, posInfo := range assembleInfo.MazeEquips {
			DumpEquipPos(logger, userId, &showBuff, posInfo.GetEquipPos().GetPos(), posInfo, assembleInfo.GetEpSuitId())
		}
		showBuff.WriteString(EndLine)
		showBuff.WriteString(fmt.Sprintf("装备套装:%d\n", assembleInfo.GetEpSuitId()))
		suitBuff, e := PackEquipSuitInfo(logger, userId, assembleInfo, effect)
		if e == nil {
			// 汇总套装属性加成
			showBuff.WriteString(suitBuff)
		}
		showBuff.WriteString(EndLine)

		ar, err := DumpDollCalcAttr(logger, userId)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		showBuff.WriteString(ar)
		showBuff.WriteString(EndLine)

		// ar, err = DumpForceAttr(logger, userId)
		// if err != nil {
		// 	writer.Write([]byte(err.Error()))
		// 	return
		// }
		// showBuff.WriteString(ar)
		// showBuff.WriteString(EndLine)

		ar, err = DumpNoForceAttr(logger, userId)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		showBuff.WriteString(ar)
		showBuff.WriteString(EndLine)

		r := showBuff.String()
		writer.Write([]byte(r))
		logger.InfoWF("LookAssembleInfo end")
	})

	gm.SafeHttpRegister(logger, "/GetEquipInfoByCfgId", func(writer http.ResponseWriter, request *http.Request) {
		userId := fkutil.ToUint64(request.Form.Get("user_id"))
		equipId := fkutil.ToInt32(request.Form.Get("equipId"))
		pos := fkutil.ToInt32(request.Form.Get("pos"))
		minLv := fkutil.ToInt32(request.Form.Get("minLv"))
		maxLv := fkutil.ToInt32(request.Form.Get("maxLv"))
		quality := fkutil.ToInt32(request.Form.Get("quality"))
		ruleId := fkutil.ToInt32(request.Form.Get("ruleId"))
		madeTime := fkutil.ToInt64(request.Form.Get("madetime"))
		guid := fkutil.ToInt64(request.Form.Get("guid"))
		lock := fkutil.ToInt32(request.Form.Get("lock"))
		showSeal := fkutil.ToInt32(request.Form.Get("showseal"))
		help := fkutil.ToInt32(request.Form.Get("help"))
		queryInUse := fkutil.ToInt32(request.Form.Get("queryInUse"))
		tmpBag := fkutil.ToInt32(request.Form.Get("tmpbag"))
		if help > 0 {
			writer.Write([]byte(CondHelp()))
			return
		}
		cd := BagCond{}
		cd.EquipId = equipId
		cd.MaxLv = maxLv
		cd.MinLv = minLv
		cd.Quality = quality
		cd.Pos = pos
		cd.RuleId = ruleId
		cd.MadeTime = madeTime
		cd.Lock = lock
		cd.Guid = guid
		cd.QueryInUse = queryInUse
		cd.TmpBag = tmpBag
		cd.ShowSeal = showSeal
		logger.SetLogId(time.Now().UnixNano())
		logger.SetUid(userId)
		logger.InfoWF("GetEquipInfoByCfgId", zap.Any("cond", cd))
		rs, err := GetEquipInfoByCfgId(logger, userId, cd)
		if err != nil {
			logger.ErrorWF("GetEquipInfoByCfgId", zap.Error(err))
			writer.Write([]byte("执行失败"))
			return
		}
		if len(rs) == 0 {
			writer.Write([]byte("查询结果为空"))
			return
		}
		writer.Write([]byte(rs))
	})

	gm.SafeHttpRegister(logger, "/GetEquipInfoByGuid", func(writer http.ResponseWriter, request *http.Request) {
		userId := fkutil.ToUint64(request.Form.Get("user_id"))
		guid := fkutil.ToInt64(request.Form.Get("guid"))
		rs, err := GetEquipInfoByGuid(logger, userId, guid)
		if err != nil {
			logger.ErrorWF("GetEquipInfoByGuid", zap.Error(err))
			writer.Write([]byte("执行失败"))
			return
		}
		writer.Write([]byte(rs))
	})

	gm.SafeHttpRegister(logger, "/FixDollEquipAttr", func(writer http.ResponseWriter, request *http.Request) {
		logger.InfoWF("FixDollEquipAttr start")

		fPath := request.Form.Get("file")
		fr := fileio.NewDefFReaderEx(logger, ",")
		err := fr.Open(fPath)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		defer fr.Close()

		var total, succ, fail int32
		fr.Range(func(logger fklog.FKLogI, line []uint64) bool {
			atomic.AddInt32(&total, 1)
			if len(line) != 1 {
				return true
			}
			userId := line[0]
			assembleInfo, err := dollassembleinfo.GetDollAssembleInfo(logger, userId)
			if err != nil {
				atomic.AddInt32(&fail, 1)
				logger.ErrorWF("FixDollEquipForce Get Assemble info fail", zap.Error(err))
				return true
			}
			_, otherAttrs, e := calcassembleattr.CalcEquipAttrs(logger, assembleInfo.MazeEquips)
			if e != nil {
				atomic.AddInt32(&fail, 1)
				return true
			}

			// 更新buff中心
			mazebuffinforedis.SaveMazeEquipBuff(logger, userId, otherAttrs)
			//	buffcenter.NotifyBuff(logger, userId, otherAttrs, constdef.ENUM_BUFF_SOURCE_DOLL_EQUIP, "maze_equip_gm_server")

			// e = dollassembleattrredis.SetDollAssembleAttr(logger, userId, constdef.AttrFieldEquip, forceAttrs)
			// if e == nil {
			// 	e = dollforcechgnotifyqueue.DollForceChgNotice(logger, userId, constdef.ForcePreviewEquip, int32(constdef.ForceChgTypeGm), "")
			// 	if e == nil {
			// 		atomic.AddInt32(&succ, 1)
			// 		return true
			// 	}
			// }
			// atomic.AddInt32(&fail, 1)
			return true
		})
		_, _ = writer.Write([]byte(fmt.Sprintf("执行结果=total:%d succ:%d fail:%d", total, succ, fail)))
	})

	gm.SafeHttpRegister(logger, "/SendOneSuitEquip", func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		userId := fkutil.ToUint64(request.Form.Get("user_id"))
		dressLv := fkutil.ToInt32(request.Form.Get("dressLv"))
		LvDis := fkutil.ToInt32(request.Form.Get("LvDis"))
		suitId := fkutil.ToInt32(request.Form.Get("suitId"))
		quality := fkutil.ToInt32(request.Form.Get("quality"))
		pos := fkutil.ToInt32(request.Form.Get("pos"))
		subType := fkutil.ToInt32(request.Form.Get("subType"))
		p := EquipParam{}
		if dressLv <= 0 {
			dollLv, e := mazeuserlevelredis.GetUserLevel(ctx, userId)
			if e != nil {
				writer.Write([]byte(e.Error()))
				return
			}
			dressLv = int32(dollLv)
		}
		p.DressLv = dressLv
		p.Quality = quality
		p.Pos = pos
		p.SuitId = suitId
		p.DisDressLv = LvDis
		p.SubType = subType
		e := CheckEquipParam(&p)
		if e != nil {
			writer.Write([]byte(e.Error()))
			return
		}
		equips, err := AddEquipByCond(context.TODO(), userId, p)
		if err != nil {
			writer.Write([]byte(fmt.Sprintf("执行失败: %s", err.Error())))
		} else {
			var bs bytes.Buffer
			bs.WriteString("执行结果:\n")
			for _, equip := range equips {
				bs.WriteString(fmt.Sprintf("装备位置:%d 装备Id:%d 装备Guid:%d Result:%s\n",
					equip.Equip.GetPos(),
					equip.Equip.GetEquipId(),
					equip.Equip.GetEquipGuid(),
					equip.Desc))
			}
			writer.Write(bs.Bytes())
		}
	})
	// // 重新初始化人偶装备
	// gm.SafeHttpRegister(logger, "/ReInitDollEquipByMap", func(writer http.ResponseWriter, request *http.Request) {
	// 	mapId := fkutil.ToUint64(request.Form.Get("mapId"))
	// 	clear := fkutil.ToBool(request.Form.Get("clear"))
	// 	leagueIDMap, err := WorldLeagueRedis.GetWorldLeagueInfo(logger, mapId)
	// 	if err != nil {
	// 		logger.ErrorWF("ReInitDollEquipByMap load world leagueInfo fail",
	// 			zap.Uint64("mapID", mapId),
	// 			zap.Error(err))
	// 		return
	// 	}
	// 	logger.InfoWF("ReInitDollEquipByMap map info",
	// 		zap.Uint64("map", mapId),
	// 		zap.Int("leagueLen", len(leagueIDMap)),
	// 	)
	//
	// 	var familyID uint64
	// 	for leagueID := range leagueIDMap {
	// 		// 取联盟下的散人家族
	// 		familyIDs, err := LeagueFamilyRedis.GetAllLeagueFamilyIDs(logger, leagueID)
	// 		if err != nil {
	// 			logger.ErrorWF("ReInitDollEquipByMap get league familyIDs fail", zap.Any("leagueID", leagueID), zap.Error(err))
	// 			continue
	// 		}
	//
	// 		for _, family := range familyIDs {
	// 			familyID = fkutil.ToUint64(family)
	// 			// 取家族下所有人
	// 			users, err := FamilyAllocUserRedis.GetAllFamilyUIDSliceFix(logger, familyID)
	// 			if err != nil {
	// 				logger.ErrorWF("ReInitDollEquipByMap get family users fail", zap.Error(err))
	// 				continue
	// 			}
	//
	// 			if len(users) == 0 {
	// 				continue
	// 			}
	//
	// 			for _, uid := range users {
	// 				if clear {
	// 					e := equipbaggm.ClearUserBag(logger, uid)
	// 					if e != nil {
	// 						continue
	// 					}
	// 				}
	// 				e := equip.ChkEquipPosUnlock(logger, uid, "gm", true)
	// 				if e != nil {
	// 					continue
	// 				}
	// 				// 初始装备套检查
	// 				e = equip.InitDollEquipSuitSeq(logger, uid)
	// 				if e != nil {
	// 					continue
	// 				}
	// 				// 处理初始化装备
	// 				e = equip.HandleDollEquipInit(logger, uid, true)
	// 				if e != nil {
	// 					continue
	// 				}
	// 			}
	// 		}
	// 	}
	//
	// 	writer.Write([]byte("ok"))
	// })

	gm.SafeHttpRegister(logger, "/ReInitDollEquipByFile", func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		fPath := request.Form.Get("file")
		fr := fileio.NewDefFReaderEx(logger, ",")
		err := fr.Open(fPath)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		defer fr.Close()

		var total, succ, fail int32
		fr.Range(func(logger fklog.FKLogI, line []uint64) bool {
			atomic.AddInt32(&total, 1)
			if len(line) != 1 {
				atomic.AddInt32(&fail, 1)
				return true
			}
			uid := line[0]
			e := equipbaggm.ClearUserBag(context.TODO(), uid)
			if e != nil {
				atomic.AddInt32(&fail, 1)
				return true
			}
			e = equip.ChkEquipPosUnlock(ctx, uid, "gm", true)
			if e != nil {
				atomic.AddInt32(&fail, 1)
				return true
			}
			// 初始装备套检查
			e = equip.InitDollEquipSuitSeq(logger, uid)
			if e != nil {
				atomic.AddInt32(&fail, 1)
				return true
			}
			// 处理初始化装备
			e = equip.HandleDollEquipInit(context.TODO(), uid, true)
			if e != nil {
				atomic.AddInt32(&fail, 1)
				return true
			}
			return true
		})
		_, _ = writer.Write([]byte(fmt.Sprintf("执行结果=total:%d succ:%d fail:%d", total, succ, fail)))

	})

	gm.SafeHttpRegister(logger, "/ReInitDollEquip", func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		uid := fkutil.ToUint64(request.Form.Get("user_id"))
		clear := fkutil.ToBool(request.Form.Get("clear"))
		if clear {
			e := equipbaggm.ClearUserBag(context.TODO(), uid)
			if e != nil {
				writer.Write([]byte(fmt.Sprintf("删除装备失败:%s", e.Error())))
				return
			}
		}
		e := equip.ChkEquipPosUnlock(ctx, uid, "gm", true)
		if e != nil {
			writer.Write([]byte(fmt.Sprintf("解锁装备位失败:%s", e.Error())))
			return
		}
		// 初始装备套检查
		e = equip.InitDollEquipSuitSeq(logger, uid)
		if e != nil {
			writer.Write([]byte(fmt.Sprintf("初始化当前套装失败:%s", e.Error())))
			return
		}
		// 处理初始化装备
		e = equip.HandleDollEquipInit(context.TODO(), uid, true)
		if e != nil {
			writer.Write([]byte(fmt.Sprintf("初始化装备失败:%s", e.Error())))
			return
		}
		writer.Write([]byte("初始化完成"))
	})

	// gm.SafeHttpRegister(logger, "/FixDollAttrByMap", func(writer http.ResponseWriter, request *http.Request) {
	// 	mapId := fkutil.ToUint64(request.Form.Get("mapId"))
	// 	fixType := fkutil.ToInt32(request.Form.Get("fixType"))
	// 	leagueIDMap, err := WorldLeagueRedis.GetWorldLeagueInfo(logger, mapId)
	// 	if err != nil {
	// 		logger.ErrorWF("FixDollAttrByMap load world leagueInfo fail",
	// 			zap.Uint64("mapID", mapId),
	// 			zap.Error(err))
	// 		return
	// 	}
	// 	logger.InfoWF("FixDollAttrByMap map info",
	// 		zap.Uint64("map", mapId),
	// 		zap.Int32("fixType", fixType),
	// 		zap.Int("leagueLen", len(leagueIDMap)),
	// 	)
	//
	// 	var familyID uint64
	// 	var succ, fail int64
	// 	for leagueID := range leagueIDMap {
	// 		// 取联盟下的散人家族
	// 		familyIDs, err := LeagueFamilyRedis.GetAllLeagueFamilyIDs(logger, leagueID)
	// 		if err != nil {
	// 			logger.ErrorWF("FixDollAttrByMap get league familyIDs fail", zap.Any("leagueID", leagueID), zap.Error(err))
	// 			continue
	// 		}
	//
	// 		for _, family := range familyIDs {
	// 			familyID = fkutil.ToUint64(family)
	// 			// 取家族下所有人
	// 			users, err := FamilyAllocUserRedis.GetAllFamilyUIDSliceFix(logger, familyID)
	// 			if err != nil {
	// 				logger.ErrorWF("FixDollAttrByMap get family users fail", zap.Error(err))
	// 				continue
	// 			}
	//
	// 			if len(users) == 0 {
	// 				continue
	// 			}
	//
	// 			for _, uid := range users {
	// 				logger.SetUid(uid)
	// 				logger.SetLogId(time.Now().UnixNano())
	// 				e := ReCalcDollEquipAttr(logger, uid, fixType)
	// 				if e == nil {
	// 					atomic.AddInt64(&succ, 1)
	// 				} else {
	// 					atomic.AddInt64(&fail, 1)
	// 				}
	// 			}
	// 		}
	// 	}
	//
	// 	writer.Write([]byte(fmt.Sprintf("succ:%d fail:%d", succ, fail)))
	// })

	// gm.SafeHttpRegister(logger, "/CalcDollAttr", func(writer http.ResponseWriter, request *http.Request) {
	// 	mapId := fkutil.ToUint64(request.Form.Get("mapId"))
	// 	leagueIDMap, err := WorldLeagueRedis.GetWorldLeagueInfo(logger, mapId)
	// 	if err != nil {
	// 		logger.ErrorWF("CalcDollAttr load world leagueInfo fail",
	// 			zap.Uint64("mapID", mapId),
	// 			zap.Error(err))
	// 		return
	// 	}
	// 	logger.InfoWF("CalcDollAttr map info",
	// 		zap.Uint64("map", mapId),
	// 		zap.Int("leagueLen", len(leagueIDMap)),
	// 	)
	//
	// 	var familyID uint64
	// 	var succ, fail int64
	// 	for leagueID := range leagueIDMap {
	// 		// 取联盟下的散人家族
	// 		familyIDs, err := LeagueFamilyRedis.GetAllLeagueFamilyIDs(logger, leagueID)
	// 		if err != nil {
	// 			logger.ErrorWF("FixDollAttrByMap get league familyIDs fail", zap.Any("leagueID", leagueID), zap.Error(err))
	// 			continue
	// 		}
	//
	// 		for _, family := range familyIDs {
	// 			familyID = fkutil.ToUint64(family)
	// 			// 取家族下所有人
	// 			users, err := FamilyAllocUserRedis.GetAllFamilyUIDSliceFix(logger, familyID)
	// 			if err != nil {
	// 				logger.ErrorWF("CalcDollAttr get family users fail", zap.Error(err))
	// 				continue
	// 			}
	//
	// 			if len(users) == 0 {
	// 				continue
	// 			}
	//
	// 			for _, uid := range users {
	// 				logger.SetUid(uid)
	// 				logger.SetLogId(time.Now().UnixNano())
	// 				e := CalcDollAttrCalc(logger, uid)
	// 				if e == nil {
	// 					atomic.AddInt64(&succ, 1)
	// 				} else {
	// 					atomic.AddInt64(&fail, 1)
	// 				}
	// 			}
	// 		}
	// 	}
	//
	// 	writer.Write([]byte(fmt.Sprintf("succ:%d fail:%d", succ, fail)))
	// })

	gm.SafeHttpRegister(logger, "/FixDollAttr", func(writer http.ResponseWriter, request *http.Request) {
		uid := fkutil.ToUint64(request.Form.Get("user_id"))
		fixType := fkutil.ToInt32(request.Form.Get("fixType"))
		logger.SetUid(uid)
		e := ReCalcDollEquipAttr(context.TODO(), uid, fixType)
		if e != nil {
			writer.Write([]byte(fmt.Sprintf("执行结果:%s", e.Error())))
		} else {
			writer.Write([]byte("执行完成"))
		}
	})

	gm.SafeHttpRegister(logger, "/FixDollAttrByFile", func(writer http.ResponseWriter, request *http.Request) {

		fPath := request.Form.Get("file")
		fixType := fkutil.ToInt32(request.Form.Get("fixType"))
		fr := fileio.NewDefFReaderEx(logger, ",")
		err := fr.Open(fPath)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		defer fr.Close()
		logger.InfoWF("FixDollAttrByFile param", zap.Int32("fixType", fixType))

		var total, succ, fail int32
		fr.Range(func(logger fklog.FKLogI, line []uint64) bool {
			atomic.AddInt32(&total, 1)
			if len(line) != 1 {
				atomic.AddInt32(&fail, 1)
				return true
			}
			uid := line[0]

			e := ReCalcDollEquipAttr(context.TODO(), uid, fixType)
			if e == nil {
				atomic.AddInt32(&succ, 1)
			} else {
				atomic.AddInt32(&fail, 1)
			}

			return true
		})
		_, _ = writer.Write([]byte(fmt.Sprintf("执行结果=total:%d succ:%d fail:%d", total, succ, fail)))

	})

	// gm.SafeHttpRegister(logger, "/UninstallEquip", func(writer http.ResponseWriter, request *http.Request) {

	// 	uid := fkutil.ToUint64(request.Form.Get("uid"))
	// 	pos := fkutil.ToInt32(request.Form.Get("pos"))
	// 	logger.SetUid(uid)
	// 	if pos != 0 {
	// 		if pos > 8 || pos < 1 {
	// 			_, _ = writer.Write([]byte("错误的pos参数,只支持[1,8]"))
	// 			return
	// 		}
	// 	}

	// 	assembleInfo, err := dollassembleinfo.GetDollAssembleInfo(logger, uid)
	// 	if err != nil {
	// 		_, _ = writer.Write([]byte(err.Error()))
	// 		return
	// 	}
	// 	if assembleInfo.GetCurSuitIndex() == 0 {
	// 		_, _ = writer.Write([]byte("OnDressEquipRQ no set cur suit"))
	// 		return
	// 	}
	// 	for _, ePos := range assembleInfo.GetMazeEquips() {
	// 		if pos != 0 && ePos.GetEquipPos().GetPos() != pos {
	// 			continue
	// 		}
	// 		rPos := ePos.GetEquipPos().GetPos()
	// 		downInfo := module.GetEquipPosInfo(assembleInfo, rPos)
	// 		if !assemble.IsAssembleEquip(downInfo) {
	// 			continue
	// 		}
	// 		downGuid := downInfo.GetEquipInfo().GetEquipGuid()
	// 		req := &DollEquip.DressDollEquipRQ{
	// 			EquipPos:          proto.Int32(int32(rPos)),
	// 			CurSuitSeq:        proto.Int32(assembleInfo.GetCurSuitIndex()),
	// 			EquipGuid:         proto.Int64(0),
	// 			ReplacedEquipGuid: proto.Int64(downGuid),
	// 		}
	// 		res := &DollEquip.DressDollEquipRS{}

	// 		err = process.OnDownEquipRQ(logger, uid, req, res)
	// 		if err != nil {
	// 			_, _ = writer.Write([]byte(err.Error()))
	// 			return
	// 		}

	// 		if res.ErrInfo.GetErrCode() != errors.NO_ERROR_CODE {
	// 			_, _ = writer.Write([]byte(res.ErrInfo.GetErrMsg()))
	// 			return
	// 		}
	// 		logger.InfoWF("UninstallEquip del pos succ", zap.Int32("pos", rPos), zap.Int64("guid", downGuid))

	// 	}
	// 	_, _ = writer.Write([]byte("ok"))
	// })

	gm.SafeHttpRegister(logger, "/FixEquipPosUnlock", func(writer http.ResponseWriter, request *http.Request) {
		uid := fkutil.ToUint64(request.Form.Get("user_id"))
		logger.SetUid(uid)
		cnt, e := UnlockPosByEquip(logger, uid)
		if e == nil {
			writer.Write([]byte(fmt.Sprintf("unlock:%d", cnt)))
		} else {
			writer.Write([]byte(e.Error()))
		}
	})

	// gm.SafeHttpRegister(logger, "/FixEquipPosUnlockByMap", func(writer http.ResponseWriter, request *http.Request) {
	// 	mapId := fkutil.ToUint64(request.Form.Get("mapId"))
	// 	leagueIDMap, err := WorldLeagueRedis.GetWorldLeagueInfo(logger, mapId)
	// 	if err != nil {
	// 		logger.ErrorWF("FixEquipPosUnlockByMap load world leagueInfo fail",
	// 			zap.Uint64("mapID", mapId),
	// 			zap.Error(err))
	// 		return
	// 	}
	// 	logger.InfoWF("FixEquipPosUnlockByMap map info",
	// 		zap.Uint64("map", mapId),
	// 		zap.Int("leagueLen", len(leagueIDMap)),
	// 	)
	//
	// 	var familyID uint64
	// 	var succ, fail int64
	// 	for leagueID := range leagueIDMap {
	// 		// 取联盟下的散人家族
	// 		familyIDs, err := LeagueFamilyRedis.GetAllLeagueFamilyIDs(logger, leagueID)
	// 		if err != nil {
	// 			logger.ErrorWF("FixEquipPosUnlockByMap get league familyIDs fail", zap.Any("leagueID", leagueID), zap.Error(err))
	// 			continue
	// 		}
	//
	// 		for _, family := range familyIDs {
	// 			familyID = fkutil.ToUint64(family)
	// 			// 取家族下所有人
	// 			users, err := FamilyAllocUserRedis.GetAllFamilyUIDSliceFix(logger, familyID)
	// 			if err != nil {
	// 				logger.ErrorWF("FixEquipPosUnlockByMap get family users fail", zap.Error(err))
	// 				continue
	// 			}
	//
	// 			if len(users) == 0 {
	// 				continue
	// 			}
	//
	// 			for _, uid := range users {
	// 				logger.SetUid(uid)
	// 				logger.SetLogId(time.Now().UnixNano())
	// 				_, e := UnlockPosByEquip(logger, uid)
	// 				if e == nil {
	// 					atomic.AddInt64(&succ, 1)
	// 				} else {
	// 					atomic.AddInt64(&fail, 1)
	// 				}
	// 			}
	// 		}
	// 	}
	//
	// 	writer.Write([]byte(fmt.Sprintf("succ:%d fail:%d", succ, fail)))
	// })

	gm.SafeHttpRegister(logger, "/FixAssembleEquipInfo", func(writer http.ResponseWriter, request *http.Request) {
		uid := fkutil.ToUint64(request.Form.Get("user_id"))
		logger.SetUid(uid)
		fixCnt, e := FixAssembleEquipInfo(context.TODO(), uid)
		if e == nil {
			writer.Write([]byte(fmt.Sprintf("fixed:%d个", fixCnt)))
		} else {
			writer.Write([]byte(e.Error()))
		}
	})

	gm.SafeHttpRegister(logger, "/GmDressBagEquip", func(writer http.ResponseWriter, request *http.Request) {
		uid := fkutil.ToUint64(request.Form.Get("user_id"))
		logger.SetUid(uid)
		pos := fkutil.ToInt32(request.Form.Get("pos"))
		guid := fkutil.ToInt64(request.Form.Get("guid"))
		e := DressEquipGm(logger, uid, pos, guid)
		if e == nil {
			writer.Write([]byte(string("ok")))
		} else {
			writer.Write([]byte(e.Error()))
		}
	})

	gm.SafeHttpRegister(logger, "/GmEquipPosLvUp", func(writer http.ResponseWriter, request *http.Request) {
		uid := fkutil.ToUint64(request.Form.Get("user_id"))
		logger.SetUid(uid)
		targetLv := fkutil.ToInt32(request.Form.Get("lv"))
		e := equip.OnGmEquipPosLvUp(context.TODO(), uid, targetLv)
		if e == nil {
			writer.Write([]byte(string("ok")))
		} else {
			writer.Write([]byte(e.Error()))
		}
	})
}
