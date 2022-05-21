package adapter

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DBAdapter struct {
	Conn *sql.DB
}

type Adapter interface {
	InTx(f func(tx *sql.Tx) error) (err error)
}

func (b *DBAdapter) InTx(f func(tx *sql.Tx) error) (err error) {
	tx, err := b.Conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = f(tx)
	return
}

func InitDB(host string, dbname string, user string, password string) (Adapter, error) {

	connFmt := `
        host=%s 
        port=6432
        dbname=%s
        user=%s
        password=%s
        sslmode=verify-full
        sslrootcert=root.crt
`

	conn, err := sql.Open("postgres", fmt.Sprintf(connFmt, host, dbname, user, password))
	if err != nil {
		return nil, err
	}

	return &DBAdapter{Conn: conn}, nil
}
