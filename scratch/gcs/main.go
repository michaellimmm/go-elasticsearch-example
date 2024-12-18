package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	bucketName = "test-bucket"
	projectID  = "test-project"
	objectName = "test-object"
	topicID    = "your-topic-id"
	subName    = "your-subscription-name"
)

func main() {
	os.Setenv("STORAGE_EMULATOR_HOST", "localhost:4443")
	os.Setenv(`PUBSUB_EMULATOR_HOST`, `localhost:8085`)

	ctx := context.Background()

	pbClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		slog.Error("failed to create pubsub client", slog.Any("error", err))
		return
	}
	defer pbClient.Close()

	topic := pbClient.Topic(topicID)
	ok, err := topic.Exists(ctx)
	if err != nil {
		slog.Error("failed to check if topic exists", slog.Any("error", err))
		return
	}
	if !ok {
		if _, err := pbClient.CreateTopic(ctx, topicID); err != nil {
			slog.Error("failed to create topic", slog.Any("error", err))
			return
		}
		slog.Info("topic created")
	}

	sub := pbClient.Subscription(subName)
	ok, err = sub.Exists(ctx)
	if err != nil {
		slog.Error("failed to check if subscription exists", slog.Any("error", err))
		return
	}
	if !ok {
		if _, err := pbClient.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{
			Topic: topic,
		}); err != nil {
			slog.Error("failed to create subscription", slog.Any("error", err))
			return
		}
		slog.Info("subscription created")
	}

	gcsClient, err := storage.NewClient(ctx, option.WithEndpoint(os.Getenv("STORAGE_EMULATOR_HOST")))
	if err != nil {
		slog.Error("failed to create storage client", slog.Any("error", err))
		return
	}
	defer gcsClient.Close()

	bucket := gcsClient.Bucket(bucketName)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrBucketNotExist) {
			if err := bucket.Create(ctx, projectID, nil); err != nil {
				slog.Error("failed to create bucket", slog.Any("error", err))
				return
			}
		} else {
			slog.Error("failed to get bucket attributes", slog.Any("error", err))
		}
	}

	wc := bucket.Object(objectName).NewWriter(ctx)
	if _, err := wc.Write([]byte("test12345678")); err != nil {
		slog.Error("failed to write to object", slog.Any("error", err))
		return
	}
	if err := wc.Close(); err != nil {
		slog.Error("failed to close writer", slog.Any("error", err))
		return
	}

	rc, err := bucket.Object(objectName).NewReader(ctx)
	if err != nil {
		slog.Error("failed to create reader", slog.Any("error", err))
		return
	}
	defer rc.Close()

	if _, err := os.Stdout.ReadFrom(rc); err != nil {
		slog.Error("failed to read from reader", slog.Any("error", err))
		return
	}

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		slog.Info("Received message", slog.String("data", string(msg.Data)))
		msg.Ack()
	})
	if err != nil {
		slog.Error("failed to receive messages", slog.Any("error", err))
		return
	}
}
