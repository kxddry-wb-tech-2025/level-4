package log

import "time"

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Entry struct {
	Level   Level
	Message string
	Time    time.Time
	Error   error
	Meta    map[string]any
}
