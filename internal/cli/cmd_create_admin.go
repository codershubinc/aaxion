package cli

import (
	"aaxion/internal/auth"
	"aaxion/internal/db"
	"fmt"
	"log"
	"os"
)

func init() {
	RegisterCommand(Command{
		Name:        "create-admin",
		Description: "Create a new admin user (e.g. aax create-admin user:pass). Use --force to overwrite.",
		Run: func(args []string) {
			if len(args) == 0 {
				fmt.Println("❌ Error: Missing credentials. Usage: aax create-admin username:password [--force]")
				os.Exit(1)
			}
			
			force := false
			if len(args) > 1 && args[1] == "--force" {
				force = true
			}

			err := db.InitDb()
			if err != nil {
				log.Fatal("❌ Got err initializing DB:", err)
			}

			auth.HandleCreateAdminCLI(args[0], force)
		},
	})
}
