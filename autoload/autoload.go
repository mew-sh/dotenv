// Package autoload automatically loads environment variables from .env file
// when imported.
//
// Usage:
//
//	import _ "github.com/mew-sh/dotenv/autoload"
//
// This will automatically load the .env file from the current directory
// when the package is imported.
package autoload

import "github.com/mew-sh/dotenv"

func init() {
	// Silently load .env file, ignoring errors
	// This follows the convention that autoload should not fail
	// if the .env file doesn't exist
	_ = dotenv.Load()
}
