package notify

import (
	"calendar/internal/models"
	"context"
	"testing"
	"time"
)

type fakeWorker struct {
	add int
	del int
}

func (f *fakeWorker) AddNotification(ctx context.Context, n models.CreateNotificationRequest) error {
	f.add++
	return nil
}
func (f *fakeWorker) DeleteAllNotificationsByEventID(ctx context.Context, id string) error {
	f.del += 2
	return nil
}
func (f *fakeWorker) DeleteNotificationByID(ctx context.Context, id string) error {
	f.del++
	return nil
}

func TestProcess_CreateAndDelete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	jobs := make(chan any, 2)
	fw := &fakeWorker{}
	svc := NewService(ctx, fw)

	go svc.Process(jobs)

	jobs <- models.CreateNotificationRequest{}
	jobs <- models.DeleteNotificationsRequest{}

	// allow goroutine to process
	time.Sleep(50 * time.Millisecond)
	cancel()

	if fw.add == 0 || fw.del == 0 {
		t.Fatalf("expected both add and delete to be invoked: add=%d del=%d", fw.add, fw.del)
	}
}
