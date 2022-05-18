package quiz

import (
	"cmkids/adapter"
	authModels "cmkids/models/auth"
	quizModels "cmkids/models/quiz"
	"encoding/json"

	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"

	"database/sql"
	"database/sql/driver"
	"fmt"
)

type QuizRepoInterface interface {
	GetPastAnswers(userID uint64) (answers []quizModels.Answer, err error)
	SetPastAnswers(userID uint64, answers []quizModels.Answer) (err error)
	GetQuestion(questionID uint64) (question quizModels.Question, err error)
	GetQuestionInTest(testID uint64, questionInTestID uint64) (question quizModels.Question, err error)
	CreateEntireQuiz(quiz quizModels.Test) (quizID uint64, err error)
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
				return quizModels.ErrQuestionNotFound
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
				return quizModels.ErrQuestionNotFound
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
				return quizModels.ErrQuestionNotFound
			}
			return fmt.Errorf("error in QuizRepo: could not get question by question_in_test_id: %s", err)
		}
		question.Answers = answers

		return nil
	})

	return question, err
}

// CreateQuiz inserts quiz (including questions and answers) into database
func (repo *QuizRepo) CreateEntireQuiz(quiz quizModels.Test) (quizID uint64, err error) {
	err = repo.conn.InTx(func(tx *sql.Tx) error {
		quizID, err := repo.createQuiz(tx, quiz)
		if err != nil {
			return err
		}

		for questionInTestID := range quiz.Questions {
			question := quiz.Questions[questionInTestID]
			question.TestID = quizID
			quiz.Questions[questionInTestID] = question
		}

		err = repo.createQuestions(tx, quiz.Questions)
		if err != nil {
			return err
		}

		rootQuestion, err := repo.GetQuestion(quizModels.QuizRootID)
		if err != nil {
			return err
		}

		firstQuizQuestion, err := repo.getQuestionInTestTx(tx, quizID, 1) // Tx since GetQuestionInTest would starts its own transaction
		if err != nil {
			return err
		}

		rootQuestion.Answers[quiz.Title] = quizModels.Answer{NextQuestionID: firstQuizQuestion.QuestionID}
		err = repo.updateQuestion(tx, rootQuestion)
		if err != nil {
			return err
		}

		return nil
	})

	return quizID, err
}

func (repo *QuizRepo) createQuiz(tx *sql.Tx, quiz quizModels.Test) (quizID uint64, err error) {
	const createQuizQuery = `INSERT INTO quiz (title, backtracking_enabled, calculate_correctness)
							 VALUES ($1, $2, $3)
							 RETURNING id`

	err = tx.QueryRow(createQuizQuery, quiz.Title, quiz.BackTrackingEnabled, quiz.CalculateCorrectness).Scan(&quizID)
	if err != nil {
		return 0, errors.Wrap(err, "Error in QuizRepo: could not create quiz metadata")
	}

	return quizID, nil
}

func (repo *QuizRepo) createQuestions(tx *sql.Tx, questionsMap map[uint64]quizModels.Question) (err error) {
	sb := sqlbuilder.NewInsertBuilder()
	sb.InsertInto("question")
	sb.Cols("question_in_test_id", "test_id", "text", "next_question_ids")
	for questionInTestID := range questionsMap {
		question := questionsMap[questionInTestID]
		sb.Values(question.QuestionInTestID, question.TestID, question.Text, cusAnswersMap(question.Answers))
	}

	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err = tx.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "Error in QuizRepo: could not create questions")
	}

	return nil
}

func (repo *QuizRepo) getQuestionInTestTx(tx *sql.Tx, testID uint64, questionInTestID uint64) (question quizModels.Question, err error) {
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
			return quizModels.Question{}, quizModels.ErrQuestionNotFound
		}
		return quizModels.Question{}, fmt.Errorf("error in QuizRepo: could not get question by question_in_test_id: %s", err)
	}
	question.Answers = answers

	return question, err
}

func (repo *QuizRepo) updateQuestion(tx *sql.Tx, question quizModels.Question) (err error) {
	sb := sqlbuilder.NewUpdateBuilder()
	sb.Update("question")
	sb.Set(
		sb.Assign("question_in_test_id", question.QuestionInTestID),
		sb.Assign("test_id", question.TestID),
		sb.Assign("text", question.Text),
		sb.Assign("next_question_ids", cusAnswersMap(question.Answers)))
	sb.Where(sb.Equal("question_id", question.QuestionID))

	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err = tx.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "Error in QuizRepo: could not update question")
	}

	return nil
}

type cusAnswersMap map[string]quizModels.Answer

// Returns the JSON-encoded representation
func (a cusAnswersMap) Value() (driver.Value, error) {
	return json.Marshal(a)
}

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
