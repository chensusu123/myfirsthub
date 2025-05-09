package raw_pkg

import (
	"encoding/binary"

	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/raw_pkg/tea"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet/fkpkg"
)

const TeaEncryption = 1

type StruSvrEsRawBaseHead struct {
	PackLen      uint16
	PackType     uint16
	SessionID    uint32
	flag         uint16
	CompressType uint8
	EsRsTime     uint64
	EsRqTime     uint64
	Data         []byte
}

const DEF_TCP_PACKHEAD_LEN_ES = 27
const (
	TfAck = 0x80
	TpAck = 0x40
)

func (pkg *StruSvrEsRawBaseHead) SetTeaflag() {
	pkg.flag = pkg.flag | TeaEncryption
}

func (pkg *StruSvrEsRawBaseHead) IsTea() bool {
	return pkg.flag&0x1 == TeaEncryption
}

// Pack pack binary data
func (pkg *StruSvrEsRawBaseHead) Pack() ([]byte, error) {
	var ret []byte

	if !pkg.IsTea() {
		pkg.PackLen = uint16(DEF_TCP_PACKHEAD_LEN_ES + len(pkg.Data))
		ret = make([]byte, pkg.PackLen)
		binary.LittleEndian.PutUint16(ret[:2], uint16(len(ret)))
		binary.LittleEndian.PutUint16(ret[2:4], pkg.flag)
		ret[4] = uint8(pkg.CompressType)
		binary.LittleEndian.PutUint32(ret[5:9], pkg.SessionID)
		binary.LittleEndian.PutUint64(ret[9:17], uint64(pkg.EsRsTime))
		binary.LittleEndian.PutUint64(ret[17:25], uint64(pkg.EsRqTime))
		binary.LittleEndian.PutUint16(ret[25:27], uint16(pkg.PackType))
		copy(ret[DEF_TCP_PACKHEAD_LEN_ES:], pkg.Data)
	} else {
		ret1 := make([]byte, 2+len(pkg.Data))
		binary.LittleEndian.PutUint16(ret1[0:2], pkg.PackType)
		copy(ret1[2:], pkg.Data)
		// tea 加密
		encData, err := tea.Encrypt(ret1)
		if err != nil {
			return pkg.Data, err
		}

		packLen := DEF_TCP_PACKHEAD_LEN_ES - 2 + uint16(len(encData))
		pkg.PackLen = packLen
		returnData := make([]byte, packLen)

		binary.LittleEndian.PutUint16(returnData[:2], uint16(len(ret)))
		binary.LittleEndian.PutUint16(returnData[2:4], pkg.flag)
		returnData[4] = uint8(pkg.CompressType)
		binary.LittleEndian.PutUint32(returnData[5:9], pkg.SessionID)
		binary.LittleEndian.PutUint64(returnData[9:17], uint64(pkg.EsRsTime))
		binary.LittleEndian.PutUint64(returnData[17:25], uint64(pkg.EsRqTime))
		copy(returnData[DEF_TCP_PACKHEAD_LEN_ES-2:], encData)
		return returnData, nil
	}

	return ret, nil
}

// UnPack pack binary data
func (pkg *StruSvrEsRawBaseHead) UnPack(data []byte) error {
	size := int(binary.LittleEndian.Uint16(data[:2]))
	if len(data) < size {
		return fkpkg.InvalidTcpPackge
	}
	pkg.flag = binary.LittleEndian.Uint16(data[2:4])
	pkg.PackLen = uint16(size)
	pkg.CompressType = uint8(data[4])
	pkg.SessionID = binary.LittleEndian.Uint32(data[5:9])
	pkg.EsRsTime = binary.LittleEndian.Uint64(data[9:17])
	pkg.EsRqTime = binary.LittleEndian.Uint64(data[17:25])
	if !pkg.IsTea() {
		pkg.PackType = binary.LittleEndian.Uint16(data[25:27])
		pkg.Data = make([]byte, size-DEF_TCP_PACKHEAD_LEN_ES)
		copy(pkg.Data, data[DEF_TCP_PACKHEAD_LEN_ES:size])
	} else {
		decData, err := tea.Decrypt(data[DEF_TCP_PACKHEAD_LEN_ES-2:])
		if err != nil {
			return err
		}
		pkg.PackType = binary.LittleEndian.Uint16(decData[0:2])
		dataSize := len(decData)
		pkg.Data = make([]byte, dataSize-2)
		copy(pkg.Data, decData[2:dataSize])
	}
	return nil
}
