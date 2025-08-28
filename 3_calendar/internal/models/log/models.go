package log

import "time"

// Level is the level of the log
type Level int

// Levels
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Entry is the log entry
type Entry struct {
	Level   Level
	Message string
	Time    time.Time
	Error   error
	Meta    map[string]any
}
