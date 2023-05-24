package flag

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// WriteMapToYAML writes a map to a YAML string
func WriteMapToYAML(data map[string]interface{}) (string, error) {
	yamlContent, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(yamlContent), nil
}

// WriteMapToJSON writes a map to a JSON string
func WriteMapToJSON(data map[string]interface{}) (string, error) {
	jsonContent, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonContent), nil
}

// WriteMapToPropertyFile writes a map to a Property file string
func WriteMapToPropertyFile(data map[string]interface{}) string {
	var sb strings.Builder
	for key, value := range data {
		sb.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}
	return sb.String()
}

// ReadJSONFile reads the contents of a JSON file from a string and returns a map[string]interface{}
func ReadJSONFile(content string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(content), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ReadYAMLFile reads the contents of a YAML file from a string and returns a map[string]interface{}
func ReadYAMLFile(content string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(content), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ReadPropertyFile reads the contents of a Property file from a string and returns a map[string]string
func ReadPropertyFile(content string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.Split(line, "#")
		keyValue := strings.TrimSpace(parts[0])
		if len(keyValue) == 0 {
			continue
		}
		parts = strings.SplitN(keyValue, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}

	return result, nil
}
