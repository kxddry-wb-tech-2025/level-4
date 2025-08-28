package log

import "time"

// Error creates a new error log entry
func Error(err error, msg string, meta map[string]any) Entry {
	return Entry{
		Level:   LevelError,
		Message: msg,
		Time:    time.Now(),
		Error:   err,
		Meta:    meta,
	}
}

// Info creates a new info log entry
func Info(msg string, meta map[string]any) Entry {
	return Entry{
		Level:   LevelInfo,
		Message: msg,
		Time:    time.Now(),
		Meta:    meta,
	}
}

// Warn creates a new warn log entry
func Warn(msg string, meta map[string]any) Entry {
	return Entry{
		Level:   LevelWarn,
		Message: msg,
		Time:    time.Now(),
		Meta:    meta,
	}
}

// Debug creates a new debug log entry
func Debug(msg string, meta map[string]any) Entry {
	return Entry{
		Level:   LevelDebug,
		Message: msg,
		Time:    time.Now(),
		Meta:    meta,
	}
}
