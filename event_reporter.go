package event_reporter

import (
	"container/ring"
	"context"
	"fmt"
	"sync"
	"time"
)

// Mode is type of event behavior
type Mode uint

const (
	// AlwaysNotify is sending when logBuffer not empty
	AlwaysNotify Mode = 1

	// BufferFull is sending when logBuffer full
	BufferFull Mode = 2
)

// ReportConfig is config of report event
type ReportConfig struct {
	Subject   string
	LogSize   int
	ResetTime time.Duration
	Senders   []Sender
	Mode      Mode
}

type event struct {
	config    *ReportConfig
	ticker    *time.Ticker
	notifier  Notifier
	logBuffer *ring.Ring
}

type EventError struct {
	Topic string
	Error error
}

// EventReporter is central struct for register events for report
type EventReporter struct {
	sync.RWMutex
	events    map[string]event
	errorChan chan EventError
}

// Add is method for adding type of event
func (er *EventReporter) Add(topic string, conf *ReportConfig) error {
	er.RLock()
	_, ok := er.events[topic]
	er.RUnlock()

	if ok {
		return fmt.Errorf("event %s already exists", topic)
	}

	ticker := time.NewTicker(conf.ResetTime)

	notifier := NewNotify()
	notifier.UseSenders(conf.Senders...)
	newEvent := event{
		config:    conf,
		ticker:    ticker,
		notifier:  notifier,
		logBuffer: ring.New(conf.LogSize),
	}
	er.Lock()
	er.events[topic] = newEvent
	er.Unlock()

	go func() {
		for {
			select {
			case <-ticker.C:
				er.Lock()
				foundEvent, _ := er.events[topic]

				count := 0
				msg := ""
				foundEvent.logBuffer.Do(func(p interface{}) {
					if p != nil {
						msg += p.(string) + "\n"
						count++
					}
				})

				if msg != "" && (foundEvent.config.Mode == AlwaysNotify ||
					(foundEvent.config.Mode == BufferFull && count == foundEvent.config.LogSize)) {
					go func() {
						err := foundEvent.notifier.Send(context.Background(), foundEvent.config.Subject, msg)
						if err != nil {
							er.errorChan <- EventError{
								Topic: topic,
								Error: err,
							}
						}
					}()
				}

				foundEvent.logBuffer = ring.New(conf.LogSize)
				er.events[topic] = foundEvent
				er.Unlock()
			}
		}
	}()

	return nil
}

// Publish is method for execute event
func (er *EventReporter) Publish(topic string, msg string) {
	er.Lock()
	defer er.Unlock()

	foundEvent, ok := er.events[topic]
	if !ok {
		return
	}

	foundEvent.logBuffer.Value = "[" + time.Now().Format("01-02-2006 15:04:05") + "] " + msg
	foundEvent.logBuffer = foundEvent.logBuffer.Next()
	er.events[topic] = foundEvent
}

func (er *EventReporter) GetErrorChan() chan EventError {
	return er.errorChan
}

func New() *EventReporter {
	return &EventReporter{
		events:    make(map[string]event),
		errorChan: make(chan EventError),
	}
}
