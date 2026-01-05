#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

echo "Building aaxion application..."
go build -o aaxion ./cmd/main.go

echo "Copying aaxion binary to /usr/local/bin/"
# Using sudo because /usr/local/bin usually requires root privileges
sudo mv aaxion /usr/local/bin/aaxion

echo "aaxion installed successfully to /usr/local/bin/aaxion"

# Create systemd service file
echo "Creating systemd service file..."
USERNAME=$(whoami)
GROUPNAME=$(id -gn)
WORKING_DIRECTORY=$(pwd)

SERVICE_FILE="aaxion.service"

cat > $SERVICE_FILE << EOL
[Unit]
Description=aaxion service
After=network.target

[Service]
User=$USERNAME
Group=$GROUPNAME
WorkingDirectory=$WORKING_DIRECTORY
ExecStart=/usr/local/bin/aaxion
Restart=always

[Install]
WantedBy=multi-user.target
EOL

echo "Successfully created $SERVICE_FILE"
echo ""
echo "To install and start the service, run the following commands:"
echo "sudo mv $SERVICE_FILE /etc/systemd/system/"
echo "sudo systemctl enable aaxion"
echo "sudo systemctl start aaxion"
echo "sudo systemctl status aaxion"
