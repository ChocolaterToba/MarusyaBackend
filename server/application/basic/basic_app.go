package basic

import (
	"cmkids/models/marusia"
	"cmkids/models/quiz"
	"fmt"
)

type BasicApp struct {
}

func NewBasicApp() *BasicApp {
	return &BasicApp{}
}

type BasicAppInterface interface {
	ProcessBasicRequest(request marusia.Request, messageID int) (response marusia.Response)
	GetBasicTest(request marusia.Request) (response marusia.Response)
	RespondToBasicAnswer(request marusia.Request) (response marusia.Response)
}

func (app *BasicApp) ProcessBasicRequest(request marusia.Request, messageID int) (response marusia.Response) {
	if messageID == 0 {
		return app.GetBasicTest(request)
	} else {
		return app.RespondToBasicAnswer(request)
	}
}

func (app *BasicApp) GetBasicTest(request marusia.Request) (response marusia.Response) {
	return marusia.Response{
		Text: "К тебе подошёл незнакомец и попросил конфету. Ты...",
		TTS:  "",
		Buttons: []marusia.Button{
			{
				Title: quiz.BASIC_TEST_YES,
			},
			{
				Title: quiz.BASIC_TEST_NO,
			},
			{
				Title: quiz.BASIC_TEST_MEXICO,
			},
		},
		EndSession: false,
	}
}

func (app *BasicApp) RespondToBasicAnswer(request marusia.Request) (response marusia.Response) {
	response = marusia.Response{EndSession: true}
	switch request.Command {
	case quiz.BASIC_TEST_YES:
		response.Text = quiz.BASIC_ANSWER_YES
	case quiz.BASIC_TEST_NO:
		response.Text = quiz.BASIC_ANSWER_NO
	case quiz.BASIC_TEST_MEXICO:
		response.Text = quiz.BASIC_ANSWER_MEXICO
	default:
		response.Text = fmt.Sprintf("Ошибочная команда: %s", request.Command)
	}
	return response
}
