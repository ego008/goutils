package splock

import "sync"

// -- simpleLock
var mutex = &sync.RWMutex{}

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
	mutex.Lock()
	defer mutex.Unlock()
	if v, ok := s.mp[k]; ok {
		return v
	}

	sl := SL{key: k}
	s.mp[k] = &sl
	return &sl
}

func (s *SimpleLock) HasLocked() (bool, string) {
	if s.mp == nil {
		return false, ""
	}
	mutex.RLock()
	defer mutex.RUnlock()
	for k := range s.mp {
		if s.mp[k].lock {
			return true, k
		}
	}
	return false, ""
}

func (l *SL) IsLocked() bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return l.lock
}

func (l *SL) Lock() {
	mutex.Lock()
	l.lock = true
	mutex.Unlock()
}

func (l *SL) UnLock() {
	mutex.Lock()
	l.lock = false
	mutex.Unlock()
}
