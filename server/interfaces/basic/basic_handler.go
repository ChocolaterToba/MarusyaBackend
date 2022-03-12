package basic

import (
	"encoding/json"
	"net/http"

	"cmkids/application/basic"
	"cmkids/models/marusia"

	"go.uber.org/zap"
)

// BasicHandler keep information about apps and cookies needed for basic package
type BasicHandler struct {
	basicApp basic.BasicAppInterface
	logger   *zap.Logger
}

func NewBasicHandler(basicApp basic.BasicAppInterface, logger *zap.Logger) *BasicHandler {
	return &BasicHandler{
		basicApp: basicApp,
		logger:   logger,
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

	// logic starts here

	output.Response = handler.basicApp.ProcessBasicRequest(input.Request, input.MessageID)

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
