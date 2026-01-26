#!/usr/bin/env bash
set -euo pipefail

# One-time server setup (run on the server as a sudo-capable user)
# Adjust PORT/USER as needed.

APP_USER=wrzapi
APP_DIR=/opt/wrzapi

sudo useradd -r -s /usr/sbin/nologin "$APP_USER" || true
sudo mkdir -p "$APP_DIR"/releases "$APP_DIR"/current
sudo chown -R "$APP_USER":"$APP_USER" "$APP_DIR"

sudo tee /etc/systemd/system/wrzapi.service >/dev/null <<'EOF'
[Unit]
Description=WRZ API Service
After=network.target

[Service]
Type=simple
User=wrzapi
WorkingDirectory=/opt/wrzapi
Environment=PORT=12088
ExecStart=/opt/wrzapi/current
Restart=on-failure
RestartSec=3

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable wrzapi