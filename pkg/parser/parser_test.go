package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yuvalwz/flags-gen/pkg/types"
)

func TestParser_ParseFile(t *testing.T) {
	// Create temporary test file
	tmpDir, err := os.MkdirTemp("", "flags-gen-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.go")
	testContent := `package main

// +flags-gen
// ServerConfig defines server configuration
type ServerConfig struct {
	// Host is the server hostname
	Host string ` + "`json:\"host\" default:\"localhost\"`" + `
	
	// Port is the server port
	Port int ` + "`json:\"port\" default:\"8080\"`" + `
	
	// Enable debug mode
	Debug bool ` + "`json:\"debug\"`" + `
	
	// Tags for filtering
	Tags []string ` + "`json:\"tags\" default:\"web,api\"`" + `
	
	// unexported field should be ignored
	secret string ` + "`json:\"secret\"`" + `
}

// Regular struct without annotation - should be ignored
type IgnoredConfig struct {
	Value string ` + "`json:\"value\"`" + `
}
`

	if err := os.WriteFile(testFile, []byte(testContent), 0o600); err != nil {
		t.Fatal(err)
	}

	parser := New()
	structs, err := parser.ParseFile(testFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(structs) != 1 {
		t.Fatalf("Expected 1 struct, got %d", len(structs))
	}

	config := structs[0]
	if config.Name != "ServerConfig" {
		t.Errorf("Expected struct name 'ServerConfig', got '%s'", config.Name)
	}

	if config.PackageName != "main" {
		t.Errorf("Expected package name 'main', got '%s'", config.PackageName)
	}

	// Should have 4 exported fields (secret is unexported)
	if len(config.Fields) != 4 {
		t.Fatalf("Expected 4 fields, got %d", len(config.Fields))
	}

	// Check first field (Host)
	hostField := config.Fields[0]
	if hostField.Name != "Host" {
		t.Errorf("Expected first field name 'Host', got '%s'", hostField.Name)
	}
	if hostField.Type != "string" {
		t.Errorf("Expected first field type 'string', got '%s'", hostField.Type)
	}
	if hostField.FlagName != "host" {
		t.Errorf("Expected flag name 'host', got '%s'", hostField.FlagName)
	}
	if hostField.DefaultValue != "localhost" {
		t.Errorf("Expected default value 'localhost', got '%v'", hostField.DefaultValue)
	}
	if hostField.Description != "Host is the server hostname" {
		t.Errorf("Expected description 'Host is the server hostname', got '%s'", hostField.Description)
	}

	// Check Port field
	portField := config.Fields[1]
	if portField.Name != "Port" || portField.Type != "int" {
		t.Errorf("Port field not parsed correctly: name=%s, type=%s", portField.Name, portField.Type)
	}
	if portField.DefaultValue != 8080 {
		t.Errorf("Expected Port default value 8080, got %v", portField.DefaultValue)
	}

	// Check Debug field
	debugField := config.Fields[2]
	if debugField.Name != "Debug" || debugField.Type != "bool" {
		t.Errorf("Debug field not parsed correctly: name=%s, type=%s", debugField.Name, debugField.Type)
	}

	// Check Tags field
	tagsField := config.Fields[3]
	if tagsField.Name != "Tags" || tagsField.Type != "[]string" {
		t.Errorf("Tags field not parsed correctly: name=%s, type=%s", tagsField.Name, tagsField.Type)
	}
	if tags, ok := tagsField.DefaultValue.([]string); !ok || len(tags) != 2 || tags[0] != "web" || tags[1] != "api" {
		t.Errorf("Tags default value not parsed correctly: %v", tagsField.DefaultValue)
	}
}

func TestParser_toKebabCase(t *testing.T) {
	parser := New()

	tests := []struct {
		input    string
		expected string
	}{
		{"ProbeAddr", "probe-addr"},
		{"EnableLeaderElection", "enable-leader-election"},
		{"V", "v"},
		{"HTTPPort", "http-port"},
		{"XMLParser", "xml-parser"},
	}

	for _, test := range tests {
		result := parser.toKebabCase(test.input)
		if result != test.expected {
			t.Errorf("toKebabCase(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestParser_hasAnnotation(_ *testing.T) {
	// This would need more complex AST setup to test properly
	// For now, we test it indirectly through ParseFile
}

func TestGetFlagMethod(t *testing.T) {
	tests := []struct {
		fieldType string
		expected  string
		exists    bool
	}{
		{"string", "StringVar", true},
		{"int", "IntVar", true},
		{"bool", "BoolVar", true},
		{"[]string", "StringSliceVar", true},
		{"time.Duration", "DurationVar", true},
		{"unsupported", "", false},
	}

	for _, test := range tests {
		method, exists := types.GetFlagMethod(test.fieldType)
		if exists != test.exists {
			t.Errorf("GetFlagMethod(%s) exists = %v, expected %v", test.fieldType, exists, test.exists)
		}
		if method != test.expected {
			t.Errorf("GetFlagMethod(%s) = %s, expected %s", test.fieldType, method, test.expected)
		}
	}
}

func TestParser_parseDefaultValue(t *testing.T) {
	parser := New()

	tests := []struct {
		value     string
		fieldType string
		expected  interface{}
	}{
		{"hello", "string", "hello"},
		{"42", "int", 42},
		{"true", "bool", true},
		{"false", "bool", false},
		{"web,api", "[]string", []string{"web", "api"}},
		{"", "[]string", []string{}},
		{"30s", "time.Duration", "30s"},
	}

	for _, test := range tests {
		result := parser.parseDefaultValue(test.value, test.fieldType)
		switch expected := test.expected.(type) {
		case []string:
			if resultSlice, ok := result.([]string); !ok {
				t.Errorf("parseDefaultValue(%s, %s) type = %T, expected []string", test.value, test.fieldType, result)
			} else if len(resultSlice) != len(expected) {
				t.Errorf("parseDefaultValue(%s, %s) length = %d, expected %d", test.value, test.fieldType, len(resultSlice), len(expected))
			} else {
				for i, v := range expected {
					if resultSlice[i] != v {
						t.Errorf("parseDefaultValue(%s, %s)[%d] = %s, expected %s", test.value, test.fieldType, i, resultSlice[i], v)
					}
				}
			}
		default:
			if result != expected {
				t.Errorf("parseDefaultValue(%s, %s) = %v, expected %v", test.value, test.fieldType, result, expected)
			}
		}
	}
}
