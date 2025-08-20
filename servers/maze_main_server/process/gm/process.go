package gm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"maze_game_server/excel/mazeenergyaffixlvv8config"
	"maze_game_server/model/tempbuffmodel"
	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/services/barrierenergyservice"
	"maze_game_server/services/tempbuffservice"

	"maze_game_server/common/function/gentradeno"
	"maze_game_server/common/function/gm"
	"maze_game_server/config/GMazeAttributeV8Cfg"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeItemsV8Cfg"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/redis/UnionIDBindRedis"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/io/redis/mazefixedbarrierredis"
	"maze_game_server/io/redis/useridredis"
	"maze_game_server/io/redis/usersection"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazecommonvalue"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/servers/maze_main_server/process/gm/cmdbattledata"
	"maze_game_server/usecase/online"

	"github.com/gorilla/schema"
	"github.com/xuri/excelize/v2"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/appconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var form = schema.NewDecoder()

func RegGm(logger fklog.FKLogI) {
	cmdbattledata.RegBattleDataGm(logger)

	gm.SafeHttpRegister(logger, "/AddExp", func(writer http.ResponseWriter, request *http.Request) {
		userId := fkutil.ToUint64(request.Form.Get("userId"))
		exp := fkutil.ToInt64(request.Form.Get("exp"))
		logger.SetLogId(time.Now().UnixNano())

		userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
		if err != nil {
			logger.ErrorWF("AddExp GetUserInfoV2 fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}

		oldLevel := userInfo.Level
		oldExp := userInfo.TotalExp

		// 更新等级经验
		err = userInfo.AddExp(exp)
		if err != nil {
			logger.ErrorWF("AddExp CalExp fail", zap.Error(err))
			writer.Write([]byte(err.Error()))
			return
		}
		err = mazeuserinfo.SetUserInfoV2(logger, userId, userInfo)
		if err != nil {
			logger.ErrorWF("AddExp SetUserInfoV2 fail", zap.Error(err))
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
				mazeuserlevelkafka.PushMazeLevelRecord(logger, levelRecord)
			}
		}()

		writer.Write([]byte("设置成功，注意尽量不要在迷宫杀怪时使用本gm"))
	})

	type SetBarrierParams struct {
		UserID    uint64 `schema:"userId,required"`
		BarrierID int32  `schema:"barrierId,required"`
		Lock      int32  `schema:"lock"` // 如果提供这个参数，则设置的关卡会记录下来，重置游戏数据也会继续生效
	}

	gm.SafeHttpRegister(logger, "/SetBarrier", func(writer http.ResponseWriter, request *http.Request) {
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

		userInfo, err := mazeuserinfo.GetUserInfoV2(logger, params.UserID)
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
		err = mazeuserinfo.SetUserInfoV2(logger, params.UserID, userInfo)
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
	})

	gm.SafeHttpRegister(logger, "/generateUser", func(writer http.ResponseWriter, request *http.Request) {
		logger.SetLogId(time.Now().UnixNano())
		AuthId := fkutil.ToUint64(request.Form.Get("AuthId"))
		ctx := request.Context()
		userID := uint64(0)

		generateUser := &GenerateUser{}
		defer func() {
			generateUser.UserId = userID
			jsonData, err := json.Marshal(generateUser)
			if err != nil {
				writer.Write([]byte(err.Error()))
				return
			}
			writer.Write(jsonData)
		}()

		if AuthId == 0 {
			generateUser.ErrorCode = 1
			generateUser.ErrorMsg = "AuthId is 0"
			return
		}

		users, err := UnionIDBindRedis.GetUsersWithUnionID(ctx, logger, uint64(AuthId))
		if err != nil {
			generateUser.ErrorCode = 1
			generateUser.ErrorMsg = err.Error()
			return
		}

		if len(users) == 0 {
			newUserID := useridredis.Generate(ctx, logger)
			if newUserID == 0 {
				generateUser.ErrorCode = 1
				generateUser.ErrorMsg = "Generate error"
				return
			}
			err = UnionIDBindRedis.AddUnionID2UserID(ctx, logger, uint64(AuthId), newUserID)
			if err != nil {
				generateUser.ErrorCode = 1
				generateUser.ErrorMsg = err.Error()
				return
			}
			err = UnionIDBindRedis.AddUserID2UnionID(ctx, logger, newUserID, uint64(AuthId))
			if err != nil {
				generateUser.ErrorCode = 1
				generateUser.ErrorMsg = err.Error()
				return
			}
			err = usersection.Set(ctx, newUserID, appconfig.GlobalConfig().Global.SectionID)
			if err != nil {
				logger.ErrorWF("usersection.Set fail",
					zap.Uint64("userID", userID),
					zap.Error(err))
			}
			userID = newUserID
		} else {
			userID = users[0]
		}
		generateUser.ErrorCode = 0
		generateUser.ErrorMsg = "success"
	})

	gm.SafeHttpRegister(logger, "/showSheet", func(writer http.ResponseWriter, request *http.Request) {
		logger.SetLogId(time.Now().UnixNano())
		showSheet := &ShowSheet{}
		defer func() {
			jsonData, err := json.Marshal(showSheet)
			if err != nil {
				writer.Write([]byte(err.Error()))
				return
			}
			writer.Write(jsonData)
		}()
		showSheet.Data = config_manager.ShowSheet()
		showSheet.ErrorCode = 0
		showSheet.ErrorMsg = "success"
	})

	gm.SafeHttpRegister(logger, "/online", func(writer http.ResponseWriter, request *http.Request) {
		logger.SetLogId(time.Now().UnixNano())
		fmt.Fprintf(writer, "会话ID 用户ID 客户端地址\n")
		online.Scan(func(id int64, s *session.Session) {
			s.RLock()
			defer s.RUnlock()
			if userID := s.UID(); userID <= 0 {
				fmt.Fprintf(writer, "%d 验证中 %s\n", s.ID(), s.RemoteAddr().String())
			} else {
				fmt.Fprintf(writer, "%d %d %s\n", s.ID(), userID, s.RemoteAddr().String())
			}
		})
	})

	gm.SafeHttpRegister(logger, "/attrs", func(writer http.ResponseWriter, request *http.Request) {
		logger.SetLogId(time.Now().UnixNano())
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
			tempBuffInfo, err = tempbuffservice.GlobalTempBuffService.GetTempBuffInfo(logger, userId, barrierId)
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
				cfg := mazeenergyaffixlvv8config.GetAffixConfig(info.BuffId)
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
	})

	gm.SafeHttpRegister(logger, "/addEnergy", func(writer http.ResponseWriter, request *http.Request) {
		logger.SetLogId(time.Now().UnixNano())
		var (
			userId = fkutil.ToUint64(request.Form.Get("userId"))
			count  = fkutil.ToInt32(request.Form.Get("count"))
		)
		if count > barrierenergyservice.GlobalBarrierEnergyService.GetEnergyMaxValue() {
			return
		}
		curEnergy, _, err := barrierenergyservice.GlobalBarrierEnergyService.AddEnergy(logger, userId, count)
		if err != nil {
			return
		}
		fmt.Fprintf(writer, "add barrier energy success, curEnergy=[%d]", curEnergy)
	})

	// 添加道具
	gm.SafeHttpRegister(logger, "/addItem", func(writer http.ResponseWriter, request *http.Request) {
		logger.SetLogId(time.Now().UnixNano())
		var (
			userId = fkutil.ToUint64(request.Form.Get("userId"))
			itemId = fkutil.ToInt32(request.Form.Get("itemId"))
			count  = fkutil.ToInt64(request.Form.Get("count"))
		)

		tradeNo := gentradeno.GetTradeNum()
		items := make([]*MazeCommon.MazeItem, 0)

		if userId <= 0 {
			fmt.Fprintf(writer, "请指定有效用户ID")
			return
		}

		if itemId <= 0 || count <= 0 {
			fmt.Fprintf(writer, "无效道具ID或道具数量")
			return
		}

		itemCfg := GMazeItemsV8Cfg.Get(itemId)
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

		errInfo := gentradeno.AddItemEx(logger, userId, 697, tradeNo, header, items...)
		if errInfo != nil {
			fmt.Fprintf(writer, "添加道具失败，错误：%s", string(errInfo.GetErrMsg()))
			logger.ErrorWF("addItem AddItemEx fail", zap.Any("errInfo", errInfo), zap.Any("ItemList", items))
			return
		}

		fmt.Fprintf(writer, "添加[%d]个道具[%s]成功", count, itemCfg.Prop_name)
	})

	gm.SafeHttpRegister(logger, "/GetExcelList", func(writer http.ResponseWriter, request *http.Request) {
		entires, err := os.ReadDir("./conf.d/data")
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}

		var records []map[string]interface{}
		for _, entry := range entires {
			if strings.Contains(entry.Name(), "black_excel") || strings.Contains(entry.Name(), "column_relation") || strings.Contains(entry.Name(), "git_version_v8") {
				continue
			}
			descRecord := make(map[string]interface{})
			if strings.ToLower(filepath.Ext(entry.Name())) == ".xlsx" {
				descRecord["excel"] = entry.Name()
				records = append(records, descRecord)
			}
		}

		output := ExcelOutput{
			Status: 0,
			Desc:   "",
			Data: DynamicData{
				List:  records,
				Total: len(records),
			},
		}
		jsonOutput, err := json.Marshal(output)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}

		writer.Write(jsonOutput)
	})

	gm.SafeHttpRegister(logger, "/GetExcelSheet", func(writer http.ResponseWriter, request *http.Request) {
		fileName := request.Form.Get("fileName")
		sheets, err := GetSheets("./conf.d/data/" + fileName)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}

		var records []map[string]interface{}
		for _, sheet := range sheets {
			descRecord := make(map[string]interface{})
			descRecord["sheetName"] = sheet
			records = append(records, descRecord)
		}

		output := ExcelOutput{
			Status: 0,
			Desc:   "",
			Data: DynamicData{
				List:  records,
				Total: len(records),
			},
		}
		jsonOutput, err := json.Marshal(output)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}

		writer.Write(jsonOutput)
	})

	gm.SafeHttpRegister(logger, "/GetExcelData", func(writer http.ResponseWriter, request *http.Request) {
		fileName := request.Form.Get("fileName")
		sheetName := request.Form.Get("sheetName")
		tableData, err := readExcelFile("./conf.d/data/"+fileName, sheetName)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}

		output := convertTableToJSON(tableData)
		jsonOutput, err := json.Marshal(output)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}

		writer.Write(jsonOutput)
	})
}

type ShowSheet struct {
	ErrorCode uint64                          `json:"errorCode"`
	ErrorMsg  string                          `json:"errorMsg"`
	Data      []config_manager.ConfigShowItem `json:"data"`
}
type GenerateUser struct {
	ErrorCode uint64 `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
	UserId    uint64 `json:"userId"`
}

func GetSheets(filePath string) ([]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开Excel文件: %w", err)
	}
	defer f.Close()
	sheets := make([]string, 0)
	for _, v := range f.GetSheetList() {
		sheets = append(sheets, v)
	}
	return sheets, nil
}

// 从Excel文件读取指定工作表数据（默认读取第一个工作表）
func readExcelFile(filePath, sheetName string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开Excel文件: %w", err)
	}
	defer f.Close()

	if sheetName == "" {
		sheets := f.GetSheetList()
		if len(sheets) == 0 {
			return nil, fmt.Errorf("Excel文件中没有工作表")
		}
		sheetName = sheets[0]
		fmt.Printf("使用工作表: %s\n", sheetName)
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表数据失败: %w", err)
	}

	return rows, nil
}

// 转换任意行列数的表格数据为指定JSON格式
func convertTableToJSON(table [][]string) ExcelOutput {
	if len(table) < 4 {
		return ExcelOutput{
			Status: 1,
			Desc:   "表格数据行数不足，至少需要4行",
			Data:   DynamicData{},
		}
	}

	keys := table[2]
	if len(keys) == 0 {
		return ExcelOutput{
			Status: 2,
			Desc:   "未找到有效键名（第三行）",
			Data:   DynamicData{},
		}
	}

	var records []map[string]interface{}

	typeRow := table[1]
	typeRecord := make(map[string]interface{})
	for i, key := range keys {
		if i < len(typeRow) {
			typeRecord[key] = typeRow[i]
		} else {
			typeRecord[key] = ""
		}
	}
	records = append(records, typeRecord)

	descRow := table[3]
	descRecord := make(map[string]interface{})
	for i, key := range keys {
		if i < len(descRow) {
			descRecord[key] = descRow[i]
		} else {
			descRecord[key] = ""
		}
	}
	records = append(records, descRecord)

	for i := 4; i < len(table); i++ {
		dataRow := table[i]
		dataRecord := make(map[string]interface{})

		for j, key := range keys {
			if j < len(dataRow) {
				dataRecord[key] = dataRow[j]
			} else {
				dataRecord[key] = ""
			}
		}

		records = append(records, dataRecord)
	}

	return ExcelOutput{
		Status: 0,
		Desc:   "",
		Data: DynamicData{
			List:  records,
			Total: len(records),
		},
	}
}

type ExcelOutput struct {
	Status int         `json:"status"`
	Desc   string      `json:"desc"`
	Data   DynamicData `json:"data"`
}

type DynamicData struct {
	List  []map[string]interface{} `json:"list"`
	Total int                      `json:"total"`
}
