package basic

import (
	quizApp "cmkids/application/quiz"
	"encoding/json"
	"fmt"
	"net/http"

	"cmkids/models/marusia"

	"go.uber.org/zap"
)

// BasicHandler keep information about apps and cookies needed for basic package
type BasicHandler struct {
	quizApp quizApp.QuizAppInterface
	logger  *zap.Logger
}

func NewBasicHandler(quizApp quizApp.QuizAppInterface, logger *zap.Logger) *BasicHandler {
	return &BasicHandler{
		quizApp: quizApp,
		logger:  logger,
	}
}

//HandleBasicRequest changes password of current user
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

	output.Response, err = handler.quizApp.ProcessBasicRequest(*input)
	if err != nil {
		output.Response = makeErrResponse(err)
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

func makeErrResponse(err error) marusia.Response {
	return marusia.Response{
		Text:       fmt.Sprintf("Произошла ошибка: %s", err.Error()),
		EndSession: true,
	}
}
