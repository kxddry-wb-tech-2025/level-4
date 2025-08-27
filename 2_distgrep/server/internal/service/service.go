package service

import (
	"errors"
	"fmt"
	"grep-server/internal/models"
	"regexp"
	"strconv"
	"strings"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Grep(req models.Request) (models.Response, error) {
	resp := models.Response{TaskID: req.ID}

	lines := req.Lines
	pattern := req.Pattern
	flags := req.Flags

	if pattern == "" {
		return resp, errors.New("empty pattern")
	}

	// Build matcher based on flags
	var re *regexp.Regexp
	var err error
	if flags.FixedString {
		// handled via strings.Contains (with case handling below)
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

	// Identify matching indices
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

	// Build merged ranges with context (-B, -A)
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
				// Extend current range if overlapping/adjacent
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

	// Assemble blocks including cross-shard context
	for _, rg := range ranges {
		start := rg.start
		end := rg.end

		// Determine how many lines we need from external context
		neededBefore := 0
		if start < 0 {
			neededBefore = -start
		}
		neededAfter := 0
		if end >= len(lines) {
			neededAfter = end - (len(lines) - 1)
		}

		// Clamp to available context
		useBefore := neededBefore
		if useBefore > len(req.BeforeContext) {
			useBefore = len(req.BeforeContext)
		}
		useAfter := neededAfter
		if useAfter > len(req.AfterContext) {
			useAfter = len(req.AfterContext)
		}

		// Prepare block lines
		blockLines := make([]string, 0)

		// Include before context from the tail of BeforeContext
		if useBefore > 0 {
			startIdx := len(req.BeforeContext) - useBefore
			blockLines = append(blockLines, req.BeforeContext[startIdx:]...)
		}

		// Include in-chunk lines
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

		// Include after context from the head of AfterContext
		if useAfter > 0 {
			blockLines = append(blockLines, req.AfterContext[:useAfter]...)
		}

		// Compute block starting absolute line number
		blockStartAbs := req.StartLineNumber + clampedStart - useBefore
		if blockStartAbs < 0 {
			blockStartAbs = 0
		}

		// Optionally prefix with line numbers
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
