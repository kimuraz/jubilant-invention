package main

import (
	"context"
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type OrchestrationAPI struct {
	redisClient *redis.Client
}

type StartProcessRequest struct {
	UserID string `json:"userId"`
}

func (api *OrchestrationAPI) StartProcess(c *gin.Context) {
	var req StartProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	found := false
	dataStore.Range(func(key, value interface{}) bool {
		if value.(UserProcess).UserID == req.UserID {
			c.JSON(http.StatusOK, value)
			found = true
			return false
		}
		return true
	})

	if found {
		return
	}
	requestID := uuid.New().String()
	userProcess := UserProcess{UserID: req.UserID, RequestID: requestID, Stage: Questions}
	dataStore.Store(requestID, userProcess)
	api.redisClient.Publish(c.Request.Context(), "questions_channel", requestID)
	c.JSON(http.StatusOK, userProcess)
}

func (api *OrchestrationAPI) SSEvents(c *gin.Context) {
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming unsupported"})
		return
	}

	requestID := c.Query("requestId")
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	var lastSentStage Stage
	for {
		if val, ok := dataStore.Load(requestID); ok {
			userProcess := val.(UserProcess)
			if userProcess.Stage != lastSentStage {
				lastSentStage = userProcess.Stage
				for _, module := range FEModulesDef {
					if module.Stage == userProcess.Stage {
						jsonData, _ := json.Marshal(module)
						dataString := "data: { \"module\": " + string(jsonData) + " }"
						c.Writer.Write([]byte(dataString + "\n\n"))
						flusher.Flush()
						break
					}
				}
			}
			select {
			case <-time.After(30 * time.Second):
				c.Writer.Write([]byte(": keep-alive\n\n"))
				flusher.Flush()
			}
		}
	}
}

func (api *OrchestrationAPI) GetAllDBData(c *gin.Context) {
	var allData []UserProcess
	dataStore.Range(func(_, value interface{}) bool {
		allData = append(allData, value.(UserProcess))
		return true
	})
	c.JSON(http.StatusOK, allData)
}

func (api *OrchestrationAPI) SubscribeToChannels() {
	ctx := context.Background()
	sub := api.redisClient.Subscribe(ctx, "questions_channel", "genres_channel", "movies_channel")
	ch := sub.Channel()

	for msg := range ch {
		api.AdvanceStage(msg.Payload)
	}
}

func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "redis:6379"})
	api := &OrchestrationAPI{redisClient: rdb}

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

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Orchestration API"})
	})
	r.POST("/start", api.StartProcess)
	r.GET("/events", api.SSEvents)
	r.GET("/backdoor/db", api.GetAllDBData)

	r.Run(":8080")
}
