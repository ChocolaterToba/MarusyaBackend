package quiz

const (
	MsgIncorrectInput = "Не смогла распознать твой ответ. Попробуй ещё раз."
	MsgSelectQuiz     = "Выбери, какой тест хочешь пройти."
	MsgQuestionRepeat = "Повторяю"
	MsgStartOverTest  = "Хорошо, начнём тест сначала."
)

const (
	QuizGetHelp       = -3
	QuizFirstQuestion = -2
	QuizQuitGame      = -1
	QuizRootID        = 0
)

var (
	Alphabet = []string{
		`{А}{"А"}`,
		`{Б}{"Бэ"}`,
		`{В}{"Вэ"}`,
		`{Г}{"Гэ"}`,
		`{Д}{"Дэ"}`,
	}
)
