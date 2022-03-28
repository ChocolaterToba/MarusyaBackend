package user

import (
	"cmkids/adapter"
	authModels "cmkids/models/auth"
	"errors"

	"database/sql"
	"fmt"
)

type UserRepoInterface interface {
	RegisterUser(username string, appID string) (userID uint64, err error)
	GetUserByUserID(userID uint64) (user authModels.User, err error)
	GetUserByAppID(appID string) (user authModels.User, err error)
	LoginUser(userID uint64, sessionID string) (err error)
	GetUserBySessionID(sessionID string) (user authModels.User, err error)
}

type UserRepo struct {
	conn adapter.Adapter
}

func NewUserRepo(conn adapter.Adapter) *UserRepo {
	return &UserRepo{conn: conn}
}

func (repo *UserRepo) RegisterUser(username string, appID string) (userID uint64, err error) {
	err = repo.conn.InTx(func(tx *sql.Tx) error {
		const query = `INSERT INTO account(username, application_id)
					   VALUES ($1, $2)
					   RETURNING user_id`

		err = tx.QueryRow(query, username, appID).Scan(&userID)
		if err != nil {
			return fmt.Errorf("error in UserRepo: could not register user: %s", err)
		}
		return nil
	})

	return userID, err
}

func (repo *UserRepo) GetUserByUserID(userID uint64) (user authModels.User, err error) {
	err = repo.conn.InTx(func(tx *sql.Tx) error {
		const query = `SELECT user_id, application_id, session_id, username
					   FROM account
					   WHERE user_id = $1`

		err = tx.QueryRow(query, userID).Scan(&user.UserID, &user.ApplicationID, &user.SessionID, &user.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return authModels.ErrUserNotFound
			}
			return fmt.Errorf("error in UserRepo: could not get user by user id: %s", err)
		}
		return nil
	})

	return user, err
}

func (repo *UserRepo) GetUserByAppID(appID string) (user authModels.User, err error) {
	err = repo.conn.InTx(func(tx *sql.Tx) error {
		const query = `SELECT user_id, application_id, session_id, username
					   FROM account
					   WHERE application_id = $1`

		err = tx.QueryRow(query, appID).Scan(&user.UserID, &user.ApplicationID, &user.SessionID, &user.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return authModels.ErrUserNotFound
			}
			return fmt.Errorf("error in UserRepo: could not get user by app id: %s", err)
		}
		return nil
	})

	return user, err
}

func (repo *UserRepo) GetUserBySessionID(sessionID string) (user authModels.User, err error) {
	err = repo.conn.InTx(func(tx *sql.Tx) error {
		const query = `SELECT user_id, application_id, session_id, username
					   FROM account
					   WHERE session_id = $1`

		err = tx.QueryRow(query, sessionID).Scan(&user.UserID, &user.ApplicationID, &user.SessionID, &user.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return authModels.ErrUserNotFound
			}
			return fmt.Errorf("error in UserRepo: could not get user by session id: %s", err)
		}
		return nil
	})

	return user, err
}

func (repo *UserRepo) LoginUser(userID uint64, sessionID string) (err error) {
	err = repo.conn.InTx(func(tx *sql.Tx) error {
		const query = `UPDATE account
					   SET session_id = $2
					   WHERE user_id = $1`

		result, err := tx.Exec(query, userID, sessionID)
		if err != nil {
			return fmt.Errorf("error in UserRepo: could not login user: %s", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return authModels.ErrUserNotFound
			}
			return fmt.Errorf("error in UserRepo: could not login user: %s", err)
		}
		if rowsAffected != 1 {
			return authModels.ErrUserNotFound
		}

		return nil
	})

	return err
}
