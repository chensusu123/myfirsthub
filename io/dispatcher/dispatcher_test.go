package dispatcher_test

import (
	"testing"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/io/dispatcher"
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
