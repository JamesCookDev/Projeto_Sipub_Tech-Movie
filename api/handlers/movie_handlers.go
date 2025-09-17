package handlers


import (
	"log"
	"net/http"
	"strconv"

	pb "github.com/jamescookdev/projeto-sipub-tech/movies-service/gen/go"
	

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateMovieRequest struct {
	Title    string `json:"title" binding:"required" example:"Interestelar"`
	Year     int32  `json:"year" binding:"required" example:"2014"`
}

type MovieHandler struct {
	MovieClient pb.MovieServiceClient
}

func NewMovieHandler(client pb.MovieServiceClient) *MovieHandler {
	return &MovieHandler{MovieClient: client}
}

// ListMovies      Lista os filmes com paginação
// @Summary      Lista os filmes com paginação
// @Description  Retorna uma lista de filmes, com a possibilidade de usar limit e offset para paginação.
// @Tags         Movies
// @Produce      json
// @Param        limit  query     int    false  "Número de resultados por página" default(20)
// @Param        offset query     int    false  "Número de resultados a pular"    default(0)
// @Success      200    {array}   pb.Movie
// @Failure      500    {object}  map[string]string{error=string}
// @Router       /movies [get]
func (h *MovieHandler) ListMovies(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	offset, _ := strconv.ParseInt(offsetStr, 10, 32)

	grpcRequest := &pb.ListMoviesRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	res, err := h.MovieClient.ListMovies(c.Request.Context(), grpcRequest)
	if err != nil {
		log.Printf("Erro ao chamar gRPC ListMovies: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar a lista de filmes."})
		return
	}
	c.JSON(http.StatusOK, res.Movies)
}

// GetMovieByID      Busca um filme por ID
// @Summary      Busca um filme por ID
// @Description  Retorna os detalhes de um filme específico baseado no seu ID.
// @Tags         Movies
// @Produce      json
// @Param        id   path      string  true  "ID do Filme" Format(mongodb-id)
// @Success      200  {object}  pb.Movie
// @Failure      404  {object}  map[string]string{error=string}
// @Failure      500  {object}  map[string]string{error=string}
// @Router       /movies/{id} [get]
func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	movieID := c.Param("id")
	grpcRequest := &pb.GetMovieRequest{Id: movieID}

	res, err := h.MovieClient.GetMovie(c.Request.Context(), grpcRequest)
	if err != nil {
		log.Printf("Erro ao chamar gRPC GetMovie: %v", err)
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Filme não encontrado."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar o filme."})
		return
	}
	c.JSON(http.StatusOK, res)
}

// CreateMovie      Cria um novo filme
// @Summary      Cria um novo filme
// @Description  Adiciona um novo filme ao banco de dados.
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        movie  body      CreateMovieRequest  true  "Dados para criar o filme"
// @Success      201    {object}  pb.Movie
// @Failure      400    {object}  map[string]string{error=string}
// @Failure      500    {object}  map[string]string{error=string}
// @Router       /movies [post]
func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var req CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	grpcRequest := &pb.CreateMovieRequest{
		Title:    req.Title,
		Year:     req.Year,
	}
	res, err := h.MovieClient.CreateMovie(c.Request.Context(), grpcRequest)
	if err != nil {
		log.Printf("Erro ao chamar gRPC CreateMovie: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar o filme."})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// DeleteMovie   Deleta um filme
// @Summary      Deleta um filme
// @Description  Remove um filme do banco de dados baseado no seu ID.
// @Tags         Movies
// @Produce      json
// @Param        id   path      string  true  "ID do Filme" Format(mongodb-id)
// @Success      200  {object}  map[string]string{message=string}
// @Failure      404  {object}  map[string]string{error=string}
// @Failure      500  {object}  map[string]string{error=string}
// @Router       /movies/{id} [delete]
func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	movieID := c.Param("id")
	grpcRequest := &pb.DeleteMovieRequest{Id: movieID}

	_, err := h.MovieClient.DeleteMovie(c.Request.Context(), grpcRequest)
	if err != nil {
		log.Printf("Erro ao chamar gRPC DeleteMovie: %v", err)
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Filme não encontrado para deletar."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar o filme."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Filme deletado com sucesso."})
}