package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg := new(sync.WaitGroup)
	var errCnt int32
	tCh := make(chan Task)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range tCh {
				if err := t(); err != nil {
					atomic.AddInt32(&errCnt, 1)
				}
			}
		}()
	}

	for _, t := range tasks {
		if atomic.LoadInt32(&errCnt) >= int32(m) {
			break
		}
		tCh <- t
	}
	close(tCh)

	wg.Wait()
	if errCnt >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func RunV2(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	wg := new(sync.WaitGroup)
	var errCnt int32
	var tCnt int32

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				tIdx := atomic.AddInt32(&tCnt, 1) - 1
				if tIdx >= int32(len(tasks)) || atomic.LoadInt32(&errCnt) >= int32(m) {
					return
				}
				if err := tasks[tIdx](); err != nil {
					atomic.AddInt32(&errCnt, 1)
				}
			}
		}()
	}

	wg.Wait()
	if errCnt >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
