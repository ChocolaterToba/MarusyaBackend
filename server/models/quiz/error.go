package quiz

import "errors"

var (
	ErrCurrentQuestionNotFound = errors.New("Не удалось найти вопрос")
	ErrNextQuestionNotFound    = errors.New("Не удалось найти подходящий вариант ответа")
	ErrChooseQuizFirst         = errors.New("Сначала выбери тест")
)

var (
	ErrNoFile                = errors.New("No files found in body")
	ErrFileTooLarge          = errors.New("Body is too large")
	ErrIncorrectQuizAmount   = errors.New("Quiz amunt must be positive")
	ErrUnsupportedFileFormat = errors.New("File format is not supported")
)
