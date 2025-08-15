package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts Tasks in n goroutines and stops its work when receiving M errors from Tasks.
func Run(tasks []Task, n, m int) error {
	ch := make(chan Task)

	var errCount int64

	var wg sync.WaitGroup

	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range ch {
				err := task()
				if err != nil {
					atomic.AddInt64(&errCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt64(&errCount) >= int64(m) {
			break
		}

		ch <- task
	}

	close(ch)
	wg.Wait()

	if atomic.LoadInt64(&errCount) >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
