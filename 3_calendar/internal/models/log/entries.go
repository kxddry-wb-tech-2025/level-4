package log

import "time"

func Error(err error, msg string, meta map[string]any) Entry {
	return Entry{
		Level:   LevelError,
		Message: msg,
		Time:    time.Now(),
		Error:   err,
		Meta:    meta,
	}
}

func Info(msg string, meta map[string]any) Entry {
	return Entry{
		Level:   LevelInfo,
		Message: msg,
		Time:    time.Now(),
		Meta:    meta,
	}
}

func Warn(msg string, meta map[string]any) Entry {
	return Entry{
		Level:   LevelWarn,
		Message: msg,
		Time:    time.Now(),
		Meta:    meta,
	}
}

func Debug(msg string, meta map[string]any) Entry {
	return Entry{
		Level:   LevelDebug,
		Message: msg,
		Time:    time.Now(),
		Meta:    meta,
	}
}
