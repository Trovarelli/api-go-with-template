package repository

import (
	"context"
	"loja-produtos/src/internal/models"
)

type ProdutosRepository interface {
	GetAll(ctx context.Context) ([]models.Produto, error)
	GetByID(ctx context.Context, id int64) (*models.Produto, error)
	Create(ctx context.Context, p *models.Produto) (int64, error)
	Update(ctx context.Context, p *models.Produto) error
	Delete(ctx context.Context, id int64) error
}
