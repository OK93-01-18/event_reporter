package event_reporter

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ReportConfig is config of report event
type ReportConfig struct {
	Subject   string
	MaxCount  int
	ResetTime time.Duration
	Senders   []Sender
}

type event struct {
	config   *ReportConfig
	ticker   *time.Ticker
	notifier Notifier
	count    int
}

// EventReporter is central struct for register events for report
type EventReporter struct {
	events map[string]event
	sync.RWMutex
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
		config:   conf,
		ticker:   ticker,
		notifier: notifier,
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
				foundEvent.count = 0
				er.events[topic] = foundEvent
				er.Unlock()
			}
		}
	}()

	return nil
}

// Publish is method for execute event
func (er *EventReporter) Publish(topic string, msg string) error {

	er.Lock()
	defer er.Unlock()

	var err error

	foundEvent, ok := er.events[topic]
	if !ok {
		return err
	}

	foundEvent.count++

	if foundEvent.count == foundEvent.config.MaxCount {
		foundEvent.count = 0
		err = foundEvent.notifier.Send(context.Background(), foundEvent.config.Subject, msg)
		foundEvent.ticker.Reset(foundEvent.config.ResetTime)
	}

	er.events[topic] = foundEvent
	return err
}

func New() *EventReporter {
	return &EventReporter{events: make(map[string]event)}
}
