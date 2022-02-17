package database

import (
	"context"
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	Db *sql.DB
}

func (db *PostgreSQL) Connect() error {
	conn, err := sql.Open("postgres", os.Getenv("dsn"))
	if err != nil {
		return err
	}
	if err = conn.Ping(); err != nil {
		return err
	}
	db.Db = conn
	return nil
}

func (db *PostgreSQL) Close() error {
	return db.Db.Close()
}

func (db *PostgreSQL) Prepare(stmt string) (*sql.Stmt, error) {
	return db.Db.Prepare(stmt)
}

func (db *PostgreSQL) PrepareContext(ctx context.Context, stmt string) (*sql.Stmt, error) {
	return db.Db.PrepareContext(ctx, stmt)
}
