package gmservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"maze_game_server/config/GMazeAttributeV8Cfg"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeEnergyAffixV8Cfg"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/io/redis/mazefixedbarrierredis"
	"maze_game_server/model/tempbuffmodel"
	"maze_game_server/module/dollassembleinfo"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/servers/maze_main_server/process/equip_gm/equipaassemblegm"
	"maze_game_server/servers/maze_main_server/process/game"
	"maze_game_server/services/tempbuffservice"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/schema"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

var EndLine = "-----------------------------------------------------------\n"
var form = schema.NewDecoder()

type SetBarrierParams struct {
	UserID    uint64 `schema:"user_id,required"`
	BarrierID int32  `schema:"barrierId,required"`
	Lock      int32  `schema:"lock"` // 如果提供这个参数，则设置的关卡会记录下来，重置游戏数据也会继续生效
}

func (s *service) SetBarrier(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	var params SetBarrierParams

	err := form.Decode(&params, request.Form)
	if err != nil {
		writer.Write([]byte("参数不正确"))
		return
	}

	logger.SetLogId(time.Now().UnixNano())

	if params.UserID <= 0 || params.BarrierID <= 0 {
		writer.Write([]byte("参数不正确"))
		return
	}

	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, params.UserID)
	if err != nil {
		logger.ErrorWF("SetBarrier GetUserInfoV2 fail", zap.Error(err))
		writer.Write([]byte(err.Error()))
		return
	}

	barrierCfg := GMazeBarriesV8Cfg.Get(params.BarrierID)
	if barrierCfg == nil {
		errCfg := errors.New("cant find barrier cfg")
		logger.ErrorWF("SetBarrier Get barrier fail", zap.Error(err))
		writer.Write([]byte(errCfg.Error()))
		return
	}

	// oldBarrier := userInfo.Barrier

	// if params.BarrierID <= oldBarrier {
	// 	fmt.Fprintf(writer, "仅支持跳过关卡，当前第%d关", oldBarrier)
	// 	return
	// }

	userInfo.SetBarrier(params.BarrierID)
	userInfo.SetPassBarrier(params.BarrierID - 1)

	// 更新设置关卡
	err = mazeuserinfo.SetUserInfoV2(ctx, params.UserID, userInfo)
	if err != nil {
		logger.ErrorWF("SetBarrier SetUserInfoV2 fail", zap.Error(err))
		writer.Write([]byte(err.Error()))
		return
	}

	// 锁定设置的关卡
	if params.Lock == 1 {
		mazefixedbarrierredis.SetUserFixedBarrierID(logger, params.UserID, params.BarrierID)
	} else {
		mazefixedbarrierredis.DelUserFixedBarrierID(logger, params.UserID)
	}

	writer.Write([]byte("设置成功，注意尽量不要在迷宫杀怪时使用本gm"))
}

func (s *service) DumpBattleData(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)
	userId := fkutil.ToUint64(request.Form.Get("user_id"))
	barrierId := fkutil.ToInt32(request.Form.Get("barrierId"))
	logger.SetLogId(time.Now().UnixNano())
	logger.SetUid(userId)
	logger.InfoWF("DumpBattleData begin")
	battleData, e := game.GetMazeBattleData(logger, userId, barrierId)
	if e != nil {
		writer.Write([]byte(e.Error()))
		return
	}
	var bs bytes.Buffer
	bs.WriteString("战斗数据:\n")
	as, _ := json.Marshal(battleData.GetAreaInfos())
	bs.WriteString(fmt.Sprintf("区域怪物数据:%s\n", string(as)))
	rc, _ := json.Marshal(battleData.GetRoleConfigInfo())
	bs.WriteString(fmt.Sprintf("角色配置数据:%s\n", string(rc)))
	em, _ := json.Marshal(battleData.GetEliteMonsterInfos())
	bs.WriteString(fmt.Sprintf("精英怪数据:%s\n", string(em)))

	tempBuffInfo, err := tempbuffservice.GlobalTempBuffService.GetTempBuffInfo(ctx, userId, barrierId)
	if err != nil {
		logger.ErrorWF("DumpBattleData GetBarrierTempBuff err", zap.Error(err))
	} else {
		tb, _ := json.Marshal(tempBuffInfo.TotalBuff)
		bs.WriteString(fmt.Sprintf("临时buff数据:%s\n", string(tb)))
	}
	attrMap := game.GetAttrType()
	bs.WriteString("属性ID<->枚举映射关系:\n")
	for k, v := range attrMap {
		var attrName string
		attrCfg := GMazeAttributeV8Cfg.Get(k)
		if attrCfg != nil {
			attrName = attrCfg.Name
		}
		bs.WriteString(fmt.Sprintf("%d(%s)->%d\n", k, attrName, v))
	}
	writer.Write(bs.Bytes())
}

func (s *service) Attrs(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

	var (
		userId    = fkutil.ToUint64(request.Form.Get("user_id"))
		barrierId = fkutil.ToInt32(request.Form.Get("barrier_id"))
	)

	userAttrMap, err := mazecalcattrredis.GetAllMazeCalcAttr(logger, userId)
	if err != nil {
		logger.ErrorWF("GetAllMazeCalcAttr nil", zap.Uint64("userId", userId), zap.Error(err))
		fmt.Fprintf(writer, "获取人物属性失败: %s\n", err.Error())
		return
	}

	var tempBuffInfo *tempbuffmodel.TempBuffInfoModel
	if barrierId > 0 {
		tempBuffInfo, err = tempbuffservice.GlobalTempBuffService.GetTempBuffInfo(ctx, userId, barrierId)
		if err != nil {
			logger.ErrorWF("GetBarrierTempBuff err", zap.Error(err))
			fmt.Fprintf(writer, "获取临时BUFF失败: %s\n", err.Error())
			return
		}
		for _, buffInfo := range tempBuffInfo.TotalBuff {
			userAttrMap[buffInfo.BuffId] += buffInfo.BuffValue
		}
	}
	type UserAttr struct {
		AttrID int32
		Value  int64
	}
	attrs := make([]UserAttr, 0, len(userAttrMap))
	for attrID, value := range userAttrMap {
		attrs = append(attrs, UserAttr{AttrID: attrID, Value: value})
	}

	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].AttrID < attrs[j].AttrID
	})

	fmt.Fprintf(writer, "----------------用户属性列表----------------\n")
	for _, attr := range attrs {
		attrCfg := GMazeAttributeV8Cfg.Get(attr.AttrID)
		if attrCfg != nil {
			switch attrCfg.Figure {
			case 1:
				fmt.Fprintf(writer, "[%d]%s: %d\n", attr.AttrID, attrCfg.Name, attr.Value)
			case 2:
				fmt.Fprintf(writer, "[%d]%s: %.4f\n", attr.AttrID, attrCfg.Name, float64(attr.Value)/10000.0)
			case 3:
				fmt.Fprintf(writer, "[%d]%s: %.7f\n", attr.AttrID, attrCfg.Name, float64(attr.Value)/1000000.0)
			default:
				fmt.Fprintf(writer, "[%d]属性值类型[%d]无效\n", attr.AttrID, attrCfg.Figure)
			}
		} else {
			fmt.Fprintf(writer, "[%d]属性配置不存在\n", attr.AttrID)
		}
	}
	if tempBuffInfo != nil {
		fmt.Fprintf(writer, "----------------临时词条列表----------------\n")
		for _, info := range tempBuffInfo.SelectedBuff {
			cfg := GMazeEnergyAffixV8Cfg.GetWithCtx(ctx, info.BuffId)
			if cfg != nil {
				fmt.Fprintf(writer, "[%d]%s 描述: %s\n", info.BuffId, cfg.Affix_name, cfg.Affix_desc)
			} else {
				fmt.Fprintf(writer, "[%d]词条配置不存在\n", info.BuffId)
			}
		}
		fmt.Fprintf(writer, "----------------临时属性列表----------------\n")
		for _, buffInfo := range tempBuffInfo.TotalBuff {
			attrCfg := GMazeAttributeV8Cfg.Get(buffInfo.BuffId)
			if attrCfg != nil {
				switch attrCfg.Figure {
				case 1:
					fmt.Fprintf(writer, "[%d]%s: %d\n", buffInfo.BuffId, attrCfg.Name, buffInfo.BuffValue)
				case 2:
					fmt.Fprintf(writer, "[%d]%s: %.4f\n", buffInfo.BuffId, attrCfg.Name, float64(buffInfo.BuffValue)/10000.0)
				case 3:
					fmt.Fprintf(writer, "[%d]%s: %.7f\n", buffInfo.BuffId, attrCfg.Name, float64(buffInfo.BuffValue)/1000000.0)
				default:
					fmt.Fprintf(writer, "[%d]属性值类型[%d]无效\n", buffInfo.BuffId, attrCfg.Figure)
				}
			} else {
				fmt.Fprintf(writer, "[%d]属性配置不存在\n", buffInfo.BuffId)
			}
		}
	}
}

func (s *service) LookAssembleInfo(writer http.ResponseWriter, request *http.Request) {
	userId := fkutil.ToUint64(request.Form.Get("user_id"))
	ctx := request.Context()
	logger := fklog.ContextAppLogger(ctx)

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
	header, err := equipaassemblegm.PackAssembleHeader(ctx, userId, assembleInfo)
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
		equipaassemblegm.DumpEquipPos(logger, userId, &showBuff, posInfo.GetEquipPos().GetPos(), posInfo, assembleInfo.GetEpSuitId())
	}
	showBuff.WriteString(EndLine)
	showBuff.WriteString(fmt.Sprintf("装备套装:%d\n", assembleInfo.GetEpSuitId()))
	suitBuff, e := equipaassemblegm.PackEquipSuitInfo(logger, userId, assembleInfo, effect)
	if e == nil {
		// 汇总套装属性加成
		showBuff.WriteString(suitBuff)
	}
	showBuff.WriteString(EndLine)

	ar, err := equipaassemblegm.DumpDollCalcAttr(logger, userId)
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

	ar, err = equipaassemblegm.DumpNoForceAttr(logger, userId)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	showBuff.WriteString(ar)
	showBuff.WriteString(EndLine)

	r := showBuff.String()
	writer.Write([]byte(r))
	logger.InfoWF("LookAssembleInfo end")
}
