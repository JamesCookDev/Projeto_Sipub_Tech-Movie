	package repository

	import (
		"context"
		"errors"
		"log"

		"github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/core/domain"
		"github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/core/ports"
		"go.mongodb.org/mongo-driver/bson"
		"go.mongodb.org/mongo-driver/bson/primitive"
		"go.mongodb.org/mongo-driver/mongo"
		"go.mongodb.org/mongo-driver/mongo/options"
	)

	// Permitem que as camadas superiores (serviço, gRPC) possam tratar os erros de forma específica.
	var (
		ErrInvalidIDFormat = errors.New("Formato de ID de filme inválido")
		ErrMovieNotFound   = errors.New("Filme não encontrado")
		ErrFetchingMovies  = errors.New("Erro ao buscar filmes")
		ErrDecodingMovies  = errors.New("Erro ao decodificar filmes")
	)


	// mongoRepository é a implementação da interface `ports.MovieRepository`.
	type mongoRepository struct {
		collection *mongo.Collection
	}

	// NewMongoRepository é o construtor para o mongoRepository.
	func NewMongoRepository(db *mongo.Database) (ports.MovieRepository, error) {
		return &mongoRepository{
			collection: db.Collection("movies"),
		}, nil
	}

	func (r *mongoRepository) Get(ctx context.Context, id string) (*domain.Movie, error) {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, ErrInvalidIDFormat
		}

		var movie domain.Movie
		err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&movie)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, ErrMovieNotFound
			}
			return nil, err
		}
		return &movie, nil
	}

	func (r *mongoRepository) GetAll(ctx context.Context, limit, offset int64) ([]domain.Movie, error) {

		findOptions := options.Find()
		findOptions.SetLimit(limit)
		findOptions.SetSkip(offset)

		cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
		if err != nil {
			log.Printf("MongoDB Find error: %v", err)
			return nil, ErrFetchingMovies
		}
		defer cursor.Close(ctx)

		var movies []domain.Movie
		if err = cursor.All(ctx, &movies); err != nil {
			log.Printf("MongoDB All error: %v", err)
			return nil, ErrDecodingMovies
		}

		if movies == nil {
			return []domain.Movie{}, nil
		}

		return movies, nil
	}

	func (r *mongoRepository) Save(ctx context.Context, movie domain.Movie) (*domain.Movie, error) {
		if movie.ID == "" {
			res, err := r.collection.InsertOne(ctx, movie)
			if err != nil {
				return nil, err
			}
			movie.ID = res.InsertedID.(primitive.ObjectID).Hex()
			return &movie, nil
		}

		objectID, err := primitive.ObjectIDFromHex(movie.ID)
		if err != nil {
			return nil, ErrInvalidIDFormat
		}

		_, err = r.collection.ReplaceOne(ctx, bson.M{"_id": objectID}, movie)
		if err != nil {
			return nil, err
		}
		return &movie, nil
	}

	func (r *mongoRepository) Delete(ctx context.Context, id string) error {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return ErrInvalidIDFormat
		}

		res, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
		if err != nil {
			return err
		}
		if res.DeletedCount == 0 {
			return ErrMovieNotFound
		}
		return nil
	}