package quiz

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

	quizMetadataLine := input[1]
	test.Title = quizMetadataLine[0]
	test.BackTrackingEnabled, err = parseBool(quizMetadataLine[1])
	if err != nil {
		return err
	}
	test.CalculateCorrectness, err = parseBool(quizMetadataLine[2])
	if err != nil {
		return err
	}

	test.Questions = make(map[uint64]Question)
	for i := 4; i < len(input); i++ {
		if len(input[i]) == 0 || input[i][0] == "" {
			break // we've exausted available rows
		}

		question := Question{}
		question.QuestionInTestID, err = strconv.ParseUint(input[i][0], 10, 64)
		if err != nil {
			return err
		}
		question.Text = input[i][1]

		question.Answers = make(map[string]Answer)
		for j := 0; j < len(input[i])/4; j++ {
			if input[i][j*4+2] == "" {
				break // we've exausted available columns
			}

			answer := Answer{}
			answer.AnswerText = input[i][j*4+3]
			answer.NextQuestionID, err = strconv.ParseUint(input[i][j*4+4], 10, 64)
			if err != nil {
				return err
			}
			answer.IsCorrect, err = parseBool(input[i][j*4+5])
			if err != nil {
				return err
			}
			question.Answers[input[i][j*4+2]] = answer
		}

		test.Questions[question.QuestionInTestID] = question
	}

	return nil
}

func parseBool(input string) (result bool, err error) {
	switch strings.ToUpper(input) {
	case "ИСТИНА", "ДА", "TRUE":
		return true, nil
	case "ЛОЖЬ", "НЕТ", "FALSE":
		return false, nil
	default:
		return false, errors.New(fmt.Sprintf("Could not parse '%s' as bool", input))
	}
}
