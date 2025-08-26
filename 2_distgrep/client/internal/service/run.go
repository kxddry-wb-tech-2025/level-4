package service

import (
	"bufio"
	"bytes"
	"client/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

func Run(pattern string, files []string, addrs []string, flags models.GrepFlags) error {
	aliveServers := make([]string, 0, len(addrs))
	for i := range addrs {
		if addrs[i] == "" {
			return errors.New("empty address found")
		}
		if !strings.Contains(addrs[i], ":") {
			return errors.New("address must contain a port (e.g., host:port)")
		}
		resp, err := http.Get("http://" + addrs[i] + "/health")
		if err == nil && (resp.StatusCode == 200 || resp.StatusCode == 204) {
			aliveServers = append(aliveServers, addrs[i])
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
	if len(aliveServers) == 0 {
		return errors.New("no alive servers found")
	}
	wg := new(sync.WaitGroup)

	for i, file := range files {
		lines, err := openFile(file)
		if err != nil {
			return err
		}

		// Create tasks with context handling
		tasks := createTasksWithContext(lines, pattern, flags, len(aliveServers), i)

		for j, addr := range aliveServers {
			for _, task := range tasks {
				wg.Add(1)
				go func(addr string, task models.Task) {
					defer wg.Done()

					task.ID = i*len(aliveServers) + j

					data, err := json.Marshal(task)
					if err != nil {
						fmt.Printf("failed to marshal request: %v\n", err)
						return
					}

					resp, err := http.Post("http://"+addr+"/task", "application/json", bytes.NewBuffer(data))
					if err != nil {
						fmt.Printf("failed to send request to %s: %v\n", addr, err)
						return
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						fmt.Printf("server %s returned status %d\n", addr, resp.StatusCode)
						return
					}
				}(addr, task)
			}
		}
	}

	wg.Wait()
	return nil
}

func createTasksWithContext(lines []string, pattern string, flags models.GrepFlags, numServers, offset int) []models.Task {
	out := make([]models.Task, 0, numServers)
	ctxB := flags.Before
	ctxA := flags.After
	for i := range numServers {
		left := i * len(lines) / numServers
		right := min((i+1)*len(lines)/numServers, len(lines))
		Lines := lines[left:right]
		before := lines[max(0, left-ctxB):left]
		after := lines[right:min(right+ctxA, len(lines))]
		out = append(out, models.Task{
			Pattern:         pattern,
			Lines:           Lines,
			ID:              i + offset*numServers,
			BeforeContext:   before,
			AfterContext:    after,
			StartLineNumber: left,
			Flags:           flags,
		})
	}
	return out
}

func openFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
