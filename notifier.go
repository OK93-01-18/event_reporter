package event_reporter

import (
	"context"
	ntf "github.com/nikoksr/notify"
)

type Notifier interface {
	UseServices(...Sender)
	Send(context.Context, string, string) error
}

type Notify struct {
	*ntf.Notify
}

func (n *Notify) UseServices(services ...Sender) {
	ntfServices := make([]ntf.Notifier, len(services))
	for _, service := range services {
		ntfServices = append(ntfServices, service)
	}
	n.Notify.UseServices(ntfServices...)
}

func (n *Notify) Send(ctx context.Context, sbj string, msg string) error {
	return n.Notify.Send(ctx, sbj, msg)
}

func NewNotify() Notifier {
	return &Notify{ntf.New()}
}
