package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/fsnotify/fsnotify"
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

	go func() {
		watchDirectory(topic, "./data/test-bucket")
	}()

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

type GCSNotification struct {
	Kind                    string `json:"kind"`
	ID                      string `json:"id"`
	SelfLink                string `json:"selfLink"`
	Name                    string `json:"name"`
	Bucket                  string `json:"bucket"`
	Generation              string `json:"generation"`
	Metageneration          string `json:"metageneration"`
	ContentType             string `json:"contentType"`
	TimeCreated             string `json:"timeCreated"`
	Updated                 string `json:"updated"`
	StorageClass            string `json:"storageClass"`
	Size                    string `json:"size"`
	TimeStorageClassUpdated string `json:"timeStorageClassUpdated"`
	EventType               string `json:"eventType"`
	NotificationMetadata    struct {
		EventType string `json:"eventType"`
		EventTime string `json:"eventTime"`
	} `json:"notificationMetadata"`
}

func sendNotificationToPubSub(ctx context.Context,
	topic *pubsub.Topic, eventType string, fileName string) error {

	// Create the GCS notification
	notification := GCSNotification{
		Kind:                    "storage#object",
		ID:                      fmt.Sprintf("%s/%s/%d", bucketName, fileName, time.Now().Unix()),
		SelfLink:                fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s/o/%s", bucketName, fileName),
		Name:                    fileName,
		Bucket:                  bucketName,
		Generation:              fmt.Sprintf("%d", time.Now().Unix()),
		Metageneration:          "1",
		ContentType:             "application/octet-stream", // Set the content type accordingly
		TimeCreated:             time.Now().Format(time.RFC3339),
		Updated:                 time.Now().Format(time.RFC3339),
		StorageClass:            "STANDARD", // Modify as needed
		Size:                    "1024",     // Adjust the size accordingly
		TimeStorageClassUpdated: time.Now().Format(time.RFC3339),
		EventType:               eventType,
		NotificationMetadata: struct {
			EventType string `json:"eventType"`
			EventTime string `json:"eventTime"`
		}{
			EventType: eventType,
			EventTime: time.Now().Format(time.RFC3339),
		},
	}

	// Convert the notification to JSON
	msgData, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %v", err)
	}

	// Create a Pub/Sub message
	pubSubMsg := pubsub.Message{
		Data: msgData,
	}

	// Publish the message to Pub/Sub
	res := topic.Publish(ctx, &pubSubMsg)

	_, err = res.Get(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get publish result: %v", err)
	}

	slog.Info("Notification sent successfully", slog.Any("notification", notification))
	return nil
}

func watchDirectory(topic *pubsub.Topic, watchDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}
	defer watcher.Close()

	// Add the directory to the watcher
	if err := watcher.Add(watchDir); err != nil {
		log.Fatalf("Error adding directory to watcher: %v", err)

	}

	slog.Info("Watching directory", slog.String("directory", watchDir))

	// Monitor events
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Filter and handle the events
			if event.Op&fsnotify.Create == fsnotify.Create {
				slog.Info("File Created", slog.String("file", event.Name))
				if err := sendNotificationToPubSub(context.Background(), topic, "OBJECT_FINALIZE", event.Name); err != nil {
					log.Fatalf("Error sending notification: %v", err)
				}
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				slog.Info("File Modified", slog.String("file", event.Name))
				if err := sendNotificationToPubSub(context.Background(), topic, "OBJECT_UPDATE", event.Name); err != nil {
					log.Fatalf("Error sending notification: %v", err)
				}
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				slog.Error("File Deleted", slog.String("file", event.Name))
				if err := sendNotificationToPubSub(context.Background(), topic, "OBJECT_DELETE", event.Name); err != nil {
					log.Fatalf("Error sending notification: %v", err)
				}
			}

		case err := <-watcher.Errors:
			log.Fatalf("Watcher error: %v", err)
		}
	}
}
