package usecase

import (
	"context"
	"github/shaolim/go-elasticsearch-example/internal/model"
)

type ItemUpsertUseCase struct {
}

func NewItemUpsertUseCase() *ItemUpsertUseCase {
	return &ItemUpsertUseCase{}
}

func (u *ItemUpsertUseCase) Execute(ctx context.Context, items []*model.Item) error {

	return nil
}
