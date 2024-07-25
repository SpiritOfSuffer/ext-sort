package pool

import "sync"

type WorkerPool struct {
	sem  chan struct{}
	wg   sync.WaitGroup
	done chan struct{}
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	return &WorkerPool{
		sem:  make(chan struct{}, maxWorkers),
		done: make(chan struct{}),
	}
}

func (wp *WorkerPool) Submit(task func()) {
	wp.sem <- struct{}{}
	wp.wg.Add(1)
	go func() {
		defer wp.wg.Done()
		defer func() { <-wp.sem }()
		task()
	}()
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.done)
}
