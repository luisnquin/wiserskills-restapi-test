package constants

import (
	"os"

	"github.com/luisnquin/restapi-technical-test/src/storage"
)

const APIVersion string = "0.0.1" // Semantic Versioning

var Persistence = func() storage.Persistence {
	persistence := os.Getenv("PERSISTENCE_NAME")

	switch persistence {
	case "PostgreSQL":
		return storage.PostgreSQL
	case "MySQL":
		return storage.MySQL
	default:
		return storage.PostgreSQL
	}
}()
