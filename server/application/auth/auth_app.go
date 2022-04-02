package auth

import (
	authModels "cmkids/models/auth"
	"cmkids/models/marusia"
	quizRepo "cmkids/repository/quiz"
	userRepo "cmkids/repository/user"

	"fmt"

	"go.uber.org/zap"
)

type AuthAppInterface interface {
	Login(input marusia.RequestBody) (response marusia.Response, finished bool, err error)
	Register(input marusia.RequestBody) (response marusia.Response, finished bool, err error)
	GetUserIDBySessionID(sessionID string) (userID uint64, err error)
}

type AuthApp struct {
	userRepo userRepo.UserRepoInterface
	quizRepo quizRepo.QuizRepoInterface
	logger   *zap.Logger
}

func NewAuthApp(userRepo userRepo.UserRepoInterface, quizRepo quizRepo.QuizRepoInterface, logger *zap.Logger) *AuthApp {
	return &AuthApp{
		userRepo: userRepo,
		quizRepo: quizRepo,
		logger:   logger}
}

// Login tries to log user in tying their sessionID to applicationID
// If err is not nil, pass it higher
// If finsihed is false, pass response and err higher
// If finished is true, use response.Text as starting point when building response text
func (app *AuthApp) Login(input marusia.RequestBody) (response marusia.Response, finished bool, err error) {
	_, err = app.GetUserIDBySessionID(input.Session.SessionID)
	if err == nil { // user is already logged in
		return marusia.Response{
			Text:       []string{authModels.MsgAlreadyLoggedIn},
			EndSession: false,
		}, false, nil
	}
	if err != authModels.ErrUserNotFound {
		return marusia.Response{}, false, err
	}

	user, err := app.userRepo.GetUserByAppID(input.Session.Application.ApplicationID)
	if err != nil {
		if err == authModels.ErrUserNotFound { // User is not registered
			return marusia.Response{
				Text:       []string{authModels.MsgRegistrationPrompt},
				EndSession: false,
			}, false, nil
		}
		return marusia.Response{}, false, err
	}

	err = app.userRepo.LoginUser(user.UserID, input.Session.SessionID)
	if err != nil {
		return marusia.Response{}, false, err
	}

	return marusia.Response{
		Text:       []string{fmt.Sprintf(authModels.MsgWelcome, user.Username)},
		EndSession: false,
	}, true, nil
}

// Register tries to register user using name provided in request
// It also logs them in subsequently
// If err is not nil, pass it higher
// If finsihed is false, pass response and err higher
// If finished is true, use response.Text as starting point when building response text
func (app *AuthApp) Register(input marusia.RequestBody) (response marusia.Response, finished bool, err error) {
	username := input.Request.OriginalUtterance // TODO: clean username
	_, err = app.userRepo.RegisterUser(username, input.Session.Application.ApplicationID)
	if err != nil {
		return marusia.Response{}, false, err
	}

	return app.Login(input)
}

func (app *AuthApp) GetUserIDBySessionID(sessionID string) (userID uint64, err error) {
	user, err := app.userRepo.GetUserBySessionID(sessionID)
	return user.UserID, err
}
