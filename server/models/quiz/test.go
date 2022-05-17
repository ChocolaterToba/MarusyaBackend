package quiz

import (
	"errors"
	"fmt"
)

type Test struct {
	TestID               uint64
	Title                string
	BackTrackingEnabled  bool
	CalculateCorrectness bool
	Questions            map[uint64]Question // uint64 is QuestionInTestID
}

type Question struct {
	QuestionID       uint64
	QuestionInTestID uint64
	TestID           uint64
	Text             string
	Answers          map[string]Answer
}

type Answer struct {
	NextQuestionID uint64 `json:"next_question_id"`
	AnswerText     string `json:"answer_text"`
	IsCorrect      bool   `json:"is_correct"`
}

func (test *Test) Parse(input [][]string) (err error) {
	if len(input) < 5 {
		return errors.New("Not enough rows")
	}
	fmt.Println(input)
	return errors.New("Not implemented yet")
}
