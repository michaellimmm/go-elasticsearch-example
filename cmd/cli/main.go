package main

import (
	"context"
	"flag"
	"fmt"
	config "github/shaolim/go-elasticsearch-example/config"
	"github/shaolim/go-elasticsearch-example/internal/usecase"
	"github/shaolim/go-elasticsearch-example/pkg/esclient"
	"os"
	"path/filepath"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/spf13/viper"
)

type Command string

const (
	CreateIndex     Command = "create-index"
	Indexing        Command = "indexing"
	MatchDocs       Command = "match-docs"
	UploadfileToGCS Command = "upload-file-to-gcs"
)

func main() {
	viper.SetConfigFile("./.env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	os.Setenv("STORAGE_EMULATOR_HOST", viper.GetString("STORAGE_EMULATOR_HOST"))
	os.Setenv(`PUBSUB_EMULATOR_HOST`, viper.GetString(`PUBSUB_EMULATOR_HOST`))
	os.Setenv("GCP_PROJECT_ID", viper.GetString("GCP_PROJECT_ID"))

	command := flag.String("command", "", "Command eg. create-index, indexing, match-docs, upload-file-to-gcs")
	filename := flag.String("file", "", "path of csv file")
	languageCode := flag.String("lang", "ja", "Language code")
	bucketName := flag.String("bucket", "test-bucket", "Bucket name")

	flag.Parse()

	if *command == "" {
		flag.PrintDefaults()
		return
	}

	switch Command(*command) {
	case CreateIndex:
		if err := createIndex(); err != nil {
			fmt.Println(err)
		}
	case Indexing:
		if *filename == "" {
			fmt.Println("filename is required to run this indexing command")
			return
		}
		if err := indexing(*languageCode, *filename); err != nil {
			fmt.Println(err)
		}
	case MatchDocs:
		if *filename == "" {
			fmt.Println("filename is required to run this match-docs command")
			return
		}

		if err := matchDocs(*filename, *languageCode); err != nil {
			fmt.Println(err)
		}
	case UploadfileToGCS:
		if *filename == "" {
			fmt.Println("filename is required to run this upload-file-to-gcs command")
			return
		}

		if err := uploadFileToGCS(*bucketName, *filename); err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Printf("unknown command: %s, valid commands: create-index, indexing, match-docs\n", *command)
	}
}

func createIndex() error {
	client := esclient.NewClient("http://localhost:9200")

	createIndexUC := usecase.NewCreateIndexUseCase(client)
	if err := createIndexUC.Execute(); err != nil {
		return fmt.Errorf("failed to create index, error: %v", err)
	}

	return nil
}

func indexing(languageCode string, filename string) error {
	client := esclient.NewClient("http://localhost:9200")

	index := config.ItemIndexJa
	if languageCode == "en" {
		index = config.ItemIndexEn
	}

	indexingUC := usecase.NewDocsInsertUseCase(client)
	if err := indexingUC.Execute(index, filename); err != nil {
		fmt.Printf("failed to indexing, error: %v\n", err)
		return err
	}
	return nil
}

func matchDocs(filename string, languageCode string) error {
	client := esclient.NewClient("http://localhost:9200")
	index := config.ItemIndexJa
	if languageCode == "en" {
		index = config.ItemIndexEn
	}

	matchDocs := usecase.NewSampleDocsUseCase(client)
	if err := matchDocs.Execute(index, filename); err != nil {
		fmt.Printf("failed to match docs, error: %v\n", err)
		return err
	}
	return nil
}

func uploadFileToGCS(bucketName, filename string) error {
	client := esclient.NewClient("http://localhost:9200")
	gcsClient, err := storage.NewClient(context.Background())
	if err != nil {
		return err
	}

	pbClient, err := pubsub.NewClient(context.Background(), os.Getenv("GCP_PROJECT_ID"))
	if err != nil {
		return err
	}

	objectName := filepath.Base(filename)

	uploadFileToGCSUC := usecase.NewUploadFileToGCSUseCase(client, gcsClient, pbClient)
	if err := uploadFileToGCSUC.Execute(context.Background(), bucketName, objectName, filename); err != nil {
		fmt.Printf("failed to upload file to GCS, error: %v\n", err)
		return err
	}
	return nil
}
