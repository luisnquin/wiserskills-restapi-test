package database

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	Db *sql.DB
}

func (db *MySQL) Connect() error {
	conn, err := sql.Open("mysql", os.Getenv("dsn"))
	if err != nil {
		return err
	}
	if err = conn.Ping(); err != nil {
		return err
	}
	db.Db = conn
	return nil
}

func (db *MySQL) Close() error {
	return db.Db.Close()
}

func (db *MySQL) Prepare(stmt string) (*sql.Stmt, error) {
	return db.Db.Prepare(stmt)
}

func (db *MySQL) PrepareContext(ctx context.Context, stmt string) (*sql.Stmt, error) {
	return db.Db.PrepareContext(ctx, stmt)
}
