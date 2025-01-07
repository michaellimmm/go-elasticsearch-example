package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"github/shaolim/kakashi/internal/model"
	"github/shaolim/kakashi/internal/usecase"
	"log/slog"

	"cloud.google.com/go/pubsub"
)

type GCSNotifConsumer struct {
	logger           *slog.Logger
	ingestionUsecase *usecase.IngestionUseCase
}

func NewGCSNotifConsumer(logger *slog.Logger, ingestionUsecase *usecase.IngestionUseCase) *GCSNotifConsumer {
	return &GCSNotifConsumer{
		logger:           logger,
		ingestionUsecase: ingestionUsecase,
	}
}

func (c *GCSNotifConsumer) Consume(ctx context.Context, msg *pubsub.Message) {
	var attr model.GscAttribute
	if err := json.Unmarshal(msg.Data, &attr); err != nil {
		c.logger.Error("failed to unmarshal, err: %v\n", err)
		// TODO: add to metrics
		msg.Ack()
		return
	}

	fmt.Printf("message: %v\n", attr)

	err := c.ingestionUsecase.Execute(ctx, attr.Bucket, attr.Name)
	if err != nil {
		msg.Ack()
		return
	}

	msg.Ack()
}
