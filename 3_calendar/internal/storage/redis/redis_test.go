package redis

import (
	"context"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	goredis "github.com/redis/go-redis/v9"
)

func setup(t *testing.T) (*Storage, func()) {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}
	client := goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
	st := &Storage{client: client}
	return st, func() { client.Close(); mr.Close() }
}

func TestEnqueueAndPopDue(t *testing.T) {
	st, cleanup := setup(t)
	defer cleanup()
	ctx := context.Background()

	id1 := "n1"
	id2 := "n2"
	if err := st.Enqueue(ctx, id1, time.Now().Add(-time.Second)); err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	if err := st.Enqueue(ctx, id2, time.Now().Add(1200*time.Millisecond)); err != nil {
		t.Fatalf("enqueue: %v", err)
	}

	vals, err := st.PopDue(ctx, keyDueZset, 10)
	if err != nil {
		t.Fatalf("pop due: %v", err)
	}
	if len(vals) != 1 || vals[0] != id1 {
		t.Fatalf("unexpected vals: %v", vals)
	}

	// after some time second becomes due
	time.Sleep(1500 * time.Millisecond)
	vals, err = st.PopDue(ctx, keyDueZset, 10)
	if err != nil {
		t.Fatalf("pop due: %v", err)
	}
	if len(vals) != 1 || vals[0] != id2 {
		t.Fatalf("unexpected vals2: %v", vals)
	}
}

func TestRemove(t *testing.T) {
	st, cleanup := setup(t)
	defer cleanup()
	ctx := context.Background()
	id := "n1"
	if err := st.Enqueue(ctx, id, time.Now().Add(time.Second)); err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	if err := st.Remove(ctx, id); err != nil {
		t.Fatalf("remove: %v", err)
	}
}
