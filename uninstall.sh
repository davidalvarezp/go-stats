#!/usr/bin/env bash
# ==============================================================================
# go-stats · Uninstallation Script for Linux
# ==============================================================================
# Author: davidalvarezp (dap) <mail@davidalvarezp.com>
# Website: https://davidalvarezp.com/go-stats
# Repository: https://github.com/davidalvarezp/go-stats
# License: MIT
# ==============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RESET='\033[0m'

if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}[ERROR] Please execute as root: sudo ./uninstall.sh${RESET}"
  exit 1
fi

echo -e "${YELLOW}[INFO] Stopping go-stats systemd service...${RESET}"
systemctl stop go-stats 2>/dev/null || true

echo -e "[INFO] Disabling auto-start on boot..."
systemctl disable go-stats 2>/dev/null || true

echo -e "[INFO] Removing /etc/systemd/system/go-stats.service..."
rm -f /etc/systemd/system/go-stats.service

echo -e "[INFO] Reloading systemd daemon..."
systemctl daemon-reload
systemctl reset-failed 2>/dev/null || true

echo -e "[INFO] Removing binary /usr/local/bin/go-stats..."
rm -f /usr/local/bin/go-stats

echo -e "${GREEN}[SUCCESS] go-stats has been completely uninstalled from this system.${RESET}"
