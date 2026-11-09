# Changelog

All notable changes to **go-stats** will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-11-09

### Added
- Native Linux kernel `/proc` parser for multi-core CPU (`/proc/stat`), RAM & Swap (`/proc/meminfo`), load average (`/proc/loadavg`), interfaces (`/proc/net/dev`), storage mounts (`/proc/mounts`), and processes (`/proc/[pid]/*`).
- Built-in embedded web dashboard on port `10590` (`http://localhost:10590`) featuring Linux terminal green and deep charcoal branding.
- Multi-core CPU matrix with individual percentage gauges and 30-second load sparklines.
- Interactive `top`/`htop` process table with live search, column sorting, state filtering (Running, Sleeping, Zombie), and POSIX signals (`SIGTERM`, `SIGKILL`).
- Systemd daemon mode and service control via `systemctl start/stop/restart/enable` and real-time `journalctl` stream.
- Zero-dependency architecture: binary compiles strictly using `go build` with no Node.js/npm requirements.
- Automated Debian/Ubuntu installer script (`install.sh`) with auto-installation of Go via APT and systemd unit registration.
- Automated uninstaller script (`uninstall.sh`).
- Full GNU Makefile with `build`, `build-prod`, `test`, `install`, `uninstall`, and `clean` targets.
- GitHub Actions CI workflow for multi-arch builds (`amd64` and `arm64`) across Go versions 1.20 to 1.24.
