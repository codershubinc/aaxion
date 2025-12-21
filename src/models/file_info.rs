use serde::Serialize;

#[derive(Serialize)]
pub struct FileInfo {
    pub name: String,
    pub is_dir: bool,
    pub size: u64,
    pub path: String,
    pub raw_path: String,
}
