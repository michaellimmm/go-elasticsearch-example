package main

import (
	"context"
	"github/shaolim/kakashi/internal/delivery/messaging"
	"github/shaolim/kakashi/internal/lib"
	"github/shaolim/kakashi/internal/usecase"
	"github/shaolim/kakashi/pkg/esclient"
	"log"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"golang.org/x/sync/errgroup"
)

func main() {
	vp := lib.NewViper()
	logger := lib.NewLogger()

	ctx := context.Background()

	pbClient, err := pubsub.NewClient(ctx, vp.GetString("GCP_PROJECT_ID"))
	if err != nil {
		log.Fatalf("failed to create pubsub client, err:%+v\n", err)
	}
	defer pbClient.Close()

	gcsClient, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create gcs client, err:%+v\n", err)
	}
	defer gcsClient.Close()

	esClient := esclient.NewClient("http://localhost:9200")

	// usecase
	ingestionUseCase := usecase.NewIngestionUseCase(vp, logger, gcsClient, getItemIngestionTopic(pbClient))
	itemUseCase := usecase.NewItemUpsertUseCase(logger, esClient)

	gcsNotifConsumer := messaging.NewGCSNotifConsumer(logger, ingestionUseCase)
	gcsNotifSubscriber := pbClient.Subscription("bucket-notification")

	itemUpsertConsumer := messaging.NewItemUpsertConsumer(logger, itemUseCase)
	itemUpsertSubscriber := pbClient.Subscription("items-upsert")

	eg := errgroup.Group{}
	eg.Go(func() error {
		return gcsNotifSubscriber.Receive(ctx, gcsNotifConsumer.Consume)
	})

	eg.Go(func() error {
		return itemUpsertSubscriber.Receive(ctx, itemUpsertConsumer.Consume)
	})

	if err := eg.Wait(); err != nil {
		log.Fatalf("failed to receive messages, err:%+v\n", err)
	}
}

func getItemIngestionTopic(client *pubsub.Client) *pubsub.Topic {
	return client.Topic("item-and-offer")
}
