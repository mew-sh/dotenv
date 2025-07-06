// Package dotenv provides functionality to load environment variables from .env files.
//
// This is a modern Go implementation inspired by the joho/godotenv library
// but updated to use latest Go features and best practices.
//
// Example usage:
//
//	// Load .env file from current directory
//	err := dotenv.Load()
//
//	// Load specific files
//	err := dotenv.Load(".env.local", ".env.production")
//
//	// Read without setting environment variables
//	env, err := dotenv.Read(".env")
//
//	// Parse from string or reader
//	env, err := dotenv.Parse(strings.NewReader("KEY=value"))
package dotenv

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Default .env filename
const DefaultEnvFile = ".env"

// Load reads the specified .env files and loads the environment variables.
// If no files are specified, it defaults to loading ".env" from the current directory.
// Existing environment variables take precedence and will not be overwritten.
func Load(filenames ...string) error {
	return load(false, filenames...)
}

// Overload reads the specified .env files and loads the environment variables.
// Unlike Load, this will overwrite existing environment variables.
func Overload(filenames ...string) error {
	return load(true, filenames...)
}

// Read reads the specified .env files and returns a map of key-value pairs
// without modifying the actual environment variables.
func Read(filenames ...string) (map[string]string, error) {
	if len(filenames) == 0 {
		filenames = []string{DefaultEnvFile}
	}

	result := make(map[string]string)

	for _, filename := range filenames {
		env, err := readFile(filename)
		if err != nil {
			return nil, err
		}

		// Merge maps, later files take precedence
		for key, value := range env {
			result[key] = value
		}
	}

	return result, nil
}

// Parse reads environment variables from an io.Reader and returns a map.
func Parse(reader io.Reader) (map[string]string, error) {
	parser := NewParser()
	return parser.Parse(reader)
}

// Unmarshal parses a .env formatted string and returns a map of key-value pairs.
func Unmarshal(data string) (map[string]string, error) {
	return Parse(strings.NewReader(data))
}

// Marshal converts a map of environment variables to .env file format.
func Marshal(env map[string]string) (string, error) {
	if len(env) == 0 {
		return "", nil
	}

	// Sort keys for consistent output
	keys := make([]string, 0, len(env))
	for key := range env {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var lines []string
	for _, key := range keys {
		value := env[key]
		line := formatEnvLine(key, value)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}

// Write serializes the environment map and writes it to a file.
func Write(env map[string]string, filename string) error {
	content, err := Marshal(env)
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file with appropriate permissions
	return os.WriteFile(filename, []byte(content+"\n"), 0644)
}

// Must is a helper that wraps Load and panics if an error occurs.
// This is useful for loading configuration at startup where failure
// should halt the program.
func Must(filenames ...string) {
	if err := Load(filenames...); err != nil {
		panic(fmt.Sprintf("dotenv: failed to load env files: %v", err))
	}
}

// load is the internal implementation for Load and Overload
func load(overload bool, filenames ...string) error {
	env, err := Read(filenames...)
	if err != nil {
		return err
	}

	for key, value := range env {
		if overload || os.Getenv(key) == "" {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("failed to set environment variable %s: %w", key, err)
			}
		}
	}

	return nil
}

// readFile reads a single .env file and returns the parsed environment variables
func readFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	return Parse(file)
}

// formatEnvLine formats a key-value pair for .env file output
func formatEnvLine(key, value string) string {
	// Simple values that don't need quoting
	if !needsQuoting(value) {
		return fmt.Sprintf("%s=%s", key, value)
	}

	// Quote and escape the value
	escaped := escapeValue(value)
	return fmt.Sprintf(`%s="%s"`, key, escaped)
}

// needsQuoting determines if a value needs to be quoted
func needsQuoting(value string) bool {
	if value == "" {
		return true
	}

	for _, char := range value {
		switch char {
		case ' ', '\t', '\n', '\r', '"', '\'', '\\', '#', '$':
			return true
		}
	}

	return false
}

// escapeValue escapes special characters in a value for double-quoted output
func escapeValue(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	value = strings.ReplaceAll(value, "\n", `\n`)
	value = strings.ReplaceAll(value, "\r", `\r`)
	value = strings.ReplaceAll(value, "\t", `\t`)
	return value
}
