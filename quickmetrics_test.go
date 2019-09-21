package qm

import (
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	Init("testkey")
}

func TestEvent(t *testing.T) {
	Event("testing.go.event.2", 12)
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
