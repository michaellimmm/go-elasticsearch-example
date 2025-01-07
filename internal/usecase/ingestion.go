package usecase

import (
	"context"
	"encoding/json"
	"github/shaolim/kakashi/internal/lib"
	"github/shaolim/kakashi/internal/model"
	"log/slog"
	"sync"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/gocarina/gocsv"
	"github.com/spf13/viper"
)

type IngestionUseCase struct {
	viper     *viper.Viper
	logger    *slog.Logger
	gcsClient *storage.Client
	publisher lib.Publisher
}

func NewIngestionUseCase(
	viper *viper.Viper,
	logger *slog.Logger,
	gcsClient *storage.Client,
	publisher lib.Publisher,
) *IngestionUseCase {
	return &IngestionUseCase{
		viper:     viper,
		logger:    logger,
		gcsClient: gcsClient,
		publisher: publisher,
	}
}

func (u *IngestionUseCase) Execute(ctx context.Context, bucketname string, filename string) error {
	rc, err := u.gcsClient.Bucket(bucketname).Object(filename).NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()

	queue := make(chan *model.Item, u.viper.GetInt("PARSER_QUEUE_SIZE"))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		u.processItem(ctx, queue)
	}()

	go func() {
		if err := gocsv.UnmarshalToChan(rc, queue); err != nil {
			u.logger.Error("failed to UnmarshalToChan", slog.Any("error", err))
			return
		}
	}()

	wg.Wait()

	return nil
}

func (u *IngestionUseCase) processItem(ctx context.Context, in <-chan *model.Item) {
	batchSize := u.viper.GetInt("PARSER_BATCH_SIZE")
	batches := make([]*model.Item, 0, batchSize)
	for item := range in {
		batches = append(batches, item)
		if len(batches) >= batchSize {
			msgData, err := json.Marshal(batches)
			if err != nil {
				u.logger.Error("failed to marshal batches", slog.Any("error", err))
				continue
			}

			pbMsg := &pubsub.Message{
				Data: msgData}
			u.publisher.Publish(ctx, pbMsg)
			u.logger.Info("published", slog.Int("batch_size", len(batches)))

			batches = make([]*model.Item, 0, viper.GetInt("PARSER_BATCH_SIZE"))
		}
	}

	if len(batches) > 0 {
		msgData, err := json.Marshal(batches)
		if err != nil {
			u.logger.Error("failed to marshal batches", slog.Any("error", err))
		}

		pbMsg := &pubsub.Message{
			Data: msgData}
		u.publisher.Publish(ctx, pbMsg)
		u.logger.Info("published", slog.Int("batch_size", len(batches)))
	}
}
