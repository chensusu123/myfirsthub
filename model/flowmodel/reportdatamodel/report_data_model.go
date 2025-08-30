package reportdatamodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"maze_game_server/pb/common/MazeGame"
	"strings"
)

const MazeUserLevelRecordTableName = "maze_cli_report_data"

type KafkaCommon = kafkacommonstruct.KafkaCommon

type ReportData struct {
	KafkaCommon
	Type        MazeGame.BattleEventType `json:"report_type,omitempty"` // 事件类型
	EventFrame  uint64                   `json:"event_frame,omitempty"`
	EventTimeMs uint64                   `json:"event_time_ms,omitempty"` // 事件发生时间
	Data        string                   `json:"report_data,omitempty"`   // 上报数据
}

func NewReportData(t MazeGame.BattleEventType, eventFrame uint64, eventTimeMs uint64, data string) *ReportData {
	res := &ReportData{
		Type:        t,
		EventFrame:  eventFrame,
		EventTimeMs: eventTimeMs,
		Data:        data,
	}
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeUserLevelRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]

	return res
}
