package qm

import (
	"testing"
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

	Event("testEvent", Fields{"string": "lalala", "number": 123, "bool": true})
	Event("testEvent", Fields{"string": "lalala", "number": 123, "bool": true})
	Event("testEvent", Fields{"string": "lalala", "number": 123, "bool": true})
	Event("testEvent", Fields{"string": "lalala", "number": 123, "bool": true})
	FlushEventsSync()
	Event("testEvent", Fields{"string": "lalala", "number": 123, "bool": true})
	Event("testEvent", Fields{"string": "lalala", "number": 123, "bool": true})
	Event("testEvent", Fields{"string": "lalala", "number": 123, "bool": true})

	FlushEventsSync()
	FlushEventsSync()
}

func TestEvent(t *testing.T) {
	Event("testEvent", Fields{"string": "lalala", "number": 123, "bool": true})
}

func TestEnable(t *testing.T) {
	// we disable event sending
	SetEnabled(false)
	// the api key is not inizalized
	// but that shouldn't matter since
	// the sending is disabled
	Event("testEvent", Fields{"string": "lalala", "number": 123, "bool": true})
}
