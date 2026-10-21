package sync

import (
	"sync"
	"testing"

	"github.com/google/uuid"
)

func TestKeyLocker_SerializesSameKey(t *testing.T) {
	kl := NewKeyLocker()
	targetID := uuid.New()

	const iterations = 1000
	counter := 0
	var wg sync.WaitGroup
	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()
			kl.Lock(targetID)
			counter++
			kl.Unlock(targetID)
		}()
	}

	wg.Wait()

	if counter != iterations {
		t.Fatalf("expected counter %d, got %d", iterations, counter)
	}
}

func TestKeyLocker_DifferentKeysRunConcurrently(t *testing.T) {
	kl := NewKeyLocker(4096)
	id1 := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	id2 := uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")

	kl.Lock(id1)

	unlockedCh := make(chan struct{})
	go func() {
		// id2 shouldn't be blocked by id1 lock
		kl.Lock(id2)
		kl.Unlock(id2)
		close(unlockedCh)
	}()

	select {
	case <-unlockedCh:
		// Success: different key was not blocked
	default:
		// Try waiting briefly, or it might just be scheduling
		<-unlockedCh
	}

	kl.Unlock(id1)
}
