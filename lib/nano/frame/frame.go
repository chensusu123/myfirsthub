package frame

import (
	"maze_game_server/lib/codec/raw_pkg"
	"maze_game_server/lib/nano/internal/message"
	"maze_game_server/lib/nano/serialize"
)

const (
	Request  = message.Request
	Notify   = message.Notify
	Response = message.Response
	Push     = message.Push
)

type Type = message.Type
type Message = message.Message

type PacketProcessor interface {

	// Encode is used to encode a custom packet protocol.
	Encode(*Message) ([]byte, error)

	// Decode is used to decode a custom packet protocol.
	Decode([]byte) ([]*Message, []*raw_pkg.StruSvrEsRawBaseHead, error)
}

type PacketCodec interface {
	// NewProcessor returns a new instance for handling package parsing.
	NewProcessor() PacketProcessor

	// Serializer returns a serializer.
	Serializer() serialize.Serializer
}
