package event_reporter

import "context"

type Sender interface {
	Send(context.Context, string, string) error
}
