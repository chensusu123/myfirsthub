package test

import (
	"database/sql"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/structsdef"
	"maze_game_server/io/kafka/dollequipassmeblekakfa"
	"maze_game_server/io/kafka/dollequipdismantlekafka"
	"maze_game_server/io/kafka/dollmazefoekafka"
	"maze_game_server/io/kafka/equipposstrengrecordkafka"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/kafka/mazecollectrecord"
	"maze_game_server/io/kafka/mazeequipbagrecord"
	"maze_game_server/io/kafka/mazemoneykafka"
	"maze_game_server/io/kafka/mazerebornkafka"
	"maze_game_server/io/kafka/mazetempbuffchgmsg"
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/mysql"
	"maze_game_server/io/mysql/flowrecord"
	"testing"
	"time"
)

var gTestLogger fklog.FKLogI

func init() {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	mysql.InitMysql()
}
func TestUserLevelRecord(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	record := &mazeuserlevelkafka.MazeUserLevelRecord{
		UserId:      1010101010,
		OldLevel:    1,
		OldTotalExp: 100,
		NewLevel:    2,
		NewTotalExp: 200,
		GroupID:     1,
		CreateTime:  time.Now().UnixMilli(),
	}
	flowrecord.SaveUserLevelRecord(gTestLogger, record)
}

func TestRebornRecord(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	record := &mazerebornkafka.MazeRebornRecord{
		UserId:      1010101010,
		Barrier:     1,
		RebornCount: 4,
		RebornCost:  "dsafdsafdsaf",
		GroupID:     1,
		CreateTime:  time.Now().UnixMilli(),
	}
	flowrecord.SaveRebornRecord(gTestLogger, record)
}

func TestMoneyRecord(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	record := &mazemoneykafka.MazeMoneyRecord{
		UserId:        1010101010,
		OldMoneyId:    1,
		OldMoneyCount: 1000,
		NewMoneyId:    1,
		NewMoneyCount: 2000,
		TradeNo:       123456789,
		ChgReason:     1,
		GroupID:       1,
		CreateTime:    time.Now().UnixMilli(),
	}
	flowrecord.SaveMoneyRecord(gTestLogger, record)
}

func TestEquipBagRecord(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	record := &mazeequipbagrecord.MazeGameEquipBagRecord{
		UserId:        1010101010,
		ChgType:       1,
		TradeNum:      123456789,
		AddEquipGuids: "1;2;3;4",
		DelEquipGuids: "5;6;7;8",
		OpType:        1,
		IsFail:        0,
		GroupID:       1,
		CreateTime:    time.Now().UnixMilli(),
	}
	flowrecord.SaveEquipBagRecord(gTestLogger, record)
}

func TestBarrierUserRecord(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	record := &mazebarrieruserkafka.MazeBarrierUserGameRecord{
		UserId:     1010101010,
		Barrier:    1,
		GameRet:    1,
		Awards:     "1;2;3;4",
		GroupID:    1,
		CreateTime: time.Now().UnixMilli(),
	}
	flowrecord.SaveBarrierUserRecord(gTestLogger, record)
}

func TestAttrChgRecord(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	record := &structsdef.MazeGameAttrChgRecord{
		UserId:     1010101010,
		GroupId:    1,
		AttrId:     1,
		AttrType:   1,
		NewVal:     20,
		OldVal:     10,
		ChgType:    1,
		ChgSubType: 2,
		ChgDesc:    "llll",
		CreateTime: time.Now().UnixMilli(),
		Extra:      "xxxxxx",
	}
	flowrecord.SaveAttrChgRecord(gTestLogger, record)
}

func TestEquipAssembleRecord(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	record := &dollequipassmeblekakfa.MazeGameEquipAssembleRecord{
		UserId:     1010101010,
		GroupId:    1,
		EquipPos:   1,
		OpType:     1,
		NewEquipId: 12,
		NewGuid:    111,
		OldEquipId: 11,
		OldGuid:    110,
		OldFElem:   "123",
		NewFElem:   "123",
		RetCode:    0,
		CodeMask:   0,
		TransID:    123456789,
		OpTime:     time.Now().UnixMilli(),
	}
	flowrecord.SaveEquipAssembleRecord(gTestLogger, record)
}

func TestEquipDismantRecord(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	record := &dollequipdismantlekafka.MazeGameEquipDismantleRecord{
		UserId:     1010101010,
		EquipGuids: "1;2;3;4",
		TradeNum:   123456789,
		Award:      "1,2,3,4",
		OpType:     1,
		IsFail:     0,
		GroupID:    1,
		CreateTime: time.Now().UnixMilli(),
	}
	flowrecord.SaveEquipDismantRecord(gTestLogger, record)
}

func TestEquipPosStrengRecord(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	record := &equipposstrengrecordkafka.EquipPosLevelUpRecord{
		UserId:       1010101010,
		GroupId:      1,
		OpTime:       time.Now().UnixMilli(),
		PosId:        1,
		OldPosLv:     1,
		NewPosLv:     1,
		OldPosSuitId: 2,
		NewPosSuitId: 2,
		TradeNo:      2,
		CostItems:    "1,2,3",
		Result:       0,
	}
	flowrecord.SaveEquipPosStrengRecord(gTestLogger, record)
}

func TestTempBuffChg(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	attrs := make([]*mazetempbuffchgmsg.AttrChgInfo, 0)
	attrs = append(attrs, &mazetempbuffchgmsg.AttrChgInfo{
		AttrId: 1,
		OldVal: 1,
		CurVal: 2,
	})
	record := &mazetempbuffchgmsg.MazeTempBuffChangeMsg{
		UserId:     1010101010,
		GroupId:    1,
		StageId:    1,
		ChgAttrs:   attrs,
		ChgType:    1,
		ChgDesc:    "测试测试",
		CreateTime: time.Now().UnixMilli(),
	}
	flowrecord.SaveTempBuffChgRecord(gTestLogger, record)
}

func TestFoeRecord(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	record := &dollmazefoekafka.DollMazeFoeRecord{
		UserId:     1010101010,
		Barrier:    1,
		Area:       1,
		Level:      1,
		AwardList:  "1,2,3",
		GroupID:    1,
		CreateTime: time.Now().UnixMilli(),
	}
	flowrecord.SaveFoeRecord(gTestLogger, record)
}

func TestMazeCollectChgRecord(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	record := &mazecollectrecord.MazeCollectChgRecord{
		UserId:        1010101010,
		OpType:        1,
		StartTime:     1,
		LastTime:      1,
		NewLastTime:   1,
		AvailableTime: 1,
		EndTime:       1,
		PeriodTime:    1,
		CollectTimes:  1,
		BarrierId:     1,
		TradeNo:       1,
		AddItems:      "1,2,3",
		RemainItems:   "4,5,6",
		RetCode:       1,
		GroupID:       1,
		CreateTime:    time.Now().UnixMilli(),
	}
	flowrecord.SaveCollectChgRecord(gTestLogger, record)
}

func TestSweepRecord(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")
	record := &mazebarrieruserkafka.MazeBarrierUserGameRecord{
		UserId:     1010101010,
		Barrier:    1,
		GameRet:    3,
		Awards:     "1,2,3",
		GroupID:    1,
		CreateTime: time.Now().UnixMilli(),
	}
	flowrecord.SaveSweepRecord(gTestLogger, record)
}

func TestSharding(t *testing.T) {
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")

	db, err := mysql.GetMysqlDb("db_maze_user_level_record_202506")
	if err != nil {
		t.Errorf("mysql.GetMysqlDb err: %v", err)
	}
	InsertTest(db, "t_maze_user_level_record_3")

	db, err = mysql.GetMysqlDb("db_maze_user_level_record_202507")
	if err != nil {
		t.Errorf("mysql.GetMysqlDb err: %v", err)
	}
	InsertTest(db, "t_maze_user_level_record_3")
}

func InsertTest(db *sql.DB, tableName string) {
	sqlStr := fmt.Sprintf("INSERT INTO %s (`user_id`, `old_level`, `old_total_exp`, `new_level`, `new_total_exp`, `group_id`, `create_time`, `server_id`) VALUES (?,?,?,?,?,?,?,?)",
		tableName,
	)
	record := &mazeuserlevelkafka.MazeUserLevelRecord{
		UserId:      1010101010,
		OldLevel:    1,
		OldTotalExp: 100,
		NewLevel:    2,
		NewTotalExp: 200,
		GroupID:     1,
		CreateTime:  time.Now().UnixMilli(),
	}
	_, err := db.Exec(sqlStr,
		record.UserId, record.OldLevel, record.OldTotalExp, record.NewLevel, record.NewTotalExp, record.GroupID, record.CreateTime, record.ServerId)
	if err != nil {
		gTestLogger.ErrorWF("SaveUserLevelRecord fail", zap.Error(err), zap.Any("flowrecord", record))
		return
	}
	gTestLogger.InfoWF("SaveUserLevelRecord succ", zap.Any("flowrecord", record))
}
