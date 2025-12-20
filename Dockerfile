FROM rust:1.83-slim-bookworm as builder

WORKDIR /app

COPY Cargo.toml Cargo.lock ./

# dummy build
RUN mkdir src && \
    echo "fn main() {}" > src/main.rs && \
    cargo build --release && \
    rm -rf src 

COPY src ./src
COPY assets ./assets

RUN touch src/main.rs && cargo build --release

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy the binary from the builder stage
COPY --from=builder /app/target/release/localdrive-rs ./server

# Copy the assets folder (HTML/CSS)
COPY --from=builder /app/assets ./assets

RUN mkdir uploads

EXPOSE 8080

CMD ["./server"]
