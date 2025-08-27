package gmservice

import (
	"fmt"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/common/tradeno"
	"maze_game_server/config/GMazeItemsV8Cfg"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/services/barrierenergyservice"
	"maze_game_server/services/itemservice"
	"net/http"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *service) AddRefreshCost(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	request.ParseForm()
	userId := fkutil.ToUint64(request.Form.Get("user_id"))
	itemId := fkutil.ToInt32(request.Form.Get("itemId"))
	count := fkutil.ToInt64(request.Form.Get("count"))
	item := &itemservice.ItemInfo{
		ItemId: itemId,
		Count:  count,
	}
	err := itemservice.GlobalItemService.AddItem(ctx, userId, itemservice.ItemOpTypeGM, tradeno.GetTradeNum(), item)
	if err != nil {
		_, _ = writer.Write([]byte(fmt.Sprintf("add cost failed itemId: %d count:%d", itemId, count)))
	}

	_, _ = writer.Write([]byte(fmt.Sprintf("add success itemId: %d count:%d", itemId, count)))
}

func (s *service) AddExp(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	userId := fkutil.ToUint64(request.Form.Get("user_id"))
	exp := fkutil.ToInt64(request.Form.Get("exp"))
	logger.SetLogId(time.Now().UnixNano())

	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "AddExp GetUserInfoV2 fail", zap.Error(err))
		writer.Write([]byte(err.Error()))
		return
	}

	oldLevel := userInfo.Level
	oldExp := userInfo.TotalExp

	// 更新等级经验
	err = userInfo.AddExp(exp)
	if err != nil {
		logger.CtxError(ctx, "AddExp CalExp fail", zap.Error(err))
		writer.Write([]byte(err.Error()))
		return
	}
	err = mazeuserinfo.SetUserInfoV2(ctx, userId, userInfo)
	if err != nil {
		logger.CtxError(ctx, "AddExp SetUserInfoV2 fail", zap.Error(err))
		writer.Write([]byte(err.Error()))
		return
	}
	mazecommonvalue.HandleUserLevelExpChg(logger, userId, userInfo.Level, userInfo.Exp, "")

	defer func() {
		if exp != 0 {
			levelRecord := &mazeuserlevelkafka.MazeUserLevelRecord{
				UserId:      userId,
				OldLevel:    int32(oldLevel),
				OldTotalExp: oldExp,
				NewLevel:    int32(userInfo.Level),
				NewTotalExp: int32(userInfo.TotalExp),
			}
			mazeuserlevelkafka.PushMazeLevelRecord(ctx, levelRecord)
		}
	}()

	writer.Write([]byte("设置成功，注意尽量不要在迷宫杀怪时使用本gm"))
}

func (s *service) AddEnergy(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	var (
		userId = fkutil.ToUint64(request.Form.Get("user_id"))
		count  = fkutil.ToInt32(request.Form.Get("count"))
	)
	if count > barrierenergyservice.GlobalBarrierEnergyService.GetEnergyMaxValue() {
		return
	}
	curEnergy, _, err := barrierenergyservice.GlobalBarrierEnergyService.AddEnergy(ctx, userId, count)
	if err != nil {
		return
	}
	fmt.Fprintf(writer, "add barrier energy success, curEnergy=[%d]", curEnergy)
}

func (s *service) AddItem(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)
	var (
		userId = fkutil.ToUint64(request.Form.Get("user_id"))
		itemId = fkutil.ToInt32(request.Form.Get("itemId"))
		count  = fkutil.ToInt64(request.Form.Get("count"))
	)

	tradeNo := tradeno.GetTradeNum()
	items := make([]*MazeCommon.MazeItem, 0)

	if userId <= 0 {
		fmt.Fprintf(writer, "请指定有效用户ID")
		return
	}

	if itemId <= 0 || count <= 0 {
		fmt.Fprintf(writer, "无效道具ID或道具数量")
		return
	}

	itemCfg := GMazeItemsV8Cfg.GetWithCtx(ctx, itemId)
	if itemCfg == nil {
		fmt.Fprintf(writer, "无效道具，请检查道具配置表：maze_items_v8【迷宫-道具】.xlsx")
		return
	}

	items = append(items, &MazeCommon.MazeItem{
		Count:  proto.Int64(count),
		ItemId: proto.Int32(itemId),
	})

	header := &Common.PacketHeader{}
	header.Sharding = proto.Int64(int64(userId))

	itemList := itemutil.ItemPb2ItemInfo(items)
	errInfo := itemservice.GlobalItemService.AddItem(ctx, userId, itemservice.ItemOpTypeGM, tradeNo, itemList...)
	if errInfo != nil {
		fmt.Fprintf(writer, "添加道具失败，错误：%s", string(errInfo.GetErrMsg()))
		logger.CtxError(ctx, "addItem AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("ItemList", items))
		return
	}

	fmt.Fprintf(writer, "添加[%d]个道具[%s]成功", count, itemCfg.Prop_name)
}
