package auth

import (
	"aaxion/internal/db"
	"fmt"
	"os"
	"strings"
)

// HandleCreateAdminCLI handles the creation of an admin user from the command line interface
func HandleCreateAdminCLI(creds string, force bool) {
	parts := strings.SplitN(creds, ":", 2)
	if len(parts) != 2 {
		fmt.Println("❌ Error: Invalid format. Please use: aax create-admin username:password")
		os.Exit(1)
	}

	username := parts[0]
	password := parts[1]

	exists, err := db.HasUsers()
	if err != nil {
		fmt.Printf("❌ Error checking database: %v\n", err)
		os.Exit(1)
	}

	if exists {
		if !force {
			fmt.Println("❌ Error: A user is already registered in the database.")
			fmt.Println("   If you want to overwrite the existing admin, append --force")
			fmt.Println("   Example: aax create-admin new_user:new_pass --force")
			os.Exit(1)
		} else {
			fmt.Println("⚠️  Force flag detected. Deleting existing admin users...")
			err = db.DeleteAllUsers()
			if err != nil {
				fmt.Printf("❌ Error clearing old users: %v\n", err)
				os.Exit(1)
			}
		}
	}

	err = db.CreateUser(username, password)
	if err != nil {
		fmt.Printf("❌ Error creating user: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Admin user '%s' created successfully!\n", username)
	fmt.Println("You can now start the server normally by running 'aax serve'")
}
