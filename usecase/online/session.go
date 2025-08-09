package online

import (
	"errors"
	"maze_game_server/lib/codec"
	"maze_game_server/lib/nano/session"
	"sync"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

var monitor = new(Monitor)

type Monitor struct {
	logger   fklog.FKLogI
	online   sync.Map
	sessions sync.Map
}

// OnCreate implements session.Monitor.
func (m *Monitor) OnCreate(s *session.Session) {
	m.sessions.Store(s.ID(), s)
	m.logger.InfoWF("Monitor OnCreate session created", zap.Int64("SessionID", s.ID()))
}

// OnClose implements session.Monitor.
func (m *Monitor) OnClose(s *session.Session, err error) {
	value, loaded := m.sessions.LoadAndDelete(s.ID())
	if loaded {
		if userID := value.(*session.Session).UID(); userID > 0 {
			m.online.Delete(uint64(userID))
		}
	}
	var lastErr string
	if err != nil {
		lastErr = err.Error()
	}
	m.logger.InfoWF("Monitor OnClose session closed", zap.Int64("SessionID", s.ID()), zap.Int64("UID", s.UID()), zap.String("lastErr", lastErr))
}

// SessionMonitor returns a session monitor.
func SessionMonitor(logger fklog.FKLogI) *Monitor {
	monitor.logger = logger
	return monitor
}

// Bind
func Bind(logger fklog.FKLogI, s *session.Session, userID uint64) (err error) {
	value, found := monitor.sessions.Load(s.ID())
	if !found {
		logger.ErrorWF("Bind session not found", zap.Error(ErrSessionNotFound), zap.Int64("ID", s.ID()), zap.Uint64("userID", userID))
		return ErrSessionNotFound
	}
	monitor.online.Store(userID, value)
	return
}

// Push
func Push(logger fklog.FKLogI, userID uint64, packetType uint16, v interface{}) (err error) {
	s, found := monitor.online.Load(userID)
	if !found {
		logger.ErrorWF("Push session not found", zap.Error(ErrSessionNotFound), zap.Uint64("userID", userID), zap.Any("v", v))
		return ErrSessionNotFound
	}
	return s.(*session.Session).ResponseMID(codec.ToMessageID(uint32(time.Now().Unix()), 0, packetType), v)
}

// Scan
func Scan(fn func(id int64, s *session.Session)) {
	monitor.sessions.Range(func(key, value interface{}) bool {
		fn(key.(int64), value.(*session.Session))
		return true
	})
}

func IsOnline(userID uint64) bool {
	_, ok := monitor.online.Load(userID)
	return ok
}
