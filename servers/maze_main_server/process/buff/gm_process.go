package buff

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazetempbuffchgmsg"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazetempbuffredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/itemmodule"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"gitlab.ifreetalk.com/maze-plate/io/redis_interface/common/CheckGM"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeTempBuffSvr"
	"go.uber.org/zap"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/31 16:30
 * @Description:
 */

func SafeHttpRegister(logger fklog.FKLogI, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		defer fkutil.CaptureException()

		logger.DebugWF("execute gm", zap.String("pattern", pattern), zap.Any("header", request.Header),
			zap.Any("host", request.Host), zap.Any("remoteAddr", request.RemoteAddr))
		if !CheckGM.CheckGMOnline(context.Background(), logger, 10000, pattern, request.RemoteAddr) {
			return
		}
		logger.WarnWF("execute gm", zap.String("pattern", pattern), zap.Any("header", request.Header),
			zap.Any("host", request.Host), zap.Any("remoteAddr", request.RemoteAddr))

		handler(writer, request)
	})
}

func InitGM(logger fklog.FKLogI) {
	// 设置buff
	SafeHttpRegister(logger, "/setMazeTempBuff", func(writer http.ResponseWriter, request *http.Request) {
		// 外网线上环境不允许使用GM
		request.ParseForm()
		userId := fkutil.ToUint64(request.Form.Get("userid"))
		stageId := fkutil.ToInt32(request.Form.Get("stageId"))
		buffs := request.Form.Get("buffs")
		logger.SetUid(userId)
		var buffList []int32
		for _, str := range strings.Split(buffs, ",") {
			if buff := fkutil.ToInt32(str); buff != 0 {
				buffList = append(buffList, buff)
			}
		}

		sort.Slice(buffList, func(i, j int) bool {
			return buffList[i] < buffList[j]
		})

		if userId == 0 || stageId == 0 || len(buffs) == 0 {
			logger.WarnWF("setMazeTempBuff args is error")
			_, _ = writer.Write([]byte("set maze temp buff args is error"))
			return
		}

		logger.InfoWF("setMazeTempBuff start", zap.Uint64("userId", userId),
			zap.Int32("stageId", stageId), zap.Int32s("buffList", buffList))

		buffInfo, err := mazetempbuffredis.GetMazeTempBuff(logger, userId, stageId)
		if err != nil {
			logger.ErrorWF("setMazeTempBuff GetMazeTempBuff", zap.Error(err))
			_, _ = writer.Write([]byte("get user buff failed"))
			return
		}

		if buffInfo == nil {
			buffInfo = &MazeTempBuffSvr.TempBuffInfo{
				BuffSequence: &MazeTempBuffSvr.BuffSequence{
					Index: proto.Int32(1),
				},
			}
		}

		// 校验选择的buff
		buffMap := make(map[int32]int32)
		for _, info := range buffInfo.GetSelectedBuff() {
			buffMap[info.GetBuffId()] += 1
		}

		var successList, failedList []string
		for _, buffId := range buffList {
			buffWeight := getOptionBuffWeightInfo(buffId, buffMap)
			if buffWeight == nil {
				failedList = append(failedList, fmt.Sprintf("%d", buffId))
				continue
			}

			buffInfo.SelectedBuff = append(buffInfo.SelectedBuff, &MazeTempBuffSvr.SelectedBuffInfo{
				BuffId: proto.Int32(buffId),
			})
			successList = append(successList, fmt.Sprintf("%d", buffId))
			buffMap[buffId] += 1
		}

		if len(successList) == 0 {
			logger.WarnWF("setMazeTempBuff optionalList is nil")
			_, _ = writer.Write([]byte("not have optional buff, failed buff:" + strings.Join(failedList, ",")))
			return
		}

		var totalMap map[int32]int64
		totalMap, buffInfo.TotalBuff = getTotalBuff(logger, buffInfo.GetSelectedBuff())
		// 更新buff信息
		err = mazetempbuffredis.SetMazeTempBuff(logger, userId, stageId, buffInfo)
		if err != nil {
			logger.ErrorWF("setMazeTempBuff SetMazeTempBuff failed", zap.Any("info", buffInfo), zap.Error(err))
			_, _ = writer.Write([]byte("save buff failed"))
			return
		}

		// 推送buff变化信息
		msg := &mazetempbuffchgmsg.MazeTempBuffChangeMsg{
			UserId:  userId,
			StageId: stageId,
			ChgType: 1,
			ChgDesc: "gm添加buff",
		}

		// 计算buff变化
		chgAttrs := make([]*mazetempbuffchgmsg.AttrChgInfo, 0, len(totalMap))
		for id, value := range totalMap {
			chgAttrs = append(chgAttrs, &mazetempbuffchgmsg.AttrChgInfo{
				AttrId: id,
				OldVal: value,
				CurVal: value,
			})
		}

		msg.ChgAttrs = chgAttrs
		_ = mazetempbuffchgmsg.PushTempBuffChangeMsg(logger, msg)
		_, _ = writer.Write([]byte("set success buff:" + strings.Join(successList, ",")))
		if len(failedList) > 0 {
			_, _ = writer.Write([]byte("failed buff:" + strings.Join(failedList, ",")))
		}
	})

	SafeHttpRegister(logger, "/addRefreshCost", func(writer http.ResponseWriter, request *http.Request) {
		// 外网线上环境不允许使用GM
		request.ParseForm()
		userId := fkutil.ToUint64(request.Form.Get("userid"))
		itemId := fkutil.ToInt32(request.Form.Get("itemId"))
		count := fkutil.ToInt64(request.Form.Get("count"))
		logger.SetUid(userId)
		items := []*MazeCommon.MazeItem{
			{
				ItemId: proto.Int32(itemId),
				Count:  proto.Int64(count),
			},
		}
		err := itemmodule.AddItems(logger, userId, itemmodule.CostRefreshType, items)
		if err != nil {
			_, _ = writer.Write([]byte(fmt.Sprintf("add cost failed itemId: %d count:%d", itemId, count)))
		}

		_, _ = writer.Write([]byte(fmt.Sprintf("add success itemId: %d count:%d", itemId, count)))
	})
}
