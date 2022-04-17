package quiz

const (
	MsgIncorrectInput = "Не смогла распознать твой ответ. Попробуй ещё раз."
	MsgSelectQuiz     = "Выбери, какой тест хочешь пройти."
	MsgQuestionRepeat = "Повторяю"
	MsgStartOverTest  = "Хорошо, начнём тест сначала."

	ParticleNot = "не"
)

const (
	QuizGetHelp       = 1000000
	QuizFirstQuestion = 1000001
	QuizQuitGame      = 1000002
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
