package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"
)

type BucketConfig struct {
	Buckets []Buckets `yaml:"buckets"`
}

type Buckets struct {
	Name          string `yaml:"name"`
	Notifications struct {
		Topic string `yaml:"topic"`
	} `yaml:"notifications"`
}

type PubsubConfig struct {
	PubsubTopics []PubsubTopic `yaml:"topics"`
}

type PubsubTopic struct {
	TopicID       string `yaml:"topic_id"`
	Subscriptions []struct {
		Name string `yaml:"name"`
	} `yaml:"subscriptions"`
}

func main() {
	// waiting for the emulator to be ready
	time.Sleep(10 * time.Second)

	fmt.Println("sync job is running")

	syncPubsub()
	syncBuckets()

	fmt.Println("deployment file already synced")
}

func syncPubsub() {
	fmt.Println("pubsub is syncing...")

	pbConfFile, err := os.ReadFile("deploy/pubsub/config.yaml")
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		panic(err)
	}
	var pbConf PubsubConfig
	if err := yaml.Unmarshal(pbConfFile, &pbConf); err != nil {
		fmt.Printf("failed to unmarshal pubsub config file: %v\n", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	projectID := os.Getenv("GCP_PROJECT_ID")
	pbClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		fmt.Printf("failed to create pubsub client: %v\n", err)
		panic(err)
	}
	defer pbClient.Close()

	currTopics := make(map[string]bool)
	topicIterator := pbClient.Topics(context.Background())
	for {
		topic, err := topicIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Printf("failed to iterate topics: %v\n", err)
			panic(err)
		}
		currTopics[topic.ID()] = false
	}

	for _, tp := range pbConf.PubsubTopics {
		if _, ok := currTopics[tp.TopicID]; ok {
			currTopics[tp.TopicID] = true
			continue
		}

		topic := pbClient.Topic(tp.TopicID)
		ok, err := topic.Exists(context.Background())
		if err != nil {
			fmt.Printf("failed to check if topic exists: %s, error: %v\n", tp.TopicID, err)
			panic(err)
		}
		if !ok {
			if _, err := pbClient.CreateTopic(context.Background(), tp.TopicID); err != nil {
				fmt.Printf("failed to create topic: %s, error: %v\n", tp.TopicID, err)
				panic(err)
			}

			fmt.Printf("created topic: %s\n", tp.TopicID)
		}

		if len(tp.Subscriptions) > 0 {
			for _, sb := range tp.Subscriptions {
				sub := pbClient.Subscription(sb.Name)
				ok, err := sub.Exists(context.Background())
				if err != nil {
					fmt.Printf("failed to check if subscription exists: %s, error: %v\n", sb.Name, err)
					panic(err)
				}
				if !ok {
					if _, err := pbClient.CreateSubscription(context.Background(), sb.Name, pubsub.SubscriptionConfig{
						Topic: topic,
					}); err != nil {
						fmt.Printf("failed to create subscription: %s, error: %v\n", sb.Name, err)
						panic(err)
					}

					fmt.Printf("created subscription: %s\n", sb.Name)
				}
			}
		}
	}

	for topicID, isUsed := range currTopics {
		if !isUsed {
			fmt.Printf("deleting unused topic: %s\n", topicID)
			topic := pbClient.Topic(topicID)
			err := topic.Delete(context.Background())
			if err != nil {
				fmt.Printf("failed to delete topic: %s, error: %v\n", topicID, err)
			}
		}
	}

	fmt.Println("pubsub synced")
}

func syncBuckets() {
	fmt.Println("buckets is syncing...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	bucketsConfFile, err := os.ReadFile("deploy/buckets/config.yaml")
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		panic(err)
	}
	var bucketsConf BucketConfig
	if err := yaml.Unmarshal(bucketsConfFile, &bucketsConf); err != nil {
		fmt.Printf("failed to unmarshal buckets config file: %v\n", err)
		panic(err)
	}

	projectID := os.Getenv("GCP_PROJECT_ID")
	storageEmulatorHost := os.Getenv("STORAGE_EMULATOR_HOST")
	gcsClient, err := storage.NewClient(ctx, option.WithEndpoint(storageEmulatorHost))
	if err != nil {
		fmt.Printf("failed to create gcs client: %v\n", err)
		panic(err)
	}
	defer gcsClient.Close()

	for _, bc := range bucketsConf.Buckets {
		bucket := gcsClient.Bucket(bc.Name)
		if _, err := bucket.Attrs(ctx); err != nil {
			if errors.Is(err, storage.ErrBucketNotExist) {
				if err := bucket.Create(ctx, projectID, nil); err != nil {
					fmt.Printf("failed to create bucket: %s, error: %v\n", bc.Name, err)
					panic(err)
				}

				fmt.Printf("created bucket: %s\n", bc.Name)
			}
		}
	}

	fmt.Println("buckets synced")
}
