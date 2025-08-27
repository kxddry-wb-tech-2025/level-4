package service

import (
	"errors"
	"fmt"
	"grep-server/internal/models"
	"regexp"
	"strconv"
	"strings"
)

// Service is the struct for the service layer
type Service struct {
}

// NewService creates a new service
func NewService() *Service {
	return &Service{}
}

// Grep is the function for the grep endpoint
func (s *Service) Grep(req models.Request) (models.Response, error) {
	resp := models.Response{TaskID: req.ID}

	lines := req.Lines
	pattern := req.Pattern
	flags := req.Flags

	if pattern == "" {
		return resp, errors.New("empty pattern")
	}

	var re *regexp.Regexp
	var err error
	if flags.FixedString {
	} else {
		pat := pattern
		if flags.IgnoreCase {
			pat = "(?i)" + pattern
		}
		re, err = regexp.Compile(pat)
		if err != nil {
			return resp, fmt.Errorf("invalid regex: %w", err)
		}
	}

	matchesLine := func(s string) bool {
		if flags.FixedString {
			if flags.IgnoreCase {
				return strings.Contains(strings.ToLower(s), strings.ToLower(pattern))
			}
			return strings.Contains(s, pattern)
		}
		return re.MatchString(s)
	}

	matched := make([]bool, len(lines))
	matchCount := 0
	for i, s := range lines {
		isMatch := matchesLine(s)
		if flags.Invert {
			isMatch = !isMatch
		}
		matched[i] = isMatch
		if isMatch {
			matchCount++
		}
	}

	if flags.CountOnly {
		resp.FoundBlocks = []models.FoundBlock{{
			StartLineNumber: 0,
			Lines:           []string{strconv.Itoa(matchCount)},
		}}
		return resp, nil
	}

	before := flags.Before
	after := flags.After

	type rng struct{ start, end int }
	ranges := make([]rng, 0)
	currentOpen := false
	var cur rng

	for i := range lines {
		if matched[i] {
			start := i - before
			end := i + after
			if !currentOpen {
				cur = rng{start: start, end: end}
				currentOpen = true
			} else {
				if start <= cur.end+1 {
					if end > cur.end {
						cur.end = end
					}
				} else {
					ranges = append(ranges, cur)
					cur = rng{start: start, end: end}
				}
			}
		}
	}
	if currentOpen {
		ranges = append(ranges, cur)
	}

	if len(ranges) == 0 {
		return resp, nil
	}

	for _, rg := range ranges {
		start := rg.start
		end := rg.end

		neededBefore := 0
		if start < 0 {
			neededBefore = -start
		}
		neededAfter := 0
		if end >= len(lines) {
			neededAfter = end - (len(lines) - 1)
		}

		useBefore := neededBefore
		if useBefore > len(req.BeforeContext) {
			useBefore = len(req.BeforeContext)
		}
		useAfter := neededAfter
		if useAfter > len(req.AfterContext) {
			useAfter = len(req.AfterContext)
		}

		blockLines := make([]string, 0)

		if useBefore > 0 {
			startIdx := len(req.BeforeContext) - useBefore
			blockLines = append(blockLines, req.BeforeContext[startIdx:]...)
		}

		clampedStart := start
		if clampedStart < 0 {
			clampedStart = 0
		}
		clampedEnd := end
		if clampedEnd > len(lines)-1 {
			clampedEnd = len(lines) - 1
		}
		if clampedStart <= clampedEnd {
			blockLines = append(blockLines, lines[clampedStart:clampedEnd+1]...)
		}

		if useAfter > 0 {
			blockLines = append(blockLines, req.AfterContext[:useAfter]...)
		}

		blockStartAbs := req.StartLineNumber + clampedStart - useBefore
		if blockStartAbs < 0 {
			blockStartAbs = 0
		}

		if flags.PrintNumbers {
			withNums := make([]string, len(blockLines))
			for i := range blockLines {
				ln := blockStartAbs + i
				withNums[i] = fmt.Sprintf("%d:%s", ln, blockLines[i])
			}
			blockLines = withNums
		}

		resp.FoundBlocks = append(resp.FoundBlocks, models.FoundBlock{
			StartLineNumber: blockStartAbs,
			Lines:           blockLines,
		})
	}

	return resp, nil
}
