package event_reporter

import (
	"context"
	ntf "github.com/nikoksr/notify"
)

// Notifier is the interface for wrapper notifiers
type Notifier interface {
	UseSenders(...Sender)
	Send(context.Context, string, string) error
}

// Notify wrapper of nikoksr/notify package
type Notify struct {
	*ntf.Notify
}

// UseSenders adds the given sender(s) to the Notifier's senders list.
func (n *Notify) UseSenders(senders ...Sender) {
	ntfServices := make([]ntf.Notifier, len(senders))
	for _, service := range senders {
		ntfServices = append(ntfServices, service)
	}
	n.Notify.UseServices(ntfServices...)
}

// Send calls the underlying notification services to send the given subject and message to their respective endpoints.
func (n *Notify) Send(ctx context.Context, sbj string, msg string) error {
	return n.Notify.Send(ctx, sbj, msg)
}

// NewNotify returns a new instance of Notifier
func NewNotify() Notifier {
	return &Notify{ntf.New()}
}
