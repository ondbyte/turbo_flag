package flag

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestEnvToMap(t *testing.T) {
	content := `# Comment line
KEY1=Value1
KEY2 = Value2
KEY3=Value3
# Another comment line
KEY4 = Value4 = ExtraValue`

	expected := map[string]string{
		"KEY1": "Value1",
		"KEY2": "Value2",
		"KEY3": "Value3",
		"KEY4": "Value4 = ExtraValue",
	}

	result, err := EnvToMap(content)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Result does not match expected value.\nExpected: %v\nGot: %v", expected, result)
	}
}

func Test_jsonnify(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "true",
			args: args{
				v: true,
			},
			want:    "true",
			wantErr: false,
		},
		{
			name: "json",
			args: args{
				v: map[string]interface{}{
					"yadu": true,
				},
			},
			want:    `{"yadu":true}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonnify(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonnify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("jsonnify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapToTOML(t *testing.T) {
	config := map[string]interface{}{
		"server": map[interface{}]interface{}{
			"address": "localhost",
			"port":    8080,
			"enabled": true,
		},
		"database": map[interface{}]interface{}{
			"host":     "127.0.0.1",
			"port":     5432,
			"name":     "mydb",
			"username": "admin",
			"password": "password123",
		},
	}

	expectedTOML := `[database]
  host = "127.0.0.1"
  name = "mydb"
  password = "password123"
  port = 5432
  username = "admin"

[server]
  address = "localhost"
  enabled = true
  port = 8080
`

	tomlString, err := MapToTOML(config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// Compare the generated TOML string with the expected TOML
	if tomlString != expectedTOML {
		t.Errorf("Unexpected TOML string. Got:\n%s\nExpected:\n%s", tomlString, expectedTOML)
	}
}

func TestMapToTOML2(t *testing.T) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"nestedObject": map[interface{}]interface{}{
			"nestedKey": "nestedValue",
		},
	}

	expectedYAML := `key1: value1
key2: value2
nestedObject:
  nestedKey: nestedValue
`

	yamlContent, err := MapToYAML(data)
	if err != nil {
		t.Errorf("Failed to write map to YAML: %v", err)
	}

	if yamlContent != expectedYAML {
		t.Errorf("Unexpected YAML content. Expected:\n%s\nGot:\n%s", expectedYAML, yamlContent)
	}

}

func TestMapToJSON(t *testing.T) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"nestedObject": map[interface{}]interface{}{
			"nestedKey": "nestedValue",
		},
	}

	expectedJSON := `{"key1":"value1","key2":"value2","nestedObject":{"nestedKey":"nestedValue"}}`

	jsonContent, err := MapToJSON(data)
	if err != nil {
		t.Errorf("Failed to write map to JSON: %v", err)
	}

	if jsonContent != expectedJSON {
		t.Errorf("Unexpected JSON content. Expected:\n%s\nGot:\n%s", expectedJSON, jsonContent)
	}
}

// compareProperty compares two Property objects for equality
func compareProperty(a, b map[interface{}]interface{}) bool {
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

func TestYAMLToMap(t *testing.T) {
	// Test YAML content
	yamlContent := `
key1: value1
key2: value2
nestedObject:
  nestedKey: nestedValue
`

	// Call ReadYAMLFile function
	result, err := YAMLToMap(yamlContent)
	if err != nil {
		t.Errorf("Failed to read YAML file: %v", err)
	}

	// Expected result
	expected, err := stringMap(map[interface{}]interface{}{
		"key1": "value1",
		"key2": "value2",
		"nestedObject": map[interface{}]interface{}{
			"nestedKey": "nestedValue",
		},
	})
	if err != nil {
		t.Fatal("stringMap error")
	}
	// Compare the actual result with the expected result
	if !compareYAML(result, expected) {
		t.Errorf("Unexpected YAML content. Expected: %v, Got: %v", expected, result)
	}
}

// compareYAML compares two YAML objects for equality
func compareYAML(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func TestJSONToMap(t *testing.T) {
	// Test JSON content
	jsonContent := `{
		"key1": "value1",
		"key2": "value2",
		"nestedObject": {
			"nestedKey": "nestedValue"
		}
	}`

	// Call ReadJSONFile function
	result, err := JSONToMap(jsonContent)
	if err != nil {
		t.Errorf("Failed to read JSON file: %v", err)
	}

	// Expected result
	expected, err := stringMap(map[interface{}]interface{}{
		"key1": "value1",
		"key2": "value2",
		"nestedObject": map[interface{}]interface{}{
			"nestedKey": "nestedValue",
		},
	})

	if err != nil {
		t.Fatal("stringMap err")
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

func TestStringMap(t *testing.T) {
	a := func() {
		m := map[interface{}]interface{}{
			"yadu": "nandan",
			"c": map[interface{}]interface{}{
				"maya": "hm",
			},
		}
		mE := map[string]interface{}{
			"yadu": "nandan",
			"c": map[string]interface{}{
				"maya": "hm",
			},
		}
		m2, err := stringMap(m)
		if !reflect.DeepEqual(m2, mE) {
			t.Fatal(err)
		}
		t.Log(m2)
	}
	b := func() {
		m := map[string]interface{}{
			"yadu": "nandan",
			"c": map[string]interface{}{
				"maya": "hm",
			},
		}
		mE := map[string]interface{}{
			"yadu": "nandan",
			"c": map[string]interface{}{
				"maya": "hm",
			},
		}
		m2, err := stringMap(m)
		if !reflect.DeepEqual(m2, mE) {
			t.Fatal(err)
		}
		t.Log(m2)
	}
	a()
	b()
}

func TestGetValueByDotNotation(t *testing.T) {
	data := map[string]interface{}{
		"foo": map[interface{}]interface{}{
			"bar": 42,
		},
		"baz": "hello",
		"y": map[interface{}]interface{}{
			"a": map[interface{}]interface{}{
				"d": map[interface{}]interface{}{
					"u": map[string]interface{}{
						"1": true,
						"2": map[string]interface{}{
							"map": true,
						},
					},
				},
			},
		},
	}

	tests := []struct {
		notation string
		expected string
		errMsg   bool
	}{
		{"foo.bar.baz", "", false},
		{"foo.bar", "42", true},
		{"baz", "hello", true},
		{"nonexistent", "", false},
		{"foo.nonexistent", "", false},
		{"y.a.d.u.1", "true", true},
		{"y.a.d.u.2", `{"map":true}`, true},
	}

	for _, test := range tests {
		value, err := getValueByDotNotation(data, test.notation)

		if value != test.expected && err != nil {
			t.Errorf("Unexpected value for notation %s. Expected: %v, Got: %v", test.notation, test.expected, value)
		}
	}
}

func TestAddValueByDotNotation(t *testing.T) {
	data := map[string]interface{}{
		"foo": map[interface{}]interface{}{
			"bar": "baz",
		},
	}
	data, err := setValueByDotNotation(data, "hello.another.key", true)
	if err != nil {
		t.Fatal(err)
	}
	v, err := getValueByDotNotation(data, "hello.another.key")
	if err != nil || v != "true" {
		t.Errorf("failed to get value : %v", err)
	}

	data, err = setValueByDotNotation(data, "foo.qux", "value")
	if err != nil {
		t.Fatal(err)
	}
	v, err = getValueByDotNotation(data, "foo.qux")
	if err != nil || v != "value" {
		t.Errorf("failed to get value : %v", err)
	}

	data, err = setValueByDotNotation(data, "hello.world", 123)
	if err != nil {
		t.Fatal(err)
	}
	v, err = getValueByDotNotation(data, "hello.world")
	if err != nil || v != "123" {
		t.Errorf("failed to get value : %v", err)
	}

}
