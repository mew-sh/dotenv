# DotEnv

[![Go Report Card](https://goreportcard.com/badge/github.com/mew-sh/dotenv)](https://goreportcard.com/report/github.com/mew-sh/dotenv)
[![GoDoc](https://godoc.org/github.com/mew-sh/dotenv?status.svg)](https://godoc.org/github.com/mew-sh/dotenv)

A modern Go library for loading environment variables from `.env` files, inspired by the [joho/godotenv](https://github.com/joho/godotenv) library but built with the latest Go features and best practices.

## Features

- 🚀 **Modern Go**: Built with Go 1.24+ features and idioms
- 📝 **Full .env support**: Comments, exports, quotes, escape sequences
- 🔄 **Variable expansion**: Support for `$VAR` and `${VAR}` syntax
- 🛡️ **Type-safe helpers**: Built-in parsing for int, bool, float
- 📦 **Zero dependencies**: Pure Go implementation
- 🎯 **Drop-in replacement**: Compatible API with existing libraries
- ⚡ **Performance focused**: Efficient parsing and minimal allocations

## Installation

```bash
go get github.com/mew-sh/dotenv
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/mew-sh/dotenv"
)

func main() {
    // Load .env file from current directory
    err := dotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    
    fmt.Println("Database URL:", os.Getenv("DATABASE_URL"))
}
```

### Auto-loading

For automatic loading during import:

```go
import _ "github.com/mew-sh/dotenv/autoload"
```

### Multiple Files

```go
// Load multiple .env files (later files take precedence)
err := dotenv.Load(".env.local", ".env.production", ".env")
```

### Reading Without Setting Environment

```go
// Read variables into a map without setting environment
env, err := dotenv.Read(".env")
if err != nil {
    log.Fatal(err)
}

fmt.Println("API Key:", env["API_KEY"])
```

### Parsing from String or Reader

```go
// From string
env, err := dotenv.Unmarshal("KEY=value\nANOTHER=value2")

// From io.Reader
file, _ := os.Open(".env")
env, err := dotenv.Parse(file)
```

### Force Override

```go
// Override existing environment variables
err := dotenv.Overload(".env")
```

### Type-Safe Parsing

```go
// Parse environment variables with type safety
port := dotenv.ParseInt("PORT", 8080)           // Default: 8080
debug := dotenv.ParseBool("DEBUG", false)       // Default: false
timeout := dotenv.ParseFloat("TIMEOUT", 30.0)   // Default: 30.0

// Required variables (panics if not set)
apiKey := dotenv.GetRequired("API_KEY")

// With default values
dbHost := dotenv.GetWithDefault("DB_HOST", "localhost")
```

### Writing .env Files

```go
env := map[string]string{
    "API_KEY": "secret-key",
    "DEBUG":   "true",
    "PORT":    "8080",
}

// Write to file
err := dotenv.Write(env, ".env.production")

// Or get as string
content, err := dotenv.Marshal(env)
```

## .env File Format

### Basic Variables

```bash
# Comments are supported
API_KEY=your-secret-key
DEBUG=true
PORT=8080

# Empty values
EMPTY_VAR=
```

### Quoted Values

```bash
# Double quotes (with escape sequence support)
MESSAGE="Hello\nWorld"
PATH="/usr/local/bin:/usr/bin"

# Single quotes (literal values, no escaping)
LITERAL='$HOME will not be expanded'
```

### Export Syntax

```bash
export NODE_ENV=production
export PATH="/usr/local/bin:$PATH"
```

### YAML-Style Syntax

```bash
database_url: postgres://user:pass@localhost/db
redis_url: redis://localhost:6379
```

### Variable Expansion

```bash
# Basic expansion
HOME_DIR=/home/user
CONFIG_FILE=${HOME_DIR}/config.json

# Environment variable expansion
PATH_EXTENDED=${PATH}:/usr/local/bin

# Nested expansion
BASE_URL=https://api.example.com
API_ENDPOINT=${BASE_URL}/v1/users
```

### Comments

```bash
# Full line comments
API_KEY=secret  # Inline comments
DATABASE_URL="postgres://localhost/myapp"  # Comments after quotes
```

### Escape Sequences (in double quotes)

```bash
MULTILINE="Line 1\nLine 2\nLine 3"
TAB_SEPARATED="Column1\tColumn2\tColumn3"
QUOTED_STRING="He said \"Hello World\""
BACKSLASH="Path\\to\\file"
```

## Advanced Usage

### Custom Parser Options

```go
// Create parser with custom options
parser := dotenv.NewParserWithOptions(false) // Disable variable expansion
env, err := parser.Parse(reader)
```

### Panic on Missing .env

```go
// Will panic if .env file cannot be loaded
dotenv.Must()  // Loads .env
dotenv.Must(".env.production")  // Loads specific file
```

### Environment-Specific Loading

```go
env := os.Getenv("APP_ENV")
if env == "" {
    env = "development"
}

// Load environment-specific files
dotenv.Load(".env." + env + ".local")
if env != "test" {
    dotenv.Load(".env.local")
}
dotenv.Load(".env." + env)
dotenv.Load(".env")  // Base .env file
```

## Error Handling

The library provides detailed error messages for common issues:

```go
env, err := dotenv.Read("config.env")
if err != nil {
    // Error includes file name and line number for parse errors
    fmt.Printf("Failed to load config: %v\n", err)
}
```

## Performance

This library is designed for performance:

- Efficient string parsing with minimal allocations
- Lazy variable expansion
- Optimized regular expressions
- Single-pass parsing

## Compatibility

This library aims to be compatible with:
- [joho/godotenv](https://github.com/joho/godotenv)
- [Ruby dotenv](https://github.com/bkeepers/dotenv)
- [Node.js dotenv](https://github.com/motdotla/dotenv)

## API Reference

### Loading Functions

- `Load(filenames ...string) error` - Load .env files into environment
- `Overload(filenames ...string) error` - Load and override existing variables
- `Must(filenames ...string)` - Load with panic on error

### Reading Functions

- `Read(filenames ...string) (map[string]string, error)` - Read without setting environment
- `Parse(reader io.Reader) (map[string]string, error)` - Parse from reader
- `Unmarshal(data string) (map[string]string, error)` - Parse from string

### Writing Functions

- `Marshal(env map[string]string) (string, error)` - Convert map to .env format
- `Write(env map[string]string, filename string) error` - Write map to file

### Type-Safe Helpers

- `ParseInt(key string, defaultValue int) int`
- `ParseBool(key string, defaultValue bool) bool`
- `ParseFloat(key string, defaultValue float64) float64`
- `GetRequired(key string) string`
- `GetWithDefault(key, defaultValue string) string`

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by [joho/godotenv](https://github.com/joho/godotenv)
- Compatible with [Ruby dotenv](https://github.com/bkeepers/dotenv) conventions
- Following [twelve-factor app](https://12factor.net/) methodology
