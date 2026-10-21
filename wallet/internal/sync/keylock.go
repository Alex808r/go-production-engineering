package sync

import (
	"hash/fnv"
	"sync"

	"github.com/google/uuid"
)

const defaultStripeCount = 4096

// KeyLocker provides striped locking by UUID to eliminate lock contention
// across different wallets while strictly serializing concurrent operations on the same wallet.
// Zero heap allocations per Lock/Unlock call.
type KeyLocker struct {
	stripes []sync.Mutex
	mask    uint32
}

// NewKeyLocker creates a new KeyLocker with the specified number of stripes (rounded to power of 2).
func NewKeyLocker(stripes ...int) *KeyLocker {
	count := defaultStripeCount
	if len(stripes) > 0 && stripes[0] > 0 {
		count = stripes[0]
	}

	// Round up to nearest power of 2 for bitwise AND masking
	power := 1
	for power < count {
		power <<= 1
	}

	return &KeyLocker{
		stripes: make([]sync.Mutex, power),
		mask:    uint32(power - 1),
	}
}

// Lock acquires the lock associated with the given UUID.
func (kl *KeyLocker) Lock(id uuid.UUID) {
	idx := kl.hash(id) & kl.mask
	kl.stripes[idx].Lock()
}

// Unlock releases the lock associated with the given UUID.
func (kl *KeyLocker) Unlock(id uuid.UUID) {
	idx := kl.hash(id) & kl.mask
	kl.stripes[idx].Unlock()
}

func (kl *KeyLocker) hash(id uuid.UUID) uint32 {
	h := fnv.New32a()
	_, _ = h.Write(id[:])
	return h.Sum32()
}
