package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"time"
)

type GenresAPI struct {
	redisClient *redis.Client
}

type Genres string

const (
	Action    Genres = "action"
	Adventure Genres = "adventure"
	Comedy    Genres = "comedy"
	Drama     Genres = "drama"
	Horror    Genres = "horror"
)

type GenreSelectionRequest struct {
	RequestID      string   `json:"requestId"`
	SelectedGenres []Genres `json:"selectedGenres"`
}

func (api *GenresAPI) GenreList(c *gin.Context) {
	c.JSON(200, gin.H{"genres": []Genres{Action, Adventure, Comedy, Drama, Horror}})
}

func (api *GenresAPI) GenreSelection(c *gin.Context) {
	var req GenreSelectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	api.redisClient.Publish(c.Request.Context(), "genres_channel", req.RequestID)
	api.redisClient.Publish(c.Request.Context(), "movies_genres_channel", req.RequestID)
	c.JSON(200, gin.H{"message": "Genres selection processed"})
}

func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "redis:6379"})

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/genres", (&GenresAPI{redisClient: rdb}).GenreList)
	r.POST("/genres", (&GenresAPI{redisClient: rdb}).GenreSelection)

	r.Run(":8082")
}
