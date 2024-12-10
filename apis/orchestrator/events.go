package main

import (
	"context"
	"encoding/json"
	"errors"
)

func checkEventRequestID(requestID string) bool {
	_, ok := dataStore.Load(requestID)
	return ok
}

type AnsweredQuestion struct {
	Question string      `json:"question"`
	Type     string      `json:"string"`
	Answer   interface{} `json:"answer"`
}

type GenericEvent struct {
	RequestID string `json:"requestId"`
}
type QuestionsAnswersEvent struct {
	*GenericEvent
	Questions []AnsweredQuestion `json:"questions"`
}

type GenresPreferencesEvent struct {
	*GenericEvent
	SelectedGenres []string `json:"selectedGenres"`
}

type MoviesPreferencesEvent struct {
	*GenericEvent
	SelectedMovies []Movie `json:"selectedMovies"`
}

func (api *OrchestrationAPI) AdvanceStage(payload string) error {
	var event GenericEvent
	var userProcess UserProcess
	json.Unmarshal([]byte(payload), &event)
	if checkEventRequestID(event.RequestID) {
		val, _ := dataStore.Load(event.RequestID)
		userProcess = val.(UserProcess)
	} else {
		return errors.New("Request ID not found")
	}

	switch userProcess.Stage {
	case Questions:
		var questionEvent QuestionsAnswersEvent
		json.Unmarshal([]byte(payload), &questionEvent)
		userProcess.Stage = Genres
		userProcess.Questions = questionEvent.Questions
		dataStore.Store(event.RequestID, userProcess)
	case Genres:
		var genresEvent GenresPreferencesEvent
		json.Unmarshal([]byte(payload), &genresEvent)
		userProcess.Stage = Movies
		userProcess.SelectedGenres = genresEvent.SelectedGenres
		api.redisClient.Publish(context.Background(), "movies_genres_channel", payload)
	case Movies:
		var moviesEvent MoviesPreferencesEvent
		json.Unmarshal([]byte(payload), &moviesEvent)
		userProcess.SelectedMovies = moviesEvent.SelectedMovies
		userProcess.Stage = End
	}

	dataStore.Store(event.RequestID, userProcess)
	return nil
}
