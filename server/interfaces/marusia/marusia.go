package marusia

import (
	quizApp "cmkids/application/quiz"
	"encoding/json"
	"fmt"
	"net/http"

	"cmkids/models/marusia"

	"go.uber.org/zap"
)

// MarusiaHandler keep information about apps and cookies needed for marusia package
type MarusiaHandler struct {
	quizApp quizApp.QuizAppInterface
	logger  *zap.Logger
}

func NewMarusiaHandler(quizApp quizApp.QuizAppInterface, logger *zap.Logger) *MarusiaHandler {
	return &MarusiaHandler{
		quizApp: quizApp,
		logger:  logger,
	}
}

//HandleMarusiaRequest changes password of current user
func (handler *MarusiaHandler) HandleMarusiaRequest(w http.ResponseWriter, r *http.Request) {
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

	output.Response, err = handler.quizApp.ProcessMarusiaRequest(*input)
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
		Text:       []string{fmt.Sprintf("Произошла ошибка: %s", err.Error())},
		EndSession: false,
	}
}
