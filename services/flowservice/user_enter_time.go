package flowservice

func (s *service) SetUserEnterTime(userID uint64, enterTime uint64) {
	s.Lock()
	defer s.Unlock()
	s.reportTime[userID] = enterTime
}

func (s *service) GetUserEnterTime(userID uint64) uint64 {
	s.RLock()
	defer s.RUnlock()
	enterTime := s.reportTime[userID]
	return enterTime
}
