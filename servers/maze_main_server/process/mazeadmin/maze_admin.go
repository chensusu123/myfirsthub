package mazeadmin

import (
	"context"
	"strconv"

	"maze_game_server/io/broadcastcli"
	"maze_game_server/usecase/online"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/admin"
)

func init() {
	admin.HandleFunc(admin.MethodGet, "/cmds/users/show", showusers)
	admin.HandleFunc(admin.MethodGet, "/cmds/users/isonline", isonline)
	admin.HandleFunc(admin.MethodGet, "/cmds/users/pushmsg", pushMsg)
	admin.HandleFunc(admin.MethodGet, "/cmds/users/broadcast", broadcast)
}

type ReturnMsg struct {
	Code int    `json:"status"`
	Msg  string `json:"desc"`
	Data any    `json:"data,omitempty"`
}

func MakeErrReturnMsg(code int, msg string) ReturnMsg {
	return ReturnMsg{
		Code: code,
		Msg:  msg,
	}
}

func MakeSuccessReturnMsg(data any) ReturnMsg {
	return ReturnMsg{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

func showusers(ctx context.Context, c *app.RequestContext) {
	userIDs := online.GetOnlineUsers()
	c.JSON(consts.StatusOK, MakeSuccessReturnMsg(userIDs))
}

func isonline(ctx context.Context, c *app.RequestContext) {
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(consts.StatusOK, MakeErrReturnMsg(400, "userID is empty"))
		return
	}
	userIDUint64, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(consts.StatusOK, MakeErrReturnMsg(400, "userID is invalid"))
		return
	}
	isOnline := online.IsOnline(uint64(userIDUint64))
	c.JSON(consts.StatusOK, MakeSuccessReturnMsg(isOnline))
}

func pushMsg(ctx context.Context, c *app.RequestContext) {
	userID := c.Query("userID")
	if userID == "" {
		c.JSON(consts.StatusOK, MakeErrReturnMsg(400, "userID is empty"))
		return
	}
	userIDUint64, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(consts.StatusOK, MakeErrReturnMsg(400, "userID is invalid"))
		return
	}
	msg := c.Query("msg")
	err = online.PushToClusterTest(ctx, uint64(userIDUint64), 222, []byte(msg))
	c.JSON(consts.StatusOK, MakeSuccessReturnMsg(err))
}

func broadcast(ctx context.Context, c *app.RequestContext) {
	userID := c.Query("broadcastID")
	if userID == "" {
		c.JSON(consts.StatusOK, MakeErrReturnMsg(400, "broadcastID is empty"))
		return
	}
	userIDUint64, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		c.JSON(consts.StatusOK, MakeErrReturnMsg(400, "broadcastID is invalid"))
		return
	}
	msg := c.Query("msg")
	broadcastcli.BroadcastTest(ctx, uint64(userIDUint64), 222, []byte(msg))
	c.JSON(consts.StatusOK, MakeSuccessReturnMsg(err))
}
