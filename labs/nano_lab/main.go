package main

import (
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/session"
	"gitlab.ifreetalk.com/maze-plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process"
	"go.uber.org/zap"
)

type (
	NaonoServiceHandler struct {
		component.Base
		fklog.FKLogI
	}
)

// Init was called to initialize the component.
func (c *NaonoServiceHandler) Init() {
	c.InfoWF("NaonoServiceHandler Init")
}

// AfterInit was called after the component is initialized.
func (c *NaonoServiceHandler) AfterInit() {
	c.InfoWF("NaonoServiceHandler AfterInit")
}

// BeforeShutdown was called before the component to shutdown.
func (c *NaonoServiceHandler) BeforeShutdown() {
	c.InfoWF("NaonoServiceHandler BeforeShutdown")
}

// Shutdown was called to shutdown the component.
func (c *NaonoServiceHandler) Shutdown() {
	c.InfoWF("NaonoServiceHandler Shutdown")
}

func (m *NaonoServiceHandler) RecvData(s *session.Session, data []byte) error {
	m.InfoWF("RecvData", zap.Any("data", len(data)))
	s.Response(data)
	return nil
}

func NewNaonoServiceHandler(logger fklog.FKLogI) *NaonoServiceHandler {
	return &NaonoServiceHandler{
		FKLogI: logger.Clone("NewNaonoServiceHandler"),
	}
}

type NanoService struct {
	fklog.FKLogI
}

func NewNanoService(logger fklog.FKLogI) *NanoService {
	return &NanoService{
		FKLogI: logger.Clone("NewNanoService"),
	}
}

func (m *NanoService) Start(config fkcore.FkConfigerI) {
	go func() {
		dataRecv := NewNaonoServiceHandler(m.FKLogI)
		components := &component.Components{}
		components.Register(dataRecv)
		tcpCfg := fkconfig.GetServerConfig()
		_, port, err := net.SplitHostPort(tcpCfg.GetTCPAddr())
		if err != nil {
			fkfmt.Println("websocket-service split ip/port failed.", tcpCfg.GetTCPAddr(), err)
			m.ErrorWF("websocket-service split ip/port failed.", zap.String("addr", tcpCfg.GetTCPAddr()), zap.Error(err))
			return
		}
		addr := ":" + port
		addr = ":5887"
		nano.Listen(addr,
			nano.WithDebugMode(),
			nano.WithIsWebsocket(true),
			// nano.WithSerializer(protobuf.NewSerializer()),
			// nano.WithSerializer(json.NewSerializer()),
			nano.WithComponents(components),
			nano.WithWSPath("/nano"),
		)
	}()
}

func initLog() {
	logConfig := fklog.LogConfig{
		LogDir:     "/tmp",
		LogLevel:   "debug",
		LogName:    "new_test",
		LogType:    "zap",
		WithCaller: true,
	}
	fklog.InitAppFkLog(&logConfig)
}

func main() {
	initLog()
	gTestLogger := fklog.AppLogger().Clone("maze_es_t")
	handler := NewNaonoServiceHandler(gTestLogger)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		nano.Listen(":5535",
			nano.WithDebugMode(),
			nano.WithIsWebsocket(true),
			// nano.WithSerializer(protobuf.NewSerializer()),
			// nano.WithSerializer(json.NewSerializer()),
			nano.WithComponents(process.Components()),
			nano.WithBytesFunc(handler.RecvData),
			nano.WithWSPath("/echo"),
		)
	}()

	gTestLogger.InfoWF("websocket-service start")
	time.Sleep(time.Second * 2)
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:5535", Path: "/echo"}

	gTestLogger.InfoWF("connecting to", zap.String("url", u.String()))

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %v %v", messageType, len(message))
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			_ = t
			err := c.WriteMessage(websocket.BinaryMessage, makeData(gTestLogger))
			if err != nil {
				log.Println("write:", err)
				return
			}
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

func makeData(logger fklog.FKLogI) []byte {
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

	logger.InfoWF("makeData", zap.Int("n", n))
	return data[0:n]
}
