package service

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestIsStateAvailable(t *testing.T) {
	// NOTE: force a race condition by commenting out the locks inside isStateAvailable()

	// test for unavailbility
	serviceStateLocker.currentServiceState = unavailable
	assert.Equal(t, unavailable, serviceStateLocker.currentServiceState)

	ok := isStateAvailable()
	assert.Equal(t, false, ok)

	// test for availability
	serviceStateLocker.currentServiceState = available
	assert.Equal(t, available, serviceStateLocker.currentServiceState)

	ok = isStateAvailable()
	assert.Equal(t, true, ok)

	// test race conditions
	const count = 20
	var wg sync.WaitGroup
	start := make(chan struct{}) // signal channel

	wg.Add(count) // #count go routines to wait for

	for i := 0; i < count; i++ {
		go func() {
			<-start // blocks code below, until channel is closed

			defer wg.Done()
			_ = isStateAvailable()
		}()
	}

	close(start) // starts executing blocked goroutines almost at the same time

	// test that read-lock inside isStateAvailable() blocks this write-lock
	serviceStateLocker.lock.Lock()
	serviceStateLocker.currentServiceState = available
	serviceStateLocker.lock.Unlock()

	wg.Wait() // wait until all goroutines finish executing
}
