package x_queue

import (
	"testing"

	"github.com/rskv-p/jtask/pkg/x_task"
)

//
// ---------- Unit Test: CreateTaskQueues ----------

// TestCreateTaskQueues verifies that tasks are correctly split into async and sequential queues.
func TestCreateTaskQueues(t *testing.T) {
	// Setup test tasks
	task1 := &x_task.Task{Name: "Async 1", IsAsync: true}
	task2 := &x_task.Task{Name: "Seq 1", IsAsync: false}
	task3 := &x_task.Task{Name: "Async 2", IsAsync: true}
	task4 := &x_task.Task{Name: "Seq 2", IsAsync: false}

	collection := &x_task.TaskCollection{
		Data: []*x_task.Task{task1, task2, task3, task4},
	}

	// Run queue creator
	queues, err := CreateTaskQueues(collection)
	if err != nil {
		t.Fatalf("CreateTaskQueues returned error: %v", err)
	}

	// Assert async queue length
	if len(queues.Async) != 2 {
		t.Errorf("expected 2 async tasks, got %d", len(queues.Async))
	}

	// Assert sequential queue length
	if len(queues.Sequential) != 2 {
		t.Errorf("expected 2 sequential tasks, got %d", len(queues.Sequential))
	}

	// Assert async queue order and content
	if queues.Async[0] != task1 || queues.Async[1] != task3 {
		t.Error("unexpected async task order or content")
	}

	// Assert sequential queue order and content
	if queues.Sequential[0] != task2 || queues.Sequential[1] != task4 {
		t.Error("unexpected sequential task order or content")
	}
}
