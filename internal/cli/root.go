package cli

import (
	"fmt"
	"os"
)

// Command represents a CLI subcommand
type Command struct {
	Name        string
	Description string
	Run         func(args []string)
}

var commands = map[string]Command{}

// RegisterCommand adds a new command to the CLI
func RegisterCommand(cmd Command) {
	commands[cmd.Name] = cmd
}

// Execute parses the CLI arguments and runs the appropriate command
func Execute() {
	if len(os.Args) < 2 {
		// Default to running the server if no command is provided
		if serveCmd, ok := commands["serve"]; ok {
			serveCmd.Run(nil)
			return
		}
	}

	commandName := os.Args[1]

	// Handle help
	if commandName == "help" || commandName == "--help" || commandName == "-h" {
		printUsage()
		os.Exit(0)
	}

	// Handle version flags
	if commandName == "--version" || commandName == "-v" {
		commandName = "version"
	}

	// Route to command
	if cmd, exists := commands[commandName]; exists {
		cmd.Run(os.Args[2:])
	} else {
		fmt.Printf("❌ Unknown command: %s\n", commandName)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("⚡️ Aaxion CLI (aax)")
	fmt.Println("\nUsage:")
	fmt.Println("  aax <command> [arguments]")
	fmt.Println("\nCommands:")
	for name, cmd := range commands {
		fmt.Printf("  %-15s %s\n", name, cmd.Description)
	}
}
