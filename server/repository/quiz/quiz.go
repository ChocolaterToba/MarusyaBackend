package quiz

import (
	"cmkids/adapter"
	quizModels "cmkids/models/quiz"
	"encoding/json"
	"errors"

	"database/sql"
	"fmt"
)

type QuizRepoInterface interface {
	GetCurrentQuestionID(userID uint64) (questionID uint64, err error)
	GetQuestion(questionID uint64) (question quizModels.Question, err error)
}

type QuizRepo struct {
	conn adapter.Adapter
}

func NewQuizRepo(conn adapter.Adapter) *QuizRepo {
	return &QuizRepo{conn: conn}
}

func (repo *QuizRepo) GetCurrentQuestionID(userID uint64) (questionID uint64, err error) {
	err = repo.conn.InTx(func(tx *sql.Tx) error {
		const query = `SELECT current_question_id
					   FROM account
					   WHERE user_id = $1`

		err = tx.QueryRow(query, userID).Scan(&questionID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return quizModels.ErrCurrentQuestionNotFound
			}
			return fmt.Errorf("error in QuizRepo: could not get currect quiestion ID: %s", err)
		}
		return nil
	})

	return questionID, err
}

func (repo *QuizRepo) GetQuestion(questionID uint64) (question quizModels.Question, err error) {
	err = repo.conn.InTx(func(tx *sql.Tx) error {
		const query = `SELECT question_id, test_id, text, next_question_ids
					   FROM question
					   WHERE question_id = $1`

		nextQuestionIDs := make(cusjsonb)
		err = tx.QueryRow(query, questionID).Scan(&question.QuestionID, &question.TestID, &question.Text, &nextQuestionIDs)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return quizModels.ErrNextQuestionNotFound
			}
			return fmt.Errorf("error in UserRepo: could not get question: %s", err)
		}
		question.NextQuestionIDs = nextQuestionIDs

		return nil
	})

	return question, err
}

type cusjsonb map[string]uint64

// Decodes a JSON-encoded value
func (a *cusjsonb) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Unmarshal from json to map[string]uint64
	if err := json.Unmarshal(b, a); err != nil {
		return err
	}
	return nil
}
