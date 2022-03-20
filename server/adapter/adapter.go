package adapter

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type DBAdapter struct {
	Conn *sql.DB
	logger *zap.Logger
}

type Adapter interface {
	InTx(f func(tx *sql.Tx) error) (err error)
}

func (b *DBAdapter) InTx(f func(tx *sql.Tx) error) (err error) {
	tx, err := b.Conn.Begin()
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // fallthrough panic after rollback on caught panic
		} else if err != nil {
			_ = tx.Rollback() // if error during computations
		} else {
			err = tx.Commit() // all good
		}
	}()
	if tx == nil {
		b.logger.Info("ADAPTER ERR")

	}
	b.logger.Info("ADAPTER")
	err = f(tx)
	return
}

func InitDB(host, pass string, logger *zap.Logger) (Adapter, error) {

	connFmt := `
        host=%s 
        port=6432
        dbname=cmkids
        user=mikhail
        password=%s
        sslmode=verify-full
        sslrootcert=/home/username/.postgresql/root.crt
`

	conn, err := sql.Open("postgres", fmt.Sprintf(connFmt, host, pass))
	if err != nil {
		return nil, err
	}

	return &DBAdapter{Conn: conn, logger: logger}, nil
}
