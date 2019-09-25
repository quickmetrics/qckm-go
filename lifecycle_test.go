package qm

import (
	"testing"
	"time"
)

func TestLifecycle(t *testing.T) {
	Init(Options{
		ApiKey:       "testkey",
		Verbose:      true,
		MaxBatchWait: 20,
	})

	Event("event", 42)
	Event("event", 41)
	EventDimension("event", "dim", 43)
	EventDimension("event", "dim", 44)
	FlushEvents()

	// // wait
	time.Sleep(time.Second * 5)
}

func BenchmarkBatching(b *testing.B) {
	Init(Options{
		ApiKey:       "testkey",
		Verbose:      true,
		MaxBatchWait: 3,
		MaxBatchSize: 10000,
	})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Event("test.benchmark", 123)
	}
}
