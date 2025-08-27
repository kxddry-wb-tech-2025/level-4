package client_test

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// buildBinaries builds the server and client binaries and returns their absolute paths.
func buildBinaries(t *testing.T) (serverBin, clientBin string) {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Dir(wd) // repo root is one level up from client/

	serverBin = filepath.Join(t.TempDir(), "server-bin")
	clientBin = filepath.Join(t.TempDir(), "client-bin")

	// Build server
	{
		cmd := exec.Command("go", "-C", filepath.Join(root, "server"), "build", "-o", serverBin, "./cmd/app")
		cmd.Env = os.Environ()
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("build server failed: %v\n%s", err, string(out))
		}
	}

	// Build client
	{
		cmd := exec.Command("go", "-C", filepath.Join(root, "client"), "build", "-o", clientBin, "./cmd/app")
		cmd.Env = os.Environ()
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("build client failed: %v\n%s", err, string(out))
		}
	}

	return serverBin, clientBin
}

// getFreePort asks the OS for a free TCP port and returns it.
func getFreePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// startServers launches n server processes on free ports and waits until their /health responds.
func startServers(t *testing.T, serverBin string, n int) ([]*exec.Cmd, []string) {
	t.Helper()
	cmds := make([]*exec.Cmd, 0, n)
	addrs := make([]string, 0, n)
	for i := 0; i < n; i++ {
		port := getFreePort(t)
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		cmd := exec.Command(serverBin, fmt.Sprintf("-port=%d", port))
		// Detach stdio but keep for debugging if needed
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			t.Fatalf("start server on %s: %v", addr, err)
		}
		cmds = append(cmds, cmd)
		addrs = append(addrs, addr)
	}

	// Wait for health endpoints
	deadline := time.Now().Add(5 * time.Second)
	for _, addr := range addrs {
		ok := false
		for time.Now().Before(deadline) {
			resp, err := http.Get("http://" + addr + "/health")
			if err == nil {
				_ = resp.Body.Close()
				if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
					ok = true
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
		if !ok {
			t.Fatalf("server %s did not become healthy in time", addr)
		}
	}

	// Ensure cleanup
	t.Cleanup(func() {
		for _, cmd := range cmds {
			if cmd.Process != nil {
				// Graceful kill
				_ = cmd.Process.Kill()
				_, _ = cmd.Process.Wait()
			}
		}
	})

	return cmds, addrs
}

// runClient executes the distributed grep client with given args and returns stdout.
func runClient(t *testing.T, clientBin string, args ...string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, clientBin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil && ctx.Err() == nil {
		// The client returns non-zero on errors only. Surface stderr in failure.
		t.Fatalf("client failed: %v\n%s", err, string(out))
	}
	return string(out)
}

// runSystemGrep runs the system grep with given args and returns stdout.
func runSystemGrep(t *testing.T, args ...string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "grep", args...)
	out, err := cmd.CombinedOutput()
	// For invert/no match, grep returns exit code 1 but still prints output.
	// We only care about stdout/combined output text equivalence, so ignore exit code.
	_ = err
	return string(out)
}

// writeTempFile writes lines to a temp file and returns its path.
func writeTempFile(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "grep-input-*.txt")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	w := bufio.NewWriter(f)
	for i, line := range lines {
		if _, err := w.WriteString(line); err != nil {
			t.Fatalf("write temp file: %v", err)
		}
		if i < len(lines)-1 {
			if _, err := w.WriteString("\n"); err != nil {
				t.Fatalf("write newline: %v", err)
			}
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("flush temp file: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close temp file: %v", err)
	}
	return f.Name()
}

// compare trims trailing whitespace for robustness and compares equality.
func compareOutputs(t *testing.T, got, want string) {
	t.Helper()
	trim := func(s string) string { return strings.TrimRight(s, "\n\r ") }
	if trim(got) != trim(want) {
		t.Fatalf("output mismatch\n--- distributed ---\n%s\n--- system grep ---\n%s", got, want)
	}
}

func TestSimpleMatch(t *testing.T) {
	serverBin, clientBin := buildBinaries(t)
	_, addrs := startServers(t, serverBin, 3)

	lines := []string{"alpha", "beta", "gamma", "alphabet"}
	file := writeTempFile(t, lines)

	// Distributed client
	clientArgs := []string{"--addrs", strings.Join(addrs, ","), "alpha", file}
	distOut := runClient(t, clientBin, clientArgs...)

	// System grep
	sysOut := runSystemGrep(t, "alpha", file)

	compareOutputs(t, distOut, sysOut)
}

func TestIgnoreCase(t *testing.T) {
	serverBin, clientBin := buildBinaries(t)
	_, addrs := startServers(t, serverBin, 3)

	lines := []string{"Alpha", "beta", "aLpHa", "Gamma"}
	file := writeTempFile(t, lines)

	distOut := runClient(t, clientBin, "--addrs", strings.Join(addrs, ","), "-i", "alpha", file)
	sysOut := runSystemGrep(t, "-i", "alpha", file)
	compareOutputs(t, distOut, sysOut)
}

func TestInvertMatch(t *testing.T) {
	serverBin, clientBin := buildBinaries(t)
	_, addrs := startServers(t, serverBin, 3)

	lines := []string{"foo", "bar", "baz", "food"}
	file := writeTempFile(t, lines)

	distOut := runClient(t, clientBin, "--addrs", strings.Join(addrs, ","), "-v", "foo", file)
	sysOut := runSystemGrep(t, "-v", "foo", file)
	compareOutputs(t, distOut, sysOut)
}

func TestFixedString(t *testing.T) {
	serverBin, clientBin := buildBinaries(t)
	_, addrs := startServers(t, serverBin, 3)

	lines := []string{"a.c", "abc", "a-c", "a c"}
	file := writeTempFile(t, lines)

	distOut := runClient(t, clientBin, "--addrs", strings.Join(addrs, ","), "-F", "a.c", file)
	sysOut := runSystemGrep(t, "-F", "a.c", file)
	compareOutputs(t, distOut, sysOut)
}

func TestCountOnly(t *testing.T) {
	serverBin, clientBin := buildBinaries(t)
	_, addrs := startServers(t, serverBin, 3)

	lines := []string{"one", "two", "two", "three"}
	file := writeTempFile(t, lines)

	distOut := runClient(t, clientBin, "--addrs", strings.Join(addrs, ","), "-c", "two", file)
	sysOut := runSystemGrep(t, "-c", "two", file)
	compareOutputs(t, distOut, sysOut)
}

func TestRegex(t *testing.T) {
	serverBin, clientBin := buildBinaries(t)
	_, addrs := startServers(t, serverBin, 3)

	lines := []string{"foo", "foa", "fob", "bar"}
	file := writeTempFile(t, lines)

	// pattern fo. should match foo, foa, fob
	distOut := runClient(t, clientBin, "--addrs", strings.Join(addrs, ","), "fo.", file)
	sysOut := runSystemGrep(t, "fo.", file)
	compareOutputs(t, distOut, sysOut)
}
