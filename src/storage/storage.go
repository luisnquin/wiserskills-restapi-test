package storage

import "github.com/luisnquin/restapi-technical-test/src/database"

type Persistence uint8

const (
	PostgreSQL Persistence = 1
	MySQL      Persistence = 2
)

func Get(persistence Persistence) database.Connecter {
	switch persistence {
	case PostgreSQL:
		return &database.PostgreSQL{}
	case MySQL:
		return &database.MySQL{}
	}
	return nil
}
