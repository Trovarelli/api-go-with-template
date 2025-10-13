package main

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"loja-produtos/src/internal/config"
	"loja-produtos/src/internal/handler"
	"loja-produtos/src/internal/repository/postgres"
	service "loja-produtos/src/internal/services"

	_ "github.com/lib/pq"
)

var temp = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	db := config.ConnectDatabase()
	defer db.Close()

	repo := postgres.NewProdutosPg(db)
	svc := service.NewProdutosService(repo)
	h := handler.NewProdutosHandler(svc)

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler(svc))
	mux.HandleFunc("/produtos/", h.ProdutoByID)
	mux.HandleFunc("/produtos", h.Produtos)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           logging(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func indexHandler(svc *service.ProdutosService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := svc.Listar(r.Context())
		if err != nil {
			http.Error(w, "erro ao buscar produtos", http.StatusInternalServerError)
			return
		}

		if err := temp.ExecuteTemplate(w, "Index", items); err != nil {
			http.Error(w, "erro ao renderizar template", http.StatusInternalServerError)
		}
	}
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
		}()
		next.ServeHTTP(w, r)
	})
}
