package event_reporter

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Config struct {
	Subject   string
	Message   string
	MaxCount  int
	ResetTime time.Duration
	Senders   []Sender
}

type Event struct {
	config   *Config
	ticker   *time.Ticker
	notifier Notifier
	count    int
}

type EventReporter struct {
	events map[string]Event
	sync.RWMutex
}

func (er *EventReporter) Add(topic string, conf *Config) error {
	er.RLock()
	_, ok := er.events[topic]
	er.RUnlock()

	if ok {
		return fmt.Errorf("event %s already exists", topic)
	}

	ticker := time.NewTicker(conf.ResetTime)

	notifier := NewNotify()
	notifier.UseServices(conf.Senders...)
	event := Event{
		config:   conf,
		ticker:   ticker,
		notifier: notifier,
	}
	er.Lock()
	er.events[topic] = event
	er.Unlock()

	go func() {
		for {
			select {
			case <-ticker.C:
				er.RLock()
				event, _ := er.events[topic]
				er.RUnlock()

				event.count = 0

				er.Lock()
				er.events[topic] = event
				er.Unlock()
			}
		}
	}()

	return nil
}

func (er *EventReporter) Publish(topic string) error {

	er.Lock()
	defer er.Unlock()

	var err error

	event, ok := er.events[topic]
	if !ok {
		return err
	}

	event.count++

	if event.count == event.config.MaxCount {
		event.count = 0
		err = event.notifier.Send(context.Background(), event.config.Message, event.config.Message)
		event.ticker.Reset(event.config.ResetTime)
	}

	er.events[topic] = event
	return err
}

func New() *EventReporter {
	return &EventReporter{events: make(map[string]Event)}
}
