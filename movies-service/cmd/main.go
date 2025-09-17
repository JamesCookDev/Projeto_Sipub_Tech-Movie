package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"log"
	"net"
	"time"


	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"


	pb "github.com/jamescookdev/projeto-sipub-tech/movies-service/gen/go"
	grpcAdapter "github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/adapters/grpc"
	mongoAdapter "github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/adapters/mongodb"
	"github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/core/services"
	domain "github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/core/domain"
)


func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}


func main() {
	mongodbURI := getEnv("MONGODB_URI", "mongodb://mongodb:27017")
	port := getEnv("PORT", ":50051")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("failed to ping mongo: %v", err)
	}

	db := client.Database("moviedb")
	log.Println("Conectado ao MongoDB")

	seedDatabase(ctx, db)	

	movieRepository, err := mongoAdapter.NewMongoRepository(db)
	if err != nil {
		log.Fatalf("failed to create mongo repository: %v", err)
	}

	movieService := services.NewMovieService(movieRepository)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	grpcServerAdapter := grpcAdapter.NewGRPCServerAdapter(movieService)
	grpcServer := grpc.NewServer()


	pb.RegisterMovieServiceServer(grpcServer, grpcServerAdapter)

	reflection.Register(grpcServer)

	fmt.Printf("gRPC server listening on %s\n", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}

type MovieSeed struct {
    ID       interface{} `json:"_id"`
    Title    string      `json:"title"`
    Year     int     	 `json:"year,string"`
}

func seedDatabase(ctx context.Context, db *mongo.Database) {
	log.Println("Iniciando a verificação do banco de dados para popular...")
    collection := db.Collection("movies")
    count, err := collection.CountDocuments(ctx, bson.M{})
    if err != nil {
		log.Fatalf("Erro ao verificar documentos na coleção: %v", err)
	}

	if count > 0 {
		log.Println("O banco de dados já contém dados. População não necessária.")
		return
	}
    
    if count > 0 {
        log.Println("O banco de dados já contém dados. População não necessária.")
        return
    }

    log.Println("Banco de dados vazio. Tentando ler o arquivo /app/data/movies.json...")
    file, err := os.ReadFile("/app/data/movies.json")
    if err != nil {
        log.Fatalf("Erro ao ler o arquivo de seed /app/movies.json: %v", err)
    }


    var movieSeeds []MovieSeed
    log.Println("Tentando decodificar o JSON...")
    if err := json.Unmarshal(file, &movieSeeds); err != nil {
        log.Fatalf("Erro ao decodificar o JSON: %v", err)
    }

    log.Printf("JSON decodificado com sucesso. Encontrados %d filmes no arquivo.", len(movieSeeds))


    var docs []interface{}
    for _, seed := range movieSeeds {
			doc := domain.Movie{
				Title:    seed.Title,
				Year:     seed.Year,
		}
        docs = append(docs, doc)
    }

	if len(docs) > 0 {
		log.Printf("Tentando inserir %d documentos convertidos do JSON no banco de dados MongoDB...", len(docs))
		_, err = collection.InsertMany(ctx, docs)
		if err != nil {
			log.Fatalf("Erro ao inserir dados no banco: %v", err)
		}
	}

    log.Printf("Banco de dados populado com sucesso com %d filmes.", len(movieSeeds))
}