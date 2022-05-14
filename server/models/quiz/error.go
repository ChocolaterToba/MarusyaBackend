package quiz

import "errors"

var (
	ErrCurrentQuestionNotFound = errors.New("Не удалось найти вопрос")
	ErrNextQuestionNotFound    = errors.New("Не удалось найти подходящий вариант ответа")
	ErrChooseQuizFirst         = errors.New("Сначала выбери тест")
)
