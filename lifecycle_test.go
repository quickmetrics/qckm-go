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

	Event("testEvent", Fields{"string": "lalala", "number": 123, "bool": true})
	Event("testEvent", Fields{"string": "lalala", "number": 123, "bool": true})
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
		Event("testEvent", Fields{"string": "lalala", "number": 123, "bool": true})
	}
}
