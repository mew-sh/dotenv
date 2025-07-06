# Usage Examples

This document provides comprehensive examples of using the dotenv library.

## Basic Examples

### 1. Simple Loading

```go
package main

import (
    "fmt"
    "os"
    "github.com/mew-sh/dotenv"
)

func main() {
    // Load .env file
    err := dotenv.Load()
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Database URL:", os.Getenv("DATABASE_URL"))
}
```

### 2. Multiple Environment Files

```go
// Load environment-specific configuration
env := os.Getenv("APP_ENV")
if env == "" {
    env = "development"
}

// Load in order of precedence (later files override earlier ones)
dotenv.Load(".env")                    // Base configuration
dotenv.Load(".env." + env)             // Environment-specific
dotenv.Load(".env." + env + ".local")  // Local overrides
dotenv.Load(".env.local")              // Global local overrides
```

## Configuration Management

### 3. Structured Configuration

```go
type Config struct {
    DatabaseURL    string
    RedisURL       string
    Port           int
    Debug          bool
    JWTSecret      string
    SMTPSettings   SMTPConfig
}

type SMTPConfig struct {
    Host     string
    Port     int
    Username string
    Password string
}

func LoadConfig() (*Config, error) {
    if err := dotenv.Load(); err != nil {
        return nil, err
    }
    
    return &Config{
        DatabaseURL:  dotenv.GetRequired("DATABASE_URL"),
        RedisURL:     dotenv.GetWithDefault("REDIS_URL", "redis://localhost:6379"),
        Port:         dotenv.ParseInt("PORT", 8080),
        Debug:        dotenv.ParseBool("DEBUG", false),
        JWTSecret:    dotenv.GetRequired("JWT_SECRET"),
        SMTPSettings: SMTPConfig{
            Host:     dotenv.GetRequired("SMTP_HOST"),
            Port:     dotenv.ParseInt("SMTP_PORT", 587),
            Username: dotenv.GetRequired("SMTP_USERNAME"),
            Password: dotenv.GetRequired("SMTP_PASSWORD"),
        },
    }, nil
}
```

### 4. Environment-aware Configuration

```go
func LoadEnvironmentConfig() error {
    environment := os.Getenv("GO_ENV")
    if environment == "" {
        environment = "development"
    }
    
    // Always load base configuration first
    if err := dotenv.Load(".env"); err != nil {
        // .env might not exist, that's okay
        fmt.Printf("No .env file found: %v\n", err)
    }
    
    // Load environment-specific configuration
    envFile := fmt.Sprintf(".env.%s", environment)
    if err := dotenv.Load(envFile); err != nil {
        return fmt.Errorf("failed to load %s: %w", envFile, err)
    }
    
    // Load local overrides (git-ignored)
    localFile := fmt.Sprintf(".env.%s.local", environment)
    if err := dotenv.Load(localFile); err != nil {
        // Local files are optional
        fmt.Printf("No local override file %s: %v\n", localFile, err)
    }
    
    return nil
}
```

## Advanced Usage

### 5. Custom Parser with Options

```go
// Disable variable expansion for security
parser := dotenv.NewParserWithOptions(false)
env, err := parser.Parse(configReader)
if err != nil {
    return err
}

// Manually validate and expand only trusted variables
for key, value := range env {
    if strings.Contains(value, "$") {
        // Log potential security issue
        log.Warnf("Variable expansion detected in %s: %s", key, value)
    }
    os.Setenv(key, value)
}
```

### 6. Reading Configuration Without Setting Environment

```go
// Read configuration into memory without affecting environment
env, err := dotenv.Read(".env", ".env.production")
if err != nil {
    return err
}

// Use configuration directly
databaseURL := env["DATABASE_URL"]
if databaseURL == "" {
    return errors.New("DATABASE_URL is required")
}

// Selectively set only certain variables
sensitiveVars := []string{"API_KEY", "JWT_SECRET", "DATABASE_PASSWORD"}
for _, key := range sensitiveVars {
    if value, exists := env[key]; exists {
        os.Setenv(key, value)
    }
}
```

### 7. Dynamic Configuration Loading

```go
func LoadDynamicConfig() error {
    // Get configuration sources from environment
    configSources := os.Getenv("CONFIG_SOURCES")
    if configSources == "" {
        configSources = ".env"
    }
    
    files := strings.Split(configSources, ",")
    for i, file := range files {
        files[i] = strings.TrimSpace(file)
    }
    
    return dotenv.Load(files...)
}
```

## Web Application Examples

### 8. HTTP Server Configuration

```go
type ServerConfig struct {
    Host            string
    Port            int
    ReadTimeout     time.Duration
    WriteTimeout    time.Duration
    ShutdownTimeout time.Duration
    TLSCert         string
    TLSKey          string
}

func LoadServerConfig() (*ServerConfig, error) {
    dotenv.Must() // Panic if .env cannot be loaded
    
    return &ServerConfig{
        Host:            dotenv.GetWithDefault("HTTP_HOST", "0.0.0.0"),
        Port:            dotenv.ParseInt("HTTP_PORT", 8080),
        ReadTimeout:     time.Duration(dotenv.ParseInt("HTTP_READ_TIMEOUT", 30)) * time.Second,
        WriteTimeout:    time.Duration(dotenv.ParseInt("HTTP_WRITE_TIMEOUT", 30)) * time.Second,
        ShutdownTimeout: time.Duration(dotenv.ParseInt("HTTP_SHUTDOWN_TIMEOUT", 10)) * time.Second,
        TLSCert:         os.Getenv("TLS_CERT_FILE"),
        TLSKey:          os.Getenv("TLS_KEY_FILE"),
    }, nil
}

func (c *ServerConfig) Address() string {
    return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *ServerConfig) IsTLSEnabled() bool {
    return c.TLSCert != "" && c.TLSKey != ""
}
```

### 9. Database Configuration with Connection Pooling

```go
type DatabaseConfig struct {
    Host            string
    Port            int
    Name            string
    Username        string
    Password        string
    SSLMode         string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
    if err := dotenv.Load(); err != nil {
        return nil, err
    }
    
    config := &DatabaseConfig{
        Host:            dotenv.GetWithDefault("DB_HOST", "localhost"),
        Port:            dotenv.ParseInt("DB_PORT", 5432),
        Name:            dotenv.GetRequired("DB_NAME"),
        Username:        dotenv.GetRequired("DB_USERNAME"),
        Password:        dotenv.GetRequired("DB_PASSWORD"),
        SSLMode:         dotenv.GetWithDefault("DB_SSL_MODE", "prefer"),
        MaxOpenConns:    dotenv.ParseInt("DB_MAX_OPEN_CONNS", 25),
        MaxIdleConns:    dotenv.ParseInt("DB_MAX_IDLE_CONNS", 5),
        ConnMaxLifetime: time.Duration(dotenv.ParseInt("DB_CONN_MAX_LIFETIME", 300)) * time.Second,
    }
    
    return config, nil
}

func (c *DatabaseConfig) ConnectionString() string {
    return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
        url.QueryEscape(c.Username),
        url.QueryEscape(c.Password),
        c.Host,
        c.Port,
        c.Name,
        c.SSLMode,
    )
}
```

## Testing Examples

### 10. Test Environment Configuration

```go
func setupTestEnvironment(t *testing.T) {
    // Save current environment
    originalEnv := os.Environ()
    
    // Clean environment for testing
    os.Clearenv()
    
    // Load test configuration
    err := dotenv.Load("testdata/.env.test")
    if err != nil {
        t.Fatalf("Failed to load test environment: %v", err)
    }
    
    // Cleanup function
    t.Cleanup(func() {
        os.Clearenv()
        for _, env := range originalEnv {
            parts := strings.SplitN(env, "=", 2)
            if len(parts) == 2 {
                os.Setenv(parts[0], parts[1])
            }
        }
    })
}

func TestWithEnvironment(t *testing.T) {
    setupTestEnvironment(t)
    
    // Test your application with controlled environment
    config, err := LoadConfig()
    if err != nil {
        t.Fatalf("Failed to load config: %v", err)
    }
    
    assert.Equal(t, "test", config.Environment)
    assert.Equal(t, "localhost:5432", config.DatabaseURL)
}
```

### 11. Mocking Environment Variables

```go
func TestConfigurationLoading(t *testing.T) {
    tests := []struct {
        name     string
        envVars  map[string]string
        expected Config
        wantErr  bool
    }{
        {
            name: "production config",
            envVars: map[string]string{
                "DATABASE_URL": "postgres://prod:secret@db.example.com/app",
                "REDIS_URL":    "redis://redis.example.com:6379",
                "DEBUG":        "false",
                "PORT":         "8080",
            },
            expected: Config{
                DatabaseURL: "postgres://prod:secret@db.example.com/app",
                RedisURL:    "redis://redis.example.com:6379",
                Debug:       false,
                Port:        8080,
            },
        },
        {
            name: "missing required variable",
            envVars: map[string]string{
                "DEBUG": "true",
                "PORT":  "3000",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create temporary .env file
            envContent := ""
            for key, value := range tt.envVars {
                envContent += fmt.Sprintf("%s=%s\n", key, value)
            }
            
            tmpFile := filepath.Join(t.TempDir(), ".env")
            err := os.WriteFile(tmpFile, []byte(envContent), 0644)
            require.NoError(t, err)
            
            // Test configuration loading
            config, err := LoadConfigFromFile(tmpFile)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            require.NoError(t, err)
            assert.Equal(t, tt.expected, *config)
        })
    }
}
```

## Command Line Tool Examples

### 12. Using the CLI Tool

```bash
# Basic usage - load .env and run command
dotenv go run main.go

# Load specific files
dotenv -f .env.production,secrets.env npm start

# Override existing environment variables
dotenv -o -f override.env python manage.py runserver

# Multiple files with precedence
dotenv -f .env,.env.local,.env.development rails server

# Run with different environments
dotenv -f .env.test go test ./...
```

### 13. Integration with Scripts

```bash
#!/bin/bash
# deploy.sh

set -e

echo "Deploying application..."

# Load production environment and run deployment
dotenv -f .env.production -f .env.secrets docker-compose up -d

# Run database migrations with production config
dotenv -f .env.production go run migrations/migrate.go

echo "Deployment complete!"
```

## Error Handling Patterns

### 14. Graceful Error Handling

```go
func LoadConfigurationGracefully() *Config {
    config := &Config{
        // Set reasonable defaults
        Port:  8080,
        Debug: false,
        Host:  "localhost",
    }
    
    // Try to load .env file, but don't fail if it doesn't exist
    if err := dotenv.Load(); err != nil {
        log.Printf("No .env file found, using defaults: %v", err)
    }
    
    // Override defaults with environment variables
    if port := dotenv.ParseInt("PORT", 0); port > 0 {
        config.Port = port
    }
    
    if host := os.Getenv("HOST"); host != "" {
        config.Host = host
    }
    
    config.Debug = dotenv.ParseBool("DEBUG", config.Debug)
    
    return config
}
```

### 15. Validation and Error Reporting

```go
func ValidateConfiguration() error {
    var errors []string
    
    // Check required variables
    required := []string{"DATABASE_URL", "JWT_SECRET", "API_KEY"}
    for _, key := range required {
        if os.Getenv(key) == "" {
            errors = append(errors, fmt.Sprintf("%s is required but not set", key))
        }
    }
    
    // Validate specific formats
    if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
        if _, err := url.Parse(dbURL); err != nil {
            errors = append(errors, fmt.Sprintf("DATABASE_URL is not a valid URL: %v", err))
        }
    }
    
    // Validate numeric ranges
    if port := dotenv.ParseInt("PORT", 0); port < 1024 || port > 65535 {
        errors = append(errors, "PORT must be between 1024 and 65535")
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("configuration validation failed:\n  %s", strings.Join(errors, "\n  "))
    }
    
    return nil
}
```

These examples demonstrate various patterns for using the dotenv library in real-world applications, from simple configuration loading to complex multi-environment setups with proper error handling and validation.
