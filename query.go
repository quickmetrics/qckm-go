package qm

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/valyala/fasthttp"
)

// QueryRequest event data
type QueryRequest struct {
	Name       string      `json:"name"`
	From       time.Time   `json:"from"`
	To         time.Time   `json:"to"`
	Items      []Item      `json:"items"`
	Conditions []Condition `json:"conditions"`
}

// ItemType what to query for
type ItemType string

const (
	TypeList       ItemType = "list"
	TypeNumber     ItemType = "number"
	TypeTimeseries ItemType = "timeseries"
)

// Agg how to aggregate numbers
type Agg string

const (
	AggAvg           Agg = "avg"
	AggSum           Agg = "sum"
	AggSumCumulative Agg = "sum_cumulative"
	AggCount         Agg = "count"
	AggCountUnique   Agg = "count_unique"
	AggMin           Agg = "min"
	AggMax           Agg = "max"
)

// Item within query
type Item struct {
	Key         string        `json:"key,omitempty"`             // the index of the column in the data
	Type        ItemType      `json:"type"`                      // list, number, timeseries
	Aggregation Agg           `json:"aggregation"`               // avg, sum, sum_cumulative, count, count_unique, min, max
	Interval    time.Duration `json:"interval,omitempty"`        // (type=timeseries) 5m, 4h
	ListOrder   string        `json:"listOrder,omitempty"`       // (type=list) asc, desc
	ListLimit   int           `json:"listLimit,omitempty"`       // amount of items to be returned in list
	ExcludeNull bool          `json:"listExcludeNull,omitempty"` // whether or not to count null values
}

type Condition struct {
	Key       string      `json:"key"`
	Operation string      `json:"op"`
	Value     interface{} `json:"value"`
}

type Result struct {
	Query   QueryRequest        `json:"query"`
	Stats   Stats               `json:"stats"`
	Results map[int]interface{} `json:"results"`
	Error   string              `json:"error"`
}

type Stats struct {
	EventsProcessed        int     `json:"eventsProcessed"`
	EventsAnalyzed         int     `json:"eventsAnalyzed"`
	MillionEventsPerSecond float64 `json:"meps"`
	DurationTotal          int64   `json:"durationTotal"`
}

func Query(data []QueryRequest) ([]Result, error) {
	if clientKey == nil || *clientKey == "" {
		return []Result{}, errors.New("missing api key, please run qm.Init() first")
	}

	body, err := json.Marshal(data)
	if err != nil {
		return []Result{}, err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(queryEndpoint)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.Header.Set(authHeader, *clientKey)
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = httpClient.Do(req, resp)
	if err != nil {
		return []Result{}, err
	}

	var out []Result
	err = json.Unmarshal(resp.Body(), &out)

	return out, nil
}
