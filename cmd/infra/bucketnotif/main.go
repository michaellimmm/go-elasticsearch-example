package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"
)

type BucketConfig struct {
	Buckets []Buckets `yaml:"buckets"`
}

type Buckets struct {
	Name          string `yaml:"name"`
	Notifications struct {
		Topic     string `yaml:"topic"`
		EventType string `yaml:"eventType"`
	} `yaml:"notifications"`
}

// since fake-gcs-server doesn't support gcs notifications,
// we just add the bucket to the watcher
func main() {
	time.Sleep(15 * time.Second)
	fmt.Println("bucket notification is running")

	bucketsConfFile, err := os.ReadFile("deploy/buckets/config.yaml")
	if err != nil {
		panic(err)
	}

	var bucketsConf BucketConfig
	if err := yaml.Unmarshal(bucketsConfFile, &bucketsConf); err != nil {
		panic(err)
	}

	projectID := os.Getenv("GCP_PROJECT_ID")
	pbClient, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		panic(err)
	}
	defer pbClient.Close()

	storageEmulatorHost := os.Getenv("STORAGE_EMULATOR_HOST")
	gcsClient, err := storage.NewClient(context.Background(), option.WithEndpoint(storageEmulatorHost))
	if err != nil {
		panic(err)
	}
	defer gcsClient.Close()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}
	defer watcher.Close()

	var eg errgroup.Group
	for _, bc := range bucketsConf.Buckets {
		topic := pbClient.Topic(bc.Notifications.Topic)

		eg.Go(func() error {
			if err := watcher.Add("/usr/share/fake-gcs-server/storage/" + bc.Name); err != nil {
				log.Fatalf("Error adding directory to watcher: %v", err)
			}

			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return fmt.Errorf("event is not valid")
					}

					filename := filepath.Base(event.Name)
					fmt.Printf("filename: %s\n", filename)
					if strings.HasPrefix(filename, ".") {
						continue
					}

					if event.Op&fsnotify.Create == fsnotify.Create && bc.Notifications.EventType == "OBJECT_FINALIZE" {
						fmt.Printf("File created, file: %s\n", event.Name)
						if err := sendNotificationToPubSub(context.Background(), topic, bc.Name, "OBJECT_FINALIZE", event.Name); err != nil {
							return err
						}
					}
					if event.Op&fsnotify.Write == fsnotify.Write && bc.Notifications.EventType == "OBJECT_UPDATE" {
						fmt.Printf("File modified, file: %s\n", event.Name)
						if err := sendNotificationToPubSub(context.Background(), topic, bc.Name, "OBJECT_UPDATE", event.Name); err != nil {
							log.Fatalf("Error sending notification: %v", err)
							return err
						}
					}
					if event.Op&fsnotify.Remove == fsnotify.Remove {
						fmt.Printf("File deleted, file: %s\n", event.Name)
						if err := sendNotificationToPubSub(context.Background(), topic, bc.Name, "OBJECT_DELETE", event.Name); err != nil {
							return err
						}
					}

				case err := <-watcher.Errors:
					return err
				}
			}
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Printf("error wait, err: %+v\n", err)
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
	topic *pubsub.Topic,
	bucketName string,
	eventType string,
	fileName string,
) error {

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
