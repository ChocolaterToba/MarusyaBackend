package quiz

import (
	"cmkids/adapter"
	authModels "cmkids/models/auth"
	quizModels "cmkids/models/quiz"
	"encoding/json"
	"errors"

	"database/sql"
	"database/sql/driver"
	"fmt"
)

type QuizRepoInterface interface {
	GetPastAnswers(userID uint64) (answers []quizModels.Answer, err error)
	SetPastAnswers(userID uint64, answers []quizModels.Answer) (err error)
	GetQuestion(questionID uint64) (question quizModels.Question, err error)
	GetQuestionInTest(testID uint64, questionInTestID uint64) (question quizModels.Question, err error)
}

type QuizRepo struct {
	conn adapter.Adapter
}

func NewQuizRepo(conn adapter.Adapter) *QuizRepo {
	return &QuizRepo{conn: conn}
}

func (repo *QuizRepo) GetPastAnswers(userID uint64) (answers []quizModels.Answer, err error) {
	err = repo.conn.InTx(func(tx *sql.Tx) error {
		const query = `SELECT past_answers
					   FROM account
					   WHERE user_id = $1`

		var result cusAnswersSlice

		err = tx.QueryRow(query, userID).Scan(&result)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return quizModels.ErrCurrentQuestionNotFound
			}
			return fmt.Errorf("error in QuizRepo: could not get past answers: %s", err)
		}

		answers = result
		return nil
	})

	return answers, err
}

func (repo *QuizRepo) SetPastAnswers(userID uint64, answers []quizModels.Answer) (err error) {
	err = repo.conn.InTx(func(tx *sql.Tx) error {
		query := `UPDATE account
				  SET past_answers = $2
				  WHERE user_id = $1`

		result, err := tx.Exec(query, userID, cusAnswersSlice(answers))
		if err != nil {
			return fmt.Errorf("error in QuizRepo: could not set past answers: %s", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("error in QuizRepo: could not set past answers: %s", err)
		}
		if rowsAffected != 1 {
			return authModels.ErrUserNotFound
		}

		return nil
	})

	return err
}

func (repo *QuizRepo) GetQuestion(questionID uint64) (question quizModels.Question, err error) {
	err = repo.conn.InTx(func(tx *sql.Tx) error {
		const query = `SELECT question_id, question_in_test_id, test_id, text, next_question_ids
					   FROM question
					   WHERE question_id = $1`

		answers := make(cusAnswersMap)
		err = tx.QueryRow(query, questionID).Scan(
			&question.QuestionID, &question.QuestionInTestID, &question.TestID,
			&question.Text, &answers,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return quizModels.ErrNextQuestionNotFound
			}
			return fmt.Errorf("error in QuizRepo: could not get question by question_id: %s", err)
		}
		question.Answers = answers

		return nil
	})

	return question, err
}

func (repo *QuizRepo) GetQuestionInTest(testID uint64, questionInTestID uint64) (question quizModels.Question, err error) {
	err = repo.conn.InTx(func(tx *sql.Tx) error {
		const query = `SELECT question_id, question_in_test_id, test_id, text, next_question_ids
					   FROM question
					   WHERE test_id = $1 and question_in_test_id = $2`

		answers := make(cusAnswersMap)
		err = tx.QueryRow(query, testID, questionInTestID).Scan(
			&question.QuestionID, &question.QuestionInTestID, &question.TestID,
			&question.Text, &answers,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return quizModels.ErrNextQuestionNotFound
			}
			return fmt.Errorf("error in QuizRepo: could not get question by question_in_test_id: %s", err)
		}
		question.Answers = answers

		return nil
	})

	return question, err
}

type cusAnswersMap map[string]quizModels.Answer

// Decodes a JSON-encoded value
func (a *cusAnswersMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Unmarshal from json to map[string]quizModels.Answer
	if err := json.Unmarshal(b, a); err != nil {
		return err
	}
	return nil
}

type cusAnswersSlice []quizModels.Answer

// Returns the JSON-encoded representation
func (a cusAnswersSlice) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Decodes a JSON-encoded value
func (a *cusAnswersSlice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	// Unmarshal from json to []quizModels.Answer
	x := make([]quizModels.Answer, 0)
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	*a = x
	return nil
}
