package family

import (
	"fmt"
	"maze_game_server/io/redis/allianceredis"
	"maze_game_server/io/redis/familyredis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/mock"
	"maze_game_server/lib/nano/session"
	"maze_game_server/model/familymodel"
	"maze_game_server/pb/common/MazeFamily"
	"maze_game_server/servers/maze_main_server/process/alliance"
	"os"
	"testing"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	fileResolver "gitlab.ifreetalk.com/maze-plate/freetk/registry/fileresolver"
	"google.golang.org/protobuf/proto"
)

var logger = log.Clone("FriendTest", 0, 0)

func TestMain(m *testing.M) {
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	_ = os.Chdir("/Users/cltx/Desktop/project/src/maze_game_server/servers/maze_main_server")
	_, err := fkserver.AppServer.Application.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = familyredis.GlobalFamilyRedis.Init(fileResolver.New("./conf.d/service.yaml"))
	err = allianceredis.GlobalAllianceRedis.Init(fileResolver.New("./conf.d/service.yaml"))
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Stdout = originalStdout
	os.Stderr = originalStderr

	m.Run()
}

func TestGetCreateFamilyCost(t *testing.T) {
	c := NewFamily()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000007)

	err := c.OnGetCreateFamilyCostRQ_10582_10583(s, &MazeFamily.GetCreateFamilyCostRQ{})

	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestOnCreateFamilyRQ_10585_10586(t *testing.T) {
	c := NewFamily()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000014)

	familyJoinType := MazeFamily.FamilyJoinType_FAMILY_JOIN_TYPE_NONE
	privilegeLevel := MazeFamily.PrivilegeLevel_FAMILY_PRIVILEGE_LEVEL_LEADER
	err := c.OnCreateFamilyRQ_10585_10586(s, &MazeFamily.CreateFamilyRQ{
		FamilyName: proto.String("test"),
		Leader: &MazeFamily.FamilyMemberInfo{
			NickName:       proto.String("test"),
			UserId:         proto.Uint64(40000014),
			PrivilegeLevel: &privilegeLevel,
		},
		AllianceId: proto.Int32(5),
		JoinType:   &familyJoinType,
	})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

}

func TestCreatAlliance(t *testing.T) {
	c := alliance.NewAlliance()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000014)

	err := c.OnCreateAllianceRQ_10610_10611(s, "test")
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestGetFamilyList(t *testing.T) {
	fmt.Printf("START-------------------------------\n")
	familyListModel, err := familymodel.LoadFamilyListModel(logger)
	if err != nil {
		t.Error(err.Error())
		return
	}
	// 查询家族列表
	familyInfoModel, err := familymodel.LoadFamilyListInfoModel(logger, familyListModel.Familys)
	if err != nil {
		t.Error(err.Error())
		return
	}

	for _, family := range familyInfoModel {
		fmt.Printf("familyID: %d, familyName: %s, familyLevel: %d, MemberCount: %d\n",
			family.FamilyID, family.FamilyName, family.FamilyLevel, family.MemberCount)
		fmt.Printf("UserList:\n")
		for _, member := range family.FamilyMembers {
			fmt.Printf("memberID: %d, memberName: %s, memberLevel: %d, memberPrivilegeLevel: %d\n",
				member.UserID, member.NickName, member.Level, member.PrivilegeLevel)
		}
		fmt.Printf("ApplyUserList:\n")
		for _, member := range family.FamilyApplyUsers {
			fmt.Printf("memberID: %d, memberName: %s, memberLevel: %d, memberPrivilegeLevel: %d\n",
				member.UserID, member.NickName, member.Level, member.PrivilegeLevel)
		}
		fmt.Printf("--------------------------------\n")
	}
	fmt.Printf("END-------------------------------\n")
}

func TestOnGetUpgradeFamilyCostRQ_10574_10575(t *testing.T) {
	c := NewFamily()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000014)

	err := c.OnGetUpgradeFamilyCostRQ_10574_10575(s, &MazeFamily.GetUpgradeFamilyCostRQ{
		FamilyId:    proto.Int32(5),
		TargetLevel: proto.Int32(2),
	})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	fmt.Println("End--------------------------------")
}

func TestOnUpgradeFamilyRQ_10576_10577(t *testing.T) {
	c := NewFamily()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000014)

	err := c.OnUpgradeFamilyRQ_10576_10577(s, &MazeFamily.UpgradeFamilyRQ{
		FamilyId:    proto.Int32(5),
		TargetLevel: proto.Int32(3),
	})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	fmt.Println("End--------------------------------")
}

func TestOnGetFamilyInfoRQ_10580_10581(t *testing.T) {
	c := NewFamily()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000014)

	err := c.OnGetFamilyInfoRQ_10580_10581(s, &MazeFamily.GetFamilyInfoRQ{
		FamilyId: proto.Int32(5),
	})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	fmt.Println("End--------------------------------")
}

func TestOnApplyFamilyRQ_10587_10588(t *testing.T) {
	c := NewFamily()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000017)

	privilegeLevel := MazeFamily.PrivilegeLevel_FAMILY_PRIVILEGE_LEVEL_MEMBER

	err := c.OnApplyFamilyRQ_10587_10588(s, &MazeFamily.ApplyFamilyRQ{
		FamilyId: proto.Int32(7),
		ApplyUser: &MazeFamily.FamilyMemberInfo{
			UserId:         proto.Uint64(40000017),
			NickName:       proto.String("test"),
			PrivilegeLevel: &privilegeLevel,
		},
	})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	fmt.Println("End--------------------------------")
}

func TestOnConfirmApplyFamilyRQ_10589_10590(t *testing.T) {
	c := NewFamily()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000014)

	privilegeLevel := MazeFamily.PrivilegeLevel_FAMILY_PRIVILEGE_LEVEL_MEMBER
	result := MazeFamily.ApplyFamilyResult_APPLY_FAMILY_RESULT_AGREE

	err := c.OnConfirmApplyFamilyRQ_10589_10590(s, &MazeFamily.ConfirmApplyFamilyRQ{
		FamilyId: proto.Int32(7),
		ApplyUser: &MazeFamily.FamilyMemberInfo{
			UserId:         proto.Uint64(40000017),
			NickName:       proto.String("test"),
			PrivilegeLevel: &privilegeLevel,
		},
		Result: &result,
	})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	fmt.Println("End--------------------------------")

}

func TestOnOperatePrivilegeRQ_10592_10593(t *testing.T) {
	c := NewFamily()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000014)

	privilegeLevel := MazeFamily.PrivilegeLevel_FAMILY_PRIVILEGE_LEVEL_MANAGER
	err := c.OnOperatePrivilegeRQ_10592_10593(s, &MazeFamily.OperatePrivilegeRQ{
		FamilyId:             proto.Int32(5),
		TargetPrivilegeLevel: &privilegeLevel,
		OperateUsers:         []uint64{40000015},
	})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	fmt.Println("End--------------------------------")
}

func TestOnKickFamilyRQ_10594_10595(t *testing.T) {
	c := NewFamily()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000014)

	err := c.OnKickFamilyRQ_10594_10595(s, &MazeFamily.KickFamilyRQ{
		FamilyId:  proto.Int32(5),
		KickUsers: []uint64{40000015},
	})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	fmt.Println("End--------------------------------")
}

func TestOnGetUserFamliyRQ_10600_10601(t *testing.T) {
	c := NewFamily()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000016)

	err := c.OnGetUserFamliyRQ_10600_10601(s, &MazeFamily.GetUserFamilyRQ{
		UserId: proto.Uint64(40000016),
	})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	fmt.Println("End--------------------------------")
}

func TestOnExitFamilyRQ_10597_10598(t *testing.T) {
	c := NewFamily()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000017)

	err := c.OnExitFamilyRQ_10597_10598(s, &MazeFamily.ExitFamilyRQ{
		FamilyId: proto.Int32(7),
	})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	fmt.Println("End--------------------------------")
}

func TestOnDissolutionFamilyRQ_10602_10603(t *testing.T) {
	c := NewFamily()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000014)

	err := c.OnDissolutionFamilyRQ_10602_10603(s, &MazeFamily.DissolutionFamilyRQ{
		FamilyId: proto.Int32(6),
	})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	fmt.Println("End--------------------------------")
}

func TestOnGetAllianceListRQ_10606_10607(t *testing.T) {
	c := alliance.NewAlliance()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000017)

	err := c.OnGetAllianceListRQ_10606_10607(s, &MazeFamily.GetAllianceListRQ{})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	fmt.Println("End--------------------------------")
}

func TestOnGetUserAllianceRQ_10604_10605(t *testing.T) {
	c := alliance.NewAlliance()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000017)

	err := c.OnGetAllianceInfoRQ_10604_10605(s, &MazeFamily.GetUserAllianceRQ{
		AllianceId: proto.Int32(5),
	})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	fmt.Println("End--------------------------------")
}

func TestOnQueryUserAllianceRQ_10608_10609(t *testing.T) {
	c := alliance.NewAlliance()
	s := session.New(mock.NewNetworkEntity())
	s.Bind(40000017)

	err := c.OnQueryUserAllianceRQ_10608_10609(s, &MazeFamily.QueryUserAllianceRQ{
		UserId: proto.Uint64(40000017),
	})

	if err != nil {
		t.Fatal(err.Error())
		return
	}

	fmt.Println("End--------------------------------")
}
