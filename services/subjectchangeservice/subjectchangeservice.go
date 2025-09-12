package subjectchangeservice

import (
	"context"
	"strconv"

	"maze_game_server/io/redis/usersubject"

	"gitlab.ifreetalk.com/maze-plate/freetk/common/commonconst"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/client/natsproduceroption"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/serverdepend"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/simpleclient/simplenatsproducer"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/nanotrace"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

var gNatsproducer database.NatsProducerI

func init() {
	natsService := "maze.usermsg.route.nats"
	natsproducer := simplenatsproducer.New(natsService, "SubjectChangeNotify")
	gNatsproducer = natsproducer
	serverdepend.RegisterDepend(natsproducer)
}

func SubjectChangeNotify(ctx context.Context, userID int64, broadcastID ...string) (err error) {
	span := nanotrace.NewSimpleTrace("SubjectChangeNotify")
	ctx = span.Start(ctx)
	defer span.Finish(ctx)
	cSpan := nanotrace.SpanFromContext(ctx)
	cSpan.SetAttributes(
		attribute.StringSlice("broadcast.id", broadcastID),
		attribute.Int64("user.id", userID),
		attribute.String("nats.subject", "maze.broadcast.subject.change.*"),
	)
	userIDStr := strconv.FormatInt(userID, 10)
	subject := "maze.broadcast.subject.change"
	err = gNatsproducer.Publish(ctx, subject, []byte(broadcastID[0]), natsproduceroption.WithSectionSubject(),
		natsproduceroption.WithTag(commonconst.NatsMsgHeaderPrimaryKey, userIDStr))
	fklog.ContextAppLogger(ctx).CtxInfo(ctx, "subjectchangeservice publish to nats",
		zap.String("subject", subject),
		zap.Strings("broadcastID", broadcastID),
		zap.Int64("userID", userID),
		zap.Error(err))
	return err
}

func Add(ctxP context.Context, userID int64, broadcastID ...string) (err error) {
	ctx := context.WithoutCancel(ctxP)
	span := nanotrace.NewSimpleTrace("SubjectAdd")
	ctx = span.Start(ctx)
	defer span.Finish(ctx)
	err = usersubject.Add(ctx, userID, broadcastID...)
	if err != nil {
		return err
	}
	err = SubjectChangeNotify(ctx, userID, broadcastID...)
	if err != nil {
		return err
	}
	return nil
}

func Del(ctxP context.Context, userID int64, broadcastID ...string) (err error) {
	ctx := context.WithoutCancel(ctxP)
	span := nanotrace.NewSimpleTrace("SubjectDel")
	ctx = span.Start(ctx)
	defer span.Finish(ctx)
	err = usersubject.Del(ctx, userID, broadcastID...)
	if err != nil {
		return err
	}
	err = SubjectChangeNotify(ctx, userID, broadcastID...)
	if err != nil {
		return err
	}
	return nil
}
