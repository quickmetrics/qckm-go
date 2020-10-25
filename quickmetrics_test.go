package qm

import (
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	Init(Options{
		ApiKey:  "testkey",
		Verbose: true,
	})
}

func TestFlush(t *testing.T) {
	Init(Options{
		ApiKey:       "testkey",
		Verbose:      true,
		MaxBatchSize: 10000,
		MaxBatchWait: 60,
	})

	Event("testing.go.flush", 42)
	Event("testing.go.flush", 42)
	Event("testing.go.flush", 42)
	Event("testing.go.flush", 42)
	Event("testing.go.flush", 42)

	FlushEventsSync()
}

func TestEvent(t *testing.T) {
	Event("testing.go.event.2", 42)
}

func TestEnable(t *testing.T) {
	// we disable event sending
	SetEnabled(false)
	// the api key is not inizalized
	// but that shouldn't matter since
	// the sending is disabled
	Event("testing.go", 42)
}

func TestEventDimension(t *testing.T) {
	EventDimension("testing.go.event", "test.dimension", 12)
}

func TestTime(t *testing.T) {
	Time(time.Now(), "testing.go.time")
}

func TestTimeDimension(t *testing.T) {
	TimeDimension(time.Now(), "testing.go.time", "test.dimension")
}
