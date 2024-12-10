package usecase

import (
	"bytes"
	_ "embed"
	"fmt"
	index "github/shaolim/go-elasticsearch-example/internal/config/index"
	"github/shaolim/go-elasticsearch-example/pkg/esclient"
	"io"
)

const (
	itemIndexEn = "item_index_en"
	itemIndexJa = "item_index_ja"
	fileSuffix  = ".json"
)

type createIndexUseCase struct {
	esClient esclient.Client
}

func NewCreateIndexUseCase(esClient esclient.Client) *createIndexUseCase {
	return &createIndexUseCase{
		esClient: esClient,
	}
}

func (c *createIndexUseCase) Execute() error {
	itemIndexEnSettings, err := c.loadJsonFile(itemIndexEn)
	if err != nil {
		fmt.Printf("failed to load json file: %s, error: %v\n", itemIndexEn, err)
		return err
	}

	err = c.createIndexIfNotExists(itemIndexEn, bytes.NewReader(itemIndexEnSettings))
	if err != nil {
		fmt.Printf("failed to create index: %s, error: %v\n", itemIndexEn, err)
		return err
	}

	itemIndexJaSettings, err := c.loadJsonFile(itemIndexJa)
	if err != nil {
		fmt.Printf("failed to load json file: %s, error: %+v\n", itemIndexJa, err)
		return err
	}

	err = c.createIndexIfNotExists(itemIndexJa, bytes.NewReader(itemIndexJaSettings))
	if err != nil {
		fmt.Printf("failed to crate index: %s, error: %+v\n", itemIndexJa, err)
		return err
	}

	return nil
}

func (c *createIndexUseCase) createIndexIfNotExists(indexName string, body io.Reader) error {
	indexRes, err := c.esClient.GetIndeces([]string{indexName}, esclient.GetIndecesWithHttpHeadOnly())
	if err != nil {
		return err
	}

	if indexRes.StatusCode == 404 {
		_, err = c.esClient.CreateIndex(indexName, body)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *createIndexUseCase) loadJsonFile(filename string) ([]byte, error) {
	data, err := index.ConfigFiles.ReadFile(filename + fileSuffix)
	if err != nil {
		return nil, err
	}

	return data, nil
}
