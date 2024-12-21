package usecase

import (
	"context"
	"github/shaolim/go-elasticsearch-example/config"
	"github/shaolim/go-elasticsearch-example/internal/model"
	"github/shaolim/go-elasticsearch-example/pkg/esclient"
	"log/slog"
)

type ItemUpsertUseCase struct {
	logger   *slog.Logger
	esClient esclient.Client
}

func NewItemUpsertUseCase(logger *slog.Logger, esClient esclient.Client) *ItemUpsertUseCase {
	return &ItemUpsertUseCase{
		logger:   logger,
		esClient: esClient,
	}
}

func (u *ItemUpsertUseCase) Execute(ctx context.Context, items []*model.Item) error {
	req := u.convItemToBulkRequest(items)
	enBulkRequest := req["en"]
	if enBulkRequest.Length() > 0 {
		res, err := u.esClient.Bulk(config.ItemIndexEn, enBulkRequest)
		if err != nil {
			u.logger.Error("failed to bulk insert", slog.String("index", config.ItemIndexEn), slog.Any("error", err))
			return err
		}

		u.logger.Info("status code", slog.String("index", config.ItemIndexEn), slog.Int("status_code", res.StatusCode))
	}

	jaBulkRequest := req["ja"]
	if jaBulkRequest.Length() > 0 {
		res, err := u.esClient.Bulk(config.ItemIndexJa, jaBulkRequest)
		if err != nil {
			u.logger.Error("failed to bulk insert", slog.String("index", config.ItemIndexJa), slog.Any("error", err))
			return err
		}

		u.logger.Info("status code", slog.String("index", config.ItemIndexJa), slog.Int("status_code", res.StatusCode))

	}

	return nil
}

func (u *ItemUpsertUseCase) convItemToBulkRequest(items []*model.Item) map[string]*esclient.BulkRequests {
	enBulkRequest := &esclient.BulkRequests{}
	jaBulkRequest := &esclient.BulkRequests{}

	result := make(map[string]*esclient.BulkRequests)
	for _, item := range items {
		docs := model.ConvertItemToItemDoc(*item)

		if item.LanguageCode == "en" {
			if item.IsDeleted() {
				enBulkRequest.Add(esclient.NewBulkDeleteRequest(docs.Sku))
			} else {
				enBulkRequest.Add(esclient.NewBulkIndexRequest().SetId(docs.Sku).SetDoc(docs))
			}
		} else if item.LanguageCode == "ja" {
			if item.IsDeleted() {
				jaBulkRequest.Add(esclient.NewBulkDeleteRequest(docs.Sku))
			} else {
				jaBulkRequest.Add(esclient.NewBulkIndexRequest().SetId(docs.Sku).SetDoc(docs))
			}
		}
	}

	result["en"] = enBulkRequest
	result["ja"] = jaBulkRequest

	return result
}
