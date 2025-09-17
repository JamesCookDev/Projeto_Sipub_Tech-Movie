package services

import (
	"context"
	"testing"
	"errors"

	"github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/core/domain"
	"github.com/jamescookdev/projeto-sipub-tech/movies-service/internal/core/ports/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateMovie_Success(t *testing.T) {

	mockRepo := new(mocks.MovieRepositoryMock)

	inputTitle := "O Auto da Compadecida"
	inputYear := 2000

	movieToSave := domain.Movie{
		Title:    inputTitle,
		Year:     inputYear,
	}

	expectedMovie := domain.Movie{
		ID:       "some_generated_id", 
		Title:    inputTitle,
		Year:     inputYear,
	}
	
	
	mockRepo.On("Save", mock.Anything, movieToSave).Return(&expectedMovie, nil)

	movieService := NewMovieService(mockRepo)

	
	result, err := movieService.CreateMovie(context.Background(), inputTitle, inputYear)

	
	assert.NoError(t, err)
	assert.NotNil(t, result) 
	assert.Equal(t, expectedMovie.ID, result.ID) 
	assert.Equal(t, expectedMovie.Title, result.Title) 

	mockRepo.AssertExpectations(t)
}

func TestCreateMovie_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MovieRepositoryMock)

	movieToSave := domain.Movie{
		Title: "O Auto da Compadecida",
		Year:  2000,
	}

	expectatedError := errors.New("database error")
	mockRepo.On("Save", mock.Anything, movieToSave).Return(nil, expectatedError)

	movieService := NewMovieService(mockRepo)
	
	result, err := movieService.CreateMovie(context.Background(), "O Auto da Compadecida", 2000)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectatedError, err)
	
	mockRepo.AssertExpectations(t)
}


func TestListMovies_Success(t *testing.T) {
	mockRepo := new(mocks.MovieRepositoryMock)

	expectedMovies := []domain.Movie{
		{ID: "id1", Title: "Filme 1", Year: 2001},
		{ID: "id2", Title: "Filme 2", Year: 2002},
	}
	
	var limit, offset int64 = 10, 0

	mockRepo.On("GetAll", mock.Anything, limit, offset).Return(expectedMovies, nil)

	movieService := NewMovieService(mockRepo)

	result, err := movieService.ListMovies(context.Background(), limit, offset)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2) 
	assert.Equal(t, expectedMovies, result) 

	mockRepo.AssertExpectations(t)
}


func TestGetMovie(t *testing.T) {
	mockRepo := new(mocks.MovieRepositoryMock)

	expectedMovie := domain.Movie{
		ID: "filme_encontrado_id", Title: "Filme Encontrado",
	}
	mockRepo.On("Get", mock.Anything, "filme_encontrado_id").Return(&expectedMovie, nil)

	expectedError := errors.New("filme não encontrado")
	mockRepo.On("Get", mock.Anything, "id_inexistente").Return(nil, expectedError)
	
	movieService := NewMovieService(mockRepo)

	testCases := []struct {
		name          string
		inputID       string
		expectedMovie *domain.Movie
		expectedError error
	}{
		{
			name:          "Sucesso - Filme Encontrado",
			inputID:       "filme_encontrado_id",
			expectedMovie: &expectedMovie,
			expectedError: nil,
		},
		{
			name:          "Falha - Filme Não Encontrado",
			inputID:       "id_inexistente",
			expectedMovie: nil,
			expectedError: expectedError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := movieService.GetMovie(context.Background(), tc.inputID)
			assert.Equal(t, tc.expectedMovie, result)
			assert.Equal(t, tc.expectedError, err)
		})
	}

	mockRepo.AssertExpectations(t)
}

func TestDeleteMovie(t *testing.T) {
	mockRepo := new(mocks.MovieRepositoryMock)
	movieService := NewMovieService(mockRepo)

	validID := "id_valido_para_deletar"
	notFoundID := "id_nao_encontrado"
	expectedNotFoundError := errors.New("filme não encontrado")


	mockRepo.On("Delete", mock.Anything, validID).Return(nil)

	mockRepo.On("Delete", mock.Anything, notFoundID).Return(expectedNotFoundError)

	testCases := []struct {
		name          string
		inputID       string
		expectedError error
	}{
		{
			name:          "Sucesso - Filme Deletado",
			inputID:       validID,
			expectedError: nil,
		},
		{
			name:          "Falha - Filme Não Encontrado",
			inputID:       notFoundID,
			expectedError: expectedNotFoundError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			err := movieService.DeleteMovie(context.Background(), tc.inputID)

			assert.Equal(t, tc.expectedError, err)
		})
	}

	mockRepo.AssertExpectations(t)
}