# Loja de Produtos (API em Go)

Um projeto de estudo em Go que expõe uma API simples para gerenciamento de produtos, usando arquitetura em camadas (handler → service → repository) e PostgreSQL como banco de dados. Inclui uma página HTML minimalista para listar produtos na rota `/`.

## Visão Geral

- Servidor HTTP com `net/http` (porta `8080`).
- Persistência em PostgreSQL via `database/sql` (repositório Postgres puro).
- Seed inicial do banco via `init.sql` (carregado pelo Docker Compose).
- Helpers para JSON (decodificação segura, erros padronizados, respostas com `Content-Type`).
- Página HTML em `templates/index.html` para visualizar a lista de produtos.

## Tecnologias

- Go 1.21+
- PostgreSQL 13 (via Docker Compose)
- `github.com/joho/godotenv` (carrega `.env`)
- `github.com/lib/pq` (driver Postgres)
- GORM é usado apenas no teste de conexão em `src/internal/config/database_test.go`.

## Arquitetura

- `src/internal/handler`: entradas HTTP (controllers/handlers)
- `src/internal/services`: regras de negócio (use-cases)
- `src/internal/repository`: portas e implementações (Postgres em `repository/postgres`)
- `src/internal/models`: modelos de domínio (ex.: `Produto`)
- `src/internal/config`: configuração e conexão ao banco
- `src/internal/helpers`: utilitários (JSON, parsing)
- `templates/`: templates HTML

## Endpoints

Base URL: `http://localhost:8080`

- GET `/` → Renderiza página HTML com a listagem de produtos.
- GET `/produtos` → Lista todos os produtos (JSON).
- GET `/produtos/{id}` → Busca produto por ID (JSON).

Modelo JSON de `Produto`:

```json
{
  "id": 1,
  "nome": "Teclado Mecânico",
  "descricao": "Switch blue",
  "preco": 199.9,
  "quantidade": 10
}
```

Exemplos (cURL):

```bash
# Listar
curl -s http://localhost:8080/produtos | jq .

# Buscar por ID
curl -s http://localhost:8080/produtos/1 | jq .
```

## Como Rodar

Pré-requisitos:

- Docker + Docker Compose
- Go 1.21+
- `make` (opcional, facilita os comandos)

1) Configure o `.env` (você pode usar `.env.example` como base):

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=12345678
POSTGRES_DB=alura_loja
PG_HOST_PORT=5432
DB_DSN="host=localhost user=postgres password=12345678 dbname=alura_loja port=5432 sslmode=disable TimeZone=UTC"
```

2) Suba apenas o banco de dados (Docker Compose aplica o `init.sql` automaticamente):

```bash
make db-up
# ou
docker compose up -d db
```

3) Rode a aplicação localmente (fora do Docker):

```bash
make run-local
# ou
go run ./src
```

4) Acesse:

- Página HTML: `http://localhost:8080/`
- API: `http://localhost:8080/produtos`

Para desligar/remover volumes do banco (cuidado: apaga os dados):

```bash
make db-down
# ou
docker compose down -v
```

Para recriar o banco (dropa volume e reaplica o seed):

```bash
make db-reseed
```

## Testes e Lint

- Testes (valida conexão com o banco; requer DB rodando e `.env`):

```bash
make test
# ou
go test ./src/internal/config -v
```

- Lint (usa container com `golangci-lint`):

```bash
make lint
```

## Estrutura do Projeto

```
.
├── src
│   ├── internal
│   │   ├── config
│   │   │   ├── database.go
│   │   │   └── database_test.go
│   │   ├── handler
│   │   │   └── produtos.go
│   │   ├── helpers
│   │   │   ├── json.go
│   │   │   └── parseIdFromPath.go
│   │   ├── models
│   │   │   └── produto.go
│   │   ├── repository
│   │   │   ├── postgres
│   │   │   │   └── produtos_pg.go
│   │   │   └── produtos.go
│   │   └── services
│   │       └── produtos.go
│   └── main.go
├── templates
│   └── index.html
├── init.sql
├── docker-compose.yml
├── Makefile
├── .env.example
└── .env
```

## Variáveis de Ambiente

- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `PG_HOST_PORT`: usados pelo Docker Compose para subir o Postgres.
- `DB_DSN`: DSN completo do Postgres, lido pela aplicação Go (ex.: `host=localhost user=postgres password=... dbname=... port=5432 sslmode=disable TimeZone=UTC`).

## Observações

- A aplicação sobe em `:8080` e loga cada requisição com latência.
- O seed (`init.sql`) cria a tabela `produtos` e insere dados de exemplo.
- A página HTML utiliza Bootstrap via CDN apenas para visualização simples.
- O repositório usa `database/sql` diretamente; o GORM é utilizado somente nos testes de conexão.

## Próximos Passos (Sugestões)

- Implementar rotas de escrita (POST/PUT/DELETE).
- Corrigir caracteres com encoding incorreto em mensagens/seed.
- Melhorar validações e mensagens de erro na camada HTTP.
- Adicionar paginação/filtros na listagem de produtos.
- Adicionar Dockerfile para rodar a aplicação Go em container.
- Criar testes de integração para os endpoints HTTP.

---

Projeto criado para fins didáticos/estudo. Sugestões e PRs são bem-vindos!

