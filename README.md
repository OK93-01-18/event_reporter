event_reporter
======

Package event_reporter sending notification by interval with event count logic


### Install

	go get github.com/OK93-01-18/event_reporter

### Example
```go
import (
    "fmt"
    "github.com/ok93-01-18/event_reporter"
    "github.com/ok93-01-18/event_reporter/senders/mattermost"
    "math/rand"
    "time"
)

// CustomError event name
const CustomError = "custom-error"

func main() {
    mmost := mattermost.New("test-user", "https://test.webhook/")
    
    // create reporter instance
    reporter := event_reporter.New()
    err := reporter.Add(CustomError, &event_reporter.ReportConfig{
        Subject:   "Test event sender",            // subject of message
        LogSize:   25,                             // event log size
        ResetTime: 20 * time.Second,               // send interval
        Senders:   []event_reporter.Sender{mmost}, // senders slice,
        Mode:      event_reporter.AlwaysNotify,    // sending mode
    })
    
    if err != nil {
        fmt.Println(err)
        return
    }
    
    // create 2 coroutine with event execute in random time
    go func() {
        for {
            time.Sleep(time.Duration(rand.Intn(2-1)+1) * time.Second)
            reporter.Publish(CustomError, "error happened")
        }
    }()
    
    go func() {
        for {
            time.Sleep(time.Duration(rand.Intn(2-1)+1) * time.Second)
            reporter.Publish(CustomError, "error happened 2")
        }
    }()
    
    // read reporter error channel
    go func() {
        for err := range reporter.GetErrorChan() {
            fmt.Printf("%s, %s\n", err.Topic, err.Error.Error())
        }
    }()
    
    time.Sleep(60 * time.Second) // wait 60 sec before close main coroutine

}
```