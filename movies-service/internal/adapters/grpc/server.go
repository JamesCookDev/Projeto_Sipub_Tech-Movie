package grpc

import (
	"context"
	"errors"
	"log"

	pb "github.com/jamescookdev/projeto-sipub-tech/movies-service/gen/go"
	repository "github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/adapters/mongodb"
	"github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/core/domain"
	"github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/core/ports"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// serverAdapter é a implementação do servidor gRPC gerado pelo Protobuf.
type serverAdapter struct {
	pb.UnimplementedMovieServiceServer 
	service ports.MovieService
}

// NewGRPCServerAdapter é o construtor do nosso adaptador.
func NewGRPCServerAdapter(service ports.MovieService) pb.MovieServiceServer {
	return &serverAdapter{service: service}
}

// mapDomainErrorToGRPCStatus é uma função auxiliar que traduz os erros internos do nosso domínio
func mapDomainErrorToGRPCStatus(err error) error {
	switch {
		case errors.Is(err, repository.ErrMovieNotFound):
			return status.Error(codes.NotFound, err.Error())
		case errors.Is(err, repository.ErrInvalidIDFormat):
			return status.Error(codes.InvalidArgument, err.Error())
		default:
			return status.Error(codes.Internal, "Um erro interno ocorreu")
	}
}
		
// GetMovie é o handler para a chamada RPC GetMovie.
func (s *serverAdapter) GetMovie(ctx context.Context, req *pb.GetMovieRequest) (*pb.Movie, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Movie ID cannot be empty")
	}

	movie, err := s.service.GetMovie(ctx, req.Id)
	if err != nil {
		return nil, mapDomainErrorToGRPCStatus(err)
	}

	return toGRPCMovie(movie), nil
}

// ListMovies é o handler para a chamada RPC ListMovies
func (s *serverAdapter) ListMovies(ctx context.Context, req *pb.ListMoviesRequest) (*pb.MovieList, error) {

	var limit int64 = 20
	if req.Limit > 0 {
		limit = int64(req.Limit)
	}

	var offset int64 = 0
	if req.Offset > 0 {
		offset = int64(req.Offset)
	}


	movies, err := s.service.ListMovies(ctx, limit, offset)
	if err != nil {
		log.Printf("Error ao listar filmes: %v", err)
		return nil, mapDomainErrorToGRPCStatus(err)
	}
	
	grpcMovies := make([]*pb.Movie, len(movies)) 
	for i, movie := range movies {
    	grpcMovies[i] = toGRPCMovie(&movie)
	}
	
	return &pb.MovieList{Movies: grpcMovies}, nil
}

// CreateMovie é o handler para a chamada RPC CreateMovie.
func (s *serverAdapter) CreateMovie(ctx context.Context, req *pb.CreateMovieRequest) (*pb.Movie, error) {
	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "Titulo não pode ser vazio")
	}

	movie, err := s.service.CreateMovie(ctx, req.Title, int(req.Year))
	if err != nil {
		return nil, mapDomainErrorToGRPCStatus(err)
	}

	return toGRPCMovie(movie), nil
}

// DeleteMovie é o handler para a chamada RPC DeleteMovie.
func (s *serverAdapter) DeleteMovie(ctx context.Context, req *pb.DeleteMovieRequest) (*pb.Empty, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Filme \"ID\" não pode ser vazio")
	}

	
	err := s.service.DeleteMovie(ctx, req.Id)
	if err != nil {
			return nil, mapDomainErrorToGRPCStatus(err)
	}

	return &pb.Empty{}, nil
}

// toGRPCMovie é uma função de conversão que traduz um struct do nosso domínio (`domain.Movie`)
func toGRPCMovie(movie *domain.Movie) *pb.Movie {
	return &pb.Movie{
		Id:       movie.ID,
		Title:    movie.Title,
		Year:     int32(movie.Year),
	}
}