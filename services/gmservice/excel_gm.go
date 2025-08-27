package gmservice

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/iancoleman/orderedmap"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
)

type ShowSheet struct {
	ErrorCode uint64                          `json:"errorCode"`
	ErrorMsg  string                          `json:"errorMsg"`
	Data      []config_manager.ConfigShowItem `json:"data"`
}

type ExcelOutput struct {
	Status int         `json:"status"`
	Desc   string      `json:"desc"`
	Data   DynamicData `json:"data"`
}

// 调整DynamicData以使用有序map
type DynamicData struct {
	List  []*orderedmap.OrderedMap `json:"list"` // 有序map切片，保证对象格式和顺序
	Total int                      `json:"total"`
}

func (s *service) ShowSheet(writer http.ResponseWriter, request *http.Request) {
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
}

func (s *service) GetExcelList(writer http.ResponseWriter, request *http.Request) {
	datas := config_manager.ShowSheet()

	var records []*orderedmap.OrderedMap
	for _, entry := range datas {
		descRecord := orderedmap.New()
		descRecord.Set("excel", entry.XlsxFile)
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

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonOutput)
}

func (s *service) GetExcelSheet(writer http.ResponseWriter, request *http.Request) {
	fileName := request.Form.Get("fileName")

	datas := config_manager.ShowSheet()

	var records []*orderedmap.OrderedMap
	for _, entry := range datas {
		if entry.XlsxFile == fileName {
			descRecord := orderedmap.New()
			descRecord.Set("sheetName", entry.XlsxSheet)
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
}

func (s *service) GetExcelData(writer http.ResponseWriter, request *http.Request) {
	fileName := request.Form.Get("fileName")
	sheetName := request.Form.Get("sheetName")
	var output ExcelOutput

	if data, ok := s.sheetDataCache[sheetName]; ok {
		output = ExcelOutput{
			Status: 0,
			Desc:   "",
			Data:   data,
		}
	} else {
		tableData, err := s.readExcelFile("./conf.d/data/"+fileName, sheetName)
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		output = s.convertTableToJSON(tableData)
	}

	jsonOutput, err := json.Marshal(output)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	writer.Write(jsonOutput)
}

func (s *service) readExcelFile(filePath, sheetName string) ([][]string, error) {
	csvFileName := fmt.Sprintf("./conf.d/data/%s.csv", sheetName)

	file, err := os.Open(csvFileName)
	if err != nil {
		return nil, fmt.Errorf("无法打开CSV文件: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取CSV数据失败: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("CSV文件为空")
	}

	return rows, nil
}

// 转换任意行列数的表格数据为指定JSON格式（保持对象格式和顺序）
func (s *service) convertTableToJSON(table [][]string) ExcelOutput {
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

	var records []*orderedmap.OrderedMap

	// 处理类型行（第二行）
	typeRecord := orderedmap.New()
	for i, key := range keys {
		val := ""
		if i < len(table[1]) {
			val = table[1][i]
		}
		typeRecord.Set(key, val) // 按顺序插入键值对
	}
	records = append(records, typeRecord)

	// 处理描述行（第四行）
	descRecord := orderedmap.New()
	for i, key := range keys {
		val := ""
		if i < len(table[3]) {
			val = table[3][i]
		}
		descRecord.Set(key, val) // 按顺序插入键值对
	}
	records = append(records, descRecord)

	// 处理数据行（从第五行开始）
	for i := 4; i < len(table); i++ {
		dataRow := table[i]
		dataRecord := orderedmap.New()

		for j, key := range keys {
			val := ""
			if j < len(dataRow) {
				val = dataRow[j]
			}
			dataRecord.Set(key, val) // 按顺序插入键值对
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
