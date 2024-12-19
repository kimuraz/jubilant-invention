package main

import "sync"

type Stage string

const (
	Questions Stage = "questions"
	Genres    Stage = "genres"
	Movies    Stage = "movies"
	End       Stage = "end"
)

var dataStore = sync.Map{}

type Movie struct {
	Title  string  `json:"title"`
	Genre  string  `json:"genre"`
	Rating float32 `json:"rating"`
}

type Question struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type UserProcess struct {
	UserID         string             `json:"userId"`
	RequestID      string             `json:"requestId"`
	Stage          Stage              `json:"stage"`
	Questions      []AnsweredQuestion `json:"questions"`
	SelectedGenres []string           `json:"selectedGenres"`
	SelectedMovies []Movie            `json:"selectedMovies"`
}

type FEModule struct {
	ModuleID  string `json:"moduleID"`
	ModuleURL string `json:"moduleURL"`
	ModuleTag string `json:"moduleTag"`
	Stage     Stage  `json:"stage"`
}

var FEModulesDef = []FEModule{
	{ModuleID: "questions", ModuleURL: "http://localhost:9000/fe/questions.js", ModuleTag: "questions-module", Stage: Questions},
	{ModuleID: "genres", ModuleURL: "http://localhost:9000/fe/genres.js", ModuleTag: "genres-module", Stage: Genres},
	{ModuleID: "movies", ModuleURL: "http://localhost:9000/fe/movies.js", ModuleTag: "movies-module", Stage: Movies},
	{ModuleID: "end", ModuleURL: "http://localhost:9000/fe/thanks.js", ModuleTag: "thanks-module", Stage: End},
}
