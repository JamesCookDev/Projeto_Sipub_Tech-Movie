// Local: api/main.go

package main

import (
	"log"

	_ "github.com/jamescookdev/projeto-sipub-tech/api/docs"
	"github.com/jamescookdev/projeto-sipub-tech/api/handlers"
	pb "github.com/jamescookdev/projeto-sipub-tech/movies-service/gen/go"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	moviesServiceAddress = "movies_service:50051"
)

// @title           API de Filmes - Microsserviços com Go e gRPC
// @version         1.0
// @description     Esta é uma API REST para consulta e gerenciamento de filmes.
// @host            localhost:8080
// @BasePath        /

//  Conexão com o Serviço de Backend (gRPC) e inicialização do servidor HTTP
func main() {
	conn, err := grpc.NewClient(moviesServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Nao foi possivel conectar ao movies-service: %v", err)
	}
	defer conn.Close()

	movieClient := pb.NewMovieServiceClient(conn)
	log.Println("Conexao com o movies-service estabelecida com sucesso!")

	// Configuração dos Handlers e Rotas
	movieHandler := handlers.NewMovieHandler(movieClient)

	router := gin.Default()

	// Definição das Rotas (Endpoints)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	movieRoutes := router.Group("/movies")
	{
		movieRoutes.GET("", movieHandler.ListMovies)
		movieRoutes.GET("/:id", movieHandler.GetMovieByID)
		movieRoutes.POST("", movieHandler.CreateMovie)
		movieRoutes.DELETE("/:id", movieHandler.DeleteMovie)
	}

	// Inicialização do Servidor HTTP
	log.Println("API Gateway rodando na porta :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Nao foi possivel iniciar o servidor da API: %v", err)
	}
}