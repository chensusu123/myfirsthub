package gmservice

import (
	"bytes"
	"fmt"
	"maze_game_server/common/function/fileio"
	"maze_game_server/common/tradeno"
	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/io/redis/mazebuffinforedis"
	"maze_game_server/io/redis/mazeequipgetnumredis"
	"maze_game_server/io/redis/mazeuserlevelredis"
	"maze_game_server/io/rpc/dollequipbagrpc"
	"maze_game_server/module/calcassembleattr"
	"maze_game_server/module/dollassembleinfo"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/servers/maze_main_server/process/equip"
	"maze_game_server/servers/maze_main_server/process/equip_gm/equipaassemblegm"
	"maze_game_server/servers/maze_main_server/process/equip_gm/equipbaggm"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *service) GetEquipInfoByCfgId(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)
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
		writer.Write([]byte(equipaassemblegm.CondHelp()))
		return
	}
	cd := equipaassemblegm.BagCond{}
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
	rs, err := equipaassemblegm.GetEquipInfoByCfgId(logger, userId, cd)
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
}

func (s *service) GetEquipInfoByGuid(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)
	userId := fkutil.ToUint64(request.Form.Get("user_id"))
	guid := fkutil.ToInt64(request.Form.Get("guid"))
	rs, err := equipaassemblegm.GetEquipInfoByGuid(logger, userId, guid)
	if err != nil {
		logger.ErrorWF("GetEquipInfoByGuid", zap.Error(err))
		writer.Write([]byte("执行失败"))
		return
	}
	writer.Write([]byte(rs))
}

func (s *service) FixDollEquipAttr(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)
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
}

func (s *service) SendOneSuitEquip(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	userId := fkutil.ToUint64(request.Form.Get("user_id"))
	dressLv := fkutil.ToInt32(request.Form.Get("dressLv"))
	LvDis := fkutil.ToInt32(request.Form.Get("LvDis"))
	suitId := fkutil.ToInt32(request.Form.Get("suitId"))
	quality := fkutil.ToInt32(request.Form.Get("quality"))
	pos := fkutil.ToInt32(request.Form.Get("pos"))
	subType := fkutil.ToInt32(request.Form.Get("subType"))
	p := equipaassemblegm.EquipParam{}
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
	e := equipaassemblegm.CheckEquipParam(&p)
	if e != nil {
		writer.Write([]byte(e.Error()))
		return
	}
	equips, err := equipaassemblegm.AddEquipByCond(ctx, userId, p)
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
}

func (s *service) ReInitDollEquipByFile(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

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
		e := equipbaggm.ClearUserBag(ctx, uid)
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
		e = equip.HandleDollEquipInit(ctx, uid, true)
		if e != nil {
			atomic.AddInt32(&fail, 1)
			return true
		}
		return true
	})
	_, _ = writer.Write([]byte(fmt.Sprintf("执行结果=total:%d succ:%d fail:%d", total, succ, fail)))

}

func (s *service) ReInitDollEquip(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	uid := fkutil.ToUint64(request.Form.Get("user_id"))
	clear := fkutil.ToBool(request.Form.Get("clear"))
	if clear {
		e := equipbaggm.ClearUserBag(ctx, uid)
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
	e = equip.HandleDollEquipInit(ctx, uid, true)
	if e != nil {
		writer.Write([]byte(fmt.Sprintf("初始化装备失败:%s", e.Error())))
		return
	}
	writer.Write([]byte("初始化完成"))
}

func (s *service) FixDollAttr(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)
	uid := fkutil.ToUint64(request.Form.Get("user_id"))
	fixType := fkutil.ToInt32(request.Form.Get("fixType"))
	logger.SetUid(uid)
	e := equipaassemblegm.ReCalcDollEquipAttr(ctx, uid, fixType)
	if e != nil {
		writer.Write([]byte(fmt.Sprintf("执行结果:%s", e.Error())))
	} else {
		writer.Write([]byte("执行完成"))
	}
}

func (s *service) FixDollAttrByFile(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

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

		e := equipaassemblegm.ReCalcDollEquipAttr(ctx, uid, fixType)
		if e == nil {
			atomic.AddInt32(&succ, 1)
		} else {
			atomic.AddInt32(&fail, 1)
		}

		return true
	})
	_, _ = writer.Write([]byte(fmt.Sprintf("执行结果=total:%d succ:%d fail:%d", total, succ, fail)))

}

func (s *service) FixEquipPosUnlock(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	uid := fkutil.ToUint64(request.Form.Get("user_id"))
	logger.SetUid(uid)
	cnt, e := equipaassemblegm.UnlockPosByEquip(logger, uid)
	if e == nil {
		writer.Write([]byte(fmt.Sprintf("unlock:%d", cnt)))
	} else {
		writer.Write([]byte(e.Error()))
	}
}

func (s *service) FixAssembleEquipInfo(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	uid := fkutil.ToUint64(request.Form.Get("user_id"))
	logger.SetUid(uid)
	fixCnt, e := equipaassemblegm.FixAssembleEquipInfo(ctx, uid)
	if e == nil {
		writer.Write([]byte(fmt.Sprintf("fixed:%d个", fixCnt)))
	} else {
		writer.Write([]byte(e.Error()))
	}
}

func (s *service) GmDressBagEquip(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	uid := fkutil.ToUint64(request.Form.Get("user_id"))
	logger.SetUid(uid)
	pos := fkutil.ToInt32(request.Form.Get("pos"))
	guid := fkutil.ToInt64(request.Form.Get("guid"))
	e := equipaassemblegm.DressEquipGm(logger, uid, pos, guid)
	if e == nil {
		writer.Write([]byte(string("ok")))
	} else {
		writer.Write([]byte(e.Error()))
	}
}

func (s *service) GmEquipPosLvUp(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	uid := fkutil.ToUint64(request.Form.Get("user_id"))
	logger.SetUid(uid)
	targetLv := fkutil.ToInt32(request.Form.Get("lv"))
	e := equip.OnGmEquipPosLvUp(ctx, uid, targetLv)
	if e == nil {
		writer.Write([]byte(string("ok")))
	} else {
		writer.Write([]byte(e.Error()))
	}
}

func (s *service) AddEquip(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	// 外网线上环境不允许使用GM
	request.ParseForm()

	userId := fkutil.ToUint64(request.Form.Get("user_id"))
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
	req.TradeNumber = proto.Uint64(tradeno.GetTradeNum())
	err := dollequipbagrpc.MazeBagAddRQ(ctx, req, res)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write([]byte("ok"))
	return
}

func (s *service) SetEquipRollScore(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	request.ParseForm()

	uid := fkutil.ToUint64(request.Form.Get("user_id"))
	equipId := fkutil.ToInt32(request.Form.Get("equipId"))
	score := fkutil.ToInt32(request.Form.Get("score"))

	if uid == 0 || equipId == 10 || score < 0 {
		writer.Write([]byte("uid 不能为0 , equipId 不能为0, score 不能小于0"))
		return
	}

	cfg := GMazeEquipInfoV8Cfg.GetWithCtx(ctx, equipId)
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
}

func (s *service) BatchAddEquip(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	// 外网线上环境不允许使用GM
	request.ParseForm()

	uid := fkutil.ToUint64(request.Form.Get("user_id"))
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
	rulesMap := equipbaggm.ParseRules(rp)
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
	req.TradeNumber = proto.Uint64(tradeno.GetTradeNum())
	err := dollequipbagrpc.MazeBagAddRQ(ctx, req, res)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write([]byte("ok"))
}

func (s *service) FixAllEquipAttrLimit(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	userId := fkutil.ToUint64(request.Form.Get("user_id"))
	logger.SetUid(userId)
	e := equipbaggm.FixAllEquipAttrLimit(logger, userId)
	if e == nil {
		writer.Write([]byte("ok"))
	} else {
		writer.Write([]byte("fail"))
	}
}
