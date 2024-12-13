package main

import (
	"github/shaolim/go-elasticsearch-example/pkg/esclient"
	"github/shaolim/go-elasticsearch-example/pkg/esclient/esquery"
	"github/shaolim/go-elasticsearch-example/utils/middleware"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

var (
	logger *slog.Logger
)

type product struct {
	Name    string  `json:"name,omitempty"`
	Price   float32 `json:"price,omitempty"`
	InStock *bool   `json:"in_stock,omitempty"`
}

func main() {
	logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

	defer func() {
		if rec := recover(); rec != nil {
			logger.Error("panic", slog.Any("error", rec),
				slog.Any("stacktrace", strings.Split(string(debug.Stack()), "\n")))
		}
	}()

	httpClient := &http.Client{
		Transport: middleware.ChainMiddleware(http.DefaultTransport,
			middleware.NewLoggingMiddleware(logger),
		),
	}

	client := esclient.NewClient("http://localhost:9200", esclient.WithHttpClient(httpClient))
	if err := ping(client); err != nil {
		return
	}

	if err := createIndexIfNotExists(client); err != nil {
		return
	}

	// bulkRequest := esclient.BulkRequests{}
	// bulkRequest.
	// 	Add(esclient.NewBulkIndexRequest().SetId("1").
	// 		SetDoc(product{Name: "Laptop Pro 15", Price: 1200, InStock: boolPtr(true)})).
	// 	Add(esclient.NewBulkIndexRequest().SetId("2").
	// 		SetDoc(product{Name: "Smartphone X", Price: 999, InStock: boolPtr(false)})).
	// 	Add(esclient.NewBulkIndexRequest().SetId("3").
	// 		SetDoc(product{Name: "Noise-Canceling Headphones", Price: 199, InStock: boolPtr(true)})).
	// 	Add(esclient.NewBulkIndexRequest().SetId("4").
	// 		SetDoc(product{Name: "Mechanical Keyboard", Price: 85, InStock: boolPtr(true)})).
	// 	Add(esclient.NewBulkIndexRequest().SetId("5").
	// 		SetDoc(product{Name: "Gaming Mouse", Price: 45, InStock: boolPtr(true)})).
	// 	Add(esclient.NewBulkIndexRequest().SetId("6").
	// 		SetDoc(product{Name: "4K Monitor", Price: 450, InStock: boolPtr(true)})).
	// 	Add(esclient.NewBulkIndexRequest().SetId("7").
	// 		SetDoc(product{Name: "Bluetooth Speaker", Price: 75, InStock: boolPtr(true)})).
	// 	Add(esclient.NewBulkIndexRequest().SetId("8").
	// 		SetDoc(product{Name: "External SSD 1TB", Price: 150, InStock: boolPtr(true)})).
	// 	Add(esclient.NewBulkIndexRequest().SetId("9").
	// 		SetDoc(product{Name: "Fitness Tracker", Price: 60, InStock: boolPtr(true)})).
	// 	Add(esclient.NewBulkIndexRequest().SetId("10").
	// 		SetDoc(product{Name: "Smart Home Hub", Price: 99, InStock: boolPtr(true)}))

	// bulkResponse, err := client.Bulk("products", bulkRequest)
	// if err != nil {
	// 	logger.Error("bulk", slog.Any("error", err))
	// 	return
	// }

	// logger.Info("bulk", slog.Any("result", bulkResponse.Result))

	searchRequest := esquery.NewSearchQueryBuilder().
		SetSize(10).
		SetFrom(0).
		SetQuery(esquery.MatchAll()).
		Build()

	searchResponse, err := client.Search("products", *searchRequest)
	if err != nil {
		logger.Error("search", slog.Any("error", err))
		return
	}

	logger.Info("search", slog.Any("result", searchResponse.Result))

	countRes, err := client.Count("products", esquery.MatchAll())
	if err != nil {
		logger.Error("search", slog.Any("error", err))
		return
	}

	logger.Info("count", slog.Any("result", countRes.Result))
}

func ping(client esclient.Client) error {
	res, err := client.Ping(esclient.PingWithHttpHeadOnly())
	if err != nil {
		logger.Error("ping", slog.Any("error", err))
		return err
	}

	if res.Result != nil {
		logger.Info("ping", slog.Any("result", res.Result))
	}

	return nil
}

func createIndexIfNotExists(client esclient.Client) error {
	indexName := "products"
	payload := `{
		"settings":{
			"index.number_of_shards":1,
			"index.number_of_replicas":1
		},
		"mappings":{
			"properties":{
				"name":{
					"type":"text"
				},
				"price":{
					"type":"float"
				},
				"in_stock":{
					"type":"boolean"
				}
			}
		}
	}`

	indexResponse, err := client.GetIndeces([]string{indexName}, esclient.GetIndecesWithHttpHeadOnly())
	if err != nil {
		if indexResponse.StatusCode != 404 {
			logger.Error("error", slog.Any("error", err))
			return err
		}
		logger.Error("error", slog.Any("error", err))
	}
	logger.Info("index exists", slog.Any("body", indexResponse.Result))

	if indexResponse.StatusCode == 404 {
		res, err := client.CreateIndex(indexName, strings.NewReader(payload))
		if err != nil {
			logger.Error("error", slog.Any("error", err))
			return err
		}
		logger.Info("index created", slog.Any("result", res))
	}

	return nil
}

func boolPtr(b bool) *bool {
	return &b
}
