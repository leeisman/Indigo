package sub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

const (
	syncMode                      = "sync"
	asyncMode                     = "async"
	defaultMaxOutstandingMessages = 20
)

type Subscriber struct {
	ctx          context.Context
	subscription *pubsub.Subscription
	pullMode     string
	errorHook    func(err error)
}

func NewSubscriber(ctx context.Context, subscription *pubsub.Subscription) *Subscriber {
	subscription.ReceiveSettings.MaxOutstandingMessages = defaultMaxOutstandingMessages
	return &Subscriber{
		ctx:          ctx,
		subscription: subscription,
	}
}

func (sb *Subscriber) Options(opts ...SubscriberOption) *Subscriber {
	for _, opt := range opts {
		opt(sb)
	}
	return sb
}

func (sb *Subscriber) Subscribe(process func(context.Context, []byte, string) error) {
	if sb.pullMode == syncMode {
		sb.subscription.ReceiveSettings.Synchronous = true
		sb.pullMessageSync(sb.ctx, process)
	}
	if sb.pullMode == asyncMode {
		sb.subscription.ReceiveSettings.Synchronous = false
		sb.pullMessageAsync(sb.ctx, process)
	}
}

func (sb *Subscriber) pullMessageSync(ctx context.Context, process func(context.Context, []byte, string) error) {
	cm := make(chan *pubsub.Message)
	go func() {
		for {
			select {
			case msg := <-cm:
				err := process(sb.ctx, msg.Data, msg.ID)
				if err != nil {
					sb.errorHook(err)
					msg.Nack()
				} else {
					msg.Ack()
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	err := sb.subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Print(msg.ID)
		cm <- msg
	})
	if err != nil {
		sb.errorHook(err)
	}
}

func (sb *Subscriber) pullMessageAsync(ctx context.Context, process func(context.Context, []byte, string) error) {
	err := sb.subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		err := process(sb.ctx, msg.Data, msg.ID)
		if err != nil {
			sb.errorHook(err)
			msg.Nack()
		} else {
			msg.Ack()
		}
	})
	if err != nil {
		sb.errorHook(err)
	}
}
