package basic

import (
	"encoding/json"
	"net/http"

	"cmkids/models/marusia"

	"go.uber.org/zap"
)

type BasicAppInterface interface {
	Activate(userID string) (response marusia.Response)
	InitIfUserNew(userID string, name string) (response marusia.Response)
	ProcessBasicRequest(request marusia.Request, messageID int) (response marusia.Response)
	GetBasicTest(request marusia.Request) (response marusia.Response)
	RespondToBasicAnswer(request marusia.Request) (response marusia.Response)
}

// BasicHandler keep information about apps and cookies needed for basic package
type BasicHandler struct {
	basicApp BasicAppInterface
	logger   *zap.Logger
}

func NewBasicHandler(basicApp BasicAppInterface, logger *zap.Logger) *BasicHandler {
	return &BasicHandler{
		basicApp: basicApp,
		logger:   logger,
	}
}

//HandleBasicRequest changes password of current user
func (handler *BasicHandler) HandleBasicRequest(w http.ResponseWriter, r *http.Request) {
	input := new(marusia.RequestBody)
	handler.logger.Info("I AM OK -2")
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		handler.logger.Info(err.Error(), zap.String("url", r.RequestURI), zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	handler.logger.Info("I AM OK -1")
	output := new(marusia.ResponseBody)

	output.Session = input.Session
	output.Version = input.Version

	// logic starts here
	if input.Session.New {
		if handler.basicApp.Activate == nil {
			handler.logger.Info("I AM NOT /1")
		}
		handler.logger.Info("I AM OK /1")
		output.Response = handler.basicApp.Activate(input.User.UserID)
	} else if input.MessageID == 1 {
		output.Response = handler.basicApp.InitIfUserNew(input.User.UserID, input.Command)
		if output.Response.Text == "" {
			output.Response = handler.basicApp.ProcessBasicRequest(input.Request, input.MessageID)
		}
	} else {
		output.Response = handler.basicApp.ProcessBasicRequest(input.Request, input.MessageID)
	}
	// logic ends here

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
