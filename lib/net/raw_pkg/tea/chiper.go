package tea

import (
	"encoding/binary"
	"errors"
	"strconv"
)

// The TEA block size in bytes.
const BlockSize = 8

type KeySizeError int

var DecryptDataSizeError error = errors.New("Decrypt Data size Error")

func (k KeySizeError) Error() string {
	return "cl_base/tea: invalid key size " + strconv.Itoa(int(k))
}

var keys = [4]uint32{2345, 1234, 3456, 6789}

// teaCipher is an instance of TEA encryption.
type teaCipher struct {
	keys [4]uint64
}

func Encrypt(src []byte) (dst []byte, err error) { return encryptBlock(keys, src) }

func Decrypt(src []byte) (dst []byte, err error) { return decryptBlock(keys, src) }

func encryptBlock(keys [4]uint32, src []byte) (dst []byte, err error) {
	dataLen := len(src)
	dataLen2 := dataLen + 2
	dataLen3 := dataLen2

	if dataLen2%8 != 0 {
		dataLen3 = dataLen2 - (dataLen2 % 8) + 8
	}

	tmpData := make([]byte, dataLen3)
	binary.LittleEndian.PutUint16(tmpData[0:2], uint16(dataLen))

	copy(tmpData[2:], src)

	diff := dataLen3 - dataLen2
	for i := 0; i < diff; i++ {
		tmpData[dataLen2+i] = byte(0)
	}

	dst = make([]byte, dataLen3)

	loopCount := dataLen3 >> 3

	for i := 0; i < loopCount; i++ {
		code(dst[i*8:i*8+8], tmpData[i*8:i*8+8], keys)
	}
	return
}

func decryptBlock(keys [4]uint32, src []byte) (dst []byte, err error) {
	dataLen := len(src)
	if dataLen%8 != 0 {
		err = DecryptDataSizeError
		return
	}

	loopCount := dataLen >> 3
	tmp := make([]byte, dataLen)

	for i := 0; i < loopCount; i++ {
		decode(tmp[i*8:i*8+8], src[i*8:i*8+8], keys)
	}
	returnDataLen := binary.LittleEndian.Uint16(tmp[0:2])

	if returnDataLen >= uint16(dataLen) {
		err = DecryptDataSizeError
		return
	}
	dst = tmp[2 : 2+returnDataLen]
	return
}

func code(dst, src []byte, k [4]uint32) {
	n := 32
	delta := uint32(0x9e3779b9)
	s := uint32(0)
	y := binary.LittleEndian.Uint32(src[0:4])
	z := binary.LittleEndian.Uint32(src[4:8])
	for i := 0; i < n; i++ {
		s += delta
		y += ((z << 4) + k[0]) ^ (z + s) ^ ((z >> 5) + k[1])
		z += ((y << 4) + k[2]) ^ (y + s) ^ ((y >> 5) + k[3])
	}
	binary.LittleEndian.PutUint32(dst[0:4], y)
	binary.LittleEndian.PutUint32(dst[4:8], z)
}

func decode(dst, src []byte, k [4]uint32) {
	n := 32
	delta := uint32(0x9e3779b9)
	s := uint32(0xC6EF3720)

	y := binary.LittleEndian.Uint32(src[0:4])
	z := binary.LittleEndian.Uint32(src[4:8])
	for i := 0; i < n; i++ {
		z -= ((y << 4) + k[2]) ^ (y + s) ^ ((y >> 5) + k[3])
		y -= ((z << 4) + k[0]) ^ (z + s) ^ ((z >> 5) + k[1])
		s -= delta
	}
	binary.LittleEndian.PutUint32(dst[0:4], y)
	binary.LittleEndian.PutUint32(dst[4:8], z)
}
