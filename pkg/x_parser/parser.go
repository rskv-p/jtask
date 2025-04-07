package x_parser

import (
	"encoding/json"
	"fmt"

	"github.com/rskv-p/jtask/pkg/x_log"
)

//
// ---------- JSON Parser Implementation ----------

// JSON implements a parser for JSON data.
type JSON struct{}

// Parser returns a new instance of JSON parser.
func Parser() *JSON {
	x_log.Debug().Msg("creating new JSON parser instance")
	return &JSON{}
}

//
// ---------- Unmarshal ----------

// Unmarshal parses a JSON byte slice into a map[string]any.
// Logs the size of input and number of parsed keys.
func (p *JSON) Unmarshal(b []byte) (map[string]any, error) {
	x_log.Debug().
		Int("bytes", len(b)).
		Msg("unmarshalling JSON")

	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		x_log.Error().
			Err(err).
			Int("bytes", len(b)).
			Msg("failed to unmarshal JSON")
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	x_log.Info().
		Int("keys", len(out)).
		Msg("successfully unmarshalled JSON")

	return out, nil
}

//
// ---------- Marshal ----------

// Marshal serializes a map[string]any into JSON byte slice.
// Logs the map size and resulting JSON size.
func (p *JSON) Marshal(o map[string]any) ([]byte, error) {
	x_log.Debug().
		Int("map_size", len(o)).
		Msg("marshalling map to JSON")

	data, err := json.Marshal(o)
	if err != nil {
		x_log.Error().
			Err(err).
			Int("map_size", len(o)).
			Msg("failed to marshal map to JSON")
		return nil, fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	// Log only the first 100 characters of the marshaled data, ensuring we don't exceed the slice bounds
	logData := string(data)
	if len(logData) > 100 {
		logData = logData[:100]
	}
	x_log.Info().
		Int("bytes", len(data)).
		Str("json_data", logData).
		Msg("successfully marshalled map to JSON") // Safely log the marshaled JSON

	return data, nil
}
