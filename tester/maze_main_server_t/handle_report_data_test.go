/*
 * @Author: majian
 * @Date: 2025-03-24 14:11:45
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-24 14:55:06
 */
package maze_main_server_t

import (
	"context"
	"fmt"
	"testing"

	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/protodef/Common"
	"gitlab.ifreetalk.com/plate/protodef/MazeGame"
)

func TestReportDataRQ(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	req := &MazeGame.ReportDataRQ{}
	data := MazeGame.ReportUserInfo{}
	data.ExpTotal = proto.Int64(375)
	data.ReportMask = proto.Int64(1)
	req.UserInfo = &data
	req.Header = &Common.PacketHeader{}
	req.Header.Session = proto.String("fdsfdsfd")
	res := &MazeGame.ReportDataRS{}
	e := game.OnReportDataRQ(ctx, 9003200130206264, req, res)
	l, _ := proto.Marshal(res)
	fmt.Println("l", len(l), "e", e, "res", res)
}
