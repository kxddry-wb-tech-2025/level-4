package service

import (
	"bufio"
	"bytes"
	"client/internal/helpers/parser"
	"client/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"
)

func Run(pattern string, files []string, addrs []string, flags models.GrepFlags) error {
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
			fmt.Printf("server %s is not alive\n", addrs[i])
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
		fmt.Println(file)
		lines, err := openFile(file)
		if err != nil {
			return err
		}

		tasks := createTasksWithContext(lines, pattern, flags, len(aliveServers), i)

		allBlocks := make([]models.FoundBlock, 0, len(tasks))
		var mu sync.Mutex

		for j, task := range tasks {
			addr := aliveServers[j%len(aliveServers)]
			wg.Add(1)
			go func(addr *models.ParsedAddr, task models.Task) {
				defer wg.Done()

				data, err := json.Marshal(task)
				if err != nil {
					fmt.Printf("failed to marshal request: %v\n", err)
					return
				}

				resp, err := http.Post("http://"+addr.Host+":"+addr.Port+"/grep", "application/json", bytes.NewBuffer(data))
				if err != nil {
					fmt.Printf("failed to send request to %s: %v\n", addr.Host+":"+addr.Port, err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					fmt.Printf("server %s returned status %d\n", addr.Host+":"+addr.Port, resp.StatusCode)
					return
				}

				var result models.Result
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					fmt.Printf("failed to decode response from %s: %v\n", addr.Host+":"+addr.Port, err)
					return
				}

				mu.Lock()
				allBlocks = append(allBlocks, result.FoundBlocks...)
				mu.Unlock()
			}(addr, task)
		}
		wg.Wait()

		printBlocksNonOverlapping(file, allBlocks, flags)
	}
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
			StartLineNumber: left + 1,
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

func printBlocksNonOverlapping(filename string, blocks []models.FoundBlock, flags models.GrepFlags) {
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

	fmt.Println(filename)
	for _, ln := range keys {
		if flags.PrintNumbers {
			fmt.Printf("%d:%s\n", ln, lineToText[ln])
		} else {
			fmt.Println(lineToText[ln])
		}
	}
}
