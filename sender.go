package event_reporter

import "context"

// Sender is the interface of sender any service
type Sender interface {
	Send(context.Context, string, string) error
}
