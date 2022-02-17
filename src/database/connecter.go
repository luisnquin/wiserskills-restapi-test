package database

import (
	"context"
	"database/sql"
)

type Connecter interface {
	Connect() error
	Close() error
	Prepare(stmt string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, stmt string) (*sql.Stmt, error)
}
