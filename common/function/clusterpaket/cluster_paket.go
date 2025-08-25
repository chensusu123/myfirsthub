package clusterpaket

import (
	"encoding/binary"
	"errors"

	"google.golang.org/protobuf/proto"
)

// ErrWrongValueType is the error used for marshal the value with protobuf encoding.
var ErrWrongValueType = errors.New("protobuf: convert on wrong type value")

func MakeClusterPacket(packetType uint16, packet interface{}) ([]byte, error) {
	pb, ok := packet.(proto.Message)
	if !ok {
		return nil, ErrWrongValueType
	}
	data, err := proto.Marshal(pb)
	if err != nil {
		return nil, err
	}
	ret := make([]byte, 2+len(data))
	binary.LittleEndian.PutUint16(ret[:2], packetType)
	copy(ret[2:], data)
	return ret, nil
}

func SplitClusterPacket(data []byte) (uint16, []byte, error) {
	if len(data) < 2 {
		return 0, nil, errors.New("data len is too short")
	}

	packetType := binary.LittleEndian.Uint16(data[:2])
	data = data[2:]

	return packetType, data, nil
}
