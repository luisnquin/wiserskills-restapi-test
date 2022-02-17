package tickets

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
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
					Message: "Bad request",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Bad request",
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
							"reason":  "The request body is empty",
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
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
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

		
		q := "SELECT EXISTS (SELECT t.id AS id, CONCAT(p.firstname, ' ',p.lastname) AS participant, e.name AS event FROM tickets AS t INNER JOIN events AS e ON e.id=t.event INNER JOIN participants AS p ON p.id=t.participant WHERE e.id = ? AND p.id = ?);"
		if constants.Persistence == storage.PostgreSQL {
			q = sqlx.Rebind(sqlx.DOLLAR, q)
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
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
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
					Message: "Bad request",
					Errors: []map[string]interface{}{
						{
							"reason": err,
							"message": "Bad request",
						},
					},
				},
			})
		}
		
		q = "INSERT INTO tickets(participant, event) VALUES(?, ?);"

		if constants.Persistence == storage.PostgreSQL {
			q = sqlx.Rebind(sqlx.DOLLAR, q)
		}
		stmt, err = db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Internal server error",
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
					Message: "Bad request",
					Errors: []map[string]interface{}{
						{
							"reason":  err,
							"message": "Bad request",
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
