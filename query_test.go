package qm

import (
	"log"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	Init(Options{
		ApiKey: "testkey",
	})

	res, err := Query([]QueryRequest{
		{
			Name: "testEvent",
			From: time.Now().Add(-time.Hour * 24 * 30),
			To:   time.Now(),
			Items: []Item{
				{
					Type:        TypeTimeseries,
					Aggregation: AggCount,
					Interval:    "24h",
				},
			},
		},
	})

	log.Println(err, res)
}
