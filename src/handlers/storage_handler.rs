use axum::{http::StatusCode, response::IntoResponse, Json};
use serde_json::Value;

use crate::utils;

pub async fn storage_info() -> impl IntoResponse {
    let handle = tokio::task::spawn_blocking(|| {
        utils::spawner::spawn_process("df", &["-h", "/home/swap"]).map_err(|e| e.to_string())
    });

    match handle.await {
        Ok(Ok(stdout)) => match df_to_json(&stdout) {
            Ok(v) => (StatusCode::OK, Json(v)).into_response(),
            Err(e) => (StatusCode::INTERNAL_SERVER_ERROR, e).into_response(),
        },
        Ok(Err(err)) => (StatusCode::INTERNAL_SERVER_ERROR, err).into_response(),
        Err(join_err) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            format!("failed to run command: {}", join_err),
        )
            .into_response(),
    }
}

fn df_to_json(output: &str) -> Result<Value, String> {
    let mut items = Vec::new();
    let mut lines = output.lines();

    // skip header (first line)
    lines.next().ok_or_else(|| "empty df output".to_string())?;

    for line in lines {
        let line = line.trim();
        if line.is_empty() {
            continue;
        }
        let parts: Vec<&str> = line.split_whitespace().collect();
        if parts.len() < 6 {
            return Err(format!("unexpected df line: {}", line));
        }

        let filesystem = parts[0].to_string();
        let size = parts[1].to_string();
        let used = parts[2].to_string();
        let avail = parts[3].to_string();
        let use_percent = parts[4].to_string();
        // mount point may contain spaces, so join the rest
        let mounted_on = parts[5..].join(" ");

        let obj = serde_json::json!({
            "filesystem": filesystem,
            "size": size,
            "used": used,
            "avail": avail,
            "use_percent": use_percent,
            "mounted_on": mounted_on
        });

        items.push(obj);
    }

    Ok(Value::Array(items))
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_json::json;

    #[test]
    fn parse_df_sample_flat() {
        let sample = "Filesystem      Size  Used Avail Use% Mounted on\n/dev/nvme0n1p4  192G  108G   75G  60% /";
        let v = df_to_json(sample).unwrap();
        let expected = json!([
            {
                "filesystem": "/dev/nvme0n1p4",
                "size": "192G",
                "used": "108G",
                "avail": "75G",
                "use_percent": "60%",
                "mounted_on": "/"
            }
        ]);
        assert_eq!(v, expected);
    }
}
