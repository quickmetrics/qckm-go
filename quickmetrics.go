package qm

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	authHeader     = "x-qm-key"
	ingestEndpoint = "https://in.qckm.io/v2/events"
	queryEndpoint  = "https://api.quickmetrics.io/v2/query"
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
	Name   string                 `json:"name"`
	Fields map[string]interface{} `json:"data"`
}

type Fields map[string]interface{}

var batcher *batch

var httpClient *fasthttp.Client
var httpQueryClient *fasthttp.Client

func Init(opt Options) {
	clientKey = &opt.ApiKey
	isEnabled = true
	isVerbose = opt.Verbose

	batcher = newBatcher(opt.MaxBatchSize, opt.MaxBatchWait, opt.BatchWorkers)

	httpClient = &fasthttp.Client{
		Name:               "go-qckm",
		MaxConnWaitTimeout: 15 * time.Second,
		WriteTimeout:       15 * time.Second,
		ReadTimeout:        time.Second,
	}

	httpQueryClient = &fasthttp.Client{
		Name:               "go-qckm",
		MaxConnWaitTimeout: 15 * time.Second,
		WriteTimeout:       15 * time.Second,
		ReadTimeout:        5 * time.Minute,
	}
}

func SetEnabled(enable bool) {
	isEnabled = enable
}

// Event sends an event with a name and optional fields
func Event(name string, fields map[string]interface{}) {
	if isEnabled {
		batcher.add(event{
			Name:   name,
			Fields: fields,
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
	sendBatch(ee)

	// let waitgroup know we're done with these events
	for i := 0; i < len(ee); i++ {
		wg.Done()
	}
}

func sendBatch(ee []event) {
	if clientKey == nil || *clientKey == "" {
		log.Println("[ERROR] missing api key, please run qm.Init() first")
		return
	}

	body, _ := json.Marshal(ee)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(ingestEndpoint)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.Header.Set(authHeader, *clientKey)
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	resp.SkipBody = !isVerbose

	start := time.Now()

	httpClient.Do(req, resp)

	if isVerbose {
		if resp != nil {
			log.Printf(
				`[INFO] Request Finished in %vms. endpoint: "%v" status: %v addr: %v body: "%v"`,
				time.Since(start).Milliseconds(),
				req.URI(),
				resp.StatusCode(),
				resp.RemoteAddr(),
				string(resp.Body()),
			)
		}
	}
}
