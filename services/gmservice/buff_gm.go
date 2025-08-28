package gmservice

import (
	"encoding/json"
	"fmt"
	"maze_game_server/config/GMazeAttributeV8Cfg"
	"maze_game_server/config/GMazeEnergyAffixV8Cfg"
	"maze_game_server/io/kafka/mazetempbuffchgmsg"
	"maze_game_server/model/gmmodel"
	"maze_game_server/model/tempbuffmodel"
	"maze_game_server/services/tempbuffservice"
	"net/http"
	"sort"
	"strings"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

func (s *service) SetMazeTempBuff(writer http.ResponseWriter, request *http.Request) {
	// 外网线上环境不允许使用GM
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

	userId := fkutil.ToUint64(request.Form.Get("user_id"))
	barrierId := fkutil.ToInt32(request.Form.Get("stageId"))
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

	if userId == 0 || barrierId == 0 || len(buffs) == 0 {
		logger.CtxError(ctx, "setMazeTempBuff args is error")
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, "set maze temp buff args is error", gmmodel.DynamicData{})
		return
	}

	logger.CtxInfo(ctx, "setMazeTempBuff start", zap.Uint64("userId", userId),
		zap.Int32("barrierId", barrierId), zap.Int32s("buffList", buffList))
	buffInfo, err := tempbuffservice.GlobalTempBuffService.GetTempBuffInfo(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "setMazeTempBuff GetMazeTempBuff", zap.Error(err))
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, "get user buff failed", gmmodel.DynamicData{})
		return
	}

	if buffInfo == nil || buffInfo.BuffSequence == nil {
		buffInfo = &tempbuffmodel.TempBuffInfoModel{
			BuffSequence: &tempbuffmodel.BuffSequence{
				Level: 1,
			},
		}
	}

	// 统计词条已经选择的次数
	selectedBuffMap := make(map[int32]int32)
	for _, info := range buffInfo.SelectedBuff {
		selectedBuffMap[info.BuffId] += 1
	}
	// 统计词条组以及选择的次数
	selectedBuffGroupMap := make(map[int32]int32)
	for _, i := range buffInfo.SelectedBuff {
		buffConfig := GMazeEnergyAffixV8Cfg.GetWithCtx(ctx, i.BuffId)
		if buffConfig == nil {
			logger.CtxWarn(ctx, "createOptionalBuffList GetAffixConfig is nil", zap.Int32("buffId", i.BuffId))
			continue
		}
		selectedBuffGroupMap[buffConfig.Affix_group_id] += 1
	}

	var errs []string
	var successList, failedList []string
	for _, buffId := range buffList {
		buffWeight := tempbuffservice.GlobalTempBuffService.GetOptionBuffWeightInfo(ctx, buffId, selectedBuffMap, selectedBuffGroupMap, map[int32]struct{}{}, 0)
		if err != nil {
			failedList = append(failedList, fmt.Sprintf("%d", buffId))
			errs = append(errs, err.Error())
			continue
		}

		_ = buffWeight

		buffInfo.SelectedBuff = append(buffInfo.SelectedBuff, &tempbuffmodel.SelectedBuffInfo{
			BuffId: buffId,
		})
		successList = append(successList, fmt.Sprintf("%d", buffId))
		selectedBuffMap[buffId] += 1
	}

	if len(successList) == 0 {
		logger.CtxWarn(ctx, "setMazeTempBuff optionalList is nil")
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, "not have optional buff, failed buff:"+strings.Join(failedList, ",")+" errs:"+strings.Join(errs, ","), gmmodel.DynamicData{})
		return
	}

	// 校验属性配置
	for _, info := range buffInfo.SelectedBuff {
		// 获取buff实际加成
		config := GMazeEnergyAffixV8Cfg.GetWithCtx(ctx, info.BuffId)
		if config != nil {
			for id := range config.Add_attr {
				attrCfg := GMazeAttributeV8Cfg.Get(id)
				if attrCfg == nil {
					outPut = *gmmodel.NewOutPut(http.StatusBadGateway, fmt.Sprintf("词条[%d]增加的属性[%d]配置无效，请检查属性配置表: maze_attribute_v8【迷宫-属性】.xlsx", info.BuffId, id), gmmodel.DynamicData{})
					return
				}
			}
		}
	}

	var totalMap map[int32]int64
	totalMap, buffInfo.TotalBuff = tempbuffservice.GlobalTempBuffService.GetTotalBuff(ctx, buffInfo.SelectedBuff)

	// 更新buff信息
	err = buffInfo.Save(ctx, userId, barrierId)
	if err != nil {
		logger.CtxError(ctx, "setMazeTempBuff SetMazeTempBuff failed", zap.Any("info", buffInfo), zap.Error(err))
		outPut = *gmmodel.NewOutPut(http.StatusBadGateway, "save buff failed", gmmodel.DynamicData{})
		return
	}

	// 推送buff变化信息
	msg := &mazetempbuffchgmsg.MazeTempBuffChangeMsg{
		UserId:  userId,
		StageId: barrierId,
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
	_ = mazetempbuffchgmsg.PushTempBuffChangeMsg(ctx, msg)

	logger.CtxInfo(ctx, "set success buff:"+strings.Join(successList, ","))
	if len(failedList) > 0 {
		logger.CtxWarn(ctx, "failed buff:"+strings.Join(failedList, ","))
	}
	outPut = *gmmodel.NewOutPut(http.StatusOK, "操作成功", gmmodel.DynamicData{})

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
}
