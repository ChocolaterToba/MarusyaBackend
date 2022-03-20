package basic

import "database/sql"

func insertNewUser(tx *sql.Tx, userID, vkId, name string) (err error) {
	const query = `
				INSERT INTO user_info (user_id, vk_id, name)
				VALUES ($1, $2, $3)`

	_, err = tx.Exec(query, userID, vkId, name)
	return err
}

func insertNewUserName(tx *sql.Tx, userID, name string) (err error) {
	const query = `
				INSERT INTO user_info (user_id, vk_id, name)
				VALUES ($1, $2, $3)`

	_, err = tx.Exec(query, userID, "", name)
	return err
}

func getUserByUserID(tx *sql.Tx, userID string) (name string, err error) {
	const query = `
				SELECT name
				from user_info
				WHERE user_id = $1`

	err = tx.QueryRow(query, userID).Scan(&name)
	return
}

