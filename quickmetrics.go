package qm

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	authHeader    = "x-qm-key"
	endpoint      = "https://qckm.io/json"
	batchEndpoint = "https://qckm.io/list"
)

var clientKey *string
var isEnabled bool
var isVerbose bool

// initialization options
type Options struct {
	ApiKey       string
	MaxBatchSize int
	MaxBatchWait time.Duration
	BatchWorkers int
	Verbose      bool
}

// event holds a single event
// ready to be sent to the qckm server
type event struct {
	Name      string
	Value     float32
	Timestamp time.Time
	Dimension string
}

// list holds a slice of listItems which can
// contain multiple events for batching
type list []listItem

type listItem struct {
	Name      string          `json:"name"`
	Dimension string          `json:"dimension,omitempty"`
	Values    [][]interface{} `json:"values"`
}

var batcher *batch

func Init(opt Options) {
	clientKey = &opt.ApiKey
	isEnabled = true
	isVerbose = opt.Verbose

	batcher = newBatcher(opt.MaxBatchSize, opt.MaxBatchWait, opt.BatchWorkers)
}

func SetEnabled(enable bool) {
	isEnabled = enable
}

// Event sends a metric with values
func Event(name string, value float32) {
	if isEnabled {
		batcher.add(event{
			Name:      name,
			Value:     value,
			Timestamp: time.Now().UTC(),
		})
	}
}

// EventDimensions sends a name, secondary dimension and value
func EventDimension(name string, dimension string, value float32) {
	if isEnabled {
		batcher.add(event{
			Name:      name,
			Dimension: dimension,
			Value:     value,
			Timestamp: time.Now().UTC(),
		})
	}
}

// Time is a helper to time functions
// pass it the star time and it'll calculate the
// duration. Alternatively pass it the current time
// and defer it at the start of your function like so:
// defer qm.Time(time.Now(), "your.metric")
func Time(start time.Time, name string) {
	if isEnabled {
		dur := float32(time.Since(start).Milliseconds())
		batcher.add(event{
			Name:      name,
			Value:     dur,
			Timestamp: time.Now().UTC(),
		})
	}
}

// TimeDimension is a helper to time functions
// pass it the star time and it'll calculate the
// duration. It also supports a secondary dimension
func TimeDimension(start time.Time, name string, dimension string) {
	if isEnabled {
		dur := float32(time.Since(start).Milliseconds())
		batcher.add(event{
			Name:      name,
			Dimension: dimension,
			Value:     dur,
			Timestamp: time.Now().UTC(),
		})
	}
}

// FlushEvents processes any events left in the queue
// and sends them to the qickmetrics server.
func FlushEvents() {
	batcher.flush()
}

// FlushEventsSync processes any events left in the queue
// and sends them to the qickmetrics server. This function
// is blocking until events are sent to ensure that the system
// doesn't shut down before then
func FlushEventsSync() {
	// flush any pending items to process them
	batcher.flush()
	// wait for all processing and network to finish
	batcher.wait()
}

func processBatch(ee []event, wg *sync.WaitGroup) {
	start := time.Now()

	holder := map[string]map[string][][]interface{}{}

	// we sort the events into a map of metric name and dimension
	for _, e := range ee {
		if holder[e.Name] == nil {
			holder[e.Name] = map[string][][]interface{}{}
		}
		if holder[e.Name][e.Dimension] == nil {
			holder[e.Name][e.Dimension] = [][]interface{}{}
		}
		holder[e.Name][e.Dimension] = append(holder[e.Name][e.Dimension], []interface{}{e.Timestamp.Unix(), e.Value})
	}

	output := list{}

	// then we process them into an array of data items
	for metricName, dimension := range holder {
		for dimensionName, values := range dimension {
			output = append(output, listItem{
				Name:      metricName,
				Dimension: dimensionName,
				Values:    values,
			})
		}
	}

	if isVerbose {
		log.Printf("[INFO] qckm-go: processed %v events in %v", len(ee), time.Since(start))
	}

	sendBatch(output)

	// let waitgroup know we're done with these events
	for i := 0; i < len(ee); i++ {
		wg.Done()
	}
}

func sendBatch(l list) {
	if clientKey == nil || *clientKey == "" {
		log.Println("[ERROR] missing api key, please run qm.Init() first")
		return
	}

	body, _ := json.Marshal(l)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(batchEndpoint)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.Header.Set(authHeader, *clientKey)
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()
	resp.SkipBody = true

	client := &fasthttp.Client{
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  time.Second,
	}

	start := time.Now()

	client.Do(req, resp)

	if isVerbose {
		if resp != nil {
			log.Printf(
				"[INFO] Events received by %v (status %v\n) in %vms",
				batchEndpoint,
				resp.StatusCode(),
				time.Since(start).Milliseconds(),
			)
		}
	}
}
