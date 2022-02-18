package participants

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/luisnquin/restapi-technical-test/src/constants"
	"github.com/luisnquin/restapi-technical-test/src/models"
	"github.com/luisnquin/restapi-technical-test/src/storage"
)

func New() echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			db      = storage.Get(constants.Persistence)
			request = new(models.Participant)
			err     error
		)

		if err = c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "participants.post",
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

		if (*request == models.Participant{}) {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "participants.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Bad Request",
							"message": "The request body data is empty",
						},
					},
				},
			})
		}

		if err = db.Connect(); err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "participants.post",
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
			q = "INSERT INTO participants(firstname, lastname, age) VALUES($1, $2, $3);"
		case storage.MySQL:
			q = "INSERT INTO participants(firstname, lastname, age) VALUES(?, ?, ?);"
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		stmt, err := db.PrepareContext(ctx, q)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "participants.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    500,
					Message: "Internal server error",
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

		r, err := stmt.ExecContext(ctx, request.Firstname, request.Lastname, request.Age)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.BadResponse{
				APIVersion: constants.APIVersion,
				Method:     "participants.post",
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
				Method:     "participants.post",
				Context:    c.Request().URL.String(),
				Error: models.Error{
					Code:    400,
					Message: "Bad Request",
					Errors: []map[string]interface{}{
						{
							"reason":  "Bad Request",
							"message": "Are you following any criteria for insertion?",
						},
					},
				},
			})
		}
		return c.JSON(http.StatusCreated, models.SuccessfulResponse{
			APIVersion: constants.APIVersion,
			Method:     "participants.post",
			Context:    c.Request().URL.String(),
		})
	}
}
