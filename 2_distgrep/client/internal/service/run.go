package service

import (
	"bufio"
	"bytes"
	"client/internal/helpers/parser"
	"client/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
)

// Run is the main function for running the grep service
func Run(pattern string, files []string, addrs []string, flags models.GrepFlags, quorum int) error {
	Err := os.Stderr
	aliveServers := make([]*models.ParsedAddr, 0, len(addrs))
	for i := range addrs {
		parsed, err := parser.ParseAddress(addrs[i], "http")
		if err != nil {
			return err
		}

		req, err := http.NewRequest("GET", "http://"+parsed.Host+":"+parsed.Port+"/health", nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent) {
			aliveServers = append(aliveServers, parsed)
		} else {
			fmt.Fprintf(Err, "server %s is not alive\n", addrs[i])
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
	if len(aliveServers) == 0 {
		return fmt.Errorf("no alive servers found")
	}

	if quorum <= 0 || quorum > len(aliveServers) {
		quorum = len(aliveServers)/2 + 1
	}
	wg := new(sync.WaitGroup)

	for i, file := range files {
		lines, err := openInput(file)
		if err != nil {
			return err
		}

		tasks := createTasksWithContext(lines, pattern, flags, len(aliveServers), i)

		allBlocks := make([]models.FoundBlock, 0, len(tasks))
		var mu sync.Mutex
		var successCount int

		for j, task := range tasks {
			addr := aliveServers[j%len(aliveServers)]
			wg.Add(1)
			go func(addr *models.ParsedAddr, task models.Task) {
				defer wg.Done()

				data, err := json.Marshal(task)
				if err != nil {
					fmt.Fprintf(Err, "failed to marshal request: %v\n", err)
					return
				}

				resp, err := http.Post("http://"+addr.Host+":"+addr.Port+"/grep", "application/json", bytes.NewBuffer(data))
				if err != nil {
					fmt.Fprintf(Err, "failed to send request to %s: %v\n", addr.Host+":"+addr.Port, err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					fmt.Fprintf(Err, "server %s returned status %d\n", addr.Host+":"+addr.Port, resp.StatusCode)
					return
				}

				var result models.Result
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					fmt.Fprintf(Err, "failed to decode response from %s: %v\n", addr.Host+":"+addr.Port, err)
					return
				}

				mu.Lock()
				allBlocks = append(allBlocks, result.FoundBlocks...)
				successCount++
				mu.Unlock()
			}(addr, task)
		}
		wg.Wait()

		if successCount >= quorum {
			printBlocksNonOverlapping(file, allBlocks, flags, len(files) > 1)
		} else {
			return fmt.Errorf("quorum not reached for %s: got %d, need %d", file, successCount, quorum)
		}
	}
	return nil
}

// createTasksWithContext creates tasks with context for each server
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
			StartLineNumber: left + 1,
			Flags:           flags,
		})
	}
	return out
}

// openFile opens a file and returns its lines
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

// openInput opens an input file or stdin
func openInput(name string) ([]string, error) {
	if name == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		lines := make([]string, 0)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		return lines, scanner.Err()
	}
	return openFile(name)
}

// printBlocksNonOverlapping prints blocks of lines non-overlapping
func printBlocksNonOverlapping(filename string, blocks []models.FoundBlock, flags models.GrepFlags, printFileName bool) {
	if len(blocks) == 0 {
		return
	}

	lineToText := make(map[int]string, 1024)
	for _, b := range blocks {
		for k, s := range b.Lines {
			ln := b.StartLineNumber + k
			if _, exists := lineToText[ln]; !exists {
				lineToText[ln] = s
			}
		}
	}

	keys := make([]int, 0, len(lineToText))
	for ln := range lineToText {
		keys = append(keys, ln)
	}
	sort.Ints(keys)

	if printFileName {
		fmt.Println(filename)
	}
	for _, ln := range keys {
		if flags.PrintNumbers {
			fmt.Printf("%d:%s\n", ln, lineToText[ln])
		} else {
			fmt.Println(lineToText[ln])
		}
	}
}
