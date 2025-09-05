package maze_es_t

import (
	"encoding/binary"
	"math/rand"
	"testing"
	"time"

	"maze_game_server/lib/net/websocket_service"

	"go.uber.org/zap"
)

func dataProcess(c *websocket_service.Client, data []byte) {
	messageLen := len(data)
	processLen := 0

	for processLen < messageLen {
		c.InfoWF("process bytes", zap.Any("processLen", processLen),
			zap.Any("messageLen", messageLen))
		n, err := c.ProcsssBytes(data[processLen:])
		if err != nil {
			c.ErrorWF("process bytes error", zap.Error(err))
			return
		}
		processLen += n
		c.InfoWF("process bytes", zap.Any("processLen", processLen),
			zap.Any("messageLen", messageLen))
	}
}

func makeData() []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(1024)
	if n < 3 {
		n = 3
	}

	data := make([]byte, n)
	for i := range data {
		// 生成 0~255 的随机字节
		data[i] = byte(rand.Intn(256))
	}
	binary.LittleEndian.PutUint16(data[0:2], uint16(n))
	gTestlogger.CtxInfo(ctx, "make data", zap.Any("data", len(data)))
	return data[0:n]
}

func randData() []byte {
	var data []byte
	for i := 0; i < 10; i++ {
		data = append(data, makeData()...)
	}
	return data
}

func TestData(t *testing.T) {
	xxx := websocket_service.Client{}
	xxx.MockLogger(gTestLogger)
	// xxx.ProcsssBytes([]byte("hello"))
	dataProcess(&xxx, randData())
}
