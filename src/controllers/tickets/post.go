package tickets

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/luisnquin/restapi-technical-test/src/constants"
	"github.com/luisnquin/restapi-technical-test/src/models"
	"github.com/luisnquin/restapi-technical-test/src/storage"
)

func NewTicket() echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			db      = storage.Get(constants.Persistence)
			request = new(models.Ticket)
			err     error
		)

		if err = c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Bad Request",
							"message": "The request body data is not valid",
						},
					},
				},
			})
		}
		if (*request == models.Ticket{}) {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason":  "The request body data is empty",
							"message": "Bad Request",
						},
					},
				},
			})
		}

		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    500,
					Message: "Internal Server Error",
					Errors: []map[string]interface{}{
						{
							"reason":  "Internal Server Error",
							"message": "Database connection failed",
						},
					},
				},
			})
		}
		defer func() {
			if err = db.Close(); err != nil {
				panic(err)
			}
		}()

		var q string
		switch constants.Persistence {
		case storage.PostgreSQL:
			q = "SELECT EXISTS (SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = $1 AND p.id = $2);"
		case storage.MySQL:
			q = "SELECT EXISTS (SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = ? AND p.id = ?);"
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    500,
					Message: "Internal Server Error",
					Errors: []map[string]interface{}{
						{
							"reason":  "Internal Server Error",
						},
					},
				},
			})
		}
		defer func() {
			if err = stmt.Close(); err != nil {
				panic(err)
			}
		}()
		
		var exists bool
		stmt.QueryRowContext(ctx, request.Event, request.Participant).Scan(&exists)
		if exists {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method: "tiokets.post",
				Context: c.Request().URL.String(),
				Error: models.Error{
					Code: 400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason": "Bad Request",
							"message": "The participant was already registered for the event previously",
						},
					},
				},
			})
		}
		
		switch constants.Persistence {
		case storage.PostgreSQL:
			q = "INSERT INTO tickets(participant, event) VALUES($1, $2);"
		case storage.MySQL:
			q = "INSERT INTO tickets(participant, event) VALUES(?, ?);"
		}

		stmt, err = db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    500,
					Message: "Internal Server Error",
					Errors: []map[string]interface{}{
						{
							"reason":  "Internal Server Error",
						},
					},
				},
			})
		}
		defer func() {
			if err = stmt.Close(); err != nil {
				panic(err)
			}
		}()

		r, err := stmt.ExecContext(ctx, request.Participant, request.Event)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Bad Request",
							"message": "The request body data was rejected, not valid",
						},
					},
				},
			})
		}

		if i, _ := r.RowsAffected(); i == 0 {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    400,
					Message: "Bad request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Your changes cannot be implemented",
							"message": "Bad request",
						},
					},
				},
			})
		}
		return c.JSON(http.StatusCreated, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "tickets.post",
			Context:    c.Request().URL.String(),
		})
	}
}
