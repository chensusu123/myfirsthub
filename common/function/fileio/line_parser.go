/*
* @Author: majian
* @Date: 2022-04-12 20:56
 */
package fileio

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/common/fkfmt"
	"strconv"
	"strings"
)

//解析用户Id
type DerLineParser struct {
}

//解析用户Id
func (m *DerLineParser) ParseLine(line string) []uint64 {
	ss := strings.Trim(line, " ")
	userId, err := strconv.ParseUint(ss, 10, 64)
	if err != nil {
		return []uint64{0}
	}
	fkfmt.Println("ss", ss, "uid", userId)
	return []uint64{userId}
}

type CommonLineParser struct {
	splitChar string //分割字符
}

func NewCommonLineParser(splitChar ...string) *CommonLineParser {
	if len(splitChar) == 1 {
		return &CommonLineParser{splitChar: splitChar[0]}
	} else {
		return &CommonLineParser{}
	}
}

//解析单行
func (m *CommonLineParser) ParseLine(line string) []uint64 {
	splitChar := m.splitChar
	if splitChar == "" {
		splitChar = " "
	}
	var ret []uint64
	//清除前后空格
	ss := strings.Trim(line, " ")
	cols := strings.Split(ss, splitChar)
	for _, col := range cols {
		colVal, err := strconv.ParseUint(col, 10, 64)
		if err != nil {
			fkfmt.Println("ParseUint error", "col", col, "err", err)
			return []uint64{0}
		}
		ret = append(ret, colVal)
	}
	return ret
}
