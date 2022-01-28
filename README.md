event_reporter
======

Package event_reporter sending notification with accumulation of error counter and distribution interval


### Install

	go get github.com/OK93-01-18/event_reporter

### Example
```go

import (
    "context"
    "fmt"
    "github.com/ok93-01-18/event_reporter"
    "math/rand"
    "sync"
    "time"
)

// event name
const CustomError = "custom-error"

func main() {
	
	// create reporter instance
    reporter := event_reporter.New()
    err := reporter.Add(CustomError, &event_reporter.ReportConfig{
        Subject:   "Сustom error", // subject of message
        MaxCount:  25, // max count event execute before send
        ResetTime: 20 * time.Second, // reset MaxCount interval
        Senders:   []event_reporter.Sender{&TestSender{}}, // senders slice
    })
    
    if err != nil {
        fmt.Println(err)
        return
    }
    
	// create 2 coroutine with event execute in random time 
	
    var wg sync.WaitGroup
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            time.Sleep(time.Duration(rand.Intn(2-1)+1) * time.Second)
            fmt.Println("error happened")
            reporter.Publish(CustomError, "error happened")
        }
    }()
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            time.Sleep(time.Duration(rand.Intn(2-1)+1) * time.Second)
            fmt.Println("error happened")
            reporter.Publish(CustomError, "error happened")
        }
    }()
    
    wg.Wait()
}

// TestSender is a simple sender for error reporter
type TestSender struct {
}

// Send is method for sending message
func (ts *TestSender) Send(ctx context.Context, subject string, msg string) error {
    fmt.Println(subject, msg)
    return nil
}
```