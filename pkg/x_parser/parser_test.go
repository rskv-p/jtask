package x_parser

import (
	"encoding/json"
	"reflect"
	"testing"
)

//
// ---------- Unit Tests for JSON Parser ----------

// TestUnmarshal verifies that JSON bytes are correctly converted into a map.
func TestUnmarshal(t *testing.T) {
	parser := Parser()

	// Sample input JSON
	input := []byte(`{"key": "value", "num": 42}`)

	// Expected result (note: numbers default to float64 in JSON)
	expected := map[string]any{
		"key": "value",
		"num": float64(42),
	}

	// Run parser
	result, err := parser.Unmarshal(input)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	// Assert equality
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unmarshal result mismatch.\nExpected: %v\nGot: %v", expected, result)
	}
}

// TestMarshal verifies that a map is correctly marshalled into JSON bytes.
func TestMarshal(t *testing.T) {
	parser := Parser()

	// Input map to convert
	input := map[string]any{
		"key": "value",
		"num": 42,
	}

	// Expected JSON output (key order may vary)
	expected := `{"key":"value","num":42}`

	// Marshal the input map
	data, err := parser.Marshal(input)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	// Unmarshal both actual and expected for safe comparison
	var resultMap map[string]any
	var expectedMap map[string]any

	if err := json.Unmarshal(data, &resultMap); err != nil {
		t.Fatalf("Unmarshal result JSON error: %v", err)
	}
	if err := json.Unmarshal([]byte(expected), &expectedMap); err != nil {
		t.Fatalf("Unmarshal expected JSON error: %v", err)
	}

	// Compare as maps (ignoring key order)
	if !reflect.DeepEqual(resultMap, expectedMap) {
		t.Errorf("Marshal result mismatch.\nExpected: %v\nGot: %v", expectedMap, resultMap)
	}
}
