# API de Filmes - Microsserviços com Go, gRPC e Docker

Este repositório contém o código-fonte de uma API RESTful para consulta e gerenciamento de filmes. O projeto foi desenvolvido utilizando uma arquitetura de microsserviços com comunicação via gRPC, seguindo os princípios da Arquitetura Hexagonal para um código limpo e desacoplado.

---

##  Visão Geral da Arquitetura

O sistema é composto por três contêineres Docker orquestrados via Docker Compose:

1. **API Gateway (`api`):** Um serviço em Go (Gin) que expõe uma interface REST pública e atua como cliente gRPC.
2. **Serviço de Filmes (`movies-service`):** Um serviço em Go que contém a lógica de negócio, se comunica com o banco de dados e expõe uma interface gRPC interna.
3. **Banco de Dados (`mongodb`):** Uma instância MongoDB para a persistência dos dados.

O fluxo de comunicação é o seguinte:

`Cliente ➔ [API Gateway (HTTP/REST)] ➔ [Serviço de Filmes (gRPC)] ➔ [MongoDB]`

A arquitetura segue o padrão **Hexagonal**, isolando as camadas de negócios das dependências externas, como o banco de dados e a comunicação gRPC.

---

##  Tecnologias Utilizadas

- **Go:** Linguagem principal para o desenvolvimento dos microsserviços.
- **Docker & Docker Compose:** Para containerização e orquestração do ambiente.
- **MongoDB:** Banco de dados NoSQL.
- **gRPC & Protobuf:** Para a comunicação interna entre os serviços.
- **Gin:** Framework web para a API Gateway.
- **Swaggo:** Ferramenta para geração automática da documentação OpenAPI (Swagger).

---

##  Como Rodar a Aplicação (Guia Rápido)

Siga os passos abaixo para executar o projeto em seu ambiente local.

### Pré-requisitos
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Configuração do Ambiente
O projeto utiliza um arquivo `.env` para gerenciar as variáveis de ambiente.

1. Clone este repositório:
    ```bash
    git clone <https://github.com/JamesCookDev/Projeto_Sipub_Tech-Movie.git>
    ```

2. Navegue até a pasta raiz do projeto.

3. Crie seu arquivo de configuração `.env` a partir do exemplo fornecido:
    ```bash
    cp .env.example .env
    ```
    > Para o ambiente de desenvolvimento padrão, nenhuma alteração é necessária no arquivo `.env`.

### Inicialização com Um Comando
Para iniciar todos os serviços com um único comando, execute:

```bash
docker compose up --build
```
Esse comando irá construir as imagens, baixar as dependências, criar os contêineres e as redes. Após a inicialização, a API estará disponível em [http://localhost:8080](http://localhost:8080).

## Parando a Aplicação

Para parar todos os contêineres, execute:

```bash
docker compose down
```

## Documentação da API (Swagger)

A API possui uma documentação interativa gerada com Swagger (OpenAPI). Com a aplicação em execução, acesse no navegador:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

Na interface do Swagger, você poderá ver todos os endpoints, seus parâmetros, schemas de dados e testar a API diretamente.

## Exemplos de Uso (cURL)

A seguir, exemplos de como interagir com a API via `curl`.

**Listando filmes com paginação (`limit=2`):**

```bash
curl "http://localhost:8080/movies?limit=2"
```

**Listando 3 filmes, pulando os 3 primeiros (segunda página):**

Este exemplo utiliza o parâmetro `offset` para buscar a próxima página de resultados.

```bash
curl "http://localhost:8080/movies?limit=3&offset=3"
```

**Criando um novo filme:**

```bash
curl -X POST http://localhost:8080/movies \
    -H "Content-Type: application/json" \
    -d '{
        "title": "Bacurau",
        "year": 2019
}'
```

> **Nota:** Copie o "id" retornado na resposta para usar nos exemplos seguintes.

**Buscando o filme criado por ID:**

```bash
# Substitua SEU_ID_AQUI pelo ID real do filme
curl http://localhost:8080/movies/SEU_ID_AQUI
```

**Deletando o filme criado:**

```bash
# Substitua SEU_ID_AQUI pelo ID real do filme
curl -X DELETE http://localhost:8080/movies/SEU_ID_AQUI
```

## 🧪 Testes

O projeto contém testes unitários para a camada de serviço, isolando a lógica de negócios com o uso de mocks. Para executar os testes, navegue até a pasta do serviço e rode o comando de teste do Go:

### Testes Unitários Sem Mocks

Para rodar os testes sem mocks (interagindo diretamente com o banco de dados), execute:

```bash
cd movies-service
go test -v ./...
```

### Testes Unitários Com Mocks

Para rodar os testes utilizando mocks (isolando a camada de serviço), execute:

```bash
cd movies-service
go test -v -run "TestMovieServiceWithMocks" ./...
```

## Estrutura do Projeto

A estrutura de pastas principal do projeto é a seguinte:

```
├── api
│   ├── Dockerfile
│   ├── docs
│   │   ├── docs.go
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   ├── go.mod
│   ├── go.sum
│   ├── handlers
│   │   └── movie_handlers.go
│   └── main.go
├── data
│   └── movies.json
├── docker-compose.yml
├── go.work
├── go.work.sum
├── movies-service
│   ├── cmd
│   │   └── main.go
│   ├── Dockerfile
│   ├── gen
│   │   └── go
│   │       ├── movies_grpc.pb.go
│   │       └── movies.pb.go
│   ├── go.mod
│   ├── go.sum
│   └── internal
│       ├── adapters
│       │   ├── grpc
│       │   │   └── server.go
│       │   └── mongodb
│       │       └── mongoRepo.go
│       └── core
│           ├── domain
│           │   └── movie.go
│           ├── ports
│           │   ├── mocks
│           │   │   └── movie_repository_mock.go
│           │   └── ports.go
│           └── services
│               ├── movie_services.go
│               └── movie_services_test.go
├── proto
│   └── movies.proto
└── ReadMe.md
```

## 