package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("if M <= 0, than return error", func(t *testing.T) {
		err := Run(make([]Task, 0, 10), 10, -1)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)

		err = Run(make([]Task, 0, 10), 10, 0)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 5
		tasks := make([]Task, 0, tasksCount)
		ttSyncCh := make(chan struct{})
		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				<-ttSyncCh
				return nil
			})
		}

		workersCount := tasksCount
		maxErrorsCount := 1

		wg := new(sync.WaitGroup)
		var err error
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = Run(tasks, workersCount, maxErrorsCount)
		}()

		require.Eventually(t, func() bool {
			return atomic.LoadInt32(&runTasksCount) == int32(tasksCount)
		}, 500*time.Millisecond, 10*time.Millisecond, "tasks were run sequentially?")

		close(ttSyncCh)
		wg.Wait()
		require.NoError(t, err)
	})
}
