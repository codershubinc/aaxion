package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func init() {
	RegisterCommand(Command{
		Name:        "service",
		Description: "Manage Aaxion systemd background service (Linux only)",
		Run: func(args []string) {
			if runtime.GOOS != "linux" {
				fmt.Println("❌ The 'service' command is only supported on Linux.")
				os.Exit(1)
			}

			if len(args) == 0 {
				fmt.Println("⚡️ Aaxion Service Manager")
				fmt.Println("Usage: aax service [install|uninstall|start|stop|restart|status]")
				return
			}

			action := args[0]
			switch action {
			case "install":
				installService()
			case "uninstall":
				uninstallService()
			case "start", "stop", "status", "restart":
				runSystemctl(action)
			default:
				fmt.Println("❌ Unknown action:", action)
			}
		},
	})
}

func installService() {
	if os.Geteuid() != 0 {
		fmt.Println("❌ Sudo/root privileges are required to install the systemd service.")
		fmt.Println("   Run: sudo aax service install")
		os.Exit(1)
	}

	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("❌ Could not determine executable path:", err)
		os.Exit(1)
	}

	// Assuming the user wants it to run under their own user, not root, for security.
	// We'll get the real user if run via sudo, or fallback to root.
	user := os.Getenv("SUDO_USER")
	if user == "" {
		user = "root"
	}

	// Look up the user's home directory to use as the working directory.
	// We need this so it doesn't try to create .aaxion.db in the root / folder.
	homeDir := "/root"
	if user != "root" {
		homeDir = "/home/" + user
	}

	serviceContent := fmt.Sprintf(`[Unit]
Description=Aaxion File Streaming Server
After=network.target

[Service]
Type=simple
User=%s
WorkingDirectory=%s
ExecStart=%s serve
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
`, user, homeDir, execPath)

	servicePath := "/etc/systemd/system/aaxion.service"
	err = os.WriteFile(servicePath, []byte(serviceContent), 0644)
	if err != nil {
		fmt.Println("❌ Failed to write service file:", err)
		os.Exit(1)
	}

	fmt.Println("✅ Created service file at", servicePath)
	
	runCommand("systemctl", "daemon-reload")
	runCommand("systemctl", "enable", "aaxion.service")
	runCommand("systemctl", "start", "aaxion.service")
	
	fmt.Println("✅ Aaxion background daemon is now running and enabled on boot.")
	fmt.Println("\nTo check the status of the server, run:")
	fmt.Println("   aax service status")
}

func uninstallService() {
	if os.Geteuid() != 0 {
		fmt.Println("❌ Sudo/root privileges are required to uninstall the systemd service.")
		fmt.Println("   Run: sudo aax service uninstall")
		os.Exit(1)
	}

	runCommand("systemctl", "stop", "aaxion.service")
	runCommand("systemctl", "disable", "aaxion.service")
	
	err := os.Remove("/etc/systemd/system/aaxion.service")
	if err != nil {
		fmt.Println("❌ Failed to remove service file (it might not exist).")
	} else {
		fmt.Println("✅ Removed /etc/systemd/system/aaxion.service")
	}
	
	runCommand("systemctl", "daemon-reload")
	fmt.Println("✅ Aaxion service completely uninstalled.")
}

func runSystemctl(action string) {
	// Status can be run without root, but start/stop requires it
	if action != "status" && os.Geteuid() != 0 {
		fmt.Println("❌ Sudo/root privileges are required to", action, "the service.")
		fmt.Printf("   Run: sudo aax service %s\n", action)
		os.Exit(1)
	}

	cmd := exec.Command("systemctl", action, "aaxion.service")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil && action != "status" {
		fmt.Println("❌ Failed to", action, "service.")
	}
}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Run() // intentionally ignoring output for background config steps
}
