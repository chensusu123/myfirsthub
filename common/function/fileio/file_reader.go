/*
* @Author: majian
* @Date: 2022-04-12 19:32
 */
package fileio

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil/workergroup"
)

type FReader struct {
	logger    fklog.FKLogI
	f         *os.File
	Filename  string
	ParserI   FileParser
	reader    *bufio.Reader
	open      bool
	counter   int32 //统计计数
	dumpRows  int32 //每隔n行统计一次
	SleepLine int32 //执行多少行后sleep n 毫秒
	SleepTime int32 //执行多少行后sleep n 毫秒
	wg        *workergroup.FkWorkGroup
}

type FileParser interface {
	ParseLine(string) []uint64 //解析单行数据
}

type HandlerFunc func(logger fklog.FKLogI, line []uint64) bool

//type AsyncHandlerFunc func(logger fklog.FKLogI, line []uint64,fr *FReader) bool

func NewFReader(logger fklog.FKLogI, parser FileParser) *FReader {
	fr := &FReader{}
	fr.ParserI = parser
	fr.dumpRows = 100
	fr.logger = logger
	return fr
}

func NewDefFReader(logger fklog.FKLogI) *FReader {
	fr := &FReader{}
	fr.ParserI = NewCommonLineParser()
	fr.dumpRows = 100
	fr.logger = logger
	return fr
}

func NewDefFReaderEx(logger fklog.FKLogI, splitChar string) *FReader {
	fr := &FReader{}
	fr.ParserI = NewCommonLineParser(splitChar)
	fr.dumpRows = 100
	fr.logger = logger
	return fr
}

var FileNameNil = errors.New("file name nil")
var FileNoOpen = errors.New("file no open")

func (fr *FReader) Open(fileName string) error {
	if fileName == "" {
		return FileNameNil
	}
	var err error
	fr.f, err = os.Open(fileName)
	if err != nil {
		return err
	}
	fr.Filename = fileName
	fr.open = true
	fr.reader = bufio.NewReader(fr.f)
	return nil
}

func (fr *FReader) Close() {
	if fr.wg != nil {
		fr.wg.Stop()
	}
	if fr.open {
		_ = fr.f.Close()
	}
	return
}

func (fr *FReader) Range(f HandlerFunc) {
	if !fr.open {
		fkfmt.Println("file no open")
		return
	}
	begin := time.Now()
	//	end:=make(chan bool)

	var sw sync.WaitGroup
	for {
		l, _, err := fr.reader.ReadLine()
		if err == io.EOF {
			//	close(end)
			fmt.Println("Range EOF")
			break
		}
		var subItems []uint64
		subItems = fr.ParserI.ParseLine(string(l))
		if len(subItems) == 0 {
			fmt.Println("Range nil row", string(l))
			continue
		}
		//延时
		if fr.SleepLine > 0 && fr.SleepTime > 0 {
			if fr.counter%fr.SleepLine == 0 {
				time.Sleep(time.Duration(fr.SleepTime) * time.Millisecond)
			}
		}

		if fr.wg != nil {
			uid := subItems[0]
			sw.Add(1)
			_ = fr.wg.SendTask(uid, func() {
				defer sw.Done()
				userLogger := fr.logger.Clone("range user")
				userLogger.SetLogId(time.Now().UnixNano())
				userLogger.SetUid(uid)
				_ = f(userLogger, subItems)
				val := atomic.AddInt32(&fr.counter, 1)
				if fr.dumpRows != 0 && val%fr.dumpRows == 0 {
					fkfmt.Println(fmt.Sprintf("run to row[%d]", val))
				}
			})
			continue
		}

		userLogger := fr.logger.Clone("range user")
		userLogger.SetLogId(time.Now().UnixNano())
		userLogger.SetUid(subItems[0])

		ok := f(userLogger, subItems)
		val := atomic.AddInt32(&fr.counter, 1)
		if fr.dumpRows != 0 && val%fr.dumpRows == 0 {
			fkfmt.Println(fmt.Sprintf("run to row[%d]", val))
		}

		if !ok {
			break
		}
	}
	if fr.wg != nil {
		sw.Wait()
	}
	fkfmt.Println("result", "cost", time.Since(begin), "total", fr.counter)
}

func (fr *FReader) Parse(line string) []uint64 {
	if !fr.open {
		fkfmt.Println("file no open")
		return nil
	}
	return fr.ParserI.ParseLine(line)
}

func (fr *FReader) ParseUserId(line string) uint64 {
	if !fr.open {
		fkfmt.Println("file no open")
		return 0
	}
	uids := fr.ParserI.ParseLine(line)
	if len(uids) != 1 {
		return 0
	}
	return uids[0]
}

func (fr *FReader) ReadLine() (s string, err error) {
	if !fr.open {
		fkfmt.Println("file no open")
		return "", FileNoOpen
	}
	l, _, err := fr.reader.ReadLine()
	if err == io.EOF {
		fmt.Println("Range EOF")
		return string(l), err
	}
	return string(l), nil
}

// 设置每隔n行输出一条统计信息
func (fr *FReader) SetDumpRow(dw int32) {
	fr.dumpRows = dw
}

func (fr *FReader) InitAsync(threadCnt int, cacheSize int) {
	fr.wg = workergroup.NewFkWrokGroup(cacheSize, threadCnt, fr.logger)
}

// 每隔waitLine 行 休眠 sleep 毫秒
func (fr *FReader) SetSleep(waitLine int32, sleep int32) {
	fr.SleepLine = waitLine
	fr.SleepTime = sleep
}
