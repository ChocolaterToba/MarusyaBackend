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
		var finished bool // to avoid variable shadowing later
		response, finished, err = app.authApp.Login(input)
		if !finished {
			return response, err
		}

		userID, err := app.authApp.GetUserIDBySessionID(input.Session.SessionID)
		if err != nil {
			return marusia.Response{}, err
		}

		currentQuestionID, err := app.quizRepo.GetCurrentQuestionID(userID)
		if err != nil {
			return marusia.Response{}, err
		}

		return app.navToQuestion(userID, currentQuestionID, response.Text)
	}

	userID, err := app.authApp.GetUserIDBySessionID(input.Session.SessionID)
	if err != nil {
		if err == authModels.ErrUserNotFound && input.Session.MessageID == 1 { // TODO: better registration input detection, beyond messageID
			var finished bool // to avoid variable shadowing later
			response, finished, err = app.authApp.Register(input)
			if !finished {
				return response, err
			}

			userID, err := app.authApp.GetUserIDBySessionID(input.Session.SessionID)
			if err != nil {
				return marusia.Response{}, err
			}

			return app.navToQuestion(userID, quizModels.QuizRootID, response.Text)
		}
		return marusia.Response{}, err
	}

	currentQuestionID, err := app.quizRepo.GetCurrentQuestionID(userID)
	if err != nil {
		return marusia.Response{}, err
	}

	currentQuestion, err := app.quizRepo.GetQuestion(currentQuestionID)
	if err != nil {
		return marusia.Response{}, err
	}

	// this technically is not supposed to happen, just in case
	if len(currentQuestion.NextQuestionIDs) == 0 { // No next questions => this question is the last one, go to root
		return app.navToQuestion(userID, quizModels.QuizRootID, append(response.Text, currentQuestion.Text))
	}

	nextQuestionID, err := getNextQuestionID(input.Request.OriginalUtterance, currentQuestion)
	if err != nil {
		if err != quizModels.ErrNextQuestionNotFound {
			return marusia.Response{}, err
		}

		return marusia.Response{
			Text:       []string{quizModels.MsgIncorrectInput, currentQuestion.Text},
			Buttons:    marusia.ToButtons(getKeys(currentQuestion.NextQuestionIDs)),
			EndSession: false,
		}, nil
	}

	var nextQuestion quizModels.Question
	switch currentQuestionID {
	case nextQuestionID: // If our destination is current question, we repeat it
		nextQuestion = currentQuestion
		response.Text = append(response.Text, quizModels.MsgQuestionRepeat)

	case quizModels.QuizRootID: // When we are in root, nextQuestionID is question_id in db
		return app.navToQuestion(userID, nextQuestionID, response.Text)

	default: // When we are not in root, nextQuestionID is internal test id or root's id
		if nextQuestionID == quizModels.QuizRootID { // root in not in any test and is handled separately
			return app.navToQuestion(userID, quizModels.QuizRootID, response.Text)
		}

		nextQuestion, err = app.quizRepo.GetQuestionInTest(currentQuestion.TestID, nextQuestionID)
		if err != nil {
			return marusia.Response{}, err
		}
	}

	if len(nextQuestion.NextQuestionIDs) == 0 { // No next questions => this question is the last one, go to root
		return app.navToQuestion(userID, quizModels.QuizRootID, append(response.Text, nextQuestion.Text))
	}

	if nextQuestion.QuestionID != currentQuestion.QuestionID {
		err = app.quizRepo.SetCurrentQuestionID(userID, nextQuestion.QuestionID)
		if err != nil {
			return marusia.Response{}, err
		}
	}

	return marusia.Response{
		Text:       append(response.Text, nextQuestion.Text),
		Buttons:    marusia.ToButtons(getKeys(nextQuestion.NextQuestionIDs)),
		EndSession: false,
	}, nil
}

func (app *QuizApp) navToQuestion(userID uint64, questionID uint64, prevText []string) (response marusia.Response, err error) {
	question, err := app.quizRepo.GetQuestion(questionID)
	if err != nil {
		return marusia.Response{}, err
	}

	err = app.quizRepo.SetCurrentQuestionID(userID, questionID)
	if err != nil {
		return marusia.Response{}, err
	}

	return marusia.Response{
		Text:       append(prevText, question.Text),
		Buttons:    marusia.ToButtons(getKeys(question.NextQuestionIDs)),
		EndSession: false,
	}, nil
}

func getNextQuestionID(userInput string, question quizModels.Question) (nextQuestionID uint64, err error) {
	userInput = strings.ToLower(userInput)

	for key := range question.NextQuestionIDs {
		if strings.ToLower(key) == userInput {
			return question.NextQuestionIDs[key], nil
		}
	}

	for _, repeatAnswer := range quizModels.AnswersRepeat {
		if strings.Contains(userInput, repeatAnswer) {
			return question.QuestionID, nil
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
