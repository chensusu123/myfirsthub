/*
* @Author: majian
* @Date: 2022-04-13 14:02
 */
package fileio

import "testing"

func TestNewFWriter1(t *testing.T) {
	fw := NewFWriter()

	err := fw.Open("2.txt")
	if err != nil {
		return
	}
	defer fw.Close()
	_ = fw.WriteLine("test1\n")
	_ = fw.WriteBytes([]byte("dsfgdsgfdf"))
}

func TestNewFWriter2(t *testing.T) {
	fw, err := NewFWriterWithOpen("3.txt")
	if err != nil {
		return
	}
	defer fw.Close()
	_ = fw.WriteLine("1111\n")
	_ = fw.WriteBytes([]byte("222\n"))
	var lines []uint64
	lines = append(lines, 6253638)
	lines = append(lines, 6253639)
	_ = fw.WriteUint64Slice(lines)
}
