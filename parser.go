package dotenv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Regular expressions for parsing
	lineRegex      = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_]*)\s*[=:]\s*(.*)$`)
	exportRegex    = regexp.MustCompile(`^\s*export\s+([A-Za-z_][A-Za-z0-9_]*)\s*[=:]\s*(.*)$`)
	expandVarRegex = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)
)

// Parser handles the parsing of .env file content
type Parser struct {
	// expandVars determines if variable expansion should be performed
	expandVars bool
	// env holds the currently parsed environment variables for expansion
	env map[string]string
}

// NewParser creates a new parser with default settings
func NewParser() *Parser {
	return &Parser{
		expandVars: true,
		env:        make(map[string]string),
	}
}

// NewParserWithOptions creates a parser with custom options
func NewParserWithOptions(expandVars bool) *Parser {
	return &Parser{
		expandVars: expandVars,
		env:        make(map[string]string),
	}
}

// Parse reads from an io.Reader and parses the .env content
func (p *Parser) Parse(reader io.Reader) (map[string]string, error) {
	result := make(map[string]string)
	p.env = result // For variable expansion

	scanner := bufio.NewScanner(reader)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Skip empty lines and comments
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := p.parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("parse error on line %d: %w", lineNumber, err)
		}

		if key != "" {
			if p.expandVars {
				value = p.expandVariables(value, result)
			}
			result[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	return result, nil
}

// parseLine parses a single line and returns key, value, and any error
func (p *Parser) parseLine(line string) (string, string, error) {
	// Remove inline comments (but not those inside quotes)
	line = p.removeInlineComment(line)

	// Handle export prefix
	if matches := exportRegex.FindStringSubmatch(line); matches != nil {
		key := matches[1]
		value := strings.TrimSpace(matches[2])
		parsedValue, err := p.parseValue(value)
		return key, parsedValue, err
	}

	// Handle regular key=value or key:value
	if matches := lineRegex.FindStringSubmatch(line); matches != nil {
		key := matches[1]
		value := strings.TrimSpace(matches[2])
		parsedValue, err := p.parseValue(value)
		return key, parsedValue, err
	}

	// If line doesn't match any pattern and isn't empty, it's an error
	if strings.TrimSpace(line) != "" {
		return "", "", fmt.Errorf("invalid line format: %q", line)
	}

	return "", "", nil
}

// parseValue parses a value, handling quotes and escaping
func (p *Parser) parseValue(value string) (string, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return "", nil
	}

	// Handle quoted values
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			quote := value[0]
			inner := value[1 : len(value)-1]

			if quote == '"' {
				// Double quotes: process escape sequences
				return p.unescapeDoubleQuoted(inner), nil
			} else {
				// Single quotes: literal value (no escape processing)
				return inner, nil
			}
		}
	}

	// Unquoted value - trim trailing whitespace and remove trailing comments
	return strings.TrimSpace(value), nil
}

// removeInlineComment removes inline comments while preserving those inside quotes
func (p *Parser) removeInlineComment(line string) string {
	inQuotes := false
	quoteChar := byte(0)

	for i := 0; i < len(line); i++ {
		char := line[i]

		if !inQuotes {
			if char == '"' || char == '\'' {
				inQuotes = true
				quoteChar = char
			} else if char == '#' {
				// Found unquoted comment, trim everything after
				return strings.TrimSpace(line[:i])
			}
		} else {
			if char == quoteChar {
				// Check if it's escaped
				if i > 0 && line[i-1] != '\\' {
					inQuotes = false
					quoteChar = 0
				}
			}
		}
	}

	return line
}

// unescapeDoubleQuoted processes escape sequences in double-quoted strings
func (p *Parser) unescapeDoubleQuoted(value string) string {
	result := strings.Builder{}

	for i := 0; i < len(value); i++ {
		if value[i] == '\\' && i+1 < len(value) {
			next := value[i+1]
			switch next {
			case 'n':
				result.WriteByte('\n')
			case 'r':
				result.WriteByte('\r')
			case 't':
				result.WriteByte('\t')
			case '\\':
				result.WriteByte('\\')
			case '"':
				result.WriteByte('"')
			case '\'':
				result.WriteByte('\'')
			default:
				// Unknown escape, keep the backslash
				result.WriteByte('\\')
				result.WriteByte(next)
			}
			i++ // Skip the next character
		} else {
			result.WriteByte(value[i])
		}
	}

	return result.String()
}

// expandVariables expands variable references in the format $VAR or ${VAR}
func (p *Parser) expandVariables(value string, env map[string]string) string {
	return expandVarRegex.ReplaceAllStringFunc(value, func(match string) string {
		var varName string

		if strings.HasPrefix(match, "${") && strings.HasSuffix(match, "}") {
			// ${VAR} format
			varName = match[2 : len(match)-1]
		} else if strings.HasPrefix(match, "$") {
			// $VAR format
			varName = match[1:]
		}

		// Look up in parsed env first, then in OS env
		if val, exists := env[varName]; exists {
			return val
		}

		if val, exists := os.LookupEnv(varName); exists {
			return val
		}

		// Variable not found, return empty string (bash behavior)
		return ""
	})
}

// ParseInt parses an environment variable as an integer
func ParseInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}

	return defaultValue
}

// ParseBool parses an environment variable as a boolean
// Recognizes: true, false, 1, 0, yes, no, on, off (case insensitive)
func ParseBool(key string, defaultValue bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if value == "" {
		return defaultValue
	}

	switch value {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		return defaultValue
	}
}

// ParseFloat parses an environment variable as a float64
func ParseFloat(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	if parsed, err := strconv.ParseFloat(value, 64); err == nil {
		return parsed
	}

	return defaultValue
}

// GetRequired gets an environment variable and panics if it's not set
func GetRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}

// GetWithDefault gets an environment variable with a default value
func GetWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
