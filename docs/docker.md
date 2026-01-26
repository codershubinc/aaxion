# Docker Deployment Guide

This guide explains how to run Aaxion using Docker for easy deployment and isolation.

## Prerequisites

- Docker installed (version 20.10 or higher)
- Docker Compose (optional, but recommended)

## Quick Start

### Using Docker

Build and run the container:

```bash
# Build the image
docker build -t aaxion:latest .

# Run the container with your file storage directory
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/home/aaxion \
  -v $(pwd)/data/db:/data \
  --name aaxion-server \
  aaxion:latest
```

### Using Docker Compose (Recommended)

The easiest way to run Aaxion is with Docker Compose:

```bash
# Create data directory with proper permissions
mkdir -p data data/db
chmod 755 data data/db

# Start the service
docker-compose up -d

# View logs
docker-compose logs -f aaxion

# Stop the service
docker-compose down
```

## Configuration

### Volumes and Permissions

**Important:** Aaxion operates on files in its home directory. The application runs as user `aaxion` (UID 1000, GID 1000).

The docker-compose.yml configures two volumes:

- `./data:/home/aaxion` - Main file storage (where Aaxion will read/write files)
- `./data/db:/data` - SQLite database storage

**To use your own directory** (e.g., `/home/yourusername/storage`):

1. **Option 1: Mount with proper permissions**
   ```yaml
   volumes:
     - /home/yourusername/storage:/home/aaxion
   ```
   Then ensure the directory is readable/writable:
   ```bash
   chmod -R 755 /home/yourusername/storage
   # Or if you need the container user to own it:
   chown -R 1000:1000 /home/yourusername/storage
   ```

2. **Option 2: Run container as your user** (avoids permission issues)
   Add to docker-compose.yml:
   ```yaml
   user: "$(id -u):$(id -g)"
   ```
   Or in docker run:
   ```bash
   docker run -d \
     -p 8080:8080 \
     -v /home/yourusername/storage:/home/aaxion \
     --user $(id -u):$(id -g) \
     aaxion:latest
   ```

3. **Option 3: Run as root** (least secure, but works for all directories)
   Add to docker-compose.yml:
   ```yaml
   user: "0:0"
   ```

You can modify these paths in `docker-compose.yml` to use different locations on your host system.

### Port Mapping

By default, Aaxion runs on port 8080. To use a different port, modify the port mapping in `docker-compose.yml`:

```yaml
ports:
  - "3000:8080"  # Maps host port 3000 to container port 8080
```

### Environment Variables

Available environment variables:

- `AAXION_DB_PATH` - Path to the SQLite database file (default: `/data/.aaxion.db`)

## Management Commands

### View Logs

```bash
# Docker Compose
docker-compose logs -f

# Docker
docker logs -f aaxion-server
```

### Restart the Service

```bash
# Docker Compose
docker-compose restart

# Docker
docker restart aaxion-server
```

### Stop the Service

```bash
# Docker Compose
docker-compose stop

# Docker
docker stop aaxion-server
```

### Remove Container and Data

```bash
# Docker Compose (keeps volumes)
docker-compose down

# Docker Compose (removes volumes)
docker-compose down -v

# Docker
docker rm -f aaxion-server
```

## Troubleshooting

### Container won't start

Check the logs:
```bash
docker-compose logs aaxion
```

### Permission Issues

The container runs as a non-root user (UID 1000). Ensure your host directories have appropriate permissions:

```bash
mkdir -p data/uploads data/db
chmod -R 755 data
```

### Network Issues

If you can't access the service, ensure the port isn't already in use:

```bash
# Check if port 8080 is in use
lsof -i :8080

# Or use netstat
netstat -tuln | grep 8080
```

## Building from Source

The Dockerfile uses a multi-stage build:

1. **Builder stage**: Compiles the Go application with CGO support for SQLite
2. **Runtime stage**: Creates a minimal Alpine Linux image with only runtime dependencies

To rebuild the image:

```bash
docker build --no-cache -t aaxion:latest .
```

## Security Features

- Runs as non-root user (UID 1000, GID 1000)
- Minimal Alpine Linux base image
- Only necessary dependencies installed
- Built-in health checks

## Advanced Configuration

### Custom Dockerfile

You can customize the Dockerfile for your needs:

- Change the Go version in the builder stage
- Modify runtime dependencies
- Adjust user permissions
- Add custom configuration

### Production Deployment

For production deployments, consider:

1. Using a reverse proxy (nginx, Caddy, Traefik)
2. Adding SSL/TLS certificates
3. Setting up automated backups for the data directory
4. Monitoring container health and resource usage
5. Using Docker secrets for sensitive configuration

Example with Caddy reverse proxy:

```yaml
version: '3.8'

services:
  aaxion:
    build: .
    container_name: aaxion-server
    volumes:
      - ./data/uploads:/data/uploads
      - ./data/db:/data
    restart: unless-stopped

  caddy:
    image: caddy:latest
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile
      - caddy_data:/data
      - caddy_config:/config
    depends_on:
      - aaxion

volumes:
  caddy_data:
  caddy_config:
```

## Support

For issues or questions:
- Check the [main README](../README.md)
- Review the [API documentation](./api.md)
- Open an issue on GitHub
