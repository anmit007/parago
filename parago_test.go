package parago

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerCount(t *testing.T) {
	var maxWorkers uint32
	targetWorkers := 3
	input := make([]int, 10)

	Map(input, func(x int) (int, error) {

		atomic.AddUint32(&maxWorkers, 1)
		defer atomic.AddUint32(&maxWorkers, ^uint32(0))
		time.Sleep(10 * time.Millisecond)
		return x, nil
	}, WithWorkers(targetWorkers))

	if maxWorkers > uint32(targetWorkers) {
		t.Errorf("Expected max %d workers, got %d", targetWorkers, maxWorkers)
	}
}

func TestGoroutineCount(t *testing.T) {

	runtime.GOMAXPROCS(1)
	initialGoroutines := runtime.NumGoroutine()

	workers := 5
	input := make([]int, 10)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		Map(input, func(x int) (int, error) {
			time.Sleep(100 * time.Millisecond)
			return x, nil
		}, WithWorkers(workers))
	}()

	time.Sleep(50 * time.Millisecond)

	peakGoroutines := runtime.NumGoroutine()
	expectedGoroutines := initialGoroutines + workers + 1

	if peakGoroutines != expectedGoroutines {
		t.Errorf("Expected %d goroutines at peak, got %d", expectedGoroutines, peakGoroutines)
	}

	wg.Wait()

	time.Sleep(200 * time.Millisecond)
	finalGoroutines := runtime.NumGoroutine()

	if finalGoroutines != initialGoroutines {
		t.Errorf("Expected %d goroutines after cleanup, got %d", initialGoroutines, finalGoroutines)
	}
}
