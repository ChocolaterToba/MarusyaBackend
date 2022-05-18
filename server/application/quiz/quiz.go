package quiz

import (
	quizModels "cmkids/models/quiz"
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"

	"github.com/pkg/errors"
)

func (app *QuizApp) AddQuizFromFile(filename string, file io.Reader) (err error) {
	extension := filepath.Ext(filename)
	switch extension {
	case ".csv":
		return app.addQuizFromCSV(file)
	default:
		return errors.Wrap(quizModels.ErrUnsupportedFileFormat, fmt.Sprintf("Could not parse %s file format", extension))
	}
}

func (app *QuizApp) addQuizFromCSV(file io.Reader) (err error) {
	csvReader := csv.NewReader(file)
	csvReader.LazyQuotes = true    // so that quotes can be used in csv file
	csvReader.FieldsPerRecord = -1 // so that we can use same file for test and questions
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {
		return errors.Wrap(err, "Cannot to parse file as CSV")
	}

	test := quizModels.Test{}
	err = test.Parse(records)
	if err != nil {
		return errors.Wrap(err, "Cannot parse test")
	}

	err = app.addQuiz(test)
	if err != nil {
		return err
	}

	return nil
}

func (app *QuizApp) addQuiz(quiz quizModels.Test) (err error) {
	_, err = app.quizRepo.CreateEntireQuiz(quiz)
	if err != nil {
		return err
	}

	return nil
}
