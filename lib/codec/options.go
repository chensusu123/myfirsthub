package codec

import (
	"maze_game_server/lib/nano/serialize"
)

type CodecOption func(codec interface{})

// WithSerializer
func WithSerializer(serializer serialize.Serializer) CodecOption {
	return func(codec interface{}) {
		switch v := codec.(type) {
		// Es packet codec
		case *EsPacketCodec:
			v.ser = serializer
		// JSON packet codec
		case *JsonPacketCodec:
			v.ser = serializer
		}
	}
}
