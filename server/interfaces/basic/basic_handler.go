package basic

import (
	"cmkids/models/marusia"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// BasicHandler keep information about apps and cookies needed for basic package
type BasicHandler struct {
	// userApp         application.UserAppInterface
	// authApp         application.AuthAppInterface
	// cookieApp       application.CookieAppInterface
	// followApp       application.FollowAppInterface
	// s3App           application.S3AppInterface
	// notificationApp application.NotificationAppInterface
	logger *zap.Logger
}

func NewBasicHandler(logger *zap.Logger) *BasicHandler {
	return &BasicHandler{
		logger: logger,
	}
}

//HandleChangePassword changes password of current user
func (handler *BasicHandler) HandleBasicRequest(w http.ResponseWriter, r *http.Request) {
	input := new(marusia.RequestBody)
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		handler.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output := new(marusia.ResponseBody)

	output.Session = input.Session
	output.Version = input.Version

	output.Response.Text = input.Command
	output.Response.EndSession = true

	responseBody, err := json.Marshal(output)
	if err != nil {
		handler.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
