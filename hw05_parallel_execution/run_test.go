package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
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

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestRunNonPositiveM(t *testing.T) {
	defer goleak.VerifyNone(t)

	makeTasks := func(outTasksRun *int64) []Task {
		count := 50

		var tasks []Task
		for range count {
			tasks = append(tasks, func() error {
				defer atomic.AddInt64(outTasksRun, 1)

				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				if rand.Intn(3) == 1 {
					return errors.New("error")
				}

				return nil
			})
		}

		return tasks
	}

	type TestCase struct {
		Name        string
		M           int
		N           int
		Tasks       func(out *int64) []Task
		ExpectedRun int64
		ExpectedErr error
	}

	testCases := []TestCase{
		{
			Name:        "zero",
			M:           0,
			N:           13,
			Tasks:       makeTasks,
			ExpectedRun: 0,
			ExpectedErr: ErrErrorsLimitExceeded,
		},
		{
			Name:        "negative",
			M:           -1,
			N:           13,
			Tasks:       makeTasks,
			ExpectedRun: 0,
			ExpectedErr: ErrErrorsLimitExceeded,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var actualRun int64
			err := Run(tc.Tasks(&actualRun), tc.N, tc.M)

			require.ErrorIs(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedRun, actualRun)
		})
	}
}
