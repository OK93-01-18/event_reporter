EventReporter
======

Package EventReporter sending notification with accumulation of error counter and distribution interval


### Install

	go get github.com/OK93-01-18/EventReporter

### Example
```go

import (
    "context"
    "fmt"
    "github.com/ok93-01-18/EventReporter"
    "math/rand"
    "sync"
    "time"
)

const CustomError = "custom-error"

func main() {
    reporter := EventReporter.New()
    err := reporter.Add(CustomError, &EventReporter.ReportConfig{
        Subject:   "Сustom error",
        Message:   "Some custom error was happened",
        MaxCount:  25,
        ResetTime: 20 * time.Second,
        Senders:   []EventReporter.Sender{&TestSender{}},
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
            reporter.Publish(CustomError)
        }
    }()
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            time.Sleep(time.Duration(rand.Intn(2-1)+1) * time.Second)
            fmt.Println("error happened")
            reporter.Publish(CustomError)
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