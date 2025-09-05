package mazemonstermodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strings"
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

const MazeMonsterRecordTableName = "maze_monster_record"

// 区域杀怪上报流水
type MazeMonsterRecordModel struct {
	KafkaCommon
	UserId         uint64 `json:"user_id"`          // 用户id
	BarrierID      uint32 `json:"barrier_id"`       // 关卡id
	AreaID         int32  `json:"area_id"`          // 区域id
	AreaIndex      int32  `json:"area_index"`       // 子区域id
	NowEquipScore  uint32 `json:"now_equip_score"`  // 当前装备分数
	NowItem1Score  uint32 `json:"now_coin_score"`   // 金币堆分数
	NowItem2Score  uint32 `json:"now_stone_score"`  // 强化石堆分数
	KillMonsterNum uint32 `json:"kill_monster_num"` // 当前杀怪数
	MonsterGuid    int64  `json:"monster_guid"`     // 怪物配置id
	MonsterPos     string `json:"monster_pos"`      // 怪物位置
	DropItems      string `json:"drop_items"`       // 掉落物品
}

func NewMazeMonsterRecordModel(userID uint64, barrierID uint32, areaID, areaIndex int32, nowEquipScore, nowItem1Score, nowItem2Score, killMonsterNum uint32, monsterGuid int64, monsterPos, dropItems string) *MazeMonsterRecordModel {
	res := &MazeMonsterRecordModel{
		UserId:         userID,
		BarrierID:      barrierID,
		AreaID:         areaID,
		AreaIndex:      areaIndex,
		NowEquipScore:  nowEquipScore,
		NowItem1Score:  nowItem1Score,
		NowItem2Score:  nowItem2Score,
		KillMonsterNum: killMonsterNum,
		MonsterGuid:    monsterGuid,
		MonsterPos:     monsterPos,
		DropItems:      dropItems,
	}
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeMonsterRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]
	return res
}
