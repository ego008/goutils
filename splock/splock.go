package splock

import "sync"

// -- simpleLock
type SL struct {
	key  string
	lock bool
}

type SimpleLock struct {
	mp sync.Map
}

func (s *SimpleLock) Init(k string) *SL {
	if v, ok := s.mp.Load(k); ok {
		return v.(*SL)
	}

	sl := SL{key: k}
	s.mp.Store(k, &sl)
	return &sl
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
