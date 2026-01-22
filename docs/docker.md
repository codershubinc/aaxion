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

# Run the container
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data/uploads:/data/uploads \
  -v $(pwd)/data/db:/data \
  --name aaxion-server \
  aaxion:latest
```

### Using Docker Compose (Recommended)

The easiest way to run Aaxion is with Docker Compose:

```bash
# Start the service
docker-compose up -d

# View logs
docker-compose logs -f aaxion

# Stop the service
docker-compose down
```

## Configuration

### Volumes

The docker-compose.yml configures two volumes:

- `./data/uploads:/data/uploads` - Stores uploaded files
- `./data/db:/data` - Stores the SQLite database

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
- Check the [main README](./README.md)
- Review the [API documentation](./docs/api.md)
- Open an issue on GitHub
