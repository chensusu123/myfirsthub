package codec

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

// ErrWrongValueType is the error used for marshal the value with protobuf encoding.
var ErrWrongValueType = errors.New("protobuf: convert on wrong type value")

type ProtobufSerializer struct {
}

// NewProtobufSerializer
func NewProtobufSerializer() *ProtobufSerializer {
	return &ProtobufSerializer{}
}

// Marshal implements serialize.Serializer.
func (p *ProtobufSerializer) Marshal(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, ErrWrongValueType
	}
	return proto.Marshal(pb)
}

// Unmarshal implements serialize.Serializer.
func (p *ProtobufSerializer) Unmarshal(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return ErrWrongValueType
	}
	return proto.Unmarshal(data, pb)
}
