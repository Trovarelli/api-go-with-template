package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDatabase() *sql.DB {
	_ = godotenv.Load(".env")

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("variável DB_DSN não definida")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("erro ao abrir conexão com o BD: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("erro ao pingar BD: %v", err)
	}

	log.Println("conectado ao banco com sucesso!")
	return db
}
