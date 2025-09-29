package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"

	"loja-produtos/src/internal/config"
	produto "loja-produtos/src/internal/models"

	_ "github.com/lib/pq"
)

var temp = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	db := config.ConnectDatabase()
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index(w, r, db)
	})

	srv := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Servidor ouvindo em :8080")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("erro no servidor: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("Servidor finalizado")
}

func index(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	items, err := listarProdutos(db)
	if err != nil {
		log.Printf("erro ao buscar produtos: %v", err)
		http.Error(w, "erro ao buscar produtos", http.StatusInternalServerError)
		return
	}

	if err := temp.ExecuteTemplate(w, "Index", items); err != nil {
		log.Printf("erro ao renderizar template: %v", err)
		http.Error(w, "erro ao renderizar template", http.StatusInternalServerError)
		return
	}
}

func listarProdutos(db *sql.DB) ([]produto.Produto, error) {
	rows, err := db.Query(`
        SELECT id, nome, descricao, preco, quantidade
        FROM items
        ORDER BY nome
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []produto.Produto
	for rows.Next() {
		var p produto.Produto
		if err := rows.Scan(&p.Id, &p.Nome, &p.Descricao, &p.Preco, &p.Quantidade); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
