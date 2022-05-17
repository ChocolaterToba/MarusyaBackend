package quiz

import (
	authApp "cmkids/application/auth"
	"cmkids/models/marusia"
	"cmkids/models/settings"
	quizRepo "cmkids/repository/quiz"
	"io"

	"go.uber.org/zap"
)

type QuizAppInterface interface {
	ProcessMarusiaRequest(input marusia.RequestBody) (response marusia.Response, err error)
	AddQuizFromFile(filename string, file io.Reader) (err error)
}

type QuizApp struct {
	authApp  authApp.AuthAppInterface
	quizRepo quizRepo.QuizRepoInterface
	config   *settings.Config
	logger   *zap.Logger
}

func NewQuizApp(authApp authApp.AuthAppInterface, quizRepo quizRepo.QuizRepoInterface, config *settings.Config, logger *zap.Logger) *QuizApp {
	return &QuizApp{
		authApp:  authApp,
		quizRepo: quizRepo,
		config:   config,
		logger:   logger,
	}
}
