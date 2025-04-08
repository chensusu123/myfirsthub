package process

import (
	"sort"

	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/addequip"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/uniqueid"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazeshopseqredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeshopmodule"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeBarriesV8Cfg"
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeShopV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/plate/protodef/Common"
	"gitlab.ifreetalk.com/plate/protodef/DollMazeShop"
	"go.uber.org/zap"
)

// 获取人偶迷宫商店物品列表请求
func OnGetMazeShopListRQ(ctx fknet.TCPContext, userId uint64, rqMsg, rsMsg proto.Message) (err error) {
	defer fkprometheus.InfoPMT("OnGetMazeShopListRQ")()

	req := rqMsg.(*DollMazeShop.GetMazeShopListRQ)
	res := rsMsg.(*DollMazeShop.GetMazeShopListRS)

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	agent := fkserver.NewUserContext(ctx.Context, uint64(userId), ctx.FKLogI)

	agent.InfoWF("OnGetMazeShopListRQ start", zap.Any("req", req))
	defer func() {
		agent.InfoWF("OnGetMazeShopListRQ end", zap.Any("res", res))
	}()

	userInfo, err := mazeuserinfo.GetUserInfoV2(agent, userId)
	if err != nil {
		agent.ErrorWF("OnGetMazeShopListRQ GetUserInfo fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	barrierId := userInfo.Barrier
	if barrierId != req.GetBarrierId() {
		barrierId = req.GetBarrierId()
	}
	barriesCfg := GMazeBarriesV8Cfg.Get(barrierId)
	if barriesCfg == nil {
		agent.ErrorWF("OnGetMazeShopListRQ curr barrier not found cfg", zap.Any("barrierId", barrierId))
		res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		return
	}
	moneyId := barriesCfg.Currency_id
	level := mazeshopmodule.GetMazeShopLv(int32(userInfo.Level), userInfo.HighArea)
	cfg := GMazeShopV8Cfg.Get(level)
	if cfg == nil {
		agent.ErrorWF("OnGetMazeShopListRQ GMazeShopV8Cfg err",
			zap.Int32("level", level))
		res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		return
	}
	mazeShopInfo, err := mazeshopseqredis.GetMazeShopInfo(agent, agent.UserID, int32(level))
	if err != nil {
		agent.ErrorWF("OnGetMazeShopListRQ GetMazeShopInfo error!", zap.Error(err), zap.Any("req", req))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	shopNumMap := make(map[int32]int32, 0)
	if mazeShopInfo != nil {
		shopNumMap = mazeShopInfo.ShopSlotNum
	}
	itemList := make([]*DollMazeShop.MazeShopInfo, 0)
	for k, v := range cfg.Equip_value_max {
		itemInfo := &DollMazeShop.MazeShopInfo{
			Index:       proto.Int32(k),
			BuyCount:    proto.Int32(shopNumMap[k]),
			MaxCount:    proto.Int32(int32(-1)),
			SingleCount: proto.Int32(1),
			Pos:         proto.Int32(k),
			MaxValue:    proto.Int64(v),
		}
		itemInfo.CostItems = append(itemInfo.CostItems, &Common.Item{
			ItemId: proto.Int32(moneyId),
			Count:  proto.Int64(int64(cfg.Buy_cost)),
		})
		itemList = append(itemList, itemInfo)
	}
	sort.Slice(itemList, func(i, j int) bool {
		return itemList[i].GetIndex() < itemList[j].GetIndex()
	})
	res.ItemList = itemList
	return nil
}

// 获取人偶迷宫商店物品列表请求
func OnMazeShopBuyRQ(ctx fknet.TCPContext, userId uint64, rqMsg, rsMsg proto.Message) (err error) {
	defer fkprometheus.InfoPMT("OnMazeShopBuyRQ")()

	req := rqMsg.(*DollMazeShop.MazeShopBuyRQ)
	res := rsMsg.(*DollMazeShop.MazeShopBuyRS)

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	userCtx := fkserver.NewUserContext(ctx.Context, uint64(userId), ctx.FKLogI)
	userCtx.InfoWF("OnMazeShopBuyRQ start", zap.Any("req", req))
	defer func() {
		userCtx.InfoWF("OnMazeShopBuyRQ end", zap.Any("res", res))
	}()
	itemType := req.GetItem().GetItemId()
	count := req.GetItem().GetCount()
	if itemType <= 0 || count <= 0 {
		userCtx.ErrorWF("OnMazeShopBuyRQ input invalid args", zap.Any("req", req))
		res.ErrInfo = errors.ARGS_NOT_MATCH.ToInfo()
		return err
	}
	userInfo, err := mazeuserinfo.GetUserInfoV2(userCtx, userId)
	if err != nil {
		userCtx.ErrorWF("OnMazeShopBuyRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	barrierId := userInfo.Barrier
	if barrierId != req.GetBarrierId() {
		userCtx.ErrorWF("OnMazeShopBuyRQ barrier not match", zap.Any("barrierId", barrierId),
			zap.Any("barrierReq", req.GetBarrierId()))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("关卡数据不匹配")
		return
	}

	barriesCfg := GMazeBarriesV8Cfg.Get(barrierId)
	if barriesCfg == nil {
		userCtx.ErrorWF("OnMazeShopBuyRQ curr barrier not found cfg", zap.Any("barrierId", barrierId))
		res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		return
	}

	moneyId := barriesCfg.Currency_id
	level := int32(userInfo.Level)
	cfg := GMazeShopV8Cfg.Get(level)
	if cfg == nil {
		userCtx.ErrorWF("OnMazeShopBuyRQ GMazeShopV8Cfg fail",
			zap.Int32("level", level))
		res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		return
	}
	mazeShopInfo, err := mazeshopmodule.GetMazeShopInfo(userCtx, userCtx.UserID, level, userInfo.HighArea)
	if err != nil || mazeShopInfo == nil {
		userCtx.ErrorWF("OnMazeShopBuyRQ HandleMazeShopSeqInit fail",
			zap.Int32("level", level))
		res.ErrInfo = errors.CONFIG_NOT_FOUND.ToInfo()
		return
	}

	//maxCount := cfg.Buy_equip_type[itemType]
	//if maxCount >= 0 {
	//	buyNum := mazeShopInfo.ShopSlotNum[itemType]
	//	if buyNum >= int32(maxCount) {
	//		userCtx.ErrorWF("OnMazeShopBuyRQ buy count to max limit", zap.Any("req", req),
	//			zap.Int32("buyNum", buyNum),
	//			zap.Int64("numMaxLimit", cfg.Buy_equip_type[itemType]))
	//		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("购买数量已达上限")
	//		return
	//	}
	//	if buyNum+int32(count) >= int32(maxCount) {
	//		userCtx.ErrorWF("OnMazeShopBuyRQ buy count to max limit", zap.Any("req", req),
	//			zap.Int32("buyNum", buyNum),
	//			zap.Int64("count", count),
	//			zap.Int64("numMaxLimit", cfg.Buy_equip_type[itemType]))
	//		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("购买数量超过上限")
	//		return
	//	}
	//}
	equipMaxMap, _, err := mazeshopmodule.GetEquipMaxMap(userCtx, userCtx.UserID, level, userInfo.HighArea)
	if err != nil {
		userCtx.ErrorWF("OnMazeShopBuyRQ GetBarrierMaxAreaId fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	equipMap, err := mazeshopmodule.GetMazeSeqEquipId(userCtx, userCtx.UserID, level, userInfo.HighArea, mazeShopInfo, int32(count), itemType, equipMaxMap)
	if err != nil {
		userCtx.ErrorWF("OnMazeShopBuyRQ GetMazeSeqEquipId fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	moneyCount := int64(cfg.Buy_cost) * count
	err = subMoney(userCtx, userCtx.UserID, moneyId, moneyCount, req.GetHeader().GetSession())
	if err != nil {
		userCtx.ErrorWF("OnMazeShopBuyRQ buy count to max limit", zap.Any("req", req),
			zap.Int32("moneyId", moneyId),
			zap.Int64("moneyCount", moneyCount),
			zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return
	}
	tradeNo := uniqueid.GenUniqueIdUInt64()
	rsAdd, err := addequip.AddEquipToBag(userCtx, userCtx.UserID, 18, tradeNo, equipMap)
	if err != nil {
		userCtx.ErrorWF("OnMazeShopBuyRQ buyShipEquipToBag err", zap.Any("req", req),
			zap.Any("equipMap", equipMap),
			zap.Int64("count", count),
			zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("添加装备失败")
		return
	}
	_ = rsAdd
	mazeShopInfo.ShopSlotNum[itemType] += int32(count)
	err1 := mazeshopseqredis.SetMazeShopInfo(userCtx, userCtx.UserID, level, mazeShopInfo)
	if err1 != nil {
		userCtx.ErrorWF("OnMazeShopBuyRQ buyShipEquipToBag err", zap.Any("req", req),
			zap.Int32("level", level),
			zap.Any("mazeShopInfo", mazeShopInfo),
			zap.Error(err1))
	}
	// PushDollMazeShopInfoLog(userCtx, userCtx.UserID, level, itemType, mazeShopInfo, rsAdd.EquipList, 18, tradeNo, 0)
	return nil
}

//func updateMazeShopSeq(logger fklog.FKLogI, mazeShopInfo *dollmazeshopseqredis.MazeShopInfo, buyShopSeq *dollmazeshopseqredis.ShopSeqInfo) (equipId int32, err error) {
//	var seqCfg *GDollMazeShopEquipListV8Cfg.DollMazeShopEquipListV8ConfigRow
//	var backCfg *GDollMazeShopEquipListV8Cfg.DollMazeShopEquipListV8ConfigRow
//	for _, cfg := range GDollMazeShopEquipListV8Cfg.GetAll() {
//		if (cfg.List_id == mazeShopInfo.SeqId) && (buyShopSeq.SlotId == cfg.Pos_id) {
//			seqCfg = cfg
//		}
//		if (cfg.List_id == mazeShopInfo.BackId) && (buyShopSeq.SlotId == cfg.Pos_id) {
//			backCfg = cfg
//		}
//	}
//	if (seqCfg == nil) || (backCfg == nil) {
//		logger.ErrorWF("GetMazeShopSeq GDollMazeShopEquipListV8Cfg err",
//			zap.Any("mazeShopInfo", mazeShopInfo),
//			zap.Any("buyShopSeq", buyShopSeq))
//		return 0, errors.New("配置不存在")
//	}
//	if buyShopSeq.BackIndex > 0 || (int(buyShopSeq.CurSeqIndex) >= len(seqCfg.Equip_id)) {
//		if int(buyShopSeq.BackIndex) >= len(backCfg.Equip_id) {
//			buyShopSeq.BackIndex = 0
//		}
//		buyShopSeq.BackIndex++
//		buyShopSeq.TotalCount++
//		equipId = backCfg.Equip_id[buyShopSeq.BackIndex-1]
//		return equipId, nil
//	}
//	buyShopSeq.CurSeqIndex++
//	buyShopSeq.TotalCount++
//	equipId = seqCfg.Equip_id[buyShopSeq.CurSeqIndex-1]
//	return equipId, nil
//}
