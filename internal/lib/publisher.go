package lib

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type Publisher interface {
	Publish(ctx context.Context, ev *pubsub.Message) *pubsub.PublishResult
}
