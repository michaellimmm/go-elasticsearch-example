package main

import (
	"context"
	"github/shaolim/go-elasticsearch-example/internal/delivery/messaging"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/spf13/viper"
)

func main() {
	vp := newViper()

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, vp.GetString("GCP_PROJECT_ID"))
	if err != nil {
		log.Fatalf("failed to create pubsub client, err:%+v\n", err)
	}

	gcsNotifConsumer := messaging.NewGCSNotifConsumer()
	subscriber := client.Subscription("bucket-notification")
	err = subscriber.Receive(ctx, gcsNotifConsumer.Consume)
	if err != nil {
		log.Fatalf("failed to receive messages, err:%+v\n", err)
	}
}

func newViper() *viper.Viper {
	vp := viper.New()
	vp.SetConfigFile("./.env")
	if err := vp.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config file: %v", err)
	}

	os.Setenv("STORAGE_EMULATOR_HOST", vp.GetString("STORAGE_EMULATOR_HOST"))
	os.Setenv(`PUBSUB_EMULATOR_HOST`, vp.GetString(`PUBSUB_EMULATOR_HOST`))
	vp.AutomaticEnv()

	return vp
}
