/*
* @Author: majian
* @Date: 2022-04-12 21:21
 */
package fileio

import (
	"gitlab.ifreetalk.com/plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"io"
	"testing"
	"time"
)

func dumpLine(logger fklog.FKLogI, line []uint64) bool {
	for _, item := range line {
		//fkfmt.Println("dump", item)
		logger.InfoWF("dump", zap.Uint64("item", item))
	}
	return true
}

//测试用例1 使用range 遍历
func TestNewFReader1(t *testing.T) {
	logger := fklog.InitStdoutLog("")
	fr := NewDefFReader(logger)
	err := fr.Open("1.txt")
	if err != nil {
		fkfmt.Println("err", err)
		return
	}
	defer fr.Close()
	fr.SetDumpRow(10)
	fr.Range(dumpLine)
}

//测试用例2 自己遍历
func TestNewFReader2(t *testing.T) {
	logger := fklog.InitStdoutLog("")
	fr := NewDefFReader(logger)
	err := fr.Open("1.txt")
	if err != nil {
		fkfmt.Println("err", err)
		return
	}
	defer fr.Close()
	for {
		l, err := fr.ReadLine()
		if err == io.EOF {
			break
		}
		var subItems []uint64
		subItems = fr.Parse(l)
		if len(subItems) == 0 {
			continue
		}
		userLogger := fr.logger.Clone("range user")
		userLogger.SetLogId(time.Now().UnixNano())
		userLogger.SetUid(subItems[0])
		ok := dumpLine(userLogger, subItems)
		if !ok {
			break
		}
	}
}
