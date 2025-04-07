package cmd

import (
	"fmt"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/rskv-p/jtask/pkg/x_log"
	"github.com/rskv-p/jtask/pkg/x_task"
	"github.com/spf13/cobra"
)

//
// ---------- Command Definition ----------

// runCmd allows selecting and running a task.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Select a task to run",
	Long:  "Select a task from the list and execute it.",
	Run: func(cmd *cobra.Command, args []string) {
		// Load tasks from the config file
		x_log.Info().
			Str("path", pathFlag).
			Msg("loading tasks from config file")

		tasks, err := x_task.LoadTasks(pathFlag)
		if err != nil {
			// Log error if task loading fails
			x_log.Error().
				Err(err).
				Str("path", pathFlag).
				Msg("failed to load tasks")
			fmt.Println("Error loading tasks:", err)
			return
		}

		// Handle case where no tasks are available
		if len(tasks.Data) == 0 {
			x_log.Warn().
				Msg("no tasks available to select")
			fmt.Println("No tasks available to select.")
			return
		}

		// Build selection options
		x_log.Info().
			Int("count", len(tasks.Data)).
			Msg("building selection options")

		var options []huh.Option[string]
		for _, task := range tasks.Data {
			options = append(options, huh.NewOption(task.Name, task.Name))
		}

		// Prompt user to select a task
		x_log.Debug().Msg("prompting user to select a task")
		var selectedTask string
		if err := huh.NewSelect[string]().
			Title("Select a task to run:").
			Options(options...).
			Value(&selectedTask).
			Run(); err != nil {
			// Log error if task selection fails
			x_log.Error().
				Err(err).
				Msg("task selection aborted")
			fmt.Println("Error selecting task:", err)
			return
		}

		// Log user selection
		x_log.Info().
			Str("task", selectedTask).
			Msg("user selected task")
		fmt.Printf("You selected: %s\n", selectedTask)

		// Find the selected task
		var selected *x_task.Task
		for _, task := range tasks.Data {
			if task.Name == selectedTask {
				selected = task
				break
			}
		}

		// Run the selected task if found
		if selected != nil {
			x_log.Info().
				Str("task", selected.Name).
				Msg("executing selected task")

			if err := executeTask(selected); err != nil {
				// Log failure if task execution fails
				x_log.Error().
					Err(err).
					Str("task", selected.Name).
					Msg("task execution failed")
				fmt.Printf("Error executing task %s: %v\n", selected.Name, err)
			} else {
				// Log success if task executes correctly
				x_log.Info().
					Str("task", selected.Name).
					Msg("task executed successfully")
				fmt.Printf("Task %s executed successfully!\n", selected.Name)
			}
		} else {
			// Log and notify if selected task was not found
			x_log.Warn().
				Str("task", selectedTask).
				Msg("selected task not found")
			fmt.Println("Selected task not found.")
		}
	},
}

// ---------- Command Initialization ----------
func init() {
	// Register 'run' command to the root command
	rootCmd.AddCommand(runCmd)
}

// ---------- Task Execution ----------
func executeTask(task *x_task.Task) error {
	// Log task execution details
	x_log.Debug().
		Str("task", task.Name).
		Interface("exec", task.Exec).
		Msg("running task command")

	var cmd *exec.Cmd
	if task.IsSudo {
		// Use sudo if needed
		cmd = exec.Command("sudo", task.Exec...)
	} else {
		cmd = exec.Command(task.Exec[0], task.Exec[1:]...)
	}

	// Capture the task output
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Log error if task execution fails
		x_log.Error().
			Err(err).
			Str("task", task.Name).
			Bytes("output", output).
			Msg("task execution failed")
		return fmt.Errorf("failed to execute task: %w, Output: %s", err, string(output))
	}

	// Print task output if configured
	if task.IsPrintOutput {
		x_log.Debug().
			Str("task", task.Name).
			Int("output_len", len(output)).
			Msg("captured task output")
		fmt.Println("Task Output:")
		fmt.Println(string(output))
	}

	// Log success if task completes successfully
	x_log.Info().
		Str("task", task.Name).
		Msg("task completed successfully")
	return nil
}
