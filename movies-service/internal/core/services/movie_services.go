package services

import (
	"context"

	"github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/core/domain"
	"github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/core/ports"
)

// movieService é a implementação concreta da interface `ports.MovieService`.
type movieService struct {
	repo ports.MovieRepository
}

// NewMovieService é o "construtor" para o nosso serviço de filmes.
func NewMovieService(repo ports.MovieRepository) ports.MovieService {
	return &movieService{repo: repo}
}

func (s *movieService) GetMovie(ctx context.Context, id string) (*domain.Movie, error) {
	return s.repo.Get(ctx, id)
}

func (s *movieService) ListMovies(ctx context.Context, limit, offset int64) ([]domain.Movie, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *movieService) CreateMovie(ctx context.Context, title string, year int) (*domain.Movie, error) {
	movie := domain.Movie{
		Title:    title,
		Year:     year,
	}
	return s.repo.Save(ctx, movie)
}

func (s *movieService) DeleteMovie(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}


