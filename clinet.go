package mq

import (
	"cloud.google.com/go/pubsub"
	"context"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context) (*pubsub.Client, error) {
	return pubsub.NewClient(ctx, ProjectID, option.WithCredentialsFile(CredentialsFile))
}
