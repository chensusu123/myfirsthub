package friend

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	fileResolver "gitlab.ifreetalk.com/maze-plate/freetk/registry/fileresolver"
	globalredis "maze_game_server/io/redis"
	"maze_game_server/lib/log"
	"maze_game_server/model/friendmodel"
	"maze_game_server/services/friendservice"
	"os"
	"testing"
)

var logger = log.Clone("FriendTest", 0, 0)

func TestMain(m *testing.M) {
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	_ = os.Chdir("C:/work/maze_game_server/servers/maze_main_server/")
	_, err := fkserver.AppServer.Application.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = globalredis.GCli.Init(fileResolver.New("./conf.d/service.yaml"))
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Stdout = originalStdout
	os.Stderr = originalStderr

	m.Run()
}

func TestAddFriendRequest(t *testing.T) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return
	}
	res, err := db.Set(context.Background(), "111", "a", 0).Result()
	res2, err := db.Incr(context.Background(), "111").Result()
	_ = res2
	//res, err := db.Get(context.Background(), "111").Result()
	if err != nil {
		return
	}
	fmt.Println(res)
	//c := NewFriendComponent()
	//s := session.New(mock.NewNetworkEntity())
	//s.Bind(40000007)
	//req := &Friend.AddFriendRequestRQ{
	//	UserId: proto.Uint64(40000001),
	//}
	//err := c.OnAddFriendRequest_10552_10553(s, req)
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}
}

func TestReceiveFriendRequestList(t *testing.T) {
	//c := NewFriendComponent()
	//s := session.New(mock.NewNetworkEntity())
	//s.Bind(40000007)
	//req := &Friend.ReceiveFriendRequestListRQ{
	//	Page:     proto.Int32(1),
	//	PageSize: proto.Int32(20),
	//}
	//err := c.OnReceiveFriendRequestList_10425_10426(s, req)
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}s
	receives, codeErr := friendservice.GlobalFriendService.ReceiveFriendRequestList(logger, 40000001, 1, 20)
	if codeErr != nil {
		fmt.Println(codeErr)
		return
	}
	fmt.Printf("%v", receives)
}

func TestSendFriendRequestList(t *testing.T) {
	sends, codeErr := friendservice.GlobalFriendService.SendFriendRequestList(logger, 40000007, 1, 20)
	if codeErr != nil {
		fmt.Println(codeErr)
		return
	}
	fmt.Printf("%v", sends)
}

func TestAcceptFriendRequest(t *testing.T) {
	codeErr := friendservice.GlobalFriendService.AcceptFriendRequest(logger, 40000001, 40000007)
	if codeErr != nil {
		fmt.Println(codeErr)
		return
	}
}

func TestRejectFriendRequest(t *testing.T) {
	codeErr := friendservice.GlobalFriendService.RejectFriendRequest(logger, 40000001, 40000007)
	if codeErr != nil {
		fmt.Println(codeErr)
		return
	}
}

func TestFriendList(t *testing.T) {
	friends, codeErr := friendservice.GlobalFriendService.FriendList(logger, 40000001, 2, 20)
	if codeErr != nil {
		fmt.Println(codeErr)
		return
	}
	fmt.Printf("%v", friends)
}

func TestAddBlacklist(t *testing.T) {
	codeErr := friendservice.GlobalFriendService.AddBlacklist(logger, 40000001, 40000007)
	if codeErr != nil {
		fmt.Println(codeErr)
		return
	}
}

func TestRemoveBlacklist(t *testing.T) {
	codeErr := friendservice.GlobalFriendService.RemoveBlacklist(logger, 40000001, 40000007)
	if codeErr != nil {
		fmt.Println(codeErr)
		return
	}
}

func TestRemoveFriend(t *testing.T) {
	codeErr := friendservice.GlobalFriendService.RemoveFriend(logger, 40000001, 40000007)
	if codeErr != nil {
		fmt.Println(codeErr)
		return
	}
}

func Test_FriendService_Cases(t *testing.T) {
	svc := friendservice.GlobalFriendService

	// TC01 A向B申请好友，正常建立申请记录
	t.Run("TC01_AddFriendRequest_Normal", func(t *testing.T) {
		err := svc.AddFriendRequest(logger, 1001, 2002)
		svc.DeleteUserAll(logger, 1001)
		svc.DeleteUserAll(logger, 2002)
		assert.Nil(t, err)
	})

	// TC02 A多次申请B，幂等无重复
	t.Run("TC02_AddFriendRequest_Idempotent", func(t *testing.T) {
		err1 := svc.AddFriendRequest(logger, 1001, 2002)
		err2 := svc.AddFriendRequest(logger, 1001, 2002)
		assert.Nil(t, err1)
		assert.Nil(t, err2)
		//svc.DeleteUserAll(logger, 1001)
		//svc.DeleteUserAll(logger, 2002)
	})

	// TC03 A已经是B好友，申请失败
	t.Run("TC03_AlreadyFriends_RejectRequest", func(t *testing.T) {
		//_ = svc.AddFriendRequest(logger, 1001, 2002)
		_ = svc.AddFriendRequest(logger, 2002, 1001)
		//_ = svc.AcceptFriendRequest(logger, 2002, 1001)
		_ = svc.AcceptFriendRequest(logger, 1001, 2002)
		err := svc.AddFriendRequest(logger, 1001, 2002)
		assert.NotNil(t, err, "应该拒绝已是好友的申请")
		//svc.DeleteUserAll(logger, 1001)
		//svc.DeleteUserAll(logger, 2002)
	})

	// TC04 B把A拉黑，申请失败
	t.Run("TC04_BlockPreventsRequest", func(t *testing.T) {
		_ = svc.AddBlacklist(logger, 2002, 1001)
		err := svc.AddFriendRequest(logger, 1001, 2002)
		assert.NotNil(t, err, "被拉黑后申请应失败")
		//svc.DeleteUserAll(logger, 1001)
		//svc.DeleteUserAll(logger, 2002)
	})

	// TC05 B同意A申请，建立好友关系
	t.Run("TC05_AcceptFriendRequest", func(t *testing.T) {
		_ = svc.AddFriendRequest(logger, 1001, 2002)
		err := svc.AcceptFriendRequest(logger, 2002, 1001)
		assert.Nil(t, err)
		// 验证好友关系存在
		friends, errList := svc.FriendList(logger, 1001, 1, 10)
		assert.Nil(t, errList)
		assert.Equal(t, 1, len(friends))
		//svc.DeleteUserAll(logger, 1001)
		//svc.DeleteUserAll(logger, 2002)
	})

	// TC06 已是好友再次同意，幂等成功
	t.Run("TC06_AcceptAlreadyFriend", func(t *testing.T) {
		_ = svc.AddFriendRequest(logger, 1001, 2002)
		_ = svc.AcceptFriendRequest(logger, 2002, 1001)
		err := svc.AcceptFriendRequest(logger, 2002, 1001)
		assert.Nil(t, err)
		//svc.DeleteUserAll(logger, 1001)
		//svc.DeleteUserAll(logger, 2002)
	})

	// TC07 无申请直接同意，拒绝操作
	t.Run("TC07_AcceptWithoutRequest_Fail", func(t *testing.T) {
		err := svc.AcceptFriendRequest(logger, 3003, 4004)
		assert.NotNil(t, err)
		svc.DeleteUserAll(logger, 1001)
		svc.DeleteUserAll(logger, 2002)
	})

	// TC08 B拒绝A申请，状态变为 rejected
	t.Run("TC08_RejectFriendRequest", func(t *testing.T) {
		_ = svc.AddFriendRequest(logger, 1001, 2002)
		err := svc.RejectFriendRequest(logger, 2002, 1001)
		assert.Nil(t, err)
		recvList, _ := svc.ReceiveFriendRequestList(logger, 2002, 1, 10)
		assert.Equal(t, friendmodel.FriendRequestStatusRejected, recvList[0].Status)
		//svc.DeleteUserAll(logger, 1001)
		//svc.DeleteUserAll(logger, 2002)
	})

	// TC09 已拒绝再拒绝，幂等
	t.Run("TC09_RejectAlreadyRejected", func(t *testing.T) {
		_ = svc.AddFriendRequest(logger, 1001, 2002)
		_ = svc.RejectFriendRequest(logger, 2002, 1001)
		err := svc.RejectFriendRequest(logger, 2002, 1001)
		assert.Nil(t, err)
		//svc.DeleteUserAll(logger, 1001)
		//svc.DeleteUserAll(logger, 2002)
	})

	// TC10 拒绝时已是好友，不清除好友，只改状态
	t.Run("TC10_RejectWhileFriend", func(t *testing.T) {
		_ = svc.AddFriendRequest(logger, 1001, 2002)
		_ = svc.AcceptFriendRequest(logger, 2002, 1001)
		err := svc.RejectFriendRequest(logger, 2002, 1001)
		assert.Nil(t, err)
		friends, _ := svc.FriendList(logger, 1001, 1, 10)
		assert.Equal(t, 1, len(friends))
		//svc.DeleteUserAll(logger, 1001)
		//svc.DeleteUserAll(logger, 2002)
	})

	// TC11 A、B同时申请对方，均成功
	t.Run("TC11_MutualRequests", func(t *testing.T) {
		err1 := svc.AddFriendRequest(logger, 1001, 2002)
		err2 := svc.AddFriendRequest(logger, 2002, 1001)
		assert.Nil(t, err1)
		assert.Nil(t, err2)
		//svc.DeleteUserAll(logger, 1001)
		//svc.DeleteUserAll(logger, 2002)
	})

	// TC12 A、B同时同意，好友关系建立成功且无重复
	t.Run("TC12_MutualAccept", func(t *testing.T) {
		_ = svc.AddFriendRequest(logger, 1001, 2002)
		_ = svc.AddFriendRequest(logger, 2002, 1001)
		err1 := svc.AcceptFriendRequest(logger, 1001, 2002)
		err2 := svc.AcceptFriendRequest(logger, 2002, 1001)
		assert.Nil(t, err1)
		assert.Nil(t, err2)
		friendsA, _ := svc.FriendList(logger, 1001, 1, 10)
		friendsB, _ := svc.FriendList(logger, 2002, 1, 10)
		assert.Equal(t, 1, len(friendsA))
		assert.Equal(t, 1, len(friendsB))
		//svc.DeleteUserAll(logger, 1001)
		//svc.DeleteUserAll(logger, 2002)
	})

	// TC13 A拉黑B，清理所有好友及申请
	t.Run("TC13_BlockAndClear", func(t *testing.T) {
		_ = svc.AddFriendRequest(logger, 1001, 2002)
		_ = svc.AcceptFriendRequest(logger, 2002, 1001)
		err := svc.AddBlacklist(logger, 1001, 2002)
		assert.Nil(t, err)
		friends, _ := svc.FriendList(logger, 1001, 1, 10)
		assert.Equal(t, 0, len(friends))
		reqList, _ := svc.ReceiveFriendRequestList(logger, 2002, 1, 10)
		assert.Equal(t, 0, len(reqList))
		//svc.DeleteUserAll(logger, 1001)
		//svc.DeleteUserAll(logger, 2002)
	})

	// TC14 A重复拉黑B，幂等
	t.Run("TC14_BlockIdempotent", func(t *testing.T) {
		err1 := svc.AddBlacklist(logger, 1001, 2002)
		err2 := svc.AddBlacklist(logger, 1001, 2002)
		assert.Nil(t, err1)
		assert.Nil(t, err2)
		svc.DeleteUserAll(logger, 1001)
		svc.DeleteUserAll(logger, 2002)
	})

	// TC15 A拉黑后B申请，失败
	t.Run("TC15_BlockPreventsRequest", func(t *testing.T) {
		_ = svc.AddBlacklist(logger, 1001, 2002)
		err := svc.AddFriendRequest(logger, 2002, 1001)
		assert.NotNil(t, err)
		//svc.DeleteUserAll(logger, 1001)
		//svc.DeleteUserAll(logger, 2002)
	})

	// TC16 查看是否为好友
	t.Run("TC16_CheckFriendship", func(t *testing.T) {
		_ = svc.AddFriendRequest(logger, 1001, 2002)
		_ = svc.AcceptFriendRequest(logger, 2002, 1001)
		list, _ := svc.FriendList(logger, 1001, 1, 10)
		assert.Equal(t, 1, len(list))
		svc.DeleteUserAll(logger, 1001)
		svc.DeleteUserAll(logger, 2002)
	})

	// TC17 查看申请状态
	t.Run("TC17_CheckRequestStatus", func(t *testing.T) {
		_ = svc.AddFriendRequest(logger, 1001, 2002)
		sendList, _ := svc.SendFriendRequestList(logger, 1001, 1, 10)
		assert.Equal(t, friendmodel.FriendRequestStatusPending, sendList[0].Status)
		svc.DeleteUserAll(logger, 1001)
		svc.DeleteUserAll(logger, 2002)
	})
}
