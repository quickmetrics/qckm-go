[![GoDoc](https://godoc.org/github.com/quickmetrics/qckm-go?status.svg)](https://godoc.org/github.com/quickmetrics/qckm-go)
## A quickmetrics client for go
#### Installation
`go get github.com/quickmetrics/qckm-go`

#### Full Example
```
package main

import "github.com/quickmetrics/qckm-go"

func init() {
  qm.Init(qm.Options{
    ApiKey: "YOUR_API_KEY"
  })
}

func main() {
  qm.Event("hello.world", 1)
}
```

#### Setup
Initialize with your API key before sending events. You'll only have to do this once in your app lifecycle.
```
  qm.Init(qm.Options{
    ApiKey:       "YOUR_API_KEY",
    MaxBatchSize: 500,
    MaxBatchWait: 5,
    BatchWorkers: 1,
    Verbose:      true,
  })
```


#### Send events

`qm.Event("your.event", 123.456)`

`qm.EventDimension("click.color", "blue", 1)`

`qm.Time(startTime, "response.time")`

`qm.TimeDimension(startTime, "response.time", "POST /login")`

And that's it!

For more info on naming conventions and examples check out our docs at https://app.quickmetrics.io/docs
