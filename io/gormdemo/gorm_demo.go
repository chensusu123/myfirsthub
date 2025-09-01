// gorm demo
package gormdemo

///Users/majiange/data/dev/go_work/maze/maze-plate/freetk/fkcore/database/nanogorm
import (
	"context"
	"maze_game_server/io/kafka/mazebarrieruserkafka"
	"maze_game_server/io/mysql"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database/nanogorm"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type GormDemo struct {
	*nanogorm.NanoGorm
}

func NewGormDemo(serviceName string, name string) *GormDemo {
	rt := &GormDemo{}
	rt.NanoGorm = nanogorm.NewNanoGorm(serviceName, name)
	return rt
}

func (g *GormDemo) SaveBarrierUserRecord(ctx context.Context, record *mazebarrieruserkafka.MazeBarrierUserGameRecord) error {
	logger := fklog.ContextAppLogger(ctx)
	db, err := g.GetGormDB()
	if err != nil {
		logger.CtxError(ctx, "GetMysqlDb fail", zap.Error(err), zap.Any("record:", record))
		return err
	}

	res := db.Table(mysql.GetFullyQualifiedTableName("test_db.test_table")).Create(record)

	return res.Error
}
