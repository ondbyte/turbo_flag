package flag

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestWriteMapToYAML(t *testing.T) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"nestedObject": map[string]interface{}{
			"nestedKey": "nestedValue",
		},
	}

	expectedYAML := `key1: value1
key2: value2
nestedObject:
  nestedKey: nestedValue
`

	yamlContent, err := WriteMapToYAML(data)
	if err != nil {
		t.Errorf("Failed to write map to YAML: %v", err)
	}

	if yamlContent != expectedYAML {
		t.Errorf("Unexpected YAML content. Expected:\n%s\nGot:\n%s", expectedYAML, yamlContent)
	}

}

func TestWriteMapToJSON(t *testing.T) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"nestedObject": map[string]interface{}{
			"nestedKey": "nestedValue",
		},
	}

	expectedJSON := `{"key1":"value1","key2":"value2","nestedObject":{"nestedKey":"nestedValue"}}`

	jsonContent, err := WriteMapToJSON(data)
	if err != nil {
		t.Errorf("Failed to write map to JSON: %v", err)
	}

	if jsonContent != expectedJSON {
		t.Errorf("Unexpected JSON content. Expected:\n%s\nGot:\n%s", expectedJSON, jsonContent)
	}
}

func TestWriteMapToPropertyFile(t *testing.T) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	expectedPropertyFile := `key1=value1
key2=value2
`

	propertyContent := WriteMapToPropertyFile(data)

	if propertyContent != expectedPropertyFile {
		t.Errorf("Unexpected Property file content. Expected:\n%s\nGot:\n%s", expectedPropertyFile, propertyContent)
	}
}
func TestReadPropertyFile(t *testing.T) {
	// Test Property file content
	propertyContent := `
# Example Property file
key1=value1
key2 = value2 # Comment after key-value pair
key3 = value3
# CommentedKey = CommentedValue
`

	// Call ReadPropertyFile function
	result, err := ReadPropertyFile(propertyContent)
	if err != nil {
		t.Errorf("Failed to read Property file: %v", err)
	}

	// Expected result
	expected := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	// Compare the actual result with the expected result
	if !compareProperty(result, expected) {
		t.Errorf("Unexpected Property file content. Expected: %v, Got: %v", expected, result)
	}
}

// compareProperty compares two Property objects for equality
func compareProperty(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if bValue, exists := b[key]; !exists || value != bValue {
			return false
		}
	}
	return true
}

func TestReadYAMLFile(t *testing.T) {
	// Test YAML content
	yamlContent := `
key1: value1
key2: value2
nestedObject:
  nestedKey: nestedValue
`

	// Call ReadYAMLFile function
	result, err := ReadYAMLFile(yamlContent)
	if err != nil {
		t.Errorf("Failed to read YAML file: %v", err)
	}

	// Expected result
	expected := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"nestedObject": map[string]interface{}{
			"nestedKey": "nestedValue",
		},
	}

	// Compare the actual result with the expected result
	if !compareYAML(result, expected) {
		t.Errorf("Unexpected YAML content. Expected: %v, Got: %v", expected, result)
	}
}

// compareYAML compares two YAML objects for equality
func compareYAML(a, b interface{}) bool {
	aStr, _ := yaml.Marshal(a)
	bStr, _ := yaml.Marshal(b)
	return string(aStr) == string(bStr)
}

func TestReadJSONFile(t *testing.T) {
	// Test JSON content
	jsonContent := `{
		"key1": "value1",
		"key2": "value2",
		"nestedObject": {
			"nestedKey": "nestedValue"
		}
	}`

	// Call ReadJSONFile function
	result, err := ReadJSONFile(jsonContent)
	if err != nil {
		t.Errorf("Failed to read JSON file: %v", err)
	}

	// Expected result
	expected := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"nestedObject": map[string]interface{}{
			"nestedKey": "nestedValue",
		},
	}

	// Compare the actual result with the expected result
	if !compareJSON(result, expected) {
		t.Errorf("Unexpected JSON content. Expected: %v, Got: %v", expected, result)
	}
}

// compareJSON compares two JSON objects for equality
func compareJSON(a, b interface{}) bool {
	aBytes, _ := json.Marshal(a)
	bBytes, _ := json.Marshal(b)
	return string(aBytes) == string(bBytes)
}
