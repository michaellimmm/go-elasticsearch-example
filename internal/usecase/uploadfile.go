package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github/shaolim/go-elasticsearch-example/pkg/esclient"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"gopkg.in/yaml.v3"
)

type uploadFileToGCSUseCase struct {
	esClient  esclient.Client
	gcsClient *storage.Client
	pbClient  *pubsub.Client
}

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

func NewUploadFileToGCSUseCase(
	esClient esclient.Client,
	gcsClient *storage.Client,
	pbClient *pubsub.Client) *uploadFileToGCSUseCase {
	return &uploadFileToGCSUseCase{
		esClient:  esClient,
		gcsClient: gcsClient,
		pbClient:  pbClient,
	}
}

func (u *uploadFileToGCSUseCase) Execute(ctx context.Context, bucketName, objectName, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return err
	}
	defer file.Close()

	bucketsConfFile, err := os.ReadFile("deploy/buckets/config.yaml")
	if err != nil {
		return fmt.Errorf("failed to read buckets config file: %v", err)
	}

	var bucketsConf BucketConfig
	if err := yaml.Unmarshal(bucketsConfFile, &bucketsConf); err != nil {
		return fmt.Errorf("failed to unmarshal buckets config file: %v", err)
	}

	bucketsMap := make(map[string]Buckets)
	for _, bc := range bucketsConf.Buckets {
		bucketsMap[bc.Name] = bc
	}

	// make sure topic is in the config
	bc, ok := bucketsMap[bucketName]
	if !ok {
		return fmt.Errorf("bucket %s not found in buckets config", bucketName)
	}

	wc := u.gcsClient.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	defer wc.Close()

	if _, err := io.Copy(wc, file); err != nil {
		fmt.Printf("failed to copy file to GCS, error: %v\n", err)
		return err
	}

	topic := u.pbClient.Topic(bc.Notifications.Topic)
	if err := u.sendNotificationToPubSub(ctx, topic, bucketName, bc.Notifications.EventType, objectName); err != nil {
		return fmt.Errorf("failed to send notification to pubsub: %v", err)
	}

	return nil
}

func (u *uploadFileToGCSUseCase) sendNotificationToPubSub(ctx context.Context,
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
