package main

import (
	"log"
	"os"

	_ "github.com/jamescookdev/projeto-sipub-tech/api/docs"
	"github.com/jamescookdev/projeto-sipub-tech/api/handlers"
	"github.com/jamescookdev/projeto-sipub-tech/api/messaging"
	pb "github.com/jamescookdev/projeto-sipub-tech/movies-service/gen/go"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title           API de Filmes - Microsserviços com Go e gRPC
// @version         1.0
// @description     Esta é uma API REST para consulta e gerenciamento de filmes.
// @host            localhost:8080
// @BasePath        /
// @schemes         http
// @contact.name    James Cook

func main() {
	moviesServiceAddress := getEnv("MOVIES_SERVICE_ADDRESS", "movies_service:50051")

	// Conexão gRPC para leituras (GETs)
	conn, err := grpc.Dial(moviesServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Nao foi possivel conectar ao movies-service: %v", err)
	}
	defer conn.Close()
	movieClient := pb.NewMovieServiceClient(conn)
	log.Println("Conexao com o movies-service estabelecida com sucesso!")

	// Publisher RabbitMQ para escritas assíncronas (POST/DELETE)
	pub, err := messaging.NewPublisher()
	if err != nil {
		log.Fatalf("Nao foi possivel conectar ao RabbitMQ: %v", err)
	}
	defer pub.Close()

	h := handlers.NewMovieHandler(movieClient, pub)
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rotas
	movieRoutes := router.Group("/movies")
	{
		movieRoutes.GET("", h.ListMovies)         
		movieRoutes.GET("/:id", h.GetMovieByID)
		movieRoutes.POST("", h.CreateMovie)       
		movieRoutes.DELETE("/:id", h.DeleteMovie)  
	}

	// Inicialização do Servidor HTTP
	apiPort := getEnv("API_PORT", ":8080")
	log.Printf("API Gateway rodando na porta %s", apiPort)
	if err := router.Run(apiPort); err != nil {
		log.Fatalf("Nao foi possivel iniciar o servidor da API: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
