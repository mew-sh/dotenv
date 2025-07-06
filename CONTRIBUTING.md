# Contributing to dotenv

We welcome contributions to the dotenv library! This document provides guidelines for contributing.

## Code of Conduct

This project follows the [Go Community Code of Conduct](https://golang.org/conduct).

## How to Contribute

### Reporting Issues

Before creating an issue, please:

1. **Search existing issues** to see if the problem has already been reported
2. **Check the documentation** to ensure you're using the library correctly
3. **Create a minimal reproduction** of the issue

When creating an issue, please include:
- Go version
- Operating system
- Clear description of the problem
- Minimal code example that reproduces the issue
- Expected vs actual behavior

### Submitting Pull Requests

1. **Fork the repository** and create a new branch from `main`
2. **Make your changes** with clear, descriptive commit messages
3. **Add tests** for any new functionality
4. **Update documentation** if needed
5. **Ensure all tests pass** by running `make test`
6. **Run the linter** with `make lint`
7. **Submit a pull request** with a clear description

## Development Setup

### Prerequisites

- Go 1.21 or later
- Make (optional, but recommended)

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/dotenv.git
cd dotenv

# Install development dependencies
make deps

# Run tests
make test

# Run all quality checks
make check
```

### Development Workflow

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** and add tests

3. **Run tests frequently**:
   ```bash
   make test
   ```

4. **Check formatting and linting**:
   ```bash
   make fmt
   make lint
   ```

5. **Run benchmarks** if you've made performance changes:
   ```bash
   make benchmark
   ```

6. **Update documentation** if needed

## Code Guidelines

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use [gofmt](https://golang.org/cmd/gofmt/) for formatting
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Testing

- **Write tests** for all new functionality
- **Maintain test coverage** at 100%
- **Use table-driven tests** where appropriate
- **Test error cases** and edge conditions
- **Add benchmarks** for performance-critical code

Example test structure:
```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "KEY=value",
            expected: "value",
            wantErr:  false,
        },
        // Add more test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ParseLine(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Documentation

- **Add godoc comments** for all public functions and types
- **Include examples** in documentation when helpful
- **Update README.md** for new features
- **Add usage examples** to EXAMPLES.md for complex features

### Compatibility

- **Maintain backward compatibility** whenever possible
- **Follow semantic versioning** for releases
- **Test against multiple Go versions** (see CI configuration)

## Architecture Guidelines

### Parser Design

The library uses a modular parser design:

- `dotenv.go` - Main API and high-level functions
- `parser.go` - Core parsing logic and utilities
- `autoload/` - Auto-loading functionality
- `cmd/dotenv/` - Command-line tool

### Performance Considerations

- **Minimize allocations** in hot paths
- **Use efficient string operations**
- **Benchmark performance-critical changes**
- **Consider memory usage** for large files

### Error Handling

- **Return descriptive errors** with context
- **Include line numbers** for parse errors
- **Use error wrapping** with `fmt.Errorf`
- **Validate inputs** early

## Feature Acceptance Criteria

We accept contributions that:

### ✅ Accepted Changes

- **Bug fixes** for existing functionality
- **Performance improvements** with benchmarks
- **Better error messages** and debugging info
- **Additional helper functions** that follow existing patterns
- **Documentation improvements** and examples
- **Test coverage improvements**
- **Compatibility improvements** with other dotenv implementations

### ❌ Generally Not Accepted

- **Breaking API changes** without strong justification
- **Features that add external dependencies**
- **Platform-specific code** without cross-platform alternatives
- **Features that significantly complicate the codebase**
- **Changes that reduce performance** without clear benefits

### 🤔 Needs Discussion

- **New major features** - please open an issue first
- **API additions** - should fit the existing design
- **Large refactoring** - discuss the approach first

## Release Process

1. **Update CHANGELOG.md** with new features and fixes
2. **Update version numbers** in relevant files
3. **Create a git tag** following semantic versioning
4. **GitHub Actions** will automatically build and test
5. **Release notes** will be generated from the changelog

## Getting Help

- **Documentation**: Check README.md and EXAMPLES.md
- **Issues**: Search existing issues or create a new one
- **Discussions**: Use GitHub Discussions for questions
- **Code**: Look at existing code and tests for examples

Thank you for contributing to dotenv! 🎉
