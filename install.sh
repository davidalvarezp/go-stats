#!/usr/bin/env bash
# ==============================================================================
# go-stats · Automated Installation Script for Linux (Debian / Ubuntu)
# ==============================================================================
# Author: davidalvarezp (dap) <mail@davidalvarezp.com>
# Website: https://davidalvarezp.com/go-stats
# Repository: https://github.com/davidalvarezp/go-stats
# License: MIT
#
# Operations performed:
# 1. Verifies superuser privileges (sudo/root).
# 2. Checks for Go runtime; automatically installs via APT if missing.
# 3. Compiles the monolithic Go binary with embedded web dashboard.
# 4. Installs the binary to /usr/local/bin/go-stats.
# 5. Configures /etc/systemd/system/go-stats.service.
# 6. Enables and starts the background daemon on port 10590.
# 7. Validates HTTP endpoint availability.
# ==============================================================================

set -e

# Terminal formatting
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
RESET='\033[0m'

PORT=${GO_STATS_PORT:-10590}
BIN_DEST="/usr/local/bin/go-stats"
SERVICE_DEST="/etc/systemd/system/go-stats.service"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo -e "${CYAN}${BOLD}"
echo "=========================================================="
echo " go-stats · Automated Native Linux Systemd Installer"
echo " Author: davidalvarezp (dap) | https://davidalvarezp.com/go-stats"
echo "=========================================================="
echo -e "${RESET}"

# 1. Verify root privileges
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}[ERROR] Superuser privileges required.${RESET}"
  echo -e "Please execute this installer with: ${YELLOW}sudo ./install.sh${RESET}"
  exit 1
fi

# 2. Check for Go compiler; install via apt if missing
echo -e "${CYAN}[1/5] Verifying toolchain dependencies...${RESET}"
if ! command -v go >/dev/null 2>&1; then
  echo -e "${YELLOW}[INFO] Go is not installed. Installing golang-go via APT...${RESET}"
  export DEBIAN_FRONTEND=noninteractive
  apt-get update -qq
  apt-get install -y --no-install-recommends golang-go curl
  echo -e "${GREEN}[OK] Go installed successfully: $(go version)${RESET}"
else
  echo -e "${GREEN}[OK] Go runtime detected:${RESET} $(go version)"
fi

# 3. Compile optimized monolithic binary
echo -e "${CYAN}[2/5] Compiling standalone go-stats binary...${RESET}"
cd "$SCRIPT_DIR"

mkdir -p bin
CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=1.0.0 -X 'main.Author=davidalvarezp (dap)'" -o bin/go-stats cmd/go-stats/main.go

if [ ! -f "bin/go-stats" ]; then
  echo -e "${RED}[ERROR] Binary compilation failed.${RESET}"
  exit 1
fi

BIN_SIZE=$(ls -lh bin/go-stats | awk '{print $5}')
echo -e "${GREEN}[OK] Binary compiled successfully (${BIN_SIZE}).${RESET}"

# 4. Install binary to /usr/local/bin
echo -e "${CYAN}[3/5] Installing binary to ${BIN_DEST}...${RESET}"
cp bin/go-stats "$BIN_DEST"
chmod 755 "$BIN_DEST"

# 5. Create systemd unit file
echo -e "${CYAN}[4/5] Configuring systemd unit at ${SERVICE_DEST}...${RESET}"
cat <<EOF > "$SERVICE_DEST"
[Unit]
Description=go-stats - Lightweight Native Linux System Monitor
Documentation=https://davidalvarezp.com/go-stats https://github.com/davidalvarezp/go-stats
After=network.target network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/usr/local/bin
ExecStart=${BIN_DEST} -port ${PORT} -daemon
Restart=always
RestartSec=3
KillMode=process

StandardOutput=journal
StandardError=journal
SyslogIdentifier=go-stats
LimitNOFILE=65535
TimeoutStopSec=10

ProtectSystem=full
ProtectHome=read-only
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
EOF

chmod 644 "$SERVICE_DEST"

# 6. Enable and start daemon
echo -e "${CYAN}[5/5] Enabling and starting go-stats service...${RESET}"
if command -v systemctl >/dev/null 2>&1; then
  systemctl daemon-reload
  systemctl enable go-stats
  systemctl restart go-stats
else
  echo -e "${YELLOW}[WARN] systemctl not found (container environment). Starting background process with nohup...${RESET}"
  nohup "${BIN_DEST}" -port "${PORT}" -daemon >/var/log/go-stats.log 2>&1 &
fi

# Allow socket listener to bind
sleep 1.5

# 7. Health verification
echo -e "${CYAN}[INFO] Checking HTTP health on http://localhost:${PORT}...${RESET}"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:${PORT}/" || echo "000")

if [ "$HTTP_CODE" = "200" ]; then
  echo -e "${GREEN}${BOLD}==========================================================${RESET}"
  echo -e "${GREEN}${BOLD}[SUCCESS] Installation completed successfully!${RESET}"
  echo -e "${GREEN}${BOLD}==========================================================${RESET}"
  echo -e "Web Dashboard: ${CYAN}${BOLD}http://localhost:${PORT}${RESET}"
  echo -e ""
  echo -e "Systemd Service Management:"
  echo -e "  - View service status:   ${YELLOW}sudo systemctl status go-stats${RESET}"
  echo -e "  - Stream live logs:      ${YELLOW}sudo journalctl -u go-stats -f${RESET}"
  echo -e "  - Restart daemon:        ${YELLOW}sudo systemctl restart go-stats${RESET}"
  echo -e "  - Stop daemon:           ${YELLOW}sudo systemctl stop go-stats${RESET}"
  echo -e "  - Complete uninstall:    ${YELLOW}sudo ./uninstall.sh${RESET}"
  echo -e ""
else
  echo -e "${YELLOW}[WARN] Service started but HTTP check returned status ${HTTP_CODE}.${RESET}"
  echo -e "Inspect diagnostics: ${YELLOW}sudo journalctl -u go-stats -xe${RESET}"
fi
