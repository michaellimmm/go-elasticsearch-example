package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type GCSNotifConsumer struct {
}

func NewGCSNotifConsumer() *GCSNotifConsumer {
	return &GCSNotifConsumer{}
}

func (c *GCSNotifConsumer) Consume(ctx context.Context, msg *pubsub.Message) {
	var m interface{}
	if err := json.Unmarshal(msg.Data, &m); err != nil {
		fmt.Printf("failed to unmarshal, err: %v\n", err)
		msg.Ack()
		return
	}

	fmt.Printf("message: %v\n", m)
	msg.Ack()
}
