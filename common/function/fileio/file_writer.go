/*
* @Author: majian
* @Date: 2022-04-12 21:03
 */
package fileio

import (
	"fmt"
	"gitlab.ifreetalk.com/plate/freetk/common/fkfmt"
	"os"
)

type FWriter struct {
	f        *os.File
	FileName string
	open     bool
}

func NewFWriter() *FWriter {
	out := new(FWriter)
	return out
}

func NewFWriterWithOpen(fileName string) (*FWriter, error) {
	out := new(FWriter)
	out.FileName = fileName
	err := out.Open(fileName)
	return out, err
}

func (m *FWriter) Open(fileName string) error {
	m.FileName = fileName
	var err error
	m.f, err = os.Create(m.FileName)
	if err != nil {
		return err
	}
	m.open = true
	return err
}

func (m *FWriter) WriteLine(line string) error {
	if !m.open {
		fkfmt.Println("file no open")
	}
	_, err := m.f.WriteString(line)
	return err
}

func (m *FWriter) WriteBytes(line []byte) error {
	if !m.open {
		fkfmt.Println("file no open")
	}
	_, err := m.f.Write(line)
	return err
}

func (m *FWriter) Close() {
	if m.open {
		m.f.Close()
	}
	return
}

func (m *FWriter) WriteUint64Slice(lines []uint64) error {
	if !m.open {
		fkfmt.Println("file no open")
	}
	for _, line := range lines {
		err := m.WriteLine(fmt.Sprintf("%d\n", line))
		if err != nil {
			return err
		}
	}
	return nil
}

//输出一行,按sep 指定分割
func (m *FWriter) WriteUint64Row(lines []uint64, sep string) error {
	if !m.open {
		fkfmt.Println("file no open")
	}
	if sep == "" {
		sep = " "
	}
	var row string
	for i, line := range lines {
		if i == len(lines)-1 {
			row += fmt.Sprintf("%d\n", line)
		} else {
			row += fmt.Sprintf("%d%s", line, sep)
		}
	}
	if row != "" {
		return m.WriteLine(row)
	}
	return nil
}
