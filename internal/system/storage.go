package system

import (
	"aaxion/internal/helpers"
	"bufio"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func GetSystemStorage(w http.ResponseWriter, r *http.Request) {
	// Get root filesystem storage info
	var stat syscall.Statfs_t
	rootPath := getRootPath()

	err := syscall.Statfs(rootPath, &stat)
	if err != nil {
		http.Error(w, "Failed to get storage information: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate storage metrics
	totalBytes := stat.Blocks * uint64(stat.Bsize)
	availableBytes := stat.Bavail * uint64(stat.Bsize)
	usedBytes := totalBytes - (stat.Bfree * uint64(stat.Bsize))

	// Get external storage devices
	externalDevices, err := getExternalStorageDevices()
	if err != nil {
		// Log error but don't fail the request
		externalDevices = []map[string]interface{}{}
	}

	response := map[string]interface{}{
		"total":            totalBytes,
		"used":             usedBytes,
		"available":        availableBytes,
		"usage_percentage": float64(usedBytes) / float64(totalBytes) * 100,
		"external_devices": externalDevices,
	}

	helpers.SetJSONResponce(w, response)
}

func GetSystemRootPath(w http.ResponseWriter, r *http.Request) {
	rootPath := getRootPath()
	helpers.SetJSONResponce(w, map[string]string{"root_path": rootPath})
}

func getRootPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(string(homeDir))
}

func getExternalStorageDevices() ([]map[string]interface{}, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var devices []map[string]interface{}
	scanner := bufio.NewScanner(file)

	// Common mount points for external storage
	externalPrefixes := []string{"/media/", "/mnt/", "/run/media/"}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 3 {
			continue
		}

		device := fields[0]
		mountPoint := fields[1]
		fsType := fields[2]

		// Skip system/virtual filesystems
		skipTypes := []string{"tmpfs", "devtmpfs", "sysfs", "proc", "devpts", "cgroup", "cgroup2", "pstore", "bpf", "configfs", "selinuxfs", "debugfs", "tracefs", "fusectl", "fuse.gvfsd-fuse", "fuse.portal", "securityfs", "hugetlbfs", "mqueue", "autofs"}
		if contains(skipTypes, fsType) {
			continue
		}

		// Check if mount point is in external storage locations
		isExternal := false
		for _, prefix := range externalPrefixes {
			if strings.HasPrefix(mountPoint, prefix) {
				isExternal = true
				break
			}
		}

		// Also include USB and removable devices
		if strings.Contains(device, "/dev/sd") || strings.Contains(device, "/dev/mmcblk") {
			if mountPoint != "/" && mountPoint != "/boot" && mountPoint != "/home" {
				isExternal = true
			}
		}

		if isExternal {
			var stat syscall.Statfs_t
			err := syscall.Statfs(mountPoint, &stat)
			if err != nil {
				continue
			}

			totalBytes := stat.Blocks * uint64(stat.Bsize)
			availableBytes := stat.Bavail * uint64(stat.Bsize)
			usedBytes := totalBytes - (stat.Bfree * uint64(stat.Bsize))

			deviceInfo := map[string]interface{}{
				"device":           device,
				"mount_point":      mountPoint,
				"filesystem_type":  fsType,
				"total":            totalBytes,
				"used":             usedBytes,
				"available":        availableBytes,
				"usage_percentage": float64(usedBytes) / float64(totalBytes) * 100,
			}

			devices = append(devices, deviceInfo)
		}
	}

	return devices, scanner.Err()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
