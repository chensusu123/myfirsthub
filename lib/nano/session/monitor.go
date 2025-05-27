package session

type Monitor interface {

	// OnCreate
	OnCreate(s *Session)

	// OnClose
	OnClose(s *Session)
}
