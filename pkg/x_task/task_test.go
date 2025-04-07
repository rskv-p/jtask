package x_task

import (
	"os"
	"testing"
)

//
// ---------- Mock JSON Data ----------

// mockTaskJSON provides a valid example task configuration in JSON format.
const mockTaskJSON = `
{
	"name": "Test Collection",
	"description": "Some tasks",
	"tasks": [
		{
			"is_async": false,
			"is_sudo": false,
			"is_print_output": false,
			"name": "Echo Hello",
			"description": "Just echo hello",
			"exec": ["echo", "hello"]
		}
	]
}`

//
// ---------- Unit Tests ----------

// TestParseFileToStruct verifies that JSON is correctly parsed into a struct.
func TestParseFileToStruct(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "tasks_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(mockTaskJSON)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	var tasks TaskCollection
	result, err := ParseFileToStruct(tmpFile.Name(), &tasks)
	if err != nil {
		t.Errorf("ParseFileToStruct returned error: %v", err)
	}

	if result.Name != "Test Collection" {
		t.Errorf("unexpected collection name: %s", result.Name)
	}

	if len(result.Data) != 1 || result.Data[0].Name != "Echo Hello" {
		t.Errorf("unexpected task data: %+v", result.Data)
	}
}

// TestLoadTasks checks if LoadTasks successfully loads tasks from a JSON file.
func TestLoadTasks(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "tasks_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(mockTaskJSON)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	tasks, err := LoadTasks(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadTasks returned error: %v", err)
	}

	if tasks.Name != "Test Collection" || len(tasks.Data) != 1 {
		t.Errorf("unexpected result from LoadTasks: %+v", tasks)
	}
}

// TestExecuteTask validates that ExecuteTask runs a basic command and captures output.
func TestExecuteTask(t *testing.T) {
	task := &Task{
		IsAsync:       false,
		IsSudo:        false,
		IsPrintOutput: true,
		Name:          "Test Echo",
		Description:   "Echo test",
		Exec:          []string{"echo", "Hello, World!"},
	}

	result, err := ExecuteTask(task)
	if err != nil {
		t.Fatalf("ExecuteTask returned error: %v", err)
	}

	if result == nil || result.Output == "" {
		t.Errorf("unexpected empty result: %+v", result)
	}

	expected := "Hello, World!\n"
	if result.Output != expected {
		t.Errorf("expected output %q, got %q", expected, result.Output)
	}
}
