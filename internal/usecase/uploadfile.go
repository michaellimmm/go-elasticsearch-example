package usecase

import (
	"context"
	"fmt"
	"github/shaolim/go-elasticsearch-example/pkg/esclient"
	"github/shaolim/go-elasticsearch-example/pkg/gcs"
	"io"
	"os"
)

type uploadFileToGCSUseCase struct {
	esClient  esclient.Client
	gcsClient *gcs.Client
}

func NewUploadFileToGCSUseCase(esClient esclient.Client, gcsClient *gcs.Client) *uploadFileToGCSUseCase {
	return &uploadFileToGCSUseCase{
		esClient:  esClient,
		gcsClient: gcsClient,
	}
}

func (u *uploadFileToGCSUseCase) Execute(ctx context.Context, bucketName, objectName, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return err
	}
	defer file.Close()

	wc := u.gcsClient.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	defer wc.Close()

	if _, err := io.Copy(wc, file); err != nil {
		fmt.Printf("failed to copy file to GCS, error: %v\n", err)
		return err
	}

	return nil
}
