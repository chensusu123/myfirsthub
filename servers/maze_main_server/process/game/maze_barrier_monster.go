package game

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/common/errors"
	"maze_game_server/common/function/addequip"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/function/maputil"
	"maze_game_server/common/tradeno"
	"maze_game_server/config/GMazeItemsV8Cfg"
	"maze_game_server/io/kafka/dollmazefoekafka"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/server/MazeEquipSvr"
	"maze_game_server/services/barrierservice"
	"maze_game_server/services/itemservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (g *Game) OnBarrierMonsterDeathRQ_10498_10499(s *session.Session, req *MazeGame.BarrierMonsterDeathRQ) (err error) {
	defer fkprometheus.InfoPMT("OnBarrierMonsterDeathRQ")()

	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeGame.BarrierMonsterDeathRS{}

	logger.CtxInfo(ctx, "OnBarrierMonsterDeathRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnBarrierMonsterDeathRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.BarrierId = req.BarrierId
	res.MonsterId = req.MonsterId
	res.OpData = req.OpData

	userId := uint64(s.UID())
	barrierID := req.GetBarrierId()
	monsterID := req.GetMonsterId()
	opData := req.GetOpData()

	// 击杀守卫后，获取守卫死亡奖励
	kongfu, equips, items, errinfo := barrierservice.Global.GuardDeath(ctx, userId, barrierID, int32(monsterID), int32(req.GetMonsterGuid()))
	if errinfo.GetErrCode() != errors.NO_ERROR_CODE {
		res.ErrInfo = errinfo
		logger.CtxError(ctx, "OnBarrierMonsterDeathRQ GuardDeath fail", zap.Error(fmt.Errorf("GuardDeath: %s", errinfo.GetErrMsg())), zap.Any("monsterID", monsterID))
		return
	}

	// 通关值
	res.Kongfu = proto.Int32(kongfu)

	tradeNo := tradeno.GetTradeNum()

	// TODO 使用equip,item服务增加奖励

	// 怪物掉落装备
	if len(equips) > 0 {
		_, err = addequip.AddEquipToBagWithOpdata(ctx, userId, int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_MONSTER_DEATH_AWARD), opData, tradeNo, equips)
		if err != nil {
			logger.CtxError(ctx, "OnBarrierMonsterDeathRQ addEquipToBag fail", zap.Error(err), zap.Any("optype", int32(MazeEquipSvr.ENUM_EQUIP_BAG_OP_TYPE_MAZE_EQUIP_FOE)),
				zap.Any("tradeNo", tradeNo), zap.Any("addEquip", equips))
		}
	}

	var bagItems []*MazeCommon.MazeItem
	// 增加掉落物品返回
	for itemID, count := range items {
		if itemID > 0 {
			itemCfg := GMazeItemsV8Cfg.GetWithCtx(ctx,itemID)
			if itemCfg == nil {
				logger.CtxInfo(ctx, "OnBarrierMonsterDeathRQ item not found", zap.Error(fmt.Errorf("item: %d not found", itemID)), zap.Any("MonsterId", req.GetMonsterId()))
			} else {
				res.Awards = append(res.Awards, &MazeCommon.MazeItem{
					ItemId: proto.Int32(itemID),
					Count:  proto.Int64(count),
				})
				// 背包道具
				if itemCfg.Is_bag == 3 {
					bagItems = append(bagItems, &MazeCommon.MazeItem{
						ItemId: proto.Int32(itemID),
						Count:  proto.Int64(count),
					})
				}
			}
		}
	}
	// 处理需要加入背包的道具
	if len(bagItems) > 0 {
		itemList := itemutil.ItemPb2ItemInfo(bagItems)
		errInfo := itemservice.GlobalItemService.AddItem(context.TODO(), userId, itemservice.ItemOpTypeMonsterDeath, tradeNo, itemList...)
		if errInfo != nil {
			logger.CtxInfo(ctx, "OnBarrierMonsterDeathRQ AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("bagItems", bagItems))
		}
	}

	// 打怪流水记录
	record := &dollmazefoekafka.DollMazeFoeRecord{
		UserId:   userId,
		Barrier:  req.GetBarrierId(),
		MasterId: req.GetMonsterId(),
		Equips:   maputil.MapToString32(equips),
	}
	awards, err := json.Marshal(res.Awards)
	if err != nil {
		logger.CtxError(ctx, "OnBarrierMonsterDeathRQ json marshal fail", zap.Error(err), zap.Any("res", res))
	}
	record.AwardList = string(awards)
	dollmazefoekafka.PushDollMazeFoeRecord(ctx, record)
	// flowrecord.SaveFoeRecord(logger, record)

	return
}
