package dispatcher_test

import (
	"context"
	"testing"
	"time"

	"maze_game_server/io/dispatcher"
)

var (
	d = dispatcher.NewDispatcher[string]()
)

func TestDispatcher(t *testing.T) {
	d.Watch(func(ctx context.Context, msg string) {
		t.Logf("watcher #1: %s", msg)
	})
	d.Watch(func(ctx context.Context, msg string) {
		t.Logf("watcher #2: %s", msg)
	})
	ctx := context.Background()
	d.Push(ctx, time.Now().String())

	time.Sleep(time.Second)
}
