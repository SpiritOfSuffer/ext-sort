package pool_test

import (
	"sync"
	"testing"
	"time"

	"ext-sort/pkg/pool"
)

func Test_WorkerPool(t *testing.T) {
	var wg sync.WaitGroup

	wp := pool.NewWorkerPool(10)
	numTasks := 100

	wg.Add(numTasks)

	for i := 0; i < numTasks; i++ {
		wp.Submit(func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
		})
	}

	wp.Wait()
	wg.Wait()
}

func Test_WorkerPool_Limit(t *testing.T) {
	var mu sync.Mutex
	var wg sync.WaitGroup

	maxWorkers := 5
	numTasks := 20

	wp := pool.NewWorkerPool(maxWorkers)
	wg.Add(numTasks)

	activeTasks := 0

	for i := 0; i < numTasks; i++ {
		wp.Submit(func() {
			defer wg.Done()
			mu.Lock()
			activeTasks++
			if activeTasks > maxWorkers {
				t.Errorf("More than %d tasks running concurrently", maxWorkers)
			}
			mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			mu.Lock()
			activeTasks--
			mu.Unlock()
		})
	}

	wp.Wait()
	wg.Wait()
}

func Test_WorkerPool_Wait(t *testing.T) {
	var wg sync.WaitGroup

	wp := pool.NewWorkerPool(10)
	numTasks := 10

	wg.Add(numTasks)

	done := make(chan struct{})
	go func() {
		wp.Wait()
		close(done)
	}()

	for i := 0; i < numTasks; i++ {
		wp.Submit(func() {
			defer wg.Done()
			time.Sleep(10 * time.Millisecond)
		})
	}

	wg.Wait()
	select {
	case <-done:
	default:
		t.Error("Channel done not closed")
	}
}
