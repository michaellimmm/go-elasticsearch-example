package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/spf13/viper"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"
)

type BucketConfig struct {
	Buckets []Buckets `yaml:"buckets"`
}

type Buckets struct {
	Name          string `yaml:"name"`
	Notifications []struct {
		Topic        string `yaml:"topic"`
		Subscription string `yaml:"subscription"`
	} `yaml:"notifications"`
}

type PubsubConfig struct {
	PubsubTopics []PubsubTopic `yaml:"topics"`
}

type PubsubTopic struct {
	TopicID        string `yaml:"topic_id"`
	Subscriptioins []struct {
		Name string `yaml:"name"`
	} `yaml:"subscriptions"`
}

func main() {
	vp := viper.New()
	vp.SetConfigFile("./.env")
	if err := vp.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config file: %v", err)
	}
	os.Setenv("STORAGE_EMULATOR_HOST", vp.GetString("STORAGE_EMULATOR_HOST"))
	os.Setenv(`PUBSUB_EMULATOR_HOST`, vp.GetString(`PUBSUB_EMULATOR_HOST`))
	vp.AutomaticEnv() // Enable automatic environment variable override

	syncPubsub(vp)
	syncBuckets(vp)
}

func syncPubsub(vp *viper.Viper) {
	pbConfFile, err := os.ReadFile("deploy/pubsub/config.yaml")
	if err != nil {
		panic(err)
	}
	var pbConf PubsubConfig
	if err := yaml.Unmarshal(pbConfFile, &pbConf); err != nil {
		panic(err)
	}

	projectID := vp.GetString("GCP_PROJECT_ID")
	pbClient, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
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
			panic(err)
		}
		if !ok {
			if _, err := pbClient.CreateTopic(context.Background(), tp.TopicID); err != nil {
				panic(err)
			}
		}

		if len(tp.Subscriptioins) > 0 {
			for _, sb := range tp.Subscriptioins {
				sub := pbClient.Subscription(sb.Name)
				ok, err := sub.Exists(context.Background())
				if err != nil {
					panic(err)
				}
				if !ok {
					if _, err := pbClient.CreateSubscription(context.Background(), sb.Name, pubsub.SubscriptionConfig{
						Topic: topic,
					}); err != nil {
						panic(err)
					}
				}
			}
		}
	}

	for topicID, isUsed := range currTopics {
		if !isUsed {
			topic := pbClient.Topic(topicID)
			err := topic.Delete(context.Background())
			if err != nil {
				fmt.Printf("failed to delete topic: %s, error: %v\n", topicID, err)
			}
		}
	}
}

func syncBuckets(vp *viper.Viper) {
	bucketsConfFile, err := os.ReadFile("deploy/buckets/config.yaml")
	if err != nil {
		panic(err)
	}
	var bucketsConf BucketConfig
	if err := yaml.Unmarshal(bucketsConfFile, &bucketsConf); err != nil {
		panic(err)
	}

	projectID := vp.GetString("GCP_PROJECT_ID")
	storageEmulatorHost := vp.GetString("STORAGE_EMULATOR_HOST")
	gcsClient, err := storage.NewClient(context.Background(), option.WithEndpoint(storageEmulatorHost))
	if err != nil {
		panic(err)
	}
	defer gcsClient.Close()

	pbClient, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		panic(err)
	}
	defer pbClient.Close()

	for _, bc := range bucketsConf.Buckets {
		bucket := gcsClient.Bucket(bc.Name)
		if err := bucket.Create(context.Background(), projectID, nil); err != nil {
			panic(err)
		}

		if len(bc.Notifications) > 0 {
			for _, nb := range bc.Notifications {
				sub := pbClient.Subscription(nb.Subscription)
				ok, err := sub.Exists(context.Background())
				if err != nil {
					panic(err)
				}
				if !ok {
					if _, err := pbClient.CreateSubscription(context.Background(), nb.Subscription, pubsub.SubscriptionConfig{
						Topic: pbClient.Topic(nb.Topic),
					}); err != nil {
						panic(err)
					}
				}
			}
		}

	}
}
