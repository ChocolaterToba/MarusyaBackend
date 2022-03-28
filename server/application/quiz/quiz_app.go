package quiz

import (
	authApp "cmkids/application/auth"
	authModels "cmkids/models/auth"
	"cmkids/models/marusia"
	quizModels "cmkids/models/quiz"
	quizRepo "cmkids/repository/quiz"
	"sort"
	"strings"

	"go.uber.org/zap"
)

type QuizAppInterface interface {
	ProcessBasicRequest(input marusia.RequestBody) (response marusia.Response, err error)
}

type QuizApp struct {
	authApp  authApp.AuthAppInterface
	quizRepo quizRepo.QuizRepoInterface
	logger   *zap.Logger
}

func NewQuizApp(authApp authApp.AuthAppInterface, quizRepo quizRepo.QuizRepoInterface, logger *zap.Logger) *QuizApp {
	return &QuizApp{
		authApp:  authApp,
		quizRepo: quizRepo,
		logger:   logger}
}

func (app *QuizApp) ProcessBasicRequest(input marusia.RequestBody) (response marusia.Response, err error) {
	if input.Session.New {
		return app.authApp.Login(input)
		// TODO: add test continuation prompt
	}

	userID, err := app.authApp.GetUserIDBySessionID(input.Session.SessionID)
	if err != nil {
		if err == authModels.ErrUserNotFound && input.Session.MessageID == 1 {
			return app.authApp.Register(input)
		}
		return marusia.Response{}, err
	}

	currentQuestionID, err := app.quizRepo.GetCurrentQuestionID(userID)
	if err != nil {
		if err == quizModels.ErrCurrentQuestionNotFound {
			return app.GetTestsRoot(input)
		}
		return marusia.Response{}, err
	}

	currentQuestion, err := app.quizRepo.GetQuestion(currentQuestionID)
	if err != nil {
		return marusia.Response{}, err
	}

	//TODO: finishing a test

	nextQuestionID, err := getNextQuestionID(input.Request.OriginalUtterance, currentQuestion.NextQuestionIDs)
	if err != nil {
		if err != quizModels.ErrNextQuestionNotFound {
			return marusia.Response{}, err
		}

		return marusia.Response{
			Text:       quizModels.MsgIncorrectInput,
			Buttons:    marusia.ToButtons(getKeys(currentQuestion.NextQuestionIDs)),
			EndSession: false,
		}, nil
	}

	nextQuestion, err := app.quizRepo.GetQuestion(nextQuestionID)
	if err != nil {
		return marusia.Response{}, err
	}

	return marusia.Response{
		Text:       nextQuestion.Text,
		Buttons:    marusia.ToButtons(getKeys(nextQuestion.NextQuestionIDs)),
		EndSession: false,
	}, nil
}

func (app *QuizApp) GetTestsRoot(input marusia.RequestBody) (response marusia.Response, err error) {
	root, err := app.quizRepo.GetQuestion(quizModels.QuizRootID)
	if err != nil {
		return marusia.Response{}, err
	}

	return marusia.Response{
		Text:       root.Text,
		Buttons:    marusia.ToButtons(getKeys(root.NextQuestionIDs)),
		EndSession: false,
	}, nil
}

func getNextQuestionID(userInput string, nextQuestions map[string]uint64) (nextQuestionID uint64, err error) {
	// TODO: ML goes here
	for key := range nextQuestions {
		if strings.ToLower(key) == strings.ToLower(userInput) {
			return nextQuestions[key], nil
		}
	}
	return 0, quizModels.ErrNextQuestionNotFound
}

func getKeys(input map[string]uint64) (keys []string) {
	keys = make([]string, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
