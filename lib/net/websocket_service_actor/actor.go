package websocket_service_actor

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/raw_pkg"
	"go.uber.org/zap"
)

type WsWriteActor struct {
	conn      *websocket.Conn
	isJson    bool
	userID    uint64
	sessionId uint64
}

func (a *WsWriteActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *uint64:
		a.userID = *msg
		ctx.Logger().Info("WsWriteActor set userID", slog.Uint64("userID", a.userID), slog.Uint64("sessionId", a.sessionId))
	case []byte:
		ctx.Logger().Info("WsWriteActor send message",
			slog.Int("messageLen", len(msg)),
			slog.Uint64("userID", a.userID),
			slog.Uint64("sessionId", a.sessionId),
		)
		a.conn.SetWriteDeadline(time.Now().Add(writeWait))
		a.conn.WriteMessage(websocket.BinaryMessage, msg)
	case *SendDataMsg:
		ctx.Logger().Info("WsWriteActor send SendDataMsg", slog.Uint64("userID", a.userID),
			slog.Uint64("sessionId", a.sessionId), slog.Uint64("PacketType", uint64(msg.PacketType)))
		var (
			data []byte
			err  error
		)
		if a.isJson {
			sendPacket := &NoramlJsonMsg{
				MsgType: int(msg.PacketType),
				Data:    msg.Data,
			}
			data, err = json.Marshal(sendPacket)
			if err != nil {
				ctx.Logger().Error("WsWriteActor SendData Marshal json failed",
					slog.Uint64("userID", a.userID),
					slog.Uint64("sessionId", a.sessionId),
					slog.Any("err", err),
					slog.Any("packetType", msg.PacketType))
				return
			}
		} else {
			dataPB, err := proto.Marshal(msg.Data.(proto.Message))
			pkg := &raw_pkg.StruSvrEsRawBaseHead{}
			pkg.SessionID = uint32(0)
			pkg.PackType = msg.PacketType
			pkg.Data = dataPB
			pkg.EsRsTime = uint64(time.Now().Unix())
			pkg.SetTeaflag()
			data, err = pkg.Pack()
			if err != nil {
				ctx.Logger().Error("WsWriteActor SendData Marshal pb failed",
					slog.Uint64("userID", a.userID),
					slog.Uint64("sessionId", a.sessionId),
					slog.Any("err", err),
					slog.Any("packetType", msg.PacketType))
				return
			}
		}
		a.conn.SetWriteDeadline(time.Now().Add(writeWait))
		a.conn.WriteMessage(websocket.BinaryMessage, data)
	}
}

func serveActorWs(ctx *app.RequestContext, logger fklog.FKLogI, actorSystem *actor.ActorSystem, connMgrPID *actor.PID, isJson bool) {
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		clientLogger := logger.Clone(fmt.Sprintf("client-addr:%s:%v", conn.RemoteAddr().String(), isJson))
		clientLogger.InfoWF("serveActorJsonWs client connected", zap.Any("addr", conn.RemoteAddr().String()), zap.Bool("isJson", isJson))
		client := &Client{conn: conn, send: make(chan []byte, 1024), FkTags: fknet.NewFkTags(), sessionId: hub.MakeSession()}
		client.SetTag("mySelf", client)
		client.SetTag("mySelfSession", client.sessionId)
		client.FKLogI = clientLogger
		client.isJson = isJson
		defer conn.Close()

		props := actor.PropsFromProducer(func() actor.Actor {
			return &WsWriteActor{conn: conn, isJson: isJson, sessionId: client.sessionId}
		})
		pid := actorSystem.Root.Spawn(props)

		client.actorSystem = actorSystem
		client.writeActor = pid
		clientProps := actor.PropsFromProducer(func() actor.Actor {
			return client
		})

		clientPID := actorSystem.Root.Spawn(clientProps)

		defer func() {
			actorSystem.Root.Send(connMgrPID, &ClientLogout{
				wPID:      pid,
				rPID:      clientPID,
				SessionId: client.sessionId,
				UserId:    client.GetUserID(),
			})
			actorSystem.Root.Stop(pid)
			actorSystem.Root.Stop(clientPID)
		}()

		// 消息循环
		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				clientLogger.ErrorWF("serveActorJsonWs read error", zap.Error(err))
				break
			}
			actorSystem.Root.Send(clientPID, msgBytes)
		}
	})
	if err != nil {
		logger.ErrorWF("serveActorJsonWs upgrade error", zap.Error(err))
	}
}
