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
					Type:        TypeNumber,
					Aggregation: AggCount,
				},
			},
		},
	})

	log.Println(err, res[0].Stats, res[0].Results[0])
}
