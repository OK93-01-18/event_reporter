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
    "time"
)

// CustomError event name
const CustomError = "custom-error"

func main() {

    // create reporter instance
    reporter := event_reporter.New()
    err := reporter.Add(CustomError, &event_reporter.ReportConfig{
        Subject:   "Сustom error",                         // subject of message
		LogSize:   25,                                     // event log max size
        ResetTime: 20 * time.Second,                       // send interval time
        Senders:   []event_reporter.Sender{&TestSender{}}, // senders slice
    })

    if err != nil {
        fmt.Println(err)
        return
    }

    // create 2 coroutine with event execute in random time
    go func() {
        for {
            time.Sleep(time.Duration(rand.Intn(2-1)+1) * time.Second)
            reporter.Publish(CustomError, "["+time.Now().Format("01-02-2006 15:04:05")+"] error happened")
        }
    }()

    go func() {
        for {
            time.Sleep(time.Duration(rand.Intn(2-1)+1) * time.Second)
            reporter.Publish(CustomError, "["+time.Now().Format("01-02-2006 15:04:05")+"] error happened 2")
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

// TestSender is a simple sender for error reporter
type TestSender struct {
}

// Send is method for sending message
func (ts *TestSender) Send(ctx context.Context, subject string, msg string) error {
    n := rand.Intn(10-1) + 1
    if n > 5 { // random success sending
        fmt.Println(subject, msg)
        return nil
    }
    return fmt.Errorf("error send to console")
}
```