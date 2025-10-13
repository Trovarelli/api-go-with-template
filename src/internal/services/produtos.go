package service

import (
	"context"
	"errors"

	"loja-produtos/src/internal/models"
	"loja-produtos/src/internal/repository"
)

var ErrProdutoNaoEncontrado = errors.New("produto não encontrado")

type ProdutosService struct {
	repo repository.ProdutosRepository
}

func NewProdutosService(repo repository.ProdutosRepository) *ProdutosService {
	return &ProdutosService{repo: repo}
}

func (s *ProdutosService) Listar(ctx context.Context) ([]models.Produto, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProdutosService) BuscarPorID(ctx context.Context, id int64) (*models.Produto, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrProdutoNaoEncontrado
	}
	return p, nil
}

func (s *ProdutosService) Criar(ctx context.Context, p *models.Produto) (int64, error) {
	if p.Nome == "" || p.Preco <= 0 {
		return 0, errors.New("dados inválidos")
	}
	return s.repo.Create(ctx, p)
}

func (s *ProdutosService) Atualizar(ctx context.Context, p *models.Produto) error {
	if p.ID == 0 {
		return errors.New("id obrigatório")
	}
	return s.repo.Update(ctx, p)
}

func (s *ProdutosService) Excluir(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
