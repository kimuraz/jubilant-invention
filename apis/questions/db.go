package main

type QuestionType string

const (
	TextQuestion    QuestionType = "text"
	BooleanQuestion QuestionType = "boolean"
)

type MovieQuestions struct {
	Question string       `json:"question"`
	Type     QuestionType `json:"type"`
}

type AnsweredQuestion struct {
	*MovieQuestions
	Answer interface{} `json:"answer"`
}

var QuestionsDB = []MovieQuestions{
	{Question: "What is your favorite movie?", Type: TextQuestion},
	{Question: "Can you watch 3 hours long movies?", Type: BooleanQuestion},
	{Question: "Are you easily scared by horror movies?", Type: BooleanQuestion},
}
