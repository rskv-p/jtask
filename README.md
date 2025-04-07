# JsonTask - Task Runner

JsonTask is a lightweight and flexible task runner designed for executing tasks defined in a JSON configuration file. It supports running multiple tasks in parallel with customizable configurations, logging, and error handling.

## Features

- Load tasks from a JSON configuration file.
- Run tasks either sequentially or in parallel.
- Limit the number of concurrent tasks using configuration.
- Detailed logging for task execution.
- Supports custom logger configuration.
- Simple command-line interface (CLI) to run tasks.

## Installation

To install JsonTask, clone the repository and build the application:

```bash
git clone https://github.com/yourusername/jtask.git
cd jtask
go build -o jtask
```

Alternatively, you can install it directly using `go get`:

```bash
go get github.com/yourusername/jtask
```

## Configuration

JsonTask uses a JSON configuration file (`config.json`) to load tasks and configuration settings. You can specify the configuration file path using the `--config` flag.

Here is an example of the default configuration (`config.json`):

```json
{
  "AppName": "JsonTask",
  "Version": "1.0.0",
  "MaxConcurrent": 5,
  "Logger": {
    "Level": "info",
    "LogFile": "logs/app.log",
    "ToConsole": true,
    "ToFile": true,
    "ColoredFile": true,
    "Style": "dark",
    "MaxSize": 10,
    "MaxBackups": 5,
    "MaxAge": 7,
    "Compress": true
  }
}
```

### Fields:

- `AppName`: The name of the application.
- `Version`: The version of the application.
- `MaxConcurrent`: The maximum number of concurrent tasks to execute.
- `Logger`: Configuration for the logger.
  - `Level`: The log level (`info`, `debug`, `warn`, `error`).
  - `LogFile`: Path to the log file.
  - `ToConsole`: Boolean to log output to the console.
  - `ToFile`: Boolean to log output to a file.
  - `ColoredFile`: Boolean to apply color formatting in the log file.
  - `Style`: Log style (`dark`, `light`).
  - `MaxSize`: Maximum log file size in MB.
  - `MaxBackups`: Maximum number of backups to keep.
  - `MaxAge`: Number of days to keep log backups.
  - `Compress`: Boolean to compress old logs.

## Usage

### Run Tasks

To execute tasks defined in your configuration, use the following command:

```bash
./jtask run
```

This will prompt you to select tasks from the list and execute them. If you want to run multiple tasks in parallel, JsonTask will handle the concurrency based on the configuration.

### Run Multiple Tasks in Parallel

To run multiple tasks in parallel, use:

```bash
./jtask runs
```

It will let you select multiple tasks to run concurrently, and you can adjust the number of parallel tasks using the `MaxConcurrent` setting in the configuration.

### Command Flags

- `--config`: Path to the configuration file (default is `_data/config.json`).
- `--help`: Show help information about the commands.

## Logging

JsonTask uses a logger to track task execution. The logger can output to both console and log files. The log level can be adjusted using the `Logger` configuration.

- Errors, warnings, and info logs are generated as tasks are executed.
- Logs can be colored based on the log level for easy reading.

## Example Output

When executing tasks, JsonTask will provide logs like:

```
INF: starting task "task1"
INF: task "task1" completed successfully
INF: starting task "task2"
ERR: task "task2" failed with error: [error details]
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
