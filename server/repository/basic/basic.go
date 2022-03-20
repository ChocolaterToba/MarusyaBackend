package basic

import (
	"cmkids/adapter"
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/zap"
)

type Repository struct {
	conn adapter.Adapter
}

func NewRepository(conn adapter.Adapter) *Repository {
	return &Repository{conn: conn}
}

func (r *Repository) InsertNewUser(userID, vkID, name string) (inserted bool, err error) {
	err = r.conn.InTx(func(tx *sql.Tx) error {
		_, err = getUserByUserID(tx, userID)
		if err != nil && !errors.Is(err, sql.ErrNoRows){
			return fmt.Errorf("error in InsertNewUser.getUserByUserID: %w", err)
		}

		inserted= false
		err = insertNewUser(tx, userID, vkID, name)
		if err != nil {
			return fmt.Errorf("error in InsertNewUser.insertNewUser: %w", err)
		}
		inserted = true
		return nil
	})

	return
}

func (r *Repository) InsertNewUserName(userID, name string) (err error) {
	err = r.conn.InTx(func(tx *sql.Tx) error {
		err = insertNewUserName(tx, userID, name)
		if err != nil {
			return fmt.Errorf("error in InsertNewUser.insertNewUserName: %w", err)
		}
		return nil
	})

	return err
}

func (r *Repository) GetUserByUserID(userID string) (name string, err error) {
	err = r.conn.InTx(func(tx *sql.Tx) error{
		name, err = getUserByUserID(tx, userID)
		if err != nil {
			return fmt.Errorf("error in InsertNewUser.insertNewUser: %w", err)
		}
		return nil
	})

	return
}

func (r *Repository) IsNewUser(userID string, logger *zap.Logger) (name string, isNew bool, err error) {
	err = r.conn.InTx(func(tx *sql.Tx) error {
		if tx == nil {
			logger.Info("IAM NIL")
			return fmt.Errorf("IAM NIL")
		}
		logger.Info("I AM OK 3")
		isNew = false
		name, err = getUserByUserID(tx, userID)
		if err != nil {
			logger.Info("I AM OK 3")
			if errors.Is(err, sql.ErrNoRows) {
				isNew = true
				return nil
			}
			return fmt.Errorf("error in InsertNewUser.insertNewUser: %w", err)
		}
		return nil
	},
	)

	return
}
