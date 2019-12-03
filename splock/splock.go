package splock

// -- simpleLock
type SL struct {
	lock bool
}

type SimpleLock struct {
	mp map[string]SL
}

func (s *SimpleLock) Init(k string) SL {
	if s.mp == nil {
		s.mp = map[string]SL{}
	}
	if v, ok := s.mp[k]; ok {
		return v
	}
	s.mp[k] = SL{}
	return s.mp[k]
}

func (sl *SL) IsLocked() bool {
	return sl.lock
}

func (sl *SL) Lock() {
	sl.lock = true
}

func (sl *SL) UnLock() {
	sl.lock = false
}
