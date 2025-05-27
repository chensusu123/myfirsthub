package online

import (
	"errors"
	"sync"
	"time"

	"maze_game_server/lib/codec"

	"github.com/lonng/nano/session"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

var monitor = new(Monitor)

type Monitor struct {
	online   sync.Map
	sessions sync.Map
}

// OnCreate implements session.Monitor.
func (m *Monitor) OnCreate(s *session.Session) {
	m.sessions.Store(s.ID(), s)
}

// OnClose implements session.Monitor.
func (m *Monitor) OnClose(s *session.Session) {
	value, loaded := m.sessions.LoadAndDelete(s.ID())
	if loaded {
		if userID := value.(*session.Session).UID(); userID > 0 {
			m.online.Delete(uint64(userID))
		}
	}
}

// SessionMonitor returns a session monitor.
func SessionMonitor() *Monitor {
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
