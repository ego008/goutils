package splock

// -- simpleLock

type SL struct {
	key  string
	lock bool
}

type SimpleLock struct {
	mp map[string]*SL
}

func (s *SimpleLock) Init(k string) *SL {
	if s.mp == nil {
		s.mp = map[string]*SL{}
	}
	if v, ok := s.mp[k]; ok {
		return v
	}

	sl := SL{key: k}
	s.mp[k] = &sl
	return &sl
}

func (l *SL) IsLocked() bool {
	return l.lock
}

func (l *SL) Lock() {
	l.lock = true
}

func (l *SL) UnLock() {
	l.lock = false
}
