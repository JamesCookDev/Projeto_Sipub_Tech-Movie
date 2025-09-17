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
- **Kubernetes (Minikube):** Orquestração de contêineres para simulação e deploy em ambiente de produção.

---

## ▶️ Como Rodar a Aplicação

Este projeto pode ser executado de duas formas principais:

1. **Docker Compose:** Para desenvolvimento local rápido.
2. **Kubernetes (Minikube):** Para simular um ambiente de produção.

Certifique-se de ter os pré-requisitos instalados antes de começar.

### Pré-requisitos

- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) (opcional, para Kubernetes)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/) (opcional, para Kubernetes)

### Configuração Inicial

1. Clone este repositório:
    ```bash
    git clone https://github.com/JamesCookDev/Projeto_Sipub_Tech-Movie
    ```
2. Acesse a pasta raiz do projeto:
    ```bash
    cd Projeto_Sipub_Tech-Movie
    ```
3. Crie o arquivo `.env` a partir do exemplo:
    ```bash
    cp .env.example .env
    ```

---

### Ambiente 1: Docker Compose

O método mais simples para rodar localmente.

#### Subindo a aplicação

Construa as imagens e inicie os serviços em segundo plano:
```bash
docker compose up --build -d
```
A API estará disponível em [http://localhost:8080](http://localhost:8080).

#### Logs dos serviços

Acompanhe os logs em tempo real:
```bash
docker compose logs -f
```

#### Parando a aplicação

Para parar todos os contêineres:
```bash
docker compose down
```

#### Limpeza completa (incluindo volumes)

Para remover os contêineres e os dados do banco:
```bash
docker compose down -v
```

---

### Ambiente 2: Kubernetes (Minikube)

Simula um ambiente de produção localmente.

> **Importante:** Certifique-se de que o Docker Compose está parado (`docker compose down`) antes de iniciar.

#### Iniciando o cluster

Inicie o Minikube:
```bash
minikube start --driver=docker
```

Configure o terminal para usar o Docker do Minikube:
```bash
eval $(minikube -p minikube docker-env)
```

#### Build das imagens no ambiente Minikube

```bash
docker compose build
```

#### Deploy dos manifestos Kubernetes

```bash
kubectl apply -f k8s/
```

#### Acompanhando os pods

```bash
kubectl get pods -w
```

#### Obtendo a URL da API Gateway

```bash
minikube service api-gateway-service --url
```
Use a URL retornada para acessar a API.

#### Parando e limpando o ambiente

Remova os recursos do cluster:
```bash
kubectl delete -f k8s/
```

Pare o Minikube:
```bash
minikube stop
```

Para deletar o cluster completamente (opcional):
```bash
minikube delete --all
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