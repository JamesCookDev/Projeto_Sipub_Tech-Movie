package ports

import (
	"context"
	"github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/core/domain"
)

// MovieRepository é a "Porta de Saída" para a persistência de dados.
type MovieRepository interface {
	Get(ctx context.Context, id string) (*domain.Movie, error)
    GetAll(ctx context.Context, limit, offset int64) ([]domain.Movie, error) 
	Save(ctx context.Context, movie domain.Movie) (*domain.Movie, error)
	Delete(ctx context.Context, id string) error

}
// MovieService é a "Porta de Entrada" para a lógica de negócio.
type MovieService interface {
	GetMovie(ctx context.Context, id string) (*domain.Movie, error)
    ListMovies(ctx context.Context, limit, offset int64) ([]domain.Movie, error) 
	CreateMovie(ctx context.Context, title string, year int) (*domain.Movie, error)
	DeleteMovie(ctx context.Context, id string) error
}