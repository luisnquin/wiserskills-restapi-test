package persistence

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Help() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome!\n\nThere are an endpoint to build the database in MySQL or PostgreSQL, just fill the request body data with your session credentials later, in the URL, set the database name of your preference and send the request\n\nFile I mean:\n\t -> [ROOT_DIR]/build.rest\n \n\nAnother option: \ncurl -X POST http://127.0.0.1:8000/persistence/build/<database-name> \\\n\t-H 'Content-Type: application/json' \\\n\t-d '{\"dbname\":\"\", \"user\":\"\", \"password\": \"\"}'  \n\nI made it fast so it may fail, in which case you will have to opt for a manual configuration, your tools are in:\n -> [ROOT_DIR]/src/database/<persistence-name>.sql, just press [Ctrl+A] and it will be ready for pasting\n -> [ROOT_DIR]/.env.example to create one customised DSN for your session")
	}
}
