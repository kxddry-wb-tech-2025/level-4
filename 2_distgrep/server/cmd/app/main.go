package main

import (
	"flag"
	"fmt"
	"grep-server/internal/delivery"
	"grep-server/internal/service"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func main() {
	var port int
	var daemon bool
	flag.BoolVar(&daemon, "d", false, "run as daemon")
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()

	if daemon {
		logDir := "./logs"
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatalf("failed to create log dir: %v", err)
		}
		tmpFile, err := os.CreateTemp(logDir, "daemon-*.log")
		if err != nil {
			log.Fatalf("failed to open temp log file: %v", err)
		}
		f := tmpFile
		logFilePath := f.Name()

		cmd := exec.Command(os.Args[0], fmt.Sprintf("-port=%d", port))
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}
		cmd.Stdout = f
		cmd.Stderr = f
		err = cmd.Start()
		if err != nil {
			log.Fatalf("failed to start daemon: %v", err)
		}
		finalPath := filepath.Join(logDir, fmt.Sprintf("%d.log", cmd.Process.Pid))
		if rerr := os.Rename(logFilePath, finalPath); rerr != nil {
			log.Printf("failed to rename log file to %s: %v", finalPath, rerr)
			finalPath = logFilePath
		}
		if _, werr := f.WriteString(fmt.Sprintf("daemon pid=%d\n", cmd.Process.Pid)); werr != nil {
			log.Printf("failed to write pid to log file: %v", werr)
		}
		_ = f.Close()
		log.Printf("daemon started with PID %d (logging to %s)", cmd.Process.Pid, finalPath)
		os.Exit(0)
	}

	srv := delivery.NewServer(service.NewService())
	for err := srv.Start(port); err != nil && port < 65535; func() {
		port++
		err = srv.Start(port)
	}() {
		log.Printf("failed to start server on port %d: %v", port, err)
	}

	log.Printf("server started on port %d", port)
}
