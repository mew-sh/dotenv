// Command dotenv loads environment variables from .env files and executes commands
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/mew-sh/dotenv"
)

var (
	envFiles    = flag.String("f", "", "comma separated paths to .env files")
	overload    = flag.Bool("o", false, "override existing environment variables")
	showHelp    = flag.Bool("h", false, "show help")
	showVersion = flag.Bool("v", false, "show version")
)

const version = "v2.0.0"

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("dotenv %s\n", version)
		return
	}

	if *showHelp || flag.NArg() == 0 {
		showUsage()
		return
	}

	// Load environment files
	var files []string
	if *envFiles != "" {
		files = strings.Split(*envFiles, ",")
		// Trim whitespace from file names
		for i, file := range files {
			files[i] = strings.TrimSpace(file)
		}
	}

	var err error
	if *overload {
		err = dotenv.Overload(files...)
	} else {
		err = dotenv.Load(files...)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env files: %v\n", err)
		os.Exit(1)
	}

	// Execute the command
	args := flag.Args()
	cmd := args[0]
	cmdArgs := args[1:]

	// Look for the command in PATH
	cmdPath, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Command not found: %s\n", cmd)
		os.Exit(127)
	}

	// Execute the command with the loaded environment
	err = syscall.Exec(cmdPath, append([]string{cmd}, cmdArgs...), os.Environ())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute command: %v\n", err)
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Printf(`dotenv %s - Load environment variables from .env files and execute commands

Usage:
  dotenv [options] COMMAND [ARGS...]

Options:
  -f FILE       comma separated paths to .env files (default: .env)
  -o            override existing environment variables
  -h            show this help message
  -v            show version

Examples:
  # Load .env and run a command
  dotenv go run main.go

  # Load specific .env files
  dotenv -f .env.local,.env.production npm start

  # Override existing environment variables
  dotenv -o -f .env.override python app.py

  # Load from multiple files (later files take precedence)
  dotenv -f .env,.env.local,.env.development rails server

Environment Files:
  If no -f flag is provided, dotenv will attempt to load .env from the current directory.
  Multiple files can be specified with comma separation.
  Files are loaded in order, with later files taking precedence for duplicate keys.

Exit Codes:
  0    Command executed successfully
  1    Error loading .env files or executing command
  127  Command not found

For more information, visit: https://github.com/mew-sh/dotenv
`, version)
}
