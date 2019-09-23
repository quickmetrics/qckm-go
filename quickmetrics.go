package qm

import (
	"encoding/json"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	authHeader = "x-qm-key"
	endpoint   = "https://qckm.io/json"
)

var clientKey *string
var isEnabled bool

type event struct {
	Name      string  `json:"name"`
	Value     float32 `json:"value"`
	Dimension string  `json:"dimension,omitempty"`
}

func Init(apiKey string) {
	clientKey = &apiKey
	isEnabled = true
}

func SetEnabled(enable bool) {
	isEnabled = enable
}

// Event sends a metric with values
func Event(name string, value float32) {
	go sendEvent(event{
		Name:  name,
		Value: value,
	})
}

// EventDimensions sends a name, secondary dimension and value
func EventDimension(name string, dimension string, value float32) {
	go sendEvent(event{
		Name:      name,
		Dimension: dimension,
		Value:     value,
	})
}

// Time is a helper to time functions
// pass it the star time and it'll calculate the
// duration. Alternatively pass it the current time
// and defer it at the start of your function like so:
// defer qm.Time(time.Now(), "your.metric")
func Time(start time.Time, name string) {
	dur := float32(time.Since(start).Milliseconds())
	go sendEvent(event{
		Name:  name,
		Value: dur,
	})
}

// TimeDimension is a helper to time functions
// pass it the star time and it'll calculate the
// duration. It also supports a secondary dimension
func TimeDimension(start time.Time, name string, dimension string) {
	dur := float32(time.Since(start).Milliseconds())
	go sendEvent(event{
		Name:      name,
		Dimension: dimension,
		Value:     dur,
	})
}

func sendEvent(e event) {
	if !isEnabled {
		return
	}

	if clientKey == nil || *clientKey == "" {
		log.Println("missing api key, please run qm.Init() first")
		return
	}

	body, _ := json.Marshal(e)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(endpoint)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.Header.Set(authHeader, *clientKey)
	req.SetBody(body)

	resp := fasthttp.AcquireResponse()
	resp.SkipBody = true

	client := &fasthttp.Client{
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
	}
	client.Do(req, resp)
}
