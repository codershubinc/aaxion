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
		Description: "Create a new admin user (e.g. aax create-admin username:password)",
		Run: func(args []string) {
			if len(args) == 0 {
				fmt.Println("❌ Error: Missing credentials. Usage: aax create-admin username:password")
				os.Exit(1)
			}
			
			err := db.InitDb()
			if err != nil {
				log.Fatal("❌ Got err initializing DB:", err)
			}

			auth.HandleCreateAdminCLI(args[0])
		},
	})
}
