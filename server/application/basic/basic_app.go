package basic

import (
	"cmkids/models/marusia"
	"cmkids/models/quiz"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

const welcomeNewUserMsg = "Привет! Я постараюсь научить тебя железным правилам детской безопасности. Как я могу к тебе обращаться?"
const welcomeMsg = "Привет, %s?"

type Writer interface {
	InsertNewUser(userID, vkID, name string) (inserted bool, err error)
	IsNewUser(userID string, logger *zap.Logger) (name string, isNew bool, err error)
	InsertNewUserName(userID, name string) (err error)
}

type BasicApp struct {
	Writer
	logger *zap.Logger
}

func NewBasicApp(writer Writer) *BasicApp {
	return &BasicApp{Writer: writer}
}

func (app *BasicApp) ProcessBasicRequest(request marusia.Request, messageID int) (response marusia.Response) {
	if messageID == 0 {
		return app.GetBasicTest(request)
	} else {
		return app.RespondToBasicAnswer(request)
	}
}

func (app *BasicApp) InitIfUserNew(userID string, name string) (response marusia.Response) {
	_, isNew, err := app.IsNewUser(userID)
	if err != nil {
		app.logger.Info(err.Error())
		return marusia.Response{
			EndSession: true,
			Text:       fmt.Sprintf("Ошибочка"),
		}
	}

	if isNew {
		err = app.InsertNewUserName(userID, name)
		if err != nil {
			app.logger.Info(err.Error())
			return marusia.Response{
				EndSession: true,
				Text:       fmt.Sprintf("Ошибочка"),
			}
		}
		return marusia.Response{
			EndSession: false,
			Text:       fmt.Sprintf(welcomeMsg, name),
		}
	}

	return marusia.Response{EndSession: false, Text: ""}
}

func (app *BasicApp) Activate(userID string) (response marusia.Response) {
	app.logger.Info("I AM OK 0")
	name, isNew, err := app.IsNewUser(userID)
	if err != nil {
		app.logger.Info(err.Error())
		return marusia.Response{
			EndSession: true,
			Text:       fmt.Sprintf("Ошибочка"),
		}
	}
	app.logger.Info("I AM OK 1")
	response.EndSession = false
	if isNew {
		response.Text = welcomeNewUserMsg
		return
	} else if !isNew && name != "" {
		response.Text = fmt.Sprintf(welcomeMsg, name)
		return
	}
	return
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
		},
		EndSession: false,
	}
}

func (app *BasicApp) RespondToBasicAnswer(request marusia.Request) (response marusia.Response) {
	response = marusia.Response{EndSession: false}
	switch strings.ToLower(request.Command) {
	case strings.ToLower(quiz.BASIC_TEST_YES):
		response.Text = quiz.BASIC_ANSWER_YES
	case strings.ToLower(quiz.BASIC_TEST_NO):
		response.Text = quiz.BASIC_ANSWER_NO
	default:
		response.Text = fmt.Sprintf("Ошибочная команда: %s", request.Command)
	}
	return response
}
