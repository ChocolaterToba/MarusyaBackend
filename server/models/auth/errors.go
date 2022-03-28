package auth

import "errors"

var (
	ErrUserNotFound = errors.New("Не удалось найти пользователя")
)

const (
	MsgAlreadyLoggedIn    = "Скилл уже был активирован."
	MsgRegistrationPrompt = "Привет! Как я могу к тебе обращаться?"
	MsgWelcome            = "Привет, %s."
)
