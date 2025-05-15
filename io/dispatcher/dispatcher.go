package dispatcher

import (
	"context"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

type msgWrapper[T any] struct {
	logger  fklog.FKLogI
	message T
}

type DispatcherOption[T any] func(*Dispatcher[T])

type Dispatcher[T any] struct {
	ch          chan msgWrapper[T]
	subscribers []func(fklog.FKLogI, T)
	cancelCtx   context.CancelFunc
}

// NewDispatcher
func NewDispatcher[T any](opts ...DispatcherOption[T]) (d *Dispatcher[T]) {
	d = &Dispatcher[T]{
		ch:          make(chan msgWrapper[T], 10),
		subscribers: make([]func(fklog.FKLogI, T), 0),
	}
	// Set options
	for _, setOpt := range opts {
		setOpt(d)
	}
	var ctx context.Context
	ctx, d.cancelCtx = context.WithCancel(context.Background())
	go d.background(ctx)
	return
}

func (d *Dispatcher[T]) background(ctx context.Context) {
	for {
		select {
		// Cancel context
		case <-ctx.Done():
			return
		case wrapper, ok := <-d.ch:
			if !ok {
				return
			}
			// 暂时先不使用WorkGroup
			for _, fn := range d.subscribers {
				go fn(wrapper.logger, wrapper.message)
			}
		}
	}
}

func (d *Dispatcher[T]) Push(logger fklog.FKLogI, msg T) {
	d.ch <- msgWrapper[T]{logger: logger, message: msg}
}

// Watch
func (d *Dispatcher[T]) Watch(fn func(logger fklog.FKLogI, msg T)) {
	d.subscribers = append(d.subscribers, fn)
}

// Close closes the dispatcher.
func (d *Dispatcher[T]) Close() {
	d.cancelCtx()
}
