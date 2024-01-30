package threadsafe

import (
	"sync"
)

// Locker is an interface that implements the Lock and Unlock methods.
type Locker sync.Locker

// RLocker is an interface that implements the RLock and RUnlock methods.
type RLocker interface {
	RLock()
	RUnlock()
}
