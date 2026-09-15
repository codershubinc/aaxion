package cli

import (
	"aaxion/internal/version"
	"fmt"
	"runtime"
)

func init() {
	RegisterCommand(Command{
		Name:        "version",
		Description: "Show Aaxion version information",
		Run: func(args []string) {
			printVersion()
		},
	})
}

func printVersion() {
	// ANSI Color Codes
	boldBlue := "\033[1;34m"
	boldCyan := "\033[1;36m"
	yellow := "\033[33m"
	gray := "\033[90m"
	reset := "\033[0m"

	fmt.Printf("%s%s\n", boldCyan, reset)

	fmt.Printf("  %sVersion:%s   %s\n", boldBlue, reset, version.Version)
	fmt.Printf("  %sCodename:%s  %s\"%s\"%s\n", boldBlue, reset, yellow, version.Codename, reset)
	fmt.Printf("  %sPlatform:%s  %s/%s\n", boldBlue, reset, runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  %sRuntime:%s   %s\n\n", boldBlue, reset, runtime.Version())

	fmt.Printf("  %s(c) 2026 CodersHub Inc. Licensed under AGPLv3.%s\n", gray, reset)
}
