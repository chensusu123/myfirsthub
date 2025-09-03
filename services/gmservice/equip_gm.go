package gmservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maze_game_server/common/tradeno"
	"maze_game_server/config/GMazeEquipInfoV8Cfg"
	"maze_game_server/io/redis/mazeequipgetnumredis"
	"maze_game_server/io/redis/mazeuserlevelredis"
	"maze_game_server/io/rpc/dollequipbagrpc"
	"maze_game_server/model/gmmodel"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/servers/maze_main_server/process/equip"
	"maze_game_server/servers/maze_main_server/process/equip_gm/equipaassemblegm"
	"maze_game_server/servers/maze_main_server/process/equip_gm/equipbaggm"
	"net/http"
	"strings"
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
	logger.CtxInfo(ctx, "GetEquipInfoByCfgId", zap.Any("cond", cd))
	rs, err := equipaassemblegm.GetEquipInfoByCfgId(ctx, userId, cd)
	if err != nil {
		logger.CtxError(ctx, "GetEquipInfoByCfgId", zap.Error(err))
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
	rs, err := equipaassemblegm.GetEquipInfoByGuid(ctx, userId, guid)
	if err != nil {
		logger.CtxError(ctx, "GetEquipInfoByGuid", zap.Error(err))
		writer.Write([]byte("执行失败"))
		return
	}
	writer.Write([]byte(rs))
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

func (s *service) ReInitDollEquip(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	// logger := fklog.ContextAppLogger(ctx)

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
	e = equip.InitDollEquipSuitSeq(ctx, uid)
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

func (s *service) FixEquipPosUnlock(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	uid := fkutil.ToUint64(request.Form.Get("user_id"))
	logger.SetUid(uid)
	cnt, e := equipaassemblegm.UnlockPosByEquip(ctx, uid)
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
	e := equipaassemblegm.DressEquipGm(ctx, uid, pos, guid)
	if e == nil {
		writer.Write([]byte(string("ok")))
	} else {
		writer.Write([]byte(e.Error()))
	}
}

func (s *service) GmEquipPosLvUp(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	var outPut gmmodel.Output
	defer func() {
		jsonOut, err := json.Marshal(outPut)
		if err != nil {
			logger.CtxError(ctx, "Post: /GmEquipPosLvUp  Marshal Fail",
				zap.Any("request", request),
				zap.Any("ouput", outPut),
				zap.Error(err),
			)
		}
		writer.Write(jsonOut)
	}()

	uid := fkutil.ToUint64(request.Form.Get("user_id"))
	logger.SetUid(uid)
	targetLv := fkutil.ToInt32(request.Form.Get("lv"))
	e := equip.OnGmEquipPosLvUp(ctx, uid, targetLv)
	if e == nil {
		outPut = *gmmodel.NewOutPut(http.StatusOK, "操作成功", gmmodel.DynamicData{})
	} else {
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, e.Error(), gmmodel.DynamicData{})
	}
}

func (s *service) AddEquip(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	var outPut gmmodel.Output
	defer func() {
		jsonOut, err := json.Marshal(outPut)
		if err != nil {
			logger.CtxError(ctx, "Post: /AddEquip  Marshal Fail",
				zap.Any("request", request),
				zap.Any("ouput", outPut),
				zap.Error(err),
			)
		}
		writer.Write(jsonOut)
	}()

	userId := fkutil.ToUint64(request.Form.Get("user_id"))
	equipId := fkutil.ToInt32(request.Form.Get("equipId"))
	ruleId := fkutil.ToInt32(request.Form.Get("ruleId"))
	subType := fkutil.ToInt32(request.Form.Get("subType"))
	headType := fkutil.ToInt32(request.Form.Get("headType"))
	tailType := fkutil.ToInt32(request.Form.Get("tailType"))

	if subType > 6 {
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, "subType 子类型无效", gmmodel.DynamicData{})
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
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, fmt.Sprintf("errMsg: %s", err.Error()), gmmodel.DynamicData{})
		return
	}

	outPut = *gmmodel.NewOutPut(http.StatusOK, "操作成功", gmmodel.DynamicData{})
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

	err := mazeequipgetnumredis.SetEquipGetNum(ctx, uid, cfg.Score_group, score)
	if err != nil {
		logger.CtxError(ctx, "SetEquipRollScore  SetEquipGetNum fail", zap.Error(err), zap.Uint64("uid", uid),
			zap.Int32("equipId", equipId), zap.Int32("score", score))
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Write([]byte("ok"))
	return
}

func (s *service) BatchAddEquip(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	var outPut gmmodel.Output
	defer func() {
		jsonOut, err := json.Marshal(outPut)
		if err != nil {
			logger.CtxError(ctx, "Post: /AddItem  Marshal Fail",
				zap.Any("request", request),
				zap.Any("ouput", outPut),
				zap.Error(err),
			)
		}
		writer.Write(jsonOut)
	}()

	uid := fkutil.ToUint64(request.Form.Get("user_id"))
	param := request.Form.Get("equips")
	rp := request.Form.Get("rules")
	if uid <= 0 {
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, "user_id不合法", gmmodel.DynamicData{})
		return
	}

	if param == "" {
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, "请指定装备参数", gmmodel.DynamicData{})
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
			cfg := GMazeEquipInfoV8Cfg.GetWithCtx(ctx, equipId)
			if cfg == nil {
				outPut = *gmmodel.NewOutPut(http.StatusBadGateway, fmt.Sprintf("equipId(%d)找不到对应的装备配置", equipId), gmmodel.DynamicData{})
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
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, "一次添加装备太多,最多100件", gmmodel.DynamicData{})
		return
	}

	req.EquipList = append(req.EquipList, equipConds...)
	res := &MazeEquipSvr.SvrAddMazeEquipRS{}
	req.OpType = proto.Int32(int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_FOE))
	req.TradeNumber = proto.Uint64(tradeno.GetTradeNum())
	err := dollequipbagrpc.MazeBagAddRQ(ctx, req, res)
	if err != nil {
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, fmt.Sprintf("errMsg: %s", err.Error()), gmmodel.DynamicData{})
		return
	}

	outPut = *gmmodel.NewOutPut(http.StatusOK, "操作成功", gmmodel.DynamicData{})
}

func (s *service) FixAllEquipAttrLimit(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	userId := fkutil.ToUint64(request.Form.Get("user_id"))
	logger.SetUid(userId)
	e := equipbaggm.FixAllEquipAttrLimit(ctx, userId)
	if e == nil {
		writer.Write([]byte("ok"))
	} else {
		writer.Write([]byte("fail"))
	}
}
