package pub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
)

type Publisher struct {
	ctx       context.Context
	topic     *pubsub.Topic
	errorHook func(err error, requestID string)
}

func NewPublish(ctx context.Context, topic *pubsub.Topic) *Publisher {
	return &Publisher{
		ctx:   ctx,
		topic: topic,
	}
}

func (p *Publisher) Publish(msg []byte, requestID string) error {
	result := p.topic.Publish(p.ctx, &pubsub.Message{
		Data: msg,
	})
	id, err := result.Get(p.ctx)
	if err != nil {
		p.errorHook(errors.WithStack(err), requestID)
		return err
	}
	log.Printf("Published a message: %v ; msg ID: %v ; requestID: %v\n", string(msg), id, requestID)
	return nil
}

func (p *Publisher) PublishJSON(msg interface{}, requestID string) error {
	b, err := json.Marshal(msg)
	if err != nil {
		p.errorHook(errors.WithStack(err), requestID)
		return errors.WithStack(err)
	}
	p.Publish(b, requestID)
	return nil
}

func (sb *Publisher) Options(opts ...PublisherOption) *Publisher {
	for _, opt := range opts {
		opt(sb)
	}
	return sb
}
