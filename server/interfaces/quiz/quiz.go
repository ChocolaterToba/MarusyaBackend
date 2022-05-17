package quiz

import (
	quizApp "cmkids/application/quiz"
	"fmt"
	"net/http"
	"strconv"

	"cmkids/models/marusia"
	quizModels "cmkids/models/quiz"

	"go.uber.org/zap"
)

// QuizHandler handles quiz-specific requests
type QuizHandler struct {
	quizApp quizApp.QuizAppInterface
	logger  *zap.Logger
}

func NewQuizHandler(quizApp quizApp.QuizAppInterface, logger *zap.Logger) *QuizHandler {
	return &QuizHandler{
		quizApp: quizApp,
		logger:  logger,
	}
}

const maxPostBodySize = 8 * 1024 * 1024 // 8mb
//HandleAddQuiz parses incoming FormData and makes qizes out of it
func (handler *QuizHandler) HandleAddQuiz(w http.ResponseWriter, r *http.Request) {
	bodySize := r.ContentLength

	if bodySize <= 0 { // No files were
		handler.logger.Info(quizModels.ErrNoFile.Error(),
			zap.String("url", r.RequestURI),
			zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if bodySize > int64(maxPostBodySize) {
		handler.logger.Info(quizModels.ErrFileTooLarge.Error(),
			zap.String("url", r.RequestURI),
			zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(bodySize)
	quizAmount, _ := strconv.ParseUint(r.FormValue("quizAmount"), 10, 64)
	if quizAmount == 0 {
		handler.logger.Info(quizModels.ErrIncorrectQuizAmount.Error(),
			zap.String("url", r.RequestURI),
			zap.String("method", r.Method))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i := uint64(1); i <= quizAmount; i++ {
		fieldName := "file" + strconv.FormatUint(i, 10)
		file, header, err := r.FormFile(fieldName)
		if err != nil {
			handler.logger.Info(err.Error(),
				zap.String("url", r.RequestURI),
				zap.String("method", r.Method))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = handler.quizApp.AddQuizFromFile(header.Filename, file)
		if err != nil {
			handler.logger.Info(err.Error(),
				zap.String("url", r.RequestURI),
				zap.String("method", r.Method))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		file.Close()
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func makeErrResponse(err error) marusia.Response {
	return marusia.Response{
		Text:       []string{fmt.Sprintf("Произошла ошибка: %s", err.Error())},
		EndSession: true,
	}
}
