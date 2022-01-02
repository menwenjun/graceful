package graceful

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func setup() {
	startOnce = sync.Once{}
}

func TestManager(t *testing.T) {
	setup()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	_ = GetManager()
}

func TestRunningJob(t *testing.T) {
	setup()
	var count int32 = 0
	m := NewManager()

	// Add job
	m.AddRunningJob(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				atomic.AddInt32(&count, 1)
				time.Sleep(100 * time.Millisecond)
			}
		}
	})

	go func() {
		time.Sleep(50 * time.Millisecond)
		m.doGracefulShutdown()
	}()

	<-m.Done()

	if atomic.LoadInt32(&count) != 1 {
		t.Errorf("count error")
	}
}

func TestRunningAndShutdownJob(t *testing.T) {
	setup()
	var count int32 = 0
	m := NewManager()

	// Add job
	m.AddRunningJob(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				atomic.AddInt32(&count, 1)
				time.Sleep(100 * time.Millisecond)
			}
		}
	})

	m.AddShutdownJob(func() error {
		atomic.AddInt32(&count, 1)
		return nil
	})

	go func() {
		time.Sleep(50 * time.Millisecond)
		m.doGracefulShutdown()
	}()

	<-m.Done()

	if atomic.LoadInt32(&count) != 2 {
		t.Errorf("count error: %v", atomic.LoadInt32(&count))
	}
}