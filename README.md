## A quickmetrics client for go
#### Installation
`go get github.com/quickmetrics/go`

#### Setup
Initialize with your API key before sending events. You'll only have to do this once in your app lifecycle.
`qm.Init("YOUR_API_KEY")`


#### Send events

`qm.Event("your.event", 123.456)`


#### Full Example
```
package main

import "github.com/quicmetrics/go"

func init() {
  qm.Init("YOUR_API_KEY")
}

func main() {
  qm.Event("hello.world", 1)
}

```
And that's it!

For more info on naming conventions and examples check out our docs at https://app.quickmetrics.io/docs
