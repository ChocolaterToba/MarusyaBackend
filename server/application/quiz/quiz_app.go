package quiz

import (
	authApp "cmkids/application/auth"
	authModels "cmkids/models/auth"
	"cmkids/models/help"
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
	if len(currentQuestion.Answers) == 0 { // No next questions => this question is the last one, go to root
		return app.navToQuestionByID(userID, quizModels.QuizRootID, append(response.Text, currentQuestion.Text), false)
	}

	answer, isTypicalNavigation, err := getFittingAnswer(input.Request.OriginalUtterance, currentQuestion)
	if err != nil {
		if err != quizModels.ErrNextQuestionNotFound {
			return marusia.Response{}, err
		}

		return app.navToQuestion(userID, currentQuestion, append(response.Text, quizModels.MsgIncorrectInput), true)
	}

	if !isTypicalNavigation {
		return app.processAbsoluteQuestionID(userID, currentQuestion, response.Text, answer.NextQuestionID)
	}

	// When we are in root, nextQuestionID is question_id in db
	if currentQuestionID == quizModels.QuizRootID {
		return app.navToQuestionByID(userID, answer.NextQuestionID, response.Text, false)
	}

	// When we are not in root, nextQuestionID is internal test id or root's id
	if answer.NextQuestionID == quizModels.QuizRootID { // root in not in any test and is handled separately
		return app.navToQuestionByID(userID, quizModels.QuizRootID, response.Text, false)
	}

	nextQuestion, err := app.quizRepo.GetQuestionInTest(currentQuestion.TestID, answer.NextQuestionID)
	if err != nil {
		return marusia.Response{}, err
	}

	if answer.NextQuestionID < currentQuestion.QuestionInTestID {
		previousQuestionID := answer.NextQuestionID
		if currentQuestion.QuestionInTestID - previousQuestionID == 1 {
			previousQuestionID = 0
		}
		response.Text = append(response.Text, fmt.Sprintf(quizModels.MsgBackToQuestionInTest, quizModels.QuestionPosition[previousQuestionID]))
	}
	if answer.AnswerText != "" {
		response.Text = append(response.Text, answer.AnswerText)
	}

	if len(nextQuestion.Answers) == 0 { // No next questions => this question is the last one, go to root
		return app.navToQuestionByID(userID, quizModels.QuizRootID, append(response.Text, nextQuestion.Text), false)
	}

	return app.navToQuestion(userID, nextQuestion, response.Text, false)
}

func (app *QuizApp) processAbsoluteQuestionID(userID uint64, currentQuestion quizModels.Question, prevText []string, nextQuestionID uint64) (response marusia.Response, err error) {
	switch nextQuestionID {
	case quizModels.QuizRepeatLastMessage:
		response.Text = append(response.Text, quizModels.MsgQuestionRepeat)
		return app.navToQuestion(userID, currentQuestion, response.Text, true)

	case quizModels.QuizFirstQuestion:
		firstQuestion, err := app.quizRepo.GetQuestionInTest(currentQuestion.TestID, 1)
		if err != nil {
			return marusia.Response{}, err
		}
		response.Text = append(response.Text, quizModels.MsgStartOverTest)
		return app.navToQuestion(userID, firstQuestion, response.Text, false)

	case quizModels.QuizGetHelp:
		response.Text = append(response.Text, help.MsgHelpMe)
		return app.navToQuestion(userID, currentQuestion, response.Text, false)

	case quizModels.QuizRootID:
		return app.navToQuestionByID(userID, quizModels.QuizRootID, response.Text, false)

	case quizModels.QuizQuitGame:
		// TODO: add logout here?
		return marusia.Response{
			Text:       []string{authModels.MsgGoodBye},
			EndSession: true,
		}, nil

	default:
		return marusia.Response{
			Text:       []string{quizModels.ErrNextQuestionNotFound.Error()},
			EndSession: false,
		}, nil
	}
}

func (app *QuizApp) navToQuestionByID(userID uint64, questionID uint64, prevText []string, isLoop bool) (response marusia.Response, err error) {
	question, err := app.quizRepo.GetQuestion(questionID)
	if err != nil {
		return marusia.Response{}, err
	}

	return app.navToQuestion(userID, question, prevText, isLoop)
}

func (app *QuizApp) navToQuestion(userID uint64, question quizModels.Question, prevText []string, isLoop bool) (response marusia.Response, err error) {
	if !isLoop {
		err = app.quizRepo.SetCurrentQuestionID(userID, question.QuestionID)
		if err != nil {
			return marusia.Response{}, err
		}
	}

	choices := getKeysFromAnswers(question.Answers)
	return marusia.Response{
		Text:       appendChoices(append(prevText, question.Text), choices),
		Buttons:    marusia.ToButtons(choices),
		EndSession: false,
	}, nil
}

func getFittingAnswer(userInput string, question quizModels.Question) (nextAnswer quizModels.Answer, isTypicalNavigation bool, err error) {
	userInput = strings.ToLower(userInput)
	userInput = strings.TrimRight(userInput, ".?!")

	// Searching for answers from db
	lastMatch, found := getLastMatch(userInput, question.Answers)
	if found {
		return lastMatch, true, nil
	}

	// searching for "repeat" and similar commands
	for _, answerRepeat := range quizModels.AnswersRepeat {
		if strings.Contains(userInput, answerRepeat) {
			return quizModels.Answer{NextQuestionID: quizModels.QuizRepeatLastMessage}, false, nil
		}
	}

	// searching for "start test again" and similar commands
	for _, answerReturnToFirstQuestion := range quizModels.AnswersReturnToFirstQuestion {
		if strings.Contains(userInput, answerReturnToFirstQuestion) {
			return quizModels.Answer{NextQuestionID: quizModels.QuizFirstQuestion}, false, nil
		}
	}

	// searching for "end test" and similar commands
	for _, answerReturnToRoot := range quizModels.AnswersReturnToRoot {
		if strings.Contains(userInput, answerReturnToRoot) {
			return quizModels.Answer{NextQuestionID: quizModels.QuizRootID}, false, nil
		}
	}

	// return to n questions back
	if question.QuestionID != quizModels.QuizRootID {
		for _, BackToQuestion := range quizModels.AnswersBackToQuestion {
			if strings.Contains(userInput, BackToQuestion) {
				for word, pos := range quizModels.AnswersIntTestPositional {
					if strings.Contains(userInput, word) {
						questionInTest := int(question.QuestionInTestID) - pos
						if questionInTest < 1 {
							questionInTest = 1
						}
						return quizModels.Answer{NextQuestionID: uint64(questionInTest)}, true, nil
					}
				}
			}
		}
	}

	for _, answerQuitGame := range quizModels.AnswersQuitGame {
		if strings.Contains(userInput, answerQuitGame) {
			return quizModels.Answer{NextQuestionID: quizModels.QuizQuitGame}, false, nil
		}
	}

	for _, helpQuestion := range help.CallHelp {
		if strings.Contains(userInput, helpQuestion) {
			return quizModels.Answer{NextQuestionID: quizModels.QuizGetHelp}, false, nil
		}
	}

	userInputTokens := strings.Fields(userInput)

	for i := len(userInputTokens) - 1; i >= 0; i-- {
		pos, exists := quizModels.AnswersPositional[userInputTokens[i]]
		if exists {
			if pos >= len(question.Answers) {
				return quizModels.Answer{}, false, quizModels.ErrNextQuestionNotFound
			}

			// if pos is valid, find corresponding answer
			return question.Answers[getKeysFromAnswers(question.Answers)[pos]], true, nil
		}
	}

	return quizModels.Answer{}, false, quizModels.ErrNextQuestionNotFound
}

func getLastMatch(userInput string, matches map[string]quizModels.Answer) (resultAnswer quizModels.Answer, found bool) {
	lastMatch := ""
	lastMatchIndex := -1
	for key := range matches {
		newMatchIndex := strings.LastIndex(userInput, strings.TrimRight(strings.ToLower(key), ".?!"))
		if newMatchIndex > lastMatchIndex {
			lastMatch = key
			lastMatchIndex = newMatchIndex
		}

		if strings.Contains(userInput, quizModels.ParticleNot) && !strings.Contains(key, quizModels.ParticleNot) {
			lastMatchIndex = -1
			newMatchIndex = strings.Index(userInput, strings.TrimRight(strings.ToLower(key), ".?!"))
			if newMatchIndex > lastMatchIndex {
				lastMatch = key
				lastMatchIndex = newMatchIndex
			}
		}
	}

	if lastMatch != "" {
		return matches[lastMatch], true
	}

	return quizModels.Answer{}, false
}

func getKeysFromAnswers(input map[string]quizModels.Answer) (keys []string) {
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
			choice = fmt.Sprintf("%s: %s", quizModels.Alphabet[i], choice)
		}

		result = append(result, choice)
	}

	return result
}
