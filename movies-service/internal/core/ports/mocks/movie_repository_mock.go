package mocks

import (
	"context"

	"github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/core/domain"
	"github.com/stretchr/testify/mock"

)

type MovieRepositoryMock struct {
	mock.Mock
}

func (m *MovieRepositoryMock) Get(ctx context.Context, id string) (*domain.Movie, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Movie), args.Error(1)
	
}

func (m *MovieRepositoryMock) GetAll(ctx context.Context, limit, offset int64) ([]domain.Movie, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Movie), args.Error(1)
}

func (m *MovieRepositoryMock) Save(ctx context.Context, movie domain.Movie) (*domain.Movie, error) {
	args := m.Called(ctx, movie)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Movie), args.Error(1)
}

func (m *MovieRepositoryMock) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}