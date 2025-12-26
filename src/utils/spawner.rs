pub fn spawn_process(command: &str, args: &[&str]) -> Result<String, std::io::Error> {
    let output = std::process::Command::new(command).args(args).output()?;
    let stdout = String::from_utf8_lossy(&output.stdout).into_owned();

    if output.status.success() {
        Ok(stdout)
    } else {
        let stderr = String::from_utf8_lossy(&output.stderr);
        Err(std::io::Error::new(
            std::io::ErrorKind::Other,
            format!(
                "command '{}' failed (status: {}): {}",
                command, output.status, stderr
            ),
        ))
    }
}
