use std::sync::Arc;
use tokio::net::UdpSocket;

pub async fn start_discovery_listener(port: u16) {
    // 1. Bind UDP to the SAME port as your web server (18875)
    // TCP 18875 and UDP 18875 do not conflict.
    let socket = match UdpSocket::bind(format!("0.0.0.0:{}", port)).await {
        Ok(s) => s,
        Err(e) => {
            eprintln!("❌ Failed to bind UDP {}: {}", port, e);
            return;
        }
    };

    socket.set_broadcast(true).unwrap();
    println!("📡 Discovery Service active on UDP {} (Same as Web)", port);

    let socket = Arc::new(socket);
    let mut buf = [0u8; 1024];

    loop {
        if let Ok((size, peer)) = socket.recv_from(&mut buf).await {
            let msg = String::from_utf8_lossy(&buf[..size]);

            if msg.trim() == "DISCOVER_MAIN_SERVER" {
                // Respond: "I am here!"
                let reply = format!("Main Server Active on {}", port);
                let _ = socket.send_to(reply.as_bytes(), peer).await;
            }
        }
    }
}
