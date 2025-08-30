package test

import (
	"maze_game_server/io/kafka/mazeuserlevelkafka"
	"maze_game_server/io/mysql/flowrecord"
	"testing"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

func BenchmarkUserLevelRecord(b *testing.B) {
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
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			flowrecord.SaveUserLevelRecord(gTestLogger, record)
		}
	})
}
