package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"loja-produtos/src/internal/models"
	"loja-produtos/src/internal/repository"
)

var _ repository.ProdutosRepository = (*ProdutosPg)(nil)

type ProdutosPg struct {
	db *sql.DB
}

func NewProdutosPg(db *sql.DB) *ProdutosPg {
	return &ProdutosPg{db: db}
}

func (r *ProdutosPg) GetAll(ctx context.Context) ([]models.Produto, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, nome, descricao, preco, quantidade
		FROM produtos
		ORDER BY nome`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Produto
	for rows.Next() {
		var p models.Produto
		if err := rows.Scan(&p.ID, &p.Nome, &p.Descricao, &p.Preco, &p.Quantidade); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *ProdutosPg) GetByID(ctx context.Context, id int64) (*models.Produto, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var p models.Produto
	err := r.db.QueryRowContext(ctx, `
		SELECT id, nome, descricao, preco, quantidade
		FROM produtos WHERE id = $1`, id).
		Scan(&p.ID, &p.Nome, &p.Descricao, &p.Preco, &p.Quantidade)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProdutosPg) Create(ctx context.Context, p *models.Produto) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var id int64
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO produtos (nome, descricao, preco, quantidade)
		VALUES ($1,$2,$3,$4)
		RETURNING id`,
		p.Nome, p.Descricao, p.Preco, p.Quantidade).
		Scan(&id)
	return id, err
}

func (r *ProdutosPg) Update(ctx context.Context, p *models.Produto) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `
		UPDATE produtos
		SET nome=$1, descricao=$2, preco=$3, quantidade=$4
		WHERE id=$5`,
		p.Nome, p.Descricao, p.Preco, p.Quantidade, p.ID)
	return err
}

func (r *ProdutosPg) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, `DELETE FROM produtos WHERE id=$1`, id)
	return err
}
