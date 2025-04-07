package x_task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/rskv-p/jtask/pkg/x_log"
	"github.com/rskv-p/jtask/pkg/x_util"
)

//
// ---------- Data Structures ----------

// TaskCollection represents a group of tasks loaded from a config file.
type TaskCollection struct {
	Name        string  `json:"name"`        // Collection name
	Description string  `json:"description"` // Description of the task collection
	Data        []*Task `json:"tasks"`       // List of tasks
}

// Task represents an individual task with execution settings.
type Task struct {
	IsAsync       bool     `json:"is_async"`        // Run in parallel
	IsSudo        bool     `json:"is_sudo"`         // Run with sudo
	IsPrintOutput bool     `json:"is_print_output"` // Capture and store output
	Name          string   `json:"name"`            // Task name
	Description   string   `json:"description"`     // Task description
	Exec          []string `json:"exec"`            // Command to execute
}

// Result contains the result of a task execution.
type Result struct {
	ID          string `json:"id"`          // Unique task ID
	Name        string `json:"name"`        // Task name
	Description string `json:"description"` // Task description
	Output      string `json:"output"`      // Captured output
}

//
// ---------- Public Functions ----------

// LoadTasks reads and parses a task definition file.
func LoadTasks(path string) (*TaskCollection, error) {
	if path == "" {
		x_log.Error().Msg("no config file path provided")
		return nil, fmt.Errorf("invalid format of tasks file")
	}

	x_log.Info().
		Str("path", path).
		Msg("loading tasks from file")

	tasks := &TaskCollection{}
	_, err := ParseFileToStruct(path, tasks)
	if err != nil {
		// Log the failure to load tasks
		x_log.Error().
			Err(err).
			Str("path", path).
			Msg("failed to load tasks")
		return nil, err
	}

	// Log the success of loading tasks
	x_log.Info().
		Int("count", len(tasks.Data)).
		Str("path", path).
		Msg("tasks loaded successfully")

	return tasks, nil
}

// ExecuteTask runs a single task and returns a result.
func ExecuteTask(t *Task) (*Result, error) {
	id, _ := x_util.RandomString(8)

	result := &Result{
		ID:          fmt.Sprintf("(async) %s", id),
		Name:        t.Name,
		Description: t.Description,
	}

	// Log the start of task execution
	x_log.Info().
		Str("task", t.Name).
		Msg("starting task execution")

	// Log task details before execution
	x_log.Debug().
		Str("task", t.Name).
		Bool("sudo", t.IsSudo).
		Bool("async", t.IsAsync).
		Bool("print_output", t.IsPrintOutput).
		Interface("exec", t.Exec).
		Msg("task execution details")

	if len(t.Exec) == 0 {
		err := fmt.Errorf("task %s has empty exec command", t.Name)
		x_log.Error().
			Err(err).
			Str("task", t.Name).
			Msg("cannot execute task")
		result.Output = err.Error()
		return result, err
	}

	// Run the task command
	var stdOut bytes.Buffer
	var cmd *exec.Cmd
	if t.IsSudo {
		cmd = exec.Command("sudo", t.Exec...)
	} else {
		cmd = exec.Command(t.Exec[0], t.Exec[1:]...)
	}

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdOut

	// Run the command and check for errors
	if err := cmd.Run(); err != nil {
		// Log failure of task execution
		x_log.Error().
			Err(err).
			Str("task", t.Name).
			Msg("task execution failed")
		result.Output = stdOut.String()
		return result, err
	}

	// Capture and log task output if configured
	if t.IsPrintOutput {
		result.Output = stdOut.String()
		x_log.Debug().
			Str("task", t.Name).
			Int("output_len", len(result.Output)).
			Msg("task output captured")
	}

	// Log task completion
	x_log.Info().
		Str("task", t.Name).
		Msg("task completed successfully")

	return result, nil
}

//
// ---------- Helper Functions ----------

// ParseFileToStruct reads JSON from file and unmarshals into a struct.
func ParseFileToStruct[T any](path string, model *T) (*T, error) {
	x_log.Debug().
		Str("path", path).
		Msg("reading file")

	content, err := os.ReadFile(path)
	if err != nil {
		// Log file read error
		x_log.Error().
			Err(err).
			Str("path", path).
			Msg("failed to read file")
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if err := json.Unmarshal(content, model); err != nil {
		// Log JSON unmarshalling error
		x_log.Error().
			Err(err).
			Str("path", path).
			Msg("invalid JSON")
		return nil, fmt.Errorf("error unmarshalling JSON data: %w", err)
	}

	// Log successful parsing
	x_log.Debug().
		Str("path", path).
		Msg("JSON parsed successfully")

	return model, nil
}
