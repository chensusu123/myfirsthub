package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/lib/net/raw_pkg"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/pb/common/SysPackDef"
)

var gAddr = "127.0.0.1:9876"

func client(logger fklog.FKLogI, addr string) {
	log.SetFlags(0)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	logger.DebugWF("connecting to", zap.String("url", u.String()))

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.ErrorWF("dial:", zap.Error(err))
		return
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			logger.DebugWF("recv message")
			messageType, message, err := c.ReadMessage()
			if err != nil {
				logger.ErrorWF("read:", zap.Error(err))
				return
			}
			logger.DebugWF("recv message",
				zap.Int("messageType", messageType),
				zap.Any("message", len(message)))

		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	loopCount := uint16(0)
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			if loopCount > 4 {
				// c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				// break
			}
			_ = t
			var sendData []byte
			if loopCount == 0 {
				sendData = makeLoginData()
			} else {
				continue
				sendData = makeOtherData()
			}

			logger.DebugWF("send message",
				zap.Any("message", len(sendData)))
			err := c.WriteMessage(websocket.BinaryMessage, sendData)
			if err != nil {
				logger.ErrorWF("write:", zap.Error(err))
				return
			}
			loopCount += 1
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func initLog() {
	logConfig := fklog.LogConfig{
		LogDir:     "/tmp",
		LogLevel:   "debug",
		LogName:    "ws_test_cli",
		LogType:    "zap",
		WithCaller: true,
	}
	fklog.InitAppFkLog(&logConfig)
}

func main() {
	initLog()
	gTestLogger := fklog.AppLogger().Clone("maze_energy_server_t")
	client(gTestLogger, gAddr)
}

func makeLoginData() []byte {
	pbPacket := SysPackDef.UserLoginRq{
		UserID:        proto.Uint64(999),
		UserAuthToken: proto.Uint64(2224),
	}
	pbData, _ := proto.Marshal(&pbPacket)
	pkg := &raw_pkg.StruSvrEsRawBaseHead{}
	pkg.PackType = 5183
	pkg.Data = pbData
	data, _ := pkg.Pack()
	return data
}

func makeOtherData() []byte {
	pbPacket := MazeGame.BarrierDeathRQ{
		BarrierId: proto.Int32(222),
	}
	pbData, _ := proto.Marshal(&pbPacket)
	pkg := &raw_pkg.StruSvrEsRawBaseHead{}
	pkg.PackType = 16212
	pkg.Data = pbData
	data, _ := pkg.Pack()
	return data
}
