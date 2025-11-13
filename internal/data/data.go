package data

import (
	"encoding/json"
	"fmt"
	"os"
)

// Element is one JSON value from the input array.
// It is left as raw JSON so any shape (object, string, number, etc.) is supported.
type Element = json.RawMessage

// Elements is the in-memory representation of the JSON array.
type Elements []Element

// LoadElements loads a JSON file that contains an array of elements.
func LoadElements(path string) (Elements, error) {
	var elements Elements

	data, err := os.ReadFile(path)
	if err != nil {
		return Elements{}, fmt.Errorf("read data: %w", err)
	}

	if err := json.Unmarshal(data, &elements); err != nil {
		return Elements{}, fmt.Errorf("unmarshal data: %w", err)
	}

	return elements, nil
}
