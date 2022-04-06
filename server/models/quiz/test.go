package quiz

type Test struct {
	TestID      uint64
	Name        string
	Description string
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
}
