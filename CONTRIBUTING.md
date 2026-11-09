# Contributing to go-stats

Thank you for your interest in contributing to **go-stats**!

## Code of Conduct
We are committed to providing a friendly, safe, and welcoming environment for everyone, regardless of experience level, gender identity, or background.

## Development Workflow

### Prerequisites
- Linux OS (Debian 11/12, Ubuntu 22.04/24.04, or similar system with `/proc` virtual filesystem)
- Go 1.20 or later

### Setting Up
1. Fork the repository on GitHub: [github.com/davidalvarezp/go-stats](https://github.com/davidalvarezp/go-stats)
2. Clone your fork locally:
   ```bash
   git clone https://github.com/<your-username>/go-stats.git
   cd go-stats
   ```
3. Run tests to verify the baseline:
   ```bash
   make test
   ```
4. Build the binary locally:
   ```bash
   make build
   ```
5. Test in foreground:
   ```bash
   ./bin/go-stats -port 10590
   ```

## Coding Guidelines
- **Zero External Dependencies**: The core collector must exclusively parse native Linux `/proc` files. Do not import heavy external telemetry libraries.
- **Pure Go Self-Contained Web UI**: The web dashboard must remain self-contained in `internal/server/ui.go` or standard embed so `go build` runs anywhere without requiring Node.js or npm.
- **Go Standards**: Format code with `gofmt -s -w .` and verify with `go vet ./...`.
- **Commit Messages**: Follow standard conventional commits (e.g., `feat: add disk I/O latency metrics`, `fix: handle empty /proc/loadavg line`).

## Submitting Pull Requests
1. Create a feature branch: `git checkout -b feat/my-enhancement`
2. Make your changes with unit tests if applicable.
3. Verify tests pass: `CGO_ENABLED=0 go test -v ./...`
4. Commit and push to your fork.
5. Open a Pull Request against the `main` branch on GitHub.

For questions or security discussions, reach out to **davidalvarezp** at [mail@davidalvarezp.com](mailto:mail@davidalvarezp.com).
