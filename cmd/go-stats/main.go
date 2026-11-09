// Package main provides the entry point for go-stats, a native Linux system
// monitor and daemon developed by davidalvarezp (dap).
//
// Website: https://davidalvarezp.com/go-stats
// GitHub:  https://github.com/davidalvarezp/go-stats
// License: MIT
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/davidalvarezp/go-stats/internal/server"
	"github.com/davidalvarezp/go-stats/internal/service"
)

// Version represents the current release version of go-stats.
const (
	Version   = "1.0.0"
	Author    = "davidalvarezp (dap)"
	SiteURL   = "https://davidalvarezp.com/go-stats"
	RepoURL   = "https://github.com/davidalvarezp/go-stats"
)

func main() {
	// CLI Flag Definitions
	// These can be modified or customized to change default runtime behavior.
	port := flag.Int("port", 10590, "HTTP server listening port for the web dashboard and REST API")
	daemon := flag.Bool("daemon", false, "Run in background daemon mode (managed by systemd)")
	install := flag.Bool("install", false, "Install binary to /usr/local/bin and configure systemd service")
	uninstall := flag.Bool("uninstall", false, "Stop service, disable, and remove systemd unit file")
	status := flag.Bool("status", false, "Query the current systemctl status of go-stats.service")
	showVersion := flag.Bool("version", false, "Print application version, author, and documentation links")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "go-stats v%s - Native Linux System Monitor & Systemd Daemon\n", Version)
		fmt.Fprintf(os.Stderr, "Author: %s | Docs: %s\n\n", Author, SiteURL)
		fmt.Fprintf(os.Stderr, "Usage: go-stats [options]\n\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Handle --version flag
	if *showVersion {
		fmt.Printf("go-stats v%s\n", Version)
		fmt.Printf("Author: %s <%s>\n", Author, "mail@davidalvarezp.com")
		fmt.Printf("Website: %s\n", SiteURL)
		fmt.Printf("GitHub:  %s\n", RepoURL)
		return
	}

	// Handle --install flag
	if *install {
		fmt.Printf("Installing go-stats as a systemd service (Port: %d)...\n", *port)
		if err := service.InstallService(*port); err != nil {
			log.Fatalf("Installation failed: %v\n", err)
		}
		fmt.Println("Service installed successfully.")
		fmt.Printf("Web Dashboard: http://localhost:%d\n", *port)
		fmt.Println("Status check:  sudo systemctl status go-stats")
		return
	}

	// Handle --uninstall flag
	if *uninstall {
		fmt.Println("Uninstalling go-stats systemd service...")
		if err := service.UninstallService(); err != nil {
			log.Fatalf("Uninstall failed: %v\n", err)
		}
		fmt.Println("go-stats service stopped and removed.")
		return
	}

	// Handle --status flag
	if *status {
		st, active, err := service.GetServiceStatus()
		if err != nil {
			fmt.Printf("go-stats.service: %s (active=%v, error=%v)\n", st, active, err)
		} else {
			fmt.Printf("go-stats.service: %s (active=%v)\n", st, active)
		}
		return
	}

	// Normal server execution
	log.Printf("========================================================\n")
	log.Printf("go-stats v%s · Native Linux System Monitor\n", Version)
	log.Printf("Author: %s | Web: %s\n", Author, SiteURL)
	log.Printf("Web Dashboard & API: http://localhost:%d\n", *port)
	log.Printf("Daemon Mode: %v\n", *daemon)
	log.Printf("========================================================\n")

	// Initialize the HTTP server with embedded UI
	srv := server.NewServer(*port)

	// Graceful shutdown handling on SIGINT or SIGTERM
	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-shutdownSig
		log.Printf("\nReceived shutdown signal (%v). Stopping server gracefully...\n", sig)
		os.Exit(0)
	}()

	// Start server (blocking call)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server terminated with error: %v\n", err)
	}
}
