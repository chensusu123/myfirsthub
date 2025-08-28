package frame_model

import (
	"github.com/google/uuid"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"sync"
	"time"
)

const (
	ROLE_COMMON = 0 //普通成员
	ROLE_MASTER = 1 //房主
)

// 玩家输入结构
type Input struct {
	Type   string `json:"type"`   // 输入类型(如: move, attack)
	Data   string `json:"data"`   // 输入数据
	Client int64  `json:"client"` // 客户端时间戳
}

// 帧数据结构
type FrameData struct {
	FrameID   int                `json:"frameId"`   // 帧ID
	Inputs    map[uint64][]Input `json:"inputs"`    // 玩家输入: key=玩家ID, value=输入列表
	Timestamp int64              `json:"timestamp"` // 时间戳
}

// 玩家连接结构
type Player struct {
	ID            uint64 // 玩家ID
	LastInputTime int64  // 最后输入时间
	//TODO 腾讯
	IsReady       bool   //玩家准备状态
	Role          int32  //0：普通成员 1：房主
	PosNum        int32  //座位号，从 0 开始
	HeadImg       string //头像 URL（用户授权才会返回）
	NickName      string //用户昵称（用户授权才会返回）
	ClientId      int32  //用户在房间内的唯一标识
	EnableToStart bool   //是否已做好游戏开始准备
	MemberExtInfo string //给第三方用的 buffer，最长 32 个字节
	Logger        fklog.FKLogI
}

// 房间结构
type Room struct {
	FrameDataList []*FrameData       // 已同步的帧数据
	CurrentFrame  int                // 当前帧号
	InputQueue    map[uint64][]Input // 输入队列
	Mutex         sync.RWMutex       // 读写锁
	Ticker        *time.Ticker       // 定时器
	Running       bool               // 房间是否运行中
	//TODO 腾讯
	RoomIdStr              string             //房间 ID
	RoomState              int32              //1：组队中，2：该房间的对局游戏已开始，3：该房间的对局游戏已结束，4：房间已销毁
	MaxMemberNum           int32              //房间最大可容纳人数
	CreateTime             int64              //创建时间
	UpdateTimeStamp        int64              //最近更新时间
	GameTick               int32              //游戏下发帧的时间间隔，单位 ms
	StartPercent           int32              //真正开始帧同步需要达到多少百分比的玩家发送了开始指令，填 50 表征 50%
	GameLastTime           int32              //游戏对局时长，单位 s
	GameVersion            int64              //第三方自定义的游戏版本号
	GameAccessInfo         string             //该房间对应的游戏的 access_info 游戏唯一标识，用于后台接口拉取对局记录
	UdpReliabilityStrategy int32              //UDP 可靠性策略， 0：全冗余 N：固定冗余 N 帧
	RoomExtInfo            string             //给第三方用的 buffer，最长 32 个字节
	Seed                   string             //游戏随机种子
	MemberMap              map[uint64]*Player //成员列表
}

// 服务器结构
type Server struct {
	Rooms      map[string]*Room  // 房间列表
	PlayerRoom map[uint64]string // 玩家所在房间映射
	Mutex      sync.RWMutex      // 读写锁
}

func NewServer() *Server {
	return &Server{
		Rooms:      make(map[string]*Room),
		PlayerRoom: make(map[uint64]string),
	}
}

func NewRoomInfo(logger fklog.FKLogI, userIdList []uint64, gameTick, startPercent, gameLastTime, udpReliabilityStrategy int32, roomExtInfo string, needSeed bool) *Room {
	roomInfo := &Room{}
	roomInfo.FrameDataList = make([]*FrameData, 0)
	roomInfo.CurrentFrame = 0
	roomInfo.InputQueue = make(map[uint64][]Input)
	roomInfo.Running = false

	roomInfo.RoomIdStr = uuid.New().String()
	roomInfo.RoomState = 1
	roomInfo.MaxMemberNum = int32(len(userIdList))
	roomInfo.CreateTime = time.Now().Unix()
	roomInfo.UpdateTimeStamp = time.Now().Unix()
	roomInfo.GameTick = gameTick //最低33ms
	roomInfo.StartPercent = startPercent
	roomInfo.GameLastTime = gameLastTime
	roomInfo.GameVersion = 0
	roomInfo.GameAccessInfo = uuid.New().String()
	roomInfo.UdpReliabilityStrategy = udpReliabilityStrategy
	roomInfo.RoomExtInfo = roomExtInfo
	if needSeed {
		//source := rand.NewSource(uint64(time.Now().UnixNano()))
		//rand.New(source)
		roomInfo.Seed = time.Now().String()
	}
	roomInfo.MemberMap = make(map[uint64]*Player, len(userIdList))
	return roomInfo
}

func NewPlayer(logger fklog.FKLogI, userId uint64, role int32) *Player {
	return &Player{
		ID:     userId,
		Logger: logger,
		Role:   role,
	}
}

func NewFrameData(FrameID int) *FrameData {
	return &FrameData{
		FrameID:   FrameID,
		Inputs:    make(map[uint64][]Input),
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
	}
}

func NewInput(message string) *Input {
	return &Input{
		Type:   "msg",
		Data:   message,
		Client: time.Now().UnixNano(),
	}
}
