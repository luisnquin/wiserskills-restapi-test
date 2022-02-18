package tickets

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/luisnquin/restapi-technical-test/src/constants"
	"github.com/luisnquin/restapi-technical-test/src/models"
	"github.com/luisnquin/restapi-technical-test/src/storage"
)

func ModifyTicketById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			db      = storage.Get(constants.Persistence)
			request = new(models.Ticket)
		)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.patch",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    422,
					Message: "Unprocessable Entity",
					Errors: []map[string]interface{}{
						{
							"reason":  "Unprocessable Entity",
							"message": "The ID parameter cannot be processed as integer",
						},
					},
				},
			})
		}

		if err = c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.patch",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason": "Bad Request",
							"message": "The request body data is not valid",
						},
					},
				},
			})
		}

		if request == new(models.Ticket) {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.patch",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
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
				Method:     "tickets.patch",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
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
				Method:     "tickets.patch",
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
				Method:     "tickets.patch",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Bad Request",
							"message": "The participant was already registered for the event previously",
						},
					},
				},
			})
		}

		switch {
		case request.Participant != 0 && request.Event != 0:
			if constants.Persistence == storage.PostgreSQL {
				q = "UPDATE tickets SET participant = $1, event = $2 WHERE id = $3;"
			} else {
				q = "UPDATE tickets SET participant = ?, event = ? WHERE id = ?;"
			}
		case request.Participant != 0:
			if constants.Persistence == storage.PostgreSQL {
				q = "UPDATE tickets SET participant = $1 WHERE id = $2;"
			} else {
				q = "UPDATE tickets SET participant = ? WHERE id = ?;"
			}
		case request.Event != 0:
			if constants.Persistence == storage.PostgreSQL {
				q = "UPDATE tickets SET event = $1 WHERE id = $2;"
			} else {
				q = "UPDATE tickets SET event = ? WHERE id = ?;"
			}
		}

		stmt, err = db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.patch",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
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

		var r sql.Result
		switch {
		case request.Participant != 0 && request.Event != 0:
			r, err = stmt.ExecContext(ctx, request.Participant, request.Event, id)
		case request.Participant != 0:
			r, err = stmt.ExecContext(ctx, request.Participant, id)
		case request.Event != 0:
			r, err = stmt.ExecContext(ctx, request.Event, id)
		}

		if err != nil {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.patch",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    400,
					Message: "Bad request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Bad request",
							"message": "The request body data or parameters was rejected, not valid",
						},
					},
				},
			})
		}
		if i, _ := r.RowsAffected(); i == 0 {
			return c.JSON(http.StatusNotFound, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "tickets.patch",
				Context:    c.Request().URL.String(),
				Params: map[string]interface{}{
					"id": id,
				},
				Error: models.Error{
					Code:    404,
					Message: "Not Found",
					Errors: []map[string]interface{}{
						{
							"reason": "Not Found",
							"message": "Ticket not found",
						},
					},
				},
			})
		}
		return c.JSON(http.StatusOK, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "tickets.patch",
			Context:    c.Request().URL.String(),
		})
	}
}
