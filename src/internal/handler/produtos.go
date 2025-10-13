package handler

import (
	"errors"
	"fmt"
	"loja-produtos/src/internal/models"
	service "loja-produtos/src/internal/services"
	"net/http"

	helpers "loja-produtos/src/internal/helpers"
)

type ProdutosHandler struct {
	svc *service.ProdutosService
}

func NewProdutosHandler(svc *service.ProdutosService) *ProdutosHandler {
	return &ProdutosHandler{svc: svc}
}

func (h *ProdutosHandler) Produtos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listar(w, r)
	case http.MethodPost:
		h.criar(w, r)
	default:
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
	}
}

func (h *ProdutosHandler) ProdutoByID(w http.ResponseWriter, r *http.Request) {
	id, ok := helpers.ParseIDFromPath(r.URL.Path, "/produtos/")
	if !ok {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.obter(w, r, id)
	case http.MethodPut:
		h.atualizar(w, r, id)
	case http.MethodDelete:
		h.excluir(w, r, id)
	default:
		http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
	}
}

func (h *ProdutosHandler) listar(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.Listar(r.Context())
	if err != nil {
		http.Error(w, "erro ao listar", http.StatusInternalServerError)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, items)
}

func (h *ProdutosHandler) obter(w http.ResponseWriter, r *http.Request, id int64) {
	p, err := h.svc.BuscarPorID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrProdutoNaoEncontrado) {
			http.Error(w, "não encontrado", http.StatusNotFound)
			return
		}
		http.Error(w, "erro ao obter", http.StatusInternalServerError)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, p)
}

func (h *ProdutosHandler) criar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var p *models.Produto
	if err := helpers.DecodeJSON(w, r, p); err != nil {
		helpers.WriteError(w, err)
		return
	}

	id, err := h.svc.Criar(r.Context(), p)
	if err != nil {
		helpers.WriteError(w, helpers.BadRequest(err.Error()))
		return
	}

	p.ID = id
	w.Header().Set("Location", fmt.Sprintf("/produtos/%d", id))
	helpers.WriteJSON(w, http.StatusCreated, p)
}
func (h *ProdutosHandler) atualizar(w http.ResponseWriter, r *http.Request, id int64) {
	var p *models.Produto
	if err := helpers.DecodeJSON(w, r, p); err != nil {
		helpers.WriteError(w, err)
		return
	}
	p.ID = id
	if err := h.svc.Atualizar(r.Context(), p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProdutosHandler) excluir(w http.ResponseWriter, r *http.Request, id int64) {
	if err := h.svc.Excluir(r.Context(), id); err != nil {
		http.Error(w, "erro ao excluir", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
