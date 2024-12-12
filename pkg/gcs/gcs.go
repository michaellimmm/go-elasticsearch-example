package gcs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/fsnotify/fsnotify"
	"google.golang.org/api/option"
)

type Client struct {
	*storage.Client
	watcher   *fsnotify.Watcher
	pbClient  *pubsub.Client
	publisher eventPublisher
}

type eventPublisher interface {
	Publish(ctx context.Context, msg *pubsub.Message) *pubsub.PublishResult
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

func NewClient(ctx context.Context, projectID string, topicID string, opts ...option.ClientOption) (*Client, error) {
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	pbClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	topic := pbClient.Topic(topicID)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		if topic, err = pbClient.CreateTopic(ctx, topicID); err != nil {
			return nil, err
		}
	}

	return &Client{
		Client:    client,
		watcher:   watcher,
		pbClient:  pbClient,
		publisher: topic,
	}, nil
}

// Watch starts watching the given bucket.
// Since fake-gcs-server doesn't support gcs notifications,
// we just add the bucket to the watcher.
func (c *Client) Watch(ctx context.Context, bucketName string) error {
	if err := c.watcher.Add("./data/" + bucketName); err != nil {
		return err
	}

	for {
		select {
		case event, ok := <-c.watcher.Events:
			if !ok {
				return fmt.Errorf("watcher closed")
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				if err := c.sendNotification(ctx, bucketName, "OBJECT_FINALIZE", event.Name); err != nil {
					return err
				}
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if err := c.sendNotification(ctx, bucketName, "OBJECT_UPDATE", event.Name); err != nil {
					return err
				}
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				if err := c.sendNotification(ctx, bucketName, "OBJECT_DELETE", event.Name); err != nil {
					return err
				}
			}
		case err := <-c.watcher.Errors:
			return err
		}
	}
}

func (c *Client) sendNotification(ctx context.Context,
	bucketName string, eventType string, fileName string) error {
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
	res := c.publisher.Publish(ctx, &pubSubMsg)

	_, err = res.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to get publish result: %v", err)
	}

	return nil
}

func (c *Client) Close() error {
	err1 := c.Client.Close()
	err2 := c.watcher.Close()
	err3 := c.pbClient.Close()
	return errors.Join(err1, err2, err3)
}
