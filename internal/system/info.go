package system

import (
	"aaxion/internal/version"
	"encoding/json"
	"net/http"
	"runtime"
)

type SystemInfo struct {
	Version  string `json:"version"`
	Codename string `json:"codename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
}

func GetSystemInfoApi(w http.ResponseWriter, r *http.Request) {
	info := SystemInfo{
		Version:  version.Version,
		Codename: version.Codename,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
