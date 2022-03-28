package auth

import (
	authModels "cmkids/models/auth"
	"cmkids/models/marusia"
	userRepo "cmkids/repository/user"

	"fmt"

	"go.uber.org/zap"
)

type AuthAppInterface interface {
	Login(input marusia.RequestBody) (response marusia.Response, err error)
	Register(input marusia.RequestBody) (response marusia.Response, err error)
	GetUserIDBySessionID(sessionID string) (userID uint64, err error)
}

type AuthApp struct {
	userRepo userRepo.UserRepoInterface
	logger   *zap.Logger
}

func NewAuthApp(userRepo userRepo.UserRepoInterface, logger *zap.Logger) *AuthApp {
	return &AuthApp{
		userRepo: userRepo,
		logger:   logger}
}

func (app *AuthApp) Login(input marusia.RequestBody) (response marusia.Response, err error) {
	_, err = app.GetUserIDBySessionID(input.Session.SessionID)
	if err == nil { // user is already logged in
		return marusia.Response{
			Text:       authModels.MsgAlreadyLoggedIn,
			EndSession: true,
		}, nil
	}
	if err != authModels.ErrUserNotFound {
		return marusia.Response{}, err
	}

	user, err := app.userRepo.GetUserByAppID(input.Session.Application.ApplicationID)
	if err != nil {
		if err == authModels.ErrUserNotFound { // User is not registered
			return marusia.Response{
				Text:       authModels.MsgRegistrationPrompt,
				EndSession: false,
			}, nil
		}
		return marusia.Response{}, err
	}

	err = app.userRepo.LoginUser(user.UserID, input.Session.SessionID)
	if err != nil {
		return marusia.Response{}, err
	}

	return marusia.Response{
		Text:       fmt.Sprintf(authModels.MsgWelcome, user.Username),
		EndSession: false,
	}, nil
}

func (app *AuthApp) Register(input marusia.RequestBody) (response marusia.Response, err error) {
	username := input.Request.OriginalUtterance // TODO: clean username
	_, err = app.userRepo.RegisterUser(username, input.Session.Application.ApplicationID)
	if err != nil {
		return marusia.Response{}, err
	}

	return app.Login(input)
}

func (app *AuthApp) GetUserIDBySessionID(sessionID string) (userID uint64, err error) {
	user, err := app.userRepo.GetUserBySessionID(sessionID)
	return user.UserID, err
}
