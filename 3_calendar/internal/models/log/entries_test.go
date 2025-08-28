package log

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	err := errors.New("boom")
	e := Error(err, "msg", map[string]any{"a": 1})
	if e.Level != LevelError || e.Message != "msg" || e.Error == nil || e.Meta["a"].(int) != 1 {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestInfo(t *testing.T) {
	e := Info("ok", map[string]any{"k": "v"})
	if e.Level != LevelInfo || e.Message != "ok" || e.Meta["k"].(string) != "v" {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestWarn(t *testing.T) {
	e := Warn("w", nil)
	if e.Level != LevelWarn || e.Message != "w" {
		t.Fatalf("unexpected entry: %+v", e)
	}
}

func TestDebug(t *testing.T) {
	e := Debug("d", nil)
	if e.Level != LevelDebug || e.Message != "d" {
		t.Fatalf("unexpected entry: %+v", e)
	}
}
