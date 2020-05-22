package mq

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"log"
	"sync"
	"time"
)

// topic controller
type Topic struct {
	ctx     context.Context
	initPub bool
	initSub bool
	topicID string
	wg      *sync.WaitGroup
}

type Option func(*Topic)

// A subID can only subscribe to a topic
func MakeSubID(topicID string) string {
	return topicID + "_" + "subscription"
}

func InitSub() Option {
	return func(t *Topic) {
		t.initSub = true
	}
}

func InitPub() Option {
	return func(t *Topic) {
		t.initPub = true
	}
}

func UseWaitGroup(wg *sync.WaitGroup) Option {
	return func(t *Topic) {
		t.wg = wg
	}
}

func (t *Topic) CreateSubscriptionIfNotExist(client *pubsub.Client, subID string) error {
	exist, err := t.IsSubscriptionExist(client, subID)
	if err != nil {
		return errors.WithStack(err)
	}
	if exist {
		return nil
	}
	_, err = client.CreateSubscription(t.ctx, subID, pubsub.SubscriptionConfig{
		Topic:       client.Topic(t.topicID),
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (t *Topic) IsSubscriptionExist(client *pubsub.Client, subID string) (bool, error) {
	it := client.Subscriptions(t.ctx)
	for {
		s, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return false, errors.WithStack(err)
		}

		if s.ID() == subID {
			log.Printf("Subscription exist: %s\n", subID)
			return true, nil
		}
	}
	return false, nil
}

func (t *Topic) IsTopicExist(client *pubsub.Client, topic string) (bool, error) {
	topicInstance := client.Topic(topic)
	return topicInstance.Exists(t.ctx)
}

func (t *Topic) CreateTopicIfNotExist(client *pubsub.Client, topic string) error {
	exist, err := t.IsTopicExist(client, topic)
	if err != nil {
		return errors.WithStack(err)
	}
	if exist {
		return nil
	}
	_, err = client.CreateTopic(t.ctx, topic)
	if err != nil {
		return errors.WithStack(err)
	}
	//gcp init topic 需要時間，不然直接pubsub會錯
	time.Sleep(time.Second * 10)
	return nil
}

func (t *Topic) DeleteTopicIfExist(client *pubsub.Client, topic string) error {
	exist, err := t.IsTopicExist(client, topic)
	if err != nil {
		return errors.WithStack(err)
	}
	if !exist {
		return nil
	}
	topicInstance := client.Topic(topic)
	if err := topicInstance.Delete(t.ctx); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
