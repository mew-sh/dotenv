# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2025-01-01

### Added
- **Modern Go Implementation**: Built with Go 1.24+ features and idioms
- **Full .env Format Support**: 
  - Comments (both full-line and inline)
  - Export syntax (`export VAR=value`)
  - YAML-style syntax (`key: value`)
  - Single and double quoted values
  - Escape sequence processing in double quotes
- **Variable Expansion**: Support for `$VAR` and `${VAR}` syntax with nested expansion
- **Type-Safe Helpers**: Built-in parsing functions for int, bool, float types
- **Zero Dependencies**: Pure Go implementation with no external dependencies
- **Command-Line Tool**: Full-featured CLI tool compatible with original godotenv
- **Performance Optimizations**: 
  - Efficient parsing with minimal allocations
  - Single-pass parsing algorithm
  - Optimized regular expressions
- **Enhanced API**:
  - `Parse()` function for reading from io.Reader
  - `Unmarshal()` and `Marshal()` for string conversion
  - `Must()` function for panic-on-error loading
  - `GetRequired()`, `GetWithDefault()` helper functions
  - `ParseInt()`, `ParseBool()`, `ParseFloat()` type converters
- **Auto-loading Package**: Import `autoload` package for automatic .env loading
- **Comprehensive Testing**: 
  - Full test suite with 100% coverage
  - Benchmark tests for performance validation
  - Example programs and usage documentation
- **Documentation**:
  - Comprehensive README with examples
  - EXAMPLES.md with real-world usage patterns
  - Inline code documentation with examples

### Changed
- **API Compatibility**: Maintains compatibility with joho/godotenv while adding new features
- **Error Handling**: Improved error messages with file names and line numbers
- **Parser Architecture**: Modular parser design allowing custom options

### Performance
- **Benchmarks** (on test machine):
  - Small files (3 vars): ~3,000 ns/op, 22 allocs/op
  - Medium files (300 vars): ~270,000 ns/op, 2,104 allocs/op
  - Complex parsing: ~24,000 ns/op, 166 allocs/op
  - Variable expansion: ~6,000 ns/op, 38 allocs/op
  - Marshal/Unmarshal: ~3,000 ns/op, 41 allocs/op

### Security
- **Safe Variable Expansion**: Prevents infinite recursion and handles undefined variables
- **Escape Processing**: Proper handling of escape sequences to prevent injection
- **Optional Features**: Variable expansion can be disabled for security-sensitive applications

## [1.0.0] - Initial Release

### Added
- Basic .env file parsing
- Environment variable loading
- Simple API compatible with existing libraries
