package client

import (
	"bytes"
	"encoding/json"

	yaml "go.yaml.in/yaml/v3"
)

// toYAML renders a value as YAML, using its JSON field names
func toYAML(value any) ([]byte, error) {
	j, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	// JSON is valid YAML. Unlike encoding/json, the YAML decoder keeps
	// integers as integers instead of turning all numbers into floats.
	var obj any
	if err := yaml.Unmarshal(j, &obj); err != nil {
		return nil, err
	}

	var out bytes.Buffer
	encoder := yaml.NewEncoder(&out)
	encoder.SetIndent(2)
	if err := encoder.Encode(obj); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
