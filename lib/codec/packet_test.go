package codec

import (
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"

	"maze_game_server/lib/codec/raw_pkg"
	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeGame"

	"google.golang.org/protobuf/proto"
)

func TestSendPacket(t *testing.T) {

	pb := &MazeGame.MazeLoginRQ{}
	pb.Header = &Common.PacketHeader{
		Session:       proto.String("2"),
		Protocol:      proto.Int(370),
		RanchProtocol: proto.Int(1),
	}
	pb.MazeVersion = proto.Int(2)

	pd, _ := proto.Marshal(pb)

	stru := raw_pkg.StruSvrEsRawBaseHead{}
	stru.PackType = 10451
	stru.SessionID = uint32(time.Now().Unix())
	stru.EsRqTime = uint64(time.Now().UnixMilli())
	// stru.Data = []byte("{\"message\":\"hello\"}")
	stru.Data = pd
	data, err := stru.Pack()
	if err != nil {
		t.Errorf("Pack error: %s", err.Error())
		return
	}

	// t.Logf("data: %s", base64.StdEncoding.EncodeToString(data))
	// return

	conn, err := net.Dial("tcp", ":5997")
	if err != nil {
		t.Errorf("Dial error: %s", err.Error())
		return
	}

	conn.Write(data)

	twoBytes := [2]byte{}

	_, err = io.ReadAtLeast(conn, twoBytes[:], 2)
	if err != nil {
		t.Errorf("ReadAtLeast error: %s", err.Error())
		return
	}

	packetLen := int(binary.LittleEndian.Uint16(twoBytes[:]))

	packet := make([]byte, packetLen)
	copy(packet, twoBytes[:])
	io.ReadAtLeast(conn, packet[2:], packetLen-2)

	stru = raw_pkg.StruSvrEsRawBaseHead{}
	err = stru.UnPack(packet)
	if err != nil {
		t.Errorf("UnPack error: %s", err.Error())
		return
	}

	t.Logf("packet: %+v", packet)
	t.Logf("stru: %+v", stru)
}
