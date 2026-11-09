// Package service manages Linux systemd daemon lifecycle and unit files for go-stats.
//
// Author: davidalvarezp (dap)
// Website: https://davidalvarezp.com/go-stats
// GitHub:  https://github.com/davidalvarezp/go-stats
package service

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// SystemdUnitPath is the default location for system-wide systemd units.
const SystemdUnitPath = "/etc/systemd/system/go-stats.service"

// BinaryInstallPath is the standard destination for the executable binary.
const BinaryInstallPath = "/usr/local/bin/go-stats"

// ServiceUnitTemplate defines the systemd configuration file content.
const ServiceUnitTemplate = `[Unit]
Description=go-stats - Lightweight Native Linux System Monitor
Documentation=https://davidalvarezp.com/go-stats https://github.com/davidalvarezp/go-stats
After=network.target network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/usr/local/bin
ExecStart=/usr/local/bin/go-stats -port %d -daemon
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
`

// InstallService writes the systemd unit file and enables the service.
func InstallService(port int) error {
	unitContent := fmt.Sprintf(ServiceUnitTemplate, port)
	err := os.WriteFile(SystemdUnitPath, []byte(unitContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s: %w (did you run with sudo?)", SystemdUnitPath, err)
	}

	// Reload systemd daemon
	if err := exec.Command("systemctl", "daemon-reload").Run(); err != nil {
		return fmt.Errorf("failed systemctl daemon-reload: %w", err)
	}

	// Enable service for autostart on boot
	if err := exec.Command("systemctl", "enable", "go-stats").Run(); err != nil {
		return fmt.Errorf("failed systemctl enable: %w", err)
	}

	// Start or restart service immediately
	if err := exec.Command("systemctl", "restart", "go-stats").Run(); err != nil {
		return fmt.Errorf("failed systemctl restart: %w", err)
	}

	return nil
}

// UninstallService stops, disables, and removes the unit file and binary.
func UninstallService() error {
	_ = exec.Command("systemctl", "stop", "go-stats").Run()
	_ = exec.Command("systemctl", "disable", "go-stats").Run()
	_ = os.Remove(SystemdUnitPath)
	_ = os.Remove(BinaryInstallPath)
	_ = exec.Command("systemctl", "daemon-reload").Run()
	return nil
}

// GetServiceStatus queries systemctl to check if the unit is active.
func GetServiceStatus() (string, bool, error) {
	out, err := exec.Command("systemctl", "is-active", "go-stats").Output()
	status := strings.TrimSpace(string(out))
	isActive := (status == "active")
	return status, isActive, err
}

// ControlService runs standard systemctl verbs (start, stop, restart, enable, disable).
func ControlService(verb string) (string, error) {
	out, err := exec.Command("systemctl", verb, "go-stats").CombinedOutput()
	return string(out), err
}
