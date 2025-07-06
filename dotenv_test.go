package dotenv

import (
	"os"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary .env file
	content := `# This is a comment
KEY1=value1
KEY2="value with spaces"
KEY3='single quoted value'
EMPTY_VAR=
NUMERIC_VAR=42
BOOLEAN_VAR=true
`

	// Clear environment
	testKeys := []string{"KEY1", "KEY2", "KEY3", "EMPTY_VAR", "NUMERIC_VAR", "BOOLEAN_VAR"}
	for _, key := range testKeys {
		os.Unsetenv(key)
	}

	// Create temp file
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	// Test Load
	err := Load(tmpFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify values
	tests := map[string]string{
		"KEY1":        "value1",
		"KEY2":        "value with spaces",
		"KEY3":        "single quoted value",
		"EMPTY_VAR":   "",
		"NUMERIC_VAR": "42",
		"BOOLEAN_VAR": "true",
	}

	for key, expected := range tests {
		if actual := os.Getenv(key); actual != expected {
			t.Errorf("Expected %s=%q, got %q", key, expected, actual)
		}
	}
}

func TestOverload(t *testing.T) {
	// Set existing environment variable
	os.Setenv("TEST_OVERLOAD", "original")
	defer os.Unsetenv("TEST_OVERLOAD")

	content := "TEST_OVERLOAD=overridden"
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	// Test that Load doesn't override
	err := Load(tmpFile)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if os.Getenv("TEST_OVERLOAD") != "original" {
		t.Error("Load should not override existing environment variables")
	}

	// Test that Overload does override
	err = Overload(tmpFile)
	if err != nil {
		t.Fatalf("Overload failed: %v", err)
	}

	if os.Getenv("TEST_OVERLOAD") != "overridden" {
		t.Error("Overload should override existing environment variables")
	}
}

func TestRead(t *testing.T) {
	content := `KEY1=value1
KEY2=value2
`
	tmpFile := createTempEnvFile(t, content)
	defer os.Remove(tmpFile)

	env, err := Read(tmpFile)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	expected := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	for key, expectedValue := range expected {
		if actualValue, exists := env[key]; !exists {
			t.Errorf("Expected key %s to exist", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected %s=%q, got %q", key, expectedValue, actualValue)
		}
	}
}

func TestParse(t *testing.T) {
	content := `# Comment
KEY1=value1
KEY2="quoted value"
KEY3='single quoted'
export EXPORTED=exported_value
KEY4: yaml_style
EMPTY=
`

	env, err := Parse(strings.NewReader(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expected := map[string]string{
		"KEY1":     "value1",
		"KEY2":     "quoted value",
		"KEY3":     "single quoted",
		"EXPORTED": "exported_value",
		"KEY4":     "yaml_style",
		"EMPTY":    "",
	}

	for key, expectedValue := range expected {
		if actualValue, exists := env[key]; !exists {
			t.Errorf("Expected key %s to exist", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected %s=%q, got %q", key, expectedValue, actualValue)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	content := "KEY1=value1\nKEY2=value2"

	env, err := Unmarshal(content)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if env["KEY1"] != "value1" || env["KEY2"] != "value2" {
		t.Error("Unmarshal did not parse correctly")
	}
}

func TestMarshal(t *testing.T) {
	env := map[string]string{
		"KEY1":    "value1",
		"KEY2":    "value with spaces",
		"KEY3":    "",
		"SPECIAL": "value\nwith\tspecial\"chars",
	}

	result, err := Marshal(env)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Parse it back to verify
	parsed, err := Unmarshal(result)
	if err != nil {
		t.Fatalf("Failed to parse marshaled content: %v", err)
	}

	for key, expected := range env {
		if actual, exists := parsed[key]; !exists {
			t.Errorf("Key %s missing after marshal/unmarshal", key)
		} else if actual != expected {
			t.Errorf("Value mismatch for %s: expected %q, got %q", key, expected, actual)
		}
	}
}

func TestVariableExpansion(t *testing.T) {
	content := `BASE=hello
EXPANDED=${BASE}_world
SIMPLE=$BASE
NESTED=${BASE}_${BASE}
UNDEFINED=${UNDEFINED_VAR}
`

	env, err := Parse(strings.NewReader(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	tests := map[string]string{
		"BASE":      "hello",
		"EXPANDED":  "hello_world",
		"SIMPLE":    "hello",
		"NESTED":    "hello_hello",
		"UNDEFINED": "",
	}

	for key, expected := range tests {
		if actual := env[key]; actual != expected {
			t.Errorf("Expected %s=%q, got %q", key, expected, actual)
		}
	}
}

func TestEscapeSequences(t *testing.T) {
	content := `NEWLINE="line1\nline2"
TAB="tab\there"
QUOTE="say \"hello\""
BACKSLASH="back\\slash"
SINGLE='no\nexpansion'
`

	env, err := Parse(strings.NewReader(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	tests := map[string]string{
		"NEWLINE":   "line1\nline2",
		"TAB":       "tab\there",
		"QUOTE":     `say "hello"`,
		"BACKSLASH": `back\slash`,
		"SINGLE":    `no\nexpansion`, // Single quotes don't expand escapes
	}

	for key, expected := range tests {
		if actual := env[key]; actual != expected {
			t.Errorf("Expected %s=%q, got %q", key, expected, actual)
		}
	}
}

func TestInlineComments(t *testing.T) {
	content := `KEY1=value1 # This is a comment
KEY2="quoted value" # Comment after quotes
KEY3=value3# No space before comment
KEY4="value with # inside quotes"
`

	env, err := Parse(strings.NewReader(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	tests := map[string]string{
		"KEY1": "value1",
		"KEY2": "quoted value",
		"KEY3": "value3",
		"KEY4": "value with # inside quotes",
	}

	for key, expected := range tests {
		if actual := env[key]; actual != expected {
			t.Errorf("Expected %s=%q, got %q", key, expected, actual)
		}
	}
}

func TestParseHelpers(t *testing.T) {
	// Set test environment variables
	os.Setenv("TEST_INT", "42")
	os.Setenv("TEST_BOOL_TRUE", "true")
	os.Setenv("TEST_BOOL_FALSE", "false")
	os.Setenv("TEST_FLOAT", "3.14")
	os.Setenv("TEST_REQUIRED", "exists")
	defer func() {
		os.Unsetenv("TEST_INT")
		os.Unsetenv("TEST_BOOL_TRUE")
		os.Unsetenv("TEST_BOOL_FALSE")
		os.Unsetenv("TEST_FLOAT")
		os.Unsetenv("TEST_REQUIRED")
	}()

	// Test ParseInt
	if ParseInt("TEST_INT", 0) != 42 {
		t.Error("ParseInt failed")
	}
	if ParseInt("NONEXISTENT", 10) != 10 {
		t.Error("ParseInt default failed")
	}

	// Test ParseBool
	if !ParseBool("TEST_BOOL_TRUE", false) {
		t.Error("ParseBool true failed")
	}
	if ParseBool("TEST_BOOL_FALSE", true) {
		t.Error("ParseBool false failed")
	}
	if !ParseBool("NONEXISTENT", true) {
		t.Error("ParseBool default failed")
	}

	// Test ParseFloat
	if ParseFloat("TEST_FLOAT", 0.0) != 3.14 {
		t.Error("ParseFloat failed")
	}
	if ParseFloat("NONEXISTENT", 1.0) != 1.0 {
		t.Error("ParseFloat default failed")
	}

	// Test GetRequired
	if GetRequired("TEST_REQUIRED") != "exists" {
		t.Error("GetRequired failed")
	}

	// Test GetWithDefault
	if GetWithDefault("TEST_REQUIRED", "default") != "exists" {
		t.Error("GetWithDefault with existing value failed")
	}
	if GetWithDefault("NONEXISTENT", "default") != "default" {
		t.Error("GetWithDefault with default failed")
	}
}

func TestErrorCases(t *testing.T) {
	// Test invalid file
	_, err := Read("nonexistent.env")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}

	// Test invalid format
	content := "INVALID LINE FORMAT"
	_, err = Parse(strings.NewReader(content))
	if err == nil {
		t.Error("Expected error for invalid format")
	}
}

// Helper function to create temporary .env file
func createTempEnvFile(t *testing.T, content string) string {
	tmpFile := t.TempDir() + "/.env"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return tmpFile
}
