package quiz

import "errors"

var (
	ErrQuestionNotFound     = errors.New("Не удалось найти вопрос")
	ErrChooseQuizFirst      = errors.New("Сначала выбери тест")
	ErrBacktrackingDisabled = errors.New("Возврат назад отключен")
)

var (
	ErrNoFile                = errors.New("No files found in body")
	ErrFileTooLarge          = errors.New("Body is too large")
	ErrIncorrectQuizAmount   = errors.New("Quiz amunt must be positive")
	ErrUnsupportedFileFormat = errors.New("File format is not supported")
)
