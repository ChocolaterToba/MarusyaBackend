package quiz

const (
	MsgIncorrectInput       = "Не смогла распознать твой ответ. Попробуй ещё раз."
	MsgSelectQuiz           = "Выбери, какой тест хочешь пройти."
	MsgQuestionRepeat       = "Повторяю"
	MsgStartOverTest        = "Хорошо, начнём тест сначала."
	MsgBackToQuestionInTest = "Хорошо, вернёмся на %s вопрос."

	ParticleNot = "не "
)

const (
	QuizGetHelp             = 1000000
	QuizFirstQuestion       = 1000001
	QuizQuitGame            = 1000002
	QuizRepeatLastMessage   = 1000003
	QuizReturnByOneQuestion = 1000004
	QuizRootID              = 0
)

var (
	Alphabet = []string{
		`{А}{"А"}`,
		`{Б}{"Бэ"}`,
		`{В}{"Вэ"}`,
		`{Г}{"Гэ"}`,
		`{Д}{"Дэ"}`,
	}

	QuestionPosition = map[uint64]string{
		1:  "первый",
		2:  "второй",
		3:  "третий",
		4:  "четвертый",
		5:  "пятый",
		6:  "шестой",
		7:  "седьмой",
		8:  "восьмой",
		9:  "девятый",
		10: "десятый",
		11: "одиннадцатый",
		12: "двенадцатый",
		13: "тринадцатый",
		14: "четырнадцатый",
		15: "пятнадцатый",
		16: "шестнадцатый",
		17: "семнадцатый",
		18: "восемнадцатый",
		19: "девятнадцатый",
		20: "двадцатый",
	}
)
