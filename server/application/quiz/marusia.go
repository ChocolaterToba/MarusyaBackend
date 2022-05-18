package quiz

import (
	authModels "cmkids/models/auth"
	"cmkids/models/help"
	"cmkids/models/marusia"
	quizModels "cmkids/models/quiz"
	"fmt"
	"sort"
	"strings"
)

func (app *QuizApp) ProcessMarusiaRequest(input marusia.RequestBody) (response marusia.Response, err error) {
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

		pastAnswers, err := app.quizRepo.GetPastAnswers(userID)
		if err != nil {
			return marusia.Response{}, err
		}

		currentQuestionID := uint64(quizModels.QuizRootID)
		if len(pastAnswers) != 0 {
			currentQuestionID = pastAnswers[len(pastAnswers)-1].NextQuestionID
		}

		return app.navToQuestionByID(userID, currentQuestionID, response.Text)
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

			return app.navToQuestionByID(userID, quizModels.QuizRootID, response.Text)
		}
		return marusia.Response{}, err
	}

	pastAnswers, err := app.quizRepo.GetPastAnswers(userID)
	if err != nil {
		return marusia.Response{}, err
	}

	currentQuestionID := uint64(quizModels.QuizRootID)
	if len(pastAnswers) != 0 {
		currentQuestionID = pastAnswers[len(pastAnswers)-1].NextQuestionID
	}

	currentQuestion, err := app.quizRepo.GetQuestion(currentQuestionID)
	if err != nil {
		return marusia.Response{}, err
	}

	// This technically is not supposed to happen, just in case
	if len(currentQuestion.Answers) == 0 { // No next questions => this question is the last one, go to root
		return app.finishQuiz(userID, pastAnswers, append(response.Text, currentQuestion.Text))
	}

	answer, isTypicalNavigation, err := getFittingAnswer(input.Request.OriginalUtterance, currentQuestion)
	if err != nil {
		if err != quizModels.ErrQuestionNotFound {
			return marusia.Response{}, err
		}

		return app.navToQuestion(userID, currentQuestion, append(response.Text, app.config.Messages.MsgIncorrectInput))
	}

	if !isTypicalNavigation {
		return app.processAbsoluteQuestionID(userID, pastAnswers, currentQuestion, response.Text, answer)
	}

	// When we are in root, nextQuestionID is question_id in db
	if currentQuestionID == quizModels.QuizRootID {
		err = app.quizRepo.SetPastAnswers(userID, []quizModels.Answer{answer})
		if err != nil {
			return marusia.Response{}, err
		}

		return app.navToQuestionByID(userID, answer.NextQuestionID, response.Text)
	}

	// When we are not in root, nextQuestionID is internal test id or root's id
	if answer.NextQuestionID == quizModels.QuizRootID { // root in not in any test and is handled separately
		return app.finishQuiz(userID, pastAnswers, append(response.Text, answer.AnswerText))
	}

	nextQuestion, err := app.quizRepo.GetQuestionInTest(currentQuestion.TestID, answer.NextQuestionID)
	if err != nil {
		return marusia.Response{}, err
	}

	// pastAnswers need to have absolute ids in them
	pastAnswers = append(pastAnswers, quizModels.Answer{NextQuestionID: nextQuestion.QuestionID, IsCorrect: answer.IsCorrect})
	if answer.AnswerText != "" {
		response.Text = append(response.Text, answer.AnswerText)
	}

	if len(nextQuestion.Answers) == 0 { // No next questions => this question is the last one, go to root
		return app.finishQuiz(userID, pastAnswers, append(response.Text, nextQuestion.Text))
	}

	err = app.quizRepo.SetPastAnswers(userID, pastAnswers)
	if err != nil {
		return marusia.Response{}, err
	}

	return app.navToQuestion(userID, nextQuestion, response.Text)
}

func (app *QuizApp) processAbsoluteQuestionID(userID uint64, pastAnswers []quizModels.Answer,
	currentQuestion quizModels.Question, prevText []string, nextAnswer quizModels.Answer) (response marusia.Response, err error) {
	switch nextAnswer.NextQuestionID {
	case quizModels.QuizRepeatLastMessage:
		response.Text = append(response.Text, app.config.Messages.MsgQuestionRepeat)
		return app.navToQuestion(userID, currentQuestion, response.Text)

	case quizModels.QuizFirstQuestion:
		// TODO: get quiz and check if backtracking is allowed
		if len(pastAnswers) == 0 {
			return marusia.Response{}, quizModels.ErrChooseQuizFirst
		}

		err = app.quizRepo.SetPastAnswers(userID, pastAnswers[:1]) // clear entire pastAnswers except for first element
		if err != nil {
			return marusia.Response{}, err
		}

		response.Text = append(response.Text, app.config.Messages.MsgStartQuizOver)
		return app.navToQuestionByID(userID, pastAnswers[0].NextQuestionID, response.Text)

	case quizModels.QuizGetHelp:
		response.Text = append(response.Text, app.config.Messages.MsgHelp)
		return app.navToQuestion(userID, currentQuestion, response.Text)

	case quizModels.QuizRootID:
		err = app.quizRepo.SetPastAnswers(userID, nil)
		if err != nil {
			return marusia.Response{}, err
		}
		return app.navToQuestionByID(userID, quizModels.QuizRootID, response.Text)

	case quizModels.QuizQuitGame:
		err = app.quizRepo.SetPastAnswers(userID, nil)
		if err != nil {
			return marusia.Response{}, err
		}

		// TODO: add logout here?
		return marusia.Response{
			Text:       []string{app.config.Messages.MsgGoodbye},
			EndSession: true,
		}, nil

	case quizModels.QuizReturnByOneQuestion:
		// TODO: get quiz and check if backtracking is allowed
		if len(pastAnswers) == 0 {
			return marusia.Response{}, quizModels.ErrChooseQuizFirst
		}

		err = app.quizRepo.SetPastAnswers(userID, pastAnswers[:len(pastAnswers)-1]) // remove last pastAnswers element
		if err != nil {
			return marusia.Response{}, err
		}

		return app.navToQuestionByID(userID, pastAnswers[len(pastAnswers)-2].NextQuestionID, response.Text)

	default:
		return app.navToQuestionByID(userID, nextAnswer.NextQuestionID, response.Text)
	}
}

func (app *QuizApp) navToQuestionByID(userID uint64, questionID uint64, prevText []string) (response marusia.Response, err error) {
	question, err := app.quizRepo.GetQuestion(questionID)
	if err != nil {
		return marusia.Response{}, err
	}

	return app.navToQuestion(userID, question, prevText)
}

func (app *QuizApp) navToQuestion(userID uint64, question quizModels.Question, prevText []string) (response marusia.Response, err error) {
	choices := getKeysFromAnswers(question.Answers)
	return marusia.Response{
		Text:       appendChoices(append(prevText, question.Text), choices),
		Buttons:    marusia.ToButtons(choices),
		EndSession: false,
	}, nil
}

// finishQuiz navigates user to root AND also collects user's statistics, if quiz so requires
func (app *QuizApp) finishQuiz(userID uint64, pastAnswers []quizModels.Answer, prevText []string) (response marusia.Response, err error) {
	correctAnswersCount := 0
	for _, answer := range pastAnswers {
		if answer.IsCorrect {
			correctAnswersCount++
		}
	}

	fmt.Printf("Верных ответов %d из %d\n", correctAnswersCount, len(pastAnswers)) // TODO: send this to CRM

	err = app.quizRepo.SetPastAnswers(userID, nil)
	return app.navToQuestionByID(userID, quizModels.QuizRootID, append(prevText, app.config.Messages.MsgFinishQuiz))
}

func getFittingAnswer(userInput string, question quizModels.Question) (nextAnswer quizModels.Answer, isTypicalNavigation bool, err error) {
	userInput = strings.ToLower(userInput)
	userInput = strings.TrimRight(userInput, ".?!")

	// Searching for answers from db
	lastMatch, found := getLastMatch(userInput, question.Answers)
	if found {
		return lastMatch, true, nil
	}

	// searching for "start test again" and similar commands
	for _, answerReturnToFirstQuestion := range quizModels.AnswersReturnToFirstQuestion {
		if strings.Contains(userInput, answerReturnToFirstQuestion) {
			return quizModels.Answer{NextQuestionID: quizModels.QuizFirstQuestion}, false, nil
		}
	}

	// searching for "repeat" and similar commands
	for _, answerRepeat := range quizModels.AnswersRepeat {
		if strings.Contains(userInput, answerRepeat) {
			return quizModels.Answer{NextQuestionID: quizModels.QuizRepeatLastMessage}, false, nil
		}
	}

	// searching for "end test" and similar commands
	for _, answerReturnToRoot := range quizModels.AnswersReturnToRoot {
		if strings.Contains(userInput, answerReturnToRoot) {
			return quizModels.Answer{NextQuestionID: quizModels.QuizRootID}, false, nil
		}
	}

	// searching for "previous question" and similar commands
	for _, answerBackToQuestion := range quizModels.AnswersBackToQuestion {
		if strings.Contains(userInput, answerBackToQuestion) {
			return quizModels.Answer{NextQuestionID: quizModels.QuizReturnByOneQuestion}, false, nil
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
				return quizModels.Answer{}, false, quizModels.ErrQuestionNotFound
			}

			// if pos is valid, find corresponding answer
			return question.Answers[getKeysFromAnswers(question.Answers)[pos]], true, nil
		}
	}

	return quizModels.Answer{}, false, quizModels.ErrQuestionNotFound
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
