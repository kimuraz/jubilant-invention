package main

import (
	"context"
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"time"
)

type MoviesAPI struct {
	redisClient *redis.Client
}

type ConfirmMoviesRequest struct {
	RequestID string  `json:"requestId"`
	Movies    []Movie `json:"movies"`
}

type GenresPreferencesEvent struct {
	RequestID      string   `json:"requestId"`
	SelectedGenres []string `json:"selectedGenres"`
}

func (api *MoviesAPI) MoviesList(c *gin.Context) {
	var genres []string
	requestId := c.Query("requestId")
	RequestIDMap.Range(func(key, value interface{}) bool {
		if key == requestId {
			genres = value.([]string)
			return false
		}
		return true
	})

	var movies []Movie
	for _, movie := range MoviesDB {
		for _, genre := range genres {
			if movie.Genre == genre {
				movies = append(movies, movie)
			}
		}
	}
}

func (api *MoviesAPI) ConfirmMovies(c *gin.Context) {
	var req ConfirmMoviesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jsonData, _ := json.Marshal(req)
	api.redisClient.Publish(c.Request.Context(), "movies_channel", jsonData)
	RequestIDMap.Delete(req.RequestID)
	c.JSON(200, gin.H{"message": "Movies selection processed"})
}

func (api *MoviesAPI) SubscribeToChannels() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	pubsub := rdb.Subscribe(context.Background(), "movies_genres_channel")
	ch := pubsub.Channel()

	for msg := range ch {
		var genresEvent GenresPreferencesEvent
		json.Unmarshal([]byte(msg.Payload), &genresEvent)
		RequestIDMap.Store(genresEvent.RequestID, genresEvent.SelectedGenres)
	}
}

func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "redis:6379"})
	api := &MoviesAPI{redisClient: rdb}
	go api.SubscribeToChannels()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/movies", (&MoviesAPI{redisClient: rdb}).MoviesList)
	r.POST("/movies", (&MoviesAPI{redisClient: rdb}).ConfirmMovies)
	r.GET("/backdoor/db", func(c *gin.Context) {
		c.JSON(200, RequestIDMap)
	})

	r.Run(":8083")
}
