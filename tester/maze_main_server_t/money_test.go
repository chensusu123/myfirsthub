package maze_main_server_t

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"maze_game_server/common/errors"
	"maze_game_server/io/rpc/mazeitemrpc"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeItemSvr"
)

func TestQueryMoney(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	rq := &MazeItemSvr.QueryItemRQ{
		UserId: proto.Uint64(9003200130206203),
		Items: []*MazeCommon.MazeItem{
			&MazeCommon.MazeItem{ItemId: proto.Int32(46200001), Count: proto.Int64(0)},
			&MazeCommon.MazeItem{ItemId: proto.Int32(46900001), Count: proto.Int64(10000)},
		},
	}

	rs := &MazeItemSvr.QueryItemRS{ErrInfo: errors.NO_ERROR}
	err := mazeitemrpc.QueryItemsRQ(ctx, rq, rs)

	fmt.Println(err)
	time.Sleep(time.Second)
}

func TestAddMoney(t *testing.T) {
	ctx := fknet.TCPContext{Context: context.Background(), FKLogI: gTestLogger}
	rq := &MazeItemSvr.AddItemRQ{
		UserId: proto.Uint64(9003200130206202),
		OpType: proto.Int32(int32(697)),
		Header: nil,
		Items: []*MazeCommon.MazeItem{
			// &MazeCommon.MazeItem{ItemId: proto.Int32(46200001), Count: proto.Int64(0)},
			&MazeCommon.MazeItem{ItemId: proto.Int32(46900001), Count: proto.Int64(10000)},
		},
	}
	rq.TradeNumber = proto.Uint64(10000)

	rs := &MazeItemSvr.AddItemRS{ErrInfo: errors.NO_ERROR}
	err := mazeitemrpc.AddItemsRQ(ctx, rq, rs)

	fmt.Println(err)
	time.Sleep(time.Second)
}
