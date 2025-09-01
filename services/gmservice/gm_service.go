package gmservice

import (
	"context"
	"maze_game_server/model/gmmodel"
	"net/http"
)

type gmService interface {
	// 注册gm接口
	SafeGETRegister(ctx context.Context, pattern string, handler func(http.ResponseWriter, *http.Request))
	SafePOSTRegister(ctx context.Context, pattern string, handler func(http.ResponseWriter, *http.Request))

	// buffGM
	//	-- 设置临时buff
	SetMazeTempBuff(writer http.ResponseWriter, request *http.Request)

	// itemGM
	AddRefreshCost(writer http.ResponseWriter, request *http.Request)
	AddExp(writer http.ResponseWriter, request *http.Request)
	AddEnergy(writer http.ResponseWriter, request *http.Request)
	AddItem(writer http.ResponseWriter, request *http.Request)

	// user
	LookAssembleInfo(writer http.ResponseWriter, request *http.Request)
	SetBarrier(writer http.ResponseWriter, request *http.Request)
	DumpBattleData(writer http.ResponseWriter, request *http.Request)
	Attrs(writer http.ResponseWriter, request *http.Request)
	SetUserLevel(writer http.ResponseWriter, request *http.Request)

	// equipGM
	GetEquipInfoByCfgId(writer http.ResponseWriter, request *http.Request)
	GetEquipInfoByGuid(writer http.ResponseWriter, request *http.Request)
	SendOneSuitEquip(writer http.ResponseWriter, request *http.Request)
	ReInitDollEquip(writer http.ResponseWriter, request *http.Request)
	GmDressBagEquip(writer http.ResponseWriter, request *http.Request)
	GmEquipPosLvUp(writer http.ResponseWriter, request *http.Request)
	AddEquip(writer http.ResponseWriter, request *http.Request)
	SetEquipRollScore(writer http.ResponseWriter, request *http.Request)
	BatchAddEquip(writer http.ResponseWriter, request *http.Request)

	// otherGM
	GenerateUser(writer http.ResponseWriter, request *http.Request)
	Online(writer http.ResponseWriter, request *http.Request)
	ClearBag(writer http.ResponseWriter, request *http.Request)
	ClearBagNotAssemble(writer http.ResponseWriter, request *http.Request)

	// excelGM
	ShowSheet(writer http.ResponseWriter, request *http.Request)
	GetExcelList(writer http.ResponseWriter, request *http.Request)
	GetExcelSheet(writer http.ResponseWriter, request *http.Request)
	GetExcelData(writer http.ResponseWriter, request *http.Request)
	readExcelFile(filePath, sheetName string) ([][]string, error)
	convertTableToJSON(table [][]string) gmmodel.Output

	// 统一注册http接口
	RegHttp(ctx context.Context)
}

var GmService gmService

type service struct {
	// 表数据缓存 todo 该数据和使用数据为同一份
	sheetDataCache map[string]gmmodel.DynamicData
}

func NewGmService() gmService {
	svr := &service{}
	svr.sheetDataCache = make(map[string]gmmodel.DynamicData)
	return svr
}

func init() {
	GmService = NewGmService()
}
