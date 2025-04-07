package x_config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rskv-p/jtask/pkg/x_log"
)

//
// ---------- Constants ----------

// Constants for config file paths and environment variables
const (
	ApplicationName = "jtask"
	DefaultConfig   = "space/cfg/jtask.json" // Default config file path
	LocalConfig     = ".jtask.json"          // Local config file path
	GlobalConfig    = "JTASK_CONFIG"         // Environment variable for custom config file path
)

//
// ---------- Default Config ----------

// defaultConfig provides default values for the fields that can be overridden in the config file.
var defaultConfig = Config{
	AppName:       ApplicationName,
	Version:       "1.0.0",        // Default version
	Logger:        x_log.Config{}, // You might want to set default logger config here
	MaxConcurrent: 5,              // Default max concurrent tasks
}

//
// ---------- Config Structure ----------

// Config is the root application config structure that holds the application settings.
type Config struct {
	AppName       string       `json:"AppName"`       // Application name
	Version       string       `json:"Version"`       // Application version
	Logger        x_log.Config `json:"Logger"`        // Logger configuration
	MaxConcurrent int          `json:"MaxConcurrent"` // max number of concurrent tasks
}

//
// ---------- LoadConfig Function ----------

// LoadConfig checks multiple locations for the configuration file and loads it.
// It looks in the following order:
// 1. Local config file (.jtask.json)
// 2. Global config environment variable (JTASK_CONFIG)
// 3. Default config file (_data/config.json)
func LoadConfig() (*Config, error) {
	// Define possible paths to check for the config file
	paths := []string{
		LocalConfig,
		os.Getenv(GlobalConfig), // Check the environment variable
		DefaultConfig,           // Default config path
	}

	// Iterate through paths and try to load the config
	for _, path := range paths {
		if path == "" {
			continue // Skip empty paths
		}

		// Log the path currently being checked
		// x_log.Debug().
		// 	Str("source", path).
		// 	Msg("checking config path")

		// Try reading and parsing the config from the path
		cfg, err := readConfig(path)
		if err == nil {
			// If config is successfully loaded, apply defaults if necessary
			applyDefaults(cfg)

			// Log the success and return
			// x_log.Debug().
			// 	Str("source", path).
			// 	Msg("config loaded successfully")
			return cfg, nil
		}

		// If file is not found, skip and continue with the next path
		if errors.Is(err, os.ErrNotExist) {
			// x_log.Debug().
			// 	Str("source", path).
			// 	Msg("config file not found, trying next path")
			continue
		}

		// Log warnings for any other errors encountered
		// x_log.Warn().
		// 	Err(err).
		// 	Str("path", path).
		// 	Msg("failed to load config")
	}

	// If no valid config file is found, return an error
	return nil, fmt.Errorf("no valid config file found")
}

//
// ---------- applyDefaults Function ----------

// applyDefaults fills missing config values from defaultConfig
func applyDefaults(cfg *Config) {
	// Apply default values if necessary
	if cfg.AppName == "" {
		cfg.AppName = defaultConfig.AppName
	}
	if cfg.Version == "" {
		cfg.Version = defaultConfig.Version
	}
	if cfg.MaxConcurrent == 0 {
		cfg.MaxConcurrent = defaultConfig.MaxConcurrent
	}

	// You can add any additional logic for default values here, if needed
}

//
// ---------- readConfig Function ----------

// readConfig reads a JSON config file from the specified path and unmarshals it into a Config object.
func readConfig(path string) (*Config, error) {
	// Read the file data
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err // Return error if file reading fails
	}

	// Unmarshal the JSON data into the Config struct
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err) // Return error if parsing fails
	}

	// Return the parsed config
	return &cfg, nil
}
