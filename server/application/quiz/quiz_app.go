package quiz

import (
	authApp "cmkids/application/auth"
	authModels "cmkids/models/auth"
	"cmkids/models/marusia"
	quizModels "cmkids/models/quiz"
	quizRepo "cmkids/repository/quiz"
	"fmt"
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

		return app.navToQuestionByID(userID, currentQuestionID, response.Text, false)
	}

	userID, err := app.authApp.GetUserIDBySessionID(input.Session.SessionID)
	if err != nil {
		if err == authModels.ErrUserNotFound && input.Session.MessageID == 1 { // TODO: better registration input detection, beyond messageID
			var finished bool // to avoid variable shadowing later
			response, finished, err = app.authApp.Register(input)
			if !finished {
				return response, err
			}

			userID, err = app.authApp.GetUserIDBySessionID(input.Session.SessionID)
			if err != nil {
				return marusia.Response{}, err
			}

			return app.navToQuestionByID(userID, quizModels.QuizRootID, response.Text, false)
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
		return app.navToQuestionByID(userID, quizModels.QuizRootID, append(response.Text, currentQuestion.Text), false)
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

	switch currentQuestionID {
	case nextQuestionID: // If our destination is current question, we repeat it
		response.Text = append(response.Text, quizModels.MsgQuestionRepeat)
		return app.navToQuestion(userID, currentQuestion, response.Text, true)

	case quizModels.QuizRootID: // When we are in root, nextQuestionID is question_id in db
		return app.navToQuestionByID(userID, nextQuestionID, response.Text, false)

	}

	// When we are not in root, nextQuestionID is internal test id or root's id
	if nextQuestionID == quizModels.QuizRootID { // root in not in any test and is handled separately
		return app.navToQuestionByID(userID, quizModels.QuizRootID, response.Text, false)
	}

	nextQuestion, err := app.quizRepo.GetQuestionInTest(currentQuestion.TestID, nextQuestionID)
	if err != nil {
		return marusia.Response{}, err
	}

	if len(nextQuestion.NextQuestionIDs) == 0 { // No next questions => this question is the last one, go to root
		return app.navToQuestionByID(userID, quizModels.QuizRootID, append(response.Text, nextQuestion.Text), false)
	}

	return app.navToQuestion(userID, nextQuestion, append(response.Text, nextQuestion.Text), false)
}

func (app *QuizApp) navToQuestionByID(userID uint64, questionID uint64, prevText []string, isLoop bool) (response marusia.Response, err error) {
	question, err := app.quizRepo.GetQuestion(questionID)
	if err != nil {
		return marusia.Response{}, err
	}

	return app.navToQuestion(userID, question, prevText, isLoop)
}

func (app *QuizApp) navToQuestion(userID uint64, question quizModels.Question, prevText []string, isLoop bool) (response marusia.Response, err error) {
	if isLoop {
		err = app.quizRepo.SetCurrentQuestionID(userID, question.QuestionID)
		if err != nil {
			return marusia.Response{}, err
		}
	}

	choices := getKeys(question.NextQuestionIDs)
	return marusia.Response{
		Text:       appendChoices(append(prevText, question.Text), choices),
		Buttons:    marusia.ToButtons(choices),
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

	for _, answerRepeat := range quizModels.AnswersRepeat {
		if strings.Contains(userInput, answerRepeat) {
			return question.QuestionID, nil
		}
	}

	// TODO: rework this somehow
	userInputPositional := userInput
	for _, infix := range quizModels.AnswersPositionalInfixes {
		if strings.HasPrefix(userInputPositional, infix+" ") {
			userInputPositional = strings.TrimPrefix(userInputPositional, infix+" ")
		}
		if strings.HasSuffix(userInputPositional, " "+infix) {
			userInputPositional = strings.TrimSuffix(userInputPositional, " "+infix)
		}
	}

	pos, exists := quizModels.AnswersPositional[userInputPositional]
	if exists {
		if pos >= len(question.NextQuestionIDs) {
			return 0, quizModels.ErrNextQuestionNotFound
		}

		// if pos is valid, find corresponding answer
		return question.NextQuestionIDs[getKeys(question.NextQuestionIDs)[pos]], nil
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

func appendChoices(text []string, choices []string) (result []string) {
	result = make([]string, 0, len(text)+len(choices))
	result = append(result, text...)

	for i, choice := range choices {
		if i < 5 {
			choice = fmt.Sprintf("%s: %s", string([]rune(quizModels.Alphabet)[i]), choice)
		}

		result = append(result, choice)
	}

	return result
}
