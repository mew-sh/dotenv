package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mew-sh/dotenv"
)

func main() {
	fmt.Println("=== DotEnv Library Example ===")
	fmt.Println()

	// Example 1: Basic loading
	fmt.Println("1. Basic Loading:")
	createExampleEnvFile()

	err := dotenv.Load("example.env")
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	} else {
		fmt.Printf("   API_KEY: %s\n", os.Getenv("API_KEY"))
		fmt.Printf("   DEBUG: %s\n", os.Getenv("DEBUG"))
		fmt.Printf("   PORT: %s\n", os.Getenv("PORT"))
	}

	// Example 2: Type-safe parsing
	fmt.Println("\n2. Type-Safe Parsing:")
	port := dotenv.ParseInt("PORT", 3000)
	debug := dotenv.ParseBool("DEBUG", false)
	timeout := dotenv.ParseFloat("TIMEOUT", 30.0)

	fmt.Printf("   Port (int): %d\n", port)
	fmt.Printf("   Debug (bool): %t\n", debug)
	fmt.Printf("   Timeout (float): %.1f\n", timeout)

	// Example 3: Reading without setting environment
	fmt.Println("\n3. Reading to Map:")
	env, err := dotenv.Read("example.env")
	if err != nil {
		log.Printf("Error reading file: %v", err)
	} else {
		fmt.Printf("   Read %d variables\n", len(env))
		for key, value := range env {
			fmt.Printf("   %s = %s\n", key, value)
		}
	}

	// Example 4: Parsing from string
	fmt.Println("\n4. Parsing from String:")
	envString := `NEW_VAR=from_string
ANOTHER_VAR="quoted value"
EXPANDED_VAR=${NEW_VAR}_expanded`

	parsed, err := dotenv.Unmarshal(envString)
	if err != nil {
		log.Printf("Error parsing string: %v", err)
	} else {
		for key, value := range parsed {
			fmt.Printf("   %s = %s\n", key, value)
		}
	}

	// Example 5: Writing env file
	fmt.Println("\n5. Writing .env File:")
	newEnv := map[string]string{
		"CREATED_BY":    "example",
		"TIMESTAMP":     "2024-01-01",
		"COMPLEX_VALUE": "value with \"quotes\" and\nnewlines",
	}

	content, err := dotenv.Marshal(newEnv)
	if err != nil {
		log.Printf("Error marshaling: %v", err)
	} else {
		fmt.Printf("   Generated content:\n%s\n", content)

		err = dotenv.Write(newEnv, "output.env")
		if err != nil {
			log.Printf("Error writing file: %v", err)
		} else {
			fmt.Println("   Written to output.env")
		}
	}

	// Example 6: Variable expansion
	fmt.Println("\n6. Variable Expansion:")
	expansionExample := `BASE_URL=https://api.example.com
API_V1=${BASE_URL}/v1
API_V2=${BASE_URL}/v2
USERS_ENDPOINT=${API_V1}/users`

	expanded, err := dotenv.Unmarshal(expansionExample)
	if err != nil {
		log.Printf("Error parsing expansion example: %v", err)
	} else {
		for key, value := range expanded {
			fmt.Printf("   %s = %s\n", key, value)
		}
	}

	// Cleanup
	os.Remove("example.env")
	os.Remove("output.env")
}

func createExampleEnvFile() {
	content := `# Example .env file
API_KEY=secret-key-123
DEBUG=true
PORT=8080
TIMEOUT=30.5

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_URL="postgres://${DB_HOST}:${DB_PORT}/${DB_NAME}"

# Features
ENABLE_FEATURE_X=yes
ENABLE_FEATURE_Y=no

# Quoted values
APP_NAME="My Awesome App"
WELCOME_MESSAGE="Welcome to \"My App\"\nEnjoy your stay!"

# Empty value
OPTIONAL_CONFIG=
`

	err := os.WriteFile("example.env", []byte(content), 0644)
	if err != nil {
		log.Printf("Error creating example file: %v", err)
	}
}
