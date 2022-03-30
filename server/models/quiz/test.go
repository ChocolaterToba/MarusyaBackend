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
	NextQuestionIDs  map[string]uint64
}
