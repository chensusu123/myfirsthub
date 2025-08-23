package mazebarrieruserrecordmodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strings"
)

const (
	GameRetSucc  = 1
	GameRetDeath = 2
	GameRetSweep = 3
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

// 用户迷宫闯关纪录
type MazeBarrierUserGameRecord struct {
	KafkaCommon
	UserId         uint64 `json:"user_id"`          // 用户id
	Barrier        int32  `json:"barrier"`          // 关卡id
	GameRet        int32  `json:"game_ret"`         // 用户闯关结果 1-通关成功 2-死亡失败 3-扫荡
	Awards         string `json:"awards"`           // 本次获得的奖励
	KillMonsterNum int64  `json:"kill_monster_num"` // 杀怪数量
}

const MazeBarrierUserRecordTableName = "maze_barrier_user_record"

func NewMazeBarrierUserGameRecord(useID uint64, barrier int32, gameRet int32, awards string, killMonsterNum int64) *MazeBarrierUserGameRecord {
	res := &MazeBarrierUserGameRecord{
		UserId:         useID,
		Barrier:        barrier,
		GameRet:        gameRet,
		Awards:         awards,
		KillMonsterNum: killMonsterNum,
	}
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeBarrierUserRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]
	return res
}
