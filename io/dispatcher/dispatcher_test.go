package dispatcher_test

import (
	"testing"
	"time"

	"gitlab.ifreetalk.com/maze/maze_game_server/io/dispatcher"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
)

var (
	d = dispatcher.NewDispatcher[string]()
)

func TestDispatcher(t *testing.T) {
	d.Watch(func(logger fklog.FKLogI, msg string) {
		t.Logf("watcher #1: %s", msg)
	})
	d.Watch(func(logger fklog.FKLogI, msg string) {
		t.Logf("watcher #2: %s", msg)
	})

	d.Push(gTestLogger, time.Now().String())

	time.Sleep(time.Second)
}
