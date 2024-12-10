package main

import (
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"time"
)

type QuestionsAPI struct {
	redisClient *redis.Client
}

type QuestionsAnswersRequest struct {
	RequestID string             `json:"requestId"`
	Questions []AnsweredQuestion `json:"questions"`
}

func (api *QuestionsAPI) QuestionsList(c *gin.Context) {
	c.JSON(200, QuestionsDB)
}

func (api *QuestionsAPI) Answers(c *gin.Context) {
	var answersRequest QuestionsAnswersRequest
	if err := c.ShouldBindJSON(&answersRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(answersRequest.Questions) != len(QuestionsDB) || !checkAnswers(answersRequest.Questions) {
		c.JSON(400, gin.H{"error": "Invalid answers"})
		return
	}

	jsonData, _ := json.Marshal(answersRequest.Questions)
	api.redisClient.Publish(c.Request.Context(), "questions_channel", jsonData)
	c.JSON(200, gin.H{"message": "Answers processed"})
}

func checkAnswers(answers []AnsweredQuestion) bool {
	if len(answers) != len(QuestionsDB) {
		return false
	}
	for _, question := range QuestionsDB {
		for _, answer := range answers {
			if question.Question == answer.Question {
				if question.Type != answer.Type {
					return false
				}
			}
		}
	}
	for _, answer := range answers {
		if answer.Type == TextQuestion {
			if strVal, ok := answer.Answer.(string); !ok || strVal == "" {
				return false
			}
		}
		if answer.Type == BooleanQuestion {
			if _, ok := answer.Answer.(bool); !ok {
				return false
			}
		}
	}
	return true
}

func main() {
	rdb := redis.NewClient(&redis.Options{Addr: "redis:6379"})
	api := &QuestionsAPI{redisClient: rdb}
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/questions", api.QuestionsList)
	r.POST("/answers", api.Answers)

	r.Run(":8081")
}
