package cmd

import (
	"sync"

	"github.com/charmbracelet/huh"
	"github.com/rskv-p/jtask/pkg/x_log"
	"github.com/rskv-p/jtask/pkg/x_task"
	"github.com/spf13/cobra"
)

var runsCmd = &cobra.Command{
	Use:   "runs",
	Short: "Select multiple tasks to run in parallel",
	Long:  "Select multiple tasks from the list and execute them in parallel.",
	Run: func(cmd *cobra.Command, args []string) {
		// Log the beginning of task loading
		x_log.Info().
			Str("path", pathFlag).
			Msg("loading tasks from file")

		// Load tasks from file
		tasks, err := x_task.LoadTasks(pathFlag)
		if err != nil {
			// Log error if tasks can't be loaded
			x_log.Error().
				Err(err).
				Str("path", pathFlag).
				Msg("failed to load tasks")
			return
		}

		// If no tasks are available, log and return
		if len(tasks.Data) == 0 {
			x_log.Info().Msg("no tasks available to select")
			return
		}

		// Log the task count and start building selection options
		x_log.Info().
			Int("count", len(tasks.Data)).
			Msg("building task selection options")

		options := createHuhOptions(tasks)

		// Variable to store selected tasks
		var selectedTasks []string

		// Prompt the user to select multiple tasks
		x_log.Debug().Msg("prompting user to select multiple tasks")
		if err := huh.NewMultiSelect[string]().
			Title("Select tasks to run in parallel:").
			Options(options...).
			Value(&selectedTasks).
			Run(); err != nil {
			// Log error if task selection fails
			x_log.Error().
				Err(err).
				Msg("task selection failed")
			return
		}

		// Log selected tasks count and names
		x_log.Info().
			Int("selected", len(selectedTasks)).
			Interface("tasks", selectedTasks).
			Msg("user selected tasks")

		// Log the selected tasks to the user
		x_log.Info().
			Strs("selected_tasks", selectedTasks).
			Msg("the following tasks were selected")

		// ---------- Parallel Task Execution ----------
		// Get max concurrent tasks from config
		maxConcurrent := cfg.MaxConcurrent
		if maxConcurrent <= 0 {
			maxConcurrent = 5 // Set a default if not specified
		}
		sem := make(chan struct{}, maxConcurrent) // Semaphore to limit concurrency

		// WaitGroup to ensure all tasks are processed
		var wg sync.WaitGroup
		for _, name := range selectedTasks {
			wg.Add(1)
			go func(taskName string) {
				defer wg.Done()

				// Use semaphore to limit concurrency
				sem <- struct{}{}
				defer func() { <-sem }()

				// Log task start
				x_log.Info().
					Str("task", taskName).
					Msg("starting task")

				// Find the task by name
				if task := findTaskByName(tasks, taskName); task != nil {
					// Execute the task
					if err := executeTask(task); err != nil {
						// Log task execution failure
						x_log.Error().
							Str("task", task.Name).
							Err(err).
							Msg("task execution failed")
					}
				} else {
					// Log if the task is not found
					x_log.Warn().
						Str("task", taskName).
						Msg("task not found")
				}
			}(name)
		}
		wg.Wait() // Wait for all tasks to be processed

		// Log after all tasks are processed
		x_log.Info().
			Int("done", len(selectedTasks)).
			Msg("all selected tasks processed")
	},
}

// ---------- Command Initialization ----------
func init() {
	rootCmd.AddCommand(runsCmd) // Register the 'runs' command
}

// ---------- Helper Functions ----------

// createHuhOptions builds selection options from tasks.
func createHuhOptions(tasks *x_task.TaskCollection) []huh.Option[string] {
	var options []huh.Option[string]
	for _, task := range tasks.Data {
		options = append(options, huh.NewOption(task.Name, task.Name)) // Add task options for selection
	}
	return options
}

// findTaskByName returns a task by its name.
func findTaskByName(tasks *x_task.TaskCollection, name string) *x_task.Task {
	for _, task := range tasks.Data {
		if task.Name == name {
			return task // Return the task if found by name
		}
	}
	return nil // Return nil if task not found
}
