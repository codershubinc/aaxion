use axum::{extract::Query, http::StatusCode, response::IntoResponse, Json};
use serde::{Deserialize, Serialize};
use std::ffi::CString;
use std::io;

/// Root upload directory (same as other handlers)
const UPLOAD_DIR: &str = "/home/swap/aaxion/";

#[derive(Serialize)]
struct StorageInfo {
    path: String,
    total_bytes: u64,
    available_bytes: u64,
    free_bytes: u64,
    used_bytes: u64,
    used_percent: f64,
    block_size: u64,
}

fn get_fs_stats(path: &str) -> Result<StorageInfo, io::Error> {
    let c_path = CString::new(path).map_err(|e| io::Error::new(io::ErrorKind::InvalidInput, e))?;

    unsafe {
        let mut stat: libc::statvfs = std::mem::zeroed();
        if libc::statvfs(c_path.as_ptr(), &mut stat) != 0 {
            return Err(io::Error::last_os_error());
        }

        let bsize = stat.f_frsize as u64;
        let total = (stat.f_blocks as u64).saturating_mul(bsize);
        let free = (stat.f_bfree as u64).saturating_mul(bsize);
        let avail = (stat.f_bavail as u64).saturating_mul(bsize);
        let used = total.saturating_sub(free);
        let used_percent = if total > 0 {
            (used as f64 / total as f64) * 100.0
        } else {
            0.0
        };

        Ok(StorageInfo {
            path: path.to_string(),
            total_bytes: total,
            available_bytes: avail,
            free_bytes: free,
            used_bytes: used,
            used_percent: (used_percent * 100.0).round() / 100.0,
            block_size: bsize,
        })
    }
}

#[derive(Deserialize)]
pub struct StorageQuery {
    /// Optional path to inspect. If omitted, uses UPLOAD_DIR.
    pub path: Option<String>,
}

pub async fn storage_info(Query(params): Query<StorageQuery>) -> axum::response::Response {
    let path = params.path.as_deref().unwrap_or(UPLOAD_DIR);

    if !std::path::Path::new(path).exists() {
        return (
            StatusCode::NOT_FOUND,
            Json(serde_json::json!({"status": "error", "message": "Path not found"})),
        )
            .into_response();
    }

    match get_fs_stats(path) {
        Ok(info) => Json(serde_json::json!({"status": "success", "data": info})).into_response(),
        Err(e) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(serde_json::json!({"status": "error", "message": e.to_string()})),
        )
            .into_response(),
    }
}

/// Returns a list of mounted filesystems with their storage stats by reading /proc/mounts
pub async fn storage_mounts() -> axum::response::Response {
    let mounts_raw = match std::fs::read_to_string("/proc/mounts") {
        Ok(s) => s,
        Err(e) => {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({"status": "error", "message": format!("Failed to read /proc/mounts: {}", e)})),
            )
            .into_response();
        }
    };

    let mut result: Vec<serde_json::Value> = Vec::new();

    for line in mounts_raw.lines() {
        let parts: Vec<&str> = line.split_whitespace().collect();
        if parts.len() < 3 {
            continue;
        }
        let device = parts[0];
        let mountpoint = parts[1];
        let fstype = parts[2];

        match get_fs_stats(mountpoint) {
            Ok(info) => result.push(serde_json::json!({"device": device, "mountpoint": mountpoint, "fstype": fstype, "stats": info})),
            Err(err) => result.push(serde_json::json!({"device": device, "mountpoint": mountpoint, "fstype": fstype, "error": err.to_string()})),
        }
    }

    Json(serde_json::json!({"status": "success", "data": result})).into_response()
}

#[derive(Serialize)]
struct DeviceStorage {
    device: String,
    mountpoint: String,
    fstype: String,
    total_bytes: u64,
    available_bytes: u64,
    free_bytes: u64,
    used_bytes: u64,
    used_percent: f64,
    block_size: u64,
}

#[derive(Serialize)]
struct SystemStorage {
    total_bytes: u64,
    available_bytes: u64,
    free_bytes: u64,
    used_bytes: u64,
    used_percent: f64,
    devices: Vec<DeviceStorage>,
}

#[derive(Deserialize)]
pub struct SystemQuery {
    /// Optional mountpoint to filter, e.g. `/` (exact match)
    pub mount: Option<String>,
}

/// Aggregates storage across physical block devices (devices starting with `/dev/`).
/// If `?mount=/` is provided, only include devices whose mountpoint equals `/`.
pub async fn storage_system(Query(params): Query<SystemQuery>) -> axum::response::Response {
    let mounts_raw = match std::fs::read_to_string("/proc/mounts") {
        Ok(s) => s,
        Err(e) => {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(serde_json::json!({"status": "error", "message": format!("Failed to read /proc/mounts: {}", e)})),
            )
            .into_response();
        }
    };

    let filter_mount = params.mount.as_deref();

    let mut seen_devices = std::collections::HashSet::new();
    let mut devices: Vec<DeviceStorage> = Vec::new();
    let mut total: u64 = 0;
    let mut free: u64 = 0;
    let mut avail: u64 = 0;

    for line in mounts_raw.lines() {
        let parts: Vec<&str> = line.split_whitespace().collect();
        if parts.len() < 3 {
            continue;
        }
        let device = parts[0];
        let mountpoint = parts[1];
        let fstype = parts[2];

        // Only consider block devices under /dev to approximate physical storage
        if !device.starts_with("/dev/") {
            continue;
        }

        // If a filter mountpoint is provided, skip anything that doesn't match exactly
        if let Some(filter) = filter_mount {
            if mountpoint != filter {
                continue;
            }
        }

        if !seen_devices.insert(device.to_string()) {
            // already processed this device
            continue;
        }

        match get_fs_stats(mountpoint) {
            Ok(info) => {
                total = total.saturating_add(info.total_bytes);
                free = free.saturating_add(info.free_bytes);
                avail = avail.saturating_add(info.available_bytes);

                devices.push(DeviceStorage {
                    device: device.to_string(),
                    mountpoint: mountpoint.to_string(),
                    fstype: fstype.to_string(),
                    total_bytes: info.total_bytes,
                    available_bytes: info.available_bytes,
                    free_bytes: info.free_bytes,
                    used_bytes: info.used_bytes,
                    used_percent: info.used_percent,
                    block_size: info.block_size,
                });
            }
            Err(_) => {
                // skip devices we can't stat
            }
        }
    }

    if filter_mount.is_some() && devices.is_empty() {
        return (
            StatusCode::NOT_FOUND,
            Json(serde_json::json!({"status": "error", "message": format!("No device found for mountpoint '{}'", filter_mount.unwrap())})),
        )
        .into_response();
    }

    let used = total.saturating_sub(free);
    let used_percent = if total > 0 {
        (used as f64 / total as f64) * 100.0
    } else {
        0.0
    };

    let system = SystemStorage {
        total_bytes: total,
        available_bytes: avail,
        free_bytes: free,
        used_bytes: used,
        used_percent: (used_percent * 100.0).round() / 100.0,
        devices,
    };

    Json(serde_json::json!({"status": "success", "data": system})).into_response()
}
