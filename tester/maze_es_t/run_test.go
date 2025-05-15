package maze_es_t

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeGame"
	"gitlab.ifreetalk.com/maze-plate/protodef/SysPackDef"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/raw_pkg"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_es/process"
	"go.uber.org/zap"
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
				// continue
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

func initService() {
	websocket_service.RegProcSimple(
		5183, &SysPackDef.UserLoginRq{},
		5184, &SysPackDef.UserLoginRs{},
		process.OnLoginRQ)

	websocket_service.RegProcSimple(
		16212, &MazeGame.BarrierDeathRQ{},
		16213, &MazeGame.BarrierDeathRS{},
		process.OnTestRQ)
}

func sendMsgToUser() {
	logger := gTestLogger.Clone("sendMsgToUser")
	for i := 0; i < 10; i++ {
		data := makeOtherData()
		websocket_service.SendBytes(logger, 999, 1, data)
		time.Sleep(time.Second * 1)
	}
}

func TestWebsocket(t *testing.T) {
	logger := gTestLogger.Clone("TestWebsocket")
	initService()
	websocket_service.MockOnInit(logger, gAddr)
	websocket_service.MockOnStart(logger)

	gTestLogger.DebugWF("a")

	go client(logger, gAddr)
	go sendMsgToUser()
	// go client(logger, gAddr)
	time.Sleep(time.Second * 120)
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
