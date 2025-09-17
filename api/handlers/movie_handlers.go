package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jamescookdev/projeto-sipub-tech/api/messaging"
	pb "github.com/jamescookdev/projeto-sipub-tech/movies-service/gen/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateMovieRequest struct {
	Title string `json:"title" binding:"required" example:"Interestelar"`
	Year  int32  `json:"year"  binding:"required" example:"2014"`
}

// Evento padronizado publicado no RabbitMQ
type MovieEvent struct {
	Action    string      `json:"action"`   // "create" | "delete"
	Data      interface{} `json:"data"`     // payload do evento
	Timestamp time.Time   `json:"timestamp"`
}

type MovieHandler struct {
	MovieClient pb.MovieServiceClient  // Leituras (GET) continuam síncronas via gRPC
	Publisher   *messaging.Publisher   // Escritas (POST/DELETE) publicam eventos
}

func NewMovieHandler(client pb.MovieServiceClient, pub *messaging.Publisher) *MovieHandler {
	return &MovieHandler{MovieClient: client, Publisher: pub}
}

// ListMovies
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
    titleQ := strings.TrimSpace(c.Query("title"))
    yearQ  := strings.TrimSpace(c.Query("year"))

    limitStr := c.DefaultQuery("limit", "10")
    offsetStr := c.DefaultQuery("offset", "0")
    limit, _ := strconv.ParseInt(limitStr, 10, 32)
    offset, _ := strconv.ParseInt(offsetStr, 10, 32)

    grpcRequest := &pb.ListMoviesRequest{
        Limit:  int32(limit),
        Offset: int32(offset),
    }

    res, err := h.MovieClient.ListMovies(c.Request.Context(), grpcRequest)
    if err != nil {
        log.Printf("Erro ao ListMovies: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar filmes"})
        return
    }

    if titleQ == "" && yearQ == "" {
        if res.Movies == nil {
            c.JSON(http.StatusOK, []any{})
            return
        }
        c.JSON(http.StatusOK, res.Movies)
        return
    }

    out := make([]*pb.Movie, 0) 
    for _, m := range res.Movies {
        if m == nil {
            continue
        }
        ok := true
        if titleQ != "" && !strings.Contains(strings.ToLower(m.Title), strings.ToLower(titleQ)) {
            ok = false
        }
        if yearQ != "" {
            if y, err := strconv.Atoi(yearQ); err == nil && int32(y) != m.Year {
                ok = false
            }
        }
        if ok {
            out = append(out, m)
        }
    }
    c.JSON(http.StatusOK, out)
}

// GetMovieByID
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

// CreateMovie (ASSÍNCRONO)
// @Summary      Solicita a criação de um novo filme (assíncrono)
// @Description  Envia um evento para criação de filme. A operação é processada em background.
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        movie  body      CreateMovieRequest  true  "Dados para criar o filme"
// @Success      202    {object}  map[string]string{message=string}
// @Failure      400    {object}  map[string]string{error=string}
// @Failure      500    {object}  map[string]string{error=string}
// @Router       /movies [post]
func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var req CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	evt := MovieEvent{
		Action:    "create",
		Data:      req,                // payload simples; o consumer mapeia para o modelo/domínio
		Timestamp: time.Now().UTC(),
	}
	body, _ := json.Marshal(evt)

	if err := h.Publisher.Publish(c.Request.Context(), h.Publisher.RoutingKeyCreated(), body); err != nil {
		log.Printf("Erro ao publicar movie.created: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao enfileirar criação"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "Solicitação de criação recebida e sendo processada."})
}

// DeleteMovie (ASSÍNCRONO)
// @Summary      Solicita a deleção de um filme (assíncrono)
// @Description  Envia um evento para deletar um filme. A operação é processada em background.
// @Tags         Movies
// @Produce      json
// @Param        id   path      string  true  "ID do Filme" Format(mongodb-id)
// @Success      202  {object}  map[string]string{message=string}
// @Failure      500  {object}  map[string]string{error=string}
// @Router       /movies/{id} [delete]
func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	movieID := c.Param("id")

	evt := MovieEvent{
		Action:    "delete",
		Data:      gin.H{"id": movieID},
		Timestamp: time.Now().UTC(),
	}
	body, _ := json.Marshal(evt)

	if err := h.Publisher.Publish(c.Request.Context(), h.Publisher.RoutingKeyDeleted(), body); err != nil {
		log.Printf("Erro ao publicar movie.deleted: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao enfileirar deleção"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "Solicitação de deleção recebida e sendo processada."})
}
