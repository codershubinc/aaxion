package aShare

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func ShowIncomingShareNotification(senderName, fileName, fileSizeStr, ext string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	title := fmt.Sprintf("⚡️ Incoming Share: %s", senderName)
	body := fmt.Sprintf("File: %s\nSize: %s | Type: %s", fileName, fileSizeStr, ext)

	cmd := exec.CommandContext(
		ctx,
		"notify-send",
		"-a", "Aaxion Share",
		"-u", "critical", // Keeps notification visible on KDE until acted upon
		"-t", "60000", // 60-second expiration timeout
		"-i", "folder-download", // Standard stock icon
		"-A", "accept=Accept", // Action button 1
		"-A", "decline=Decline", // Action button 2
		title,
		body,
	)
	cmd.Env = os.Environ()

	out, err := cmd.Output()
	if err != nil {
		return "timeout", err
	}

	result := strings.TrimSpace(string(out))
	switch result {
	case "accept":
		return "accepted", nil
	case "decline":
		return "declined", nil
	default:
		return "dismissed", nil
	}
}

func ShowTransferCompleteNotification(fileName, saveDir string) {
	go func() {
		cmd := exec.Command("notify-send",
			"-a", "Aaxion Share",
			"-i", "document-save",
			"-t", "8000",
			"-A", "open=Open Folder",
			"✅ Transfer Complete",
			fmt.Sprintf("Received: %s\nSaved in: %s", fileName, saveDir),
		)
		out, err := cmd.Output()
		if err == nil && strings.TrimSpace(string(out)) == "open" {
			// Opens Dolphin / default file manager
			exec.Command("xdg-open", saveDir).Start()
		}
	}()
}
