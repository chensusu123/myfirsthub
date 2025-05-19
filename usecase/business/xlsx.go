package business

import (
	"errors"
	"strings"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"

	"go.uber.org/zap"
)

// 读取golang组配置文件
func loadGoTeamXlsx(logger fklog.FKLogI, file, sheetName string, rows [][]string) (cache *excelFullCache, err error) {
	// 4行
	if len(rows) < 4 {
		return
	}
	// 第一行空. 是标记
	// 第二行是类型
	rowType := rows[1]
	// 第三行是字段名
	rowField := rows[2]
	// 第四行是注释
	rowComment := rows[3]
	ignore := false
	fieldName := ""
	checkField := make(map[string]bool)

	var fieldConfig []int

	cache = &excelFullCache{}
	cache.logger = logger
	cache.FieldMap = make(map[string]int)

	for k := 0; k < len(rowType); k++ {
		if k >= len(rowType) {
			logger.ErrorWF("files sheet head data has empty.", zap.String("file", file),
				zap.String("sheet", sheetName), zap.String("row", "1"), zap.Int("column", k))
			continue
		}
		if k >= len(rowField) {
			logger.ErrorWF("files sheet head data has empty.", zap.String("file", file),
				zap.String("sheet", sheetName), zap.String("row", "2"), zap.Int("column", k))
			continue
		}
		if k >= len(rowComment) {
			logger.ErrorWF("files sheet head data has empty.", zap.String("file", file),
				zap.String("sheet", sheetName), zap.String("row", "3"), zap.Int("column", k))
			continue
		}
		// // 检测有效性
		// if rowType[k] == nil ||
		// 	rowField[k] == nil ||
		// 	rowComment[k] == nil {
		// 	logger.WarnWF(" files sheet head data has empty.", zap.String("sheet", sheetName), zap.Int("column", k))
		// 	break
		// }

		ignore = checkIgnoreFiled(rowType[k], true)
		if ignore {
			logger.WarnWF(" ignore column data(type ignore).", zap.String("file", file),
				zap.String("sheet", sheetName),
				zap.Int("column", k), zap.String("type", rowType[k]))
			continue
		}
		fieldName = rowField[k]
		if fieldName == "" {
			logger.ErrorWF(" get column field_name (empty).", zap.String("file", file),
				zap.String("sheet", sheetName),
				zap.Int("column", k), zap.Strings("field_row", rowField), zap.Error(err))
			continue
		}
		if _, ok := checkField[fieldName]; ok {
			logger.ErrorWF(" invliad column name(repeated).", zap.String("file", file),
				zap.String("sheet", sheetName),
				zap.Int("column", k), zap.String("field_name", fieldName))
			err = errors.New("invliad column name(repeated)")
			continue
		}
		checkField[fieldName] = true
		// 保存缓存数据
		cache.Fields = append(cache.Fields, fieldName)
		cache.FieldMap[fieldName] = len(cache.Fields) - 1
		fieldConfig = append(fieldConfig, k)
	}
	// 写入数据
	checkEmptyRow := true
	for rowIndex := 3 + 1; rowIndex < len(rows); rowIndex++ {
		var rowData []string
		cells := rows[rowIndex]
		for _, column := range fieldConfig {
			if rowIndex >= len(rows) || column >= len(cells) {
				rowData = append(rowData, "")
			} else {
				rowData = append(rowData, cells[column])
			}
			// logger.DebugWF("debug data.", zap.String("file", file),
			// 	zap.String("sheet", sheetName), zap.Int("index", k), zap.String("column", cache.Fields[k]),
			// 	zap.String("data", rowData[k]))
		}
		checkEmptyRow = true
		for _, column := range rowData {
			if column != "" {
				checkEmptyRow = false
				break
			}
		}
		if checkEmptyRow {
			logger.WarnWF("ignore empty row.", zap.String("file", file),
				zap.String("sheet", sheetName),
				zap.Int("row", rowIndex+1))
			continue
		}
		cache.Data = append(cache.Data, rowData)
	}
	return
}

// 解析配置文件sheet
func parseExcellSheet(logger fklog.FKLogI, file, sheetName string, rows [][]string) (cache *excelFullCache, err error) {
	if len(rows) < 3 || len(rows[0]) < 1 {
		logger.ErrorWF("read invadlid xlsx sheet.", zap.String("file", file),
			zap.String("sheet", sheetName))
		return
	}
	defer func() {
		logger.WarnWF("read show sheet.", zap.String("file", file),
			zap.String("sheet", sheetName), zap.Any("fields", cache.Fields), zap.Any("data", len(cache.Data)))
	}()
	firstRow := rows[0]
	// 检测是否golang组的配置
	if firstRow[0] == "paipaiworld_config" {
		return loadGoTeamXlsx(logger, file, sheetName, rows)
	}

	// 第一行是注释,第二行是类型.第三行是字段名
	rowComment := firstRow
	rowType := rows[1]
	rowField := rows[2]

	fieldName := ""
	checkField := make(map[string]bool)

	var fieldConfig []int

	cache = &excelFullCache{}
	cache.logger = logger
	cache.FieldMap = make(map[string]int)

	// 遍历所有列. 构建表头
	for k := 0; k < len(rows[0]); k++ {

		if k >= len(rowType) {
			logger.ErrorWF("files sheet head data has empty.", zap.String("file", file),
				zap.String("sheet", sheetName),
				zap.String("row", "1"), zap.Int("column", k))
			continue
		}
		if k >= len(rowField) {
			logger.ErrorWF("files sheet head data has empty.", zap.String("file", file),
				zap.String("sheet", sheetName),
				zap.String("row", "2"), zap.Int("column", k))
			continue
		}
		if k >= len(rowComment) {
			logger.ErrorWF("files sheet head data has empty.", zap.String("file", file),
				zap.String("sheet", sheetName),
				zap.String("row", "0"), zap.Int("column", k))
			continue
		}
		// // 检测有效性
		// if rowComment[k] == nil ||
		// 	rowType[k] == nil ||
		// 	rowField[k] == nil {
		// 	logger.WarnWF(" files sheet head data has nil.", zap.String("sheet", sheetName), zap.Int("column", k))
		// 	continue
		// }
		if rowComment[k] == "" &&
			rowType[k] == "" &&
			rowField[k] == "" {
			logger.WarnWF(" files sheet head data has empty.", zap.String("file", file),
				zap.String("sheet", sheetName), zap.Int("column", k))
			continue
		}
		fieldName = rowField[k]
		if fieldName == "" {
			logger.ErrorWF("files sheet head empty column in field row.", zap.String("file", file),
				zap.String("sheet", sheetName), zap.Int("column", k))
			continue
		}
		//
		// sfs = append(sfs, field)
		if _, ok := checkField[fieldName]; ok {
			logger.ErrorWF(" invliad column name(repeated).", zap.String("file", file),
				zap.String("sheet", sheetName),
				zap.Int("column", k), zap.String("field_name", fieldName))
			err = errors.New("invliad column name(repeated)")
			continue
		}

		checkField[fieldName] = true
		// 保存缓存数据
		cache.Fields = append(cache.Fields, fieldName)
		cache.FieldMap[fieldName] = len(cache.Fields) - 1
		fieldConfig = append(fieldConfig, k)
	}
	// 写入数据
	checkEmptyRow := true
	for rowIndex := 3; rowIndex < len(rows); rowIndex++ {
		var rowData []string
		cells := rows[rowIndex]
		for _, column := range fieldConfig {
			if rowIndex >= len(rows) || column >= len(cells) {
				rowData = append(rowData, "")
			} else {
				rowData = append(rowData, cells[column])
			}
		}
		checkEmptyRow = true
		for _, column := range rowData {
			if column != "" {
				checkEmptyRow = false
				break
			}
		}
		if checkEmptyRow {
			logger.WarnWF("ignore empty row.", zap.String("file", file),
				zap.String("sheet", sheetName),
				zap.Int("row", rowIndex+1))
			continue
		}
		cache.Data = append(cache.Data, rowData)
	}

	return
}

func checkIgnoreFiled(tsrc string, goTeam bool) bool {
	// Go项目组 对于类型组织为 TYPE_[TAG]
	// TYPE 类型
	//   基本类型: INT,LONG,FLOAT,STRING
	//	 符合类型: LIST,MAP
	// TAG 标记
	//   C	客户端
	//   S  服务端
	//   N  策划用说明
	//   K  主键程序查询用

	// 这只只处理 S,K 标记. 其他都忽略
	if goTeam {
		// 先转换小写.防止书写错误
		ts := strings.Split(strings.ToLower(tsrc), "_")
		for _, v := range ts[1:] {
			if v == "s" || v == "k" {
				return false
			}
		}
		return true
	}
	return false
}

func checkSheetName(name string) (string, bool) {
	for _, ch := range name {
		if ch >= 256 {
			return name, false
		}
	}
	return name, true
}
