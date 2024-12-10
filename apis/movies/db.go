package main

import "sync"

type Movie struct {
	Title  string  `json:"title"`
	Genre  string  `json:"genre"`
	Rating float32 `json:"rating"`
}

var RequestIDMap = sync.Map{}

var MoviesDB = []Movie{
	Movie{Title: "The Shawshank Redemption", Genre: "Drama", Rating: 9.3},
	Movie{Title: "The Godfather", Genre: "Crime", Rating: 9.2},
	Movie{Title: "The Dark Knight", Genre: "Action", Rating: 9.0},
	Movie{Title: "The Lord of the Rings: The Return of the King", Genre: "Adventure", Rating: 8.9},
	Movie{Title: "Pulp Fiction", Genre: "Crime", Rating: 8.9},
	Movie{Title: "Schindler's List", Genre: "Biography", Rating: 8.9},
	Movie{Title: "Forrest Gump", Genre: "Drama", Rating: 8.8},
	Movie{Title: "Inception", Genre: "Action", Rating: 8.8},
}
