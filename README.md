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

const CustomError = "custom-error"

func main() {
    reporter := event_reporter.New()
    err := reporter.Add(CustomError, &event_reporter.ReportConfig{
        Subject:   "Сustom error",
        MaxCount:  25,
        ResetTime: 20 * time.Second,
        Senders:   []event_reporter.Sender{&TestSender{}},
    })
    
    if err != nil {
        fmt.Println(err)
        return
    }
    
    var wg sync.WaitGroup
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            time.Sleep(time.Duration(rand.Intn(2-1)+1) * time.Second)
            fmt.Println("error happened")
            reporter.Publish(CustomError, fmt.Errorf("error happened"))
        }
    }()
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            time.Sleep(time.Duration(rand.Intn(2-1)+1) * time.Second)
            fmt.Println("error happened")
            reporter.Publish(CustomError, fmt.Errorf("error happened"))
        }
    }()
    
    wg.Wait()
}

type TestSender struct {
}

func (ts *TestSender) Send(ctx context.Context, subject string, msg string) error {
    fmt.Println(subject, msg)
    return nil
}
```