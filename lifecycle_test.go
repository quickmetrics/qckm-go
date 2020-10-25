package qm

import (
	"testing"
)

func TestLifecycle(t *testing.T) {
	Init(Options{
		ApiKey:       "testkey",
		Verbose:      true,
		MaxBatchWait: 20,
	})

	Event("api.request.time", 42)
	Event("api.request.time", 41)
	FlushEventsSync()

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
