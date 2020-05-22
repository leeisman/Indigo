package mq

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/leeisman/Indigo/inside/pub"
	"github.com/leeisman/Indigo/inside/sub"
	"log"
)

// pub/sub message
type MQ struct {
	subscriber *sub.Subscriber
	publisher  *pub.Publisher
	ctx        context.Context
	cancel     context.CancelFunc
	client     *pubsub.Client
	logger     Logger
	Topic      *Topic
}

func NewMq(ctx context.Context) (*MQ, error) {
	client, err := NewClient(ctx)
	ctx, cacnel := context.WithCancel(ctx)
	if err != nil {
		return nil, err
	}
	return &MQ{
		ctx:    ctx,
		client: client,
		logger: &defaultLogger{},
		cancel: cacnel,
	}, nil
}

type Logger interface {
	InfoMsg(ctx context.Context, msg string, v ...interface{})
	ErrMsg(ctx context.Context, msg string, v ...interface{})
}

type defaultLogger struct{}

func (d *defaultLogger) InfoMsg(ctx context.Context, msg string, v ...interface{}) {
	log.Printf(msg, v)
}

func (d *defaultLogger) ErrMsg(ctx context.Context, msg string, v ...interface{}) {
	log.Printf(msg, v)
}

func (mq *MQ) SetLogger(logger Logger) {
	mq.logger = logger
}

func (mq *MQ) Publisher() *pub.Publisher {
	if mq.publisher == nil {
		mq.publisher = pub.NewPublish(mq.ctx, mq.client.Topic(mq.Topic.topicID))
	}
	return mq.publisher
}

func (mq *MQ) Subscriber() *sub.Subscriber {
	if mq.subscriber == nil {
		subscription := mq.client.Subscription(MakeSubID(mq.Topic.topicID))
		mq.subscriber = sub.NewSubscriber(mq.ctx, subscription)
	}
	return mq.subscriber
}

func (mq *MQ) Stop() {
	mq.cancel()
}
