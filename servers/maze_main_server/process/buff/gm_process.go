package buff

import (
	"fmt"
	"maze_game_server/io/kafka/mazetempbuffchgmsg"
	"maze_game_server/io/redis/mazetempbuffredis"
	"maze_game_server/module/itemmodule"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/server/MazeTempBuffSvr"
	"net/http"
	"sort"
	"strings"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/31 16:30
 * @Description:
 */

func SafeHttpRegister(logger fklog.FKLogI, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	appConfig := appconfig.GlobalConfig()
	// /s4/AddExp
	pattern = "/s" + appConfig.Global.SectionID + pattern
	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		defer fkutil.CaptureException()
		logger.DebugWF("execute gm", zap.String("pattern", pattern), zap.Any("header", request.Header),
			zap.Any("host", request.Host), zap.Any("remoteAddr", request.RemoteAddr))
		// if !CheckGM.CheckGMOnline(context.Background(), logger, 10000, pattern, request.RemoteAddr) {
		// 	return
		// }
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

		var errs []string
		var successList, failedList []string
		for _, buffId := range buffList {
			buffWeight, err := getOptionBuffWeightInfo(buffId, buffMap)
			if err != nil {
				failedList = append(failedList, fmt.Sprintf("%d", buffId))
				errs = append(errs, err.Error())
				continue
			}

			_ = buffWeight

			buffInfo.SelectedBuff = append(buffInfo.SelectedBuff, &MazeTempBuffSvr.SelectedBuffInfo{
				BuffId: proto.Int32(buffId),
			})
			successList = append(successList, fmt.Sprintf("%d", buffId))
			buffMap[buffId] += 1
		}

		if len(successList) == 0 {
			logger.WarnWF("setMazeTempBuff optionalList is nil")
			_, _ = writer.Write([]byte("not have optional buff, failed buff:" + strings.Join(failedList, ",") + " errs:" + strings.Join(errs, ",")))
			return
		}

		var totalMap map[int32]int64
		totalMap, buffInfo.TotalBuff = GetTotalBuff(logger, buffInfo.GetSelectedBuff())
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

		// // buff中心
		// forceAttr, err := GetSelectBuffForceAttr(buffInfo.TotalBuff)
		// if err != nil {
		// 	logger.ErrorWF("setMazeTempBuff GetSelectBuffForceAttr failed", zap.Error(err))
		// 	_, _ = writer.Write([]byte("\n更新人物属性失败:" + err.Error()))
		// 	return
		// }
		// attrDb := &MazeBuffData.MazeBuffDb{
		// 	MazeRealBuffs: PackMazeBuff(forceAttr),
		// 	// MazeShowBuffs: PackMazeBuff(showBuff), // todo 现在暂时没有展示武力值
		// }

		// err = mazebuffinforedis.SaveMazeBuffInfo(logger, userId, constdef.MazeBuffSrcSelectBuffForce, attrDb)
		// if err != nil {
		// 	logger.ErrorWF("setMazeTempBuff SaveMazeBuffInfo failed", zap.Uint64("userId", userId), zap.Error(err))
		// 	_, _ = writer.Write([]byte("\n更新人物属性失败:" + err.Error()))
		// 	return
		// }

		// // 推送属性计算消息
		// calcAttrNotify := &structsdef.MazeCalcAttrNotifyMsg{
		// 	UserId: userId,
		// 	// FromServer: fmt.Sprintf("%d %s", fkconfig.EnvVal.ServerType, fkconfig.EnvVal.AppName),
		// 	ChgType: constdef.MazeBuffChgForceValue,
		// 	Session: "buff",
		// 	BuffSrc: constdef.MazeBuffSrcSelectBuffForce,
		// }
		// err = mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, calcAttrNotify)
		// if err != nil {
		// 	_, _ = writer.Write([]byte("\n更新人物属性失败:" + err.Error()))
		// }
		return
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
