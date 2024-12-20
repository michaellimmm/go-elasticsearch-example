package messaging

import (
	"context"
	"encoding/json"
	"github/shaolim/go-elasticsearch-example/internal/model"
	"github/shaolim/go-elasticsearch-example/internal/usecase"
	"log/slog"

	"cloud.google.com/go/pubsub"
)

type ItemUpsertConsumer struct {
	logger      *slog.Logger
	itemUseCase *usecase.ItemUpsertUseCase
}

func NewItemUpsertConsumer(logger *slog.Logger, itemUseCase *usecase.ItemUpsertUseCase) *ItemUpsertConsumer {
	return &ItemUpsertConsumer{
		logger:      logger,
		itemUseCase: itemUseCase,
	}
}

func (c *ItemUpsertConsumer) Consume(ctx context.Context, msg *pubsub.Message) {
	var items []*model.Item
	if err := json.Unmarshal(msg.Data, &items); err != nil {
		c.logger.Error("failed to unmarshal, err: %v\n", err)
		msg.Ack()
		return
	}

	err := c.itemUseCase.Execute(ctx, items)
	if err != nil {
		c.logger.Error("failed to execute item usecase", slog.Any("error", err))
		msg.Ack()
		return
	}

	msg.Ack()
}
