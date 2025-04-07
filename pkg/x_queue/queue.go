package x_queue

import (
	"github.com/rskv-p/jtask/pkg/x_log"
	"github.com/rskv-p/jtask/pkg/x_task"
)

//
// ---------- Task Queues ----------

// TaskQueues holds tasks separated into async and sequential execution queues.
type TaskQueues struct {
	Async      []*x_task.Task `json:"async"`      // Tasks to run concurrently
	Sequential []*x_task.Task `json:"sequential"` // Tasks to run one-by-one
}

//
// ---------- Public API ----------

// CreateTaskQueues separates tasks into async and sequential queues.
func CreateTaskQueues(tasks *x_task.TaskCollection) (*TaskQueues, error) {
	x_log.Info().
		Int("total", len(tasks.Data)).
		Msg("creating task queues")

	var (
		async = make([]*x_task.Task, 0)
		seq   = make([]*x_task.Task, 0)
	)

	// Separate tasks into async and sequential queues
	for _, t := range tasks.Data {
		if t.IsAsync {
			async = append(async, t)
			// Log the task added to async queue
			x_log.Debug().
				Str("task", t.Name).
				Msg("added to async queue")
		} else {
			seq = append(seq, t)
			// Log the task added to sequential queue
			x_log.Debug().
				Str("task", t.Name).
				Msg("added to sequential queue")
		}
	}

	// Log the count of tasks in each queue
	x_log.Info().
		Int("async", len(async)).
		Int("sequential", len(seq)).
		Msg("task queues created")

	// Log the names of the tasks in each queue for better tracking
	x_log.Debug().
		Str("async_tasks", getTaskNames(async)).
		Str("sequential_tasks", getTaskNames(seq)).
		Msg("tasks grouped into async and sequential queues")

	return &TaskQueues{
		Async:      async,
		Sequential: seq,
	}, nil
}

// ---------- Helper Function ----------

// getTaskNames returns a comma-separated string of task names.
func getTaskNames(tasks []*x_task.Task) string {
	names := []string{}
	for _, t := range tasks {
		names = append(names, t.Name)
	}
	return stringJoin(names, ", ")
}

// stringJoin is a helper to join strings with a separator.
func stringJoin(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
