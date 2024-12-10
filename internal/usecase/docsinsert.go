package usecase

import (
	"fmt"
	"github/shaolim/go-elasticsearch-example/internal/model"
	"github/shaolim/go-elasticsearch-example/pkg/esclient"
	"os"
	"sync"

	"github.com/gocarina/gocsv"
)

type docsInsertUseCase struct {
	esClient esclient.Client
}

func NewDocsInsertUseCase(esClient esclient.Client) *docsInsertUseCase {
	return &docsInsertUseCase{
		esClient: esClient,
	}
}

func (u *docsInsertUseCase) Execute(indexname string, filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	queue := make(chan *model.Item, 1000)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		u.processItem(indexname, queue)
	}()

	go func() {
		if err := gocsv.UnmarshalToChan(file, queue); err != nil {
			fmt.Printf("failed to UnmarshalToChan: %+v\n", err)
			return
		}
	}()

	wg.Wait()

	return nil
}

func (u *docsInsertUseCase) processItem(indexname string, in <-chan *model.Item) {
	batches := make([]*model.Item, 0, 100)
	for item := range in {
		batches = append(batches, item)
		if len(batches) >= 100 {
			req := u.convItemToBulkRequest(batches)
			res, err := u.esClient.Bulk(indexname, req)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(res.StatusCode)

			if res.IsError() {
				fmt.Printf("failed to bulk: %+v\n", res.Result)
			}

			batches = make([]*model.Item, 0, 100)
		}
	}

	if len(batches) > 0 {
		req := u.convItemToBulkRequest(batches)
		res, err := u.esClient.Bulk(indexname, req)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(res.StatusCode)

		if res.IsError() {
			fmt.Printf("failed to bulk: %+v\n", res.Result)
		}
	}
}

func (u *docsInsertUseCase) convItemToBulkRequest(items []*model.Item) *esclient.BulkRequests {
	bulkRequest := &esclient.BulkRequests{}
	for _, item := range items {
		docs := model.ConvertItemToItemDoc(*item)
		if item.IsDeleted() {
			bulkRequest.Add(esclient.NewBulkDeleteRequest(docs.Sku))
		} else {
			bulkRequest.Add(esclient.NewBulkIndexRequest().SetId(docs.Sku).SetDoc(docs))
		}
	}
	return bulkRequest
}
