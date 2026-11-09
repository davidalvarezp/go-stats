# go-stats

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Debian%20%7C%20Ubuntu-E95420?style=flat&logo=linux)](https://davidalvarezp.com/go-stats)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub release](https://img.shields.io/badge/Release-v1.0.0-blue.svg)](https://github.com/davidalvarezp/go-stats/releases)
[![Author](https://img.shields.io/badge/Author-dap-black.svg)](https://davidalvarezp.com)

A high-performance, lightweight Linux system monitor and background daemon written in **Go**.

`go-stats` runs silently as a native **systemd** service (`go-stats.service`) and serves a real-time web dashboard on port **`10590`** (`http://localhost:10590`). It provides the iconic observability metrics of `top` and `htop` through a clean, modern interface styled in a Linux terminal aesthetic (deep charcoal black with terminal green accents).

The compiled binary is **100% self-contained**: the web dashboard is embedded directly into the Go binary with zero external runtime dependencies. No Node.js, npm, or external assets are required on target hosts.

---

## Key Features

- **Native Linux Kernel `/proc` Parsing**: Zero third-party telemetry agents. Direct parsing of `/proc/stat`, `/proc/meminfo`, `/proc/loadavg`, `/proc/net/dev`, `/proc/mounts`, and `/proc/[pid]/*`.
- **Systemd Daemon Control**: Complete lifecycle control via `systemctl` (`start`, `stop`, `restart`, `enable`, `status`) and real-time log streaming via `journalctl`.
- **Embedded Web Dashboard on Port 10590**: Serves an interactive single-page dashboard at `http://localhost:10590` with real-time sub-second updates.
- **Terminal Green & Charcoal Branding**: Designed around a Linux terminal aesthetic. All colors are exposed as CSS custom variables at the top of the stylesheet for easy customization.
- **Top & Htop Process Management**: Full process table with live search, column sorting, state filtering (Running, Sleeping, Zombie), and interactive process termination (`SIGTERM`, `SIGKILL`).
- **Multi-Core CPU Matrix**: Individual core percentage meters, load breakdown (user, system, idle, iowait), and a 30-second rolling sparkline history.
- **Minimal Footprint**: Consumes under 15 MB of resident memory (RES) and less than 0.2% CPU under continuous monitoring.

---

## Quick Installation (Debian / Ubuntu)

Clone the repository and run the automated installer:

```bash
git clone https://github.com/davidalvarezp/go-stats.git
cd go-stats
sudo ./install.sh
```

### What `install.sh` Does Automatically:
1. Verifies root privileges (`sudo`).
2. Checks for the Go compiler. If not installed, it automatically installs `golang-go` via APT.
3. Compiles the monolithic binary with `-ldflags="-s -w"` optimization.
4. Installs the executable to `/usr/local/bin/go-stats`.
5. Configures the systemd unit file at `/etc/systemd/system/go-stats.service`.
6. Enables automatic startup on boot (`systemctl enable go-stats`) and starts the daemon (`systemctl restart go-stats`).
7. Performs an automated health check on `http://localhost:10590`.

Once installed, open your browser:
👉 **`http://localhost:10590`**

---

## Systemd Service Management

`go-stats` is managed through standard systemd commands:

```bash
# Check daemon status and resource utilization
sudo systemctl status go-stats

# Stream daemon logs in real time
sudo journalctl -u go-stats -f

# Restart the service
sudo systemctl restart go-stats

# Stop the daemon
sudo systemctl stop go-stats

# Start the daemon
sudo systemctl start go-stats

# Disable autostart on boot
sudo systemctl disable go-stats

# Completely uninstall from system
sudo ./uninstall.sh
```

---

## Manual Build (GNU Makefile)

You can build and run `go-stats` manually using standard Go toolchains:

```bash
# Build optimized static binary
make build-prod

# Run in foreground on default port 10590
make run

# Run unit tests
make test

# Run code linter
make lint

# Install to system
sudo make install
```

---

## Command-Line Flags

```text
Usage: go-stats [options]

Options:
  -port int       HTTP server listening port (default: 10590)
  -daemon         Run in background daemon mode managed by systemd (default: false)
  -install        Install binary to /usr/local/bin and configure systemd service
  -uninstall      Stop service, disable autostart, and remove files
  -status         Query current systemctl status of go-stats.service
  -version        Print application version and author metadata
```

Example with custom listening port:
```bash
go-stats -port 8080 -daemon
```

---

## REST API Specification

The embedded HTTP server exposes structured JSON endpoints for automation and external monitoring tools:

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/metrics` | System snapshot: multi-core CPU, RAM, Swap, Load, Mounts, Network |
| `GET` | `/api/processes` | List of active Linux processes scanned from `/proc/[pid]/*` |
| `POST` | `/api/process/signal` | Send POSIX signals (`{"pid": 1234, "signal": "SIGTERM"}`) |
| `GET` | `/api/systemd/status` | Current status of `go-stats.service` via `systemctl` |
| `POST` | `/api/systemd/action` | Execute action (`{"action": "restart"}`) |
| `GET` | `/api/health` | Service healthcheck endpoint |

### Example Metrics Response (`GET /api/metrics`):
```json
{
  "timestamp": 1794189000000,
  "hostname": "debian-host",
  "os_name": "Debian GNU/Linux 12 (bookworm)",
  "kernel": "Linux 6.1.0-28-amd64",
  "uptime_string": "12d 04h 18m",
  "cpu_count": 4,
  "cpu": {
    "overall": 8.4,
    "user": 5.1,
    "system": 2.3,
    "idle": 91.6,
    "per_core": [9.1, 7.8, 8.4, 8.3]
  },
  "memory": {
    "total": 16777216000,
    "used": 4194304000,
    "available": 12582912000,
    "used_percent": 25.0
  }
}
```

---

## Branding & CSS Customization

The web dashboard is styled with a terminal dark-mode theme by default, featuring custom CSS variables defined at the top of `internal/server/ui.go`:

```css
:root {
  /* Canvas and surfaces */
  --bg-canvas: #090a0c;             /* Main background */
  --bg-surface: #121418;            /* Card backgrounds */
  --bg-surface-elevated: #1a1d24;   /* Inputs, hover states */
  --border-subtle: #242933;         /* Hairline dividers */

  /* Terminal green branding */
  --color-accent-green: #22c55e;    /* Primary highlight color */
  --color-accent-glow: rgba(34, 197, 94, 0.15);
  
  /* Telemetry meters */
  --color-cyan: #06b6d4;            /* Disk and virtual memory */
  --color-indigo: #818cf8;          /* Swap metrics */
  --color-warning: #eab308;         /* Load warnings */
  --color-danger: #ef4444;          /* Critical thresholds and errors */
}
```

To modify color schemes, edit the CSS variables block in `internal/server/ui.go` and run `make build-prod`.

---

## Project Structure

```text
go-stats/
├── .github/
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.md          # Standard bug report template
│   │   └── feature_request.md     # Feature request template
│   └── workflows/
│       ├── ci.yml                 # Multi-arch CI pipeline (amd64, arm64)
│       └── release.yml            # Automated GitHub Release builder
├── cmd/
│   └── go-stats/
│       └── main.go                # CLI entry point, flag parsing, signal trap
├── internal/
│   ├── collector/
│   │   ├── proc_linux.go          # Native /proc virtual filesystem parser
│   │   └── proc_test.go           # Unit tests for parser functions
│   ├── server/
│   │   ├── server.go              # HTTP server, REST routes, CORS middleware
│   │   └── ui.go                  # Embedded single-page dashboard HTML/CSS/JS
│   └── service/
│       └── systemd.go             # Systemd integration, unit manager
├── go-stats.service               # Systemd service unit template
├── install.sh                     # Automated Debian/Ubuntu installer script
├── uninstall.sh                   # Clean uninstaller script
├── Makefile                       # Standard GNU Makefile build rules
├── go.mod                         # Go module definition (github.com/davidalvarezp/go-stats)
├── CHANGELOG.md                   # Version history
├── CONTRIBUTING.md                # Contribution guidelines
├── LICENSE                        # MIT License
└── README.md                      # Project documentation
```

---

## Author & Contact

**davidalvarezp (dap)**
- **Website**: [https://davidalvarezp.com](https://davidalvarezp.com)
- **Tool Page**: [https://davidalvarezp.com/go-stats](https://davidalvarezp.com/go-stats)
- **GitHub**: [@davidalvarezp](https://github.com/davidalvarezp)
- **Email**: [mail@davidalvarezp.com](mailto:mail@davidalvarezp.com)

---

## License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.
