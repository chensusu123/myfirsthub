package reportdatamodel

import (
	"context"
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"maze_game_server/io/redis/reportredis"
	"maze_game_server/pb/common/MazeGame"
	"strings"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

const MazeUserLevelRecordTableName = "maze_cli_report_data"

type KafkaCommon = kafkacommonstruct.KafkaCommon

type ReportData struct {
	KafkaCommon
	Type        MazeGame.BattleEventType `json:"report_type,omitempty"`   // 事件类型
	UserID      uint64                   `json:"user_id,omitempty"`       // 用户id
	EnterTimeMs uint64                   `json:"enter_time_ms,omitempty"` // 进入关卡时间
	EventFrame  uint64                   `json:"event_frame,omitempty"`   // 上报帧
	EventTimeMs uint64                   `json:"event_time_ms,omitempty"` // 事件发生时间
	Data        string                   `json:"report_data,omitempty"`   // 上报数据
}

func NewReportData(ctx context.Context, t MazeGame.BattleEventType, userID uint64, eventFrame uint64, eventTimeMs uint64, data string) *ReportData {
	logger := fklog.ContextAppLogger(ctx)
	res := &ReportData{
		Type:        t,
		UserID:      userID,
		EventFrame:  eventFrame,
		EventTimeMs: eventTimeMs,
		Data:        data,
	}
	enterTime, err := reportredis.GetEnterTime(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "GetEnterTime Fail",
			zap.Uint64("userID", userID),
		)
		enterTime = 0
	}

	res.EnterTimeMs = enterTime
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeUserLevelRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]

	return res
}
