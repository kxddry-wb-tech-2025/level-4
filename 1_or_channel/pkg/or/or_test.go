package or

import (
	"testing"
	"time"
)

// go test or

func sigAfter(d time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		time.Sleep(d)
	}()
	return ch
}

func TestOr_NoChannels(t *testing.T) {
	if ch := Or(); ch != nil {
		t.Errorf("expected nil channel when no input channels are provided")
	}
}

func TestOr_SingleChannel(t *testing.T) {
	sig := sigAfter(10 * time.Millisecond)
	start := time.Now()

	<-Or(sig)

	if elapsed := time.Since(start); elapsed < 10*time.Millisecond {
		t.Errorf("Or closed too early with single channel")
	}
}

func TestOr_TwoChannels_FirstFires(t *testing.T) {
	sig1 := sigAfter(5 * time.Millisecond)
	sig2 := sigAfter(50 * time.Millisecond)
	start := time.Now()

	<-Or(sig1, sig2)

	if elapsed := time.Since(start); elapsed > 20*time.Millisecond {
		t.Errorf("Or did not close early enough, elapsed=%v", elapsed)
	}
}

func TestOr_TwoChannels_SecondFires(t *testing.T) {
	sig1 := sigAfter(50 * time.Millisecond)
	sig2 := sigAfter(5 * time.Millisecond)
	start := time.Now()

	<-Or(sig1, sig2)

	if elapsed := time.Since(start); elapsed > 20*time.Millisecond {
		t.Errorf("Or did not close early enough, elapsed=%v", elapsed)
	}
}

func TestOr_MultipleChannels(t *testing.T) {
	sigs := []<-chan struct{}{
		sigAfter(80 * time.Millisecond),
		sigAfter(150 * time.Millisecond),
		sigAfter(10 * time.Millisecond),
		sigAfter(200 * time.Millisecond),
	}
	start := time.Now()

	<-Or(sigs...)

	if elapsed := time.Since(start); elapsed > 30*time.Millisecond {
		t.Errorf("Or did not close at earliest signal, elapsed=%v", elapsed)
	}
}

func TestOr_Stress(t *testing.T) {
	// Many channels, one closes quickly
	channels := make([]<-chan struct{}, 1000)
	for i := 0; i < 999; i++ {
		ch := make(chan struct{})
		channels[i] = ch
	}
	// last one closes after short time
	channels[999] = sigAfter(5 * time.Millisecond)

	select {
	case <-Or(channels...):
	case <-time.After(50 * time.Millisecond):
		t.Errorf("Or timed out with large number of channels")
	}
}

func TestOr_ClosePropagation(t *testing.T) {
	// if one channel is already closed, Or should return immediately
	alreadyClosed := make(chan struct{})
	close(alreadyClosed)

	start := time.Now()
	<-Or(alreadyClosed, sigAfter(1*time.Hour))

	if elapsed := time.Since(start); elapsed > 10*time.Millisecond {
		t.Errorf("Or did not return immediately with pre-closed channel, elapsed=%v", elapsed)
	}
}
